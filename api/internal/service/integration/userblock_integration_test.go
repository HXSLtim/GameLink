package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	userrepository "gamelink/internal/repository/user"
	userblockrepo "gamelink/internal/repository/userblock"
	userblockservice "gamelink/internal/service/userblock"

	"gorm.io/gorm"
)

// ============================================================================
// UserBlock Service Integration Tests
// ============================================================================

// Helper function to create a service instance
func newUserBlockService(db *gorm.DB) *userblockservice.UserBlockService {
	blockRepo := userblockrepo.NewUserBlockRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	return userblockservice.NewUserBlockService(blockRepo, userRepo)
}

// ============================================================================
// Block Operation Tests
// ============================================================================

func TestUserBlockService_Block_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	// Create users
	blocker := CreateUniqueTestUser(t, db, "blocker")
	blocked := CreateUniqueTestUser(t, db, "blocked")

	// Block user
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Rude behavior",
	}

	result, err := svc.Block(ctx, input)
	require.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, blocker.ID, result.BlockerID)
	assert.Equal(t, blocked.ID, result.BlockedID)
	assert.Equal(t, model.BlockStatusActive, result.Status)
	assert.Equal(t, "Rude behavior", result.Reason)
	assert.False(t, result.BlockedAt.IsZero())
}

func TestUserBlockService_Block_PlayerBlocksUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	// Create player and user
	playerUser := CreateUniqueTestUser(t, db, "player_user")
	_ = CreateTestPlayer(t, db, playerUser)
	regularUser := CreateUniqueTestUser(t, db, "regular_user")

	// Player blocks user
	input := userblockservice.BlockInput{
		BlockerID:   playerUser.ID,
		BlockerType: model.BlockUserTypePlayer,
		BlockedID:   regularUser.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Harassment",
	}

	result, err := svc.Block(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.BlockUserTypePlayer, result.BlockerType)
	assert.Equal(t, model.BlockUserTypeUser, result.BlockedType)
}

func TestUserBlockService_Block_UserBlocksPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	// Create user and player
	regularUser := CreateUniqueTestUser(t, db, "regular_user")
	playerUser := CreateUniqueTestUser(t, db, "player_user")
	_ = CreateTestPlayer(t, db, playerUser)

	// User blocks player
	input := userblockservice.BlockInput{
		BlockerID:   regularUser.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   playerUser.ID,
		BlockedType: model.BlockUserTypePlayer,
		Reason:      "Bad service",
	}

	result, err := svc.Block(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.BlockUserTypeUser, result.BlockerType)
	assert.Equal(t, model.BlockUserTypePlayer, result.BlockedType)
}

func TestUserBlockService_Block_CannotBlockSelf(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "self_blocker")

	input := userblockservice.BlockInput{
		BlockerID:   user.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Self block",
	}

	_, err := svc.Block(ctx, input)
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrCannotBlockSelf, err)
}

func TestUserBlockService_Block_UserNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "blocker")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   99999, // Non-existent user
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Test",
	}

	_, err := svc.Block(ctx, input)
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrUserNotFound, err)
}

func TestUserBlockService_Block_AlreadyBlocked(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "blocker")
	blocked := CreateUniqueTestUser(t, db, "blocked")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "First block",
	}

	// First block should succeed
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Second block should fail
	input.Reason = "Second block"
	_, err = svc.Block(ctx, input)
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrAlreadyBlocked, err)
}

func TestUserBlockService_Block_ReblockAfterUnblock(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "blocker")
	blocked := CreateUniqueTestUser(t, db, "blocked")

	// Block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "First block",
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Unblock
	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)

	// Re-block should succeed
	input.Reason = "Re-block"
	_, err = svc.Block(ctx, input)
	require.NoError(t, err)
}

// ============================================================================
// Unblock Operation Tests
// ============================================================================

func TestUserBlockService_Unblock_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "unblock_blocker")
	blocked := CreateUniqueTestUser(t, db, "unblock_blocked")

	// Create block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Will unblock",
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Unblock
	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)

	// Verify status changed
	updated, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BlockStatusCanceled, updated.Status)
	assert.NotNil(t, updated.CanceledAt)
	assert.Nil(t, updated.CanceledBy) // User cancel doesn't track who
}

func TestUserBlockService_Unblock_NotBlocked(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "user1")
	user2 := CreateUniqueTestUser(t, db, "user2")

	err := svc.Unblock(ctx, user1.ID, user2.ID)
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrNotBlocked, err)
}

