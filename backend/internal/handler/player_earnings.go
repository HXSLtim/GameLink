package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/earnings"
)

// RegisterPlayerEarningsRoutes 注册陪玩师端收益管理路由
func RegisterPlayerEarningsRoutes(router gin.IRouter, svc *earnings.EarningsService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/earnings")
	group.Use(authMiddleware) // 需要认证
	{
		group.GET("/summary", func(c *gin.Context) { getEarningsSummaryHandler(c, svc) })
		group.GET("/trend", func(c *gin.Context) { getEarningsTrendHandler(c, svc) })
		group.POST("/withdraw", func(c *gin.Context) { requestWithdrawHandler(c, svc) })
		group.GET("/withdraw-history", func(c *gin.Context) { getWithdrawHistoryHandler(c, svc) })
	}
}

// getEarningsSummaryHandler 获取收益概览
// @Summary      获取收益概览
// @Description  获取陪玩师收益概览
// @Tags         Player - Earnings
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[earnings.EarningsSummaryResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/earnings/summary [get]
func getEarningsSummaryHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	resp, err := svc.GetEarningsSummary(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[earnings.EarningsSummaryResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getEarningsTrendHandler 获取收益趋势
// @Summary      获取收益趋势
// @Description  获取陪玩师收益趋势
// @Tags         Player - Earnings
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        days           query     int     true  "天数（7/30/90）"
// @Success      200            {object}  model.APIResponse[earnings.EarningsTrendResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/earnings/trend [get]
func getEarningsTrendHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	days, err := strconv.Atoi(c.Query("days"))
	if err != nil || days < 7 || days > 90 {
		respondError(c, http.StatusBadRequest, ErrInvalidParameter)
		return
	}

	resp, err := svc.GetEarningsTrend(c.Request.Context(), userID, days)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[earnings.EarningsTrendResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// requestWithdrawHandler 申请提现
// @Summary      申请提现
// @Description  申请提现
// @Tags         Player - Earnings
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                     true  "Bearer {token}"
// @Param        request        body      earnings.WithdrawRequest   true  "提现信息"
// @Success      200            {object}  model.APIResponse[earnings.WithdrawResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/earnings/withdraw [post]
func requestWithdrawHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	var req earnings.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.RequestWithdraw(c.Request.Context(), userID, req)
	if err != nil {
		if err == earnings.ErrInsufficientBalance {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[earnings.WithdrawResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "提现申请已提交",
		Data:    *resp,
	})
}

// getWithdrawHistoryHandler 获取提现记录
// @Summary      获取提现记录
// @Description  获取提现记录
// @Tags         Player - Earnings
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[earnings.WithdrawHistoryResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/earnings/withdraw-history [get]
func getWithdrawHistoryHandler(c *gin.Context, svc *earnings.EarningsService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := svc.GetWithdrawHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[earnings.WithdrawHistoryResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}
