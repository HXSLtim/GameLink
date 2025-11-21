package stats

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gamelink/internal/repository"
	"gamelink/internal/repository/mocks"
)

// TestNewStatsService 测试构造函数。
func TestNewStatsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)

	svc := NewStatsService(repo)

	if svc == nil {
		t.Fatal("NewStatsService returned nil")
	}

	if svc.repo != repo {
		t.Error("repo not set correctly")
	}
}

// TestDashboard 测试获取仪表板数据。
func TestDashboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取仪表板数据", func(t *testing.T) {
		expectedDashboard := repository.Dashboard{
			TotalUsers:           1000,
			TotalPlayers:         500,
			TotalOrders:          2000,
			TotalGames:           100,
			TotalPaidAmountCents: 5000000, // 50000.00 元 = 5000000 分
			OrdersByStatus:       map[string]int64{"completed": 1500, "pending": 500},
			PaymentsByStatus:     map[string]int64{"paid": 1800, "pending": 200},
		}

		repo.EXPECT().
			Dashboard(ctx).
			Return(expectedDashboard, nil)

		dashboard, err := svc.Dashboard(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if dashboard.TotalUsers != 1000 {
			t.Errorf("Expected TotalUsers 1000, got %d", dashboard.TotalUsers)
		}

		if dashboard.TotalPaidAmountCents != 5000000 {
			t.Errorf("Expected TotalPaidAmountCents 5000000, got %d", dashboard.TotalPaidAmountCents)
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			Dashboard(ctx).
			Return(repository.Dashboard{}, errors.New("database error"))

		dashboard, err := svc.Dashboard(ctx)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if dashboard.TotalUsers != 0 {
			t.Error("Expected empty dashboard on error")
		}
	})
}

// TestRevenueTrend 测试获取收入趋势。
func TestRevenueTrend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取7天收入趋势", func(t *testing.T) {
		expectedTrend := []repository.DateValue{
			{Date: "2025-10-24", Value: 1000},
			{Date: "2025-10-25", Value: 1500},
			{Date: "2025-10-26", Value: 2000},
		}

		repo.EXPECT().
			RevenueTrend(ctx, 7).
			Return(expectedTrend, nil)

		trend, err := svc.RevenueTrend(ctx, 7)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(trend) != 3 {
			t.Errorf("Expected 3 data points, got %d", len(trend))
		}

		if trend[0].Value != 1000 {
			t.Errorf("Expected first value 1000, got %d", trend[0].Value)
		}
	})

	t.Run("成功获取30天收入趋势", func(t *testing.T) {
		repo.EXPECT().
			RevenueTrend(ctx, 30).
			Return([]repository.DateValue{}, nil)

		trend, err := svc.RevenueTrend(ctx, 30)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(trend) != 0 {
			t.Errorf("Expected 0 data points, got %d", len(trend))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			RevenueTrend(ctx, 7).
			Return(nil, errors.New("database error"))

		trend, err := svc.RevenueTrend(ctx, 7)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if trend != nil {
			t.Error("Expected nil trend on error")
		}
	})
}

// TestUserGrowth 测试获取用户增长。
func TestUserGrowth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取用户增长", func(t *testing.T) {
		expectedGrowth := []repository.DateValue{
			{Date: "2025-10-24", Value: 10},
			{Date: "2025-10-25", Value: 15},
			{Date: "2025-10-26", Value: 20},
		}

		repo.EXPECT().
			UserGrowth(ctx, 7).
			Return(expectedGrowth, nil)

		growth, err := svc.UserGrowth(ctx, 7)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(growth) != 3 {
			t.Errorf("Expected 3 data points, got %d", len(growth))
		}
	})

	t.Run("空数据", func(t *testing.T) {
		repo.EXPECT().
			UserGrowth(ctx, 90).
			Return([]repository.DateValue{}, nil)

		growth, err := svc.UserGrowth(ctx, 90)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(growth) != 0 {
			t.Errorf("Expected 0 data points, got %d", len(growth))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			UserGrowth(ctx, 7).
			Return(nil, errors.New("database error"))

		growth, err := svc.UserGrowth(ctx, 7)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if growth != nil {
			t.Error("Expected nil growth on error")
		}
	})
}

