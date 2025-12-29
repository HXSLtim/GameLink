package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/userblock"
)

// ============================================================================
// UserBlock CRUD Tests
// ============================================================================

func TestUserBlockRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "blocker")
	blocked := CreateUniqueTestUser(t, db, "blocked")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Test block reason",
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}

	err := repo.Create(ctx, block)
	require.NoError(t, err)
	assert.NotZero(t, block.ID)
}

func TestUserBlockRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "get_blocker")
	blocked := CreateUniqueTestUser(t, db, "get_blocked")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Reason:      "Get test",
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	got, err := repo.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, blocker.ID, got.BlockerID)
	assert.Equal(t, blocked.ID, got.BlockedID)
	assert.Equal(t, "Get test", got.Reason)
}

func TestUserBlockRepository_GetByBlockerAndBlocked(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "pair_blocker")
	blocked := CreateUniqueTestUser(t, db, "pair_blocked")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	got, err := repo.GetByBlockerAndBlocked(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)
	assert.Equal(t, block.ID, got.ID)
}

func TestUserBlockRepository_IsBlocked_Bidirectional(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "bidir_user1")
	user2 := CreateUniqueTestUser(t, db, "bidir_user2")

	// Initially not blocked
	blocked, err := repo.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.False(t, blocked)

	// User1 blocks User2
	block := &model.UserBlock{
		BlockerID:   user1.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user2.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	// Now blocked (either direction)
	blocked, err = repo.IsBlocked(ctx, user1.ID, user2.ID)
	require.NoError(t, err)
	assert.True(t, blocked)

	// Also blocked in reverse direction check
	blocked, err = repo.IsBlocked(ctx, user2.ID, user1.ID)
	require.NoError(t, err)
	assert.True(t, blocked)
}

func TestUserBlockRepository_IsBlockedBy(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "isblocked_blocker")
	blocked := CreateUniqueTestUser(t, db, "isblocked_blocked")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	// blocked is blocked by blocker
	isBlocked, err := repo.IsBlockedBy(ctx, blocker.ID, blocked.ID)
	require.NoError(t, err)
	assert.True(t, isBlocked)

	// blocker is NOT blocked by blocked (direction matters)
	isBlocked, err = repo.IsBlockedBy(ctx, blocked.ID, blocker.ID)
	require.NoError(t, err)
	assert.False(t, isBlocked)
}

func TestUserBlockRepository_ListByBlockerID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "list_blocker")

	// Block multiple users
	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "list_blocked")
		block := &model.UserBlock{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Status:      model.BlockStatusActive,
			BlockedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, block))
	}

	blocks, err := repo.ListByBlockerID(ctx, blocker.ID, nil)
	require.NoError(t, err)
	assert.Len(t, blocks, 3)
}

func TestUserBlockRepository_ListByBlockerID_WithStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "status_blocker")

	// Create active block
	blocked1 := CreateUniqueTestUser(t, db, "status_blocked1")
	block1 := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked1.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block1))

	// Create canceled block
	blocked2 := CreateUniqueTestUser(t, db, "status_blocked2")
	block2 := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked2.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusCanceled,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block2))

	// Filter by active status
	activeStatus := model.BlockStatusActive
	blocks, err := repo.ListByBlockerID(ctx, blocker.ID, &activeStatus)
	require.NoError(t, err)
	assert.Len(t, blocks, 1)
	assert.Equal(t, model.BlockStatusActive, blocks[0].Status)
}

func TestUserBlockRepository_UpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "upd_blocker")
	blocked := CreateUniqueTestUser(t, db, "upd_blocked")
	admin := CreateUniqueTestUser(t, db, "upd_admin")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	// Cancel by admin
	err := repo.UpdateStatus(ctx, block.ID, model.BlockStatusAdminCanceled, &admin.ID, "Admin canceled")
	require.NoError(t, err)

	got, err := repo.Get(ctx, block.ID)
	require.NoError(t, err)
	assert.Equal(t, model.BlockStatusAdminCanceled, got.Status)
	assert.Equal(t, "Admin canceled", got.AdminRemark)
	assert.NotNil(t, got.CanceledAt)
}

func TestUserBlockRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "del_blocker")
	blocked := CreateUniqueTestUser(t, db, "del_blocked")

	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	err := repo.Delete(ctx, block.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, block.ID)
	assert.Error(t, err)
}

func TestUserBlockRepository_GetBlockedUserIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	blocker := CreateUniqueTestUser(t, db, "ids_blocker")
	var expectedIDs []uint64

	for i := 0; i < 3; i++ {
		blocked := CreateUniqueTestUser(t, db, "ids_blocked")
		expectedIDs = append(expectedIDs, blocked.ID)
		block := &model.UserBlock{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Status:      model.BlockStatusActive,
			BlockedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, block))
	}

	ids, err := repo.GetBlockedUserIDs(ctx, blocker.ID)
	require.NoError(t, err)
	assert.Len(t, ids, 3)
}

func TestUserBlockRepository_GetAllBlockRelatedUserIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "related_user")
	blockedByUser := CreateUniqueTestUser(t, db, "related_blocked")
	blockerOfUser := CreateUniqueTestUser(t, db, "related_blocker")

	// User blocks someone
	block1 := &model.UserBlock{
		BlockerID:   user.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blockedByUser.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block1))

	// Someone blocks user
	block2 := &model.UserBlock{
		BlockerID:   blockerOfUser.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   user.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block2))

	ids, err := repo.GetAllBlockRelatedUserIDs(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, blockedByUser.ID)
	assert.Contains(t, ids, blockerOfUser.ID)
}

func TestUserBlockRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	// Create multiple blocks
	for i := 0; i < 5; i++ {
		blocker := CreateUniqueTestUser(t, db, "paged_blocker")
		blocked := CreateUniqueTestUser(t, db, "paged_blocked")
		block := &model.UserBlock{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Status:      model.BlockStatusActive,
			BlockedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, block))
	}

	blocks, total, err := repo.ListPaged(ctx, repository.UserBlockListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.LessOrEqual(t, len(blocks), 3)
}

func TestUserBlockRepository_CountByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	// Create active blocks
	for i := 0; i < 2; i++ {
		blocker := CreateUniqueTestUser(t, db, "count_blocker")
		blocked := CreateUniqueTestUser(t, db, "count_blocked")
		block := &model.UserBlock{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Status:      model.BlockStatusActive,
			BlockedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, block))
	}

	// Create canceled block
	blocker := CreateUniqueTestUser(t, db, "count_blocker_c")
	blocked := CreateUniqueTestUser(t, db, "count_blocked_c")
	block := &model.UserBlock{
		BlockerID:   blocker.ID,
		BlockerType: model.BlockUserTypeUser,
		BlockedID:   blocked.ID,
		BlockedType: model.BlockUserTypeUser,
		Status:      model.BlockStatusCanceled,
		BlockedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, block))

	stats, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.BlockStatusActive], int64(2))
	assert.GreaterOrEqual(t, stats[model.BlockStatusCanceled], int64(1))
}

func TestUserBlockRepository_GetActiveCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := userblock.NewUserBlockRepository(db)
	ctx := context.Background()

	// Create active blocks
	for i := 0; i < 3; i++ {
		blocker := CreateUniqueTestUser(t, db, "active_blocker")
		blocked := CreateUniqueTestUser(t, db, "active_blocked")
		block := &model.UserBlock{
			BlockerID:   blocker.ID,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   blocked.ID,
			BlockedType: model.BlockUserTypeUser,
			Status:      model.BlockStatusActive,
			BlockedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, block))
	}

	count, err := repo.GetActiveCount(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}
