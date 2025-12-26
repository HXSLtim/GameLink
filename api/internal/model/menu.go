package model

// Menu represents a front-end route/menu entry.
type Menu struct {
	Base
	Name        string  `json:"name" gorm:"size:128;not null"`
	Path        string  `json:"path" gorm:"size:255;not null;index"`
	Component   string  `json:"component" gorm:"size:255"`
	Icon        string  `json:"icon" gorm:"size:64"`
	ParentID    *uint64 `json:"parentId" gorm:"index"`
	Order       int     `json:"order" gorm:"default:0"`
	Hidden      bool    `json:"hidden" gorm:"default:false"`
	Permission  string  `json:"permission" gorm:"size:128"` // optional permission code required
	Redirect    string  `json:"redirect" gorm:"size:255"`
	Description string  `json:"description" gorm:"size:255"`
	Children    []Menu  `json:"children,omitempty" gorm:"-"` // 非数据库字段，用于构建树形结构
}

func (Menu) TableName() string { return "menus" }
