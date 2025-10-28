package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type OperationLogRepository struct{ db *gorm.DB }

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository { return &OperationLogRepository{db: db} }

func (r *OperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *OperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, page, pageSize int) ([]model.OperationLog, int64, error) {
    page = repository.NormalizePage(page)
    pageSize = repository.NormalizePageSize(pageSize)
    offset := (page - 1) * pageSize
    q := r.db.WithContext(ctx).Model(&model.OperationLog{}).Where("entity_type = ? AND entity_id = ?", entityType, entityID)
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    var rows []model.OperationLog
    if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil { return nil, 0, err }
    return rows, total, nil
}

var _ repository.OperationLogRepository = (*OperationLogRepository)(nil)

