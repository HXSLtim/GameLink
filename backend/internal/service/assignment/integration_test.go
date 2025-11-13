package assignment

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Mock repositories for testing
type mockDisputeRepository struct {
	disputes map[uint64]*model.OrderDispute
	nextID   uint64
}

func newMockDisputeRepository() *mockDisputeRepository {
	return &mockDisputeRepository{
		disputes: make(map[uint64]*model.OrderDispute),
		nextID:   1,
	}
}

func (m *mockDisputeRepository) Create(ctx context.Context, dispute *model.OrderDispute) error {
	dispute.ID = m.nextID
	m.nextID++
	m.disputes[dispute.ID] = dispute
	return nil
}

func (m *mockDisputeRepository) Get(ctx context.Context, id uint64) (*model.OrderDispute, error) {
	if d, ok := m.disputes[id]; ok {
		return d, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockDisputeRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
	for _, d := range m.disputes {
		if d.OrderID == orderID {
			return d, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockDisputeRepository) Update(ctx context.Context, dispute *model.OrderDispute) error {
	m.disputes[dispute.ID] = dispute
	return nil
}

func (m *mockDisputeRepository) List(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
	var result []model.OrderDispute
	for _, d := range m.disputes {
		result = append(result, *d)
	}
	return result, int64(len(result)), nil
}

func (m *mockDisputeRepository) ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	var result []model.OrderDispute
	for _, d := range m.disputes {
		if d.Status == model.DisputeStatusPending && d.AssignedToUserID == nil {
			result = append(result, *d)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockDisputeRepository) ListSLABreached(ctx context.Context) ([]model.OrderDispute, error) {
	var result []model.OrderDispute
	for _, d := range m.disputes {
		if !d.SLABreached && d.SLADeadline != nil && time.Now().After(*d.SLADeadline) {
			result = append(result, *d)
		}
	}
	return result, nil
}

func (m *mockDisputeRepository) MarkSLABreached(ctx context.Context, disputeID uint64) error {
	if d, ok := m.disputes[disputeID]; ok {
		d.SLABreached = true
		now := time.Now()
		d.SLABreachedAt = &now
	}
	return nil
}

func (m *mockDisputeRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.disputes, id)
	return nil
}

func (m *mockDisputeRepository) CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error) {
	count := 0
	for _, d := range m.disputes {
		if d.Status == status {
			count++
		}
	}
	return int64(count), nil
}

func (m *mockDisputeRepository) GetPendingCount(ctx context.Context) (int64, error) {
	count := 0
	for _, d := range m.disputes {
		if d.Status == model.DisputeStatusPending {
			count++
		}
	}
	return int64(count), nil
}

// Mock other repositories
type mockOrderRepository struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

type mockUserRepository struct {
	users map[uint64]*model.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uint64]*model.User),
	}
}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockOperationLogRepository struct{}

func (m *mockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	return nil
}

