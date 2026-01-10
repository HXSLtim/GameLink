// Package integration provides integration tests for withdraw batch operations.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/settlementcompany"
	withdrawrepo "gamelink/internal/repository/withdraw"
	withdrawservice "gamelink/internal/service/withdraw"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// BatchApprove Tests
// ============================================================================

// TestWithdrawService_BatchApprove_MultiplePending_Batch tests batch approving multiple pending withdrawals
func TestWithdrawService_BatchApprove_MultiplePending_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_batch_approve_multi")
	playerUser := CreateUniqueTestUser(t, db, "player_batch_approve_multi")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create multiple pending withdrawals
	var withdrawIDs []uint64
	for i := 0; i < 5; i++ {
		withdraw := CreateTestWithdraw(t, db, testPlayer, int64(10000+i*1000), model.WithdrawStatusPending)
		withdrawIDs = append(withdrawIDs, withdraw.ID)
	}

	// Batch approve
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Remark:      "Batch approval for test",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 5)
	assert.Empty(t, result.FailedItems)

	// Verify database state for each withdrawal
	for i, id := range withdrawIDs {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.WithdrawStatusApproved, withdraw.Status)
		assert.Equal(t, adminUser.ID, *withdraw.ProcessedBy)
		assert.NotNil(t, withdraw.ProcessedAt)
		assert.Equal(t, "Batch approval for test", withdraw.AdminRemark)
		assert.Equal(t, int64(10000+i*1000), withdraw.AmountCents)
	}
}

// TestWithdrawService_BatchApprove_PartialStatusAllowed_Batch tests batch approve with mixed valid/invalid statuses
func TestWithdrawService_BatchApprove_PartialStatusAllowed_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_partial_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_partial_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdrawals with different statuses
	pending1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	pending2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	approved := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	rejected := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusRejected)
	completed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusCompleted)
	failed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusFailed)

	withdrawIDs := []uint64{pending1.ID, pending2.ID, approved.ID, rejected.ID, completed.ID, failed.ID}

	// Batch approve
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Remark:      "Partial status test",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount) // Only pending withdrawals
	assert.Equal(t, 4, result.FailedCount)  // approved, rejected, completed, failed
	assert.Len(t, result.SuccessIDs, 2)
	assert.Len(t, result.FailedItems, 4)

	// Verify SuccessIDs contains only pending withdrawals
	successIDMap := make(map[uint64]bool)
	for _, id := range result.SuccessIDs {
		successIDMap[id] = true
	}
	assert.True(t, successIDMap[pending1.ID])
	assert.True(t, successIDMap[pending2.ID])
	assert.False(t, successIDMap[approved.ID])
	assert.False(t, successIDMap[rejected.ID])

	// Verify FailedItems contains correct IDs with status error messages
	failedIDMap := make(map[uint64]string)
	for _, item := range result.FailedItems {
		failedIDMap[item.ID] = item.Message
	}
	assert.Contains(t, failedIDMap, approved.ID)
	assert.Contains(t, failedIDMap, rejected.ID)
	assert.Contains(t, failedIDMap, completed.ID)
	assert.Contains(t, failedIDMap, failed.ID)

	// Verify message contains status information
	for _, msg := range failedIDMap {
		assert.Contains(t, msg, "cannot approve withdrawal with status")
	}

	// Verify database state - pending should be approved
	withdraw1, _ := withdrawRepo.Get(ctx, pending1.ID)
	assert.Equal(t, model.WithdrawStatusApproved, withdraw1.Status)

	withdraw2, _ := withdrawRepo.Get(ctx, pending2.ID)
	assert.Equal(t, model.WithdrawStatusApproved, withdraw2.Status)

	// Other statuses should remain unchanged
	approvedWithdraw, _ := withdrawRepo.Get(ctx, approved.ID)
	assert.Equal(t, model.WithdrawStatusApproved, approvedWithdraw.Status)

	rejectedWithdraw, _ := withdrawRepo.Get(ctx, rejected.ID)
	assert.Equal(t, model.WithdrawStatusRejected, rejectedWithdraw.Status)
}

