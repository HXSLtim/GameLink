package model

import "time"

// CommissionRule 抽成规则
type CommissionRule struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"type:varchar(128);not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	Type             string    `gorm:"type:varchar(32);not null" json:"type"` // default/special/gift
	Rate             int       `gorm:"not null" json:"rate"`                  // 抽成比例（百分比，例如20表示20%）
	IsActive         bool      `gorm:"default:true" json:"isActive"`
	GameID           *uint64   `gorm:"index" json:"gameId"`     // 特定游戏的抽成
	PlayerID         *uint64   `gorm:"index" json:"playerId"`   // 特定陪玩师的抽成
	ServiceType      *string   `gorm:"type:varchar(64)" json:"serviceType"` // 特定服务类型的抽成
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CommissionRule) TableName() string {
	return "commission_rules"
}

// CommissionRecord 抽成记录
type CommissionRecord struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID            uint64    `gorm:"not null;index" json:"orderId"`
	PlayerID           uint64    `gorm:"not null;index" json:"playerId"`
	TotalAmountCents   int64     `gorm:"not null" json:"totalAmountCents"`       // 订单总金额
	CommissionRate     int       `gorm:"not null" json:"commissionRate"`         // 抽成比例
	CommissionCents    int64     `gorm:"not null" json:"commissionCents"`        // 平台抽成金额
	PlayerIncomeCents  int64     `gorm:"not null" json:"playerIncomeCents"`      // 陪玩师收入
	SettlementStatus   string    `gorm:"type:varchar(32);not null;default:'pending'" json:"settlementStatus"` // pending/settled
	SettlementMonth    string    `gorm:"type:varchar(7);index" json:"settlementMonth"` // YYYY-MM
	SettledAt          *time.Time `json:"settledAt"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CommissionRecord) TableName() string {
	return "commission_records"
}

// MonthlySettlement 月度结算
type MonthlySettlement struct {
	ID                    uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID              uint64    `gorm:"not null;index" json:"playerId"`
	SettlementMonth       string    `gorm:"type:varchar(7);not null;index" json:"settlementMonth"` // YYYY-MM
	TotalOrderCount       int64     `gorm:"not null" json:"totalOrderCount"`
	TotalAmountCents      int64     `gorm:"not null" json:"totalAmountCents"`
	TotalCommissionCents  int64     `gorm:"not null" json:"totalCommissionCents"`
	TotalIncomeCents      int64     `gorm:"not null" json:"totalIncomeCents"`
	BonusCents            int64     `gorm:"default:0" json:"bonusCents"`        // 奖金
	FinalIncomeCents      int64     `gorm:"not null" json:"finalIncomeCents"`   // 最终收入
	Status                string    `gorm:"type:varchar(32);not null;default:'pending'" json:"status"` // pending/confirmed/paid
	IncomeRank            *int      `json:"incomeRank"`                         // 收入排名
	OrderRank             *int      `json:"orderRank"`                          // 订单数排名
	QualityRank           *int      `json:"qualityRank"`                        // 质量排名
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	SettledAt             *time.Time `json:"settledAt"`
}

// TableName 指定表名
func (MonthlySettlement) TableName() string {
	return "monthly_settlements"
}

