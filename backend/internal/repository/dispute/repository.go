package dispute

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// gormDisputeRepository uses GORM to manage order disputes.
type gormDisputeRepository struct {
	db *gorm.DB
}

// NewDisputeRepository creates a new dispute repository instance.
func NewDisputeRepository(db *gorm.DB) repository.DisputeRepository {
	return &gormDisputeRepository{db: db}
}

// Create inserts a new dispute.
func (r *gormDisputeRepository) Create(ctx context.Context, dispute *model.OrderDispute) error {
	return r.db.WithContext(ctx).Create(dispute).Error
}

// Get retrieves a dispute by ID.
func (r *gormDisputeRepository) Get(ctx context.Context, id uint64) (*model.OrderDispute, error) {
	var dispute model.OrderDispute
	if err := r.db.WithContext(ctx).First(&dispute, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &dispute, nil
}

// GetByOrderID retrieves a dispute by order ID.
func (r *gormDisputeRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
	var dispute model.OrderDispute
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&dispute).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &dispute, nil
}

// Update updates an existing dispute.
func (r *gormDisputeRepository) Update(ctx context.Context, dispute *model.OrderDispute) error {
	return r.db.WithContext(ctx).Save(dispute).Error
}

// List returns a page of disputes with filters applied.
func (r *gormDisputeRepository) List(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.OrderDispute{})

	// Apply filters
	if len(opts.Statuses) > 0 {
		query = query.Where("status IN ?", opts.Statuses)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.AssignedToUserID != nil {
		query = query.Where("assigned_to_user_id = ?", *opts.AssignedToUserID)
	}
	if opts.SLABreached != nil {
		query = query.Where("sla_breached = ?", *opts.SLABreached)
	}
	if opts.DateFrom != nil {
		query = query.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("created_at <= ?", *opts.DateTo)
	}
	if trimmed := strings.TrimSpace(opts.Keyword); trimmed != "" {
		like := "%" + trimmed + "%"
		query = query.Where("reason LIKE ? OR description LIKE ?", like, like)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Normalize pagination
	page := repository.NormalizePage(opts.Page)
	pageSize := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * pageSize

	// Fetch disputes
	var disputes []model.OrderDispute
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&disputes).Error; err != nil {
		return nil, 0, err
	}

	return disputes, total, nil
}

// ListPendingAssignment returns disputes pending assignment (status = pending and not assigned).
func (r *gormDisputeRepository) ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	query := r.db.WithContext(ctx).
		Where("status = ?", model.DisputeStatusPending).
		Where("assigned_to_user_id IS NULL")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	var disputes []model.OrderDispute
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&disputes).Error; err != nil {
		return nil, 0, err
	}

	return disputes, total, nil
}

// ListSLABreached returns disputes that have breached SLA.
func (r *gormDisputeRepository) ListSLABreached(ctx context.Context) ([]model.OrderDispute, error) {
	var disputes []model.OrderDispute
	if err := r.db.WithContext(ctx).
		Where("sla_breached = ?", false).
		Where("sla_deadline < ?", time.Now()).
		Where("status NOT IN ?", []model.DisputeStatus{model.DisputeStatusResolved, model.DisputeStatusRejected, model.DisputeStatusCanceled}).
		Find(&disputes).Error; err != nil {
		return nil, err
	}
	return disputes, nil
}

// MarkSLABreached marks a dispute as SLA breached.
func (r *gormDisputeRepository) MarkSLABreached(ctx context.Context, disputeID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OrderDispute{}).
		Where("id = ?", disputeID).
		Updates(map[string]interface{}{
			"sla_breached":    true,
			"sla_breached_at": now,
		}).Error
}

// Delete soft-deletes a dispute.
func (r *gormDisputeRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.OrderDispute{}, id).Error
}

// CountByStatus returns the count of disputes by status.
func (r *gormDisputeRepository) CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.OrderDispute{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPendingCount returns the count of pending disputes.
func (r *gormDisputeRepository) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.OrderDispute{}).
		Where("status = ?", model.DisputeStatusPending).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
