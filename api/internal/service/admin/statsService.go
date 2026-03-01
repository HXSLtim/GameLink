package admin

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/repository"
)

// StatsService 聚合统计查询。
type StatsService struct {
	repo                 repository.StatsRepository
	userBehaviorRepo     repository.UserBehaviorRepository
	userLoginHistoryRepo repository.UserLoginHistoryRepository
}

func NewStatsService(repo repository.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

// SetUserBehaviorRepo 设置用户行为仓储
func (s *StatsService) SetUserBehaviorRepo(repo repository.UserBehaviorRepository) {
	s.userBehaviorRepo = repo
}

// SetUserLoginHistoryRepo 设置登录历史仓储
func (s *StatsService) SetUserLoginHistoryRepo(repo repository.UserLoginHistoryRepository) {
	s.userLoginHistoryRepo = repo
}

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

// UserBehaviorStatsResponse 用户行为统计响应
type UserBehaviorStatsResponse struct {
	DAU            int64   `json:"dau"`            // 日活跃用户
	AvgOnlineTime  string  `json:"avgOnlineTime"`  // 平均在线时长
	AvgConsumption float64 `json:"avgConsumption"` // 人均消费
}

// UserDistributionResponse 用户分布响应
type UserDistributionResponse struct {
	ByRegion []RegionData `json:"byRegion"` // 地域分布
	ByAge    []AgeData    `json:"byAge"`    // 年龄分布
}

// RegionData 地域数据
type RegionData struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// AgeData 年龄数据
type AgeData struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// UserBehaviorStats 用户行为统计
// 返回DAU、平均在线时长、人均消费等指标
func (s *StatsService) UserBehaviorStats(ctx context.Context) (*UserBehaviorStatsResponse, error) {
	metrics, err := s.repo.UserBehaviorStats(ctx)
	if err != nil {
		return nil, err
	}

	return &UserBehaviorStatsResponse{
		DAU:            metrics.DAU,
		AvgOnlineTime:  formatOnlineDuration(metrics.AvgOnlineDurationSecond),
		AvgConsumption: metrics.AvgConsumptionCents / 100.0,
	}, nil
}

// UserActivityTrend 用户活动趋势
// days: 统计最近N天的数据
func (s *StatsService) UserActivityTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 7
	}
	return s.repo.UserActivityTrend(ctx, days)
}

// UserDistribution 用户分布统计
func (s *StatsService) UserDistribution(ctx context.Context) (*UserDistributionResponse, error) {
	metrics, err := s.repo.UserDistribution(ctx)
	if err != nil {
		return nil, err
	}

	regionData := make([]RegionData, 0, len(metrics.ByRegion))
	for _, item := range metrics.ByRegion {
		regionData = append(regionData, RegionData{Name: item.Name, Value: item.Value})
	}

	ageData := make([]AgeData, 0, len(metrics.ByAge))
	for _, item := range metrics.ByAge {
		ageData = append(ageData, AgeData{Name: item.Name, Value: item.Value})
	}

	return &UserDistributionResponse{
		ByRegion: regionData,
		ByAge:    ageData,
	}, nil
}

func formatOnlineDuration(seconds float64) string {
	if seconds <= 0 {
		return "0m"
	}

	minutes := int(seconds+30) / 60
	if minutes <= 0 {
		minutes = 1
	}

	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	remainMinutes := minutes % 60
	if remainMinutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, remainMinutes)
}
