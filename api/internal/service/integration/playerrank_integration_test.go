// Package integration provides integration tests for services.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/playerrank"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerRankRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_create")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_create")
	rank := CreateTestGameRank(t, db, game, "Diamond", 5, 5000)

	record := &model.PlayerRankRecord{
		PlayerID:       player.ID,
		GameID:         game.ID,
		RankID:         rank.ID,
		Status:         model.PlayerRankStatusPending,
		ScreenshotURLs: "[]",
	}

	err := repo.Create(ctx, record)
	require.NoError(t, err)
	assert.NotZero(t, record.ID)
}

func TestPlayerRankRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_get")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_get")
	rank := CreateTestGameRank(t, db, game, "Gold", 3, 3000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	result, err := repo.Get(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, record.ID, result.ID)
	assert.Equal(t, player.ID, result.PlayerID)
}

func TestPlayerRankRepository_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	_, err := repo.Get(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerRankRepository_GetWithRelations(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_relations")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_relations")
	rank := CreateTestGameRank(t, db, game, "Platinum", 4, 4000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusVerified)

	result, err := repo.GetWithRelations(ctx, record.ID)
	require.NoError(t, err)
	assert.NotNil(t, result.Player)
	assert.NotNil(t, result.Game)
	assert.NotNil(t, result.Rank)
	assert.Equal(t, player.ID, result.Player.ID)
	assert.Equal(t, game.ID, result.Game.ID)
	assert.Equal(t, rank.ID, result.Rank.ID)
}

func TestPlayerRankRepository_GetByPlayerAndGame(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_by_pg")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_by_pg")
	rank := CreateTestGameRank(t, db, game, "Master", 6, 6000)
	CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusVerified)

	result, err := repo.GetByPlayerAndGame(ctx, player.ID, game.ID)
	require.NoError(t, err)
	assert.Equal(t, player.ID, result.PlayerID)
	assert.Equal(t, game.ID, result.GameID)
}

func TestPlayerRankRepository_ListByPlayerID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_list_player")
	player := CreateTestPlayer(t, db, playerUser)
	game1 := CreateTestGame(t, db, "pr_game_list1")
	game2 := CreateTestGame(t, db, "pr_game_list2")
	rank1 := CreateTestGameRank(t, db, game1, "Gold", 3, 3000)
	rank2 := CreateTestGameRank(t, db, game2, "Diamond", 5, 5000)

	CreateTestPlayerRankRecord(t, db, player, game1, rank1, model.PlayerRankStatusVerified)
	CreateTestPlayerRankRecord(t, db, player, game2, rank2, model.PlayerRankStatusPending)

	records, err := repo.ListByPlayerID(ctx, player.ID)
	require.NoError(t, err)
	assert.Len(t, records, 2)
}

func TestPlayerRankRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	// Create multiple records
	for i := 0; i < 5; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_paged_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("pr_game_paged_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, game, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)
	}

	records, total, err := repo.ListPaged(ctx, repository.PlayerRankListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, records, 3)
}

func TestPlayerRankRepository_ListPaged_FilterByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser1 := CreateUniqueTestUser(t, db, "pr_status1")
	player1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "pr_status2")
	player2 := CreateTestPlayer(t, db, playerUser2)
	game := CreateTestGame(t, db, "pr_game_status")
	rank := CreateTestGameRank(t, db, game, "Gold", 3, 3000)

	CreateTestPlayerRankRecord(t, db, player1, game, rank, model.PlayerRankStatusPending)
	CreateTestPlayerRankRecord(t, db, player2, game, rank, model.PlayerRankStatusVerified)

	status := model.PlayerRankStatusPending
	records, total, err := repo.ListPaged(ctx, repository.PlayerRankListOptions{
		Status:   &status,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, model.PlayerRankStatusPending, records[0].Status)
}

func TestPlayerRankRepository_ListPending(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_pending")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_pending")
	rank := CreateTestGameRank(t, db, game, "Silver", 2, 2000)

	CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	records, total, err := repo.ListPending(ctx, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, r := range records {
		assert.Equal(t, model.PlayerRankStatusPending, r.Status)
	}
}

func TestPlayerRankRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_update")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_update")
	rank := CreateTestGameRank(t, db, game, "Bronze", 1, 1000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	// Update
	record.Status = model.PlayerRankStatusVerified
	record.Remark = "Approved"
	err := repo.Update(ctx, record)
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusVerified, updated.Status)
	assert.Equal(t, "Approved", updated.Remark)
}

func TestPlayerRankRepository_UpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pr_admin")
	playerUser := CreateUniqueTestUser(t, db, "pr_status_update")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_status_update")
	rank := CreateTestGameRank(t, db, game, "Gold", 3, 3000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	// Approve
	err := repo.UpdateStatus(ctx, record.ID, model.PlayerRankStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusVerified, updated.Status)
	assert.NotNil(t, updated.VerifiedBy)
	assert.Equal(t, adminUser.ID, *updated.VerifiedBy)
}

func TestPlayerRankRepository_UpdateStatus_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pr_admin_reject")
	playerUser := CreateUniqueTestUser(t, db, "pr_reject")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_reject")
	rank := CreateTestGameRank(t, db, game, "Diamond", 5, 5000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	// Reject
	err := repo.UpdateStatus(ctx, record.ID, model.PlayerRankStatusRejected, &adminUser.ID, "Invalid screenshot")
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusRejected, updated.Status)
	assert.Equal(t, "Invalid screenshot", updated.RejectReason)
}

func TestPlayerRankRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_delete")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_delete")
	rank := CreateTestGameRank(t, db, game, "Silver", 2, 2000)
	record := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	err := repo.Delete(ctx, record.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, record.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerRankRepository_CountByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	// Create records with different statuses
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_count_pending_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("pr_game_count_pending_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, game, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)
	}

	for i := 0; i < 2; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_count_verified_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		game := CreateTestGame(t, db, fmt.Sprintf("pr_game_count_verified_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, game, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusVerified)
	}

	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, counts[model.PlayerRankStatusPending], int64(3))
	assert.GreaterOrEqual(t, counts[model.PlayerRankStatusVerified], int64(2))
}

func TestPlayerRankRepository_GetPendingCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playerrank.NewPlayerRankRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pr_pending_count")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "pr_game_pending_count")
	rank := CreateTestGameRank(t, db, game, "Gold", 3, 3000)
	CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)

	count, err := repo.GetPendingCount(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}
