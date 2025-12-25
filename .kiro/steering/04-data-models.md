# 数据模型 - 核心模块

> ⚠️ **文档维护提醒**：修改 `backend/internal/model/` 后必须同步更新此文档。

## 相关文档

- [04a-marketing-models.md](./04a-marketing-models.md) - VIP、优惠券、充值、活动、推荐
- [04b-team-models.md](./04b-team-models.md) - 团队系统
- [04c-enums-indexes.md](./04c-enums-indexes.md) - 枚举类型和数据库索引
- [04d-notification-models.md](./04d-notification-models.md) - 通知系统

## 维护规范

1. **及时更新**：新增/修改/删除模型字段时，同步更新此文档
2. **版本追踪**：重大变更在 PROGRESS.md 中记录
3. **审核机制**：PR 中涉及模型变更时，检查此文档是否同步更新
4. **定期审计**：每月检查文档与代码的一致性

---

## 基础模型

### Base（通用基类）

```go
type Base struct {
    ID        uint64         `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"` // 软删除
    ExtJSON   string         `json:"extJson,omitempty"`   // 扩展字段（JSON格式）
}
```

> **ExtJSON 用途**：预留扩展字段，用于存储临时性或实验性数据，避免频繁修改表结构。

---

## 用户与权限模块

### User（用户）

| 字段 | 类型 | 说明 |
|------|------|------|
| Phone | string | 手机号，唯一索引 |
| Email | string | 邮箱，唯一索引 |
| PasswordHash | string | 密码哈希（不返回前端） |
| Name | string | 用户名 |
| AvatarURL | string | 头像URL |
| Role | Role | 主要角色：user/player/admin |
| Status | UserStatus | 状态：active/suspended/banned |
| BanReason | string | 封禁原因 |
| BannedAt | *time.Time | 封禁时间 |
| BannedBy | *uint64 | 封禁操作人ID |
| LastLoginAt | *time.Time | 最后登录时间 |
| Roles | []RoleModel | 多角色关联（RBAC） |
| Wallet | *Wallet | 用户钱包 |
| VipLevelID | *uint64 | 当前VIP等级ID |
| VipUnlocked | bool | VIP是否已解锁 |
| VipExp | int64 | VIP经验（累计消费，分） |
| TotalRechargeCents | int64 | 累计充值（分） |
| VipUnlockedAt | *time.Time | VIP解锁时间 |
| VipExpireAt | *time.Time | VIP过期时间（nil=永久） |
| LastMonthlyCouponAt | *time.Time | 上次发放月度券时间 |

### AuthConfig（认证系统配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `access_token_expire_hours` - Access Token 有效期（小时），默认24
- `refresh_token_expire_days` - Refresh Token 有效期（天），默认7
- `max_login_devices` - 最大同时登录设备数，默认5
- `password_min_length` - 密码最小长度，默认8
- `password_require_uppercase` - 密码需要大写字母，默认true
- `password_require_lowercase` - 密码需要小写字母，默认true
- `password_require_number` - 密码需要数字，默认true
- `sms_code_expire_minutes` - 短信验证码有效期（分钟），默认5
- `email_code_expire_minutes` - 邮箱验证码有效期（分钟），默认10

### UserLoginLog（用户登录日志）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID |
| LoginType | LoginType | 登录方式：password/sms/email/oauth |
| DeviceID | string | 设备标识 |
| DeviceType | string | 设备类型：web/ios/android/mini |
| IP | string | 登录IP |
| UserAgent | string | User-Agent |
| Location | string | 登录地点（IP解析） |
| LoginAt | time.Time | 登录时间 |
| LogoutAt | *time.Time | 登出时间 |
| TokenID | string | Token标识（用于踢出） |

### RoleModel（角色）

| 字段 | 类型 | 说明 |
|------|------|------|
| Slug | string | 角色标识（唯一）：superAdmin/admin/player/user |
| Name | string | 角色名称 |
| Description | string | 角色描述 |
| IsSystem | bool | 是否系统角色（不可删除） |
| ParentID | *uint64 | 父角色ID（继承） |
| Priority | int | 优先级（权限冲突解决，数字越大优先级越高） |
| Level | int | 继承层级（最大5层） |
| Permissions | []Permission | 关联权限 |

### Permission（权限）

| 字段 | 类型 | 说明 |
|------|------|------|
| Method | HTTPMethod | HTTP方法：GET/POST/PUT/PATCH/DELETE |
| Path | string | API路径 |
| Code | string | 语义化标识，格式：module.resource.action |
| Group | string | API分组 |
| Description | string | 权限描述 |
| ParentID | *uint64 | 父权限ID（树形结构） |
| IsSystem | bool | 是否系统权限 |

### Menu（菜单）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 菜单名称 |
| Path | string | 路由路径 |
| Component | string | 前端组件名 |
| Icon | string | 图标 |
| ParentID | *uint64 | 父菜单ID |
| Order | int | 排序 |
| Hidden | bool | 是否隐藏 |
| Permission | string | 所需权限码 |

### UserBlock（用户拉黑）

> 用户/陪玩师之间的拉黑关系，支持双向拉黑，实现沙箱隔离效果

| 字段 | 类型 | 说明 |
|------|------|------|
| BlockerID | uint64 | 拉黑发起人ID |
| BlockerType | BlockUserType | 发起人类型：user/player |
| BlockedID | uint64 | 被拉黑人ID |
| BlockedType | BlockUserType | 被拉黑人类型：user/player |
| Reason | string | 拉黑原因（可选） |
| Status | BlockStatus | 状态：active/canceled/admin_canceled |
| BlockedAt | time.Time | 拉黑时间 |
| CanceledAt | *time.Time | 取消时间 |
| CanceledBy | *uint64 | 取消人ID（管理员强制解除时） |
| AdminRemark | string | 管理员备注 |

**拉黑效果：**
- 无法发消息：双方无法互相发送消息
- 列表隐藏：被拉黑方在对方的列表中不可见
- 订单房间隔离：下单后的聊天房间中，被拉黑的陪玩师对用户不可见

**拉黑与进行中订单的交互规则：**
- 下单后拉黑：已进行中的订单继续执行，但聊天消息对拉黑方不可见
- 拉黑不影响订单完成：陪玩师仍可完成服务
- 拉黑不影响评价：用户仍可对被拉黑的陪玩师评价
- 拉黑后不可再下单：用户无法向被拉黑的陪玩师下单

**双向拉黑实现：**
- 每次拉黑创建一条记录（单向）
- A 拉黑 B：创建 BlockerID=A, BlockedID=B
- B 拉黑 A：创建 BlockerID=B, BlockedID=A（独立记录）
- 查询时需检查双向：WHERE (BlockerID=A AND BlockedID=B) OR (BlockerID=B AND BlockedID=A)

**取消拉黑：**
- 只能取消自己发起的拉黑
- A 取消对 B 的拉黑，不影响 B 对 A 的拉黑（如有）

---

## 用户与权限业务流程

### 注册方式

| 方式 | 说明 | 状态 |
|------|------|------|
| 手机号 | 手机号 + 验证码注册 | ✅ 已实现 |
| 邮箱 | 邮箱 + 验证码注册 | ✅ 已实现 |
| 第三方OAuth | 微信/QQ/微博等 | 🔜 预留 |

### 登录方式

| 方式 | 说明 | 状态 |
|------|------|------|
| 密码登录 | 手机号/邮箱 + 密码 | ✅ 已实现 |
| 验证码登录 | 手机号/邮箱 + 验证码 | ✅ 已实现 |
| 第三方OAuth | 微信/QQ/微博等 | 🔜 预留 |

### 密码规则

```
【密码要求】
- 最小长度：8位
- 必须包含：大写字母 + 小写字母 + 数字
- 示例：Abcdefg1358

