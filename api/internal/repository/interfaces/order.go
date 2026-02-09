package interfaces

import (
	"context"
	"time"

	"gamelink/internal/model"
)

// OrderReader 只负责读取单个订单
type OrderReader interface {
	Get(ctx context.Context, id uint64) (*model.Order, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]model.Order, error)
}

// OrderWriter 只负责写入/删除订单
type OrderWriter interface {
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uint64) error
	// UpdateWithCondition 原子性更新订单,仅当满足条件时才更新
	// 返回是否成功更新(如果条件不满足则返回false,nil)
	UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error)
}

// OrderQuery 只负责按条件查询订单列表
type OrderQuery interface {
	List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
}

// OrderReadWriter 组合读写接口
type OrderReadWriter interface {
	OrderReader
	OrderWriter
}

// OrderRepository 聚合读写与查询能力
// 这是仓储对外使用的主要接口
type OrderRepository interface {
	OrderReadWriter
	OrderQuery
}

// OrderListOptions 订单列表查询选项
// 字段命名和含义需与现有调用方保持一致
type OrderListOptions struct {
	Page     int
	PageSize int

	UserID   *uint64
	PlayerID *uint64
	GameID   *uint64

	Statuses []model.OrderStatus

	DateFrom *time.Time
	DateTo   *time.Time

	Keyword string
}
