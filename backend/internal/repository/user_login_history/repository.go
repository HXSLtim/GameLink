package user_login_history

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormUserLoginHistoryRepository struct {
	db *gorm.DB
}

// NewUserLoginHistoryRepository 创建登录历史仓储实例
func NewUserLoginHistoryRepository(db *gorm.DB) repository.UserLoginHistoryRepository {
	return &gormUserLoginHistoryRepository{db: db}
}

// Create 创建登录历史记录
func (r *gormUserLoginHistoryRepository) Create(ctx context.Context, history *model.UserLoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetByUserID 获取用户登录历史列表（分页）
func (r *gormUserLoginHistoryRepository) GetByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserLoginHistory, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	var total int64
	query := r.db.WithContext(ctx).Model(&model.UserLoginHistory{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var histories []model.UserLoginHistory
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// GetByUserIDAndDate 获取用户指定时间范围内的登录历史
func (r *gormUserLoginHistoryRepository) GetByUserIDAndDate(ctx context.Context, userID uint64, dateFrom, dateTo time.Time) ([]model.UserLoginHistory, error) {
	var histories []model.UserLoginHistory

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, dateFrom, dateTo).
		Order("created_at DESC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}
