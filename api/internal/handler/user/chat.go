package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/apierr"
)

// ChatMessage 聊天消息模型（类型别名）
type ChatMessage = model.ChatMessage

type chatGroupSummaryResponse struct {
	ID             uint64                `json:"id"`
	GroupType      model.ChatGroupType   `json:"groupType"`
	Type           model.ChatGroupType   `json:"type"`
	Name           string                `json:"name"`
	GroupName      string                `json:"groupName"`
	Avatar         string                `json:"avatar,omitempty"`
	AvatarURL      string                `json:"avatarUrl,omitempty"`
	Description    string                `json:"description,omitempty"`
	OrderID        *uint64               `json:"orderId,omitempty"`
	RelatedOrderID *uint64               `json:"relatedOrderId,omitempty"`
	MemberCount    int                   `json:"memberCount"`
	UnreadCount    int                   `json:"unreadCount"`
	TargetUserID   *uint64               `json:"targetUserId,omitempty"`
	TargetNickname string                `json:"targetNickname,omitempty"`
	TargetIsOnline bool                  `json:"targetIsOnline"`
	LastMessage    *chatMessageItem      `json:"lastMessage,omitempty"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
	IsActive       bool                  `json:"isActive"`
	RoomStatus     model.ChatGroupStatus `json:"roomStatus"`
}

type chatMessageItem struct {
	ID             uint64                `json:"id"`
	GroupID        uint64                `json:"groupId"`
	SenderID       uint64                `json:"senderId"`
	SenderNickname string                `json:"senderNickname,omitempty"`
	SenderAvatar   string                `json:"senderAvatar,omitempty"`
	MessageType    model.ChatMessageType `json:"messageType"`
	Content        string                `json:"content"`
	ImageURL       string                `json:"imageUrl,omitempty"`
	Status         string                `json:"status"`
	OrderID        *uint64               `json:"orderId,omitempty"`
	CreatedAt      string                `json:"createdAt"`
}

type chatGroupMemberItem struct {
	ID                uint64               `json:"id"`
	GroupID           uint64               `json:"groupId"`
	UserID            uint64               `json:"userId"`
	Role              model.ChatMemberRole `json:"role"`
	Nickname          string               `json:"nickname,omitempty"`
	JoinedAt          string               `json:"joinedAt"`
	LastReadAt        *string              `json:"lastReadAt,omitempty"`
	LastReadMessageID *uint64              `json:"lastReadMessageId,omitempty"`
	IsMuted           bool                 `json:"isMuted"`
	MutedUntil        *string              `json:"mutedUntil,omitempty"`
	MuteReason        string               `json:"muteReason,omitempty"`
	IsActive          bool                 `json:"isActive"`
	IsReady           bool                 `json:"isReady"`
}

type chatGroupDetailResponse struct {
	chatGroupSummaryResponse
	CreatedBy      uint64                `json:"createdBy"`
	MaxMembers     int                   `json:"maxMembers"`
	CurrentMembers int                   `json:"currentMembers"`
	AutoDestroy    bool                  `json:"autoDestroy"`
	DeactivatedAt  *string               `json:"deactivatedAt,omitempty"`
	GameID         *uint64               `json:"gameId,omitempty"`
	IsPrivate      bool                  `json:"isPrivate"`
	RelatedTeamID  *uint64               `json:"relatedTeamId,omitempty"`
	RelatedLFGID   *uint64               `json:"relatedLfgId,omitempty"`
	StartedAt      *string               `json:"startedAt,omitempty"`
	FinishedAt     *string               `json:"finishedAt,omitempty"`
	VoiceEnabled   bool                  `json:"voiceEnabled"`
	VoiceRoomID    string                `json:"voiceRoomId,omitempty"`
	VoiceProvider  string                `json:"voiceProvider,omitempty"`
	VoiceSDKAppID  uint64                `json:"voiceSdkAppId"`
	VoiceStartedAt *string               `json:"voiceStartedAt,omitempty"`
	VoiceEndedAt   *string               `json:"voiceEndedAt,omitempty"`
	VoiceDuration  int                   `json:"voiceDuration"`
	VoiceRecordURL string                `json:"voiceRecordUrl,omitempty"`
	VoiceMaxMember int                   `json:"voiceMaxMembers"`
	Members        []chatGroupMemberItem `json:"members"`
}

// RegisterChatRoutes 注册用户端聊天相关路由。
func RegisterChatRoutes(router gin.IRouter, svc *chatservice.ChatService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/chat")
	group.Use(authMiddleware)
	group.POST("/groups", func(c *gin.Context) { createChatGroupHandler(c, svc) })
	group.GET("/groups", func(c *gin.Context) { listChatGroupsHandler(c, svc) })
	group.GET("/groups/:id", func(c *gin.Context) { getChatGroupHandler(c, svc) })
	group.POST("/groups/:id/join", func(c *gin.Context) { joinChatGroupHandler(c, svc) })
	group.POST("/groups/:id/leave", func(c *gin.Context) { leaveChatGroupHandler(c, svc) })
	group.POST("/groups/:id/read", func(c *gin.Context) { markChatGroupReadHandler(c, svc) })
	group.GET("/groups/:id/messages", func(c *gin.Context) { listChatMessagesHandler(c, svc) })
	group.POST("/groups/:id/messages", func(c *gin.Context) { sendChatMessageHandler(c, svc) })
	group.POST("/messages/:id/report", func(c *gin.Context) { reportChatMessageHandler(c, svc) })
}

type createChatGroupRequest struct {
	TargetUserID uint64  `json:"targetUserId" binding:"required"`
	GroupType    string  `json:"groupType" binding:"required"` // private | order
	OrderID      *uint64 `json:"orderId,omitempty"`
}

// createChatGroupHandler 创建私聊群组
// @Summary      创建私聊群组
// @Description  创建与陪玩师私聊或订单群组
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                 true  "Bearer {token}"
// @Param        request        body      createChatGroupRequest true  "Create group request"
// @Success      200            {object}  chatGroupDetailResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups [post]
func createChatGroupHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	var req createChatGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	groupType := model.ChatGroupType(strings.ToLower(strings.TrimSpace(req.GroupType)))
	if groupType != model.ChatGroupTypePrivate && groupType != model.ChatGroupTypeOrder {
		respondError(c, http.StatusBadRequest, "unsupported group type")
		return
	}

	group, err := svc.CreateGroup(c.Request.Context(), userID, req.TargetUserID, chatservice.CreateGroupRequest{
		TargetUserID: req.TargetUserID,
		GroupType:    groupType,
		OrderID:      req.OrderID,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	groupData := toChatGroupDetail(group, userID)

	respondJSON[chatGroupDetailResponse](c, http.StatusOK, model.APIResponse[chatGroupDetailResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    groupData,
	})
}

// getChatGroupHandler 获取群组详情
// @Summary      获取群组详情
// @Description  获取群组详情（包含成员列表）
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        id             path      int     true   "Group ID"
// @Success      200            {object}  chatGroupDetailResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      403            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id} [get]
func getChatGroupHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	if _, err := svc.EnsureMembership(c.Request.Context(), groupID, userID); err != nil {
		respondError(c, http.StatusForbidden, err.Error())
		return
	}
	group, err := svc.GetGroupDetail(c.Request.Context(), groupID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	groupData := toChatGroupDetail(group, userID)
	respondJSON[chatGroupDetailResponse](c, http.StatusOK, model.APIResponse[chatGroupDetailResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    groupData,
	})
}

// joinChatGroupHandler 加入公共频道
// @Summary      加入公共频道
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        id             path      int     true   "Group ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id}/join [post]
func joinChatGroupHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	if err := svc.JoinGroupWithUser(c.Request.Context(), groupID, userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "joined",
	})
}

// leaveChatGroupHandler 离开群组
// @Summary      离开群组
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        id             path      int     true   "Group ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id}/leave [post]
func leaveChatGroupHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	if err := svc.LeaveGroup(c.Request.Context(), groupID, userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "left",
	})
}

type markChatReadRequest struct {
	MessageID uint64 `json:"messageId" binding:"required"`
}

// markChatGroupReadHandler 标记消息已读
// @Summary      标记消息已读
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string              true  "Bearer {token}"
// @Param        id             path      int                 true  "Group ID"
// @Param        request        body      markChatReadRequest true  "Last read message id"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id}/read [post]
func markChatGroupReadHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var req markChatReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := svc.MarkRead(c.Request.Context(), groupID, userID, req.MessageID); err != nil {
		resp.Error(c, err)
		return
	}
	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "read",
	})
}

type reportMessageRequest struct {
	Reason   string `json:"reason"`
	Evidence string `json:"evidence"`
}

// reportChatMessageHandler 举报聊天消息
// @Summary      举报聊天消息
// @Description  举报不当聊天消息，包括违规内容、骚扰等
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer {token}"
// @Param        id             path      int                      true  "Message ID"
// @Param        request        body      reportMessageRequest     true  "Report reason and evidence"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/messages/{id}/report [post]
func reportChatMessageHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	messageID, err := parseUintParam(c, "id")
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
	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "reported",
	})
}

// listChatGroupsHandler 获取聊天群组列表
// @Summary      获取聊天群组列表
// @Description  获取当前用户的所有聊天群组
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "Page number"
// @Param        pageSize       query     int     false  "Page size"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups [get]
func listChatGroupsHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	groups, total, err := svc.ListUserGroups(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]chatGroupSummaryResponse, 0, len(groups))
	for i := range groups {
		items = append(items, toChatGroupSummary(&groups[i], userID))
	}

	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"groups": items,
			"items":  items,
			"total":  total,
		},
	})
}

// listChatMessagesHandler 获取聊天消息列表
// @Summary      获取聊天消息列表
// @Description  获取指定聊天群组的消息列表，支持分页和前后翻页
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        id             path      int     true   "Group ID"
// @Param        page           query     int     false  "Page number"
// @Param        pageSize       query     int     false  "Page size"
// @Param        beforeId       query     int     false  "Load messages before this ID"
// @Param        afterId        query     int     false  "Load messages after this ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      403            {object}  model.ErrorResponse
// @Failure      410            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id}/messages [get]
func listChatMessagesHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
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
	groupDetail, _ := svc.GetGroupDetail(c.Request.Context(), groupID)
	messageItems := make([]chatMessageItem, 0, len(messages))
	for i := range messages {
		messageItems = append(messageItems, toChatMessageItem(&messages[i], groupDetail))
	}

	respondJSON[any](c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"messages": messageItems,
			"items":    messageItems,
			"total":    total,
		},
	})
}

type sendMessageRequest struct {
	Content     string  `json:"content"`
	MessageType string  `json:"messageType"`
	ImageURL    string  `json:"imageUrl"`
	ReplyToID   *uint64 `json:"replyToId"`
}

// sendChatMessageHandler 发送聊天消息
// @Summary      发送聊天消息
// @Description  在指定聊天群组中发送消息，支持文本、图片、文件等类型
// @Tags         User - Chat
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string               true  "Bearer {token}"
// @Param        id             path      int                  true  "Group ID"
// @Param        request        body      sendMessageRequest   true  "Message content"
// @Success      201            {object}  ChatMessage
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      403            {object}  model.ErrorResponse
// @Failure      410            {object}  model.ErrorResponse
// @Failure      429            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /user/chat/groups/{id}/messages [post]
func sendChatMessageHandler(c *gin.Context, svc *chatservice.ChatService) {
	userID := getUserIDFromContext(c)
	groupID, err := parseUintParam(c, "id")
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
	groupDetail, _ := svc.GetGroupDetail(c.Request.Context(), groupID)
	item := toChatMessageItem(msg, groupDetail)

	respondJSON[chatMessageItem](c, http.StatusCreated, model.APIResponse[chatMessageItem]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    item,
	})
}

func toChatGroupSummary(group *model.ChatGroup, currentUserID uint64) chatGroupSummaryResponse {
	if group == nil {
		return chatGroupSummaryResponse{}
	}
	name := strings.TrimSpace(group.GroupName)
	if name == "" {
		name = "聊天"
	}

	memberCount := 0
	var targetUserID *uint64
	targetNickname := ""
	for _, member := range group.Members {
		if member.IsActive {
			memberCount++
		}
		if member.UserID == currentUserID {
			continue
		}
		if targetUserID == nil {
			id := member.UserID
			targetUserID = &id
			targetNickname = strings.TrimSpace(member.Nickname)
		}
	}

	if targetNickname == "" && targetUserID != nil {
		targetNickname = "用户" + strconv.FormatUint(*targetUserID, 10)
	}

	createdAt := group.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	updatedAt := group.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")

	return chatGroupSummaryResponse{
		ID:             group.ID,
		GroupType:      group.GroupType,
		Type:           group.GroupType,
		Name:           name,
		GroupName:      name,
		Avatar:         group.AvatarURL,
		AvatarURL:      group.AvatarURL,
		Description:    group.Description,
		OrderID:        group.RelatedOrderID,
		RelatedOrderID: group.RelatedOrderID,
		MemberCount:    memberCount,
		UnreadCount:    0,
		TargetUserID:   targetUserID,
		TargetNickname: targetNickname,
		TargetIsOnline: false,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		IsActive:       group.IsActive,
		RoomStatus:     group.RoomStatus,
	}
}

func toChatMessageItem(msg *model.ChatMessage, group *model.ChatGroup) chatMessageItem {
	if msg == nil {
		return chatMessageItem{}
	}
	content := strings.TrimSpace(msg.Content)
	if content == "" && msg.ImageURL != "" {
		content = msg.ImageURL
	}

	senderNickname := ""
	if group != nil {
		for _, member := range group.Members {
			if member.UserID == msg.SenderID {
				senderNickname = strings.TrimSpace(member.Nickname)
				break
			}
		}
	}
	if senderNickname == "" {
		senderNickname = "用户" + strconv.FormatUint(msg.SenderID, 10)
	}

	status := "sent"
	if msg.IsDeleted {
		status = "deleted"
	} else if msg.AuditStatus == model.ChatMessageAuditPending {
		status = "pending"
	}

	var orderID *uint64
	if group != nil && group.RelatedOrderID != nil {
		id := *group.RelatedOrderID
		orderID = &id
	}

	return chatMessageItem{
		ID:             msg.ID,
		GroupID:        msg.GroupID,
		SenderID:       msg.SenderID,
		SenderNickname: senderNickname,
		SenderAvatar:   "",
		MessageType:    msg.MessageType,
		Content:        content,
		ImageURL:       msg.ImageURL,
		Status:         status,
		OrderID:        orderID,
		CreatedAt:      msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toChatGroupDetail(group *model.ChatGroup, currentUserID uint64) chatGroupDetailResponse {
	if group == nil {
		return chatGroupDetailResponse{}
	}

	summary := toChatGroupSummary(group, currentUserID)
	members := make([]chatGroupMemberItem, 0, len(group.Members))
	for i := range group.Members {
		member := group.Members[i]
		var lastReadAt *string
		if member.LastReadAt != nil {
			value := member.LastReadAt.Format("2006-01-02T15:04:05Z07:00")
			lastReadAt = &value
		}
		var mutedUntil *string
		if member.MutedUntil != nil {
			value := member.MutedUntil.Format("2006-01-02T15:04:05Z07:00")
			mutedUntil = &value
		}
		members = append(members, chatGroupMemberItem{
			ID:                member.ID,
			GroupID:           member.GroupID,
			UserID:            member.UserID,
			Role:              member.Role,
			Nickname:          strings.TrimSpace(member.Nickname),
			JoinedAt:          member.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastReadAt:        lastReadAt,
			LastReadMessageID: member.LastReadMessageID,
			IsMuted:           member.IsMuted,
			MutedUntil:        mutedUntil,
			MuteReason:        member.MuteReason,
			IsActive:          member.IsActive,
			IsReady:           member.IsReady,
		})
	}

	var deactivatedAt *string
	if group.DeactivatedAt != nil {
		value := group.DeactivatedAt.Format("2006-01-02T15:04:05Z07:00")
		deactivatedAt = &value
	}
	var startedAt *string
	if group.StartedAt != nil {
		value := group.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		startedAt = &value
	}
	var finishedAt *string
	if group.FinishedAt != nil {
		value := group.FinishedAt.Format("2006-01-02T15:04:05Z07:00")
		finishedAt = &value
	}
	var voiceStartedAt *string
	if group.VoiceStartedAt != nil {
		value := group.VoiceStartedAt.Format("2006-01-02T15:04:05Z07:00")
		voiceStartedAt = &value
	}
	var voiceEndedAt *string
	if group.VoiceEndedAt != nil {
		value := group.VoiceEndedAt.Format("2006-01-02T15:04:05Z07:00")
		voiceEndedAt = &value
	}

	return chatGroupDetailResponse{
		chatGroupSummaryResponse: summary,
		CreatedBy:                group.CreatedBy,
		MaxMembers:               group.MaxMembers,
		CurrentMembers:           group.CurrentMembers,
		AutoDestroy:              group.AutoDestroy,
		DeactivatedAt:            deactivatedAt,
		GameID:                   group.GameID,
		IsPrivate:                group.IsPrivate,
		RelatedTeamID:            group.RelatedTeamID,
		RelatedLFGID:             group.RelatedLFGID,
		StartedAt:                startedAt,
		FinishedAt:               finishedAt,
		VoiceEnabled:             group.VoiceEnabled,
		VoiceRoomID:              group.VoiceRoomID,
		VoiceProvider:            group.VoiceProvider,
		VoiceSDKAppID:            group.VoiceSDKAppID,
		VoiceStartedAt:           voiceStartedAt,
		VoiceEndedAt:             voiceEndedAt,
		VoiceDuration:            group.VoiceDuration,
		VoiceRecordURL:           group.VoiceRecordURL,
		VoiceMaxMember:           group.VoiceMaxMembers,
		Members:                  members,
	}
}
