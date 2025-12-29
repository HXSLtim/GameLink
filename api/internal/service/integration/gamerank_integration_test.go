// Package integration provides service-level integration tests for GameRank and PlayerRank modules.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/game"
	gamerankrepo "gamelink/internal/repository/gamerank"
	"gamelink/internal/repository/player"
	playerrankrepo "gamelink/internal/repository/playerrank"
	gameranksvc "gamelink/internal/service/gamerank"
	playerranksvc "gamelink/internal/service/playerrank"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== GameRank Service Tests ====================

func TestGameRankService_Create_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_create")

	// Execute
	rank, err := service.Create(ctx, gameranksvc.CreateInput{
		GameID:      testGame.ID,
		Name:        "Diamond",
		Level:       5,
		PriceCents:  5000,
		IconURL:     "https://example.com/diamond.png",
		Color:       "#00D9FF",
		Description: "Diamond tier players",
		SortOrder:   5,
		IsActive:    true,
	})

	// Assert
	require.NoError(t, err)
	assert.NotZero(t, rank.ID)
	assert.Equal(t, testGame.ID, rank.GameID)
	assert.Equal(t, "Diamond", rank.Name)
	assert.Equal(t, 5, rank.Level)
	assert.Equal(t, int64(5000), rank.PriceCents)
}

func TestGameRankService_Create_GameNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	// Execute
	rank, err := service.Create(ctx, gameranksvc.CreateInput{
		GameID:     99999,
		Name:       "Test",
		Level:      1,
		PriceCents: 1000,
	})

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrGameNotFound)
	assert.Nil(t, rank)
}

func TestGameRankService_Create_EmptyName(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_empty_name")

	// Execute
	rank, err := service.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "",
		Level:      1,
		PriceCents: 1000,
	})

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrValidation)
	assert.Nil(t, rank)
}

func TestGameRankService_Create_NegativePrice(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_negative_price")

	// Execute
	rank, err := service.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Test",
		Level:      1,
		PriceCents: -100,
	})

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrValidation)
	assert.Nil(t, rank)
}

func TestGameRankService_Get_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_get")
	testRank := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)

	// Execute
	rank, err := service.Get(ctx, testRank.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testRank.ID, rank.ID)
	assert.Equal(t, testGame.ID, rank.GameID)
	assert.NotNil(t, rank.Game)
	assert.Equal(t, "Gold", rank.Name)
}

func TestGameRankService_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	// Execute
	rank, err := service.Get(ctx, 99999)

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrNotFound)
	assert.Nil(t, rank)
}

func TestGameRankService_List_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_list")
	CreateTestGameRank(t, db, testGame, "Bronze", 1, 1000)
	CreateTestGameRank(t, db, testGame, "Silver", 2, 2000)
	CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)

	// Execute
	ranks, err := service.List(ctx)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(ranks), 3)
}

func TestGameRankService_ListByGameID_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame1 := CreateTestGame(t, db, "test_game_list1")
	testGame2 := CreateTestGame(t, db, "test_game_list2")

	CreateTestGameRank(t, db, testGame1, "Bronze", 1, 1000)
	CreateTestGameRank(t, db, testGame1, "Silver", 2, 2000)
	CreateTestGameRank(t, db, testGame2, "Gold", 3, 3000)

	// Execute
	ranks, err := service.ListByGameID(ctx, testGame1.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, ranks, 2)
	// Verify sorted by level
	assert.Equal(t, 1, ranks[0].Level)
	assert.Equal(t, 2, ranks[1].Level)
}

