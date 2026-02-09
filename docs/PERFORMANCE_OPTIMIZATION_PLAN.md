# GameLink 后端性能优化计划

**创建者：** Backend-Lead
**日期：** 2026-02-09
**任务：** #53 - 后端代码优化和功能开发
**优先级：** P1

---

## 执行摘要

本文档记录了 GameLink 后端系统的性能优化计划和实施步骤。

**目标：**
- API 响应时间减少 20-30%
- 数据库查询效率提升 30-50%
- 系统吞吐量提升
- 代码测试覆盖率提升到 80%+

---

## 1. 性能分析

### 1.1 当前状态

**数据库：** PostgreSQL
**ORM：** GORM
**缓存：** Redis
**架构：** 分层架构（Handler → Service → Repository）

**已识别的潜在问题：**
- 105 个文件包含数据库查询操作
- 可能存在 N+1 查询问题
- 可能缺少必要的数据库索引
- 可能存在未优化的关联查询

### 1.2 性能瓶颈识别

**第一阶段：代码审查**

需要审查的文件类型：
1. **Service 层**：`internal/service/**/*.go`
2. **Repository 层**：`internal/repository/**/*.go`
3. **Handler 层**：`internal/handler/**/*.go`

**重点关注的模式：**
- N+1 查询：循环中执行查询
- 未预加载的关联数据
- 缺少索引的查询条件
- 未使用缓存的重复查询

---

## 2. 优化方案

### 2.1 数据库查询优化

#### A. 索引优化

**审查现有索引：**
```sql
-- 检查现有索引
SELECT
    tablename,
    indexname,
    indexdef
FROM
    pg_indexes
WHERE
    schemaname = 'public'
ORDER BY
    tablename, indexname;
```

**需要添加的索引：**
```sql
-- 订单表索引
CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_player_status ON orders(player_id, status) WHERE player_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders(status, created_at) WHERE status != 'completed';
CREATE INDEX IF NOT EXISTS idx_orders_game_status ON orders(game_id, status) WHERE game_id IS NOT NULL;

-- 支付表索引
CREATE INDEX IF NOT EXISTS idx_payments_order_status ON payments(order_id, status);
CREATE INDEX IF NOT EXISTS idx_payments_user_status ON payments(user_id, status);
CREATE INDEX IF NOT EXISTS idx_payments_method_status ON payments(method, status) WHERE status = 'pending';

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, created_at);

-- 聊天消息表索引
CREATE INDEX IF NOT EXISTS idx_chat_messages_group_time ON chat_messages(group_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_group_unread ON chat_messages(group_id, is_read) WHERE is_read = false;
```

#### B. 查询优化

**预加载关联数据：**
```go
// 优化前（N+1 查询）
orders, _ := s.orders.List(ctx, opts)
for _, order := range orders {
    user, _ := s.users.Get(ctx, order.UserID)  // N+1 问题
    player, _ := s.players.Get(ctx, order.PlayerID)  // N+1 问题
    // ...
}

// 优化后（使用 Joins 或 Preload）
orders, _ := s.orders.ListWithPreload(ctx, opts, "User", "Player")
```

**批量查询：**
```go
// 优化前
for _, orderID := range orderIDs {
    order, _ := s.orders.Get(ctx, orderID)
}

// 优化后
orders, _ := s.orders.ListByIDs(ctx, orderIDs)
```

#### C. 分页优化

**使用游标分页替代 OFFSET：**
```go
// 优化前（OFFSET 在大数据集时性能差）
orders, _ := s.orders.List(ctx, ListOptions{
    Page:     1000,
    PageSize: 20,
})

// 优化后（使用游标分页）
orders, _ := s.orders.ListByCursor(ctx, cursor, 20)
```

### 2.2 缓存优化

#### A. Redis 缓存策略

**缓存热点数据：**
- 用户信息（TTL: 5分钟）
- 陪玩师信息（TTL: 5分钟）
- 游戏列表（TTL: 10分钟）
- 排行榜（TTL: 1分钟）
- 统计数据（TTL: 按需）

**缓存实现：**
```go
// 伪代码示例
func (s *Service) GetUser(ctx context.Context, id uint64) (*model.User, error) {
    // 1. 尝试从缓存获取
    cached, err := s.cache.Get(ctx, fmt.Sprintf("user:%d", id))
    if err == nil {
        var user model.User
        json.Unmarshal(cached, &user)
        return &user, nil
    }

    // 2. 从数据库获取
    user, err := s.users.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, fmt.Sprintf("user:%d", id), data, 5*time.Minute)

    return user, nil
}
```

