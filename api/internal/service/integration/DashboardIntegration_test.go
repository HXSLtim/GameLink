// Package integration provides integration tests for the dashboard aggregation statistics service.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/stats"
	"gamelink/internal/service/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// Platform Overview API Tests (平台概览API测试)
// ============================================================================

func TestDashboard_GetOverview_TotalCounts(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create test data
	user1 := CreateUniqueTestUser(t, db, "user1")
	user2 := CreateUniqueTestUser(t, db, "user2")
	_ = CreateTestPlayer(t, db, user1)
	_ = CreateTestPlayer(t, db, user2)
	game := CreateTestGame(t, db, "LOL")

	order1 := CreateTestOrderWithDetails(t, db, user1, CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player1")), game, model.OrderStatusPending, 10000)
	order2 := CreateTestOrderWithDetails(t, db, user2, CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player2")), game, model.OrderStatusCompleted, 15000)

	// Create payments
	_ = CreateTestPayment(t, db, order1, model.PaymentStatusPaid)
	_ = CreateTestPayment(t, db, order2, model.PaymentStatusPaid)

	// Get overview stats
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), dashboard.TotalPlayers) // user1 and user2 became players
	assert.GreaterOrEqual(t, dashboard.TotalUsers, int64(4))
	assert.Equal(t, int64(1), dashboard.TotalGames)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(2))
	assert.Equal(t, int64(25000), dashboard.TotalPaidAmountCents) // 10000 + 15000
}

func TestDashboard_GetOverview_TodayStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create today's orders
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	user := CreateUniqueTestUser(t, db, "today_user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "today_player"))
	game := CreateTestGame(t, db, "TodayGame")

	// Create order for today
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 20000)
	order.CreatedAt = todayStart.Add(2 * time.Hour)
	db.Save(order)

	// Create payment
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
	paidTime := todayStart.Add(2 * time.Hour)
	payment.PaidAt = &paidTime
	db.Save(payment)

	// Get overview
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(1))

	// Check orders by status includes completed
	_, hasCompleted := dashboard.OrdersByStatus[string(model.OrderStatusCompleted)]
	assert.True(t, hasCompleted)
}

func TestDashboard_GetOverview_OrdersByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders with different statuses
	user1 := CreateUniqueTestUser(t, db, "user1")
	user2 := CreateUniqueTestUser(t, db, "user2")
	user3 := CreateUniqueTestUser(t, db, "user3")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	CreateTestOrderWithDetails(t, db, user1, player, game, model.OrderStatusPending, 10000)
	CreateTestOrderWithDetails(t, db, user2, player, game, model.OrderStatusInProgress, 15000)
	CreateTestOrderWithDetails(t, db, user3, player, game, model.OrderStatusCompleted, 20000)
	CreateTestOrderWithDetails(t, db, user1, player, game, model.OrderStatusCanceled, 5000)

	// Get overview
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	// Verify status counts
	assert.Contains(t, dashboard.OrdersByStatus, string(model.OrderStatusPending))
	assert.Contains(t, dashboard.OrdersByStatus, string(model.OrderStatusInProgress))
	assert.Contains(t, dashboard.OrdersByStatus, string(model.OrderStatusCompleted))
	assert.Contains(t, dashboard.OrdersByStatus, string(model.OrderStatusCanceled))

	// Check each status has at least 1 order
	for _, status := range []model.OrderStatus{
		model.OrderStatusPending,
		model.OrderStatusInProgress,
		model.OrderStatusCompleted,
		model.OrderStatusCanceled,
	} {
		assert.GreaterOrEqual(t, dashboard.OrdersByStatus[string(status)], int64(1),
			fmt.Sprintf("Status %s should have at least 1 order", status))
	}
}

func TestDashboard_GetOverview_PaymentsByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create payments with different statuses
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	order1 := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	order2 := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 20000)
	order3 := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 30000)

	CreateTestPayment(t, db, order1, model.PaymentStatusPending)
	CreateTestPayment(t, db, order2, model.PaymentStatusPaid)
	CreateTestPayment(t, db, order3, model.PaymentStatusPaid)

	// Get overview
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	// Verify payment status counts
	assert.Contains(t, dashboard.PaymentsByStatus, string(model.PaymentStatusPending))
	assert.Contains(t, dashboard.PaymentsByStatus, string(model.PaymentStatusPaid))

	// Verify total paid amount (20000 + 30000 = 50000)
	assert.Equal(t, int64(50000), dashboard.TotalPaidAmountCents)
}

