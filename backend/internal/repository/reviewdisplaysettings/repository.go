package reviewdisplaysettings

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 评价展示设置仓储实现
type Repository struct {
	db *gorm.DB
}

// New 创建评价展示设置仓储实例
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Get 获取当前设置（单例模式，只有一条记录）
// 如果不存在则返回默认设置
func (r *Repository) Get(ctx context.Context) (*model.ReviewDisplaySettings, error) {
	var settings model.ReviewDisplaySettings
	err := r.db.WithContext(ctx).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回默认设置
			return model.DefaultReviewDisplaySettings(), nil
		}
		return nil, err
	}
	return &settings, nil
}

// Save 保存设置（创建或更新）
// 使用 upsert 模式，确保只有一条记录
func (r *Repository) Save(ctx context.Context, settings *model.ReviewDisplaySettings) error {
	// 验证设置
	if err := settings.Validate(); err != nil {
		return err
	}

	// 强制使用ID=1，确保单例模式
	settings.ID = 1

	// 使用 upsert 模式
	return r.db.WithContext(ctx).Save(settings).Error
}

// Ensure Repository implements the interface
var _ repository.ReviewDisplaySettingsRepository = (*Repository)(nil)
