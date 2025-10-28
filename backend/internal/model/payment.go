package model

import (
	"encoding/json"
	"time"
)

// PaymentMethod enumerates supported payment channels.
type PaymentMethod string

// PaymentMethod values enumerate supported channels.
const (
	PaymentMethodWeChat PaymentMethod = "wechat"
	PaymentMethodAlipay PaymentMethod = "alipay"
)

// PaymentStatus enumerates payment states.
type PaymentStatus string

// PaymentStatus values enumerate payment states.
const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Payment records a payment attempt/result for an order.
type Payment struct {
	Base
	OrderID         uint64          `json:"order_id" gorm:"not null;index"`
	UserID          uint64          `json:"user_id" gorm:"not null;index"`
	Method          PaymentMethod   `json:"method" gorm:"size:32"`
	AmountCents     int64           `json:"amount_cents"`
	Currency        Currency        `json:"currency,omitempty" gorm:"type:char(3)"` // default CNY
	Status          PaymentStatus   `json:"status" gorm:"size:32;index"`
	ProviderTradeNo string          `json:"provider_trade_no,omitempty" gorm:"size:128"`
	ProviderRaw     json.RawMessage `json:"provider_raw,omitempty" gorm:"type:json"` // provider response payload
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	RefundedAt      *time.Time      `json:"refunded_at,omitempty"`

	// Relations + FKs
	Order Order `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:OrderID;references:ID"`
	User  User  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
}
