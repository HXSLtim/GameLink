package model

// UserTag 用户标签
// 用于为用户打标分类，如 VIP、黑名单、测试用户等
type UserTag struct {
    Base
    Name        string `json:"name" gorm:"size:64;uniqueIndex;not null;comment:标签名称（唯一）"`
    Color       string `json:"color" gorm:"size:7;comment:颜色值，格式如#FF6B6B"`
    Description string `json:"description,omitempty" gorm:"type:text;comment:标签说明"`
    CreatedBy   uint64 `json:"createdBy" gorm:"column:created_by;not null;index;comment:创建人用户ID"`
}

// TableName 指定表名
func (UserTag) TableName() string { return "user_tags" }

