# 雪季猫业务整合计划

> 将雪季猫游戏陪玩平台业务逻辑整合到 GameLink 项目

## 一、业务对比分析

### 1.1 已有功能（可复用）

| 雪季猫功能 | GameLink 现有 | 状态 | 备注 |
|------------|---------------|------|------|
| 用户系统 | User 模型 | ✅ 已有 | 需扩展 VIP 字段 |
| 陪玩师管理 | Player 模型 | ✅ 已有 | 基本一致 |
| 订单系统 | Order 模型 | ✅ 已有 | 需扩展状态和字段 |
| 支付系统 | Payment 模型 | ✅ 已有 | 需扩展组合支付 |
| 钱包余额 | Wallet 模型 | ✅ 已有 | 基本一致 |
| 游戏分类 | Game 模型 | ✅ 已有 | 需扩展排序字段 |
| 服务项目 | ServiceItem 模型 | ✅ 已有 | 需大量扩展 |

### 1.2 需新增功能

| 雪季猫功能 | 优先级 | 状态 | 说明 |
|------------|--------|------|------|
| VIP 会员系统 | P0 | ⬜ 待开发 | 等级、折扣、月度券 |
| 优惠券系统 | P0 | ⬜ 待开发 | 满减/折扣券 |
| 充值系统 | P0 | ⬜ 待开发 | 档位、赠金、赠券 |
| 单陪/双陪定价 | P1 | ⬜ 待开发 | 服务项目扩展 |
| 游戏档案 | P1 | ⬜ 待开发 | 用户快捷下单 |
| 老板预存绑定 | P2 | ⬜ 待开发 | 上游余额同步 |
| 使用限制 | P2 | ⬜ 待开发 | 每日/每周/总量限购 |

---

## 二、数据模型扩展

### 2.1 User 模型扩展

```go
// 新增字段
VipLevel        string     // VIP等级: normal/vip1/vip2/vip3/svip1/svip2/svip3-6
VipUnlocked     bool       // VIP是否已解锁
MaxExp          int64      // 累计消费（分）
TotalRecharge   int64      // 累计充值（分）
BossDepositID   *uint64    // 绑定的老板预存ID（上游）
```

### 2.2 新增 VipLevel 模型

```go
type VipLevel struct {
    Base
    Title                string  // 等级名称
    PriceCents           int64   // 升级所需累计消费（分）
    OrderDiscount        float64 // 下单永久折扣 (0.98 = 98折)
    MonthlyCouponDiscount float64 // 每月折扣券折扣率
    MonthlyCouponCount   int     // 每月折扣券数量 (-1=无限)
    Benefits             string  // 其他权益描述(JSON)
    SortOrder            int     // 排序
}
```

### 2.3 新增 Coupon 模型

```go
type Coupon struct {
    Base
    UserID      uint64          // 所属用户
    Name        string          // 优惠券名称
    Type        CouponType      // 类型: normal/daily/weekly
    Mode        CouponMode      // 模式: deduct/discount
    Amount      int64           // 满减金额（分）或折扣率*100
    MinAmount   int64           // 最低消费门槛（分）
    ProjectIDs  string          // 限定项目ID(JSON数组)，空=通用
    Source      CouponSource    // 来源: system/vip_monthly/recharge/activity/manual
    State       CouponState     // 状态: 0未使用/1已使用/2已过期/3已锁定
    LockedBy    *uint64         // 锁定订单ID
    StartDate   time.Time       // 生效时间
    EndDate     time.Time       // 过期时间
}
```

### 2.4 新增 RechargeOption 模型

```go
type RechargeOption struct {
    Base
    AmountCents                int64   // 充值金额（分）
    BonusCents                 int64   // 赠送金额（分）
    BonusDiscountCouponCount   int     // 赠送折扣券数量
    BonusDiscountCouponRate    float64 // 赠送折扣券折扣率
    BonusDeductionCouponCount  int     // 赠送满减券数量
    BonusDeductionCouponAmount int64   // 赠送满减券金额（分）
    BonusDeductionCouponMin    int64   // 赠送满减券门槛（分）
    BonusCouponValidityDays    int     // 赠送券有效期(天)
    Description                string  // 档位描述
    Status                     bool    // 是否启用
    SortOrder                  int     // 排序
}
```

