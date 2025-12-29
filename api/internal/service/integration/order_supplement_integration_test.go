// Package integration provides supplementary integration tests for order service.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
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

// TestOrderService_CreateMultiPlayerOrder tests creating an order with multiple players.
func TestOrderService_CreateMultiPlayerOrder(t *testing.T) {
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
	commissionRepo := commission.NewCommissionRepository(db)

	svc := order.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "multi_order_user")
	playerUser1 := CreateUniqueTestUser(t, db, "multi_player1")
	testPlayer1 := CreateTestPlayer(t, db, playerUser1)
	testGame := CreateTestGame(t, db, "multi_game")
	serviceItem := CreateTestServiceItem(t, db, testGame, "Multi Service", 5000)

	// Set player hourly rate
	testPlayer1.HourlyRateCents = 5000
	db.Save(testPlayer1)

	// Create multi-player order
	now := time.Now().Add(time.Hour)
	serviceID := serviceItem.ID
	req := order.CreateOrderRequest{
		PlayerID:       testPlayer1.ID,
		GameID:         testGame.ID,
		ServiceID:      &serviceID,
		Title:          "Multi Player Order",
		Description:    "Order requiring multiple players",
		ScheduledStart: &now,
		DurationHours:  2,
	}

	resp, err := svc.CreateOrder(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.OrderID)

	// Update order to set RequiredPlayers (this is set at order level, not in request)
	var createdOrder model.Order
	err = db.First(&createdOrder, resp.OrderID).Error
	require.NoError(t, err)
	createdOrder.RequiredPlayers = 2
	createdOrder.CurrentPlayers = 0
	err = db.Save(&createdOrder).Error
	require.NoError(t, err)

	// Verify order has RequiredPlayers set
	err = db.First(&createdOrder, resp.OrderID).Error
	require.NoError(t, err)
	assert.Equal(t, 2, createdOrder.RequiredPlayers)
	assert.Equal(t, 0, createdOrder.CurrentPlayers)
}

// TestOrderService_OrderWithOrderItems tests order with multiple order items (seats).
func TestOrderService_OrderWithOrderItems(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "items_user")
	playerUser := CreateUniqueTestUser(t, db, "items_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "items_game")

	// Create order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	testOrder.RequiredPlayers = 2
	db.Save(testOrder)

	// Create order items (seats)
	item1 := CreateTestOrderItem(t, db, testOrder, 1, 5000, model.OrderItemStatusPending)
	item2 := CreateTestOrderItem(t, db, testOrder, 2, 5000, model.OrderItemStatusPending)

	// Verify items created
	var items []model.OrderItem
	err := db.Where("order_id = ?", testOrder.ID).Find(&items).Error
	require.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, item1.Slot, 1)
	assert.Equal(t, item2.Slot, 2)
}

// TestOrderService_PlayerJoinOrder tests a player joining an order.
func TestOrderService_PlayerJoinOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "join_user")
	playerUser := CreateUniqueTestUser(t, db, "join_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "join_game")

	// Create order with items
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	testOrder.RequiredPlayers = 2
	testOrder.CurrentPlayers = 0
	db.Save(testOrder)

	item := CreateTestOrderItem(t, db, testOrder, 1, 5000, model.OrderItemStatusPending)

	// Player joins order
	orderPlayer := CreateTestOrderPlayer(t, db, testOrder, item, testPlayer, model.OrderPlayerStatusJoined)

	// Update order item status
	item.PlayerID = &testPlayer.ID
	item.Status = model.OrderItemStatusMatched
	db.Save(item)

	// Update order current players
	testOrder.CurrentPlayers = 1
	db.Save(testOrder)

	// Verify
	var updatedOrder model.Order
	err := db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 1, updatedOrder.CurrentPlayers)

	var updatedItem model.OrderItem
	err = db.First(&updatedItem, item.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderItemStatusMatched, updatedItem.Status)
	assert.NotNil(t, updatedItem.PlayerID)

	// Verify order player record
	var op model.OrderPlayer
	err = db.First(&op, orderPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderPlayerStatusJoined, op.Status)
}

// TestOrderService_GiftOrder tests creating a gift order.
func TestOrderService_GiftOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "gift_user")
	playerUser := CreateUniqueTestUser(t, db, "gift_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "gift_game")

	// Create gift service item
	giftItem := CreateTestServiceItem(t, db, testGame, "Gift Item", 1000)
	giftItem.SubCategory = model.SubCategoryGift
	db.Save(giftItem)

	// Create gift order
	now := time.Now()
	testOrder := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:           "GIFT" + time.Now().Format("20060102150405"),
		UserID:            testUser.ID,
		ItemID:            giftItem.ID,
		RecipientPlayerID: &testPlayer.ID,
		Quantity:          1,
		UnitPriceCents:    1000,
		TotalPriceCents:   1000,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusCompleted, // Gift orders complete immediately
		GiftMessage:       "Happy Birthday!",
		IsAnonymous:       false,
		DeliveredAt:       &now,
		OrderConfig:       "{}",
	}
	err := db.Create(testOrder).Error
	require.NoError(t, err)

	// Verify gift order
	var savedOrder model.Order
	err = db.First(&savedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.True(t, savedOrder.IsGiftOrder())
	assert.Equal(t, model.OrderStatusCompleted, savedOrder.Status)
	assert.Equal(t, "Happy Birthday!", savedOrder.GiftMessage)
	assert.NotNil(t, savedOrder.RecipientPlayerID)
}

