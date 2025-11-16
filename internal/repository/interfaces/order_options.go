package interfaces

import (
	"time"

	"gamelink/internal/model"
)

// OrderListOptions contains filtering options for order queries.
type OrderListOptions struct {
	Page     int
	PageSize int
	UserID   *uint64
	PlayerID *uint64
	GameID   *uint64
	Statuses []model.OrderStatus
	Keyword  string
	DateFrom *time.Time
	DateTo   *time.Time
}
