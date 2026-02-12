// Package integration provides supplementary integration tests for dispute service.
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

// TestDisputeService_CreateChatSnapshot tests creating chat snapshot when dispute is initiated.
func TestDisputeService_CreateChatSnapshot(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "snapshot_user")
	playerUser := CreateUniqueTestUser(t, db, "snapshot_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "snapshot_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create chat group for order
	chatGroup := CreateTestChatGroup(t, db, "Order Chat", model.ChatGroupTypeOrder, &testOrder.ID)

	// Add some messages
	CreateTestChatMessage(t, db, chatGroup.ID, testUser.ID, "Hello, I have an issue")
	CreateTestChatMessage(t, db, chatGroup.ID, playerUser.ID, "What's the problem?")
	CreateTestChatMessage(t, db, chatGroup.ID, testUser.ID, "Service was not as expected")

	// Create dispute
	dispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)

	// Create chat snapshot
	snapshot := CreateTestChatSnapshot(t, db, dispute.ID, testOrder.ID, chatGroup.ID)

	// Link snapshot to dispute
	dispute.ChatSnapshotID = &snapshot.ID
	err := db.Save(dispute).Error
	require.NoError(t, err)

	// Verify
	var savedDispute model.OrderDispute
	err = db.First(&savedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, savedDispute.ChatSnapshotID)

	var savedSnapshot model.ChatSnapshot
	err = db.First(&savedSnapshot, snapshot.ID).Error
	require.NoError(t, err)
	assert.Equal(t, dispute.ID, savedSnapshot.DisputeID)
	assert.Equal(t, testOrder.ID, savedSnapshot.OrderID)
	assert.NotEmpty(t, savedSnapshot.Messages)
}

// TestDisputeService_DoubleCSAssignment tests double customer service assignment.
func TestDisputeService_DoubleCSAssignment(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "double_cs_user")
	playerUser := CreateUniqueTestUser(t, db, "double_cs_player")
	originalCS := CreateUniqueTestUser(t, db, "original_cs")
	assignedCS := CreateUniqueTestUser(t, db, "assigned_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "double_cs_game")

	// Create completed order with original CS
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute with original CS
	dispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	dispute.OriginalServiceID = &originalCS.ID
	err := db.Save(dispute).Error
	require.NoError(t, err)

	// Assign independent CS
	req := order.AssignDisputeRequest{
		DisputeID:         dispute.ID,
		AssignedServiceID: assignedCS.ID,
		ActorUserID:       assignedCS.ID,
	}

	err = svc.AssignDispute(ctx, req)
	require.NoError(t, err)

	// Verify double CS assignment
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, updatedDispute.OriginalServiceID)
	assert.NotNil(t, updatedDispute.AssignedServiceID)
	assert.Equal(t, originalCS.ID, *updatedDispute.OriginalServiceID)
	assert.Equal(t, assignedCS.ID, *updatedDispute.AssignedServiceID)
	assert.NotEqual(t, *updatedDispute.OriginalServiceID, *updatedDispute.AssignedServiceID)
}

// TestDisputeService_UploadEvidence tests uploading evidence for dispute.
func TestDisputeService_UploadEvidence(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "evidence_user")
	playerUser := CreateUniqueTestUser(t, db, "evidence_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "evidence_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute with evidence
	dispute := &model.OrderDispute{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:       testOrder.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Status:        model.DisputeStatusPending,
		Reason:        "Service quality issue",
		EvidenceText:  "The player did not complete the service as promised",
		EvidenceURLs: model.EvidenceURLArray{
			"https://example.com/evidence1.jpg",
			"https://example.com/evidence2.jpg",
			"https://example.com/evidence3.jpg",
		},
	}
	err := db.Create(dispute).Error
	require.NoError(t, err)

	// Verify evidence
	var savedDispute model.OrderDispute
	err = db.First(&savedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.Len(t, savedDispute.EvidenceURLs, 3)
	assert.NotEmpty(t, savedDispute.EvidenceText)
}

// TestDisputeService_EvidenceLimit tests evidence upload limit (max 5).
func TestDisputeService_EvidenceLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "limit_user")
	playerUser := CreateUniqueTestUser(t, db, "limit_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "limit_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute with max evidence (5)
	dispute := &model.OrderDispute{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:       testOrder.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Status:        model.DisputeStatusPending,
		Reason:        "Multiple issues",
		EvidenceURLs: model.EvidenceURLArray{
			"https://example.com/evidence1.jpg",
			"https://example.com/evidence2.jpg",
			"https://example.com/evidence3.jpg",
			"https://example.com/evidence4.jpg",
			"https://example.com/evidence5.jpg",
		},
	}
	err := db.Create(dispute).Error
	require.NoError(t, err)

	// Verify max 5 evidence
	var savedDispute model.OrderDispute
	err = db.First(&savedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.Len(t, savedDispute.EvidenceURLs, 5)
}

// TestDisputeService_DisputeTemplate tests using dispute template.
func TestDisputeService_DisputeTemplate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create dispute templates
	template1 := CreateTestDisputeTemplate(t, db, "service_quality", "服务质量问题", model.DisputeInitiatorUser)
	template2 := CreateTestDisputeTemplate(t, db, "bad_attitude", "态度问题", model.DisputeInitiatorUser)
	template3 := CreateTestDisputeTemplate(t, db, "user_not_cooperative", "用户不配合", model.DisputeInitiatorPlayer)

	// Verify templates
	var templates []model.DisputeTemplate
	err := db.Where("is_active = ?", true).Find(&templates).Error
	require.NoError(t, err)
	assert.Len(t, templates, 3)

	// Verify user templates
	var userTemplates []model.DisputeTemplate
	err = db.Where("initiator_type = ? AND is_active = ?", model.DisputeInitiatorUser, true).Find(&userTemplates).Error
	require.NoError(t, err)
	assert.Len(t, userTemplates, 2)
	assert.Equal(t, template1.Code, userTemplates[0].Code)
	assert.Equal(t, template2.Code, userTemplates[1].Code)

	// Verify player templates
	var playerTemplates []model.DisputeTemplate
	err = db.Where("initiator_type = ? AND is_active = ?", model.DisputeInitiatorPlayer, true).Find(&playerTemplates).Error
	require.NoError(t, err)
	assert.Len(t, playerTemplates, 1)
	assert.Equal(t, template3.Code, playerTemplates[0].Code)
}