// TestOrdersByStatus 测试按状态获取订单。
func TestOrdersByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取订单状态统计", func(t *testing.T) {
		expectedStats := map[string]int64{
			"pending":   100,
			"completed": 500,
			"cancelled": 50,
		}

		repo.EXPECT().
			OrdersByStatus(ctx).
			Return(expectedStats, nil)

		stats, err := svc.OrdersByStatus(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(stats) != 3 {
			t.Errorf("Expected 3 statuses, got %d", len(stats))
		}

		if stats["pending"] != 100 {
			t.Errorf("Expected pending 100, got %d", stats["pending"])
		}

		if stats["completed"] != 500 {
			t.Errorf("Expected completed 500, got %d", stats["completed"])
		}
	})

	t.Run("空统计", func(t *testing.T) {
		repo.EXPECT().
			OrdersByStatus(ctx).
			Return(map[string]int64{}, nil)

		stats, err := svc.OrdersByStatus(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(stats) != 0 {
			t.Errorf("Expected 0 statuses, got %d", len(stats))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			OrdersByStatus(ctx).
			Return(nil, errors.New("database error"))

		stats, err := svc.OrdersByStatus(ctx)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if stats != nil {
			t.Error("Expected nil stats on error")
		}
	})
}

// TestTopPlayers 测试获取顶级玩家。
func TestTopPlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取前10名玩家", func(t *testing.T) {
		expectedPlayers := []repository.PlayerTop{
			{PlayerID: 1, Nickname: "Player1", RatingAverage: 4.8, RatingCount: 100},
			{PlayerID: 2, Nickname: "Player2", RatingAverage: 4.5, RatingCount: 80},
			{PlayerID: 3, Nickname: "Player3", RatingAverage: 4.2, RatingCount: 60},
		}

		repo.EXPECT().
			TopPlayers(ctx, 10).
			Return(expectedPlayers, nil)

		players, err := svc.TopPlayers(ctx, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(players) != 3 {
			t.Errorf("Expected 3 players, got %d", len(players))
		}

		if players[0].RatingAverage != 4.8 {
			t.Errorf("Expected first player rating 4.8, got %f", players[0].RatingAverage)
		}
	})

	t.Run("获取前5名玩家", func(t *testing.T) {
		expectedPlayers := []repository.PlayerTop{
			{PlayerID: 1, Nickname: "Player1", RatingAverage: 4.8, RatingCount: 100},
		}

		repo.EXPECT().
			TopPlayers(ctx, 5).
			Return(expectedPlayers, nil)

		players, err := svc.TopPlayers(ctx, 5)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(players) != 1 {
			t.Errorf("Expected 1 player, got %d", len(players))
		}
	})

	t.Run("空结果", func(t *testing.T) {
		repo.EXPECT().
			TopPlayers(ctx, 10).
			Return([]repository.PlayerTop{}, nil)

		players, err := svc.TopPlayers(ctx, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(players) != 0 {
			t.Errorf("Expected 0 players, got %d", len(players))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			TopPlayers(ctx, 10).
			Return(nil, errors.New("database error"))

		players, err := svc.TopPlayers(ctx, 10)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if players != nil {
			t.Error("Expected nil players on error")
		}
	})
}

// TestAuditOverview 测试审计概览。
func TestAuditOverview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取审计概览", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -7)
		to := now

		expectedEntity := map[string]int64{
			"user":   100,
			"order":  200,
			"player": 150,
		}

		expectedAction := map[string]int64{
			"create": 200,
			"update": 150,
			"delete": 100,
		}

		repo.EXPECT().
			AuditOverview(ctx, &from, &to).
			Return(expectedEntity, expectedAction, nil)

		entityStats, actionStats, err := svc.AuditOverview(ctx, &from, &to)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(entityStats) != 3 {
			t.Errorf("Expected 3 entities, got %d", len(entityStats))
		}

		if len(actionStats) != 3 {
			t.Errorf("Expected 3 actions, got %d", len(actionStats))
		}

		if entityStats["user"] != 100 {
			t.Errorf("Expected user count 100, got %d", entityStats["user"])
		}

		if actionStats["create"] != 200 {
			t.Errorf("Expected create count 200, got %d", actionStats["create"])
		}
	})

	t.Run("使用nil时间范围", func(t *testing.T) {
		repo.EXPECT().
			AuditOverview(ctx, nil, nil).
			Return(map[string]int64{}, map[string]int64{}, nil)

		entityStats, actionStats, err := svc.AuditOverview(ctx, nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(entityStats) != 0 {
			t.Errorf("Expected 0 entities, got %d", len(entityStats))
		}

		if len(actionStats) != 0 {
			t.Errorf("Expected 0 actions, got %d", len(actionStats))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -7)
		to := now

		repo.EXPECT().
			AuditOverview(ctx, &from, &to).
			Return(nil, nil, errors.New("database error"))

		entityStats, actionStats, err := svc.AuditOverview(ctx, &from, &to)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if entityStats != nil {
			t.Error("Expected nil entity stats on error")
		}

		if actionStats != nil {
			t.Error("Expected nil action stats on error")
		}
	})
}

// TestAuditTrend 测试审计趋势。
func TestAuditTrend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockStatsRepository(ctrl)
	svc := NewStatsService(repo)

	ctx := context.Background()

	t.Run("成功获取审计趋势（按实体）", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -7)
		to := now

		expectedTrend := []repository.DateValue{
			{Date: "2025-10-24", Value: 10},
			{Date: "2025-10-25", Value: 15},
			{Date: "2025-10-26", Value: 20},
		}

		repo.EXPECT().
			AuditTrend(ctx, &from, &to, "user", "").
			Return(expectedTrend, nil)

		trend, err := svc.AuditTrend(ctx, &from, &to, "user", "")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(trend) != 3 {
			t.Errorf("Expected 3 data points, got %d", len(trend))
		}

		if trend[0].Value != 10 {
			t.Errorf("Expected first value 10, got %d", trend[0].Value)
		}
	})

	t.Run("成功获取审计趋势（按实体和操作）", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -7)
		to := now

		expectedTrend := []repository.DateValue{
			{Date: "2025-10-24", Value: 5},
			{Date: "2025-10-25", Value: 8},
		}

		repo.EXPECT().
			AuditTrend(ctx, &from, &to, "order", "create").
			Return(expectedTrend, nil)

		trend, err := svc.AuditTrend(ctx, &from, &to, "order", "create")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(trend) != 2 {
			t.Errorf("Expected 2 data points, got %d", len(trend))
		}
	})

	t.Run("空趋势", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -30)
		to := now

		repo.EXPECT().
			AuditTrend(ctx, &from, &to, "", "").
			Return([]repository.DateValue{}, nil)

		trend, err := svc.AuditTrend(ctx, &from, &to, "", "")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(trend) != 0 {
			t.Errorf("Expected 0 data points, got %d", len(trend))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		now := time.Now()
		from := now.AddDate(0, 0, -7)
		to := now

		repo.EXPECT().
			AuditTrend(ctx, &from, &to, "user", "create").
			Return(nil, errors.New("database error"))

		trend, err := svc.AuditTrend(ctx, &from, &to, "user", "create")
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if trend != nil {
			t.Error("Expected nil trend on error")
		}
	})
}
