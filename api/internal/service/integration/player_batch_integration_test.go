// Package integration provides integration tests for player batch operations.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	"gorm.io/gorm"
)

// setupPlayerBatchAdminService creates an AdminService instance for player batch tests.
func setupPlayerBatchAdminService(t *testing.T, db *gorm.DB) *admin.AdminService {
	t.Helper()

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	orders := implementations.NewOrderRepository(db)
	payments := payment.NewPaymentRepository(db)
	serviceItems := serviceitem.NewServiceItemRepository(db)

	// For player batch tests, we don't need all repositories
	var roles repository.RoleRepository = nil
	var permissions repository.PermissionRepository = nil
	var menus repository.MenuRepository = nil
	var stats repository.StatsRepository = nil
	var wallets repository.WalletRepository = nil

	return admin.NewAdminService(games, users, playersRepo, orders, payments, roles,
		serviceItems, permissions, menus, stats, wallets, gamecategory.NewGameCategoryRepository(db), nil)
}

// ============================================================================
// BatchUpdatePlayerVerificationStatus Tests
// ============================================================================

// TestAdminService_BatchUpdateVerificationStatus_Approve tests batch approving pending players.
func TestAdminService_BatchUpdateVerificationStatus_Approve(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user for verifiedBy
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create test players with pending status
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "pending_player")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationPending
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Add non-existent player (should fail)
	playerIDs = append(playerIDs, 99999)

	// Batch approve
	remark := "Batch approve for testing"
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, playerIDs, model.VerificationVerified, &adminID, remark)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, resp.TotalCount)
	assert.Equal(t, 3, resp.SuccessCount)
	assert.Equal(t, 1, resp.FailedCount)
	assert.Len(t, resp.SuccessItems, 3)
	assert.Len(t, resp.FailedItems, 1)

	// Verify failed item contains non-existent player
	assert.Equal(t, uint64(99999), resp.FailedItems[0].ID)
	assert.Contains(t, resp.FailedItems[0].Message, "player not found")

	// Verify database updates
	for _, id := range playerIDs[:3] {
		var p model.Player
		err := db.First(&p, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.VerificationVerified, p.VerificationStatus)
		assert.NotNil(t, p.VerifiedAt)
		assert.Equal(t, &adminID, p.VerifiedBy)
		assert.Equal(t, remark, p.VerifyRemark)
		assert.Empty(t, p.RejectReason)
	}
}

// TestAdminService_BatchUpdateVerificationStatus_Reject tests batch rejecting pending players.
func TestAdminService_BatchUpdateVerificationStatus_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create test players with pending status
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "reject_player")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationPending
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Batch reject
	rejectReason := "Documentation incomplete"
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, playerIDs, model.VerificationRejected, &adminID, rejectReason)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 3, resp.TotalCount)
	assert.Equal(t, 3, resp.SuccessCount)
	assert.Equal(t, 0, resp.FailedCount)

	// Verify database updates
	for _, id := range playerIDs {
		var p model.Player
		err := db.First(&p, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.VerificationRejected, p.VerificationStatus)
		assert.Nil(t, p.VerifiedAt)
		assert.Nil(t, p.VerifiedBy)
		assert.Empty(t, p.VerifyRemark)
		assert.Equal(t, rejectReason, p.RejectReason)
	}
}

// TestAdminService_BatchUpdateVerificationStatus_ResetToPending tests resetting rejected players back to pending.
func TestAdminService_BatchUpdateVerificationStatus_ResetToPending(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")

	// Create test players with rejected status
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "reset_player")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationRejected
		player.RejectReason = "Previous rejection"
		player.VerifiedBy = &adminUser.ID
		player.VerifiedAt = &[]time.Time{time.Now()}[0]
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Batch reset to pending
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, playerIDs, model.VerificationPending, nil, "Allow re-submission")
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 3, resp.TotalCount)
	assert.Equal(t, 3, resp.SuccessCount)

	// Verify database updates - all verification fields should be cleared
	for _, id := range playerIDs {
		var p model.Player
		err := db.First(&p, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.VerificationPending, p.VerificationStatus)
		assert.Nil(t, p.VerifiedAt, "VerifiedAt should be cleared")
		assert.Nil(t, p.VerifiedBy, "VerifiedBy should be cleared")
		assert.Empty(t, p.VerifyRemark, "VerifyRemark should be cleared")
		assert.Empty(t, p.RejectReason, "RejectReason should be cleared")
	}
}

