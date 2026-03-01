package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	conversationservice "gamelink/internal/service/conversation"
	"gamelink/pkg/apierr"
)

type customerServiceHandler struct {
	conversationSvc *conversationservice.Service
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
	conversationSvc *conversationservice.Service,
	authMiddleware gin.HandlerFunc,
) {
	if conversationSvc == nil {
		return
	}

	h := &customerServiceHandler{
		conversationSvc: conversationSvc,
	}

	group := router.Group("/customer-service")
	group.Use(authMiddleware)
	group.GET("/session", h.getSession)
	group.GET("/messages", h.listMessages)
	group.POST("/messages", h.sendMessage)
	group.GET("/conversations", h.listConversations)
	group.POST("/conversations", h.createConversation)
	group.GET("/conversations/:id/messages", h.listConversationMessages)
	group.POST("/conversations/:id/messages", h.sendConversationMessage)
	group.DELETE("/conversations/:id", h.closeConversation)
}

// getSession 获取当前用户的客服会话（不存在时自动创建）。
func (h *customerServiceHandler) getSession(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	group, agent, err := h.conversationSvc.EnsureSession(c.Request.Context(), userID)
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
			IsOnline: agent != nil && agent.Status == model.UserStatusActive,
		},
	})
}

// listMessages 获取当前客服会话历史消息（兼容旧接口）。
func (h *customerServiceHandler) listMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	page, pageSize := parsePagination(c)
	beforeID := parseBeforeMessageID(c)

	group, _, err := h.conversationSvc.EnsureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	messages, total, err := h.conversationSvc.ListMessages(c.Request.Context(), userID, group.ID, page, pageSize, beforeID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	normalized := make([]customerServiceMessageResponse, 0, len(messages))
	for _, item := range messages {
		normalized = append(normalized, toCustomerServiceMessageResponse(item))
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

// sendMessage 在当前客服会话中发送消息（兼容旧接口）。
func (h *customerServiceHandler) sendMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	var req sendCustomerServiceMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	group, _, err := h.conversationSvc.EnsureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	message, err := h.conversationSvc.SendMessage(c.Request.Context(), userID, group.ID, req.Content)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[customerServiceMessageResponse](c, http.StatusCreated, model.APIResponse[customerServiceMessageResponse]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    toCustomerServiceMessageResponse(*message),
	})
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

func toCustomerServiceMessageResponse(item model.ConversationMessage) customerServiceMessageResponse {
	return customerServiceMessageResponse{
		ID:          item.ID,
		GroupID:     item.GroupID,
		SenderID:    item.SenderID,
		Content:     item.Content,
		MessageType: item.MessageType,
		IsMe:        item.IsMe,
		CreatedAt:   item.CreatedAt,
	}
}
