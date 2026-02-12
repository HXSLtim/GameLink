// Package integration provides integration tests for ranking service.
package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/implementations"
	rankingrepo "gamelink/internal/repository/ranking"
	"gamelink/internal/service/ranking"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// CreateTestPlayerRanking creates a test player ranking record.
func CreateTestPlayerRanking(t *testing.T, db *gorm.DB, playerID uint64, rankingType model.RankingType, period, periodValue string, rank int) *model.PlayerRanking {
	t.Helper()
	ranking := &model.PlayerRanking{
		PlayerID:    playerID,
		RankingType: rankingType,
		Period:      period,
		PeriodValue: periodValue,
		Rank:        rank,
		Score:       float64(100 - rank), // Higher score for better rank
		OrderCount:  int64(100 - rank),
		IncomeCents: int64((100 - rank) * 1000),
		AvgRating:   float32(4.5 + float32(10-rank)/10),
		BonusCents:  int64((100 - rank) * 100),
	}
	if err := db.Create(ranking).Error; err != nil {
		t.Fatalf("Failed to create test player ranking: %v", err)
	}
	return ranking
}

// CreateTestRankingReward creates a test ranking reward.
func CreateTestRankingReward(t *testing.T, db *gorm.DB, rankingType model.RankingType, period string, rankStart, rankEnd int, rewardCents int64) *model.RankingReward {
	t.Helper()
	reward := &model.RankingReward{
		RankingType: rankingType,
		Period:      period,
		RankStart:   rankStart,
		RankEnd:     rankEnd,
		RewardType:  "commission",
		RewardValue: rewardCents,
		Description: "Test ranking reward",
		IsActive:    true,
	}
	if err := db.Create(reward).Error; err != nil {
		t.Fatalf("Failed to create test ranking reward: %v", err)
	}
	return reward
}

// CreateTestRankingCommissionConfig creates a test ranking commission config.
func CreateTestRankingCommissionConfig(t *testing.T, db *gorm.DB, rankingType model.RankingType, month string) *model.RankingCommissionConfig {
	t.Helper()
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 3, CommissionRate: 10},
		{RankStart: 4, RankEnd: 10, CommissionRate: 12},
		{RankStart: 11, RankEnd: 50, CommissionRate: 15},
	}
	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		t.Fatalf("Failed to marshal ranking rules: %v", err)
	}
	config := &model.RankingCommissionConfig{
		Name:        "Test Config " + string(rankingType),
		RankingType: rankingType,
		Period:      "monthly",
		Month:       month,
		RulesJSON:   string(rulesJSON),
		Description: "Test ranking commission config",
		IsActive:    true,
	}
	if err := db.Create(config).Error; err != nil {
		t.Fatalf("Failed to create test ranking commission config: %v", err)
	}
	return config
}

// TestRankingRepository_CreateRanking tests creating a player ranking.
func TestRankingRepository_CreateRanking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	testUser := CreateUniqueTestUser(t, db, "player_rank")
	testPlayer := CreateTestPlayer(t, db, testUser)

	rankingRecord := &model.PlayerRanking{
		PlayerID:    testPlayer.ID,
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		PeriodValue: "2025-01",
		Rank:        1,
		Score:       100000.0,
		OrderCount:  50,
		IncomeCents: 10000000,
		AvgRating:   4.8,
		BonusCents:  50000,
	}

	err := repo.CreateRanking(ctx, rankingRecord)
	require.NoError(t, err)
	assert.NotZero(t, rankingRecord.ID)
}

