package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

// Errors specific to chat domain.
var (
	ErrNotFound        = repository.ErrNotFound
	ErrNotMember       = apierr.Forbidden("用户不是聊天组成员")
	ErrInactiveGroup   = apierr.BadRequest("聊天组已停用")
	ErrMessageTooLarge = apierr.BadRequest("消息超出长度限制")
	ErrThrottled       = apierr.TooManyRequests("消息发送过于频繁，请稍后再试")
)

// SendMessageInput represents payload for sending chat messages.
type SendMessageInput struct {
	GroupID     uint64
	SenderID    uint64
	Content     string
	MessageType model.ChatMessageType
	ReplyToID   *uint64
	ImageURL    string
}

// ApproveMessage sets audit status to approved.
func (s *ChatService) ApproveMessage(ctx context.Context, messageID uint64, moderatorID uint64) error {
	return s.messages.UpdateAuditStatus(ctx, messageID, model.ChatMessageAuditApproved, &moderatorID, "")
}

// RejectMessage sets audit status to rejected with reason.
func (s *ChatService) RejectMessage(ctx context.Context, messageID uint64, moderatorID uint64, reason string) error {
	if reason == "" {
		reason = "rejected"
	}
	return s.messages.UpdateAuditStatus(ctx, messageID, model.ChatMessageAuditRejected, &moderatorID, reason)
}

// ReportMessage creates a report record for moderation.
func (s *ChatService) ReportMessage(ctx context.Context, reporterID, messageID uint64, reason, evidence string) error {
	if s.reports == nil {
		return apierr.InternalError("举报仓库未配置")
	}
	if reason == "" {
		reason = "unspecified"
	}
	report := &model.ChatReport{
		MessageID:  messageID,
		ReporterID: reporterID,
		Reason:     reason,
		Evidence:   evidence,
		Status:     "pending",
	}
	if err := s.reports.Create(ctx, report); err != nil {
		return apierr.InternalError("创建举报记录失败")
	}
	return nil
}

// ListMessagesOptions defines pagination options for chat history.
type ListMessagesOptions struct {
	Page     int
	PageSize int
	BeforeID *uint64
	AfterID  *uint64
}

// ChatService aggregates chat domain logic.
type ChatService struct {
	groups   repository.ChatGroupRepository
	members  repository.ChatMemberRepository
	messages repository.ChatMessageRepository
	reports  repository.ChatReportRepository
	cache    cache.Cache
}

// NewChatService constructs a ChatService instance.
func NewChatService(
	groups repository.ChatGroupRepository,
	members repository.ChatMemberRepository,
	messages repository.ChatMessageRepository,
	reports repository.ChatReportRepository,
	cache cache.Cache,
) *ChatService {
	return &ChatService{
		groups:   groups,
		members:  members,
		messages: messages,
		reports:  reports,
		cache:    cache,
	}
}

// ListUserGroups returns groups joined by the user with pagination.
func (s *ChatService) ListUserGroups(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	groups, total, err := s.groups.ListByUser(ctx, userID, repository.ChatGroupListOptions{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, apierr.InternalError("获取用户聊天组失败").WithDetails(err.Error())
	}
	return groups, total, nil
}

// EnsureMembership verifies that the user belongs to group.
func (s *ChatService) EnsureMembership(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
	member, err := s.members.Get(ctx, groupID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotMember
		}
		return nil, apierr.InternalError("获取聊天成员信息失败").WithDetails(err.Error())
	}
	if !member.IsActive {
		return nil, ErrNotMember
	}
	return member, nil
}

// ListMessages returns chat history for a group.
func (s *ChatService) ListMessages(ctx context.Context, userID, groupID uint64, opts ListMessagesOptions) ([]model.ChatMessage, int64, error) {
	if _, err := s.EnsureMembership(ctx, groupID, userID); err != nil {
		return nil, 0, err
	}

	group, err := s.groups.Get(ctx, groupID)
	if err != nil {
		return nil, 0, apierr.InternalError("获取聊天组失败").WithDetails(err.Error())
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 50
	}

	listOpts := repository.ChatMessageListOptions{
		GroupID:  groupID,
		Page:     opts.Page,
		PageSize: opts.PageSize,
		BeforeID: opts.BeforeID,
		AfterID:  opts.AfterID,
	}
	// 在公共群只返回已审核通过的消息
	if group.GroupType == model.ChatGroupTypePublic {
		listOpts.AuditStatuses = []model.ChatMessageAuditStatus{model.ChatMessageAuditApproved}
	}

	messages, total, err := s.messages.ListByGroup(ctx, listOpts)
	if err != nil {
		return nil, 0, apierr.InternalError("获取聊天消息失败").WithDetails(err.Error())
	}
	return messages, total, nil
}

