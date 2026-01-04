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

// TestOrderService_CreateOrder_MinimumDuration tests order creation with minimum duration
func TestOrderService_CreateOrder_MinimumDuration(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
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
		Title:          "Minimum Duration Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  0.5, // Minimum allowed duration
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.OrderID)
}

// TestOrderService_CreateOrder_MaximumDuration tests order creation with maximum duration
func TestOrderService_CreateOrder_MaximumDuration(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
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
		Title:          "Maximum Duration Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  24, // Maximum allowed duration
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.OrderID)
}

// TestOrderService_calculateOrderPricing_ZeroHourlyRate tests pricing with zero hourly rate
func TestOrderService_calculateOrderPricing_ZeroHourlyRate(t *testing.T) {
	player := createTestPlayer(100, 200)
	player.HourlyRateCents = 0

	req := CreateOrderRequest{
		DurationHours: 2,
	}

	service := &OrderService{}

	totalPrice, commissionCents, playerIncomeCents := service.calculateOrderPricing(player, req)

	assert.Equal(t, int64(0), totalPrice)
	assert.Equal(t, int64(0), commissionCents)
	assert.Equal(t, int64(0), playerIncomeCents)
}

// TestOrderService_calculateOrderPricing_VeryHighRate tests pricing with very high hourly rate
func TestOrderService_calculateOrderPricing_VeryHighRate(t *testing.T) {
	player := createTestPlayer(100, 200)
	player.HourlyRateCents = 100000 // 1000 CNY/hour

	req := CreateOrderRequest{
		DurationHours: 5,
	}

	service := &OrderService{}

	totalPrice, commissionCents, playerIncomeCents := service.calculateOrderPricing(player, req)

	assert.Equal(t, int64(500000), totalPrice)       // 100000 * 5
	assert.Equal(t, int64(100000), commissionCents)  // 20% of 500000
	assert.Equal(t, int64(400000), playerIncomeCents) // 500000 - 100000
}

// TestOrderService_GetMyOrders_ZeroPage tests pagination with page 0
func TestOrderService_GetMyOrders_ZeroPage(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			// Should default to page 1
			assert.Equal(t, 1, opts.Page)
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

	req := MyOrderListRequest{
		Page: 0,
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Orders)
}

// TestOrderService_GetMyOrders_ExceedsMaxPageSize tests pagination with page size exceeding maximum
func TestOrderService_GetMyOrders_ExceedsMaxPageSize(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			// Should default to 20 (max allowed)
			assert.Equal(t, 20, opts.PageSize)
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

	req := MyOrderListRequest{
		PageSize: 200, // Exceeds max of 100
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestOrderService_CancelOrder_ConfirmedWithRefund tests canceling confirmed order triggers refund
func TestOrderService_CancelOrder_ConfirmedWithRefund(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(123)

	var updatedOrder *model.Order
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				UserID:      userID,
				Status:      model.OrderStatusConfirmed,
				TotalPriceCents: 10000,
			}, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	// Mock payment repository to return paid payment
	paymentTime := time.Now()
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{
					Base:        model.Base{ID: 1},
					Status:      model.PaymentStatusPaid,
					PaidAt:      &paymentTime,
				},
			}, 1, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CancelOrderRequest{
		Reason: "Changed mind",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(10000), updatedOrder.RefundAmountCents)
	assert.Equal(t, "用户取消订单", updatedOrder.RefundReason)
	assert.NotNil(t, updatedOrder.RefundedAt)
}

// TestOrderService_CancelOrder_PendingNoRefund tests canceling pending order doesn't trigger refund
func TestOrderService_CancelOrder_PendingNoRefund(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(123)

	var updatedOrder *model.Order
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				UserID:      userID,
				Status:      model.OrderStatusPending,
				TotalPriceCents: 10000,
			}, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{}, 0, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := CancelOrderRequest{
		Reason: "Changed mind",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status) // Not refunded
	assert.Equal(t, "Changed mind", updatedOrder.CancelReason)
	assert.Nil(t, updatedOrder.RefundedAt)
}

