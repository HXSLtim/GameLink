package admin

import (
	"context"
	"time"

	"gamelink/internal/repository"
)

// StatsService 聚合统计查询。
type StatsService struct {
	repo                repository.StatsRepository
	userBehaviorRepo    repository.UserBehaviorRepository
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
	DAU            int64   `json:"dau"`           // 日活跃用户
	AvgOnlineTime  string  `json:"avgOnlineTime"` // 平均在线时长
	AvgConsumption float64 `json:"avgConsumption"`// 人均消费
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
	// 简化实现:返回模拟数据,实际应从数据库查询
	// TODO: 实现真实的数据统计逻辑
	return &UserBehaviorStatsResponse{
		DAU:            1128,
		AvgOnlineTime:  "45m",
		AvgConsumption: 128.50,
	}, nil
}

// UserActivityTrend 用户活动趋势
// days: 统计最近N天的数据
func (s *StatsService) UserActivityTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 7
	}

	// 简化实现:返回模拟数据
	// TODO: 从UserBehavior表统计每日活跃用户数
	result := []repository.DateValue{
		{Date: time.Now().AddDate(0, 0, -6).Format("2006-01-02"), Value: 4000},
		{Date: time.Now().AddDate(0, 0, -5).Format("2006-01-02"), Value: 3000},
		{Date: time.Now().AddDate(0, 0, -4).Format("2006-01-02"), Value: 2000},
		{Date: time.Now().AddDate(0, 0, -3).Format("2006-01-02"), Value: 2780},
		{Date: time.Now().AddDate(0, 0, -2).Format("2006-01-02"), Value: 1890},
		{Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02"), Value: 2390},
		{Date: time.Now().Format("2006-01-02"), Value: 3490},
	}

	return result, nil
}

// UserDistribution 用户分布统计
func (s *StatsService) UserDistribution(ctx context.Context) (*UserDistributionResponse, error) {
	// 简化实现:返回模拟数据
	// TODO: 从User表统计地域和年龄分布
	return &UserDistributionResponse{
		ByRegion: []RegionData{
			{Name: "北京", Value: 400},
			{Name: "上海", Value: 300},
			{Name: "广州", Value: 300},
			{Name: "深圳", Value: 200},
		},
		ByAge: []AgeData{
			{Name: "18-25岁", Value: 500},
			{Name: "26-30岁", Value: 400},
			{Name: "31-35岁", Value: 300},
		},
	}, nil
}
