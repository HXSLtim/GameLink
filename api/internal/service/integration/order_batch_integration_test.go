// Package integration provides integration tests for admin batch order operations.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/user"
	"gamelink/internal/service/admin"
)

// setupOrderBatchAdminService creates an AdminService instance for order batch tests.
func setupOrderBatchAdminService(t *testing.T, db *gorm.DB) *admin.AdminService {
	t.Helper()

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	orders := implementations.NewOrderRepository(db)
	payments := payment.NewPaymentRepository(db)
	serviceItems := serviceitem.NewServiceItemRepository(db)

	// For order batch tests, we don't need all repositories
	var roles repository.RoleRepository = nil
	var permissions repository.PermissionRepository = nil
	var menus repository.MenuRepository = nil
	var stats repository.StatsRepository = nil
	var wallets repository.WalletRepository = nil

	return admin.NewAdminService(games, users, players, orders, payments, roles,
		serviceItems, permissions, menus, stats, wallets, gamecategory.NewGameCategoryRepository(db), nil)
}

// ============================================================================
// BatchCancelOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchCancelOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_user")
	playerUser := CreateUniqueTestUser(t, db, "cancel_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_cancel")

	// Create pending orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch cancel
	response, err := svc.BatchCancelOrders(ctx, orderIDs, "Test reason", "Admin note")
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
	assert.Len(t, response.SuccessItems, 3)
	assert.Len(t, response.FailedItems, 0)

	// Verify orders are canceled
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusCanceled, order.Status)
		assert.Equal(t, "Test reason", order.CancelReason)
	}
}

func TestAdminOrderBatch_BatchCancelOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_partial")
	playerUser := CreateUniqueTestUser(t, db, "cancel_partial_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_partial")

	// Create orders with different statuses
	pendingOrder1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	pendingOrder2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	completedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	inProgressOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)

	orderIDs := []uint64{pendingOrder1.ID, pendingOrder2.ID, completedOrder.ID, inProgressOrder.ID, 99999}

	// Batch cancel
	response, err := svc.BatchCancelOrders(ctx, orderIDs, "Test reason", "")
	require.NoError(t, err)
	assert.Equal(t, 5, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 3, response.FailedCount)
	assert.Len(t, response.SuccessItems, 2)
	assert.Len(t, response.FailedItems, 3)

	// Verify failed items contain correct error messages
	failedIDs := make(map[uint64]string)
	for _, item := range response.FailedItems {
		failedIDs[item.ID] = item.Message
	}
	assert.Contains(t, failedIDs, completedOrder.ID)
	assert.Contains(t, failedIDs[completedOrder.ID], "cannot cancel order with status")
	assert.Contains(t, failedIDs, inProgressOrder.ID)
	assert.Contains(t, failedIDs, inProgressOrder.ID, "cannot cancel order with status")
	assert.Contains(t, failedIDs, uint64(99999))
	assert.Contains(t, failedIDs[99999], "order not found")

	// Verify pending orders are canceled
	var canceledOrder1, canceledOrder2 model.Order
	require.NoError(t, db.First(&canceledOrder1, pendingOrder1.ID).Error)
	assert.Equal(t, model.OrderStatusCanceled, canceledOrder1.Status)
	require.NoError(t, db.First(&canceledOrder2, pendingOrder2.ID).Error)
	assert.Equal(t, model.OrderStatusCanceled, canceledOrder2.Status)
}

func TestAdminOrderBatch_BatchCancelOrders_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Empty order IDs list
	response, err := svc.BatchCancelOrders(ctx, []uint64{}, "Test reason", "")
	require.NoError(t, err)
	assert.Equal(t, 0, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
}

func TestAdminOrderBatch_BatchCancelOrders_AllNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// All non-existent order IDs
	orderIDs := []uint64{99999, 99998, 99997}

	response, err := svc.BatchCancelOrders(ctx, orderIDs, "Test reason", "")
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 3, response.FailedCount)
	assert.Len(t, response.FailedItems, 3)

	// Verify all failed items have "order not found" message
	for _, item := range response.FailedItems {
		assert.Contains(t, item.Message, "order not found")
	}
}

