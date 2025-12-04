// Package alert provides alert repository implementation.
package alert

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// gormAlertRepository implements model.AlertRepository using GORM.
type gormAlertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new alert repository.
func NewAlertRepository(db *gorm.DB) model.AlertRepository {
	return &gormAlertRepository{db: db}
}

// Create creates a new alert.
func (r *gormAlertRepository) Create(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetByID retrieves an alert by ID.
func (r *gormAlertRepository) GetByID(ctx context.Context, id uint) (*model.Alert, error) {
	var alert model.Alert
	if err := r.db.WithContext(ctx).First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

// List retrieves alerts with options.
func (r *gormAlertRepository) List(ctx context.Context, opts model.AlertQueryOptions) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	q := r.db.WithContext(ctx).Model(&model.Alert{})

	// Apply filters
	if opts.Level != "" {
		q = q.Where("level = ?", opts.Level)
	}
	if opts.Type != "" {
		q = q.Where("type = ?", opts.Type)
	}
	if opts.IsRead != nil {
		q = q.Where("is_read = ?", *opts.IsRead)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		endOfDay := opts.DateTo.Add(24 * time.Hour)
		q = q.Where("created_at < ?", endOfDay)
	}

	// Count total
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// MarkAsRead marks an alert as read.
func (r *gormAlertRepository) MarkAsRead(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Alert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// BatchMarkAsRead marks multiple alerts as read.
func (r *gormAlertRepository) BatchMarkAsRead(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Alert{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// GetUnreadCount returns the count of unread alerts.
func (r *gormAlertRepository) GetUnreadCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Alert{}).
		Where("is_read = ?", false).
		Count(&count).Error
	return count, err
}