// TestDisputeService_DisputeWithinTimeWindow tests dispute within 7-day window.
func TestDisputeService_DisputeWithinTimeWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "window_user")
	playerUser := CreateUniqueTestUser(t, db, "window_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "window_game")

	// Create completed order within 7 days
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	completedAt := time.Now().Add(-3 * 24 * time.Hour) // 3 days ago
	testOrder.CompletedAt = &completedAt
	db.Save(testOrder)

	// Check if dispute can be initiated
	canDispute := model.CanInitiateDispute(testOrder)
	assert.True(t, canDispute)

	// Create dispute
	dispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	assert.NotZero(t, dispute.ID)
}

// TestDisputeService_DisputeOutsideTimeWindow tests dispute outside 7-day window.
func TestDisputeService_DisputeOutsideTimeWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "outside_window_user")
	playerUser := CreateUniqueTestUser(t, db, "outside_window_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "outside_window_game")

	// Create completed order outside 7 days
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	completedAt := time.Now().Add(-10 * 24 * time.Hour) // 10 days ago
	testOrder.CompletedAt = &completedAt
	db.Save(testOrder)

	// Check if dispute can be initiated
	canDispute := model.CanInitiateDispute(testOrder)
	assert.False(t, canDispute)
}

// TestDisputeService_DisputeDuringService tests dispute during service (in_progress).
func TestDisputeService_DisputeDuringService(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "during_service_user")
	playerUser := CreateUniqueTestUser(t, db, "during_service_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "during_service_game")

	// Create in_progress order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)

	// Check if dispute can be initiated during service
	canDispute := model.CanInitiateDispute(testOrder)
	assert.True(t, canDispute)
}

// TestDisputeService_SLADeadlineCalculation tests SLA deadline calculation.
func TestDisputeService_SLADeadlineCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "sla_calc_user")
	playerUser := CreateUniqueTestUser(t, db, "sla_calc_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "sla_calc_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create dispute with SLA deadline (30 minutes)
	slaDeadline := time.Now().Add(30 * time.Minute)
	dispute := &model.OrderDispute{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:       testOrder.ID,
		InitiatorID:   testUser.ID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Status:        model.DisputeStatusPending,
		Reason:        "Test dispute",
		SLADeadline:   &slaDeadline,
		EvidenceURLs:  model.EvidenceURLArray{},
	}
	err := db.Create(dispute).Error
	require.NoError(t, err)

	// Verify SLA
	var savedDispute model.OrderDispute
	err = db.First(&savedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, savedDispute.SLADeadline)
	assert.False(t, savedDispute.IsOverSLA())
	assert.Greater(t, savedDispute.GetSLARemaining(), int64(0))
}

// TestDisputeService_ResolveWithRefund tests resolving dispute with full refund.
func TestDisputeService_ResolveWithRefund(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "refund_resolve_user")
	playerUser := CreateUniqueTestUser(t, db, "refund_resolve_player")
	csUser := CreateUniqueTestUser(t, db, "refund_resolve_cs")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "refund_resolve_game")

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create and assign dispute
	dispute := CreateTestDispute(t, db, testOrder, testUser.ID, model.DisputeInitiatorUser)
	dispute.Status = model.DisputeStatusAssigned
	dispute.AssignedServiceID = &csUser.ID
	db.Save(dispute)

	// Resolve with refund
	req := order.ResolveDisputeRequest{
		DisputeID:     dispute.ID,
		Resolution:    model.ResolutionRefund,
		ResolveRemark: "Approved full refund",
		ActorUserID:   csUser.ID,
	}

	err := svc.ResolveDispute(ctx, req)
	require.NoError(t, err)

	// Verify
	var updatedDispute model.OrderDispute
	err = db.First(&updatedDispute, dispute.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.DisputeStatusResolved, updatedDispute.Status)
	assert.Equal(t, model.ResolutionRefund, updatedDispute.Resolution)

	var updatedOrder model.Order
	err = db.First(&updatedOrder, testOrder.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
}
