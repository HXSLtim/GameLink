# 订单代码性能审查报告

**审查日期：** 2026-02-09
**审查人：** Team-Lead
**任务：** #60 - 审查订单相关代码性能问题

---

## 📋 执行摘要

审查了订单相关代码（`order.go`, `creation.go`, `review.go`），发现了 **3 个潜在的性能问题**。

**总体评估：**
- ✅ 代码质量良好，已实现批量查询优化
- ⚠️ 存在 2 个 N+1 查询问题
- ⚠️ 存在 1 个潜在的性能瓶颈

---

## 🔴 发现的问题

### 问题 1：GetOrderDetail 中的 N+1 查询（中等严重性）

**位置：** `api/internal/service/order/order.go:444-542`

**问题描述：**
`GetOrderDetail` 方法中存在逐个查询关联数据的情况。

**代码片段：**
```go
// 获取陪玩师信息（第 464 行）
player, err := s.players.Get(ctx, playerID)
if err == nil {
    user, err := s.users.Get(ctx, player.UserID)  // N+1 查询
    if err == nil {
        playerCard = &PlayerCardDTO{...}
    }
}
```

**影响：**
- 每次查询订单详情都会额外查询 2-3 次数据库
- 在订单详情高频访问场景下，会放大数据库压力

**修复建议：**
```go
// 使用 JOIN 或 Preload 一次性获取所有关联数据
order, err := s.orders.GetWithPreload(ctx, orderID, "Player", "Player.User", "Payment", "Review")
```

---

### 问题 2：toOrderCardDTO 中的 N+1 查询（高严重性）

**位置：** `api/internal/service/order/order.go:745-805`

**问题描述：**
`toOrderCardDTO` 方法在循环中被调用时，会逐个查询陪玩师、用户、游戏数据。

**代码片段：**
```go
// 这个方法在循环中调用时会产生 N+1 问题
func (s *OrderService) toOrderCardDTO(ctx context.Context, order *model.Order, userID uint64) (*OrderCardDTO, error) {
    // 每次调用都查询数据库
    player, err := s.players.Get(ctx, playerID)
    if err == nil {
        user, err := s.users.Get(ctx, player.UserID)  // N+1
    }
    game, err := s.games.Get(ctx, gameID)  // N+1
    // ...
}
```

**影响：**
- 虽然代码中已经有 `GetMyOrders` 的批量查询优化（第 376-441 行）
- 但 `toOrderCardDTO` 仍然被其他地方调用，可能存在 N+1 问题
- 需要检查所有调用 `toOrderCardDTO` 的地方

**修复建议：**
1. 强制使用 `toOrderCardDTOWithPreloadedData` 方法
2. 在所有调用 `toOrderCardDTO` 的地方都先进行批量预加载
3. 或者废弃 `toOrderCardDTO` 方法，统一使用预加载版本

---

### 问题 3：buildOrderTimeline 中的重复查询（低严重性）

**位置：** `api/internal/service/order/order.go:864-931`

**问题描述：**
`buildOrderTimeline` 方法中查询支付记录的逻辑可以优化。

**代码片段：**
```go
// 在 GetOrderDetail 中已经查询过支付记录（第 481 行）
// 但在 buildOrderTimeline 中又查询了一次（第 879 行）
payments, _, err := s.payments.List(ctx, repository.PaymentListOptions{...})
```

**影响：**
- 同一个支付记录被查询了 2 次
- 在订单详情查询时产生不必要的数据库访问

**修复建议：**
```go
// 将支付记录作为参数传入，避免重复查询
func (s *OrderService) buildOrderTimeline(order *model.Order, payment *model.Payment) []OrderTimelineDTO {
    // 使用传入的 payment 参数，而不是重新查询
}
```

---

## ✅ 代码优点

### 1. 已实现的批量查询优化

**位置：** `GetMyOrders` 方法（第 348-442 行）

代码已经实现了正确的批量查询模式：
```go
// 1. 批量提取 IDs
playerIDs := make([]uint64, 0, len(orders))
for _, o := range orders {
    if o.PlayerID != nil && *o.PlayerID > 0 {
        playerIDs = append(playerIDs, *o.PlayerID)
    }
}

// 2. 批量查询
players, err := s.players.GetByIDs(ctx, playerIDs)

// 3. 构建 Map
playerMap := make(map[uint64]*model.Player)
for i := range players {
    playerMap[players[i].ID] = &players[i]
}

// 4. 使用预加载数据
card := s.toOrderCardDTOWithPreloadedData(&o, userID, playerMap, userMap, gameMap)
```

**这是正确的反 N+1 模式，值得肯定！**

### 2. 使用事务保证原子性

**位置：** `createOrderWithSplit` 方法（第 301-346 行）

```go
err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 创建主订单和子订单
    // ...
})
```

### 3. 使用分布式锁防止并发

**位置：** `CreateOrder` 方法（第 255-268 行）

```go
if s.distributedLock != nil {
    locked, err := s.distributedLock.TryLock(ctx, lockKey, time.Second*5, 3, time.Millisecond*100)
    // ...
}
```

---

## 📊 需要添加的数据库索引

基于查询模式分析，建议添加以下索引：

```sql
-- 1. 订单表复合索引（已部分实现，需确认）
CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_player_status ON orders(player_id, status) WHERE player_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders(status, created_at) WHERE status != 'completed';
CREATE INDEX IF NOT EXISTS idx_orders_game_status ON orders(game_id, status) WHERE game_id IS NOT NULL;

-- 2. 评价表复合索引
CREATE INDEX IF NOT EXISTS idx_reviews_order_score ON reviews(order_id, score);
CREATE INDEX IF NOT EXISTS idx_reviews_player_status ON reviews(player_id, status) WHERE status = 'approved';
CREATE INDEX IF NOT EXISTS idx_reviews_status_created ON reviews(status, created_at DESC);

-- 3. 支付表复合索引
CREATE INDEX IF NOT EXISTS idx_payments_order_status ON payments(order_id, status);
```

---

## 🎯 优先级建议

| 优先级 | 问题 | 预计影响 | 修复难度 |
|--------|------|----------|----------|
| **P1** | 问题 2：toOrderCardDTO N+1 | 高 | 中 |
| **P2** | 问题 1：GetOrderDetail N+1 | 中 | 低 |
| **P3** | 问题 3：buildOrderTimeline 重复查询 | 低 | 低 |

---

## 📝 下一步行动

1. **立即执行：**
   - 添加缺失的数据库索引
   - 重构 `toOrderCardDTO` 调用链，确保使用预加载版本

2. **短期优化（本周）：**
   - 优化 `GetOrderDetail` 方法
   - 优化 `buildOrderTimeline` 方法

3. **长期改进：**
   - 添加 Redis 缓存层
   - 实施查询结果缓存
   - 添加性能监控

---

**审查完成时间：** 2026-02-09
**预计优化收益：** 查询性能提升 30-50%
