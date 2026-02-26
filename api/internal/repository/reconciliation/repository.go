package reconciliation

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository provides reconciliation persistence.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a reconciliation repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List returns reconciliations with filters and pagination.
func (r *Repository) List(ctx context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)

	q := r.db.WithContext(ctx).Model(&model.Reconciliation{})
	if opts.Type != nil {
		q = q.Where("type = ?", *opts.Type)
	}
	if opts.Status != nil {
		q = q.Where("status = ?", *opts.Status)
	}
	if opts.DateFrom != nil {
		q = q.Where("reconciliation_date >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("reconciliation_date <= ?", *opts.DateTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Reconciliation
	if err := q.Order("reconciliation_date DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Get returns a reconciliation by id, optionally preloading details.
func (r *Repository) Get(ctx context.Context, id uint64, withDetails bool) (*model.Reconciliation, error) {
	q := r.db.WithContext(ctx)
	if withDetails {
		q = q.Preload("Details", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_no ASC, id ASC")
		})
	}

	var rec model.Reconciliation
	if err := q.First(&rec, id).Error; err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &rec, nil
}

// Create creates a reconciliation and its details in one transaction.
func (r *Repository) Create(ctx context.Context, rec *model.Reconciliation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		details := rec.Details
		rec.Details = nil
		if err := tx.Create(rec).Error; err != nil {
			return err
		}

		if len(details) > 0 {
			for i := range details {
				details[i].ReconciliationID = rec.ID
			}
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		rec.Details = details
		return nil
	})
}

// Execute performs transaction-safe reconciliation execution and status transition.
func (r *Repository) Execute(ctx context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
	var out model.Reconciliation

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rec model.Reconciliation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&rec, id).Error; err != nil {
			return repository.WrapNotFound(err)
		}

		if rec.Status != model.ReconciliationStatusPending {
			return fmt.Errorf("%w: current=%s", repository.ErrInvalidStatusTransition, rec.Status)
		}

		now := time.Now()
		processedBy := opts.ProcessedBy

		if err := tx.Model(&model.Reconciliation{}).
			Where("id = ?", rec.ID).
			Updates(map[string]any{
				"status":       model.ReconciliationStatusProgress,
				"processed_by": processedBy,
				"processed_at": now,
			}).Error; err != nil {
			return err
		}

		var summary struct {
			TotalRecords     int64
			MatchedRecords   int64
			DifferenceAmount int64
		}
		if err := tx.Model(&model.ReconciliationDetail{}).
			Select(`
				COUNT(*) AS total_records,
				COALESCE(SUM(CASE WHEN difference_amount = 0 THEN 1 ELSE 0 END), 0) AS matched_records,
				COALESCE(SUM(ABS(difference_amount)), 0) AS difference_amount
			`).
			Where("reconciliation_id = ?", rec.ID).
			Scan(&summary).Error; err != nil {
			return err
		}

		finalStatus := model.ReconciliationStatusSuccess
		if opts.ForceStatus != nil {
			finalStatus = *opts.ForceStatus
		} else if summary.TotalRecords == 0 ||
			summary.TotalRecords != summary.MatchedRecords ||
			summary.DifferenceAmount != 0 {
			finalStatus = model.ReconciliationStatusException
		}

		if err := tx.Model(&model.Reconciliation{}).
			Where("id = ?", rec.ID).
			Updates(map[string]any{
				"status":            finalStatus,
				"processed_by":      processedBy,
				"processed_at":      now,
				"total_records":     int(summary.TotalRecords),
				"matched_records":   int(summary.MatchedRecords),
				"difference_amount": summary.DifferenceAmount,
			}).Error; err != nil {
			return err
		}

		if err := tx.Preload("Details", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_no ASC, id ASC")
		}).First(&out, rec.ID).Error; err != nil {
			return repository.WrapNotFound(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &out, nil
}
