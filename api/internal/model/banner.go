package model

// BannerType banner 类型
type BannerType string

const (
	BannerTypeLink    BannerType = "link"    // 跳转链接
	BannerTypePreview BannerType = "preview" // 图片预览
)

// Banner 首页轮播图
type Banner struct {
	Base
	Title       string     `json:"title" gorm:"size:128"`                                  // 标题（PC 端 Hero 区域用）
	Description string     `json:"description" gorm:"size:256"`                            // 描述（PC 端 Hero 区域用）
	ImageURL    string     `json:"imageUrl" gorm:"column:image_url;size:512;not null"`     // 图片地址
	Type        BannerType `json:"type" gorm:"size:32;default:'link'"`                     // 类型
	Link        string     `json:"link,omitempty" gorm:"size:512"`                         // 跳转链接（type=link 时）
	ActionText  string     `json:"actionText,omitempty" gorm:"column:action_text;size:64"` // 按钮文字（PC 端 Hero 区域用）
	SortOrder   int        `json:"sortOrder" gorm:"column:sort_order;default:0"`           // 排序（越小越靠前）
	IsVisible   bool       `json:"isVisible" gorm:"column:is_visible;default:true"`        // 是否展示
}

// TableName 指定表名
func (Banner) TableName() string {
	return "banners"
}
