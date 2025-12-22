package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	contentservice "gamelink/internal/service/content"
)

// ListNotificationsRequest 获取通知列表请求
type ListNotificationsRequest struct {
	Page     int      `form:"page"`
	PageSize int      `form:"pageSize"`
	Unread   bool     `form:"unread"`
	Priority []string `form:"priority"`
}

// MarkNotificationsReadRequest 标记通知为已读请求
type MarkNotificationsReadRequest struct {
	IDs []uint64 `json:"ids"`
}

// NotificationListResponse 通知列表响应（类型别名）
type NotificationListResponse = contentservice.NotificationListResponse

// RegisterRoutes 注册通知中心路由。
func RegisterRoutes(router gin.IRouter, svc *contentservice.NotificationService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/notifications")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { listNotificationsHandler(c, svc) })
	group.POST("/read", func(c *gin.Context) { markNotificationsReadHandler(c, svc) })
	group.POST("/read-all", func(c *gin.Context) { markAllNotificationsReadHandler(c, svc) })
	group.GET("/unread-count", func(c *gin.Context) { unreadCountHandler(c, svc) })
	group.DELETE("/:id", func(c *gin.Context) { deleteNotificationHandler(c, svc) })
}

// listNotificationsHandler 获取通知列表
// @Summary      获取通知列表
// @Description  获取当前用户的通知列表，支持分页和过滤
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "Page number"
// @Param        pageSize       query     int     false  "Page size"
// @Param        unread         query     bool    false  "Filter unread only"
// @Param        priority       query     array   false  "Filter by priority"
// @Success      200            {object}  model.APIResponse[NotificationListResponse]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /notifications [get]
func listNotificationsHandler(c *gin.Context, svc *contentservice.NotificationService) {
	userID := getUserIDFromContext(c)
	var req ListNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	priorities := make([]model.NotificationPriority, 0, len(req.Priority))
	for _, p := range req.Priority {
		priorities = append(priorities, model.NotificationPriority(p))
	}
	resp, err := svc.ListNotifications(c.Request.Context(), userID, contentservice.NotificationListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		UnreadOnly: req.Unread,
		Priorities: priorities,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[*contentservice.NotificationListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    resp,
	})
}

// markAllNotificationsReadHandler 标记所有通知为已读
// @Summary      标记所有通知为已读
// @Description  将当前用户的所有未读通知标记为已读
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string           true  "Bearer {token}"
// @Success      200            {object}  model.SuccessResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /notifications/read-all [post]
func markAllNotificationsReadHandler(c *gin.Context, svc *contentservice.NotificationService) {
	userID := getUserIDFromContext(c)
	if err := svc.MarkAllNotificationsRead(c.Request.Context(), userID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "所有通知已标记为已读",
	})
}

// markNotificationsReadHandler 标记通知为已读
// @Summary      标记通知为已读
// @Description  批量标记指定通知为已读状态
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      MarkNotificationsReadRequest   true  "Notification IDs to mark as read"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /notifications/read [post]
func markNotificationsReadHandler(c *gin.Context, svc *contentservice.NotificationService) {
	userID := getUserIDFromContext(c)
	var body MarkNotificationsReadRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := svc.MarkNotificationsRead(c.Request.Context(), userID, body.IDs); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "已标记为已读",
	})
}

// unreadCountHandler 获取未读通知数量
// @Summary      获取未读通知数量
// @Description  获取当前用户的未读通知总数
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[map[string]int64]
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /notifications/unread-count [get]
func unreadCountHandler(c *gin.Context, svc *contentservice.NotificationService) {
	userID := getUserIDFromContext(c)
	count, err := svc.GetUnreadNotificationCount(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[map[string]int64]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]int64{"unread": count},
	})
}

// deleteNotificationHandler 删除通知
// @Summary      删除通知
// @Description  删除指定ID的通知
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      uint64  true  "Notification ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /notifications/{id} [delete]
func deleteNotificationHandler(c *gin.Context, svc *contentservice.NotificationService) {
	userID := getUserIDFromContext(c)

	var uri struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := svc.DeleteNotification(c.Request.Context(), userID, uri.ID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "通知已删除",
	})
}

func getUserIDFromContext(c *gin.Context) uint64 {
	return resp.GetUserID(c)
}

// respondJSON is an alias for resp.JSON for backward compatibility.
func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	resp.JSON(c, status, payload)
}

// respondError is an alias for resp.ErrorMsg for backward compatibility.
func respondError(c *gin.Context, status int, msg string) {
	resp.ErrorMsg(c, status, msg)
}
