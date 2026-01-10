package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"

	"gorm.io/gorm"
)

// ============================================================================
// Helper Functions for Game Tests
// ============================================================================

// CreateTestGameWithKey creates a test game with the given key and name.
func CreateTestGameWithKey(t *testing.T, db *gorm.DB, key, name string) *model.Game {
	t.Helper()
	game := &model.Game{
		Base:        model.Base{ExtJSON: "{}"},
		Key:         key,
		Name:        name,
		Category:    "test",
		IconURL:     "https://example.com/" + key + ".png",
		Description: "Test game description for " + name,
		IsActive:    true,
		SortOrder:   0,
	}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return game
}

// CreateTestGameWithCategoryID creates a test game associated with a category ID.
func CreateTestGameWithCategoryID(t *testing.T, db *gorm.DB, key, name string, categoryID *uint64) *model.Game {
	t.Helper()
	game := &model.Game{
		Base:        model.Base{ExtJSON: "{}"},
		Key:         key,
		Name:        name,
		CategoryID:  categoryID,
		IconURL:     "https://example.com/" + key + ".png",
		Description: "Test game description",
		IsActive:    true,
		SortOrder:   0,
	}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return game
}

// CreateTestGameWithStatus creates a test game with specific active status.
func CreateTestGameWithStatus(t *testing.T, db *gorm.DB, key, name string, isActive bool) *model.Game {
	t.Helper()
	game := &model.Game{
		Base:      model.Base{ExtJSON: "{}"},
		Key:       key,
		Name:      name,
		Category:  "test",
		IsActive:  isActive,
		SortOrder: 0,
	}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return game
}

// GetGameCount returns the total number of games.
func GetGameCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.Game{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count games: %v", err)
	}
	return count
}

// GetActiveGameCount returns the number of active games.
func GetActiveGameCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.Game{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count active games: %v", err)
	}
	return count
}

// GetGamesByCategoryName returns games with the given category name.
func GetGamesByCategoryName(t *testing.T, db *gorm.DB, categoryName string) []model.Game {
	t.Helper()
	var games []model.Game
	if err := db.Where("category = ?", categoryName).Find(&games).Error; err != nil {
		t.Fatalf("Failed to get games by category: %v", err)
	}
	return games
}

// ============================================================================
// Repository Tests - Basic CRUD
// ============================================================================

func TestGameRepository_Create_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	gameObj := &model.Game{
		Key:         fmt.Sprintf("test_game_%d", time.Now().UnixNano()),
		Name:        "Test Game",
		Category:    "moba",
		IconURL:     "https://example.com/icon.png",
		Description: "A test game",
	}

	err := repo.Create(ctx, gameObj)
	require.NoError(t, err)
	assert.NotZero(t, gameObj.ID)
}

func TestGameRepository_Create_DuplicateKey(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	key := fmt.Sprintf("dup_key_%d", time.Now().UnixNano())

	// Create first game
	game1 := &model.Game{
		Key:      key,
		Name:     "First Game",
		Category: "test",
		IsActive: true,
	}
	err := repo.Create(ctx, game1)
	require.NoError(t, err)

	// Try to create duplicate key
	game2 := &model.Game{
		Key:      key,
		Name:     "Second Game",
		Category: "test",
		IsActive: true,
	}
	err = repo.Create(ctx, game2)
	assert.Error(t, err)
}

func TestGameRepository_Get_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	gameObj := &model.Game{
		Key:      fmt.Sprintf("get_game_%d", time.Now().UnixNano()),
		Name:     "Get Test Game",
		Category: "fps",
	}
	require.NoError(t, repo.Create(ctx, gameObj))

	got, err := repo.Get(ctx, gameObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Test Game", got.Name)
	assert.Equal(t, "fps", got.Category)
}

func TestGameRepository_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, 999999)
	assert.Error(t, err)
}

