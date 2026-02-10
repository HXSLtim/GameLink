package common

import (
	"context"
	"sync"
	"time"

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
	db             *gorm.DB
	defaultTimeout time.Duration
	once           sync.Once
	root           *Repos // cached non-transactional repos
}

// NewUnitOfWork creates a UnitOfWork from the root *gorm.DB.
func NewUnitOfWork(db *gorm.DB) *UnitOfWork { return &UnitOfWork{db: db} }

// SetDefaultTimeout configures a default context timeout applied to every
// transactional call. A zero value disables the timeout.
func (u *UnitOfWork) SetDefaultTimeout(d time.Duration) { u.defaultTimeout = d }

// Repos returns repositories bound to the root (non-transactional) DB handle.
// The result is cached after the first call.
func (u *UnitOfWork) Repos() *Repos {
	u.once.Do(func() {
		u.root = buildRepos(u.db)
	})
	return u.root
}

// WithTx runs fn within a database transaction. If fn returns an error the
// transaction is rolled back; otherwise it is committed.
// When a default timeout is configured the context is wrapped with
// context.WithTimeout before entering the transaction.
func (u *UnitOfWork) WithTx(ctx context.Context, fn func(r *Repos) error) error {
	if u.defaultTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, u.defaultTimeout)
		defer cancel()
	}
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(buildRepos(tx))
	})
}

// buildRepos constructs a Repos struct bound to the given *gorm.DB handle.
func buildRepos(db *gorm.DB) *Repos {
	return &Repos{
		Games:          gameRepo.NewGameRepository(db),
		GameCategories: gamecategoryrepo.NewGameCategoryRepository(db),
		Users:          user.NewUserRepository(db),
		Players:        user.NewPlayerRepository(db),
		Orders:         orderrepo.NewOrderRepository(db),
		Payments:       order.NewPaymentRepository(db),
		Tags:           user.NewPlayerTagRepository(db),
		OpLogs:         adminrepo.NewOperationLogRepository(db),
		Reviews:        order.NewReviewRepository(db),
		ReviewReports:  reviewreport.NewReviewReportRepository(db),
		ReviewReplies:  reviewreply.NewReviewReplyRepository(db),
		Notifications:  notification.NewNotificationRepository(db),
	}
}

// MustGetGameCategoryRepo returns the GameCategoryRepository (for use in service methods)
func (r *Repos) MustGetGameCategoryRepo() repository.GameCategoryRepository {
	if r.GameCategories == nil {
		panic("GameCategories repository not initialized")
	}
	return r.GameCategories
}
