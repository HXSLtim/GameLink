package admin

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// AdminChatHandler handles admin chat management endpoints.
type AdminChatHandler struct {
	groups   repository.ChatGroupRepository
	members  repository.ChatMemberRepository
	messages repository.ChatMessageRepository
	users    repository.UserRepository
}

// NewAdminChatHandler creates a new admin chat handler.
func NewAdminChatHandler(
	groups repository.ChatGroupRepository,
	members repository.ChatMemberRepository,
	messages repository.ChatMessageRepository,
	users repository.UserRepository,
) *AdminChatHandler {
	return &AdminChatHandler{
		groups:   groups,
		members:  members,
		messages: messages,
		users:    users,
	}
}

// conversationResponse maps a ChatGroup to the frontend's "conversation" format.
type conversationResponse struct {
	ID                 uint64  `json:"id"`
	Type               string  `json:"type"`
	OrderID            *uint64 `json:"orderId,omitempty"`
	OrderNo            string  `json:"orderNo,omitempty"`
	UserID             *uint64 `json:"userId,omitempty"`
	UserName           string  `json:"userName,omitempty"`
	UserAvatar         string  `json:"userAvatar,omitempty"`
	PlayerID           *uint64 `json:"playerId,omitempty"`
	PlayerName         string  `json:"playerName,omitempty"`
	PlayerAvatar       string  `json:"playerAvatar,omitempty"`
	MessageCount       int     `json:"messageCount"`
	LastMessageAt      *string `json:"lastMessageAt,omitempty"`
	LastMessageContent string  `json:"lastMessageContent,omitempty"`
	Status             string  `json:"status"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt,omitempty"`
}

func (h *AdminChatHandler) mapGroupToConversation(g model.ChatGroup) conversationResponse {
	convType := "group_chat"
	if g.GroupType == model.ChatGroupTypeOrder || g.GroupType == model.ChatGroupTypePrivate {
		convType = "user_order"
	}

	status := "active"
	if !g.IsActive {
		status = "closed"
	}

	cr := conversationResponse{
		ID:        g.ID,
		Type:      convType,
		Status:    status,
		CreatedAt: g.CreatedAt.Format(time.RFC3339),
	}
	if !g.UpdatedAt.IsZero() {
		t := g.UpdatedAt.Format(time.RFC3339)
		cr.UpdatedAt = t
	}
	if g.RelatedOrderID != nil {
		cr.OrderID = g.RelatedOrderID
	}

	// Extract user/player from members
	for _, m := range g.Members {
		if m.Role == model.ChatMemberRoleOwner {
			uid := m.UserID
			cr.UserID = &uid
			cr.UserName = m.Nickname
		} else if m.Role == model.ChatMemberRoleMember && cr.PlayerID == nil {
			pid := m.UserID
			cr.PlayerID = &pid
			cr.PlayerName = m.Nickname
		}
	}

	cr.MessageCount = len(g.Members) // placeholder; real count can be expensive

	return cr
}

// ListConversations returns a paginated list of all chat conversations.
// GET /admin/chat/conversations
func (h *AdminChatHandler) ListConversations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	opts := repository.AdminChatGroupListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	// Type filter
	if t := c.Query("type"); t != "" {
		switch t {
		case "user_order":
			gt := model.ChatGroupTypeOrder
			opts.GroupType = &gt
		case "group_chat":
			gt := model.ChatGroupTypePublic
			opts.GroupType = &gt
		}
	}

	// Status filter
	if s := c.Query("status"); s != "" {
		switch s {
		case "active":
			b := true
			opts.IsActive = &b
		case "closed":
			b := false
			opts.IsActive = &b
		}
	}

	// Keyword filter
	if k := c.Query("keyword"); k != "" {
		opts.Keyword = k
	}

	// Order ID filter
	if oid := c.Query("orderId"); oid != "" {
		if id, err := strconv.ParseUint(oid, 10, 64); err == nil {
			opts.RelatedOrderID = &id
		}
	}

	// User ID filter
	if uid := c.Query("userId"); uid != "" {
		if id, err := strconv.ParseUint(uid, 10, 64); err == nil {
			opts.UserID = &id
		}
	}

	groups, total, err := h.groups.ListAll(c.Request.Context(), opts)
	if err != nil {
		resp.InternalError(c, "加载会话列表失败")
		return
	}

	items := make([]conversationResponse, 0, len(groups))
	for _, g := range groups {
		items = append(items, h.mapGroupToConversation(g))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	resp.OK(c, gin.H{
		"items": items,
		"total": total,
		"pagination": model.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      int(total),
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	})
}

