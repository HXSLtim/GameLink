package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/logging"
)

// --- Payment management ---

// UpdatePaymentInput 调整支付状态。
type UpdatePaymentInput struct {
	Status              model.PaymentStatus
	ProviderTradeNo     string
	ProviderRaw         json.RawMessage
	PaidAt              *time.Time
	RefundedAt          *time.Time
	RefundedAmountCents *int64 // 已退款金额（分）
}

// CreatePaymentInput 创建支付记录。
type CreatePaymentInput struct {
	OrderID     uint64
	UserID      uint64
	Method      model.PaymentMethod
	AmountCents int64
	Currency    model.Currency
	ProviderRaw json.RawMessage
}

// CreatePayment 新建支付记录，默认状态 pending。
func (s *AdminService) CreatePayment(ctx context.Context, in CreatePaymentInput) (*model.Payment, error) {
	if in.OrderID == 0 || in.UserID == 0 || in.AmountCents <= 0 || !model.IsValidCurrency(in.Currency) {
		return nil, ErrValidation
	}
	if in.Method == "" {
		return nil, ErrValidation
	}
	if _, err := s.orders.Get(ctx, in.OrderID); err != nil {
		return nil, WrapError(err, "get order")
	}
	if _, err := s.users.Get(ctx, in.UserID); err != nil {
		return nil, mapUserError(err)
	}
	pay := &model.Payment{
		OrderID:     in.OrderID,
		UserID:      in.UserID,
		Method:      in.Method,
		AmountCents: in.AmountCents,
		Currency:    in.Currency,
		Status:      model.PaymentStatusPending,
		ProviderRaw: in.ProviderRaw,
	}
	if err := s.payments.Create(ctx, pay); err != nil {
		return nil, WrapError(err, "create payment")
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), pay.ID, string(model.OpActionCreate), map[string]any{"status": pay.Status, "method": pay.Method})
	return pay, nil
}

// CapturePaymentInput 确认入账。
type CapturePaymentInput struct {
	ProviderTradeNo string
	ProviderRaw     json.RawMessage
	PaidAt          *time.Time
}

// CapturePayment 将支付置为 paid。
func (s *AdminService) CapturePayment(ctx context.Context, id uint64, in CapturePaymentInput) (*model.Payment, error) {
	pay, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}
	if !isAllowedPaymentTransition(pay.Status, model.PaymentStatusPaid) {
		return nil, ErrValidation
	}
	pay.Status = model.PaymentStatusPaid
	pay.ProviderTradeNo = strings.TrimSpace(in.ProviderTradeNo)
	pay.ProviderRaw = in.ProviderRaw
	if in.PaidAt != nil {
		pay.PaidAt = in.PaidAt
	} else {
		now := time.Now().UTC()
		pay.PaidAt = &now
	}
	if err := s.payments.Update(ctx, pay); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), pay.ID, string(model.OpActionCapture), map[string]any{"trade_no": pay.ProviderTradeNo})
	return pay, nil
}

// ListPayments 列出支付记录。
func (s *AdminService) ListPayments(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, *model.Pagination, error) {
	normalized := opts
	normalized.Page = repository.NormalizePage(opts.Page)
	normalized.PageSize = repository.NormalizePageSize(opts.PageSize)

	payments, total, err := s.payments.List(ctx, normalized)
	if err != nil {
		return nil, nil, err
	}

	pagination := buildPagination(normalized.Page, normalized.PageSize, total)
	return payments, &pagination, nil
}

// GetPayment 获取支付详情。
func (s *AdminService) GetPayment(ctx context.Context, id uint64) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}
	return payment, nil
}

// GetPaymentWithRelations 获取支付详情及关联的订单和用户信息。
func (s *AdminService) GetPaymentWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	payment, err := s.payments.GetWithRelations(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment with relations")
	}
	return payment, nil
}

// GetPaymentsByOrderID 根据订单ID获取所有支付记录。
func (s *AdminService) GetPaymentsByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	payments, err := s.payments.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, WrapError(err, "get payments by order id")
	}
	return payments, nil
}

