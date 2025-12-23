# 数据模型 - 枚举类型和数据库索引

> 所有枚举类型定义和数据库索引规范

## 相关文档

- [04-data-models.md](./04-data-models.md) - 核心模块
- [04a-marketing-models.md](./04a-marketing-models.md) - 营销模块
- [04b-team-models.md](./04b-team-models.md) - 团队系统
- [04d-notification-models.md](./04d-notification-models.md) - 通知系统

---

## 枚举类型

### 用户相关

#### 用户状态 (UserStatus)
`active` | `suspended` | `banned`

#### 认证状态 (VerificationStatus)
`pending` | `verified` | `rejected`

---

### 订单相关

#### 订单状态 (OrderStatus)
`pending` | `confirmed` | `in_progress` | `completed` | `canceled` | `refunded`

#### 订单明细状态 (OrderItemStatus)
`pending` | `matched` | `completed` | `canceled`

#### 订单陪玩师状态 (OrderPlayerStatus)
`joined` | `left` | `completed`

#### 争议状态 (DisputeStatus)
`pending` | `assigned` | `mediating` | `resolved` | `rejected` | `canceled`

#### 货币类型 (Currency)
`CNY` | `USD` | `EUR`

---

### 支付相关

#### 支付状态 (PaymentStatus)
`pending` | `paid` | `failed` | `refunded`

#### 支付方式 (PaymentMethod)
`wechat` | `alipay` | `wallet` | `combined`

#### 充值状态 (RechargeStatus)
`pending` | `paid` | `failed` | `refunded` | `canceled`

---

### 优惠券相关

#### 优惠券类型 (CouponType)
`deduct` | `discount`

#### 优惠券适用范围 (CouponScope)
`all` | `game` | `item`

#### 优惠券来源 (CouponSource)
`new_user` | `link` | `vip` | `recharge` | `activity` | `manual`

#### 优惠券状态 (CouponState)
`available` | `locked` | `used` | `expired` | `deleted`

---

### 服务项目相关

#### 使用限制类型 (UsageLimitType)
`none` | `once` | `daily` | `weekly` | `monthly`

---

### 团队相关

#### 团队状态 (TeamStatus)
`active` | `busy` | `inactive`

#### 团队成员角色 (TeamMemberRole)
`leader` | `member`

#### 团队成员状态 (TeamMemberStatus)
`active` | `left` | `kicked`

---

### 活动相关

#### 活动状态 (ActivityStatus)
`draft` | `preheat` | `active` | `paused` | `ended` | `canceled`

#### 活动类型 (ActivityType)
`coupon` | `discount` | `gift`

---

### 推荐相关

#### 推荐类型 (ReferralType)
`user_to_user` | `player_to_player` | `user_to_player`

#### 推荐状态 (ReferralStatus)
`pending` | `completed` | `rewarded` | `expired` | `canceled`

#### 奖励类型 (RewardType)
`cash` | `coupon` | `points`

#### 奖励发放状态 (ReferralRewardStatus)
`pending` | `issued` | `failed`

---

### 通知相关

#### 通知类型 (NotificationType)
`order_status` | `vip_expire` | `coupon_expire` | `activity_start` | `activity_end` | `system` | `promotion`（预留） | `chat`（预留）

#### 通知渠道 (NotificationChannel)
`in_app` | `push` | `sms`（预留） | `wechat`（预留） | `email`（预留）

#### 通知状态 (NotificationStatus)
`pending` | `sent` | `read` | `failed` | `canceled`

---

## 数据库索引

### 唯一索引

| 表 | 字段 |
|---|---|
| User | Phone |
| User | Email |
| Order | OrderNo |
| Player | UserID |
| Game | Key |
| RoleModel | Slug |
| Permission | Method+Path |
| Permission | Code |
| ServiceItem | ItemCode |
| VipLevel | Slug |
| VipConfig | ConfigKey |
| CouponTemplate | ClaimLink |
| RechargeRecord | OrderNo |
| ReferralConfig | ConfigKey |
| ReferralCode | Code |
| ActivityDailyStats | (ActivityID, StatsDate) |
| NotificationTemplate | Code |
| NotificationConfig | ConfigKey |
| UserNotificationSetting | UserID |

### 复合索引

#### 订单相关
- Order: (UserID, Status, CreatedAt)
- Order: (PlayerID, Status)
- OrderItem: (OrderID, Status)
- OrderPlayer: (OrderID, PlayerID)
- OrderPlayer: (OrderItemID, PlayerID)

#### 用户相关
- User: (Status, LastLoginAt)
- User: VipLevelID

#### 优惠券相关
- Coupon: (UserID, State)
- Coupon: (TemplateID)
- Coupon: (ExpireAt)
- Coupon: (LockedByOrderID)
- Coupon: (UsedOrderID)

#### 充值相关
- RechargeRecord: (UserID)
- RechargeRecord: (Status)
- RechargeRecord: (MerchantOrderNo)
- RechargeRecord: (ProviderTradeNo)
- RechargeRecord: (MerchantID)
- RechargeRecord: (PaymentChannel)

#### 活动相关
- Activity: (Status, StartAt, EndAt)
- Activity: (Type, Status)
- ActivityReward: (ActivityID)
- ActivityParticipation: (ActivityID, UserID)
- ActivityParticipation: (UserID, ClaimedAt)

#### 推荐相关
- ReferralCode: (UserID, Type)
- Referral: (ReferrerID, Status)
- Referral: (RefereeID)
- ReferralReward: (UserID, Status)

#### 通知相关
- UserNotification: (UserID, Status)
- UserNotification: (UserID, Type)
- UserNotification: (TemplateID)
- UserNotification: (Channel)
- UserNotification: (RelatedType, RelatedID)
- NotificationSchedule: (ScheduleAt, Status)
- NotificationSchedule: (Type)

---

## 变更日志

| 日期 | 变更内容 |
|------|----------|
| 2024-12-24 | 新增通知系统：NotificationTemplate、UserNotification、UserNotificationSetting、NotificationConfig、NotificationSchedule 模型 |
| 2024-12-24 | 新增活动系统：Activity、ActivityReward、ActivityParticipation、ActivityDailyStats 模型 |
| 2024-12-24 | 新增推荐/邀请系统（预留）：ReferralConfig、ReferralCode、Referral、ReferralReward 模型 |
| 2024-12-24 | Base 模型新增 ExtJSON 扩展字段 |
| 2024-12-24 | 新增团队系统：Team、TeamMember、TeamInvite 模型 |
| 2024-12-24 | 新增订单补差价业务流程文档（升级座位/加座位） |
| 2024-12-24 | 新增充值系统：RechargeOption、RechargeRecord 模型 |
| 2024-12-24 | VipLevel 月度券改为绑定优惠券模板ID |
| 2024-12-24 | 新增 VIP 会员系统：VipLevel、VipConfig 模型，User VIP 字段 |
| 2024-12-24 | 新增优惠券系统：CouponTemplate、Coupon 模型 |
| 2024-12-24 | ServiceItem 新增 VipPriceCents、UsageLimitType、UsageLimitCount、MaxPerOrder 字段 |
| 2024-12-24 | 新增多人订单支持：OrderItem、OrderPlayer 模型 |
| 2024-12-24 | ServiceItem/Order 新增 RequiredPlayers 字段 |
| 2024-12-24 | ChatGroup 新增语音服务预留字段 |
| 2024-12-24 | Review 新增 OrderItemID 字段支持多人评价 |
| 2024-12-21 | 初始化数据模型文档 |
