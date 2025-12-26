// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/service/statistics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatisticsService_UpdateUserStatistics(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := statistics.NewService(db)
	ctx := context.Background()

	// Create test user
	user := CreateTestUser(t, db, "stats_user")

	// Create test player
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "stats_player"))

	// Create completed orders
	for i := 0; i < 3; i++ {
		order := CreateTestOrder(t, db, user, player, model.OrderStatusCompleted)
		order.TotalPriceCents = 5000 // 50 yuan each
		db.Save(order)
	}

	// Create canceled order
	CreateTestOrder(t, db, user, player, model.OrderStatusCanceled)

	// Update statistics
	err := svc.UpdateUserStatistics(ctx, user.ID)
	require.NoError(t, err)

	// Verify statistics
	var stats model.UserStatistics
	err = db.Where("user_id = ?", user.ID).First(&stats).Error
	require.NoError(t, err)

	assert.Equal(t, user.ID, stats.UserID)
	assert.Equal(t, 4, stats.TotalOrderCount)
	assert.Equal(t, 3, stats.CompletedOrderCount)
	assert.Equal(t, 1, stats.CanceledOrderCount)
}

func TestStatisticsService_UpdatePlayerStatistics(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := statistics.NewService(db)
	ctx := context.Background()

	// Create test user and player
	user := CreateTestUser(t, db, "player_stats_user")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "player_stats_player"))

	// Create completed orders for the player
	for i := 0; i < 5; i++ {
		order := &model.Order{
			OrderNo:           "PSTATS" + string(rune('A'+i)),
			UserID:            user.ID,
			PlayerID:          &player.ID,
			TotalPriceCents:   10000,
			PlayerIncomeCents: 8000,
			CommissionCents:   2000,
			Status:            model.OrderStatusCompleted,
			Currency:          model.CurrencyCNY,
		}
		db.Create(order)
	}

	// Update statistics
	err := svc.UpdatePlayerStatistics(ctx, player.ID)
	require.NoError(t, err)

	// Verify statistics
	var stats model.PlayerStatistics
	err = db.Where("player_id = ?", player.ID).First(&stats).Error
	require.NoError(t, err)

	assert.Equal(t, player.ID, stats.PlayerID)
	assert.Equal(t, 5, stats.TotalOrderCount)
	assert.Equal(t, 5, stats.CompletedOrderCount)
	assert.Equal(t, int64(40000), stats.TotalEarningsCents) // 5 * 8000
}

func TestStatisticsService_UpdatePlatformDailyStatistics(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := statistics.NewService(db)
	ctx := context.Background()

	// Create test data for today
	today := time.Now().Truncate(24 * time.Hour)

	// Create users
	user1 := CreateTestUser(t, db, "platform_user1")
	user2 := CreateTestUser(t, db, "platform_user2")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "platform_player"))

	// Create orders for today
	for i, user := range []*model.User{user1, user2} {
		order := &model.Order{
			OrderNo:         "PLATFORM" + string(rune('A'+i)),
			UserID:          user.ID,
			PlayerID:        &player.ID,
			TotalPriceCents: 10000,
			CommissionCents: 2000,
			Status:          model.OrderStatusCompleted,
			Currency:        model.CurrencyCNY,
		}
		order.CreatedAt = today.Add(time.Duration(i) * time.Hour)
		db.Create(order)
	}

	// Update platform statistics
	err := svc.UpdatePlatformDailyStatistics(ctx, today)
	require.NoError(t, err)

	// Verify statistics
	var stats model.PlatformStatistics
	err = db.Where("stat_date = ?", today).First(&stats).Error
	require.NoError(t, err)

	assert.Equal(t, 2, stats.DailyOrderCount)
	assert.Equal(t, 2, stats.DailyCompletedCount)
	assert.Equal(t, int64(20000), stats.DailyGMVCents)
}

func TestEventHooks_OnOrderCompleted(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	svc := statistics.NewService(db)
	hooks := statistics.NewEventHooks(svc)
	ctx := context.Background()

	// Create test data
	user := CreateTestUser(t, db, "hook_user")
	player := CreateTestPlayer(t, db, CreateTestUser(t, db, "hook_player"))

	// Create service item
	item := &model.ServiceItem{
		ItemCode:       "HOOK_ITEM",
		Name:           "Test Item",
		Category:       "escort",
		BasePriceCents: 5000,
		IsActive:       true,
	}
	db.Create(item)

	// Create order
	order := &model.Order{
		OrderNo:         "HOOK_ORDER",
		UserID:          user.ID,
		PlayerID:        &player.ID,
		ItemID:          item.ID,
		TotalPriceCents: 10000,
		Status:          model.OrderStatusCompleted,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	// Trigger hook (runs in goroutine, so we need to wait)
	hooks.OnOrderCompleted(ctx, order)

	// Wait for async processing
	time.Sleep(500 * time.Millisecond)

	// Verify user statistics were updated
	var userStats model.UserStatistics
	err := db.Where("user_id = ?", user.ID).First(&userStats).Error
	require.NoError(t, err)
	assert.Equal(t, 1, userStats.TotalOrderCount)

	// Verify player statistics were updated
	var playerStats model.PlayerStatistics
	err = db.Where("player_id = ?", player.ID).First(&playerStats).Error
	require.NoError(t, err)
	assert.Equal(t, 1, playerStats.TotalOrderCount)
}

func TestTagEvaluator_EvaluateUserTags(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test user
	user := CreateTestUser(t, db, "tag_user")

	// Create user statistics
	stats := &model.UserStatistics{
		UserID:              user.ID,
		TotalOrderCount:     100,
		CompletedOrderCount: 95,
		TotalSpentCents:     1000000, // 10000 yuan
	}
	db.Create(stats)

	// Create tag threshold (if model exists)
	// This test verifies the evaluator can query without errors
	evaluator := statistics.NewTagEvaluator(db)
	tagIDs, err := evaluator.EvaluateUserTags(ctx, user.ID)

	// Should not error even if no thresholds configured
	require.NoError(t, err)
	// May return empty if no thresholds match
	assert.NotNil(t, tagIDs)
}