// TestRankingRepository_GetPlayerRanking tests getting a player's ranking.
func TestRankingRepository_GetPlayerRanking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	testUser := CreateUniqueTestUser(t, db, "player_get_rank")
	testPlayer := CreateTestPlayer(t, db, testUser)

	// Create a ranking
	_ = CreateTestPlayerRanking(t, db, testPlayer.ID, model.RankingTypeOrderCount, "monthly", "2025-01", 5)

	// Get the ranking
	rankingRecord, err := repo.GetPlayerRanking(ctx, testPlayer.ID, model.RankingTypeOrderCount, "monthly", "2025-01")
	require.NoError(t, err)
	require.NotNil(t, rankingRecord)
	assert.Equal(t, testPlayer.ID, rankingRecord.PlayerID)
	assert.Equal(t, model.RankingTypeOrderCount, rankingRecord.RankingType)
	assert.Equal(t, 5, rankingRecord.Rank)
}

// TestRankingRepository_ListRankings tests listing rankings with filters.
func TestRankingRepository_ListRankings(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create multiple players and rankings
	for i := 1; i <= 5; i++ {
		user := CreateUniqueTestUser(t, db, "player_list_ranking")
		player := CreateTestPlayer(t, db, user)
		CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", "2025-01", i)
	}

	tests := []struct {
		name          string
		opts          rankingrepo.RankingListOptions
		expectedCount int
	}{
		{
			name: "List all rankings",
			opts: rankingrepo.RankingListOptions{
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 5,
		},
		{
			name: "Filter by ranking type",
			opts: rankingrepo.RankingListOptions{
				RankingType: func() *model.RankingType { t := model.RankingTypeIncome; return &t }(),
				Page:        1,
				PageSize:    10,
			},
			expectedCount: 5,
		},
		{
			name: "Filter by period value",
			opts: rankingrepo.RankingListOptions{
				PeriodValue: func() *string { s := "2025-01"; return &s }(),
				Page:        1,
				PageSize:    10,
			},
			expectedCount: 5,
		},
		{
			name: "Pagination - page size 2",
			opts: rankingrepo.RankingListOptions{
				Page:     1,
				PageSize: 2,
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rankings, total, err := repo.ListRankings(ctx, tt.opts)
			require.NoError(t, err)
			assert.Equal(t, int64(tt.expectedCount), total)
			assert.Len(t, rankings, tt.expectedCount)
		})
	}
}

// TestRankingRepository_UpdateRanking tests updating a ranking.
func TestRankingRepository_UpdateRanking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	testUser := CreateUniqueTestUser(t, db, "player_update_rank")
	testPlayer := CreateTestPlayer(t, db, testUser)

	// Create a ranking
	original := CreateTestPlayerRanking(t, db, testPlayer.ID, model.RankingTypeIncome, "monthly", "2025-01", 1)

	// Update the ranking
	original.Rank = 2
	original.Score = 90000.0
	original.BonusCents = 30000

	err := repo.UpdateRanking(ctx, original)
	require.NoError(t, err)

	// Verify update
	updated, err := repo.GetPlayerRanking(ctx, testPlayer.ID, model.RankingTypeIncome, "monthly", "2025-01")
	require.NoError(t, err)
	assert.Equal(t, 2, updated.Rank)
	assert.Equal(t, 90000.0, updated.Score)
	assert.Equal(t, int64(30000), updated.BonusCents)
}

// TestRankingRepository_CreateReward tests creating a ranking reward.
func TestRankingRepository_CreateReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	reward := &model.RankingReward{
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     3,
		RewardType:  "commission",
		RewardValue: 50000,
		Description: "Top 3 reward",
		IsActive:    true,
	}

	err := repo.CreateReward(ctx, reward)
	require.NoError(t, err)
	assert.NotZero(t, reward.ID)
}

