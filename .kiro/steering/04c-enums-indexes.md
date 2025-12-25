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

#### 用户角色 (Role)
`user` | `player` | `admin`

#### 用户状态 (UserStatus)
`active` | `suspended` | `banned`

#### 登录方式 (LoginType)
`password` | `sms` | `email` | `oauth`（预留）

#### 认证状态 (VerificationStatus)
`pending` | `verified` | `rejected`

#### 陪玩师段位认证状态 (PlayerRankStatus)
`pending` | `verified` | `rejected` | `revoked` | `expired`

#### 实名认证状态 (CertificationStatus)
`pending` | `verified` | `rejected`

#### 拉黑用户类型 (BlockUserType)
`user` | `player`

#### 拉黑状态 (BlockStatus)
`active` | `canceled` | `admin_canceled`

#### 陪玩师在线状态 (PlayerOnlineStatus)
`online` | `offline` | `busy`

---

### 订单相关

#### 订单状态 (OrderStatus)
`pending` | `confirmed` | `in_progress` | `completed` | `canceled` | `refunded` | `disputed`

**状态说明：**
- `pending` - 待支付：订单创建后等待用户支付
- `confirmed` - 已确认：支付成功，等待陪玩师接单
- `in_progress` - 进行中：陪玩师已接单，服务进行中
- `completed` - 已完成：服务完成
- `canceled` - 已取消：订单被取消（超时/用户取消/系统取消）
- `refunded` - 已退款：订单已退款
- `disputed` - 争议中：订单存在争议，等待处理

**状态流转图：**
```
pending → canceled（支付超时/用户取消）
pending → confirmed（支付成功）
confirmed → canceled（接单超时/用户取消）
confirmed → in_progress（陪玩师接单完成）
in_progress → completed（服务完成）
in_progress → disputed（发起争议）
in_progress → canceled（双方协商取消）
completed → refunded（7天内退款审核通过）
completed → disputed（7天内发起争议）
disputed → completed（争议驳回，恢复原状态）
disputed → refunded（争议通过，全额退款）
```

#### 订单明细状态 (OrderItemStatus)
`pending` | `matched` | `completed` | `canceled`

#### 订单陪玩师状态 (OrderPlayerStatus)
`joined` | `left` | `completed`

#### 争议发起人类型 (DisputeInitiatorType)
`user` | `player`

#### 争议类型 (DisputeType)
`service_quality` | `bad_attitude` | `incomplete_service` | `user_not_cooperative` | `user_harassment` | `other`

#### 争议状态 (DisputeStatus)
`pending` | `assigned` | `mediating` | `closed`

**状态说明：**
- `pending` - 待处理：争议刚发起
- `assigned` - 已分配：客服已分配
- `mediating` - 调解中：客服正在处理
- `closed` - 已关闭：争议处理完成（查看 Resolution 获取结果）

#### 争议处理决定 (DisputeResolution)
`full_refund` | `dismissed`

**决定说明：**
- `full_refund` - 支持申诉方，全额退款
- `dismissed` - 驳回申诉，维持原状

#### 货币类型 (Currency)
`CNY` | `USD` | `EUR`

#### 订单超时类型 (OrderTimeoutType)
`accept_timeout`

#### 订单超时处理动作 (OrderTimeoutAction)
`canceled` | `refunded` | `notified`

#### 客服分配状态 (ServiceAssignmentStatus)
`assigned` | `joined` | `left` | `completed`

---

### 支付相关

#### 支付状态 (PaymentStatus)
`pending` | `paid` | `failed` | `refunded`

#### 支付方式 (PaymentMethod)
`wechat` | `alipay` | `wallet` | `combined`

#### 充值状态 (RechargeStatus)
`pending` | `paid` | `failed` | `refunded` | `canceled`

#### 结算状态 (SettlementStatus)
`pending` | `disputed` | `settled`

**状态说明：**
- `pending` - T+7 售后期内，收入冻结中
- `disputed` - 发生争议，收入继续冻结
- `settled` - T+7 结束且无争议，收入已解冻

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

#### 服务项目子类别 (ServiceItemSubCategory)
`solo` | `team` | `gift`

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

### 评价相关

#### 评价状态 (ReviewStatus)
`pending` | `approved` | `rejected` | `deleted`

#### 评价回复状态 (ReviewReplyStatus)
`pending` | `approved` | `rejected`

#### 申诉状态 (AppealStatus)
`pending` | `approved` | `rejected`

---

### 聊天相关

#### 聊天群组类型 (ChatGroupType)
`public` | `order`