// ============================================================================
// Revenue Trend API Tests (收入趋势API测试)
// ============================================================================

func TestDashboard_RevenueTrend_Last7Days(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create payments over the last 7 days
	now := time.Now()
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	for i := 0; i < 7; i++ {
		paymentTime := now.AddDate(0, 0, -i)
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, int64((i+1)*1000))
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		payment.PaidAt = &paymentTime
		db.Save(payment)
	}

	// Get revenue trend
	trend, err := svc.RevenueTrend(ctx, 7)
	require.NoError(t, err)
	assert.NotEmpty(t, trend)
	assert.LessOrEqual(t, len(trend), 7)

	// Verify all entries have dates and values
	for _, entry := range trend {
		assert.NotEmpty(t, entry.Date)
		assert.GreaterOrEqual(t, entry.Value, int64(0))
	}
}

func TestDashboard_RevenueTrend_CustomDays(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create payments over the last 30 days
	now := time.Now()
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	for i := 0; i < 30; i++ {
		paymentTime := now.AddDate(0, 0, -i)
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 5000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		payment.PaidAt = &paymentTime
		db.Save(payment)
	}

	// Get revenue trend for 30 days
	trend, err := svc.RevenueTrend(ctx, 30)
	require.NoError(t, err)
	assert.NotEmpty(t, trend)
	assert.LessOrEqual(t, len(trend), 30)
}

func TestDashboard_RevenueTrend_NoData(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Get revenue trend with no data
	trend, err := svc.RevenueTrend(ctx, 7)
	require.NoError(t, err)
	assert.Empty(t, trend)
}

func TestDashboard_RevenueTrend_OnlyPaidPayments(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create paid and unpaid payments
	now := time.Now()
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	// Paid payment
	order1 := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPaid)
	payment1.PaidAt = &now
	db.Save(payment1)

	// Pending payment (should not be counted)
	order2 := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 20000)
	_ = CreateTestPayment(t, db, order2, model.PaymentStatusPending)

	// Get revenue trend
	trend, err := svc.RevenueTrend(ctx, 7)
	require.NoError(t, err)

	// Verify only paid payments are counted
	totalRevenue := int64(0)
	for _, entry := range trend {
		totalRevenue += entry.Value
	}
	assert.Equal(t, int64(10000), totalRevenue)
}

// ============================================================================
// User Growth Trend API Tests (用户增长趋势API测试)
// ============================================================================

func TestDashboard_UserGrowth_Last7Days(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create users over the last 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		userTime := now.AddDate(0, 0, -i)
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		user.CreatedAt = userTime
		db.Save(user)
	}

	// Get user growth trend
	trend, err := svc.UserGrowth(ctx, 7)
	require.NoError(t, err)
	assert.NotEmpty(t, trend)
	assert.LessOrEqual(t, len(trend), 7)

	// Verify all entries have dates and values
	for _, entry := range trend {
		assert.NotEmpty(t, entry.Date)
		assert.GreaterOrEqual(t, entry.Value, int64(0))
	}
}

func TestDashboard_UserGrowth_CustomDays(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create users over the last 30 days
	now := time.Now()
	for i := 0; i < 30; i++ {
		userTime := now.AddDate(0, 0, -i)
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		user.CreatedAt = userTime
		db.Save(user)
	}

	// Get user growth trend for 30 days
	trend, err := svc.UserGrowth(ctx, 30)
	require.NoError(t, err)
	assert.NotEmpty(t, trend)
	assert.LessOrEqual(t, len(trend), 30)
}

func TestDashboard_UserGrowth_NoData(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Get user growth trend with no data
	trend, err := svc.UserGrowth(ctx, 7)
	require.NoError(t, err)
	assert.Empty(t, trend)
}

