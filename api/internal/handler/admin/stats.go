package admin

import (
	"strconv"
	"time"

	_ "gamelink/internal/model" // Imported for Swagger annotations

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	statsrepo "gamelink/internal/repository/stats"
	statsservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// StatsHandler 统计数据Handler
type StatsHandler struct {
	svc *statsservice.StatsService
}

// NewStatsHandler 创建统计Handler
func NewStatsHandler(svc *statsservice.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// RegisterStatsAnalysisRoutes 注册统计分析和仪表板路由
func RegisterStatsAnalysisRoutes(router gin.IRouter, db *gorm.DB) {
	statsRepo := statsrepo.NewStatsRepository(db)
	h := NewStatsHandler(statsservice.NewStatsService(statsRepo))

	group := router
	// 仪表板概览
	group.GET("/stats/dashboard", h.Dashboard)
	// 收入趋势
	group.GET("/stats/revenue-trend", h.RevenueTrend)
	// 用户增长
	group.GET("/stats/user-growth", h.UserGrowth)
	// 订单统计
	group.GET("/stats/orders", h.OrdersSummary)
	// 顶级陪玩师
	group.GET("/stats/top-players", h.TopPlayers)
	// 审计概览
	group.GET("/stats/audit/overview", h.AuditOverview)
	// 审计趋势
	group.GET("/stats/audit/trend", h.AuditTrend)
}

// Dashboard 获取仪表板数据
// @Summary      仪表板数据
// @Description  获取平台统计数据总览
// @Tags         Admin - Stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/stats/dashboard [get]
func (h *StatsHandler) Dashboard(c *gin.Context) {
	dashboard, err := h.svc.Dashboard(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, dashboard)
}

// RevenueTrend 获取收入趋势
// @Summary      收入趋势
// @Description  获取指定天数的收入趋势
// @Tags         Admin - Stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        days  query     int  false  "天数" default(7)
// @Success      200   {object}  model.SuccessResponse
// @Failure      400   {object}  model.ErrorResponse
// @Failure      401   {object}  model.ErrorResponse
// @Failure      500   {object}  model.ErrorResponse
// @Router       /admin/stats/revenue-trend [get]
func (h *StatsHandler) RevenueTrend(c *gin.Context) {
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if val, err := strconv.Atoi(daysStr); err == nil && val > 0 {
			days = val
		}
	}

	trend, err := h.svc.RevenueTrend(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, trend)
}

// UserGrowth 获取用户增长趋势
// @Summary      用户增长趋势
// @Description  获取指定天数的用户增长趋势
// @Tags         Admin - Stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        days  query     int  false  "天数" default(7)
// @Success      200   {object}  model.SuccessResponse
// @Failure      400   {object}  model.ErrorResponse
// @Failure      401   {object}  model.ErrorResponse
// @Failure      500   {object}  model.ErrorResponse
// @Router       /admin/stats/user-growth [get]
func (h *StatsHandler) UserGrowth(c *gin.Context) {
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if val, err := strconv.Atoi(daysStr); err == nil && val > 0 {
			days = val
		}
	}

	trend, err := h.svc.UserGrowth(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, trend)
}

// OrdersSummary 获取订单状态汇总
// @Summary      订单状态汇总
// @Description  获取各状态订单数量统计
// @Tags         Admin - Stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/stats/orders [get]
func (h *StatsHandler) OrdersSummary(c *gin.Context) {
	stats, err := h.svc.OrdersByStatus(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, stats)
}

// TopPlayers 获取顶级陪玩师
// @Summary      顶级陪玩师
// @Description  获取收入最高的陪玩师列表
// @Tags         Admin - Stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "数量限制" default(10)
// @Success      200    {object}  model.SuccessResponse
// @Failure      400    {object}  model.ErrorResponse
// @Failure      401    {object}  model.ErrorResponse
// @Failure      500    {object}  model.ErrorResponse
// @Router       /admin/stats/top-players [get]
func (h *StatsHandler) TopPlayers(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	players, err := h.svc.TopPlayers(c.Request.Context(), limit)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, players)
}

// AuditOverview 获取审计概览
// @Summary      审计概览
// @Description  获取审计日志统计概览
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        from           query     string  false  "开始日期"
// @Param        to             query     string  false  "结束日期"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/stats/audit/overview [get]
func (h *StatsHandler) AuditOverview(c *gin.Context) {
	var from, to *time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	entityStats, actionStats, err := h.svc.AuditOverview(c.Request.Context(), from, to)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, gin.H{
		"entityStats": entityStats,
		"actionStats": actionStats,
	})
}

// AuditTrend 获取审计趋势
// @Summary      审计趋势
// @Description  获取审计日志时间趋势
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        from           query     string  false  "开始日期"
// @Param        to             query     string  false  "结束日期"
// @Param        entity         query     string  false  "实体类型"
// @Param        action         query     string  false  "操作类型"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Failure      500            {object}  model.ErrorResponse
// @Router       /admin/stats/audit/trend [get]
func (h *StatsHandler) AuditTrend(c *gin.Context) {
	var from, to *time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	entity := c.Query("entity")
	action := c.Query("action")

	trend, err := h.svc.AuditTrend(c.Request.Context(), from, to, entity, action)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, trend)
}
