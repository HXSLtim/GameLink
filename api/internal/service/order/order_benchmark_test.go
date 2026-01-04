package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/repository"
	"gamelink/internal/service/integration"
)

// BenchmarkGetMyOrders_WithBatchQuery 测试批量查询的性能 (优化后)
func BenchmarkGetMyOrders_WithBatchQuery(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	// Setup test service
	orderRepo := repository.NewOrderRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	commissionRepo := repository.NewCommissionRepository(db)

	orderService := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)

	ctx := context.Background()

	// Create test user
	user := integration.CreateTestUser(b, db, "benchuser")
	game := integration.CreateTestGame(b, db, "Benchmark Game")
	player := integration.CreateTestPlayer(b, db, user.ID)

	// Create multiple orders for benchmarking
	numOrders := 50
	for i := 0; i < numOrders; i++ {
		scheduledStart := time.Now().Add(time.Hour * time.Duration(i))
		scheduledEnd := scheduledStart.Add(time.Hour * 2)

		order := &model.Order{
			OrderNo:         fmt.Sprintf("BNCH-%d", i),
			UserID:          user.ID,
			GameID:          &game.ID,
			PlayerID:        &player.ID,
			Title:           fmt.Sprintf("Benchmark Order %d", i),
			Description:     "Test order for benchmarking",
			TotalPriceCents: 5000,
			Status:          model.OrderStatusCompleted,
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		if err := db.Create(order).Error; err != nil {
			b.Fatalf("Failed to create test order: %v", err)
		}
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := MyOrderListRequest{
			Page:     1,
			PageSize: numOrders,
		}
		_, err := orderService.GetMyOrders(ctx, user.ID, req)
		if err != nil {
			b.Fatalf("GetMyOrders failed: %v", err)
		}
	}
}

// BenchmarkGetMyOrders_WithoutBatchQuery 模拟 N+1 查询的性能 (优化前)
// 这个测试手动模拟 N+1 查询，用于对比性能
func BenchmarkGetMyOrders_WithoutBatchQuery(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	ctx := context.Background()

	// Create test user
	user := integration.CreateTestUser(b, db, "benchuser2")
	game := integration.CreateTestGame(b, db, "Benchmark Game 2")
	player := integration.CreateTestPlayer(b, db, user.ID)

	// Create multiple orders for benchmarking
	numOrders := 50
	orders := make([]*model.Order, numOrders)
	for i := 0; i < numOrders; i++ {
		scheduledStart := time.Now().Add(time.Hour * time.Duration(i))
		scheduledEnd := scheduledStart.Add(time.Hour * 2)

		order := &model.Order{
			OrderNo:         fmt.Sprintf("BNCH2-%d", i),
			UserID:          user.ID,
			GameID:          &game.ID,
			PlayerID:        &player.ID,
			Title:           fmt.Sprintf("Benchmark Order %d", i),
			Description:     "Test order for benchmarking",
			TotalPriceCents: 5000,
			Status:          model.OrderStatusCompleted,
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		if err := db.Create(order).Error; err != nil {
			b.Fatalf("Failed to create test order: %v", err)
		}
		orders[i] = order
	}

	playerRepo := repository.NewPlayerRepository(db)
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)

	// Run benchmark - simulating N+1 queries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟 N+1 查询：对每个订单都查询关联数据
		for _, order := range orders {
			// N+1 查询 1: 查询陪玩师
			_, _ = playerRepo.Get(ctx, *order.PlayerID)

			// N+1 查询 2: 查询用户
			player, _ := playerRepo.Get(ctx, *order.PlayerID)
			if player != nil {
				_, _ = userRepo.Get(ctx, player.UserID)
			}

			// N+1 查询 3: 查询游戏
			_, _ = gameRepo.Get(ctx, *order.GameID)
		}
	}
}

