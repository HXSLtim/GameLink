package chat

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewChatMemberRepository creates chat member repository implementation.
func NewChatMemberRepository(db *gorm.DB) repository.ChatMemberRepository {
	return &chatMemberRepository{db: db}
}

type chatMemberRepository struct {
	db *gorm.DB
}

func (r *chatMemberRepository) Add(ctx context.Context, member *model.ChatGroupMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *chatMemberRepository) AddBatch(ctx context.Context, members []*model.ChatGroupMember) error {
	if len(members) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&members).Error
}

func (r *chatMemberRepository) Update(ctx context.Context, member *model.ChatGroupMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *chatMemberRepository) Remove(ctx context.Context, groupID, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.ChatGroupMember{}).Error
}

func (r *chatMemberRepository) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
	var member model.ChatGroupMember
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}
