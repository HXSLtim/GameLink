// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/service/kpi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKPIService_GetOverview(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create test data
	user := CreateTestUser(t, db, "kpi_user")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "kpi_player"))

	// Create orders
	for i := 0; i < 5; i++ {
		order := CreateTestOrder(t, db, user, player, model.OrderStatusCompleted)
		order.TotalPriceCents = 10000
		db.Save(order)

		// Create payment
		CreateTestPayment(t, db, order, model.PaymentStatusPaid)
	}

	// Get overview
	params := kpi.QueryParams{
		Period:  kpi.PeriodMonth,
		Compare: kpi.CompareMoM,
	}
	overview, err := svc.GetOverview(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, overview)

	// Verify GMV
	assert.Equal(t, "GMV", overview.GMV.Name)
	assert.Equal(t, float64(500), overview.GMV.CurrentValue) // 5 * 100 yuan

	// Verify Orders
	assert.Equal(t, "订单数", overview.Orders.Name)
	assert.Equal(t, float64(5), overview.Orders.CurrentValue)

	// Verify New Users
	assert.Equal(t, "新用户", overview.NewUsers.Name)
	assert.GreaterOrEqual(t, overview.NewUsers.CurrentValue, float64(2)) // At least 2 users created
}

func TestKPIService_GetTrend(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create test data
	user := CreateTestUser(t, db, "trend_user")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "trend_player"))

	// Create orders for the past week
	for i := 0; i < 7; i++ {
		order := &model.Order{
			OrderNo:         "TREND" + string(rune('A'+i)),
			UserID:          user.ID,
			PlayerID:        &player.ID,
			TotalPriceCents: 10000,
			Status:          model.OrderStatusCompleted,
			Currency:        model.CurrencyCNY,
		}
		order.CreatedAt = time.Now().AddDate(0, 0, -i)
		db.Create(order)
	}

	// Get trend
	params := kpi.QueryParams{
		Period: kpi.PeriodWeek,
	}
	trend, err := svc.GetTrend(ctx, "orders", params)
	require.NoError(t, err)
	require.NotNil(t, trend)
	assert.NotEmpty(t, trend)

	// Verify trend points have dates
	for _, point := range trend {
		assert.NotEmpty(t, point.Date)
	}
}

func TestKPIService_GetTargets(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create KPI targets
	now := time.Now()
	targets := []*model.KPITarget{
		{
			MetricName:  "gmv",
			PeriodType:  "monthly",
			TargetValue: 100000,
			StartDate:   now.AddDate(0, 0, -30),
			EndDate:     now.AddDate(0, 0, 30),
		},
		{
			MetricName:  "orders",
			PeriodType:  "monthly",
			TargetValue: 1000,
			StartDate:   now.AddDate(0, 0, -30),
			EndDate:     now.AddDate(0, 0, 30),
		},
	}
	for _, target := range targets {
		db.Create(target)
	}

	// Get all targets
	result, err := svc.GetTargets(ctx, "", "")
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Get targets by period type
	result, err = svc.GetTargets(ctx, "monthly", "")
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Get targets by metric name
	result, err = svc.GetTargets(ctx, "", "gmv")
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "gmv", result[0].MetricName)
}

func TestKPIService_CreateTarget(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	now := time.Now()
	target := &model.KPITarget{
		MetricName:  "new_users",
		PeriodType:  "weekly",
		TargetValue: 500,
		StartDate:   now,
		EndDate:     now.AddDate(0, 0, 7),
	}

	err := svc.CreateTarget(ctx, target)
	require.NoError(t, err)
	assert.NotZero(t, target.ID)

	// Verify in database
	var saved model.KPITarget
	err = db.First(&saved, target.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "new_users", saved.MetricName)
	assert.Equal(t, float64(500), saved.TargetValue)
}

func TestKPIService_UpdateTarget(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create target
	now := time.Now()
	target := &model.KPITarget{
		MetricName:  "dau",
		PeriodType:  "daily",
		TargetValue: 1000,
		StartDate:   now,
		EndDate:     now.AddDate(0, 0, 1),
	}
	db.Create(target)

	// Update target
	err := svc.UpdateTarget(ctx, uint(target.ID), map[string]interface{}{
		"target_value": 2000,
	})
	require.NoError(t, err)

	// Verify update
	var updated model.KPITarget
	err = db.First(&updated, target.ID).Error
	require.NoError(t, err)
	assert.Equal(t, float64(2000), updated.TargetValue)
}

func TestKPIService_DeleteTarget(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create target
	now := time.Now()
	target := &model.KPITarget{
		MetricName:  "retention",
		PeriodType:  "monthly",
		TargetValue: 50,
		StartDate:   now,
		EndDate:     now.AddDate(0, 1, 0),
	}
	db.Create(target)

	// Delete target
	err := svc.DeleteTarget(ctx, uint(target.ID))
	require.NoError(t, err)

	// Verify deletion
	var deleted model.KPITarget
	err = db.First(&deleted, target.ID).Error
	assert.Error(t, err) // Should not find
}

func TestKPIService_CalculateChangeRateAndTrend(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := kpi.NewKPIService(db)
	ctx := context.Background()

	// Create historical data (previous month)
	user := CreateTestUser(t, db, "change_user")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "change_player"))

	// Create orders for previous month
	lastMonth := time.Now().AddDate(0, -1, 0)
	for i := 0; i < 3; i++ {
		order := &model.Order{
			OrderNo:         "PREV" + string(rune('A'+i)),
			UserID:          user.ID,
			PlayerID:        &player.ID,
			TotalPriceCents: 10000,
			Status:          model.OrderStatusCompleted,
			Currency:        model.CurrencyCNY,
		}
		order.CreatedAt = lastMonth
		db.Create(order)

		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      user.ID,
			AmountCents: 10000,
			Method:      model.PaymentMethodWeChat,
			Status:      model.PaymentStatusPaid,
		}
		payment.CreatedAt = lastMonth
		db.Create(payment)
	}

	// Create orders for current month
	for i := 0; i < 5; i++ {
		order := CreateTestOrder(t, db, user, player, model.OrderStatusCompleted)
		order.TotalPriceCents = 10000
		db.Save(order)
		CreateTestPayment(t, db, order, model.PaymentStatusPaid)
	}

	// Get overview with MoM comparison
	params := kpi.QueryParams{
		Period:  kpi.PeriodMonth,
		Compare: kpi.CompareMoM,
	}
	overview, err := svc.GetOverview(ctx, params)
	require.NoError(t, err)

	// Current should be higher than previous
	assert.Greater(t, overview.Orders.CurrentValue, overview.Orders.PreviousValue)
	// Trend should be "up"
	assert.Equal(t, "up", overview.Orders.Trend)
}