// SendMessage persists chat message and returns saved entity.
func (s *ChatService) SendMessage(ctx context.Context, input SendMessageInput) (*model.ChatMessage, error) {
	if input.Content == "" && input.ImageURL == "" {
		return nil, ErrMessageTooLarge
	}
	if len([]rune(input.Content)) > 2000 {
		return nil, ErrMessageTooLarge
	}

	if _, err := s.EnsureMembership(ctx, input.GroupID, input.SenderID); err != nil {
		return nil, err
	}

	group, err := s.groups.Get(ctx, input.GroupID)
	if err != nil {
		return nil, apierr.InternalError("获取聊天组失败").WithDetails(err.Error())
	}
	if !group.IsActive {
		return nil, ErrInactiveGroup
	}

	// 公共群发言限流
	if group.GroupType == model.ChatGroupTypePublic {
		key := fmt.Sprintf("chat:throttle:g:%d:u:%d", input.GroupID, input.SenderID)
		if _, ok, _ := s.cache.Get(ctx, key); ok {
			return nil, ErrThrottled
		}
		// 默认30s冷却
		_ = s.cache.Set(ctx, key, "1", 30*time.Second)
	}

	msg := &model.ChatMessage{
		GroupID:     input.GroupID,
		SenderID:    input.SenderID,
		Content:     input.Content,
		MessageType: input.MessageType,
		ReplyToID:   input.ReplyToID,
		ImageURL:    input.ImageURL,
		Metadata:    "{}",
	}

	// 公共群消息默认 pending，订单群直接 approved（如需严格也可全部 pending）
	if group.GroupType == model.ChatGroupTypePublic {
		msg.AuditStatus = model.ChatMessageAuditPending
	} else {
		msg.AuditStatus = model.ChatMessageAuditApproved
	}

	if err := s.messages.Create(ctx, msg); err != nil {
		return nil, apierr.InternalError("创建聊天消息失败").WithDetails(err.Error())
	}

	return msg, nil
}

// JoinGroup marks user as active member of group (creates if needed).
func (s *ChatService) JoinGroup(ctx context.Context, groupID, userID uint64, nickname string) error {
	group, err := s.groups.Get(ctx, groupID)
	if err != nil {
		return apierr.InternalError("获取聊天组失败").WithDetails(err.Error())
	}
	if !group.IsActive {
		return ErrInactiveGroup
	}

	member, err := s.members.Get(ctx, groupID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			m := &model.ChatGroupMember{
				GroupID:  groupID,
				UserID:   userID,
				Role:     "member",
				Nickname: nickname,
				JoinedAt: time.Now(),
				IsActive: true,
			}
			return s.members.Add(ctx, m)
		}
		return apierr.InternalError("获取聊天成员信息失败").WithDetails(err.Error())
	}

	// Reactivate existing member.
	member.IsActive = true
	member.Nickname = nickname
	return s.members.Update(ctx, member)
}

// LeaveGroup marks the user's membership as inactive.
func (s *ChatService) LeaveGroup(ctx context.Context, groupID, userID uint64) error {
	member, err := s.members.Get(ctx, groupID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotMember
		}
		return apierr.InternalError("获取聊天成员信息失败").WithDetails(err.Error())
	}
	member.IsActive = false
	return s.members.Update(ctx, member)
}

// MarkRead updates last read pointer for membership.
func (s *ChatService) MarkRead(ctx context.Context, groupID, userID, messageID uint64) error {
	member, err := s.members.Get(ctx, groupID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotMember
		}
		return apierr.InternalError("获取聊天成员信息失败").WithDetails(err.Error())
	}
	member.LastReadMessageID = &messageID
	now := time.Now()
	member.LastReadAt = &now
	return s.members.Update(ctx, member)
}
