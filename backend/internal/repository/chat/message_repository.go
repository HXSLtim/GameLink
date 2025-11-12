package chat

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewChatMessageRepository creates chat message repository implementation.
func NewChatMessageRepository(db *gorm.DB) repository.ChatMessageRepository {
    return &chatMessageRepository{db: db}
}

type chatMessageRepository struct {
	db *gorm.DB
}

func (r *chatMessageRepository) Create(ctx context.Context, message *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *chatMessageRepository) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error {
	if len(messages) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(messages).Error
}

func (r *chatMessageRepository) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) {
	base := r.db.WithContext(ctx).Model(&model.ChatMessage{})

	applyFilters := func(db *gorm.DB) *gorm.DB {
		db = db.Where("group_id = ?", opts.GroupID)
		if opts.MessageType != nil {
			db = db.Where("message_type = ?", *opts.MessageType)
		}
		if opts.BeforeID != nil {
			db = db.Where("id < ?", *opts.BeforeID)
		}
		if opts.AfterID != nil {
			db = db.Where("id > ?", *opts.AfterID)
		}
		if opts.DateFrom != nil {
			db = db.Where("created_at >= ?", *opts.DateFrom)
		}
		if opts.DateTo != nil {
			db = db.Where("created_at <= ?", *opts.DateTo)
		}
		if len(opts.AuditStatuses) > 0 {
			db = db.Where("audit_status IN ?", opts.AuditStatuses)
		}
		return db
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	countQuery := applyFilters(base.Session(&gorm.Session{}))
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := applyFilters(base.Session(&gorm.Session{}))
	var messages []model.ChatMessage
	if err := dataQuery.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *chatMessageRepository) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) {
	var message model.ChatMessage
	if err := r.db.WithContext(ctx).First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *chatMessageRepository) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.ChatMessage{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_deleted": true,
			"updated_at": time.Now(),
		}).Error
}

func (r *chatMessageRepository) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) {
	base := r.db.WithContext(ctx).Model(&model.ChatMessage{})

	applyFilters := func(db *gorm.DB) *gorm.DB {
		if opts.GroupID != nil {
			db = db.Where("group_id = ?", *opts.GroupID)
		}
		if opts.SenderID != nil {
			db = db.Where("sender_id = ?", *opts.SenderID)
		}
		if opts.AuditStatus != nil {
			db = db.Where("audit_status = ?", *opts.AuditStatus)
		}
		if opts.DateFrom != nil {
			db = db.Where("created_at >= ?", *opts.DateFrom)
		}
		if opts.DateTo != nil {
			db = db.Where("created_at <= ?", *opts.DateTo)
		}
		return db
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	countQuery := applyFilters(base.Session(&gorm.Session{}))
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := applyFilters(base.Session(&gorm.Session{}))
	var messages []model.ChatMessage
	if err := dataQuery.Order("created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *chatMessageRepository) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error {
	updates := map[string]any{
		"audit_status":  status,
		"updated_at":    time.Now(),
		"reject_reason": reason,
	}
	if moderatorID != nil {
		now := time.Now()
		updates["moderated_by"] = *moderatorID
		updates["moderated_at"] = now
	} else {
		updates["moderated_by"] = nil
		updates["moderated_at"] = nil
	}
	return r.db.WithContext(ctx).
		Model(&model.ChatMessage{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *chatMessageRepository) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
    if len(groupIDs) == 0 { return nil }
    return r.db.WithContext(ctx).Unscoped().
        Where("group_id IN ?", groupIDs).
        Delete(&model.ChatMessage{}).Error
}
