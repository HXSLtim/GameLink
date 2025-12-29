package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// Mock implementations for dispute service testing
type MockDisputeRepository struct {
	createDispute       func(ctx context.Context, dispute *model.OrderDispute) error
	getDispute          func(ctx context.Context, id uint64) (*model.OrderDispute, error)
	getDisputeByOrderID func(ctx context.Context, orderID uint64) (*model.OrderDispute, error)
	updateDispute       func(ctx context.Context, dispute *model.OrderDispute) error
	listDisputes        func(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error)
	listSLABreached     func(ctx context.Context) ([]model.OrderDispute, error)
	markSLABreached     func(ctx context.Context, disputeID uint64) error
	getStats            func(ctx context.Context) (map[string]int64, error)
	listPendingAssignment func(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error)
}

func (m *MockDisputeRepository) Create(ctx context.Context, dispute *model.OrderDispute) error {
	if m.createDispute != nil {
		return m.createDispute(ctx, dispute)
	}
	return nil
}

func (m *MockDisputeRepository) Get(ctx context.Context, id uint64) (*model.OrderDispute, error) {
	if m.getDispute != nil {
		return m.getDispute(ctx, id)
	}
	return nil, nil
}

func (m *MockDisputeRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
	if m.getDisputeByOrderID != nil {
		return m.getDisputeByOrderID(ctx, orderID)
	}
	return nil, nil
}

func (m *MockDisputeRepository) Update(ctx context.Context, dispute *model.OrderDispute) error {
	if m.updateDispute != nil {
		return m.updateDispute(ctx, dispute)
	}
	return nil
}

func (m *MockDisputeRepository) List(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
	if m.listDisputes != nil {
		return m.listDisputes(ctx, opts)
	}
	return nil, 0, nil
}

func (m *MockDisputeRepository) ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
	if m.listPendingAssignment != nil {
		return m.listPendingAssignment(ctx, page, pageSize)
	}
	if m.listDisputes != nil {
		return m.listDisputes(ctx, repository.DisputeListOptions{Page: page, PageSize: pageSize, Statuses: []model.DisputeStatus{model.DisputeStatusPending}})
	}
	return nil, 0, nil
}

func (m *MockDisputeRepository) ListSLABreached(ctx context.Context) ([]model.OrderDispute, error) {
	if m.listSLABreached != nil {
		return m.listSLABreached(ctx)
	}
	return nil, nil
}

func (m *MockDisputeRepository) MarkSLABreached(ctx context.Context, disputeID uint64) error {
	if m.markSLABreached != nil {
		return m.markSLABreached(ctx, disputeID)
	}
	return nil
}

func (m *MockDisputeRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockDisputeRepository) CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error) { return 0, nil }
func (m *MockDisputeRepository) GetPendingCount(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockDisputeRepository) GetStats(ctx context.Context) (map[string]int64, error) {
	if m.getStats != nil {
		return m.getStats(ctx)
	}
	return map[string]int64{}, nil
}

type MockOperationLogRepository struct {
	appendLog func(ctx context.Context, log *model.OperationLog) error
}

func (m *MockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	if m.appendLog != nil {
		return m.appendLog(ctx, log)
	}
	return nil
}

