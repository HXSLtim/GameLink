package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/apierr"
)

var (
	// ErrDisputeNotFound dispute not found
	ErrDisputeNotFound = repository.ErrNotFound
	// ErrDisputeValidation validation failed
	ErrDisputeValidation = apierr.BadRequest("验证失败")
	// ErrDisputeInvalidStatus invalid dispute status
	ErrDisputeInvalidStatus = apierr.BadRequest("纠纷状态无效")
	// ErrDisputeUnauthorized unauthorized operation
	ErrDisputeUnauthorized = apierr.Unauthorized("无权限操作")
	// ErrDisputeSLAExpired SLA deadline has passed
	ErrDisputeSLAExpired = apierr.BadRequest("SLA期限已过期")
	// ErrOrderNotFound order not found
	ErrOrderNotFound = apierr.NotFound("订单不存在")
	// ErrDisputeExists dispute already exists for this order
	ErrDisputeExists = apierr.Conflict("该订单已存在纠纷")
	// ErrCannotInitiateDispute cannot initiate dispute for this order
	ErrCannotInitiateDispute = apierr.BadRequest("无法为该订单发起纠纷")
)

// DisputeService handles dispute operations
type DisputeService struct {
	disputes       repository.DisputeRepository
	orders         repoiface.OrderReadWriter
	users          repository.UserRepository
	players        repository.PlayerRepository
	operationLogs  repository.OperationLogRepository
	notifications  repository.NotificationRepository
	payments       repository.PaymentRepository
	defaultSLAMins int // default SLA in minutes (30)
}

// NewDisputeService creates a new dispute service
func NewDisputeService(
	disputes repository.DisputeRepository,
	orders repoiface.OrderReadWriter,
	users repository.UserRepository,
	operationLogs repository.OperationLogRepository,
	notifications repository.NotificationRepository,
	payments repository.PaymentRepository,
) *DisputeService {
	return &DisputeService{
		disputes:       disputes,
		orders:         orders,
		users:          users,
		operationLogs:  operationLogs,
		notifications:  notifications,
		payments:       payments,
		defaultSLAMins: 30,
	}
}

// NewDisputeServiceWithPlayers creates a new dispute service with player repository
func NewDisputeServiceWithPlayers(
	disputes repository.DisputeRepository,
	orders repoiface.OrderReadWriter,
	users repository.UserRepository,
	players repository.PlayerRepository,
	operationLogs repository.OperationLogRepository,
	notifications repository.NotificationRepository,
	payments repository.PaymentRepository,
) *DisputeService {
	return &DisputeService{
		disputes:       disputes,
		orders:         orders,
		users:          users,
		players:        players,
		operationLogs:  operationLogs,
		notifications:  notifications,
		payments:       payments,
		defaultSLAMins: 30,
	}
}

// InitiateDisputeRequest represents a request to initiate a dispute
type InitiateDisputeRequest struct {
	OrderID       uint64
	InitiatorID   uint64
	InitiatorType model.DisputeInitiatorType
	Type          model.DisputeType
	Reason        string
	EvidenceText  string
	EvidenceURLs  []string
}

// InitiateDisputeResponse represents the response after initiating a dispute
type InitiateDisputeResponse struct {
	DisputeID   uint64
	TraceID     string
	SLADeadline *time.Time
}

