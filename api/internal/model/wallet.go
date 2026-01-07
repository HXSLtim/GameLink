package model

import "time"

// Wallet 用户钱包，记录用户的余额信息
// @Description 用户钱包模型，存储用户的账户余额和冻结金额
// @Example {"id": 1, "userId": 1001, "balanceCents": 100000, "frozenCents": 5000, "version": 1, "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-15T12:30:00Z"}
//
// Optimistic Locking:
//   - Version field implements optimistic locking to prevent concurrent update conflicts
//   - Every update must include WHERE version = ? AND increment version
//   - If affected rows = 0, it means concurrent modification occurred, should retry
type Wallet struct {
	Base
	// 用户ID
	// @Example 1001
	UserID uint64 `json:"userId" gorm:"column:user_id;uniqueIndex;not null"`
	// 可用余额，单位为分
	// @Example 100000
	BalanceCents int64 `json:"balanceCents" gorm:"column:balance_cents;default:0"`
	// 冻结金额，单位为分
	// @Example 5000
	FrozenCents int64 `json:"frozenCents" gorm:"column:frozen_cents;default:0"`
	// 乐观锁版本号，每次更新自增
	// @Example 1
	Version int64 `json:"version" gorm:"column:version;default:1;not null"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (Wallet) TableName() string {
	return "wallets"
}
