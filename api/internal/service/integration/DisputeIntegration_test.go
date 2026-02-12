// Package integration provides integration tests for dispute service.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/notification"
	"gamelink/internal/repository/operationlog"
	orderrepo "gamelink/internal/repository/order"
	paymentrepo "gamelink/internal/repository/payment"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/service/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisputeService_InitiateDispute(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	playerRepo := userrepo.NewPlayerRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeServiceWithPlayers(disputeRepo, orderRepo, userRepo, playerRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "dispute_user")
	playerUser := CreateUniqueTestUser(t, db, "dispute_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_dispute")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Initiate dispute as user
	req := order.InitiateDisputeRequest{
		OrderID:       testOrder.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "服务态度不好",
		EvidenceText:  "陪玩师在游戏中表现不专业",
		EvidenceURLs:  []string{"https://example.com/evidence1.jpg"},
	}

	resp, err := svc.InitiateDispute(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.DisputeID)
	assert.NotEmpty(t, resp.TraceID)
	assert.NotNil(t, resp.SLADeadline)

	// Verify order status is updated
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusDisputed, updatedOrder.Status)
	assert.True(t, updatedOrder.HasDispute)
}

func TestDisputeService_InitiateDispute_ByPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	playerRepo := userrepo.NewPlayerRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeServiceWithPlayers(disputeRepo, orderRepo, userRepo, playerRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "dispute_user2")
	playerUser := CreateUniqueTestUser(t, db, "dispute_player2")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_dispute2")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Initiate dispute as player - use playerUser.ID (user ID) since InitiatorID references User table
	req := order.InitiateDisputeRequest{
		OrderID:       testOrder.ID,
		InitiatorID:   playerUser.ID,
		InitiatorType: model.DisputeInitiatorPlayer,
		Type:          model.DisputeTypeUserNotCooperative,
		Reason:        "用户不配合",
		EvidenceText:  "用户在游戏中不听指挥",
	}

	resp, err := svc.InitiateDispute(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.DisputeID)
}

func TestDisputeService_InitiateDispute_Unauthorized(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "dispute_user3")
	otherUser := CreateUniqueTestUser(t, db, "dispute_other")
	playerUser := CreateUniqueTestUser(t, db, "dispute_player3")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_dispute3")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Try to initiate dispute as unrelated user
	req := order.InitiateDisputeRequest{
		OrderID:       testOrder.ID,
		InitiatorID:   otherUser.ID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        "Unauthorized dispute",
	}

	_, err := svc.InitiateDispute(ctx, req)
	assert.Error(t, err)
}

func TestDisputeService_AssignDispute(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "assign_user")
	playerUser := CreateUniqueTestUser(t, db, "assign_player")
	csUser := CreateUniqueTestUser(t, db, "assign_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_assign")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Assign dispute to customer service
	req := order.AssignDisputeRequest{
		DisputeID:         testDispute.ID,
		AssignedServiceID: csUser.ID,
		ActorUserID:       csUser.ID,
	}

	err := svc.AssignDispute(ctx, req)
	require.NoError(t, err)

	// Verify dispute is assigned
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusAssigned, updatedDispute.Status)
	assert.NotNil(t, updatedDispute.AssignedServiceID)
	assert.Equal(t, csUser.ID, *updatedDispute.AssignedServiceID)
}

func TestDisputeService_ResolveDispute_Refund(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "resolve_user")
	playerUser := CreateUniqueTestUser(t, db, "resolve_player")
	csUser := CreateUniqueTestUser(t, db, "resolve_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_resolve")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create and assign dispute
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusAssigned
	testDispute.AssignedServiceID = &csUser.ID
	db.Save(testDispute)

	// Resolve dispute with refund
	req := order.ResolveDisputeRequest{
		DisputeID:     testDispute.ID,
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "服务质量问题，同意退款",
		ActorUserID:   csUser.ID,
	}

	err := svc.ResolveDispute(ctx, req)
	require.NoError(t, err)

	// Verify dispute is resolved
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusResolved, updatedDispute.Status)
	assert.Equal(t, model.ResolutionRefund, updatedDispute.Resolution)
	assert.NotNil(t, updatedDispute.ResolvedAt)

	// Verify order is refunded
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(10000), updatedOrder.RefundAmountCents)
}

