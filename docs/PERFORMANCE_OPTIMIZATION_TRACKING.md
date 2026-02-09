# 后端性能优化追踪记录

**负责人：** Backend-Lead
**任务：** #53 - 后端代码优化和功能开发
**开始时间：** 2026-02-09
**状态：** 进行中

---

## 审查记录

### 2026-02-09 - 初始审查

#### 审查范围
- `internal/service/` - 业务逻辑层
- `internal/repository/` - 数据访问层
- `internal/handler/` - HTTP 处理层

#### 审查方法
1. 搜索数据库查询模式（.Find, .Where, .Preload）
2. 检查循环中的查询（N+1 问题）
3. 检查缺少索引的查询条件
4. 检查缓存使用情况

#### 发现 1：数据库查询分散

**位置：** 105 个文件包含数据库查询

**分析：**
- 查询分散在多个 Repository 层
- 需要检查是否有重复查询
- 需要检查是否有批量查询机会

**建议：**
- ✅ 审查 Repository 层的实现
- ✅ 添加批量查询方法
- ✅ 优化常用查询路径

#### 发现 2：需要添加索引

**建议的索引：**
```sql
-- 订单表索引
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
CREATE INDEX idx_orders_player_status ON orders(player_id, status) WHERE player_id IS NOT NULL;
CREATE INDEX idx_orders_status_created ON orders(status, created_at) WHERE status != 'completed';

-- 支付表索引
CREATE INDEX idx_payments_order_status ON payments(order_id, status);
CREATE INDEX idx_payments_user_status ON payments(user_id, status);
CREATE INDEX idx_payments_method_status ON payments(method, status) WHERE status = 'pending';

-- 聊天消息表索引
CREATE INDEX idx_chat_messages_group_time ON chat_messages(group_id, created_at DESC);
CREATE INDEX idx_chat_messages_group_unread ON chat_messages(group_id, is_read) WHERE is_read = false;
```

#### 发现 3：缓存机会

**需要缓存的热点数据：**
1. 用户信息（TTL: 5分钟）
2. 陪玩师信息（TTL: 5分钟）
3. 游戏列表（TTL: 10分钟）
4. 排行榜（TTL: 1分钟）

**建议：**
- 使用 Redis 作为缓存层
- 实现缓存穿透保护
- 实现缓存雪崩保护

#### 发现 4：测试覆盖率

**当前测试覆盖率：** 约 60%（估算）

**目标：** 提升到 80%+

**需要添加测试的模块：**
- Service 层业务逻辑
- Repository 层数据访问
- Handler 层请求处理

---

## 待审查文件列表

### 高优先级（核心业务）

1. `internal/service/order/order.go` - 订单服务
2. `internal/service/user/user.go` - 用户服务
3. `internal/service/player/` - 陪玩师服务
4. `internal/service/payment/payment.go` - 支付服务

### 中优先级（频繁调用）

5. `internal/repository/order/` - 订单仓储
6. `internal/repository/user/` - 用户仓储
7. `internal/repository/player/` - 陪玩师仓储

### 低优先级（偶尔调用）

8. 其他 Service 层
9. 其他 Repository 层

---

## 优化建议

### 短期（1-2 天）

1. **添加数据库索引**（30分钟）
   - 审查和执行上面的索引创建语句
   - 验证索引是否被使用

2. **审查 N+1 查询**（2-3小时）
   - 逐个审查核心服务文件
   - 识别循环中的查询
   - 提出优化方案

3. **添加缓存层**（2-3小时）
   - 定义缓存接口
   - 实现 Redis 缓存
   - 选择热点数据进行缓存

### 中期（3-5 天）

4. **重构优化代码**（2-3天）
   - 优化查询逻辑
   - 实现批量查询
   - 预加载关联数据

5. **添加单元测试**（2-3天）
   - 核心业务逻辑测试
   - 数据访问层测试
   - 提高覆盖率到 80%+

### 长期（1 周+）

6. **性能监控**
   - 添加性能监控
   - 建立性能基准
   - 持续优化

---

## 下一步行动

- [x] 创建性能优化计划
- [x] 创建追踪记录文档
- [x] 开始代码审查
- [x] 识别性能瓶颈 - **发现严重 N+1 查询问题**
- [ ] 实施优化方案

