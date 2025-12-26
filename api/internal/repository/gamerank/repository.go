package gamerank

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormGameRankRepository 游戏段位配置仓储实现
type gormGameRankRepository struct {
	db *gorm.DB
}

// NewGameRankRepository 创建游戏段位配置仓储
func NewGameRankRepository(db *gorm.DB) repository.GameRankRepository {
	return &gormGameRankRepository{db: db}
}

// Create 创建游戏段位
func (r *gormGameRankRepository) Create(ctx context.Context, rank *model.GameRank) error {
	return r.db.WithContext(ctx).Create(rank).Error
}

// Get 根据ID获取游戏段位
func (r *gormGameRankRepository) Get(ctx context.Context, id uint64) (*model.GameRank, error) {
	var rank model.GameRank
	if err := r.db.WithContext(ctx).First(&rank, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rank, nil
}

// GetWithGame 获取段位及关联的游戏信息
func (r *gormGameRankRepository) GetWithGame(ctx context.Context, id uint64) (*model.GameRank, error) {
	var rank model.GameRank
	if err := r.db.WithContext(ctx).Preload("Game").First(&rank, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rank, nil
}

// List 获取所有游戏段位
func (r *gormGameRankRepository) List(ctx context.Context) ([]model.GameRank, error) {
	var ranks []model.GameRank
	if err := r.db.WithContext(ctx).Order("game_id ASC, level ASC").Find(&ranks).Error; err != nil {
		return nil, err
	}
	return ranks, nil
}

// ListByGameID 根据游戏ID获取段位列表
func (r *gormGameRankRepository) ListByGameID(ctx context.Context, gameID uint64) ([]model.GameRank, error) {
	var ranks []model.GameRank
	if err := r.db.WithContext(ctx).Where("game_id = ?", gameID).Order("level ASC").Find(&ranks).Error; err != nil {
		return nil, err
	}
	return ranks, nil
}

// ListPaged 分页获取游戏段位
func (r *gormGameRankRepository) ListPaged(ctx context.Context, opts repository.GameRankListOptions) ([]model.GameRank, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.GameRank{})

	// 游戏ID筛选
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}

	// 关键词搜索
	if opts.Keyword != "" {
		likePattern := "%" + opts.Keyword + "%"
		query = query.Where("name ILIKE ?", likePattern)
	}

	// 是否启用筛选
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var ranks []model.GameRank
	if err := query.Preload("Game").Order("game_id ASC, level ASC").Offset(offset).Limit(pageSize).Find(&ranks).Error; err != nil {
		return nil, 0, err
	}

	return ranks, total, nil
}

// Update 更新游戏段位
func (r *gormGameRankRepository) Update(ctx context.Context, rank *model.GameRank) error {
	tx := r.db.WithContext(ctx).Model(rank).Updates(map[string]any{
		"game_id":     rank.GameID,
		"name":        rank.Name,
		"level":       rank.Level,
		"price_cents": rank.PriceCents,
		"icon_url":    rank.IconURL,
		"color":       rank.Color,
		"description": rank.Description,
		"sort_order":  rank.SortOrder,
		"is_active":   rank.IsActive,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除游戏段位
func (r *gormGameRankRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.GameRank{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// BatchDelete 批量删除游戏段位
func (r *gormGameRankRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Delete(&model.GameRank{}, ids)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// BatchUpdateStatus 批量更新启用状态
func (r *gormGameRankRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Model(&model.GameRank{}).Where("id IN ?", ids).Update("is_active", isActive)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
