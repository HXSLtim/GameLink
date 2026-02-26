package admin

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// BatchCapturePaymentsRequest 批量收款请求
type BatchCapturePaymentsRequest struct {
	PaymentIDs      []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
	ProviderTradeNo string   `json:"providerTradeNo,omitempty"`
	PaidAt          *string  `json:"paidAt,omitempty"`
}

// BatchRefundPaymentsRequest 批量退款请求
type BatchRefundPaymentsRequest struct {
	PaymentIDs []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
	Reason     string   `json:"reason" binding:"required,max=500"`
	RefundedAt *string  `json:"refundedAt,omitempty"`
}

// BatchCancelPaymentsRequest 批量取消支付请求
type BatchCancelPaymentsRequest struct {
	PaymentIDs []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
}

// BatchUpdatePaymentsStatusRequest 批量更新支付状态请求
type BatchUpdatePaymentsStatusRequest struct {
	PaymentIDs []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
	Status     string   `json:"status" binding:"required"`
}

// BatchCapturePayments 批量确认支付入账
func (h *PaymentHandler) BatchCapturePayments(c *gin.Context) {
	var req BatchCapturePaymentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	paidAt, err := parseRFC3339Ptr(req.PaidAt)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid paid at time"))
		return
	}

	result, err := h.svc.BatchCapture(contextWithActor(c), adminservice.BatchCaptureRequest{
		PaymentIDs:      req.PaymentIDs,
		ProviderTradeNo: req.ProviderTradeNo,
		PaidAt:          paidAt,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch capture payments failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchRefundPayments 批量退款
func (h *PaymentHandler) BatchRefundPayments(c *gin.Context) {
	var req BatchRefundPaymentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	refundedAt, err := parseRFC3339Ptr(req.RefundedAt)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid refunded at time"))
		return
	}

	result, err := h.svc.BatchRefund(contextWithActor(c), adminservice.BatchRefundRequest{
		PaymentIDs: req.PaymentIDs,
		Reason:     req.Reason,
		RefundedAt: refundedAt,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch refund payments failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchCancelPayments 批量取消支付
func (h *PaymentHandler) BatchCancelPayments(c *gin.Context) {
	var req BatchCancelPaymentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	result, err := h.svc.BatchCancel(contextWithActor(c), adminservice.BatchCancelRequest{
		PaymentIDs: req.PaymentIDs,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch cancel payments failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchUpdatePaymentsStatus 批量更新支付状态
func (h *PaymentHandler) BatchUpdatePaymentsStatus(c *gin.Context) {
	var req BatchUpdatePaymentsStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	status := model.PaymentStatus(strings.ToLower(strings.TrimSpace(req.Status)))
	result, err := h.svc.BatchUpdateStatus(contextWithActor(c), adminservice.BatchUpdateStatusRequest{
		PaymentIDs: req.PaymentIDs,
		Status:     status,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch update payment status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}
