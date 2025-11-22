package order

import (
	"context"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlayerRepositoryForConcurrency 专为并发测试设计的陪玩师仓库
type mockPlayerRepositoryForConcurrency struct {
	players map[uint64]*model.Player
	mu      sync.Mutex
}

func newMockPlayerRepositoryForConcurrency() *mockPlayerRepositoryForConcurrency {
	return &mockPlayerRepositoryForConcurrency{
		players: make(map[uint64]*model.Player),
	}
}

func (m *mockPlayerRepositoryForConcurrency) ensureInitialized() {
	if m.players == nil {
		m.players = make(map[uint64]*model.Player)
	}
}

func (m *mockPlayerRepositoryForConcurrency) AddPlayer(id, userID uint64, nickname string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.ensureInitialized()
	m.players[id] = &model.Player{
		Base:            model.Base{ID: id},
		UserID:          userID,
		Nickname:        nickname,
		HourlyRateCents: 10000, // 100元/小时
	}
}

func (m *mockPlayerRepositoryForConcurrency) Create(ctx context.Context, player *model.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.players[player.ID] = player
	return nil
}

func (m *mockPlayerRepositoryForConcurrency) Get(ctx context.Context, id uint64) (*model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	player, exists := m.players[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return player, nil
}

func (m *mockPlayerRepositoryForConcurrency) Update(ctx context.Context, player *model.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.players[player.ID]; !exists {
		return repository.ErrNotFound
	}
	m.players[player.ID] = player
	return nil
}

func (m *mockPlayerRepositoryForConcurrency) Delete(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.players, id)
	return nil
}

func (m *mockPlayerRepositoryForConcurrency) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, player := range m.players {
		if player.UserID == userID {
			return player, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockPlayerRepositoryForConcurrency) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.Player
	for _, player := range m.players {
		result = append(result, *player)
	}
	
	// 如果没有陪玩师，返回一个默认的
	if len(result) == 0 {
		result = []model.Player{
			{
				Base:            model.Base{ID: 1},
				UserID:          100,
				Nickname:        "Default Player",
				HourlyRateCents: 10000,
			},
		}
	}
	
	total := int64(len(result))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(result) {
		return []model.Player{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

func (m *mockPlayerRepositoryForConcurrency) List(ctx context.Context) ([]model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.Player
	for _, player := range m.players {
		result = append(result, *player)
	}
	return result, nil
}

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

func (m *mockConcurrentOrderRepository) Delete(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.orders, id)
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
	playerRepo := &mockPlayerRepositoryForConcurrency{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 创建测试订单
	ctx := context.Background()
	userID := uint64(1)
	
	scheduledStart := time.Now().Add(time.Hour)
	order := &model.Order{
		UserID:         userID,
		Title:          "Test Order",
		Description:    "Test Description",
		Status:         model.OrderStatusConfirmed,
		ScheduledStart: &scheduledStart,
		TotalPriceCents: 1000,
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
		playerRepo.AddPlayer(uint64(i+1), playerUserID, "Player")
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
	
	for _, err := range results {
		if err == nil {
			successCount++
		} else {
			// 失败的请求应该是因为订单状态已变更或陪玩师不存在
			t.Logf("AcceptOrder failed: %v", err)
		}
	}
	
	// 只有一个请求应该成功
	assert.Equal(t, 1, successCount, "Only one player should successfully accept the order")
	
	// 验证订单状态
	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status)
	// 验证陪玩师ID是否设置
	assert.NotNil(t, updatedOrder.PlayerID)
}

// TestConcurrentOrderCreation 测试并发订单创建
func TestConcurrentOrderCreation(t *testing.T) {
	// 设置并发测试的 goroutine 数量
	concurrency := 20
	
	// 创建模拟仓库
	orderRepo := newMockConcurrentOrderRepository()
	playerRepo := &mockPlayerRepositoryForConcurrency{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 准备测试数据
	ctx := context.Background()
	
	// 并发创建订单
	var wg sync.WaitGroup
	results := make([]*CreateOrderResponse, concurrency)
	errors := make([]error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			scheduledStart := time.Now().Add(time.Hour + time.Duration(index)*time.Minute)
			req := CreateOrderRequest{
				PlayerID:       1,
				GameID:         1,
				Title:          "Concurrent Order",
				Description:    "Test concurrent order creation",
				ScheduledStart: &scheduledStart,
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
	playerRepo := &mockPlayerRepositoryForConcurrency{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	// 创建测试订单
	ctx := context.Background()
	scheduledStart := time.Now().Add(time.Hour)
	order := &model.Order{
		UserID:         1,
		Title:          "Race Condition Test Order",
		Description:    "Test Description",
		Status:         model.OrderStatusConfirmed,
		ScheduledStart: &scheduledStart,
		TotalPriceCents: 1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 添加陪玩师
	playerRepo.AddPlayer(1, 100, "Player A")
	playerRepo.AddPlayer(2, 101, "Player B")
	
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



