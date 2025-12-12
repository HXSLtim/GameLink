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

// InitiateDisputeRequest represents a request to initiate a dispute
type InitiateDisputeRequest struct {
	OrderID      uint64
	UserID       uint64
	Reason       string
	Description  string
	EvidenceURLs []string
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
	if req.OrderID == 0 || req.UserID == 0 {
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

	// Verify order belongs to user
	if order.UserID != req.UserID {
		return nil, ErrDisputeUnauthorized
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
		OrderID:      req.OrderID,
		UserID:       req.UserID,
		Status:       model.DisputeStatusPending,
		Reason:       req.Reason,
		Description:  req.Description,
		EvidenceURLs: req.EvidenceURLs,
		SLADeadline:  &slaDeadline,
		TraceID:      traceID,
	}

	if err := s.disputes.Create(ctx, dispute); err != nil {
		return nil, err
	}

	// Update order
	order.HasDispute = true
	// Note: DisputeID field removed from Order model to avoid circular dependency
	// The relationship is maintained through OrderDispute.OrderID
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionInitiateDispute, "User initiated dispute", traceID, &req.UserID)

	return &InitiateDisputeResponse{
		DisputeID:   dispute.ID,
		TraceID:     traceID,
		SLADeadline: &slaDeadline,
	}, nil
}

// AssignDisputeRequest represents a request to assign a dispute to a customer service representative
type AssignDisputeRequest struct {
	DisputeID        uint64
	AssignedToUserID uint64
	Source           model.AssignmentSource
	ActorUserID      uint64 // who is making this assignment
}

// AssignDispute assigns a dispute to a customer service representative
func (s *DisputeService) AssignDispute(ctx context.Context, req AssignDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 || req.AssignedToUserID == 0 {
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
	assignedUser, err := s.users.Get(ctx, req.AssignedToUserID)
	if err != nil {
		return err
	}
	if assignedUser == nil {
		return apierr.NotFound("分配的用户不存在")
	}

	// Update dispute
	now := time.Now()
	dispute.Status = model.DisputeStatusAssigned
	dispute.AssignedToUserID = &req.AssignedToUserID
	dispute.AssignmentSource = req.Source
	dispute.AssignedAt = &now

	if err := s.disputes.Update(ctx, dispute); err != nil {
		return err
	}

	// Log operation
	metadata := fmt.Sprintf("Assigned to user %d via %s", req.AssignedToUserID, req.Source)
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionAssignDispute, metadata, dispute.TraceID, &req.ActorUserID)

	// Send notification to assigned user
	s.sendNotification(ctx, req.AssignedToUserID, "New Dispute Assignment",
		fmt.Sprintf("You have been assigned dispute #%d", dispute.ID), dispute.TraceID)

	return nil
}

// ResolveDisputeRequest represents a request to resolve a dispute
type ResolveDisputeRequest struct {
	DisputeID        uint64
	Resolution       model.DisputeResolution
	ResolutionAmount int64 // in cents
	ResolutionNotes  string
	ActorUserID      uint64 // who is resolving this
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
	dispute.Status = model.DisputeStatusResolved
	dispute.Resolution = req.Resolution
	dispute.ResolutionAmount = req.ResolutionAmount
	dispute.ResolutionNotes = req.ResolutionNotes
	dispute.ResolvedAt = &now
	dispute.ResolvedByUserID = &req.ActorUserID

	if err := s.disputes.Update(ctx, dispute); err != nil {
		return err
	}

	// Handle resolution based on decision
	if req.Resolution == model.ResolutionRefund {
		if err := s.processRefund(ctx, order, dispute, req.ResolutionAmount, &req.ActorUserID); err != nil {
			return err
		}
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionResolveDispute,
		fmt.Sprintf("Resolved with %s decision", req.Resolution), dispute.TraceID, &req.ActorUserID)

	// Send notification to user
	s.sendNotification(ctx, dispute.UserID, "Dispute Resolved",
		fmt.Sprintf("Your dispute #%d has been resolved", dispute.ID), dispute.TraceID)

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
	dispute.AssignedToUserID = nil
	dispute.AssignmentSource = ""
	dispute.AssignedAt = nil
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
			"SLA breached", dispute.TraceID, dispute.AssignedToUserID)

		// Send alert notification
		s.sendNotification(ctx, *dispute.AssignedToUserID, "SLA Breached",
			fmt.Sprintf("Dispute #%d has exceeded SLA deadline", dispute.ID), dispute.TraceID)
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

// Helper functions

func (s *DisputeService) processRefund(ctx context.Context, order *model.Order, dispute *model.OrderDispute, amount int64, actorID *uint64) error {
	// Update order status
	order.Status = model.OrderStatusRefunded
	order.RefundAmountCents = amount
	order.RefundReason = fmt.Sprintf("Dispute resolution: %s", dispute.ResolutionNotes)
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