// TestAdminService_BatchUpdateVerificationStatus_InvalidTransition tests invalid status transitions.
func TestAdminService_BatchUpdateVerificationStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create player with verified status
	verifiedUser := CreateUniqueTestUser(t, db, "verified_player")
	verifiedPlayer := CreateTestPlayer(t, db, verifiedUser)
	verifiedPlayer.VerificationStatus = model.VerificationVerified
	verifiedPlayer.VerifiedAt = &[]time.Time{time.Now()}[0]
	db.Save(verifiedPlayer)

	// Try to set to pending (valid - revoke)
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{verifiedPlayer.ID}, model.VerificationPending, nil, "Revoking")
	require.NoError(t, err)
	assert.Equal(t, 1, resp.SuccessCount)

	// Verify it's now pending
	var p model.Player
	err = db.First(&p, verifiedPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationPending, p.VerificationStatus)

	// Try invalid transition: pending -> pending (no change allowed)
	resp, err = svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{verifiedPlayer.ID}, model.VerificationPending, nil, "No change")
	require.NoError(t, err)
	assert.Equal(t, 1, resp.FailedCount)
	assert.Contains(t, resp.FailedItems[0].Message, "cannot transition")
}

// TestAdminService_BatchUpdateVerificationStatus_MixedStates tests batch update with mixed player states.
func TestAdminService_BatchUpdateVerificationStatus_MixedStates(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	var playerIDs []uint64

	// Create pending players (can be approved)
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "mixed_pending")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationPending
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Create already verified player (cannot be approved again)
	verifiedUser := CreateUniqueTestUser(t, db, "mixed_verified")
	verifiedPlayer := CreateTestPlayer(t, db, verifiedUser)
	verifiedPlayer.VerificationStatus = model.VerificationVerified
	verifiedPlayer.VerifiedAt = &[]time.Time{time.Now()}[0]
	db.Save(verifiedPlayer)
	playerIDs = append(playerIDs, verifiedPlayer.ID)

	// Create non-existent player
	playerIDs = append(playerIDs, 99999)

	// Batch approve all
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, playerIDs, model.VerificationVerified, &adminID, "Mixed test")
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, resp.TotalCount)
	assert.Equal(t, 2, resp.SuccessCount) // Only pending players can be approved
	assert.Equal(t, 2, resp.FailedCount)  // Verified + non-existent

	// Verify failed items
	var failedIDs []uint64
	for _, item := range resp.FailedItems {
		failedIDs = append(failedIDs, item.ID)
	}
	assert.Contains(t, failedIDs, verifiedPlayer.ID)
	assert.Contains(t, failedIDs, uint64(99999))
}

// TestAdminService_BatchUpdateVerificationStatus_EmptyList tests empty player ID list.
func TestAdminService_BatchUpdateVerificationStatus_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Empty list should succeed with zero counts
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{}, model.VerificationVerified, nil, "")
	require.NoError(t, err)
	assert.Equal(t, 0, resp.TotalCount)
	assert.Equal(t, 0, resp.SuccessCount)
	assert.Equal(t, 0, resp.FailedCount)
}

// TestAdminService_BatchUpdateVerificationStatus_InvalidStatus tests invalid verification status.
func TestAdminService_BatchUpdateVerificationStatus_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create player
	testUser := CreateUniqueTestUser(t, db, "invalid_status_player")
	player := CreateTestPlayer(t, db, testUser)

	// Try invalid status
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, "invalid_status", nil, "")
	require.Error(t, err)
	assert.Nil(t, resp)
}

// ============================================================================
// BatchRevokePlayerCertification Tests
// ============================================================================

