package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// PaymentMethod enumerates supported payment channels.
// @Enum wechat, alipay, wallet, combined
type PaymentMethod string

// PaymentMethod values enumerate supported channels.
const (
	PaymentMethodWeChat   PaymentMethod = "wechat"
	PaymentMethodAlipay   PaymentMethod = "alipay"
	PaymentMethodWallet   PaymentMethod = "wallet"   // 纯钱包支付
	PaymentMethodCombined PaymentMethod = "combined" // 组合支付（钱包+第三方）
)

// PaymentStatus enumerates payment states.
// @Enum pending, paid, failed, refunded
type PaymentStatus string

// PaymentStatus values enumerate payment states.
const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Payment records a payment attempt/result for an order.
//
// Covering Index Notes:
//   - idx_payments_user_status_created_covering: PostgreSQL covering index for payment history queries
//   - Created via migration: api/migrations/0001_add_covering_indexes.sql
//   - Covers: SELECT id, amount_cents, payment_method FROM payments WHERE user_id = ? AND status IN (?) ORDER BY created_at DESC
//   - Index columns: (user_id, status, created_at DESC)
//   - INCLUDE columns: (id, amount_cents, payment_method, provider_trade_no)
//   - Benefit: Index-only scan for payment history, no heap fetch needed
type Payment struct {
	Base
	OrderID             uint64          `json:"orderId" gorm:"column:order_id;not null;index"`
	UserID              uint64          `json:"userId" gorm:"column:user_id;not null;index"`
	Method              PaymentMethod   `json:"method" gorm:"size:32"`
	AmountCents         int64           `json:"amountCents" gorm:"column:amount_cents"`
	Currency            Currency        `json:"currency,omitempty" gorm:"type:char(3)"` // default CNY
	Status              PaymentStatus   `json:"status" gorm:"size:32;index"`
	ProviderTradeNo     string          `json:"providerTradeNo,omitempty" gorm:"column:provider_trade_no;size:128"`
	ProviderRaw         json.RawMessage `json:"providerRaw,omitempty" gorm:"column:provider_raw;type:json"` // provider response payload
	PaidAt              *time.Time      `json:"paidAt,omitempty" gorm:"column:paid_at"`
	RefundedAt          *time.Time      `json:"refundedAt,omitempty" gorm:"column:refunded_at"`
	RefundedAmountCents int64           `json:"refundedAmountCents" gorm:"column:refunded_amount_cents;default:0"`     // 已退款金额（分）
	CollectionEntityID  *uint64         `json:"collectionEntityId,omitempty" gorm:"column:collection_entity_id;index"` // 收款主体ID
	MerchantNo          string          `json:"merchantNo,omitempty" gorm:"column:merchant_no;size:64"`                // 实际使用的商户号

	// 组合支付字段
	WalletAmountCents     int64         `json:"walletAmountCents" gorm:"column:wallet_amount_cents;default:0"`          // 钱包支付金额（分）
	ThirdPartyMethod      PaymentMethod `json:"thirdPartyMethod,omitempty" gorm:"column:third_party_method;size:32"`    // 第三方支付方式（组合支付时使用）
	ThirdPartyAmountCents int64         `json:"thirdPartyAmountCents" gorm:"column:third_party_amount_cents;default:0"` // 第三方支付金额（分）

	// Relations + FKs
	Order Order `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:OrderID;references:ID"`
	User  User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
}

// HasRequiredFields checks if the payment record has all required fields populated.
// This is used for property testing to verify payment record completeness.
func (p *Payment) HasRequiredFields() bool {
	return p.OrderID != 0 &&
		p.UserID != 0 &&
		p.Method != "" &&
		p.Status != "" &&
		!p.CreatedAt.IsZero()
}

// RemainingRefundableAmount returns the amount that can still be refunded.
func (p *Payment) RemainingRefundableAmount() int64 {
	if p.Status != PaymentStatusPaid {
		return 0
	}
	return p.AmountCents - p.RefundedAmountCents
}

// RefundValidationError represents a refund validation error with specific reason.
type RefundValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *RefundValidationError) Error() string {
	return e.Message
}

// Refund validation error codes
const (
	RefundErrCodeInvalidAmount    = "INVALID_REFUND_AMOUNT"
	RefundErrCodeExceedsRemaining = "EXCEEDS_REMAINING_AMOUNT"
	RefundErrCodeInvalidStatus    = "INVALID_PAYMENT_STATUS"
	RefundErrCodeInvalidPrecision = "INVALID_AMOUNT_PRECISION"
	RefundErrCodeAlreadyRefunded  = "ALREADY_FULLY_REFUNDED"
)

// ValidateRefundAmount validates the refund amount against the payment.
// Returns nil if valid, or a RefundValidationError if invalid.
// Requirements: 9.1, 9.2, 9.3
func (p *Payment) ValidateRefundAmount(amountCents int64) error {
	// 9.1: Validate refund amount is positive
	if amountCents <= 0 {
		return &RefundValidationError{
			Code:    RefundErrCodeInvalidAmount,
			Message: "refund amount must be positive",
		}
	}

	// Check payment status - only paid payments can be refunded
	if p.Status != PaymentStatusPaid {
		return &RefundValidationError{
			Code:    RefundErrCodeInvalidStatus,
			Message: fmt.Sprintf("payment status must be paid, current: %s", p.Status),
		}
	}

	// Check if already fully refunded
	remaining := p.RemainingRefundableAmount()
	if remaining <= 0 {
		return &RefundValidationError{
			Code:    RefundErrCodeAlreadyRefunded,
			Message: "payment has already been fully refunded",
		}
	}

	// 9.2: Validate refund amount does not exceed remaining refundable amount
	if amountCents > remaining {
		return &RefundValidationError{
			Code:    RefundErrCodeExceedsRemaining,
			Message: fmt.Sprintf("refund amount %d exceeds remaining refundable amount %d", amountCents, remaining),
		}
	}

	return nil
}