func (m *MockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

func (m *MockOperationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

type MockNotificationRepository struct {
	createNotification func(ctx context.Context, event *model.NotificationEvent) error
}

func (m *MockNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	if m.createNotification != nil {
		return m.createNotification(ctx, event)
	}
	return nil
}

func (m *MockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}

func (m *MockNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error { return nil }
func (m *MockNotificationRepository) MarkAllRead(ctx context.Context, userID uint64) error { return nil }
func (m *MockNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) { return 0, nil }
func (m *MockNotificationRepository) Delete(ctx context.Context, userID uint64, id uint64) error { return nil }

// Helper functions
func createTestOrderForDispute(orderID uint64, userID uint64, status model.OrderStatus) *model.Order {
	playerID := uint64(100)
	gameID := uint64(1)
	now := time.Now()

	return &model.Order{
		Base:             model.Base{ID: orderID, CreatedAt: now, UpdatedAt: now},
		OrderNo:          model.GenerateEscortOrderNo(),
		UserID:           userID,
		ItemID:           1,
		PlayerID:         &playerID,
		GameID:           &gameID,
		Quantity:         1,
		UnitPriceCents:   5000,
		TotalPriceCents:  5000,
		CommissionCents:  1000,
		PlayerIncomeCents: 4000,
		Currency:         model.CurrencyCNY,
		Status:           status,
		Title:            "Test Order",
		HasDispute:       false,
	}
}

func createTestDispute(disputeID uint64, orderID uint64, initiatorID uint64) *model.OrderDispute {
	slaDeadline := time.Now().Add(30 * time.Minute)
	traceID := "test-trace-123"
	now := time.Now()

	return &model.OrderDispute{
		Base:           model.Base{ID: disputeID, CreatedAt: now, UpdatedAt: now},
		OrderID:        orderID,
		InitiatorID:    initiatorID,
		InitiatorType:  model.DisputeInitiatorUser,
		Type:           model.DisputeTypeServiceQuality,
		Status:         model.DisputeStatusPending,
		Reason:         "Service quality issue",
		EvidenceText:   "Test evidence",
		EvidenceURLs:   model.EvidenceURLArray{"https://example.com/evidence.jpg"},
		SLADeadline:    &slaDeadline,
		TraceID:        traceID,
	}
}

// TestDisputeService_InitiateDispute_Success tests successful dispute initiation
func TestDisputeService_InitiateDispute_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	// Use in_progress status to allow dispute initiation
	testOrder := createTestOrderForDispute(orderID, userID, model.OrderStatusInProgress)

	var createdDispute *model.OrderDispute
	var updatedOrder *model.Order

	disputes := &MockDisputeRepository{
		getDisputeByOrderID: func(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
			return nil, repository.ErrNotFound // No existing dispute
		},
		createDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			createdDispute = dispute
			dispute.ID = 123
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	users := &MockUserRepository{}
	players := &MockPlayerRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeServiceWithPlayers(disputes, orders, users, players, operationLogs, notifications, payments)

	req := InitiateDisputeRequest{
		OrderID:       orderID,
		InitiatorID:   userID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Service quality issue",
		EvidenceText:  "Test evidence",
		EvidenceURLs:  []string{"https://example.com/evidence.jpg"},
	}

	resp, err := service.InitiateDispute(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.DisputeID)
	assert.NotEmpty(t, resp.TraceID)
	assert.NotNil(t, resp.SLADeadline)
	assert.NotNil(t, createdDispute)
	assert.Equal(t, orderID, createdDispute.OrderID)
	assert.Equal(t, userID, createdDispute.InitiatorID)
	assert.NotNil(t, updatedOrder)
	assert.True(t, updatedOrder.HasDispute)
	assert.Equal(t, model.OrderStatusDisputed, updatedOrder.Status)
}

// TestDisputeService_InitiateDispute_OrderNotFound tests dispute initiation with non-existent order
func TestDisputeService_InitiateDispute_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(999)

	disputes := &MockDisputeRepository{}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := InitiateDisputeRequest{
		OrderID:       orderID,
		InitiatorID:   userID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Test reason",
	}

	resp, err := service.InitiateDispute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrOrderNotFound, err)
}

// TestDisputeService_InitiateDispute_Unauthorized tests dispute initiation by unauthorized user
func TestDisputeService_InitiateDispute_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	wrongUserID := uint64(999)
	orderID := uint64(1)

	testOrder := createTestOrderForDispute(orderID, userID, model.OrderStatusCompleted)

	disputes := &MockDisputeRepository{}
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := InitiateDisputeRequest{
		OrderID:       orderID,
		InitiatorID:   wrongUserID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Test reason",
	}

	resp, err := service.InitiateDispute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrDisputeUnauthorized, err)
}

