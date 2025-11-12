package chat

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewChatGroupRepository creates a chat group repository implementation.
func NewChatGroupRepository(db *gorm.DB) repository.ChatGroupRepository {
    return &chatGroupRepository{db: db}

}

func (r *chatGroupRepository) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
    if limit <= 0 || limit > 1000 {
        limit = 100
    }
    var groups []model.ChatGroup
    if err := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
        Where("group_type = ? AND is_active = ? AND deactivated_at IS NOT NULL AND deactivated_at < ?", model.ChatGroupTypeOrder, false, cutoff).
        Order("deactivated_at ASC").
        Limit(limit).
        Find(&groups).Error; err != nil {
        return nil, err
    }
    return groups, nil
}

func (r *chatGroupRepository) DeleteByIDs(ctx context.Context, ids []uint64) error {
    if len(ids) == 0 {
        return nil
    }
    return r.db.WithContext(ctx).Unscoped().
        Where("id IN ?", ids).
        Delete(&model.ChatGroup{}).Error
}

type chatGroupRepository struct {
	db *gorm.DB
}

func (r *chatGroupRepository) Create(ctx context.Context, group *model.ChatGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *chatGroupRepository) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *chatGroupRepository) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("related_order_id = ?", orderID).
		First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *chatGroupRepository) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroup{}).
		Joins("JOIN chat_group_members AS m ON m.group_id = chat_groups.id AND m.user_id = ?", userID)

	if opts.GroupType != nil {
		tx = tx.Where("chat_groups.group_type = ?", *opts.GroupType)
	}
	if !opts.IncludeInactive {
		tx = tx.Where("chat_groups.is_active = ?", true)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("chat_groups.group_name LIKE ? OR chat_groups.description LIKE ?", like, like)
	}
	if opts.RelatedOrderID != nil {
		tx = tx.Where("chat_groups.related_order_id = ?", *opts.RelatedOrderID)
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	totalTx := tx.Session(&gorm.Session{NewDB: true})
	var total int64
	if err := totalTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := tx.
		Preload("Members").
		Order("chat_groups.updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *chatGroupRepository) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatGroupMember{}).
		Where("group_id = ?", groupID)
	if opts.Role != "" {
		tx = tx.Where("role = ?", opts.Role)
	}
	if opts.Keyword != "" {
		like := fmt.Sprintf("%%%s%%", opts.Keyword)
		tx = tx.Where("nickname LIKE ?", like)
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	totalTx := tx.Session(&gorm.Session{NewDB: true})
	var total int64
	if err := totalTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var members []model.ChatGroupMember
	if err := tx.
		Order("joined_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *chatGroupRepository) Update(ctx context.Context, group *model.ChatGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *chatGroupRepository) Deactivate(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":       false,
			"deactivated_at":  now,
		}).Error
}