func TestUserBlockService_Unblock_AlreadyCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "already_blocker")
	blocked := CreateUniqueTestUser(t, db, "already_blocked")

	// Create and cancel block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)

	// Try to unblock again
	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrNotBlocked, err)
}

// ============================================================================
// AdminUnblock Operation Tests
// ============================================================================

func TestUserBlockService_AdminUnblock_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "admin_blocker")
	blocked := CreateUniqueTestUser(t, db, "admin_blocked")
	admin := CreateUniqueTestUser(t, db, "admin")

	// Create block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Will be removed by admin",
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Admin unblock
	err = svc.AdminUnblock(ctx, block.ID, admin.ID, "Investigated and resolved")
	require.NoError(t, err)

	// Verify
	updated, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BlockStatusAdminCanceled, updated.Status)
	assert.NotNil(t, updated.CanceledAt)
	assert.NotNil(t, updated.CanceledBy)
	assert.Equal(t, admin.ID, *updated.CanceledBy)
	assert.Equal(t, "Investigated and resolved", updated.AdminRemark)
}

func TestUserBlockService_AdminUnblock_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	admin := CreateUniqueTestUser(t, db, "admin")

	err := svc.AdminUnblock(ctx, 99999, admin.ID, "Test")
	require.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestUserBlockService_AdminUnblock_AlreadyCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "cancel_blocker")
	blocked := CreateUniqueTestUser(t, db, "cancel_blocked")
	admin := CreateUniqueTestUser(t, db, "admin")

	// Create and cancel block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)

	// Admin tries to unblock already canceled block
	err = svc.AdminUnblock(ctx, block.ID, admin.ID, "Admin retry")
	require.Error(t, err)
	assert.Equal(t, userblockservice.ErrNotBlocked, err)
}

// ============================================================================
// IsBlocked / IsBlockedBy Tests
// ============================================================================

func TestUserBlockService_IsBlocked_Bidirectional(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "bidir1")
	user2 := CreateUniqueTestUser(t, db, "bidir2")

	// Initially not blocked
	blocked, err := svc.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.False(t, blocked)

	// User1 blocks User2
	input := userblockservice.BlockInput{
		BlockerID:   user1.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user2.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err = svc.Block(ctx, input)
	require.NoError(t, err)

	// Check IsBlocked (bidirectional)
	blocked, err = svc.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.True(t, blocked)

	blocked, err = svc.IsBlocked(ctx, user2.ID, user1.ID)
	require.NoError(t, err)
	assert.True(t, blocked) // Should work both ways
}

func TestUserBlockService_IsBlockedBy_Directional(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "dir_blocker")
	blocked := CreateUniqueTestUser(t, db, "dir_blocked")

	// Block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Check direction
	isBlocked, err := svc.IsBlockedBy(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)
	assert.True(t, isBlocked) // blocker blocked blocked

	isBlocked, err = svc.IsBlockedBy(ctx, blocked.ID, blocker.ID)
	require.NoError(t, err)
	assert.False(t, isBlocked) // blocked did NOT block blocker
}

func TestUserBlockService_IsBlocked_NotCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "notcancel1")
	user2 := CreateUniqueTestUser(t, db, "notcancel2")

	// Block and then unblock
	input := userblockservice.BlockInput{
		BlockerID:   user1.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user2.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	err = svc.Unblock(ctx, user1.ID, user2.ID)
	require.NoError(t, err)

	// Should no longer be blocked
	blocked, err := svc.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.False(t, blocked)
}

// ============================================================================
// List Operations Tests
// ============================================================================

func TestUserBlockService_ListByBlocker_ActiveOnly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "list_blocker")

	// Create active blocks
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "list_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Reason:      "Active block",
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// Create canceled block
	canceledBlocked := CreateUniqueTestUser(t, db, "canceled_blocked")
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   canceledBlocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)
	svc.Unblock(ctx, blocker.ID, canceledBlocked.ID)

	// List active only
	blocks, err := svc.ListByBlocker(ctx, blocker.ID, true)
	require.NoError(t, err)
	assert.Len(t, blocks, 3)

	// List all
	blocks, err = svc.ListByBlocker(ctx, blocker.ID, false)
	require.NoError(t, err)
	assert.Len(t, blocks, 4)
}

