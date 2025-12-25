package model

import "time"

// SettlementStatus 结算状态
type SettlementStatus string

const (
	SettlementStatusPending  SettlementStatus = "pending"  // T+7 售后期内，收入冻结中
	SettlementStatusDisputed SettlementStatus = "disputed" // 发生争议，收入继续冻结
	SettlementStatusSettled  SettlementStatus = "settled"  // T+7 结束且无争议，收入已解冻
)

// CommissionRuleType 抽成规则类型
type CommissionRuleType string

const (
	CommissionRuleTypeDefault CommissionRuleType = "default" // 默认规则
	CommissionRuleTypeSpecial CommissionRuleType = "special" // 特殊规则
	CommissionRuleTypeGift    CommissionRuleType = "gift"    // 礼物规则
)

// CommissionRule 抽成规则
type CommissionRule struct {
	ID          uint64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string             `gorm:"type:varchar(128);not null" json:"name"`
	Description string             `gorm:"type:text" json:"description"`
	Type        CommissionRuleType `gorm:"type:varchar(32);not null;index" json:"type"` // default/special/gift
	Rate        int                `gorm:"not null" json:"rate"`                        // 抽成比例（百分比，例如20表示20%）
	IsActive    bool               `gorm:"default:true;index" json:"isActive"`
	GameID      *uint64            `gorm:"index" json:"gameId"`                 // 特定游戏的抽成
	PlayerID    *uint64            `gorm:"index" json:"playerId"`               // 特定陪玩师的抽成
	ServiceType *string            `gorm:"type:varchar(64)" json:"serviceType"` // 特定服务类型的抽成
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (CommissionRule) TableName() string {
	return "commission_rules"
}

// CommissionRecord 抽成记录
type CommissionRecord struct {
	ID                uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID           uint64           `gorm:"not null;index" json:"orderId"`
	PlayerID          uint64           `gorm:"not null;index" json:"playerId"`
	TotalAmountCents  int64            `gorm:"not null" json:"totalAmountCents"`                                                          // 订单总金额
	CommissionRate    int              `gorm:"not null" json:"commissionRate"`                                                            // 抽成比例
	CommissionCents   int64            `gorm:"not null" json:"commissionCents"`                                                           // 平台抽成金额
	PlayerIncomeCents int64            `gorm:"not null" json:"playerIncomeCents"`                                                         // 陪玩师收入
	SettlementStatus  SettlementStatus `gorm:"type:varchar(32);not null;default:'pending';index" json:"settlementStatus"`                 // pending/disputed/settled
	SettlementMonth   string           `gorm:"type:varchar(7);index;index:idx_player_settlement_month,priority:2" json:"settlementMonth"` // YYYY-MM，复合索引
	SettledAt         *time.Time       `json:"settledAt"`
	CreatedAt         time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime" json:"updatedAt"`

	// 复合索引：PlayerID + SettlementMonth
	_ struct{} `gorm:"index:idx_player_settlement_month,priority:1" json:"-"` // PlayerID 在上面已定义
}

// TableName 指定表名
func (CommissionRecord) TableName() string {
	return "commission_records"
}

// MonthlySettlementStatus 月度结算状态
type MonthlySettlementStatus string

const (
	MonthlySettlementStatusPending   MonthlySettlementStatus = "pending"   // 待确认
	MonthlySettlementStatusConfirmed MonthlySettlementStatus = "confirmed" // 已确认
	MonthlySettlementStatusPaid      MonthlySettlementStatus = "paid"      // 已支付
)

// MonthlySettlement 月度结算
type MonthlySettlement struct {
	ID                   uint64                  `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID             uint64                  `gorm:"not null;uniqueIndex:idx_player_month,priority:1" json:"playerId"`
	SettlementMonth      string                  `gorm:"type:varchar(7);not null;uniqueIndex:idx_player_month,priority:2;index" json:"settlementMonth"` // YYYY-MM
	TotalOrderCount      int64                   `gorm:"not null" json:"totalOrderCount"`
	TotalAmountCents     int64                   `gorm:"not null" json:"totalAmountCents"`
	TotalCommissionCents int64                   `gorm:"not null" json:"totalCommissionCents"`
	TotalIncomeCents     int64                   `gorm:"not null" json:"totalIncomeCents"`
	BonusCents           int64                   `gorm:"default:0" json:"bonusCents"`                                     // 奖金
	FinalIncomeCents     int64                   `gorm:"not null" json:"finalIncomeCents"`                                // 最终收入
	Status               MonthlySettlementStatus `gorm:"type:varchar(32);not null;default:'pending';index" json:"status"` // pending/confirmed/paid
	IncomeRank           *int                    `json:"incomeRank"`                                                      // 收入排名
	OrderRank            *int                    `json:"orderRank"`                                                       // 订单数排名
	QualityRank          *int                    `json:"qualityRank"`                                                     // 质量排名
	CreatedAt            time.Time               `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time               `gorm:"autoUpdateTime" json:"updatedAt"`
	SettledAt            *time.Time              `json:"settledAt"`
}

// TableName 指定表名
func (MonthlySettlement) TableName() string {
	return "monthly_settlements"
}
