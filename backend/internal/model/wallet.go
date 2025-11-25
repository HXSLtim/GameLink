package model

import "time"

// Wallet 用户钱包
type Wallet struct {
	Base
	UserID       uint64    `json:"userId" gorm:"column:user_id;uniqueIndex;not null"`
	BalanceCents int64     `json:"balanceCents" gorm:"column:balance_cents;default:0"`
	FrozenCents  int64     `json:"frozenCents" gorm:"column:frozen_cents;default:0"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (Wallet) TableName() string {
	return "wallets"
}
