// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/service/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "order_user")
	playerUser := CreateUniqueTestUser(t, db, "order_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_order")

	// Create a service item
	serviceItem := CreateTestServiceItem(t, db, testGame, "Default Service", 5000)

	// Set player hourly rate
	testPlayer.HourlyRateCents = 5000
	db.Save(testPlayer)

	// Create order
	now := time.Now().Add(time.Hour)
	serviceID := serviceItem.ID
	req := order.CreateOrderRequest{
		PlayerID:       testPlayer.ID,
		GameID:         testGame.ID,
		ServiceID:      &serviceID,
		Title:          "Test Order",
		Description:    "Test description",
		ScheduledStart: &now,
		DurationHours:  2,
	}

	resp, err := svc.CreateOrder(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.OrderID)
	assert.True(t, resp.NeedPayment)
	assert.Greater(t, resp.PriceCents, int64(0))
}

func TestOrderService_GetMyOrders(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "myorders_user")
	playerUser := CreateUniqueTestUser(t, db, "myorders_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_myorders")

	// Create multiple orders
	for i := 0; i < 5; i++ {
		CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	}
	for i := 0; i < 3; i++ {
		CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 5000)
	}

	// Get all orders
	resp, err := svc.GetMyOrders(ctx, testUser.ID, order.MyOrderListRequest{
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(8), resp.Total)
	assert.Len(t, resp.Orders, 8)

	// Get only completed orders
	resp, err = svc.GetMyOrders(ctx, testUser.ID, order.MyOrderListRequest{
		Status:   "completed",
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), resp.Total)
}

func TestOrderService_CancelOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_user")
	playerUser := CreateUniqueTestUser(t, db, "cancel_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_cancel")

	// Create pending order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Cancel order
	err := svc.CancelOrder(ctx, testUser.ID, testOrder.ID, order.CancelOrderRequest{
		Reason: "Test cancellation",
	})
	require.NoError(t, err)

	// Verify order is canceled
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status)
	assert.Equal(t, "Test cancellation", updatedOrder.CancelReason)
}

func TestOrderService_CancelOrder_Unauthorized(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_auth_user")
	otherUser := CreateUniqueTestUser(t, db, "cancel_auth_other")
	playerUser := CreateUniqueTestUser(t, db, "cancel_auth_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_cancel_auth")

	// Create order for testUser
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Try to cancel with different user
	err := svc.CancelOrder(ctx, otherUser.ID, testOrder.ID, order.CancelOrderRequest{
		Reason: "Unauthorized cancellation",
	})
	assert.Error(t, err)
}

func TestOrderService_CompleteOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create default commission rule
	CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "complete_user")
	playerUser := CreateUniqueTestUser(t, db, "complete_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_complete")

	// Create in_progress order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)

	// Complete order
	err := svc.CompleteOrder(ctx, testUser.ID, testOrder.ID)
	require.NoError(t, err)

	// Verify order is completed
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.CompletedAt)
}

func TestOrderService_AcceptOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "accept_user")
	playerUser := CreateUniqueTestUser(t, db, "accept_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_accept")

	// Create confirmed order (waiting for player to accept)
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	// Clear player ID to simulate unassigned order
	db.Model(&testOrder).Update("player_id", nil)

	// Accept order
	err := svc.AcceptOrder(ctx, playerUser.ID, testOrder.ID)
	require.NoError(t, err)

	// Verify order is in progress
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.PlayerID)
	assert.Equal(t, testPlayer.ID, *updatedOrder.PlayerID)
}

func TestOrderService_GetAvailableOrders(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "available_user")
	playerUser := CreateUniqueTestUser(t, db, "available_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_available")

	// Create confirmed orders (available for players)
	for i := 0; i < 5; i++ {
		CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	}

	// Get available orders
	orders, total, err := svc.GetAvailableOrders(ctx, order.AvailableOrdersRequest{
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, orders, 5)
}

func TestOrderService_GetOrderDetail(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	svc := order.NewOrderService(order.OrderDeps{Orders: orderRepo, Players: playerRepo, Users: userRepo, Games: gameRepo, Payments: paymentRepo, Reviews: reviewRepo, Commissions: commissionRepo})

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "detail_user")
	playerUser := CreateUniqueTestUser(t, db, "detail_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_detail")

	// Create order with payment
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	CreateTestPayment(t, db, testOrder, model.PaymentStatusPaid)
	CreateTestReview(t, db, testOrder, model.Rating5)

	// Get order detail
	resp, err := svc.GetOrderDetail(ctx, testUser.ID, testOrder.ID)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, testOrder.ID, resp.Order.ID)
	assert.NotNil(t, resp.Player)
	assert.NotNil(t, resp.Payment)
	assert.NotNil(t, resp.Review)
	assert.NotEmpty(t, resp.Timeline)
}
