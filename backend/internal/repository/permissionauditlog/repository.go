// Package permissionauditlog provides data access for permission audit logs.
package permissionauditlog

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository defines the interface for permission audit log data access.
type Repository interface {
	// Create creates a single audit log entry.
	Create(ctx context.Context, log *model.PermissionAuditLog) error
	// CreateBatch creates multiple audit log entries in a single transaction.
	CreateBatch(ctx context.Context, logs []*model.PermissionAuditLog) error
	// List retrieves audit logs with filtering and pagination.
	List(ctx context.Context, opts ListOptions) ([]model.PermissionAuditLog, int64, error)
	// Get retrieves a single audit log by ID.
	Get(ctx context.Context, id uint64) (*model.PermissionAuditLog, error)
	// DeleteBefore deletes audit logs older than the specified time (for archiving).
	DeleteBefore(ctx context.Context, before time.Time) (int64, error)
	// CountByDateRange counts audit logs within a date range.
	CountByDateRange(ctx context.Context, from, to time.Time) (int64, error)
}

// ListOptions defines filtering options for audit log queries.
type ListOptions struct {
	Page       int
	PageSize   int
	OperatorID *uint64
	TargetType *model.AuditTargetType
	TargetID   *uint64
	Action     *model.AuditAction
	DateFrom   *time.Time
	DateTo     *time.Time
	RequestID  string
	IPAddress  string
}

// gormRepository implements Repository using GORM.
type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new permission audit log repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Create creates a single audit log entry.
func (r *gormRepository) Create(ctx context.Context, log *model.PermissionAuditLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// CreateBatch creates multiple audit log entries in a single transaction.
// This is optimized for batch inserts from the async writer.
func (r *gormRepository) CreateBatch(ctx context.Context, logs []*model.PermissionAuditLog) error {
	if len(logs) == 0 {
		return nil
	}

	now := time.Now()
	for _, log := range logs {
		if log.CreatedAt.IsZero() {
			log.CreatedAt = now
		}
	}

	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}

// List retrieves audit logs with filtering and pagination.
func (r *gormRepository) List(ctx context.Context, opts ListOptions) ([]model.PermissionAuditLog, int64, error) {
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	q := r.db.WithContext(ctx).Model(&model.PermissionAuditLog{})

	// Apply filters
	if opts.OperatorID != nil {
		q = q.Where("operator_id = ?", *opts.OperatorID)
	}
	if opts.TargetType != nil {
		q = q.Where("target_type = ?", *opts.TargetType)
	}
	if opts.TargetID != nil {
		q = q.Where("target_id = ?", *opts.TargetID)
	}
	if opts.Action != nil {
		q = q.Where("action = ?", *opts.Action)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}
	if opts.RequestID != "" {
		q = q.Where("request_id = ?", opts.RequestID)
	}
	if opts.IPAddress != "" {
		q = q.Where("ip_address = ?", opts.IPAddress)
	}

	// Count total
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch results
	var logs []model.PermissionAuditLog
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// Get retrieves a single audit log by ID.
func (r *gormRepository) Get(ctx context.Context, id uint64) (*model.PermissionAuditLog, error) {
	var log model.PermissionAuditLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &log, nil
}

// DeleteBefore deletes audit logs older than the specified time.
// Returns the number of deleted records.
func (r *gormRepository) DeleteBefore(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&model.PermissionAuditLog{})
	return result.RowsAffected, result.Error
}

// CountByDateRange counts audit logs within a date range.
func (r *gormRepository) CountByDateRange(ctx context.Context, from, to time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PermissionAuditLog{}).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Count(&count).Error
	return count, err
}
