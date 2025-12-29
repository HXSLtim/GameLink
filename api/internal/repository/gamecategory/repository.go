/**
 * @file game category repository
 * @description 游戏分类数据访问层实现
 */

package gamecategory

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormGameCategoryRepository 使用 GORM 实现游戏分类管理
type gormGameCategoryRepository struct {
	db *gorm.DB
}

// NewGameCategoryRepository 创建 GORM 仓储实例
func NewGameCategoryRepository(db *gorm.DB) repository.GameCategoryRepository {
	return &gormGameCategoryRepository{db: db}
}

// Create 创建游戏分类
func (r *gormGameCategoryRepository) Create(ctx context.Context, category *model.GameCategory) error {
	// Check if category name already exists
	var existing model.GameCategory
	err := r.db.WithContext(ctx).Where("name = ?", category.Name).First(&existing).Error
	if err == nil {
		return errors.New("category name already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.WithContext(ctx).Create(category).Error
}

// Get 获取游戏分类详情
func (r *gormGameCategoryRepository) Get(ctx context.Context, id uint64) (*model.GameCategory, error) {
	var category model.GameCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

// GetByName 根据名称获取游戏分类
func (r *gormGameCategoryRepository) GetByName(ctx context.Context, name string) (*model.GameCategory, error) {
	var category model.GameCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

// List 获取游戏分类列表
func (r *gormGameCategoryRepository) List(ctx context.Context, opts repository.GameCategoryListOptions) ([]*model.GameCategory, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.GameCategory{})

	// 关键词搜索（匹配 name, description）
	if opts.Keyword != "" {
		searchPattern := "%" + opts.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// 启用状态过滤
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var categories []*model.GameCategory
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// Update 更新游戏分类
func (r *gormGameCategoryRepository) Update(ctx context.Context, category *model.GameCategory) error {
	// Check if category exists
	var existing model.GameCategory
	if err := r.db.WithContext(ctx).First(&existing, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	// If name is being changed, check for duplicates
	if category.Name != existing.Name {
		var duplicate model.GameCategory
		err := r.db.WithContext(ctx).Where("name = ? AND id != ?", category.Name, category.ID).First(&duplicate).Error
		if err == nil {
			return errors.New("category name already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	tx := r.db.WithContext(ctx).Model(category).Where("id = ?", category.ID).Updates(map[string]any{
		"name":        category.Name,
		"description": category.Description,
		"icon_url":    category.IconURL,
		"sort_order":  category.SortOrder,
		"is_active":   category.IsActive,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除游戏分类
func (r *gormGameCategoryRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.GameCategory{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// BatchDelete 批量删除游戏分类
func (r *gormGameCategoryRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	tx := r.db.WithContext(ctx).Delete(&model.GameCategory{}, ids)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// BatchUpdateStatus 批量更新游戏分类启用状态
func (r *gormGameCategoryRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	if len(ids) == 0 {
		return nil
	}
	tx := r.db.WithContext(ctx).Model(&model.GameCategory{}).Where("id IN ?", ids).Update("is_active", isActive)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Exists 检查分类是否存在
func (r *gormGameCategoryRepository) Exists(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.GameCategory{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountGames 统计分类下的游戏数量
// Note: Currently Game model uses a Category string field, not CategoryID.
// This method counts games where category field matches the category name.
func (r *gormGameCategoryRepository) CountGames(ctx context.Context, categoryID uint64) (int64, error) {
	// First get the category to obtain its name
	var category model.GameCategory
	if err := r.db.WithContext(ctx).First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, repository.ErrNotFound
		}
		return 0, err
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Game{}).Where("category = ?", category.Name).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountServiceItems 统计分类下的服务项目数量
// Note: Currently ServiceItem uses a Category string field, not CategoryID.
// This method counts service items where category field matches the category name.
func (r *gormGameCategoryRepository) CountServiceItems(ctx context.Context, categoryID uint64) (int64, error) {
	// First get the category to obtain its name
	var category model.GameCategory
	if err := r.db.WithContext(ctx).First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, repository.ErrNotFound
		}
		return 0, err
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ServiceItem{}).Where("category = ?", category.Name).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
