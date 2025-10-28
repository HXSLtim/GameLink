package repository

import (
    "context"

    "gamelink/internal/model"
)

// OperationLogRepository 管理操作日志。
type OperationLogRepository interface {
    Append(ctx context.Context, log *model.OperationLog) error
    // ListByEntity 按实体分页返回日志与总数。
    ListByEntity(ctx context.Context, entityType string, entityID uint64, page, pageSize int) ([]model.OperationLog, int64, error)
}