func TestGameRepository_Update_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	gameObj := &model.Game{
		Key:      fmt.Sprintf("update_game_%d", time.Now().UnixNano()),
		Name:     "Update Test",
		Category: "rpg",
	}
	require.NoError(t, repo.Create(ctx, gameObj))

	// Update
	gameObj.Name = "Updated Name"
	gameObj.Category = "mmorpg"
	err := repo.Update(ctx, gameObj)
	require.NoError(t, err)

	// Verify
	got, err := repo.Get(ctx, gameObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "mmorpg", got.Category)
}

func TestGameRepository_Delete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	gameObj := &model.Game{
		Key:      fmt.Sprintf("delete_game_%d", time.Now().UnixNano()),
		Name:     "Delete Test",
		Category: "card",
	}
	require.NoError(t, repo.Create(ctx, gameObj))

	err := repo.Delete(ctx, gameObj.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = repo.Get(ctx, gameObj.ID)
	assert.Error(t, err)
}

func TestGameRepository_List_All(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create multiple games
	for i := 0; i < 3; i++ {
		gameObj := &model.Game{
			Key:      fmt.Sprintf("list_game_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("List Game %d", i),
			Category: "moba",
		}
		require.NoError(t, repo.Create(ctx, gameObj))
	}

	games, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(games), 3)
}

func TestGameRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games
	for i := 0; i < 5; i++ {
		gameObj := &model.Game{
			Key:      fmt.Sprintf("paged_game_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("Paged Game %d", i),
			Category: "fps",
		}
		require.NoError(t, repo.Create(ctx, gameObj))
	}

	games, total, err := repo.ListPaged(ctx, 1, 3)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.LessOrEqual(t, len(games), 3)
}

func TestGameRepository_ListPagedWithFilter(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games with specific names
	gameObj := &model.Game{
		Key:      fmt.Sprintf("filter_unique_%d", time.Now().UnixNano()),
		Name:     "UniqueFilterName",
		Category: "rpg",
	}
	require.NoError(t, repo.Create(ctx, gameObj))

	games, total, err := repo.ListPagedWithFilter(ctx, 1, 10, "UniqueFilter")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(games), 1)
}

func TestGameRepository_BatchDelete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	var ids []uint64
	for i := 0; i < 3; i++ {
		gameObj := &model.Game{
			Key:      fmt.Sprintf("batch_del_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("Batch Delete %d", i),
			Category: "card",
		}
		require.NoError(t, repo.Create(ctx, gameObj))
		ids = append(ids, gameObj.ID)
	}

	affected, err := repo.BatchDelete(ctx, ids)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)

	// Verify all deleted
	for _, id := range ids {
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)
	}
}

// ============================================================================
// Batch Operations Tests
// ============================================================================

func TestGameRepository_BatchUpdateStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create multiple active games
	var ids []uint64
	for i := 0; i < 3; i++ {
		gameObj := &model.Game{
			Key:      fmt.Sprintf("status_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("Status Game %d", i),
			Category: "test",
			IsActive: true,
		}
		require.NoError(t, repo.Create(ctx, gameObj))
		ids = append(ids, gameObj.ID)
	}

	// Batch update to inactive
	affected, err := repo.BatchUpdateStatus(ctx, ids, false)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)

	// Verify all are now inactive
	for _, id := range ids {
		game, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.False(t, game.IsActive)
	}
}

func TestGameRepository_BatchUpdateStatus_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Empty list should not error
	affected, err := repo.BatchUpdateStatus(ctx, []uint64{}, false)
	require.NoError(t, err)
	assert.Equal(t, int64(0), affected)
}

func TestGameRepository_BatchUpdateSortOrder_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games
	game1 := CreateTestGameWithKey(t, db, "sort1", "Game 1")
	game2 := CreateTestGameWithKey(t, db, "sort2", "Game 2")
	game3 := CreateTestGameWithKey(t, db, "sort3", "Game 3")

	// Update sort orders
	updates := map[uint64]int{
		game1.ID: 10,
		game2.ID: 20,
		game3.ID: 30,
	}

	affected, err := repo.BatchUpdateSortOrder(ctx, updates)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)

	// Verify updates
	g1, _ := repo.Get(ctx, game1.ID)
	g2, _ := repo.Get(ctx, game2.ID)
	g3, _ := repo.Get(ctx, game3.ID)

	assert.Equal(t, 10, g1.SortOrder)
	assert.Equal(t, 20, g2.SortOrder)
	assert.Equal(t, 30, g3.SortOrder)
}

