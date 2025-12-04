package user

import (
    "context"

    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type userTagRepository struct{ db *gorm.DB }

// NewUserTagRepository 创建用户标签仓储
func NewUserTagRepository(db *gorm.DB) repository.UserTagRepository {
    return &userTagRepository{db: db}
}

// CreateTag 创建标签
func (r *userTagRepository) CreateTag(ctx context.Context, tag *model.UserTag) error {
    return r.db.WithContext(ctx).Create(tag).Error
}

// GetTag 获取标签
func (r *userTagRepository) GetTag(ctx context.Context, id uint64) (*model.UserTag, error) {
    var tag model.UserTag
    err := r.db.WithContext(ctx).First(&tag, id).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    return &tag, nil
}

// ListTags 标签列表
func (r *userTagRepository) ListTags(ctx context.Context) ([]model.UserTag, error) {
    var tags []model.UserTag
    err := r.db.WithContext(ctx).Order("created_at DESC").Find(&tags).Error
    return tags, err
}

// UpdateTag 更新标签
func (r *userTagRepository) UpdateTag(ctx context.Context, tag *model.UserTag) error {
    tx := r.db.WithContext(ctx).Model(&model.UserTag{}).Where("id = ?", tag.ID).Updates(map[string]any{
        "name":        tag.Name,
        "color":       tag.Color,
        "description": tag.Description,
    })
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    return nil
}

// DeleteTag 删除标签
func (r *userTagRepository) DeleteTag(ctx context.Context, id uint64) error {
    tx := r.db.WithContext(ctx).Delete(&model.UserTag{}, id)
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    return nil
}

// AddTagToUser 为用户添加标签（幂等）
func (r *userTagRepository) AddTagToUser(ctx context.Context, userID uint64, tagID uint64) error {
    rel := &model.UserTagRelation{UserID: userID, TagID: tagID}
    return r.db.WithContext(ctx).Where("user_id = ? AND tag_id = ?", userID, tagID).FirstOrCreate(rel).Error
}

// RemoveTagFromUser 移除用户标签
func (r *userTagRepository) RemoveTagFromUser(ctx context.Context, userID uint64, tagID uint64) error {
    tx := r.db.WithContext(ctx).Where("user_id = ? AND tag_id = ?", userID, tagID).Delete(&model.UserTagRelation{})
    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        return repository.ErrNotFound
    }
    return nil
}

// GetUserTags 获取用户的所有标签
func (r *userTagRepository) GetUserTags(ctx context.Context, userID uint64) ([]model.UserTag, error) {
    var tags []model.UserTag
    err := r.db.WithContext(ctx).
        Table((model.UserTag{}).TableName()).
        Joins("JOIN user_tag_relations utr ON utr.tag_id = user_tags.id").
        Where("utr.user_id = ?", userID).
        Order("user_tags.created_at DESC").
        Find(&tags).Error
    return tags, err
}

// BatchSetUserTags 批量设置用户标签（事务）
func (r *userTagRepository) BatchSetUserTags(ctx context.Context, userID uint64, tagIDs []uint64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("user_id = ?", userID).Delete(&model.UserTagRelation{}).Error; err != nil {
            return err
        }
        if len(tagIDs) > 0 {
            relations := make([]model.UserTagRelation, len(tagIDs))
            for i, tagID := range tagIDs {
                relations[i] = model.UserTagRelation{UserID: userID, TagID: tagID}
            }
            if err := tx.Create(&relations).Error; err != nil {
                return err
            }
        }
        return nil
    })
}

// GetUsersByTag 根据标签查询用户（分页）
func (r *userTagRepository) GetUsersByTag(ctx context.Context, tagID uint64, page, pageSize int) ([]model.User, int64, error) {
    page = repository.NormalizePage(page)
    pageSize = repository.NormalizePageSize(pageSize)
    offset := (page - 1) * pageSize

    q := r.db.WithContext(ctx).
        Model(&model.User{}).
        Joins("JOIN user_tag_relations utr ON utr.user_id = users.id").
        Where("utr.tag_id = ?", tagID)

    var total int64
    if err := q.Distinct("users.id").Count(&total).Error; err != nil {
        return nil, 0, err
    }

    var users []model.User
    if err := q.Order("users.created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    return users, total, nil
}

