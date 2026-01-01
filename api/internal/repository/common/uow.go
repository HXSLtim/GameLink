package common

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/repository"
	adminrepo "gamelink/internal/repository/admin"
	notification "gamelink/internal/repository/content"
	gameRepo "gamelink/internal/repository/game"
	gamecategoryrepo "gamelink/internal/repository/gamecategory"
	orderrepo "gamelink/internal/repository/implementations"
	repoiface "gamelink/internal/repository/interfaces"
	order "gamelink/internal/repository/order"
	"gamelink/internal/repository/reviewreply"
	"gamelink/internal/repository/reviewreport"
	"gamelink/internal/repository/user"
)

// Repos bundles repository interfaces bound to a specific DB (tx) handle.
type Repos struct {
	Games          repository.GameRepository
	GameCategories repository.GameCategoryRepository
	Users          repository.UserRepository
	Players        repository.PlayerRepository
	Orders         repoiface.OrderRepository
	Payments       repository.PaymentRepository
	Tags           repository.PlayerTagRepository
	OpLogs         repository.OperationLogRepository
	Reviews        repository.ReviewRepository
	ReviewReports  repository.ReviewReportRepository
	ReviewReplies  repository.ReviewReplyRepository
	Notifications  repository.NotificationRepository
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
			Games:          gameRepo.NewGameRepository(tx),
			GameCategories: gamecategoryrepo.NewGameCategoryRepository(tx),
			Users:          user.NewUserRepository(tx),
			Players:        user.NewPlayerRepository(tx),
			Orders:         orderrepo.NewOrderRepository(tx),
			Payments:       order.NewPaymentRepository(tx),
			Tags:           user.NewPlayerTagRepository(tx),
			OpLogs:         adminrepo.NewOperationLogRepository(tx),
			Reviews:        order.NewReviewRepository(tx),
			ReviewReports:  reviewreport.NewReviewReportRepository(tx),
			ReviewReplies:  reviewreply.NewReviewReplyRepository(tx),
			Notifications:  notification.NewNotificationRepository(tx),
		}
		return fn(r)
	})
}

// MustGetGameCategoryRepo returns the GameCategoryRepository (for use in service methods)
func (r *Repos) MustGetGameCategoryRepo() repository.GameCategoryRepository {
	if r.GameCategories == nil {
		panic("GameCategories repository not initialized")
	}
	return r.GameCategories
}
