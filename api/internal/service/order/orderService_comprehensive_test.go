package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
)

// TestOrderService_CreateOrder_TeamType tests team order creation
func TestOrderService_CreateOrder_TeamType(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	var createdOrder *model.Order
	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
			createdOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Team Order",
		Description:    "Need 2 players",
		ScheduledStart: &scheduledStart,
		DurationHours:  2,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdOrder)
	assert.Equal(t, uint64(123), resp.OrderID)
}

// TestOrderService_CreateOrder_WithServiceID tests order creation with service ID
func TestOrderService_CreateOrder_WithServiceID(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	serviceID := uint64(10)
	scheduledStart := time.Now().Add(time.Hour)

	var createdOrder *model.Order
	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
			createdOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceID,
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, serviceID, createdOrder.ItemID)
}

// TestOrderService_CreateOrder_DatabaseError tests order creation with database error
func TestOrderService_CreateOrder_DatabaseError(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			return assert.AnError
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestOrderService_CreateOrder_ValidateGameNotFound tests validation when game doesn't exist
func TestOrderService_CreateOrder_ValidateGameNotFound(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)
	gameID := uint64(999)

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := &OrderService{
		players: players,
		games:   games,
	}

	req := CreateOrderRequest{
		PlayerID: playerID,
		GameID:   gameID,
	}

	player, err := service.validateCreateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, player)
}

// TestOrderService_calculateOrderPricing_EdgeCases tests pricing calculation edge cases
func TestOrderService_calculateOrderPricing_EdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		hourlyRate         int64
		durationHours      float32
		expectedTotal      int64
		expectedCommission int64
		expectedIncome     int64
	}{
		{
			name:               "Minimum duration (0.5 hours)",
			hourlyRate:         5000,
			durationHours:      0.5,
			expectedTotal:      2500,
			expectedCommission: 500,
			expectedIncome:     2000,
		},
		{
			name:               "Maximum duration (24 hours)",
			hourlyRate:         5000,
			durationHours:      24,
			expectedTotal:      120000,
			expectedCommission: 24000,
			expectedIncome:     96000,
		},
		{
			name:               "Fractional duration (1.5 hours)",
			hourlyRate:         5000,
			durationHours:      1.5,
			expectedTotal:      7500,
			expectedCommission: 1500,
			expectedIncome:     6000,
		},
		{
			name:               "Low hourly rate",
			hourlyRate:         2000,
			durationHours:      1,
			expectedTotal:      2000,
			expectedCommission: 400,
			expectedIncome:     1600,
		},
		{
			name:               "High hourly rate",
			hourlyRate:         10000,
			durationHours:      1,
			expectedTotal:      10000,
			expectedCommission: 2000,
			expectedIncome:     8000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := createTestPlayer(100, 200)
			player.HourlyRateCents = tt.hourlyRate

			req := CreateOrderRequest{
				DurationHours: tt.durationHours,
			}

			service := &OrderService{}
			totalPrice, commissionCents, playerIncomeCents := service.calculateOrderPricing(player, req)

			assert.Equal(t, tt.expectedTotal, totalPrice)
			assert.Equal(t, tt.expectedCommission, commissionCents)
			assert.Equal(t, tt.expectedIncome, playerIncomeCents)
		})
	}
}

