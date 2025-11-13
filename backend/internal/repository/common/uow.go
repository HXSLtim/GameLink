package common

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/repository"
	"gamelink/internal/repository/game"
	operationlog "gamelink/internal/repository/operation_log"
        "gamelink/internal/repository/order"
        orderdispute "gamelink/internal/repository/orderdispute"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	playertag "gamelink/internal/repository/player_tag"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
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
        Disputes repository.OrderDisputeRepository
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
			Games:    game.NewGameRepository(tx),
			Users:    user.NewUserRepository(tx),
			Players:  player.NewPlayerRepository(tx),
			Orders:   order.NewOrderRepository(tx),
			Payments: payment.NewPaymentRepository(tx),
			Tags:     playertag.NewPlayerTagRepository(tx),
			OpLogs:   operationlog.NewOperationLogRepository(tx),
                        Reviews:  review.NewReviewRepository(tx),
                        Disputes: orderdispute.NewRepository(tx),
                }
		return fn(r)
	})
}
