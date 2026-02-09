package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	notificationservice "gamelink/internal/service/notification"
	userservice "gamelink/internal/service/user"
	"gamelink/pkg/apierr"
)

// SettingsHandler handles user settings endpoints.
type SettingsHandler struct {
	settingsSvc     *userservice.SettingsService
	notificationSvc *notificationservice.SettingsService
}

// NewSettingsHandler creates settings handler.
func NewSettingsHandler(settingsSvc *userservice.SettingsService, notificationSvc *notificationservice.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsSvc:     settingsSvc,
		notificationSvc: notificationSvc,
	}
}

// RegisterSettingsRoutes registers settings routes.
func RegisterSettingsRoutes(rg *gin.RouterGroup, settingsSvc *userservice.SettingsService, notificationSvc *notificationservice.SettingsService, authMiddleware gin.HandlerFunc) {
	h := NewSettingsHandler(settingsSvc, notificationSvc)

	settings := rg.Group("/settings")
	settings.Use(authMiddleware)
	{
		settings.GET("", h.GetSettings)
		settings.PUT("", h.UpdateSettings)
	}

	notifySettings := rg.Group("/notification-settings")
	notifySettings.Use(authMiddleware)
	{
		notifySettings.GET("", h.GetNotificationSettings)
		notifySettings.PUT("", h.UpdateNotificationSettings)
	}
}

// GetSettings 获取用户设置
// @Summary 获取用户设置
// @Tags 用户-设置
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userservice.SettingsPayload
// @Router /user/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID := resp.GetUserID(c)
	payload, err := h.settingsSvc.GetSettings(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取用户设置失败").WithDetails(err.Error()))
		return
	}
	resp.OK(c, payload)
}

// UpdateSettings 更新用户设置
// @Summary 更新用户设置
// @Tags 用户-设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body userservice.SettingsPayload true "settings payload"
// @Success 200 {object} userservice.SettingsPayload
// @Router /user/settings [put]
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	userID := resp.GetUserID(c)
	var payload userservice.SettingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	updated, err := h.settingsSvc.UpdateSettings(c.Request.Context(), userID, payload)
	if err != nil {
		resp.Error(c, apierr.InternalError("更新用户设置失败").WithDetails(err.Error()))
		return
	}
	resp.OK(c, updated)
}

// GetNotificationSettings 获取通知设置
// @Summary 获取通知设置
// @Tags 用户-设置
// @Produce json
// @Security BearerAuth
// @Success 200 {object} notificationservice.SettingsPayload
// @Router /user/notification-settings [get]
func (h *SettingsHandler) GetNotificationSettings(c *gin.Context) {
	userID := resp.GetUserID(c)
	payload, err := h.notificationSvc.GetSettings(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取通知设置失败").WithDetails(err.Error()))
		return
	}
	resp.OK(c, payload)
}

// UpdateNotificationSettings 更新通知设置
// @Summary 更新通知设置
// @Tags 用户-设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body notificationservice.SettingsPayload true "notification settings payload"
// @Success 200 {object} notificationservice.SettingsPayload
// @Router /user/notification-settings [put]
func (h *SettingsHandler) UpdateNotificationSettings(c *gin.Context) {
	userID := resp.GetUserID(c)
	var payload notificationservice.SettingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	updated, err := h.notificationSvc.UpdateSettings(c.Request.Context(), userID, payload)
	if err != nil {
		resp.Error(c, apierr.InternalError("更新通知设置失败").WithDetails(err.Error()))
		return
	}
	resp.OK(c, updated)
}
