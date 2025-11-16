package interfaces

import (
	"context"

	"gamelink/internal/model"
)

// OrderQuery groups methods that scan or paginate orders.
type OrderQuery interface {
	List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
}
