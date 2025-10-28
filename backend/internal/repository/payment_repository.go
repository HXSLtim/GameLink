package repository

import (
	"context"
	"time"

	"gamelink/internal/model"
)

// PaymentListOptions 定义支付查询参数。
type PaymentListOptions struct {
	Page     int
	PageSize int
	Statuses []model.PaymentStatus
	Methods  []model.PaymentMethod
	UserID   *uint64
	OrderID  *uint64
	DateFrom *time.Time
	DateTo   *time.Time
}

// PaymentRepository 管理支付记录。
type PaymentRepository interface {
    Create(ctx context.Context, payment *model.Payment) error
    List(ctx context.Context, opts PaymentListOptions) ([]model.Payment, int64, error)
    Get(ctx context.Context, id uint64) (*model.Payment, error)
    Update(ctx context.Context, payment *model.Payment) error
    Delete(ctx context.Context, id uint64) error
}