func TestGameRankService_ListPaged_WithFilters(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_paged")
	rank1 := CreateTestGameRank(t, db, testGame, "Diamond", 5, 5000)
	rank2 := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)
	rank2.IsActive = false
	db.Save(rank2)

	// Test filter by keyword
	ranks, pagination, err := service.ListPaged(ctx, repository.GameRankListOptions{
		Keyword:  "Dia",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, pagination.Total)
	assert.Equal(t, "Diamond", ranks[0].Name)

	// Test filter by active
	isActive := true
	ranks, pagination, err = service.ListPaged(ctx, repository.GameRankListOptions{
		GameID:   &testGame.ID,
		IsActive: &isActive,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, pagination.Total)
	assert.Equal(t, rank1.ID, ranks[0].ID)
}

func TestGameRankService_Update_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_update")
	rank := CreateTestGameRank(t, db, testGame, "Bronze", 1, 1000)

	// Execute
	updated, err := service.Update(ctx, rank.ID, gameranksvc.UpdateInput{
		Name:        "Updated Bronze",
		Level:       2,
		PriceCents:  1500,
		IconURL:     "https://example.com/new.png",
		Color:       "#FF0000",
		Description: "Updated description",
		SortOrder:   10,
		IsActive:    false,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, rank.ID, updated.ID)
	assert.Equal(t, "Updated Bronze", updated.Name)
	assert.Equal(t, 2, updated.Level)
	assert.Equal(t, int64(1500), updated.PriceCents)
	assert.False(t, updated.IsActive)
}

func TestGameRankService_Update_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	// Execute
	_, err := service.Update(ctx, 99999, gameranksvc.UpdateInput{
		Name:       "Test",
		Level:      1,
		PriceCents: 1000,
	})

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrNotFound)
}

func TestGameRankService_Update_EmptyName(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_empty_name_update")
	rank := CreateTestGameRank(t, db, testGame, "Test", 1, 1000)

	// Execute
	_, err := service.Update(ctx, rank.ID, gameranksvc.UpdateInput{
		Name:       "",
		Level:      1,
		PriceCents: 1000,
	})

	// Assert
	assert.ErrorIs(t, err, gameranksvc.ErrValidation)
}

func TestGameRankService_Delete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_delete")
	rank := CreateTestGameRank(t, db, testGame, "ToDelete", 1, 1000)

	// Execute
	err := service.Delete(ctx, rank.ID)

	// Assert
	require.NoError(t, err)
	_, err = service.Get(ctx, rank.ID)
	assert.ErrorIs(t, err, gameranksvc.ErrNotFound)
}

func TestGameRankService_BatchDelete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_batch_delete")
	rank1 := CreateTestGameRank(t, db, testGame, "Del1", 1, 1000)
	rank2 := CreateTestGameRank(t, db, testGame, "Del2", 2, 2000)
	rank3 := CreateTestGameRank(t, db, testGame, "Keep", 3, 3000)

	// Execute
	affected, err := service.BatchDelete(ctx, []uint64{rank1.ID, rank2.ID})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	// Verify rank3 still exists
	_, err = service.Get(ctx, rank3.ID)
	assert.NoError(t, err)
}

func TestGameRankService_BatchUpdateStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "test_game_batch_status")
	rank1 := CreateTestGameRank(t, db, testGame, "Status1", 1, 1000)
	rank2 := CreateTestGameRank(t, db, testGame, "Status2", 2, 2000)

	// Execute
	affected, err := service.BatchUpdateStatus(ctx, []uint64{rank1.ID, rank2.ID}, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	// Verify
	updated1, _ := service.Get(ctx, rank1.ID)
	updated2, _ := service.Get(ctx, rank2.ID)
	assert.False(t, updated1.IsActive)
	assert.False(t, updated2.IsActive)
}

// ==================== PlayerRank Service Tests ====================

func TestPlayerRankService_Apply_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_apply_game")
	testRank := CreateTestGameRank(t, db, testGame, "Diamond", 5, 5000)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID:       testPlayer.ID,
		GameID:         testGame.ID,
		RankID:         testRank.ID,
		ScreenshotURLs: `["https://example.com/screenshot.jpg"]`,
		Remark:         "Please verify my rank",
	})

	// Assert
	require.NoError(t, err)
	assert.NotZero(t, record.ID)
	assert.Equal(t, testPlayer.ID, record.PlayerID)
	assert.Equal(t, testGame.ID, record.GameID)
	assert.Equal(t, testRank.ID, record.RankID)
	assert.Equal(t, model.PlayerRankStatusPending, record.Status)
}

func TestPlayerRankService_Apply_PlayerNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	testGame := CreateTestGame(t, db, "pr_apply_player_game")
	rank := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: 99999,
		GameID:   testGame.ID,
		RankID:   rank.ID,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrPlayerNotFound)
	assert.Nil(t, record)
}

