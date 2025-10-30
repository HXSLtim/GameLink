package order

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
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

func TestCreateOrder(t *testing.T) {
	svc := NewOrderService(
		newMockOrderRepository(),
		&mockPlayerRepository{},
		&mockUserRepository{},
		&mockGameRepository{},
		&mockPaymentRepository{},
		&mockReviewRepository{},
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
	)

	// Create a test order
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1, CreatedAt: now},
		UserID:         1,
		PlayerID:       1,
		GameID:         1,
		Title:          "Test Order",
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a test order
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         1,
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a test order owned by user 2
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2, // Different user
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a test order
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1, CreatedAt: now},
		UserID:         1,
		PlayerID:       1,
		GameID:         1,
		Title:          "Test Order",
		Description:    "Test description",
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create an in-progress order
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         1,
		PlayerID:       1,
		Status:         model.OrderStatusInProgress,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a pending order (can't complete directly)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         1,
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a confirmed order (ready to be accepted)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2,                          // Different user
		Status:         model.OrderStatusConfirmed, // Must be confirmed to accept
		PriceCents:     10000,
		ScheduledStart: &now,
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

	if updatedOrder.PlayerID == 0 {
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
	)

	// Create a pending order (not yet in progress)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         1,
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create an in-progress order assigned to player 1
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2,
		PlayerID:       1,
		Status:         model.OrderStatusInProgress,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create an order assigned to player 2
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         3,
		PlayerID:       2, // Different player
		Status:         model.OrderStatusInProgress,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a completed order (cannot be canceled)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         1,
		Status:         model.OrderStatusCompleted,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create an order for user 2
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2, // Different user
		Status:         model.OrderStatusInProgress,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create an order for user 2
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1, CreatedAt: now},
		UserID:         2, // Different user
		PlayerID:       3,
		GameID:         1,
		Title:          "Test Order",
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a pending order (not confirmed yet)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2,
		Status:         model.OrderStatusPending, // Wrong status for accepting
		PriceCents:     10000,
		ScheduledStart: &now,
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
	)

	// Create a pending order (not in-progress)
	now := time.Now()
	order := &model.Order{
		Base:           model.Base{ID: 1},
		UserID:         2,
		PlayerID:       1,
		Status:         model.OrderStatusPending, // Wrong status for completing
		PriceCents:     10000,
		ScheduledStart: &now,
	}
	orderRepo.orders[1] = order

	// Try to complete (should fail)
	err := svc.CompleteOrderByPlayer(context.Background(), 1, 1)

	if err == nil {
		t.Error("expected error for invalid status transition")
	}
}