func TestGameRepository_BatchUpdateCategory_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games with different categories
	game1 := CreateTestGameWithKey(t, db, "cat1", "Game 1")
	game1.Category = "old_cat"
	db.Save(game1)

	game2 := CreateTestGameWithKey(t, db, "cat2", "Game 2")
	game2.Category = "old_cat"
	db.Save(game2)

	ids := []uint64{game1.ID, game2.ID}

	// Batch update category
	affected, err := repo.BatchUpdateCategory(ctx, ids, "new_cat")
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	// Verify updates
	g1, _ := repo.Get(ctx, game1.ID)
	g2, _ := repo.Get(ctx, game2.ID)

	assert.Equal(t, "new_cat", g1.Category)
	assert.Equal(t, "new_cat", g2.Category)
}

// ============================================================================
// Game-GameCategory Relationship Tests
// ============================================================================

func TestGame_GameCategoryRelationship_CreateWithCategoryID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create category
	category := CreateTestCategory(t, db, "Action")

	// Create game with category ID
	game := &model.Game{
		Key:        fmt.Sprintf("game_cat_%d", time.Now().UnixNano()),
		Name:       "Action Game",
		CategoryID: &category.ID,
		IsActive:   true,
	}
	require.NoError(t, db.Create(game).Error)

	// Verify relationship
	assert.Equal(t, category.ID, *game.CategoryID)
}

func TestGame_GameCategoryRelationship_CategoryHasManyGames(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create category
	category := CreateTestCategory(t, db, "Strategy")

	// Create multiple games in this category
	game1 := CreateTestGameWithCategoryID(t, db, "strat1", "Strategy Game 1", &category.ID)
	game2 := CreateTestGameWithCategoryID(t, db, "strat2", "Strategy Game 2", &category.ID)
	game3 := CreateTestGameWithCategoryID(t, db, "strat3", "Strategy Game 3", &category.ID)

	// Verify all games have the same category ID
	assert.Equal(t, category.ID, *game1.CategoryID)
	assert.Equal(t, category.ID, *game2.CategoryID)
	assert.Equal(t, category.ID, *game3.CategoryID)
}

func TestGame_GameCategoryRelationship_UpdateGameCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create two categories
	cat1 := CreateTestCategory(t, db, "Category1")
	cat2 := CreateTestCategory(t, db, "Category2")

	// Create game in first category
	game := CreateTestGameWithCategoryID(t, db, "movecat", "Move Category Game", &cat1.ID)
	assert.Equal(t, cat1.ID, *game.CategoryID)

	// Update to second category
	game.CategoryID = &cat2.ID
	require.NoError(t, db.Save(game).Error)

	// Verify update
	var fetched model.Game
	require.NoError(t, db.First(&fetched, game.ID).Error)
	assert.Equal(t, cat2.ID, *fetched.CategoryID)
}

func TestGame_GameCategoryRelationship_DeleteCategoryDoesNotDeleteGames(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamecategory.NewGameCategoryRepository(db)

	// Create category
	category := CreateTestCategory(t, db, "Adventure")

	// Create games in this category
	game1 := CreateTestGameWithCategoryID(t, db, "adv1", "Adventure Game 1", &category.ID)
	game2 := CreateTestGameWithCategoryID(t, db, "adv2", "Adventure Game 2", &category.ID)

	// Delete category
	err := repo.Delete(ctx, category.ID)
	require.NoError(t, err)

	// Verify category is deleted
	_, err = repo.Get(ctx, category.ID)
	assert.Error(t, err)

	// Verify games still exist (soft delete with SET NULL)
	var games []model.Game
	err = db.Where("id IN ?", []uint64{game1.ID, game2.ID}).Find(&games).Error
	require.NoError(t, err)
	assert.Len(t, games, 2)
}