// InitiateDispute creates a new dispute for an order
func (s *DisputeService) InitiateDispute(ctx context.Context, req InitiateDisputeRequest) (*InitiateDisputeResponse, error) {
	// Validate request
	if req.OrderID == 0 || req.InitiatorID == 0 {
		return nil, ErrDisputeValidation
	}
	if req.Reason == "" {
		return nil, apierr.BadRequest("纠纷原因不能为空")
	}

	// Get order
	order, err := s.orders.Get(ctx, req.OrderID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Verify initiator is related to the order (user or player)
	if req.InitiatorType == model.DisputeInitiatorUser && order.UserID != req.InitiatorID {
		return nil, ErrDisputeUnauthorized
	}
	// For player initiator, check if they are the assigned player
	// InitiatorID is a user ID, so we need to check if the player's UserID matches
	if req.InitiatorType == model.DisputeInitiatorPlayer {
		if order.PlayerID == nil {
			return nil, ErrDisputeUnauthorized
		}
		// If we have a player repository, use it to verify
		if s.players != nil {
			player, err := s.players.Get(ctx, *order.PlayerID)
			if err != nil || player == nil || player.UserID != req.InitiatorID {
				return nil, ErrDisputeUnauthorized
			}
		} else {
			// Fallback: assume InitiatorID is the player record ID (legacy behavior)
			if *order.PlayerID != req.InitiatorID {
				return nil, ErrDisputeUnauthorized
			}
		}
	}

	// Check if dispute can be initiated
	if !model.CanInitiateDispute(order) {
		return nil, ErrCannotInitiateDispute
	}

	// Check if dispute already exists
	existingDispute, err := s.disputes.GetByOrderID(ctx, req.OrderID)
	if err == nil && existingDispute != nil {
		return nil, ErrDisputeExists
	}
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	// Generate trace ID
	traceID := uuid.New().String()

	// Calculate SLA deadline
	slaDeadline := time.Now().Add(time.Duration(s.defaultSLAMins) * time.Minute)

	// Create dispute
	dispute := &model.OrderDispute{
		OrderID:       req.OrderID,
		InitiatorID:   req.InitiatorID,
		InitiatorType: req.InitiatorType,
		Type:          req.Type,
		Status:        model.DisputeStatusPending,
		Reason:        req.Reason,
		EvidenceText:  req.EvidenceText,
		EvidenceURLs:  req.EvidenceURLs,
		SLADeadline:   &slaDeadline,
		TraceID:       traceID,
	}

	if err := s.disputes.Create(ctx, dispute); err != nil {
		return nil, err
	}

	// Update order
	order.HasDispute = true
	order.Status = model.OrderStatusDisputed
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionInitiateDispute, "Dispute initiated", traceID, &req.InitiatorID)

	return &InitiateDisputeResponse{
		DisputeID:   dispute.ID,
		TraceID:     traceID,
		SLADeadline: &slaDeadline,
	}, nil
}

// AssignDisputeRequest represents a request to assign a dispute to a customer service representative
type AssignDisputeRequest struct {
	DisputeID         uint64
	AssignedServiceID uint64  // 分配的独立客服ID
	OriginalServiceID *uint64 // 原客服ID（可选）
	ActorUserID       uint64  // who is making this assignment
}

// AssignDispute assigns a dispute to a customer service representative
func (s *DisputeService) AssignDispute(ctx context.Context, req AssignDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 || req.AssignedServiceID == 0 {
		return ErrDisputeValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute can be assigned
	if dispute.Status != model.DisputeStatusPending {
		return apierr.BadRequest("纠纷状态不是待处理，无法分配")
	}

	// Verify assigned user exists and has appropriate role
	assignedUser, err := s.users.Get(ctx, req.AssignedServiceID)
	if err != nil {
		return err
	}
	if assignedUser == nil {
		return apierr.NotFound("分配的客服不存在")
	}

	// Update dispute
	dispute.Status = model.DisputeStatusAssigned
	dispute.AssignedServiceID = &req.AssignedServiceID
	if req.OriginalServiceID != nil {
		dispute.OriginalServiceID = req.OriginalServiceID
	}

	if err := s.disputes.Update(ctx, dispute); err != nil {
		return err
	}

	// Log operation
	metadata := fmt.Sprintf("Assigned to service user %d", req.AssignedServiceID)
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionAssignDispute, metadata, dispute.TraceID, &req.ActorUserID)

	// Send notification to assigned user
	s.sendNotification(ctx, req.AssignedServiceID, "新争议分配",
		fmt.Sprintf("您已被分配处理争议 #%d", dispute.ID), dispute.TraceID)

	return nil
}

