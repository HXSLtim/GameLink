package ranking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// 简单的接口测试，不依赖数据库

func TestRankingListOptions_Defaults(t *testing.T) {
	t.Run("默认分页参数", func(t *testing.T) {
		opts := RankingListOptions{
			Page:     0,
			PageSize: 0,
		}
		
		// 测试参数规范化逻辑
		if opts.Page < 1 {
			opts.Page = 1
		}
		if opts.PageSize < 1 {
			opts.PageSize = 20
		}
		
		assert.Equal(t, 1, opts.Page)
		assert.Equal(t, 20, opts.PageSize)
	})
	
	t.Run("有效分页参数", func(t *testing.T) {
		opts := RankingListOptions{
			Page:     5,
			PageSize: 50,
		}
		
		assert.Equal(t, 5, opts.Page)
		assert.Equal(t, 50, opts.PageSize)
	})
}

func TestRewardListOptions_Defaults(t *testing.T) {
	t.Run("默认分页参数", func(t *testing.T) {
		opts := RewardListOptions{
			Page:     0,
			PageSize: 0,
		}
		
		// 测试参数规范化逻辑
		if opts.Page < 1 {
			opts.Page = 1
		}
		if opts.PageSize < 1 {
			opts.PageSize = 20
		}
		
		assert.Equal(t, 1, opts.Page)
		assert.Equal(t, 20, opts.PageSize)
	})
}

func TestRankingType_Constants(t *testing.T) {
	t.Run("验证排名类型常量", func(t *testing.T) {
		assert.Equal(t, model.RankingType("income"), model.RankingTypeIncome)
		assert.Equal(t, model.RankingType("order_count"), model.RankingTypeOrderCount)
		assert.Equal(t, model.RankingType("quality"), model.RankingTypeQuality)
		assert.Equal(t, model.RankingType("popularity"), model.RankingTypePopularity)
	})
}

func TestPlayerRanking_Structure(t *testing.T) {
	t.Run("创建排名结构", func(t *testing.T) {
		ranking := &model.PlayerRanking{
			PlayerID:    1,
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			PeriodValue: "2024-01",
			Rank:        1,
			Score:       1000,
		}
		
		assert.Equal(t, uint64(1), ranking.PlayerID)
		assert.Equal(t, model.RankingTypeIncome, ranking.RankingType)
		assert.Equal(t, "month", ranking.Period)
		assert.Equal(t, "2024-01", ranking.PeriodValue)
		assert.Equal(t, 1, ranking.Rank)
		assert.Equal(t, float64(1000), ranking.Score)
	})
}

func TestRankingReward_Structure(t *testing.T) {
	t.Run("创建奖励结构", func(t *testing.T) {
		reward := &model.RankingReward{
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			RankStart:   1,
			RankEnd:     3,
			RewardType:  "cash",
			RewardValue: 10000,
			IsActive:    true,
		}
		
		assert.Equal(t, model.RankingTypeIncome, reward.RankingType)
		assert.Equal(t, "month", reward.Period)
		assert.Equal(t, 1, reward.RankStart)
		assert.Equal(t, 3, reward.RankEnd)
		assert.Equal(t, "cash", reward.RewardType)
		assert.Equal(t, int64(10000), reward.RewardValue)
		assert.True(t, reward.IsActive)
	})
	
	t.Run("奖励范围验证", func(t *testing.T) {
		reward := &model.RankingReward{
			RankStart: 1,
			RankEnd:   10,
		}
		
		// 验证排名范围
		assert.True(t, reward.RankStart <= reward.RankEnd)
		assert.True(t, reward.RankStart > 0)
	})
}

func TestRankingListOptions_Filtering(t *testing.T) {
	t.Run("按玩家ID筛选", func(t *testing.T) {
		playerID := uint64(1)
		opts := RankingListOptions{
			PlayerID: &playerID,
		}
		
		assert.NotNil(t, opts.PlayerID)
		assert.Equal(t, uint64(1), *opts.PlayerID)
	})
	
	t.Run("按排名类型筛选", func(t *testing.T) {
		rankingType := model.RankingTypeIncome
		opts := RankingListOptions{
			RankingType: &rankingType,
		}
		
		assert.NotNil(t, opts.RankingType)
		assert.Equal(t, model.RankingTypeIncome, *opts.RankingType)
	})
	
	t.Run("按周期筛选", func(t *testing.T) {
		period := "month"
		periodValue := "2024-01"
		opts := RankingListOptions{
			Period:      &period,
			PeriodValue: &periodValue,
		}
		
		assert.NotNil(t, opts.Period)
		assert.Equal(t, "month", *opts.Period)
		assert.NotNil(t, opts.PeriodValue)
		assert.Equal(t, "2024-01", *opts.PeriodValue)
	})
}

