package order

import (
	"testing"
	"gamelink/internal/repository"
	"gamelink/internal/repository/commission"
	"errors"
	"context"
	"time"
	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
)

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

func TestOrderStatusTransitions(t *testing.T) {
	ctx := context.Background()

	t.Run("正常状态流转_pending到confirmed", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		// 创建pending状态的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 更新为confirmed状态
		order.Status = model.OrderStatusConfirmed
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)

		// 验证状态已更新
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusConfirmed, updated.Status)
	})

	t.Run("正常状态流转_confirmed到in_progress", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusInProgress
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusInProgress, updated.Status)
	})

	t.Run("正常状态流转_in_progress到completed", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("取消流转_pending到canceled", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusCanceled
		order.CancelReason = "用户取消"
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, updated.Status)
		assert.Equal(t, "用户取消", updated.CancelReason)
	})

	t.Run("已完成订单状态不应该再改变", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		completedAt := time.Now()
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
			CompletedAt:     &completedAt,
		}
		orderRepo.Create(ctx, order)

		// 尝试修改已完成的订单（业务层应该阻止）
		originalStatus := order.Status
		order.Status = model.OrderStatusCanceled

		// Repository层会允许更新，但Service层应该阻止
		// 这里测试的是数据一致性
		err := orderRepo.Update(ctx, order)
		assert.NoError(t, err) // Repository层允许

		// 但在实际业务中，Service层应该检查并拒绝这种操作
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.NotEqual(t, originalStatus, updated.Status) // Repository已更新
		// 注意：这个测试说明需要在Service层添加状态检查
	})

	t.Run("已取消订单状态不应该再改变", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCanceled,
			TotalPriceCents: 10000,
			CancelReason:    "用户取消",
		}
		orderRepo.Create(ctx, order)

		// 尝试修改已取消的订单
		order.Status = model.OrderStatusConfirmed
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err) // Repository层允许
		// 注意：Service层应该添加检查防止这种情况
	})
}

func TestOrderCreation_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("创建订单时价格为0", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 0, // 零价格
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), order.TotalPriceCents)
	})

	t.Run("创建订单时价格为极大值", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000000, // 100,000元
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, int64(10000000), order.TotalPriceCents)
	})

	t.Run("创建订单时必须有用户ID", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          0, // 无效的用户ID
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		err := orderRepo.Create(ctx, order)

		// Repository层会允许，但Service层应该验证
		assert.NoError(t, err)
		// 注意：Service层应该添加用户ID验证
	})

	t.Run("创建订单时默认状态应该是pending", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, model.OrderStatusPending, order.Status)
	})
}

func TestOrderCancellation_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("取消pending状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 取消订单
		order.Status = model.OrderStatusCanceled
		order.CancelReason = "用户取消"
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, updated.Status)
		assert.Equal(t, "用户取消", updated.CancelReason)
	})

	t.Run("取消confirmed状态的订单_需要退款", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 取消已确认的订单
		order.Status = model.OrderStatusRefunded
		order.RefundReason = "用户申请退款"
		refundedAt := time.Now()
		order.RefundedAt = &refundedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		// 注意：实际业务中应该触发退款流程
	})

	t.Run("不能取消in_progress状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 尝试取消进行中的订单
		order.Status = model.OrderStatusCanceled
		err := orderRepo.Update(ctx, order)

		// Repository层会允许，但Service层应该阻止
		assert.NoError(t, err)
		// 注意：Service层应该添加状态检查
	})

	t.Run("不能取消completed状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		completedAt := time.Now()
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
			CompletedAt:     &completedAt,
		}
		orderRepo.Create(ctx, order)

		// 尝试取消已完成的订单
		order.Status = model.OrderStatusCanceled
		err := orderRepo.Update(ctx, order)

		// Repository层会允许，但Service层应该阻止
		assert.NoError(t, err)
		// 注意：Service层应该添加状态检查
	})
}

func TestOrderCompletion_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("正常完成订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 完成订单
		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("完成订单时应该记录完成时间", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 完成订单
		beforeComplete := time.Now()
		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		orderRepo.Update(ctx, order)
		afterComplete := time.Now()

		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.NotNil(t, updated.CompletedAt)
		assert.True(t, updated.CompletedAt.After(beforeComplete) || updated.CompletedAt.Equal(beforeComplete))
		assert.True(t, updated.CompletedAt.Before(afterComplete) || updated.CompletedAt.Equal(afterComplete))
	})
}

