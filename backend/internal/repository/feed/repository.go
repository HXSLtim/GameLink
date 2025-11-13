package feed

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

const defaultFeedPageSize = 20
const maxFeedPageSize = 50

// NewFeedRepository creates a GORM implementation of repository.FeedRepository.
func NewFeedRepository(db *gorm.DB) repository.FeedRepository {
	return &gormFeedRepository{db: db}
}

type gormFeedRepository struct {
	db *gorm.DB
}

func (r *gormFeedRepository) Create(ctx context.Context, feed *model.Feed) error {
	return r.db.WithContext(ctx).Create(feed).Error
}

func (r *gormFeedRepository) Get(ctx context.Context, id uint64) (*model.Feed, error) {
	var feed model.Feed
	if err := r.db.WithContext(ctx).Preload("Images").First(&feed, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &feed, nil
}

func (r *gormFeedRepository) List(ctx context.Context, opts repository.FeedListOptions) ([]model.Feed, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = defaultFeedPageSize
	}
	if limit > maxFeedPageSize {
		limit = maxFeedPageSize
	}

	query := r.db.WithContext(ctx).Model(&model.Feed{}).Preload("Images").Order("id DESC")
	if opts.CursorBefore != nil {
		query = query.Where("id < ?", *opts.CursorBefore)
	}
	if opts.AuthorID != nil {
		query = query.Where("author_id = ?", *opts.AuthorID)
	}
	if len(opts.Visibility) > 0 {
		query = query.Where("visibility IN ?", opts.Visibility)
	}
	if opts.OnlyApproved {
		query = query.Where("moderation_status = ?", model.FeedModerationApproved)
	}

	var feeds []model.Feed
	if err := query.Limit(limit).Find(&feeds).Error; err != nil {
		return nil, err
	}
	return feeds, nil
}

func (r *gormFeedRepository) UpdateModeration(ctx context.Context, feedID uint64, status model.FeedModerationStatus, note string, manual bool) error {
	updates := map[string]any{
		"moderation_status": status,
		"moderation_note":   note,
	}
	if manual {
		updates["manual_moderated_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	} else {
		updates["auto_moderated_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}
	tx := r.db.WithContext(ctx).Model(&model.Feed{}).Where("id = ?", feedID).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *gormFeedRepository) CreateReport(ctx context.Context, report *model.FeedReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}
