package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/stats"
)

// StatsHandler 统计数据Handler
type StatsHandler struct {
	svc *stats.StatsService
}

// NewStatsHandler 创建统计Handler
func NewStatsHandler(svc *stats.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// Dashboard 获取仪表板数据
// @Summary      仪表板数据
// @Description  获取平台统计数据总览
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
// @Router       /admin/stats/dashboard [get]
func (h *StatsHandler) Dashboard(c *gin.Context) {
	dashboard, err := h.svc.Dashboard(c.Request.Context())
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    dashboard,
	})
}

// RevenueTrend 获取收入趋势
// @Summary      收入趋势
// @Description  获取指定天数的收入趋势
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        days           query     int     false  "天数" default(7)
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    trend,
	})
}

// UserGrowth 获取用户增长趋势
// @Summary      用户增长趋势
// @Description  获取指定天数的用户增长趋势
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        days           query     int     false  "天数" default(7)
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    trend,
	})
}

// OrdersSummary 获取订单状态汇总
// @Summary      订单状态汇总
// @Description  获取各状态订单数量统计
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
// @Router       /admin/stats/orders [get]
func (h *StatsHandler) OrdersSummary(c *gin.Context) {
	stats, err := h.svc.OrdersByStatus(c.Request.Context())
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    stats,
	})
}

// TopPlayers 获取顶级陪玩师
// @Summary      顶级陪玩师
// @Description  获取收入最高的陪玩师列表
// @Tags         Admin - Stats
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        limit          query     int     false  "数量限制" default(10)
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    players,
	})
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
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"entityStats": entityStats,
			"actionStats": actionStats,
		},
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
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      500            {object}  model.APIResponse[any]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    trend,
	})
}