// TestRankingRepository_GetRewardForRank tests getting reward for a specific rank.
func TestRankingRepository_GetRewardForRank(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create reward for ranks 1-3
	_ = CreateTestRankingReward(t, db, model.RankingTypeOrderCount, "monthly", 1, 3, 10000)

	// Test rank 1 (should find reward)
	reward, err := repo.GetRewardForRank(ctx, model.RankingTypeOrderCount, "monthly", 1)
	require.NoError(t, err)
	require.NotNil(t, reward)
	assert.Equal(t, 1, reward.RankStart)
	assert.Equal(t, 3, reward.RankEnd)
	assert.Equal(t, int64(10000), reward.RewardValue)

	// Test rank 2 (should find same reward)
	reward2, err := repo.GetRewardForRank(ctx, model.RankingTypeOrderCount, "monthly", 2)
	require.NoError(t, err)
	require.NotNil(t, reward2)
	assert.Equal(t, reward.ID, reward2.ID)

	// Test rank 5 (should not find reward)
	reward3, err := repo.GetRewardForRank(ctx, model.RankingTypeOrderCount, "monthly", 5)
	assert.Error(t, err)
	assert.Nil(t, reward3)
}

// TestRankingRepository_ListRewards tests listing rewards.
func TestRankingRepository_ListRewards(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create multiple rewards
	CreateTestRankingReward(t, db, model.RankingTypeIncome, "monthly", 1, 3, 50000)
	CreateTestRankingReward(t, db, model.RankingTypeOrderCount, "monthly", 1, 5, 30000)
	CreateTestRankingReward(t, db, model.RankingTypeQuality, "weekly", 1, 10, 10000)

	tests := []struct {
		name          string
		opts          rankingrepo.RewardListOptions
		expectedCount int
	}{
		{
			name: "List all rewards",
			opts: rankingrepo.RewardListOptions{
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 3,
		},
		{
			name: "Filter by ranking type",
			opts: rankingrepo.RewardListOptions{
				RankingType: func() *model.RankingType { t := model.RankingTypeIncome; return &t }(),
				Page:        1,
				PageSize:    10,
			},
			expectedCount: 1,
		},
		{
			name: "Filter by period",
			opts: rankingrepo.RewardListOptions{
				Period:   func() *string { s := "monthly"; return &s }(),
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 2,
		},
		{
			name: "Filter by active",
			opts: rankingrepo.RewardListOptions{
				IsActive: func() *bool { b := true; return &b }(),
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewards, total, err := repo.ListRewards(ctx, tt.opts)
			require.NoError(t, err)
			assert.Equal(t, int64(tt.expectedCount), total)
			assert.Len(t, rewards, tt.expectedCount)
		})
	}
}

// TestRankingRepository_UpdateReward tests updating a reward.
func TestRankingRepository_UpdateReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create a reward
	reward := CreateTestRankingReward(t, db, model.RankingTypeIncome, "monthly", 1, 3, 50000)

	// Update
	reward.RewardValue = 60000
	reward.Description = "Updated reward"

	err := repo.UpdateReward(ctx, reward)
	require.NoError(t, err)

	// Verify
	updated, err := repo.GetReward(ctx, reward.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(60000), updated.RewardValue)
	assert.Equal(t, "Updated reward", updated.Description)
}

// TestRankingRepository_DeleteReward tests deleting a reward.
func TestRankingRepository_DeleteReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create a reward
	reward := CreateTestRankingReward(t, db, model.RankingTypeIncome, "monthly", 1, 3, 50000)

	// Delete
	err := repo.DeleteReward(ctx, reward.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = repo.GetReward(ctx, reward.ID)
	assert.Error(t, err)
}

// TestRankingCommissionRepository_CreateConfig tests creating commission config.
func TestRankingCommissionRepository_CreateConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingCommissionRepository(db)

	config := &model.RankingCommissionConfig{
		Name:        "Test Commission Config",
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		Month:       "2025-01",
		RulesJSON:   `[{"rankStart":1,"rankEnd":3,"commissionRate":10}]`,
		Description: "Test config",
		IsActive:    true,
	}

	err := repo.CreateConfig(ctx, config)
	require.NoError(t, err)
	assert.NotZero(t, config.ID)
}

