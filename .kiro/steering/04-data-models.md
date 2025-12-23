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

### RoleModel（角色）

| 字段 | 类型 | 说明 |
|------|------|------|
| Slug | string | 角色标识（唯一）：superAdmin/admin/player/user |
| Name | string | 角色名称 |
| Description | string | 角色描述 |
| IsSystem | bool | 是否系统角色（不可删除） |
| ParentID | *uint64 | 父角色ID（继承） |
| Priority | int | 优先级（权限冲突解决） |
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

---

## 陪玩师模块

### Player（陪玩师）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 关联用户ID（唯一） |
| Nickname | string | 昵称 |
| Bio | string | 个人简介 |
| Rank | string | 段位 |
| RatingAverage | float32 | 平均评分（0-5） |
| RatingCount | uint32 | 评价数量 |
| HourlyRateCents | int64 | 时薪（分） |
| MainGameID | uint64 | 主要游戏ID |
| VerificationStatus | VerificationStatus | 认证状态：pending/verified/rejected |
| VerifiedAt | *time.Time | 审核时间 |
| VerifiedBy | *uint64 | 审核人ID |
| RejectReason | string | 拒绝原因 |

---

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
| Status | string | pending/matched/completed/canceled |
| PlayerID | *uint64 | 接单的陪玩师ID |

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
| Status | string | joined/left/completed |

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

## 游戏模块

### Game（游戏）

| 字段 | 类型 | 说明 |
|------|------|------|
| Key | string | 游戏标识（唯一），如 lol/dota2 |
| Name | string | 游戏名称 |
| Category | string | 分类：moba/fps 等 |
| IconURL | string | 图标URL |

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
| Content | string | 评价内容 |
| Status | ReviewStatus | 状态：pending/approved/rejected/deleted |
| IsReported | bool | 是否被举报 |
| Images | StringArray | 评价图片URL数组 |

---

## 争议模块

### OrderDispute（订单争议）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderID | uint64 | 订单ID |
| UserID | uint64 | 发起用户ID |
| Status | DisputeStatus | 状态（见枚举） |
| Reason | string | 争议原因 |
| EvidenceURLs | EvidenceURLArray | 证据截图URL列表 |
| AssignedToUserID | *uint64 | 指派客服ID |
| SLADeadline | *time.Time | SLA截止时间（默认30分钟） |
| Resolution | DisputeResolution | 处理决定（见枚举） |
| ResolutionAmount | int64 | 退款金额（分） |

---

## 聊天模块

### ChatGroup（聊天群组）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupName | string | 群组名称 |
| GroupType | ChatGroupType | 类型：public/order |
| RelatedOrderID | *uint64 | 关联订单ID |
| CreatedBy | uint64 | 创建者ID |
| MaxMembers | int | 最大成员数（默认100） |
| IsActive | bool | 是否激活 |
| AutoDestroy | bool | 订单完成后自动销毁 |
| VoiceEnabled | bool | 是否启用语音（预留） |
| VoiceRoomID | string | 语音房间ID（第三方服务） |
| VoiceProvider | string | 语音服务商：agora/tencent/zego |

### ChatMessage（聊天消息）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupID | uint64 | 群组ID |
| SenderID | uint64 | 发送者ID |
| Content | string | 消息内容 |
| MessageType | ChatMessageType | 类型：text/image/file/system |
| AuditStatus | ChatMessageAuditStatus | 审核状态 |

---

## 金额处理规范

所有金额字段使用 **分（Cents）** 为单位存储：

```go
TotalPriceCents int64  // 存储 10000 表示 ¥100.00

// 前端显示时转换
displayPrice := float64(totalPriceCents) / 100
```