func TestUserBlockService_ListByBlocked(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blockedUser := CreateUniqueTestUser(t, db, "blocked_user")

	// Multiple users block the same person
	for i := 0; i < 3; i++ {
		blocker := CreateUniqueTestUser(t, db, "blocker")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blockedUser.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// List who blocked this user
	blocks, err := svc.ListByBlocked(ctx, blockedUser.ID, true)
	require.NoError(t, err)
	assert.Len(t, blocks, 3)
}

func TestUserBlockService_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "paged_blocker")

	// Create 5 blocks
	for i := 0; i < 5; i++ {
		blocked := CreateUniqueTestUser(t, db, "paged_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// First page
	blocks, pagination, err := svc.ListPaged(ctx, repository.UserBlockListOptions{
		Page:     1,
		PageSize: 2,
	})
	require.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.GreaterOrEqual(t, pagination.Total, 5)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 2, pagination.PageSize)

	// Second page
	blocks, pagination, err = svc.ListPaged(ctx, repository.UserBlockListOptions{
		Page:     2,
		PageSize: 2,
	})
	require.NoError(t, err)
	assert.Len(t, blocks, 2)
}

func TestUserBlockService_ListPaged_WithFilters(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "filter_blocker")

	// Create active blocks
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "filter_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// Filter by blockerID
	filterBlockerID := blocker.ID
	blocks, _, err := svc.ListPaged(ctx, repository.UserBlockListOptions{
		Page:        1,
		PageSize:    10,
		BlockerID:   &filterBlockerID,
		BlockerType: nil,
		Status:      nil,
	})
	require.NoError(t, err)
	assert.Len(t, blocks, 3)

	// Filter by status
	activeStatus := model.BlockStatusActive
	blocks, _, err = svc.ListPaged(ctx, repository.UserBlockListOptions{
		Page:     1,
		PageSize: 10,
		Status:   &activeStatus,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(blocks), 3)
}

// ============================================================================
// GetBlockedUserIDs / GetBlockerUserIDs Tests
// ============================================================================

func TestUserBlockService_GetBlockedUserIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "ids_blocker")
	var expectedIDs []uint64

	// Block multiple users
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "ids_blocked")
		expectedIDs = append(expectedIDs, blocked.ID)
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// Get blocked user IDs
	ids, err := svc.GetBlockedUserIDs(ctx, blocker.ID)
	require.NoError(t, err)
	assert.Len(t, ids, 3)
	for _, expectedID := range expectedIDs {
		assert.Contains(t, ids, expectedID)
	}
}

func TestUserBlockService_GetBlockerUserIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blockedUser := CreateUniqueTestUser(t, db, "getblocker_blocked")
	var expectedIDs []uint64

	// Multiple users block the same person
	for i := 0; i < 3; i++ {
		blocker := CreateUniqueTestUser(t, db, "getblocker_blocker")
		expectedIDs = append(expectedIDs, blocker.ID)
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blockedUser.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// Get blocker user IDs
	ids, err := svc.GetBlockerUserIDs(ctx, blockedUser.ID)
	require.NoError(t, err)
	assert.Len(t, ids, 3)
	for _, expectedID := range expectedIDs {
		assert.Contains(t, ids, expectedID)
	}
}

func TestUserBlockService_GetAllBlockRelatedUserIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "related_user")
	blockedByUser := CreateUniqueTestUser(t, db, "related_blocked")
	blockerOfUser := CreateUniqueTestUser(t, db, "related_blocker")

	// User blocks someone
	input1 := userblockservice.BlockInput{
		BlockerID:   user.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blockedByUser.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input1)
	require.NoError(t, err)

	// Someone blocks user
	input2 := userblockservice.BlockInput{
		BlockerID:   blockerOfUser.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err = svc.Block(ctx, input2)
	require.NoError(t, err)

	// Get all related user IDs (both directions)
	ids, err := svc.GetAllBlockRelatedUserIDs(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, blockedByUser.ID)
	assert.Contains(t, ids, blockerOfUser.ID)
}

// ============================================================================
// Statistics Tests
// ============================================================================

func TestUserBlockService_GetStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	// Create active blocks
	for i := 0; i < 2; i++ {
		blocker := CreateUniqueTestUser(t, db, "stats_blocker")
		blocked := CreateUniqueTestUser(t, db, "stats_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	// Create canceled block
	blocker := CreateUniqueTestUser(t, db, "stats_cancel_blocker")
	blocked := CreateUniqueTestUser(t, db, "stats_cancel_blocked")
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)
	svc.Unblock(ctx, blocker.ID, blocked.ID)

	// Get stats
	stats, err := svc.GetStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.BlockStatusActive], int64(2))
	assert.GreaterOrEqual(t, stats[model.BlockStatusCanceled], int64(1))

	// Create admin canceled block
	admin := CreateUniqueTestUser(t, db, "admin")
	err = svc.AdminUnblock(ctx, block.ID, admin.ID, "Admin cancel")
	require.NoError(t, err)

	stats, err = svc.GetStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.BlockStatusAdminCanceled], int64(1))
}