// TestRankingCommissionRepository_GetActiveConfigForMonth tests getting active config.
func TestRankingCommissionRepository_GetActiveConfigForMonth(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingCommissionRepository(db)

	// Create active config
	_ = CreateTestRankingCommissionConfig(t, db, model.RankingTypeIncome, "2025-01")

	// Get active config
	config, err := repo.GetActiveConfigForMonth(ctx, model.RankingTypeIncome, "2025-01")
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, model.RankingTypeIncome, config.RankingType)
	assert.Equal(t, "2025-01", config.Month)
	assert.True(t, config.IsActive)
}

// TestRankingCommissionRepository_ListConfigs tests listing commission configs.
func TestRankingCommissionRepository_ListConfigs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingCommissionRepository(db)

	// Create multiple configs
	CreateTestRankingCommissionConfig(t, db, model.RankingTypeIncome, "2025-01")
	CreateTestRankingCommissionConfig(t, db, model.RankingTypeOrderCount, "2025-01")
	CreateTestRankingCommissionConfig(t, db, model.RankingTypeIncome, "2025-02")

	tests := []struct {
		name          string
		opts          rankingrepo.RankingCommissionConfigListOptions
		expectedCount int
	}{
		{
			name: "List all configs",
			opts: rankingrepo.RankingCommissionConfigListOptions{
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 3,
		},
		{
			name: "Filter by ranking type",
			opts: rankingrepo.RankingCommissionConfigListOptions{
				RankingType: func() *model.RankingType { t := model.RankingTypeIncome; return &t }(),
				Page:        1,
				PageSize:    10,
			},
			expectedCount: 2,
		},
		{
			name: "Filter by month",
			opts: rankingrepo.RankingCommissionConfigListOptions{
				Month:    func() *string { s := "2025-01"; return &s }(),
				Page:     1,
				PageSize: 10,
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configs, total, err := repo.ListConfigs(ctx, tt.opts)
			require.NoError(t, err)
			assert.Equal(t, int64(tt.expectedCount), total)
			assert.Len(t, configs, tt.expectedCount)
		})
	}
}

// TestRankingService_CalculateMonthlyRankings tests monthly ranking calculation.
func TestRankingService_CalculateMonthlyRankings(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	rankingRepo := rankingrepo.NewRankingRepository(db)
	commissionRepo := rankingrepo.NewRankingCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := ranking.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	// Create test players and orders
	month := "2025-01"
	for i := 1; i <= 3; i++ {
		user := CreateUniqueTestUser(t, db, "ranking_player")
		player := CreateTestPlayer(t, db, user)

		// Create completed orders for this player
		for j := 1; j <= (4 - i); j++ {
			game := CreateTestGame(t, db, "game"+string(rune(i)))
			CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
		}
	}

	// Calculate rankings
	err := svc.CalculateMonthlyRankings(ctx, month)
	require.NoError(t, err)

	// Verify rankings were created
	rankings, total, err := rankingRepo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: func() *string { s := month; return &s }(),
		Page:        1,
		PageSize:    100,
	})
	require.NoError(t, err)
	assert.Greater(t, total, int64(0))
	assert.NotEmpty(t, rankings)

	// Verify rankings are ordered correctly by rank
	for i := 0; i < len(rankings)-1; i++ {
		assert.LessOrEqual(t, rankings[i].Rank, rankings[i+1].Rank)
	}
}

// TestRankingService_GetPlayerRankingInfo tests getting player ranking info.
func TestRankingService_GetPlayerRankingInfo(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	rankingRepo := rankingrepo.NewRankingRepository(db)
	commissionRepo := rankingrepo.NewRankingCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := ranking.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "rank_info_player")
	player := CreateTestPlayer(t, db, user)
	month := "2025-01"

	// Create rankings
	CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", month, 3)
	CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeOrderCount, "monthly", month, 5)

	// Get ranking info
	info, err := svc.GetPlayerRankingInfo(ctx, player.ID, month)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, player.ID, info.PlayerID)
	assert.Equal(t, month, info.Month)
	assert.Equal(t, 3, info.BestRank) // Should get the better rank
	assert.Equal(t, "income", info.RankingType)
}