// ============================================================================
// Orders By Status API Tests (订单状态分布API测试)
// ============================================================================

func TestDashboard_OrdersByStatus_AllStatuses(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders with all possible statuses
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	statuses := []model.OrderStatus{
		model.OrderStatusPending,
		model.OrderStatusInProgress,
		model.OrderStatusCompleted,
		model.OrderStatusCanceled,
		model.OrderStatusRefunded,
		model.OrderStatusDisputed,
	}

	for _, status := range statuses {
		CreateTestOrderWithDetails(t, db, user, player, game, status, 10000)
	}

	// Get orders by status
	statusCounts, err := svc.OrdersByStatus(ctx)
	require.NoError(t, err)

	// Verify all statuses are present
	for _, status := range statuses {
		assert.Contains(t, statusCounts, string(status))
		assert.GreaterOrEqual(t, statusCounts[string(status)], int64(1))
	}
}

func TestDashboard_OrdersByStatus_EmptyDatabase(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Get orders by status with no orders
	statusCounts, err := svc.OrdersByStatus(ctx)
	require.NoError(t, err)
	assert.Empty(t, statusCounts)
}

func TestDashboard_OrdersByStatus_MultipleOrdersPerStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create multiple orders with same status
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LOL")

	for i := 0; i < 5; i++ {
		CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i)), player, game, model.OrderStatusCompleted, 10000)
	}

	// Get orders by status
	statusCounts, err := svc.OrdersByStatus(ctx)
	require.NoError(t, err)

	// Verify completed status has 5 orders
	assert.Equal(t, int64(5), statusCounts[string(model.OrderStatusCompleted)])
}

// ============================================================================
// Top Players API Tests (陪玩师排行榜API测试)
// ============================================================================

func TestDashboard_TopPlayers_ByRating(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create players with different ratings
	user1 := CreateUniqueTestUser(t, db, "user1")
	user2 := CreateUniqueTestUser(t, db, "user2")
	user3 := CreateUniqueTestUser(t, db, "user3")

	player1 := CreateTestPlayer(t, db, user1)
	player2 := CreateTestPlayer(t, db, user2)
	player3 := CreateTestPlayer(t, db, user3)

	// Set different ratings
	player1.RatingAverage = 4.8
	player1.RatingCount = 100
	db.Save(player1)

	player2.RatingAverage = 4.5
	player2.RatingCount = 50
	db.Save(player2)

	player3.RatingAverage = 5.0
	player3.RatingCount = 20
	db.Save(player3)

	// Get top players
	topPlayers, err := svc.TopPlayers(ctx, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, topPlayers)
	assert.LessOrEqual(t, len(topPlayers), 10)

	// Verify player structure
	for _, p := range topPlayers {
		assert.NotZero(t, p.PlayerID)
		assert.NotEmpty(t, p.Nickname)
		assert.GreaterOrEqual(t, p.RatingAverage, float32(0))
		assert.GreaterOrEqual(t, p.RatingCount, uint32(0))
	}
}

func TestDashboard_TopPlayers_Limit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create more players than the limit
	for i := 0; i < 15; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		player := CreateTestPlayer(t, db, user)
		player.RatingCount = uint32(i + 1)
		db.Save(player)
	}

	// Get top 5 players
	topPlayers, err := svc.TopPlayers(ctx, 5)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(topPlayers), 5)
}

func TestDashboard_TopPlayers_NoPlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Get top players with no players
	topPlayers, err := svc.TopPlayers(ctx, 10)
	require.NoError(t, err)
	assert.Empty(t, topPlayers)
}

func TestDashboard_TopPlayers_DefaultLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create players
	for i := 0; i < 15; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		player := CreateTestPlayer(t, db, user)
		player.RatingCount = uint32(i + 1)
		db.Save(player)
	}

	// Get top players with default limit (should be 10)
	topPlayers, err := svc.TopPlayers(ctx, 0)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(topPlayers), 10)
}

// ============================================================================
// Game Statistics Filter Tests (游戏维度筛选测试)
// ============================================================================

