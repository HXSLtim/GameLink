# 数据模型 - 营销模块

> VIP会员、优惠券、充值、活动、推荐系统

## 相关文档

- [04-data-models.md](./04-data-models.md) - 核心模块
- [04b-team-models.md](./04b-team-models.md) - 团队系统
- [04c-enums-indexes.md](./04c-enums-indexes.md) - 枚举类型和数据库索引
- [04d-notification-models.md](./04d-notification-models.md) - 通知系统

---

## VIP 会员模块

### VipLevel（VIP等级配置）

> 完全可配置的VIP等级，支持从数据库动态管理

| 字段 | 类型 | 说明 |
|------|------|------|
| Slug | string | 等级标识（唯一，如 vip1/vip2/svip1） |
| Title | string | 等级名称（显示用） |
| ExpRequired | int64 | 升级所需累计消费/经验（分） |
| OrderDiscount | float64 | 下单永久折扣（0.98=98折，1.0=无折扣） |
| MonthlyCouponTemplateID | *uint64 | 月度券模板ID |
| MonthlyCouponCount | int | 每月发放数量 |
| IconURL | string | 等级图标 |
| Color | string | 等级颜色（前端展示） |
| Benefits | string | 其他权益描述（JSON） |
| SortOrder | int | 排序（越小越靠前） |
| IsDefault | bool | 是否默认等级（新用户解锁后） |
| IsActive | bool | 是否启用 |

### VipConfig（VIP系统配置）

> 全局VIP配置，键值对形式

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `unlock_by_consume` - 累计消费解锁门槛（分）
- `unlock_by_recharge` - 累计充值解锁门槛（分）
- `expire_days` - VIP过期天数（0=永久）

### VIP 业务流程

```
【VIP 解锁】
用户消费/充值 → 检查门槛
  ├── 累计消费 >= unlock_by_consume → 解锁VIP
  └── 累计充值 >= unlock_by_recharge → 解锁VIP

【VIP 升级】
订单完成 → VipExp += 订单金额
  └── 检查是否达到下一等级 ExpRequired

【月度券发放】
每月1号（或用户首次登录时检查）
  └── LastMonthlyCouponAt 不在本月 → 发放券
```

---

## 折扣叠加规则

> 统一的折扣计算规则，适用于 VIP 折扣、普通优惠券、活动券

```
【全局规则】
- VIP 永久折扣与优惠券互斥，不叠加
- 多张优惠券之间互斥，不叠加
- 系统自动选择最优折扣方案

【折扣优先级计算】
1. 计算 VIP 永久折扣金额
2. 计算每张可用优惠券的折扣金额
3. 取最大折扣金额的方案

【活动券特殊规则】
- Activity.AllowVipStack = true：活动券可与 VIP 折扣叠加（特例）
- Activity.AllowVipStack = false：活动券与 VIP 折扣择优使用（默认）

【折扣计算示例】
订单金额：¥100
VIP 折扣：98折 → 优惠 ¥2
优惠券A：满100减10 → 优惠 ¥10
优惠券B：9折 → 优惠 ¥10
→ 系统选择优惠券A或B（优惠 ¥10 > VIP ¥2）

【活动券叠加示例】
订单金额：¥100
VIP 折扣：98折 → 优惠 ¥2
活动券（AllowVipStack=true）：满100减5 → 优惠 ¥5
→ 可叠加：¥100 × 0.98 - ¥5 = ¥93（总优惠 ¥7）
```

---

## 优惠券模块

### CouponTemplate（优惠券模板/配置）

> 优惠券批量发放的配置模板

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 券名称 |
| Type | CouponType | 类型：deduct（满减）/discount（折扣） |
| Source | CouponSource | 来源（见枚举） |
| Description | string | 描述 |
| MinAmountCents | int64 | 最低消费门槛（分），0=无门槛 |
| DeductAmountCents | int64 | 满减金额（分）- 满减券用 |
| DiscountRate | float64 | 折扣率（0.9=9折）- 折扣券用 |
| MaxDiscountCents | int64 | 最大折扣金额（分）- 折扣券用，0=无上限 |
| Scope | CouponScope | 适用范围：all/game/item |
| GameIDs | string | 指定游戏ID列表（JSON数组） |
| ItemIDs | string | 指定服务项目ID列表（JSON数组） |
| ValidityType | string | days=固定天数，fixed=固定截止日期 |
| ValidityDays | int | 有效天数（领取后） |
| FixedExpireAt | *time.Time | 固定截止日期 |
| TotalCount | int | 总发放数量（0=无限制） |
| ClaimedCount | int | 已领取数量 |
| PerUserLimit | int | 每人限领数量 |
| ClaimLink | string | 领取链接码（链接领取用） |
| IsActive | bool | 是否启用 |

### Coupon（用户优惠券）

> 用户持有的具体优惠券实例