【密码存储】
- 使用 bcrypt 哈希存储
- 不存储明文密码
```

### 用户状态管理

```
【状态流转】
active ←→ suspended ←→ banned

【封禁流程】
管理员发起封禁 → 选择封禁类型
  ├── suspended（暂停）：临时封禁，可解封
  └── banned（永久）：永久封禁

【封禁记录】
User.Status = suspended/banned
User.BanReason = "封禁原因"
User.BannedAt = now
User.BannedBy = 管理员ID

【解封流程】
管理员发起解封 → User.Status = active
  └── 清空封禁相关字段
```

### Token 管理

```
【双Token机制】
- Access Token：短期有效（默认24小时，可配置）
- Refresh Token：长期有效（默认7天，可配置）

【刷新流程】
Access Token 过期 → 使用 Refresh Token 获取新 Token
  ├── Refresh Token 有效 → 返回新的 Access Token + Refresh Token
  └── Refresh Token 过期 → 需要重新登录

【多设备登录】
- 支持多设备同时登录
- 最大设备数可配置（默认5个）
- 超过限制 → 踢出最早登录的设备

【设备管理】
- 用户可查看所有登录设备
- 用户可主动踢出指定设备
- 管理员可强制踢出用户所有设备
```

### RBAC 权限系统

```
【角色继承】
子角色继承父角色的所有权限
  └── 最大继承层级：5层

【权限冲突解决】
同一用户有多个角色时：
  └── 按 Priority 优先级决定（数字越大优先级越高）

【系统角色】
- superAdmin：超级管理员（不可删除）
- admin：普通管理员
- player：陪玩师
- user：普通用户

【权限检查流程】
请求到达 → 获取用户角色 → 合并所有角色权限
  → 检查请求路径是否在权限列表中
  → 有权限 → 放行
  → 无权限 → 返回 403
```

### 动态菜单

```
【菜单生成】
用户登录 → 获取用户角色 → 获取角色关联的菜单
  → 过滤隐藏菜单 → 按 Order 排序 → 返回菜单树

【菜单权限】
每个菜单关联一个权限码（Permission）
  └── 用户无该权限 → 菜单不显示
```

## 陪玩师模块

### Player（陪玩师）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 关联用户ID（唯一） |
| Nickname | string | 昵称 |
| Bio | string | 个人简介 |
| Rank | string | 段位（冗余，主段位名称） |
| RatingAverage | float32 | 平均评分（0-5） |
| RatingCount | uint32 | 评价数量 |
| HourlyRateCents | int64 | 时薪（分，冗余） |
| MainGameID | uint64 | 主要游戏ID |
| VerificationStatus | VerificationStatus | 入驻认证状态：pending/verified/rejected |
| VerifiedAt | *time.Time | 审核时间 |
| VerifiedBy | *uint64 | 审核人ID |
| RejectReason | string | 拒绝原因 |
| OnlineStatus | PlayerOnlineStatus | 在线状态：online/offline/busy |
| AcceptingOrders | bool | 接单开关（true=在线接单） |
| LastOnlineAt | *time.Time | 最后在线时间 |

### GameRank（游戏段位配置）

> 平台自定义的游戏段位，每个游戏可以有多个段位，用于陪玩师资质认证和定价

| 字段 | 类型 | 说明 |
|------|------|------|
| GameID | uint64 | 游戏ID |
| Name | string | 段位名称（青铜、白银、黄金、铂金、钻石、大师、王者） |
| Level | int | 段位等级（数字，用于排序，越大越高） |
| PriceCents | int64 | 该段位定价（分） |
| IconURL | string | 段位图标 |
| Color | string | 段位颜色（前端展示） |
| Description | string | 段位描述 |
| SortOrder | int | 排序 |
| IsActive | bool | 是否启用 |

### PlayerRankRecord（陪玩师段位认证）

> 陪玩师的游戏段位认证记录，一个陪玩师可以有多个游戏的多个段位

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID |
| GameID | uint64 | 游戏ID |
| RankID | uint64 | 段位ID |
| Status | PlayerRankStatus | 状态：pending/verified/rejected/revoked/expired |
| ScreenshotURLs | string | 段位截图（JSON数组） |
| VerifiedAt | *time.Time | 认证时间 |
| VerifiedBy | *uint64 | 审核人ID |
| RejectReason | string | 拒绝/撤销原因 |
| ExpireAt | *time.Time | 过期时间（预留，定期复审） |
| Remark | string | 备注 |

### PlayerCertification（陪玩师实名认证）

> 陪玩师的实名认证信息

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID（唯一） |
| RealName | string | 真实姓名 |
| IDCardNo | string | 身份证号（加密存储） |
| IDCardFrontURL | string | 身份证正面照片 |
| IDCardBackURL | string | 身份证背面照片 |
| Status | CertificationStatus | 状态：pending/verified/rejected |
| VerifiedAt | *time.Time | 认证时间 |
| VerifiedBy | *uint64 | 审核人ID |
| RejectReason | string | 拒绝原因 |
| PhotoURL | string | 个人照片（预留） |
| VoiceURL | string | 语音介绍（预留） |
| ExtJSON | string | 扩展字段（预留其他认证资料）

---

## 陪玩师业务流程

### 入驻方式

| 方式 | 说明 |
|------|------|
| 独立注册 | 新用户直接申请成为陪玩师 |
| 用户转陪玩师 | 已有用户申请转为陪玩师身份 |

### 入驻审核流程

```
【入驻流程】
用户提交入驻申请 → 联系管理员 → 管理员安排考核
  ├── 考核通过 → VerificationStatus = verified
  │     ├── User.Role = player
  │     └── 创建 Player 记录
  └── 考核不通过 → VerificationStatus = rejected
        └── 记录拒绝原因
```

### 实名认证流程

```
【实名认证】
陪玩师提交实名信息 → 上传身份证照片 → 管理员审核
  ├── 审核通过 → PlayerCertification.Status = verified
  └── 审核拒绝 → 通知陪玩师，可重新提交
```

### 段位认证流程

```
【段位认证】
陪玩师申请段位认证 → 联系管理员 → 管理员对接考官
  → 考官进行段位考核 → 考核结果
  ├── 通过 → PlayerRankRecord.Status = verified
  │     └── 可接该段位的订单
  └── 不通过 → PlayerRankRecord.Status = rejected
        └── 可重新申请

【认证有效期】
- 长期有效，无过期时间
- 客诉可能导致降级重考（revoked）

【降级处理】
客诉成立 → PlayerRankRecord.Status = revoked
  └── 陪玩师需重新申请认证
```

### 陪玩师状态管理

```
【在线状态规则】
┌─────────────────────────────────────────────────────────┐
│ 场景                    │ 对普通用户显示 │ 对被拉黑用户显示 │
├─────────────────────────────────────────────────────────┤
│ 接单开关开 + 无订单     │ 在线 (online)  │ 忙碌 (busy)      │
│ 接单开关开 + 有订单     │ 忙碌 (busy)    │ 忙碌 (busy)      │
│ 接单开关关              │ 离线 (offline) │ 离线 (offline)   │
└─────────────────────────────────────────────────────────┘

