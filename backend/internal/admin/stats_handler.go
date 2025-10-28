package admin

import (
    "net/http"
    "strconv"

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
    if err != nil { c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success:false, Code:500, Message: err.Error()}); return }
    c.JSON(http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "OK", Data: d })
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
    if err != nil { c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success:false, Code:500, Message: err.Error()}); return }
    c.JSON(http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "OK", Data: items })
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
    if err != nil { c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success:false, Code:500, Message: err.Error()}); return }
    c.JSON(http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "OK", Data: items })
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
    if err != nil { c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success:false, Code:500, Message: err.Error()}); return }
    c.JSON(http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "OK", Data: m })
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
    if err != nil { c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success:false, Code:500, Message: err.Error()}); return }
    c.JSON(http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "OK", Data: items })
}

func parseIntDefault(s string, def int) int {
    if s == "" { return def }
    if v, err := strconv.Atoi(s); err == nil && v > 0 { return v }
    return def
}

