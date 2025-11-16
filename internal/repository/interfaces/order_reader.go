package interfaces

import (
	"context"

	"gamelink/internal/model"
)

// OrderReader exposes read-only operations for retrieving single orders.
type OrderReader interface {
	Get(ctx context.Context, id uint64) (*model.Order, error)
}
