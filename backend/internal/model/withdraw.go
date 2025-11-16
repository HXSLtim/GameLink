package model

import (
	"time"
)

// WithdrawStatus defines withdrawal status
type WithdrawStatus string

const (
	// WithdrawStatusPending withdrawal is pending processing
	WithdrawStatusPending WithdrawStatus = "pending"
	// WithdrawStatusApproved withdrawal has been approved
	WithdrawStatusApproved WithdrawStatus = "approved"
	// WithdrawStatusRejected withdrawal has been rejected
	WithdrawStatusRejected WithdrawStatus = "rejected"
	// WithdrawStatusCompleted withdrawal has been completed
	WithdrawStatusCompleted WithdrawStatus = "completed"
	// WithdrawStatusFailed withdrawal processing failed
	WithdrawStatusFailed WithdrawStatus = "failed"
)

// WithdrawMethod defines withdrawal methods
type WithdrawMethod string

const (
	// WithdrawMethodAlipay Alipay payment method
	WithdrawMethodAlipay WithdrawMethod = "alipay"
	// WithdrawMethodWeChat WeChat payment method
	WithdrawMethodWeChat WithdrawMethod = "wechat"
	// WithdrawMethodBank bank transfer method
	WithdrawMethodBank WithdrawMethod = "bank"
)

// Withdraw represents a withdrawal request
type Withdraw struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID    uint64         `gorm:"not null;index" json:"playerId"`
	UserID      uint64         `gorm:"not null;index" json:"userId"` // redundant field for convenient queries
	AmountCents int64          `gorm:"not null" json:"amountCents"`  // withdrawal amount in cents
	Method      WithdrawMethod `gorm:"type:varchar(32);not null" json:"method"`
	AccountInfo string         `gorm:"type:varchar(255);not null" json:"accountInfo"` // account info (encrypted storage)
	Status      WithdrawStatus `gorm:"type:varchar(32);not null;default:'pending'" json:"status"`
	RejectReason string        `gorm:"type:text" json:"rejectReason"`    // reason for rejection
	AdminRemark  string        `gorm:"type:text" json:"adminRemark"`     // admin remarks
	ProcessedBy  *uint64       `gorm:"index" json:"processedBy"`         // processor ID
	ProcessedAt  *time.Time    `json:"processedAt"`                      // processing time
	CompletedAt  *time.Time    `json:"completedAt"`                      // completion time
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for Withdraw model
func (Withdraw) TableName() string {
	return "withdraws"
}