func TestGame_GameCategoryRelationship_QueryGamesByCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create categories
	cat1 := CreateTestCategory(t, db, "RPG")
	cat2 := CreateTestCategory(t, db, "FPS")

	// Create games in different categories
	CreateTestGameWithCategoryID(t, db, "rpg1", "RPG Game 1", &cat1.ID)
	CreateTestGameWithCategoryID(t, db, "rpg2", "RPG Game 2", &cat1.ID)
	CreateTestGameWithCategoryID(t, db, "fps1", "FPS Game 1", &cat2.ID)

	// Query games by category ID
	var rpgGames []model.Game
	err := db.Where("category_id = ?", cat1.ID).Find(&rpgGames).Error
	require.NoError(t, err)
	assert.Len(t, rpgGames, 2)

	var fpsGames []model.Game
	err = db.Where("category_id = ?", cat2.ID).Find(&fpsGames).Error
	require.NoError(t, err)
	assert.Len(t, fpsGames, 1)
}

func TestGame_GameCategoryRelationship_CategoryWithNullCategoryID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create game without category
	game := &model.Game{
		Key:        fmt.Sprintf("nocat_%d", time.Now().UnixNano()),
		Name:       "Game Without Category",
		CategoryID: nil,
		IsActive:   true,
	}
	require.NoError(t, db.Create(game).Error)

	// Verify category ID is nil
	assert.Nil(t, game.CategoryID)

	// Verify we can query games without category
	var games []model.Game
	err := db.Where("category_id IS NULL").Find(&games).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(games), 1)
}

// ============================================================================
// Validation and Uniqueness Tests
// ============================================================================

func TestGame_KeyUniqueness(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	key := fmt.Sprintf("unique_key_%d", time.Now().UnixNano())

	// Create first game with unique key
	game1 := &model.Game{
		Key:      key,
		Name:     "First Game",
		Category: "test",
		IsActive: true,
	}
	err := repo.Create(ctx, game1)
	require.NoError(t, err)

	// Try to create second game with same key
	game2 := &model.Game{
		Key:      key,
		Name:     "Second Game",
		Category: "test",
		IsActive: true,
	}
	err = repo.Create(ctx, game2)
	assert.Error(t, err)
}

func TestGame_KeyValidation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Key is required (database constraint)
	game := &model.Game{
		Key:      "",
		Name:     "Test Game",
		Category: "test",
	}

	err := repo.Create(ctx, game)
	assert.Error(t, err)
}

func TestGame_NameValidation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Name is required (database constraint)
	game := &model.Game{
		Key:      fmt.Sprintf("noname_%d", time.Now().UnixNano()),
		Name:     "",
		Category: "test",
	}

	err := repo.Create(ctx, game)
	assert.Error(t, err)
}

func TestGame_IsActiveToggle(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create active game
	game := &model.Game{
		Key:      fmt.Sprintf("toggle_%d", time.Now().UnixNano()),
		Name:     "Toggle Game",
		Category: "test",
		IsActive: true,
	}
	require.NoError(t, repo.Create(ctx, game))
	assert.True(t, game.IsActive)

	// Toggle to inactive
	_, err := repo.BatchUpdateStatus(ctx, []uint64{game.ID}, false)
	require.NoError(t, err)

	// Verify
	fetched, _ := repo.Get(ctx, game.ID)
	assert.False(t, fetched.IsActive)

	// Toggle back to active
	_, err = repo.BatchUpdateStatus(ctx, []uint64{game.ID}, true)
	require.NoError(t, err)

	// Verify
	fetched, _ = repo.Get(ctx, game.ID)
	assert.True(t, fetched.IsActive)
}

// ============================================================================
// Sort Order Tests
// ============================================================================

func TestGame_SortOrder_Default(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create games without specifying sort order
	game1 := CreateTestGameWithKey(t, db, "sort_def1", "Default 1")
	game2 := CreateTestGameWithKey(t, db, "sort_def2", "Default 2")

	// Both should have default sort order 0
	assert.Equal(t, 0, game1.SortOrder)
	assert.Equal(t, 0, game2.SortOrder)
}

