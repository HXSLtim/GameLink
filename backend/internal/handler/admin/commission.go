package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/commission"
)

// RegisterCommissionRoutes Register admin commission management routes
func RegisterCommissionRoutes(router gin.IRouter, svc *commission.CommissionService, scheduler interface{ TriggerSettlement(string) error }) {
	group := router.Group("/admin/commission")
	{
		// 抽成规则管理
		group.POST("/rules", func(c *gin.Context) { createCommissionRuleHandler(c, svc) })
		group.PUT("/rules/:id", func(c *gin.Context) { updateCommissionRuleHandler(c, svc) })

		// 月度结算
		group.POST("/settlements/trigger", func(c *gin.Context) { triggerSettlementHandler(c, scheduler) })
		group.GET("/stats", func(c *gin.Context) { getPlatformStatsHandler(c, svc) })
	}
}

// createCommissionRuleHandler 创建抽成规则
// @Summary      创建抽成规则
// @Description  管理员创建抽成规�?
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                    true  "Bearer {token}"
// @Param        request        body      commission.CreateCommissionRuleRequest  true  "抽成规则信息"
// @Success      200            {object}  model.APIResponse[model.CommissionRule]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/commission/rules [post]
func createCommissionRuleHandler(c *gin.Context, svc *commission.CommissionService) {
	var req commission.CreateCommissionRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	rule, err := svc.CreateCommissionRule(c.Request.Context(), req)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[model.CommissionRule]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Commission rule created successfully",
		Data:    *rule,
	})
}

// updateCommissionRuleHandler 更新抽成规则
// @Summary      更新抽成规则
// @Description  管理员更新抽成规�?
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                    true  "Bearer {token}"
// @Param        id             path      int                                       true  "规则ID"
// @Param        request        body      commission.UpdateCommissionRuleRequest  true  "更新信息"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/commission/rules/{id} [put]
func updateCommissionRuleHandler(c *gin.Context, svc *commission.CommissionService) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	var req commission.UpdateCommissionRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = svc.UpdateCommissionRule(c.Request.Context(), id, req)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Commission rule updated successfully",
	})
}

// triggerSettlementHandler 手动触发月度结算
// @Summary      手动触发月度结算
// @Description  管理员手动触发指定月份的结算
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/commission/settlements/trigger [post]
func triggerSettlementHandler(c *gin.Context, scheduler interface{ TriggerSettlement(string) error }) {
	month := c.Query("month")
	if month == "" {
		// 默认结算上个�?
		lastMonth := time.Now().AddDate(0, -1, 0)
		month = lastMonth.Format("2006-01")
	}

	err := scheduler.TriggerSettlement(month)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Settlement triggered successfully for month: " + month,
	})
}

// getPlatformStatsHandler 获取平台统计
// @Summary      获取平台统计
// @Description  管理员查看平台月度统计数�?
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  model.APIResponse[commission.PlatformStatsResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/commission/stats [get]
func getPlatformStatsHandler(c *gin.Context, svc *commission.CommissionService) {
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	stats, err := svc.GetPlatformStats(c.Request.Context(), month)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[commission.PlatformStatsResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *stats,
	})
}
