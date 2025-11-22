package payment

import (
	"context"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockConcurrentPaymentRepository 模拟并发场景的支付仓库
type mockConcurrentPaymentRepository struct {
	payments map[uint64]*model.Payment
	orders   map[uint64]*model.Order
	mu       sync.Mutex
}

func newMockConcurrentPaymentRepository() *mockConcurrentPaymentRepository {
	return &mockConcurrentPaymentRepository{
		payments: make(map[uint64]*model.Payment),
		orders:   make(map[uint64]*model.Order),
	}
}

func (m *mockConcurrentPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	payment.ID = uint64(len(m.payments) + 1)
	payment.CreatedAt = time.Now()
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockConcurrentPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	payment, exists := m.payments[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return payment, nil
}

func (m *mockConcurrentPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.payments[payment.ID]; !exists {
		return repository.ErrNotFound
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockConcurrentPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var result []model.Payment
	for _, p := range m.payments {
		if opts.OrderID != nil && p.OrderID != *opts.OrderID {
			continue
		}
		if opts.UserID != nil && p.UserID != *opts.UserID {
			continue
		}
		result = append(result, *p)
	}
	
	total := int64(len(result))
	start := (opts.Page - 1) * opts.PageSize
	end := start + opts.PageSize
	if start >= len(result) {
		return []model.Payment{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	
	return result[start:end], total, nil
}

func (m *mockConcurrentPaymentRepository) Delete(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.payments, id)
	return nil
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

func (m *mockConcurrentOrderRepository) Delete(ctx context.Context, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.orders, id)
	return nil
}

// TestConcurrentCreatePayment 测试并发创建支付场景
func TestConcurrentCreatePayment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent payment test in short mode")
	}
	
	// 设置并发测试的 goroutine 数量
	concurrency := 10
	
	// 创建模拟仓库
	paymentRepo := newMockConcurrentPaymentRepository()
	orderRepo := newMockConcurrentOrderRepository()
	
	// 创建支付服务
	service := NewPaymentService(paymentRepo, orderRepo)
	
	// 创建测试订单
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	
	order := &model.Order{
		ID:          orderID,
		UserID:      userID,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusPending,
		PriceCents:  1000,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 并发创建支付
	var wg sync.WaitGroup
	results := make([]*CreatePaymentResponse, concurrency)
	errors := make([]error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			req := CreatePaymentRequest{
				OrderID: orderID,
				Method:  model.PaymentMethodWeChat,
			}
			
			results[index], errors[index] = service.CreatePayment(ctx, userID, req)
		}(i)
	}
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	duplicateCount := 0
	
	for i := 0; i < concurrency; i++ {
		if errors[i] == nil {
			successCount++
			assert.NotNil(t, results[i])
			assert.Greater(t, results[i].PaymentID, uint64(0))
		} else {
			// 检查是否是因为重复支付的错误
			if errors[i] == ErrOrderAlreadyPaid {
				duplicateCount++
			}
		}
	}
	
	// 只有一个支付应该成功创建
	assert.Equal(t, 1, successCount, "Only one payment should be created successfully")
	assert.Equal(t, concurrency-1, duplicateCount, "Other attempts should fail with duplicate payment error")
	
	// 验证支付记录数量
	payments, total, err := paymentRepo.List(ctx, repository.PaymentListOptions{
		OrderID:  &orderID,
		Page:     1,
		PageSize: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(payments))
}

// TestRaceConditionPaymentStatus 测试支付状态竞态条件
func TestRaceConditionPaymentStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}
	
	// 创建模拟仓库
	paymentRepo := newMockConcurrentPaymentRepository()
	orderRepo := newMockConcurrentOrderRepository()
	
	// 创建支付服务
	service := NewPaymentService(paymentRepo, orderRepo)
	
	// 创建测试订单和支付
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	
	order := &model.Order{
		ID:          orderID,
		UserID:      userID,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusPending,
		PriceCents:  1000,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 创建支付
	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}
	
	paymentResp, err := service.CreatePayment(ctx, userID, req)
	require.NoError(t, err)
	require.NotNil(t, paymentResp)
	
	// 并发查询支付状态
	var wg sync.WaitGroup
	results := make([]*PaymentStatusResponse, 10)
	errors := make([]error, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errors[index] = service.GetPaymentStatus(ctx, userID, paymentResp.PaymentID)
		}(i)
	}
	
	wg.Wait()
	
	// 验证所有查询都应该成功
	for i := 0; i < 10; i++ {
		assert.NoError(t, errors[i])
		assert.NotNil(t, results[i])
		assert.Equal(t, paymentResp.PaymentID, results[i].PaymentID)
		assert.Equal(t, orderID, results[i].OrderID)
	}
}