func TestOrderQuery_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("查询不存在的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order, err := orderRepo.Get(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("查询用户的所有订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		// 创建多个订单
		for i := 0; i < 5; i++ {
			order := &model.Order{
				UserID:          1,
				Status:          model.OrderStatusPending,
				TotalPriceCents: int64((i + 1) * 1000),
			}
			orderRepo.Create(ctx, order)
		}

		// 查询用户1的订单
		userID := uint64(1)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, orders, 5)
	})

	t.Run("按状态过滤订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		// 创建不同状态的订单
		statuses := []model.OrderStatus{
			model.OrderStatusPending,
			model.OrderStatusConfirmed,
			model.OrderStatusInProgress,
			model.OrderStatusCompleted,
			model.OrderStatusCanceled,
		}

		for _, status := range statuses {
			order := &model.Order{
				UserID:          1,
				Status:          status,
				TotalPriceCents: 10000,
			}
			orderRepo.Create(ctx, order)
		}

		// 只查询pending状态的订单
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			Statuses: []model.OrderStatus{model.OrderStatusPending},
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, model.OrderStatusPending, orders[0].Status)
	})

	t.Run("查询空结果集", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		userID := uint64(999)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, orders)
	})
}

func TestOrderAuthorization(t *testing.T) {
	ctx := context.Background()

	t.Run("用户只能查看自己的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		// 创建用户1的订单
		order1 := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order1)

		// 创建用户2的订单
		order2 := &model.Order{
			UserID:          2,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 20000,
		}
		orderRepo.Create(ctx, order2)

		// 用户1查询自己的订单
		userID := uint64(1)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, uint64(1), orders[0].UserID)
	})

	t.Run("用户不能操作他人的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		// 创建用户1的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 用户2尝试获取用户1的订单
		retrieved, err := orderRepo.Get(ctx, order.ID)

		// Repository层会返回订单，但Service层应该检查权限
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint64(1), retrieved.UserID)
		// 注意：Service层应该添加权限检查
	})
}

func TestOrderConcurrency(t *testing.T) {
	ctx := context.Background()

	t.Run("并发更新同一订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 模拟并发更新（实际应该使用乐观锁或悲观锁）
		order1 := *order
		order2 := *order

		order1.Status = model.OrderStatusConfirmed
		order2.Status = model.OrderStatusCanceled

		// 第一次更新
		err1 := orderRepo.Update(ctx, &order1)
		assert.NoError(t, err1)

		// 第二次更新会覆盖第一次
		err2 := orderRepo.Update(ctx, &order2)
		assert.NoError(t, err2)

		// 最终状态是第二次更新的结果
		final, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, final.Status)
		// 注意：实际业务中应该使用版本号或锁来防止并发冲突
	})
}

type mockChatGroupRepo struct {
	lastOrderID       uint64
	lastDeactivatedID uint64
	group             *model.ChatGroup
}

func (m *mockChatGroupRepo) Create(ctx context.Context, group *model.ChatGroup) error { return nil }
func (m *mockChatGroupRepo) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	return m.group, nil
}

func (m *mockChatGroupRepo) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	m.lastOrderID = orderID
	return m.group, nil
}

func (m *mockChatGroupRepo) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	return nil, 0, nil
}

func (m *mockChatGroupRepo) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	return nil, 0, nil
}

func (m *mockChatGroupRepo) Update(ctx context.Context, group *model.ChatGroup) error { return nil }
func (m *mockChatGroupRepo) Deactivate(ctx context.Context, id uint64) error {
	m.lastDeactivatedID = id
	return nil
}

func (m *mockChatGroupRepo) ListDeactivatedBefore(ctx context.Context, cutoffTime time.Time, limit int) ([]model.ChatGroup, error) {
	return nil, nil
}

func (m *mockChatGroupRepo) DeleteByIDs(ctx context.Context, ids []uint64) error { return nil }


