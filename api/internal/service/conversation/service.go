package conversation

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// ConversationRepository defines persistence operations required by Service.
type ConversationRepository interface {
	Create(ctx context.Context, group *model.ChatGroup, members []*model.ChatGroupMember) error
	Get(ctx context.Context, conversationID uint64) (*model.ChatGroup, error)
	ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error)
	Update(ctx context.Context, group *model.ChatGroup) error
	Delete(ctx context.Context, conversationID uint64) error
	FindActivePrivateByParticipants(ctx context.Context, userID uint64, agentIDs []uint64) (*model.ChatGroup, error)
	ListMessages(ctx context.Context, groupID uint64, page, pageSize int, beforeID *uint64) ([]model.ChatMessage, int64, error)
	CreateMessage(ctx context.Context, message *model.ChatMessage) error
	CountActiveByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error)
	IsMember(ctx context.Context, groupID, userID uint64) (bool, error)
}

// UserStore defines user lookups needed by Service.
type UserStore interface {
	Get(ctx context.Context, id uint64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

// RoleStore defines role lookup operations needed by Service.
type RoleStore interface {
	GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error)
	GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error)
}

// HubPublisher is the subset of websocket hub operations used by Service.
type HubPublisher interface {
	BroadcastToGroup(message []byte, groupID uint64)
	BroadcastToUser(message []byte, userID uint64)
}

// Service encapsulates customer-service conversation business logic.
type Service struct {
	repo  ConversationRepository
	users UserStore
	roles RoleStore
	hub   HubPublisher
}

// NewService creates a conversation service.
func NewService(repo ConversationRepository, users UserStore, roles RoleStore, hub HubPublisher) *Service {
	return &Service{
		repo:  repo,
		users: users,
		roles: roles,
		hub:   hub,
	}
}

// EnsureSession finds existing active conversation or creates a new one with least-loaded agent.
func (s *Service) EnsureSession(ctx context.Context, userID uint64) (*model.ChatGroup, *model.User, error) {
	agentIDs, err := s.resolveServiceAgentIDs(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if len(agentIDs) == 0 {
		return nil, nil, apierr.NotFound("no customer service agent available")
	}

	existing, err := s.repo.FindActivePrivateByParticipants(ctx, userID, agentIDs)
	if err != nil {
		return nil, nil, apierr.InternalError("failed to find conversation").WithDetails(err.Error())
	}
	if existing != nil {
		agentID := pickConversationAgentID(existing, userID, agentIDs)
		agent := s.getUser(agentID)
		return existing, agent, nil
	}

	agentLoads, err := s.repo.CountActiveByAgentIDs(ctx, agentIDs)
	if err != nil {
		return nil, nil, apierr.InternalError("failed to calculate agent load").WithDetails(err.Error())
	}
	assignedAgentID := chooseLeastLoadedAgent(agentIDs, agentLoads)
	assignedAgent := s.getUser(assignedAgentID)
	if assignedAgent == nil {
		return nil, nil, apierr.NotFound("no active customer service agent available")
	}

	group := &model.ChatGroup{
		GroupName:  "在线客服",
		GroupType:  model.ChatGroupTypePrivate,
		CreatedBy:  userID,
		MaxMembers: 2,
		IsActive:   true,
	}
	members := []*model.ChatGroupMember{
		{
			GroupID:  group.ID,
			UserID:   userID,
			Role:     model.ChatMemberRoleOwner,
			Nickname: "用户",
			JoinedAt: time.Now(),
			IsActive: true,
		},
		{
			GroupID:  group.ID,
			UserID:   assignedAgentID,
			Role:     model.ChatMemberRoleMember,
			Nickname: s.resolveNickname(assignedAgent),
			JoinedAt: time.Now(),
			IsActive: true,
		},
	}
	if err := s.repo.Create(ctx, group, members); err != nil {
		return nil, nil, apierr.InternalError("failed to create conversation").WithDetails(err.Error())
	}

	loaded, err := s.repo.Get(ctx, group.ID)
	if err != nil {
		return nil, nil, apierr.InternalError("failed to load conversation").WithDetails(err.Error())
	}
	return loaded, assignedAgent, nil
}

// ListConversations lists conversations visible to user.
func (s *Service) ListConversations(ctx context.Context, userID uint64, page, pageSize int) ([]model.Conversation, int64, error) {
	groups, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, apierr.InternalError("failed to list conversations").WithDetails(err.Error())
	}

	conversations := make([]model.Conversation, 0, len(groups))
	for _, group := range groups {
		agentID := pickConversationAgentID(&group, userID, nil)
		agent := s.getUser(agentID)

		lastMessages, _, listErr := s.repo.ListMessages(ctx, group.ID, 1, 1, nil)
		var lastMessage *model.ChatMessage
		if listErr == nil && len(lastMessages) > 0 {
			lastMessage = &lastMessages[0]
		}

		conversations = append(conversations, buildConversation(&group, agent, userID, lastMessage))
	}

	return conversations, total, nil
}

