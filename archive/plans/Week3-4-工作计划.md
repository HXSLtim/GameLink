# 📋 Week 3-4 工作计划

制定时间: 2025-11-10 15:30  
目标: 覆盖率 49.5% → 60%

---

## 🎯 总体目标

### 覆盖率提升
- **当前**: 49.5%
- **目标**: 60%
- **提升**: +10.5%

### 测试策略调整
- **之前**: 轻量级测试，专注结构验证
- **现在**: 深度测试，使用mock service，测试业务逻辑
- **方法**: 真实的集成测试 + 业务场景测试

---

## 📊 当前状态分析

### 测试数量
- 当前: 184个
- 目标: 220-240个
- 新增: 36-56个

### 覆盖率瓶颈
1. **Handler层**: 45% → 需要mock service深度测试
2. **Ranking模块**: 0% → 完全未覆盖
3. **集成测试**: 框架已建立，需要实现真实逻辑

---

## 📅 Week 3-4 详细计划

### Week 3: 深度Handler测试 (20-25个)

#### Day 1-2: User Handler深度测试 (8-10个)
```
目标: 使用mock service测试实际业务逻辑

文件: user/*_test.go (扩展)
- order_test.go: 深度测试下单流程
- payment_test.go: 深度测试支付流程
- review_test.go: 深度测试评价流程

预期覆盖率: +2-3%
```

#### Day 3-4: Player Handler深度测试 (8-10个)
```
文件: player/*_test.go (扩展)
- order_test.go: 深度测试接单流程
- earnings_test.go: 深度测试收益计算
- profile_test.go: 深度测试资料更新

预期覆盖率: +2-3%
```

#### Day 5: Admin Handler深度测试 (4-5个)
```
文件: admin/*_test.go (扩展)
- order_test.go: 深度测试订单管理
- withdraw_test.go: 深度测试提现审核

预期覆盖率: +1-2%
```

### Week 4: Ranking模块 + 真实集成测试 (16-31个)

#### Day 1-2: Ranking Service测试 (10-15个)
```
文件: service/ranking/*_test.go (新建)
- 排名计算测试
- 排名更新测试
- 排名查询测试
- 边界情况测试

预期覆盖率: +3-4%
```

#### Day 3-4: 真实集成测试 (6-10个)
```
文件: service/integration_test.go (重写)
- 完整下单→支付→完成流程
- 完整提现→审核→到账流程
- 完整评价→更新评分流程
- 错误处理和回滚测试

预期覆盖率: +2-3%
```

#### Day 5: 边界情况和性能测试 (0-6个)
```
- 并发测试
- 大数据量测试
- 边界值测试

预期覆盖率: +0-1%
```

---

## 🔧 技术方法

### 1. 使用Mock Service
```go
// 示例：深度测试Handler
func TestUserOrder_CreateWithMockService(t *testing.T) {
    mockOrderService := new(MockOrderService)
    mockOrderService.On("CreateOrder", mock.Anything, mock.Anything).
        Return(&order.CreateOrderResponse{OrderID: 123}, nil)
    
    handler := NewOrderHandler(mockOrderService)
    // 测试实际业务逻辑
}
```

### 2. 真实集成测试
```go
// 示例：多层协作测试
func TestIntegration_CompleteOrderFlow(t *testing.T) {
    // 创建真实的repository和service
    orderRepo := setupTestOrderRepo()
    paymentRepo := setupTestPaymentRepo()
    
    orderSvc := order.NewOrderService(orderRepo, ...)
    paymentSvc := payment.NewPaymentService(paymentRepo, ...)
    
    // 测试完整流程
    // 1. 创建订单
    // 2. 创建支付
    // 3. 支付回调
    // 4. 订单完成
}
```

### 3. 业务场景测试
```go
// 示例：测试实际业务场景
func TestBusinessScenario_UserOrderAndPay(t *testing.T) {
    // 场景：用户下单并支付
    // 1. 用户选择服务
    // 2. 创建订单（检查库存、价格）
    // 3. 创建支付（检查金额）
    // 4. 支付成功（更新订单状态）
    // 5. 通知玩家（发送消息）
}
```

---

## 📈 预期成果

### 覆盖率提升路径
| Week | 新增测试 | 覆盖率 | 提升 |
|------|---------|--------|------|
| Week 1-2 | 184个 | 49.5% | 基准 |
| Week 3 | +20-25个 | 54-55% | +4.5-5.5% |
| Week 4 | +16-31个 | 58-60% | +4-5% |
| **总计** | **220-240个** | **58-60%** | **+8.5-10.5%** |

### 分层覆盖率预期
| 层级 | 当前 | Week 3 | Week 4 | 提升 |
|------|------|--------|--------|------|
| Handler | 45% | 52-55% | 55-60% | +10-15% |
| Service | 72% | 75-78% | 78-82% | +6-10% |
| Repository | 84% | 85-86% | 86-88% | +2-4% |

---

## 🎯 成功标准

### 必达目标 ✅
- [ ] 新增测试 ≥ 36个
- [ ] 覆盖率 ≥ 58%
- [ ] 通过率 = 100%
- [ ] Ranking模块覆盖率 > 0%

### 期望目标 🎯
- [ ] 新增测试 ≥ 45个
- [ ] 覆盖率 ≥ 60%
- [ ] Handler层 ≥ 55%
- [ ] Service层 ≥ 78%

### 超额目标 🌟
- [ ] 新增测试 ≥ 56个
- [ ] 覆盖率 ≥ 62%
- [ ] Handler层 ≥ 60%
- [ ] Ranking模块 ≥ 50%

---

## 💡 关键策略

### 1. 深度优先
- 不追求测试数量
- 专注测试质量和深度
- 每个测试都要测试实际业务逻辑

### 2. Mock Service
- 使用testify/mock
- Mock所有外部依赖
- 测试各种场景（成功、失败、边界）

### 3. 真实集成
- 创建测试数据库
- 使用真实的repository
- 测试完整的业务流程

### 4. 业务场景
- 从用户角度设计测试
- 测试实际使用场景
- 覆盖关键业务路径

---

## ⚠️ 注意事项

### 避免的问题
1. ❌ 不要创建太多轻量级测试
2. ❌ 不要忽略业务逻辑
3. ❌ 不要跳过错误处理测试
4. ❌ 不要忽略边界情况

### 应该做的
1. ✅ 使用mock service
2. ✅ 测试业务逻辑
3. ✅ 测试错误处理
4. ✅ 测试边界情况
5. ✅ 创建真实集成测试

---

## 🚀 立即开始

### 第一步
创建深度测试框架和mock结构

### 第二步
扩展User Handler测试，使用mock service

### 第三步
创建Ranking模块测试

### 第四步
重写集成测试，使用真实逻辑

---

**Week 3-4 开始！** 🚀  
**目标60%覆盖率！** 📈  
**深度测试，质量优先！** 💪
