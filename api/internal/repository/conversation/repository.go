package conversation

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository provides persistence for customer-service conversations.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a conversation repository backed by GORM.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create persists a conversation group and its members.
func (r *Repository) Create(ctx context.Context, group *model.ChatGroup, members []*model.ChatGroupMember) error {
	if group == nil {
		return errors.New("group is required")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}
		if len(members) == 0 {
			return nil
		}

		for _, member := range members {
			if member == nil {
				continue
			}
			member.GroupID = group.ID
		}

		return tx.Create(members).Error
	})
}

// Get returns conversation by ID with members loaded.
func (r *Repository) Get(ctx context.Context, conversationID uint64) (*model.ChatGroup, error) {
	var group model.ChatGroup
	err := r.db.WithContext(ctx).
		Preload("Members").
		First(&group, conversationID).
		Error
	if err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &group, nil
}

// ListByUser lists private conversations where user is an active member.
func (r *Repository) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	base := r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Joins("JOIN chat_group_members cgm ON cgm.group_id = chat_groups.id").
		Where("cgm.user_id = ? AND cgm.is_active = ?", userID, true).
		Where("chat_groups.group_type = ?", model.ChatGroupTypePrivate)

	var total int64
	if err := base.Distinct("chat_groups.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []model.ChatGroup
	if err := base.
		Distinct("chat_groups.id").
		Preload("Members").
		Order("chat_groups.updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// Update updates conversation metadata.
func (r *Repository) Update(ctx context.Context, group *model.ChatGroup) error {
	if group == nil {
		return errors.New("group is required")
	}
	return r.db.WithContext(ctx).Save(group).Error
}

// Delete closes conversation (soft close via is_active=false).
func (r *Repository) Delete(ctx context.Context, conversationID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Where("id = ?", conversationID).
		Updates(map[string]any{
			"is_active":      false,
			"deactivated_at": &now,
			"updated_at":     now,
		}).Error
}

// FindActivePrivateByParticipants finds latest active private conversation between user and any agent.
func (r *Repository) FindActivePrivateByParticipants(ctx context.Context, userID uint64, agentIDs []uint64) (*model.ChatGroup, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}

	var groups []model.ChatGroup
	err := r.db.WithContext(ctx).
		Model(&model.ChatGroup{}).
		Joins("JOIN chat_group_members um ON um.group_id = chat_groups.id AND um.user_id = ? AND um.is_active = ?", userID, true).
		Joins("JOIN chat_group_members am ON am.group_id = chat_groups.id AND am.user_id IN ? AND am.is_active = ?", agentIDs, true).
		Where("chat_groups.group_type = ? AND chat_groups.is_active = ?", model.ChatGroupTypePrivate, true).
		Preload("Members").
		Order("chat_groups.updated_at DESC").
		Limit(1).
		Find(&groups).
		Error
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, nil
	}
	return &groups[0], nil
}

// ListMessages returns paged messages for a conversation.
func (r *Repository) ListMessages(ctx context.Context, groupID uint64, page, pageSize int, beforeID *uint64) ([]model.ChatMessage, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	base := r.db.WithContext(ctx).
		Model(&model.ChatMessage{}).
		Where("group_id = ?", groupID)
	if beforeID != nil && *beforeID > 0 {
		base = base.Where("id < ?", *beforeID)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var messages []model.ChatMessage
	if err := base.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// CreateMessage stores a chat message under conversation.
func (r *Repository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	if message == nil {
		return errors.New("message is required")
	}
	return r.db.WithContext(ctx).Create(message).Error
}

// CountActiveByAgentIDs counts active private conversations per agent.
func (r *Repository) CountActiveByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(agentIDs))
	for _, item := range agentIDs {
		result[item] = 0
	}
	if len(agentIDs) == 0 {
		return result, nil
	}

	type row struct {
		UserID uint64
		Total  int64
	}

	var rows []row
	err := r.db.WithContext(ctx).
		Table("chat_group_members cgm").
		Select("cgm.user_id, COUNT(DISTINCT cgm.group_id) AS total").
		Joins("JOIN chat_groups cg ON cg.id = cgm.group_id").
		Where("cgm.user_id IN ? AND cgm.is_active = ?", agentIDs, true).
		Where("cg.group_type = ? AND cg.is_active = ?", model.ChatGroupTypePrivate, true).
		Group("cgm.user_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.UserID] = item.Total
	}
	return result, nil
}

// IsMember checks whether user is an active member in conversation.
func (r *Repository) IsMember(ctx context.Context, groupID, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ChatGroupMember{}).
		Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