| 字段 | 类型 | 说明 |
|------|------|------|
| TemplateID | uint64 | 模板ID |
| UserID | uint64 | 用户ID |
| State | CouponState | 状态（见枚举） |
| Name | string | 券名称（冗余） |
| Type | CouponType | 类型（冗余） |
| Source | CouponSource | 来源（冗余） |
| MinAmountCents | int64 | 最低消费门槛（冗余） |
| DeductAmountCents | int64 | 满减金额（冗余） |
| DiscountRate | float64 | 折扣率（冗余） |
| MaxDiscountCents | int64 | 最大折扣金额（冗余） |
| Scope | CouponScope | 适用范围（冗余） |
| GameIDs | string | 指定游戏ID（冗余） |
| ItemIDs | string | 指定服务项目ID（冗余） |
| ClaimedAt | *time.Time | 领取时间 |
| ExpireAt | time.Time | 过期时间 |
| UsedAt | *time.Time | 使用时间 |
| LockedByOrderID | *uint64 | 锁定的订单ID |
| LockedAt | *time.Time | 锁定时间 |
| UsedOrderID | *uint64 | 使用的订单ID |
| DiscountCents | int64 | 实际折扣金额（分） |

### 优惠券业务流程

```
【优惠券类型】
满减券（deduct）：订单 >= MinAmountCents，减 DeductAmountCents
折扣券（discount）：折扣 = 订单金额 × (1 - DiscountRate)，上限 MaxDiscountCents

【优惠券来源】
new_user | link | vip | recharge | activity | manual

【优惠券生命周期】
领取(available) → 下单选择(locked) → 支付成功(used) → 退款(available)
```

---

## 充值模块

### RechargeOption（充值档位配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 档位名称 |
| AmountCents | int64 | 充值金额（分） |
| BonusCents | int64 | 赠送金额（分） |
| TotalCents | int64 | 实际到账（分）= AmountCents + BonusCents |
| OriginalCents | *int64 | 原价（分），用于显示划线价 |
| DiscountPercent | *int | 折扣百分比（显示用） |
| Description | string | 描述 |
| Tag | string | 标签（如"推荐"、"热门"） |
| IconURL | string | 图标URL |
| SortOrder | int | 排序 |
| IsActive | bool | 是否启用 |
| IsRecommended | bool | 是否推荐 |
| CouponTemplateID | *uint64 | 赠送优惠券模板ID |
| CouponCount | int | 赠送优惠券数量 |
| MinVipLevel | *uint64 | 最低VIP等级要求 |
| PerUserLimit | int | 每人限购次数（0=无限制） |
| TotalLimit | int | 总限购次数（0=无限制） |
| PurchaseCount | int | 已购买次数 |

### RechargeRecord（充值记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID |
| OptionID | *uint64 | 档位ID（自定义金额时为nil） |
| AmountCents | int64 | 充值金额（分） |
| BonusCents | int64 | 赠送金额（分） |
| TotalCents | int64 | 实际到账（分） |
| Status | RechargeStatus | 状态（见枚举） |
| OrderNo | string | 内部订单号（唯一） |
| MerchantOrderNo | string | 商户订单号（提交给支付渠道） |
| ProviderTradeNo | string | 第三方交易号（支付渠道返回） |
| MerchantID | string | 商户号（收款分流） |
| CollectionEntity | string | 收款主体 |
| PaymentChannel | string | 支付渠道：wechat/alipay |
| PaymentMethod | string | 支付方式：wechat_h5/wechat_mini/alipay_h5 等 |
| PaidAt | *time.Time | 支付时间 |
| ExpireAt | *time.Time | 过期时间 |
| RefundedAt | *time.Time | 退款时间 |
| RefundAmountCents | int64 | 退款金额（分） |
| RefundReason | string | 退款原因 |
| RefundProviderNo | string | 退款第三方单号 |
| CouponIssued | bool | 优惠券是否已发放 |
| CouponIDs | string | 发放的优惠券ID列表（JSON数组） |
| ClientIP | string | 客户端IP |
| UserAgent | string | User-Agent |
| DeviceInfo | string | 设备信息（JSON） |
| Remark | string | 备注 |

### 充值业务流程

```
【充值下单】
用户选择档位 → 创建 RechargeRecord → 根据 RoutingRule 选择收款主体

【支付回调】
支付渠道回调 → 验证签名 → 更新记录 → 触发后续处理

【充值成功后处理】
1. 钱包充值：Wallet.BalanceCents += TotalCents
2. 累计充值更新：User.TotalRechargeCents += AmountCents
3. VIP 解锁检查
4. 优惠券发放（如果档位配置了）

【退款】
管理员发起退款 → 调用支付渠道退款接口 → 扣减钱包余额
```

---

## 活动模块

> 限时活动系统，支持优惠券发放，复用 CouponTemplate

