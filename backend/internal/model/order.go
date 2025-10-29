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
	UserID            uint64      `json:"userId" gorm:"column:user_id;not null;index"`      // customer
	PlayerID          uint64      `json:"playerId,omitempty" gorm:"column:player_id;index"` // assigned player (optional until accepted)
	GameID            uint64      `json:"gameId" gorm:"column:game_id;not null;index"`
	Title             string      `json:"title" gorm:"size:128"`
	Description       string      `json:"description,omitempty" gorm:"type:text"`
	Status            OrderStatus `json:"status" gorm:"size:32;index"`
	PriceCents        int64       `json:"priceCents" gorm:"column:price_cents;check:price_cents >= 0"`
	Currency          Currency    `json:"currency,omitempty" gorm:"type:char(3)"` // default CNY
	ScheduledStart    *time.Time  `json:"scheduledStart,omitempty" gorm:"column:scheduled_start"`
	ScheduledEnd      *time.Time  `json:"scheduledEnd,omitempty" gorm:"column:scheduled_end"`
	CancelReason      string      `json:"cancelReason,omitempty" gorm:"column:cancel_reason;type:text"`
	StartedAt         *time.Time  `json:"startedAt,omitempty" gorm:"column:started_at"`
	CompletedAt       *time.Time  `json:"completedAt,omitempty" gorm:"column:completed_at"`
	RefundAmountCents int64       `json:"refundAmountCents,omitempty" gorm:"column:refund_amount_cents;default:0"`
	RefundReason      string      `json:"refundReason,omitempty" gorm:"column:refund_reason;type:text"`
	RefundedAt        *time.Time  `json:"refundedAt,omitempty" gorm:"column:refunded_at"`

	// Relations + FKs
	User   User    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
	Player *Player `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PlayerID;references:ID"`
	Game   Game    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:GameID;references:ID"`
}