// TestWithdrawService_BatchApprove_WithNonExistentIDs_Batch tests batch approve with some non-existent IDs
func TestWithdrawService_BatchApprove_WithNonExistentIDs_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_nonexistent_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_nonexistent_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create some real withdrawals
	pending1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	pending2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

	// Add non-existent IDs
	nonExistentID1 := uint64(999991)
	nonExistentID2 := uint64(999992)

	withdrawIDs := []uint64{pending1.ID, pending2.ID, nonExistentID1, nonExistentID2}

	// Batch approve
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Remark:      "Non-existent ID test",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount) // Only real ones
	assert.Equal(t, 2, result.FailedCount)  // Non-existent ones
	assert.Len(t, result.SuccessIDs, 2)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items contain non-existent IDs
	failedIDMap := make(map[uint64]bool)
	for _, item := range result.FailedItems {
		failedIDMap[item.ID] = true
		assert.Contains(t, item.Message, "not found")
	}
	assert.True(t, failedIDMap[nonExistentID1])
	assert.True(t, failedIDMap[nonExistentID2])
}

// TestWithdrawService_BatchApprove_EmptyList_Batch tests batch approve with empty ID list
func TestWithdrawService_BatchApprove_EmptyList_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Try to approve empty list
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: []uint64{},
		ProcessedBy: 1,
		Remark:      "Empty list test",
	}

	_, err := svc.BatchApprove(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestWithdrawService_BatchApprove_ExceedsLimit_Batch tests batch approve exceeds maximum limit
func TestWithdrawService_BatchApprove_ExceedsLimit_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create a list with 101 IDs (exceeds 100 limit)
	withdrawIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		withdrawIDs[i] = uint64(i + 1)
	}

	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: 1,
		Remark:      "Exceeds limit test",
	}

	_, err := svc.BatchApprove(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 100 withdrawals")
}

// TestWithdrawService_BatchApprove_VerifyTimestamps_Batch tests that ProcessedAt is set correctly
func TestWithdrawService_BatchApprove_VerifyTimestamps_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_timestamps_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_timestamps_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create pending withdrawals
	withdraw1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	withdraw2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

	// Record time before operation
	beforeTime := time.Now()

	// Batch approve
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: []uint64{withdraw1.ID, withdraw2.ID},
		ProcessedBy: adminUser.ID,
		Remark:      "Timestamp test",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Record time after operation
	afterTime := time.Now()

	// Verify ProcessedAt timestamps are set correctly
	for _, id := range []uint64{withdraw1.ID, withdraw2.ID} {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, withdraw.ProcessedAt)
		assert.True(t, withdraw.ProcessedAt.After(beforeTime) || withdraw.ProcessedAt.Equal(beforeTime))
		assert.True(t, withdraw.ProcessedAt.Before(afterTime) || withdraw.ProcessedAt.Equal(afterTime))
	}
}

// ============================================================================
// BatchReject Tests
// ============================================================================

// TestWithdrawService_BatchReject_MultiplePending_Batch tests batch rejecting multiple pending withdrawals
func TestWithdrawService_BatchReject_MultiplePending_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_batch_reject_multi")
	playerUser := CreateUniqueTestUser(t, db, "player_batch_reject_multi")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create multiple pending withdrawals
	var withdrawIDs []uint64
	for i := 0; i < 5; i++ {
		withdraw := CreateTestWithdraw(t, db, testPlayer, int64(10000+i*1000), model.WithdrawStatusPending)
		withdrawIDs = append(withdrawIDs, withdraw.ID)
	}

	// Batch reject
	req := &withdrawservice.BatchRejectRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Reason:      "Invalid bank account information",
	}

	result, err := svc.BatchReject(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 5)
	assert.Empty(t, result.FailedItems)

	// Verify database state for each withdrawal
	for i, id := range withdrawIDs {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.WithdrawStatusRejected, withdraw.Status)
		assert.Equal(t, adminUser.ID, *withdraw.ProcessedBy)
		assert.NotNil(t, withdraw.ProcessedAt)
		assert.Equal(t, "Invalid bank account information", withdraw.RejectReason)
		assert.Equal(t, int64(10000+i*1000), withdraw.AmountCents)
	}
}

