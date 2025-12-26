package model

// UserTagRelation 用户与标签关联
// 记录用户与标签的多对多关系
type UserTagRelation struct {
	Base
	UserID uint64 `json:"userId" gorm:"column:user_id;index:idx_user_tag,priority:1;not null;comment:用户ID"`
	TagID  uint64 `json:"tagId" gorm:"column:tag_id;index:idx_user_tag,priority:2;not null;comment:标签ID"`
	IsAuto bool   `json:"isAuto" gorm:"column:is_auto;default:false;comment:是否自动标签（基于统计指标）"`
}

// TableName 指定表名
func (UserTagRelation) TableName() string { return "user_tag_relations" }