// TestConcurrentPaymentCancellation 测试并发支付取消
func TestConcurrentPaymentCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent payment cancellation test in short mode")
	}
	
	// 创建模拟仓库
	paymentRepo := newMockConcurrentPaymentRepository()
	orderRepo := newMockConcurrentOrderRepository()
	
	// 创建支付服务
	service := NewPaymentService(paymentRepo, orderRepo)
	
	// 创建测试订单和支付
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	
	order := &model.Order{
		ID:          orderID,
		UserID:      userID,
		Title:       "Test Order",
		Description: "Test Description",
		Status:      model.OrderStatusPending,
		PriceCents:  1000,
	}
	
	err := orderRepo.Create(ctx, order)
	require.NoError(t, err)
	
	// 创建支付
	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}
	
	paymentResp, err := service.CreatePayment(ctx, userID, req)
	require.NoError(t, err)
	require.NotNil(t, paymentResp)
	
	// 并发取消支付
	var wg sync.WaitGroup
	results := make([]error, 5)
	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = service.CancelPayment(ctx, userID, paymentResp.PaymentID)
		}(i)
	}
	
	wg.Wait()
	
	// 验证结果
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}
	
	// 只有一个取消操作应该成功
	assert.Equal(t, 1, successCount, "Only one cancellation should succeed")
}

// BenchmarkCreatePayment 支付创建性能基准测试
func BenchmarkCreatePayment(b *testing.B) {
	// 创建模拟仓库
	paymentRepo := newMockConcurrentPaymentRepository()
	orderRepo := newMockConcurrentOrderRepository()
	
	// 创建支付服务
	service := NewPaymentService(paymentRepo, orderRepo)
	
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	
	// 创建测试订单
	order := &model.Order{
		ID:          orderID,
		UserID:      userID,
		Title:       "Benchmark Order",
		Description: "Benchmark Description",
		Status:      model.OrderStatusPending,
		PriceCents:  1000,
	}
	
	err := orderRepo.Create(ctx, order)
	if err != nil {
		b.Fatal(err)
	}
	
	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		// 每次创建新的订单ID
		currentOrderID := uint64(i + 1)
		currentReq := req
		currentReq.OrderID = currentOrderID
		
		// 创建新订单
		newOrder := &model.Order{
			ID:          currentOrderID,
			UserID:      userID,
			Title:       "Benchmark Order",
			Description: "Benchmark Description",
			Status:      model.OrderStatusPending,
			PriceCents:  1000,
		}
		
		err := orderRepo.Create(ctx, newOrder)
		if err != nil {
			b.Fatal(err)
		}
		
		// 创建支付
		_, err = service.CreatePayment(ctx, userID, currentReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetPaymentStatus 支付状态查询性能基准测试
func BenchmarkGetPaymentStatus(b *testing.B) {
	// 创建模拟仓库
	paymentRepo := newMockConcurrentPaymentRepository()
	orderRepo := newMockConcurrentOrderRepository()
	
	// 创建支付服务
	service := NewPaymentService(paymentRepo, orderRepo)
	
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	
	// 创建测试订单和支付
	order := &model.Order{
		ID:          orderID,
		UserID:      userID,
		Title:       "Benchmark Order",
		Description: "Benchmark Description",
		Status:      model.OrderStatusPending,
		PriceCents:  1000,
	}
	
	err := orderRepo.Create(ctx, order)
	if err != nil {
		b.Fatal(err)
	}
	
	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}
	
	paymentResp, err := service.CreatePayment(ctx, userID, req)
	if err != nil {
		b.Fatal(err)
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := service.GetPaymentStatus(ctx, userID, paymentResp.PaymentID)
		if err != nil {
			b.Fatal(err)
		}
	}
}