// TestWithdrawService_BatchReject_PartialStatusAllowed_Batch tests batch reject with mixed valid/invalid statuses
func TestWithdrawService_BatchReject_PartialStatusAllowed_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_reject_partial_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_reject_partial_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdrawals with different statuses
	pending1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	pending2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	approved := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	rejected := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusRejected)
	completed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusCompleted)

	withdrawIDs := []uint64{pending1.ID, pending2.ID, approved.ID, rejected.ID, completed.ID}

	// Batch reject
	req := &withdrawservice.BatchRejectRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Reason:      "Account verification failed",
	}

	result, err := svc.BatchReject(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount) // Only pending withdrawals
	assert.Equal(t, 3, result.FailedCount)  // approved, rejected, completed
	assert.Len(t, result.SuccessIDs, 2)
	assert.Len(t, result.FailedItems, 3)

	// Verify SuccessIDs contains only pending withdrawals
	successIDMap := make(map[uint64]bool)
	for _, id := range result.SuccessIDs {
		successIDMap[id] = true
	}
	assert.True(t, successIDMap[pending1.ID])
	assert.True(t, successIDMap[pending2.ID])

	// Verify database state - pending should be rejected
	withdraw1, _ := withdrawRepo.Get(ctx, pending1.ID)
	assert.Equal(t, model.WithdrawStatusRejected, withdraw1.Status)
	assert.Equal(t, "Account verification failed", withdraw1.RejectReason)

	// Other statuses should remain unchanged
	approvedWithdraw, _ := withdrawRepo.Get(ctx, approved.ID)
	assert.Equal(t, model.WithdrawStatusApproved, approvedWithdraw.Status)
	assert.Empty(t, approvedWithdraw.RejectReason)
}

// TestWithdrawService_BatchReject_WithNonExistentIDs_Batch tests batch reject with some non-existent IDs
func TestWithdrawService_BatchReject_WithNonExistentIDs_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_reject_nonexist_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_reject_nonexist_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create some real withdrawals
	pending1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

	// Add non-existent IDs
	nonExistentID := uint64(888881)

	withdrawIDs := []uint64{pending1.ID, nonExistentID}

	// Batch reject
	req := &withdrawservice.BatchRejectRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Reason:      "Test with non-existent ID",
	}

	result, err := svc.BatchReject(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 1)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	assert.Equal(t, nonExistentID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "not found")
}

// TestWithdrawService_BatchReject_EmptyList_Batch tests batch reject with empty ID list
func TestWithdrawService_BatchReject_EmptyList_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Try to reject empty list
	req := &withdrawservice.BatchRejectRequest{
		WithdrawIDs: []uint64{},
		ProcessedBy: 1,
		Reason:      "Empty list test",
	}

	_, err := svc.BatchReject(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestWithdrawService_BatchReject_ExceedsLimit_Batch tests batch reject exceeds maximum limit
func TestWithdrawService_BatchReject_ExceedsLimit_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create a list with 101 IDs (exceeds 100 limit)
	withdrawIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		withdrawIDs[i] = uint64(i + 1)
	}

	req := &withdrawservice.BatchRejectRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: 1,
		Reason:      "Exceeds limit test",
	}

	_, err := svc.BatchReject(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 100 withdrawals")
}

// TestWithdrawService_BatchReject_DifferentReasons_Batch tests batch reject with various rejection reasons
func TestWithdrawService_BatchReject_DifferentReasons_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_reasons_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_reasons_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create pending withdrawals
	withdraw1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	withdraw2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	withdraw3 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

	reasons := []string{
		"Invalid bank account",
		"Account name mismatch",
		"Bank account frozen",
	}

	// Test each reason separately
	for i, withdrawID := range []uint64{withdraw1.ID, withdraw2.ID, withdraw3.ID} {
		req := &withdrawservice.BatchRejectRequest{
			WithdrawIDs: []uint64{withdrawID},
			ProcessedBy: adminUser.ID,
			Reason:      reasons[i],
		}

		result, err := svc.BatchReject(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount)

		withdraw, _ := withdrawRepo.Get(ctx, withdrawID)
		assert.Equal(t, reasons[i], withdraw.RejectReason)
	}
}

// ============================================================================
// BatchComplete Tests
// ============================================================================

// TestWithdrawService_BatchComplete_MultipleApproved_Batch tests batch completing multiple approved withdrawals
func TestWithdrawService_BatchComplete_MultipleApproved_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_batch_complete_multi")
	playerUser := CreateUniqueTestUser(t, db, "player_batch_complete_multi")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create multiple approved withdrawals
	var withdrawIDs []uint64
	for i := 0; i < 5; i++ {
		withdraw := CreateTestWithdraw(t, db, testPlayer, int64(10000+i*1000), model.WithdrawStatusApproved)
		withdrawIDs = append(withdrawIDs, withdraw.ID)
	}

	// Batch complete
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: withdrawIDs,
	}

	result, err := svc.BatchComplete(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 5)
	assert.Empty(t, result.FailedItems)

	// Verify database state for each withdrawal
	for _, id := range withdrawIDs {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.WithdrawStatusCompleted, withdraw.Status)
		assert.NotNil(t, withdraw.CompletedAt)
		// ProcessedBy should be set to adminUser.ID if not already set
		assert.Equal(t, adminUser.ID, *withdraw.ProcessedBy)
	}
}