#### 聊天消息类型 (ChatMessageType)
`text` | `image` | `file` | `system` | `voice`（预留） | `emoji`（预留）

#### 聊天消息审核状态 (ChatMessageAuditStatus)
`pending` | `approved` | `rejected`

#### 聊天成员角色 (ChatMemberRole)
`owner` | `admin` | `member`

---

### 通知相关

#### 通知类型 (NotificationType)
`order_status` | `vip_expire` | `coupon_expire` | `activity_start` | `activity_end` | `system` | `promotion`（预留） | `chat`（预留） | `review`（评价提醒） | `review_reply`（评价回复）

#### 通知渠道 (NotificationChannel)
`in_app` | `push` | `sms`（预留） | `wechat`（预留） | `email`（预留）

#### 通知状态 (NotificationStatus)
`pending` | `sent` | `read` | `failed` | `canceled`

---

### 敏感词相关

#### 敏感词分类 (SensitiveWordCategory)
`politics` | `porn` | `abuse` | `ad` | `other`

#### 敏感词匹配类型 (SensitiveWordMatchType)
`exact` | `fuzzy` | `regex`

---

### 内容管理相关

#### 内容类型 (ContentType)
`announcement` | `help` | `agreement` | `dynamic`

#### 内容状态 (ContentStatus)
`draft` | `pending` | `published` | `rejected` | `archived`

---

### 提现相关

#### 提现方式 (WithdrawMethod)
`alipay` | `wechat` | `bank`

#### 提现状态 (WithdrawStatus)
`pending` | `approved` | `rejected` | `processing` | `completed` | `failed`

---

### 排名相关

#### 排名类型 (RankingType)
`income` | `order_count` | `quality` | `popularity`

#### 排名周期 (RankingPeriod)
`daily` | `weekly` | `monthly` | `yearly`

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
| PlayerCertification | PlayerID |
| AuthConfig | ConfigKey |

### 复合索引

#### 订单相关
- Order: (UserID, Status, CreatedAt)
- Order: (PlayerID, Status)
- Order: (RecipientPlayerID)
- OrderItem: (OrderID, Status)
- OrderItem: (ReviewID)
- OrderPlayer: (OrderID, PlayerID)
- OrderPlayer: (OrderItemID, PlayerID)

#### 用户相关
- User: (Status, LastLoginAt)
- User: VipLevelID
- User: (BannedBy)
- UserLoginLog: (UserID, LoginAt)
- UserLoginLog: (DeviceID)
- UserLoginLog: (TokenID)
- UserLoginLog: (LoginType)

#### 优惠券相关
- Coupon: (UserID, State)
- Coupon: (TemplateID)
- Coupon: (ExpireAt)
- Coupon: (Source)
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

#### 陪玩师等级/认证相关
- GameRank: (GameID)
- PlayerRankRecord: (PlayerID)
- PlayerRankRecord: (GameID)
- PlayerRankRecord: (RankID)
- PlayerRankRecord: (Status)
- PlayerCertification: (Status)

#### 订单超时相关
- OrderTimeoutConfig: (ConfigKey) - 唯一索引
- OrderTimeoutLog: (OrderID)
- OrderTimeoutLog: (TimeoutType)
- OrderServiceAssignment: (OrderID)
- OrderServiceAssignment: (ServiceUserID)
- OrderServiceAssignment: (ChatGroupID)
- OrderServiceAssignment: (Status)

#### 用户拉黑相关
- UserBlock: (BlockerID)
- UserBlock: (BlockedID)
- UserBlock: (BlockerType)
- UserBlock: (BlockedType)
- UserBlock: (Status)

#### 评价相关
- Review: (OrderID)
- Review: (UserID)
- Review: (PlayerID)
- Review: (Status)
- Review: (OrderItemID)
- Review: (ExpireAt)
- ReviewReply: (ReviewID)
- ReviewReply: (PlayerID)
- ReviewAppeal: (ReviewID)
- ReviewAppeal: (PlayerID)
- ReviewAppeal: (Status)

#### 争议相关
- OrderDispute: (OrderID)
- OrderDispute: (InitiatorID)
- OrderDispute: (Status)
- OrderDispute: (OriginalServiceID)
- OrderDispute: (AssignedServiceID)
- DisputeTemplate: (Code) - 唯一索引
- DisputeTemplate: (InitiatorType)
- ChatSnapshot: (DisputeID)
- ChatSnapshot: (OrderID)

