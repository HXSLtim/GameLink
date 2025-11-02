package model

import "time"

// RankingType 排名类型
type RankingType string

const (
	// RankingTypeIncome 收入排名
	RankingTypeIncome RankingType = "income"
	// RankingTypeOrderCount 订单数量排名
	RankingTypeOrderCount RankingType = "order_count"
	// RankingTypeQuality 服务质量排名
	RankingTypeQuality RankingType = "quality"
	// RankingTypePopularity 人气排名
	RankingTypePopularity RankingType = "popularity"
)

// PlayerRanking 陪玩师排名
type PlayerRanking struct {
	ID           uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID     uint64      `gorm:"not null;index" json:"playerId"`
	RankingType  RankingType `gorm:"type:varchar(32);not null;index" json:"rankingType"`
	Period       string      `gorm:"type:varchar(32);not null;index" json:"period"` // daily/weekly/monthly/yearly
	PeriodValue  string      `gorm:"type:varchar(32);not null;index" json:"periodValue"` // YYYY-MM-DD, YYYY-WW, YYYY-MM
	Rank         int         `gorm:"not null" json:"rank"`
	Score        float64     `gorm:"not null" json:"score"` // 排名分数
	OrderCount   int64       `json:"orderCount"`
	IncomeCents  int64       `json:"incomeCents"`
	AvgRating    float32     `json:"avgRating"`
	BonusCents   int64       `gorm:"default:0" json:"bonusCents"` // 排名奖金
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (PlayerRanking) TableName() string {
	return "player_rankings"
}

// RankingCommissionConfig 排名抽成配置（使用JSON序列化存储阶梯规则）
type RankingCommissionConfig struct {
	ID          uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string      `gorm:"type:varchar(128);not null" json:"name"`
	RankingType RankingType `gorm:"type:varchar(32);not null;index" json:"rankingType"` // income/order_count
	Period      string      `gorm:"type:varchar(32);not null" json:"period"`            // monthly
	Month       string      `gorm:"type:varchar(7);not null;index" json:"month"`        // YYYY-MM
	RulesJSON   string      `gorm:"type:text;not null" json:"rulesJson"`                // JSON序列化的阶梯规则
	Description string      `gorm:"type:text" json:"description"`
	IsActive    bool        `gorm:"default:true;index" json:"isActive"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (RankingCommissionConfig) TableName() string {
	return "ranking_commission_configs"
}

// RankingCommissionRule 排名抽成规则（用于JSON序列化）
// 示例：[{"rankStart":1,"rankEnd":3,"commissionRate":10},{"rankStart":4,"rankEnd":8,"commissionRate":12}]
type RankingCommissionRule struct {
	RankStart      int `json:"rankStart"`      // 排名开始
	RankEnd        int `json:"rankEnd"`        // 排名结束  
	CommissionRate int `json:"commissionRate"` // 抽成比例（%）
}

// RankingReward 排名奖励（奖金）
type RankingReward struct {
	ID          uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	RankingType RankingType `gorm:"type:varchar(32);not null;index" json:"rankingType"`
	Period      string      `gorm:"type:varchar(32);not null" json:"period"`
	RankStart   int         `gorm:"not null" json:"rankStart"` // 排名开始（例如1）
	RankEnd     int         `gorm:"not null" json:"rankEnd"`   // 排名结束（例如10）
	RewardType  string      `gorm:"type:varchar(32);not null" json:"rewardType"` // fixed/percentage
	RewardValue int64       `gorm:"not null" json:"rewardValue"` // 固定金额（分）或百分比
	Description string      `gorm:"type:text" json:"description"`
	IsActive    bool        `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (RankingReward) TableName() string {
	return "ranking_rewards"
}

