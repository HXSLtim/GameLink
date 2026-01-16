package presence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormPresenceRepository implements PlayerPresenceRepository.
type gormPresenceRepository struct {
	db *gorm.DB
}

// NewPresenceRepository creates a new presence repository.
func NewPresenceRepository(db *gorm.DB) repository.PlayerPresenceRepository {
	return &gormPresenceRepository{db: db}
}

// Create inserts a new player presence record.
func (r *gormPresenceRepository) Create(ctx context.Context, presence *model.PlayerPresence) error {
	return r.db.WithContext(ctx).Create(presence).Error
}

// Get returns a presence by id.
func (r *gormPresenceRepository) Get(ctx context.Context, id uint64) (*model.PlayerPresence, error) {
	var presence model.PlayerPresence
	if err := r.db.WithContext(ctx).First(&presence, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &presence, nil
}

// GetByPlayerID returns presence by player id.
func (r *gormPresenceRepository) GetByPlayerID(ctx context.Context, playerID uint64) (*model.PlayerPresence, error) {
	var presence model.PlayerPresence
	if err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&presence).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &presence, nil
}

// GetWithPlayer returns presence with player relation preloaded.
func (r *gormPresenceRepository) GetWithPlayer(ctx context.Context, id uint64) (*model.PlayerPresence, error) {
	var presence model.PlayerPresence
	if err := r.db.WithContext(ctx).Preload("Player").Preload("Player.User").First(&presence, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &presence, nil
}

// Update updates a presence record.
func (r *gormPresenceRepository) Update(ctx context.Context, presence *model.PlayerPresence) error {
	tx := r.db.WithContext(ctx).Model(presence).Updates(map[string]any{
		"status":            presence.Status,
		"current_game_id":   presence.CurrentGameID,
		"current_game_name": presence.CurrentGameName,
		"custom_status":     presence.CustomStatus,
		"current_order_id":  presence.CurrentOrderID,
		"current_room_id":   presence.CurrentRoomID,
		"last_heartbeat_at": presence.LastHeartbeatAt,
		"device_type":       presence.DeviceType,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateStatus updates only the status field.
func (r *gormPresenceRepository) UpdateStatus(ctx context.Context, playerID uint64, status model.PlayerPresenceStatus) error {
	tx := r.db.WithContext(ctx).Model(&model.PlayerPresence{}).
		Where("player_id = ?", playerID).
		Updates(map[string]any{
			"status":            status,
			"last_heartbeat_at": time.Now(),
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateHeartbeat updates the last heartbeat timestamp.
func (r *gormPresenceRepository) UpdateHeartbeat(ctx context.Context, playerID uint64) error {
	tx := r.db.WithContext(ctx).Model(&model.PlayerPresence{}).
		Where("player_id = ?", playerID).
		Update("last_heartbeat_at", time.Now())
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a presence record.
func (r *gormPresenceRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.PlayerPresence{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ListByPlayerIDs returns presences for multiple players.
func (r *gormPresenceRepository) ListByPlayerIDs(ctx context.Context, playerIDs []uint64) ([]model.PlayerPresence, error) {
	if len(playerIDs) == 0 {
		return []model.PlayerPresence{}, nil
	}
	var presences []model.PlayerPresence
	if err := r.db.WithContext(ctx).
		Preload("Player").
		Where("player_id IN ?", playerIDs).
		Find(&presences).Error; err != nil {
		return nil, err
	}
	return presences, nil
}

// ListOnline returns online presences with pagination.
func (r *gormPresenceRepository) ListOnline(ctx context.Context, opts repository.PlayerPresenceListOptions) ([]model.PlayerPresence, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.PlayerPresence{})

	// Filter by statuses (default to online statuses)
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	} else {
		// Default: exclude offline and invisible
		query = query.Where("status NOT IN ?", []model.PlayerPresenceStatus{
			model.PresenceOffline,
			model.PresenceInvisible,
		})
	}

	if opts.GameID != nil {
		query = query.Where("current_game_id = ?", *opts.GameID)
	}
	if opts.DeviceType != "" {
		query = query.Where("device_type = ?", opts.DeviceType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var presences []model.PlayerPresence
	if err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Player.User").
		Where("status NOT IN ?", []model.PlayerPresenceStatus{
			model.PresenceOffline,
			model.PresenceInvisible,
		}).
		Order("last_heartbeat_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&presences).Error; err != nil {
		return nil, 0, err
	}
	return presences, total, nil
}

// ListByStatus returns presences by status with pagination.
func (r *gormPresenceRepository) ListByStatus(ctx context.Context, status model.PlayerPresenceStatus, page, pageSize int) ([]model.PlayerPresence, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.PlayerPresence{}).Where("status = ?", status)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var presences []model.PlayerPresence
	if err := r.db.WithContext(ctx).
		Preload("Player").
		Where("status = ?", status).
		Order("last_heartbeat_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&presences).Error; err != nil {
		return nil, 0, err
	}
	return presences, total, nil
}

// ListStalePresences returns presences with heartbeat older than threshold.
func (r *gormPresenceRepository) ListStalePresences(ctx context.Context, threshold time.Time) ([]model.PlayerPresence, error) {
	var presences []model.PlayerPresence
	if err := r.db.WithContext(ctx).
		Where("last_heartbeat_at < ?", threshold).
		Where("status NOT IN ?", []model.PlayerPresenceStatus{
			model.PresenceOffline,
			model.PresenceInvisible,
		}).
		Find(&presences).Error; err != nil {
		return nil, err
	}
	return presences, nil
}

// BatchUpdateOffline marks multiple players as offline.
func (r *gormPresenceRepository) BatchUpdateOffline(ctx context.Context, playerIDs []uint64) error {
	if len(playerIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&model.PlayerPresence{}).
		Where("player_id IN ?", playerIDs).
		Update("status", model.PresenceOffline).Error
}

// CountByStatus returns count of presences grouped by status.
func (r *gormPresenceRepository) CountByStatus(ctx context.Context) (map[model.PlayerPresenceStatus]int64, error) {
	type result struct {
		Status model.PlayerPresenceStatus
		Count  int64
	}
	var results []result
	if err := r.db.WithContext(ctx).Model(&model.PlayerPresence{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.PlayerPresenceStatus]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// CountOnline returns count of online presences.
func (r *gormPresenceRepository) CountOnline(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.PlayerPresence{}).
		Where("status NOT IN ?", []model.PlayerPresenceStatus{
			model.PresenceOffline,
			model.PresenceInvisible,
		}).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
