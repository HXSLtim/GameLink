package scheduler

import "context"

// CommissionService defines the interface for commission settlement operations
type CommissionService interface {
	SettleMonth(ctx context.Context, month string) error
}