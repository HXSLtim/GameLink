package player

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/service/earnings"
)

// RegisterEarningsRoutes 注册陪玩师端收益管理路由
func RegisterEarningsRoutes(router gin.IRouter, svc *earnings.EarningsService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/earnings")
	group.Use(authMiddleware) // 需要认证
	group.GET("/summary", func(c *gin.Context) { getEarningsSummaryHandler(c, svc) })
	group.GET("/trend", func(c *gin.Context) { getEarningsTrendHandler(c, svc) })
	group.POST("/withdraw", func(c *gin.Context) { requestWithdrawHandler(c, svc) })
	group.GET("/withdraw-history", func(c *gin.Context) { getWithdrawHistoryHandler(c, svc) })
}

// getEarningsSummaryHandler 获取收益概览
// @Summary      获取收益概览
// @Description  获取陪玩师收益概览
// @Tags         Player - Earnings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.APIResponse[earnings.EarningsSummaryResponse]
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/earnings/summary [get]
func getEarningsSummaryHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	resp, err := svc.GetEarningsSummary(c.Request.Context(), userID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取收益概览失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// getEarningsTrendHandler 获取收益趋势
// @Summary      获取收益趋势
// @Description  获取陪玩师收益趋势数据
// @Tags         Player - Earnings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        days  query  int  true  "天数范围(7-90)"
// @Success      200   {object}  model.APIResponse[earnings.EarningsTrendResponse]
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      500   {object}  apierr.APIError
// @Router       /player/earnings/trend [get]
func getEarningsTrendHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	days, err := strconv.Atoi(c.Query("days"))
	if err != nil || days < 7 || days > 90 {
		respondAPIError(c, apierr.BadRequest("days参数无效，必须在7-90之间"))
		return
	}

	resp, err := svc.GetEarningsTrend(c.Request.Context(), userID, days)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取收益趋势失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// requestWithdrawHandler 申请提现
// @Summary      申请提现
// @Description  申请提现
// @Tags         Player - Earnings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  earnings.WithdrawRequest  true  "提现信息"
// @Success      200      {object}  model.APIResponse[earnings.WithdrawResponse]
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      403      {object}  apierr.APIError
// @Failure      422      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/earnings/withdraw [post]
func requestWithdrawHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	var req earnings.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.RequestWithdraw(c.Request.Context(), userID, req)
	if err != nil {
		if err == earnings.ErrInsufficientBalance {
			respondAPIError(c, apierr.BadRequest("余额不足"))
			return
		}
		if err == earnings.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("申请提现失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "提现申请已提交", *resp)
}

// getWithdrawHistoryHandler 获取提现记录
// @Summary      获取提现记录
// @Description  获取提现记录
// @Tags         Player - Earnings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码" default(1)
// @Param        pageSize  query     int  false  "每页数量" default(20)
// @Success      200       {object}  model.APIResponse[earnings.WithdrawHistoryResponse]
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /player/earnings/withdraw-history [get]
func getWithdrawHistoryHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := svc.GetWithdrawHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取提现记录失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}
