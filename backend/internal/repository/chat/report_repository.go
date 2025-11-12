package chat

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type chatReportRepository struct{ db *gorm.DB }

func NewChatReportRepository(db *gorm.DB) repository.ChatReportRepository {
	return &chatReportRepository{db: db}
}

func (r *chatReportRepository) Create(ctx context.Context, report *model.ChatReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *chatReportRepository) Get(ctx context.Context, id uint64) (*model.ChatReport, error) {
	var rep model.ChatReport
	if err := r.db.WithContext(ctx).First(&rep, id).Error; err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *chatReportRepository) Update(ctx context.Context, report *model.ChatReport) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *chatReportRepository) List(ctx context.Context, opts repository.ChatReportListOptions) ([]model.ChatReport, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.ChatReport{})
	if opts.Status != "" {
		tx = tx.Where("status = ?", opts.Status)
	}
	if opts.ReporterID != nil {
		tx = tx.Where("reporter_id = ?", *opts.ReporterID)
	}
	if opts.MessageID != nil {
		tx = tx.Where("message_id = ?", *opts.MessageID)
	}
	if opts.DateFrom != nil {
		tx = tx.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		tx = tx.Where("created_at <= ?", *opts.DateTo)
	}
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.ChatReport
	if err := tx.Order("id DESC").
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