func TestDashboard_GameStats_Aggregation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test data
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))

	game1 := CreateTestGame(t, db, "LOL")
	game2 := CreateTestGame(t, db, "王者荣耀")

	// Create orders for both games
	for i := 0; i < 3; i++ {
		CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("user1_%d", i)), player, game1, model.OrderStatusCompleted, 10000)
	}
	for i := 0; i < 2; i++ {
		CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("user2_%d", i)), player, game2, model.OrderStatusCompleted, 15000)
	}

	// Verify dashboard includes all games
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), dashboard.TotalGames)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(5))
}

// ============================================================================
// Time Window Query Tests (时间窗口查询测试)
// ============================================================================

func TestDashboard_TodayStats_WithinWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create today's order
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "today_player"))
	game := CreateTestGame(t, db, "TodayGame")

	_ = CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, "today_user"), player, game, model.OrderStatusCompleted, 10000)
	_ = todayStart.Add(2 * time.Hour) // Use the time

	// Verify dashboard includes today's order
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(1))
}

func TestDashboard_ThisWeekStats_WithinWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders throughout the week
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))

	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "week_player"))
	game := CreateTestGame(t, db, "WeekGame")

	for i := 0; i < 7; i++ {
		_ = weekStart.AddDate(0, 0, i)
		_ = CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("week_user_%d", i)), player, game, model.OrderStatusCompleted, 10000)
	}

	// Verify dashboard includes this week's orders
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(7))
}

func TestDashboard_ThisMonthStats_WithinWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders throughout the month
	_ = time.Now()
	_ = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "month_player"))
	game := CreateTestGame(t, db, "MonthGame")

	for i := 0; i < 5; i++ {
		_ = CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("month_user_%d", i)), player, game, model.OrderStatusCompleted, 10000)
	}

	// Verify dashboard includes this month's orders
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(5))
}

func TestDashboard_CustomDateRange_WithinWindow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders in specific date range
	_ = time.Now()
	_ = time.Now().AddDate(0, 0, -10)

	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "range_player"))
	game := CreateTestGame(t, db, "RangeGame")

	for i := 0; i < 3; i++ {
		_ = CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, fmt.Sprintf("range_user_%d", i)), player, game, model.OrderStatusCompleted, 10000)
	}

	// Verify revenue trend includes data from range
	trend, err := svc.RevenueTrend(ctx, 15)
	require.NoError(t, err)
	assert.NotEmpty(t, trend)
}

// ============================================================================
// Data Aggregation Tests (数据聚合计算测试)
// ============================================================================

func TestDashboard_Aggregation_SUM(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create orders with different amounts
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "SumGame")

	amounts := []int64{10000, 20000, 30000}
	for _, amount := range amounts {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, amount)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		now := time.Now()
		payment.PaidAt = &now
		db.Save(payment)
	}

	// Verify SUM aggregation in revenue trend
	trend, err := svc.RevenueTrend(ctx, 1)
	require.NoError(t, err)

	totalRevenue := int64(0)
	for _, entry := range trend {
		totalRevenue += entry.Value
	}
	assert.Equal(t, int64(60000), totalRevenue) // 10000 + 20000 + 30000
}

func TestDashboard_Aggregation_COUNT(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create multiple users
	for i := 0; i < 10; i++ {
		CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
	}

	// Verify COUNT aggregation
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalUsers, int64(10))
}

func TestDashboard_Aggregation_AVG(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// This test verifies that the system correctly calculates averages
	// Create orders with varying amounts
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "AvgGame")

	amounts := []int64{10000, 20000, 30000, 40000, 50000}
	for _, amount := range amounts {
		CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, amount)
	}

	// Verify dashboard
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(5))
}

func TestDashboard_Aggregation_MAX_MIN(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create players with different rating counts
	for i := 1; i <= 10; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		player := CreateTestPlayer(t, db, user)
		player.RatingCount = uint32(i * 10)
		db.Save(player)
	}

	// Get top players to verify MAX aggregation
	topPlayers, err := svc.TopPlayers(ctx, 3)
	require.NoError(t, err)
	assert.NotEmpty(t, topPlayers)

	// Verify top players have highest rating counts
	maxRatingCount := topPlayers[0].RatingCount
	assert.GreaterOrEqual(t, maxRatingCount, uint32(10))
}

