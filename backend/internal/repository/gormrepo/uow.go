package gormrepo

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/repository"
)

// Repos bundles repository interfaces bound to a specific DB (tx) handle.
type Repos struct {
    Games    repository.GameRepository
    Users    repository.UserRepository
    Players  repository.PlayerRepository
    Orders   repository.OrderRepository
    Payments repository.PaymentRepository
    Tags     repository.PlayerTagRepository
    OpLogs   repository.OperationLogRepository
    Reviews  repository.ReviewRepository
}

// UnitOfWork provides a simple transaction wrapper for GORM repositories.
type UnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork creates a UnitOfWork from the root *gorm.DB.
func NewUnitOfWork(db *gorm.DB) *UnitOfWork { return &UnitOfWork{db: db} }

// WithTx runs fn within a database transaction. If fn returns an error the
// transaction is rolled back; otherwise it is committed.
func (u *UnitOfWork) WithTx(ctx context.Context, fn func(r *Repos) error) error {
    return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        r := &Repos{
            Games:    NewGameRepository(tx),
            Users:    NewUserRepository(tx),
            Players:  NewPlayerRepository(tx),
            Orders:   NewOrderRepository(tx),
            Payments: NewPaymentRepository(tx),
            Tags:     NewPlayerTagRepository(tx),
            OpLogs:   NewOperationLogRepository(tx),
            Reviews:  NewReviewRepository(tx),
        }
        return fn(r)
    })
}
