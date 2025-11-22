package order

import (
	"context"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockConcurrentOrderRepository 模拟并发场景的订单仓库
type mockConcurrentOrderRepository struct {
	orders map[uint64]*model.Order
	mu     sync.Mutex
}

func newMockConcurrentOrderRepository() *mockConcurrentOrderRepository {
	return &mockConcurrentOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockConcurrentOrderRepository) Create(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockConcurrentOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	order, exists := m.orders[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return order, nil
}

func (m *mockConcurrentOrderRepository) Update(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.orders[order.ID]; !exists {
		return repository.ErrNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockConcurrentOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.Order
	for _, o := range m.orders {
		if len(opts.Statuses) > 0 {
			match := false
			for _, s := range opts.Statuses {
				if o.Status == s {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		if opts.UserID != nil && o.UserID != *opts.UserID {
			continue
		}
		result = append(result, *o)
	}
	
	total := int64(len(result))
	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start >= len(result) {
		return []model.Order{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

// TestConcurrentAcceptOrder 测试并发接单场景
func TestConcurrentAcceptOrder(t *testing.T) {
	// 设置并发测试的 goroutine 数量
	concurrency := 10
	
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := newMockPlayerRepository()
	userRepo := newMockUserRepository()
	gameRepo := newMockGameRepository()
	paymentRepo := newMockPaymentRepository()
	reviewRepo := newMockReviewRepository()
	commissionRepo := newMockCommissionRepository()
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 创建测试订单
	ctx := context.Background()
	userID := uint64(1)
	
	order := &model.Order{
		UserID:      userID,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusConfirmed,
		ScheduledStart: time.Now().Add(time.Hour),
		DurationHours: 2.0,
		PriceCents: 1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 创建多个陪玩师
	var playerUserIDs []uint64
	for i := 0; i < concurrency; i++ {
		playerUserID := uint64(100 + i)
		playerUserIDs = append(playerUserIDs, playerUserID)
		
		player := &model.Player{
			ID:     uint64(i + 1),
			UserID: playerUserID,
			Nickname: "Player " + string(rune('A'+i)),
			HourlyPriceCents: 500,
			GameID: 1,
		}
		playerRepo.players[player.ID] = player
	}
	
	// 并发测试接单
	var wg sync.WaitGroup
	results := make([]error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int, playerUserID uint64) {
			defer wg.Done()
			results[index] = service.AcceptOrder(ctx, playerUserID, order.ID)
		}(i, playerUserIDs[i])
	}
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	var firstSuccessPlayerID uint64
	
	for i, err := range results {
		if err == nil {
			successCount++
			firstSuccessPlayerID = playerUserIDs[i]
		} else {
			// 失败的请求应该是因为订单状态已变更
			assert.Contains(t, err.Error(), "invalid order status")
		}
	}
	
	// 只有一个请求应该成功
	assert.Equal(t, 1, successCount, "Only one player should successfully accept the order")
	
	// 验证订单状态
	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status)
	assert.Equal(t, firstSuccessPlayerID, updatedOrder.GetPlayerUserID())
}

// TestConcurrentOrderCreation 测试并发订单创建
func TestConcurrentOrderCreation(t *testing.T) {
	// 设置并发测试的 goroutine 数量
	concurrency := 20
	
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := newMockPlayerRepository()
	userRepo := newMockUserRepository()
	gameRepo := newMockGameRepository()
	paymentRepo := newMockPaymentRepository()
	reviewRepo := newMockReviewRepository()
	commissionRepo := newMockCommissionRepository()
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 准备测试数据
	ctx := context.Background()
	player := &model.Player{
		ID:     1,
		UserID: 100,
		Nickname: "Test Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	playerRepo.players[1] = player
	
	game := &model.Game{
		ID:   1,
		Name: "Test Game",
	}
	gameRepo.games[1] = game
	
	// 并发创建订单
	var wg sync.WaitGroup
	results := make([]*CreateOrderResponse, concurrency)
	errors := make([]error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			req := CreateOrderRequest{
				PlayerID:       1,
				GameID:         1,
				Title:          "Concurrent Order " + string(rune('A'+index)),
				Description:    "Test concurrent order creation",
				ScheduledStart: time.Now().Add(time.Hour + time.Duration(index)*time.Minute),
				DurationHours:  1.0,
			}
			
			results[index], errors[index] = service.CreateOrder(ctx, uint64(index+1), req)
		}(i)
	}
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	for i := 0; i < concurrency; i++ {
		if errors[i] == nil {
			successCount++
			assert.NotNil(t, results[i])
			assert.Greater(t, results[i].OrderID, uint64(0))
		} else {
			t.Logf("Order creation %d failed: %v", i, errors[i])
		}
	}
	
	// 所有订单都应该创建成功
	assert.Equal(t, concurrency, successCount, "All concurrent orders should be created successfully")
	
	// 验证订单数量
	orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
		Page: 1,
		PageSize: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(concurrency), total)
	assert.Equal(t, concurrency, len(orders))
}

// TestRaceConditionAcceptOrder 测试接单竞态条件（使用 -race 标志运行）
func TestRaceConditionAcceptOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}
	
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := newMockPlayerRepository()
	userRepo := newMockUserRepository()
	gameRepo := newMockGameRepository()
	paymentRepo := newMockPaymentRepository()
	reviewRepo := newMockReviewRepository()
	commissionRepo := newMockCommissionRepository()
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 创建测试订单
	ctx := context.Background()
	order := &model.Order{
		UserID:      1,
		Title:       "Race Condition Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusConfirmed,
		ScheduledStart: time.Now().Add(time.Hour),
		DurationHours: 2.0,
		PriceCents: 1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 创建两个陪玩师
	player1 := &model.Player{
		ID:     1,
		UserID: 100,
		Nickname: "Player A",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	player2 := &model.Player{
		ID:     2,
		UserID: 101,
		Nickname: "Player B", 
		HourlyPriceCents: 500,
		GameID: 1,
	}
	playerRepo.players[1] = player1
	playerRepo.players[2] = player2
	
	// 并发尝试接单
	var wg sync.WaitGroup
	results := make([]error, 2)
	
	wg.Add(2)
	go func() {
		defer wg.Done()
		results[0] = service.AcceptOrder(ctx, 100, order.ID)
	}()
	
	go func() {
		defer wg.Done()
		results[1] = service.AcceptOrder(ctx, 101, order.ID)
	}()
	
	wg.Wait()
	
	// 验证竞态条件处理
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}
	
	// 只有一个应该成功
	assert.Equal(t, 1, successCount, "Only one player should successfully accept the order")
	
	// 验证最终状态
	finalOrder, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, finalOrder.Status)
}

// BenchmarkAcceptOrder 接单性能基准测试
func BenchmarkAcceptOrder(b *testing.B) {
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := newMockPlayerRepository()
	userRepo := newMockUserRepository()
	gameRepo := newMockGameRepository()
	paymentRepo := newMockPaymentRepository()
	reviewRepo := newMockReviewRepository()
	commissionRepo := newMockCommissionRepository()
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	ctx := context.Background()
	
	// 创建陪玩师
	player := &model.Player{
		ID:     1,
		UserID: 100,
		Nickname: "Benchmark Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	playerRepo.players[1] = player
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		// 创建订单
		order := &model.Order{
			UserID:      1,
			Title:       "Benchmark Order",
			Description: "Benchmark Description",
			Status:      model.OrderStatusConfirmed,
			ScheduledStart: time.Now().Add(time.Hour),
			DurationHours: 2.0,
			PriceCents: 1000,
			CommissionCents: 200,
			PlayerIncomeCents: 800,
		}
		
		err := orderRepo.Create(ctx, order)
		if err != nil {
			b.Fatal(err)
		}
		
		// 接单操作
		err = service.AcceptOrder(ctx, 100, order.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCreateOrder 订单创建性能基准测试
func BenchmarkCreateOrder(b *testing.B) {
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := newMockPlayerRepository()
	userRepo := newMockUserRepository()
	gameRepo := newMockGameRepository()
	paymentRepo := newMockPaymentRepository()
	reviewRepo := newMockReviewRepository()
	commissionRepo := newMockCommissionRepository()
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	ctx := context.Background()
	userID := uint64(1)
	
	// 准备测试数据
	player := &model.Player{
		ID:     1,
		UserID: 100,
		Nickname: "Benchmark Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	playerRepo.players[1] = player
	
	game := &model.Game{
		ID:   1,
		Name: "Benchmark Game",
	}
	gameRepo.games[1] = game
	
	req := CreateOrderRequest{
		PlayerID:       1,
		GameID:         1,
		Title:          "Benchmark Order",
		Description:    "Benchmark Description",
		ScheduledStart: time.Now().Add(time.Hour),
		DurationHours:  2.0,
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := service.CreateOrder(ctx, userID, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}