### Activity（活动）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 活动名称 |
| Description | string | 活动描述 |
| Type | ActivityType | 活动类型 |
| Status | ActivityStatus | 活动状态 |
| CoverURL | string | 活动封面图 |
| BannerURL | string | 活动Banner图 |
| PreheatAt | *time.Time | 预热开始时间 |
| StartAt | time.Time | 活动开始时间 |
| EndAt | time.Time | 活动结束时间 |
| TotalLimit | int | 总参与次数限制（0=无限制） |
| DailyLimit | int | 每日参与次数限制（0=无限制） |
| PerUserLimit | int | 每人参与次数限制 |
| TotalParticipants | int | 总参与人数 |
| TodayParticipants | int | 今日参与人数 |
| TotalClaimed | int | 总领取次数 |
| AllowVipStack | bool | 是否允许与VIP折扣叠加 |
| Rules | string | 活动规则说明 |
| SortOrder | int | 排序 |
| IsVisible | bool | 是否在前端展示 |

### ActivityReward（活动奖励配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ActivityID | uint64 | 活动ID |
| CouponTemplateID | uint64 | 优惠券模板ID |
| CouponCount | int | 每次发放数量 |
| Probability | int | 发放概率（1-100，预留抽奖用） |
| TotalStock | int | 总库存（0=无限制） |
| RemainingStock | int | 剩余库存 |
| SortOrder | int | 排序 |

### ActivityParticipation（活动参与记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| ActivityID | uint64 | 活动ID |
| UserID | uint64 | 用户ID |
| RewardID | uint64 | 奖励配置ID |
| CouponIDs | string | 发放的优惠券ID列表（JSON数组） |
| ClaimedAt | time.Time | 领取时间 |
| ClientIP | string | 客户端IP |

### ActivityDailyStats（活动每日统计）

| 字段 | 类型 | 说明 |
|------|------|------|
| ActivityID | uint64 | 活动ID |
| StatsDate | time.Time | 统计日期 |
| Participants | int | 当日参与人数 |
| ClaimCount | int | 当日领取次数 |

### 活动业务流程

```
【活动生命周期】
draft → preheat → active → ended
         ↓
       paused（可恢复）
         ↓
      canceled（不可恢复）

【预热期】
PreheatAt <= now < StartAt → 活动可见，但不可参与

【用户领取流程】
用户点击领取 → 检查限制（活动状态、每人限制、每日限制、总量限制、库存）
  → 所有检查通过 → 发放优惠券 → 创建参与记录 → 更新统计

【VIP叠加配置】
AllowVipStack = true：活动券可与VIP折扣叠加
AllowVipStack = false：活动券与VIP折扣择优使用（互斥）

【多活动处理】
同时存在多个活动 → 领取的券视作常规优惠券 → 下单时择优使用（不叠加）
```

---

## 推荐/邀请模块（预留）

> 邀请系统预留模型，支持多级分销和多种奖励类型

### ReferralConfig（推荐配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `enabled` - 是否启用推荐系统
- `expire_days` - 邀请链接过期天数
- `max_level` - 最大推荐层级（1=直推，2=二级分销）
- `user_reward_type` - 用户邀请奖励类型
- `user_reward_amount` - 用户邀请奖励金额（分）
- `player_reward_type` - 陪玩师邀请奖励类型
- `player_reward_amount` - 陪玩师邀请奖励金额（分）

### ReferralCode（邀请码）

| 字段 | 类型 | 说明 |
|------|------|------|
| Code | string | 邀请码（唯一） |
| UserID | uint64 | 所属用户ID |
| Type | ReferralType | 推荐类型 |
| IsActive | bool | 是否启用 |
| ExpireAt | *time.Time | 过期时间（nil=永久） |
| UseCount | int | 使用次数 |
| MaxUse | int | 最大使用次数（0=无限制） |

### Referral（推荐记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| ReferrerID | uint64 | 推荐人ID |
| RefereeID | uint64 | 被推荐人ID |
| CodeID | *uint64 | 使用的邀请码ID |
| Type | ReferralType | 推荐类型 |
| Level | int | 推荐层级（1=直推，2=二级） |
| Status | ReferralStatus | 状态 |
| CompletedAt | *time.Time | 完成时间 |
| RewardType | RewardType | 奖励类型 |
| RewardAmountCents | int64 | 奖励金额（分） |
| RewardedAt | *time.Time | 奖励发放时间 |
| RewardNote | string | 奖励备注 |
| RefereeCondition | string | 完成条件：registered/first_order/first_recharge |

### ReferralReward（推荐奖励记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| ReferralID | uint64 | 推荐记录ID |
| UserID | uint64 | 获得奖励的用户ID |
| Type | RewardType | 奖励类型 |
| AmountCents | int64 | 奖励金额（分） |
| CouponID | *uint64 | 发放的优惠券ID |
| Status | ReferralRewardStatus | 状态：pending/issued/failed |
| IssuedAt | *time.Time | 发放时间 |
| FailureReason | string | 失败原因 |

### 推荐业务流程

```
【邀请码生成】
用户/陪玩师 → 生成专属邀请码

【被邀请人注册】
新用户通过邀请码注册 → 创建推荐记录 → ReferralCode.UseCount++

【完成条件检查】
registered - 注册即完成
first_order - 首单完成
first_recharge - 首次充值

【奖励发放】
推荐完成 → 发放奖励（cash/coupon/points）

【多级分销】（max_level > 1 时）
A 邀请 B，B 邀请 C → C 完成条件 → B 获得一级奖励，A 获得二级奖励
```