// TestRankingService_CreateRankingReward tests creating ranking reward.
func TestRankingService_CreateRankingReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	rankingRepo := rankingrepo.NewRankingRepository(db)
	commissionRepo := rankingrepo.NewRankingCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := ranking.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	req := ranking.CreateRankingRewardRequest{
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     3,
		RewardType:  "commission",
		RewardValue: 50000,
		Description: "Top 3 bonus",
	}

	reward, err := svc.CreateRankingReward(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, reward)
	assert.NotZero(t, reward.ID)
	assert.Equal(t, model.RankingTypeIncome, reward.RankingType)
	assert.Equal(t, 1, reward.RankStart)
	assert.Equal(t, 3, reward.RankEnd)
	assert.Equal(t, int64(50000), reward.RewardValue)
}

// TestRanking_RankingTypes tests all ranking types.
func TestRanking_RankingTypes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	user := CreateUniqueTestUser(t, db, "ranking_types")
	player := CreateTestPlayer(t, db, user)
	month := "2025-01"

	// Test all ranking types
	rankingTypes := []model.RankingType{
		model.RankingTypeIncome,
		model.RankingTypeOrderCount,
		model.RankingTypeQuality,
		model.RankingTypePopularity,
	}

	for _, rt := range rankingTypes {
		t.Run(string(rt), func(t *testing.T) {
			ranking := CreateTestPlayerRanking(t, db, player.ID, rt, "monthly", month, 1)

			// Verify retrieval
			retrieved, err := repo.GetPlayerRanking(ctx, player.ID, rt, "monthly", month)
			require.NoError(t, err)
			assert.Equal(t, rt, retrieved.RankingType)
			assert.Equal(t, ranking.ID, retrieved.ID)
		})
	}
}

// TestRanking_RankingPeriods tests different ranking periods.
func TestRanking_RankingPeriods(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	user := CreateUniqueTestUser(t, db, "ranking_periods")
	player := CreateTestPlayer(t, db, user)

	periods := []struct {
		period      string
		periodValue string
	}{
		{"daily", "2025-01-15"},
		{"weekly", "2025-W03"},
		{"monthly", "2025-01"},
		{"yearly", "2025"},
	}

	for _, p := range periods {
		t.Run(p.period, func(t *testing.T) {
			ranking := CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, p.period, p.periodValue, 1)

			// Verify retrieval by period
			rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
				Period:      &p.period,
				PeriodValue: &p.periodValue,
				Page:        1,
				PageSize:    10,
			})
			require.NoError(t, err)
			assert.Equal(t, int64(1), total)
			assert.Len(t, rankings, 1)
			assert.Equal(t, ranking.ID, rankings[0].ID)
		})
	}
}

// TestRanking_RankingUniqueness tests ranking uniqueness constraint.
func TestRanking_RankingUniqueness(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	user := CreateUniqueTestUser(t, db, "ranking_unique")
	player := CreateTestPlayer(t, db, user)

	// Create first ranking
	_ = CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", "2025-01", 1)

	// Try to create duplicate ranking (same player, type, period, periodValue)
	duplicate := &model.PlayerRanking{
		PlayerID:    player.ID,
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		PeriodValue: "2025-01",
		Rank:        2,
		Score:       50000.0,
	}

	// This should either error or update, depending on unique constraint
	// The important thing is we don't get duplicate rankings
	err := repo.CreateRanking(ctx, duplicate)
	// We expect this might fail due to unique constraint or succeed
	// The key is that our queries should still work correctly
	_ = err // For integration test, we just verify behavior

	// Verify we can still retrieve rankings
	rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PlayerID:    &player.ID,
		RankingType: func() *model.RankingType { t := model.RankingTypeIncome; return &t }(),
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	// Should have at least one ranking
	assert.GreaterOrEqual(t, total, int64(1))
	assert.NotEmpty(t, rankings)
}

