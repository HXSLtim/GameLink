package repository

import (
    "context"
    "time"

    "gamelink/internal/model"
)

// ReviewListOptions 评价列表筛选。
type ReviewListOptions struct {
    Page     int
    PageSize int
    OrderID  *uint64
    UserID   *uint64
    PlayerID *uint64
    DateFrom *time.Time
    DateTo   *time.Time
}

// ReviewRepository 评价仓储接口。
type ReviewRepository interface {
    List(ctx context.Context, opts ReviewListOptions) ([]model.Review, int64, error)
    Get(ctx context.Context, id uint64) (*model.Review, error)
    Create(ctx context.Context, r *model.Review) error
    Update(ctx context.Context, r *model.Review) error
    Delete(ctx context.Context, id uint64) error
}