// CanRefund checks if the payment can be refunded (any amount).
func (p *Payment) CanRefund() bool {
	return p.Status == PaymentStatusPaid && p.RemainingRefundableAmount() > 0
}

// IsFullyRefunded checks if the payment has been fully refunded.
func (p *Payment) IsFullyRefunded() bool {
	return p.RefundedAmountCents >= p.AmountCents
}

// IsPartiallyRefunded checks if the payment has been partially refunded.
func (p *Payment) IsPartiallyRefunded() bool {
	return p.RefundedAmountCents > 0 && p.RefundedAmountCents < p.AmountCents
}

// PaymentStatusTransition represents a status transition for a payment.
type PaymentStatusTransition struct {
	From PaymentStatus
	To   PaymentStatus
}

// ValidPaymentTransitions defines all valid payment status transitions.
// Requirements: 2.3, 2.4, 8.2
// Valid transitions:
// - pending -> paid (payment successful)
// - pending -> failed (payment failed)
// - paid -> refunded (full refund)
var ValidPaymentTransitions = map[PaymentStatusTransition]bool{
	{From: PaymentStatusPending, To: PaymentStatusPaid}:   true,
	{From: PaymentStatusPending, To: PaymentStatusFailed}: true,
	{From: PaymentStatusPaid, To: PaymentStatusRefunded}:  true,
	// Same status transitions are always valid (no-op)
	{From: PaymentStatusPending, To: PaymentStatusPending}:   true,
	{From: PaymentStatusPaid, To: PaymentStatusPaid}:         true,
	{From: PaymentStatusFailed, To: PaymentStatusFailed}:     true,
	{From: PaymentStatusRefunded, To: PaymentStatusRefunded}: true,
}

// IsValidStatusTransition checks if a status transition is valid.
// Requirements: 2.3, 2.4, 8.2
func IsValidStatusTransition(from, to PaymentStatus) bool {
	transition := PaymentStatusTransition{From: from, To: to}
	return ValidPaymentTransitions[transition]
}

// ValidateStatusTransition validates if the payment can transition to the new status.
// Returns nil if valid, or an error if invalid.
func (p *Payment) ValidateStatusTransition(newStatus PaymentStatus) error {
	if !IsValidStatusTransition(p.Status, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", p.Status, newStatus)
	}
	return nil
}

// GetAllowedTransitions returns all valid next statuses for the current status.
func (p *Payment) GetAllowedTransitions() []PaymentStatus {
	var allowed []PaymentStatus
	allStatuses := []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusPaid,
		PaymentStatusFailed,
		PaymentStatusRefunded,
	}
	for _, status := range allStatuses {
		if IsValidStatusTransition(p.Status, status) && status != p.Status {
			allowed = append(allowed, status)
		}
	}
	return allowed
}

// RefundRecord represents a single refund transaction for a payment.
// A payment can have multiple partial refunds.
type RefundRecord struct {
	Base
	PaymentID       uint64       `json:"paymentId" gorm:"column:payment_id;not null;index"`
	OrderID         uint64       `json:"orderId" gorm:"column:order_id;not null;index"`
	UserID          uint64       `json:"userId" gorm:"column:user_id;not null;index"`
	AmountCents     int64        `json:"amountCents" gorm:"column:amount_cents;not null"`
	Reason          string       `json:"reason,omitempty" gorm:"type:text"`
	Status          RefundStatus `json:"status" gorm:"size:32;index"`
	ProviderTradeNo string       `json:"providerTradeNo,omitempty" gorm:"column:provider_trade_no;size:128"`
	OperatorID      *uint64      `json:"operatorId,omitempty" gorm:"column:operator_id;index"` // Admin who initiated the refund
	RefundedAt      *time.Time   `json:"refundedAt,omitempty" gorm:"column:refunded_at"`
	Note            string       `json:"note,omitempty" gorm:"type:text"` // Internal note

	// Relations
	Payment Payment `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:PaymentID;references:ID"`
}

// RefundStatus enumerates refund states.
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
	RefundStatusFailed    RefundStatus = "failed"
)

// RefundRequest represents a request to refund a payment.
type RefundRequest struct {
	PaymentID   uint64  `json:"paymentId" binding:"required"`
	AmountCents int64   `json:"amountCents" binding:"required,gt=0"`
	Reason      string  `json:"reason" binding:"required"`
	Note        string  `json:"note,omitempty"`
	OperatorID  *uint64 `json:"operatorId,omitempty"`
}

// RefundResponse represents the result of a refund operation.
type RefundResponse struct {
	RefundRecord    *RefundRecord `json:"refundRecord"`
	Payment         *Payment      `json:"payment"`
	RemainingAmount int64         `json:"remainingAmount"`
}
