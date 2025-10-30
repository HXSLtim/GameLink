package stats

import (
    "context"
    "time"

    "gamelink/internal/repository"
)

// StatsService 聚合统计查询。
type StatsService struct {
    repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository) *StatsService { return &StatsService{repo: repo} }

func (s *StatsService) Dashboard(ctx context.Context) (repository.Dashboard, error) {
    return s.repo.Dashboard(ctx)
}

func (s *StatsService) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
    return s.repo.RevenueTrend(ctx, days)
}

func (s *StatsService) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) {
    return s.repo.UserGrowth(ctx, days)
}

func (s *StatsService) OrdersByStatus(ctx context.Context) (map[string]int64, error) {
    return s.repo.OrdersByStatus(ctx)
}

func (s *StatsService) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) {
    return s.repo.TopPlayers(ctx, limit)
}

// AuditOverview returns counts grouped by entity and action within a time window.
func (s *StatsService) AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error) {
    return s.repo.AuditOverview(ctx, from, to)
}

// AuditTrend returns daily counts filtered by time and optional entity/action.
func (s *StatsService) AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]repository.DateValue, error) {
    return s.repo.AuditTrend(ctx, from, to, entity, action)
}