// GetConversation returns detail for a single conversation.
// GET /admin/chat/conversations/:id
func (h *AdminChatHandler) GetConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.BadRequest(c, "无效的会话ID")
		return
	}

	group, err := h.groups.GetWithRelations(c.Request.Context(), id)
	if err != nil {
		resp.NotFound(c, "会话不存在")
		return
	}

	resp.OK(c, h.mapGroupToConversation(*group))
}

// CloseConversation deactivates a conversation.
// POST /admin/chat/conversations/:id/close
func (h *AdminChatHandler) CloseConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.BadRequest(c, "无效的会话ID")
		return
	}

	if err := h.groups.Deactivate(c.Request.Context(), id); err != nil {
		resp.InternalError(c, "关闭会话失败")
		return
	}

	resp.OK(c, gin.H{"message": "会话已关闭"})
}

// ReopenConversation reactivates a closed conversation.
// POST /admin/chat/conversations/:id/reopen
func (h *AdminChatHandler) ReopenConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.BadRequest(c, "无效的会话ID")
		return
	}

	if err := h.groups.Reactivate(c.Request.Context(), id); err != nil {
		resp.InternalError(c, "重新打开会话失败")
		return
	}

	resp.OK(c, gin.H{"message": "会话已重新打开"})
}

// BatchCloseConversations deactivates multiple conversations.
// POST /admin/chat/conversations/batch-close
func (h *AdminChatHandler) BatchCloseConversations(c *gin.Context) {
	var req struct {
		ConversationIds []uint64 `json:"conversationIds"`
		Reason          string   `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "无效的请求参数")
		return
	}

	ctx := c.Request.Context()
	var failed int
	for _, id := range req.ConversationIds {
		if err := h.groups.Deactivate(ctx, id); err != nil {
			failed++
		}
	}

	resp.OK(c, gin.H{
		"total":   len(req.ConversationIds),
		"success": len(req.ConversationIds) - failed,
		"failed":  failed,
	})
}

// ListConversationMessages returns messages for a specific conversation.
// GET /admin/chat/conversations/:id/messages
func (h *AdminChatHandler) ListConversationMessages(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.BadRequest(c, "无效的会话ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	opts := repository.ChatMessageListOptions{
		GroupID:  groupID,
		Page:     page,
		PageSize: pageSize,
	}

	messages, total, err := h.messages.ListByGroup(c.Request.Context(), opts)
	if err != nil {
		resp.InternalError(c, "加载消息列表失败")
		return
	}

	resp.OK(c, gin.H{
		"items": messages,
		"total": total,
	})
}

// ListAllMessages returns all messages across conversations (for moderation).
// GET /admin/chat/messages
func (h *AdminChatHandler) ListAllMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var groupID *uint64
	if gid := c.Query("conversationId"); gid != "" {
		if id, err := strconv.ParseUint(gid, 10, 64); err == nil {
			groupID = &id
		}
	}

	var senderID *uint64
	if sid := c.Query("senderId"); sid != "" {
		if id, err := strconv.ParseUint(sid, 10, 64); err == nil {
			senderID = &id
		}
	}

	opts := repository.ChatMessageModerationListOptions{
		Page:     page,
		PageSize: pageSize,
		GroupID:  groupID,
		SenderID: senderID,
	}

	messages, total, err := h.messages.ListForModeration(c.Request.Context(), opts)
	if err != nil {
		resp.InternalError(c, "加载消息列表失败")
		return
	}

	resp.OK(c, gin.H{
		"items": messages,
		"total": total,
	})
}

// DeleteMessage soft-deletes a chat message.
// DELETE /admin/chat/messages/:id
func (h *AdminChatHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.BadRequest(c, "无效的消息ID")
		return
	}

	// Get admin user ID from context
	adminID := extractUserID(c)

	if err := h.messages.MarkDeleted(c.Request.Context(), id, adminID); err != nil {
		resp.InternalError(c, "删除消息失败")
		return
	}

	resp.Deleted(c)
}

// BatchDeleteMessages soft-deletes multiple messages.
// POST /admin/chat/messages/batch-delete
func (h *AdminChatHandler) BatchDeleteMessages(c *gin.Context) {
	var req struct {
		MessageIds []uint64 `json:"messageIds"`
		Reason     string   `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "无效的请求参数")
		return
	}

	adminID := extractUserID(c)

	ctx := c.Request.Context()
	var failed int
	for _, id := range req.MessageIds {
		if err := h.messages.MarkDeleted(ctx, id, adminID); err != nil {
			failed++
		}
	}

	resp.OK(c, gin.H{
		"total":   len(req.MessageIds),
		"success": len(req.MessageIds) - failed,
		"failed":  failed,
	})
}

