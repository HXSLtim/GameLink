package ordertimeout

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrOrderNotFound 订单不存在
	ErrOrderNotFound = errors.New("order not found")
	// ErrAlreadyAssigned 已分配客服
	ErrAlreadyAssigned = errors.New("order already has service assignment")
	// ErrNoAvailableService 无可用客服
	ErrNoAvailableService = errors.New("no available service user")
)

// OrderTimeoutService 订单超时处理服务
type OrderTimeoutService struct {
	repo   repository.OrderTimeoutRepository
	orders repository.OrderRepository
	users  repository.UserRepository
}

// NewOrderTimeoutService 创建订单超时处理服务
func NewOrderTimeoutService(
	repo repository.OrderTimeoutRepository,
	orders repository.OrderRepository,
	users repository.UserRepository,
) *OrderTimeoutService {
	return &OrderTimeoutService{
		repo:   repo,
		orders: orders,
		users:  users,
	}
}

// ============================================================================
// 配置管理
// ============================================================================

// GetConfig 获取配置
func (s *OrderTimeoutService) GetConfig(ctx context.Context, key string) (*model.OrderTimeoutConfig, error) {
	return s.repo.GetConfig(ctx, key)
}

// ListConfigs 获取所有配置
func (s *OrderTimeoutService) ListConfigs(ctx context.Context) ([]model.OrderTimeoutConfig, error) {
	return s.repo.ListConfigs(ctx)
}

// SaveConfig 保存配置
func (s *OrderTimeoutService) SaveConfig(ctx context.Context, key, value, description string) error {
	config := &model.OrderTimeoutConfig{
		ConfigKey:   key,
		ConfigValue: value,
		Description: description,
	}
	return s.repo.SaveConfig(ctx, config)
}

// GetPaymentTimeoutMinutes 获取支付超时时间（分钟）
func (s *OrderTimeoutService) GetPaymentTimeoutMinutes(ctx context.Context) (int, error) {
	config, err := s.repo.GetConfig(ctx, model.PaymentTimeoutMinutes)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 30, nil // 默认30分钟
		}
		return 0, err
	}
	minutes, err := strconv.Atoi(config.ConfigValue)
	if err != nil {
		return 30, nil
	}
	return minutes, nil
}

// GetAcceptTimeoutMinutes 获取接单超时时间（分钟）
func (s *OrderTimeoutService) GetAcceptTimeoutMinutes(ctx context.Context) (int, error) {
	config, err := s.repo.GetConfig(ctx, model.OrderAcceptTimeoutMinutes)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 30, nil // 默认30分钟
		}
		return 0, err
	}
	minutes, err := strconv.Atoi(config.ConfigValue)
	if err != nil {
		return 30, nil
	}
	return minutes, nil
}

// IsAutoCancelEnabled 是否启用自动取消
func (s *OrderTimeoutService) IsAutoCancelEnabled(ctx context.Context) (bool, error) {
	config, err := s.repo.GetConfig(ctx, model.AutoCancelEnabled)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return true, nil // 默认启用
		}
		return false, err
	}
	return config.ConfigValue == "true", nil
}

// IsAutoRefundEnabled 是否启用自动退款
func (s *OrderTimeoutService) IsAutoRefundEnabled(ctx context.Context) (bool, error) {
	config, err := s.repo.GetConfig(ctx, model.AutoRefundEnabled)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return true, nil // 默认启用
		}
		return false, err
	}
	return config.ConfigValue == "true", nil
}

// IsAutoAssignServiceEnabled 是否启用自动分配客服
func (s *OrderTimeoutService) IsAutoAssignServiceEnabled(ctx context.Context) (bool, error) {
	config, err := s.repo.GetConfig(ctx, model.AutoAssignServiceEnabled)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return true, nil // 默认启用
		}
		return false, err
	}
	return config.ConfigValue == "true", nil
}

// ============================================================================
// 超时日志
// ============================================================================

// LogTimeoutInput 记录超时日志输入
type LogTimeoutInput struct {
	OrderID           uint64
	TimeoutType       model.OrderTimeoutType
	Action            model.OrderTimeoutAction
	RefundAmountCents int64
	RefundRecordID    *uint64
	Remark            string
}

// LogTimeout 记录超时日志
func (s *OrderTimeoutService) LogTimeout(ctx context.Context, input LogTimeoutInput) (*model.OrderTimeoutLog, error) {
	// 验证订单是否存在
	if _, err := s.orders.Get(ctx, input.OrderID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}

	log := &model.OrderTimeoutLog{
		OrderID:           input.OrderID,
		TimeoutType:       input.TimeoutType,
		TimeoutAt:         time.Now(),
		Action:            input.Action,
		RefundAmountCents: input.RefundAmountCents,
		RefundRecordID:    input.RefundRecordID,
		Remark:            input.Remark,
	}

	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, fmt.Errorf("create log: %w", err)
	}

	return log, nil
}

// GetLog 获取超时日志
func (s *OrderTimeoutService) GetLog(ctx context.Context, id uint64) (*model.OrderTimeoutLog, error) {
	return s.repo.GetLogWithOrder(ctx, id)
}

// ListLogsByOrderID 根据订单ID获取超时日志
func (s *OrderTimeoutService) ListLogsByOrderID(ctx context.Context, orderID uint64) ([]model.OrderTimeoutLog, error) {
	return s.repo.ListLogsByOrderID(ctx, orderID)
}