// ListMessages returns paginated conversation messages for user.
func (s *Service) ListMessages(
	ctx context.Context,
	userID uint64,
	conversationID uint64,
	page int,
	pageSize int,
	beforeID *uint64,
) ([]model.ConversationMessage, int64, error) {
	if ok, err := s.repo.IsMember(ctx, conversationID, userID); err != nil {
		return nil, 0, apierr.InternalError("failed to verify conversation member").WithDetails(err.Error())
	} else if !ok {
		return nil, 0, apierr.Forbidden("forbidden")
	}

	messages, total, err := s.repo.ListMessages(ctx, conversationID, page, pageSize, beforeID)
	if err != nil {
		return nil, 0, apierr.InternalError("failed to list conversation messages").WithDetails(err.Error())
	}

	result := make([]model.ConversationMessage, 0, len(messages))
	for _, item := range messages {
		result = append(result, toConversationMessage(item, userID))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, total, nil
}

// SendMessage creates and broadcasts a customer-service message.
func (s *Service) SendMessage(ctx context.Context, userID uint64, conversationID uint64, content string) (*model.ConversationMessage, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, apierr.BadRequest("content is required")
	}

	if ok, err := s.repo.IsMember(ctx, conversationID, userID); err != nil {
		return nil, apierr.InternalError("failed to verify conversation member").WithDetails(err.Error())
	} else if !ok {
		return nil, apierr.Forbidden("forbidden")
	}

	group, err := s.repo.Get(ctx, conversationID)
	if err != nil {
		return nil, apierr.NotFound("conversation not found")
	}
	if !group.IsActive {
		return nil, apierr.BadRequest("conversation is closed")
	}

	message := &model.ChatMessage{
		GroupID:     conversationID,
		SenderID:    userID,
		Content:     trimmed,
		MessageType: model.ChatMessageTypeText,
		AuditStatus: model.ChatMessageAuditApproved,
	}
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, apierr.InternalError("failed to send message").WithDetails(err.Error())
	}

	s.broadcastMessage(group, message)

	dto := toConversationMessage(*message, userID)
	return &dto, nil
}

// CreateConversation ensures a session and posts initial message.
func (s *Service) CreateConversation(ctx context.Context, userID uint64, content string) (model.Conversation, *model.ConversationMessage, error) {
	group, agent, err := s.EnsureSession(ctx, userID)
	if err != nil {
		return model.Conversation{}, nil, err
	}

	message, err := s.SendMessage(ctx, userID, group.ID, content)
	if err != nil {
		return model.Conversation{}, nil, err
	}

	conv := buildConversation(group, agent, userID, nil)
	conv.LastMessage = message.Content
	if parsed, parseErr := time.Parse(time.RFC3339, message.CreatedAt); parseErr == nil {
		conv.LastMessageAt = &parsed
	}
	return conv, message, nil
}

// CloseConversation marks conversation inactive.
func (s *Service) CloseConversation(ctx context.Context, userID, conversationID uint64) (*model.Conversation, error) {
	if ok, err := s.repo.IsMember(ctx, conversationID, userID); err != nil {
		return nil, apierr.InternalError("failed to verify conversation member").WithDetails(err.Error())
	} else if !ok {
		return nil, apierr.Forbidden("forbidden")
	}

	group, err := s.repo.Get(ctx, conversationID)
	if err != nil {
		return nil, apierr.NotFound("conversation not found")
	}
	if err := s.repo.Delete(ctx, conversationID); err != nil {
		return nil, apierr.InternalError("failed to close conversation").WithDetails(err.Error())
	}

	group.IsActive = false
	agentID := pickConversationAgentID(group, userID, nil)
	agent := s.getUser(agentID)
	conv := buildConversation(group, agent, userID, nil)
	return &conv, nil
}