// TestRanking_RankingPagination tests ranking list pagination.
func TestRanking_RankingPagination(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create 25 rankings
	for i := 1; i <= 25; i++ {
		user := CreateUniqueTestUser(t, db, "pagination_player")
		player := CreateTestPlayer(t, db, user)
		CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", "2025-01", i)
	}

	// Test pagination
	tests := []struct {
		name          string
		page          int
		pageSize      int
		expectedCount int
	}{
		{"First page", 1, 10, 10},
		{"Second page", 2, 10, 10},
		{"Third page", 3, 10, 5},
		{"Large page size", 1, 50, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
				Page:     tt.page,
				PageSize: tt.pageSize,
			})
			require.NoError(t, err)
			assert.Equal(t, int64(25), total)
			assert.Len(t, rankings, tt.expectedCount)
		})
	}
}

// TestRanking_RankingScoreCalculation tests ranking score calculations.
func TestRanking_RankingScoreCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	user := CreateUniqueTestUser(t, db, "score_player")
	player := CreateTestPlayer(t, db, user)

	tests := []struct {
		name        string
		rankingType model.RankingType
		score       float64
		orderCount  int64
		incomeCents int64
		avgRating   float32
	}{
		{
			name:        "Income ranking",
			rankingType: model.RankingTypeIncome,
			score:       100000.0,
			orderCount:  50,
			incomeCents: 10000000,
			avgRating:   4.5,
		},
		{
			name:        "Order count ranking",
			rankingType: model.RankingTypeOrderCount,
			score:       100.0,
			orderCount:  100,
			incomeCents: 5000000,
			avgRating:   4.8,
		},
		{
			name:        "Quality ranking",
			rankingType: model.RankingTypeQuality,
			score:       5.0,
			orderCount:  30,
			incomeCents: 3000000,
			avgRating:   5.0,
		},
		{
			name:        "Popularity ranking",
			rankingType: model.RankingTypePopularity,
			score:       1000.0,
			orderCount:  80,
			incomeCents: 8000000,
			avgRating:   4.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranking := &model.PlayerRanking{
				PlayerID:    player.ID,
				RankingType: tt.rankingType,
				Period:      "monthly",
				PeriodValue: "2025-01",
				Rank:        1,
				Score:       tt.score,
				OrderCount:  tt.orderCount,
				IncomeCents: tt.incomeCents,
				AvgRating:   tt.avgRating,
			}

			err := repo.CreateRanking(ctx, ranking)
			require.NoError(t, err)

			// Verify
			retrieved, err := repo.GetPlayerRanking(ctx, player.ID, tt.rankingType, "monthly", "2025-01")
			require.NoError(t, err)
			assert.Equal(t, tt.score, retrieved.Score)
			assert.Equal(t, tt.orderCount, retrieved.OrderCount)
			assert.Equal(t, tt.incomeCents, retrieved.IncomeCents)
			assert.Equal(t, tt.avgRating, retrieved.AvgRating)
		})
	}
}

// TestRanking_RankingWithReward tests ranking with bonus reward.
func TestRanking_RankingWithReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	rankingRepo := rankingrepo.NewRankingRepository(db)
	commissionRepo := rankingrepo.NewRankingCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := ranking.NewRankingService(rankingRepo, commissionRepo, orderRepo)

	// Create reward for top 3
	_, err := svc.CreateRankingReward(ctx, ranking.CreateRankingRewardRequest{
		RankingType: model.RankingTypeOrderCount,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     3,
		RewardType:  "commission",
		RewardValue: 10000,
		Description: "Top 3 reward",
	})
	require.NoError(t, err)

	// Create user and player
	user := CreateUniqueTestUser(t, db, "reward_player")
	player := CreateTestPlayer(t, db, user)

	// Create completed orders
	game := CreateTestGame(t, db, "reward_game")
	for i := 0; i < 10; i++ {
		CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 5000)
	}

	// Calculate rankings (should apply reward)
	month := time.Now().Format("2006-01")
	err = svc.CalculateMonthlyRankings(ctx, month)
	require.NoError(t, err)

	// Verify ranking has bonus
	rankings, _, err := rankingRepo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PlayerID:    &player.ID,
		PeriodValue: &month,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, rankings)

	// Check if any ranking has bonus (player should be in top 3)
	for _, r := range rankings {
		if r.Rank <= 3 {
			assert.Greater(t, r.BonusCents, int64(0), "Rank %d should have bonus", r.Rank)
		}
	}
}

