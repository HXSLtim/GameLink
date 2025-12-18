package integration

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	orderservice "gamelink/internal/service/order"
)

// mockConcurrentOrderRepo 模拟支持并发安全的订单repository
type mockConcurrentOrderRepo struct {
	orders       map[uint64]*model.Order
	mu           sync.Mutex
	updateCount  int32 // 原子计数器,记录成功更新次数
	getCallCount int32 // 记录Get调用次数
}

func newMockConcurrentOrderRepo() *mockConcurrentOrderRepo {
	return &mockConcurrentOrderRepo{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockConcurrentOrderRepo) Create(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if order.ID == 0 {
		order.ID = uint64(len(m.orders) + 1)
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockConcurrentOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	atomic.AddInt32(&m.getCallCount, 1)

	m.mu.Lock()
	defer m.mu.Unlock()

	order, ok := m.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}

	// 返回副本,避免外部直接修改
	orderCopy := *order
	return &orderCopy, nil
}

func (m *mockConcurrentOrderRepo) Update(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.orders[order.ID]; !ok {
		return errors.New("order not found")
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockConcurrentOrderRepo) Delete(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.orders, id)
	return nil
}

func (m *mockConcurrentOrderRepo) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}

// UpdateWithCondition 实现原子性条件更新
func (m *mockConcurrentOrderRepo) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	order, ok := m.orders[orderID]
	if !ok {
		return false, errors.New("order not found")
	}

	// 检查状态是否匹配
	if order.Status != expectedStatus {
		return false, nil
	}

	// 应用更新
	if playerID, ok := updates["player_id"].(uint64); ok {
		order.PlayerID = &playerID
	}
	if status, ok := updates["status"].(model.OrderStatus); ok {
		order.Status = status
	}
	if startedAt, ok := updates["started_at"].(*time.Time); ok {
		order.StartedAt = startedAt
	}

	atomic.AddInt32(&m.updateCount, 1)
	return true, nil
}

// mockPlayerRepo 模拟玩家repository
type mockPlayerRepo struct {
	players []*model.Player
}

func newMockPlayerRepo(players []*model.Player) *mockPlayerRepo {
	return &mockPlayerRepo{players: players}
}

func (m *mockPlayerRepo) Create(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	for _, p := range m.players {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("player not found")
}

func (m *mockPlayerRepo) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	result := make([]model.Player, len(m.players))
	for i, p := range m.players {
		result[i] = *p
	}
	return result, int64(len(result)), nil
}

func (m *mockPlayerRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	for _, p := range m.players {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, errors.New("player not found")
}

func (m *mockPlayerRepo) List(ctx context.Context) ([]model.Player, error) {
	result := make([]model.Player, len(m.players))
	for i, p := range m.players {
		result[i] = *p
	}
	return result, nil
}

func (m *mockPlayerRepo) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *mockPlayerRepo) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	return 0, nil
}

func (m *mockPlayerRepo) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return 0, nil
}

// mockUserRepo 模拟用户repository
type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error      { return nil }
func (m *mockUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) { return nil, nil }
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepo) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *mockUserRepo) List(ctx context.Context) ([]model.User, error) {
	return nil, nil
}
func (m *mockUserRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) ListWithFilters(ctx context.Context, filters repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}
func (m *mockUserRepo) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	return nil, nil
}

// TestAcceptOrderConcurrency 测试多个陪玩师同时接单的并发安全性
func TestAcceptOrderConcurrency(t *testing.T) {
	ctx := context.Background()

	// 创建测试数据
	const numPlayers = 5
	players := make([]*model.Player, numPlayers)
	playerUserIDs := make([]uint64, numPlayers)

	for i := 0; i < numPlayers; i++ {
		userID := uint64(i + 100)
		playerUserIDs[i] = userID
		players[i] = &model.Player{
			Base:               model.Base{ID: uint64(i + 1)},
			UserID:             userID,
			Nickname:           "Player",
			VerificationStatus: model.VerificationVerified,
		}
	}

	// 创建mock repositories
	orderRepo := newMockConcurrentOrderRepo()
	playerRepo := newMockPlayerRepo(players)
	userRepo := &mockUserRepo{}

	// 创建待接单的订单
	scheduledStart := time.Now().Add(1 * time.Hour)
	itemID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		ItemID:          itemID,
		Status:          model.OrderStatusConfirmed,
		Quantity:        1,
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
		Currency:        "CNY",
		Title:           "Test Order",
		ScheduledStart:  &scheduledStart,
	}
	assert.NoError(t, orderRepo.Create(ctx, order))

	// 创建service (参数顺序: orders, users, players, payments, serviceItems, games, reviews)
	service := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, nil, nil, nil, nil)

	// 并发测试: 多个陪玩师同时尝试接单
	var wg sync.WaitGroup
	results := make([]error, numPlayers)
	startSignal := make(chan struct{})

	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// 等待开始信号,确保所有goroutine尽可能同时开始
			<-startSignal
			results[idx] = service.AcceptOrder(ctx, playerUserIDs[idx], order.ID)
		}(i)
	}

	// 发送开始信号,开始并发测试
	close(startSignal)
	wg.Wait()

	// 验证结果
	successCount := 0
	failCount := 0

	for i, err := range results {
		if err == nil {
			successCount++
			t.Logf("Player %d (UserID=%d): 成功接单", i, playerUserIDs[i])
		} else {
			failCount++
			if err == orderservice.ErrInvalidTransition {
				t.Logf("Player %d (UserID=%d): 接单失败 - 订单状态已变更", i, playerUserIDs[i])
			} else {
				t.Logf("Player %d (UserID=%d): 接单失败 - %v", i, playerUserIDs[i], err)
			}
		}
	}

	// 关键断言: 只有1个成功
	assert.Equal(t, 1, successCount, "应该只有一个陪玩师成功接单")
	assert.Equal(t, numPlayers-1, failCount, "其他陪玩师应该都失败")

	// 验证订单最终状态
	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status, "订单状态应该是进行中")
	assert.NotNil(t, updatedOrder.PlayerID, "订单应该有陪玩师ID")
	assert.NotNil(t, updatedOrder.StartedAt, "订单应该有开始时间")

	// 验证只执行了一次成功的更新
	updateCount := atomic.LoadInt32(&orderRepo.updateCount)
	assert.Equal(t, int32(1), updateCount, "应该只有一次成功的更新操作")

	t.Logf("总计: %d个陪玩师, %d个成功, %d个失败", numPlayers, successCount, failCount)
	t.Logf("数据库更新次数: %d", updateCount)
}

