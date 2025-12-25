package ranking

import (
	"context"
	"errors"

	"gamelink/internal/model"

	"gorm.io/gorm"
)

var (
	// ErrNotFound 资源不存在错误
	ErrNotFound = errors.New("resource not found")
)

// RankingCommissionRepository 排名抽成配置仓储
type RankingCommissionRepository interface {
	// 配置管理
	CreateConfig(ctx context.Context, config *model.RankingCommissionConfig) error
	GetConfig(ctx context.Context, id uint64) (*model.RankingCommissionConfig, error)
	GetActiveConfigForMonth(ctx context.Context, rankingType model.RankingType, month string) (*model.RankingCommissionConfig, error)
	ListConfigs(ctx context.Context, opts RankingCommissionConfigListOptions) ([]model.RankingCommissionConfig, int64, error)
	UpdateConfig(ctx context.Context, config *model.RankingCommissionConfig) error
	DeleteConfig(ctx context.Context, id uint64) error
}

// RankingCommissionConfigListOptions 查询选项
type RankingCommissionConfigListOptions struct {
	RankingType *model.RankingType
	Month       *string
	IsActive    *bool
	Page        int
	PageSize    int
}

type rankingCommissionRepository struct {
	db *gorm.DB
}

// NewRankingCommissionRepository 创建排名抽成配置仓储
func NewRankingCommissionRepository(db *gorm.DB) RankingCommissionRepository {
	return &rankingCommissionRepository{db: db}
}

// CreateConfig 创建配置
func (r *rankingCommissionRepository) CreateConfig(ctx context.Context, config *model.RankingCommissionConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// GetConfig 获取配置
func (r *rankingCommissionRepository) GetConfig(ctx context.Context, id uint64) (*model.RankingCommissionConfig, error) {
	var config model.RankingCommissionConfig
	err := r.db.WithContext(ctx).First(&config, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &config, nil
}

// GetActiveConfigForMonth 获取指定月份的激活配置
func (r *rankingCommissionRepository) GetActiveConfigForMonth(ctx context.Context, rankingType model.RankingType, month string) (*model.RankingCommissionConfig, error) {
	var config model.RankingCommissionConfig
	err := r.db.WithContext(ctx).
		Where("ranking_type = ? AND month = ? AND is_active = ?", rankingType, month, true).
		First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &config, nil
}

// ListConfigs 查询配置列表
func (r *rankingCommissionRepository) ListConfigs(ctx context.Context, opts RankingCommissionConfigListOptions) ([]model.RankingCommissionConfig, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RankingCommissionConfig{})

	if opts.RankingType != nil {
		query = query.Where("ranking_type = ?", *opts.RankingType)
	}
	if opts.Month != nil {
		query = query.Where("month = ?", *opts.Month)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	var configs []model.RankingCommissionConfig
	err := query.Order("month DESC, created_at DESC").
		Offset(offset).Limit(opts.PageSize).Find(&configs).Error
	if err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// UpdateConfig 更新配置
func (r *rankingCommissionRepository) UpdateConfig(ctx context.Context, config *model.RankingCommissionConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// DeleteConfig 删除配置
func (r *rankingCommissionRepository) DeleteConfig(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.RankingCommissionConfig{}, id).Error
}