// TestOrderService_CancelOrder_PendingToCanceled tests cancellation from pending status
func TestOrderService_CancelOrder_PendingToCanceled(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	// CancelOrder now uses UpdateWithCondition (atomic update) instead of Update
	var capturedUpdates map[string]any
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateWithCondition: func(ctx context.Context, oID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			capturedUpdates = updates
			return true, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CancelOrderRequest{
		Reason: "No longer needed",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.NotNil(t, capturedUpdates)
	assert.Equal(t, model.OrderStatusCanceled, capturedUpdates["status"])
	assert.Equal(t, "No longer needed", capturedUpdates["cancel_reason"])
}

// TestOrderService_CancelOrder_NotFound tests cancellation of non-existent order
func TestOrderService_CancelOrder_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(999)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CancelOrderRequest{
		Reason: "Test",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	assert.Error(t, err)
}

// TestOrderService_CancelOrder_DatabaseError tests cancellation with database error
func TestOrderService_CancelOrder_DatabaseError(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateWithCondition: func(ctx context.Context, oID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			return false, assert.AnError
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CancelOrderRequest{
		Reason: "Test",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	assert.Error(t, err)
}

// TestOrderService_CompleteOrder_NotFound tests completion of non-existent order
func TestOrderService_CompleteOrder_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(999)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrder(ctx, userID, orderID)

	assert.Error(t, err)
}

// TestOrderService_CompleteOrder_DatabaseError tests completion with database error
func TestOrderService_CompleteOrder_DatabaseError(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusInProgress)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return assert.AnError
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrder(ctx, userID, orderID)

	assert.Error(t, err)
}

// TestOrderService_CompleteOrder_Unauthorized tests completion by unauthorized user
func TestOrderService_CompleteOrder_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	wrongUserID := uint64(999)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusInProgress)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrder(ctx, wrongUserID, orderID)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestOrderService_CompleteOrderByPlayer_NotFound tests player completion of non-existent order
func TestOrderService_CompleteOrderByPlayer_NotFound(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	orderID := uint64(999)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(100, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	assert.Error(t, err)
}

// TestOrderService_CompleteOrderByPlayer_DatabaseError tests player completion with database error
func TestOrderService_CompleteOrderByPlayer_DatabaseError(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return assert.AnError
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	assert.Error(t, err)
}

// TestOrderService_AcceptOrder_NotFound tests acceptance of non-existent order
func TestOrderService_AcceptOrder_NotFound(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	orderID := uint64(999)

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(100, playerUserID), nil
		},
	}

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, id uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			return false, repository.ErrNotFound
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.AcceptOrder(ctx, playerUserID, orderID)

	assert.Error(t, err)
}

// TestOrderService_GetOrderDetail_PlayerView tests order detail retrieval by player
func TestOrderService_GetOrderDetail_PlayerView(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(1) // Same as the test order user ID
	orderID := uint64(1)
	userID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusCompleted)
	now := time.Now()
	testOrder.CompletedAt = &now

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, playerUserID), nil
		},
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(100, playerUserID), nil
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	paymentTime := time.Now()
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, Method: model.PaymentMethodAlipay, AmountCents: 5000, Status: model.PaymentStatusPaid, PaidAt: &paymentTime},
			}, 1, nil
		},
	}

	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Player (as user) can view the order because they are the order user
	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, orderID, resp.Order.ID)
}

// TestOrderService_GetMyOrders_RepositoryError tests order listing with repository error
func TestOrderService_GetMyOrders_RepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return nil, 0, assert.AnError
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := MyOrderListRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestOrderService_GetAvailableOrders_WithGameFilter tests available orders with game filter
func TestOrderService_GetAvailableOrders_WithGameFilter(t *testing.T) {
	ctx := context.Background()
	gameID := uint64(5)

	var capturedOpts repoiface.OrderListOptions
	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			capturedOpts = opts
			return []model.Order{*createTestOrder(1, 1, model.OrderStatusConfirmed)}, 1, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		GameID:   &gameID,
		Page:     1,
		PageSize: 10,
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, ordersResp)
	assert.Equal(t, &gameID, capturedOpts.GameID)
	assert.Equal(t, int64(1), total)
}

// TestOrderService_GetAvailableOrders_RepositoryError tests available orders with repository error
func TestOrderService_GetAvailableOrders_RepositoryError(t *testing.T) {
	ctx := context.Background()

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return nil, 0, assert.AnError
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		Page:     1,
		PageSize: 10,
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, ordersResp)
	assert.Equal(t, int64(0), total)
}

// TestOrderService_GetAvailableOrders_DefaultPagination tests default pagination for available orders
func TestOrderService_GetAvailableOrders_DefaultPagination(t *testing.T) {
	ctx := context.Background()

	var capturedOpts repoiface.OrderListOptions
	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			capturedOpts = opts
			return []model.Order{}, 0, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		Page:     0, // Should default to 1
		PageSize: 0, // Should default to 20
	}

	_, _, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, capturedOpts.Page)
	assert.Equal(t, 20, capturedOpts.PageSize)
}

