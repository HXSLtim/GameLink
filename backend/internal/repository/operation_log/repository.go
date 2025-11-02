package operationlog

import (
    "context"

    "gorm.io/gorm"

    "gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormOperationLogRepository struct{ db *gorm.DB }

func NewOperationLogRepository(db *gorm.DB) repository.OperationLogRepository { return &gormOperationLogRepository{db: db} }

func (r *gormOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *gormOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
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
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    var rows []model.OperationLog
    if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil { return nil, 0, err }
    return rows, total, nil
}