func TestPlayerRankService_Apply_GameNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_game_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   99999,
		RankID:   1,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrGameNotFound)
	assert.Nil(t, record)
}

func TestPlayerRankService_Apply_RankNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_rank_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_apply_rank_game")

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame.ID,
		RankID:   99999,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrRankNotFound)
	assert.Nil(t, record)
}

func TestPlayerRankService_Apply_RankNotBelongToGame(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_mismatch_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame1 := CreateTestGame(t, db, "pr_apply_mismatch_game1")
	testGame2 := CreateTestGame(t, db, "pr_apply_mismatch_game2")
	rank := CreateTestGameRank(t, db, testGame1, "Gold", 3, 3000)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute - try to apply with rank from game1 to game2
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame2.ID,
		RankID:   rank.ID,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrValidation)
	assert.Nil(t, record)
}

func TestPlayerRankService_Apply_AlreadyApplied(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_dup_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_apply_dup_game")
	rank := CreateTestGameRank(t, db, testGame, "Silver", 2, 2000)

	// Create existing pending record
	CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute - try to apply again
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame.ID,
		RankID:   rank.ID,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrAlreadyApplied)
	assert.Nil(t, record)
}

func TestPlayerRankService_Apply_AfterRejected(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_apply_rejected_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_apply_rejected_game")
	rank := CreateTestGameRank(t, db, testGame, "Platinum", 4, 4000)

	// Create existing rejected record
	CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusRejected)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute - should be able to reapply after rejection
	record, err := service.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame.ID,
		RankID:   rank.ID,
	})

	// Assert
	require.NoError(t, err)
	assert.NotZero(t, record.ID)
	assert.Equal(t, model.PlayerRankStatusPending, record.Status)
}

func TestPlayerRankService_Verify_Approve(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	adminUser := CreateUniqueTestUser(t, db, "pr_verify_admin")
	playerUser := CreateUniqueTestUser(t, db, "pr_verify_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_verify_game")
	rank := CreateTestGameRank(t, db, testGame, "Master", 6, 6000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	updated, err := service.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:     record.ID,
		Status:       model.PlayerRankStatusVerified,
		VerifiedBy:   adminUser.ID,
		RejectReason: "",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusVerified, updated.Status)
	assert.NotNil(t, updated.VerifiedBy)
	assert.Equal(t, adminUser.ID, *updated.VerifiedBy)
	assert.NotNil(t, updated.VerifiedAt)

	// Verify player rank was updated
	updatedPlayer, _ := playersRepo.Get(ctx, testPlayer.ID)
	assert.Equal(t, rank.Name, updatedPlayer.Rank)
	assert.Equal(t, rank.PriceCents, updatedPlayer.HourlyRateCents)
}

func TestPlayerRankService_Verify_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	adminUser := CreateUniqueTestUser(t, db, "pr_reject_admin")
	playerUser := CreateUniqueTestUser(t, db, "pr_reject_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_reject_game")
	rank := CreateTestGameRank(t, db, testGame, "Diamond", 5, 5000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	updated, err := service.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:     record.ID,
		Status:       model.PlayerRankStatusRejected,
		VerifiedBy:   adminUser.ID,
		RejectReason: "Screenshot unclear",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusRejected, updated.Status)
	assert.Equal(t, "Screenshot unclear", updated.RejectReason)
}

func TestPlayerRankService_Verify_Revoke(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	adminUser := CreateUniqueTestUser(t, db, "pr_revoke_admin")
	playerUser := CreateUniqueTestUser(t, db, "pr_revoke_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_revoke_game")
	rank := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusVerified)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	updated, err := service.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:     record.ID,
		Status:       model.PlayerRankStatusRevoked,
		VerifiedBy:   adminUser.ID,
		RejectReason: "Rank verification failed",
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusRevoked, updated.Status)
	assert.Equal(t, "Rank verification failed", updated.RejectReason)
}

func TestPlayerRankService_Verify_InvalidStatusTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	adminUser := CreateUniqueTestUser(t, db, "pr_invalid_admin")
	playerUser := CreateUniqueTestUser(t, db, "pr_invalid_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_invalid_game")
	rank := CreateTestGameRank(t, db, testGame, "Bronze", 1, 1000)

	// Create rejected record and try to verify it again
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusRejected)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute - try to change rejected to verified (invalid transition)
	_, err := service.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   record.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})

	// Assert
	assert.ErrorIs(t, err, playerranksvc.ErrInvalidStatus)
}