【接单开关】
- AcceptingOrders = true → 在线，可接单
- AcceptingOrders = false → 离线/忙碌，不接单

【长时间订单】
- 部分订单持续数天（如包天服务）
- 订单进行中 → 状态为忙碌
```

### 陪玩师定价

```
【定价规则】
- 系统统一定价，按段位定价
- 价格来源：GameRank.PriceCents
- 陪玩师不可自定义价格

【价格获取】
用户下单 → 根据选择的段位 → 获取 GameRank.PriceCents
```

### 陪玩师评分

```
【评分计算】
新评价提交 → 重新计算平均分
  Player.RatingAverage = 所有评价分数总和 / 评价数量
  Player.RatingCount++

【评分影响】（预留）
- 排名权重
- 推荐权重
- 首页展示优先级
```

## 订单模块

### Order（订单/主订单）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderNo | string | 订单号（唯一） |
| UserID | uint64 | 下单用户ID |
| ItemID | uint64 | 服务项目ID（向后兼容，单人订单） |
| PlayerID | *uint64 | 服务陪玩师ID（向后兼容，单人订单） |
| RecipientPlayerID | *uint64 | 礼物接收陪玩师ID |
| Quantity | int | 数量 |
| UnitPriceCents | int64 | 单价（分） |
| TotalPriceCents | int64 | 总价（分） |
| CommissionCents | int64 | 平台抽成（分） |
| PlayerIncomeCents | int64 | 陪玩师收入（分） |
| Currency | Currency | 货币：CNY/USD/EUR |
| Status | OrderStatus | 状态（见枚举） |
| GameID | *uint64 | 游戏ID |
| ScheduledStart | *time.Time | 预约开始时间 |
| CompletedAt | *time.Time | 完成时间 |
| HasDispute | bool | 是否有争议 |
| RequiredPlayers | int | 需要的陪玩师数量（默认1） |
| CurrentPlayers | int | 当前已接单的陪玩师数量 |

### OrderItem（订单明细）

> 支持一个订单包含多个服务项目/座位，每个座位可以是不同等级的陪玩师

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 主订单ID |
| ItemID | uint64 | 服务项目ID |
| Slot | int | 座位号（1, 2, 3...） |
| UnitPriceCents | int64 | 单价（分） |
| Quantity | int | 数量 |
| TotalCents | int64 | 小计（分） |
| CommissionRate | float64 | 抽成比例 |
| Status | OrderItemStatus | 状态：pending/matched/completed/canceled |
| PlayerID | *uint64 | 接单的陪玩师ID |
| ReviewID | *uint64 | 关联的评价ID（nil=未评价） |

### OrderPlayer（订单陪玩师记录）

> 记录陪玩师接单详情，用于收入分配和状态追踪

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| OrderItemID | uint64 | 订单明细ID |
| PlayerID | uint64 | 陪玩师ID |
| JoinedAt | time.Time | 接单时间 |
| TeamID | *uint64 | 团队ID（团队接单时） |
| IncomeCents | int64 | 该陪玩师收入（分） |
| CommissionCents | int64 | 该陪玩师抽成（分） |
| Status | OrderPlayerStatus | 状态：joined/left/completed |

### ServiceItem（服务项目）

| 字段 | 类型 | 说明 |
|------|------|------|
| ItemCode | string | 项目编码（唯一） |
| Name | string | 项目名称 |
| Category | string | 分类：escort |
| SubCategory | ServiceItemSubCategory | 子类别：solo/team/gift |
| GameID | *uint64 | 游戏ID |
| BasePriceCents | int64 | 基础价格（分） |
| ServiceHours | int | 服务时长（小时） |
| CommissionRate | float64 | 抽成比例（默认0.20） |
| IsActive | bool | 是否启用 |
| RequiredPlayers | int | 需要的陪玩师数量（默认1） |
| VipPriceCents | *int64 | VIP专属价格（分） |
| UsageLimitType | UsageLimitType | 限制类型：none/once/daily/weekly/monthly |
| UsageLimitCount | int | 限制次数（0=无限制） |
| MaxPerOrder | int | 单次购买数量限制（0=无限制） |

### OrderTimeoutConfig（订单超时配置）

> 订单超时相关的系统配置

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `payment_timeout_minutes` - 支付超时时间（分钟），默认30
- `order_accept_timeout_minutes` - 接单超时时间（分钟），默认30
- `auto_cancel_enabled` - 是否启用自动取消，默认true
- `auto_refund_enabled` - 是否启用自动退款，默认true
- `auto_assign_service_enabled` - 接单后是否自动分配客服，默认true

### OrderTimeoutLog（订单超时日志）

> 记录订单超时处理的日志

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| TimeoutType | OrderTimeoutType | 超时类型：accept_timeout |
| TimeoutAt | time.Time | 超时时间 |
| Action | OrderTimeoutAction | 处理动作：canceled/refunded/notified |
| RefundAmountCents | int64 | 退款金额（分） |
| RefundRecordID | *uint64 | 退款记录ID |
| Remark | string | 备注 |

### OrderServiceAssignment（订单客服分配）

> 记录订单接单后自动分配的客服

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| ServiceUserID | uint64 | 客服用户ID |
| ChatGroupID | *uint64 | 聊天群组ID |
| Status | ServiceAssignmentStatus | 状态：assigned/joined/left/completed |
| AssignedAt | time.Time | 分配时间 |
| JoinedAt | *time.Time | 加入房间时间 |
| LeftAt | *time.Time | 离开时间 |
| AssignType | string | 分配方式：auto/manual |
| Remark | string | 备注 |

---

## 订单业务流程

### 多人订单下单流程

```
1. 用户选择服务项目和座位配置
   ├── 座位1: 王者陪玩 ¥50
   └── 座位2: 钻石陪玩 ¥30
   
2. 创建订单
   ├── Order (TotalPrice=80, RequiredPlayers=2, CurrentPlayers=0)
   ├── OrderItem (Slot=1, ItemID=王者, Price=50, Status=pending)
   └── OrderItem (Slot=2, ItemID=钻石, Price=30, Status=pending)
   
3. 陪玩师接单（防超卖：数据库行锁）
   ├── 王者陪玩师接单 → OrderItem1.PlayerID=100, Status=matched
   └── 钻石陪玩师接单 → OrderItem2.PlayerID=200, Status=matched
   
4. 人齐检查
   └── CurrentPlayers == RequiredPlayers → Order.Status = in_progress
   
5. 服务完成
   └── Order.Status = completed
   
6. 用户评价（一键评价 + 可选单独评价）
   └── 批量创建 Review 记录
```

### 订单补差价（升级/加座位）

> 服务进行中（in_progress）用户可以升级座位等级或新增座位

**触发时机：** 陪玩师进入后（Order.Status = in_progress）

**支持操作：**
- 升级座位：钻石陪玩 → 王者陪玩
- 加座位：从 2 人变 3 人

```
【升级座位流程】
用户发起升级 → 计算差价 → 创建 Payment → 支付成功 → 更新订单

