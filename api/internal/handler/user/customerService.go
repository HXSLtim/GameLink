package user

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/apierr"
)

type customerServiceHandler struct {
	chatSvc  *chatservice.ChatService
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

type customerServiceSessionResponse struct {
	GroupID  uint64                       `json:"groupId"`
	Agent    customerServiceAgentResponse `json:"agent"`
	IsOnline bool                         `json:"isOnline"`
}

type customerServiceAgentResponse struct {
	UserID    uint64 `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	Status    string `json:"status"`
}

type customerServiceMessageResponse struct {
	ID          uint64                `json:"id"`
	GroupID     uint64                `json:"groupId"`
	SenderID    uint64                `json:"senderId"`
	Content     string                `json:"content"`
	MessageType model.ChatMessageType `json:"messageType"`
	IsMe        bool                  `json:"isMe"`
	CreatedAt   string                `json:"createdAt"`
}

type sendCustomerServiceMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

func RegisterCustomerServiceRoutes(
	router gin.IRouter,
	chatSvc *chatservice.ChatService,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	authMiddleware gin.HandlerFunc,
) {
	if chatSvc == nil || userRepo == nil {
		return
	}

	h := &customerServiceHandler{
		chatSvc:  chatSvc,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}

	group := router.Group("/customer-service")
	group.Use(authMiddleware)
	group.GET("/session", h.getSession)
	group.GET("/messages", h.listMessages)
	group.POST("/messages", h.sendMessage)
}

// getSession 获取当前用户的客服会话（不存在时自动创建）。
func (h *customerServiceHandler) getSession(c *gin.Context) {
	userID := getUserIDFromContext(c)

	group, agent, err := h.ensureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[customerServiceSessionResponse](c, http.StatusOK, model.APIResponse[customerServiceSessionResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: customerServiceSessionResponse{
			GroupID:  group.ID,
			Agent:    toCustomerServiceAgent(agent),
			IsOnline: agent.Status == model.UserStatusActive,
		},
	})
}

// listMessages 获取客服会话历史消息。
func (h *customerServiceHandler) listMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	page, pageSize := parsePagination(c)

	var beforeID *uint64
	if raw := strings.TrimSpace(c.Query("beforeId")); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			beforeID = &parsed
		}
	}

	group, _, err := h.ensureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	messages, total, err := h.chatSvc.ListMessages(c.Request.Context(), userID, group.ID, chatservice.ListMessagesOptions{
		Page:     page,
		PageSize: pageSize,
		BeforeID: beforeID,
	})
	if err != nil {
		respondAPIError(c, err)
		return
	}

	normalized := make([]customerServiceMessageResponse, 0, len(messages))
	for _, item := range messages {
		content := strings.TrimSpace(item.Content)
		if content == "" && item.ImageURL != "" {
			content = item.ImageURL
		}
		normalized = append(normalized, customerServiceMessageResponse{
			ID:          item.ID,
			GroupID:     item.GroupID,
			SenderID:    item.SenderID,
			Content:     content,
			MessageType: item.MessageType,
			IsMe:        item.SenderID == userID,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ID < normalized[j].ID
	})

	if len(normalized) > 0 {
		last := normalized[len(normalized)-1]
		_ = h.chatSvc.MarkRead(c.Request.Context(), group.ID, userID, last.ID)
	}

	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"groupId":   group.ID,
			"messages":  normalized,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"hasMore":   int64(page*pageSize) < total,
			"isSession": true,
		},
	})
}

// sendMessage 在客服会话中发送文本消息。
func (h *customerServiceHandler) sendMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req sendCustomerServiceMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		respondAPIError(c, apierr.BadRequest("content is required"))
		return
	}

	group, _, err := h.ensureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	msg, err := h.chatSvc.SendMessage(c.Request.Context(), chatservice.SendMessageInput{
		GroupID:     group.ID,
		SenderID:    userID,
		Content:     content,
		MessageType: model.ChatMessageTypeText,
	})
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[customerServiceMessageResponse](c, http.StatusCreated, model.APIResponse[customerServiceMessageResponse]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data: customerServiceMessageResponse{
			ID:          msg.ID,
			GroupID:     msg.GroupID,
			SenderID:    msg.SenderID,
			Content:     msg.Content,
			MessageType: msg.MessageType,
			IsMe:        true,
			CreatedAt:   msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

func (h *customerServiceHandler) ensureSession(ctx context.Context, userID uint64) (*model.ChatGroup, *model.User, error) {
	agentIDs, err := h.resolveServiceAgentIDs(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if len(agentIDs) == 0 {
		return nil, nil, apierr.NotFound("no customer service agent available")
	}

	agentSet := make(map[uint64]struct{}, len(agentIDs))
	for _, item := range agentIDs {
		agentSet[item] = struct{}{}
	}

	groups, _, err := h.chatSvc.ListUserGroups(ctx, userID, 1, 100)
	if err != nil {
		return nil, nil, err
	}

	for _, group := range groups {
		if group.GroupType != model.ChatGroupTypePrivate || !group.IsActive {
			continue
		}

		var matchedAgentID uint64
		for _, member := range group.Members {
			if member.UserID == userID {
				continue
			}
			if _, ok := agentSet[member.UserID]; ok {
				matchedAgentID = member.UserID
				break
			}
		}

		if matchedAgentID == 0 {
			continue
		}

		agent, getErr := h.userRepo.Get(ctx, matchedAgentID)
		if getErr != nil {
			continue
		}

		return &group, agent, nil
	}

	selectedAgentID := agentIDs[0]
	created, err := h.chatSvc.CreateGroup(ctx, userID, selectedAgentID, chatservice.CreateGroupRequest{
		TargetUserID: selectedAgentID,
		GroupType:    model.ChatGroupTypePrivate,
	})
	if err != nil {
		return nil, nil, err
	}

	agent, err := h.userRepo.Get(ctx, selectedAgentID)
	if err != nil {
		return nil, nil, err
	}

	return created, agent, nil
}

func (h *customerServiceHandler) resolveServiceAgentIDs(ctx context.Context, requesterID uint64) ([]uint64, error) {
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

	if h.roleRepo != nil {
		for _, item := range rolePriority {
			role, err := h.roleRepo.GetBySlug(ctx, item.Slug)
			if err != nil || role == nil {
				continue
			}

			userIDs, err := h.roleRepo.GetUserIDsByRoleID(ctx, role.ID)
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

	if len(priorityByUser) == 0 {
		fallbackEmails := []string{"cs.agent@gamelink.com", "cs.leader@gamelink.com"}
		for idx, email := range fallbackEmails {
			user, err := h.userRepo.GetByEmail(ctx, email)
			if err != nil || user == nil || user.ID == requesterID {
				continue
			}
			priorityByUser[user.ID] = 10 + idx
		}
	}

	ranked := make([]rankedAgent, 0, len(priorityByUser))
	for userID, priority := range priorityByUser {
		user, err := h.userRepo.Get(ctx, userID)
		if err != nil || user == nil {
			continue
		}
		if user.Status != model.UserStatusActive {
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

func toCustomerServiceAgent(user *model.User) customerServiceAgentResponse {
	if user == nil {
		return customerServiceAgentResponse{}
	}
	nickname := strings.TrimSpace(user.Nickname)
	if nickname == "" {
		nickname = strings.TrimSpace(user.Name)
	}
	if nickname == "" {
		nickname = "在线客服"
	}

	return customerServiceAgentResponse{
		UserID:    user.ID,
		Nickname:  nickname,
		AvatarURL: user.AvatarURL,
		Status:    string(user.Status),
	}
}
