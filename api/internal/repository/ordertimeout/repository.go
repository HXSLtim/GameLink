package ordertimeout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// Repository 订单超时仓储实现
type Repository struct {
	db *gorm.DB
}

// NewOrderTimeoutRepository 创建订单超时仓储
func NewOrderTimeoutRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 订单超时配置
// ============================================================================

// GetConfig 获取配置
func (r *Repository) GetConfig(ctx context.Context, key string) (*model.OrderTimeoutConfig, error) {
	var config model.OrderTimeoutConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get config: %w", err)
	}
	return &config, nil
}

// ListConfigs 获取所有配置
func (r *Repository) ListConfigs(ctx context.Context) ([]model.OrderTimeoutConfig, error) {
	var configs []model.OrderTimeoutConfig
	if err := r.db.WithContext(ctx).Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("list configs: %w", err)
	}
	return configs, nil
}

// SaveConfig 保存配置（创建或更新）
func (r *Repository) SaveConfig(ctx context.Context, config *model.OrderTimeoutConfig) error {
	var existing model.OrderTimeoutConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", config.ConfigKey).First(&existing).Error
	if err == nil {
		// 更新
		existing.ConfigValue = config.ConfigValue
		existing.Description = config.Description
		return r.db.WithContext(ctx).Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建
		return r.db.WithContext(ctx).Create(config).Error
	}
	return fmt.Errorf("save config: %w", err)
}

// DeleteConfig 删除配置
func (r *Repository) DeleteConfig(ctx context.Context, key string) error {
	result := r.db.WithContext(ctx).Where("config_key = ?", key).Delete(&model.OrderTimeoutConfig{})
	if result.Error != nil {
		return fmt.Errorf("delete config: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================================
// 订单超时日志
// ============================================================================

// CreateLog 创建超时日志
func (r *Repository) CreateLog(ctx context.Context, log *model.OrderTimeoutLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create log: %w", err)
	}
	return nil
}

// GetLog 获取超时日志
func (r *Repository) GetLog(ctx context.Context, id uint64) (*model.OrderTimeoutLog, error) {
	var log model.OrderTimeoutLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get log: %w", err)
	}
	return &log, nil
}

// GetLogWithOrder 获取超时日志（含订单信息）
func (r *Repository) GetLogWithOrder(ctx context.Context, id uint64) (*model.OrderTimeoutLog, error) {
	var log model.OrderTimeoutLog
	if err := r.db.WithContext(ctx).Preload("Order").First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get log with order: %w", err)
	}
	return &log, nil
}

// ListLogsByOrderID 根据订单ID获取超时日志
func (r *Repository) ListLogsByOrderID(ctx context.Context, orderID uint64) ([]model.OrderTimeoutLog, error) {
	var logs []model.OrderTimeoutLog
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("list logs by order: %w", err)
	}
	return logs, nil
}

// ListLogsPaged 分页获取超时日志
func (r *Repository) ListLogsPaged(ctx context.Context, opts repository.OrderTimeoutLogListOptions) ([]model.OrderTimeoutLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.OrderTimeoutLog{})

	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.TimeoutType != nil {
		query = query.Where("timeout_type = ?", *opts.TimeoutType)
	}
	if opts.Action != nil {
		query = query.Where("action = ?", *opts.Action)
	}
	if opts.DateFrom != nil {
		query = query.Where("timeout_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("timeout_at <= ?", *opts.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count logs: %w", err)
	}

	var logs []model.OrderTimeoutLog
	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Order").Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list logs: %w", err)
	}

	return logs, total, nil
}

// GetLogStats 获取超时日志统计
func (r *Repository) GetLogStats(ctx context.Context) (map[model.OrderTimeoutType]int64, error) {
	var results []struct {
		TimeoutType model.OrderTimeoutType
		Count       int64
	}
	if err := r.db.WithContext(ctx).Model(&model.OrderTimeoutLog{}).
		Select("timeout_type, COUNT(*) as count").
		Group("timeout_type").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("get log stats: %w", err)
	}

	stats := make(map[model.OrderTimeoutType]int64)
	for _, r := range results {
		stats[r.TimeoutType] = r.Count
	}
	return stats, nil
}

// ============================================================================
// 客服分配记录
// ============================================================================

// CreateAssignment 创建客服分配记录
func (r *Repository) CreateAssignment(ctx context.Context, assignment *model.OrderServiceAssignment) error {
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return fmt.Errorf("create assignment: %w", err)
	}
	return nil
}