// ============================================================================
// BatchConfirmOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchConfirmOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "confirm_user")
	playerUser := CreateUniqueTestUser(t, db, "confirm_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_confirm")

	// Create pending orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch confirm
	response, err := svc.BatchConfirmOrders(ctx, orderIDs, "Admin confirmed")
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify orders are confirmed
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusConfirmed, order.Status)
	}
}

func TestAdminOrderBatch_BatchConfirmOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "confirm_partial")
	playerUser := CreateUniqueTestUser(t, db, "confirm_partial_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_confirm_partial")

	// Create orders with different statuses
	pendingOrder1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	pendingOrder2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	completedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	orderIDs := []uint64{pendingOrder1.ID, pendingOrder2.ID, completedOrder.ID, 88888}

	response, err := svc.BatchConfirmOrders(ctx, orderIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 4, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 2, response.FailedCount)
}

func TestAdminOrderBatch_BatchConfirmOrders_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "confirm_invalid")
	playerUser := CreateUniqueTestUser(t, db, "confirm_invalid_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_confirm_invalid")

	// Create orders with statuses that cannot be confirmed
	inProgressOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)
	canceledOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)

	orderIDs := []uint64{inProgressOrder.ID, canceledOrder.ID}

	response, err := svc.BatchConfirmOrders(ctx, orderIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 2, response.FailedCount)

	// Verify error messages
	for _, item := range response.FailedItems {
		assert.Contains(t, item.Message, "cannot confirm order with status")
	}
}

// ============================================================================
// BatchCompleteOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchCompleteOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create default commission rule
	CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "complete_user")
	playerUser := CreateUniqueTestUser(t, db, "complete_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_complete")

	// Create in_progress orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch complete
	response, err := svc.BatchCompleteOrders(ctx, orderIDs, "Admin completed")
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify orders are completed
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusCompleted, order.Status)
		assert.NotNil(t, order.CompletedAt)
	}
}

func TestAdminOrderBatch_BatchCompleteOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create default commission rule
	CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "complete_partial")
	playerUser := CreateUniqueTestUser(t, db, "complete_partial_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_complete_partial")

	// Create orders with different statuses
	inProgressOrder1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)
	inProgressOrder2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)
	pendingOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	orderIDs := []uint64{inProgressOrder1.ID, inProgressOrder2.ID, pendingOrder.ID}

	response, err := svc.BatchCompleteOrders(ctx, orderIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 1, response.FailedCount)

	// Verify error message for pending order
	failedIDs := make(map[uint64]string)
	for _, item := range response.FailedItems {
		failedIDs[item.ID] = item.Message
	}
	assert.Contains(t, failedIDs, pendingOrder.ID)
	assert.Contains(t, failedIDs[pendingOrder.ID], "cannot complete order with status")
}

func TestAdminOrderBatch_BatchCompleteOrders_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "complete_invalid")
	playerUser := CreateUniqueTestUser(t, db, "complete_invalid_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_complete_invalid")

	// Create orders with statuses that cannot be completed
	pendingOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	completedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	orderIDs := []uint64{pendingOrder.ID, completedOrder.ID}

	response, err := svc.BatchCompleteOrders(ctx, orderIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 2, response.FailedCount)
}

// ============================================================================
// BatchRefundOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchRefundOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "refund_user")
	playerUser := CreateUniqueTestUser(t, db, "refund_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_refund")

	// Create completed orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch refund
	input := admin.BatchRefundInput{
		Reason: "Customer request",
		Note:   "Approved by admin",
	}
	response, err := svc.BatchRefundOrders(ctx, orderIDs, input)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify orders are refunded
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusRefunded, order.Status)
		assert.Equal(t, "Customer request", order.RefundReason)
		assert.NotNil(t, order.RefundedAt)
	}
}

