package ranking

import (
    "context"
    "testing"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    rankingrepo "gamelink/internal/repository/ranking"

    "github.com/stretchr/testify/require"
)

type fakeRankingRepo struct {
    created []*model.PlayerRanking
    list    []model.PlayerRanking
    reward  *model.RankingReward
}

func (f *fakeRankingRepo) CreateRanking(ctx context.Context, r *model.PlayerRanking) error { f.created = append(f.created, r); return nil }
func (f *fakeRankingRepo) GetPlayerRanking(ctx context.Context, playerID uint64, rt model.RankingType, period, pv string) (*model.PlayerRanking, error) {
    return nil, repository.ErrNotFound
}
func (f *fakeRankingRepo) ListRankings(ctx context.Context, opts rankingrepo.RankingListOptions) ([]model.PlayerRanking, int64, error) {
    return f.list, int64(len(f.list)), nil
}
func (f *fakeRankingRepo) UpdateRanking(ctx context.Context, r *model.PlayerRanking) error { return nil }
func (f *fakeRankingRepo) CreateReward(ctx context.Context, reward *model.RankingReward) error { return nil }
func (f *fakeRankingRepo) GetReward(ctx context.Context, id uint64) (*model.RankingReward, error) { return nil, repository.ErrNotFound }
func (f *fakeRankingRepo) ListRewards(ctx context.Context, opts rankingrepo.RewardListOptions) ([]model.RankingReward, int64, error) { return nil, 0, nil }
func (f *fakeRankingRepo) UpdateReward(ctx context.Context, reward *model.RankingReward) error { return nil }
func (f *fakeRankingRepo) DeleteReward(ctx context.Context, id uint64) error { return nil }
func (f *fakeRankingRepo) GetRewardForRank(ctx context.Context, rt model.RankingType, period string, rank int) (*model.RankingReward, error) { return f.reward, nil }

type fakeOrderRepo struct{ orders []model.Order }