【加座位流程】
用户发起加座 → 创建新 OrderItem → 创建 Payment → 支付成功 → 等待接单
```

---

## 支付模块

### Payment（支付记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| UserID | uint64 | 用户ID |
| Method | PaymentMethod | 支付方式（见枚举） |
| AmountCents | int64 | 支付金额（分） |
| Status | PaymentStatus | 状态（见枚举） |
| ProviderTradeNo | string | 第三方交易号 |
| PaidAt | *time.Time | 支付时间 |
| RefundedAmountCents | int64 | 已退款金额（分） |
| WalletAmountCents | int64 | 钱包支付金额（组合支付） |

### Wallet（钱包）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID（唯一） |
| BalanceCents | int64 | 可用余额（分） |
| FrozenCents | int64 | 冻结金额（分） |

### Withdraw（提现）

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID |
| AmountCents | int64 | 提现金额（分） |
| Method | WithdrawMethod | 提现方式：alipay/wechat/bank |
| Status | WithdrawStatus | 状态（见枚举） |
| TaxDeductedCents | int64 | 代扣个税（分） |
| ActualAmountCents | int64 | 实际到账（分） |

---

## 支付/钱包业务流程

### 支付方式

| 方式 | 说明 |
|------|------|
| wechat | 纯微信支付 |
| alipay | 纯支付宝支付 |
| wallet | 纯钱包余额支付 |
| combined | 组合支付（钱包 + 微信/支付宝） |

### 支付流程

```
【纯第三方支付】
用户下单 → 选择微信/支付宝 → 调用支付接口 → 等待回调
  → 支付成功 → 更新 Payment.Status = paid → 更新订单状态
  → 支付超时（可配置） → 自动取消订单

【纯钱包支付】
用户下单 → 选择钱包支付 → 检查余额是否充足
  → 充足 → 扣减余额 → Payment.Status = paid → 更新订单状态
  → 不足 → 提示充值或选择其他支付方式

【组合支付】
用户下单 → 选择组合支付 → 输入钱包支付金额
  → 扣减钱包余额（WalletAmountCents）
  → 剩余金额调用第三方支付
  → 第三方支付成功 → Payment.Status = paid
  → 第三方支付失败/超时 → 退还钱包余额
```

### 钱包余额来源

| 来源 | 说明 |
|------|------|
| 充值 | 用户主动充值 |
| 订单退款 | 退款退回钱包部分 |
| 陪玩师收入 | 订单完成后 T+7 进入（扣除平台抽成） |
| 活动奖励 | 活动发放的现金奖励 |

### 冻结余额场景

```
【提现冻结】
陪玩师申请提现 → 冻结提现金额
  → 提现成功 → 扣减冻结金额
  → 提现失败/取消 → 解冻，返还可用余额
```

### 退款流程

```
【退款与争议规则】
┌─────────────────────────────────────────────────────────────────┐
│ 时间窗口          │ 可用方式           │ 说明                    │
├─────────────────────────────────────────────────────────────────┤
│ 订单完成后7天内   │ 申请退款或发起争议 │ 退款需管理员审核         │
│                   │                    │ 争议由客服介入处理       │
├─────────────────────────────────────────────────────────────────┤
│ 超过7天           │ 不支持退款/争议    │ 售后期结束               │
└─────────────────────────────────────────────────────────────────┘

【退款规则】
- 订单完成后 7 天内：用户可申请退款（管理员审核）
- 退款被拒绝：7 天内可发起争议，由客服协调处理
- 超过 7 天：不支持退款/争议
- 不支持部分退款，只能全额退款
- 退款需要管理员审核

【退款执行顺序（组合支付）】
1. 先退第三方支付部分（微信/支付宝原路退回）
2. 确认第三方退款成功
3. 再退钱包支付部分（返还钱包余额）
4. 更新订单状态为 refunded

【退款流程】
用户申请退款 → 管理员审核
  → 审核通过 → 执行退款（按上述顺序）
  → 审核拒绝 → 通知用户（24小时内可发起争议）
```

### 陪玩师收入结算

```
【结算规则】
- 订单完成后，收入进入陪玩师钱包的 FrozenCents（冻结余额）
- T+7 售后冻结期，期间不可提现
- 平台抽成在进入钱包时扣除

【收入状态追踪】
使用 CommissionRecord.SettlementStatus 追踪收入状态：
  - pending：T+7 售后期内，收入冻结中
  - disputed：发生争议，收入继续冻结
  - settled：T+7 结束且无争议，收入解冻

【计算公式】
陪玩师收入 = 订单金额 × (1 - 平台抽成比例)

【售后期处理】
- T+7 期间发生退款：从陪玩师 FrozenCents 扣回收入
- T+7 期间发生争议：CommissionRecord.SettlementStatus = disputed，收入继续冻结
- T+7 结束无问题：FrozenCents → BalanceCents，收入解冻可提现

【示例】
订单金额：¥100，平台抽成：20%
陪玩师收入：¥100 × 80% = ¥80
订单完成 → ¥80 进入 FrozenCents（冻结状态）
T+7 后 → ¥80 从 FrozenCents 转入 BalanceCents，可提现
```

### 提现流程

```
【提现条件】
- 只有陪玩师可以提现
- 只能提现已解冻的余额（T+7 售后期结束）
- 提现方式：支付宝/微信/银行卡

【提现流程】
陪玩师申请提现 → 检查可提现余额
  → 余额充足 → 冻结提现金额 → 管理员审核
  → 审核通过 → 执行打款 → 扣减冻结金额
  → 审核拒绝/打款失败 → 解冻，返还可用余额

【税务处理】
提现时代扣个税（如适用）
实际到账 = 提现金额 - 代扣个税
```

---

## 游戏模块

### Game（游戏）

| 字段 | 类型 | 说明 |
|------|------|------|
| Key | string | 游戏标识（唯一），如 lol/dota2 |
| Name | string | 游戏名称 |
| Category | string | 分类（管理员自定义）：moba/fps/rpg/card 等 |
| IconURL | string | 图标URL |
| CoverURL | string | 封面图URL |
| Description | string | 游戏描述 |
| IsActive | bool | 是否上架 |
| SortOrder | int | 排序 |

---

## 游戏业务流程

### 游戏管理

```
【游戏分类】
- 分类由管理员自定义（moba/fps/rpg/card 等）
- 不需要审核流程

【上下架规则】
- 后台手动上下架
- IsActive = true → 前端可见
- IsActive = false → 前端隐藏
```

### 游戏与段位关系

```
【关系说明】
- 游戏与段位：一对多（一个游戏有多个段位）
- 段位与陪玩师：多对多（一个陪玩师可拥有多个游戏的多个段位）

【段位定价】
- 段位决定价格
- 价格来源：GameRank.PriceCents
- 系统统一定价，陪玩师不可自定义

【段位配置示例】
王者荣耀:
  ├── 青铜 → ¥20/小时
  ├── 白银 → ¥25/小时
  ├── 黄金 → ¥30/小时
  ├── 铂金 → ¥35/小时
  ├── 钻石 → ¥40/小时
  ├── 大师 → ¥50/小时
  └── 王者 → ¥60/小时
```

---

## 服务项目模块

### ServiceItem（服务项目）业务流程

```
【服务项目类型】
- solo（单人陪玩）：RequiredPlayers = 1
- team（组队陪玩）：RequiredPlayers > 1
- gift（礼物）：直接到陪玩师钱包

【类型判断优先级】
1. SubCategory = gift → 礼物订单（RequiredPlayers 忽略）
2. RequiredPlayers = 1 → 单人陪玩
3. RequiredPlayers > 1 → 组队陪玩

SubCategory 字段用于前端展示和快速筛选，礼物订单必须通过此字段判断。

【礼物订单】
用户下单礼物 → 支付成功 → 金额进入陪玩师钱包
  ├── 不需要服务流程
  ├── 收入即时到账（无 T+7 售后期）
  ├── 礼物订单不支持退款
  └── 平台抽成在进入钱包时扣除