func TestGame_SortOrder_Custom(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games with custom sort order
	game1 := &model.Game{
		Key:       fmt.Sprintf("custom1_%d", time.Now().UnixNano()),
		Name:      "Custom 1",
		Category:  "test",
		SortOrder: 10,
		IsActive:  true,
	}
	game2 := &model.Game{
		Key:       fmt.Sprintf("custom2_%d", time.Now().UnixNano()),
		Name:      "Custom 2",
		Category:  "test",
		SortOrder: 20,
		IsActive:  true,
	}
	require.NoError(t, repo.Create(ctx, game1))
	require.NoError(t, repo.Create(ctx, game2))

	// Verify sort orders
	assert.Equal(t, 10, game1.SortOrder)
	assert.Equal(t, 20, game2.SortOrder)
}

func TestGame_ListOrderedBySortOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create games with different sort orders
	game1 := CreateTestGameWithKey(t, db, "order1", "Order 1")
	game1.SortOrder = 30
	db.Save(game1)

	game2 := CreateTestGameWithKey(t, db, "order2", "Order 2")
	game2.SortOrder = 10
	db.Save(game2)

	game3 := CreateTestGameWithKey(t, db, "order3", "Order 3")
	game3.SortOrder = 20
	db.Save(game3)

	// List games (should be ordered by created_at desc by default)
	var games []model.Game
	err := db.Order("sort_order ASC").Find(&games).Error
	require.NoError(t, err)

	// Verify sort order
	assert.Equal(t, 10, games[0].SortOrder)
	assert.Equal(t, 20, games[1].SortOrder)
	assert.Equal(t, 30, games[2].SortOrder)
}

// ============================================================================
// Filter and Search Tests
// ============================================================================

func TestGame_FilterByKeyword_MatchesName(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with unique name
	game := &model.Game{
		Key:      fmt.Sprintf("search_%d", time.Now().UnixNano()),
		Name:     "UniqueSearchName",
		Category: "test",
	}
	require.NoError(t, repo.Create(ctx, game))

	// Search by keyword in name
	games, total, err := repo.ListPagedWithFilter(ctx, 1, 10, "UniqueSearch")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(games), 1)
	assert.Contains(t, games[0].Name, "UniqueSearch")
}

func TestGame_FilterByKeyword_MatchesKey(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with unique key
	game := &model.Game{
		Key:      "special_key_game",
		Name:     "Some Game",
		Category: "test",
	}
	require.NoError(t, repo.Create(ctx, game))

	// Search by keyword in key
	games, total, err := repo.ListPagedWithFilter(ctx, 1, 10, "special_key")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(games), 1)
	assert.Contains(t, games[0].Key, "special_key")
}

func TestGame_FilterByKeyword_MatchesCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with specific category
	game := &model.Game{
		Key:      fmt.Sprintf("cat_search_%d", time.Now().UnixNano()),
		Name:     "Game Name",
		Category: "UniqueCategory",
	}
	require.NoError(t, repo.Create(ctx, game))

	// Search by keyword in category
	games, total, err := repo.ListPagedWithFilter(ctx, 1, 10, "UniqueCategory")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(games), 1)
}

func TestGame_FilterByKeyword_MatchesDescription(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with unique description
	game := &model.Game{
		Key:         fmt.Sprintf("desc_search_%d", time.Now().UnixNano()),
		Name:        "Game Name",
		Category:    "test",
		Description: "This is a unique description text",
	}
	require.NoError(t, repo.Create(ctx, game))

	// Search by keyword in description
	games, total, err := repo.ListPagedWithFilter(ctx, 1, 10, "unique description")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(games), 1)
	assert.Contains(t, games[0].Description, "unique description")
}