func TestAdminOrderBatch_BatchRefundOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "refund_partial")
	playerUser := CreateUniqueTestUser(t, db, "refund_partial_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_refund_partial")

	// Create orders with different statuses
	completedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	CreateTestPayment(t, db, completedOrder, model.PaymentStatusPaid)

	inProgressOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)
	CreateTestPayment(t, db, inProgressOrder, model.PaymentStatusPaid)

	confirmedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	CreateTestPayment(t, db, confirmedOrder, model.PaymentStatusPaid)

	canceledOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)

	orderIDs := []uint64{completedOrder.ID, inProgressOrder.ID, confirmedOrder.ID, canceledOrder.ID}

	input := admin.BatchRefundInput{
		Reason: "Test refund",
	}
	response, err := svc.BatchRefundOrders(ctx, orderIDs, input)
	require.NoError(t, err)
	assert.Equal(t, 4, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 1, response.FailedCount)

	// Verify canceled order failed
	failedIDs := make(map[uint64]string)
	for _, item := range response.FailedItems {
		failedIDs[item.ID] = item.Message
	}
	assert.Contains(t, failedIDs, canceledOrder.ID)
	assert.Contains(t, failedIDs[canceledOrder.ID], "cannot refund order with status")
}

func TestAdminOrderBatch_BatchRefundOrders_WithCustomAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "refund_amount")
	playerUser := CreateUniqueTestUser(t, db, "refund_amount_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_refund_amount")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	// Refund with custom amount
	customAmount := int64(5000) // Half refund
	input := admin.BatchRefundInput{
		Reason:      "Partial refund",
		AmountCents: &customAmount,
	}
	response, err := svc.BatchRefundOrders(ctx, []uint64{order.ID}, input)
	require.NoError(t, err)
	assert.Equal(t, 1, response.SuccessCount)

	// Verify partial refund amount
	var updatedOrder model.Order
	err = db.First(&updatedOrder, order.ID).Error
	require.NoError(t, err)
	assert.Equal(t, customAmount, updatedOrder.RefundAmountCents)
}

// ============================================================================
// BatchDeleteOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchDeleteOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "delete_user")
	playerUser := CreateUniqueTestUser(t, db, "delete_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_delete")

	// Create orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch delete (soft delete)
	response, err := svc.BatchDeleteOrders(ctx, orderIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify orders are soft deleted (deleted_at should be set)
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.Unscoped().First(&order, orderID).Error
		require.NoError(t, err)
		assert.NotNil(t, order.DeletedAt)
	}
}

func TestAdminOrderBatch_BatchDeleteOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "delete_partial")
	playerUser := CreateUniqueTestUser(t, db, "delete_partial_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_delete_partial")

	// Create orders
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)

	orderIDs := []uint64{order1.ID, order2.ID, 77777}

	response, err := svc.BatchDeleteOrders(ctx, orderIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 1, response.FailedCount)

	// Verify failed item
	assert.Len(t, response.FailedItems, 1)
	assert.Equal(t, uint64(77777), response.FailedItems[0].ID)
	assert.Contains(t, response.FailedItems[0].Message, "delete failed")
}

// ============================================================================
// BatchUpdateOrderStatus Tests
// ============================================================================

func TestAdminOrderBatch_BatchUpdateOrderStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "update_status_user")
	playerUser := CreateUniqueTestUser(t, db, "update_status_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_update_status")

	// Create pending orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	now := time.Now()
	input := admin.BatchUpdateStatusInput{
		Status:      model.OrderStatusConfirmed,
		Note:        "Status updated by admin",
		StartedAt:   &now,
		CompletedAt: nil,
	}

	response, err := svc.BatchUpdateOrderStatus(ctx, orderIDs, input)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify status updates
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusConfirmed, order.Status)
	}
}

func TestAdminOrderBatch_BatchUpdateOrderStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "invalid_transition")
	playerUser := CreateUniqueTestUser(t, db, "invalid_transition_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_invalid_transition")

	// Create completed order and try to set to pending (invalid transition)
	completedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	input := admin.BatchUpdateStatusInput{
		Status: model.OrderStatusPending,
		Note:   "Invalid transition",
	}

	response, err := svc.BatchUpdateOrderStatus(ctx, []uint64{completedOrder.ID}, input)
	require.NoError(t, err)
	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 1, response.FailedCount)

	// Verify error message
	assert.Len(t, response.FailedItems, 1)
	assert.Contains(t, response.FailedItems[0].Message, "invalid status transition")
}

