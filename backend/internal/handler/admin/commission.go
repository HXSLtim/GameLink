package admin

import (
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/service/commission"
)

// PlatformStatsResponse 平台统计响应（类型别名）
type PlatformStatsResponse = commission.PlatformStatsResponse

// RegisterCommissionRoutes Register admin commission management routes
func RegisterCommissionRoutes(router gin.IRouter, svc *commission.CommissionService, scheduler interface{ TriggerSettlement(string) error }) {
	group := router.Group("/commission")
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
// @Description  创建一个新的抽成规则
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request        body      commission.CreateCommissionRuleRequest  true  "抽成规则信息"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/commission/rules [post]
func createCommissionRuleHandler(c *gin.Context, svc *commission.CommissionService) {
	var req commission.CreateCommissionRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	rule, err := svc.CreateCommissionRule(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, *rule)
}

// updateCommissionRuleHandler 更新抽成规则
// @Summary      更新抽成规则
// @Description  更新一个已存在的抽成规则
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path      int                                       true  "规则ID"
// @Param        request        body      commission.UpdateCommissionRuleRequest  true  "更新信息"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/commission/rules/{id} [put]
func updateCommissionRuleHandler(c *gin.Context, svc *commission.CommissionService) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req commission.UpdateCommissionRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	err := svc.UpdateCommissionRule(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Commission rule updated successfully")
}

// triggerSettlementHandler 手动触发月度结算
// @Summary      手动触发月度结算
// @Description  管理员手动触发指定月份的结算
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/commission/settlements/trigger [post]
func triggerSettlementHandler(c *gin.Context, scheduler interface{ TriggerSettlement(string) error }) {
	month := c.Query("month")
	if month == "" {
		// 默认结算上个月
		lastMonth := time.Now().AddDate(0, -1, 0)
		month = lastMonth.Format("2006-01")
	}

	err := scheduler.TriggerSettlement(month)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Settlement triggered successfully for month: "+month)
}

// getPlatformStatsHandler 获取平台统计
// @Summary      获取平台统计
// @Description  获取指定月份的平台统计数据，包括总收入、总抽成等
// @Tags         Admin - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  model.PlatformStatsResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/commission/stats [get]
func getPlatformStatsHandler(c *gin.Context, svc *commission.CommissionService) {
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	stats, err := svc.GetPlatformStats(c.Request.Context(), month)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, *stats)
}