func TestDisputeService_ResolveDispute_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "reject_user")
	playerUser := CreateUniqueTestUser(t, db, "reject_player")
	csUser := CreateUniqueTestUser(t, db, "reject_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_reject")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create and assign dispute
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusAssigned
	testDispute.AssignedServiceID = &csUser.ID
	db.Save(testDispute)

	// Resolve dispute with rejection
	req := order.ResolveDisputeRequest{
		DisputeID:     testDispute.ID,
		Resolution:    model.ResolutionReject,
		ResolveRemark: "争议理由不成立，驳回",
		ActorUserID:   csUser.ID,
	}

	err := svc.ResolveDispute(ctx, req)
	require.NoError(t, err)

	// Verify dispute is rejected
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusRejected, updatedDispute.Status)
	assert.Equal(t, model.ResolutionReject, updatedDispute.Resolution)

	// Verify order is restored to completed
	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.False(t, updatedOrder.HasDispute)
}

func TestDisputeService_ListPendingDisputes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "pending_user")
	playerUser := CreateUniqueTestUser(t, db, "pending_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_pending")

	// Create multiple pending disputes
	for i := 0; i < 5; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	}

	// List pending disputes
	disputes, total, err := svc.ListPendingDisputes(ctx, 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, disputes, 5)
}

func TestDisputeService_GetDisputeDetail(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "detail_user")
	playerUser := CreateUniqueTestUser(t, db, "detail_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_detail_dispute")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Get dispute detail
	dispute, err := svc.GetDisputeDetail(ctx, testDispute.ID)
	require.NoError(t, err)
	require.NotNil(t, dispute)
	assert.Equal(t, testDispute.ID, dispute.ID)
	assert.Equal(t, testDispute.OrderID, dispute.OrderID)
	assert.Equal(t, testDispute.InitiatorID, dispute.InitiatorID)
}

func TestDisputeService_CheckSLABreaches(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "sla_user")
	playerUser := CreateUniqueTestUser(t, db, "sla_player")
	csUser := CreateUniqueTestUser(t, db, "sla_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_sla")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute with past SLA deadline
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	pastTime := time.Now().Add(-1 * time.Hour)
	testDispute.SLADeadline = &pastTime
	testDispute.Status = model.DisputeStatusAssigned
	testDispute.AssignedServiceID = &csUser.ID
	db.Save(testDispute)

	// Check and mark SLA breaches
	err := svc.CheckAndMarkSLABreaches(ctx)
	require.NoError(t, err)

	// Verify dispute is marked as SLA breached
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.True(t, updatedDispute.SLABreached)
	assert.NotNil(t, updatedDispute.SLABreachedAt)
}

// ============================================================================
// Batch Operations Tests
// ============================================================================

func TestDisputeService_BatchAssignDisputes_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_assign_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_assign_player")
	csUser := CreateUniqueTestUser(t, db, "batch_assign_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_assign")

	// Create multiple pending disputes
	var disputeIDs []uint64
	for i := 0; i < 3; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
		disputeIDs = append(disputeIDs, testDispute.ID)
	}

	// Batch assign disputes
	req := order.BatchAssignDisputesRequest{
		DisputeIDs:        disputeIDs,
		AssignedServiceID: csUser.ID,
		ActorUserID:       csUser.ID,
	}

	result, err := svc.BatchAssignDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.Errors)

	// Verify all disputes are assigned
	for _, disputeID := range disputeIDs {
		var dispute model.OrderDispute
		err = db.First(&dispute, disputeID).Error
		require.NoError(t, err)
		assert.Equal(t, model.DisputeStatusAssigned, dispute.Status)
		assert.NotNil(t, dispute.AssignedServiceID)
		assert.Equal(t, csUser.ID, *dispute.AssignedServiceID)
	}
}

func TestDisputeService_BatchAssignDisputes_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_assign_pf_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_assign_pf_player")
	csUser := CreateUniqueTestUser(t, db, "batch_assign_pf_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_assign_pf")

	// Create 2 pending disputes and 1 already assigned dispute
	var disputeIDs []uint64
	testOrder1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute1 := CreateTestDispute(t, db, testOrder1, testUser.ID, model.DisputeInitiatorUser)
	disputeIDs = append(disputeIDs, testDispute1.ID)

	testOrder2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute2 := CreateTestDispute(t, db, testOrder2, testUser.ID, model.DisputeInitiatorUser)
	disputeIDs = append(disputeIDs, testDispute2.ID)

	testOrder3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute3 := CreateTestDispute(t, db, testOrder3, testUser.ID, model.DisputeInitiatorUser)
	testDispute3.Status = model.DisputeStatusAssigned
	db.Save(testDispute3)
	disputeIDs = append(disputeIDs, testDispute3.ID)

	// Batch assign disputes
	req := order.BatchAssignDisputesRequest{
		DisputeIDs:        disputeIDs,
		AssignedServiceID: csUser.ID,
		ActorUserID:       csUser.ID,
	}

	result, err := svc.BatchAssignDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Message, "成功")
	assert.Contains(t, result.Message, "失败")
}