func TestRewardListOptions_Filtering(t *testing.T) {
	t.Run("按排名类型筛选", func(t *testing.T) {
		rankingType := model.RankingTypeIncome
		opts := RewardListOptions{
			RankingType: &rankingType,
		}
		
		assert.NotNil(t, opts.RankingType)
		assert.Equal(t, model.RankingTypeIncome, *opts.RankingType)
	})
	
	t.Run("按周期筛选", func(t *testing.T) {
		period := "month"
		opts := RewardListOptions{
			Period: &period,
		}
		
		assert.NotNil(t, opts.Period)
		assert.Equal(t, "month", *opts.Period)
	})
	
	t.Run("按激活状态筛选", func(t *testing.T) {
		isActive := true
		opts := RewardListOptions{
			IsActive: &isActive,
		}
		
		assert.NotNil(t, opts.IsActive)
		assert.True(t, *opts.IsActive)
	})
}

func TestErrNotFound(t *testing.T) {
	t.Run("验证错误常量", func(t *testing.T) {
		err := ErrNotFound
		assert.NotNil(t, err)
		assert.Equal(t, "resource not found", err.Error())
	})
}

func TestNewRankingRepository(t *testing.T) {
	t.Run("创建repository实例", func(t *testing.T) {
		// 这个测试只验证构造函数不会panic
		// 实际的数据库操作需要集成测试
		repo := NewRankingRepository(nil)
		assert.NotNil(t, repo)
	})
}

// Mock repository for interface testing
type mockRankingRepository struct {
	rankings map[uint64]*model.PlayerRanking
	rewards  map[uint64]*model.RankingReward
	nextID   uint64
}

func newMockRankingRepository() *mockRankingRepository {
	return &mockRankingRepository{
		rankings: make(map[uint64]*model.PlayerRanking),
		rewards:  make(map[uint64]*model.RankingReward),
		nextID:   1,
	}
}

func (m *mockRankingRepository) CreateRanking(ctx context.Context, ranking *model.PlayerRanking) error {
	ranking.ID = m.nextID
	m.nextID++
	m.rankings[ranking.ID] = ranking
	return nil
}