【礼物订单特殊处理】
- RequiredPlayers = 1（固定值，但不参与匹配逻辑）
- 无需陪玩师接单流程
- 支付成功即完成
- Order.Status 直接设为 completed

【礼物订单售后规则】
- 不支持退款：礼物一经送出，不可退款
- 不支持争议：礼物订单不进入争议流程
- 误操作处理：如用户误操作送错礼物，可联系客服申诉
  ├── 客服核实后可人工处理（特殊情况）
  └── 需提供充分证据证明误操作
- 收入即时到账：陪玩师收入直接进入 BalanceCents（可用余额）
```

### 服务项目定价

```
【定价规则】
- 价格由系统统一设定
- 陪玩师不可自定义价格
- 价格来源：GameRank.PriceCents（按段位定价）

【VIP价格】
- VipPriceCents 字段预留
- VIP用户可享受专属价格
- 未设置时使用基础价格
```

### 使用限制

```
【限制类型】
- none：无限制
- once：仅可购买一次（如新人专享）
- daily：每日限购
- weekly：每周限购
- monthly：每月限购

【限制对象】
限制针对用户，非陪玩师

【限制示例】
新人首单优惠：UsageLimitType = once, UsageLimitCount = 1
每日特惠：UsageLimitType = daily, UsageLimitCount = 3
```

---

## 敏感词模块

### SensitiveWord（敏感词）

| 字段 | 类型 | 说明 |
|------|------|------|
| Word | string | 敏感词 |
| Category | SensitiveWordCategory | 分类：politics/porn/abuse/ad/other |
| MatchType | SensitiveWordMatchType | 匹配类型：exact/fuzzy/regex |
| Severity | SensitiveWordSeverity | 严重程度：low/medium/high |
| Replacement | string | 替换内容（默认 ***） |
| IsActive | bool | 是否启用 |
| CreatedBy | uint64 | 创建人ID |

---

## 敏感词业务流程

### 敏感词来源

```
【来源】
- 管理员手动添加
- 支持批量导入/导出（CSV/Excel）

【分类】
- politics：政治敏感
- porn：色情
- abuse：辱骂
- ad：广告
- other：其他
```

### 敏感词匹配规则

```
【匹配类型】
- exact（精确匹配）：完全匹配词汇
- fuzzy（模糊匹配）：包含即命中
- regex（正则匹配）：正则表达式匹配

【匹配优先级】
regex > exact > fuzzy
```

### 命中后处理

```
【处理流程】
内容提交 → 敏感词检测
  ├── 未命中 → 直接通过
  └── 命中 → 拦截，进入人工审核队列

【人工审核】
审核员查看内容 → 做出决定
  ├── 通过 → 内容发布
  ├── 拒绝 → 通知用户，内容不发布
  └── 替换 → 一键将敏感词替换为 ***，内容发布
```

---

## 内容管理模块

### Content（内容）

| 字段 | 类型 | 说明 |
|------|------|------|
| Title | string | 标题 |
| Type | ContentType | 类型：announcement/help/agreement/dynamic |
| Content | string | 内容（富文本/Markdown） |
| CoverURL | string | 封面图URL |
| AuthorID | uint64 | 作者ID |
| Status | ContentStatus | 状态：draft/pending/published/rejected/archived |
| PublishAt | *time.Time | 发布时间（定时发布） |
| ExpireAt | *time.Time | 过期时间 |
| ViewCount | int | 浏览次数 |
| SortOrder | int | 排序 |
| IsTop | bool | 是否置顶 |
| TargetType | string | 目标类型：all/user/player |

---

## 内容管理业务流程

### 内容类型

| 类型 | 说明 | 示例 |
|------|------|------|
| announcement | 公告 | 系统公告、活动通知 |
| help | 帮助文档 | 使用指南、FAQ |
| agreement | 协议 | 用户协议、隐私政策 |
| dynamic | 动态 | 陪玩师动态、平台动态 |

### 内容发布流程

```
【发布流程】
创建内容 → 保存草稿 → 提交审核 → 审核通过 → 发布

【状态流转】
draft（草稿）→ pending（待审核）→ published（已发布）
                    ↓
               rejected（已拒绝）→ draft（修改后重新提交）

【定时发布】
设置 PublishAt → 审核通过后 → 到达发布时间 → 自动发布

【内容审核】
- 所有内容需要审核
- 敏感词过滤
- 管理员审核通过后发布
```

### 陪玩师动态

```
【动态发布】
陪玩师发布动态 → 敏感词过滤 → 审核
  ├── 审核通过 → 发布到个人主页
  └── 审核拒绝 → 通知陪玩师

【动态内容】
- 文字 + 图片
- 支持多图（最多9张）
- 可设置公开/仅粉丝可见
```

---

## 统计模块

### Statistics（统计数据）

> 各模块可量化数据的统计汇总

### PlatformStatistics（平台每日统计）

> 代码中表名为 `platform_statistics`，按日期汇总平台数据

| 字段 | 类型 | 说明 |
|------|------|------|
| StatDate | time.Time | 统计日期（唯一索引） |
| DailyOrderCount | int | 日订单数 |
| DailyCompletedCount | int | 日完成订单数 |
| DailyCanceledCount | int | 日取消订单数 |
| DailyGMVCents | int64 | 日GMV（分） |
| DailyCommissionCents | int64 | 日抽成（分） |
| DailyRefundAmountCents | int64 | 日退款金额（分） |
| DailyNewUserCount | int | 日新增用户数 |
| DailyActiveUserCount | int | 日活跃用户数 |
| DailyPayingUserCount | int | 日付费用户数 |
| DailyNewPlayerCount | int | 日新增陪玩师数 |
| DailyActivePlayerCount | int | 日活跃陪玩师数 |
| DailyRechargeCents | int64 | 日充值金额（分） |
| DailyWithdrawCents | int64 | 日提现金额（分） |
| DailyRechargeCount | int | 日充值笔数 |
| DailyWithdrawCount | int | 日提现笔数 |
| DailyDisputeCount | int | 日争议数 |
| DailyResolvedCount | int | 日解决争议数 |
| DailySLABreachCount | int | 日SLA超时数 |

### PlayerStatistics（陪玩师统计）

> 代码中表名为 `player_statistics`，累计统计（非按日期）

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID（唯一索引） |
| TotalEarningsCents | int64 | 累计收入（分） |
| TotalCommissionCents | int64 | 累计被抽成（分） |
| TotalWithdrawCents | int64 | 累计提现（分） |
| PendingWithdrawCents | int64 | 待提现（分） |
| TotalOrderCount | int | 累计接单数 |
| CompletedOrderCount | int | 完成订单数 |
| CanceledOrderCount | int | 取消订单数 |
| RefundOrderCount | int | 退款订单数 |
| TotalServiceMinutes | int | 累计服务时长（分钟） |
| AvgResponseTimeSec | int | 平均响应时间（秒） |
| AvgOrderAmountCents | int64 | 平均订单金额（分） |
| TotalCustomerCount | int | 累计服务客户数 |
| RepeatCustomerCount | int | 回头客数量 |
| RepeatOrderRate | float32 | 复购率 |
| DisputeCount | int | 被投诉次数 |
| DisputeWonCount | int | 投诉胜诉次数 |
| DisputeLostCount | int | 投诉败诉次数 |
| GiftReceivedCount | int | 收到礼物数 |
| GiftReceivedAmountCents | int64 | 收到礼物金额（分） |
| FirstOrderAt | *time.Time | 首次接单时间 |
| LastOrderAt | *time.Time | 最后接单时间 |
| LastActiveAt | *time.Time | 最后活跃时间 |

### UserStatistics（用户统计）

> 代码中表名为 `user_statistics`，用户消费行为统计

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID（唯一索引） |
| TotalSpentCents | int64 | 累计消费金额（分） |
| TotalOrderCount | int | 累计订单数 |
| CompletedOrderCount | int | 完成订单数 |
| CanceledOrderCount | int | 取消订单数 |
| RefundOrderCount | int | 退款订单数 |
| AvgOrderAmountCents | int64 | 平均订单金额（分） |
| DisputeCount | int | 发起争议次数 |
| DisputeWonCount | int | 争议胜诉次数 |
| DisputeLostCount | int | 争议败诉次数 |
| TotalRechargeCents | int64 | 累计充值金额（分） |
| RechargeCount | int | 充值次数 |
| AvgRechargeAmountCents | int64 | 平均充值金额（分） |
| ReviewCount | int | 评价次数 |
| AvgReviewScore | float32 | 平均评分 |
| FirstOrderAt | *time.Time | 首次下单时间 |
| LastOrderAt | *time.Time | 最后下单时间 |
| LastRechargeAt | *time.Time | 最后充值时间 |

---

## 统计业务流程

### 统计维度

```
【平台统计】
- 用户统计：新增、活跃、留存
- 订单统计：数量、金额、完成率、取消率
- 收入统计：充值、提现、平台抽成
- 争议统计：数量、处理时长、解决率

