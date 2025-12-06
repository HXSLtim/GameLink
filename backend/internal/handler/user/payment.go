package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	paymentservice "gamelink/internal/service/payment"
	"gamelink/pkg/apierr"
)

type CreatePaymentResponse = paymentservice.CreatePaymentResponse
type Payment = model.Payment
type PaymentStatusResponse = paymentservice.PaymentStatusResponse

func RegisterPaymentRoutes(router gin.IRouter, svc *paymentservice.PaymentService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/payments")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { listPaymentsHandler(c, svc) })
	group.POST("", func(c *gin.Context) { createPaymentHandler(c, svc) })
	group.GET("/:id", func(c *gin.Context) { getPaymentStatusHandler(c, svc) })
	group.POST("/:id/cancel", func(c *gin.Context) { cancelPaymentHandler(c, svc) })
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID       uint64 `json:"orderId" binding:"required"`
	PaymentMethod string `json:"paymentMethod" binding:"required"`
}

// createPaymentHandler 创建支付
// @Summary      创建支付
// @Description  为订单创建支付
// @Tags         User - Payments
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      CreatePaymentRequest    true  "创建支付请求"
// @Success      200            {object}  model.SuccessResponse
// @Router       /user/payments [post]
func createPaymentHandler(c *gin.Context, svc *paymentservice.PaymentService) {
	userID := getUserIDFromContext(c)

	var req paymentservice.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.CreatePayment(c.Request.Context(), userID, req)
	if err != nil {
		if err == paymentservice.ErrOrderAlreadyPaid {
			respondAPIError(c, apierr.Conflict(err.Error()))
			return
		}
		if err == paymentservice.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("创建支付失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "支付创建成功", *resp)
}

// getPaymentStatusHandler 获取支付状态
// @Summary      获取支付状态
// @Tags         User - Payments
// @Security     BearerAuth
// @Param        id    path      uint64  true  "支付ID"
// @Success      200   {object}  model.APIResponse[Payment]
// @Router       /user/payments/{id} [get]
func getPaymentStatusHandler(c *gin.Context, svc *paymentservice.PaymentService) {
	paymentID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	resp, err := svc.GetPaymentStatus(c.Request.Context(), paymentID)
	if err != nil {
		if err == paymentservice.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("获取支付状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// cancelPaymentHandler 取消支付
// @Summary      取消支付
// @Tags         User - Payments
// @Security     BearerAuth
// @Param        id    path      uint64  true  "支付ID"
// @Success      200   {object}  model.SuccessResponse
// @Router       /user/payments/{id}/cancel [post]
func cancelPaymentHandler(c *gin.Context, svc *paymentservice.PaymentService) {
	userID := getUserIDFromContext(c)

	paymentID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	if err := svc.CancelPayment(c.Request.Context(), userID, paymentID); err != nil {
		if err == paymentservice.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("取消支付失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "支付已取消", struct{}{})
}

// listPaymentsHandler 获取当前用户的支付列表
// @Summary      获取支付列表
// @Tags         User - Payments
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        status     query  string  false  "状态过滤"
// @Param        method     query  string  false  "支付方式"
// @Param        dateFrom   query  string  false  "开始时间 RFC3339"
// @Param        dateTo     query  string  false  "结束时间 RFC3339"
// @Success      200        {object}  model.SuccessResponse
// @Router       /user/payments [get]
func listPaymentsHandler(c *gin.Context, svc *paymentservice.PaymentService) {
	userID := getUserIDFromContext(c)
	opts := buildPaymentListOptionsFromQuery(c)
	opts.UserID = &userID
	payments, total, err := svc.List(c.Request.Context(), opts)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取支付列表失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, "OK", struct {
		Items []model.Payment `json:"items"`
		Total int64           `json:"total"`
	}{Items: payments, Total: total})
}
func buildPaymentListOptionsFromQuery(c *gin.Context) repository.PaymentListOptions {
	opts := repository.PaymentListOptions{
		Page:     parseIntWithDefault(c.Query("page"), 1),
		PageSize: parseIntWithDefault(c.Query("pageSize"), 20),
	}
	if status := c.Query("status"); status != "" {
		s := model.PaymentStatus(status)
		opts.Status = &s
	}
	if method := c.Query("method"); method != "" {
		m := model.PaymentMethod(method)
		opts.Method = &m
	}
	if df := c.Query("dateFrom"); df != "" {
		if t, err := time.Parse(time.RFC3339, df); err == nil {
			opts.DateFrom = &t
		}
	}
	if dt := c.Query("dateTo"); dt != "" {
		if t, err := time.Parse(time.RFC3339, dt); err == nil {
			opts.DateTo = &t
		}
	}
	return opts
}

func parseIntWithDefault(val string, def int) int {
	if val == "" {
		return def
	}
	if v, err := strconv.Atoi(val); err == nil && v > 0 {
		return v
	}
	return def
}
