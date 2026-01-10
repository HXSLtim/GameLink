# UnitOfWork 使用指南

## 概述

`pkg/uow/work.go` 提供了工作单元模式，用于管理跨多个 Repository 的事务。

## 基本使用

### 1. 简单事务操作

```go
import (
    "gamelink/pkg/uow"
    "gorm.io/gorm"
)

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    var order *Order

    // 使用工作单元确保所有操作在同一事务中
    err := uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
        // 通过工作单元获取仓储
        orderRepo := repository.NewOrderRepository(w.DB())
        paymentRepo := repository.NewPaymentRepository(w.DB())

        // 创建订单
        order = &Order{
            UserID:          req.UserID,
            PlayerID:        req.PlayerID,
            TotalPriceCents: req.TotalPriceCents,
            Status:          model.OrderStatusPending,
        }
        if err := orderRepo.Create(ctx, order); err != nil {
            return fmt.Errorf("create order: %w", err)
        }

        // 创建支付记录
        payment := &Payment{
            OrderID:       order.ID,
            AmountCents:   req.TotalPriceCents,
            PaymentMethod: req.PaymentMethod,
            Status:        model.PaymentStatusPending,
        }
        if err := paymentRepo.Create(ctx, payment); err != nil {
            return fmt.Errorf("create payment: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return order, nil
}
```

### 2. 只读操作（不需要事务）

```go
func (s *OrderService) GetOrder(ctx context.Context, id uint64) (*Order, error) {
    var order *Order

    err := uow.NewWork(ctx, s.db).Execute(func(w *uow.Work) error {
        orderRepo := repository.NewOrderRepository(w.DB())
        var err error
        order, err = orderRepo.Get(ctx, id)
        return err
    })

    return order, err
}
```

### 3. 使用辅助函数

```go
// 使用 WithTransaction 简化事务操作
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint64, status string) error {
    return uow.WithTransaction(ctx, s.db, func(ctx context.Context, tx *gorm.DB) error {
        orderRepo := repository.NewOrderRepository(tx)

        // 更新订单状态
        return orderRepo.UpdateStatus(ctx, orderID, status)
    })
}

// 使用 TransactionScope
func (s *OrderService) ProcessOrder(ctx context.Context, orderID uint64) error {
    return uow.TransactionScope(ctx, s.db, func(ctx context.Context) error {
        // 所有数据库操作在这个作用域内是事务性的
        orderRepo := repository.NewOrderRepository(/* ... */)
        paymentRepo := repository.NewPaymentRepository(/* ... */)

        // ... 业务逻辑

        return nil
    })
}
```

### 4. 带错误处理的完整示例

```go
func (s *OrderService) CreateOrderWithCommission(
    ctx context.Context,
    req CreateOrderRequest,
) (*Order, error) {
    var order *Order

    err := uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
        orderRepo := repository.NewOrderRepository(w.DB())
        paymentRepo := repository.NewPaymentRepository(w.DB())
        commissionRepo := repository.NewCommissionRepository(w.DB())

        // 1. 创建订单
        order = &Order{
            UserID:          req.UserID,
            PlayerID:        req.PlayerID,
            TotalPriceCents: req.TotalPriceCents,
            Status:          model.OrderStatusPending,
        }
        if err := orderRepo.Create(ctx, order); err != nil {
            return fmt.Errorf("create order: %w", err)
        }

        // 2. 创建支付记录
        payment := &Payment{
            OrderID:       order.ID,
            AmountCents:   req.TotalPriceCents,
            PaymentMethod: req.PaymentMethod,
            Status:        model.PaymentStatusPending,
        }
        if err := paymentRepo.Create(ctx, payment); err != nil {
            return fmt.Errorf("create payment: %w", err)
        }

        // 3. 计算并记录佣金
        commissionRate := 0.20 // 从配置获取
        commission := &Commission{
            OrderID:          order.ID,
            PlayerID:         order.PlayerID,
            BaseRate:         commissionRate,
            AmountCents:      int64(float64(req.TotalPriceCents) * commissionRate),
            SettlementStatus: model.SettlementStatusPending,
        }
        if err := commissionRepo.Create(ctx, commission); err != nil {
            return fmt.Errorf("record commission: %w", err)
        }

        // 4. 更新陪玩师统计
        playerRepo := repository.NewPlayerRepository(w.DB())
        if err := playerRepo.IncrementOrderCount(ctx, order.PlayerID); err != nil {
            return fmt.Errorf("update player stats: %w", err)
        }

        return nil
    })

    if err != nil {
        // 事务已自动回滚
        return nil, fmt.Errorf("create order: %w", err)
    }

    return order, nil
}
```

## 高级用法

### 1. 嵌套事务处理

```go
func (s *OrderService) BatchCreateOrders(
    ctx context.Context,
    requests []CreateOrderRequest,
) ([]*Order, error) {
    var orders []*Order

    err := uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
        orderRepo := repository.NewOrderRepository(w.DB())

        for _, req := range requests {
            order := &Order{
                UserID:          req.UserID,
                PlayerID:        req.PlayerID,
                TotalPriceCents: req.TotalPriceCents,
                Status:          model.OrderStatusPending,
            }
            if err := orderRepo.Create(ctx, order); err != nil {
                return fmt.Errorf("create order for user %d: %w", req.UserID, err)
            }
            orders = append(orders, order)
        }

        return nil
    })

    return orders, err
}
```

