package payment

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// PaymentRepository ä½¿ç¨ GORM ç®¡çæ¯ä»è®°å½ã?
type gormPaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository åå»ºå®ä¾ã?
func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &gormPaymentRepository{db: db}
}

// Create inserts a new payment row.
func (r *gormPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// List returns a page of payments and the total count with filters applied.
func (r *gormPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Payment{})

	// Status filters
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	// Method filters
	if len(opts.Methods) > 0 {
		query = query.Where("method IN ?", opts.Methods)
	}
	// User filter
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	// Order filter
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	// Date range filters
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at <= ?", *opts.DateTo)
	}
	// Collection entity filter
	if opts.CollectionEntityID != nil {
		query = query.Where("collection_entity_id = ?", *opts.CollectionEntityID)
	}
	// Merchant number filter
	if opts.MerchantNo != "" {
		query = query.Where("merchant_no = ?", opts.MerchantNo)
	}
	// Provider trade number filter (partial match)
	if opts.ProviderTradeNo != "" {
		query = query.Where("provider_trade_no LIKE ?", "%"+opts.ProviderTradeNo+"%")
	}
	// Amount range filters
	if opts.MinAmountCents != nil {
		query = query.Where("amount_cents >= ?", *opts.MinAmountCents)
	}
	if opts.MaxAmountCents != nil {
		query = query.Where("amount_cents <= ?", *opts.MaxAmountCents)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	var payments []model.Payment
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// Get returns a payment by id.
func (r *gormPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).First(&payment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// Update updates editable fields of a payment.
func (r *gormPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	tx := r.db.WithContext(ctx).Model(payment).Where("id = ?", payment.ID).Updates(map[string]any{
		"status":                payment.Status,
		"provider_trade_no":     payment.ProviderTradeNo,
		"provider_raw":          payment.ProviderRaw,
		"paid_at":               payment.PaidAt,
		"refunded_at":           payment.RefundedAt,
		"refunded_amount_cents": payment.RefundedAmountCents,
		"collection_entity_id":  payment.CollectionEntityID,
		"merchant_no":           payment.MerchantNo,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a payment by id.
func (r *gormPaymentRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Payment{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetWithRelations returns a payment by id with preloaded Order and User relations.
func (r *gormPaymentRepository) GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("User").
		First(&payment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID returns all payments for a given order ID.
func (r *gormPaymentRepository) GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	var payments []model.Payment
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
