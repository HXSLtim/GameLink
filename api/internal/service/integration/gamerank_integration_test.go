// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/gamerank"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameRankRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_create")

	rank := &model.GameRank{
		GameID:     game.ID,
		Name:       "Diamond",
		Level:      5,
		PriceCents: 5000,
		IsActive:   true,
	}

	err := repo.Create(ctx, rank)
	require.NoError(t, err)
	assert.NotZero(t, rank.ID)
}

func TestGameRankRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_get")
	rank := CreateTestGameRank(t, db, game, "Gold", 3, 3000)

	result, err := repo.Get(ctx, rank.ID)
	require.NoError(t, err)
	assert.Equal(t, rank.ID, result.ID)
	assert.Equal(t, "Gold", result.Name)
	assert.Equal(t, 3, result.Level)
}

func TestGameRankRepository_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)

	_, err := repo.Get(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGameRankRepository_GetWithGame(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_with")
	rank := CreateTestGameRank(t, db, game, "Platinum", 4, 4000)

	result, err := repo.GetWithGame(ctx, rank.ID)
	require.NoError(t, err)
	assert.Equal(t, rank.ID, result.ID)
	assert.NotNil(t, result.Game)
	assert.Equal(t, game.ID, result.Game.ID)
}

func TestGameRankRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_list")

	// Create multiple ranks
	CreateTestGameRank(t, db, game, "Bronze", 1, 1000)
	CreateTestGameRank(t, db, game, "Silver", 2, 2000)
	CreateTestGameRank(t, db, game, "Gold", 3, 3000)

	ranks, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(ranks), 3)
}

func TestGameRankRepository_ListByGameID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game1 := CreateTestGame(t, db, "rank_game_list1")
	game2 := CreateTestGame(t, db, "rank_game_list2")

	// Create ranks for game1
	CreateTestGameRank(t, db, game1, "Bronze", 1, 1000)
	CreateTestGameRank(t, db, game1, "Silver", 2, 2000)

	// Create ranks for game2
	CreateTestGameRank(t, db, game2, "Gold", 3, 3000)

	ranks, err := repo.ListByGameID(ctx, game1.ID)
	require.NoError(t, err)
	assert.Len(t, ranks, 2)

	// Verify sorted by level
	assert.Equal(t, 1, ranks[0].Level)
	assert.Equal(t, 2, ranks[1].Level)
}

func TestGameRankRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_paged")

	// Create multiple ranks
	for i := 1; i <= 5; i++ {
		CreateTestGameRank(t, db, game, "Rank"+string(rune('A'+i-1)), i, int64(i*1000))
	}

	// Test pagination
	ranks, total, err := repo.ListPaged(ctx, repository.GameRankListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, ranks, 3)

	// Test filter by game
	ranks, total, err = repo.ListPaged(ctx, repository.GameRankListOptions{
		GameID:   &game.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
}

func TestGameRankRepository_ListPaged_FilterByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_keyword")

	CreateTestGameRank(t, db, game, "Diamond", 5, 5000)
	CreateTestGameRank(t, db, game, "Master", 6, 6000)

	ranks, total, err := repo.ListPaged(ctx, repository.GameRankListOptions{
		Keyword:  "Dia",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Diamond", ranks[0].Name)
}

func TestGameRankRepository_ListPaged_FilterByActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_active")

	rank1 := CreateTestGameRank(t, db, game, "Active", 1, 1000)
	rank2 := CreateTestGameRank(t, db, game, "Inactive", 2, 2000)

	// Deactivate rank2
	rank2.IsActive = false
	db.Save(rank2)

	isActive := true
	ranks, total, err := repo.ListPaged(ctx, repository.GameRankListOptions{
		GameID:   &game.ID,
		IsActive: &isActive,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, rank1.ID, ranks[0].ID)
}

func TestGameRankRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_update")
	rank := CreateTestGameRank(t, db, game, "Bronze", 1, 1000)

	// Update
	rank.Name = "Updated Bronze"
	rank.PriceCents = 1500
	err := repo.Update(ctx, rank)
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, rank.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Bronze", updated.Name)
	assert.Equal(t, int64(1500), updated.PriceCents)
}

func TestGameRankRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_delete")
	rank := CreateTestGameRank(t, db, game, "ToDelete", 1, 1000)

	err := repo.Delete(ctx, rank.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, rank.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGameRankRepository_BatchDelete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_batch_del")

	rank1 := CreateTestGameRank(t, db, game, "Del1", 1, 1000)
	rank2 := CreateTestGameRank(t, db, game, "Del2", 2, 2000)
	rank3 := CreateTestGameRank(t, db, game, "Keep", 3, 3000)

	affected, err := repo.BatchDelete(ctx, []uint64{rank1.ID, rank2.ID})
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	// Verify rank3 still exists
	_, err = repo.Get(ctx, rank3.ID)
	require.NoError(t, err)
}

func TestGameRankRepository_BatchUpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := gamerank.NewGameRankRepository(db)
	game := CreateTestGame(t, db, "rank_game_batch_status")

	rank1 := CreateTestGameRank(t, db, game, "Status1", 1, 1000)
	rank2 := CreateTestGameRank(t, db, game, "Status2", 2, 2000)

	// Deactivate both
	affected, err := repo.BatchUpdateStatus(ctx, []uint64{rank1.ID, rank2.ID}, false)
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	// Verify
	updated1, _ := repo.Get(ctx, rank1.ID)
	updated2, _ := repo.Get(ctx, rank2.ID)
	assert.False(t, updated1.IsActive)
	assert.False(t, updated2.IsActive)
}
