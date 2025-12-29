package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/service/external"
	"gamelink/pkg/apierr"
)

// RefundService handles refund operations for payments.
// Requirements: 2.1, 2.2, 2.3, 2.4
type RefundService struct {
	payments  repository.PaymentRepository
	refunds   repository.RefundRecordRepository
	orders    repoiface.OrderReadWriter
	wallets   repository.WalletRepository
	opLogs    repository.OperationLogRepository
	providers map[model.PaymentMethod]ProviderClient
}

// NewRefundService creates a new RefundService instance.
func NewRefundService(
	payments repository.PaymentRepository,
	refunds repository.RefundRecordRepository,
	orders repoiface.OrderReadWriter,
) *RefundService {
	return &RefundService{
		payments: payments,
		refunds:  refunds,
		orders:   orders,
		providers: map[model.PaymentMethod]ProviderClient{
			model.PaymentMethodWeChat: wechatProvider{},
			model.PaymentMethodAlipay: alipayProvider{},
		},
	}
}

// SetWalletRepository injects wallet repository for refund credit.
func (s *RefundService) SetWalletRepository(repo repository.WalletRepository) {
	s.wallets = repo
}

// SetOperationLogRepository injects operation log repository for audit logging.
func (s *RefundService) SetOperationLogRepository(repo repository.OperationLogRepository) {
	s.opLogs = repo
}

// SetExternalConfig configures external API credentials for payment providers
func (s *RefundService) SetExternalConfig(cfg *external.Config) {
	factory := NewProviderFactory(cfg)
	s.providers = factory.CreateProviders()
}

// ProcessRefund processes a refund request with validation.
// Requirements: 2.1, 2.2, 2.3, 2.4, 9.1, 9.2, 9.3
func (s *RefundService) ProcessRefund(ctx context.Context, req model.RefundRequest) (*model.RefundResponse, error) {
	// Get the payment record
	payment, err := s.payments.Get(ctx, req.PaymentID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apierr.NotFound("payment not found")
		}
		return nil, apierr.InternalError("failed to get payment").WithDetails(err.Error())
	}

	// Validate refund amount using the model's validation method
	// Requirements: 9.1, 9.2, 9.3
	if err := payment.ValidateRefundAmount(req.AmountCents); err != nil {
		if refundErr, ok := err.(*model.RefundValidationError); ok {
			return nil, apierr.BadRequest(refundErr.Message).WithDetails(refundErr.Code)
		}
		return nil, apierr.BadRequest("invalid refund amount")
	}

	// Create refund record
	now := time.Now()
	refundRecord := &model.RefundRecord{
		PaymentID:   req.PaymentID,
		OrderID:     payment.OrderID,
		UserID:      payment.UserID,
		AmountCents: req.AmountCents,
		Reason:      req.Reason,
		Status:      model.RefundStatusPending,
		OperatorID:  req.OperatorID,
		Note:        req.Note,
	}

	if err := s.refunds.Create(ctx, refundRecord); err != nil {
		return nil, apierr.InternalError("failed to create refund record").WithDetails(err.Error())
	}

	// Process refund with payment provider (mock for now)
	client, ok := s.providers[payment.Method]
	if !ok {
		client = genericProvider{}
	}

	tradeNo, raw, refundedAt, err := client.Refund(ctx, payment, req.Reason)
	if err != nil {
		// Update refund record as failed
		refundRecord.Status = model.RefundStatusFailed
		refundRecord.Note = fmt.Sprintf("Provider error: %s", err.Error())
		_ = s.refunds.Update(ctx, refundRecord)

		// Log the failed refund attempt
		s.logRefundOperation(ctx, payment.ID, refundRecord.ID, "refund_failed", req.OperatorID, map[string]any{
			"amount_cents": req.AmountCents,
			"reason":       req.Reason,
			"error":        err.Error(),
		})

		return nil, apierr.InternalError("refund processing failed").WithDetails(err.Error())
	}

	// Update refund record as processed
	refundRecord.Status = model.RefundStatusProcessed
	refundRecord.ProviderTradeNo = tradeNo
	refundRecord.RefundedAt = &refundedAt
	if err := s.refunds.Update(ctx, refundRecord); err != nil {
		return nil, apierr.InternalError("failed to update refund record").WithDetails(err.Error())
	}

	// Update payment record
	payment.RefundedAmountCents += req.AmountCents
	payment.ProviderRaw = raw

	// Check if fully refunded
	if payment.IsFullyRefunded() {
		payment.Status = model.PaymentStatusRefunded
		payment.RefundedAt = &now
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, apierr.InternalError("failed to update payment").WithDetails(err.Error())
	}

	// Update order if fully refunded
	if payment.IsFullyRefunded() {
		order, err := s.orders.Get(ctx, payment.OrderID)
		if err == nil {
			order.Status = model.OrderStatusRefunded
			order.RefundAmountCents = payment.RefundedAmountCents
			order.RefundReason = req.Reason
			order.RefundedAt = &now
			_ = s.orders.Update(ctx, order)
		}
	}

	// Credit wallet if configured
	if s.wallets != nil {
		if err := s.creditWallet(ctx, payment.UserID, req.AmountCents); err != nil {
			// Log but don't fail the refund
			s.logRefundOperation(ctx, payment.ID, refundRecord.ID, "wallet_credit_failed", req.OperatorID, map[string]any{
				"amount_cents": req.AmountCents,
				"error":        err.Error(),
			})
		}
	}

	// Log successful refund operation
	s.logRefundOperation(ctx, payment.ID, refundRecord.ID, string(model.OpActionRefund), req.OperatorID, map[string]any{
		"amount_cents":      req.AmountCents,
		"reason":            req.Reason,
		"remaining_amount":  payment.RemainingRefundableAmount(),
		"is_full_refund":    payment.IsFullyRefunded(),
		"provider_trade_no": tradeNo,
	})

	return &model.RefundResponse{
		RefundRecord:    refundRecord,
		Payment:         payment,
		RemainingAmount: payment.RemainingRefundableAmount(),
	}, nil
}