func TestPlayerRankService_Get_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_get_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_get_game")
	rank := CreateTestGameRank(t, db, testGame, "Silver", 2, 2000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusVerified)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	result, err := service.Get(ctx, record.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, record.ID, result.ID)
	assert.NotNil(t, result.Player)
	assert.NotNil(t, result.Game)
	assert.NotNil(t, result.Rank)
}

func TestPlayerRankService_ListByPlayerID_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_list_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame1 := CreateTestGame(t, db, "pr_list_game1")
	testGame2 := CreateTestGame(t, db, "pr_list_game2")
	rank1 := CreateTestGameRank(t, db, testGame1, "Gold", 3, 3000)
	rank2 := CreateTestGameRank(t, db, testGame2, "Diamond", 5, 5000)

	CreateTestPlayerRankRecord(t, db, testPlayer, testGame1, rank1, model.PlayerRankStatusVerified)
	CreateTestPlayerRankRecord(t, db, testPlayer, testGame2, rank2, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	records, err := service.ListByPlayerID(ctx, testPlayer.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, records, 2)
}

func TestPlayerRankService_ListPending_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_pending_%d", i))
		testPlayer := CreateTestPlayer(t, db, playerUser)
		testGame := CreateTestGame(t, db, fmt.Sprintf("pr_pending_game_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, testGame, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)
	}

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	records, pagination, err := service.ListPending(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, pagination.Total, int64(3))
	for _, r := range records {
		assert.Equal(t, model.PlayerRankStatusPending, r.Status)
	}
}

func TestPlayerRankService_ListPaged_WithFilters(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser1 := CreateUniqueTestUser(t, db, "pr_paged_user1")
	player1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "pr_paged_user2")
	player2 := CreateTestPlayer(t, db, playerUser2)
	testGame := CreateTestGame(t, db, "pr_paged_game")
	rank := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)

	CreateTestPlayerRankRecord(t, db, player1, testGame, rank, model.PlayerRankStatusPending)
	CreateTestPlayerRankRecord(t, db, player2, testGame, rank, model.PlayerRankStatusVerified)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Test filter by status
	status := model.PlayerRankStatusVerified
	records, pagination, err := service.ListPaged(ctx, repository.PlayerRankListOptions{
		Status:   &status,
		Page:     1,
		PageSize: 10,
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, pagination.Total)
	assert.Equal(t, model.PlayerRankStatusVerified, records[0].Status)
}

func TestPlayerRankService_Delete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_delete_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_delete_game")
	rank := CreateTestGameRank(t, db, testGame, "Bronze", 1, 1000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	err := service.Delete(ctx, record.ID)

	// Assert
	require.NoError(t, err)
	_, err = service.Get(ctx, record.ID)
	assert.ErrorIs(t, err, playerranksvc.ErrNotFound)
}

func TestPlayerRankService_GetStats_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_stats_pending_%d", i))
		testPlayer := CreateTestPlayer(t, db, playerUser)
		testGame := CreateTestGame(t, db, fmt.Sprintf("pr_stats_game_pending_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, testGame, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)
	}

	for i := 0; i < 2; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pr_stats_verified_%d", i))
		testPlayer := CreateTestPlayer(t, db, playerUser)
		testGame := CreateTestGame(t, db, fmt.Sprintf("pr_stats_game_verified_%d_%d", i, time.Now().UnixNano()))
		rank := CreateTestGameRank(t, db, testGame, "Rank", i+1, int64((i+1)*1000))
		CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusVerified)
	}

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	stats, err := service.GetStats(ctx)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats[model.PlayerRankStatusPending], int64(3))
	assert.GreaterOrEqual(t, stats[model.PlayerRankStatusVerified], int64(2))
}