func TestAdminOrderBatch_BatchUpdateOrderStatus_MixedValidTransitions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "mixed_transition")
	playerUser := CreateUniqueTestUser(t, db, "mixed_transition_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_mixed_transition")

	// Create orders with different statuses
	pendingOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	confirmedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	// Try to move all to in_progress
	// pending -> confirmed (valid)
	// confirmed -> in_progress (valid)
	input := admin.BatchUpdateStatusInput{
		Status: model.OrderStatusInProgress,
	}

	orderIDs := []uint64{pendingOrder.ID, confirmedOrder.ID}

	response, err := svc.BatchUpdateOrderStatus(ctx, orderIDs, input)
	require.NoError(t, err)
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, 1, response.SuccessCount) // Only confirmed -> in_progress is valid
	assert.Equal(t, 1, response.FailedCount)  // pending -> in_progress is invalid

	// Verify confirmed order succeeded
	var updatedConfirmedOrder model.Order
	err = db.First(&updatedConfirmedOrder, confirmedOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedConfirmedOrder.Status)

	// Verify pending order stayed the same
	var updatedPendingOrder model.Order
	err = db.First(&updatedPendingOrder, pendingOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusPending, updatedPendingOrder.Status)
}

func TestAdminOrderBatch_BatchUpdateOrderStatus_WithTimestamps(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "timestamps_user")
	playerUser := CreateUniqueTestUser(t, db, "timestamps_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_timestamps")

	confirmedOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	// Create default commission rule
	CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	now := time.Now()
	input := admin.BatchUpdateStatusInput{
		Status:      model.OrderStatusCompleted,
		CompletedAt: &now,
		Note:        "Completed with timestamp",
	}

	response, err := svc.BatchUpdateOrderStatus(ctx, []uint64{confirmedOrder.ID}, input)
	require.NoError(t, err)
	assert.Equal(t, 1, response.SuccessCount)

	// Verify timestamps were set
	var updatedOrder model.Order
	err = db.First(&updatedOrder, confirmedOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.CompletedAt)
}

// ============================================================================
// BatchAssignOrders Tests
// ============================================================================

func TestAdminOrderBatch_BatchAssignOrders_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "assign_user")
	playerUser1 := CreateUniqueTestUser(t, db, "assign_player1")
	testPlayer1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "assign_player2")
	testPlayer2 := CreateTestPlayer(t, db, playerUser2)
	testGame := CreateTestGame(t, db, "test_game_assign")

	// Create pending orders without player assignment
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer1, testGame, model.OrderStatusPending, 10000)
		// Clear player ID to simulate unassigned order
		db.Model(&order).Update("player_id", nil)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch assign to player2
	response, err := svc.BatchAssignOrders(ctx, orderIDs, testPlayer2.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify player assignments
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.NotNil(t, order.PlayerID)
		assert.Equal(t, testPlayer2.ID, *order.PlayerID)
	}
}

func TestAdminOrderBatch_BatchAssignOrders_PlayerNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "assign_not_found")
	playerUser := CreateUniqueTestUser(t, db, "assign_not_found_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_assign_not_found")

	// Create pending orders
	var orderIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		db.Model(&order).Update("player_id", nil)
		orderIDs = append(orderIDs, order.ID)
	}

	// Try to assign to non-existent player
	response, err := svc.BatchAssignOrders(ctx, orderIDs, 99999)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 3, response.FailedCount)

	// Verify all failed with "player not found"
	for _, item := range response.FailedItems {
		assert.Contains(t, item.Message, "player not found")
	}
}