func (f *fakeOrderRepo) Create(ctx context.Context, order *model.Order) error { return nil }
func (f *fakeOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { return f.orders, int64(len(f.orders)), nil }
func (f *fakeOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) { return nil, repository.ErrNotFound }
func (f *fakeOrderRepo) Update(ctx context.Context, order *model.Order) error { return nil }
func (f *fakeOrderRepo) Delete(ctx context.Context, id uint64) error { return nil }

type fakeRankingCommissionRepo struct{}

func (f *fakeRankingCommissionRepo) CreateConfig(ctx context.Context, config *model.RankingCommissionConfig) error { return nil }
func (f *fakeRankingCommissionRepo) GetConfig(ctx context.Context, id uint64) (*model.RankingCommissionConfig, error) { return nil, repository.ErrNotFound }
func (f *fakeRankingCommissionRepo) GetActiveConfigForMonth(ctx context.Context, rt model.RankingType, month string) (*model.RankingCommissionConfig, error) {
    return nil, repository.ErrNotFound
}
func (f *fakeRankingCommissionRepo) ListConfigs(ctx context.Context, opts rankingrepo.RankingCommissionConfigListOptions) ([]model.RankingCommissionConfig, int64, error) {
    return nil, 0, nil
}
func (f *fakeRankingCommissionRepo) UpdateConfig(ctx context.Context, config *model.RankingCommissionConfig) error { return nil }
func (f *fakeRankingCommissionRepo) DeleteConfig(ctx context.Context, id uint64) error { return nil }

func TestSortByOrderCount(t *testing.T) {
    players := []*PlayerMonthStats{{PlayerID: 1, OrderCount: 2}, {PlayerID: 2, OrderCount: 5}, {PlayerID: 3, OrderCount: 1}}
    sortByOrderCount(players)
    require.Equal(t, uint64(2), players[0].PlayerID)
    require.Equal(t, uint64(1), players[1].PlayerID)
    require.Equal(t, uint64(3), players[2].PlayerID)
}

func TestSortByIncome(t *testing.T) {
    players := []*PlayerMonthStats{{PlayerID: 1, TotalIncome: 200}, {PlayerID: 2, TotalIncome: 500}, {PlayerID: 3, TotalIncome: 100}}
    sortByIncome(players)
    require.Equal(t, uint64(2), players[0].PlayerID)
    require.Equal(t, uint64(1), players[1].PlayerID)
    require.Equal(t, uint64(3), players[2].PlayerID)
}

func TestSaveOrderCountRankings(t *testing.T) {
    repo := &fakeRankingRepo{}
    srv := NewRankingService(repo, &fakeRankingCommissionRepo{}, &fakeOrderRepo{})
    stats := map[uint64]*PlayerMonthStats{1: {PlayerID: 1, OrderCount: 2}, 2: {PlayerID: 2, OrderCount: 5}, 3: {PlayerID: 3, OrderCount: 1}}
    err := srv.saveOrderCountRankings(context.Background(), "2025-10", stats)
    require.NoError(t, err)
    require.Len(t, repo.created, 3)
    require.Equal(t, model.RankingTypeOrderCount, repo.created[0].RankingType)
    require.Equal(t, 1, repo.created[0].Rank)
    require.Equal(t, float64(5), repo.created[0].Score)
}

func TestSaveIncomeRankings_WithReward(t *testing.T) {
    repo := &fakeRankingRepo{reward: &model.RankingReward{RewardValue: 1000}}
    srv := NewRankingService(repo, &fakeRankingCommissionRepo{}, &fakeOrderRepo{})
    stats := map[uint64]*PlayerMonthStats{1: {PlayerID: 1, TotalIncome: 200}, 2: {PlayerID: 2, TotalIncome: 500}}
    err := srv.saveIncomeRankings(context.Background(), "2025-10", stats)
    require.NoError(t, err)
    require.Len(t, repo.created, 2)
    require.Equal(t, model.RankingTypeIncome, repo.created[0].RankingType)
    require.Equal(t, int64(1000), repo.created[0].BonusCents)
}

func TestGetPlayerRankingInfo(t *testing.T) {
    repo := &fakeRankingRepo{list: []model.PlayerRanking{{PlayerID: 8, RankingType: model.RankingTypeIncome, Period: "monthly", PeriodValue: "2025-10", Rank: 3}, {PlayerID: 8, RankingType: model.RankingTypeOrderCount, Period: "monthly", PeriodValue: "2025-10", Rank: 2}}}
    srv := NewRankingService(repo, &fakeRankingCommissionRepo{}, &fakeOrderRepo{})
    info, err := srv.GetPlayerRankingInfo(context.Background(), 8, "2025-10")
    require.NoError(t, err)
    require.Equal(t, 2, info.BestRank)
    require.Equal(t, "order_count", info.RankingType)
}

func TestCalculateMonthlyRankings_FiltersGiftOrders(t *testing.T) {
    repo := &fakeRankingRepo{}
    giftRecipient := uint64(99)
    p1 := uint64(1)
    orders := []model.Order{{Status: model.OrderStatusCompleted, PlayerID: &p1, TotalPriceCents: 100}, {Status: model.OrderStatusCompleted, RecipientPlayerID: &giftRecipient, PlayerID: &p1, TotalPriceCents: 200}}
    srv := NewRankingService(repo, &fakeRankingCommissionRepo{}, &fakeOrderRepo{orders: orders})
    err := srv.CalculateMonthlyRankings(context.Background(), "2025-10")
    require.NoError(t, err)
    require.Len(t, repo.created, 2)
    require.Equal(t, model.RankingTypeOrderCount, repo.created[0].RankingType)
    require.Equal(t, model.RankingTypeIncome, repo.created[1].RankingType)
    require.Equal(t, float64(1), repo.created[0].Score)
    require.Equal(t, float64(100), repo.created[1].Score)
}