// ResolveDisputeRequest represents a request to resolve a dispute
type ResolveDisputeRequest struct {
	DisputeID     uint64
	Resolution    model.DisputeResolution
	ResolveRemark string
	ActorUserID   uint64 // who is resolving this
}

// ResolveDispute resolves a dispute with a decision
func (s *DisputeService) ResolveDispute(ctx context.Context, req ResolveDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 {
		return ErrDisputeValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute can be resolved
	if dispute.Status == model.DisputeStatusResolved || dispute.Status == model.DisputeStatusRejected || dispute.Status == model.DisputeStatusCanceled {
		return apierr.BadRequest("纠纷已解决、已拒绝或已取消，无法再次处理")
	}

	// Get order
	order, err := s.orders.Get(ctx, dispute.OrderID)
	if err != nil {
		return err
	}

	// Update dispute
	now := time.Now()
	dispute.Resolution = req.Resolution
	dispute.ResolveRemark = req.ResolveRemark
	dispute.ResolvedAt = &now
	dispute.ResolvedBy = &req.ActorUserID

	// Set status based on resolution
	if req.Resolution == model.ResolutionReject {
		dispute.Status = model.DisputeStatusRejected
	} else {
		dispute.Status = model.DisputeStatusResolved
	}

	if err := s.disputes.Update(ctx, dispute); err != nil {
		return err
	}

	// Handle resolution based on decision
	if req.Resolution == model.ResolutionRefund {
		// Full refund
		if err := s.processRefund(ctx, order, dispute, order.TotalPriceCents, &req.ActorUserID); err != nil {
			return err
		}
	} else if req.Resolution == model.ResolutionReject {
		// Restore order status
		order.Status = model.OrderStatusCompleted
		order.HasDispute = false
		if err := s.orders.Update(ctx, order); err != nil {
			return err
		}
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionResolveDispute,
		fmt.Sprintf("Resolved with %s decision", req.Resolution), dispute.TraceID, &req.ActorUserID)

	// Send notification to initiator
	s.sendNotification(ctx, dispute.InitiatorID, "争议处理结果",
		fmt.Sprintf("您的争议 #%d 已处理完成", dispute.ID), dispute.TraceID)

	return nil
}

// RollbackDisputeRequest represents a request to rollback an assignment
type RollbackDisputeRequest struct {
	DisputeID      uint64
	RollbackReason string
	ActorUserID    uint64
}

// RollbackDisputeAssignment rolls back a dispute assignment
func (s *DisputeService) RollbackDisputeAssignment(ctx context.Context, req RollbackDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 {
		return ErrDisputeValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute is assigned
	if dispute.Status != model.DisputeStatusAssigned && dispute.Status != model.DisputeStatusMediating {
		return apierr.BadRequest("纠纷状态不是已分配或调解中，无法回滚分配")
	}

	// Update dispute
	now := time.Now()
	dispute.Status = model.DisputeStatusPending
	dispute.AssignedServiceID = nil
	dispute.OriginalServiceID = nil
	dispute.RolledBackAt = &now
	dispute.RolledBackByUserID = &req.ActorUserID
	dispute.RollbackReason = req.RollbackReason

	if err := s.disputes.Update(ctx, dispute); err != nil {
		return err
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionRollbackDispute,
		fmt.Sprintf("Rolled back: %s", req.RollbackReason), dispute.TraceID, &req.ActorUserID)

	return nil
}

// CheckAndMarkSLABreaches checks for disputes that have breached SLA and marks them
func (s *DisputeService) CheckAndMarkSLABreaches(ctx context.Context) error {
	breachedDisputes, err := s.disputes.ListSLABreached(ctx)
	if err != nil {
		return err
	}

	for _, dispute := range breachedDisputes {
		if err := s.disputes.MarkSLABreached(ctx, dispute.ID); err != nil {
			continue // Log but continue processing others
		}

		// Log operation
		s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionUpdateStatus,
			"SLA breached", dispute.TraceID, dispute.AssignedServiceID)

		// Send alert notification to assigned service user
		if dispute.AssignedServiceID != nil {
			s.sendNotification(ctx, *dispute.AssignedServiceID, "SLA 超时警告",
				fmt.Sprintf("争议 #%d 已超过 SLA 截止时间", dispute.ID), dispute.TraceID)
		}
	}

	return nil
}

// GetDisputeDetail retrieves detailed information about a dispute
func (s *DisputeService) GetDisputeDetail(ctx context.Context, disputeID uint64) (*model.OrderDispute, error) {
	return s.disputes.Get(ctx, disputeID)
}

// ListPendingDisputes lists disputes pending assignment
func (s *DisputeService) ListPendingDisputes(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	return s.disputes.ListPendingAssignment(ctx, page, pageSize)
}

// ListDisputesByStatus lists disputes filtered by status
func (s *DisputeService) ListDisputesByStatus(ctx context.Context, statuses []model.DisputeStatus, page, pageSize int) ([]model.OrderDispute, int64, error) {
	opts := repository.DisputeListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
	}
	return s.disputes.List(ctx, opts)
}

