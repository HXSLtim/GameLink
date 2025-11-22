package order

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"

	"github.com/stretchr/testify/require"
)

// BenchmarkCreateOrder 订单创建性能基准测试
func BenchmarkCreateOrder(b *testing.B) {
	// 创建模拟仓库
	orderRepo := &mockOrderRepository{}
	playerRepo := &mockPlayerRepository{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	ctx := context.Background()
	userID := uint64(1)
	
	// 准备测试数据
	scheduledStart := time.Now().Add(time.Hour)
	req := CreateOrderRequest{
		PlayerID:       1,
		GameID:         1,
		Title:          "Benchmark Order",
		Description:    "Benchmark Description",
		ScheduledStart: &scheduledStart,
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

// BenchmarkAcceptOrder 接单性能基准测试
func BenchmarkAcceptOrder(b *testing.B) {
	// 创建模拟仓库
	orderRepo := &mockConcurrentOrderRepository{}
	playerRepo := &mockPlayerRepositoryForConcurrency{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	ctx := context.Background()
	
	// 添加陪玩师
	playerRepo.AddPlayer(1, 100, "Benchmark Player")
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		// 创建订单
		scheduledStart := time.Now().Add(time.Hour)
		order := &model.Order{
			UserID:         1,
			Title:          "Benchmark Order",
			Description:    "Benchmark Description",
			Status:         model.OrderStatusConfirmed,
			ScheduledStart: &scheduledStart,
			TotalPriceCents: 1000,
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

// BenchmarkGetOrder 订单查询性能基准测试
func BenchmarkGetOrder(b *testing.B) {
	// 创建模拟仓库
	orderRepo := &mockOrderRepository{}
	playerRepo := &mockPlayerRepository{}
	userRepo := &mockUserRepository{}
	gameRepo := &mockGameRepository{}
	paymentRepo := &mockPaymentRepository{}
	reviewRepo := &mockReviewRepository{}
	commissionRepo := &mockCommissionRepository{}
	
	// 创建订单服务
	service := NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	
	ctx := context.Background()
	userID := uint64(1)
	
	// 预创建一些订单
	for i := 0; i < 100; i++ {
		scheduledStart := time.Now().Add(time.Hour)
		req := CreateOrderRequest{
			PlayerID:       1,
			GameID:         1,
			Title:          "Test Order",
			Description:    "Test Description",
			ScheduledStart: &scheduledStart,
		}
		
		_, err := service.CreateOrder(ctx, userID, req)
		require.NoError(b, err)
	}
	
	// 重置基准测试计时器
	b.ResetTimer()
	
	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := service.GetMyOrders(ctx, userID, MyOrderListRequest{
			Page:     1,
			PageSize: 20,
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}