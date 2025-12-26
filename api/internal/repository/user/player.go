package user

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// PlayerRepository å®ç°éªç©èµæä»å¨ã?
type gormPlayerRepository struct {
	db *gorm.DB
}

// NewPlayerRepository åå»ºéªç©ä»å¨ã?
func NewPlayerRepository(db *gorm.DB) repository.PlayerRepository {
	return &gormPlayerRepository{db: db}
}

// List returns all players ordered by creation time.
func (r *gormPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	var players []model.Player
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

// ListPaged returns a page of players and the total count.
func (r *gormPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Player{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var players []model.Player
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&players).Error; err != nil {
		return nil, 0, err
	}
	return players, total, nil
}

// Get returns a player by id.
func (r *gormPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	if err := r.db.WithContext(ctx).First(&player, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &player, nil
}

// GetByUserID returns player by bound user id.
func (r *gormPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	var player model.Player
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&player).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &player, nil
}

// Create inserts a new player.
func (r *gormPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

// Update updates editable fields of a player.
func (r *gormPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	tx := r.db.WithContext(ctx).Model(player).Updates(map[string]any{
		"nickname":            player.Nickname,
		"bio":                 player.Bio,
		"rating_average":      player.RatingAverage,
		"rating_count":        player.RatingCount,
		"hourly_rate_cents":   player.HourlyRateCents,
		"main_game_id":        player.MainGameID,
		"verification_status": player.VerificationStatus,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a player by id.
func (r *gormPlayerRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Player{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ListPagedWithFilter returns a page of players with keyword and status filter.
func (r *gormPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Player{})

	// Apply keyword filter (search in nickname and bio)
	if keyword != "" {
		likePattern := "%" + keyword + "%"
		query = query.Where("nickname ILIKE ? OR bio ILIKE ?", likePattern, likePattern)
	}

	// Apply status filter
	if status != nil {
		query = query.Where("verification_status = ?", *status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var players []model.Player
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&players).Error; err != nil {
		return nil, 0, err
	}
	return players, total, nil
}

// BatchUpdateStatus updates verification status for multiple players.
func (r *gormPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Model(&model.Player{}).Where("id IN ?", ids).Update("verification_status", status)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// BatchDelete soft-deletes multiple players by ids.
func (r *gormPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := r.db.WithContext(ctx).Delete(&model.Player{}, ids)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