// ListDisputesRequest represents a request to list disputes
type ListDisputesRequest struct {
	Page     int
	PageSize int
	Status   string
	OrderNo  string
}

// ListDisputes lists disputes with optional filters
func (s *DisputeService) ListDisputes(ctx context.Context, req ListDisputesRequest) ([]model.OrderDispute, int64, error) {
	opts := repository.DisputeListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Add status filter if provided
	if req.Status != "" {
		opts.Statuses = []model.DisputeStatus{model.DisputeStatus(req.Status)}
	}

	// Add order number filter if provided
	if req.OrderNo != "" {
		opts.OrderNo = req.OrderNo
	}

	return s.disputes.List(ctx, opts)
}

// GetDisputeStats returns dispute statistics by status
func (s *DisputeService) GetDisputeStats(ctx context.Context) (map[string]int64, error) {
	return s.disputes.GetStats(ctx)
}

// Helper functions

func (s *DisputeService) processRefund(ctx context.Context, order *model.Order, dispute *model.OrderDispute, amount int64, actorID *uint64) error {
	// Update order status
	order.Status = model.OrderStatusRefunded
	order.RefundAmountCents = amount
	order.RefundReason = fmt.Sprintf("争议处理退款: %s", dispute.ResolveRemark)
	now := time.Now()
	order.RefundedAt = &now

	if err := s.orders.Update(ctx, order); err != nil {
		return err
	}

	// Log refund operation
	s.logOperation(ctx, model.OpEntityOrder, order.ID, model.OpActionRefund,
		fmt.Sprintf("Refund processed: %d cents", amount), dispute.TraceID, actorID)

	return nil
}

func (s *DisputeService) logOperation(ctx context.Context, entityType model.OperationEntityType, entityID uint64, action model.OperationAction, reason string, traceID string, actorID *uint64) {
	log := &model.OperationLog{
		EntityType:  string(entityType),
		EntityID:    entityID,
		Action:      string(action),
		Reason:      reason,
		TraceID:     traceID,
		ActorUserID: actorID,
	}

	// Try to add trace ID to metadata if possible
	// For now, we'll just log the operation
	_ = s.operationLogs.Append(ctx, log)
}

func (s *DisputeService) sendNotification(ctx context.Context, userID uint64, title, message, traceID string) {
	event := &model.NotificationEvent{
		UserID:   userID,
		Title:    title,
		Message:  message,
		Channel:  "web",
		Priority: model.NotificationPriorityHigh,
		ReadAt:   nil,
	}

	_ = s.notifications.Create(ctx, event)
}