// TestAdminService_BatchRevokeCertification_Success tests batch revoking verified players.
func TestAdminService_BatchRevokeCertification_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create verified players
	var playerIDs []uint64
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "revoke_player")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationVerified
		now := time.Now()
		player.VerifiedAt = &now
		player.VerifiedBy = &adminID
		player.VerifyRemark = "Original approval"
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Add non-existent player
	playerIDs = append(playerIDs, 99999)

	// Batch revoke
	revokeReason := "Policy violation - fraudulent documents"
	resp, err := svc.BatchRevokePlayerCertification(ctx, playerIDs, revokeReason, &adminID)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, resp.TotalCount)
	assert.Equal(t, 3, resp.SuccessCount)
	assert.Equal(t, 1, resp.FailedCount)
	assert.Len(t, resp.SuccessItems, 3)
	assert.Len(t, resp.FailedItems, 1)

	// Verify failed item
	assert.Equal(t, uint64(99999), resp.FailedItems[0].ID)
	assert.Contains(t, resp.FailedItems[0].Message, "player not found")

	// Verify database updates
	for _, id := range playerIDs[:3] {
		var p model.Player
		err := db.First(&p, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.VerificationPending, p.VerificationStatus, "Status should be pending after revoke")
		assert.Nil(t, p.VerifiedAt, "VerifiedAt should be cleared")
		assert.Nil(t, p.VerifiedBy, "VerifiedBy should be cleared")
		assert.Empty(t, p.VerifyRemark, "VerifyRemark should be cleared")
		assert.Contains(t, p.RejectReason, "认证已撤销", "RejectReason should contain revocation message")
		assert.Contains(t, p.RejectReason, revokeReason, "RejectReason should contain reason")
	}
}

// TestAdminService_BatchRevokeCertification_NotVerified tests revoking non-verified players.
func TestAdminService_BatchRevokeCertification_NotVerified(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	var playerIDs []uint64

	// Create pending player
	pendingUser := CreateUniqueTestUser(t, db, "pending_revoke")
	pendingPlayer := CreateTestPlayer(t, db, pendingUser)
	pendingPlayer.VerificationStatus = model.VerificationPending
	db.Save(pendingPlayer)
	playerIDs = append(playerIDs, pendingPlayer.ID)

	// Create rejected player
	rejectedUser := CreateUniqueTestUser(t, db, "rejected_revoke")
	rejectedPlayer := CreateTestPlayer(t, db, rejectedUser)
	rejectedPlayer.VerificationStatus = model.VerificationRejected
	rejectedPlayer.RejectReason = "Original rejection"
	db.Save(rejectedPlayer)
	playerIDs = append(playerIDs, rejectedPlayer.ID)

	// Try to revoke (should fail - only verified players can be revoked)
	resp, err := svc.BatchRevokePlayerCertification(ctx, playerIDs, "Test revoke", &adminID)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 2, resp.TotalCount)
	assert.Equal(t, 0, resp.SuccessCount)
	assert.Equal(t, 2, resp.FailedCount)

	// Verify error messages
	for _, item := range resp.FailedItems {
		assert.Contains(t, item.Message, "cannot revoke player with status")
	}

	// Verify database unchanged
	var pendingP, rejectedP model.Player
	err = db.First(&pendingP, pendingPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationPending, pendingP.VerificationStatus)

	err = db.First(&rejectedP, rejectedPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationRejected, rejectedP.VerificationStatus)
}

// TestAdminService_BatchRevokeCertification_EmptyList tests empty player ID list.
func TestAdminService_BatchRevokeCertification_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Empty list should succeed with zero counts
	resp, err := svc.BatchRevokePlayerCertification(ctx, []uint64{}, "Test reason", nil)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.TotalCount)
	assert.Equal(t, 0, resp.SuccessCount)
	assert.Equal(t, 0, resp.FailedCount)
}