// TestAcceptOrderConcurrencyStress 压力测试: 20个陪玩师同时抢单
func TestAcceptOrderConcurrencyStress(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	ctx := context.Background()
	const numPlayers = 20

	// 创建测试数据
	players := make([]*model.Player, numPlayers)
	playerUserIDs := make([]uint64, numPlayers)

	for i := 0; i < numPlayers; i++ {
		userID := uint64(i + 200)
		playerUserIDs[i] = userID
		players[i] = &model.Player{
			Base:               model.Base{ID: uint64(i + 1)},
			UserID:             userID,
			Nickname:           "StressPlayer",
			VerificationStatus: model.VerificationVerified,
		}
	}

	orderRepo := newMockConcurrentOrderRepo()
	playerRepo := newMockPlayerRepo(players)
	userRepo := &mockUserRepo{}

	scheduledStart := time.Now().Add(2 * time.Hour)
	itemID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          1,
		ItemID:          itemID,
		Status:          model.OrderStatusConfirmed,
		Quantity:        1,
		UnitPriceCents:  20000,
		TotalPriceCents: 20000,
		Currency:        "CNY",
		Title:           "Stress Test Order",
		ScheduledStart:  &scheduledStart,
	}
	assert.NoError(t, orderRepo.Create(ctx, order))

	service := orderservice.NewOrderService(orderRepo, playerRepo, userRepo, nil, nil, nil, nil)

	// 高并发测试
	var wg sync.WaitGroup
	results := make([]error, numPlayers)
	startSignal := make(chan struct{})

	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-startSignal
			results[idx] = service.AcceptOrder(ctx, playerUserIDs[idx], order.ID)
		}(i)
	}

	close(startSignal)
	wg.Wait()

	// 统计结果
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}

	// 关键断言
	assert.Equal(t, 1, successCount, "在高并发场景下,应该只有一个陪玩师成功接单")

	updateCount := atomic.LoadInt32(&orderRepo.updateCount)
	assert.Equal(t, int32(1), updateCount, "应该只执行了一次更新操作")

	t.Logf("压力测试完成: %d个并发请求, %d个成功, 数据库更新%d次", numPlayers, successCount, updateCount)
}

// TestUpdateWithConditionAtomicity 测试UpdateWithCondition的原子性
func TestUpdateWithConditionAtomicity(t *testing.T) {
	ctx := context.Background()
	orderRepo := newMockConcurrentOrderRepo()

	// 创建订单
	itemID := uint64(1)
	order := &model.Order{
		Base:            model.Base{ID: 1},
		ItemID:          itemID,
		Status:          model.OrderStatusConfirmed,
		Quantity:        1,
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
	}
	assert.NoError(t, orderRepo.Create(ctx, order))

	// 10个goroutine同时尝试更新
	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make([]bool, numGoroutines)
	errors := make([]error, numGoroutines)
	startSignal := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-startSignal

			now := time.Now()
			updated, err := orderRepo.UpdateWithCondition(
				ctx,
				order.ID,
				model.OrderStatusConfirmed,
				map[string]any{
					"player_id":  uint64(idx + 300),
					"status":     model.OrderStatusInProgress,
					"started_at": &now,
				},
			)
			results[idx] = updated
			errors[idx] = err
		}(i)
	}

	close(startSignal)
	wg.Wait()

	// 验证结果
	successCount := 0
	for i, updated := range results {
		assert.NoError(t, errors[i], "不应该有数据库错误")
		if updated {
			successCount++
			t.Logf("Goroutine %d: 更新成功", i)
		}
	}

	// 关键断言: 只有一个成功
	assert.Equal(t, 1, successCount, "应该只有一个goroutine成功更新")

	// 验证更新次数
	updateCount := atomic.LoadInt32(&orderRepo.updateCount)
	assert.Equal(t, int32(1), updateCount, "应该只执行了一次更新")

	// 验证最终状态
	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStatusInProgress, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.PlayerID)
}