func TestGame_FilterByActiveStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create active and inactive games
	activeGame1 := CreateTestGameWithStatus(t, db, "active1", "Active 1", true)
	activeGame2 := CreateTestGameWithStatus(t, db, "active2", "Active 2", true)
	inactiveGame1 := CreateTestGameWithStatus(t, db, "inactive1", "Inactive 1", false)
	inactiveGame2 := CreateTestGameWithStatus(t, db, "inactive2", "Inactive 2", false)

	// Count active games
	var activeCount int64
	err := db.Model(&model.Game{}).Where("is_active = ?", true).Count(&activeCount).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, activeCount, int64(2))

	// Count inactive games
	var inactiveCount int64
	err = db.Model(&model.Game{}).Where("is_active = ?", false).Count(&inactiveCount).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, inactiveCount, int64(2))

	// Verify specific games
	var g1, g2, g3, g4 model.Game
	db.First(&g1, activeGame1.ID)
	db.First(&g2, activeGame2.ID)
	db.First(&g3, inactiveGame1.ID)
	db.First(&g4, inactiveGame2.ID)

	assert.True(t, g1.IsActive)
	assert.True(t, g2.IsActive)
	assert.False(t, g3.IsActive)
	assert.False(t, g4.IsActive)
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestGame_CreateWithLongDescription(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with very long description (text field can handle it)
	longDesc := ""
	for i := 0; i < 1000; i++ {
		longDesc += "A"
	}

	game := &model.Game{
		Key:         fmt.Sprintf("long_desc_%d", time.Now().UnixNano()),
		Name:        "Long Description Game",
		Category:    "test",
		Description: longDesc,
		IsActive:    true,
	}

	err := repo.Create(ctx, game)
	require.NoError(t, err)
	assert.NotZero(t, game.ID)
}

func TestGame_CreateWithSpecialCharactersInName(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with special characters
	game := &model.Game{
		Key:         fmt.Sprintf("special_%d", time.Now().UnixNano()),
		Name:        "Game with 中文 Characters & Symbols!",
		Category:    "test",
		Description: "Test description",
		IsActive:    true,
	}

	err := repo.Create(ctx, game)
	require.NoError(t, err)
	assert.NotZero(t, game.ID)
}

func TestGame_UpdateWithNilFields(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game with all fields populated
	game := &model.Game{
		Key:         fmt.Sprintf("nil_test_%d", time.Now().UnixNano()),
		Name:        "Test Game",
		Category:    "test",
		IconURL:     "https://example.com/icon.png",
		Description: "Test description",
		IsActive:    true,
	}
	require.NoError(t, repo.Create(ctx, game))

	// Update with empty string values (should be trimmed to empty)
	game.IconURL = ""
	game.Description = ""
	err := repo.Update(ctx, game)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.Get(ctx, game.ID)
	require.NoError(t, err)
	assert.Equal(t, "", fetched.IconURL)
	assert.Equal(t, "", fetched.Description)
}

func TestGame_DeleteNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Try to delete non-existent game
	err := repo.Delete(ctx, 99999)
	assert.Error(t, err)
}

func TestGame_UpdateNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Try to update non-existent game
	game := &model.Game{
		Base:     model.Base{ID: 99999},
		Key:      "non_existent",
		Name:     "Non Existent",
		Category: "test",
	}
	err := repo.Update(ctx, game)
	assert.Error(t, err)
}

func TestGame_BatchUpdateStatusWithNonExistentIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create one valid game
	game := CreateTestGameWithKey(t, db, "valid", "Valid Game")

	// Mix valid and invalid IDs
	ids := []uint64{game.ID, 99999, 99998}

	// Batch update should succeed for valid IDs only or fail entirely
	// depending on implementation
	affected, err := repo.BatchUpdateStatus(ctx, ids, false)
	// Current implementation updates what it can
	require.NoError(t, err)
	assert.GreaterOrEqual(t, affected, int64(1))
}

// ============================================================================
// Performance Tests
// ============================================================================

func TestGame_LargeDataset(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create large dataset
	for i := 0; i < 50; i++ {
		game := &model.Game{
			Key:      fmt.Sprintf("perf_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("Performance Game %d", i),
			Category: "test",
			IsActive: i%2 == 0,
		}
		require.NoError(t, repo.Create(ctx, game))
	}

	// Test list performance
	start := time.Now()
	games, total, err := repo.ListPaged(ctx, 1, 25)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, int64(50), total)
	assert.Len(t, games, 25)
	assert.Less(t, duration.Milliseconds(), int64(1000)) // Should complete in < 1s
}

