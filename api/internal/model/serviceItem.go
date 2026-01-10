package model

import "time"

// ServiceItemSubCategory 服务子类别
type ServiceItemSubCategory string

const (
	// SubCategorySolo 单人护航
	SubCategorySolo ServiceItemSubCategory = "solo"
	// SubCategoryTeam 团队护航
	SubCategoryTeam ServiceItemSubCategory = "team"
	// SubCategoryGift 礼物
	SubCategoryGift ServiceItemSubCategory = "gift"
)

// UsageLimitType 使用限制类型
type UsageLimitType string

const (
	UsageLimitNone    UsageLimitType = "none"    // 无限制
	UsageLimitOnce    UsageLimitType = "once"    // 只能用一次（终身）
	UsageLimitDaily   UsageLimitType = "daily"   // 每日限制
	UsageLimitWeekly  UsageLimitType = "weekly"  // 每周限制
	UsageLimitMonthly UsageLimitType = "monthly" // 每月限制
)

// ServiceItem 服务项目表 (统一管理所有服务类型，包括礼物)
type ServiceItem struct {
	ID             uint64                 `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemCode       string                 `gorm:"type:varchar(32);not null;uniqueIndex" json:"itemCode"`
	Name           string                 `gorm:"type:varchar(128);not null" json:"name"`
	Description    string                 `gorm:"type:text" json:"description"`
	Category       string                 `gorm:"type:varchar(32);not null;default:'escort';index" json:"category"` // 统一为 'escort'
	SubCategory    ServiceItemSubCategory `gorm:"type:varchar(32);not null;index" json:"subCategory"`               // solo/team/gift
	GameID         *uint64                `gorm:"index" json:"gameId"`
	CategoryID     *uint64                `gorm:"index" json:"categoryId,omitempty"` // 游戏分类ID
	PlayerID       *uint64                `gorm:"index" json:"playerId"`
	RankLevel      string                 `gorm:"type:varchar(32)" json:"rankLevel"`
	BasePriceCents int64                  `gorm:"not null;default:0" json:"basePriceCents"`
	ServiceHours   int                    `gorm:"not null;default:0" json:"serviceHours"`               // 服务时长（小时），礼物为0
	CommissionRate float64                `gorm:"type:decimal(5,2);default:0.20" json:"commissionRate"` // 抽成比例
	MinUsers       int                    `gorm:"default:1" json:"minUsers"`
	MaxPlayers     int                    `gorm:"default:1" json:"maxPlayers"`
	Tags           string                 `gorm:"type:json;default:'[]'" json:"tags"` // JSON数组
	IconURL        string                 `gorm:"type:varchar(255)" json:"iconUrl"`
	IsActive       bool                   `gorm:"default:true;index" json:"isActive"`
	SortOrder      int                    `gorm:"default:0" json:"sortOrder"`
	CreatedAt      time.Time              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"autoUpdateTime" json:"updatedAt"`

	// 多人服务配置
	RequiredPlayers int `gorm:"column:required_players;default:1" json:"requiredPlayers"` // 需要的陪玩师数量

	// VIP 专属价（预留）
	VipPriceCents *int64 `gorm:"column:vip_price_cents" json:"vipPriceCents,omitempty"` // VIP专属价格（分），nil表示无专属价

	// 使用限制配置
	UsageLimitType  UsageLimitType `gorm:"column:usage_limit_type;size:32;default:'none'" json:"usageLimitType"` // 限制类型
	UsageLimitCount int            `gorm:"column:usage_limit_count;default:0" json:"usageLimitCount"`            // 限制次数（0=无限制）
	MaxPerOrder     int            `gorm:"column:max_per_order;default:0" json:"maxPerOrder"`                    // 单次购买数量限制（0=无限制）

	// Relations
	GameCategory *GameCategory `gorm:"foreignKey:CategoryID" json:"gameCategory,omitempty"` // 游戏分类关联
}

// TableName 指定表名
func (ServiceItem) TableName() string {
	return "service_items"
}

// IsGift 判断是否为礼物
func (s *ServiceItem) IsGift() bool {
	return s.SubCategory == SubCategoryGift
}

// CalculateCommission 计算抽成
func (s *ServiceItem) CalculateCommission(quantity int) (platformCommission, playerIncome int64) {
	totalAmount := s.BasePriceCents * int64(quantity)
	platformCommission = int64(float64(totalAmount) * s.CommissionRate)
	playerIncome = totalAmount - platformCommission
	return platformCommission, playerIncome
}
