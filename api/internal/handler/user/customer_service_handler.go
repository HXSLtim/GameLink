package user

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/apierr"
)

func (h *customerServiceHandler) listConversations(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	page, pageSize := parsePagination(c)
	beforeID := parseBeforeMessageID(c)

	group, agent, err := h.ensureSession(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	messages, total, err := h.loadConversationMessages(c.Request.Context(), userID, group.ID, page, pageSize, beforeID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		_ = h.chatSvc.MarkRead(c.Request.Context(), group.ID, userID, lastMessage.ID)
	}

	conversation := h.buildConversation(group, agent, userID, messages)
	payload := model.ConversationListResponse{
		Conversations: []model.Conversation{conversation},
		Messages:      messages,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		HasMore:       int64(page*pageSize) < total,
	}

	respondJSON[model.ConversationListResponse](c, http.StatusOK, model.APIResponse[model.ConversationListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    payload,
	})
}

func (h *customerServiceHandler) createConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	var req model.ConversationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		respondAPIError(c, apierr.BadRequest("content is required"))
		return
	}

	group, agent, err := h.ensureSession(c.Request.Context(), userID)
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

	message := toConversationMessage(*msg, userID)
	conversation := h.buildConversation(group, agent, userID, []model.ConversationMessage{message})

	respondJSON[any](c, http.StatusCreated, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data: gin.H{
			"conversation": conversation,
			"message":      message,
		},
	})
}

func (h *customerServiceHandler) loadConversationMessages(
	ctx context.Context,
	userID uint64,
	groupID uint64,
	page int,
	pageSize int,
	beforeID *uint64,
) ([]model.ConversationMessage, int64, error) {
	rawMessages, total, err := h.chatSvc.ListMessages(ctx, userID, groupID, chatservice.ListMessagesOptions{
		Page:     page,
		PageSize: pageSize,
		BeforeID: beforeID,
	})
	if err != nil {
		return nil, 0, err
	}

	messages := make([]model.ConversationMessage, 0, len(rawMessages))
	for _, item := range rawMessages {
		messages = append(messages, toConversationMessage(item, userID))
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	return messages, total, nil
}

func (h *customerServiceHandler) buildConversation(
	group *model.ChatGroup,
	agent *model.User,
	userID uint64,
	messages []model.ConversationMessage,
) model.Conversation {
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

	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		conversation.LastMessage = lastMessage.Content
		if parsed, err := time.Parse(time.RFC3339, lastMessage.CreatedAt); err == nil {
			conversation.LastMessageAt = &parsed
		}
	}

	return conversation
}

func toConversationMessage(item model.ChatMessage, userID uint64) model.ConversationMessage {
	content := strings.TrimSpace(item.Content)
	if content == "" && item.ImageURL != "" {
		content = item.ImageURL
	}

	return model.ConversationMessage{
		ID:          item.ID,
		GroupID:     item.GroupID,
		SenderID:    item.SenderID,
		Content:     content,
		MessageType: item.MessageType,
		IsMe:        item.SenderID == userID,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}

func parseBeforeMessageID(c *gin.Context) *uint64 {
	raw := strings.TrimSpace(c.Query("beforeId"))
	if raw == "" {
		return nil
	}

	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || parsed == 0 {
		return nil
	}
	return &parsed
}