### 2.5 新增 RechargeRecord 模型

```go
type RechargeRecord struct {
    Base
    UserID          uint64          // 用户ID
    OptionID        uint64          // 充值档位ID
    AmountCents     int64           // 充值金额（分）
    BonusCents      int64           // 赠送金额（分）
    Status          RechargeStatus  // 状态: pending/paid/refunded
    PaymentMethod   PaymentMethod   // 支付方式
    ProviderTradeNo string          // 第三方交易号
    PaidAt          *time.Time      // 支付时间
    RefundedAt      *time.Time      // 退款时间
}
```

### 2.6 新增 GameProfile 模型

```go
type GameProfile struct {
    Base
    UserID         uint64  // 用户ID
    GameID         uint64  // 游戏分类ID
    GameNickname   string  // 游戏昵称
    GameAccountID  string  // 游戏账号ID
    PlatformType   string  // 平台类型: pc/mobile
    GameServer     string  // 游戏区服
    PhoneNumber    string  // 联系电话
    IsDefault      bool    // 是否默认
}
```

### 2.7 ServiceItem 模型扩展

```go
// 新增字段
ShowHandlerType    bool    // 是否显示单陪/双陪选项
SingleExtraPrice   int64   // 单陪加价金额（分）
UsageLimitType     string  // 限制类型: none/daily/weekly/total
UsageLimitCount    int     // 限制次数
Rating             float32 // 项目评分(1-5)
ServiceCount       int64   // 累计服务次数
SortOrder          int     // 排序权重
```

### 2.8 Order 模型扩展

```go
// 新增字段
HandlerType        string  // 陪玩类型: single/double
CouponID           *uint64 // 使用的优惠券ID
CouponDiscount     int64   // 优惠券折扣金额（分）
VipDiscount        int64   // VIP折扣金额（分）
FinalPriceCents    int64   // 最终支付金额（分）
BalancePaidCents   int64   // 余额支付金额（分）
WechatPaidCents    int64   // 微信支付金额（分）
GameProfileID      *uint64 // 游戏档案ID
Remark             string  // 备注
DispatchedAt       *time.Time // 派单时间
```

### 2.9 订单状态扩展

```go
// 雪季猫状态映射
const (
    OrderStatusPending    = "pending"     // 0 未支付
    OrderStatusPaid       = "paid"        // 1 已支付（等待派单）
    OrderStatusDispatched = "dispatched"  // 5 已派单（服务中）
    OrderStatusCompleted  = "completed"   // 2 已完成
    OrderStatusRefunded   = "refunded"    // 3 已退款
    OrderStatusCanceled   = "canceled"    // 4 已取消
)
```

---

## 三、整合进度跟踪

### 3.1 Phase 1: 核心模型扩展（P0）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| User 模型添加 VIP 字段 | ⬜ | - |
| 创建 VipLevel 模型 | ⬜ | - |
| 创建 Coupon 模型 | ⬜ | - |
| 创建 RechargeOption 模型 | ⬜ | - |
| 创建 RechargeRecord 模型 | ⬜ | - |
| ServiceItem 模型扩展 | ⬜ | - |
| Order 模型扩展 | ⬜ | - |
| 数据库迁移脚本 | ⬜ | - |

### 3.2 Phase 2: VIP 会员系统（P0）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| VIP 等级配置 CRUD | ⬜ | - |
| VIP 解锁逻辑 | ⬜ | - |
| VIP 等级计算 | ⬜ | - |
| VIP 永久折扣 | ⬜ | - |
| VIP 月度券发放 | ⬜ | - |
| VIP 升级预测通知 | ⬜ | - |

### 3.3 Phase 3: 优惠券系统（P0）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| 优惠券 CRUD | ⬜ | - |
| 优惠券验证逻辑 | ⬜ | - |
| 优惠券锁定/释放 | ⬜ | - |
| 优惠券过期检查（定时任务） | ⬜ | - |
| 每日/每周券发放（定时任务） | ⬜ | - |
| 优惠券折扣计算 | ⬜ | - |

