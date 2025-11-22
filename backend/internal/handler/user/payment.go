package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/service/payment"
)

// CreatePaymentResponse 创建支付响应（类型别名）
type CreatePaymentResponse = payment.CreatePaymentResponse

// PaymentStatusResponse 支付状态响应（类型别名）
type PaymentStatusResponse = payment.PaymentStatusResponse

// RegisterPaymentRoutes 注册用户端支付路由
func RegisterPaymentRoutes(router gin.IRouter, svc *payment.PaymentService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/payments")
	group.Use(authMiddleware) // 需要认证
	group.POST("", func(c *gin.Context) { createPaymentHandler(c, svc) })
	group.GET("/:id", func(c *gin.Context) { getPaymentStatusHandler(c, svc) })
	group.POST("/:id/cancel", func(c *gin.Context) { cancelPaymentHandler(c, svc) })
}

// createPaymentHandler 创建支付
// @Summary      创建支付
// @Description  为订单创建支付
// @Tags         User - Payments
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      payment.CreatePaymentRequest    true  "创建支付请求"
// @Success      200            {object}  model.APIResponse[model.Payment]
// @Failure      400            {object}  apierr.APIError
// @Failure      401            {object}  apierr.APIError
// @Failure      404            {object}  apierr.APIError
// @Failure      409            {object}  apierr.APIError
// @Failure      500            {object}  apierr.APIError
// @Router       /user/payments [post]
func createPaymentHandler(c *gin.Context, svc *payment.PaymentService) {
	userID := getUserIDFromContext(c)

	var req payment.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.CreatePayment(c.Request.Context(), userID, req)
	if err != nil {
		if err == payment.ErrOrderAlreadyPaid {
			respondAPIError(c, apierr.Conflict(err.Error()))
			return
		}
		if err == payment.ErrValidation {
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
// @Description  获取支付状态
// @Tags         User - Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      uint64  true  "支付ID"
// @Success      200   {object}  model.APIResponse[model.Payment]
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Failure      500   {object}  apierr.APIError
// @Router       /user/payments/{id} [get]
func getPaymentStatusHandler(c *gin.Context, svc *payment.PaymentService) {
	paymentID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	resp, err := svc.GetPaymentStatus(c.Request.Context(), paymentID)
	if err != nil {
		if err == payment.ErrNotFound {
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
// @Description  取消支付
// @Tags         User - Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      uint64  true  "支付ID"
// @Success      200   {object}  model.APIResponse[any]
// @Failure      400   {object}  apierr.APIError
// @Failure      401   {object}  apierr.APIError
// @Failure      403   {object}  apierr.APIError
// @Failure      404   {object}  apierr.APIError
// @Failure      500   {object}  apierr.APIError
// @Router       /user/payments/{id}/cancel [post]
func cancelPaymentHandler(c *gin.Context, svc *payment.PaymentService) {
	userID := getUserIDFromContext(c)

	paymentID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	if err := svc.CancelPayment(c.Request.Context(), userID, paymentID); err != nil {
		if err == payment.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("取消支付失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "支付已取消", struct{}{})
}