func TestDisputeService_BatchAssignDisputes_NonexistentDispute(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_assign_nd_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_assign_nd_player")
	csUser := CreateUniqueTestUser(t, db, "batch_assign_nd_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_assign_nd")

	// Create one pending dispute
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Batch assign with one valid and one invalid dispute ID
	req := order.BatchAssignDisputesRequest{
		DisputeIDs:        []uint64{testDispute.ID, 999999},
		AssignedServiceID: csUser.ID,
		ActorUserID:       csUser.ID,
	}

	result, err := svc.BatchAssignDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Error, "不存在")
}

func TestDisputeService_BatchUpdateDisputesStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_status_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_status_player")
	csUser := CreateUniqueTestUser(t, db, "batch_status_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_status")

	// Create multiple assigned disputes
	var disputeIDs []uint64
	for i := 0; i < 3; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
		testDispute.Status = model.DisputeStatusAssigned
		db.Save(testDispute)
		disputeIDs = append(disputeIDs, testDispute.ID)
	}

	// Batch update status to mediating
	req := order.BatchUpdateDisputesStatusRequest{
		DisputeIDs:  disputeIDs,
		Status:      model.DisputeStatusMediating,
		ActorUserID: csUser.ID,
	}

	result, err := svc.BatchUpdateDisputesStatus(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all disputes are updated
	for _, disputeID := range disputeIDs {
		var dispute model.OrderDispute
		err = db.First(&dispute, disputeID).Error
		require.NoError(t, err)
		assert.Equal(t, model.DisputeStatusMediating, dispute.Status)
	}
}

func TestDisputeService_BatchUpdateDisputesStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_status_it_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_status_it_player")
	csUser := CreateUniqueTestUser(t, db, "batch_status_it_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_status_it")

	// Create one pending dispute (cannot transition directly to mediating)
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Try to update status from pending to mediating (invalid transition)
	req := order.BatchUpdateDisputesStatusRequest{
		DisputeIDs:  []uint64{testDispute.ID},
		Status:      model.DisputeStatusMediating,
		ActorUserID: csUser.ID,
	}

	result, err := svc.BatchUpdateDisputesStatus(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Error, "无法从")
	assert.Contains(t, result.Errors[0].Error, "转换到")
}

func TestDisputeService_BatchUpdateDisputesStatus_Cancel(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_status_c_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_status_c_player")
	csUser := CreateUniqueTestUser(t, db, "batch_status_c_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_status_c")

	// Create assigned dispute
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusAssigned
	db.Save(testDispute)

	// Cancel dispute
	req := order.BatchUpdateDisputesStatusRequest{
		DisputeIDs:  []uint64{testDispute.ID},
		Status:      model.DisputeStatusCanceled,
		ActorUserID: csUser.ID,
	}

	result, err := svc.BatchUpdateDisputesStatus(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify dispute is canceled
	var dispute model.OrderDispute
	err = db.First(&dispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusCanceled, dispute.Status)
}

func TestDisputeService_BatchCloseDisputes_Refund(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_close_rf_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_close_rf_player")
	csUser := CreateUniqueTestUser(t, db, "batch_close_rf_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_close_rf")

	// Create multiple assigned disputes
	var disputeIDs []uint64
	for i := 0; i < 2; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
		testDispute.Status = model.DisputeStatusAssigned
		db.Save(testDispute)
		disputeIDs = append(disputeIDs, testDispute.ID)
	}

	// Batch close with refund
	req := order.BatchCloseDisputesRequest{
		DisputeIDs:    disputeIDs,
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "服务质量问题，批量退款",
		ActorUserID:   csUser.ID,
	}

	result, err := svc.BatchCloseDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all disputes are resolved and orders refunded
	for _, disputeID := range disputeIDs {
		var dispute model.OrderDispute
		err = db.First(&dispute, disputeID).Error
		require.NoError(t, err)
		assert.Equal(t, model.DisputeStatusResolved, dispute.Status)
		assert.Equal(t, model.ResolutionRefund, dispute.Resolution)
		assert.NotNil(t, dispute.ResolvedAt)
		assert.Equal(t, csUser.ID, *dispute.ResolvedBy)

		// Verify order status
		var order model.Order
		err = db.First(&order, dispute.OrderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusRefunded, order.Status)
		assert.Equal(t, int64(10000), order.RefundAmountCents)
	}
}

