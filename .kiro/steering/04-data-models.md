# 数据模型

> ⚠️ **文档维护提醒**：修改 `backend/internal/model/` 后必须同步更新此文档。

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
}
```

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

### Order（订单）

| 字段 | 类型 | 说明 |
|------|------|------|
| OrderNo | string | 订单号（唯一） |
| UserID | uint64 | 下单用户ID |
| ItemID | uint64 | 服务项目ID |
| PlayerID | *uint64 | 服务陪玩师ID |
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

### ChatMessage（聊天消息）

| 字段 | 类型 | 说明 |
|------|------|------|
| GroupID | uint64 | 群组ID |
| SenderID | uint64 | 发送者ID |
| Content | string | 消息内容 |
| MessageType | ChatMessageType | 类型：text/image/file/system |
| AuditStatus | ChatMessageAuditStatus | 审核状态 |

---

## 枚举类型

### 用户状态 (UserStatus)
`active` | `suspended` | `banned`

### 订单状态 (OrderStatus)
`pending` | `confirmed` | `in_progress` | `completed` | `canceled` | `refunded`

### 支付状态 (PaymentStatus)
`pending` | `paid` | `failed` | `refunded`

### 支付方式 (PaymentMethod)
`wechat` | `alipay` | `wallet` | `combined`

### 认证状态 (VerificationStatus)
`pending` | `verified` | `rejected`

### 争议状态 (DisputeStatus)
`pending` | `assigned` | `mediating` | `resolved` | `rejected` | `canceled`

### 货币类型 (Currency)
`CNY` | `USD` | `EUR`

---

## 金额处理规范

所有金额字段使用 **分（Cents）** 为单位存储：

```go
TotalPriceCents int64  // 存储 10000 表示 ¥100.00

// 前端显示时转换
displayPrice := float64(totalPriceCents) / 100
```

---

## 数据库索引

### 唯一索引
- User: Phone, Email
- Order: OrderNo
- Player: UserID
- Game: Key
- RoleModel: Slug
- Permission: Method+Path, Code
- ServiceItem: ItemCode

### 复合索引
- Order: (UserID, Status, CreatedAt)
- Order: (PlayerID, Status)
- User: (Status, LastLoginAt)

---

## 变更日志

| 日期 | 变更内容 |
|------|----------|
| 2024-12-21 | 初始化数据模型文档 |
