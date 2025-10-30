package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/payment"
)

// RegisterUserPaymentRoutes 注册用户端支付路由
func RegisterUserPaymentRoutes(router gin.IRouter, svc *payment.PaymentService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/payments")
	group.Use(authMiddleware) // 需要认证
	{
		group.POST("", func(c *gin.Context) { createPaymentHandler(c, svc) })
		group.GET("/:id", func(c *gin.Context) { getPaymentStatusHandler(c, svc) })
		group.POST("/:id/cancel", func(c *gin.Context) { cancelPaymentHandler(c, svc) })
	}
}

// createPaymentHandler 创建支付
// @Summary      创建支付
// @Description  为订单创建支付
// @Tags         User - Payments
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      payment.CreatePaymentRequest    true  "创建支付请求"
// @Success      200            {object}  model.APIResponse[payment.CreatePaymentResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/payments [post]
func createPaymentHandler(c *gin.Context, svc *payment.PaymentService) {
	userID := getUserIDFromContext(c)

	var req payment.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.CreatePayment(c.Request.Context(), userID, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[payment.CreatePaymentResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "支付创建成功",
		Data:    *resp,
	})
}

// getPaymentStatusHandler 查询支付状态
// @Summary      查询支付状态
// @Description  查询支付状态
// @Tags         User - Payments
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "支付ID"
// @Success      200            {object}  model.APIResponse[payment.PaymentStatusResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      404            {object}  model.APIResponse[any]
// @Router       /user/payments/{id} [get]
func getPaymentStatusHandler(c *gin.Context, svc *payment.PaymentService) {
	idStr := c.Param("id")
	paymentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidID)
		return
	}

	resp, err := svc.GetPaymentStatus(c.Request.Context(), paymentID)
	if err != nil {
		if err == payment.ErrNotFound {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[payment.PaymentStatusResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// cancelPaymentHandler 取消支付
// @Summary      取消支付
// @Description  取消支付
// @Tags         User - Payments
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "支付ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/payments/{id}/cancel [post]
func cancelPaymentHandler(c *gin.Context, svc *payment.PaymentService) {
	userID := getUserIDFromContext(c)

	idStr := c.Param("id")
	paymentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidID)
		return
	}

	if err := svc.CancelPayment(c.Request.Context(), userID, paymentID); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "支付已取消",
	})
}
