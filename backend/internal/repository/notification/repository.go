package notification

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewNotificationRepository returns a GORM-based notification repository.
func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &gormNotificationRepository{db: db}
}

type gormNotificationRepository struct {
	db *gorm.DB
}

func (r *gormNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.NotificationEvent{}).Where("user_id = ?", opts.UserID)
	if opts.Unread != nil {
		if *opts.Unread {
			query = query.Where("read_at IS NULL")
		} else {
			query = query.Where("read_at IS NOT NULL")
		}
	}
	if len(opts.Priority) > 0 {
		query = query.Where("priority IN ?", opts.Priority)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.NotificationEvent
	if err := query.Order("priority DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *gormNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	tx := r.db.WithContext(ctx).Model(&model.NotificationEvent{}).
		Where("user_id = ?", userID)
	if len(ids) > 0 {
		tx = tx.Where("id IN ?", ids)
	}
	tx = tx.Where("read_at IS NULL").Update("read_at", gorm.Expr("CURRENT_TIMESTAMP"))
	return tx.Error
}

func (r *gormNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NotificationEvent{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&count).Error
	return count, err
}

func (r *gormNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}
