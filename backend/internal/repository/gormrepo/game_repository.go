package gormrepo

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// GameRepository 使用 GORM 实现游戏管理。
type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository 创建 GORM 仓储实例。
func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// List returns all games ordered by creation time.
func (r *GameRepository) List(ctx context.Context) ([]model.Game, error) {
	var games []model.Game
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

// ListPaged returns a page of games and the total count.
func (r *GameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Game{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var games []model.Game
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&games).Error; err != nil {
		return nil, 0, err
	}
	return games, total, nil
}

// Get returns a game by id.
func (r *GameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	var game model.Game
	err := r.db.WithContext(ctx).First(&game, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &game, nil
}

// Create inserts a new game.
func (r *GameRepository) Create(ctx context.Context, game *model.Game) error {
	return r.db.WithContext(ctx).Create(game).Error
}

// Update updates editable fields of a game.
func (r *GameRepository) Update(ctx context.Context, game *model.Game) error {
	tx := r.db.WithContext(ctx).Model(game).Updates(map[string]any{
		"key":         game.Key,
		"name":        game.Name,
		"category":    game.Category,
		"icon_url":    game.IconURL,
		"description": game.Description,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a game by id.
func (r *GameRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Game{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

var _ repository.GameRepository = (*GameRepository)(nil)