// BenchmarkGetAvailableOrders_WithBatchQuery 测试批量查询的性能 (优化后)
func BenchmarkGetAvailableOrders_WithBatchQuery(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	// Setup test service
	orderRepo := repository.NewOrderRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	commissionRepo := repository.NewCommissionRepository(db)

	orderService := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)

	ctx := context.Background()

	// Create test data
	game := integration.CreateTestGame(b, db, "Game")

	numOrders := 50
	for i := 0; i < numOrders; i++ {
		user := integration.CreateUniqueTestUser(b, db, fmt.Sprintf("availuser%d", i))
		scheduledStart := time.Now().Add(time.Hour * time.Duration(i))
		scheduledEnd := scheduledStart.Add(time.Hour * 2)

		order := &model.Order{
			OrderNo:         fmt.Sprintf("AVAIL-%d", i),
			UserID:          user.ID,
			GameID:          &game.ID,
			PlayerID:        nil, // No player assigned yet
			Title:           fmt.Sprintf("Available Order %d", i),
			Description:     "Test available order",
			TotalPriceCents: 5000,
			Status:          model.OrderStatusConfirmed, // Confirmed but no player
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		if err := db.Create(order).Error; err != nil {
			b.Fatalf("Failed to create test order: %v", err)
		}
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := AvailableOrderListRequest{
			Page:     1,
			PageSize: numOrders,
			GameID:   &game.ID,
		}
		_, _, err := orderService.GetAvailableOrders(ctx, req)
		if err != nil {
			b.Fatalf("GetAvailableOrders failed: %v", err)
		}
	}
}

// BenchmarkGetAvailableOrders_WithoutBatchQuery 模拟 N+1 查询的性能 (优化前)
func BenchmarkGetAvailableOrders_WithoutBatchQuery(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	ctx := context.Background()

	// Create test data
	game := integration.CreateTestGame(b, db, "Game2")

	numOrders := 50
	orders := make([]*model.Order, numOrders)
	users := make([]*model.User, numOrders)

	for i := 0; i < numOrders; i++ {
		user := integration.CreateUniqueTestUser(b, db, fmt.Sprintf("availuser2%d", i))
		users[i] = user
		scheduledStart := time.Now().Add(time.Hour * time.Duration(i))
		scheduledEnd := scheduledStart.Add(time.Hour * 2)

		order := &model.Order{
			OrderNo:         fmt.Sprintf("AVAIL2-%d", i),
			UserID:          user.ID,
			GameID:          &game.ID,
			PlayerID:        nil,
			Title:           fmt.Sprintf("Available Order %d", i),
			Description:     "Test available order",
			TotalPriceCents: 5000,
			Status:          model.OrderStatusConfirmed,
			ScheduledStart:  &scheduledStart,
			ScheduledEnd:    &scheduledEnd,
		}
		if err := db.Create(order).Error; err != nil {
			b.Fatalf("Failed to create test order: %v", err)
		}
		orders[i] = order
	}

	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)

	// Run benchmark - simulating N+1 queries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟 N+1 查询：对每个订单都查询关联数据
		for _, order := range orders {
			// N+1 查询 1: 查询用户
			_, _ = userRepo.Get(ctx, order.UserID)

			// N+1 查询 2: 查询游戏
			_, _ = gameRepo.Get(ctx, *order.GameID)
		}
	}
}

// BenchmarkRepositoryGetByIDs 测试批量查询方法本身的性能
func BenchmarkRepositoryGetByIDs(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	ctx := context.Background()

	// Create test players
	numPlayers := 100
	playerIDs := make([]uint64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		user := integration.CreateUniqueTestUser(b, db, fmt.Sprintf("playeruser%d", i))
		player := integration.CreateTestPlayer(b, db, user.ID)
		playerIDs[i] = player.ID
	}

	playerRepo := repository.NewPlayerRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := playerRepo.GetByIDs(ctx, playerIDs)
		if err != nil {
			b.Fatalf("GetByIDs failed: %v", err)
		}
	}
}

// BenchmarkRepositoryGet_Single 测试单个查询的性能 (用于对比)
func BenchmarkRepositoryGet_Single(b *testing.B) {
	integration.SkipIfNoTestDB(b)
	db := integration.SetupTestDB(b)

	ctx := context.Background()

	// Create test players
	numPlayers := 100
	playerIDs := make([]uint64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		user := integration.CreateUniqueTestUser(b, db, fmt.Sprintf("singleuser%d", i))
		player := integration.CreateTestPlayer(b, db, user.ID)
		playerIDs[i] = player.ID
	}

	playerRepo := repository.NewPlayerRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟 N 次单个查询
		for _, id := range playerIDs {
			_, _ = playerRepo.Get(ctx, id)
		}
	}
}
