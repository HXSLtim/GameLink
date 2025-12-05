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

func (r *gormFeedRepository) ListPaged(ctx context.Context, opts repository.FeedPagedListOptions) ([]model.Feed, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.Feed{}).Preload("Images").Preload("Category")

	if opts.AuthorID != nil {
		q = q.Where("author_id = ?", *opts.AuthorID)
	}
	if opts.CategoryID != nil {
		q = q.Where("category_id = ?", *opts.CategoryID)
	}
	if opts.Keyword != "" {
		q = q.Where("content LIKE ?", "%"+opts.Keyword+"%")
	}
	if opts.ModerationStatus != nil {
		q = q.Where("moderation_status = ?", *opts.ModerationStatus)
	}
	if opts.Visibility != nil {
		q = q.Where("visibility = ?", *opts.Visibility)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var feeds []model.Feed
	if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&feeds).Error; err != nil {
		return nil, 0, err
	}

	return feeds, total, nil
}

func (r *gormFeedRepository) Update(ctx context.Context, feed *model.Feed) error {
	tx := r.db.WithContext(ctx).Model(feed).Where("id = ?", feed.ID).Updates(feed)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *gormFeedRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Feed{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *gormFeedRepository) UpdateModeration(ctx context.Context, feedID uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error {
	updates := map[string]any{
		"moderation_status": status,
		"moderation_note":   note,
	}
	if moderatorID != nil {
		updates["moderated_by"] = *moderatorID
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

func (r *gormFeedRepository) BatchUpdateModeration(ctx context.Context, feedIDs []uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error {
	if len(feedIDs) == 0 {
		return nil
	}
	updates := map[string]any{
		"moderation_status":   status,
		"moderation_note":     note,
		"manual_moderated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}
	if moderatorID != nil {
		updates["moderated_by"] = *moderatorID
	}
	return r.db.WithContext(ctx).Model(&model.Feed{}).Where("id IN ?", feedIDs).Updates(updates).Error
}

func (r *gormFeedRepository) CreateReport(ctx context.Context, report *model.FeedReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *gormFeedRepository) GetReport(ctx context.Context, id uint64) (*model.FeedReport, error) {
	var report model.FeedReport
	if err := r.db.WithContext(ctx).First(&report, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &report, nil
}

func (r *gormFeedRepository) ListReports(ctx context.Context, opts repository.FeedReportListOptions) ([]model.FeedReport, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.FeedReport{})

	if opts.FeedID != nil {
		q = q.Where("feed_id = ?", *opts.FeedID)
	}
	if opts.ReporterID != nil {
		q = q.Where("reporter_id = ?", *opts.ReporterID)
	}
	if opts.Status != nil {
		q = q.Where("status = ?", *opts.Status)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reports []model.FeedReport
	if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *gormFeedRepository) UpdateReport(ctx context.Context, report *model.FeedReport) error {
	tx := r.db.WithContext(ctx).Model(report).Where("id = ?", report.ID).Updates(report)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *gormFeedRepository) CountByStatus(ctx context.Context) (map[model.FeedModerationStatus]int64, error) {
	type result struct {
		Status model.FeedModerationStatus
		Count  int64
	}
	var results []result
	if err := r.db.WithContext(ctx).Model(&model.Feed{}).
		Select("moderation_status as status, count(*) as count").
		Group("moderation_status").
		Find(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[model.FeedModerationStatus]int64)
	for _, res := range results {
		counts[res.Status] = res.Count
	}
	return counts, nil
}

func (r *gormFeedRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	var results []repository.DateValue
	if err := r.db.WithContext(ctx).Model(&model.Feed{}).
		Select("DATE(created_at) as date, count(*) as value").
		Where("created_at >= DATE_SUB(CURRENT_DATE, INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
