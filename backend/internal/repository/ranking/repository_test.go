package ranking

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gamelink/internal/model"
)

func setupRankingDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.PlayerRanking{}, &model.RankingReward{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		t.Cleanup(func() {
			_ = sqlDB.Close()
		})
	}

	return db
}

func TestRankingRepository_CreateAndGetRanking(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingRepository(db)
	ctx := context.Background()

	ranking := &model.PlayerRanking{
		PlayerID:    1,
		RankingType: model.RankingTypeOrderCount,
		Period:      "monthly",
		PeriodValue: "2025-01",
		Rank:        1,
		Score:       10,
		OrderCount:  10,
	}

	require.NoError(t, repo.CreateRanking(ctx, ranking))
	require.NotZero(t, ranking.ID)

	got, err := repo.GetPlayerRanking(ctx, ranking.PlayerID, ranking.RankingType, ranking.Period, ranking.PeriodValue)
	require.NoError(t, err)
	assert.Equal(t, ranking.PlayerID, got.PlayerID)
	assert.Equal(t, ranking.Score, got.Score)
	assert.Equal(t, ranking.OrderCount, got.OrderCount)

	_, err = repo.GetPlayerRanking(ctx, 9999, ranking.RankingType, ranking.Period, ranking.PeriodValue)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingRepository_ListRankingsWithFilters(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingRepository(db)
	ctx := context.Background()

	period := "monthly"
	value := "2025-01"
	rt := model.RankingTypeOrderCount

	for i := 0; i < 25; i++ {
		require.NoError(t, repo.CreateRanking(ctx, &model.PlayerRanking{
			PlayerID:    uint64(i + 1),
			RankingType: rt,
			Period:      period,
			PeriodValue: value,
			Rank:        i + 1,
			Score:       float64(i + 1),
		}))
	}

	t.Run("pagination", func(t *testing.T) {
		opts := RankingListOptions{
			RankingType: &rt,
			Period:      &period,
			PeriodValue: &value,
			Page:        2,
			PageSize:    10,
		}

		rankings, total, err := repo.ListRankings(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		require.Len(t, rankings, 10)
		assert.Equal(t, 11, rankings[0].Rank)
		assert.Equal(t, 20, rankings[len(rankings)-1].Rank)
	})

	t.Run("filter by player with default pagination", func(t *testing.T) {
		playerID := uint64(5)
		opts := RankingListOptions{
			PlayerID:    &playerID,
			Period:      &period,
			PeriodValue: &value,
			Page:        0,
			PageSize:    0,
		}

		result, total, err := repo.ListRankings(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, result, 1)
		assert.Equal(t, playerID, result[0].PlayerID)
		assert.Equal(t, 5, result[0].Rank)
	})
}

func TestRankingRepository_ListRankings_DefaultPaginationAndOrdering(t *testing.T) {
    db := setupRankingDB(t)
    repo := NewRankingRepository(db)
    ctx := context.Background()

    period := "monthly"
    value := "2025-01"
    rt := model.RankingTypeOrderCount

    for i := 0; i < 25; i++ {
        require.NoError(t, repo.CreateRanking(ctx, &model.PlayerRanking{
            PlayerID:    uint64(i + 1),
            RankingType: rt,
            Period:      period,
            PeriodValue: value,
            Rank:        i + 1,
            Score:       float64(i + 1),
        }))
    }

    opts := RankingListOptions{RankingType: &rt, Period: &period, PeriodValue: &value, Page: -1, PageSize: -1}
    rankings, total, err := repo.ListRankings(ctx, opts)
    require.NoError(t, err)
    assert.Equal(t, int64(25), total)
    require.Len(t, rankings, 20)
    assert.Equal(t, 1, rankings[0].Rank)
    assert.Equal(t, 20, rankings[len(rankings)-1].Rank)
}

func TestRankingRepository_ListRewardsOrdering(t *testing.T) {
    db := setupRankingDB(t)
    repo := NewRankingRepository(db)
    ctx := context.Background()

    rt := model.RankingTypeIncome
    require.NoError(t, repo.CreateReward(ctx, &model.RankingReward{RankingType: rt, Period: "monthly", RankStart: 5, RankEnd: 10, RewardType: "cash", RewardValue: 100}))
    require.NoError(t, repo.CreateReward(ctx, &model.RankingReward{RankingType: rt, Period: "monthly", RankStart: 1, RankEnd: 3, RewardType: "cash", RewardValue: 200}))
    require.NoError(t, repo.CreateReward(ctx, &model.RankingReward{RankingType: rt, Period: "monthly", RankStart: 4, RankEnd: 4, RewardType: "cash", RewardValue: 150}))

    m := "monthly"
    list, total, err := repo.ListRewards(ctx, RewardListOptions{RankingType: &rt, Period: &m, Page: 1, PageSize: 10})
    require.NoError(t, err)
    assert.Equal(t, int64(3), total)
    require.Len(t, list, 3)
    assert.Equal(t, 1, list[0].RankStart)
    assert.Equal(t, 4, list[1].RankStart)
    assert.Equal(t, 5, list[2].RankStart)
}

func TestRankingRepository_UpdateRanking(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingRepository(db)
	ctx := context.Background()

	ranking := &model.PlayerRanking{
		PlayerID:    99,
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		PeriodValue: "2025-02",
		Rank:        1,
		Score:       100,
		IncomeCents: 1000,
	}

	require.NoError(t, repo.CreateRanking(ctx, ranking))

	ranking.Score = 150
	ranking.IncomeCents = 2500

	require.NoError(t, repo.UpdateRanking(ctx, ranking))

	got, err := repo.GetPlayerRanking(ctx, ranking.PlayerID, ranking.RankingType, ranking.Period, ranking.PeriodValue)
	require.NoError(t, err)
	assert.Equal(t, float64(150), got.Score)
	assert.Equal(t, int64(2500), got.IncomeCents)
}

func TestRankingRepository_RewardsLifecycle(t *testing.T) {
	db := setupRankingDB(t)
	repo := NewRankingRepository(db)
	ctx := context.Background()

	reward := &model.RankingReward{
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     3,
		RewardType:  "cash",
		RewardValue: 5000,
		IsActive:    true,
	}

	require.NoError(t, repo.CreateReward(ctx, reward))
	require.NotZero(t, reward.ID)

	stored, err := repo.GetReward(ctx, reward.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5000), stored.RewardValue)

	got, err := repo.GetRewardForRank(ctx, reward.RankingType, reward.Period, 2)
	require.NoError(t, err)
	assert.Equal(t, reward.ID, got.ID)

	reward.RewardValue = 8800
	require.NoError(t, repo.UpdateReward(ctx, reward))

	active := true
	list, total, err := repo.ListRewards(ctx, RewardListOptions{
		RankingType: &reward.RankingType,
		Period:      &reward.Period,
		IsActive:    &active,
		Page:        1,
		PageSize:    10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, int64(8800), list[0].RewardValue)

	reward.IsActive = false
	require.NoError(t, repo.UpdateReward(ctx, reward))

	inactive := false
	inactiveList, total, err := repo.ListRewards(ctx, RewardListOptions{
		IsActive: &inactive,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, inactiveList, 1)
	assert.False(t, inactiveList[0].IsActive)

	require.NoError(t, repo.DeleteReward(ctx, reward.ID))

	_, err = repo.GetReward(ctx, reward.ID)
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetRewardForRank(ctx, reward.RankingType, reward.Period, 2)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestRankingRepository_GetRewardForRank_NotFound(t *testing.T) {
    db := setupRankingDB(t)
    repo := NewRankingRepository(db)
    ctx := context.Background()
    rt := model.RankingTypeIncome
    require.NoError(t, repo.CreateReward(ctx, &model.RankingReward{RankingType: rt, Period: "monthly", RankStart: 1, RankEnd: 3, RewardType: "cash", RewardValue: 1000, IsActive: true}))
    _, err := repo.GetRewardForRank(ctx, rt, "monthly", 100)
    require.ErrorIs(t, err, ErrNotFound)
}
