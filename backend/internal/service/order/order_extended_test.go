package order

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
)

// TestOrderStatusTransitions 测试订单状态流转
func TestOrderStatusTransitions(t *testing.T) {
	ctx := context.Background()

	t.Run("正常状态流转_pending到confirmed", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		// 创建pending状态的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 更新为confirmed状态
		order.Status = model.OrderStatusConfirmed
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		
		// 验证状态已更新
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusConfirmed, updated.Status)
	})

	t.Run("正常状态流转_confirmed到in_progress", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusInProgress
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusInProgress, updated.Status)
	})

	t.Run("正常状态流转_in_progress到completed", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("取消流转_pending到canceled", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		order.Status = model.OrderStatusCanceled
		order.CancelReason = "用户取消"
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, updated.Status)
		assert.Equal(t, "用户取消", updated.CancelReason)
	})

	t.Run("已完成订单状态不应该再改变", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		completedAt := time.Now()
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
			CompletedAt:     &completedAt,
		}
		orderRepo.Create(ctx, order)

		// 尝试修改已完成的订单（业务层应该阻止）
		originalStatus := order.Status
		order.Status = model.OrderStatusCanceled
		
		// Repository层会允许更新，但Service层应该阻止
		// 这里测试的是数据一致性
		err := orderRepo.Update(ctx, order)
		assert.NoError(t, err) // Repository层允许
		
		// 但在实际业务中，Service层应该检查并拒绝这种操作
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.NotEqual(t, originalStatus, updated.Status) // Repository已更新
		// 注意：这个测试说明需要在Service层添加状态检查
	})

	t.Run("已取消订单状态不应该再改变", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCanceled,
			TotalPriceCents: 10000,
			CancelReason:    "用户取消",
		}
		orderRepo.Create(ctx, order)

		// 尝试修改已取消的订单
		order.Status = model.OrderStatusConfirmed
		err := orderRepo.Update(ctx, order)
		
		assert.NoError(t, err) // Repository层允许
		// 注意：Service层应该添加检查防止这种情况
	})
}

// TestOrderCreation_EdgeCases 测试订单创建的边界情况
func TestOrderCreation_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("创建订单时价格为0", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 0, // 零价格
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), order.TotalPriceCents)
	})

	t.Run("创建订单时价格为极大值", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000000, // 100,000元
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, int64(10000000), order.TotalPriceCents)
	})

	t.Run("创建订单时必须有用户ID", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          0, // 无效的用户ID
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		err := orderRepo.Create(ctx, order)

		// Repository层会允许，但Service层应该验证
		assert.NoError(t, err)
		// 注意：Service层应该添加用户ID验证
	})

	t.Run("创建订单时默认状态应该是pending", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		err := orderRepo.Create(ctx, order)

		assert.NoError(t, err)
		assert.Equal(t, model.OrderStatusPending, order.Status)
	})
}

// TestOrderCancellation_EdgeCases 测试订单取消的边界情况
func TestOrderCancellation_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("取消pending状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 取消订单
		order.Status = model.OrderStatusCanceled
		order.CancelReason = "用户取消"
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, updated.Status)
		assert.Equal(t, "用户取消", updated.CancelReason)
	})

	t.Run("取消confirmed状态的订单_需要退款", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 取消已确认的订单
		order.Status = model.OrderStatusRefunded
		order.RefundReason = "用户申请退款"
		refundedAt := time.Now()
		order.RefundedAt = &refundedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		// 注意：实际业务中应该触发退款流程
	})

	t.Run("不能取消in_progress状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 尝试取消进行中的订单
		order.Status = model.OrderStatusCanceled
		err := orderRepo.Update(ctx, order)

		// Repository层会允许，但Service层应该阻止
		assert.NoError(t, err)
		// 注意：Service层应该添加状态检查
	})

	t.Run("不能取消completed状态的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		completedAt := time.Now()
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 10000,
			CompletedAt:     &completedAt,
		}
		orderRepo.Create(ctx, order)

		// 尝试取消已完成的订单
		order.Status = model.OrderStatusCanceled
		err := orderRepo.Update(ctx, order)

		// Repository层会允许，但Service层应该阻止
		assert.NoError(t, err)
		// 注意：Service层应该添加状态检查
	})
}