func TestPlayerRankService_GetPendingCount_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerUser := CreateUniqueTestUser(t, db, "pr_pending_count_user")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "pr_pending_count_game")
	rank := CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)
	CreateTestPlayerRankRecord(t, db, testPlayer, testGame, rank, model.PlayerRankStatusPending)

	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)
	ranksRepo := gamerankrepo.NewGameRankRepository(db)
	playersRepo := player.NewPlayerRepository(db)
	gamesRepo := game.NewGameRepository(db)
	service := playerranksvc.NewPlayerRankService(recordsRepo, ranksRepo, playersRepo, gamesRepo)

	// Execute
	count, err := service.GetPendingCount(ctx)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

// ==================== GameRank and PlayerRank Integration Tests ====================

func TestGameRank_PlayerRank_FullWorkflow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)

	gameRankService := gameranksvc.NewGameRankService(rankRepo, gameRepo)
	playerRankService := playerranksvc.NewPlayerRankService(recordsRepo, rankRepo, playerRepo, gameRepo)

	// Step 1: Admin creates game
	testGame := &model.Game{
		Key:      "lol",
		Name:     "League of Legends",
		Category: "moba",
		IsActive: true,
	}
	err := gameRepo.Create(ctx, testGame)
	require.NoError(t, err)

	// Step 2: Admin creates game ranks
	_, err = gameRankService.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Bronze",
		Level:      1,
		PriceCents: 2000,
		IsActive:   true,
	})
	require.NoError(t, err)

	_, err = gameRankService.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Silver",
		Level:      2,
		PriceCents: 3000,
		IsActive:   true,
	})
	require.NoError(t, err)

	_, err = gameRankService.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Gold",
		Level:      3,
		PriceCents: 5000,
		IsActive:   true,
	})
	require.NoError(t, err)

	// Step 3: Player registers
	playerUser := CreateUniqueTestUser(t, db, "workflow_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Step 4: Player applies for rank certification
	ranks, _ := gameRankService.ListByGameID(ctx, testGame.ID)
	require.Len(t, ranks, 3)

	record, err := playerRankService.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID:       testPlayer.ID,
		GameID:         testGame.ID,
		RankID:         ranks[2].ID, // Apply for Gold
		ScreenshotURLs: `["https://example.com/gold_rank.jpg"]`,
		Remark:         "I reached Gold rank",
	})
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusPending, record.Status)

	// Step 5: Admin approves the certification
	adminUser := CreateUniqueTestUser(t, db, "workflow_admin")
	approved, err := playerRankService.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   record.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusVerified, approved.Status)

	// Step 6: Verify player's rank and hourly rate updated
	updatedPlayer, err := playerRepo.Get(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, "Gold", updatedPlayer.Rank)
	assert.Equal(t, int64(5000), updatedPlayer.HourlyRateCents)
}

func TestGameRank_PriceAffectsPlayerHourlyRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)

	gameRankService := gameranksvc.NewGameRankService(rankRepo, gameRepo)
	playerRankService := playerranksvc.NewPlayerRankService(recordsRepo, rankRepo, playerRepo, gameRepo)

	// Create game with different rank prices
	testGame := CreateTestGame(t, db, "price_test_game")

	bronzeRank, _ := gameRankService.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Bronze",
		Level:      1,
		PriceCents: 2000,
		IsActive:   true,
	})

	goldRank, _ := gameRankService.Create(ctx, gameranksvc.CreateInput{
		GameID:     testGame.ID,
		Name:       "Gold",
		Level:      3,
		PriceCents: 6000,
		IsActive:   true,
	})

	// Player applies for Bronze
	playerUser := CreateUniqueTestUser(t, db, "price_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	record, err := playerRankService.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID:       testPlayer.ID,
		GameID:         testGame.ID,
		RankID:         bronzeRank.ID,
		ScreenshotURLs: `["https://example.com/bronze.jpg"]`,
	})
	require.NoError(t, err)

	// Admin approves
	adminUser := CreateUniqueTestUser(t, db, "price_admin")
	_, err = playerRankService.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   record.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})
	require.NoError(t, err)

	// Verify player's hourly rate matches Bronze rank
	updatedPlayer, err := playerRepo.Get(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2000), updatedPlayer.HourlyRateCents)

	// Player upgrades to Gold
	newRecord, err := playerRankService.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID:       testPlayer.ID,
		GameID:         testGame.ID,
		RankID:         goldRank.ID,
		ScreenshotURLs: `["https://example.com/gold.jpg"]`,
	})
	require.NoError(t, err)

	_, err = playerRankService.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   newRecord.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})
	require.NoError(t, err)

	// Verify player's hourly rate updated to Gold rank
	updatedPlayer, err = playerRepo.Get(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(6000), updatedPlayer.HourlyRateCents)
	assert.Equal(t, "Gold", updatedPlayer.Rank)
}

