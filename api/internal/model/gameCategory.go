package model

import "gorm.io/gorm"

// GameCategory represents a game category classification.
// @Description 游戏分类模型，用于对游戏进行分类管理，如MOBA、FPS、RPG等
// @Example {"id": 1, "name": "MOBA", "description": "多人在线战术竞技游戏", "iconUrl": "https://example.com/moba.png", "sortOrder": 1, "isActive": true, "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}
type GameCategory struct {
	Base
	// 分类名称（唯一）
	// @Example MOBA
	Name string `json:"name" gorm:"column:name;type:varchar(50);uniqueIndex;not null" binding:"required,min=1,max=50"`
	// 分类描述
	// @Example 多人在线战术竞技游戏，如英雄联盟、DOTA2等
	Description string `json:"description,omitempty" gorm:"column:description;type:text" binding:"max=1000"`
	// 分类图标URL
	// @Example https://example.com/categories/moba.png
	IconURL string `json:"iconUrl,omitempty" gorm:"column:icon_url;type:varchar(255)" binding:"omitempty,url,max=255"`
	// 排序权重，数值越小越靠前
	// @Example 1
	SortOrder int `json:"sortOrder" gorm:"column:sort_order;default:0;index"`
	// 是否启用
	// @Example true
	IsActive bool `json:"isActive" gorm:"column:is_active;default:true;index"`

	// Relations
	// 该分类下的游戏列表（非持久化字段）
	Games []Game `json:"games,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	// 该分类下的服务项目列表（非持久化字段）
	ServiceItems []ServiceItem `json:"serviceItems,omitempty" gorm:"-"` // ServiceItem currently uses Category string field
}

// TableName specifies the table name for GameCategory.
func (GameCategory) TableName() string {
	return "game_categories"
}

// BeforeCreate GORM hook called before creating a new GameCategory.
func (gc *GameCategory) BeforeCreate(tx *gorm.DB) error {
	// Set default sort order if not provided
	if gc.SortOrder == 0 {
		var maxSortOrder int
		tx.Model(&GameCategory{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSortOrder)
		gc.SortOrder = maxSortOrder + 1
	}
	return nil
}