// GetAssignment 获取客服分配记录
func (r *Repository) GetAssignment(ctx context.Context, id uint64) (*model.OrderServiceAssignment, error) {
	var assignment model.OrderServiceAssignment
	if err := r.db.WithContext(ctx).First(&assignment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get assignment: %w", err)
	}
	return &assignment, nil
}

// GetAssignmentWithRelations 获取客服分配记录（含关联）
func (r *Repository) GetAssignmentWithRelations(ctx context.Context, id uint64) (*model.OrderServiceAssignment, error) {
	var assignment model.OrderServiceAssignment
	if err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("ServiceUser").
		Preload("ChatGroup").
		First(&assignment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get assignment with relations: %w", err)
	}
	return &assignment, nil
}

// GetAssignmentByOrderID 根据订单ID获取客服分配记录
func (r *Repository) GetAssignmentByOrderID(ctx context.Context, orderID uint64) (*model.OrderServiceAssignment, error) {
	var assignment model.OrderServiceAssignment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get assignment by order: %w", err)
	}
	return &assignment, nil
}

// ListAssignmentsByServiceUser 根据客服ID获取分配记录
func (r *Repository) ListAssignmentsByServiceUser(ctx context.Context, serviceUserID uint64, status *model.ServiceAssignmentStatus) ([]model.OrderServiceAssignment, error) {
	query := r.db.WithContext(ctx).Where("service_user_id = ?", serviceUserID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var assignments []model.OrderServiceAssignment
	if err := query.Preload("Order").Order("assigned_at DESC").Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("list assignments by service user: %w", err)
	}
	return assignments, nil
}

// ListAssignmentsPaged 分页获取客服分配记录
func (r *Repository) ListAssignmentsPaged(ctx context.Context, opts repository.ServiceAssignmentListOptions) ([]model.OrderServiceAssignment, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.OrderServiceAssignment{})

	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.ServiceUserID != nil {
		query = query.Where("service_user_id = ?", *opts.ServiceUserID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.AssignType != "" {
		query = query.Where("assign_type = ?", opts.AssignType)
	}
	if opts.DateFrom != nil {
		query = query.Where("assigned_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("assigned_at <= ?", *opts.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count assignments: %w", err)
	}

	var assignments []model.OrderServiceAssignment
	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Order").Preload("ServiceUser").Order("assigned_at DESC").Offset(offset).Limit(opts.PageSize).Find(&assignments).Error; err != nil {
		return nil, 0, fmt.Errorf("list assignments: %w", err)
	}

	return assignments, total, nil
}

// UpdateAssignment 更新客服分配记录
func (r *Repository) UpdateAssignment(ctx context.Context, assignment *model.OrderServiceAssignment) error {
	if err := r.db.WithContext(ctx).Save(assignment).Error; err != nil {
		return fmt.Errorf("update assignment: %w", err)
	}
	return nil
}

// UpdateAssignmentStatus 更新客服分配状态
func (r *Repository) UpdateAssignmentStatus(ctx context.Context, id uint64, status model.ServiceAssignmentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	switch status {
	case model.ServiceAssignmentStatusJoined:
		updates["joined_at"] = now
	case model.ServiceAssignmentStatusLeft, model.ServiceAssignmentStatusCompleted:
		updates["left_at"] = now
	}

	result := r.db.WithContext(ctx).Model(&model.OrderServiceAssignment{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update assignment status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteAssignment 删除客服分配记录
func (r *Repository) DeleteAssignment(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.OrderServiceAssignment{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete assignment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetAssignmentStats 获取客服分配统计
func (r *Repository) GetAssignmentStats(ctx context.Context) (map[model.ServiceAssignmentStatus]int64, error) {
	var results []struct {
		Status model.ServiceAssignmentStatus
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&model.OrderServiceAssignment{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("get assignment stats: %w", err)
	}

	stats := make(map[model.ServiceAssignmentStatus]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}
	return stats, nil
}

// GetActiveAssignmentCount 获取活跃分配数量（assigned 或 joined 状态）
func (r *Repository) GetActiveAssignmentCount(ctx context.Context, serviceUserID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.OrderServiceAssignment{}).
		Where("service_user_id = ? AND status IN ?", serviceUserID, []model.ServiceAssignmentStatus{
			model.ServiceAssignmentStatusAssigned,
			model.ServiceAssignmentStatusJoined,
		}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("get active assignment count: %w", err)
	}
	return count, nil
}