func (m *mockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type mockNotificationRepository struct{}

func (m *mockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}

func (m *mockNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return nil
}

func (m *mockNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

func (m *mockNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	return nil
}

type mockPaymentRepository struct{}

func (m *mockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}

func (m *mockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

// Integration Tests

func TestInitiateDisputeFullFlow(t *testing.T) {
	ctx := context.Background()

	// Setup
	disputeRepo := newMockDisputeRepository()
	orderRepo := newMockOrderRepository()
	userRepo := newMockUserRepository()
	opLogRepo := &mockOperationLogRepository{}
	notifRepo := &mockNotificationRepository{}
	paymentRepo := &mockPaymentRepository{}

	svc := NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notifRepo, paymentRepo)

	// Create test order
	now := time.Now()
	order := &model.Order{
		Base:     model.Base{ID: 1, CreatedAt: now},
		OrderNo:  "ORD001",
		UserID:   1,
		Status:   model.OrderStatusCompleted,
		CompletedAt: &now,
	}
	orderRepo.Create(ctx, order)

	// Create test user
	user := &model.User{
		Base: model.Base{ID: 1},
		Phone: "13800000001",
	}
	userRepo.users[1] = user

	// Test: Initiate dispute
	resp, err := svc.InitiateDispute(ctx, InitiateDisputeRequest{
		OrderID:      1,
		UserID:       1,
		Reason:       "Service not provided",
		Description:  "Player did not show up",
		EvidenceURLs: []string{"https://example.com/screenshot.jpg"},
	})

	if err != nil {
		t.Fatalf("InitiateDispute failed: %v", err)
	}

	if resp.DisputeID == 0 {
		t.Fatal("Expected non-zero dispute ID")
	}

	if resp.TraceID == "" {
		t.Fatal("Expected non-empty trace ID")
	}

	if resp.SLADeadline == nil {
		t.Fatal("Expected SLA deadline")
	}

	// Verify dispute created
	dispute, err := disputeRepo.Get(ctx, resp.DisputeID)
	if err != nil {
		t.Fatalf("Failed to get dispute: %v", err)
	}

	if dispute.Status != model.DisputeStatusPending {
		t.Errorf("Expected pending status, got %s", dispute.Status)
	}

	if dispute.OrderID != 1 {
		t.Errorf("Expected order ID 1, got %d", dispute.OrderID)
	}

	// Verify order updated
	updatedOrder, _ := orderRepo.Get(ctx, 1)
	if !updatedOrder.HasDispute {
		t.Fatal("Expected order to have dispute flag set")
	}

	t.Log("✓ InitiateDispute test passed")
}

func TestAssignAndResolveDispute(t *testing.T) {
	ctx := context.Background()

	// Setup
	disputeRepo := newMockDisputeRepository()
	orderRepo := newMockOrderRepository()
	userRepo := newMockUserRepository()
	opLogRepo := &mockOperationLogRepository{}
	notifRepo := &mockNotificationRepository{}
	paymentRepo := &mockPaymentRepository{}

	svc := NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notifRepo, paymentRepo)

	// Create test data
	now := time.Now()
	order := &model.Order{
		Base:        model.Base{ID: 1, CreatedAt: now},
		OrderNo:     "ORD001",
		UserID:      1,
		Status:      model.OrderStatusCompleted,
		CompletedAt: &now,
		TotalPriceCents: 10000,
	}
	orderRepo.Create(ctx, order)

	csRep := &model.User{Base: model.Base{ID: 2}, Phone: "13800000002"}
	userRepo.users[2] = csRep

	// Create dispute
	slaDeadline := time.Now().Add(30 * time.Minute)
	dispute := &model.OrderDispute{
		Base:        model.Base{ID: 1},
		OrderID:     1,
		UserID:      1,
		Status:      model.DisputeStatusPending,
		Reason:      "Service not provided",
		SLADeadline: &slaDeadline,
		TraceID:     "trace-001",
	}
	disputeRepo.Create(ctx, dispute)

	// Test: Assign dispute
	err := svc.AssignDispute(ctx, AssignDisputeRequest{
		DisputeID:        1,
		AssignedToUserID: 2,
		Source:           model.AssignmentSourceSystem,
		ActorUserID:      999,
	})

	if err != nil {
		t.Fatalf("AssignDispute failed: %v", err)
	}

	// Verify assignment
	assignedDispute, _ := disputeRepo.Get(ctx, 1)
	if assignedDispute.Status != model.DisputeStatusAssigned {
		t.Errorf("Expected assigned status, got %s", assignedDispute.Status)
	}

	if assignedDispute.AssignedToUserID == nil || *assignedDispute.AssignedToUserID != 2 {
		t.Fatal("Expected dispute to be assigned to user 2")
	}

	// Test: Resolve dispute
	err = svc.ResolveDispute(ctx, ResolveDisputeRequest{
		DisputeID:        1,
		Resolution:       model.ResolutionRefund,
		ResolutionAmount: 10000,
		ResolutionNotes:  "Full refund approved",
		ActorUserID:      2,
	})

	if err != nil {
		t.Fatalf("ResolveDispute failed: %v", err)
	}

	// Verify resolution
	resolvedDispute, _ := disputeRepo.Get(ctx, 1)
	if resolvedDispute.Status != model.DisputeStatusResolved {
		t.Errorf("Expected resolved status, got %s", resolvedDispute.Status)
	}

	if resolvedDispute.Resolution != model.ResolutionRefund {
		t.Errorf("Expected refund resolution, got %s", resolvedDispute.Resolution)
	}

	// Verify order updated
	refundedOrder, _ := orderRepo.Get(ctx, 1)
	if refundedOrder.Status != model.OrderStatusRefunded {
		t.Errorf("Expected refunded order status, got %s", refundedOrder.Status)
	}

	t.Log("✓ AssignAndResolveDispute test passed")
}

func TestRollbackAssignment(t *testing.T) {
	ctx := context.Background()

	// Setup
	disputeRepo := newMockDisputeRepository()
	orderRepo := newMockOrderRepository()
	userRepo := newMockUserRepository()
	opLogRepo := &mockOperationLogRepository{}
	notifRepo := &mockNotificationRepository{}
	paymentRepo := &mockPaymentRepository{}

	svc := NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notifRepo, paymentRepo)

	// Create assigned dispute
	now := time.Now()
	userID := uint64(2)
	slaDeadline := time.Now().Add(30 * time.Minute)
	dispute := &model.OrderDispute{
		Base:             model.Base{ID: 1},
		OrderID:          1,
		UserID:           1,
		Status:           model.DisputeStatusAssigned,
		AssignedToUserID: &userID,
		AssignmentSource: model.AssignmentSourceSystem,
		AssignedAt:       &now,
		SLADeadline:      &slaDeadline,
		TraceID:          "trace-001",
	}
	disputeRepo.Create(ctx, dispute)

	// Test: Rollback assignment
	err := svc.RollbackAssignment(ctx, RollbackAssignmentRequest{
		DisputeID:      1,
		RollbackReason: "Assigned user unavailable",
		ActorUserID:    999,
	})

	if err != nil {
		t.Fatalf("RollbackAssignment failed: %v", err)
	}

	// Verify rollback
	rolledBackDispute, _ := disputeRepo.Get(ctx, 1)
	if rolledBackDispute.Status != model.DisputeStatusPending {
		t.Errorf("Expected pending status after rollback, got %s", rolledBackDispute.Status)
	}

	if rolledBackDispute.AssignedToUserID != nil {
		t.Fatal("Expected assignment to be cleared")
	}

	if rolledBackDispute.RolledBackAt == nil {
		t.Fatal("Expected rollback timestamp")
	}

	t.Log("✓ RollbackAssignment test passed")
}

