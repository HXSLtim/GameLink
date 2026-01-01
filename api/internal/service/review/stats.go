package review

import (
	"context"

	"gamelink/internal/repository"
)

// ReviewStatsService 评价统计服务
type ReviewStatsService struct {
	reviews repository.ReviewRepository
}

// NewReviewStatsService 创建评价统计服务
func NewReviewStatsService(reviews repository.ReviewRepository) *ReviewStatsService {
	return &ReviewStatsService{
		reviews: reviews,
	}
}

// GetReviewStatsResponse 评价统计响应
type GetReviewStatsResponse struct {
	TotalReviews       int64         `json:"totalReviews"`
	AverageRating      float64       `json:"averageRating"`
	RatingDistribution map[int]int64 `json:"ratingDistribution"`
}

// GetReviewTrendResponse 评价趋势响应
type GetReviewTrendResponse struct {
	Trend []repository.DateValue `json:"trend"`
}

// GetTopPlayersResponse 陪玩师排行响应
type GetTopPlayersResponse struct {
	Players []repository.PlayerReviewStats `json:"players"`
}

// GetGameStatsResponse 游戏统计响应
type GetGameStatsResponse struct {
	Games []repository.GameReviewStats `json:"games"`
}

// GetReviewStats 获取评价统计概览
// 需求: 6.1 - 显示总评价数、平均评分、各评分段分布
func (s *ReviewStatsService) GetReviewStats(ctx context.Context) (*GetReviewStatsResponse, error) {
	stats, err := s.reviews.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &GetReviewStatsResponse{
		TotalReviews:       stats.TotalReviews,
		AverageRating:      stats.AverageRating,
		RatingDistribution: stats.RatingDistribution,
	}, nil
}

// GetReviewTrend 获取评价趋势（最近30天）
// 需求: 6.2 - 使用折线图显示最近30天的评价数量趋势
func (s *ReviewStatsService) GetReviewTrend(ctx context.Context, days int) (*GetReviewTrendResponse, error) {
	if days <= 0 {
		days = 30
	}

	trend, err := s.reviews.GetTrend(ctx, days)
	if err != nil {
		return nil, err
	}

	return &GetReviewTrendResponse{
		Trend: trend,
	}, nil
}

// GetTopPlayers 获取陪玩师排行榜
// 需求: 6.3 - 显示评价最多的陪玩师和评分最高的陪玩师排行榜
func (s *ReviewStatsService) GetTopPlayers(ctx context.Context, limit int, sortBy string) (*GetTopPlayersResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	var players []repository.PlayerReviewStats
	var err error

	switch sortBy {
	case "rating":
		// 按评分排序
		players, err = s.reviews.GetTopPlayersByRating(ctx, limit)
	default:
		// 默认按评价数量排序
		players, err = s.reviews.GetTopPlayersByReviewCount(ctx, limit)
	}

	if err != nil {
		return nil, err
	}

	return &GetTopPlayersResponse{
		Players: players,
	}, nil
}

// GetGameStats 获取按游戏统计的评价数据
// 需求: 6.4 - 显示各游戏的评价数量和平均评分
func (s *ReviewStatsService) GetGameStats(ctx context.Context) (*GetGameStatsResponse, error) {
	games, err := s.reviews.GetGameStats(ctx)
	if err != nil {
		return nil, err
	}

	return &GetGameStatsResponse{
		Games: games,
	}, nil
}
