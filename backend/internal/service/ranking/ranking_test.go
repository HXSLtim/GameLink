package ranking

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoRanking "gamelink/internal/repository/ranking"
)

func TestRankingService_CalculateMonthlyRankings(t *testing.T) {
	ctx := context.Background()
	rankRepo := newFakeRankingRepository()
	orderRepo := &fakeOrderRepository{
		orders: []model.Order{
			newTestOrder(1, 1, 1_000),
			newTestOrder(1, 2, 500),
			newTestOrder(2, 3, 300),
			giftOrder(3, 4, 200),
		},
	}

	// 只为收入排名第一设置奖励，便于检测 BonusCents
	rankRepo.rewards = []*model.RankingReward{
		{
			ID:          1,
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			RankStart:   1,
			RankEnd:     1,
			RewardValue: 8888,
			IsActive:    true,
		},
	}

	service := NewRankingService(rankRepo, nil, orderRepo)

	require.NoError(t, service.CalculateMonthlyRankings(ctx, "2025-01"))

	// 确认订单仓储的筛选条件正确
	require.NotNil(t, orderRepo.lastOpts.DateFrom)
	require.Equal(t, "2025-01-01", orderRepo.lastOpts.DateFrom.Format("2006-01-02"))
	require.NotNil(t, orderRepo.lastOpts.DateTo)
	require.Equal(t, "2025-02-01", orderRepo.lastOpts.DateTo.Format("2006-01-02"))
	require.Len(t, orderRepo.lastOpts.Statuses, 1)
	assert.Equal(t, model.OrderStatusCompleted, orderRepo.lastOpts.Statuses[0])

	// CalculateMonthlyRankings 会分别写入订单数和收入排名
	assert.Len(t, rankRepo.createdRankings, 4)

	var incomeLeader *model.PlayerRanking
	for _, r := range rankRepo.createdRankings {
		if r.RankingType == model.RankingTypeIncome && r.Rank == 1 {
			incomeLeader = r
		}
		if r.RankingType == model.RankingTypeOrderCount && r.PlayerID == 2 {
			assert.Equal(t, 2, r.Rank)
			assert.Equal(t, float64(1), r.Score)
		}
	}

	require.NotNil(t, incomeLeader, "expected income leader in created rankings")
	assert.Equal(t, uint64(1), incomeLeader.PlayerID)
	assert.Equal(t, int64(1500), incomeLeader.IncomeCents)
	assert.Equal(t, int64(8888), incomeLeader.BonusCents)

	// Gift 订单应被排除
	for _, r := range rankRepo.createdRankings {
		assert.NotEqual(t, uint64(4), r.PlayerID)
	}
}

func TestRankingService_GetPlayerRankingInfo(t *testing.T) {
	ctx := context.Background()
	rankRepo := newFakeRankingRepository()
	playerID := uint64(42)
	rankRepo.listResults[playerID] = []model.PlayerRanking{
		{PlayerID: playerID, RankingType: model.RankingTypeOrderCount, Rank: 5},
		{PlayerID: playerID, RankingType: model.RankingTypeIncome, Rank: 2},
	}

	service := NewRankingService(rankRepo, nil, &fakeOrderRepository{})

	info, err := service.GetPlayerRankingInfo(ctx, playerID, "2025-01")
	require.NoError(t, err)
	assert.Equal(t, playerID, info.PlayerID)
	assert.Equal(t, 2, info.BestRank)
	assert.Equal(t, "income", info.RankingType)
}

func TestRankingService_CreateRankingReward(t *testing.T) {
	ctx := context.Background()
	rankRepo := newFakeRankingRepository()
	service := NewRankingService(rankRepo, nil, &fakeOrderRepository{})

	req := CreateRankingRewardRequest{
		RankingType: model.RankingTypeIncome,
		Period:      "monthly",
		RankStart:   1,
		RankEnd:     3,
		RewardType:  "commission",
		RewardValue: 3_000,
		Description: "Top 3 bonus",
	}

	reward, err := service.CreateRankingReward(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, reward)
	assert.True(t, reward.IsActive)

	require.Len(t, rankRepo.rewards, 1)
	assert.Equal(t, req.RankEnd, int(rankRepo.rewards[0].RankEnd))
	assert.Equal(t, req.RewardValue, rankRepo.rewards[0].RewardValue)
}