// TestWithdrawService_BatchComplete_PartialStatusAllowed_Batch tests batch complete with mixed valid/invalid statuses
func TestWithdrawService_BatchComplete_PartialStatusAllowed_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_complete_partial_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_complete_partial_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdrawals with different statuses
	approved1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	approved2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	pending := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	rejected := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusRejected)
	completed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusCompleted)

	withdrawIDs := []uint64{approved1.ID, approved2.ID, pending.ID, rejected.ID, completed.ID}

	// Batch complete
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: withdrawIDs,
	}

	result, err := svc.BatchComplete(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount) // Only approved withdrawals
	assert.Equal(t, 3, result.FailedCount)  // pending, rejected, completed
	assert.Len(t, result.SuccessIDs, 2)
	assert.Len(t, result.FailedItems, 3)

	// Verify SuccessIDs contains only approved withdrawals
	successIDMap := make(map[uint64]bool)
	for _, id := range result.SuccessIDs {
		successIDMap[id] = true
	}
	assert.True(t, successIDMap[approved1.ID])
	assert.True(t, successIDMap[approved2.ID])

	// Verify database state - approved should be completed
	withdraw1, _ := withdrawRepo.Get(ctx, approved1.ID)
	assert.Equal(t, model.WithdrawStatusCompleted, withdraw1.Status)

	withdraw2, _ := withdrawRepo.Get(ctx, approved2.ID)
	assert.Equal(t, model.WithdrawStatusCompleted, withdraw2.Status)

	// Other statuses should remain unchanged
	pendingWithdraw, _ := withdrawRepo.Get(ctx, pending.ID)
	assert.Equal(t, model.WithdrawStatusPending, pendingWithdraw.Status)

	rejectedWithdraw, _ := withdrawRepo.Get(ctx, rejected.ID)
	assert.Equal(t, model.WithdrawStatusRejected, rejectedWithdraw.Status)
}

// TestWithdrawService_BatchComplete_WithNonExistentIDs_Batch tests batch complete with some non-existent IDs
func TestWithdrawService_BatchComplete_WithNonExistentIDs_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_complete_nonexist_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_complete_nonexist_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create some real withdrawals
	approved1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	approved2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)

	// Add non-existent IDs
	nonExistentID1 := uint64(777771)
	nonExistentID2 := uint64(777772)

	withdrawIDs := []uint64{approved1.ID, approved2.ID, nonExistentID1, nonExistentID2}

	// Batch complete
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: withdrawIDs,
	}

	result, err := svc.BatchComplete(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount) // Only real ones
	assert.Equal(t, 2, result.FailedCount)  // Non-existent ones
	assert.Len(t, result.SuccessIDs, 2)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items contain non-existent IDs
	failedIDMap := make(map[uint64]bool)
	for _, item := range result.FailedItems {
		failedIDMap[item.ID] = true
		assert.Contains(t, item.Message, "not found")
	}
	assert.True(t, failedIDMap[nonExistentID1])
	assert.True(t, failedIDMap[nonExistentID2])
}

// TestWithdrawService_BatchComplete_EmptyList_Batch tests batch complete with empty ID list
func TestWithdrawService_BatchComplete_EmptyList_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Try to complete empty list
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: []uint64{},
	}

	_, err := svc.BatchComplete(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestWithdrawService_BatchComplete_ExceedsLimit_Batch tests batch complete exceeds maximum limit
func TestWithdrawService_BatchComplete_ExceedsLimit_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create a list with 101 IDs (exceeds 100 limit)
	withdrawIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		withdrawIDs[i] = uint64(i + 1)
	}

	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: withdrawIDs,
	}

	_, err := svc.BatchComplete(ctx, req, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 100 withdrawals")
}

