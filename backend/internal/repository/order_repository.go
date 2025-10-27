package repository

import (
	"context"
	"time"

	"gamelink/internal/model"
)

// OrderListOptions 定义订单查询与分页参数。
type OrderListOptions struct {
	Page     int
	PageSize int
	Statuses []model.OrderStatus
	UserID   *uint64
	PlayerID *uint64
	GameID   *uint64
	DateFrom *time.Time
	DateTo   *time.Time
	Keyword  string
}

// OrderRepository 管理订单生命周期。
type OrderRepository interface {
	List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
	Get(ctx context.Context, id uint64) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uint64) error
}
