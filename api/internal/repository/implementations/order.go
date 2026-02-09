package implementations

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
)

// gormOrderRepository 使用 GORM 管理订单。
type gormOrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建实例。
func NewOrderRepository(db *gorm.DB) repoiface.OrderRepository {
	return &gormOrderRepository{db: db}
}

// Create inserts a new order.
func (r *gormOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// List returns a page of orders and the total count with filters applied.
func (r *gormOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
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
	// 使用 Preload 避免 N+1 查询问题
	if err := query.
		Preload("User").
		Preload("Player").
		Preload("Player.User").
		Preload("Game").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Get returns an order by id with related data preloaded.
func (r *gormOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Player").
		Preload("Player.User").
		Preload("Game").
		First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByIDs returns multiple orders by their IDs with related data preloaded.
// This method is optimized for batch operations to avoid N+1 query problems.
func (r *gormOrderRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Order, error) {
	if len(ids) == 0 {
		return []model.Order{}, nil
	}

	var orders []model.Order
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Player").
		Preload("Player.User").
		Preload("Game").
		Where("id IN ?", ids).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// Update updates editable fields of an order.
func (r *gormOrderRepository) Update(ctx context.Context, order *model.Order) error {
	tx := r.db.WithContext(ctx).Model(order).Where("id = ?", order.ID).Updates(map[string]any{
		"player_id":           order.PlayerID,
		"recipient_player_id": order.RecipientPlayerID,
		"game_id":             order.GameID,
		"status":              order.Status,
		"quantity":            order.Quantity,
		"unit_price_cents":    order.UnitPriceCents,
		"total_price_cents":   order.TotalPriceCents,
		"commission_cents":    order.CommissionCents,
		"player_income_cents": order.PlayerIncomeCents,
		"currency":            order.Currency,
		"title":               order.Title,
		"description":         order.Description,
		"scheduled_start":     order.ScheduledStart,
		"scheduled_end":       order.ScheduledEnd,
		"started_at":          order.StartedAt,
		"completed_at":        order.CompletedAt,
		"cancel_reason":       order.CancelReason,
		"refund_amount_cents": order.RefundAmountCents,
		"refund_reason":       order.RefundReason,
		"refunded_at":         order.RefundedAt,
		"gift_message":        order.GiftMessage,
		"is_anonymous":        order.IsAnonymous,
		"delivered_at":        order.DeliveredAt,
		"has_dispute":         order.HasDispute,
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
func (r *gormOrderRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Order{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateWithCondition 原子性更新订单,仅当状态匹配时才更新
// 使用数据库层面的WHERE条件确保原子性,避免并发竞态条件
//
// 参数:
//   - orderID: 订单ID
//   - expectedStatus: 期望的当前状态(仅当订单处于此状态时才更新)
//   - updates: 要更新的字段map
//
// 返回:
//   - bool: 是否成功更新(false表示状态不匹配,true表示更新成功)
//   - error: 数据库错误
func (r *gormOrderRepository) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	// 使用WHERE子句确保原子性: UPDATE ... WHERE id = ? AND status = ?
	// 这样即使多个goroutine同时执行,也只有一个能成功更新
	tx := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND status = ?", orderID, expectedStatus).
		Updates(updates)

	if tx.Error != nil {
		return false, tx.Error
	}

	// RowsAffected = 0 表示条件不满足(状态已变更或订单不存在)
	return tx.RowsAffected > 0, nil
}