### 2. 条件性事务

```go
func (s *OrderService) ProcessOrder(ctx context.Context, orderID uint64, force bool) error {
    if force {
        // 强制模式：使用事务
        return uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
            // ... 事务性操作
            return nil
        })
    }

    // 普通模式：不使用事务
    return uow.NewWork(ctx, s.db).Execute(func(w *uow.Work) error {
        // ... 非事务性操作
        return nil
    })
}
```

### 3. 保存点（Savepoint）

GORM 支持保存点，可以在事务中创建部分回滚点：

```go
func (s *OrderService) BatchProcessWithPartialFailure(
    ctx context.Context,
    orderIDs []uint64,
) error {
    return uow.NewWork(ctx, s.db).Commit(func(w *uow.Work) error {
        orderRepo := repository.NewOrderRepository(w.DB())

        for i, orderID := range orderIDs {
            // 创建保存点
            savepoint := fmt.Sprintf("sp_%d", i)

            err := w.DB().Transaction(func(tx *gorm.DB) error {
                // 尝试处理订单
                if err := orderRepo.Process(ctx, orderID); err != nil {
                    // 回滚到保存点
                    tx.RollbackTo(savepoint)
                    // 记录错误但继续处理其他订单
                    log.Printf("Order %d processing failed: %v", orderID, err)
                    return nil
                }
                return nil
            })

            if err != nil {
                return err
            }
        }

        return nil
    })
}
```

## 迁移现有代码

### 迁移前（无事务管理）

```go
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    // 创建订单
    order, err := s.orderRepo.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    // 创建支付（如果失败，订单已创建但支付未创建，数据不一致！）
    _, err = s.paymentRepo.Create(ctx, order.ID, req.Amount)
    if err != nil {
        return nil, err
    }

    return order, nil
}
```

### 迁移后（使用工作单元）

```go
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    var order *Order

    err := uow.ExecuteInUnitOfWork(ctx, s.db, func(w *uow.Work) error {
        orderRepo := repository.NewOrderRepository(w.DB())
        paymentRepo := repository.NewPaymentRepository(w.DB())

        // 创建订单
        var err error
        order, err = orderRepo.Create(ctx, req)
        if err != nil {
            return err
        }

        // 创建支付（如果失败，事务会自动回滚订单）
        _, err = paymentRepo.Create(ctx, order.ID, req.Amount)
        return err
    })

    return order, err
}
```

## 注意事项

1. **不要在事务中执行长时间操作**
   - 事务会锁定数据库资源
   - 避免在事务中调用外部 API
   - 避免在事务中进行繁重的计算

2. **保持事务简短**
   - 只在必要的时候使用事务
   - 事务范围应尽可能小
   - 完成数据库操作后立即提交

3. **错误处理**
   - 事务中返回错误会自动回滚
   - 返回 nil 会自动提交
   - 不要在事务中使用 `defer` 来处理错误

4. **只读查询**
   - 只读查询不需要使用事务
   - 使用 `Execute` 方法而不是 `Commit`

5. **嵌套事务**
   - GORM 不支持真正的嵌套事务
   - 内层事务实际上是保存点
   - 建议扁平化事务结构

## 测试

```go
func TestOrderService_CreateOrder(t *testing.T) {
    db := setupTestDB(t)
    service := NewOrderService(db)

    ctx := context.Background()

    t.Run("成功创建订单", func(t *testing.T) {
        req := CreateOrderRequest{
            UserID:          1,
            PlayerID:        2,
            TotalPriceCents: 10000,
        }

        order, err := service.CreateOrder(ctx, req)
        assert.NoError(t, err)
        assert.NotNil(t, order)
        assert.Equal(t, uint64(1), order.UserID)

        // 验证数据库状态
        // 由于在事务中，这里需要单独查询
        savedOrder, err := service.GetOrder(ctx, order.ID)
        assert.NoError(t, err)
        assert.Equal(t, order.Status, savedOrder.Status)
    })

    t.Run("支付创建失败时订单回滚", func(t *testing.T) {
        // 模拟支付失败场景
        req := CreateOrderRequest{
            UserID:          1,
            PlayerID:        999, // 不存在的陪玩师
            TotalPriceCents: 10000,
        }

        _, err := service.CreateOrder(ctx, req)
        assert.Error(t, err)

        // 验证订单未被创建（事务已回滚）
        orders, _ := service.ListOrders(ctx, 1)
        assert.Empty(t, orders)
    })
}
```

## 环境变量

以下环境变量可以控制工作单元行为：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `DB_TX_TIMEOUT` | 30s | 事务超时时间 |
| `DB_TX_RETRY` | 0 | 事务重试次数 |
| `DB_LOG_SQL` | false | 是否记录 SQL 日志 |

## 最佳实践

1. **明确事务边界**
   - 一次事务只完成一个业务操作
   - 不要在一个事务中处理多个不相关的操作

2. **避免大事务**
   - 事务越大，锁定时间越长
   - 考虑拆分为多个小事务

3. **使用合适的隔离级别**
   - 默认使用 READ COMMITTED
   - 特殊场景可使用 SERIALIZABLE

4. **监控事务性能**
   - 记录慢事务
   - 设置事务超时
   - 监控死锁
