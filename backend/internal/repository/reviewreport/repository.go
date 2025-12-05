package reviewreport

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormReviewReportRepository struct {
	db *gorm.DB
}

// NewReviewReportRepository creates a new review report repository
func NewReviewReportRepository(db *gorm.DB) repository.ReviewReportRepository {
	return &gormReviewReportRepository{db: db}
}

// Create creates a new review report
func (r *gormReviewReportRepository) Create(ctx context.Context, report *model.ReviewReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

// Get retrieves a review report by ID
func (r *gormReviewReportRepository) Get(ctx context.Context, id uint64) (*model.ReviewReport, error) {
	var report model.ReviewReport
	if err := r.db.WithContext(ctx).First(&report, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &report, nil
}

// List retrieves review reports with filtering and pagination
func (r *gormReviewReportRepository) List(ctx context.Context, opts repository.ReviewReportListOptions) ([]model.ReviewReport, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.ReviewReport{})

	// Apply filters
	if opts.ReviewID != nil {
		q = q.Where("review_id = ?", *opts.ReviewID)
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

	// Get total count
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	var reports []model.ReviewReport
	if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// Update updates a review report
func (r *gormReviewReportRepository) Update(ctx context.Context, report *model.ReviewReport) error {
	tx := r.db.WithContext(ctx).Model(report).Where("id = ?", report.ID).Updates(map[string]any{
		"status":        report.Status,
		"handled_by":    report.HandledBy,
		"handled_at":    report.HandledAt,
		"handling_note": report.HandlingNote,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}