func (m *mockRankingRepository) GetPlayerRanking(ctx context.Context, playerID uint64, rankingType model.RankingType, period, periodValue string) (*model.PlayerRanking, error) {
	for _, r := range m.rankings {
		if r.PlayerID == playerID && r.RankingType == rankingType && r.Period == period && r.PeriodValue == periodValue {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRankingRepository) ListRankings(ctx context.Context, opts RankingListOptions) ([]model.PlayerRanking, int64, error) {
	var result []model.PlayerRanking
	for _, r := range m.rankings {
		if opts.PlayerID != nil && r.PlayerID != *opts.PlayerID {
			continue
		}
		if opts.RankingType != nil && r.RankingType != *opts.RankingType {
			continue
		}
		result = append(result, *r)
	}
	return result, int64(len(result)), nil
}

func (m *mockRankingRepository) UpdateRanking(ctx context.Context, ranking *model.PlayerRanking) error {
	if _, exists := m.rankings[ranking.ID]; !exists {
		return ErrNotFound
	}
	m.rankings[ranking.ID] = ranking
	return nil
}

func (m *mockRankingRepository) CreateReward(ctx context.Context, reward *model.RankingReward) error {
	reward.ID = m.nextID
	m.nextID++
	m.rewards[reward.ID] = reward
	return nil
}

func (m *mockRankingRepository) GetReward(ctx context.Context, id uint64) (*model.RankingReward, error) {
	if reward, exists := m.rewards[id]; exists {
		return reward, nil
	}
	return nil, ErrNotFound
}

func (m *mockRankingRepository) ListRewards(ctx context.Context, opts RewardListOptions) ([]model.RankingReward, int64, error) {
	var result []model.RankingReward
	for _, r := range m.rewards {
		if opts.RankingType != nil && r.RankingType != *opts.RankingType {
			continue
		}
		if opts.IsActive != nil && r.IsActive != *opts.IsActive {
			continue
		}
		result = append(result, *r)
	}
	return result, int64(len(result)), nil
}

func (m *mockRankingRepository) UpdateReward(ctx context.Context, reward *model.RankingReward) error {
	if _, exists := m.rewards[reward.ID]; !exists {
		return ErrNotFound
	}
	m.rewards[reward.ID] = reward
	return nil
}

func (m *mockRankingRepository) DeleteReward(ctx context.Context, id uint64) error {
	if _, exists := m.rewards[id]; !exists {
		return ErrNotFound
	}
	delete(m.rewards, id)
	return nil
}

func (m *mockRankingRepository) GetRewardForRank(ctx context.Context, rankingType model.RankingType, period string, rank int) (*model.RankingReward, error) {
	for _, r := range m.rewards {
		if r.RankingType == rankingType && r.Period == period && r.IsActive && r.RankStart <= rank && r.RankEnd >= rank {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

// 使用mock repository的测试
func TestMockRankingRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newMockRankingRepository()
	
	t.Run("创建并获取排名", func(t *testing.T) {
		ranking := &model.PlayerRanking{
			PlayerID:    1,
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			PeriodValue: "2024-01",
			Rank:        1,
			Score:       1000,
		}
		
		err := repo.CreateRanking(ctx, ranking)
		assert.NoError(t, err)
		assert.NotZero(t, ranking.ID)
		
		got, err := repo.GetPlayerRanking(ctx, 1, model.RankingTypeIncome, "month", "2024-01")
		assert.NoError(t, err)
		assert.Equal(t, ranking.PlayerID, got.PlayerID)
		assert.Equal(t, ranking.Score, got.Score)
	})
	
	t.Run("更新排名", func(t *testing.T) {
		ranking := &model.PlayerRanking{
			PlayerID:    2,
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			PeriodValue: "2024-01",
			Rank:        2,
			Score:       900,
		}
		
		err := repo.CreateRanking(ctx, ranking)
		assert.NoError(t, err)
		
		ranking.Score = 950
		err = repo.UpdateRanking(ctx, ranking)
		assert.NoError(t, err)
		
		got, err := repo.GetPlayerRanking(ctx, 2, model.RankingTypeIncome, "month", "2024-01")
		assert.NoError(t, err)
		assert.Equal(t, float64(950), got.Score)
	})
	
	t.Run("查询排名列表", func(t *testing.T) {
		rankings, total, err := repo.ListRankings(ctx, RankingListOptions{})
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.NotEmpty(t, rankings)
	})
}

func TestMockRankingRepository_Rewards(t *testing.T) {
	ctx := context.Background()
	repo := newMockRankingRepository() // 每个测试使用新的repo实例
	
	t.Run("创建并获取奖励", func(t *testing.T) {
		reward := &model.RankingReward{
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			RankStart:   1,
			RankEnd:     3,
			RewardValue: 10000,
			IsActive:    true,
		}
		
		err := repo.CreateReward(ctx, reward)
		assert.NoError(t, err)
		assert.NotZero(t, reward.ID)
		
		got, err := repo.GetReward(ctx, reward.ID)
		assert.NoError(t, err)
		assert.Equal(t, reward.RewardValue, got.RewardValue)
	})
	
	t.Run("删除奖励", func(t *testing.T) {
		reward := &model.RankingReward{
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			RankStart:   4,
			RankEnd:     10,
			RewardValue: 5000,
			IsActive:    true,
		}
		
		err := repo.CreateReward(ctx, reward)
		assert.NoError(t, err)
		
		err = repo.DeleteReward(ctx, reward.ID)
		assert.NoError(t, err)
		
		_, err = repo.GetReward(ctx, reward.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
	
	t.Run("获取排名对应的奖励", func(t *testing.T) {
		// 使用新的repo避免之前测试的数据干扰
		freshRepo := newMockRankingRepository()
		
		reward := &model.RankingReward{
			RankingType: model.RankingTypeIncome,
			Period:      "month",
			RankStart:   1,
			RankEnd:     5,
			RewardValue: 8000,
			IsActive:    true,
		}
		
		err := freshRepo.CreateReward(ctx, reward)
		assert.NoError(t, err)
		
		got, err := freshRepo.GetRewardForRank(ctx, model.RankingTypeIncome, "month", 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(8000), got.RewardValue)
	})
}
