package model

import (
	"time"
)

// WithdrawStatus 提现状态
type WithdrawStatus string

const (
	// WithdrawStatusPending 待处理
	WithdrawStatusPending WithdrawStatus = "pending"
	// WithdrawStatusApproved 已批准
	WithdrawStatusApproved WithdrawStatus = "approved"
	// WithdrawStatusRejected 已拒绝
	WithdrawStatusRejected WithdrawStatus = "rejected"
	// WithdrawStatusCompleted 已完成
	WithdrawStatusCompleted WithdrawStatus = "completed"
	// WithdrawStatusFailed 失败
	WithdrawStatusFailed WithdrawStatus = "failed"
)

// WithdrawMethod 提现方式
type WithdrawMethod string

const (
	// WithdrawMethodAlipay 支付宝
	WithdrawMethodAlipay WithdrawMethod = "alipay"
	// WithdrawMethodWeChat 微信
	WithdrawMethodWeChat WithdrawMethod = "wechat"
	// WithdrawMethodBank 银行卡
	WithdrawMethodBank WithdrawMethod = "bank"
)

// Withdraw 提现记录
type Withdraw struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID    uint64         `gorm:"not null;index" json:"playerId"`
	UserID      uint64         `gorm:"not null;index" json:"userId"` // 冗余字段，方便查询
	AmountCents int64          `gorm:"not null" json:"amountCents"`  // 提现金额（分）
	Method      WithdrawMethod `gorm:"type:varchar(32);not null" json:"method"`
	AccountInfo string         `gorm:"type:varchar(255);not null" json:"accountInfo"` // 账号信息（加密存储）
	Status      WithdrawStatus `gorm:"type:varchar(32);not null;default:'pending'" json:"status"`
	RejectReason string        `gorm:"type:text" json:"rejectReason"`    // 拒绝原因
	AdminRemark  string        `gorm:"type:text" json:"adminRemark"`     // 管理员备注
	ProcessedBy  *uint64       `gorm:"index" json:"processedBy"`         // 处理人ID
	ProcessedAt  *time.Time    `json:"processedAt"`                      // 处理时间
	CompletedAt  *time.Time    `json:"completedAt"`                      // 完成时间
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (Withdraw) TableName() string {
	return "withdraws"
}

