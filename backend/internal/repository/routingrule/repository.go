package routingrule

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// RoutingRuleRepository 分流规则仓储接口
// Requirements: 16.5
type RoutingRuleRepository interface {
	// Rule CRUD operations
	Create(ctx context.Context, rule *model.RoutingRule) error
	Get(ctx context.Context, id uint64) (*model.RoutingRule, error)
	Update(ctx context.Context, rule *model.RoutingRule) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, opts ListOptions) ([]model.RoutingRule, int64, error)
	ToggleStatus(ctx context.Context, id uint64, status model.RuleStatus) error

	// Rule matching operations
	// Requirements: 16.2 - 按优先级顺序获取所有活跃规则
	ListActiveByPriority(ctx context.Context) ([]model.RoutingRule, error)

	// History operations
	CreateHistory(ctx context.Context, history *model.RoutingRuleHistory) error
	GetHistory(ctx context.Context, ruleID uint64) ([]model.RoutingRuleHistory, error)

	// Routing log operations
	CreateRoutingLog(ctx context.Context, log *model.RoutingLog) error
	GetRoutingLogByPayment(ctx context.Context, paymentID uint64) (*model.RoutingLog, error)
	ListRoutingLogs(ctx context.Context, opts RoutingLogListOptions) ([]model.RoutingLog, int64, error)
}

// ListOptions 分流规则查询选项
type ListOptions struct {
	Status         *model.RuleStatus
	TargetEntityID *uint64
	Keyword        string
	Page           int
	PageSize       int
}

// RoutingLogListOptions 分流日志查询选项
type RoutingLogListOptions struct {
	PaymentID          *uint64
	OrderID            *uint64
	CollectionEntityID *uint64
	IsDefault          *bool
	IsFallback         *bool
	Page               int
	PageSize           int
}

type routingRuleRepository struct {
	db *gorm.DB
}

// NewRoutingRuleRepository 创建分流规则仓储
func NewRoutingRuleRepository(db *gorm.DB) RoutingRuleRepository {
	return &routingRuleRepository{db: db}
}

// Create 创建分流规则
func (r *routingRuleRepository) Create(ctx context.Context, rule *model.RoutingRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// Get 获取分流规则
func (r *routingRuleRepository) Get(ctx context.Context, id uint64) (*model.RoutingRule, error) {
	var rule model.RoutingRule
	err := r.db.WithContext(ctx).
		Preload("TargetEntity").
		First(&rule, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// Update 更新分流规则
func (r *routingRuleRepository) Update(ctx context.Context, rule *model.RoutingRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// Delete 删除分流规则
func (r *routingRuleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.RoutingRule{}, id).Error
}

// List 查询分流规则列表
// Requirements: 16.5
func (r *routingRuleRepository) List(ctx context.Context, opts ListOptions) ([]model.RoutingRule, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RoutingRule{})

	// 过滤条件
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.TargetEntityID != nil {
		query = query.Where("target_entity_id = ?", *opts.TargetEntityID)
	}
	if opts.Keyword != "" {
		keyword := "%" + opts.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	// 查询数据（按优先级排序）
	var rules []model.RoutingRule
	err := query.
		Preload("TargetEntity").
		Order("priority ASC, created_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// ToggleStatus 切换分流规则状态
// Requirements: 16.4
func (r *routingRuleRepository) ToggleStatus(ctx context.Context, id uint64, status model.RuleStatus) error {
	return r.db.WithContext(ctx).Model(&model.RoutingRule{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ListActiveByPriority 按优先级顺序获取所有活跃规则
// Requirements: 16.2 - Property 15: 收款分流规则优先级
func (r *routingRuleRepository) ListActiveByPriority(ctx context.Context) ([]model.RoutingRule, error) {
	var rules []model.RoutingRule
	err := r.db.WithContext(ctx).
		Where("status = ?", model.RuleStatusActive).
		Preload("TargetEntity").
		Order("priority ASC").
		Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// CreateHistory 创建分流规则修改历史
func (r *routingRuleRepository) CreateHistory(ctx context.Context, history *model.RoutingRuleHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistory 获取分流规则修改历史
func (r *routingRuleRepository) GetHistory(ctx context.Context, ruleID uint64) ([]model.RoutingRuleHistory, error) {
	var histories []model.RoutingRuleHistory
	err := r.db.WithContext(ctx).
		Where("routing_rule_id = ?", ruleID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// CreateRoutingLog 创建分流日志
// Requirements: 17.3 - Property 17: 收款分流记录完整性
func (r *routingRuleRepository) CreateRoutingLog(ctx context.Context, log *model.RoutingLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetRoutingLogByPayment 根据支付ID获取分流日志
func (r *routingRuleRepository) GetRoutingLogByPayment(ctx context.Context, paymentID uint64) (*model.RoutingLog, error) {
	var log model.RoutingLog
	err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Preload("CollectionEntity").
		Preload("MatchedRule").
		First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &log, nil
}

// ListRoutingLogs 查询分流日志列表
func (r *routingRuleRepository) ListRoutingLogs(ctx context.Context, opts RoutingLogListOptions) ([]model.RoutingLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RoutingLog{})

	// 过滤条件
	if opts.PaymentID != nil {
		query = query.Where("payment_id = ?", *opts.PaymentID)
	}
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.CollectionEntityID != nil {
		query = query.Where("collection_entity_id = ?", *opts.CollectionEntityID)
	}
	if opts.IsDefault != nil {
		query = query.Where("is_default = ?", *opts.IsDefault)
	}
	if opts.IsFallback != nil {
		query = query.Where("is_fallback = ?", *opts.IsFallback)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	// 查询数据
	var logs []model.RoutingLog
	err := query.
		Preload("CollectionEntity").
		Preload("MatchedRule").
		Order("created_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
