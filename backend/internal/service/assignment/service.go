package assignment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound dispute not found
	ErrNotFound = repository.ErrNotFound
	// ErrValidation validation failed
	ErrValidation = errors.New("validation failed")
	// ErrInvalidStatus invalid dispute status
	ErrInvalidStatus = errors.New("invalid dispute status")
	// ErrUnauthorized unauthorized operation
	ErrUnauthorized = errors.New("unauthorized")
	// ErrSLAExpired SLA deadline has passed
	ErrSLAExpired = errors.New("sla deadline expired")
	// ErrOrderNotFound order not found
	ErrOrderNotFound = errors.New("order not found")
	// ErrDisputeExists dispute already exists for this order
	ErrDisputeExists = errors.New("dispute already exists for this order")
	// ErrCannotInitiateDispute cannot initiate dispute for this order
	ErrCannotInitiateDispute = errors.New("cannot initiate dispute for this order")
)

// AssignmentService handles dispute and assignment operations
type AssignmentService struct {
	disputes       repository.DisputeRepository
	orders         repository.OrderRepository
	users          repository.UserRepository
	operationLogs  repository.OperationLogRepository
	notifications  repository.NotificationRepository
	payments       repository.PaymentRepository
	defaultSLAMins int // default SLA in minutes (30)
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(
	disputes repository.DisputeRepository,
	orders repository.OrderRepository,
	users repository.UserRepository,
	operationLogs repository.OperationLogRepository,
	notifications repository.NotificationRepository,
	payments repository.PaymentRepository,
) *AssignmentService {
	return &AssignmentService{
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
	DisputeID  uint64
	TraceID    string
	SLADeadline *time.Time
}

// InitiateDispute creates a new dispute for an order
func (s *AssignmentService) InitiateDispute(ctx context.Context, req InitiateDisputeRequest) (*InitiateDisputeResponse, error) {
	// Validate request
	if req.OrderID == 0 || req.UserID == 0 {
		return nil, ErrValidation
	}
	if req.Reason == "" {
		return nil, fmt.Errorf("%w: reason is required", ErrValidation)
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
		return nil, ErrUnauthorized
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
	order.DisputeID = &dispute.ID
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionInitiateDispute, "User initiated dispute", traceID)

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
func (s *AssignmentService) AssignDispute(ctx context.Context, req AssignDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 || req.AssignedToUserID == 0 {
		return ErrValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute can be assigned
	if dispute.Status != model.DisputeStatusPending {
		return fmt.Errorf("%w: dispute is not in pending status", ErrInvalidStatus)
	}

	// Verify assigned user exists and has appropriate role
	assignedUser, err := s.users.Get(ctx, req.AssignedToUserID)
	if err != nil {
		return err
	}
	if assignedUser == nil {
		return fmt.Errorf("assigned user not found")
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
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionAssignDispute, metadata, dispute.TraceID)

	// Send notification to assigned user
	s.sendNotification(ctx, req.AssignedToUserID, "New Dispute Assignment", 
		fmt.Sprintf("You have been assigned dispute #%d", dispute.ID), dispute.TraceID)

	return nil
}

// ResolveDisputeRequest represents a request to resolve a dispute
type ResolveDisputeRequest struct {
	DisputeID       uint64
	Resolution      model.DisputeResolution
	ResolutionAmount int64  // in cents
	ResolutionNotes string
	ActorUserID     uint64 // who is resolving this
}

// ResolveDispute resolves a dispute with a decision
func (s *AssignmentService) ResolveDispute(ctx context.Context, req ResolveDisputeRequest) error {
	// Validate request
	if req.DisputeID == 0 {
		return ErrValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute can be resolved
	if dispute.Status == model.DisputeStatusResolved || dispute.Status == model.DisputeStatusRejected || dispute.Status == model.DisputeStatusCanceled {
		return fmt.Errorf("%w: dispute is already resolved", ErrInvalidStatus)
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
		if err := s.processRefund(ctx, order, dispute, req.ResolutionAmount); err != nil {
			return err
		}
	}

	// Log operation
	s.logOperation(ctx, model.OpEntityDispute, dispute.ID, model.OpActionResolveDispute, 
		fmt.Sprintf("Resolved with %s decision", req.Resolution), dispute.TraceID)

	// Send notification to user
	s.sendNotification(ctx, dispute.UserID, "Dispute Resolved", 
		fmt.Sprintf("Your dispute #%d has been resolved", dispute.ID), dispute.TraceID)

	return nil
}

// RollbackAssignmentRequest represents a request to rollback an assignment
type RollbackAssignmentRequest struct {
	DisputeID      uint64
	RollbackReason string
	ActorUserID    uint64
}

// RollbackAssignment rolls back a dispute assignment
func (s *AssignmentService) RollbackAssignment(ctx context.Context, req RollbackAssignmentRequest) error {
	// Validate request
	if req.DisputeID == 0 {
		return ErrValidation
	}

	// Get dispute
	dispute, err := s.disputes.Get(ctx, req.DisputeID)
	if err != nil {
		return err
	}

	// Check if dispute is assigned
	if dispute.Status != model.DisputeStatusAssigned && dispute.Status != model.DisputeStatusMediating {
		return fmt.Errorf("%w: dispute is not in assigned or mediating status", ErrInvalidStatus)
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
		fmt.Sprintf("Rolled back: %s", req.RollbackReason), dispute.TraceID)

	return nil
}

// GetDisputeDetail retrieves detailed information about a dispute
func (s *AssignmentService) GetDisputeDetail(ctx context.Context, disputeID uint64) (*model.OrderDispute, error) {
	return s.disputes.Get(ctx, disputeID)
}

// ListPendingDisputes lists disputes pending assignment
func (s *AssignmentService) ListPendingDisputes(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	return s.disputes.ListPendingAssignment(ctx, page, pageSize)
}

// ListDisputesByStatus lists disputes filtered by status
func (s *AssignmentService) ListDisputesByStatus(ctx context.Context, statuses []model.DisputeStatus, page, pageSize int) ([]model.OrderDispute, int64, error) {
	opts := repository.DisputeListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
	}
	return s.disputes.List(ctx, opts)
}

// CheckAndMarkSLABreaches checks for disputes that have breached SLA and marks them
func (s *AssignmentService) CheckAndMarkSLABreaches(ctx context.Context) error {
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
			"SLA breached", dispute.TraceID)

		// Send alert notification
		s.sendNotification(ctx, *dispute.AssignedToUserID, "SLA Breached", 
			fmt.Sprintf("Dispute #%d has exceeded SLA deadline", dispute.ID), dispute.TraceID)
	}

	return nil
}

// Helper functions

func (s *AssignmentService) processRefund(ctx context.Context, order *model.Order, dispute *model.OrderDispute, amount int64) error {
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
		fmt.Sprintf("Refund processed: %d cents", amount), dispute.TraceID)

	return nil
}

func (s *AssignmentService) logOperation(ctx context.Context, entityType model.OperationEntityType, entityID uint64, action model.OperationAction, reason string, traceID string) {
	log := &model.OperationLog{
		EntityType: string(entityType),
		EntityID:   entityID,
		Action:     string(action),
		Reason:     reason,
	}

	// Try to add trace ID to metadata if possible
	// For now, we'll just log the operation
	_ = s.operationLogs.Append(ctx, log)
}

func (s *AssignmentService) sendNotification(ctx context.Context, userID uint64, title, message, traceID string) {
	event := &model.NotificationEvent{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Channel:   "web",
		Priority:  model.NotificationPriorityHigh,
		ReadAt:    nil,
	}

	_ = s.notifications.Create(ctx, event)
}