// TestOrderService_recordCommissionAsync_NoPlayerID tests commission recording with no player ID
func TestOrderService_recordCommissionAsync_NoPlayerID(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusCompleted)
	testOrder.PlayerID = nil // No player assigned

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.recordCommissionAsync(ctx, orderID)

	assert.Error(t, err)
}

// TestOrderService_recordCommissionAsync_OrderNotFound tests commission recording for non-existent order
func TestOrderService_recordCommissionAsync_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(999)

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.recordCommissionAsync(ctx, orderID)

	assert.Error(t, err)
}

// TestOrderService_recordCommissionAsync_CreateError tests commission recording with create error
func TestOrderService_recordCommissionAsync_CreateError(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 20}, nil
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			return assert.AnError
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.recordCommissionAsync(ctx, orderID)

	assert.Error(t, err)
}

// TestOrderService_recordCommissionAsync_WithPlayerSpecificRule tests commission recording with player-specific rule
func TestOrderService_recordCommissionAsync_WithPlayerSpecificRule(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	var createdRecord *model.CommissionRecord
	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 15}, nil // Player-specific rate
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			createdRecord = record
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.recordCommissionAsync(ctx, orderID)

	require.NoError(t, err)
	assert.NotNil(t, createdRecord)
	assert.Equal(t, 15, createdRecord.CommissionRate)
	assert.Equal(t, int64(750), createdRecord.CommissionCents) // 15% of 5000
}

// TestOrderService_toOrderCardDTO_NoPlayer tests order card DTO when no player assigned
func TestOrderService_toOrderCardDTO_NoPlayer(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrder := createTestOrder(1, userID, model.OrderStatusPending)
	testOrder.PlayerID = nil

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, testOrder, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "", card.PlayerNickname)
	assert.Equal(t, "", card.PlayerAvatar)
}

// TestOrderService_toOrderCardDTO_NoGame tests order card DTO when game not found
func TestOrderService_toOrderCardDTO_NoGame(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrder := createTestOrder(1, userID, model.OrderStatusPending)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return nil, repository.ErrNotFound
		},
	}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, testOrder, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "", card.GameName)
}

// TestOrderService_buildOrderTimeline_CanceledOrder tests timeline for canceled order
func TestOrderService_buildOrderTimeline_CanceledOrder(t *testing.T) {
	testOrder := createTestOrder(1, 1, model.OrderStatusCanceled)
	testOrder.CancelReason = "User request"

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	timeline := service.buildOrderTimeline(testOrder)

	assert.GreaterOrEqual(t, len(timeline), 2) // At least created and canceled
	cancelFound := false
	for _, item := range timeline {
		if item.Status == string(model.OrderStatusCanceled) {
			cancelFound = true
			assert.Contains(t, item.Message, "User request")
		}
	}
	assert.True(t, cancelFound, "Canceled status should be in timeline")
}

// TestOrderService_buildOrderTimeline_RefundedOrder tests timeline for refunded order
func TestOrderService_buildOrderTimeline_RefundedOrder(t *testing.T) {
	testOrder := createTestOrder(1, 1, model.OrderStatusRefunded)
	now := time.Now()
	testOrder.RefundedAt = &now

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	timeline := service.buildOrderTimeline(testOrder)

	assert.GreaterOrEqual(t, len(timeline), 2)
	refundFound := false
	for _, item := range timeline {
		if item.Status == string(model.OrderStatusRefunded) {
			refundFound = true
		}
	}
	assert.True(t, refundFound, "Refunded status should be in timeline")
}

// TestOrderService_buildOrderForCreation_WithoutServiceID tests order building without service ID
func TestOrderService_buildOrderForCreation_WithoutServiceID(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)

	scheduledStart := time.Now().Add(time.Hour)

	req := CreateOrderRequest{
		PlayerID: playerID,
		GameID:   gameID,
		// No ServiceID
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  2,
	}

	order := service.buildOrderForCreation(userID, req, 10000, 2000, 8000)

	assert.NotNil(t, order)
	assert.Equal(t, uint64(0), order.ItemID) // Should be 0 when no ServiceID provided
}