【陪玩师统计】
- 接单统计：数量、完成率
- 收入统计：总收入、日/周/月收入
- 评价统计：评分、评价数
- 服务统计：服务时长

```

### 统计周期

```
【周期类型】
- 日统计：每日凌晨自动汇总前一天数据
- 周统计：每周一汇总上周数据
- 月统计：每月1日汇总上月数据

【实时统计】
- 今日数据：实时计算
- 历史数据：从统计表读取
```

### 统计数据用途

```
【管理后台仪表盘】
- 平台概览：今日/本周/本月核心指标
- 趋势图表：订单、收入、用户增长趋势
- 排行榜：陪玩师收入排行、订单排行

【陪玩师收益报表】
- 收入明细：日/周/月收入
- 订单统计：接单数、完成率
- 评价统计：评分趋势

【运营优化】
- 用户行为分析
- 转化率分析
- 活动效果评估
```

---

## 评价模块

### Review（评价）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| OrderItemID | *uint64 | 订单明细ID（多人订单时关联具体座位） |
| UserID | uint64 | 评价者ID |
| PlayerID | uint64 | 被评价陪玩师ID |
| Score | Rating | 评分（1-5） |
| Content | string | 评价内容（不填使用默认模板） |
| Status | ReviewStatus | 状态：pending/approved/rejected/deleted |
| IsReported | bool | 是否被举报 |
| Images | StringArray | 评价图片URL数组（最多3张） |
| IsPublic | bool | 是否公开（默认false） |
| IsAnonymous | bool | 是否匿名 |
| EditCount | int | 修改次数（最多3次） |
| LastEditAt | *time.Time | 最后修改时间 |
| ExpireAt | time.Time | 评价截止时间（订单完成后7天） |

### ReviewReply（评价回复）

> 陪玩师对评价的回复

| 字段 | 类型 | 说明 |
|------|------|------|
| ReviewID | uint64 | 评价ID |
| PlayerID | uint64 | 陪玩师ID |
| Content | string | 回复内容 |
| ReplyCount | int | 回复次数（最多3次） |
| Status | ReviewReplyStatus | 状态：pending/approved/rejected |

### ReviewAppeal（评价申诉）

> 陪玩师对差评的申诉

| 字段 | 类型 | 说明 |
|------|------|------|
| ReviewID | uint64 | 评价ID |
| PlayerID | uint64 | 陪玩师ID |
| Reason | string | 申诉原因 |
| EvidenceURLs | string | 证据截图（JSON数组） |
| Status | AppealStatus | 状态：pending/approved/rejected |
| HandledBy | *uint64 | 处理人ID |
| HandledAt | *time.Time | 处理时间 |
| HandleRemark | string | 处理备注 |

---

## 评价业务流程

### 评价时机

```
【评价窗口】
订单完成 → 开启7天评价窗口
  └── 订单完成后弱提醒用户评价

【超时处理】
7天内未评价 → 评价窗口关闭，不再可评价
```

### 评价内容

| 项目 | 规则 |
|------|------|
| 评分 | 1-5 星，必填 |
| 文字 | 选填，不填使用默认模板（如"用户默认好评"） |
| 图片 | 选填，最多3张 |
| 公开 | 默认不公开 |
| 匿名 | 可选匿名评价 |

### 多人订单评价

```
【一键评价】
用户选择统一评分 → 为所有陪玩师创建相同评价
  └── 每个 OrderItem 对应一条 Review 记录

【单独评价】
用户选择单独评价 → 分别为每个陪玩师打分
  └── 可以给不同陪玩师不同评分和内容

【部分评价规则】
- 7天评价窗口内，用户可以分多次评价
- 每次评价可以只评价部分陪玩师
- 已评价的陪玩师不可重复评价（但可修改）
- 评价窗口关闭时，未评价的陪玩师不自动好评
- 未评价的陪玩师不计入其评分统计

【评价状态追踪】
OrderItem.ReviewID：关联的评价ID
  ├── nil：未评价
  └── 有值：已评价
```

### 评价审核

```
【审核流程】
用户提交评价 → 敏感词过滤
  ├── 命中敏感词 → 人工审核队列
  │     ├── 审核通过 → 发布评价
  │     └── 审核拒绝 → 通知用户
  └── 未命中 → 直接发布（陪玩师可申诉）

【申诉流程】
陪玩师发起申诉 → 提交申诉原因和证据 → 人工审核
  ├── 申诉成功 → 评价删除/隐藏，通知用户
  └── 申诉失败 → 通知陪玩师
```

### 差评处理

```
【差评定义】
1 星评价 = 差评

【差评通知】
差评发布 → 通知陪玩师
  └── 陪玩师可选择回复或申诉

【差评影响】（预留字段）
- 是否影响陪玩师评分：Player.RatingAverage
- 是否影响排名/推荐权重：预留接口
```

### 评价修改

```
【修改规则】
- 修改次数：最多3次
- 修改时间：7天内（评价窗口期内）
- 修改后：重新进入审核流程

【修改后状态变更】
Review.Status = pending（重新进入审核）
Review.EditCount++
Review.LastEditAt = now

【审核通过后】
Review.Status = approved
```

### 陪玩师回复

```
【回复规则】
- 回复次数：最多3次
- 回复审核：需要敏感词过滤

【回复流程】
陪玩师提交回复 → 敏感词过滤 → 发布/人工审核
```

### 评价展示

```
【默认展示】
- 评价默认不公开（仅用户和陪玩师可见）
- 用户可选择公开评价