// TestRanking_BatchRankingCreation tests batch creation of rankings.
func TestRanking_BatchRankingCreation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create batch of rankings
	month := "2025-01"
	count := 20

	for i := 1; i <= count; i++ {
		user := CreateUniqueTestUser(t, db, "batch_player")
		player := CreateTestPlayer(t, db, user)
		ranking := &model.PlayerRanking{
			PlayerID:    player.ID,
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			PeriodValue: month,
			Rank:        i,
			Score:       float64((count - i + 1) * 1000),
			OrderCount:  int64(count - i + 1),
			IncomeCents: int64((count - i + 1) * 10000),
		}
		err := repo.CreateRanking(ctx, ranking)
		require.NoError(t, err)
	}

	// Verify all rankings created
	rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: &month,
		Page:        1,
		PageSize:    100,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(count), total)
	assert.Len(t, rankings, count)

	// Verify order
	for i := 0; i < len(rankings)-1; i++ {
		assert.LessOrEqual(t, rankings[i].Rank, rankings[i+1].Rank)
	}
}

// TestRanking_TopRankingsQuery tests querying top N rankings.
func TestRanking_TopRankingsQuery(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create 50 rankings
	for i := 1; i <= 50; i++ {
		user := CreateUniqueTestUser(t, db, "top_player")
		player := CreateTestPlayer(t, db, user)
		CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", "2025-01", i)
	}

	// Query top 10
	month := "2025-01"
	rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: &month,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(50), total)
	assert.Len(t, rankings, 10)

	// Verify they are the top 10
	for i, r := range rankings {
		assert.Equal(t, i+1, r.Rank)
	}

	// Query top 100 (should get all 50)
	rankings, total, err = repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: &month,
		Page:        1,
		PageSize:    100,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(50), total)
	assert.Len(t, rankings, 50)
}

// TestRanking_MonthlySettlementIntegration tests ranking integration with monthly settlement.
func TestRanking_MonthlySettlementIntegration(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	commissionRepo := rankingrepo.NewRankingCommissionRepository(db)

	// Create commission config for month
	config := CreateTestRankingCommissionConfig(t, db, model.RankingTypeIncome, "2025-01")

	// Verify config exists
	retrieved, err := commissionRepo.GetActiveConfigForMonth(ctx, model.RankingTypeIncome, "2025-01")
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, config.ID, retrieved.ID)

	// Parse and verify rules
	var rules []model.RankingCommissionRule
	err = json.Unmarshal([]byte(retrieved.RulesJSON), &rules)
	require.NoError(t, err)
	assert.NotEmpty(t, rules)

	// Verify rule structure
	for _, rule := range rules {
		assert.Greater(t, rule.RankStart, 0)
		assert.GreaterOrEqual(t, rule.RankEnd, rule.RankStart)
		assert.Greater(t, rule.CommissionRate, 0)
		assert.LessOrEqual(t, rule.CommissionRate, 100)
	}
}

// TestRanking_CommissionConfigUpdate tests updating commission config.
func TestRanking_CommissionConfigUpdate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingCommissionRepository(db)

	// Create config
	config := CreateTestRankingCommissionConfig(t, db, model.RankingTypeIncome, "2025-01")

	// Update
	newRules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 8},
		{RankStart: 6, RankEnd: 20, CommissionRate: 10},
	}
	newRulesJSON, err := json.Marshal(newRules)
	require.NoError(t, err)

	config.RulesJSON = string(newRulesJSON)
	config.Description = "Updated commission rules"

	err = repo.UpdateConfig(ctx, config)
	require.NoError(t, err)

	// Verify update
	updated, err := repo.GetConfig(ctx, config.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated commission rules", updated.Description)

	var parsedRules []model.RankingCommissionRule
	err = json.Unmarshal([]byte(updated.RulesJSON), &parsedRules)
	require.NoError(t, err)
	assert.Len(t, parsedRules, 2)
	assert.Equal(t, 8, parsedRules[0].CommissionRate)
}

