package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	playerservice "gamelink/internal/service/player"
)

// RegisterStatsRoutes registers stats routes.
func RegisterStatsRoutes(router gin.IRouter, svc *playerservice.PlayerService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/stats")
	group.Use(authMiddleware)
	group.GET("/today", func(c *gin.Context) { getTodayStatsHandler(c, svc) })
	group.GET("/overview", func(c *gin.Context) { getOverviewStatsHandler(c, svc) })
}

// getTodayStatsHandler 获取今日统计
// @Summary 获取今日统计
// @Tags 陪玩师-统计
// @Produce json
// @Security BearerAuth
// @Success 200 {object} playerservice.PlayerStatsToday
// @Router /player/stats/today [get]
func getTodayStatsHandler(c *gin.Context, svc *playerservice.PlayerService) {
	userID := resp.GetUserID(c)
	stats, err := svc.GetTodayStats(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, stats)
}

// getOverviewStatsHandler 获取总览统计
// @Summary 获取总览统计
// @Tags 陪玩师-统计
// @Produce json
// @Security BearerAuth
// @Success 200 {object} playerservice.PlayerStatsOverview
// @Router /player/stats/overview [get]
func getOverviewStatsHandler(c *gin.Context, svc *playerservice.PlayerService) {
	userID := resp.GetUserID(c)
	stats, err := svc.GetOverviewStats(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, stats)
}
