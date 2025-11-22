package commission

import (
	"context"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	repoiface "gamelink/internal/repository/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockConcurrentCommissionRepository 模拟并发场景的佣金仓库
type mockConcurrentCommissionRepository struct {
	rules    map[uint64]*model.CommissionRule
	records  map[uint64]*model.CommissionRecord
	mu       sync.Mutex
}

func newMockConcurrentCommissionRepository() *mockConcurrentCommissionRepository {
	return &mockConcurrentCommissionRepository{
		rules:   make(map[uint64]*model.CommissionRule),
		records: make(map[uint64]*model.CommissionRecord),
	}
}

func (m *mockConcurrentCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	rule.ID = uint64(len(m.rules) + 1)
	m.rules[rule.ID] = rule
	return nil
}

func (m *mockConcurrentCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	rule, exists := m.rules[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return rule, nil
}

func (m *mockConcurrentCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, rule := range m.rules {
		// 检查是否是默认规则 - 根据Type字段判断
		if rule.Type == "default" {
			return rule, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockConcurrentCommissionRepository) GetRuleForOrder(ctx context.Context, gameID, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 简化实现：返回第一个规则
	for _, rule := range m.rules {
		return rule, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockConcurrentCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.rules[rule.ID]; !exists {
		return repository.ErrNotFound
	}
	m.rules[rule.ID] = rule
	return nil
}

func (m *mockConcurrentCommissionRepository) DeleteRule(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.rules, id)
	return nil
}

func (m *mockConcurrentCommissionRepository) ListRules(ctx context.Context, opts repository.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.CommissionRule
	for _, rule := range m.rules {
		result = append(result, *rule)
	}
	
	total := int64(len(result))
	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start >= len(result) {
		return []model.CommissionRule{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

func (m *mockConcurrentCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查是否已存在
	for _, existing := range m.records {
		if existing.OrderID == record.OrderID {
			return ErrAlreadyRecorded
		}
	}
	
	record.ID = uint64(len(m.records) + 1)
	m.records[record.ID] = record
	return nil
}

func (m *mockConcurrentCommissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	record, exists := m.records[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return record, nil
}

func (m *mockConcurrentCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, record := range m.records {
		if record.OrderID == orderID {
			return record, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockConcurrentCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil // 简化实现
}

func (m *mockConcurrentCommissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound // 简化实现
}

func (m *mockConcurrentCommissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound // 简化实现
}

func (m *mockConcurrentCommissionRepository) ListSettlements(ctx context.Context, opts interface{}) ([]model.MonthlySettlement, int64, error) {
	return []model.MonthlySettlement{}, 0, nil // 简化实现
}

func (m *mockConcurrentCommissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil // 简化实现
}

func (m *mockConcurrentCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (interface{}, error) {
	return nil, nil // 简化实现
}

func (m *mockConcurrentCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil // 简化实现
}

func (m *mockConcurrentCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.records[record.ID]; !exists {
		return repository.ErrNotFound
	}
	m.records[record.ID] = record
	return nil
}

func (m *mockConcurrentCommissionRepository) ListRecords(ctx context.Context, opts repository.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.CommissionRecord
	for _, record := range m.records {
		if opts.PlayerID != nil && record.PlayerID == *opts.PlayerID {
			result = append(result, *record)
		} else if opts.PlayerID == nil {
			result = append(result, *record)
		}
	}
	
	total := int64(len(result))
	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start >= len(result) {
		return []model.CommissionRecord{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

// mockConcurrentOrderReader 模拟并发场景的订单读取器
type mockConcurrentOrderReader struct {
	orders map[uint64]*model.Order
	mu     sync.Mutex
}

func newMockConcurrentOrderReader() *mockConcurrentOrderReader {
	return &mockConcurrentOrderReader{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockConcurrentOrderReader) Get(ctx context.Context, id uint64) (*model.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	order, exists := m.orders[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return order, nil
}

func (m *mockConcurrentOrderReader) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
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

// mockConcurrentPlayerRepository 模拟并发场景的陪玩师仓库
type mockConcurrentPlayerRepository struct {
	players map[uint64]*model.Player
	mu      sync.Mutex
}

func newMockConcurrentPlayerRepository() *mockConcurrentPlayerRepository {
	return &mockConcurrentPlayerRepository{
		players: make(map[uint64]*model.Player),
	}
}

func (m *mockConcurrentPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	player.ID = uint64(len(m.players) + 1)
	m.players[player.ID] = player
	return nil
}

func (m *mockConcurrentPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	player, exists := m.players[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return player, nil
}

func (m *mockConcurrentPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, player := range m.players {
		if player.UserID == userID {
			return player, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockConcurrentPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.Player
	for _, player := range m.players {
		result = append(result, *player)
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

// TestConcurrentRecordCommission 测试并发记录佣金场景
func TestConcurrentRecordCommission(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent commission test in short mode")
	}
	
	// 设置并发测试的 goroutine 数量
	concurrency := 10
	
	// 创建模拟仓库
	commissionRepo := newMockConcurrentCommissionRepository()
	orderRepo := newMockConcurrentOrderReader()
	playerRepo := newMockConcurrentPlayerRepository()
	
	// 创建佣金服务
	service := NewCommissionService(commissionRepo, orderRepo, playerRepo)
	
	// 创建默认佣金规则
	ctx := context.Background()
	defaultRule := &model.CommissionRule{
		Name:        "Default Rule",
		IsDefault:   true,
		Rate:        0.2,
		MinAmount:   100,
		MaxAmount:   10000,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}
	
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)
	
	// 创建测试订单和陪玩师
	orderID := uint64(1)
	playerID := uint64(1)
	playerUserID := uint64(100)
	
	order := &model.Order{
		ID:          orderID,
		UserID:      1,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusCompleted,
		PriceCents:  1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err = orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 设置订单的陪玩师
	order.SetPlayerID(playerID)
	
	player := &model.Player{
		ID:     playerID,
		UserID: playerUserID,
		Nickname: "Test Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	
	err = playerRepo.Create(ctx, player)
	require.NoError(t, err)
	
	// 并发记录佣金
	var wg sync.WaitGroup
	results := make([]error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = service.RecordCommission(ctx, orderID)
		}(i)
	}
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	duplicateCount := 0
	
	for _, err := range results {
		if err == nil {
			successCount++
		} else if err == ErrAlreadyRecorded {
			duplicateCount++
		}
	}
	
	// 只有一个请求应该成功
	assert.Equal(t, 1, successCount, "Only one commission record should be created")
	assert.Equal(t, concurrency-1, duplicateCount, "Other attempts should fail with already recorded error")
	
	// 验证佣金记录数量
	record, err := commissionRepo.GetRecordByOrderID(ctx, orderID)
	require.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, orderID, record.OrderID)
	assert.Equal(t, playerID, record.PlayerID)
}

// TestConcurrentCalculateCommission 测试并发计算佣金
func TestConcurrentCalculateCommission(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent commission calculation test in short mode")
	}
	
	// 创建模拟仓库
	commissionRepo := newMockConcurrentCommissionRepository()
	orderRepo := newMockConcurrentOrderReader()
	playerRepo := newMockConcurrentPlayerRepository()
	
	// 创建佣金服务
	service := NewCommissionService(commissionRepo, orderRepo, playerRepo)
	
	// 创建默认佣金规则
	ctx := context.Background()
	defaultRule := &model.CommissionRule{
		Name:        "Default Rule",
		IsDefault:   true,
		Rate:        0.2,
		MinAmount:   100,
		MaxAmount:   10000,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}
	
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)
	
	// 创建测试订单
	orderID := uint64(1)
	order := &model.Order{
		ID:          orderID,
		UserID:      1,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusCompleted,
		PriceCents:  1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err = orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 并发计算佣金
	var wg sync.WaitGroup
	results := make([]*CommissionCalculation, 10)
	errors := make([]error, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errors[index] = service.CalculateCommission(ctx, orderID)
		}(i)
	}
	
	wg.Wait()
	
	// 验证所有计算都应该成功且结果一致
	var expectedResult *CommissionCalculation
	for i := 0; i < 10; i++ {
		assert.NoError(t, errors[i])
		assert.NotNil(t, results[i])
		
		if expectedResult == nil {
			expectedResult = results[i]
		} else {
			// 验证结果一致性
			assert.Equal(t, expectedResult.TotalAmountCents, results[i].TotalAmountCents)
			assert.Equal(t, expectedResult.CommissionRate, results[i].CommissionRate)
			assert.Equal(t, expectedResult.CommissionCents, results[i].CommissionCents)
			assert.Equal(t, expectedResult.PlayerIncomeCents, results[i].PlayerIncomeCents)
		}
	}
}

// TestRaceConditionCommissionSettlement 测试佣金结算竞态条件
func TestRaceConditionCommissionSettlement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}
	
	// 创建模拟仓库
	commissionRepo := newMockConcurrentCommissionRepository()
	orderRepo := newMockConcurrentOrderReader()
	playerRepo := newMockConcurrentPlayerRepository()
	
	// 创建佣金服务
	service := NewCommissionService(commissionRepo, orderRepo, playerRepo)
	
	// 创建默认佣金规则
	ctx := context.Background()
	defaultRule := &model.CommissionRule{
		Name:        "Default Rule",
		IsDefault:   true,
		Rate:        0.2,
		MinAmount:   100,
		MaxAmount:   10000,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}
	
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)
	
	// 创建测试订单和陪玩师
	orderID := uint64(1)
	playerID := uint64(1)
	
	order := &model.Order{
		ID:          orderID,
		UserID:      1,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusCompleted,
		PriceCents:  1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err = orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 设置订单的陪玩师
	order.SetPlayerID(playerID)
	
	player := &model.Player{
		ID:     playerID,
		UserID: 100,
		Nickname: "Test Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	
	err = playerRepo.Create(ctx, player)
	require.NoError(t, err)
	
	// 先记录佣金
	err = service.RecordCommission(ctx, orderID)
	require.NoError(t, err)
	
	// 并发尝试结算
	var wg sync.WaitGroup
	results := make([]error, 2)
	
	wg.Add(2)
	go func() {
		defer wg.Done()
		results[0] = service.SettlePlayerCommission(ctx, playerID, time.Now().Format("2006-01"))
	}()
	
	go func() {
		defer wg.Done()
		results[1] = service.SettlePlayerCommission(ctx, playerID, time.Now().Format("2006-01"))
	}()
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}
	
	// 两个结算操作都应该成功（幂等性）
	assert.Equal(t, 2, successCount, "Both settlement operations should succeed")
}

// BenchmarkCalculateCommission 佣金计算性能基准测试
func BenchmarkCalculateCommission(b *testing.B) {
	// 创建模拟仓库
	commissionRepo := newMockConcurrentCommissionRepository()
	orderRepo := newMockConcurrentOrderReader()
	playerRepo := newMockConcurrentPlayerRepository()
	
	// 创建佣金服务
	service := NewCommissionService(commissionRepo, orderRepo, playerRepo)
	
	ctx := context.Background()
	
	// 创建默认佣金规则
	defaultRule := &model.CommissionRule{
		Name:        "Default Rule",
		IsDefault:   true,
		Rate:        0.2,
		MinAmount:   100,
		MaxAmount:   10000,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}
	
	err := commissionRepo.CreateRule(ctx, defaultRule)
	if err != nil {
		b.Fatal(err)
	}
	
	// 创建测试订单
	orderID := uint64(1)
	order := &model.Order{
		ID:          orderID,
		UserID:      1,
		Title:       "Benchmark Order",
		Description: "Benchmark Description",
		Status:      model.OrderStatusCompleted,
		PriceCents:  1000,
		CommissionCents: 200,
		PlayerIncomeCents: 800,
	}
	
	err = orderRepo.Create(ctx, order)
	if err != nil {
		b.Fatal(err)
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := service.CalculateCommission(ctx, orderID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRecordCommission 佣金记录性能基准测试
func BenchmarkRecordCommission(b *testing.B) {
	// 创建模拟仓库
	commissionRepo := newMockConcurrentCommissionRepository()
	orderRepo := newMockConcurrentOrderReader()
	playerRepo := newMockConcurrentPlayerRepository()
	
	// 创建佣金服务
	service := NewCommissionService(commissionRepo, orderRepo, playerRepo)
	
	ctx := context.Background()
	
	// 创建默认佣金规则
	defaultRule := &model.CommissionRule{
		Name:        "Default Rule",
		IsDefault:   true,
		Rate:        0.2,
		MinAmount:   100,
		MaxAmount:   10000,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}
	
	err := commissionRepo.CreateRule(ctx, defaultRule)
	if err != nil {
		b.Fatal(err)
	}
	
	// 创建陪玩师
	player := &model.Player{
		ID:     1,
		UserID: 100,
		Nickname: "Benchmark Player",
		HourlyPriceCents: 500,
		GameID: 1,
	}
	
	err = playerRepo.Create(ctx, player)
	if err != nil {
		b.Fatal(err)
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		// 每次创建新的订单
		orderID := uint64(i + 1)
		
		order := &model.Order{
			ID:          orderID,
			UserID:      1,
			Title:       "Benchmark Order",
			Description: "Benchmark Description",
			Status:      model.OrderStatusCompleted,
			PriceCents:  1000,
			CommissionCents: 200,
			PlayerIncomeCents: 800,
		}
		
		err := orderRepo.Create(ctx, order)
		if err != nil {
			b.Fatal(err)
		}
		
		// 设置订单的陪玩师
		order.SetPlayerID(1)
		
		// 记录佣金
		err = service.RecordCommission(ctx, orderID)
		if err != nil {
			b.Fatal(err)
		}
	}
}