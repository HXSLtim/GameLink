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
)

// ============================================================================
// Game CRUD Tests
// ============================================================================

func TestGameRepository_Create(t *testing.T) {
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

func TestGameRepository_Get(t *testing.T) {
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

func TestGameRepository_Update(t *testing.T) {
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

func TestGameRepository_Delete(t *testing.T) {
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

func TestGameRepository_List(t *testing.T) {
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

func TestGameRepository_BatchDelete(t *testing.T) {
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
