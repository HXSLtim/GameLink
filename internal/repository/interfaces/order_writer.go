package interfaces

import (
	"context"

	"gamelink/internal/model"
)

// OrderWriter encapsulates write operations that mutate orders.
type OrderWriter interface {
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uint64) error
}
