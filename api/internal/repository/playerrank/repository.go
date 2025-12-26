package playerrank

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormPlayerRankRepository 陪玩师段位认证仓储实现
type gormPlayerRankRepository struct {
	db *gorm.DB
}

// NewPlayerRankRepository 创建陪玩师段位认证仓储
func NewPlayerRankRepository(db *gorm.DB) repository.PlayerRankRepository {
	return &gormPlayerRankRepository{db: db}
}

// Create 创建段位认证记录
func (r *gormPlayerRankRepository) Create(ctx context.Context, record *model.PlayerRankRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// Get 根据ID获取段位认证记录
func (r *gormPlayerRankRepository) Get(ctx context.Context, id uint64) (*model.PlayerRankRecord, error) {
	var record model.PlayerRankRecord
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// GetWithRelations 获取段位认证记录及关联信息
func (r *gormPlayerRankRepository) GetWithRelations(ctx context.Context, id uint64) (*model.PlayerRankRecord, error) {
	var record model.PlayerRankRecord
	if err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Game").
		Preload("Rank").
		Preload("Verifier").
		First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// GetByPlayerAndGame 根据陪玩师ID和游戏ID获取认证记录
func (r *gormPlayerRankRepository) GetByPlayerAndGame(ctx context.Context, playerID, gameID uint64) (*model.PlayerRankRecord, error) {
	var record model.PlayerRankRecord
	if err := r.db.WithContext(ctx).
		Where("player_id = ? AND game_id = ?", playerID, gameID).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// ListByPlayerID 根据陪玩师ID获取所有认证记录
func (r *gormPlayerRankRepository) ListByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerRankRecord, error) {
	var records []model.PlayerRankRecord
	if err := r.db.WithContext(ctx).
		Preload("Game").
		Preload("Rank").
		Where("player_id = ?", playerID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// ListPaged 分页获取段位认证记录
func (r *gormPlayerRankRepository) ListPaged(ctx context.Context, opts repository.PlayerRankListOptions) ([]model.PlayerRankRecord, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.PlayerRankRecord{})

	// 陪玩师ID筛选
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}

	// 游戏ID筛选
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	// 多状态筛选
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []model.PlayerRankRecord
	if err := query.
		Preload("Player").
		Preload("Game").
		Preload("Rank").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// ListPending 获取待审核的段位认证记录
func (r *gormPlayerRankRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerRankRecord, int64, error) {
	return r.ListPaged(ctx, repository.PlayerRankListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   ptrPlayerRankStatus(model.PlayerRankStatusPending),
	})
}

// Update 更新段位认证记录
func (r *gormPlayerRankRepository) Update(ctx context.Context, record *model.PlayerRankRecord) error {
	tx := r.db.WithContext(ctx).Model(record).Updates(map[string]any{
		"rank_id":         record.RankID,
		"status":          record.Status,
		"screenshot_urls": record.ScreenshotURLs,
		"verified_at":     record.VerifiedAt,
		"verified_by":     record.VerifiedBy,
		"reject_reason":   record.RejectReason,
		"expire_at":       record.ExpireAt,
		"remark":          record.Remark,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateStatus 更新认证状态
func (r *gormPlayerRankRepository) UpdateStatus(ctx context.Context, id uint64, status model.PlayerRankStatus, verifiedBy *uint64, rejectReason string) error {
	updates := map[string]any{
		"status":        status,
		"reject_reason": rejectReason,
	}
	if verifiedBy != nil {
		updates["verified_by"] = *verifiedBy
	}
	if status == model.PlayerRankStatusVerified {
		updates["verified_at"] = gorm.Expr("NOW()")
	}

	tx := r.db.WithContext(ctx).Model(&model.PlayerRankRecord{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除段位认证记录
func (r *gormPlayerRankRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.PlayerRankRecord{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// CountByStatus 统计各状态的认证数量
func (r *gormPlayerRankRepository) CountByStatus(ctx context.Context) (map[model.PlayerRankStatus]int64, error) {
	type result struct {
		Status model.PlayerRankStatus
		Count  int64
	}
	var results []result
	if err := r.db.WithContext(ctx).
		Model(&model.PlayerRankRecord{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.PlayerRankStatus]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// GetPendingCount 获取待审核数量
func (r *gormPlayerRankRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.PlayerRankRecord{}).
		Where("status = ?", model.PlayerRankStatusPending).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ptrPlayerRankStatus 返回状态指针
func ptrPlayerRankStatus(s model.PlayerRankStatus) *model.PlayerRankStatus {
	return &s
}