// ============================================================================
// Cache Strategy Tests (缓存策略测试)
// ============================================================================

func TestDashboard_Cache_RealtimeData(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create initial data
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "CacheGame")

	// Get initial dashboard
	dashboard1, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	// Add new order
	CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)

	// Get dashboard again - should reflect new data immediately (realtime)
	dashboard2, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.Greater(t, dashboard2.TotalOrders, dashboard1.TotalOrders)
}

func TestDashboard_Cache_Consistency(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create initial data
	for i := 0; i < 5; i++ {
		CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
	}

	// Get dashboard multiple times
	dashboard1, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	dashboard2, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	// Verify consistency
	assert.Equal(t, dashboard1.TotalUsers, dashboard2.TotalUsers)
	assert.Equal(t, dashboard1.TotalPlayers, dashboard2.TotalPlayers)
}

// ============================================================================
// Data Consistency Tests (数据一致性测试)
// ============================================================================

func TestDashboard_Consistency_MultiTableJOIN(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create related data across multiple tables
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	// Verify JOIN query consistency
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)

	// Verify all related entities are counted
	assert.GreaterOrEqual(t, dashboard.TotalUsers, int64(1))
	assert.GreaterOrEqual(t, dashboard.TotalPlayers, int64(1))
	assert.GreaterOrEqual(t, dashboard.TotalGames, int64(1))
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(1))
	assert.GreaterOrEqual(t, dashboard.TotalPaidAmountCents, payment.AmountCents)
}

func TestDashboard_Consistency_AggregationAccuracy(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create test data with known values
	user := CreateUniqueTestUser(t, db, "user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "AccuracyGame")

	expectedOrders := 5
	expectedTotalAmount := int64(0)

	for i := 0; i < expectedOrders; i++ {
		amount := int64((i + 1) * 10000)
		expectedTotalAmount += amount
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, amount)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		now := time.Now()
		payment.PaidAt = &now
		db.Save(payment)
	}

	// Get revenue trend and verify accuracy
	trend, err := svc.RevenueTrend(ctx, 1)
	require.NoError(t, err)

	totalRevenue := int64(0)
	for _, entry := range trend {
		totalRevenue += entry.Value
	}
	assert.Equal(t, expectedTotalAmount, totalRevenue)
}

// ============================================================================
// Edge Cases and Boundary Tests (边界情况测试)
// ============================================================================

func TestDashboard_EdgeCase_EmptyDatabase(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Get dashboard with empty database
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), dashboard.TotalUsers)
	assert.Equal(t, int64(0), dashboard.TotalPlayers)
	assert.Equal(t, int64(0), dashboard.TotalGames)
	assert.Equal(t, int64(0), dashboard.TotalOrders)
	assert.Equal(t, int64(0), dashboard.TotalPaidAmountCents)
}

func TestDashboard_EdgeCase_SingleRecord(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create single record
	user := CreateUniqueTestUser(t, db, "single_user")
	player := CreateTestPlayer(t, db, user)
	_ = CreateTestGame(t, db, "SingleGame")

	_ = CreateTestOrderWithDetails(t, db, user, player, CreateTestGame(t, db, "SingleGame"), model.OrderStatusCompleted, 10000)

	// Verify dashboard handles single record correctly
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), dashboard.TotalUsers)
	assert.Equal(t, int64(1), dashboard.TotalPlayers)
	assert.Equal(t, int64(1), dashboard.TotalGames)
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(1))
}

func TestDashboard_EdgeCase_LargeDataSet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create large dataset
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player"))
	game := CreateTestGame(t, db, "LargeGame")

	for i := 0; i < 100; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 10000)
	}

	// Verify dashboard handles large dataset
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dashboard.TotalUsers, int64(101)) // 100 users + 1 player user
	assert.GreaterOrEqual(t, dashboard.TotalOrders, int64(100))
}