---

## 🔴 严重性能问题：N+1 查询（2026-02-09 更新）

### 问题 1：批量订单操作 - orderBatch.go

**位置：** `internal/service/admin/orderBatch.go`

**影响函数：**
- `BatchCancelOrders()` (line 28-95)
- `BatchConfirmOrders()` (line 98-152)
- `BatchCompleteOrders()` (line 155-209)
- `BatchRefundOrders()` (line 220-282)
- `BatchUpdateOrderStatus()` (line 321-388)
- `BatchAssignOrders()` (line 391-432)

**问题代码模式：**
```go
for _, orderID := range orderIDs {
    order, err := s.orders.Get(ctx, orderID)  // 🔴 N+1 查询！
    if err != nil {
        // error handling
        continue
    }
    // process order...
}
```

**性能影响：**
- 批量 100 个订单 = 100 次独立数据库查询
- 预估总耗时：1000-5000ms（1-5秒）
- 优化后：~50-100ms（单次批量查询）
- **性能提升：10-50倍** ⚡

**根本原因：**
- ❌ OrderRepository 接口缺少 `GetByIDs(ctx, []uint64)` 方法
- ✅ Game、User、Player 仓库已有 `GetByIDs()` 方法
- ⚠️ Order 仓库需要实现批量查询

**优化方案：**

1. **添加 OrderRepository 批量查询接口：**
```go
// internal/repository/interfaces/order.go
type OrderReader interface {
    Get(ctx context.Context, id uint64) (*model.Order, error)
    GetByIDs(ctx context.Context, ids []uint64) ([]model.Order, error)  // 新增
}
```

2. **实现批量查询方法：**
```go
// internal/repository/implementations/order.go
func (r *gormOrderRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Order, error) {
    var orders []model.Order
    err := r.db.WithContext(ctx).
        Where("id IN ?", ids).
        Find(&orders).Error
    return orders, err
}
```

3. **重构批量操作函数：**
```go
// 优化前
for _, orderID := range orderIDs {
    order, err := s.orders.Get(ctx, orderID)  // N+1
    // ...
}

// 优化后
orders, err := s.orders.GetByIDs(ctx, orderIDs)  // 批量查询
if err != nil {
    return nil, err
}
orderMap := make(map[uint64]*model.Order)
for i := range orders {
    orderMap[orders[i].ID] = &orders[i]
}
for _, orderID := range orderIDs {
    order, exists := orderMap[orderID]
    if !exists {
        // handle not found
        continue
    }
    // process order...
}
```

### 问题 2：其他发现 - 29 个服务文件存在类似问题

**发现位置：**
- `admin.go`: 支付批量操作（4处）
- `roleService.go`: 用户角色批量操作（2处）
- `adminFeed.go`: 内容批量删除（2处）
- `menuService.go`: 菜单批量操作（2处）
- `playerBatch.go`: 陪玩师批量操作（3处）
- `gamecategory/service.go`: 分类批量操作（1处）
- `presence/service.go`: 在线状态批量查询（1处）

**影响评估：**
- 总计：15+ 批量操作函数存在 N+1 查询问题
- 预估总性能损失：每操作浪费 90-98% 数据库查询时间

---

## 优化优先级

### P0 - 立即优化（今天）

1. ✅ **审查完成** - 识别 N+1 查询问题
2. ⏳ **实现 OrderRepository.GetByIDs()**
3. ⏳ **重构 orderBatch.go 中的 6 个批量操作函数**
4. ⏳ **添加单元测试验证优化效果**

### P1 - 明天优化

5. ⏳ **优化 admin.go 中的支付批量操作**
6. ⏳ **优化 playerBatch.go 中的陪玩师批量操作**
7. ⏳ **优化其他发现的批量操作**

### P2 - 本周完成

8. ⏳ **添加数据库索引**
9. ⏳ **实现热点数据缓存**
10. ⏳ **性能测试验证**

---

**状态：** 🔴 发现严重性能问题 - 准备实施优化
**最后更新：** 2026-02-09 10:30
**预计性能提升：** 批量操作性能提升 10-50倍
