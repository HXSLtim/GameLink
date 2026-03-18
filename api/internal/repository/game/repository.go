package game

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/cache"
)

const (
	cacheKeyGamesList = "games:list:all"
	cacheTTLGames     = 1 * time.Hour
)

// gormGameRepository 使用 GORM 实现游戏管理。
type gormGameRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGameRepository 创建 GORM 仓储实例。
func NewGameRepository(db *gorm.DB) repository.GameRepository {
	return &gormGameRepository{db: db, cache: nil}
}

// NewGameRepositoryWithCache 创建带缓存的 GORM 仓储实例。
func NewGameRepositoryWithCache(db *gorm.DB, cache cache.Cache) repository.GameRepository {
	return &gormGameRepository{db: db, cache: cache}
}

// List returns all games ordered by creation time.
func (r *gormGameRepository) List(ctx context.Context) ([]model.Game, error) {
	// Try cache first
	if r.cache != nil {
		if cached, ok, _ := r.cache.Get(ctx, cacheKeyGamesList); ok {
			var games []model.Game
			if err := json.Unmarshal([]byte(cached), &games); err == nil {
				return games, nil
			}
		}
	}

	// Cache miss, query database
	var games []model.Game
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&games).Error; err != nil {
		return nil, err
	}

	// Cache the result
	r.cacheGamesList(ctx, games)
	return games, nil
}

// ListPaged returns a page of games and the total count.
func (r *gormGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
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

// ListPagedWithFilter returns a page of games with keyword filter.
func (r *gormGameRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Game{})

	// 关键词搜索（匹配 name, key, category, description）
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name ILIKE ? OR key ILIKE ? OR category ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

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
func (r *gormGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
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

// GetByIDs returns games by a list of IDs.
func (r *gormGameRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Game, error) {
	if len(ids) == 0 {
		return []model.Game{}, nil
	}
	var games []model.Game
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

// Create inserts a new game.
func (r *gormGameRepository) Create(ctx context.Context, game *model.Game) error {
	if err := r.db.WithContext(ctx).Create(game).Error; err != nil {
		return err
	}
	// Invalidate cache
	r.invalidateCache(ctx)
	return nil
}

// Update updates editable fields of a game.
func (r *gormGameRepository) Update(ctx context.Context, game *model.Game) error {
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
	// Invalidate cache
	r.invalidateCache(ctx)
	return nil
}

// Delete soft-deletes a game by id.
func (r *gormGameRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Game{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	// Invalidate cache
	r.invalidateCache(ctx)
	return nil
}

// BatchDelete soft-deletes multiple games by ids.
func (r *gormGameRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Delete(&model.Game{}, ids)
	if tx.Error != nil {
		return 0, tx.Error
	}
	// Invalidate cache
	r.invalidateCache(ctx)
	return tx.RowsAffected, nil
}

// BatchUpdateStatus updates is_active for multiple games.
func (r *gormGameRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Model(&model.Game{}).Where("id IN ?", ids).Update("is_active", isActive)
	if tx.Error != nil {
		return 0, tx.Error
	}
	// Invalidate cache
	r.invalidateCache(ctx)
	return tx.RowsAffected, nil
}

// BatchUpdateSortOrder updates sort_order for multiple games.
func (r *gormGameRepository) BatchUpdateSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	var totalAffected int64
	for id, sortOrder := range updates {
		tx := r.db.WithContext(ctx).Model(&model.Game{}).Where("id = ?", id).Update("sort_order", sortOrder)
		if tx.Error != nil {
			return totalAffected, tx.Error
		}
		totalAffected += tx.RowsAffected
	}

	// Invalidate cache
	r.invalidateCache(ctx)
	return totalAffected, nil
}

// BatchUpdateCategory updates category for multiple games.
func (r *gormGameRepository) BatchUpdateCategory(ctx context.Context, ids []uint64, category string) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Model(&model.Game{}).Where("id IN ?", ids).Update("category", category)
	if tx.Error != nil {
		return 0, tx.Error
	}
	// Invalidate cache
	r.invalidateCache(ctx)
	return tx.RowsAffected, nil
}

// cacheGamesList caches the games list.
func (r *gormGameRepository) cacheGamesList(ctx context.Context, games []model.Game) {
	if r.cache == nil {
		return
	}
	data, err := json.Marshal(games)
	if err != nil {
		return
	}
	_ = r.cache.Set(ctx, cacheKeyGamesList, string(data), cacheTTLGames)
}

// invalidateCache clears the games list cache.
func (r *gormGameRepository) invalidateCache(ctx context.Context) {
	if r.cache == nil {
		return
	}
	_ = r.cache.Delete(ctx, cacheKeyGamesList)
}