【匿名评价】
- 用户选择匿名 → 陪玩师看不到评价者信息
- 管理后台可见真实用户
```

---

## 争议模块

### OrderDispute（订单争议）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| InitiatorID | uint64 | 发起人ID |
| InitiatorType | DisputeInitiatorType | 发起人类型：user/player |
| Type | DisputeType | 争议类型 |
| Status | DisputeStatus | 状态（见枚举） |
| Reason | string | 争议原因（用户填写） |
| EvidenceURLs | string | 证据截图URL列表（JSON数组，最多5张） |
| EvidenceText | string | 文字证据说明 |
| ChatSnapshotID | *uint64 | 聊天记录快照ID |
| OriginalServiceID | *uint64 | 原客服ID |
| AssignedServiceID | *uint64 | 分配的独立客服ID |
| SLADeadline | *time.Time | SLA截止时间（30分钟） |
| Resolution | DisputeResolution | 处理决定：full_refund/rejected |
| ResolvedBy | *uint64 | 处理人ID |
| ResolvedAt | *time.Time | 处理时间 |
| ResolveRemark | string | 处理备注 |

### DisputeTemplate（争议类型模板）

> 预设的争议类型，用户选择后再填写具体原因

| 字段 | 类型 | 说明 |
|------|------|------|
| Code | string | 模板编码（唯一） |
| Name | string | 模板名称 |
| InitiatorType | DisputeInitiatorType | 适用发起人：user/player/both |
| Description | string | 描述说明 |
| SortOrder | int | 排序 |
| IsActive | bool | 是否启用 |

**预设模板：**

| 编码 | 名称 | 适用方 |
|------|------|--------|
| `service_quality` | 服务质量问题 | 用户 |
| `bad_attitude` | 态度问题 | 用户 |
| `incomplete_service` | 未完成服务 | 用户 |
| `user_not_cooperative` | 用户不配合/不听指挥 | 陪玩师 |
| `user_harassment` | 用户骚扰 | 陪玩师 |
| `other` | 其他 | 双方 |

### ChatSnapshot（聊天记录快照）

> 争议发起时自动保存的聊天记录快照

| 字段 | 类型 | 说明 |
|------|------|------|
| DisputeID | uint64 | 争议ID |
| OrderID | uint64 | 订单ID |
| ChatGroupID | uint64 | 聊天群组ID |
| Messages | string | 聊天记录（JSON数组） |
| SnapshotAt | time.Time | 快照时间 |

---

## 争议业务流程

### 争议发起条件

| 条件 | 说明 |
|------|------|
| 发起人 | 用户或陪玩师均可发起 |
| 时间窗口 | 订单进行中 + 订单完成后7天内 |
| 限制 | 同一订单只能有一个进行中的争议 |

### 争议类型

**用户可选类型：**
- 服务质量问题
- 态度问题
- 未完成服务
- 其他

**陪玩师可选类型：**
- 用户不配合/不听指挥
- 用户骚扰
- 其他

### 争议证据

| 类型 | 限制 |
|------|------|
| 截图 | 最多5张 |
| 文字说明 | 无限制 |
| 聊天记录 | 自动获取当前订单聊天内容 |

### 客服分配规则

```
【双客服机制】
争议发起 → 分配两个客服
  ├── 原客服：订单原有的客服（如有）
  └── 独立客服：与该订单无关的客服（保证公正）

【分配逻辑】
1. 查找订单原客服 → OriginalServiceID
2. 从客服池中选择一个与该订单无关的客服 → AssignedServiceID
3. 两个客服共同处理争议
```

### 争议处理流程

```
【状态流转】
pending → assigned → mediating → closed

【详细流程】
1. 用户/陪玩师发起争议
   ├── 选择争议类型（模板）
   ├── 填写具体原因
   ├── 上传证据截图（最多5张）
   └── 系统自动保存聊天记录快照

2. 系统处理
   ├── 创建 OrderDispute 记录
   ├── Order.HasDispute = true
   ├── Order.Status = disputed（争议中）
   ├── 冻结陪玩师收入（如已结算）
   └── 分配双客服

3. 客服介入（SLA 30分钟）
   ├── 查看争议详情和证据
   ├── 查看聊天记录快照
   └── 与双方沟通调解

4. 处理决定
   ├── 全额退款 → Resolution = full_refund
   │     ├── 执行退款流程
   │     ├── 扣回陪玩师收入（如已结算）
   │     └── Order.Status = refunded
   └── 驳回 → Resolution = rejected
         ├── 解冻陪玩师收入
         └── Order.Status 恢复原状态

5. 通知双方处理结果
```

### 争议期间订单状态

```
【订单状态变更】
原状态（in_progress/completed）→ disputed（争议中）

【争议期间限制】
- 订单不可取消
- 订单不可评价
- 陪玩师收入冻结

【争议结束后】
- 全额退款 → Order.Status = refunded
- 驳回 → Order.Status 恢复为争议前状态
```

### 争议与收入的关系

```
【收入冻结】
争议发起 → 检查陪玩师收入状态
  ├── 未结算：标记为争议冻结，暂不结算
  └── 已结算（T+7内）：冻结对应金额

【争议结束】
全额退款 → 扣回陪玩师收入
驳回 → 解冻收入，正常结算
```

---

## 聊天模块

### ChatGroup（聊天群组）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupName | string | 群组名称 |
| GroupType | ChatGroupType | 类型：public/order |
| RelatedOrderID | *uint64 | 关联订单ID（订单群组） |
| CreatedBy | uint64 | 创建者ID |
| MaxMembers | int | 最大成员数（默认100） |
| IsActive | bool | 是否激活 |
| AutoDestroy | bool | 订单完成后自动销毁 |
| DeactivatedAt | *time.Time | 停用时间 |
| MessageRetentionDays | int | 消息保留天数（默认30） |
| VoiceEnabled | bool | 是否启用语音（预留） |
| VoiceRoomID | string | 语音房间ID（第三方服务） |
| VoiceProvider | string | 语音服务商：agora/tencent/zego |

### ChatMessage（聊天消息）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupID | uint64 | 群组ID |
| SenderID | uint64 | 发送者ID |
| Content | string | 消息内容 |
| MessageType | ChatMessageType | 类型：text/image/file/system/voice/emoji |
| AuditStatus | ChatMessageAuditStatus | 审核状态：pending/approved/rejected |
| VoiceURL | string | 语音消息URL（预留） |
| VoiceDuration | int | 语音时长秒数（预留） |
| EmojiCode | string | 表情包编码（预留） |

### ChatConfig（聊天系统配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `public_chat_rate_limit_seconds` - 公共聊天室发言间隔（秒），默认5
- `order_chat_rate_limit_seconds` - 订单群聊发言间隔（秒），默认0（无限制）
- `message_retention_days` - 消息保留天数，默认30
- `max_image_size_mb` - 图片最大大小（MB），默认5
- `max_file_size_mb` - 文件最大大小（MB），默认10

### ChatGroupMember（聊天群组成员）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupID | uint64 | 群组ID |
| UserID | uint64 | 用户ID |
| Role | ChatMemberRole | 角色：owner/admin/member |
| JoinedAt | time.Time | 加入时间 |
| LeftAt | *time.Time | 离开时间 |
| IsMuted | bool | 是否被禁言 |
| MutedUntil | *time.Time | 禁言截止时间 |

---

## 聊天业务流程

> **注意**：`social.go` 中已有简单的 `Notification` 模型（表名 `notifications`），用于基础通知场景。详见 [04d-notification-models.md](./04d-notification-models.md) 中"与现有 Notification 的关系"章节。

### 聊天类型

| 类型 | 说明 | 私聊支持 |
|------|------|----------|
| public | 公共聊天室 | ❌ 不支持私聊 |
| order | 订单群聊 | ❌ 不支持私聊 |

> **注意**：平台不支持用户与陪玩师私聊，所有沟通必须在公共聊天室或订单群聊中进行

### 公共聊天室

```
【创建规则】
- 管理员可创建公共聊天室（频道）
- 所有用户可加入公共聊天室

