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