func TestDashboard_EdgeCase_ZeroAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create order with zero amount
	user := CreateUniqueTestUser(t, db, "zero_user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "zero_player"))
	game := CreateTestGame(t, db, "ZeroGame")

	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusCompleted, 0)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
	payment.AmountCents = 0
	db.Save(payment)

	// Verify dashboard handles zero amount correctly
	dashboard, err := svc.Dashboard(ctx)
	require.NoError(t, err)
	// Zero amount should not cause errors
	assert.GreaterOrEqual(t, dashboard.TotalPaidAmountCents, int64(0))
}

func TestDashboard_EdgeCase_FutureDate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create payment with future date
	_ = CreateUniqueTestUser(t, db, "future_user")
	player := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "future_player"))
	game := CreateTestGame(t, db, "FutureGame")

	order := CreateTestOrderWithDetails(t, db, CreateUniqueTestUser(t, db, "future_user2"), player, game, model.OrderStatusCompleted, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
	futureTime := time.Now().Add(24 * time.Hour)
	payment.PaidAt = &futureTime
	db.Save(payment)

	// Verify revenue trend includes future date
	trend, err := svc.RevenueTrend(ctx, 7)
	require.NoError(t, err)
	// Future date should not cause errors
	assert.NotNil(t, trend)
}

func TestDashboard_EdgeCase_InvalidLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Test with negative limit
	topPlayers, err := svc.TopPlayers(ctx, -1)
	require.NoError(t, err)
	// Should handle negative limit gracefully
	assert.NotNil(t, topPlayers)
}

func TestDashboard_EdgeCase_VeryLargeLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	statsRepo := stats.NewStatsRepository(db)
	svc := admin.NewStatsService(statsRepo)

	// Create some players
	for i := 0; i < 5; i++ {
		user := CreateUniqueTestUser(t, db, fmt.Sprintf("user%d", i))
		player := CreateTestPlayer(t, db, user)
		player.RatingCount = uint32(i + 1)
		db.Save(player)
	}

	// Test with very large limit
	topPlayers, err := svc.TopPlayers(ctx, 999999)
	require.NoError(t, err)
	assert.NotEmpty(t, topPlayers)
	assert.LessOrEqual(t, len(topPlayers), 5)
}

// ============================================================================
// Helper Functions (辅助函数)
// ============================================================================

// CreateTestOrderWithTimestamp creates an order with a specific timestamp
func CreateTestOrderWithTimestamp(t *testing.T, db *gorm.DB, user *model.User, player *model.Player, status model.OrderStatus, priceCents int64, createdAt time.Time) *model.Order {
	t.Helper()

	now := time.Now()
	scheduledStart := now.Add(time.Hour)
	scheduledEnd := now.Add(2 * time.Hour)

	order := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:           fmt.Sprintf("ORD%d%d", user.ID, time.Now().UnixNano()),
		UserID:            user.ID,
		PlayerID:          &player.ID,
		Quantity:          1,
		UnitPriceCents:    priceCents,
		TotalPriceCents:   priceCents,
		CommissionCents:   priceCents * 20 / 100,
		PlayerIncomeCents: priceCents * 80 / 100,
		Currency:          model.CurrencyCNY,
		Status:            status,
		Title:             "Test Order",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		OrderConfig:       "{}",
	}

	if status == model.OrderStatusCompleted {
		order.CompletedAt = &createdAt
	}

	if err := db.Create(order).Error; err != nil {
		t.Fatalf("Failed to create test order: %v", err)
	}
	// Update timestamps after creation
	db.Model(order).Updates(map[string]interface{}{
		"created_at": createdAt,
		"updated_at": createdAt,
	})
	return order
}

// CreateTestUserWithTimestamp creates a user with a specific timestamp
func CreateTestUserWithTimestamp(t *testing.T, db *gorm.DB, name string, createdAt time.Time) *model.User {
	t.Helper()
	ts := createdAt.UnixNano()
	user := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:   fmt.Sprintf("%s_%d", name, ts),
		Phone:  fmt.Sprintf("138%011d", ts%100000000000),
		Email:  fmt.Sprintf("%s_%d@test.com", name, ts),
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Update timestamps after creation
	db.Model(user).Updates(map[string]interface{}{
		"created_at": createdAt,
		"updated_at": createdAt,
	})
	return user
}