// GetRefundHistory returns all refund records for a payment.
// Requirements: 2.5
func (s *RefundService) GetRefundHistory(ctx context.Context, paymentID uint64) ([]model.RefundRecord, error) {
	// Verify payment exists
	_, err := s.payments.Get(ctx, paymentID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apierr.NotFound("payment not found")
		}
		return nil, apierr.InternalError("failed to get payment").WithDetails(err.Error())
	}

	records, err := s.refunds.ListByPaymentID(ctx, paymentID)
	if err != nil {
		return nil, apierr.InternalError("failed to get refund history").WithDetails(err.Error())
	}

	return records, nil
}

// GetRefundsByOrder returns all refund records for an order.
func (s *RefundService) GetRefundsByOrder(ctx context.Context, orderID uint64) ([]model.RefundRecord, error) {
	records, err := s.refunds.ListByOrderID(ctx, orderID)
	if err != nil {
		return nil, apierr.InternalError("failed to get refunds by order").WithDetails(err.Error())
	}
	return records, nil
}

// creditWallet credits the refund amount to user's wallet.
func (s *RefundService) creditWallet(ctx context.Context, userID uint64, amount int64) error {
	w, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if err != repository.ErrNotFound {
			return err
		}
		w = &model.Wallet{UserID: userID}
	}
	w.BalanceCents += amount
	return s.wallets.Save(ctx, w)
}

// logRefundOperation logs a refund operation for audit purposes.
func (s *RefundService) logRefundOperation(ctx context.Context, paymentID, refundID uint64, action string, operatorID *uint64, metadata map[string]any) {
	if s.opLogs == nil {
		return
	}

	metadata["refund_id"] = refundID

	var raw []byte
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			raw = b
		}
	}

	log := &model.OperationLog{
		EntityType:   string(model.OpEntityPayment),
		EntityID:     paymentID,
		Action:       action,
		ActorUserID:  operatorID,
		MetadataJSON: raw,
	}

	_ = s.opLogs.Append(ctx, log)
}