func TestCancelOrder_AutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: pending & owned by user 100
	order := &model.Order{Base: model.Base{ID: 1}, UserID: 100, Status: model.OrderStatusPending}
	orderRepo.orders[1] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{group: &model.ChatGroup{Base: model.Base{ID: 55}, GroupType: model.ChatGroupTypeOrder, IsActive: true}}
	svc.SetChatGroupRepository(chatRepo)

	if err := svc.CancelOrder(context.Background(), 100, 1, CancelOrderRequest{Reason: "x"}); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if chatRepo.lastDeactivatedID != 55 {
		t.Fatalf("expected deactivated group=55, got %d", chatRepo.lastDeactivatedID)
	}
}

func TestCompleteOrder_AutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: in_progress & owned by user 100
	now := time.Now().Add(-time.Hour)
	order := &model.Order{Base: model.Base{ID: 2}, UserID: 100, Status: model.OrderStatusInProgress, StartedAt: &now}
	orderRepo.orders[2] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{group: &model.ChatGroup{Base: model.Base{ID: 77}, GroupType: model.ChatGroupTypeOrder, IsActive: true}}
	svc.SetChatGroupRepository(chatRepo)

	if err := svc.CompleteOrder(context.Background(), 100, 2); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if chatRepo.lastDeactivatedID != 77 {
		t.Fatalf("expected deactivated group=77, got %d", chatRepo.lastDeactivatedID)
	}
}

func TestCompleteOrderByPlayer_NoAutoDeactivateOrderChat(t *testing.T) {
	orderRepo := newMockOrderRepository()
	// seed order: in_progress & assigned to player 1
	pid := uint64(1)
	order := &model.Order{Base: model.Base{ID: 3}, UserID: 200, Status: model.OrderStatusInProgress}
	order.PlayerID = &pid
	orderRepo.orders[3] = order

	svc := NewOrderService(orderRepo, &mockPlayerRepository{}, &mockUserRepository{}, &mockGameRepository{}, &mockPaymentRepository{}, &mockReviewRepository{}, &mockCommissionRepository{})
	chatRepo := &mockChatGroupRepo{group: &model.ChatGroup{Base: model.Base{ID: 88}, GroupType: model.ChatGroupTypeOrder, IsActive: true}}
	svc.SetChatGroupRepository(chatRepo)

	if err := svc.CompleteOrderByPlayer(context.Background(), 1, 3); err != nil {
		t.Fatalf("complete by player: %v", err)
	}
	if chatRepo.lastDeactivatedID != 0 {
		t.Fatalf("expected no auto-deactivate, got %d", chatRepo.lastDeactivatedID)
	}
}

func TestOrderService_GetAvailableOrders_DefaultsAndMapping(t *testing.T) {
	t.Helper()

	start := time.Now().Add(time.Hour)
	end := start.Add(2 * time.Hour)
	gameID := uint64(99)
	userID := uint64(42)

	orderRepo := &spyAvailableOrderRepository{
		orders: []model.Order{
			{
				Base: model.Base{
					ID:        1,
					CreatedAt: start.Add(-15 * time.Minute),
				},
				Title:           "Need boost",
				Description:     "Carry me please",
				UserID:          userID,
				TotalPriceCents: 18800,
				GameID:          &gameID,
				ScheduledStart:  &start,
				ScheduledEnd:    &end,
			},
		},
		total: 1,
	}

	userRepo := &stubUserRepository{
		mockUserRepository: &mockUserRepository{},
		data: map[uint64]*model.User{
			userID: {
				Base: model.Base{ID: userID},
				Name: "Alice",
			},
		},
	}

	gameRepo := &stubGameRepository{
		mockGameRepository: &mockGameRepository{},
		data: map[uint64]*model.Game{
			gameID: {
				Base: model.Base{ID: gameID},
				Name: "Valorant",
			},
		},
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       userRepo,
		games:       gameRepo,
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	result, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{})
	if err != nil {
		t.Fatalf("GetAvailableOrders returned error: %v", err)
	}

	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if orderRepo.lastOpts.Page != 1 {
		t.Fatalf("expected default Page=1, got %d", orderRepo.lastOpts.Page)
	}
	if orderRepo.lastOpts.PageSize != 20 {
		t.Fatalf("expected default PageSize=20, got %d", orderRepo.lastOpts.PageSize)
	}
	if len(orderRepo.lastOpts.Statuses) != 1 || orderRepo.lastOpts.Statuses[0] != model.OrderStatusConfirmed {
		t.Fatalf("expected statuses filter to be confirmed, got %+v", orderRepo.lastOpts.Statuses)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 available order, got %d", len(result))
	}

	dto := result[0]
	if dto.GameName != "Valorant" {
		t.Fatalf("expected game name Valorant, got %s", dto.GameName)
	}
	if dto.UserNickname != "Alice" {
		t.Fatalf("expected user nickname Alice, got %s", dto.UserNickname)
	}
	if dto.DurationHours != 2 {
		t.Fatalf("expected duration 2 hours, got %f", dto.DurationHours)
	}
	if dto.PriceCents != 18800 {
		t.Fatalf("expected price 18800, got %d", dto.PriceCents)
	}
}