// TestOrderService_OrderStatusTransitions tests order status transitions.
func TestOrderService_OrderStatusTransitions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "status_user")
	playerUser := CreateUniqueTestUser(t, db, "status_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "status_game")

	tests := []struct {
		name       string
		fromStatus model.OrderStatus
		toStatus   model.OrderStatus
		valid      bool
	}{
		{"pending_to_confirmed", model.OrderStatusPending, model.OrderStatusConfirmed, true},
		{"confirmed_to_in_progress", model.OrderStatusConfirmed, model.OrderStatusInProgress, true},
		{"in_progress_to_completed", model.OrderStatusInProgress, model.OrderStatusCompleted, true},
		{"pending_to_canceled", model.OrderStatusPending, model.OrderStatusCanceled, true},
		{"completed_to_refunded", model.OrderStatusCompleted, model.OrderStatusRefunded, true},
		{"in_progress_to_disputed", model.OrderStatusInProgress, model.OrderStatusDisputed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, tt.fromStatus, 10000)

			// Transition status
			testOrder.Status = tt.toStatus
			if tt.toStatus == model.OrderStatusCompleted {
				now := time.Now()
				testOrder.CompletedAt = &now
			}
			err := db.Save(testOrder).Error
			require.NoError(t, err)

			// Verify
			var updatedOrder model.Order
			err = db.First(&updatedOrder, testOrder.ID).Error
			require.NoError(t, err)
			assert.Equal(t, tt.toStatus, updatedOrder.Status)
		})
	}
}

// TestOrderService_OrderWithDispute tests order with dispute flag.
func TestOrderService_OrderWithDispute(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "dispute_order_user")
	playerUser := CreateUniqueTestUser(t, db, "dispute_order_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "dispute_order_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute
	CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Update order dispute flag
	testOrder.HasDispute = true
	testOrder.Status = model.OrderStatusDisputed
	db.Save(testOrder)

	// Verify
	var updatedOrder model.Order
	err := db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.True(t, updatedOrder.HasDispute)
	assert.Equal(t, model.OrderStatusDisputed, updatedOrder.Status)
}

// TestOrderService_OrderCommissionCalculation tests commission calculation for orders.
func TestOrderService_OrderCommissionCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "commission_user")
	playerUser := CreateUniqueTestUser(t, db, "commission_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "commission_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create commission record
	record := CreateTestCommissionRecord(t, db, testOrder.ID, testPlayer.ID, 10000, model.SettlementStatusPending)

	// Verify commission calculation
	assert.Equal(t, int64(10000), record.TotalAmountCents)
	assert.Equal(t, 20, record.CommissionRate)
	assert.Equal(t, int64(2000), record.CommissionCents)
	assert.Equal(t, int64(8000), record.PlayerIncomeCents)
	assert.Equal(t, model.SettlementStatusPending, record.SettlementStatus)
}

// TestOrderService_OrderRefund tests order refund process.
func TestOrderService_OrderRefund(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "refund_order_user")
	playerUser := CreateUniqueTestUser(t, db, "refund_order_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "refund_order_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Process refund
	now := time.Now()
	testOrder.Status = model.OrderStatusRefunded
	testOrder.RefundAmountCents = testOrder.TotalPriceCents
	testOrder.RefundReason = "Customer requested refund"
	testOrder.RefundedAt = &now
	err := db.Save(testOrder).Error
	require.NoError(t, err)

	// Verify
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(10000), updatedOrder.RefundAmountCents)
	assert.NotNil(t, updatedOrder.RefundedAt)
}

// TestOrderService_OrderCancellation tests order cancellation.
func TestOrderService_OrderCancellation(t *testing.T) {
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
	commissionRepo := commission.NewCommissionRepository(db)

	svc := order.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_order_user")
	playerUser := CreateUniqueTestUser(t, db, "cancel_order_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "cancel_order_game")

	// Create pending order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Cancel order
	err := svc.CancelOrder(ctx, testUser.ID, testOrder.ID, order.CancelOrderRequest{
		Reason: "Changed my mind",
	})
	require.NoError(t, err)

	// Verify
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status)
	assert.Equal(t, "Changed my mind", updatedOrder.CancelReason)
}