#### B. 查询结果缓存

**缓存列表查询结果：**
```go
// 缓存首页数据（TTL: 1分钟）
func (s *Service) GetDashboardData(ctx context.Context) (*DashboardData, error) {
    cacheKey := "dashboard:summary"

    // 尝试从缓存获取
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var data DashboardData
        json.Unmarshal(cached, &data)
        return &data, nil
    }

    // 生成数据
    data := s.generateDashboardData(ctx)

    // 缓存结果
    bytes, _ := json.Marshal(data)
    s.cache.Set(ctx, cacheKey, bytes, 1*time.Minute)

    return data, nil
}
```

### 2.3 并发优化

#### A. Goroutine 池

**使用 worker pool 处理并发任务：**
```go
// 使用 ants 或 tunny
pool, _ := ants.NewPool(100, 10*time.Second)
defer pool.Release()

for _, task := range tasks {
    task := task
    pool.Submit(func() {
        processTask(task)
    })
}
```

#### B. 数据库连接池优化

**GORM 连接池配置：**
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    ConnPool: db.Resolve(10, 100, 30*time.Minute), // 最大 100 个连接
})
```

### 2.4 代码质量提升

#### A. 单元测试

**目标覆盖率：80%+**

**需要测试的模块：**
- Service 层业务逻辑
- Repository 层数据访问
- Handler 层请求处理
- 工具函数和辅助方法

**测试策略：**
```bash
# 运行测试并查看覆盖率
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### B. 代码重构

**重构目标：**
- 降低圈复杂度（目标：≤10）
- 提取重复代码
- 改进命名
- 添加注释

---

## 3. 实施计划

### 3.1 第一阶段：性能分析和优化（2-3 天）

**Day 1：**
- ✅ 审查数据库查询
- ✅ 识别 N+1 查询
- ✅ 分析慢查询日志

**Day 2：**
- ⏳ 添加数据库索引
- ⏳ 优化关联查询
- ⏳ 实施缓存策略

**Day 3：**
- ⏳ 性能测试
- ⏳ 对比优化前后效果
- ⏳ 调整优化方案

### 3.2 第二阶段：代码质量提升（2-3 天）

**Day 4-5：**
- ⏳ 添加单元测试
- ⏳ 重构复杂函数
- ⏳ 改进错误处理

**Day 6：**
- ⏳ 代码审查
- ⏳ 文档更新
- ⏳ 最终测试

### 3.3 第三阶段：安全加固（1 天）

**Day 7：**
- ⏳ SQL 注入审查
- ⏳ 输入验证检查
- ⏳ 速率限制配置审查

---

## 4. 预期成果

### 4.1 性能指标

| 指标 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| API 平均响应时间 | ~200ms | ~140-160ms | 20-30% |
| 数据库查询时间 | ~50ms | ~25-35ms | 30-50% |
| 系统吞吐量 | ~100 req/s | ~130-150 req/s | 30-50% |

### 4.2 代码质量

| 指标 | 优化前 | 优化后 |
|------|--------|--------|
| 测试覆盖率 | ~60% | ~80%+ |
| 圈复杂度 >10 的函数 | 未知 | 0 |
| 代码注释覆盖率 | ~30% | ~70% |

---

## 5. 监控和验证

### 5.1 性能监控

**使用 GORM 的日志功能：**
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

**启用慢查询日志：**
```sql
-- 在 postgresql.conf 中设置
log_min_duration_statement = 100ms  -- 记录超过 100ms 的查询
```

### 5.2 性能测试

**使用 Go 的基准测试：**
```go
func BenchmarkOrderService_ListOrders(b *testing.B) {
    service := setupTestService()
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.List(ctx, OrderListOptions{
            Page:     1,
            PageSize: 20,
        })
    }
}
```

---

## 6. 注意事项

### 6.1 避免过度优化

- 不要过早优化
- 基于性能分析数据做决策
- 优先优化热点代码路径

### 6.2 保持向后兼容

- 不修改 API 接口签名
- 不修改数据库结构
- 保持现有功能正常工作

### 6.3 渐进式优化

- 小步迭代，持续改进
- 每次优化后验证效果
- 记录优化过程和结果

---

## 7. 下一步行动

1. ✅ 创建性能优化计划（本文档）
2. ⏳ 开始代码审查
3. ⏳ 识别性能瓶颈
4. ⏳ 实施优化方案
5. ⏳ 验证优化效果

---

**文档状态：** 初始版本
**创建时间：** 2026-02-09
**最后更新：** 2026-02-09
**负责人：** Backend-Lead
