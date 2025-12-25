package model

// ContentCategoryStatus 内容分类状态
type ContentCategoryStatus string

const (
	ContentCategoryStatusActive   ContentCategoryStatus = "active"   // 启用
	ContentCategoryStatusInactive ContentCategoryStatus = "inactive" // 禁用
)

// Valid 检查内容分类状态是否合法
func (s ContentCategoryStatus) Valid() bool {
	switch s {
	case ContentCategoryStatusActive, ContentCategoryStatusInactive:
		return true
	default:
		return false
	}
}

// ContentCategory 内容分类模型
// @Description 内容分类模型，用于组织和管理动态内容
type ContentCategory struct {
	Base
	// 分类名称
	// @Example 游戏攻略
	Name string `json:"name" gorm:"column:name;type:varchar(64);uniqueIndex;not null"`
	// 分类描述
	// @Example 分享游戏攻略和技巧
	Description string `json:"description,omitempty" gorm:"column:description;type:text"`
	// 排序权重，数值越小越靠前
	// @Example 1
	SortOrder int `json:"sortOrder" gorm:"column:sort_order;default:0;index"`
	// 分类状态
	// @Enum active, inactive
	// @Example active
	Status ContentCategoryStatus `json:"status" gorm:"column:status;type:varchar(16);default:'active';index"`
	// 分类图标URL
	// @Example https://example.com/icon.png
	IconURL string `json:"iconUrl,omitempty" gorm:"column:icon_url;type:varchar(255)"`
	// 该分类下的动态数量（非持久化字段）
	FeedCount int64 `json:"feedCount,omitempty" gorm:"-"`
}

// TableName 指定表名
func (ContentCategory) TableName() string { return "content_categories" }