// ============================================================================
// Batch Operations
// ============================================================================

// BatchOperationResult represents the result of a batch operation
type BatchOperationResult struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	SuccessCount int                    `json:"successCount"`
	FailedCount  int                    `json:"failedCount"`
	Errors       []BatchOperationError  `json:"errors,omitempty"`
}

// BatchOperationError represents an error that occurred during batch operation
type BatchOperationError struct {
	DisputeID uint64 `json:"disputeId"`
	Error     string `json:"error"`
}

// BatchAssignDisputesRequest represents a request to batch assign disputes
type BatchAssignDisputesRequest struct {
	DisputeIDs         []uint64 `json:"disputeIds" binding:"required,min=1,max=100"`
	AssignedServiceID  uint64   `json:"assignedServiceId" binding:"required"`
	OriginalServiceID  *uint64  `json:"originalServiceId,omitempty"`
	ActorUserID        uint64   `json:"actorUserId"`
}

// BatchAssignDisputes assigns multiple disputes to a customer service representative
func (s *DisputeService) BatchAssignDisputes(ctx context.Context, req BatchAssignDisputesRequest) (*BatchOperationResult, error) {
	// Validate request
	if len(req.DisputeIDs) == 0 {
		return nil, ErrDisputeValidation
	}
	if len(req.DisputeIDs) > 100 {
		return nil, apierr.BadRequest("批量操作最多支持100个纠纷")
	}
	if req.AssignedServiceID == 0 {
		return nil, ErrDisputeValidation
	}

	// Verify assigned user exists
	assignedUser, err := s.users.Get(ctx, req.AssignedServiceID)
	if err != nil {
		return nil, err
	}
	if assignedUser == nil {
		return nil, apierr.NotFound("分配的客服不存在")
	}

	result := &BatchOperationResult{
		Success:      true,
		SuccessCount: 0,
		FailedCount:  0,
		Errors:       make([]BatchOperationError, 0),
	}

	// Process each dispute
	for _, disputeID := range req.DisputeIDs {
		// Get dispute
		dispute, err := s.disputes.Get(ctx, disputeID)
		if err != nil {
			if err == repository.ErrNotFound {
				result.FailedCount++
				result.Errors = append(result.Errors, BatchOperationError{
					DisputeID: disputeID,
					Error:     "纠纷不存在",
				})
				continue
			}
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("获取纠纷失败: %v", err),
			})
			continue
		}

		// Check if dispute can be assigned
		if dispute.Status != model.DisputeStatusPending {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("纠纷状态为%s，无法分配", dispute.Status),
			})
			continue
		}

		// Update dispute
		dispute.Status = model.DisputeStatusAssigned
		dispute.AssignedServiceID = &req.AssignedServiceID
		if req.OriginalServiceID != nil {
			dispute.OriginalServiceID = req.OriginalServiceID
		}

		if err := s.disputes.Update(ctx, dispute); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("更新失败: %v", err),
			})
			continue
		}

		// Log operation
		s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionAssignDispute,
			fmt.Sprintf("Batch assigned to service user %d", req.AssignedServiceID), dispute.TraceID, &req.ActorUserID)

		// Send notification
		s.sendNotification(ctx, req.AssignedServiceID, "新争议分配",
			fmt.Sprintf("您已被分配处理争议 #%d", dispute.ID), dispute.TraceID)

		result.SuccessCount++
	}

	// Generate message
	if result.FailedCount == 0 {
		result.Message = fmt.Sprintf("批量分配成功，共%d个纠纷", result.SuccessCount)
	} else if result.SuccessCount == 0 {
		result.Success = false
		result.Message = "批量分配失败"
	} else {
		result.Message = fmt.Sprintf("批量分配完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	return result, nil
}

// BatchUpdateDisputesStatusRequest represents a request to batch update dispute status
type BatchUpdateDisputesStatusRequest struct {
	DisputeIDs  []uint64             `json:"disputeIds" binding:"required,min=1,max=100"`
	Status      model.DisputeStatus  `json:"status" binding:"required"`
	ActorUserID uint64               `json:"actorUserId"`
}

// BatchUpdateDisputesStatus updates status for multiple disputes
func (s *DisputeService) BatchUpdateDisputesStatus(ctx context.Context, req BatchUpdateDisputesStatusRequest) (*BatchOperationResult, error) {
	// Validate request
	if len(req.DisputeIDs) == 0 {
		return nil, ErrDisputeValidation
	}
	if len(req.DisputeIDs) > 100 {
		return nil, apierr.BadRequest("批量操作最多支持100个纠纷")
	}

	// Validate status transition
	validStatuses := map[model.DisputeStatus]bool{
		model.DisputeStatusAssigned:  true,
		model.DisputeStatusMediating: true,
		model.DisputeStatusCanceled:  true,
	}
	if !validStatuses[req.Status] {
		return nil, apierr.BadRequest("批量更新不支持的目标状态")
	}

	result := &BatchOperationResult{
		Success:      true,
		SuccessCount: 0,
		FailedCount:  0,
		Errors:       make([]BatchOperationError, 0),
	}

	// Process each dispute
	for _, disputeID := range req.DisputeIDs {
		// Get dispute
		dispute, err := s.disputes.Get(ctx, disputeID)
		if err != nil {
			if err == repository.ErrNotFound {
				result.FailedCount++
				result.Errors = append(result.Errors, BatchOperationError{
					DisputeID: disputeID,
					Error:     "纠纷不存在",
				})
				continue
			}
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("获取纠纷失败: %v", err),
			})
			continue
		}

		// Validate status transition
		if !canTransitionTo(dispute.Status, req.Status) {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("无法从%s转换到%s", dispute.Status, req.Status),
			})
			continue
		}

		// Update dispute
		oldStatus := dispute.Status
		dispute.Status = req.Status

		if err := s.disputes.Update(ctx, dispute); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("更新失败: %v", err),
			})
			continue
		}

		// Log operation
		s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionUpdateStatus,
			fmt.Sprintf("Batch status updated: %s -> %s", oldStatus, req.Status), dispute.TraceID, &req.ActorUserID)

		result.SuccessCount++
	}

	// Generate message
	if result.FailedCount == 0 {
		result.Message = fmt.Sprintf("批量更新状态成功，共%d个纠纷", result.SuccessCount)
	} else if result.SuccessCount == 0 {
		result.Success = false
		result.Message = "批量更新状态失败"
	} else {
		result.Message = fmt.Sprintf("批量更新状态完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	return result, nil
}

// BatchCloseDisputesRequest represents a request to batch close disputes
type BatchCloseDisputesRequest struct {
	DisputeIDs    []uint64                 `json:"disputeIds" binding:"required,min=1,max=100"`
	Resolution    model.DisputeResolution  `json:"resolution" binding:"required"`
	ResolveRemark string                   `json:"resolveRemark" binding:"required"`
	ActorUserID   uint64                   `json:"actorUserId"`
}

// BatchCloseDisputes closes multiple disputes with resolution
func (s *DisputeService) BatchCloseDisputes(ctx context.Context, req BatchCloseDisputesRequest) (*BatchOperationResult, error) {
	// Validate request
	if len(req.DisputeIDs) == 0 {
		return nil, ErrDisputeValidation
	}
	if len(req.DisputeIDs) > 100 {
		return nil, apierr.BadRequest("批量操作最多支持100个纠纷")
	}
	if req.ResolveRemark == "" {
		return nil, apierr.BadRequest("处理备注不能为空")
	}

	// Validate resolution
	validResolutions := map[model.DisputeResolution]bool{
		model.ResolutionRefund:  true,
		model.ResolutionPartial: true,
		model.ResolutionReject:  true,
	}
	if !validResolutions[req.Resolution] {
		return nil, apierr.BadRequest("批量关闭不支持的处理决定")
	}

	result := &BatchOperationResult{
		Success:      true,
		SuccessCount: 0,
		FailedCount:  0,
		Errors:       make([]BatchOperationError, 0),
	}

	// Process each dispute
	for _, disputeID := range req.DisputeIDs {
		// Get dispute
		dispute, err := s.disputes.Get(ctx, disputeID)
		if err != nil {
			if err == repository.ErrNotFound {
				result.FailedCount++
				result.Errors = append(result.Errors, BatchOperationError{
					DisputeID: disputeID,
					Error:     "纠纷不存在",
				})
				continue
			}
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("获取纠纷失败: %v", err),
			})
			continue
		}

		// Check if dispute can be resolved
		if dispute.Status == model.DisputeStatusResolved || dispute.Status == model.DisputeStatusRejected || dispute.Status == model.DisputeStatusCanceled {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("纠纷状态为%s，无法再次处理", dispute.Status),
			})
			continue
		}

		// Get order for refund processing
		order, err := s.orders.Get(ctx, dispute.OrderID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("获取订单失败: %v", err),
			})
			continue
		}

		// Update dispute
		now := time.Now()
		dispute.Resolution = req.Resolution
		dispute.ResolveRemark = req.ResolveRemark
		dispute.ResolvedAt = &now
		dispute.ResolvedBy = &req.ActorUserID

		// Set status based on resolution
		if req.Resolution == model.ResolutionReject {
			dispute.Status = model.DisputeStatusRejected
		} else {
			dispute.Status = model.DisputeStatusResolved
		}

		if err := s.disputes.Update(ctx, dispute); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BatchOperationError{
				DisputeID: disputeID,
				Error:     fmt.Sprintf("更新纠纷失败: %v", err),
			})
			continue
		}

		// Handle resolution based on decision
		if req.Resolution == model.ResolutionRefund {
			// Full refund
			if err := s.processRefund(ctx, order, dispute, order.TotalPriceCents, &req.ActorUserID); err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, BatchOperationError{
					DisputeID: disputeID,
					Error:     fmt.Sprintf("处理退款失败: %v", err),
				})
				continue
			}
		} else if req.Resolution == model.ResolutionReject {
			// Restore order status
			order.Status = model.OrderStatusCompleted
			order.HasDispute = false
			if err := s.orders.Update(ctx, order); err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, BatchOperationError{
					DisputeID: disputeID,
					Error:     fmt.Sprintf("恢复订单状态失败: %v", err),
				})
				continue
			}
		}

		// Log operation
		s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionResolveDispute,
			fmt.Sprintf("Batch resolved with %s decision", req.Resolution), dispute.TraceID, &req.ActorUserID)

		// Send notification to initiator
		s.sendNotification(ctx, dispute.InitiatorID, "争议处理结果",
			fmt.Sprintf("您的争议 #%d 已处理完成", dispute.ID), dispute.TraceID)

		result.SuccessCount++
	}

	// Generate message
	if result.FailedCount == 0 {
		result.Message = fmt.Sprintf("批量关闭成功，共%d个纠纷", result.SuccessCount)
	} else if result.SuccessCount == 0 {
		result.Success = false
		result.Message = "批量关闭失败"
	} else {
		result.Message = fmt.Sprintf("批量关闭完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	return result, nil
}

// canTransitionTo checks if a status transition is valid
func canTransitionTo(from, to model.DisputeStatus) bool {
	validTransitions := map[model.DisputeStatus][]model.DisputeStatus{
		model.DisputeStatusPending: {
			model.DisputeStatusAssigned,
			model.DisputeStatusCanceled,
		},
		model.DisputeStatusAssigned: {
			model.DisputeStatusMediating,
			model.DisputeStatusCanceled,
		},
		model.DisputeStatusMediating: {
			model.DisputeStatusAssigned,
			model.DisputeStatusCanceled,
		},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
