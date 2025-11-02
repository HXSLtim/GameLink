package ranking

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/ranking"
	"gamelink/internal/repository/order"
)

var (
	// ErrNotFound 排名不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
)

// RankingService 排名服务
type RankingService struct {
	rankings    ranking.RankingRepository
	commissions ranking.RankingCommissionRepository
	orders      order.OrderRepository
}

// NewRankingService 创建排名服务
func NewRankingService(
	rankings ranking.RankingRepository,
	commissions ranking.RankingCommissionRepository,
	orders order.OrderRepository,
) *RankingService {
	return &RankingService{
		rankings:    rankings,
		commissions: commissions,
		orders:      orders,
	}
}

// CalculateMonthlyRankings 计算月度排名
// 注意：礼物订单不计入单量和金额排名
func (s *RankingService) CalculateMonthlyRankings(ctx context.Context, month string) error {
	// 获取该月所有已完成订单（排除礼物订单）
	monthStart, _ := time.Parse("2006-01", month)
	monthEnd := monthStart.AddDate(0, 1, 0)

	orders, _, err := s.orders.List(ctx, repository.OrderListOptions{
		DateFrom: &monthStart,
		DateTo:   &monthEnd,
		Statuses: []model.OrderStatus{model.OrderStatusCompleted},
		Page:     1,
		PageSize: 100000,
	})
	if err != nil {
		return err
	}

	// 按陪玩师统计（排除礼物订单）
	playerStats := make(map[uint64]*PlayerMonthStats)
	for _, order := range orders {
		// 跳过礼物订单（不计入排名）
		if order.IsGiftOrder() {
			continue
		}

		playerID := order.GetPlayerID()
		if playerID == 0 {
			continue
		}

		stats, exists := playerStats[playerID]
		if !exists {
			stats = &PlayerMonthStats{PlayerID: playerID}
			playerStats[playerID] = stats
		}

		stats.OrderCount++
		stats.TotalIncome += order.TotalPriceCents
	}

	// 按订单数量排名
	if err := s.saveOrderCountRankings(ctx, month, playerStats); err != nil {
		return err
	}

	// 按金额排名
	if err := s.saveIncomeRankings(ctx, month, playerStats); err != nil {
		return err
	}

	return nil
}

// PlayerMonthStats 玩家月度统计
type PlayerMonthStats struct {
	PlayerID    uint64
	OrderCount  int64
	TotalIncome int64
}

// saveOrderCountRankings 保存订单数量排名
func (s *RankingService) saveOrderCountRankings(ctx context.Context, month string, stats map[uint64]*PlayerMonthStats) error {
	// 转换为切片并排序
	players := make([]*PlayerMonthStats, 0, len(stats))
	for _, stat := range stats {
		players = append(players, stat)
	}

	// 按订单数量排序
	sortByOrderCount(players)

	// 保存排名（只保存前20名）
	limit := 20
	if len(players) < limit {
		limit = len(players)
	}

	for i := 0; i < limit; i++ {
		ranking := &model.PlayerRanking{
			PlayerID:    players[i].PlayerID,
			RankingType: model.RankingTypeOrderCount,
			Period:      "monthly",
			PeriodValue: month,
			Rank:        i + 1,
			Score:       float64(players[i].OrderCount),
			OrderCount:  players[i].OrderCount,
		}

		// 检查是否有排名奖励
		reward, _ := s.rankings.GetRewardForRank(ctx, model.RankingTypeOrderCount, "monthly", ranking.Rank)
		if reward != nil {
			ranking.BonusCents = reward.RewardValue
		}

		if err := s.rankings.CreateRanking(ctx, ranking); err != nil {
			return err
		}
	}

	return nil
}

// saveIncomeRankings 保存收入排名
func (s *RankingService) saveIncomeRankings(ctx context.Context, month string, stats map[uint64]*PlayerMonthStats) error {
	// 转换为切片并排序
	players := make([]*PlayerMonthStats, 0, len(stats))
	for _, stat := range stats {
		players = append(players, stat)
	}

	// 按收入排序
	sortByIncome(players)

	// 保存排名（只保存前20名）
	limit := 20
	if len(players) < limit {
		limit = len(players)
	}

	for i := 0; i < limit; i++ {
		ranking := &model.PlayerRanking{
			PlayerID:    players[i].PlayerID,
			RankingType: model.RankingTypeIncome,
			Period:      "monthly",
			PeriodValue: month,
			Rank:        i + 1,
			Score:       float64(players[i].TotalIncome),
			IncomeCents: players[i].TotalIncome,
		}

		// 检查是否有排名奖励
		reward, _ := s.rankings.GetRewardForRank(ctx, model.RankingTypeIncome, "monthly", ranking.Rank)
		if reward != nil {
			ranking.BonusCents = reward.RewardValue
		}

		if err := s.rankings.CreateRanking(ctx, ranking); err != nil {
			return err
		}
	}

	return nil
}

// GetPlayerRankingInfo 获取陪玩师排名信息（用于抽成计算）
func (s *RankingService) GetPlayerRankingInfo(ctx context.Context, playerID uint64, month string) (*PlayerRankingInfo, error) {
	rankings, _, err := s.rankings.ListRankings(ctx, ranking.RankingListOptions{
		PlayerID:    &playerID,
		PeriodValue: &month,
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		return nil, err
	}

	info := &PlayerRankingInfo{
		PlayerID: playerID,
		Month:    month,
	}

	// 查找单量排名和金额排名中最好的
	bestRank := 999
	for _, ranking := range rankings {
		if ranking.Rank < bestRank {
			bestRank = ranking.Rank
			info.BestRank = ranking.Rank
			info.RankingType = string(ranking.RankingType)
		}
	}

	return info, nil
}

// PlayerRankingInfo 玩家排名信息
type PlayerRankingInfo struct {
	PlayerID    uint64
	Month       string
	BestRank    int
	RankingType string
}

// sortByOrderCount 按订单数量排序
func sortByOrderCount(players []*PlayerMonthStats) {
	// 冒泡排序（简化实现）
	n := len(players)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if players[j].OrderCount < players[j+1].OrderCount {
				players[j], players[j+1] = players[j+1], players[j]
			}
		}
	}
}

// sortByIncome 按收入排序
func sortByIncome(players []*PlayerMonthStats) {
	// 冒泡排序（简化实现）
	n := len(players)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if players[j].TotalIncome < players[j+1].TotalIncome {
				players[j], players[j+1] = players[j+1], players[j]
			}
		}
	}
}

// CreateRankingRewardRequest 创建排名奖励规则请求
type CreateRankingRewardRequest struct {
	RankingType model.RankingType `json:"rankingType" binding:"required"`
	Period      string             `json:"period" binding:"required"`
	RankStart   int                `json:"rankStart" binding:"required,min=1"`
	RankEnd     int                `json:"rankEnd" binding:"required,min=1"`
	RewardType  string             `json:"rewardType" binding:"required,oneof=commission"`
	RewardValue int64              `json:"rewardValue" binding:"required"`
	Description string             `json:"description"`
}

// CreateRankingReward 创建排名奖励规则（管理员）
func (s *RankingService) CreateRankingReward(ctx context.Context, req CreateRankingRewardRequest) (*model.RankingReward, error) {
	reward := &model.RankingReward{
		RankingType: req.RankingType,
		Period:      req.Period,
		RankStart:   req.RankStart,
		RankEnd:     req.RankEnd,
		RewardType:  req.RewardType,
		RewardValue: req.RewardValue,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.rankings.CreateReward(ctx, reward); err != nil {
		return nil, err
	}

	return reward, nil
}