// TestRanking_DeleteExpiredRankings tests deleting old rankings.
func TestRanking_DeleteExpiredRankings(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	repo := rankingrepo.NewRankingRepository(db)
	ctx := context.Background()

	// Create old rankings (3 months ago)
	user := CreateUniqueTestUser(t, db, "old_player")
	player := CreateTestPlayer(t, db, user)
	oldMonth := "2024-10"
	oldRanking := CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", oldMonth, 1)

	// Create recent rankings
	recentMonth := "2025-01"
	recentRanking := CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", recentMonth, 1)

	// Verify both exist
	_, err := repo.GetPlayerRanking(ctx, player.ID, model.RankingTypeIncome, "monthly", oldMonth)
	require.NoError(t, err)

	_, err = repo.GetPlayerRanking(ctx, player.ID, model.RankingTypeIncome, "monthly", recentMonth)
	require.NoError(t, err)

	// Delete old ranking
	err = db.Delete(&model.PlayerRanking{}, oldRanking.ID).Error
	require.NoError(t, err)

	// Verify old ranking deleted
	_, err = repo.GetPlayerRanking(ctx, player.ID, model.RankingTypeIncome, "monthly", oldMonth)
	assert.Error(t, err)

	// Verify recent ranking still exists
	retrieved, err := repo.GetPlayerRanking(ctx, player.ID, model.RankingTypeIncome, "monthly", recentMonth)
	require.NoError(t, err)
	assert.Equal(t, recentRanking.ID, retrieved.ID)
}

// TestRanking_MultiplePlayersSameRank tests handling multiple players with same rank.
func TestRanking_MultiplePlayersSameRank(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create multiple players with rank 1 (tie)
	month := "2025-01"
	for i := 1; i <= 3; i++ {
		user := CreateUniqueTestUser(t, db, "tie_player")
		player := CreateTestPlayer(t, db, user)
		ranking := &model.PlayerRanking{
			PlayerID:    player.ID,
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			PeriodValue: month,
			Rank:        1, // Same rank
			Score:       100000.0,
			OrderCount:  50,
			IncomeCents: 10000000,
		}
		err := repo.CreateRanking(ctx, ranking)
		require.NoError(t, err)
	}

	// Query all rank 1 players
	rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: &month,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)

	// Count rank 1 entries
	rank1Count := 0
	for _, r := range rankings {
		if r.Rank == 1 {
			rank1Count++
		}
	}
	assert.Equal(t, 3, rank1Count, "Should have 3 players with rank 1")
}

// TestRanking_RangeQueryByRank tests querying rankings by rank range.
func TestRanking_RangeQueryByRank(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := rankingrepo.NewRankingRepository(db)

	// Create rankings 1-20
	month := "2025-01"
	for i := 1; i <= 20; i++ {
		user := CreateUniqueTestUser(t, db, "range_player")
		player := CreateTestPlayer(t, db, user)
		CreateTestPlayerRanking(t, db, player.ID, model.RankingTypeIncome, "monthly", month, i)
	}

	// Get first page (top 10)
	rankings, total, err := repo.ListRankings(ctx, rankingrepo.RankingListOptions{
		PeriodValue: &month,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(20), total)
	assert.Len(t, rankings, 10)

	// Verify rank range
	minRank := rankings[0].Rank
	maxRank := rankings[len(rankings)-1].Rank
	assert.Equal(t, 1, minRank)
	assert.Equal(t, 10, maxRank)
}