// UpdatePayment 更新支付状态。
func (s *AdminService) UpdatePayment(ctx context.Context, id uint64, input UpdatePaymentInput) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}

	if !isValidPaymentStatus(input.Status) {
		return nil, apierr.BadRequest("invalid payment status")
	}

	if !isAllowedPaymentTransition(payment.Status, input.Status) {
		return nil, apierr.BadRequest("invalid payment status transition")
	}

	payment.Status = input.Status
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.PaidAt = input.PaidAt
	payment.RefundedAt = input.RefundedAt
	if input.RefundedAmountCents != nil {
		payment.RefundedAmountCents = *input.RefundedAmountCents
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, WrapError(err, "update payment")
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	payAction := model.OpActionUpdateStatus
	if input.Status == model.PaymentStatusRefunded {
		payAction = model.OpActionRefund
	}
	s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(payAction), map[string]any{"status": payment.Status})
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("payment_status_changed", slog.Uint64("payment_id", payment.ID), slog.String("status", string(payment.Status)), slog.String("request_id", rid))
	} else {
		slog.Info("payment_status_changed", slog.Uint64("payment_id", payment.ID), slog.String("status", string(payment.Status)))
	}
	return payment, nil
}

// UpdatePaymentWithRefund processes a refund with amount validation and logging.
// Requirements: 2.1, 2.2, 2.3, 2.4, 9.1, 9.2, 9.3
func (s *AdminService) UpdatePaymentWithRefund(ctx context.Context, id uint64, input UpdatePaymentInput, refundAmount int64, reason string, operatorID *uint64) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}

	// Validate refund amount
	if err := payment.ValidateRefundAmount(refundAmount); err != nil {
		return nil, apierr.BadRequest(err.Error())
	}

	// Validate status transition if status is changing
	if input.Status != "" && input.Status != payment.Status {
		if !isAllowedPaymentTransition(payment.Status, input.Status) {
			return nil, apierr.BadRequest("invalid payment status transition")
		}
		payment.Status = input.Status
	}

	// Update refunded amount
	payment.RefundedAmountCents += refundAmount
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.RefundedAt = input.RefundedAt

	// Check if fully refunded and update status
	if payment.IsFullyRefunded() && payment.Status != model.PaymentStatusRefunded {
		payment.Status = model.PaymentStatusRefunded
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, WrapError(err, "update payment")
	}

	s.invalidateCache(ctx, cacheKeyPayments)

	// 同步订单退款聚合字段（累计退款金额、退款状态）
	if err := s.syncOrderRefundSummary(ctx, payment.OrderID, reason, input.RefundedAt); err != nil {
		return nil, WrapError(err, "sync order refund summary")
	}

	// Log the refund operation with detailed metadata
	s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionRefund), map[string]any{
		"refund_amount_cents":  refundAmount,
		"total_refunded_cents": payment.RefundedAmountCents,
		"remaining_cents":      payment.RemainingRefundableAmount(),
		"reason":               reason,
		"is_full_refund":       payment.IsFullyRefunded(),
		"operator_id":          operatorID,
		"status":               payment.Status,
	})

	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("payment_refunded",
			slog.Uint64("payment_id", payment.ID),
			slog.Int64("refund_amount", refundAmount),
			slog.Int64("total_refunded", payment.RefundedAmountCents),
			slog.String("status", string(payment.Status)),
			slog.String("request_id", rid))
	} else {
		slog.Info("payment_refunded",
			slog.Uint64("payment_id", payment.ID),
			slog.Int64("refund_amount", refundAmount),
			slog.Int64("total_refunded", payment.RefundedAmountCents),
			slog.String("status", string(payment.Status)))
	}

	return payment, nil
}

// DeletePayment 删除支付记录。
func (s *AdminService) DeletePayment(ctx context.Context, id uint64) error {
	if err := s.payments.Delete(ctx, id); err != nil {
		return WrapError(err, "delete payment")
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), id, string(model.OpActionDelete), nil)
	return nil
}

// ============================================================================
// Batch Payment Operations
// ============================================================================

// BatchCaptureResult 批量收款操作结果
type BatchCaptureResult struct {
	SuccessCount int                   `json:"successCount"`
	FailedCount  int                   `json:"failedCount"`
	FailedIDs    []uint64              `json:"failedIds,omitempty"`
	Errors       []BatchOperationError `json:"errors,omitempty"`
}

// BatchOperationError 批量操作错误详情
type BatchOperationError struct {
	PaymentID uint64 `json:"paymentId"`
	Message   string `json:"message"`
}