// ListLogsPaged 分页获取超时日志
func (s *OrderTimeoutService) ListLogsPaged(ctx context.Context, opts repository.OrderTimeoutLogListOptions) ([]model.OrderTimeoutLog, *model.Pagination, error) {
	logs, total, err := s.repo.ListLogsPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return logs, pagination, nil
}

// GetLogStats 获取超时日志统计
func (s *OrderTimeoutService) GetLogStats(ctx context.Context) (map[model.OrderTimeoutType]int64, error) {
	return s.repo.GetLogStats(ctx)
}

// ============================================================================
// 客服分配
// ============================================================================

// AssignServiceInput 分配客服输入
type AssignServiceInput struct {
	OrderID       uint64
	ServiceUserID uint64
	ChatGroupID   *uint64
	AssignType    string // auto/manual
	Remark        string
}

// AssignService 分配客服
func (s *OrderTimeoutService) AssignService(ctx context.Context, input AssignServiceInput) (*model.OrderServiceAssignment, error) {
	// 验证订单是否存在
	if _, err := s.orders.Get(ctx, input.OrderID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}

	// 检查是否已分配
	existing, err := s.repo.GetAssignmentByOrderID(ctx, input.OrderID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyAssigned
	}

	// 验证客服用户是否存在
	if _, err := s.users.Get(ctx, input.ServiceUserID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: service user not found", ErrValidation)
		}
		return nil, fmt.Errorf("get service user: %w", err)
	}

	assignType := input.AssignType
	if assignType == "" {
		assignType = "manual"
	}

	assignment := &model.OrderServiceAssignment{
		OrderID:       input.OrderID,
		ServiceUserID: input.ServiceUserID,
		ChatGroupID:   input.ChatGroupID,
		Status:        model.ServiceAssignmentStatusAssigned,
		AssignedAt:    time.Now(),
		AssignType:    assignType,
		Remark:        input.Remark,
	}

	if err := s.repo.CreateAssignment(ctx, assignment); err != nil {
		return nil, fmt.Errorf("create assignment: %w", err)
	}

	return assignment, nil
}

// GetAssignment 获取客服分配记录
func (s *OrderTimeoutService) GetAssignment(ctx context.Context, id uint64) (*model.OrderServiceAssignment, error) {
	return s.repo.GetAssignmentWithRelations(ctx, id)
}

// GetAssignmentByOrderID 根据订单ID获取客服分配记录
func (s *OrderTimeoutService) GetAssignmentByOrderID(ctx context.Context, orderID uint64) (*model.OrderServiceAssignment, error) {
	return s.repo.GetAssignmentByOrderID(ctx, orderID)
}

// ListAssignmentsByServiceUser 根据客服ID获取分配记录
func (s *OrderTimeoutService) ListAssignmentsByServiceUser(ctx context.Context, serviceUserID uint64, status *model.ServiceAssignmentStatus) ([]model.OrderServiceAssignment, error) {
	return s.repo.ListAssignmentsByServiceUser(ctx, serviceUserID, status)
}

// ListAssignmentsPaged 分页获取客服分配记录
func (s *OrderTimeoutService) ListAssignmentsPaged(ctx context.Context, opts repository.ServiceAssignmentListOptions) ([]model.OrderServiceAssignment, *model.Pagination, error) {
	assignments, total, err := s.repo.ListAssignmentsPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return assignments, pagination, nil
}

// UpdateAssignmentStatus 更新客服分配状态
func (s *OrderTimeoutService) UpdateAssignmentStatus(ctx context.Context, id uint64, status model.ServiceAssignmentStatus) error {
	return s.repo.UpdateAssignmentStatus(ctx, id, status)
}

// JoinAssignment 客服加入订单
func (s *OrderTimeoutService) JoinAssignment(ctx context.Context, id uint64) error {
	return s.repo.UpdateAssignmentStatus(ctx, id, model.ServiceAssignmentStatusJoined)
}

// LeaveAssignment 客服离开订单
func (s *OrderTimeoutService) LeaveAssignment(ctx context.Context, id uint64) error {
	return s.repo.UpdateAssignmentStatus(ctx, id, model.ServiceAssignmentStatusLeft)
}

// CompleteAssignment 完成客服分配
func (s *OrderTimeoutService) CompleteAssignment(ctx context.Context, id uint64) error {
	return s.repo.UpdateAssignmentStatus(ctx, id, model.ServiceAssignmentStatusCompleted)
}

// DeleteAssignment 删除客服分配记录
func (s *OrderTimeoutService) DeleteAssignment(ctx context.Context, id uint64) error {
	return s.repo.DeleteAssignment(ctx, id)
}

// GetAssignmentStats 获取客服分配统计
func (s *OrderTimeoutService) GetAssignmentStats(ctx context.Context) (map[model.ServiceAssignmentStatus]int64, error) {
	return s.repo.GetAssignmentStats(ctx)
}

// GetActiveAssignmentCount 获取客服活跃分配数量
func (s *OrderTimeoutService) GetActiveAssignmentCount(ctx context.Context, serviceUserID uint64) (int64, error) {
	return s.repo.GetActiveAssignmentCount(ctx, serviceUserID)
}
