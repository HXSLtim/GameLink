package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	chatservice "gamelink/internal/service/chat"
)

// RegisterChatRoutes 注册用户端聊天相关路由。
func RegisterChatRoutes(router gin.IRouter, svc *chatservice.ChatService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/chat")
	group.Use(authMiddleware)
	group.GET("/groups", func(c *gin.Context) { listChatGroupsHandler(c, svc) })
	group.GET("/groups/:id/messages", func(c *gin.Context) { listChatMessagesHandler(c, svc) })
	group.POST("/groups/:id/messages", func(c *gin.Context) { sendChatMessageHandler(c, svc) })
	group.POST("/messages/:id/report", func(c *gin.Context) { reportChatMessageHandler(c, svc) })
}

type reportMessageRequest struct {
    Reason   string `json:"reason"`
    Evidence string `json:"evidence"`
}

func reportChatMessageHandler(c *gin.Context, svc *chatservice.ChatService) {
    userID := getUserIDFromContext(c)
    messageID, err := parseUintFromParam(c, "id")
    if err != nil {
        respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
        return
    }
    var req reportMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, err.Error())
        return
    }
    if err := svc.ReportMessage(c.Request.Context(), userID, messageID, req.Reason, req.Evidence); err != nil {
        respondError(c, http.StatusInternalServerError, err.Error())
        return
    }
    respondJSON(c, http.StatusOK, model.APIResponse[any]{
        Success: true,
        Code:    http.StatusOK,
        Message: "reported",
    })
}

func listChatGroupsHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	groups, total, err := svc.ListUserGroups(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"groups": groups,
			"total":  total,
		},
	})
}

func listChatMessagesHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintFromParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	var beforeID *uint64
	if val := strings.TrimSpace(c.Query("beforeId")); val != "" {
		if parsed, parseErr := strconv.ParseUint(val, 10, 64); parseErr == nil {
			beforeID = &parsed
		}
	}

	var afterID *uint64
	if val := strings.TrimSpace(c.Query("afterId")); val != "" {
		if parsed, parseErr := strconv.ParseUint(val, 10, 64); parseErr == nil {
			afterID = &parsed
		}
	}

	messages, total, err := svc.ListMessages(c.Request.Context(), userID, groupID, chatservice.ListMessagesOptions{
		Page:     page,
		PageSize: pageSize,
		BeforeID: beforeID,
		AfterID:  afterID,
	})
	if err != nil {
		switch err {
		case chatservice.ErrNotMember:
			respondError(c, http.StatusForbidden, err.Error())
		case chatservice.ErrInactiveGroup:
			respondError(c, http.StatusGone, err.Error())
		default:
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"messages": messages,
			"total":    total,
		},
	})
}

type sendMessageRequest struct {
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
	ImageURL    string `json:"imageUrl"`
	ReplyToID   *uint64 `json:"replyToId"`
}

func sendChatMessageHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintFromParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	messageType := model.ChatMessageTypeText
	if req.MessageType != "" {
		switch req.MessageType {
		case "text":
			messageType = model.ChatMessageTypeText
		case "image":
			messageType = model.ChatMessageTypeImage
		case "file":
			messageType = model.ChatMessageTypeFile
		case "system":
			messageType = model.ChatMessageTypeSystem
		default:
			respondError(c, http.StatusBadRequest, "unsupported message type")
			return
		}
	}

	msg, err := svc.SendMessage(c.Request.Context(), chatservice.SendMessageInput{
		GroupID:     groupID,
		SenderID:    userID,
		Content:     req.Content,
		MessageType: messageType,
		ReplyToID:   req.ReplyToID,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		switch err {
		case chatservice.ErrNotMember:
			respondError(c, http.StatusForbidden, err.Error())
		case chatservice.ErrInactiveGroup:
			respondError(c, http.StatusGone, err.Error())
		case chatservice.ErrMessageTooLarge:
			respondError(c, http.StatusBadRequest, err.Error())
		case chatservice.ErrThrottled:
			respondError(c, http.StatusTooManyRequests, err.Error())
		default:
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(c, http.StatusCreated, model.APIResponse[*model.ChatMessage]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    msg,
	})
}

func parseUintFromParam(c *gin.Context, name string) (uint64, error) {
	value := c.Param(name)
	return strconv.ParseUint(value, 10, 64)
}
