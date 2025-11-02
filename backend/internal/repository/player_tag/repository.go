package playertag

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormPlayerTagRepository struct{ db *gorm.DB }

func NewPlayerTagRepository(db *gorm.DB) repository.PlayerTagRepository {
	return &gormPlayerTagRepository{db: db}
}

func (r *gormPlayerTagRepository) GetTags(ctx context.Context, playerID uint64) ([]string, error) {
	items, err := r.List(ctx, playerID)
	if err != nil {
		return nil, err
	}
	tags := make([]string, len(items))
	for i, item := range items {
		tags[i] = item.Tag
	}
	return tags, nil
}

func (r *gormPlayerTagRepository) List(ctx context.Context, playerID uint64) ([]model.PlayerSkillTag, error) {
	var items []model.PlayerSkillTag
	if err := r.db.WithContext(ctx).Where("player_id = ?", playerID).Order("tag ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormPlayerTagRepository) ReplaceTags(ctx context.Context, playerID uint64, tags []string) error {
	// normalize + de-dup
	m := make(map[string]struct{}, len(tags))
	norm := make([]string, 0, len(tags))
	for _, t := range tags {
		v := strings.TrimSpace(strings.ToLower(t))
		if v == "" {
			continue
		}
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		norm = append(norm, v)
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// delete existing
		if err := tx.Where("player_id = ?", playerID).Delete(&model.PlayerSkillTag{}).Error; err != nil {
			return err
		}
		// insert new
		if len(norm) == 0 {
			return nil
		}
		rows := make([]model.PlayerSkillTag, 0, len(norm))
		for _, v := range norm {
			rows = append(rows, model.PlayerSkillTag{PlayerID: playerID, Tag: v})
		}
		return tx.Create(&rows).Error
	})
}