#### 聊天相关
- ChatGroup: (GroupType)
- ChatGroup: (RelatedOrderID)
- ChatGroup: (IsActive)
- ChatMessage: (GroupID, CreatedAt)
- ChatMessage: (SenderID)
- ChatMessage: (AuditStatus)
- ChatGroupMember: (GroupID, UserID) - 唯一索引
- ChatGroupMember: (UserID)
- ChatConfig: (ConfigKey) - 唯一索引

#### 敏感词相关
- SensitiveWord: (Word) - 唯一索引
- SensitiveWord: (Category)
- SensitiveWord: (MatchType)
- SensitiveWord: (IsActive)

#### 内容管理相关
- Content: (Type, Status)
- Content: (AuthorID)
- Content: (PublishAt)
- Content: (Status)

#### 统计相关
- DailyStatistics: (StatsDate) - 唯一索引
- PlayerStatistics: (PlayerID, StatsDate) - 唯一索引
- PlayerStatistics: (PlayerID)
- PlayerStatistics: (StatsDate)

#### 抽成/排名相关
- CommissionRule: (Type)
- CommissionRule: (GameID)
- CommissionRule: (PlayerID)
- CommissionRule: (IsActive)
- RankingCommissionConfig: (RankingType, Period)
- RankingCommissionConfig: (Month)
- RankingCommissionConfig: (IsActive)
- CommissionRecord: (OrderID)
- CommissionRecord: (PlayerID)
- CommissionRecord: (SettlementStatus)
- CommissionRecord: (SettlementMonth)
- CommissionRecord: (PlayerID, SettlementMonth)
- MonthlySettlement: (PlayerID, SettlementMonth) - 唯一索引
- MonthlySettlement: (Status)
- PlayerRanking: (PlayerID, RankingType, Period, PeriodValue) - 唯一索引
- PlayerRanking: (RankingType, Period, PeriodValue)
- PlayerRanking: (Rank)
- RankingReward: (RankingType, Period)
- RankingReward: (IsActive)

#### 团队相关
- Team: (LeaderID)
- Team: (Status)
- TeamMember: (TeamID, PlayerID) - 唯一索引
- TeamMember: (PlayerID, Status)
- TeamMember: (TeamID, Status)
- TeamInvite: (TeamID, PlayerID)
- TeamInvite: (PlayerID, Status)
- TeamInvite: (ExpireAt)

---

## 变更日志

| 日期 | 变更内容 |
|------|----------|
| 2024-12-25 | 补充缺失索引：Review.ExpireAt、Coupon.Source |
| 2024-12-25 | 业务逻辑问题修复：订单状态流转图、退款与争议规则、收入结算逻辑、多人订单评价、团队接单规则、拉黑交互规则、折扣叠加规则、礼物订单售后规则 |
| 2024-12-25 | 补充缺失枚举：Role、ServiceItemSubCategory、WithdrawMethod、WithdrawStatus、RankingType、RankingPeriod |
| 2024-12-25 | 补充缺失索引：团队相关索引、CommissionRecord (PlayerID, SettlementMonth) 复合索引 |
| 2024-12-25 | 完善游戏/服务项目/敏感词/内容管理/统计模块：新增枚举类型和索引 |
| 2024-12-25 | 完善用户/认证模块：新增 LoginType 枚举、AuthConfig 唯一索引、UserLoginLog 索引 |
| 2024-12-25 | 完善聊天模块业务流程：公共聊天室、订单群聊、消息审核、限流配置、30天记录保留、导出功能 |
| 2024-12-25 | 完善陪玩师模块业务流程：入驻审核、实名认证、段位认证、在线状态管理、接单开关、系统统一定价 |
| 2024-12-25 | 完善争议模块业务流程：双方可发起、双客服机制、争议类型模板、聊天记录快照、全额退款/驳回处理 |
| 2024-12-25 | 完善评价模块业务流程：新增 ReviewReply、ReviewAppeal 模型，评价7天窗口、修改3次、回复3次、差评申诉等规则 |
| 2024-12-25 | 完善支付/钱包业务流程：组合支付、退款顺序、T+7售后期、陪玩师收入结算等 |
| 2024-12-25 | 新增项目管理 steering 规则（06-project-management.md），AI 辅助维护模块状态 |
| 2024-12-25 | 新增用户拉黑系统：UserBlock 模型，支持双向拉黑和沙箱隔离 |
| 2024-12-25 | 新增订单超时处理系统：OrderTimeoutConfig、OrderTimeoutLog、OrderServiceAssignment 模型 |
| 2024-12-25 | 新增陪玩师等级/认证系统：GameRank、PlayerRankRecord、PlayerCertification 模型 |
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
