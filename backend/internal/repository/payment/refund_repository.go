package payment

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormRefundRecordRepository implements RefundRecordRepository using GORM.
type gormRefundRecordRepository struct {
	db *gorm.DB
}

// NewRefundRecordRepository creates a new RefundRecordRepository instance.
func NewRefundRecordRepository(db *gorm.DB) repository.RefundRecordRepository {
	return &gormRefundRecordRepository{db: db}
}

// Create inserts a new refund record.
func (r *gormRefundRecordRepository) Create(ctx context.Context, record *model.RefundRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// Get returns a refund record by id.
func (r *gormRefundRecordRepository) Get(ctx context.Context, id uint64) (*model.RefundRecord, error) {
	var record model.RefundRecord
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

// Update updates a refund record.
func (r *gormRefundRecordRepository) Update(ctx context.Context, record *model.RefundRecord) error {
	tx := r.db.WithContext(ctx).Model(record).Where("id = ?", record.ID).Updates(map[string]any{
		"status":            record.Status,
		"provider_trade_no": record.ProviderTradeNo,
		"refunded_at":       record.RefundedAt,
		"note":              record.Note,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ListByPaymentID returns all refund records for a given payment ID.
func (r *gormRefundRecordRepository) ListByPaymentID(ctx context.Context, paymentID uint64) ([]model.RefundRecord, error) {
	var records []model.RefundRecord
	if err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// ListByOrderID returns all refund records for a given order ID.
func (r *gormRefundRecordRepository) ListByOrderID(ctx context.Context, orderID uint64) ([]model.RefundRecord, error) {
	var records []model.RefundRecord
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
