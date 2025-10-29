package model

import "time"

// OrderStatus defines lifecycle states for an order.
type OrderStatus string

// OrderStatus values define the lifecycle of an order.
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCanceled   OrderStatus = "canceled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// Order represents a play session request and assignment.
type Order struct {
	Base
	UserID            uint64      `json:"user_id" gorm:"not null;index"`    // customer
	PlayerID          uint64      `json:"player_id,omitempty" gorm:"index"` // assigned player (optional until accepted)
	GameID            uint64      `json:"game_id" gorm:"not null;index"`
	Title             string      `json:"title" gorm:"size:128"`
	Description       string      `json:"description,omitempty" gorm:"type:text"`
	Status            OrderStatus `json:"status" gorm:"size:32;index"`
	PriceCents        int64       `json:"price_cents" gorm:"check:price_cents >= 0"`
	Currency          Currency    `json:"currency,omitempty" gorm:"type:char(3)"` // default CNY
	ScheduledStart    *time.Time  `json:"scheduled_start,omitempty"`
	ScheduledEnd      *time.Time  `json:"scheduled_end,omitempty"`
	CancelReason      string      `json:"cancel_reason,omitempty" gorm:"type:text"`
	StartedAt         *time.Time  `json:"started_at,omitempty"`
	CompletedAt       *time.Time  `json:"completed_at,omitempty"`
	RefundAmountCents int64       `json:"refund_amount_cents,omitempty" gorm:"default:0"`
	RefundReason      string      `json:"refund_reason,omitempty" gorm:"type:text"`
	RefundedAt        *time.Time  `json:"refunded_at,omitempty"`

	// Relations + FKs
	User   User    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
	Player *Player `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PlayerID;references:ID"`
	Game   Game    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:GameID;references:ID"`
}
