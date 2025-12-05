package userbehavior

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormUserBehaviorRepository struct {
	db *gorm.DB
}

// NewUserBehaviorRepository 创建用户行为仓储实例
func NewUserBehaviorRepository(db *gorm.DB) repository.UserBehaviorRepository {
	return &gormUserBehaviorRepository{db: db}
}

// Create 创建用户行为记录
func (r *gormUserBehaviorRepository) Create(ctx context.Context, behavior *model.UserBehavior) error {
	return r.db.WithContext(ctx).Create(behavior).Error
}

// GetUserBehaviors 获取用户行为列表（分页）
func (r *gormUserBehaviorRepository) GetUserBehaviors(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserBehavior, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	var total int64
	query := r.db.WithContext(ctx).Model(&model.UserBehavior{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var behaviors []model.UserBehavior
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&behaviors).Error; err != nil {
		return nil, 0, err
	}

	return behaviors, total, nil
}

// GetBehaviorStats 获取行为统计数据
// days: 统计最近N天的数据
// 返回：行为类型 -> 次数的映射
func (r *gormUserBehaviorRepository) GetBehaviorStats(ctx context.Context, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 7
	}

	since := time.Now().AddDate(0, 0, -days+1)

	type ActionCount struct {
		Action string
		Count  int64
	}

	var results []ActionCount
	err := r.db.WithContext(ctx).
		Model(&model.UserBehavior{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("action").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Action] = r.Count
	}

	return stats, nil
}