func TestDisputeService_BatchCloseDisputes_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_close_rj_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_close_rj_player")
	csUser := CreateUniqueTestUser(t, db, "batch_close_rj_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_close_rj")

	// Create assigned dispute
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusAssigned
	db.Save(testDispute)

	// Batch close with rejection
	req := order.BatchCloseDisputesRequest{
		DisputeIDs:    []uint64{testDispute.ID},
		Resolution:    model.ResolutionReject,
		ResolveRemark: "争议理由不成立，批量驳回",
		ActorUserID:   csUser.ID,
	}

	result, err := svc.BatchCloseDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify dispute is rejected and order is restored
	var dispute model.OrderDispute
	err = db.First(&dispute, testDispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusRejected, dispute.Status)
	assert.Equal(t, model.ResolutionReject, dispute.Resolution)

	var order model.Order
	err = db.First(&order, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)
	assert.False(t, order.HasDispute)
}

func TestDisputeService_BatchCloseDisputes_AlreadyResolved(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_close_ar_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_close_ar_player")
	csUser := CreateUniqueTestUser(t, db, "batch_close_ar_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_close_ar")

	// Create already resolved dispute
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusResolved
	db.Save(testDispute)

	// Try to close already resolved dispute
	req := order.BatchCloseDisputesRequest{
		DisputeIDs:    []uint64{testDispute.ID},
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "Already resolved",
		ActorUserID:   csUser.ID,
	}

	result, err := svc.BatchCloseDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Error, "无法再次处理")
}

func TestDisputeService_BatchCloseDisputes_MixedResolutions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_close_mr_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_close_mr_player")
	csUser := CreateUniqueTestUser(t, db, "batch_close_mr_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_close_mr")

	// Create multiple assigned disputes
	var disputeIDs []uint64
	for i := 0; i < 3; i++ {
		testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
		testDispute.Status = model.DisputeStatusAssigned
		db.Save(testDispute)
		disputeIDs = append(disputeIDs, testDispute.ID)
	}

	// Batch close with partial resolution
	req := order.BatchCloseDisputesRequest{
		DisputeIDs:    disputeIDs,
		Resolution:    model.ResolutionPartial,
		ResolveRemark: "部分退款处理",
		ActorUserID:   csUser.ID,
	}

	result, err := svc.BatchCloseDisputes(ctx, req)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all disputes are resolved with partial resolution
	for _, disputeID := range disputeIDs {
		var dispute model.OrderDispute
		err = db.First(&dispute, disputeID).Error
		require.NoError(t, err)
		assert.Equal(t, model.DisputeStatusResolved, dispute.Status)
		assert.Equal(t, model.ResolutionPartial, dispute.Resolution)
	}
}

func TestDisputeService_BatchAssignDisputes_TooManyDisputes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	csUser := CreateUniqueTestUser(t, db, "batch_assign_tm_cs")

	// Create 101 dispute IDs (exceeds max)
	var disputeIDs []uint64
	for i := 0; i < 101; i++ {
		disputeIDs = append(disputeIDs, uint64(i+1))
	}

	// Try to batch assign too many disputes
	req := order.BatchAssignDisputesRequest{
		DisputeIDs:        disputeIDs,
		AssignedServiceID: csUser.ID,
		ActorUserID:       csUser.ID,
	}

	_, err := svc.BatchAssignDisputes(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "最多支持")
}

func TestDisputeService_BatchCloseDisputes_EmptyRemark(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	orderRepo := implementations.NewOrderRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	disputeRepo := orderrepo.NewDisputeRepository(db)
	operationLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := paymentrepo.NewPaymentRepository(db)

	svc := order.NewDisputeService(disputeRepo, orderRepo, userRepo, operationLogRepo, notificationRepo, paymentRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "batch_close_er_user")
	playerUser := CreateUniqueTestUser(t, db, "batch_close_er_player")
	csUser := CreateUniqueTestUser(t, db, "batch_close_er_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "test_game_batch_close_er")

	// Create assigned dispute
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	testDispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	testDispute.Status = model.DisputeStatusAssigned
	db.Save(testDispute)

	// Try to close with empty remark
	req := order.BatchCloseDisputesRequest{
		DisputeIDs:    []uint64{testDispute.ID},
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "", // Empty remark
		ActorUserID:   csUser.ID,
	}

	_, err := svc.BatchCloseDisputes(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "备注不能为空")
}