func (s *Service) resolveServiceAgentIDs(ctx context.Context, requesterID uint64) ([]uint64, error) {
	type rankedAgent struct {
		UserID    uint64
		Priority  int
		CreatedAt int64
	}

	priorityByUser := make(map[uint64]int)
	rolePriority := []struct {
		Slug     string
		Priority int
	}{
		{Slug: string(model.RoleSlugCSAgent), Priority: 1},
		{Slug: string(model.RoleSlugCSLeader), Priority: 2},
		{Slug: string(model.RoleSlugCustomerService), Priority: 3},
	}

	if s.roles != nil {
		for _, item := range rolePriority {
			role, err := s.roles.GetBySlug(ctx, item.Slug)
			if err != nil || role == nil {
				continue
			}
			userIDs, err := s.roles.GetUserIDsByRoleID(ctx, role.ID)
			if err != nil {
				continue
			}
			for _, userID := range userIDs {
				if userID == 0 || userID == requesterID {
					continue
				}
				if existing, ok := priorityByUser[userID]; !ok || item.Priority < existing {
					priorityByUser[userID] = item.Priority
				}
			}
		}
	}

	if len(priorityByUser) == 0 && s.users != nil {
		for idx, email := range []string{"cs.agent@gamelink.com", "cs.leader@gamelink.com"} {
			user, err := s.users.GetByEmail(ctx, email)
			if err != nil || user == nil || user.ID == requesterID {
				continue
			}
			priorityByUser[user.ID] = 10 + idx
		}
	}

	ranked := make([]rankedAgent, 0, len(priorityByUser))
	for userID, priority := range priorityByUser {
		user := s.getUser(userID)
		if user == nil || user.Status != model.UserStatusActive {
			continue
		}
		ranked = append(ranked, rankedAgent{
			UserID:    userID,
			Priority:  priority,
			CreatedAt: user.CreatedAt.Unix(),
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Priority != ranked[j].Priority {
			return ranked[i].Priority < ranked[j].Priority
		}
		if ranked[i].CreatedAt != ranked[j].CreatedAt {
			return ranked[i].CreatedAt < ranked[j].CreatedAt
		}
		return ranked[i].UserID < ranked[j].UserID
	})

	result := make([]uint64, 0, len(ranked))
	for _, item := range ranked {
		result = append(result, item.UserID)
	}
	return result, nil
}

func (s *Service) broadcastMessage(group *model.ChatGroup, message *model.ChatMessage) {
	if s.hub == nil || group == nil || message == nil {
		return
	}

	payload := map[string]any{
		"type": "conversation_message",
		"data": map[string]any{
			"groupId": group.ID,
			"message": map[string]any{
				"id":          message.ID,
				"groupId":     message.GroupID,
				"senderId":    message.SenderID,
				"content":     message.Content,
				"messageType": message.MessageType,
				"createdAt":   message.CreatedAt.Format(time.RFC3339),
			},
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	s.hub.BroadcastToGroup(bytes, group.ID)
	for _, member := range group.Members {
		if !member.IsActive {
			continue
		}
		s.hub.BroadcastToUser(bytes, member.UserID)
	}
}

func (s *Service) getUser(id uint64) *model.User {
	if s.users == nil || id == 0 {
		return nil
	}
	user, err := s.users.Get(context.Background(), id)
	if err != nil {
		return nil
	}
	return user
}

func buildConversation(group *model.ChatGroup, agent *model.User, userID uint64, lastMessage *model.ChatMessage) model.Conversation {
	conversation := model.Conversation{
		ID:            group.ID,
		GroupID:       group.ID,
		UserID:        userID,
		Status:        model.ConversationStatusActive,
		IsAgentOnline: agent != nil && agent.Status == model.UserStatusActive,
		CreatedAt:     group.CreatedAt,
		UpdatedAt:     group.UpdatedAt,
	}
	if !group.IsActive {
		conversation.Status = model.ConversationStatusClosed
	}

	if agent != nil {
		conversation.AgentID = agent.ID
		conversation.AgentAvatar = agent.AvatarURL
		conversation.AgentName = strings.TrimSpace(agent.Nickname)
		if conversation.AgentName == "" {
			conversation.AgentName = strings.TrimSpace(agent.Name)
		}
		if conversation.AgentName == "" {
			conversation.AgentName = "在线客服"
		}
	}

	if lastMessage != nil {
		conversation.LastMessage = resolveMessageContent(*lastMessage)
		t := lastMessage.CreatedAt
		conversation.LastMessageAt = &t
	}

	return conversation
}

func toConversationMessage(item model.ChatMessage, userID uint64) model.ConversationMessage {
	return model.ConversationMessage{
		ID:          item.ID,
		GroupID:     item.GroupID,
		SenderID:    item.SenderID,
		Content:     resolveMessageContent(item),
		MessageType: item.MessageType,
		IsMe:        item.SenderID == userID,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}

func resolveMessageContent(item model.ChatMessage) string {
	content := strings.TrimSpace(item.Content)
	if content == "" && item.ImageURL != "" {
		return item.ImageURL
	}
	return content
}

func pickConversationAgentID(group *model.ChatGroup, userID uint64, preferredAgentIDs []uint64) uint64 {
	if group == nil {
		return 0
	}
	if len(preferredAgentIDs) > 0 {
		set := make(map[uint64]struct{}, len(preferredAgentIDs))
		for _, item := range preferredAgentIDs {
			set[item] = struct{}{}
		}
		for _, member := range group.Members {
			if member.UserID == userID {
				continue
			}
			if _, ok := set[member.UserID]; ok {
				return member.UserID
			}
		}
	}

	for _, member := range group.Members {
		if member.UserID != userID {
			return member.UserID
		}
	}
	return 0
}

func chooseLeastLoadedAgent(agentIDs []uint64, loads map[uint64]int64) uint64 {
	if len(agentIDs) == 0 {
		return 0
	}

	selected := agentIDs[0]
	selectedLoad := loads[selected]
	for _, agentID := range agentIDs[1:] {
		load := loads[agentID]
		if load < selectedLoad {
			selected = agentID
			selectedLoad = load
		}
	}
	return selected
}

func (s *Service) resolveNickname(user *model.User) string {
	if user == nil {
		return "在线客服"
	}
	if nickname := strings.TrimSpace(user.Nickname); nickname != "" {
		return nickname
	}
	if name := strings.TrimSpace(user.Name); name != "" {
		return name
	}
	return "在线客服"
}
