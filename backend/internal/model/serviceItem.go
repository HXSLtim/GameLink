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

// ServiceItem 服务项目表 (统一管理所有服务类型，包括礼物)
type ServiceItem struct {
	ID             uint64                 `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemCode       string                 `gorm:"type:varchar(32);not null;uniqueIndex" json:"itemCode"`
	Name           string                 `gorm:"type:varchar(128);not null" json:"name"`
	Description    string                 `gorm:"type:text" json:"description"`
	Category       string                 `gorm:"type:varchar(32);not null;default:'escort';index" json:"category"` // 统一为 'escort'
	SubCategory    ServiceItemSubCategory `gorm:"type:varchar(32);not null;index" json:"subCategory"`               // solo/team/gift
	GameID         *uint64                `gorm:"index" json:"gameId"`
	PlayerID       *uint64                `gorm:"index" json:"playerId"`
	RankLevel      string                 `gorm:"type:varchar(32)" json:"rankLevel"`
	BasePriceCents int64                  `gorm:"not null;default:0" json:"basePriceCents"`
	ServiceHours   int                    `gorm:"not null;default:0" json:"serviceHours"`               // 服务时长（小时），礼物为0
	CommissionRate float64                `gorm:"type:decimal(5,2);default:0.20" json:"commissionRate"` // 抽成比例
	MinUsers       int                    `gorm:"default:1" json:"minUsers"`
	MaxPlayers     int                    `gorm:"default:1" json:"maxPlayers"`
	Tags           string                 `gorm:"type:json" json:"tags"` // JSON数组
	IconURL        string                 `gorm:"type:varchar(255)" json:"iconUrl"`
	IsActive       bool                   `gorm:"default:true;index" json:"isActive"`
	SortOrder      int                    `gorm:"default:0" json:"sortOrder"`
	CreatedAt      time.Time              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"autoUpdateTime" json:"updatedAt"`
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