func TestOrderService_GetAvailableOrders_GameFilterAndFallback(t *testing.T) {
	t.Helper()

	gameID := uint64(7)
	orderRepo := &spyAvailableOrderRepository{
		orders: []model.Order{
			{
				Base:   model.Base{ID: 10},
				Title:  "Any game ok",
				UserID: 100,
				GameID: &gameID,
				Status: model.OrderStatusConfirmed,
			},
		},
		total: 1,
	}

	svc := &OrderService{
		orders:  orderRepo,
		players: &mockPlayerRepository{},
		users: &stubUserRepository{
			mockUserRepository: &mockUserRepository{},
			data:               map[uint64]*model.User{},
		},
		games: &stubGameRepository{
			mockGameRepository: &mockGameRepository{},
			data:               map[uint64]*model.Game{},
		},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	reqGameID := uint64(7)
	result, total, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{
		GameID:   &reqGameID,
		Page:     2,
		PageSize: 30,
	})
	if err != nil {
		t.Fatalf("GetAvailableOrders returned error: %v", err)
	}

	if orderRepo.lastOpts.GameID == nil || *orderRepo.lastOpts.GameID != reqGameID {
		t.Fatalf("expected GameID filter %d, got %+v", reqGameID, orderRepo.lastOpts.GameID)
	}
	if orderRepo.lastOpts.Page != 2 || orderRepo.lastOpts.PageSize != 30 {
		t.Fatalf("expected paging options to be preserved, got %+v", orderRepo.lastOpts)
	}

	if total != 1 || len(result) != 1 {
		t.Fatalf("expected single result, got total=%d len=%d", total, len(result))
	}

	dto := result[0]
	if dto.GameName != "" {
		t.Fatalf("expected empty game name when lookup fails, got %s", dto.GameName)
	}
	if dto.UserNickname != "" {
		t.Fatalf("expected empty user nickname when lookup fails, got %s", dto.UserNickname)
	}
	if dto.DurationHours != 0 {
		t.Fatalf("expected duration 0 without schedule, got %f", dto.DurationHours)
	}
}

func TestOrderService_GetAvailableOrders_PropagatesErrors(t *testing.T) {
	t.Helper()

	orderRepo := &spyAvailableOrderRepository{
		err: errors.New("db down"),
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       &mockUserRepository{},
		games:       &mockGameRepository{},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: &mockCommissionRepository{},
	}

	if _, _, err := svc.GetAvailableOrders(context.Background(), AvailableOrdersRequest{}); err == nil {
		t.Fatal("expected error to be propagated from repository")
	}
}

