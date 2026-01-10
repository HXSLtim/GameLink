package player

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	userservice "gamelink/internal/service/user"
	"gamelink/pkg/apierr"
)

// Swagger DTOs
type EarningsSummaryResponseSwagger struct {
	TodayEarnings    int64 `json:"todayEarnings" example:"10000" description:"today's earnings (cents)"`
	MonthEarnings    int64 `json:"monthEarnings" example:"150000" description:"current month earnings (cents)"`
	TotalEarnings    int64 `json:"totalEarnings" example:"500000" description:"total earnings (cents)"`
	AvailableBalance int64 `json:"availableBalance" example:"20000" description:"available balance (cents)"`
	PendingBalance   int64 `json:"pendingBalance" example:"5000" description:"pending settlement (cents)"`
	WithdrawTotal    int64 `json:"withdrawTotal" example:"300000" description:"total withdrawn (cents)"`
}

type DailyEarningSwagger struct {
	Date       string `json:"date" example:"2024-01-15" description:"YYYY-MM-DD"`
	Earnings   int64  `json:"earnings" example:"5000" description:"earnings of the day (cents)"`
	OrderCount int    `json:"orderCount" example:"3" description:"orders completed"`
}

type EarningsTrendResponseSwagger struct {
	Trend []DailyEarningSwagger `json:"trend" description:"daily earning trend"`
}

type WithdrawResponseSwagger struct {
	WithdrawID uint64 `json:"withdrawId" example:"12345" description:"withdrawal ID"`
	Status     string `json:"status" example:"pending" description:"pending/processing/completed/failed"`
}

type WithdrawRecordSwagger struct {
	ID          uint64  `json:"id" example:"12345" description:"record ID"`
	AmountCents int64   `json:"amountCents" example:"10000" description:"amount in cents"`
	Method      string  `json:"method" example:"alipay" description:"alipay/wechat/bank"`
	Status      string  `json:"status" example:"completed" description:"pending/processing/completed/failed"`
	CreatedAt   string  `json:"createdAt" example:"2024-01-15T10:30:00Z" description:"created at"`
	ProcessedAt *string `json:"processedAt,omitempty" example:"2024-01-16T14:20:00Z" description:"processed at"`
}

type WithdrawHistoryResponseSwagger struct {
	Records []WithdrawRecordSwagger `json:"records" description:"withdrawal records"`
	Total   int64                   `json:"total" example:"100" description:"total records"`
}

// Swagger envelopes (avoid generics in annotations)
type EarningsSummaryAPIResponseSwagger struct {
	Success    bool                           `json:"success"`
	Code       int                            `json:"code"`
	Message    string                         `json:"message"`
	Data       EarningsSummaryResponseSwagger `json:"data"`
	Pagination *model.Pagination              `json:"pagination,omitempty"`
	TraceID    string                         `json:"traceId,omitempty"`
}

type EarningsTrendAPIResponseSwagger struct {
	Success    bool                         `json:"success"`
	Code       int                          `json:"code"`
	Message    string                       `json:"message"`
	Data       EarningsTrendResponseSwagger `json:"data"`
	Pagination *model.Pagination            `json:"pagination,omitempty"`
	TraceID    string                       `json:"traceId,omitempty"`
}

type WithdrawAPIResponseSwagger struct {
	Success    bool                    `json:"success"`
	Code       int                     `json:"code"`
	Message    string                  `json:"message"`
	Data       WithdrawResponseSwagger `json:"data"`
	Pagination *model.Pagination       `json:"pagination,omitempty"`
	TraceID    string                  `json:"traceId,omitempty"`
}

type WithdrawHistoryAPIResponseSwagger struct {
	Success    bool                           `json:"success"`
	Code       int                            `json:"code"`
	Message    string                         `json:"message"`
	Data       WithdrawHistoryResponseSwagger `json:"data"`
	Pagination *model.Pagination              `json:"pagination,omitempty"`
	TraceID    string                         `json:"traceId,omitempty"`
}

// RegisterEarningsRoutes 注册陪玩师端收益管理路由
func RegisterEarningsRoutes(router gin.IRouter, svc *userservice.EarningsService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/earnings")
	group.Use(authMiddleware)
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
// @Success      200  {object}  EarningsSummaryAPIResponseSwagger
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/earnings/summary [get]
func getEarningsSummaryHandler(c *gin.Context, svc *userservice.EarningsService) {
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
// @Success      200   {object}  EarningsTrendAPIResponseSwagger
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      500   {object}  apierr.APIError
// @Router       /player/earnings/trend [get]
func getEarningsTrendHandler(c *gin.Context, svc *userservice.EarningsService) {
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
// @Param        request  body  userservice.WithdrawRequest  true  "提现信息"
// @Success      200      {object}  WithdrawAPIResponseSwagger
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      403      {object}  apierr.APIError
// @Failure      422      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/earnings/withdraw [post]
//
// ✅ 资金安全修复: 增强错误处理,支持新的验证错误类型
func requestWithdrawHandler(c *gin.Context, svc *userservice.EarningsService) {
	userID := getUserIDFromContext(c)

	var req userservice.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.RequestWithdraw(c.Request.Context(), userID, req)
	if err != nil {
		// ✅ 处理各种验证错误
		switch err {
		case userservice.ErrInsufficientBalance:
			respondAPIError(c, apierr.BadRequest("余额不足"))
			return
		case userservice.ErrDailyLimitExceeded:
			respondAPIError(c, apierr.BadRequest("超过每日提现限额"))
			return
		case userservice.ErrMonthlyLimitExceeded:
			respondAPIError(c, apierr.BadRequest("超过每月提现限额"))
			return
		case userservice.ErrPendingWithdrawExists:
			respondAPIError(c, apierr.BadRequest("存在待处理的提现申请,请等待处理完成后再试"))
			return
		case userservice.ErrValidation:
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}

		// 处理其他包含详细信息的错误
		errMsg := err.Error()
		if len(errMsg) > 0 && (err != userservice.ErrNotFound && err != userservice.ErrUnauthorized) {
			respondAPIError(c, apierr.BadRequest(errMsg))
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
// @Success      200       {object}  WithdrawHistoryAPIResponseSwagger
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /player/earnings/withdraw-history [get]
func getWithdrawHistoryHandler(c *gin.Context, svc *userservice.EarningsService) {
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
