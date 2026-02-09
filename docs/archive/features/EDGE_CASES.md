# GameLink 边缘场景处理文档

> **文档版本**: v1.0
> **创建日期**: 2025-01-05
> **用途**: 描述系统边缘场景和异常处理机制

---

## 目录

1. [订单系统边缘场景](#1-订单系统边缘场景)
2. [支付系统边缘场景](#2-支付系统边缘场景)
3. [聊天系统边缘场景](#3-聊天系统边缘场景)
4. [争议处理边缘场景](#4-争议处理边缘场景)
5. [提现系统边缘场景](#5-提现系统边缘场景)
6. [系统故障场景](#6-系统故障场景)
7. [安全边缘场景](#7-安全边缘场景)

---

## 1. 订单系统边缘场景

### 1.1 订单超时未接单

**场景描述**: 用户下单支付后，订单长时间无人接单

**触发条件**:
- 订单状态为 `confirmed` (已支付)
- 超过配置的超时时间（默认30分钟）
- 仍无陪玩师接单

**处理流程**:

```
1. 系统检测到订单超时
2. 自动触发超时处理:
   - 状态变更为 canceled
   - 自动全额退款
   - 记录超时日志到 order_timeout_logs
3. 发送通知:
   - 用户: "订单超时未接单，已自动退款"
   - 管理员: 超时订单告警
```

**API响应**:
```json
{
  "success": true,
  "message": "订单超时已取消",
  "data": {
    "order_id": 12345,
    "status": "canceled",
    "refund_amount_cents": 5000,
    "reason": "timeout_no_accept"
  }
}
```

**相关代码**:
- `api/internal/service/orderTimeoutService.go` - 超时检测
- `api/internal/model/orderTimeoutLog.go` - 超时日志

---

### 1.2 陪玩师接单后失联

**场景描述**: 陪玩师接单后，长时间未开始服务或无响应

**触发条件**:
- 订单状态为 `confirmed` (已接单)
- 陪玩师未点击"开始服务"
- 超过服务开始时间15分钟

**处理流程**:

```
1. 系统检测到异常
2. 发送催单通知给陪玩师
3. 若30分钟内仍无响应:
   - 用户可选择"取消订单并退款"
   - 或"重新分配陪玩师"
4. 记录到陪玩师信誉分
```

**用户操作**:
```typescript
// 用户端 - 催单按钮
<Button onClick={handleUrgePlayer}>
  催单
</Button>

// 用户端 - 取消并退款
<Button onClick={handleCancelAndRefund}>
  取消订单并退款
</Button>
```

---

### 1.3 团队订单人员不足

**场景描述**: 用户创建team订单，但实际接单人数不足

**触发条件**:
- 订单类型为 `team`
- 已接单人数 < 订单所需人数
- 超过配置等待时间

**处理流程**:

```
1. 系统持续显示"等待更多陪玩师加入"
2. 用户可选择:
   - 继续等待
   - 降低人数要求(转solo订单)
   - 取消订单全额退款
3. 陪玩师可选择:
   - 继续等待队友
   - 独自接单(降低收费)
```

**订单转换逻辑**:
```go
// Team转Solo订单
func ConvertTeamToSolo(order *Order) error {
    if order.Type != "team" {
        return errors.New("not a team order")
    }

    // 更新订单信息
    order.Type = "solo"
    order.RequiredPlayers = 1
    // 按比例调整价格
    order.PriceCents = order.PriceCents / order.CurrentPlayers

    return nil
}
```

---

### 1.4 订单服务时间异常

**场景描述**: 实际服务时间与预订时间不符

**子场景**:

| 场景 | 触发条件 | 处理方式 |
|------|---------|---------|
| **服务时间不足** | 实际时长 < 预订时长的50% | 按实际时长计费退款差价 |
| **服务时间超长** | 实际时长 > 预订时长 | 陪玩师可选择追加收费或免费 |
| **提前结束** | 双方协商一致 | 按实际时长结算，无需争议 |
| **延时服务** | 双方协商一致 | 创建补充订单 |

**处理流程**:
```javascript
// 陪玩师端 - 完成订单时确认实际时长
const actualDuration = await confirmActualDuration();
if (actualDuration < bookedDuration * 0.5) {
  // 自动部分退款
  const refundAmount = calculateRefund(bookedDuration - actualDuration);
  await processPartialRefund(orderId, refundAmount);
}
```

---

### 1.5 重复下单

**场景描述**: 用户短时间内创建多个相同订单

**触发条件**:
- 同一用户
- 同一游戏/服务项目
- 时间间隔 < 5分钟

**处理流程**:

```
1. 检测到重复下单
2. 弹窗提示:
   "检测到您已有相同订单在进行中，是否继续?"
3. 用户选择:
   - 取消: 返回已有订单
   - 继续: 创建新订单(记录原因)
```

---

## 2. 支付系统边缘场景

### 2.1 支付超时未确认

**场景描述**: 用户发起支付后，支付平台长时间未回调

**触发条件**:
- 订单状态为 `pending`
- 发起支付 > 5分钟
- 未收到支付回调

**处理流程**:

```
1. 前端显示"支付确认中..."
2. 后端主动查询支付状态:
   - 调用支付平台查询接口
   - 若已成功: 更新订单状态
   - 若未支付: 提示用户"支付未完成，请重试"
3. 支付订单15分钟自动过期
```

**代码实现**:
```go
// 主动查询支付状态
func QueryPaymentStatus(orderID uint64) error {
    order := GetOrder(orderID)
    if time.Since(order.CreatedAt) > 15*time.Minute {
        return errors.New("payment expired")
    }

    // 查询支付平台
    status, err := paymentGateway.Query(order.PaymentID)
    if err != nil {
        return err
    }

    if status == "success" {
        return ConfirmPayment(orderID)
    }
    return nil
}
```

---

### 2.2 支付金额不匹配

**场景描述**: 支付回调金额与订单金额不一致

**触发条件**:
- 支付回调金额 ≠ 订单金额
- 可能原因: 优惠券/秒杀/恶意篡改

**处理流程**:

```
1. 验证签名和金额
2. 若金额不匹配:
   - 记录异常日志
   - 订单状态不变更
   - 发送告警给管理员
   - 人工核查处理
3. 用户端提示: "支付异常，请联系客服"
```

---

### 2.3 组合支付部分失败

**场景描述**: 用户使用余额+第三方支付组合支付，第三方支付失败

**场景**:
- 订单金额: ¥100
- 余额支付: ¥30 (成功)
- 第三方支付: ¥70 (失败)

**处理流程**:

```
1. 第三方支付失败
2. 自动退款已扣余额: ¥30
3. 订单状态保持 pending
4. 提示用户: "第三方支付失败，已退还余额到钱包"
5. 用户可重新选择支付方式
```

---

### 2.4 重复支付

**场景描述**: 用户因网络问题重复发起支付

**检测机制**:
```go
// 幂等性检查
func CheckDuplicatePayment(userID, orderID uint64) error {
    var existing Payment
    err := db.Where("user_id = ? AND order_id = ?", userID, orderID).First(&existing).Error
    if err == nil {
        if existing.Status == "success" {
            return errors.New("payment already completed")
        }
    }
    return nil
}
```

**处理流程**:
```
1. 检测到重复支付
2. 若第一笔已成功:
   - 自动退款第二笔
   - 通知用户"支付成功，重复支付已退款"
3. 若两笔都失败:
   - 提示用户重新支付
```

---

### 2.5 退款边缘场景

#### 2.5.1 退款超时未到账

**场景**: 退款成功但用户未收到

**处理**:
```
1. 查询支付平台退款状态
2. 若平台已退款: 提示用户查看银行/支付宝
3. 若平台未退款: 重新发起退款
4. 24小时未到账: 转人工处理
```

#### 2.5.2 部分退款

**场景**: 订单已完成但服务有瑕疵，用户要求部分退款

**处理**:
```
1. 争议处理通过
2. 计算退款比例: 30%
3. 执行部分退款
4. 更新订单状态: completed → partially_refunded
5. 调整陪玩师收益: 扣除对应比例
```

---

## 3. 聊天系统边缘场景

### 3.1 消息发送失败

**场景描述**: 用户发送消息失败（网络/服务器原因）

**重试机制**:
```typescript
// 消息队列重试
const messageQueue = new MessageQueue({
  maxRetries: 3,
  backoff: 'exponential',
  onFailed: (msg) => {
    // 3次重试失败后
    notifyUser('消息发送失败，请检查网络');
    saveToLocalStorage(msg); // 本地缓存
  }
});
```

---

### 3.2 敏感词过滤误杀

**场景描述**: 正常内容被误判为敏感词

**处理流程**:

```
1. 消息被拦截
2. 提示: "消息包含敏感内容，是否申诉?"
3. 用户申诉:
   - 消息进入人工审核队列
   - 审核通过: 放行消息
   - 审核失败: 保持拦截
```

---

### 3.3 订单群聊冲突

**场景描述**: 用户拉黑了陪玩师，但存在进行中的订单

**处理策略**:
```
1. 拉黑不生效于进行中的订单
2. 订单群聊保持可用
3. 订单完成后:
   - 自动解散群聊
   - 拉黑关系生效
   - 双方无法再发消息
```

---

### 3.4 WebSocket连接断开

**场景描述**: 用户网络波动导致WebSocket断开

**自动重连**:
```typescript
// WebSocket重连策略
const ws = new WebSocket({
  url: WS_URL,
  reconnectInterval: 1000, // 初始1秒
  maxReconnectInterval: 30000, // 最大30秒
  reconnectDecay: 1.5, // 指数退避
  onReconnect: () => {
    // 重连成功后拉取离线消息
    fetchOfflineMessages();
  }
});
```

---

## 4. 争议处理边缘场景

### 4.1 争议超时未处理

**场景描述**: 用户/陪玩师发起争议后，客服30分钟未响应

**SLA机制**:
```
1. 争议创建时间记录
2. 倒计时30分钟
3. 超时自动升级:
   - 分配给上级客服
   - 发送告警给客服主管
   - 记录SLA违约
4. 若60分钟仍无响应:
   - 自动按用户诉求处理
   - 或全额退款(保护用户优先)
```

---

### 4.2 双方同时发起争议

**场景描述**: 用户和陪玩师几乎同时发起争议

**处理流程**:

```
1. 系统检测到重复争议
2. 合并为一个争议工单
3. 通知双方:
   "争议已合并，客服将统一处理"
4. 聊天快照包含双方提交的证据
```

---

### 4.3 争议证据不足

**场景描述**: 争议双方无法提供充分证据

**处理原则**:
```
1. 默认保护用户(消费者):
   - 若无法判定责任，倾向用户
   - 但有恶意投诉记录的用户除外
2. 陪玩师信誉保护:
   - 高信誉陪玩师(好评率>95%)给予更多信任
3. 争议处理记录:
   - 影响双方未来信誉分
   - 恶意争议方降低信誉
```

---

### 4.4 争议期间订单状态

**场景描述**: 争议进行中，订单资金如何处理

**资金冻结**:
```
1. 争议创建后:
   - 订单状态: completed → disputed
   - 陪玩师收益: 冻结
   - 无法提现
2. 争议解决:
   - 用户胜: 全额退款
   - 陪玩师胜: 解冻收益
   - 部分支持: 按比例退款
```

---

## 5. 提现系统边缘场景

### 5.1 提现到失败银行卡

**场景描述**: 陪玩师绑定错误的银行卡信息

**处理流程**:

```
1. 打款失败
2. 银行返回错误码(如: 卡号不存在/已销户)
3. 系统:
   - 提现状态: pending → failed
   - 金额退回到陪玩师余额
   - 通知: "提现失败，请检查收款信息"
4. 陪玩师:
   - 修改收款信息
   - 重新申请提现
```

---

### 5.2 提现金额异常

**子场景**:

| 场景 | 触发条件 | 处理 |
|------|---------|------|
| **超额提现** | 提现金额 > 可用余额 | 前端校验 + 后端二次校验 |
| **最小金额** | 提现金额 < 最低提现额(¥100) | 提示用户增加金额 |
| **日限额** | 当日累计 > 日限额(¥10000) | 提示明日再提或申请提高额度 |
| **未结收益** | 提现冻结期内收益 | 前端显示"可用余额"，排除冻结金额 |

---

### 5.3 提现审核争议

**场景描述**: 陪玩师对提现审核结果有异议

**申诉流程**:
```
1. 陪玩师点击"申诉"
2. 填写申诉理由
3. 转交财务主管复核
4. 复核结果:
   - 维持原判: 说明原因
   - 撤销驳回: 通过提现
```

---

## 6. 系统故障场景

### 6.1 数据库连接失败

**降级策略**:
```go
// 数据库健康检查
func CheckDBHealth() error {
    if err := db.DB().Ping(); err != nil {
        // 启用只读缓存模式
        enableReadOnlyMode()
        // 发送告警
        alert.Send("数据库连接失败")
        return err
    }
    return nil
}
```

**用户端提示**:
```
"系统繁忙，请稍后再试"
```

---

### 6.2 Redis缓存故障

**降级处理**:
```
1. 检测到Redis不可用
2. 自动切换到内存缓存
3. 性能下降但服务可用
4. 告警管理员
```

**代码**:
```go
// 缓存降级
func GetWithCache(key string) (string, error) {
    // 尝试Redis
    val, err := redis.Get(key)
    if err != nil {
        // 降级到内存缓存
        return memoryCache.Get(key)
    }
    return val, nil
}
```

---

### 6.3 第三方服务超时

**场景**: 支付/短信/推送服务超时

**熔断机制**:
```go
// 熔断器模式
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.isOpen() {
        return errors.New("circuit breaker is open")
    }

    err := fn()
    if err != nil {
        cb.recordFailure()
        if cb.failures >= cb.maxFailures {
            cb.open() // 开启熔断
        }
    }
    return err
}
```

---

### 6.4 消息队列积压

**场景**: WebSocket消息发送队列积压

**处理**:
```
1. 监控队列长度
2. 队列 > 1000条:
   - 限流: 只发送最近消息
   - 批量发送
   - 丢弃低优先级消息
3. 队列 > 10000条:
   - 告警
   - 增加消费者
```

---

## 7. 安全边缘场景

### 7.1 暴力登录攻击

**场景**: 攻击者尝试大量密码组合

**防护机制**:
```go
// 登录频率限制
func CheckLoginRateLimit(ip string) error {
    key := fmt.Sprintf("login:rate:%s", ip)
    count, _ := redis.Incr(key)
    redis.Expire(key, 5*time.Minute)

    if count > 5 {
        return errors.New("too many login attempts")
    }
    return nil
}

// 账户锁定
func CheckAccountLock(userID uint64) error {
    key := fmt.Sprintf("login:failed:%d", userID)
    count, _ := redis.Get(key)

    if count >= 5 {
        return errors.New("account locked for 30 minutes")
    }
    return nil
}
```

---

### 7.2 SQL注入防护

**场景**: 搜索框输入恶意SQL代码

**防护**:
```go
// 使用参数化查询
func SearchPlayers(keyword string) ([]Player, error) {
    var players []Player
    // ✅ 安全: 参数化查询
    err := db.Where("nickname LIKE ?", "%"+keyword+"%").Find(&players).Error

    // ❌ 不安全: 字符串拼接
    // err := db.Where("nickname LIKE '%" + keyword + "%'").Find(&players).Error

    return players, err
}
```

---

### 7.3 XSS攻击防护

**场景**: 用户输入包含恶意脚本代码

**防护**:
```go
// 输入过滤
func SanitizeInput(input string) string {
    // 移除HTML标签
    reg := regexp.MustCompile(`<[^>]*>`)
    cleaned := reg.ReplaceAllString(input, "")

    // 转义特殊字符
    cleaned = html.EscapeString(cleaned)

    return cleaned
}
```

---

### 7.4 CSRF攻击防护

**场景**: 攻击者诱导用户访问恶意页面

**防护**:
```go
// CSRF Token验证
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-CSRF-Token")
        expectedToken := GetCSRFToken(c)

        if token != expectedToken {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

### 7.5 恶意订单检测

**场景**: 用户频繁下单取消或发起争议

**检测逻辑**:
```go
// 异常行为检测
func DetectAbnormalBehavior(userID uint64) error {
    // 检查最近7天订单
    orders := GetRecentOrders(userID, 7)

    cancelRate := calculateCancelRate(orders)
    disputeRate := calculateDisputeRate(orders)

    if cancelRate > 0.5 || disputeRate > 0.3 {
        // 标记为可疑用户
        MarkSuspiciousUser(userID)
        // 限制下单
        return errors.New("account restricted due to abnormal behavior")
    }

    return nil
}
```

---

## 8. 附录

### 8.1 错误码规范

| 错误码 | 含义 | HTTP状态码 |
|--------|------|-----------|
| 4001 | 参数错误 | 400 |
| 4002 | 重复操作 | 400 |
| 4003 | 资源不存在 | 404 |
| 4004 | 权限不足 | 403 |
| 5001 | 系统繁忙 | 500 |
| 5002 | 第三方服务异常 | 502 |
| 5003 | 数据库错误 | 500 |

---

### 8.2 监控告警阈值

| 指标 | 阈值 | 级别 |
|------|------|------|
| API响应时间 | P95 > 500ms | 🟡 Warning |
| API错误率 | > 5% | 🟡 Warning |
| 数据库连接池 | 使用率 > 80% | 🟡 Warning |
| Redis连接 | 失败 > 10次/分钟 | 🔴 Critical |
| 订单超时率 | > 10% | 🔴 Critical |
| 争议SLA违约 | > 5单/天 | 🔴 Critical |

---

### 8.3 相关文档

| 文档 | 路径 |
|------|------|
| 用户旅程地图 | `docs/USER_JOURNEY.md` |
| API映射表 | `docs/API_PRD_MAPPING.md` |
| 数据模型 | `.kiro/steering/04-data-models.md` |

---

*文档维护：GameLink 开发团队*
*最后更新：2025-01-05*
