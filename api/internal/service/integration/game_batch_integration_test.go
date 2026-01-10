// Package integration provides integration tests for game batch operations.
package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/service/admin"
)

// ============================================================================
// Game Batch Operations Tests
// ============================================================================

// createTestAdminService creates an admin service with only game repository populated.
func createTestAdminService(db *gorm.DB, gameRepo repository.GameRepository) *admin.AdminService {
	return admin.NewAdminService(
		gameRepo,
		nil, // users
		nil, // players
		nil, // orders
		nil, // payments
		nil, // roles
		nil, // serviceItems
		nil, // permissions
		nil, // menus
		nil, // stats
		nil, // wallets
		gamecategory.NewGameCategoryRepository(db),
		nil, // cache
	)
}

// CreateTestGameWithCategory creates a test game with specified category.
func CreateTestGameWithCategory(t *testing.T, db *gorm.DB, name, category string) *model.Game {
	t.Helper()
	game := &model.Game{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Key:      name,
		Name:     name,
		Category: category,
		IsActive: true,
	}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return game
}

func TestGameBatch_BatchUpdateGamesStatusWithResponse_BatchEnable(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create test games - some active, some inactive
	var gameIDs []uint64
	for i := 0; i < 3; i++ {
		game := CreateTestGameWithCategory(t, db, "enable_game_"+string(rune('A'+i)), "moba")
		if i < 2 {
			// First 2 games are inactive
			game.IsActive = false
			require.NoError(t, db.Save(game).Error)
		}
		gameIDs = append(gameIDs, game.ID)
	}

	// Batch enable all games
	response, err := adminSvc.BatchUpdateGamesStatusWithResponse(ctx, gameIDs, true)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
	assert.Len(t, response.SuccessItems, 3)
	assert.Len(t, response.FailedItems, 0)

	// Verify database state - all games should be active
	for _, id := range gameIDs {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.True(t, game.IsActive, "Game %d should be active", id)
	}
}

func TestGameBatch_BatchUpdateGamesStatusWithResponse_BatchDisable(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create test games - some active, some inactive
	var gameIDs []uint64
	for i := 0; i < 3; i++ {
		game := CreateTestGameWithCategory(t, db, "disable_game_"+string(rune('A'+i)), "fps")
		if i < 2 {
			// First 2 games are active
			game.IsActive = true
			require.NoError(t, db.Save(game).Error)
		}
		gameIDs = append(gameIDs, game.ID)
	}

	// Batch disable all games
	response, err := adminSvc.BatchUpdateGamesStatusWithResponse(ctx, gameIDs, false)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
	assert.Len(t, response.SuccessItems, 3)
	assert.Len(t, response.FailedItems, 0)

	// Verify database state - all games should be inactive
	for _, id := range gameIDs {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.False(t, game.IsActive, "Game %d should be inactive", id)
	}
}

func TestGameBatch_BatchUpdateGamesStatusWithResponse_PartialGamesNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create 2 valid games
	var gameIDs []uint64
	for i := 0; i < 2; i++ {
		game := CreateTestGameWithCategory(t, db, "partial_game_"+string(rune('A'+i)), "rpg")
		gameIDs = append(gameIDs, game.ID)
	}

	// Add non-existent game IDs
	gameIDs = append(gameIDs, 99998, 99999)

	// Batch update status
	response, err := adminSvc.BatchUpdateGamesStatusWithResponse(ctx, gameIDs, true)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 2, response.FailedCount)
	assert.Len(t, response.SuccessItems, 2)
	assert.Len(t, response.FailedItems, 2)

	// Verify failed items contain correct IDs and error messages
	failedIDs := make([]uint64, 0)
	for _, item := range response.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "game not found")
	}
	assert.Contains(t, failedIDs, uint64(99998))
	assert.Contains(t, failedIDs, uint64(99999))

	// Verify valid games were updated
	for _, id := range gameIDs[:2] {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.True(t, game.IsActive)
	}
}

func TestGameBatch_BatchUpdateGamesStatusWithResponse_EmptyIDList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Batch update with empty ID list
	response, err := adminSvc.BatchUpdateGamesStatusWithResponse(ctx, []uint64{}, true)

	// Should return error response (not nil error)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGameBatch_BatchUpdateGamesCategory_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create test games with different categories
	var gameIDs []uint64
	categories := []string{"moba", "fps", "rpg"}
	for i := 0; i < 3; i++ {
		game := CreateTestGameWithCategory(t, db, "category_game_"+string(rune('A'+i)), categories[i])
		gameIDs = append(gameIDs, game.ID)
	}

	// Batch update all games to "strategy" category
	newCategory := "strategy"
	response, err := adminSvc.BatchUpdateGamesCategory(ctx, gameIDs, newCategory)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)
	assert.Len(t, response.SuccessItems, 3)
	assert.Len(t, response.FailedItems, 0)

	// Verify database state - all games should have new category
	for _, id := range gameIDs {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newCategory, game.Category, "Game %d should have category %s", id, newCategory)
	}
}

func TestGameBatch_BatchUpdateGamesCategory_PartialGamesNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create 2 valid games
	var gameIDs []uint64
	for i := 0; i < 2; i++ {
		game := CreateTestGameWithCategory(t, db, "partial_cat_game_"+string(rune('A'+i)), "moba")
		gameIDs = append(gameIDs, game.ID)
	}

	// Add non-existent game IDs
	gameIDs = append(gameIDs, 88888, 99999)

	// Batch update category
	newCategory := "action"
	response, err := adminSvc.BatchUpdateGamesCategory(ctx, gameIDs, newCategory)
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, 4, response.TotalCount)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 2, response.FailedCount)
	assert.Len(t, response.SuccessItems, 2)
	assert.Len(t, response.FailedItems, 2)

	// Verify failed items contain correct IDs and error messages
	failedIDs := make([]uint64, 0)
	for _, item := range response.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "game not found")
	}
	assert.Contains(t, failedIDs, uint64(88888))
	assert.Contains(t, failedIDs, uint64(99999))

	// Verify valid games were updated
	for _, id := range gameIDs[:2] {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newCategory, game.Category)
	}
}

func TestGameBatch_BatchUpdateGamesCategory_EmptyIDList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Batch update with empty ID list
	response, err := adminSvc.BatchUpdateGamesCategory(ctx, []uint64{}, "strategy")

	// Should return error response (not nil error)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGameBatch_BatchUpdateGamesCategory_WhitespaceCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create admin service
	gameRepo := game.NewGameRepository(db)
	adminSvc := createTestAdminService(db, gameRepo)
	ctx := context.Background()

	// Create test games
	var gameIDs []uint64
	for i := 0; i < 3; i++ {
		game := CreateTestGameWithCategory(t, db, "whitespace_cat_"+string(rune('A'+i)), "moba")
		gameIDs = append(gameIDs, game.ID)
	}

	// Batch update with whitespace-only category (should fail)
	response, err := adminSvc.BatchUpdateGamesCategory(ctx, gameIDs, "   ")
	require.NoError(t, err)

	// Verify response - all should fail due to empty category after trim
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 0, response.SuccessCount)
	assert.Equal(t, 3, response.FailedCount)
	assert.Len(t, response.FailedItems, 3)

	// Verify error messages
	for _, item := range response.FailedItems {
		assert.Contains(t, item.Message, "category cannot be empty")
	}

	// Verify database state unchanged
	for _, id := range gameIDs {
		game, err := gameRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "moba", game.Category) // Original category unchanged
	}
}