func TestSLABreachDetection(t *testing.T) {
	ctx := context.Background()

	// Setup
	disputeRepo := newMockDisputeRepository()
	orderRepo := newMockOrderRepository()
	userRepo := newMockUserRepository()
	opLogRepo := &mockOperationLogRepository{}
	notifRepo := &mockNotificationRepository{}
	paymentRepo := &mockPaymentRepository{}

	svc := NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notifRepo, paymentRepo)

	// Create dispute with past SLA deadline
	pastDeadline := time.Now().Add(-5 * time.Minute)
	userID := uint64(2)
	dispute := &model.OrderDispute{
		Base:             model.Base{ID: 1},
		OrderID:          1,
		UserID:           1,
		Status:           model.DisputeStatusAssigned,
		AssignedToUserID: &userID,
		SLADeadline:      &pastDeadline,
		SLABreached:      false,
		TraceID:          "trace-001",
	}
	disputeRepo.Create(ctx, dispute)

	// Test: Check and mark SLA breaches
	err := svc.CheckAndMarkSLABreaches(ctx)
	if err != nil {
		t.Fatalf("CheckAndMarkSLABreaches failed: %v", err)
	}

	// Verify SLA marked as breached
	breachedDispute, _ := disputeRepo.Get(ctx, 1)
	if !breachedDispute.SLABreached {
		t.Fatal("Expected dispute to be marked as SLA breached")
	}

	if breachedDispute.SLABreachedAt == nil {
		t.Fatal("Expected SLA breach timestamp")
	}

	t.Log("✓ SLABreachDetection test passed")
}
