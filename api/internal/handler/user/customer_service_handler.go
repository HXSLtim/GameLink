package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
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

	conversations, total, err := h.conversationSvc.ListConversations(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	var messages []model.ConversationMessage
	if len(conversations) > 0 {
		primary := conversations[0]
		messages, _, err = h.conversationSvc.ListMessages(c.Request.Context(), userID, primary.GroupID, page, pageSize, beforeID)
		if err != nil {
			respondAPIError(c, err)
			return
		}
	}

	payload := model.ConversationListResponse{
		Conversations: conversations,
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

	conversation, message, err := h.conversationSvc.CreateConversation(c.Request.Context(), userID, req.Content)
	if err != nil {
		respondAPIError(c, err)
		return
	}

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

func (h *customerServiceHandler) listConversationMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	conversationID, err := parseUintParam(c, "id")
	if err != nil || conversationID == 0 {
		respondAPIError(c, apierr.BadRequest("invalid conversation id"))
		return
	}

	page, pageSize := parsePagination(c)
	beforeID := parseBeforeMessageID(c)
	messages, total, err := h.conversationSvc.ListMessages(c.Request.Context(), userID, conversationID, page, pageSize, beforeID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"conversationId": conversationID,
			"messages":       messages,
			"total":          total,
			"page":           page,
			"pageSize":       pageSize,
			"hasMore":        int64(page*pageSize) < total,
		},
	})
}

func (h *customerServiceHandler) sendConversationMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	conversationID, err := parseUintParam(c, "id")
	if err != nil || conversationID == 0 {
		respondAPIError(c, apierr.BadRequest("invalid conversation id"))
		return
	}

	var req sendCustomerServiceMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	message, err := h.conversationSvc.SendMessage(c.Request.Context(), userID, conversationID, req.Content)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[model.ConversationMessage](c, http.StatusCreated, model.APIResponse[model.ConversationMessage]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    *message,
	})
}

func (h *customerServiceHandler) closeConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("unauthorized"))
		return
	}

	conversationID, err := parseUintParam(c, "id")
	if err != nil || conversationID == 0 {
		respondAPIError(c, apierr.BadRequest("invalid conversation id"))
		return
	}

	conversation, err := h.conversationSvc.CloseConversation(c.Request.Context(), userID, conversationID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	respondJSON[model.Conversation](c, http.StatusOK, model.APIResponse[model.Conversation]{
		Success: true,
		Code:    http.StatusOK,
		Message: "closed",
		Data:    *conversation,
	})
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