// TestWithdrawService_BatchComplete_VerifyTimestamps_Batch tests that CompletedAt is set correctly
func TestWithdrawService_BatchComplete_VerifyTimestamps_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_complete_time_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_complete_time_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create approved withdrawals
	approved1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	approved2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)

	// Record time before operation
	beforeTime := time.Now()

	// Batch complete
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: []uint64{approved1.ID, approved2.ID},
	}

	result, err := svc.BatchComplete(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Record time after operation
	afterTime := time.Now()

	// Verify CompletedAt timestamps are set correctly
	for _, id := range []uint64{approved1.ID, approved2.ID} {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, withdraw.CompletedAt)
		assert.True(t, withdraw.CompletedAt.After(beforeTime) || withdraw.CompletedAt.Equal(beforeTime))
		assert.True(t, withdraw.CompletedAt.Before(afterTime) || withdraw.CompletedAt.Equal(afterTime))
		assert.Equal(t, model.WithdrawStatusCompleted, withdraw.Status)
	}
}

// TestWithdrawService_BatchComplete_PreservesProcessedBy_Batch tests that ProcessedBy is preserved if already set
func TestWithdrawService_BatchComplete_PreservesProcessedBy_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser1 := CreateUniqueTestUser(t, db, "admin_complete_preserve1_batch")
	adminUser2 := CreateUniqueTestUser(t, db, "admin_complete_preserve2_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_complete_preserve_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create approved withdrawals with ProcessedBy already set
	now := time.Now()
	approved1 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	approved1.ProcessedBy = &adminUser1.ID
	approved1.ProcessedAt = &now
	withdrawRepo.Update(ctx, approved1)

	approved2 := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	// Leave ProcessedBy nil for this one

	// Batch complete with different admin
	req := &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: []uint64{approved1.ID, approved2.ID},
	}

	result, err := svc.BatchComplete(ctx, req, adminUser2.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Verify ProcessedBy is preserved for approved1, set for approved2
	withdraw1, _ := withdrawRepo.Get(ctx, approved1.ID)
	assert.Equal(t, adminUser1.ID, *withdraw1.ProcessedBy) // Should preserve original

	withdraw2, _ := withdrawRepo.Get(ctx, approved2.ID)
	assert.Equal(t, adminUser2.ID, *withdraw2.ProcessedBy) // Should set to current admin
}

// ============================================================================
// Edge Cases and Complex Scenarios
// ============================================================================

// TestWithdrawService_BatchOperations_AllStatuses_Batch tests all withdrawal statuses in batch operations
func TestWithdrawService_BatchOperations_AllStatuses_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_all_statuses_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_all_statuses_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create withdrawals with all possible statuses
	pending := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)
	approved := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)
	rejected := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusRejected)
	completed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusCompleted)
	failed := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusFailed)

	// Test BatchApprove - only pending should succeed
	t.Run("BatchApprove_AllStatuses", func(t *testing.T) {
		req := &withdrawservice.BatchApproveRequest{
			WithdrawIDs: []uint64{pending.ID, approved.ID, rejected.ID, completed.ID, failed.ID},
			ProcessedBy: adminUser.ID,
			Remark:      "Test all statuses",
		}

		result, err := svc.BatchApprove(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount) // Only pending
		assert.Equal(t, 4, result.FailedCount)
	})

	// Test BatchReject - only pending should succeed
	t.Run("BatchReject_AllStatuses", func(t *testing.T) {
		// Reset pending status for next test
		pending.Status = model.WithdrawStatusPending
		withdrawRepo.Update(ctx, pending)

		req := &withdrawservice.BatchRejectRequest{
			WithdrawIDs: []uint64{pending.ID, approved.ID, rejected.ID, completed.ID, failed.ID},
			ProcessedBy: adminUser.ID,
			Reason:      "Test all statuses",
		}

		result, err := svc.BatchReject(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount) // Only pending
		assert.Equal(t, 4, result.FailedCount)
	})

	// Test BatchComplete - only approved should succeed
	t.Run("BatchComplete_AllStatuses", func(t *testing.T) {
		req := &withdrawservice.BatchCompleteRequest{
			WithdrawIDs: []uint64{pending.ID, approved.ID, rejected.ID, completed.ID, failed.ID},
		}

		result, err := svc.BatchComplete(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount) // Only approved
		assert.Equal(t, 4, result.FailedCount)
	})
}