### 3.4 Phase 4: 充值系统（P0）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| 充值档位配置 CRUD | ⬜ | - |
| 充值订单创建 | ⬜ | - |
| 充值支付回调 | ⬜ | - |
| 赠金发放 | ⬜ | - |
| 赠券发放 | ⬜ | - |
| 充值退款 | ⬜ | - |

### 3.5 Phase 5: 订单流程优化（P1）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| 单陪/双陪定价 | ⬜ | - |
| 价格计算（VIP折扣 vs 优惠券） | ⬜ | - |
| 组合支付（余额+微信） | ⬜ | - |
| 订单状态流转优化 | ⬜ | - |
| 订单超时取消（定时任务） | ⬜ | - |
| 退款流程（余额/微信/组合） | ⬜ | - |
| 累计消费(maxExp)流转 | ⬜ | - |

### 3.6 Phase 6: 辅助功能（P1-P2）

| 任务 | 状态 | 完成日期 |
|------|------|----------|
| 游戏档案 CRUD | ⬜ | - |
| 服务项目使用限制 | ⬜ | - |
| 老板预存绑定 | ⬜ | - |
| 统一余额服务 | ⬜ | - |

---

## 四、API 接口规划

### 4.1 用户端 API

```
POST   /api/v1/user/orders              # 创建订单
POST   /api/v1/user/orders/:id/pay      # 发起支付
POST   /api/v1/user/orders/:id/cancel   # 取消订单
GET    /api/v1/user/orders              # 订单列表
GET    /api/v1/user/orders/:id          # 订单详情

GET    /api/v1/user/coupons             # 优惠券列表
GET    /api/v1/user/coupons/available   # 可用优惠券（下单时）

GET    /api/v1/user/vip                 # VIP信息
GET    /api/v1/user/vip/levels          # VIP等级列表

POST   /api/v1/user/recharge            # 创建充值订单
GET    /api/v1/user/recharge/options    # 充值档位列表

GET    /api/v1/user/game-profiles       # 游戏档案列表
POST   /api/v1/user/game-profiles       # 创建游戏档案
PUT    /api/v1/user/game-profiles/:id   # 更新游戏档案
DELETE /api/v1/user/game-profiles/:id   # 删除游戏档案
```

### 4.2 管理端 API

```
# VIP 管理
GET    /api/v1/admin/vip/levels         # VIP等级列表
POST   /api/v1/admin/vip/levels         # 创建VIP等级
PUT    /api/v1/admin/vip/levels/:id     # 更新VIP等级

# 优惠券管理
GET    /api/v1/admin/coupons            # 优惠券列表
POST   /api/v1/admin/coupons/batch      # 批量发放优惠券

# 充值管理
GET    /api/v1/admin/recharge/options   # 充值档位列表
POST   /api/v1/admin/recharge/options   # 创建充值档位
PUT    /api/v1/admin/recharge/options/:id # 更新充值档位
POST   /api/v1/admin/recharge/:id/refund # 充值退款

# 订单管理
POST   /api/v1/admin/orders/:id/dispatch # 派单
POST   /api/v1/admin/orders/:id/complete # 完成订单
POST   /api/v1/admin/orders/:id/refund   # 订单退款
```

---

## 五、配置项

### 5.1 系统配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `order_cancel_timeout` | 未支付订单取消时限(分钟) | 30 |
| `refund_timeout_hours` | 退款申请时限(小时) | 168 (7天) |
| `vip_unlock_threshold` | VIP充值解锁门槛(分) | 配置 |
| `vip_unlock_consume_threshold` | VIP消费解锁门槛(分) | 配置 |

---

## 六、注意事项

1. **金额单位**：所有金额使用分(Cents)存储，与现有 GameLink 保持一致
2. **并发控制**：退款、支付回调需要防重复处理
3. **优惠互斥**：VIP折扣和优惠券二选一，取优惠更大的
4. **状态流转**：严格按照状态机流转，防止非法状态变更
5. **定时任务**：优惠券过期、订单超时、月度券发放等需要定时任务支持

---

## 变更日志

| 日期 | 变更内容 |
|------|----------|
| 2024-12-24 | 创建整合计划文档 |