【发言规则】
- 所有人都可以发言
- 限流：每 N 秒只能发一条消息（可配置，默认5秒）
- 敏感词过滤

【成员管理】
- 管理员可禁言用户
- 管理员可踢出用户
```

### 订单群聊

```
【创建时机】
订单支付成功 → 自动创建订单群聊
  └── 初始成员：用户（老板）+ 客服

【成员变化】
陪玩师接单 → 自动加入群聊
  └── 群聊成员：用户 + 客服 + 陪玩师（1~N个）

【群聊生命周期】
订单支付 → 创建群聊
  ↓
订单进行中 → 群聊活跃
  ↓
订单完成/取消 → 群聊标记为销毁
  ↓
30天后 → 聊天记录清理
```

### 消息类型

| 类型 | 说明 | 状态 |
|------|------|------|
| text | 文字消息 | ✅ 已实现 |
| image | 图片消息 | ✅ 已实现 |
| file | 文件消息 | ✅ 已实现 |
| system | 系统消息 | ✅ 已实现 |
| voice | 语音消息 | 🔜 预留 |
| emoji | 表情包 | 🔜 预留 |

### 消息审核

```
【审核流程】
用户发送消息 → 敏感词过滤
  ├── 未命中敏感词 → 直接发送，AuditStatus = approved
  └── 命中敏感词 → 拦截，AuditStatus = pending
        └── 进入人工审核队列
              ├── 审核通过 → 发送消息
              └── 审核拒绝 → 通知用户
```

### 聊天记录

```
【保存规则】
- 默认保存30天（可配置）
- 订单群聊：订单完成后30天清理
- 公共聊天室：滚动清理30天前的消息

【导出功能】
- 用户可导出自己参与的聊天记录
- 管理员可导出任意聊天记录
- 导出格式：JSON/CSV
```

### 限流规则

```
【公共聊天室】
- 发言间隔：默认5秒（可配置）
- 超过限制 → 提示"发言太频繁，请稍后再试"

【订单群聊】
- 发言间隔：默认无限制
- 可根据需要配置
```

---

## 抽成计算规则

### 三维度抽成体系

平台抽成由三个维度共同决定，最终抽成 = 基础抽成 + 排名调整：

| 维度 | 数据来源 | 说明 |
|------|----------|------|
| 项目抽成 | `ServiceItem.CommissionRate` | 服务项目级别的基础抽成（默认 20%） |
| 陪玩师个人抽成 | `CommissionRule.PlayerID` | 特定陪玩师的专属抽成比例 |
| 上月排名抽成 | `RankingCommissionConfig` | 根据上月排名的阶梯式抽成调整 |

### CommissionRule（抽成规则）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 规则名称 |
| Description | string | 规则描述 |
| Type | string | 类型：default/special/gift |
| Rate | int | 抽成比例（百分比，20=20%） |
| IsActive | bool | 是否启用 |
| GameID | *uint64 | 特定游戏的抽成（可选） |
| PlayerID | *uint64 | 特定陪玩师的抽成（可选） |
| ServiceType | *string | 特定服务类型的抽成（可选） |

### RankingCommissionConfig（排名抽成配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 配置名称 |
| RankingType | RankingType | 排名类型：income/order_count |
| Period | string | 周期：monthly |
| Month | string | 月份（YYYY-MM） |
| RulesJSON | string | 阶梯规则（JSON数组） |
| Description | string | 描述 |
| IsActive | bool | 是否启用 |

**阶梯规则示例：**
```json
[
  {"rankStart": 1, "rankEnd": 10, "commissionRate": 5},
  {"rankStart": 11, "rankEnd": 50, "commissionRate": 3},
  {"rankStart": 51, "rankEnd": 100, "commissionRate": 2}
]
```

### 抽成计算流程

```
【抽成计算优先级】
1. 检查陪玩师是否有个人专属抽成（CommissionRule.PlayerID）
   └── 有 → 使用个人抽成比例
   └── 无 → 继续下一步

2. 检查服务项目抽成（ServiceItem.CommissionRate）
   └── 使用项目设置的抽成比例（默认 20%）

3. 检查上月排名抽成调整（RankingCommissionConfig）
   └── 根据陪玩师上月排名，应用阶梯式抽成调整
   └── 排名越高，抽成比例越低（激励机制）

【计算公式】
最终抽成比例 = 基础抽成比例 - 排名抽成减免
陪玩师收入 = 订单金额 × (1 - 最终抽成比例)
平台收入 = 订单金额 × 最终抽成比例

【示例】
订单金额：¥100
项目基础抽成：20%
陪玩师上月排名：第5名（排名抽成减免 5%）
最终抽成：20% - 5% = 15%
陪玩师收入：¥100 × 85% = ¥85
平台收入：¥100 × 15% = ¥15
```

### CommissionRecord（抽成记录）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| PlayerID | uint64 | 陪玩师ID |
| TotalAmountCents | int64 | 订单总金额（分） |
| CommissionRate | int | 抽成比例 |
| CommissionCents | int64 | 平台抽成金额（分） |
| PlayerIncomeCents | int64 | 陪玩师收入（分） |
| SettlementStatus | SettlementStatus | 结算状态：pending/disputed/settled |
| SettlementMonth | string | 结算月份（YYYY-MM） |
| SettledAt | *time.Time | 结算时间 |

### MonthlySettlement（月度结算）

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID |
| SettlementMonth | string | 结算月份（YYYY-MM） |
| TotalOrderCount | int64 | 总订单数 |
| TotalAmountCents | int64 | 总金额（分） |
| TotalCommissionCents | int64 | 总抽成（分） |
| TotalIncomeCents | int64 | 总收入（分） |
| BonusCents | int64 | 奖金（分） |
| FinalIncomeCents | int64 | 最终收入（分） |
| Status | MonthlySettlementStatus | 状态：pending/confirmed/paid |
| IncomeRank | *int | 收入排名 |
| OrderRank | *int | 订单数排名 |
| QualityRank | *int | 质量排名 |

---

## 排名系统

### PlayerRanking（陪玩师排名）

| 字段 | 类型 | 说明 |
|------|------|------|
| PlayerID | uint64 | 陪玩师ID |
| RankingType | RankingType | 排名类型：income/order_count/quality/popularity |
| Period | string | 周期：daily/weekly/monthly/yearly |
| PeriodValue | string | 周期值（YYYY-MM-DD/YYYY-WW/YYYY-MM） |
| Rank | int | 排名 |
| Score | float64 | 排名分数 |
| OrderCount | int64 | 订单数 |
| IncomeCents | int64 | 收入（分） |
| AvgRating | float32 | 平均评分 |
| BonusCents | int64 | 排名奖金（分） |

### RankingReward（排名奖励）

| 字段 | 类型 | 说明 |
|------|------|------|
| RankingType | RankingType | 排名类型 |
| Period | string | 周期 |
| RankStart | int | 排名开始 |
| RankEnd | int | 排名结束 |
| RewardType | string | 奖励类型：fixed/percentage |
| RewardValue | int64 | 奖励值（固定金额分/百分比） |
| Description | string | 描述 |
| IsActive | bool | 是否启用 |

---

## 金额处理规范

所有金额字段使用 **分（Cents）** 为单位存储：

```go
TotalPriceCents int64  // 存储 10000 表示 ¥100.00

// 前端显示时转换
displayPrice := float64(totalPriceCents) / 100
```
