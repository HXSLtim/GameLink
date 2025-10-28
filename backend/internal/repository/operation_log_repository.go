package repository

import (
    "context"
    "time"

    "gamelink/internal/model"
)

// OperationLogRepository 管理操作日志。
type OperationLogRepository interface {
    Append(ctx context.Context, log *model.OperationLog) error
    // ListByEntity 按实体分页返回日志与总数，可选过滤。
    ListByEntity(ctx context.Context, entityType string, entityID uint64, opts OperationLogListOptions) ([]model.OperationLog, int64, error)
}

// OperationLogListOptions 日志查询过滤。
type OperationLogListOptions struct {
    Page       int
    PageSize   int
    Action     string   // 单一动作过滤（可扩展为多值）
    ActorUserID *uint64
    DateFrom   *time.Time
    DateTo     *time.Time
}
