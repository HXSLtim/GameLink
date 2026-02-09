package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/apierr"
)

// RegisterScheduleRoutes registers schedule routes.
func RegisterScheduleRoutes(router gin.IRouter, svc *playerservice.ScheduleService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/schedule")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { getScheduleHandler(c, svc) })
	group.PUT("", func(c *gin.Context) { updateScheduleHandler(c, svc) })
}

// getScheduleHandler 获取接单时间设置
// @Summary 获取接单时间设置
// @Tags 陪玩师-排班
// @Produce json
// @Security BearerAuth
// @Success 200 {object} playerservice.SchedulePayload
// @Router /player/schedule [get]
func getScheduleHandler(c *gin.Context, svc *playerservice.ScheduleService) {
	userID := resp.GetUserID(c)
	payload, err := svc.GetSchedule(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, payload)
}

// updateScheduleHandler 更新接单时间设置
// @Summary 更新接单时间设置
// @Tags 陪玩师-排班
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body playerservice.SchedulePayload true "schedule payload"
// @Success 200 {object} playerservice.SchedulePayload
// @Router /player/schedule [put]
func updateScheduleHandler(c *gin.Context, svc *playerservice.ScheduleService) {
	userID := resp.GetUserID(c)
	var payload playerservice.SchedulePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	updated, err := svc.UpdateSchedule(c.Request.Context(), userID, payload)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, updated)
}