// BatchCaptureRequest 批量收款请求
type BatchCaptureRequest struct {
	PaymentIDs      []uint64   `json:"paymentIds" binding:"required,min=1,max=500"`
	ProviderTradeNo string     `json:"providerTradeNo,omitempty"`
	PaidAt          *time.Time `json:"paidAt,omitempty"`
}

// BatchCapture 批量收款 - 将多个pending状态的支付标记为已支付
// 业务规则：
// 1. 只能处理pending状态的支付
// 2. 支付会设置为paid状态
// 3. 返回成功/失败计数及错误详情
func (s *AdminService) BatchCapture(ctx context.Context, req BatchCaptureRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	paidAt := req.PaidAt
	if paidAt == nil {
		now := time.Now().UTC()
		paidAt = &now
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能capture pending状态的支付
		if payment.Status != model.PaymentStatusPending {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for capture: %s (expected: pending)", payment.Status),
			})
			continue
		}

		// 更新支付状态
		payment.Status = model.PaymentStatusPaid
		payment.PaidAt = paidAt
		if req.ProviderTradeNo != "" {
			payment.ProviderTradeNo = strings.TrimSpace(req.ProviderTradeNo)
		} else {
			payment.ProviderTradeNo = fmt.Sprintf("batch_capture_%d_%d", paymentID, time.Now().Unix())
		}

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionCapture), map[string]any{
			"batch_operation": true,
			"trade_no":        payment.ProviderTradeNo,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchRefundRequest 批量退款请求
type BatchRefundRequest struct {
	PaymentIDs []uint64   `json:"paymentIds" binding:"required,min=1,max=500"`
	Reason     string     `json:"reason" binding:"required,max=500"`
	RefundedAt *time.Time `json:"refundedAt,omitempty"`
}

// BatchRefund 批量退款 - 退款多个已支付的支付
// 业务规则：
// 1. 只能退款paid状态的支付
// 2. 全额退款，状态更新为refunded
// 3. 订单状态也会更新为refunded
func (s *AdminService) BatchRefund(ctx context.Context, req BatchRefundRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, apierr.BadRequest("refund reason is required")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	refundedAt := req.RefundedAt
	if refundedAt == nil {
		now := time.Now().UTC()
		refundedAt = &now
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能退款paid状态的支付
		if payment.Status != model.PaymentStatusPaid {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for refund: %s (expected: paid)", payment.Status),
			})
			continue
		}

		// 检查是否已经全额退款
		if payment.IsFullyRefunded() {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   "payment is already fully refunded",
			})
			continue
		}

		// 更新支付状态
		payment.Status = model.PaymentStatusRefunded
		payment.RefundedAt = refundedAt
		payment.RefundedAmountCents = payment.AmountCents
		payment.ProviderTradeNo = fmt.Sprintf("batch_refund_%d_%d", paymentID, time.Now().Unix())

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		// 更新关联订单状态
		order, err := s.orders.Get(ctx, payment.OrderID)
		if err == nil {
			order.Status = model.OrderStatusRefunded
			order.RefundAmountCents = payment.AmountCents
			order.RefundReason = req.Reason
			order.RefundedAt = refundedAt
			_ = s.orders.Update(ctx, order)
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionRefund), map[string]any{
			"batch_operation":     true,
			"refund_amount_cents": payment.AmountCents,
			"reason":              req.Reason,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchCancelRequest 批量取消支付请求
type BatchCancelRequest struct {
	PaymentIDs []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
}

// BatchCancel 批量取消支付 - 取消多个pending状态的支付
// 业务规则：
// 1. 只能取消pending状态的支付
// 2. 支付状态更新为failed
func (s *AdminService) BatchCancel(ctx context.Context, req BatchCancelRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能取消pending状态的支付
		if payment.Status != model.PaymentStatusPending {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for cancel: %s (expected: pending)", payment.Status),
			})
			continue
		}

		// 更新支付状态为failed（表示已取消）
		payment.Status = model.PaymentStatusFailed

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionCancel), map[string]any{
			"batch_operation": true,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchUpdateStatusRequest 批量更新支付状态请求
type BatchUpdateStatusRequest struct {
	PaymentIDs []uint64            `json:"paymentIds" binding:"required,min=1,max=500"`
	Status     model.PaymentStatus `json:"-" binding:"-"` // Not used for binding, set from handler
}

// BatchUpdateStatus 批量更新支付状态
// 业务规则：
// 1. 验证状态转换是否有效
// 2. 只允许有效的状态转换
func (s *AdminService) BatchUpdateStatus(ctx context.Context, req BatchUpdateStatusRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	// 验证状态是否有效
	if !isValidPaymentStatus(req.Status) {
		return nil, apierr.BadRequest("invalid payment status")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态转换
		if !isAllowedPaymentTransition(payment.Status, req.Status) {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status transition from %s to %s", payment.Status, req.Status),
			})
			continue
		}

		// 如果是转换为paid，设置PaidAt
		if req.Status == model.PaymentStatusPaid && payment.PaidAt == nil {
			now := time.Now().UTC()
			payment.PaidAt = &now
		}

		// 如果是转换为refunded，设置RefundedAt
		if req.Status == model.PaymentStatusRefunded && payment.RefundedAt == nil {
			now := time.Now().UTC()
			payment.RefundedAt = &now
			payment.RefundedAmountCents = payment.AmountCents
		}

		// 更新支付状态
		payment.Status = req.Status

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		action := model.OpActionUpdateStatus
		if req.Status == model.PaymentStatusPaid {
			action = model.OpActionCapture
		} else if req.Status == model.PaymentStatusRefunded {
			action = model.OpActionRefund
		} else if req.Status == model.PaymentStatusFailed {
			action = model.OpActionCancel
		}
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(action), map[string]any{
			"batch_operation": true,
			"old_status":      string(payment.Status),
			"new_status":      string(req.Status),
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// GetPaymentLogs returns operation logs for a payment.
// Requirements: 2.5
func (s *AdminService) GetPaymentLogs(ctx context.Context, paymentID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	if s.tx == nil {
		return nil, 0, apierr.InternalError("transaction manager not configured")
	}
	logs, total, err := s.repos().OpLogs.ListByEntity(ctx, string(model.OpEntityPayment), paymentID, opts)
	if err != nil {
		return nil, 0, WrapError(err, "get payment logs")
	}
	return logs, total, nil
}

func (s *AdminService) syncOrderRefundSummary(ctx context.Context, orderID uint64, reason string, refundedAt *time.Time) error {
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return WrapError(err, "get order")
	}

	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return WrapError(err, "list payments by order")
	}

	var totalRefundedCents int64
	for _, pay := range payments {
		totalRefundedCents += clampRefundedAmount(pay.RefundedAmountCents, pay.AmountCents)
	}
	if totalRefundedCents > order.TotalPriceCents {
		totalRefundedCents = order.TotalPriceCents
	}

	order.RefundAmountCents = totalRefundedCents
	if totalRefundedCents > 0 {
		trimmedReason := strings.TrimSpace(reason)
		if trimmedReason != "" {
			order.RefundReason = trimmedReason
		}
		if refundedAt != nil {
			order.RefundedAt = refundedAt
		} else if order.RefundedAt == nil {
			now := time.Now().UTC()
			order.RefundedAt = &now
		}
		if order.TotalPriceCents > 0 && totalRefundedCents >= order.TotalPriceCents {
			order.Status = model.OrderStatusRefunded
		}
	}

	if err := s.orders.Update(ctx, order); err != nil {
		return WrapError(err, "update order")
	}

	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionRefund), map[string]any{
		"refund_amount_cents": order.RefundAmountCents,
		"refund_reason":       order.RefundReason,
		"status":              order.Status,
	})
	return nil
}


func isValidPaymentStatus(status model.PaymentStatus) bool {
	switch status {
	case model.PaymentStatusPending, model.PaymentStatusPaid, model.PaymentStatusFailed, model.PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func isAllowedPaymentTransition(prev, next model.PaymentStatus) bool {
	if prev == next {
		return true
	}
	switch prev {
	case model.PaymentStatusPending:
		return next == model.PaymentStatusPaid || next == model.PaymentStatusFailed || next == model.PaymentStatusRefunded
	case model.PaymentStatusPaid:
		return next == model.PaymentStatusRefunded
	case model.PaymentStatusFailed, model.PaymentStatusRefunded:
		return false
	default:
		return false
	}
}