func TestAdminOrderBatch_BatchAssignOrders_PartialSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "assign_partial")
	playerUser1 := CreateUniqueTestUser(t, db, "assign_partial_player1")
	testPlayer1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "assign_partial_player2")
	testPlayer2 := CreateTestPlayer(t, db, playerUser2)
	testGame := CreateTestGame(t, db, "test_game_assign_partial")

	// Create pending orders
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer1, testGame, model.OrderStatusPending, 10000)
	db.Model(&order1).Update("player_id", nil)

	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer1, testGame, model.OrderStatusPending, 10000)
	db.Model(&order2).Update("player_id", nil)

	// Add non-existent order
	orderIDs := []uint64{order1.ID, order2.ID, 66666}

	// Batch assign
	response, err := svc.BatchAssignOrders(ctx, orderIDs, testPlayer2.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 1, response.FailedCount)

	// Verify failed item
	assert.Len(t, response.FailedItems, 1)
	assert.Equal(t, uint64(66666), response.FailedItems[0].ID)
}

func TestAdminOrderBatch_BatchAssignOrders_EmptyOrderList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create a valid player
	playerUser := CreateUniqueTestUser(t, db, "assign_empty")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Empty order IDs list
	response, err := svc.BatchAssignOrders(ctx, []uint64{}, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

func TestAdminOrderBatch_BatchOperations_EmptyLists(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Test all batch operations with empty lists
	cancelResp, err := svc.BatchCancelOrders(ctx, []uint64{}, "reason", "")
	require.NoError(t, err)
	assert.Equal(t, 0, cancelResp.TotalCount)

	confirmResp, err := svc.BatchConfirmOrders(ctx, []uint64{}, "")
	require.NoError(t, err)
	assert.Equal(t, 0, confirmResp.TotalCount)

	completeResp, err := svc.BatchCompleteOrders(ctx, []uint64{}, "")
	require.NoError(t, err)
	assert.Equal(t, 0, completeResp.TotalCount)

	deleteResp, err := svc.BatchDeleteOrders(ctx, []uint64{})
	require.NoError(t, err)
	assert.Equal(t, 0, deleteResp.TotalCount)

	updateInput := admin.BatchUpdateStatusInput{Status: model.OrderStatusConfirmed}
	updateResp, err := svc.BatchUpdateOrderStatus(ctx, []uint64{}, updateInput)
	require.NoError(t, err)
	assert.Equal(t, 0, updateResp.TotalCount)
}

func TestAdminOrderBatch_BatchOperations_AllOrdersNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	nonExistentIDs := []uint64{99999, 88888, 77777}

	// Test cancel
	cancelResp, err := svc.BatchCancelOrders(ctx, nonExistentIDs, "reason", "")
	require.NoError(t, err)
	assert.Equal(t, 3, cancelResp.TotalCount)
	assert.Equal(t, 0, cancelResp.SuccessCount)
	assert.Equal(t, 3, cancelResp.FailedCount)

	// Test confirm
	confirmResp, err := svc.BatchConfirmOrders(ctx, nonExistentIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 3, confirmResp.TotalCount)
	assert.Equal(t, 0, confirmResp.SuccessCount)

	// Test complete
	completeResp, err := svc.BatchCompleteOrders(ctx, nonExistentIDs, "")
	require.NoError(t, err)
	assert.Equal(t, 3, completeResp.TotalCount)
	assert.Equal(t, 0, completeResp.SuccessCount)

	// Test delete
	deleteResp, err := svc.BatchDeleteOrders(ctx, nonExistentIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, deleteResp.TotalCount)
	assert.Equal(t, 0, deleteResp.SuccessCount)
}

func TestAdminOrderBatch_BatchOperations_LargeBatch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "large_batch_user")
	playerUser := CreateUniqueTestUser(t, db, "large_batch_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_large_batch")

	// Create 50 pending orders
	var orderIDs []uint64
	for i := 0; i < 50; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		orderIDs = append(orderIDs, order.ID)
	}

	// Batch cancel 50 orders
	response, err := svc.BatchCancelOrders(ctx, orderIDs, "Large batch test", "")
	require.NoError(t, err)
	assert.Equal(t, 50, response.TotalCount)
	assert.Equal(t, 50, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify all are canceled
	for _, orderID := range orderIDs {
		var order model.Order
		err := db.First(&order, orderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusCanceled, order.Status)
	}
}
