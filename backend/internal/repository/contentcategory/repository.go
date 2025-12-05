package contentcategory

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormContentCategoryRepository struct {
	db *gorm.DB
}

// NewContentCategoryRepository 创建内容分类仓储实例
func NewContentCategoryRepository(db *gorm.DB) repository.ContentCategoryRepository {
	return &gormContentCategoryRepository{db: db}
}

// Create 创建内容分类（带唯一性验证）
func (r *gormContentCategoryRepository) Create(ctx context.Context, category *model.ContentCategory) error {
	// 检查分类名称是否已存在
	var existing model.ContentCategory
	err := r.db.WithContext(ctx).Where("name = ?", category.Name).First(&existing).Error
	if err == nil {
		return errors.New("分类名称已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.WithContext(ctx).Create(category).Error
}

// Get 获取内容分类详情
func (r *gormContentCategoryRepository) Get(ctx context.Context, id uint64) (*model.ContentCategory, error) {
	var category model.ContentCategory
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

// GetByName 根据名称获取内容分类
func (r *gormContentCategoryRepository) GetByName(ctx context.Context, name string) (*model.ContentCategory, error) {
	var category model.ContentCategory
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

// List 列出内容分类（支持搜索和状态筛选）
func (r *gormContentCategoryRepository) List(ctx context.Context, opts repository.ContentCategoryListOptions) ([]model.ContentCategory, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.ContentCategory{})

	// 关键词搜索
	if opts.Keyword != "" {
		q = q.Where("name LIKE ? OR description LIKE ?", "%"+opts.Keyword+"%", "%"+opts.Keyword+"%")
	}

	// 状态筛选
	if opts.Status != nil {
		q = q.Where("status = ?", *opts.Status)
	}

	// 统计总数
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	var categories []model.ContentCategory
	if err := q.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(size).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// Update 更新内容分类
func (r *gormContentCategoryRepository) Update(ctx context.Context, category *model.ContentCategory) error {
	// 检查是否存在
	var existing model.ContentCategory
	if err := r.db.WithContext(ctx).First(&existing, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	// 如果修改了名称，检查新名称是否与其他分类重复
	if category.Name != existing.Name {
		var duplicate model.ContentCategory
		err := r.db.WithContext(ctx).Where("name = ? AND id != ?", category.Name, category.ID).First(&duplicate).Error
		if err == nil {
			return errors.New("分类名称已存在")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新分类
	tx := r.db.WithContext(ctx).Model(category).Where("id = ?", category.ID).Updates(map[string]any{
		"name":        category.Name,
		"description": category.Description,
		"sort_order":  category.SortOrder,
		"status":      category.Status,
		"icon_url":    category.IconURL,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除内容分类
func (r *gormContentCategoryRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.ContentCategory{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetFeedCount 获取分类下的动态数量
func (r *gormContentCategoryRepository) GetFeedCount(ctx context.Context, categoryID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Feed{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// MigrateFeeds 将分类下的动态迁移到另一个分类
func (r *gormContentCategoryRepository) MigrateFeeds(ctx context.Context, fromCategoryID, toCategoryID uint64) error {
	return r.db.WithContext(ctx).Model(&model.Feed{}).
		Where("category_id = ?", fromCategoryID).
		Update("category_id", toCategoryID).Error
}