func TestUserBlockService_GetActiveCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	// Create active blocks
	for i := 0; i < 3; i++ {
		blocker := CreateUniqueTestUser(t, db, "active_blocker")
		blocked := CreateUniqueTestUser(t, db, "active_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		_, err := svc.Block(ctx, input)
		require.NoError(t, err)
	}

	count, err := svc.GetActiveCount(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}

// ============================================================================
// Get Operation Tests
// ============================================================================

func TestUserBlockService_Get_WithRelations(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "get_blocker")
	blocked := CreateUniqueTestUser(t, db, "get_blocked")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Test with relations",
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Get with relations
	result, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, block.ID, result.ID)
	assert.NotNil(t, result.Blocker)
	assert.NotNil(t, result.Blocked)
	assert.Equal(t, blocker.Name, result.Blocker.Name)
	assert.Equal(t, blocked.Name, result.Blocked.Name)
}

func TestUserBlockService_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	_, err := svc.Get(ctx, 99999)
	require.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
}

// ============================================================================
// Delete Operation Tests
// ============================================================================

func TestUserBlockService_Delete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "del_blocker")
	blocked := CreateUniqueTestUser(t, db, "del_blocked")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Delete
	err = svc.Delete(ctx, block.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.Get(ctx, block.ID)
	assert.Error(t, err)
}

// ============================================================================
// Batch Operations Tests
// ============================================================================

func TestUserBlockService_BatchBlock_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "batch_blocker")
	var items []userblockservice.BlockInputItemForBatch

	// Create batch items
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "batch_blocked")
		items = append(items, userblockservice.BlockInputItemForBatch{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Reason:      "Batch block",
		})
	}

	result, err := svc.BatchBlock(ctx, items)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)
}

func TestUserBlockService_BatchBlock_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "partial_blocker")
	blocked1 := CreateUniqueTestUser(t, db, "partial_blocked1")

	// Block first user
	input1 := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked1.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	_, err := svc.Block(ctx, input1)
	require.NoError(t, err)

	items := []userblockservice.BlockInputItemForBatch{
		{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked1.ID, // Already blocked
			BlockedType: model.BlockUserTypeUser,
		},
		{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   99999, // Not found
			BlockedType: model.BlockUserTypeUser,
		},
		{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocker.ID, // Cannot block self
			BlockedType: model.BlockUserTypeUser,
		},
	}

	result, err := svc.BatchBlock(ctx, items)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 3, result.FailedCount)
	assert.Len(t, result.FailedIDs, 3)
}

func TestUserBlockService_BatchUnblock_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "batchun_blocker")
	admin := CreateUniqueTestUser(t, db, "batchun_admin")
	var blockIDs []uint64

	// Create blocks
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "batchun_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		block, err := svc.Block(ctx, input)
		require.NoError(t, err)
		blockIDs = append(blockIDs, block.ID)
	}

	result, err := svc.BatchUnblock(ctx, blockIDs, admin.ID, "Batch unblock")
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserBlockService_BatchUnblock_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "partialun_blocker")
	blocked := CreateUniqueTestUser(t, db, "partialun_blocked")
	admin := CreateUniqueTestUser(t, db, "partialun_admin")

	// Create one block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	ids := []uint64{
		block.ID, // Valid
		99999,    // Not found
		99998,    // Not found
	}

	result, err := svc.BatchUnblock(ctx, ids, admin.ID, "Batch partial")
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)
	assert.Contains(t, result.FailedIDs, uint64(99999))
	assert.Contains(t, result.FailedIDs, uint64(99998))
}

func TestUserBlockService_BatchDelete_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "batchdel_blocker")
	var blockIDs []uint64

	// Create blocks
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "batchdel_blocked")
		input := userblockservice.BlockInput{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
		}
		block, err := svc.Block(ctx, input)
		require.NoError(t, err)
		blockIDs = append(blockIDs, block.ID)
	}

	result, err := svc.BatchDelete(ctx, blockIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserBlockService_BatchDelete_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "partialdel_blocker")
	blocked := CreateUniqueTestUser(t, db, "partialdel_blocked")

	// Create one block
	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	ids := []uint64{
		block.ID, // Valid
		99999,    // Not found
		99998,    // Not found
	}

	result, err := svc.BatchDelete(ctx, ids)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
}

// ============================================================================
// Edge Cases and Complex Scenarios
// ============================================================================

