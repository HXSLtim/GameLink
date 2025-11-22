package notification

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
	notificationservice "gamelink/internal/service/notification"
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
type NotificationListResponse = notificationservice.ListResponse

// RegisterRoutes 注册通知中心路由。
func RegisterRoutes(router gin.IRouter, svc *notificationservice.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/notifications")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { listNotificationsHandler(c, svc) })
	group.POST("/read", func(c *gin.Context) { markNotificationsReadHandler(c, svc) })
	group.GET("/unread-count", func(c *gin.Context) { unreadCountHandler(c, svc) })
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
func listNotificationsHandler(c *gin.Context, svc *notificationservice.Service) {
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
	resp, err := svc.List(c.Request.Context(), userID, notificationservice.ListRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		UnreadOnly: req.Unread,
		Priorities: priorities,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondError(c, http.StatusBadRequest, err.Error())
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(c, http.StatusOK, model.APIResponse[*notificationservice.ListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    resp,
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
func markNotificationsReadHandler(c *gin.Context, svc *notificationservice.Service) {
	userID := getUserIDFromContext(c)
	var body MarkNotificationsReadRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := svc.MarkRead(c.Request.Context(), userID, body.IDs); err != nil {
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
func unreadCountHandler(c *gin.Context, svc *notificationservice.Service) {
	userID := getUserIDFromContext(c)
	count, err := svc.GetUnreadCount(c.Request.Context(), userID)
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

func getUserIDFromContext(c *gin.Context) uint64 {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		return 0
	}
	return userID
}

func respondJSON[T any](c *gin.Context, status int, payload model.APIResponse[T]) {
	if payload.TraceID == "" {
		if rid, ok := c.Get("request_id"); ok {
			if ridStr, ok := rid.(string); ok {
				payload.TraceID = ridStr
			}
		}
	}
	c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, msg string) {
	respondJSON(c, status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: msg,
	})
}
