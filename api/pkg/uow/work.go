package uow

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Work 工作单元，管理事务和仓储
type Work struct {
	db    *gorm.DB
	repos map[string]interface{}
	ctx   context.Context
}

// NewWork 创建工作单元
func NewWork(ctx context.Context, db *gorm.DB) *Work {
	return &Work{
		db:    db,
		repos: make(map[string]interface{}),
		ctx:   ctx,
	}
}

// DB 返回底层数据库连接
func (w *Work) DB() *gorm.DB {
	return w.db
}

// Context 返回上下文
func (w *Work) Context() context.Context {
	return w.ctx
}

// Commit 在事务中执行操作
//
// 使用示例:
//
//	err := uow.NewWork(ctx, db).Commit(func(w *uow.Work) error {
//	    orderRepo := w.OrderRepo()
//	    paymentRepo := w.PaymentRepo()
//	    // 所有操作在同一事务中
//	    return nil
//	})
func (w *Work) Commit(fn func(*Work) error) error {
	return w.db.WithContext(w.ctx).Transaction(func(tx *gorm.DB) error {
		w.db = tx
		return fn(w)
	})
}

// Execute 执行操作（不在事务中，用于只读操作）
func (w *Work) Execute(fn func(*Work) error) error {
	return fn(w)
}

// Rollback 手动回滚事务（通常不需要，Commit 会自动处理）
func (w *Work) Rollback() error {
	// GORM 事务会自动回滚，这里不需要手动实现
	// 保留接口以便扩展
	return nil
}

// GetRepository 获取仓储（通用方法）
func (w *Work) GetRepository(name string, factory func(*gorm.DB) interface{}) interface{} {
	if r, ok := w.repos[name]; ok {
		return r
	}

	repo := factory(w.db)
	w.repos[name] = repo
	return repo
}

// 以下是常用仓储的获取方法
// 这些方法需要在具体的 repository 实现后才能使用

// OrderRepo 获取订单仓储
// func (w *Work) OrderRepo() OrderRepository {
//     return w.GetRepository("order", func(db *gorm.DB) interface{} {
//         return repository.NewOrderRepository(db)
//     }).(OrderRepository)
// }

// PaymentRepo 获取支付仓储
// func (w *Work) PaymentRepo() PaymentRepository {
//     return w.GetRepository("payment", func(db *gorm.DB) interface{} {
//         return repository.NewPaymentRepository(db)
//     }).(PaymentRepository)
// }

// UserRepo 获取用户仓储
// func (w *Work) UserRepo() UserRepository {
//     return w.GetRepository("user", func(db *gorm.DB) interface{} {
//         return repository.NewUserRepository(db)
//     }).(UserRepository)
// }

// PlayerRepo 获取陪玩师仓储
// func (w *Work) PlayerRepo() PlayerRepository {
//     return w.GetRepository("player", func(db *gorm.DB) interface{} {
//         return repository.NewPlayerRepository(db)
//     }).(PlayerRepository)
// }

// CommissionRepo 获取佣金仓储
// func (w *Work) CommissionRepo() CommissionRepository {
//     return w.GetRepository("commission", func(db *gorm.DB) interface{} {
//         return repository.NewCommissionRepository(db)
//     }).(CommissionRepository)
// }

// TransactionScope 在事务作用域内执行函数（辅助方法）
func TransactionScope(ctx context.Context, db *gorm.DB, fn func(context.Context) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

// InTransaction 检查上下文是否在事务中
func InTransaction(ctx context.Context) bool {
	return ctx.Value("tx") != nil
}

// GetTx 从上下文获取事务（如果在事务中）
func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	return tx, ok
}

// Transactional 事务装饰器，用于确保函数在事务中执行
func Transactional(db *gorm.DB) func(next func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			ctxWithTx := context.WithValue(ctx, "tx", tx)
			return next(ctxWithTx)
		})
	}
}

// WithTransaction 在事务中执行操作（简化版）
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(context.Context, *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

// WorkFunc 工作单元函数类型
type WorkFunc func(*Work) error

// ExecuteInUnitOfWork 在工作单元中执行操作
func ExecuteInUnitOfWork(ctx context.Context, db *gorm.DB, fn WorkFunc) error {
	work := NewWork(ctx, db)
	return work.Commit(fn)
}

// RepositoryNotFoundError 仓储未找到错误
type RepositoryNotFoundError struct {
	Name string
}

func (e *RepositoryNotFoundError) Error() string {
	return fmt.Sprintf("repository not found: %s", e.Name)
}

// IsRepositoryNotFound 检查错误是否为仓储未找到
func IsRepositoryNotFound(err error) bool {
	_, ok := err.(*RepositoryNotFoundError)
	return ok
}