// TestWithdrawService_BatchOperations_MultiplePlayers_Batch tests batch operations across multiple players
func TestWithdrawService_BatchOperations_MultiplePlayers_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_multi_player_batch")

	// Create multiple players with pending withdrawals
	var allWithdrawIDs []uint64
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, "player_multi"+string(rune('0'+i))+"_batch")
		testPlayer := CreateTestPlayer(t, db, playerUser)
		for j := 0; j < 2; j++ {
			withdraw := CreateTestWithdraw(t, db, testPlayer, int64(10000+i*1000+j*500), model.WithdrawStatusPending)
			allWithdrawIDs = append(allWithdrawIDs, withdraw.ID)
		}
	}

	// Batch approve all
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: allWithdrawIDs,
		ProcessedBy: adminUser.ID,
		Remark:      "Multi-player batch approval",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 6, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all withdrawals are approved
	for _, id := range allWithdrawIDs {
		withdraw, err := withdrawRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.WithdrawStatusApproved, withdraw.Status)
	}
}

// TestWithdrawService_BatchOperations_SingleWithdrawal_Batch tests batch operations with single withdrawal
func TestWithdrawService_BatchOperations_SingleWithdrawal_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_single_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_single_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Test single withdrawal for each operation
	t.Run("BatchApprove_Single", func(t *testing.T) {
		withdraw := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

		req := &withdrawservice.BatchApproveRequest{
			WithdrawIDs: []uint64{withdraw.ID},
			ProcessedBy: adminUser.ID,
			Remark:      "Single approval",
		}

		result, err := svc.BatchApprove(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount)
		assert.Equal(t, 0, result.FailedCount)
		assert.Len(t, result.SuccessIDs, 1)
		assert.Equal(t, withdraw.ID, result.SuccessIDs[0])
	})

	t.Run("BatchReject_Single", func(t *testing.T) {
		withdraw := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusPending)

		req := &withdrawservice.BatchRejectRequest{
			WithdrawIDs: []uint64{withdraw.ID},
			ProcessedBy: adminUser.ID,
			Reason:      "Single rejection",
		}

		result, err := svc.BatchReject(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount)
		assert.Equal(t, 0, result.FailedCount)
		assert.Len(t, result.SuccessIDs, 1)
		assert.Equal(t, withdraw.ID, result.SuccessIDs[0])
	})

	t.Run("BatchComplete_Single", func(t *testing.T) {
		withdraw := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusApproved)

		req := &withdrawservice.BatchCompleteRequest{
			WithdrawIDs: []uint64{withdraw.ID},
		}

		result, err := svc.BatchComplete(ctx, req, adminUser.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount)
		assert.Equal(t, 0, result.FailedCount)
		assert.Len(t, result.SuccessIDs, 1)
		assert.Equal(t, withdraw.ID, result.SuccessIDs[0])
	})
}

// TestWithdrawService_BatchOperations_MaximumLimit_Batch tests batch operations at the maximum limit (100)
func TestWithdrawService_BatchOperations_MaximumLimit_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_max_limit_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_max_limit_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create exactly 100 pending withdrawals
	var withdrawIDs []uint64
	for i := 0; i < 100; i++ {
		withdraw := CreateTestWithdraw(t, db, testPlayer, int64(10000+i*100), model.WithdrawStatusPending)
		withdrawIDs = append(withdrawIDs, withdraw.ID)
	}

	// Batch approve all 100
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: withdrawIDs,
		ProcessedBy: adminUser.ID,
		Remark:      "Maximum limit test",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 100, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 100)
}

// TestWithdrawService_BatchOperations_FailedStatus_Batch tests handling of failed status withdrawals
func TestWithdrawService_BatchOperations_FailedStatus_Batch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	withdrawRepo := withdrawrepo.NewWithdrawRepository(db)
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	svc := withdrawservice.NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	// Create test data
	adminUser := CreateUniqueTestUser(t, db, "admin_failed_batch")
	playerUser := CreateUniqueTestUser(t, db, "player_failed_batch")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create failed withdrawal
	failedWithdraw := CreateTestWithdraw(t, db, testPlayer, 10000, model.WithdrawStatusFailed)

	// Try to approve failed withdrawal
	req := &withdrawservice.BatchApproveRequest{
		WithdrawIDs: []uint64{failedWithdraw.ID},
		ProcessedBy: adminUser.ID,
		Remark:      "Test failed status",
	}

	result, err := svc.BatchApprove(ctx, req, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
	assert.Contains(t, result.FailedItems[0].Message, "cannot approve withdrawal with status: failed")

	// Verify status remains failed
	withdraw, _ := withdrawRepo.Get(ctx, failedWithdraw.ID)
	assert.Equal(t, model.WithdrawStatusFailed, withdraw.Status)
}
