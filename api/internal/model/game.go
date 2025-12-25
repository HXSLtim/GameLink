package model

// Game represents a supported game and its metadata.
// @Description 游戏配置模型，表示平台支持的游戏及其元数据
// @Example {"id": 1, "key": "lol", "name": "英雄联盟", "category": "MOBA", "iconUrl": "https://example.com/lol.png", "description": "热门MOBA游戏", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}
type Game struct {
	Base
	// 游戏唯一标识key，如 lol, dota2
	// @Example lol
	Key string `json:"key" gorm:"size:64;uniqueIndex"`
	// 游戏名称
	// @Example 英雄联盟
	Name string `json:"name" gorm:"size:128"`
	// 游戏分类，如 moba/fps
	// @Example MOBA
	Category string `json:"category,omitempty" gorm:"size:64;index"`
	// 游戏图标URL
	// @Example https://example.com/game-icon.png
	IconURL string `json:"iconUrl,omitempty" gorm:"column:icon_url;size:255"`
	// 游戏封面图URL
	// @Example https://example.com/game-cover.png
	CoverURL string `json:"coverUrl,omitempty" gorm:"column:cover_url;size:255"`
	// 游戏描述
	// @Example 英雄联盟是一款由Riot Games开发的MOBA类游戏
	Description string `json:"description,omitempty" gorm:"type:text"`
	// 是否上架
	// @Example true
	IsActive bool `json:"isActive" gorm:"column:is_active;default:true;index"`
	// 排序（越小越靠前）
	// @Example 0
	SortOrder int `json:"sortOrder" gorm:"column:sort_order;default:0;index"`
}
