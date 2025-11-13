package notification

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
	notificationservice "gamelink/internal/service/notification"
)

// RegisterRoutes 注册通知中心路由。
func RegisterRoutes(router gin.IRouter, svc *notificationservice.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/notifications")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { listNotificationsHandler(c, svc) })
	group.POST("/read", func(c *gin.Context) { markNotificationsReadHandler(c, svc) })
	group.GET("/unread-count", func(c *gin.Context) { unreadCountHandler(c, svc) })
}

func listNotificationsHandler(c *gin.Context, svc *notificationservice.Service) {
	userID := getUserIDFromContext(c)
	var req struct {
		Page     int      `form:"page"`
		PageSize int      `form:"pageSize"`
		Unread   bool     `form:"unread"`
		Priority []string `form:"priority"`
	}
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

func markNotificationsReadHandler(c *gin.Context, svc *notificationservice.Service) {
	userID := getUserIDFromContext(c)
	var body struct {
		IDs []uint64 `json:"ids"`
	}
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
