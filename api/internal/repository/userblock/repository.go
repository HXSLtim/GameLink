package userblock

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 用户拉黑仓储实现
type Repository struct {
	db *gorm.DB
}

// NewUserBlockRepository 创建用户拉黑仓储
func NewUserBlockRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 创建拉黑记录
func (r *Repository) Create(ctx context.Context, block *model.UserBlock) error {
	return r.db.WithContext(ctx).Create(block).Error
}

// Get 获取拉黑记录
func (r *Repository) Get(ctx context.Context, id uint64) (*model.UserBlock, error) {
	var block model.UserBlock
	if err := r.db.WithContext(ctx).First(&block, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &block, nil
}

// GetWithRelations 获取拉黑记录（含关联用户信息）
func (r *Repository) GetWithRelations(ctx context.Context, id uint64) (*model.UserBlock, error) {
	var block model.UserBlock
	if err := r.db.WithContext(ctx).
		Preload("Blocker").
		Preload("Blocked").
		Preload("Canceler").
		First(&block, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &block, nil
}

// GetByBlockerAndBlocked 根据拉黑双方获取记录
func (r *Repository) GetByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error) {
	var block model.UserBlock
	if err := r.db.WithContext(ctx).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		First(&block).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &block, nil
}

// GetActiveByBlockerAndBlocked 获取生效中的拉黑记录
func (r *Repository) GetActiveByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error) {
	var block model.UserBlock
	if err := r.db.WithContext(ctx).
		Where("blocker_id = ? AND blocked_id = ? AND status = ?", blockerID, blockedID, model.BlockStatusActive).
		First(&block).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &block, nil
}

// IsBlocked 检查是否存在拉黑关系（任一方向）
func (r *Repository) IsBlocked(ctx context.Context, userID1, userID2 uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Where("status = ? AND ((blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?))",
			model.BlockStatusActive, userID1, userID2, userID2, userID1).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsBlockedBy 检查 blockedID 是否被 blockerID 拉黑
func (r *Repository) IsBlockedBy(ctx context.Context, blockerID, blockedID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Where("blocker_id = ? AND blocked_id = ? AND status = ?", blockerID, blockedID, model.BlockStatusActive).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListByBlockerID 获取用户拉黑的列表
func (r *Repository) ListByBlockerID(ctx context.Context, blockerID uint64, status *model.BlockStatus) ([]model.UserBlock, error) {
	query := r.db.WithContext(ctx).
		Preload("Blocked").
		Where("blocker_id = ?", blockerID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var blocks []model.UserBlock
	if err := query.Order("blocked_at DESC").Find(&blocks).Error; err != nil {
		return nil, err
	}
	return blocks, nil
}

// ListByBlockedID 获取被拉黑的列表
func (r *Repository) ListByBlockedID(ctx context.Context, blockedID uint64, status *model.BlockStatus) ([]model.UserBlock, error) {
	query := r.db.WithContext(ctx).
		Preload("Blocker").
		Where("blocked_id = ?", blockedID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var blocks []model.UserBlock
	if err := query.Order("blocked_at DESC").Find(&blocks).Error; err != nil {
		return nil, err
	}
	return blocks, nil
}

// ListPaged 分页获取拉黑记录
func (r *Repository) ListPaged(ctx context.Context, opts repository.UserBlockListOptions) ([]model.UserBlock, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.UserBlock{})

	// 筛选条件
	if opts.BlockerID != nil {
		query = query.Where("blocker_id = ?", *opts.BlockerID)
	}
	if opts.BlockedID != nil {
		query = query.Where("blocked_id = ?", *opts.BlockedID)
	}
	if opts.BlockerType != nil {
		query = query.Where("blocker_type = ?", *opts.BlockerType)
	}
	if opts.BlockedType != nil {
		query = query.Where("blocked_type = ?", *opts.BlockedType)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.DateFrom != nil {
		query = query.Where("blocked_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("blocked_at <= ?", *opts.DateTo)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var blocks []model.UserBlock
	offset := (opts.Page - 1) * opts.PageSize
	if err := query.
		Preload("Blocker").
		Preload("Blocked").
		Order("blocked_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&blocks).Error; err != nil {
		return nil, 0, err
	}

	return blocks, total, nil
}

// Update 更新拉黑记录
func (r *Repository) Update(ctx context.Context, block *model.UserBlock) error {
	return r.db.WithContext(ctx).Save(block).Error
}

// UpdateStatus 更新拉黑状态
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status model.BlockStatus, canceledBy *uint64, adminRemark string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == model.BlockStatusCanceled || status == model.BlockStatusAdminCanceled {
		now := time.Now()
		updates["canceled_at"] = &now
	}

	if canceledBy != nil {
		updates["canceled_by"] = canceledBy
	}

	if adminRemark != "" {
		updates["admin_remark"] = adminRemark
	}

	result := r.db.WithContext(ctx).Model(&model.UserBlock{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除拉黑记录（软删除）
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.UserBlock{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetBlockedUserIDs 获取用户拉黑的所有用户ID列表
func (r *Repository) GetBlockedUserIDs(ctx context.Context, blockerID uint64) ([]uint64, error) {
	var ids []uint64
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Where("blocker_id = ? AND status = ?", blockerID, model.BlockStatusActive).
		Pluck("blocked_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// GetBlockerUserIDs 获取拉黑该用户的所有用户ID列表
func (r *Repository) GetBlockerUserIDs(ctx context.Context, blockedID uint64) ([]uint64, error) {
	var ids []uint64
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Where("blocked_id = ? AND status = ?", blockedID, model.BlockStatusActive).
		Pluck("blocker_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// GetAllBlockRelatedUserIDs 获取与用户有拉黑关系的所有用户ID（双向）
func (r *Repository) GetAllBlockRelatedUserIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	var ids []uint64

	// 获取该用户拉黑的
	blockedIDs, err := r.GetBlockedUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取拉黑该用户的
	blockerIDs, err := r.GetBlockerUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 合并去重
	idSet := make(map[uint64]struct{})
	for _, id := range blockedIDs {
		idSet[id] = struct{}{}
	}
	for _, id := range blockerIDs {
		idSet[id] = struct{}{}
	}

	for id := range idSet {
		ids = append(ids, id)
	}

	return ids, nil
}

// CountByStatus 按状态统计拉黑记录数量
func (r *Repository) CountByStatus(ctx context.Context) (map[model.BlockStatus]int64, error) {
	type result struct {
		Status model.BlockStatus
		Count  int64
	}

	var results []result
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[model.BlockStatus]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}
	return stats, nil
}

// GetActiveCount 获取生效中的拉黑记录数量
func (r *Repository) GetActiveCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.UserBlock{}).
		Where("status = ?", model.BlockStatusActive).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