// ---- test helpers ----------------------------------------------------------------

type fakeOrderRepository struct {
	orders   []model.Order
	lastOpts repository.OrderListOptions
}

var _ repository.OrderRepository = (*fakeOrderRepository)(nil)

func (f *fakeOrderRepository) Create(context.Context, *model.Order) error {
	return errors.New("not implemented")
}

func (f *fakeOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	f.lastOpts = opts
	return f.orders, int64(len(f.orders)), nil
}

func (f *fakeOrderRepository) Get(context.Context, uint64) (*model.Order, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeOrderRepository) Update(context.Context, *model.Order) error {
	return errors.New("not implemented")
}

func (f *fakeOrderRepository) Delete(context.Context, uint64) error {
	return errors.New("not implemented")
}

type fakeRankingRepository struct {
	createdRankings []*model.PlayerRanking
	rewards         []*model.RankingReward
	listResults     map[uint64][]model.PlayerRanking
}

var _ repoRanking.RankingRepository = (*fakeRankingRepository)(nil)

func newFakeRankingRepository() *fakeRankingRepository {
	return &fakeRankingRepository{
		listResults: make(map[uint64][]model.PlayerRanking),
	}
}

func (f *fakeRankingRepository) CreateRanking(ctx context.Context, ranking *model.PlayerRanking) error {
	copy := *ranking
	f.createdRankings = append(f.createdRankings, &copy)
	return nil
}

func (f *fakeRankingRepository) GetPlayerRanking(context.Context, uint64, model.RankingType, string, string) (*model.PlayerRanking, error) {
	return nil, repoRanking.ErrNotFound
}

func (f *fakeRankingRepository) ListRankings(ctx context.Context, opts repoRanking.RankingListOptions) ([]model.PlayerRanking, int64, error) {
	if opts.PlayerID == nil {
		return nil, 0, nil
	}
	list := f.listResults[*opts.PlayerID]
	results := make([]model.PlayerRanking, len(list))
	copy(results, list)
	return results, int64(len(results)), nil
}

func (f *fakeRankingRepository) UpdateRanking(context.Context, *model.PlayerRanking) error {
	return nil
}

func (f *fakeRankingRepository) CreateReward(ctx context.Context, reward *model.RankingReward) error {
	copy := *reward
	if copy.ID == 0 {
		copy.ID = uint64(len(f.rewards) + 1)
		reward.ID = copy.ID
	}
	f.rewards = append(f.rewards, &copy)
	return nil
}

func (f *fakeRankingRepository) GetReward(ctx context.Context, id uint64) (*model.RankingReward, error) {
	for _, r := range f.rewards {
		if r.ID == id {
			copy := *r
			return &copy, nil
		}
	}
	return nil, repoRanking.ErrNotFound
}

func (f *fakeRankingRepository) ListRewards(context.Context, repoRanking.RewardListOptions) ([]model.RankingReward, int64, error) {
	return nil, 0, nil
}

func (f *fakeRankingRepository) UpdateReward(context.Context, *model.RankingReward) error { return nil }

func (f *fakeRankingRepository) DeleteReward(context.Context, uint64) error { return nil }

func (f *fakeRankingRepository) GetRewardForRank(ctx context.Context, rankingType model.RankingType, period string, rank int) (*model.RankingReward, error) {
	for _, reward := range f.rewards {
		if reward.RankingType == rankingType && reward.Period == period && reward.IsActive &&
			reward.RankStart <= rank && reward.RankEnd >= rank {
			copy := *reward
			return &copy, nil
		}
	}
	return nil, repoRanking.ErrNotFound
}

func newTestOrder(playerID uint64, orderNo int, amount int64) model.Order {
	order := model.Order{
		OrderNo:         fmt.Sprintf("NO-%d", orderNo),
		TotalPriceCents: amount,
		Status:          model.OrderStatusCompleted,
	}
	order.SetPlayerID(playerID)
	return order
}

func giftOrder(playerID, orderNo uint64, amount int64) model.Order {
	order := newTestOrder(playerID, int(orderNo), amount)
	order.PlayerID = nil
	order.RecipientPlayerID = &playerID
	return order
}
