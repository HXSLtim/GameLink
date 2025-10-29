package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
)

// StatsHandler 管理统计接口。
type StatsHandler struct{ svc *service.StatsService }

func NewStatsHandler(s *service.StatsService) *StatsHandler { return &StatsHandler{svc: s} }

// Dashboard
// @Summary      Dashboard 概览
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/dashboard [get]
func (h *StatsHandler) Dashboard(c *gin.Context) {
	d, err := h.svc.Dashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK", Data: d})
}

// RevenueTrend
// @Summary      收入趋势（日）
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Param        days  query  int  false  "天数，默认7"
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/revenue-trend [get]
func (h *StatsHandler) RevenueTrend(c *gin.Context) {
	days := parseIntDefault(c.Query("days"), 7)
	items, err := h.svc.RevenueTrend(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK", Data: items})
}

// UserGrowth
// @Summary      用户增长（日）
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Param        days  query  int  false  "天数，默认7"
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/user-growth [get]
func (h *StatsHandler) UserGrowth(c *gin.Context) {
	days := parseIntDefault(c.Query("days"), 7)
	items, err := h.svc.UserGrowth(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK", Data: items})
}

// OrdersSummary
// @Summary      订单状态统计
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/orders [get]
func (h *StatsHandler) OrdersSummary(c *gin.Context) {
	m, err := h.svc.OrdersByStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK", Data: m})
}

// TopPlayers
// @Summary      Top 陪玩师排行
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Param        limit  query  int  false  "数量，默认10"
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/top-players [get]
func (h *StatsHandler) TopPlayers(c *gin.Context) {
	limit := parseIntDefault(c.Query("limit"), 10)
	items, err := h.svc.TopPlayers(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK", Data: items})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return def
}

// AuditOverview
// @Summary      审计总览（按实体/动作汇总）
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Param        dateFrom   query     string    false  "开始时间"
// @Param        dateTo     query     string    false  "结束时间"
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/audit/overview [get]
func (h *StatsHandler) AuditOverview(c *gin.Context) {
	from, err := queryTimePtr(c, "date_from")
	if err != nil {
		c.JSON(400, model.APIResponse[any]{Success: false, Code: 400, Message: "invalid date_from"})
		return
	}
	to, err := queryTimePtr(c, "date_to")
	if err != nil {
		c.JSON(400, model.APIResponse[any]{Success: false, Code: 400, Message: "invalid date_to"})
		return
	}
	byEntity, byAction, err := h.svc.AuditOverview(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(500, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, model.APIResponse[map[string]any]{Success: true, Code: 200, Message: "OK", Data: map[string]any{"byEntity": byEntity, "byAction": byAction}})
}

// AuditTrend
// @Summary      审计趋势（日）
// @Tags         Admin/Stats
// @Security     BearerAuth
// @Produce      json
// @Param        dateFrom   query     string    false  "开始时间"
// @Param        dateTo     query     string    false  "结束时间"
// @Param        entity     query  string  false  "实体类型" Enums(order,payment,player,game,review,user)
// @Param        action     query  string  false  "动作"
// @Success      200  {object}  map[string]any
// @Router       /admin/stats/audit/trend [get]
func (h *StatsHandler) AuditTrend(c *gin.Context) {
	from, err := queryTimePtr(c, "date_from")
	if err != nil {
		c.JSON(400, model.APIResponse[any]{Success: false, Code: 400, Message: "invalid date_from"})
		return
	}
	to, err := queryTimePtr(c, "date_to")
	if err != nil {
		c.JSON(400, model.APIResponse[any]{Success: false, Code: 400, Message: "invalid date_to"})
		return
	}
	entity := strings.TrimSpace(c.Query("entity"))
	action := strings.TrimSpace(c.Query("action"))
	items, err := h.svc.AuditTrend(c.Request.Context(), from, to, entity, action)
	if err != nil {
		c.JSON(500, model.APIResponse[any]{Success: false, Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, model.APIResponse[any]{Success: true, Code: 200, Message: "OK", Data: items})
}