func TestUserBlockService_MutualBlock_SameUsers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "mutual1")
	user2 := CreateUniqueTestUser(t, db, "mutual2")

	// User1 blocks User2
	input1 := userblockservice.BlockInput{
		BlockerID:   user1.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user2.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "User1 blocks User2",
	}
	_, err := svc.Block(ctx, input1)
	require.NoError(t, err)

	// User2 blocks User1 (should be allowed - different direction)
	input2 := userblockservice.BlockInput{
		BlockerID:   user2.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user1.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "User2 blocks User1",
	}
	_, err = svc.Block(ctx, input2)
	require.NoError(t, err)

	// Both should be blocked from each other
	blocked, err := svc.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.True(t, blocked)

	// Each should have their own block record
	blocks1, err := svc.ListByBlocker(ctx, user1.ID, true)
	require.NoError(t, err)
	assert.Len(t, blocks1, 1)

	blocks2, err := svc.ListByBlocker(ctx, user2.ID, true)
	require.NoError(t, err)
	assert.Len(t, blocks2, 1)
}

func TestUserBlockService_StatusTransition_ActiveToCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "trans_blocker")
	blocked := CreateUniqueTestUser(t, db, "trans_blocked")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Verify initial status
	assert.Equal(t, model.BlockStatusActive, block.Status)

	// User cancels
	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)

	updated, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BlockStatusCanceled, updated.Status)
	assert.NotNil(t, updated.CanceledAt)
}

func TestUserBlockService_StatusTransition_ActiveToAdminCanceled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "admtrans_blocker")
	blocked := CreateUniqueTestUser(t, db, "admtrans_blocked")
	admin := CreateUniqueTestUser(t, db, "admin")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Admin cancels
	err = svc.AdminUnblock(ctx, block.ID, admin.ID, "Admin intervention")
	require.NoError(t, err)

	updated, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BlockStatusAdminCanceled, updated.Status)
	assert.NotNil(t, updated.CanceledAt)
	assert.NotNil(t, updated.CanceledBy)
	assert.Equal(t, admin.ID, *updated.CanceledBy)
}

func TestUserBlockService_ListByBlocker_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "empty_blocker")

	blocks, err := svc.ListByBlocker(ctx, user.ID, true)
	require.NoError(t, err)
	assert.Empty(t, blocks)
}

func TestUserBlockService_WithPlayers_BothDirections(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	player1User := CreateUniqueTestUser(t, db, "player1")
	_ = CreateTestPlayer(t, db, player1User)

	player2User := CreateUniqueTestUser(t, db, "player2")
	_ = CreateTestPlayer(t, db, player2User)

	// Player1 blocks Player2
	input := userblockservice.BlockInput{
		BlockerID:   player1User.ID,
		BlockerType: model.BlockUserTypePlayer,
		BlockedID:   player2User.ID,
		BlockedType: model.BlockUserTypePlayer,
		Reason:      "Competition blocking",
	}
	_, err := svc.Block(ctx, input)
	require.NoError(t, err)

	// Verify
	isBlocked, err := svc.IsBlocked(ctx, player1User.ID, player2User.ID)
	require.NoError(t, err)
	assert.True(t, isBlocked)

	isBlockedBy, err := svc.IsBlockedBy(ctx, player1User.ID, player2User.ID)
	require.NoError(t, err)
	assert.True(t, isBlockedBy)
}

func TestUserBlockService_CreatedAtTimestamp(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "time_blocker")
	blocked := CreateUniqueTestUser(t, db, "time_blocked")

	beforeBlock := time.Now()

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	afterBlock := time.Now()

	// Verify BlockedAt is set
	assert.False(t, block.BlockedAt.IsZero())
	assert.True(t, block.BlockedAt.After(beforeBlock) || block.BlockedAt.Equal(beforeBlock))
	assert.True(t, block.BlockedAt.Before(afterBlock) || block.BlockedAt.Equal(afterBlock))
}

func TestUserBlockService_CancelTimestamp(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := newUserBlockService(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "canceltime_blocker")
	blocked := CreateUniqueTestUser(t, db, "canceltime_blocked")

	input := userblockservice.BlockInput{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
	}
	block, err := svc.Block(ctx, input)
	require.NoError(t, err)

	beforeCancel := time.Now()
	err = svc.Unblock(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)
	afterCancel := time.Now()

	updated, err := svc.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.NotNil(t, updated.CanceledAt)
	assert.True(t, updated.CanceledAt.After(beforeCancel) || updated.CanceledAt.Equal(beforeCancel))
	assert.True(t, updated.CanceledAt.Before(afterCancel) || updated.CanceledAt.Equal(afterCancel))
}
