package vip

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository VIP仓储实现
type Repository struct {
	db *gorm.DB
}

// NewVipRepository 创建VIP仓储
func NewVipRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// VIP等级管理
// ============================================================================

// CreateLevel 创建VIP等级
func (r *Repository) CreateLevel(ctx context.Context, level *model.VipLevel) error {
	return r.db.WithContext(ctx).Create(level).Error
}

// GetLevel 获取VIP等级
func (r *Repository) GetLevel(ctx context.Context, id uint64) (*model.VipLevel, error) {
	var level model.VipLevel
	if err := r.db.WithContext(ctx).First(&level, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &level, nil
}

// GetLevelBySlug 根据Slug获取VIP等级
func (r *Repository) GetLevelBySlug(ctx context.Context, slug string) (*model.VipLevel, error) {
	var level model.VipLevel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &level, nil
}

// GetDefaultLevel 获取默认VIP等级
func (r *Repository) GetDefaultLevel(ctx context.Context) (*model.VipLevel, error) {
	var level model.VipLevel
	if err := r.db.WithContext(ctx).Where("is_default = ? AND is_active = ?", true, true).First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &level, nil
}

// ListLevels 获取所有VIP等级
func (r *Repository) ListLevels(ctx context.Context) ([]model.VipLevel, error) {
	var levels []model.VipLevel
	if err := r.db.WithContext(ctx).Order("sort_order ASC, exp_required ASC").Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

// ListActiveLevels 获取所有启用的VIP等级
func (r *Repository) ListActiveLevels(ctx context.Context) ([]model.VipLevel, error) {
	var levels []model.VipLevel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC, exp_required ASC").Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

// ListLevelsPaged 分页获取VIP等级
func (r *Repository) ListLevelsPaged(ctx context.Context, opts repository.VipLevelListOptions) ([]model.VipLevel, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.VipLevel{})

	// 关键词搜索
	if opts.Keyword != "" {
		likePattern := "%" + opts.Keyword + "%"
		query = query.Where("title ILIKE ? OR slug ILIKE ?", likePattern, likePattern)
	}

	// 是否启用筛选
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var levels []model.VipLevel
	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Order("sort_order ASC, exp_required ASC").Offset(offset).Limit(opts.PageSize).Find(&levels).Error; err != nil {
		return nil, 0, err
	}

	return levels, total, nil
}

// UpdateLevel 更新VIP等级
func (r *Repository) UpdateLevel(ctx context.Context, level *model.VipLevel) error {
	result := r.db.WithContext(ctx).Model(level).Updates(map[string]any{
		"slug":                       level.Slug,
		"title":                      level.Title,
		"exp_required":               level.ExpRequired,
		"order_discount":             level.OrderDiscount,
		"monthly_coupon_template_id": level.MonthlyCouponTemplateID,
		"monthly_coupon_count":       level.MonthlyCouponCount,
		"icon_url":                   level.IconURL,
		"color":                      level.Color,
		"benefits":                   level.Benefits,
		"sort_order":                 level.SortOrder,
		"is_default":                 level.IsDefault,
		"is_active":                  level.IsActive,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteLevel 删除VIP等级
func (r *Repository) DeleteLevel(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.VipLevel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// SetDefaultLevel 设置默认VIP等级（清除其他默认）
func (r *Repository) SetDefaultLevel(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 清除所有默认
		if err := tx.Model(&model.VipLevel{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新默认
		result := tx.Model(&model.VipLevel{}).Where("id = ?", id).Update("is_default", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return repository.ErrNotFound
		}
		return nil
	})
}

// GetLevelByExp 根据经验值获取对应的VIP等级
func (r *Repository) GetLevelByExp(ctx context.Context, exp int64) (*model.VipLevel, error) {
	var level model.VipLevel
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND exp_required <= ?", true, exp).
		Order("exp_required DESC").
		First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &level, nil
}

// BatchUpdateLevelStatus 批量更新VIP等级状态
func (r *Repository) BatchUpdateLevelStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.WithContext(ctx).Model(&model.VipLevel{}).Where("id IN ?", ids).Update("is_active", isActive)
	return result.RowsAffected, result.Error
}

// BatchDeleteLevels 批量删除VIP等级
func (r *Repository) BatchDeleteLevels(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.WithContext(ctx).Delete(&model.VipLevel{}, ids)
	return result.RowsAffected, result.Error
}

// ============================================================================
// VIP配置管理
// ============================================================================

// GetConfig 获取VIP配置
func (r *Repository) GetConfig(ctx context.Context, key string) (*model.VipConfig, error) {
	var config model.VipConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &config, nil
}

// ListConfigs 获取所有VIP配置
func (r *Repository) ListConfigs(ctx context.Context) ([]model.VipConfig, error) {
	var configs []model.VipConfig
	if err := r.db.WithContext(ctx).Order("config_key ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// SaveConfig 保存VIP配置（创建或更新）
func (r *Repository) SaveConfig(ctx context.Context, config *model.VipConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// DeleteConfig 删除VIP配置
func (r *Repository) DeleteConfig(ctx context.Context, key string) error {
	result := r.db.WithContext(ctx).Where("config_key = ?", key).Delete(&model.VipConfig{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}
