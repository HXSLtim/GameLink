package order

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/commission"
)

// Mock repositories (reusing some from player service tests)
type mockOrderRepository struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		// Filter by statuses
		if len(opts.Statuses) > 0 {
			match := false
			for _, s := range opts.Statuses {
				if o.Status == s {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		// Filter by UserID
		if opts.UserID != nil && o.UserID != *opts.UserID {
			continue
		}
		result = append(result, *o)
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return repository.ErrNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

type mockPlayerRepository struct{}

func (m *mockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *mockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{
		{
			Base:            model.Base{ID: 1},
			UserID:          1,
			Nickname:        "TestPlayer",
			HourlyRateCents: 10000,
		},
	}, 1, nil
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{
		Base:            model.Base{ID: id},
		UserID:          1,
		Nickname:        "TestPlayer",
		HourlyRateCents: 10000,
	}, nil
}

func (m *mockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{
		Base:            model.Base{ID: 1},
		UserID:          userID,
		Nickname:        "TestPlayer",
		HourlyRateCents: 10000,
	}, nil
}

type mockUserRepository struct{}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	return []model.User{}, nil
}

func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{
		Base: model.Base{ID: id},
		Name: "TestUser",
	}, nil
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

type mockGameRepository struct{}

func (m *mockGameRepository) List(ctx context.Context) ([]model.Game, error) {
	return []model.Game{}, nil
}

func (m *mockGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return []model.Game{}, 0, nil
}

func (m *mockGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	return &model.Game{
		Base: model.Base{ID: id},
		Name: "TestGame",
	}, nil
}

func (m *mockGameRepository) Create(ctx context.Context, game *model.Game) error {
	return nil
}

func (m *mockGameRepository) Update(ctx context.Context, game *model.Game) error {
	return nil
}

func (m *mockGameRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockPaymentRepository struct{}

func (m *mockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return []model.Payment{}, 0, nil
}

func (m *mockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return &model.Payment{}, nil
}

func (m *mockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockReviewRepository struct{}

func (m *mockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	return []model.Review{}, 0, nil
}

func (m *mockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	return &model.Review{}, nil
}

func (m *mockReviewRepository) Create(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepository) Update(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockCommissionRepository struct{}

func (m *mockCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}

func (m *mockCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}

func (m *mockCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}

func (m *mockCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}

func (m *mockCommissionRepository) ListRules(ctx context.Context, opts commission.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return []model.CommissionRule{}, 0, nil
}

func (m *mockCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}

func (m *mockCommissionRepository) DeleteRule(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

func (m *mockCommissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return nil, nil
}

func (m *mockCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return nil, nil
}

func (m *mockCommissionRepository) ListRecords(ctx context.Context, opts commission.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return []model.CommissionRecord{}, 0, nil
}

func (m *mockCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

func (m *mockCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

func (m *mockCommissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return nil, nil
}

func (m *mockCommissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return nil, nil
}

func (m *mockCommissionRepository) ListSettlements(ctx context.Context, opts commission.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return []model.MonthlySettlement{}, 0, nil
}

func (m *mockCommissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

func (m *mockCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (*commission.MonthlyStats, error) {
	return &commission.MonthlyStats{}, nil
}

func (m *mockCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}

func TestCreateOrder(t *testing.T) {
	svc := NewOrderService(
		newMockOrderRepository(),
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	now := time.Now().Add(24 * time.Hour)
	resp, err := svc.CreateOrder(context.Background(), 1, CreateOrderRequest{
		PlayerID:       1,
		GameID:         1,
		Title:          "Test Order",
		Description:    "Test description",
		ScheduledStart: &now,
		DurationHours:  2.0,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.PriceCents != 20000 { // 10000 * 2
		t.Errorf("expected 20000, got %d", resp.PriceCents)
	}
}

func TestGetMyOrders(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a test order
	now := time.Now()
	playerID := uint64(1)
	gameID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          1,
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	resp, err := svc.GetMyOrders(context.Background(), 1, MyOrderListRequest{
		Page:     1,
		PageSize: 20,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if len(resp.Orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(resp.Orders))
	}
}

func TestCancelOrder(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a test order
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	err := svc.CancelOrder(context.Background(), 1, 1, CancelOrderRequest{
		Reason: "Test cancellation",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedOrder := orderRepo.orders[1]
	if updatedOrder.Status != model.OrderStatusCanceled {
		t.Errorf("expected canceled status, got %s", updatedOrder.Status)
	}

	if updatedOrder.CancelReason != "Test cancellation" {
		t.Errorf("expected 'Test cancellation', got %s", updatedOrder.CancelReason)
	}
}

func TestCancelOrderUnauthorized(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a test order owned by user 2
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2, // Different user
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// User 1 tries to cancel user 2's order (should fail)
	err := svc.CancelOrder(context.Background(), 1, 1, CancelOrderRequest{
		Reason: "Test cancellation",
	})

	if err == nil {
		t.Error("expected error when unauthorized, got nil")
	}
}

func TestGetOrderDetail(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a test order
	now := time.Now()
	playerID := uint64(1)
	gameID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          1,
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Description:     "Test description",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	resp, err := svc.GetOrderDetail(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Order.Title != "Test Order" {
		t.Errorf("expected 'Test Order', got %s", resp.Order.Title)
	}
}

func TestGetMyOrdersWithStatusFilter(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create multiple orders with different statuses
	now := time.Now()
	orderRepo.orders[1] = &model.Order{
		Base:   model.Base{ID: 1, CreatedAt: now},
		UserID: 1,
		Status: model.OrderStatusPending,
	}
	orderRepo.orders[2] = &model.Order{
		Base:   model.Base{ID: 2, CreatedAt: now},
		UserID: 1,
		Status: model.OrderStatusCompleted,
	}

	// Filter by pending status
	resp, err := svc.GetMyOrders(context.Background(), 1, MyOrderListRequest{
		Page:     1,
		PageSize: 20,
		Status:   "pending",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(resp.Orders))
	}

	if resp.Orders[0].Status != model.OrderStatusPending {
		t.Errorf("expected pending status, got %s", resp.Orders[0].Status)
	}
}

// ---- Additional Tests for Better Coverage ----

func TestCompleteOrder(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create an in-progress order
	now := time.Now()
	playerID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		PlayerID:        &playerID,
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	err := svc.CompleteOrder(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedOrder := orderRepo.orders[1]
	if updatedOrder.Status != model.OrderStatusCompleted {
		t.Errorf("expected completed status, got %s", updatedOrder.Status)
	}
}

func TestCompleteOrder_InvalidTransition(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a pending order (can't complete directly)
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	err := svc.CompleteOrder(context.Background(), 1, 1)

	if err != ErrInvalidTransition {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
}

func TestAcceptOrder_Success(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a confirmed order (ready to be accepted)
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2,                          // Different user
		Status:          model.OrderStatusConfirmed, // Must be confirmed to accept
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Player (user 1) accepts the order
	err := svc.AcceptOrder(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedOrder := orderRepo.orders[1]
	if updatedOrder.Status != model.OrderStatusInProgress { // After accepting, it should be in-progress
		t.Errorf("expected in-progress status, got %s", updatedOrder.Status)
	}

	if updatedOrder.PlayerID == nil || *updatedOrder.PlayerID == 0 {
		t.Error("expected player ID to be set")
	}
}

func TestCompleteOrder_InvalidStatus(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a pending order (not yet in progress)
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Try to complete (should fail - must be in-progress first)
	err := svc.CompleteOrder(context.Background(), 1, 1)

	if err == nil {
		t.Error("expected error for invalid status transition")
	}
}

func TestCompleteOrderByPlayer(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create an in-progress order assigned to player 1
	now := time.Now()
	playerID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2,
		PlayerID:        &playerID,
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Player 1 completes the order
	err := svc.CompleteOrderByPlayer(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedOrder := orderRepo.orders[1]
	if updatedOrder.Status != model.OrderStatusCompleted {
		t.Errorf("expected completed status, got %s", updatedOrder.Status)
	}
}

func TestCompleteOrderByPlayer_Unauthorized(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create an order assigned to player 2
	now := time.Now()
	playerID := uint64(2)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          3,
		PlayerID:        &playerID, // Different player
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Player 1 tries to complete player 2's order (should fail)
	err := svc.CompleteOrderByPlayer(context.Background(), 1, 1)

	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGetMyOrders_EmptyList(t *testing.T) {
	svc := NewOrderService(
		newMockOrderRepository(),
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	resp, err := svc.GetMyOrders(context.Background(), 1, MyOrderListRequest{
		Page:     1,
		PageSize: 20,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if len(resp.Orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(resp.Orders))
	}
}

func TestCancelOrder_InvalidStatus(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a completed order (cannot be canceled)
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	err := svc.CancelOrder(context.Background(), 1, 1, CancelOrderRequest{
		Reason: "Test",
	})

	if err == nil {
		t.Error("expected error when canceling completed order")
	}
}

func TestCompleteOrder_Unauthorized(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create an order for user 2
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2, // Different user
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// User 1 tries to complete user 2's order
	err := svc.CompleteOrder(context.Background(), 1, 1)

	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGetOrderDetail_Unauthorized(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create an order for user 2
	now := time.Now()
	playerID := uint64(3)
	gameID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          2, // Different user
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// User 1 tries to view user 2's order (not their order and not their player order)
	_, err := svc.GetOrderDetail(context.Background(), 1, 1)

	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGetOrderDetail_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Try to get non-existent order
	_, err := svc.GetOrderDetail(context.Background(), 1, 9999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCancelOrder_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Try to cancel non-existent order
	err := svc.CancelOrder(context.Background(), 1, 9999, CancelOrderRequest{
		Reason: "Test",
	})

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCompleteOrder_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Try to complete non-existent order
	err := svc.CompleteOrder(context.Background(), 1, 9999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCompleteOrderByPlayer_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Try to complete non-existent order
	err := svc.CompleteOrderByPlayer(context.Background(), 1, 9999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAcceptOrder_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Try to accept non-existent order
	err := svc.AcceptOrder(context.Background(), 1, 9999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAcceptOrder_InvalidStatus(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a pending order (not confirmed yet)
	now := time.Now()
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2,
		Status:          model.OrderStatusPending, // Wrong status for accepting
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Try to accept (should fail)
	err := svc.AcceptOrder(context.Background(), 1, 1)

	if err != ErrInvalidTransition {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
}

func TestCompleteOrderByPlayer_InvalidStatus(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	// Create a pending order (not in-progress)
	now := time.Now()
	playerID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          2,
		PlayerID:        &playerID,
		Status:          model.OrderStatusPending, // Wrong status for completing
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	// Try to complete (should fail)
	err := svc.CompleteOrderByPlayer(context.Background(), 1, 1)

	if err == nil {
		t.Error("expected error for invalid status transition")
	}
}

// TestGetOrderDetail_WithPayment 测试获取订单详情（包含支付信息）
func TestGetOrderDetail_WithPayment(t *testing.T) {
	orderRepo := newMockOrderRepository()
	now := time.Now()
	playerID := uint64(1)
	gameID := uint64(1)
	paidAt := now.Add(1 * time.Hour)
	
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          1,
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	paymentRepo := &mockPaymentRepositoryWithData{
		payments: []model.Payment{
			{
				Base:        model.Base{ID: 1, CreatedAt: now},
				OrderID:     1,
				UserID:      1,
				AmountCents: 10000,
				Status:      model.PaymentStatusPaid,
				Method:      model.PaymentMethodWeChat,
				PaidAt:      &paidAt,
			},
		},
	}

	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		paymentRepo,
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	resp, err := svc.GetOrderDetail(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Payment == nil {
		t.Error("expected payment information")
	}

	if resp.Payment.Status != model.PaymentStatusPaid {
		t.Errorf("expected payment status 'paid', got '%s'", resp.Payment.Status)
	}
}

// TestGetOrderDetail_WithReview 测试获取订单详情（包含评价信息）
func TestGetOrderDetail_WithReview(t *testing.T) {
	orderRepo := newMockOrderRepository()
	now := time.Now()
	playerID := uint64(1)
	gameID := uint64(1)
	
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          1,
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
	}
	orderRepo.orders[1] = order

	reviewRepo := &mockReviewRepositoryWithData{
		reviews: []model.Review{
			{
				Base:     model.Base{ID: 1, CreatedAt: now},
				OrderID:  1,
				UserID:   1,
				PlayerID: playerID,
				Score:    5,
				Content:  "Great service!",
			},
		},
	}

	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		reviewRepo,
		&mockCommissionRepository{},
	)

	resp, err := svc.GetOrderDetail(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Review == nil {
		t.Error("expected review information")
	}

	if resp.Review.Rating != 5 {
		t.Errorf("expected rating 5, got %d", resp.Review.Rating)
	}
}

// TestGetOrderDetail_WithTimeline 测试获取订单详情（包含时间线）
func TestGetOrderDetail_WithTimeline(t *testing.T) {
	orderRepo := newMockOrderRepository()
	now := time.Now()
	playerID := uint64(1)
	gameID := uint64(1)
	startedAt := now.Add(1 * time.Hour)
	completedAt := now.Add(2 * time.Hour)
	
	order := &model.Order{
		Base:            model.Base{ID: 1, CreatedAt: now},
		UserID:          1,
		PlayerID:        &playerID,
		GameID:          &gameID,
		Title:           "Test Order",
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: 10000,
		ScheduledStart:  &now,
		StartedAt:       &startedAt,
		CompletedAt:     &completedAt,
	}
	orderRepo.orders[1] = order

	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	resp, err := svc.GetOrderDetail(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Timeline) == 0 {
		t.Error("expected timeline items")
	}

	// 验证时间线包含创建、开始、完成等事件
	hasCreated := false
	hasStarted := false
	hasCompleted := false
	for _, item := range resp.Timeline {
		switch item.Status {
		case string(model.OrderStatusPending):
			hasCreated = true
		case string(model.OrderStatusInProgress):
			hasStarted = true
		case string(model.OrderStatusCompleted):
			hasCompleted = true
		}
	}

	if !hasCreated {
		t.Error("expected 'created' event in timeline")
	}
	if !hasStarted {
		t.Error("expected 'started' event in timeline")
	}
	if !hasCompleted {
		t.Error("expected 'completed' event in timeline")
	}
}

// TestCancelOrder_EdgeCases 测试取消订单的边界条件
func TestCancelOrder_EdgeCases(t *testing.T) {
	orderRepo := newMockOrderRepository()
	svc := NewOrderService(
		orderRepo,
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
		&mockCommissionRepository{},
	)

	now := time.Now()

	t.Run("取消已支付的订单", func(t *testing.T) {
		order := &model.Order{
			Base:            model.Base{ID: 1},
			UserID:          1,
			Status:          model.OrderStatusConfirmed, // 已确认（通常已支付）
			TotalPriceCents: 10000,
			ScheduledStart:  &now,
		}
		orderRepo.orders[1] = order

		err := svc.CancelOrder(context.Background(), 1, 1, CancelOrderRequest{
			Reason: "Change of mind",
		})

		// 已支付的订单应该可以取消（但可能需要退款）
		if err != nil {
			t.Logf("Cancel paid order returned: %v (may be expected)", err)
		}
	})

	t.Run("取消已完成的订单应该失败", func(t *testing.T) {
		order := &model.Order{
			Base:            model.Base{ID: 2},
			UserID:          1,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
			ScheduledStart:  &now,
		}
		orderRepo.orders[2] = order

		err := svc.CancelOrder(context.Background(), 1, 2, CancelOrderRequest{
			Reason: "Too late",
		})

		if err == nil {
			t.Error("expected error when canceling completed order")
		}
	})
}

// mockPaymentRepositoryWithData 提供数据的mock支付仓库
type mockPaymentRepositoryWithData struct {
	payments []model.Payment
}

func (m *mockPaymentRepositoryWithData) Create(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepositoryWithData) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	var filtered []model.Payment
	for _, p := range m.payments {
		if opts.OrderID != nil && p.OrderID != *opts.OrderID {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockPaymentRepositoryWithData) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	for _, p := range m.payments {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockPaymentRepositoryWithData) Update(ctx context.Context, payment *model.Payment) error {
	return nil
}

func (m *mockPaymentRepositoryWithData) Delete(ctx context.Context, id uint64) error {
	return nil
}

type mockReviewRepositoryWithData struct {
	reviews []model.Review
}

func (m *mockReviewRepositoryWithData) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	var filtered []model.Review
	for _, r := range m.reviews {
		if opts.OrderID != nil && r.OrderID != *opts.OrderID {
			continue
		}
		if opts.PlayerID != nil && r.PlayerID != *opts.PlayerID {
			continue
		}
		if opts.UserID != nil && r.UserID != *opts.UserID {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockReviewRepositoryWithData) Get(ctx context.Context, id uint64) (*model.Review, error) {
	for _, r := range m.reviews {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockReviewRepositoryWithData) Create(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepositoryWithData) Update(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *mockReviewRepositoryWithData) Delete(ctx context.Context, id uint64) error {
	return nil
}