// MuteUser mutes a user in a conversation.
// POST /admin/chat/users/mute
func (h *AdminChatHandler) MuteUser(c *gin.Context) {
	var req struct {
		ConversationID uint64 `json:"conversationId" binding:"required"`
		UserID         uint64 `json:"userId" binding:"required"`
		Duration       int    `json:"duration"` // minutes
		Reason         string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "无效的请求参数")
		return
	}

	adminID := extractUserID(c)

	member, err := h.members.Get(c.Request.Context(), req.ConversationID, req.UserID)
	if err != nil {
		resp.NotFound(c, "用户不是该会话的成员")
		return
	}

	member.IsMuted = true
	muteUntil := time.Now().Add(time.Duration(req.Duration) * time.Minute)
	member.MutedUntil = &muteUntil
	member.MutedBy = &adminID
	member.MuteReason = req.Reason

	if err := h.members.Update(c.Request.Context(), member); err != nil {
		resp.InternalError(c, "禁言失败")
		return
	}

	resp.OK(c, gin.H{"message": "禁言成功"})
}

// UnmuteUser unmutes a user in a conversation.
// POST /admin/chat/users/unmute
func (h *AdminChatHandler) UnmuteUser(c *gin.Context) {
	var req struct {
		ConversationID uint64 `json:"conversationId" binding:"required"`
		UserID         uint64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "无效的请求参数")
		return
	}

	member, err := h.members.Get(c.Request.Context(), req.ConversationID, req.UserID)
	if err != nil {
		resp.NotFound(c, "用户不是该会话的成员")
		return
	}

	member.IsMuted = false
	member.MutedUntil = nil
	member.MutedBy = nil
	member.MuteReason = ""

	if err := h.members.Update(c.Request.Context(), member); err != nil {
		resp.InternalError(c, "解除禁言失败")
		return
	}

	resp.OK(c, gin.H{"message": "已解除禁言"})
}

// GetChatStats returns chat statistics overview.
// GET /admin/chat/stats
func (h *AdminChatHandler) GetChatStats(c *gin.Context) {
	ctx := c.Request.Context()

	totalGroups, _ := h.groups.CountAll(ctx)

	// Count active groups
	activeTrue := true
	activeOpts := repository.AdminChatGroupListOptions{
		IsActive: &activeTrue,
		Page:     1,
		PageSize: 1,
	}
	_, activeCount, _ := h.groups.ListAll(ctx, activeOpts)

	resp.OK(c, gin.H{
		"overview": gin.H{
			"totalConversations":  totalGroups,
			"activeConversations": activeCount,
			"totalMessages":       0,
			"todayMessages":       0,
			"totalUsers":          0,
			"onlineUsers":         0,
		},
		"trends": []gin.H{},
	})
}

// extractUserID extracts user_id from gin context.
func extractUserID(c *gin.Context) uint64 {
	if uid, ok := c.Get("user_id"); ok {
		switch v := uid.(type) {
		case float64:
			return uint64(v)
		case uint64:
			return v
		case int:
			return uint64(v)
		case int64:
			return uint64(v)
		}
	}
	return 0
}

// RegisterAdminChatRoutes registers admin chat management routes.
func RegisterAdminChatRoutes(
	group gin.IRouter,
	chatHandler *AdminChatHandler,
	pm interface {
		RequirePermission(method model.HTTPMethod, path string) gin.HandlerFunc
	},
) {
	// Conversation management
	group.GET("/chat/conversations", chatHandler.ListConversations)
	group.GET("/chat/conversations/:id", chatHandler.GetConversation)
	group.POST("/chat/conversations/:id/close", chatHandler.CloseConversation)
	group.POST("/chat/conversations/:id/reopen", chatHandler.ReopenConversation)
	group.POST("/chat/conversations/batch-close", chatHandler.BatchCloseConversations)

	// Message management
	group.GET("/chat/conversations/:id/messages", chatHandler.ListConversationMessages)
	group.GET("/chat/messages", chatHandler.ListAllMessages)
	group.DELETE("/chat/messages/:id", chatHandler.DeleteMessage)
	group.POST("/chat/messages/batch-delete", chatHandler.BatchDeleteMessages)

	// User management (mute/unmute)
	group.POST("/chat/users/mute", chatHandler.MuteUser)
	group.POST("/chat/users/unmute", chatHandler.UnmuteUser)

	// Stats
	group.GET("/chat/stats", chatHandler.GetChatStats)
	group.GET("/chat/stats/conversations", chatHandler.GetChatStats)
	group.GET("/chat/stats/messages", chatHandler.GetChatStats)
	group.GET("/chat/stats/user-activity", chatHandler.GetChatStats)
}