// TestAdminService_BatchRevokeCertification_MixedStates tests revoke with mixed player states.
func TestAdminService_BatchRevokeCertification_MixedStates(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	var playerIDs []uint64

	// Create verified players (can be revoked)
	for i := 0; i < 2; i++ {
		user := CreateUniqueTestUser(t, db, "mixed_revoke_verified")
		player := CreateTestPlayer(t, db, user)
		player.VerificationStatus = model.VerificationVerified
		now := time.Now()
		player.VerifiedAt = &now
		player.VerifiedBy = &adminID
		db.Save(player)
		playerIDs = append(playerIDs, player.ID)
	}

	// Create pending player (cannot be revoked)
	pendingUser := CreateUniqueTestUser(t, db, "mixed_revoke_pending")
	pendingPlayer := CreateTestPlayer(t, db, pendingUser)
	pendingPlayer.VerificationStatus = model.VerificationPending
	db.Save(pendingPlayer)
	playerIDs = append(playerIDs, pendingPlayer.ID)

	// Add non-existent player
	playerIDs = append(playerIDs, 99999)

	// Batch revoke
	resp, err := svc.BatchRevokePlayerCertification(ctx, playerIDs, "Mixed revoke test", &adminID)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, resp.TotalCount)
	assert.Equal(t, 2, resp.SuccessCount) // Only verified players
	assert.Equal(t, 2, resp.FailedCount)  // Pending + non-existent

	// Verify revoked players
	var revokedPlayers []model.Player
	for _, id := range playerIDs[:2] {
		var p model.Player
		err := db.First(&p, id).Error
		require.NoError(t, err)
		revokedPlayers = append(revokedPlayers, p)
		assert.Equal(t, model.VerificationPending, p.VerificationStatus)
		assert.Nil(t, p.VerifiedAt)
		assert.Nil(t, p.VerifiedBy)
	}

	// Verify pending player unchanged
	var pendingP model.Player
	err = db.First(&pendingP, pendingPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationPending, pendingP.VerificationStatus)
}

// ============================================================================
// Edge Cases and Error Handling
// ============================================================================

// TestAdminService_PlayerBatchOperations_AllNonExistent tests batch operations with all non-existent players.
func TestAdminService_PlayerBatchOperations_AllNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	playerIDs := []uint64{99999, 88888, 77777}

	// Test BatchUpdateVerificationStatus
	resp1, err := svc.BatchUpdatePlayerVerificationStatus(ctx, playerIDs, model.VerificationVerified, nil, "")
	require.NoError(t, err)
	assert.Equal(t, 3, resp1.TotalCount)
	assert.Equal(t, 0, resp1.SuccessCount)
	assert.Equal(t, 3, resp1.FailedCount)
	for _, item := range resp1.FailedItems {
		assert.Contains(t, item.Message, "player not found")
	}

	// Test BatchRevokeCertification
	resp2, err := svc.BatchRevokePlayerCertification(ctx, playerIDs, "Test", nil)
	require.NoError(t, err)
	assert.Equal(t, 3, resp2.TotalCount)
	assert.Equal(t, 0, resp2.SuccessCount)
	assert.Equal(t, 3, resp2.FailedCount)
	for _, item := range resp2.FailedItems {
		assert.Contains(t, item.Message, "player not found")
	}
}

// TestAdminService_PlayerBatchOperations_VerifyThenRevoke tests verifying then immediately revoking.
func TestAdminService_PlayerBatchOperations_VerifyThenRevoke(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create pending player
	user := CreateUniqueTestUser(t, db, "verify_revoke_player")
	player := CreateTestPlayer(t, db, user)
	player.VerificationStatus = model.VerificationPending
	db.Save(player)

	// Verify the player
	resp1, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationVerified, &adminID, "Approve for test")
	require.NoError(t, err)
	assert.Equal(t, 1, resp1.SuccessCount)

	// Verify the player is now verified
	var p model.Player
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationVerified, p.VerificationStatus)
	assert.NotNil(t, p.VerifiedAt)

	// Now revoke the certification
	resp2, err := svc.BatchRevokePlayerCertification(ctx, []uint64{player.ID}, "Immediate revoke for test", &adminID)
	require.NoError(t, err)
	assert.Equal(t, 1, resp2.SuccessCount)

	// Verify the player is now pending and verification fields are cleared
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationPending, p.VerificationStatus)
	assert.Nil(t, p.VerifiedAt)
	assert.Nil(t, p.VerifiedBy)
	assert.Empty(t, p.VerifyRemark)
	assert.Contains(t, p.RejectReason, "认证已撤销")
}

