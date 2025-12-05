package admin

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type operationLogRepository struct{ db *gorm.DB }

func NewOperationLogRepository(db *gorm.DB) repository.OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize
	q := r.db.WithContext(ctx).Model(&model.OperationLog{}).Where("entity_type = ? AND entity_id = ?", entityType, entityID)
	if opts.Action != "" {
		q = q.Where("action = ?", opts.Action)
	}
	if opts.ActorUserID != nil {
		q = q.Where("actor_user_id = ?", *opts.ActorUserID)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.OperationLog
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *operationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize
	q := r.db.WithContext(ctx).Model(&model.OperationLog{})
	
	if opts.EntityType != "" {
		q = q.Where("entity_type = ?", opts.EntityType)
	}
	if opts.EntityID != nil {
		q = q.Where("entity_id = ?", *opts.EntityID)
	}
	if opts.Action != "" {
		q = q.Where("action = ?", opts.Action)
	}
	if opts.ActorUserID != nil {
		q = q.Where("actor_user_id = ?", *opts.ActorUserID)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}
	
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.OperationLog
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}