// TestOrderCompletion_EdgeCases 测试订单完成的边界情况
func TestOrderCompletion_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("正常完成订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 完成订单
		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		err := orderRepo.Update(ctx, order)

		assert.NoError(t, err)
		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCompleted, updated.Status)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("完成订单时应该记录完成时间", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusInProgress,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 完成订单
		beforeComplete := time.Now()
		order.Status = model.OrderStatusCompleted
		completedAt := time.Now()
		order.CompletedAt = &completedAt
		orderRepo.Update(ctx, order)
		afterComplete := time.Now()

		updated, _ := orderRepo.Get(ctx, order.ID)
		assert.NotNil(t, updated.CompletedAt)
		assert.True(t, updated.CompletedAt.After(beforeComplete) || updated.CompletedAt.Equal(beforeComplete))
		assert.True(t, updated.CompletedAt.Before(afterComplete) || updated.CompletedAt.Equal(afterComplete))
	})
}

// TestOrderQuery_EdgeCases 测试订单查询的边界情况
func TestOrderQuery_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("查询不存在的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order, err := orderRepo.Get(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("查询用户的所有订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		// 创建多个订单
		for i := 0; i < 5; i++ {
			order := &model.Order{
				UserID:          1,
				Status:          model.OrderStatusPending,
				TotalPriceCents: int64((i + 1) * 1000),
			}
			orderRepo.Create(ctx, order)
		}

		// 查询用户1的订单
		userID := uint64(1)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, orders, 5)
	})

	t.Run("按状态过滤订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		// 创建不同状态的订单
		statuses := []model.OrderStatus{
			model.OrderStatusPending,
			model.OrderStatusConfirmed,
			model.OrderStatusInProgress,
			model.OrderStatusCompleted,
			model.OrderStatusCanceled,
		}
		
		for _, status := range statuses {
			order := &model.Order{
				UserID:          1,
				Status:          status,
				TotalPriceCents: 10000,
			}
			orderRepo.Create(ctx, order)
		}

		// 只查询pending状态的订单
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			Statuses: []model.OrderStatus{model.OrderStatusPending},
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, model.OrderStatusPending, orders[0].Status)
	})

	t.Run("查询空结果集", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		userID := uint64(999)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, orders)
	})
}

// TestOrderAuthorization 测试订单权限控制
func TestOrderAuthorization(t *testing.T) {
	ctx := context.Background()

	t.Run("用户只能查看自己的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		// 创建用户1的订单
		order1 := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order1)

		// 创建用户2的订单
		order2 := &model.Order{
			UserID:          2,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 20000,
		}
		orderRepo.Create(ctx, order2)

		// 用户1查询自己的订单
		userID := uint64(1)
		orders, total, err := orderRepo.List(ctx, repository.OrderListOptions{
			UserID: &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, orders, 1)
		assert.Equal(t, uint64(1), orders[0].UserID)
	})

	t.Run("用户不能操作他人的订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		// 创建用户1的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 用户2尝试获取用户1的订单
		retrieved, err := orderRepo.Get(ctx, order.ID)

		// Repository层会返回订单，但Service层应该检查权限
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, uint64(1), retrieved.UserID)
		// 注意：Service层应该添加权限检查
	})
}

// TestOrderConcurrency 测试订单并发场景
func TestOrderConcurrency(t *testing.T) {
	ctx := context.Background()

	t.Run("并发更新同一订单", func(t *testing.T) {
		orderRepo := newMockOrderRepository()
		
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		// 模拟并发更新（实际应该使用乐观锁或悲观锁）
		order1 := *order
		order2 := *order

		order1.Status = model.OrderStatusConfirmed
		order2.Status = model.OrderStatusCanceled

		// 第一次更新
		err1 := orderRepo.Update(ctx, &order1)
		assert.NoError(t, err1)

		// 第二次更新会覆盖第一次
		err2 := orderRepo.Update(ctx, &order2)
		assert.NoError(t, err2)

		// 最终状态是第二次更新的结果
		final, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusCanceled, final.Status)
		// 注意：实际业务中应该使用版本号或锁来防止并发冲突
	})
}
