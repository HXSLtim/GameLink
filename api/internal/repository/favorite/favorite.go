package favorite

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 收藏仓储
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建收藏仓储
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建收藏
func (r *Repository) Create(ctx context.Context, fav *model.Favorite) error {
	return r.db.WithContext(ctx).Create(fav).Error
}

// Delete 删除收藏
func (r *Repository) Delete(ctx context.Context, userID, playerID uint64) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND player_id = ?", userID, playerID).
		Delete(&model.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Exists 检查是否已收藏
func (r *Repository) Exists(ctx context.Context, userID, playerID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Favorite{}).
		Where("user_id = ? AND player_id = ?", userID, playerID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListByUserID 获取用户的收藏列表
func (r *Repository) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).
		Preload("Player").
		Order("created_at DESC").
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// GetByUserAndPlayer 获取指定收藏
func (r *Repository) GetByUserAndPlayer(ctx context.Context, userID, playerID uint64) (*model.Favorite, error) {
	var fav model.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND player_id = ?", userID, playerID).
		First(&fav).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &fav, nil
}