// TestOrderService_CompleteOrder_AlreadyCompleted tests completing already completed order
func TestOrderService_CompleteOrder_AlreadyCompleted(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(123)

	completedTime := time.Now()
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				UserID:      userID,
				Status:      model.OrderStatusCompleted,
				CompletedAt: &completedTime,
			}, nil
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
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_CompleteOrder_CannotCompletePendingOrder tests completing pending order is not allowed
func TestOrderService_CompleteOrder_CannotCompletePendingOrder(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(123)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:   model.Base{ID: id},
				UserID: userID,
				Status: model.OrderStatusPending,
			}, nil
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
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_AcceptOrder_OrderAlreadyInProgress tests accepting order that's already in progress
func TestOrderService_AcceptOrder_OrderAlreadyInProgress(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(100)
	orderID := uint64(123)

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, oid uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			// Simulate order already in progress (not confirmed)
			return false, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(1, 200), nil
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
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_AcceptOrder_UpdateFailure tests accepting order when update fails
func TestOrderService_AcceptOrder_UpdateFailure(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(100)
	orderID := uint64(123)

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, oid uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			return false, assert.AnError
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(1, 200), nil
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

// TestOrderService_CompleteOrderByPlayer_AlreadyCompleted tests player completing already completed order
func TestOrderService_CompleteOrderByPlayer_AlreadyCompleted(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(100)
	orderID := uint64(123)
	playerID := uint64(1)

	completedTime := time.Now()
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return &model.Order{
				Base:        model.Base{ID: id},
				PlayerID:    &playerID,
				Status:      model.OrderStatusCompleted,
				CompletedAt: &completedTime,
			}, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return &model.Player{Base: model.Base{ID: playerID}}, nil
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
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_GetAvailableOrders_WithPagination tests available orders with pagination
func TestOrderService_GetAvailableOrders_WithPagination(t *testing.T) {
	ctx := context.Background()
	gameID := uint64(1)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			assert.Equal(t, 2, opts.Page)
			assert.Equal(t, 10, opts.PageSize)
			assert.Equal(t, &gameID, opts.GameID)

			return []model.Order{
				{Base: model.Base{ID: 1}, Title: "Order 1"},
				{Base: model.Base{ID: 2}, Title: "Order 2"},
			}, 2, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return &model.User{Name: "Test User"}, nil
		},
	}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return &model.Game{Name: "Test Game"}, nil
		},
	}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		GameID:   &gameID,
		Page:     2,
		PageSize: 10,
	}

	availableOrders, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 2, len(availableOrders))
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "Order 1", availableOrders[0].Title)
	assert.Equal(t, "Order 2", availableOrders[1].Title)
}


// TestOrderService_buildOrderForCreation_CalculatesScheduledEnd tests that scheduled end is calculated correctly
func TestOrderService_buildOrderForCreation_CalculatesScheduledEnd(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)

	scheduledStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  2.5,
	}

	order := service.buildOrderForCreation(userID, req, 10000, 2000, 8000)

	assert.NotNil(t, order)
	assert.NotNil(t, order.ScheduledEnd)

	expectedEnd := scheduledStart.Add(2 * time.Hour + 30*time.Minute)
	assert.WithinDuration(t, expectedEnd, *order.ScheduledEnd, time.Second)
}

// TestOrderService_recordCommissionAsync_WithExistingCommission tests that existing commission is not duplicated
func TestOrderService_recordCommissionAsync_WithExistingCommission(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(123)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			playerID := uint64(100)
			return &model.Order{
				Base:             model.Base{ID: id},
				PlayerID:         &playerID,
				TotalPriceCents:  10000,
			}, nil
		},
	}

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, oid uint64) (*model.CommissionRecord, error) {
			// Return existing record
			return &model.CommissionRecord{
				ID:      1,
				OrderID: oid,
			}, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Should not error and should not create duplicate
	err := service.recordCommissionAsync(ctx, orderID)

	assert.NoError(t, err)
}

// TestOrderService_toOrderCardDTO_WithCompletedOrder tests DTO conversion for completed order
func TestOrderService_toOrderCardDTO_WithCompletedOrder(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	completedTime := time.Now()
	order := &model.Order{
		Base:         model.Base{ID: 123, CreatedAt: time.Now()},
		UserID:       userID,
		Status:       model.OrderStatusCompleted,
		CompletedAt:  &completedTime,
		TotalPriceCents: 10000,
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:     model.Base{ID: id},
				Nickname: "Test Player",
			}, nil
		},
	}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return &model.User{AvatarURL: "avatar.jpg"}, nil
		},
	}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return &model.Game{Name: "Test Game"}, nil
		},
	}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				{Base: model.Base{ID: 1}},
			}, 1, nil
		},
	}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, order, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, order.ID, card.ID)
	assert.Equal(t, model.OrderStatusCompleted, card.Status)
	assert.False(t, card.CanComplete) // Already completed, cannot complete again
	assert.False(t, card.CanReview)   // Already reviewed
}

// TestOrderService_toOrderCardDTO_WithInProgressOrder tests DTO conversion for in-progress order
func TestOrderService_toOrderCardDTO_WithInProgressOrder(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	startTime := time.Now()
	order := &model.Order{
		Base:         model.Base{ID: 123, CreatedAt: time.Now()},
		UserID:       userID,
		Status:       model.OrderStatusInProgress,
		StartedAt:    &startTime,
		TotalPriceCents: 10000,
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, order, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, model.OrderStatusInProgress, card.Status)
	assert.True(t, card.CanComplete)  // In progress orders can be completed
	assert.False(t, card.CanReview)   // Can't review until completed
}

// TestOrderService_validateCreateOrder_GameRepositoryError tests validation when game repository errors
func TestOrderService_validateCreateOrder_GameRepositoryError(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)
	gameID := uint64(1)

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return nil, assert.AnError
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

// TestOrderService_GetOrderDetail_MissingPaymentInfo tests order detail when payment info is missing
func TestOrderService_GetOrderDetail_MissingPaymentInfo(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(123)
	gameID := uint64(5)

	playerID := uint64(100)
	order := &model.Order{
		Base:   model.Base{ID: orderID},
		UserID: userID,
		PlayerID: &playerID,
		GameID:   &gameID,
		Status:  model.OrderStatusCompleted,
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return order, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return &model.User{AvatarURL: "avatar.jpg"}, nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	// Return empty payment list
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{}, 0, nil
		},
	}

	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	detail, err := service.GetOrderDetail(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Nil(t, detail.Payment) // Payment should be nil
}