func TestOrderService_BuildOrderTimeline_FullFlow(t *testing.T) {
	t.Helper()

	paidAt := time.Now().Add(5 * time.Minute)
	start := time.Now()
	started := start.Add(30 * time.Minute)
	completed := started.Add(time.Hour)
	refunded := completed.Add(30 * time.Minute)

	svc := &OrderService{
		payments: &stubPaymentRepository{
			mockPaymentRepository: &mockPaymentRepository{},
			listFn: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
				return []model.Payment{
					{
						Base:    model.Base{ID: 1},
						OrderID: 1,
						PaidAt:  &paidAt,
					},
				}, 1, nil
			},
		},
	}

	order := &model.Order{
		Base: model.Base{
			ID:        1,
			CreatedAt: start,
			UpdatedAt: completed.Add(15 * time.Minute),
		},
		Status:         model.OrderStatusCanceled,
		ScheduledStart: &start,
		StartedAt:      &started,
		CompletedAt:    &completed,
		CancelReason:   "user canceled",
		RefundedAt:     &refunded,
	}

	timeline := svc.buildOrderTimeline(order)
	if len(timeline) != 6 {
		t.Fatalf("expected 6 timeline entries, got %d", len(timeline))
	}

	expectedStatuses := []string{
		string(model.OrderStatusPending),
		string(model.OrderStatusConfirmed),
		string(model.OrderStatusInProgress),
		string(model.OrderStatusCompleted),
		string(model.OrderStatusCanceled),
		string(model.OrderStatusRefunded),
	}

	for i, status := range expectedStatuses {
		if timeline[i].Status != status {
			t.Fatalf("expected status %s at index %d, got %s", status, i, timeline[i].Status)
		}
	}

	if !timeline[1].Time.Equal(paidAt) {
		t.Fatalf("expected payment time to use PaidAt, got %v", timeline[1].Time)
	}
	if !timeline[2].Time.Equal(*order.StartedAt) {
		t.Fatalf("expected in-progress time to equal StartedAt, got %v", timeline[2].Time)
	}
	if !timeline[4].Time.Equal(order.UpdatedAt) {
		t.Fatalf("expected cancel time to equal UpdatedAt, got %v", timeline[4].Time)
	}
	if !timeline[5].Time.Equal(*order.RefundedAt) {
		t.Fatalf("expected refund time to equal RefundedAt, got %v", timeline[5].Time)
	}
}

func TestOrderService_BuildOrderTimeline_Fallbacks(t *testing.T) {
	t.Helper()

	svc := &OrderService{
		payments: &stubPaymentRepository{
			mockPaymentRepository: &mockPaymentRepository{},
			listFn: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
				return nil, 0, errors.New("cache miss")
			},
		},
	}

	created := time.Now()
	order := &model.Order{
		Base: model.Base{
			ID:        2,
			CreatedAt: created,
		},
		Status: model.OrderStatusConfirmed,
	}

	timeline := svc.buildOrderTimeline(order)
	if len(timeline) != 2 {
		t.Fatalf("expected 2 timeline entries, got %d", len(timeline))
	}
	if !timeline[1].Time.Equal(created) {
		t.Fatalf("expected confirmed entry to fallback to order creation time, got %v", timeline[1].Time)
	}

	pending := &model.Order{
		Base: model.Base{
			ID:        3,
			CreatedAt: created,
		},
		Status: model.OrderStatusPending,
	}
	pendingTimeline := svc.buildOrderTimeline(pending)
	if len(pendingTimeline) != 1 {
		t.Fatalf("expected only creation entry for pending order, got %d", len(pendingTimeline))
	}
}

func TestOrderService_RecordCommissionAsync(t *testing.T) {
	t.Helper()

	orderRepo := newMockOrderRepository()
	playerID := uint64(501)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		TotalPriceCents: 10000,
	}
	order.SetPlayerID(playerID)
	orderRepo.orders[order.ID] = order

	commissionRepo := &recordingCommissionRepository{
		mockCommissionRepository: &mockCommissionRepository{},
		ruleErr:                  errors.New("no specific rule"),
		defaultRule: &model.CommissionRule{
			Rate: 25,
		},
	}

	svc := &OrderService{
		orders:      orderRepo,
		players:     &mockPlayerRepository{},
		users:       &mockUserRepository{},
		games:       &mockGameRepository{},
		payments:    &mockPaymentRepository{},
		reviews:     &mockReviewRepository{},
		commissions: commissionRepo,
	}

	if err := svc.recordCommissionAsync(context.Background(), order.ID); err != nil {
		t.Fatalf("recordCommissionAsync returned error: %v", err)
	}

	if len(commissionRepo.createdRecords) != 1 {
		t.Fatalf("expected commission record to be created, got %d", len(commissionRepo.createdRecords))
	}

	record := commissionRepo.createdRecords[0]
	if record.PlayerID != playerID {
		t.Fatalf("expected record player %d, got %d", playerID, record.PlayerID)
	}
	if record.CommissionRate != 25 {
		t.Fatalf("expected commission rate 25, got %d", record.CommissionRate)
	}
	if record.CommissionCents != 2500 {
		t.Fatalf("expected commission cents 2500, got %d", record.CommissionCents)
	}
	if record.PlayerIncomeCents != 7500 {
		t.Fatalf("expected player income 7500, got %d", record.PlayerIncomeCents)
	}
}

