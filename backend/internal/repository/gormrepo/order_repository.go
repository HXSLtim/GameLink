package gormrepo

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// OrderRepository 使用 GORM 管理订单。
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建实例。
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// List returns a page of orders and the total count with filters applied.
func (r *OrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Order{})

	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at <= ?", *opts.DateTo)
	}
	if trimmed := strings.TrimSpace(opts.Keyword); trimmed != "" {
		like := "%" + trimmed + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	var orders []model.Order
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Get returns an order by id.
func (r *OrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

// Update updates editable fields of an order.
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	tx := r.db.WithContext(ctx).Model(order).Where("id = ?", order.ID).Updates(map[string]any{
		"status":          order.Status,
		"price_cents":     order.PriceCents,
		"currency":        order.Currency,
		"scheduled_start": order.ScheduledStart,
		"scheduled_end":   order.ScheduledEnd,
		"cancel_reason":   order.CancelReason,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes an order by id.
func (r *OrderRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Order{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

var _ repository.OrderRepository = (*OrderRepository)(nil)