// TestAdminService_PlayerBatchOperations_RejectThenApprove tests rejecting then approving a player.
func TestAdminService_PlayerBatchOperations_RejectThenApprove(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create pending player
	user := CreateUniqueTestUser(t, db, "reject_approve_player")
	player := CreateTestPlayer(t, db, user)
	player.VerificationStatus = model.VerificationPending
	db.Save(player)

	// Reject the player
	resp1, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationRejected, &adminID, "Incomplete documentation")
	require.NoError(t, err)
	assert.Equal(t, 1, resp1.SuccessCount)

	// Verify the player is rejected
	var p model.Player
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationRejected, p.VerificationStatus)
	assert.Equal(t, "Incomplete documentation", p.RejectReason)

	// Reset to pending
	resp2, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationPending, nil, "Documentation resubmitted")
	require.NoError(t, err)
	assert.Equal(t, 1, resp2.SuccessCount)

	// Verify the player is pending and rejection reason is cleared
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationPending, p.VerificationStatus)
	assert.Empty(t, p.RejectReason)

	// Now approve
	resp3, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationVerified, &adminID, "Documentation verified")
	require.NoError(t, err)
	assert.Equal(t, 1, resp3.SuccessCount)

	// Verify the player is verified
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationVerified, p.VerificationStatus)
	assert.NotNil(t, p.VerifiedAt)
	assert.Equal(t, &adminID, p.VerifiedBy)
	assert.Equal(t, "Documentation verified", p.VerifyRemark)
}

// TestAdminService_PlayerBatchOperations_LongRemark tests handling of long remark strings.
func TestAdminService_PlayerBatchOperations_LongRemark(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create admin user
	adminUser := CreateUniqueTestUser(t, db, "admin")
	var adminID uint64 = adminUser.ID

	// Create player
	user := CreateUniqueTestUser(t, db, "long_remark_player")
	player := CreateTestPlayer(t, db, user)
	player.VerificationStatus = model.VerificationPending
	db.Save(player)

	// Create a very long remark (500+ chars)
	longRemark := "This is a very long rejection reason that contains a lot of details. "
	for i := 0; i < 10; i++ {
		longRemark += " Additional information about why this player's verification was rejected. "
	}

	// Should handle long remark
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationRejected, &adminID, longRemark)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.SuccessCount)

	// Verify the remark was stored (may be truncated by DB column size)
	var p model.Player
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationRejected, p.VerificationStatus)
	// The remark may be truncated to 500 chars (DB column size)
	assert.NotEmpty(t, p.RejectReason)
}

// TestAdminService_PlayerBatchOperations_NilVerifiedBy tests operations with nil verifiedBy.
func TestAdminService_PlayerBatchOperations_NilVerifiedBy(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := setupPlayerBatchAdminService(t, db)
	ctx := context.Background()

	// Create player
	user := CreateUniqueTestUser(t, db, "nil_admin_player")
	player := CreateTestPlayer(t, db, user)
	player.VerificationStatus = model.VerificationPending
	db.Save(player)

	// Approve without admin (nil verifiedBy)
	resp, err := svc.BatchUpdatePlayerVerificationStatus(ctx, []uint64{player.ID}, model.VerificationVerified, nil, "System approval")
	require.NoError(t, err)
	assert.Equal(t, 1, resp.SuccessCount)

	// Verify the player is approved but VerifiedBy is nil
	var p model.Player
	err = db.First(&p, player.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationVerified, p.VerificationStatus)
	assert.NotNil(t, p.VerifiedAt)
	assert.Nil(t, p.VerifiedBy, "VerifiedBy should be nil when no admin specified")
	assert.Equal(t, "System approval", p.VerifyRemark)
}