func TestGame_ConcurrentUpdates(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create a game
	game := &model.Game{
		Key:      fmt.Sprintf("concurrent_%d", time.Now().UnixNano()),
		Name:     "Concurrent Test",
		Category: "test",
		IsActive: true,
	}
	require.NoError(t, repo.Create(ctx, game))

	// Update concurrently
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			g := &model.Game{
				Base:     model.Base{ID: game.ID},
				Key:      game.Key,
				Name:     fmt.Sprintf("Updated %d", index),
				Category: "test",
			}
			_ = repo.Update(ctx, g)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify final state (last write wins)
	fetched, _ := repo.Get(ctx, game.ID)
	assert.NotNil(t, fetched)
}

// ============================================================================
// Timestamp Tests
// ============================================================================

func TestGame_Timestamps(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game
	beforeCreate := time.Now()
	game := &model.Game{
		Key:      fmt.Sprintf("timestamp_%d", time.Now().UnixNano()),
		Name:     "Timestamp Test",
		Category: "test",
	}
	require.NoError(t, repo.Create(ctx, game))
	afterCreate := time.Now()

	// Verify CreatedAt
	assert.True(t, game.CreatedAt.After(beforeCreate) || game.CreatedAt.Equal(beforeCreate))
	assert.True(t, game.CreatedAt.Before(afterCreate) || game.CreatedAt.Equal(afterCreate))

	// Update game
	beforeUpdate := time.Now()
	game.Name = "Updated Name"
	err := repo.Update(ctx, game)
	require.NoError(t, err)
	afterUpdate := time.Now()

	// Verify UpdatedAt changed
	fetched, _ := repo.Get(ctx, game.ID)
	assert.True(t, fetched.UpdatedAt.After(beforeUpdate) || fetched.UpdatedAt.Equal(beforeUpdate))
	assert.True(t, fetched.UpdatedAt.Before(afterUpdate) || fetched.UpdatedAt.Equal(afterUpdate))
}

// ============================================================================
// Soft Delete Tests
// ============================================================================

func TestGame_SoftDelete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create game
	game := &model.Game{
		Key:      fmt.Sprintf("soft_del_%d", time.Now().UnixNano()),
		Name:     "Soft Delete Test",
		Category: "test",
	}
	require.NoError(t, repo.Create(ctx, game))

	// Soft delete
	err := repo.Delete(ctx, game.ID)
	require.NoError(t, err)

	// Verify cannot be found with Get
	_, err = repo.Get(ctx, game.ID)
	assert.Error(t, err)

	// Verify still exists in database with deleted_at set
	var deletedAt *time.Time
	err = db.Model(&model.Game{}).
		Select("deleted_at").
		Unscoped().
		Where("id = ?", game.ID).
		Scan(&deletedAt).Error
	require.NoError(t, err)
	assert.NotNil(t, deletedAt)

	// Verify can be found with Unscoped
	var game2 model.Game
	err = db.Unscoped().First(&game2, game.ID).Error
	require.NoError(t, err)
	assert.Equal(t, game.ID, game2.ID)
	assert.NotNil(t, game2.DeletedAt)
}

func TestGame_SoftDeleteBatch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create games
	var ids []uint64
	for i := 0; i < 3; i++ {
		game := &model.Game{
			Key:      fmt.Sprintf("batch_soft_%d_%d", i, time.Now().UnixNano()),
			Name:     fmt.Sprintf("Batch Soft %d", i),
			Category: "test",
		}
		require.NoError(t, repo.Create(ctx, game))
		ids = append(ids, game.ID)
	}

	// Batch delete
	affected, err := repo.BatchDelete(ctx, ids)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)

	// Verify all are soft deleted
	for _, id := range ids {
		// Should not be found with normal Get
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)

		// Should exist in database with deleted_at set
		var deletedAt *time.Time
		err = db.Model(&model.Game{}).
			Select("deleted_at").
			Unscoped().
			Where("id = ?", id).
			Scan(&deletedAt).Error
		require.NoError(t, err)
		assert.NotNil(t, deletedAt)
	}
}
