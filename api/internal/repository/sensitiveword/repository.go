package sensitiveword

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormSensitiveWordRepository struct {
	db *gorm.DB
}

// NewSensitiveWordRepository 创建敏感词仓储实例
func NewSensitiveWordRepository(db *gorm.DB) repository.SensitiveWordRepository {
	return &gormSensitiveWordRepository{db: db}
}

// Create 创建敏感词（带唯一性验证）
func (r *gormSensitiveWordRepository) Create(ctx context.Context, word *model.SensitiveWord) error {
	// 检查敏感词是否已存在
	var existing model.SensitiveWord
	err := r.db.WithContext(ctx).Where("word = ?", word.Word).First(&existing).Error
	if err == nil {
		return errors.New("敏感词已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建敏感词
	return r.db.WithContext(ctx).Create(word).Error
}

// Get 获取敏感词详情
func (r *gormSensitiveWordRepository) Get(ctx context.Context, id uint64) (*model.SensitiveWord, error) {
	var word model.SensitiveWord
	err := r.db.WithContext(ctx).First(&word, id).Error
	return repository.HandleGetError(&word, err)
}

// List 列出敏感词（支持搜索和分类筛选）
func (r *gormSensitiveWordRepository) List(ctx context.Context, opts repository.SensitiveWordListOptions) ([]model.SensitiveWord, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.SensitiveWord{})

	// 关键词搜索
	if opts.Keyword != "" {
		q = q.Where("word LIKE ?", "%"+opts.Keyword+"%")
	}

	// 分类筛选
	if opts.Category != nil {
		q = q.Where("category = ?", *opts.Category)
	}

	// 严重程度筛选
	if opts.Severity != nil {
		q = q.Where("severity = ?", *opts.Severity)
	}

	// 统计总数
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	var words []model.SensitiveWord
	if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&words).Error; err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

// Update 更新敏感词
func (r *gormSensitiveWordRepository) Update(ctx context.Context, word *model.SensitiveWord) error {
	// 检查是否存在
	var existing model.SensitiveWord
	if err := r.db.WithContext(ctx).First(&existing, word.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	// 如果修改了敏感词内容，检查新内容是否与其他敏感词重复
	if word.Word != existing.Word {
		var duplicate model.SensitiveWord
		err := r.db.WithContext(ctx).Where("word = ? AND id != ?", word.Word, word.ID).First(&duplicate).Error
		if err == nil {
			return errors.New("敏感词已存在")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新敏感词
	tx := r.db.WithContext(ctx).Model(word).Where("id = ?", word.ID).Updates(map[string]any{
		"word":     word.Word,
		"category": word.Category,
		"severity": word.Severity,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除敏感词
func (r *gormSensitiveWordRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.SensitiveWord{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetAll 获取所有敏感词（用于检测）
func (r *gormSensitiveWordRepository) GetAll(ctx context.Context) ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	if err := r.db.WithContext(ctx).Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}