func TestGameRank_ActiveStatusFilter(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "active_test_game")

	// Create active and inactive ranks
	activeRank := CreateTestGameRank(t, db, testGame, "Active Rank", 1, 3000)
	inactiveRank := CreateTestGameRank(t, db, testGame, "Inactive Rank", 2, 4000)
	inactiveRank.IsActive = false
	db.Save(inactiveRank)

	// List all ranks
	allRanks, err := service.ListByGameID(ctx, testGame.ID)
	require.NoError(t, err)
	assert.Len(t, allRanks, 2)

	// List only active ranks
	isActive := true
	ranks, pagination, err := service.ListPaged(ctx, repository.GameRankListOptions{
		GameID:   &testGame.ID,
		IsActive: &isActive,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, pagination.Total)
	assert.Equal(t, activeRank.ID, ranks[0].ID)
}

func TestPlayerRank_MultipleGamesOnePlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	recordsRepo := playerrankrepo.NewPlayerRankRepository(db)

	playerRankService := playerranksvc.NewPlayerRankService(recordsRepo, rankRepo, playerRepo, gameRepo)

	// Create multiple games
	testGame1 := CreateTestGame(t, db, "multi_game1")
	testGame2 := CreateTestGame(t, db, "multi_game2")

	rank1 := CreateTestGameRank(t, db, testGame1, "Diamond", 5, 5000)
	rank2 := CreateTestGameRank(t, db, testGame2, "Gold", 3, 3000)

	// Player applies for ranks in both games
	playerUser := CreateUniqueTestUser(t, db, "multi_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	record1, err := playerRankService.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame1.ID,
		RankID:   rank1.ID,
	})
	require.NoError(t, err)

	record2, err := playerRankService.Apply(ctx, playerranksvc.ApplyInput{
		PlayerID: testPlayer.ID,
		GameID:   testGame2.ID,
		RankID:   rank2.ID,
	})
	require.NoError(t, err)

	// Approve both
	adminUser := CreateUniqueTestUser(t, db, "multi_admin")
	_, err = playerRankService.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   record1.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})
	require.NoError(t, err)

	_, err = playerRankService.Verify(ctx, playerranksvc.VerifyInput{
		RecordID:   record2.ID,
		Status:     model.PlayerRankStatusVerified,
		VerifiedBy: adminUser.ID,
	})
	require.NoError(t, err)

	// Verify player has certifications for both games
	records, err := playerRankService.ListByPlayerID(ctx, testPlayer.ID)
	require.NoError(t, err)
	assert.Len(t, records, 2)
}

func TestGameRank_RankLevelOrdering(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	gameRepo := game.NewGameRepository(db)
	rankRepo := gamerankrepo.NewGameRankRepository(db)
	service := gameranksvc.NewGameRankService(rankRepo, gameRepo)

	testGame := CreateTestGame(t, db, "order_test_game")

	// Create ranks in random order
	CreateTestGameRank(t, db, testGame, "Diamond", 5, 5000)
	CreateTestGameRank(t, db, testGame, "Bronze", 1, 1000)
	CreateTestGameRank(t, db, testGame, "Gold", 3, 3000)
	CreateTestGameRank(t, db, testGame, "Silver", 2, 2000)

	// List should be ordered by level
	ranks, err := service.ListByGameID(ctx, testGame.ID)
	require.NoError(t, err)
	assert.Len(t, ranks, 4)

	// Verify ordering
	assert.Equal(t, 1, ranks[0].Level)
	assert.Equal(t, "Bronze", ranks[0].Name)
	assert.Equal(t, 2, ranks[1].Level)
	assert.Equal(t, "Silver", ranks[1].Name)
	assert.Equal(t, 3, ranks[2].Level)
	assert.Equal(t, "Gold", ranks[2].Name)
	assert.Equal(t, 5, ranks[3].Level)
	assert.Equal(t, "Diamond", ranks[3].Name)
}
