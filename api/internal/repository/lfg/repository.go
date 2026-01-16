package lfg

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormLFGRepository implements LFGRequestRepository.
type gormLFGRepository struct {
	db *gorm.DB
}

// NewLFGRepository creates a new LFG repository.
func NewLFGRepository(db *gorm.DB) repository.LFGRequestRepository {
	return &gormLFGRepository{db: db}
}

// Create inserts a new LFG request.
func (r *gormLFGRepository) Create(ctx context.Context, request *model.LFGRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// Get returns a request by id.
func (r *gormLFGRepository) Get(ctx context.Context, id uint64) (*model.LFGRequest, error) {
	var request model.LFGRequest
	if err := r.db.WithContext(ctx).First(&request, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &request, nil
}

// GetWithRelations returns a request with all relations preloaded.
func (r *gormLFGRepository) GetWithRelations(ctx context.Context, id uint64) (*model.LFGRequest, error) {
	var request model.LFGRequest
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Game").
		Preload("MatchedRoom").
		First(&request, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &request, nil
}

// GetActiveByUserID returns the active (pending) request for a user.
func (r *gormLFGRepository) GetActiveByUserID(ctx context.Context, userID uint64) (*model.LFGRequest, error) {
	var request model.LFGRequest
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.LFGPending).
		First(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &request, nil
}

// Update updates a request.
func (r *gormLFGRepository) Update(ctx context.Context, request *model.LFGRequest) error {
	tx := r.db.WithContext(ctx).Model(request).Updates(map[string]any{
		"request_type":     request.RequestType,
		"title":            request.Title,
		"description":      request.Description,
		"required_players": request.RequiredPlayers,
		"min_rank":         request.MinRank,
		"max_price_cents":  request.MaxPriceCents,
		"status":           request.Status,
		"expires_at":       request.ExpiresAt,
		"matched_room_id":  request.MatchedRoomID,
		"matched_at":       request.MatchedAt,
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
func (r *gormLFGRepository) UpdateStatus(ctx context.Context, id uint64, status model.LFGRequestStatus) error {
	tx := r.db.WithContext(ctx).Model(&model.LFGRequest{}).Where("id = ?", id).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateMatched updates the matched room and status.
func (r *gormLFGRepository) UpdateMatched(ctx context.Context, id uint64, roomID uint64) error {
	now := time.Now()
	tx := r.db.WithContext(ctx).Model(&model.LFGRequest{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":          model.LFGMatched,
			"matched_room_id": roomID,
			"matched_at":      &now,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a request.
func (r *gormLFGRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.LFGRequest{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// List returns requests with pagination and filters.
func (r *gormLFGRepository) List(ctx context.Context, opts repository.LFGRequestListOptions) ([]model.LFGRequest, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.LFGRequest{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}
	if opts.RequestType != nil {
		query = query.Where("request_type = ?", *opts.RequestType)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	if opts.MinRank != "" {
		query = query.Where("min_rank = ?", opts.MinRank)
	}
	if opts.MaxPrice != nil {
		query = query.Where("max_price_cents <= ?", *opts.MaxPrice)
	}
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at <= ?", *opts.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var requests []model.LFGRequest
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Game").
		Where(query).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

// ListByUserID returns requests by user id.
func (r *gormLFGRepository) ListByUserID(ctx context.Context, userID uint64, status *model.LFGRequestStatus) ([]model.LFGRequest, error) {
	query := r.db.WithContext(ctx).Preload("Game").Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var requests []model.LFGRequest
	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// ListByGameID returns requests by game id with pagination.
func (r *gormLFGRepository) ListByGameID(ctx context.Context, gameID uint64, status *model.LFGRequestStatus, page, pageSize int) ([]model.LFGRequest, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.LFGRequest{}).Where("game_id = ?", gameID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var requests []model.LFGRequest
	findQuery := r.db.WithContext(ctx).Preload("User").Preload("Game").Where("game_id = ?", gameID)
	if status != nil {
		findQuery = findQuery.Where("status = ?", *status)
	}
	if err := findQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

// ListPending returns pending requests with pagination.
func (r *gormLFGRepository) ListPending(ctx context.Context, gameID *uint64, page, pageSize int) ([]model.LFGRequest, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.LFGRequest{}).
		Where("status = ?", model.LFGPending).
		Where("expires_at > ?", time.Now())

	if gameID != nil {
		query = query.Where("game_id = ?", *gameID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var requests []model.LFGRequest
	findQuery := r.db.WithContext(ctx).
		Preload("User").
		Preload("Game").
		Where("status = ?", model.LFGPending).
		Where("expires_at > ?", time.Now())

	if gameID != nil {
		findQuery = findQuery.Where("game_id = ?", *gameID)
	}
	if err := findQuery.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

// ListExpired returns expired pending requests.
func (r *gormLFGRepository) ListExpired(ctx context.Context, limit int) ([]model.LFGRequest, error) {
	var requests []model.LFGRequest
	if err := r.db.WithContext(ctx).
		Where("status = ?", model.LFGPending).
		Where("expires_at <= ?", time.Now()).
		Limit(limit).
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// BatchExpire marks multiple requests as expired.
func (r *gormLFGRepository) BatchExpire(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&model.LFGRequest{}).
		Where("id IN ?", ids).
		Update("status", model.LFGExpired).Error
}

// FindMatchingRequests finds requests that could match with the given request.
func (r *gormLFGRepository) FindMatchingRequests(ctx context.Context, request *model.LFGRequest, limit int) ([]model.LFGRequest, error) {
	query := r.db.WithContext(ctx).
		Preload("User").
		Where("id != ?", request.ID).
		Where("game_id = ?", request.GameID).
		Where("status = ?", model.LFGPending).
		Where("expires_at > ?", time.Now())

	// Match opposite request types
	if request.RequestType == model.LFGFindPlayer {
		query = query.Where("request_type = ?", model.LFGFindTeam)
	} else {
		query = query.Where("request_type = ?", model.LFGFindPlayer)
	}

	// Price matching (if specified)
	if request.MaxPriceCents > 0 {
		query = query.Where("max_price_cents <= ? OR max_price_cents = 0", request.MaxPriceCents)
	}

	var requests []model.LFGRequest
	if err := query.Order("created_at ASC").Limit(limit).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// CountByStatus returns count of requests grouped by status.
func (r *gormLFGRepository) CountByStatus(ctx context.Context) (map[model.LFGRequestStatus]int64, error) {
	type result struct {
		Status model.LFGRequestStatus
		Count  int64
	}
	var results []result
	if err := r.db.WithContext(ctx).Model(&model.LFGRequest{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.LFGRequestStatus]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// CountPending returns count of pending requests.
func (r *gormLFGRepository) CountPending(ctx context.Context, gameID *uint64) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.LFGRequest{}).
		Where("status = ?", model.LFGPending).
		Where("expires_at > ?", time.Now())

	if gameID != nil {
		query = query.Where("game_id = ?", *gameID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
