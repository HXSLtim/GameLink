package player

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	repoiface "gamelink/internal/repository/interfaces"
)

const (
	playerStatsTodayTTL    = 30 * time.Second
	playerStatsOverviewTTL = 2 * time.Minute
)

// PlayerStatsToday represents today's stats.
type PlayerStatsToday struct {
	OrderCount    int64   `json:"orderCount"`
	EarningsCents int64   `json:"earningsCents"`
	RatingAverage float32 `json:"ratingAverage"`
}

// PlayerStatsOverview represents overall stats.
type PlayerStatsOverview struct {
	TotalOrders        int64   `json:"totalOrders"`
	TotalEarningsCents int64   `json:"totalEarningsCents"`
	RatingAverage      float32 `json:"ratingAverage"`
}

// GetTodayStats returns today's order count, earnings, and rating.
func (s *PlayerService) GetTodayStats(ctx context.Context, userID uint64) (*PlayerStatsToday, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24*time.Hour - time.Nanosecond)

	cacheKey := fmt.Sprintf("player:stats:today:%d:%s", player.ID, start.Format("20060102"))
	if s.cache != nil {
		if cached, ok, _ := s.cache.Get(ctx, cacheKey); ok {
			var payload PlayerStatsToday
			if err := json.Unmarshal([]byte(cached), &payload); err == nil {
				return &payload, nil
			}
		}
	}

	opts := repoiface.OrderListOptions{
		PlayerID: &player.ID,
		DateFrom: &start,
		DateTo:   &end,
		Page:     1,
		PageSize: 1,
	}
	_, total, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	earnings, err := s.sumEarnings(ctx, player.ID, &start, &end)
	if err != nil {
		return nil, err
	}

	payload := &PlayerStatsToday{
		OrderCount:    total,
		EarningsCents: earnings,
		RatingAverage: player.RatingAverage,
	}
	_ = s.cacheStats(ctx, cacheKey, *payload, playerStatsTodayTTL)
	return payload, nil
}

// GetOverviewStats returns total orders, earnings, and rating.
func (s *PlayerService) GetOverviewStats(ctx context.Context, userID uint64) (*PlayerStatsOverview, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("player:stats:overview:%d", player.ID)
	if s.cache != nil {
		if cached, ok, _ := s.cache.Get(ctx, cacheKey); ok {
			var payload PlayerStatsOverview
			if err := json.Unmarshal([]byte(cached), &payload); err == nil {
				return &payload, nil
			}
		}
	}
	opts := repoiface.OrderListOptions{
		PlayerID: &player.ID,
		Page:     1,
		PageSize: 1,
	}
	_, total, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	earnings, err := s.sumEarnings(ctx, player.ID, nil, nil)
	if err != nil {
		return nil, err
	}
	payload := &PlayerStatsOverview{
		TotalOrders:        total,
		TotalEarningsCents: earnings,
		RatingAverage:      player.RatingAverage,
	}
	_ = s.cacheStats(ctx, cacheKey, *payload, playerStatsOverviewTTL)
	return payload, nil
}

func (s *PlayerService) sumEarnings(ctx context.Context, playerID uint64, dateFrom, dateTo *time.Time) (int64, error) {
	statuses := []model.OrderStatus{model.OrderStatusCompleted}
	pageSize := 200
	page := 1
	var sum int64
	for {
		opts := repoiface.OrderListOptions{
			PlayerID: &playerID,
			Statuses: statuses,
			DateFrom: dateFrom,
			DateTo:   dateTo,
			Page:     page,
			PageSize: pageSize,
		}
		orders, _, err := s.orders.List(ctx, opts)
		if err != nil {
			return 0, err
		}
		for _, order := range orders {
			sum += order.PlayerIncomeCents
		}
		if len(orders) < pageSize {
			break
		}
		page++
	}
	return sum, nil
}

func (s *PlayerService) cacheStats(ctx context.Context, key string, payload any, ttl time.Duration) error {
	if s.cache == nil {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return s.cache.Set(ctx, key, string(data), ttl)
}
