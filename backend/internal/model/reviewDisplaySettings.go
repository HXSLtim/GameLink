package model

import (
	"time"
)

// ReviewSortBy 评价排序方式
type ReviewSortBy string

const (
	ReviewSortByTime   ReviewSortBy = "time"    // 按时间排序
	ReviewSortByScore  ReviewSortBy = "score"   // 按评分排序
	ReviewSortByLikes  ReviewSortBy = "likes"   // 按点赞数排序
)

// Valid 检查排序方式是否合法
func (s ReviewSortBy) Valid() bool {
	switch s {
	case ReviewSortByTime, ReviewSortByScore, ReviewSortByLikes:
		return true
	default:
		return false
	}
}

// ReviewDisplaySettings 评价展示设置模型
// @Description 评价展示设置，用于配置前端评价的显示方式
type ReviewDisplaySettings struct {
	ID            uint64       `json:"id" gorm:"primaryKey"`
	// 排序方式：time/score/likes
	// @Enum time, score, likes
	// @Example time
	SortBy        ReviewSortBy `json:"sortBy" gorm:"column:sort_by;type:varchar(20);default:'time'"`
	// 最低评分阈值（1-5），低于此评分的评价不显示
	// @Example 1
	MinScore      int          `json:"minScore" gorm:"column:min_score;type:tinyint;default:1"`
	// 是否显示匿名评价
	// @Example true
	ShowAnonymous bool         `json:"showAnonymous" gorm:"column:show_anonymous;default:true"`
	// 每页显示数量
	// @Example 20
	PageSize      int          `json:"pageSize" gorm:"column:page_size;type:int;default:20"`
	// 创建时间
	CreatedAt     time.Time    `json:"createdAt" gorm:"column:created_at"`
	// 更新时间
	UpdatedAt     time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (ReviewDisplaySettings) TableName() string {
	return "review_display_settings"
}

// DefaultReviewDisplaySettings 返回默认的评价展示设置
func DefaultReviewDisplaySettings() *ReviewDisplaySettings {
	return &ReviewDisplaySettings{
		ID:            1, // 单例配置，固定ID为1
		SortBy:        ReviewSortByTime,
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      20,
	}
}

// Validate 验证设置是否合法
func (s *ReviewDisplaySettings) Validate() error {
	if !s.SortBy.Valid() {
		return ErrInvalidReviewSortBy
	}
	if s.MinScore < 1 || s.MinScore > 5 {
		return ErrInvalidMinScore
	}
	if s.PageSize < 1 || s.PageSize > 100 {
		return ErrInvalidPageSize
	}
	return nil
}

// ReviewDisplaySettingsErrors 评价展示设置相关错误
var (
	ErrInvalidReviewSortBy = &ValidationError{Field: "sortBy", Message: "invalid sort by value, must be one of: time, score, likes"}
	ErrInvalidMinScore     = &ValidationError{Field: "minScore", Message: "min score must be between 1 and 5"}
	ErrInvalidPageSize     = &ValidationError{Field: "pageSize", Message: "page size must be between 1 and 100"}
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
