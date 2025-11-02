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

	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	if len(opts.Methods) > 0 {
		query = query.Where("method IN ?", opts.Methods)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at <= ?", *opts.DateTo)
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
		"status":            payment.Status,
		"provider_trade_no": payment.ProviderTradeNo,
		"provider_raw":      payment.ProviderRaw,
		"paid_at":           payment.PaidAt,
		"refunded_at":       payment.RefundedAt,
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