// TestDisputeService_InitiateDispute_DisputeExists tests dispute initiation when dispute already exists
func TestDisputeService_InitiateDispute_DisputeExists(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	// Use in_progress status to allow dispute initiation
	testOrder := createTestOrderForDispute(orderID, userID, model.OrderStatusInProgress)
	existingDispute := createTestDispute(1, orderID, userID)

	disputes := &MockDisputeRepository{
		getDisputeByOrderID: func(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
			return existingDispute, nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := InitiateDisputeRequest{
		OrderID:       orderID,
		InitiatorID:   userID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Test reason",
	}

	resp, err := service.InitiateDispute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrDisputeExists, err)
}

// TestDisputeService_InitiateDispute_InvalidInput tests dispute initiation with invalid input
func TestDisputeService_InitiateDispute_InvalidInput(t *testing.T) {
	ctx := context.Background()

	disputes := &MockDisputeRepository{}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	tests := []struct {
		name    string
		req     InitiateDisputeRequest
		wantErr error
	}{
		{
			name: "Empty order ID",
			req: InitiateDisputeRequest{
				OrderID:       0,
				InitiatorID:   1,
				InitiatorType: model.DisputeInitiatorUser,
				Reason:        "Test",
			},
			wantErr: ErrDisputeValidation,
		},
		{
			name: "Empty initiator ID",
			req: InitiateDisputeRequest{
				OrderID:       1,
				InitiatorID:   0,
				InitiatorType: model.DisputeInitiatorUser,
				Reason:        "Test",
			},
			wantErr: ErrDisputeValidation,
		},
		{
			name: "Empty reason",
			req: InitiateDisputeRequest{
				OrderID:       1,
				InitiatorID:   1,
				InitiatorType: model.DisputeInitiatorUser,
				Reason:        "",
			},
			wantErr: apierr.BadRequest("纠纷原因不能为空"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.InitiateDispute(ctx, tt.req)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}

// TestDisputeService_AssignDispute_Success tests successful dispute assignment
func TestDisputeService_AssignDispute_Success(t *testing.T) {
	ctx := context.Background()
	disputeID := uint64(1)
	serviceUserID := uint64(100)

	testDispute := createTestDispute(disputeID, 1, 1)

	var updatedDispute *model.OrderDispute
	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			return testDispute, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			updatedDispute = dispute
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := AssignDisputeRequest{
		DisputeID:         disputeID,
		AssignedServiceID: serviceUserID,
		ActorUserID:       200,
	}

	err := service.AssignDispute(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedDispute)
	assert.Equal(t, model.DisputeStatusAssigned, updatedDispute.Status)
	assert.Equal(t, serviceUserID, *updatedDispute.AssignedServiceID)
}

// TestDisputeService_AssignDispute_InvalidStatus tests dispute assignment with invalid status
func TestDisputeService_AssignDispute_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	disputeID := uint64(1)

	testDispute := createTestDispute(disputeID, 1, 1)
	testDispute.Status = model.DisputeStatusResolved

	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			return testDispute, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := AssignDisputeRequest{
		DisputeID:         disputeID,
		AssignedServiceID: 100,
		ActorUserID:       200,
	}

	err := service.AssignDispute(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "纠纷状态不是待处理")
}

// TestDisputeService_ResolveDispute_Refund tests successful dispute resolution with refund
func TestDisputeService_ResolveDispute_Refund(t *testing.T) {
	ctx := context.Background()
	disputeID := uint64(1)
	orderID := uint64(1)

	testDispute := createTestDispute(disputeID, orderID, 1)
	testOrder := createTestOrderForDispute(orderID, 1, model.OrderStatusDisputed)

	var updatedDispute *model.OrderDispute
	var updatedOrder *model.Order

	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			return testDispute, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			updatedDispute = dispute
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := ResolveDisputeRequest{
		DisputeID:     disputeID,
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "Full refund approved",
		ActorUserID:   100,
	}

	err := service.ResolveDispute(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedDispute)
	assert.Equal(t, model.DisputeStatusResolved, updatedDispute.Status)
	assert.Equal(t, model.ResolutionRefund, updatedDispute.Resolution)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(5000), updatedOrder.RefundAmountCents)
	assert.NotNil(t, updatedOrder.RefundedAt)
}

// TestDisputeService_ResolveDispute_Reject tests dispute resolution rejection
func TestDisputeService_ResolveDispute_Reject(t *testing.T) {
	ctx := context.Background()
	disputeID := uint64(1)
	orderID := uint64(1)

	testDispute := createTestDispute(disputeID, orderID, 1)
	testOrder := createTestOrderForDispute(orderID, 1, model.OrderStatusDisputed)

	var updatedDispute *model.OrderDispute
	var updatedOrder *model.Order

	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			return testDispute, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			updatedDispute = dispute
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := ResolveDisputeRequest{
		DisputeID:     disputeID,
		Resolution:    model.ResolutionReject,
		ResolveRemark: "Dispute rejected",
		ActorUserID:   100,
	}

	err := service.ResolveDispute(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedDispute)
	assert.Equal(t, model.DisputeStatusRejected, updatedDispute.Status)
	assert.Equal(t, model.ResolutionReject, updatedDispute.Resolution)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.False(t, updatedOrder.HasDispute)
}

// TestDisputeService_RollbackDisputeAssignment_Success tests successful dispute assignment rollback
func TestDisputeService_RollbackDisputeAssignment_Success(t *testing.T) {
	ctx := context.Background()
	disputeID := uint64(1)

	testDispute := createTestDispute(disputeID, 1, 1)
	testDispute.Status = model.DisputeStatusAssigned
	serviceUserID := uint64(100)
	testDispute.AssignedServiceID = &serviceUserID

	var updatedDispute *model.OrderDispute
	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			return testDispute, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			updatedDispute = dispute
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := RollbackDisputeRequest{
		DisputeID:      disputeID,
		RollbackReason: "Mistaken assignment",
		ActorUserID:    200,
	}

	err := service.RollbackDisputeAssignment(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedDispute)
	assert.Equal(t, model.DisputeStatusPending, updatedDispute.Status)
	assert.Nil(t, updatedDispute.AssignedServiceID)
	assert.Nil(t, updatedDispute.OriginalServiceID)
	assert.NotNil(t, updatedDispute.RolledBackAt)
	assert.Equal(t, "Mistaken assignment", updatedDispute.RollbackReason)
}

// TestDisputeService_CheckAndMarkSLABreaches tests SLA breach checking and marking
func TestDisputeService_CheckAndMarkSLABreaches(t *testing.T) {
	ctx := context.Background()

	breachedDispute := createTestDispute(1, 1, 1)
	slaPast := time.Now().Add(-1 * time.Hour)
	breachedDispute.SLADeadline = &slaPast
	breachedDispute.Status = model.DisputeStatusAssigned

	var markedDisputeIDs []uint64
	disputes := &MockDisputeRepository{
		listSLABreached: func(ctx context.Context) ([]model.OrderDispute, error) {
			return []model.OrderDispute{*breachedDispute}, nil
		},
		markSLABreached: func(ctx context.Context, disputeID uint64) error {
			markedDisputeIDs = append(markedDisputeIDs, disputeID)
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	err := service.CheckAndMarkSLABreaches(ctx)

	require.NoError(t, err)
	assert.Equal(t, []uint64{1}, markedDisputeIDs)
}

// TestDisputeService_ListPendingDisputes tests listing pending disputes
func TestDisputeService_ListPendingDisputes(t *testing.T) {
	ctx := context.Background()

	pendingDisputes := []model.OrderDispute{
		*createTestDispute(1, 1, 1),
		*createTestDispute(2, 2, 2),
	}

	disputes := &MockDisputeRepository{
		listPendingAssignment: func(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error) {
			return pendingDisputes, 2, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	result, total, err := service.ListPendingDisputes(ctx, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
}

// TestDisputeService_ListDisputesByStatus tests listing disputes by status
func TestDisputeService_ListDisputesByStatus(t *testing.T) {
	ctx := context.Background()

	filteredDisputes := []model.OrderDispute{
		*createTestDispute(1, 1, 1),
	}

	var listOpts repository.DisputeListOptions
	disputes := &MockDisputeRepository{
		listDisputes: func(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
			listOpts = opts
			return filteredDisputes, 1, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	statuses := []model.DisputeStatus{model.DisputeStatusPending}
	result, total, err := service.ListDisputesByStatus(ctx, statuses, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(1), total)
	assert.Equal(t, statuses, listOpts.Statuses)
}

// TestDisputeService_ListDisputes tests listing disputes with filters
func TestDisputeService_ListDisputes(t *testing.T) {
	ctx := context.Background()

	filteredDisputes := []model.OrderDispute{
		*createTestDispute(1, 1, 1),
	}

	var listOpts repository.DisputeListOptions
	disputes := &MockDisputeRepository{
		listDisputes: func(ctx context.Context, opts repository.DisputeListOptions) ([]model.OrderDispute, int64, error) {
			listOpts = opts
			return filteredDisputes, 1, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := ListDisputesRequest{
		Page:     1,
		PageSize: 10,
		Status:   "pending",
		OrderNo:  "TEST123",
	}

	result, total, err := service.ListDisputes(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "pending", string(listOpts.Statuses[0]))
	assert.Equal(t, "TEST123", listOpts.OrderNo)
}

// TestDisputeService_GetDisputeStats tests getting dispute statistics
func TestDisputeService_GetDisputeStats(t *testing.T) {
	ctx := context.Background()

	expectedStats := map[string]int64{
		"pending":  10,
		"assigned": 5,
		"resolved": 20,
	}

	disputes := &MockDisputeRepository{
		getStats: func(ctx context.Context) (map[string]int64, error) {
			return expectedStats, nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	stats, err := service.GetDisputeStats(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedStats, stats)
}

// TestDisputeService_BatchAssignDisputes_Success tests successful batch dispute assignment
func TestDisputeService_BatchAssignDisputes_Success(t *testing.T) {
	ctx := context.Background()
	serviceUserID := uint64(100)

	dispute1 := createTestDispute(1, 1, 1)
	dispute2 := createTestDispute(2, 2, 2)

	getCount := 0
	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			getCount++
			if id == 1 {
				return dispute1, nil
			}
			return dispute2, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := BatchAssignDisputesRequest{
		DisputeIDs:        []uint64{1, 2},
		AssignedServiceID: serviceUserID,
		ActorUserID:       200,
	}

	result, err := service.BatchAssignDisputes(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Equal(t, "批量分配成功，共2个纠纷", result.Message)
}

// TestDisputeService_BatchAssignDisputes_PartialFailure tests batch assignment with partial failures
func TestDisputeService_BatchAssignDisputes_PartialFailure(t *testing.T) {
	ctx := context.Background()
	serviceUserID := uint64(100)

	dispute1 := createTestDispute(1, 1, 1)
	dispute2 := createTestDispute(2, 2, 2)
	dispute2.Status = model.DisputeStatusResolved // Can't be assigned

	getCount := 0
	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			getCount++
			if id == 1 {
				return dispute1, nil
			}
			return dispute2, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := BatchAssignDisputesRequest{
		DisputeIDs:        []uint64{1, 2},
		AssignedServiceID: serviceUserID,
		ActorUserID:       200,
	}

	result, err := service.BatchAssignDisputes(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.Success) // Overall success since at least one succeeded
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Equal(t, 1, len(result.Errors))
	assert.Contains(t, result.Errors[0].Error, "无法分配")
}

// TestDisputeService_BatchUpdateDisputesStatus_Success tests successful batch status update
func TestDisputeService_BatchUpdateDisputesStatus_Success(t *testing.T) {
	ctx := context.Background()

	dispute1 := createTestDispute(1, 1, 1)
	dispute2 := createTestDispute(2, 2, 2)
	// Set dispute status to assigned so they can transition to mediating
	dispute1.Status = model.DisputeStatusAssigned
	dispute2.Status = model.DisputeStatusAssigned

	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			if id == 1 {
				return dispute1, nil
			}
			return dispute2, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			return nil
		},
	}

	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := BatchUpdateDisputesStatusRequest{
		DisputeIDs:  []uint64{1, 2},
		Status:      model.DisputeStatusMediating,
		ActorUserID: 100,
	}

	result, err := service.BatchUpdateDisputesStatus(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// TestDisputeService_BatchCloseDisputes_Success tests successful batch closing disputes
func TestDisputeService_BatchCloseDisputes_Success(t *testing.T) {
	ctx := context.Background()

	dispute1 := createTestDispute(1, 1, 1)
	order1 := createTestOrderForDispute(1, 1, model.OrderStatusDisputed)

	dispute2 := createTestDispute(2, 2, 2)
	order2 := createTestOrderForDispute(2, 2, model.OrderStatusDisputed)

	disputes := &MockDisputeRepository{
		getDispute: func(ctx context.Context, id uint64) (*model.OrderDispute, error) {
			if id == 1 {
				return dispute1, nil
			}
			return dispute2, nil
		},
		updateDispute: func(ctx context.Context, dispute *model.OrderDispute) error {
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == 1 {
				return order1, nil
			}
			return order2, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return nil
		},
	}

	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := BatchCloseDisputesRequest{
		DisputeIDs:    []uint64{1, 2},
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "Batch refund",
		ActorUserID:   100,
	}

	result, err := service.BatchCloseDisputes(ctx, req)

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// TestDisputeService_BatchCloseDisputes_InvalidResolution tests batch closing with invalid resolution
func TestDisputeService_BatchCloseDisputes_InvalidResolution(t *testing.T) {
	ctx := context.Background()

	disputes := &MockDisputeRepository{}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	req := BatchCloseDisputesRequest{
		DisputeIDs:    []uint64{1, 2},
		Resolution:    "invalid_resolution",
		ResolveRemark: "Test",
		ActorUserID:   100,
	}

	result, err := service.BatchCloseDisputes(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test_canTransitionTo tests dispute status transition validation
func Test_canTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     model.DisputeStatus
		to       model.DisputeStatus
		expected bool
	}{
		{
			name:     "Pending to Assigned",
			from:     model.DisputeStatusPending,
			to:       model.DisputeStatusAssigned,
			expected: true,
		},
		{
			name:     "Pending to Canceled",
			from:     model.DisputeStatusPending,
			to:       model.DisputeStatusCanceled,
			expected: true,
		},
		{
			name:     "Assigned to Mediating",
			from:     model.DisputeStatusAssigned,
			to:       model.DisputeStatusMediating,
			expected: true,
		},
		{
			name:     "Assigned to Canceled",
			from:     model.DisputeStatusAssigned,
			to:       model.DisputeStatusCanceled,
			expected: true,
		},
		{
			name:     "Mediating to Assigned",
			from:     model.DisputeStatusMediating,
			to:       model.DisputeStatusAssigned,
			expected: true,
		},
		{
			name:     "Pending to Resolved",
			from:     model.DisputeStatusPending,
			to:       model.DisputeStatusResolved,
			expected: false,
		},
		{
			name:     "Resolved to Assigned",
			from:     model.DisputeStatusResolved,
			to:       model.DisputeStatusAssigned,
			expected: false,
		},
		{
			name:     "Invalid transition",
			from:     model.DisputeStatusCanceled,
			to:       model.DisputeStatusMediating,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canTransitionTo(tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDisputeService_processRefund tests refund processing
func TestDisputeService_processRefund(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrderForDispute(orderID, 1, model.OrderStatusDisputed)
	testDispute := createTestDispute(1, orderID, 1)
	testDispute.ResolveRemark = "Quality issue refund"

	var updatedOrder *model.Order
	orders := &MockOrderRepository{
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	operationLogs := &MockOperationLogRepository{}
	disputes := &MockDisputeRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	actorID := uint64(100)
	err := service.processRefund(ctx, testOrder, testDispute, 5000, &actorID)

	require.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(5000), updatedOrder.RefundAmountCents)
	assert.Contains(t, updatedOrder.RefundReason, "Quality issue refund")
	assert.NotNil(t, updatedOrder.RefundedAt)
}

// TestDisputeService_logOperation tests operation logging
func TestDisputeService_logOperation(t *testing.T) {
	ctx := context.Background()

	var loggedOp *model.OperationLog
	operationLogs := &MockOperationLogRepository{
		appendLog: func(ctx context.Context, log *model.OperationLog) error {
			loggedOp = log
			return nil
		},
	}

	disputes := &MockDisputeRepository{}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	notifications := &MockNotificationRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	actorID := uint64(100)
	service.logOperation(ctx, model.OpEntityDispute, 1, model.OpActionInitiateDispute, "Test", "trace-123", &actorID)

	assert.NotNil(t, loggedOp)
	assert.Equal(t, "dispute", loggedOp.EntityType)
	assert.Equal(t, uint64(1), loggedOp.EntityID)
	assert.Equal(t, "initiate_dispute", loggedOp.Action)
	assert.Equal(t, "Test", loggedOp.Reason)
	assert.Equal(t, "trace-123", loggedOp.TraceID)
	assert.Equal(t, &actorID, loggedOp.ActorUserID)
}

// TestDisputeService_sendNotification tests notification sending
func TestDisputeService_sendNotification(t *testing.T) {
	ctx := context.Background()

	var sentNotification *model.NotificationEvent
	notifications := &MockNotificationRepository{
		createNotification: func(ctx context.Context, event *model.NotificationEvent) error {
			sentNotification = event
			return nil
		},
	}

	disputes := &MockDisputeRepository{}
	orders := &MockOrderRepository{}
	users := &MockUserRepository{}
	operationLogs := &MockOperationLogRepository{}
	payments := &MockPaymentRepository{}

	service := NewDisputeService(disputes, orders, users, operationLogs, notifications, payments)

	service.sendNotification(ctx, 1, "Test Title", "Test Message", "trace-123")

	assert.NotNil(t, sentNotification)
	assert.Equal(t, uint64(1), sentNotification.UserID)
	assert.Equal(t, "Test Title", sentNotification.Title)
	assert.Equal(t, "Test Message", sentNotification.Message)
	assert.Equal(t, model.NotificationPriorityHigh, sentNotification.Priority)
}
