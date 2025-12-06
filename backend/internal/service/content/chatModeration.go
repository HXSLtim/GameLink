package content

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/sensitiveword"
)

// ChatModerationService 聊天内容审核服务
type ChatModerationService struct {
	messageRepo   repository.ChatMessageRepository
	memberRepo    repository.ChatMemberRepository
	sensitiveWord *sensitiveword.SensitiveWordService
	opLogRepo     repository.OperationLogRepository
}

// NewChatModerationService 创建聊天内容审核服务
func NewChatModerationService(
	messageRepo repository.ChatMessageRepository,
	memberRepo repository.ChatMemberRepository,
	sensitiveWord *sensitiveword.SensitiveWordService,
	opLogRepo repository.OperationLogRepository,
) *ChatModerationService {
	return &ChatModerationService{
		messageRepo:   messageRepo,
		memberRepo:    memberRepo,
		sensitiveWord: sensitiveWord,
		opLogRepo:     opLogRepo,
	}
}

// ChatMessageDTO 聊天消息DTO
type ChatMessageDTO struct {
	ID           uint64                       `json:"id"`
	GroupID      uint64                       `json:"groupId"`
	SenderID     uint64                       `json:"senderId"`
	Content      string                       `json:"content"`
	MessageType  model.ChatMessageType        `json:"messageType"`
	AuditStatus  model.ChatMessageAuditStatus `json:"auditStatus"`
	RejectReason string                       `json:"rejectReason,omitempty"`
	CreatedAt    string                       `json:"createdAt"`
	// 敏感词检测结果
	HasSensitiveWords  bool   `json:"hasSensitiveWords,omitempty"`
	HighlightedContent string `json:"highlightedContent,omitempty"`
}

// ListMessagesRequest 列出消息请求
type ListMessagesRequest struct {
	Page        int                           `form:"page"`
	PageSize    int                           `form:"pageSize"`
	GroupID     *uint64                       `form:"groupId"`
	SenderID    *uint64                       `form:"senderId"`
	AuditStatus *model.ChatMessageAuditStatus `form:"auditStatus"`
	DateFrom    *time.Time                    `form:"dateFrom"`
	DateTo      *time.Time                    `form:"dateTo"`
}

// ListMessagesResponse 列出消息响应
type ListMessagesResponse struct {
	Messages []ChatMessageDTO `json:"messages"`
	Total    int64            `json:"total"`
}

// MuteUserRequest 禁言用户请求
type MuteUserRequest struct {
	GroupID  uint64 `json:"groupId" binding:"required"`
	UserID   uint64 `json:"userId" binding:"required"`
	Duration int    `json:"duration" binding:"required,min=1"` // 禁言时长（分钟）
	Reason   string `json:"reason"`
}

// ListMessages 列出聊天消息（管理员）
func (s *ChatModerationService) ListMessages(ctx context.Context, req ListMessagesRequest) (*ListMessagesResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	messages, total, err := s.messageRepo.ListForModeration(ctx, repository.ChatMessageModerationListOptions{
		Page:        req.Page,
		PageSize:    req.PageSize,
		GroupID:     req.GroupID,
		SenderID:    req.SenderID,
		AuditStatus: req.AuditStatus,
		DateFrom:    req.DateFrom,
		DateTo:      req.DateTo,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]ChatMessageDTO, 0, len(messages))
	for _, m := range messages {
		dto := s.toMessageDTO(&m)
		// 检测敏感词
		if s.sensitiveWord != nil && m.MessageType == model.ChatMessageTypeText {
			result, _ := s.sensitiveWord.DetectSensitiveWords(ctx, sensitiveword.DetectSensitiveWordsRequest{
				Content: m.Content,
			})
			if result != nil {
				dto.HasSensitiveWords = result.HasSensitiveWords
				dto.HighlightedContent = result.HighlightedContent
			}
		}
		dtos = append(dtos, *dto)
	}

	return &ListMessagesResponse{
		Messages: dtos,
		Total:    total,
	}, nil
}

// DeleteMessage 删除聊天消息
func (s *ChatModerationService) DeleteMessage(ctx context.Context, messageID uint64, moderatorID uint64, reason string) error {
	if err := s.messageRepo.MarkDeleted(ctx, messageID, moderatorID); err != nil {
		return err
	}

	// 记录操作日志
	s.logOperation(ctx, messageID, moderatorID, model.OpActionDeleteMessage, reason)
	return nil
}

// ApproveMessage 批准消息
func (s *ChatModerationService) ApproveMessage(ctx context.Context, messageID uint64, moderatorID uint64) error {
	if err := s.messageRepo.UpdateAuditStatus(ctx, messageID, model.ChatMessageAuditApproved, &moderatorID, ""); err != nil {
		return err
	}

	s.logOperation(ctx, messageID, moderatorID, model.OpActionApprove, "")
	return nil
}

// RejectMessage 拒绝消息
func (s *ChatModerationService) RejectMessage(ctx context.Context, messageID uint64, moderatorID uint64, reason string) error {
	if reason == "" {
		return ErrAdminValidation
	}

	if err := s.messageRepo.UpdateAuditStatus(ctx, messageID, model.ChatMessageAuditRejected, &moderatorID, reason); err != nil {
		return err
	}

	s.logOperation(ctx, messageID, moderatorID, model.OpActionReject, reason)
	return nil
}

// MuteUser 禁言用户
func (s *ChatModerationService) MuteUser(ctx context.Context, req MuteUserRequest, moderatorID uint64) error {
	if req.Duration < 1 {
		return ErrAdminValidation
	}

	// 获取群成员
	member, err := s.memberRepo.Get(ctx, req.GroupID, req.UserID)
	if err != nil {
		return err
	}

	// 设置禁言
	mutedUntil := time.Now().Add(time.Duration(req.Duration) * time.Minute)
	member.IsMuted = true
	member.MutedUntil = &mutedUntil
	member.MutedBy = &moderatorID
	member.MuteReason = req.Reason

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return err
	}

	// 记录操作日志
	s.logOperation(ctx, req.UserID, moderatorID, model.OpActionMuteUser, req.Reason)
	return nil
}

// UnmuteUser 解除禁言
func (s *ChatModerationService) UnmuteUser(ctx context.Context, groupID, userID, moderatorID uint64) error {
	member, err := s.memberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return err
	}

	member.IsMuted = false
	member.MutedUntil = nil
	member.MutedBy = nil
	member.MuteReason = ""

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return err
	}

	s.logOperation(ctx, userID, moderatorID, model.OpActionUnmuteUser, "")
	return nil
}

func (s *ChatModerationService) logOperation(ctx context.Context, entityID uint64, actorID uint64, action model.OperationAction, reason string) {
	if s.opLogRepo == nil {
		return
	}
	_ = s.opLogRepo.Append(ctx, &model.OperationLog{
		EntityType:  string(model.OpEntityChatMessage),
		EntityID:    entityID,
		ActorUserID: &actorID,
		Action:      string(action),
		Reason:      reason,
	})
}

func (s *ChatModerationService) toMessageDTO(msg *model.ChatMessage) *ChatMessageDTO {
	return &ChatMessageDTO{
		ID:           msg.ID,
		GroupID:      msg.GroupID,
		SenderID:     msg.SenderID,
		Content:      msg.Content,
		MessageType:  msg.MessageType,
		AuditStatus:  msg.AuditStatus,
		RejectReason: msg.RejectReason,
		CreatedAt:    msg.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