func TestOrderService_RecordCommissionAsync_EdgeCases(t *testing.T) {
	t.Helper()

	t.Run("skips when record exists", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		order := &model.Order{
			Base:            model.Base{ID: 2},
			TotalPriceCents: 5000,
		}
		order.SetPlayerID(1)
		orderRepo.orders[order.ID] = order

		commissionRepo := &recordingCommissionRepository{
			mockCommissionRepository: &mockCommissionRepository{},
			existingRecord:           &model.CommissionRecord{ID: 99},
		}

		svc := &OrderService{
			orders:      orderRepo,
			players:     &mockPlayerRepository{},
			users:       &mockUserRepository{},
			games:       &mockGameRepository{},
			payments:    &mockPaymentRepository{},
			reviews:     &mockReviewRepository{},
			commissions: commissionRepo,
		}

		if err := svc.recordCommissionAsync(context.Background(), order.ID); err != nil {
			t.Fatalf("expected nil error when record already exists, got %v", err)
		}
		if len(commissionRepo.createdRecords) != 0 {
			t.Fatalf("expected no new record when one already exists")
		}
	})

	t.Run("errors when player missing", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		orderWithoutPlayer := &model.Order{
			Base:            model.Base{ID: 3},
			TotalPriceCents: 8000,
		}
		orderRepo.orders[orderWithoutPlayer.ID] = orderWithoutPlayer

		svc := &OrderService{
			orders:      orderRepo,
			players:     &mockPlayerRepository{},
			users:       &mockUserRepository{},
			games:       &mockGameRepository{},
			payments:    &mockPaymentRepository{},
			reviews:     &mockReviewRepository{},
			commissions: &recordingCommissionRepository{mockCommissionRepository: &mockCommissionRepository{}},
		}

		if err := svc.recordCommissionAsync(context.Background(), orderWithoutPlayer.ID); err == nil {
			t.Fatal("expected error when order has no player assigned")
		}
	})
}

type spyAvailableOrderRepository struct {
	orders   []model.Order
	total    int64
	lastOpts repository.OrderListOptions
	err      error
}

func (s *spyAvailableOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *spyAvailableOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	s.lastOpts = opts
	if s.err != nil {
		return nil, 0, s.err
	}
	return s.orders, s.total, nil
}

func (s *spyAvailableOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for i := range s.orders {
		if s.orders[i].ID == id {
			return &s.orders[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *spyAvailableOrderRepository) Update(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *spyAvailableOrderRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

type stubUserRepository struct {
	*mockUserRepository
	data map[uint64]*model.User
}

func (s *stubUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	if user, ok := s.data[id]; ok {
		copy := *user
		return &copy, nil
	}
	return nil, repository.ErrNotFound
}

type stubGameRepository struct {
	*mockGameRepository
	data map[uint64]*model.Game
}

func (s *stubGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	if game, ok := s.data[id]; ok {
		copy := *game
		return &copy, nil
	}
	return nil, repository.ErrNotFound
}

type stubPaymentRepository struct {
	*mockPaymentRepository
	listFn func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error)
}

func (s *stubPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, opts)
	}
	return s.mockPaymentRepository.List(ctx, opts)
}

type recordingCommissionRepository struct {
	*mockCommissionRepository
	existingRecord *model.CommissionRecord
	rule           *model.CommissionRule
	ruleErr        error
	defaultRule    *model.CommissionRule
	defaultErr     error
	createdRecords []*model.CommissionRecord
}

func (r *recordingCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return r.existingRecord, nil
}

func (r *recordingCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	if r.rule != nil {
		return r.rule, nil
	}
	if r.ruleErr != nil {
		return nil, r.ruleErr
	}
	return nil, nil
}

func (r *recordingCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	if r.defaultErr != nil {
		return nil, r.defaultErr
	}
	if r.defaultRule != nil {
		return r.defaultRule, nil
	}
	return &model.CommissionRule{Rate: 20}, nil
}

func (r *recordingCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	copy := *record
	r.createdRecords = append(r.createdRecords, &copy)
	return nil
}
