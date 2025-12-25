# 数据模型 - 通知模块

> 消息通知系统，支持站内消息、App推送，预留短信/微信/邮件渠道

## 相关文档

- [04-data-models.md](./04-data-models.md) - 核心模块
- [04a-marketing-models.md](./04a-marketing-models.md) - 营销模块
- [04b-team-models.md](./04b-team-models.md) - 团队系统
- [04c-enums-indexes.md](./04c-enums-indexes.md) - 枚举类型和数据库索引

---

## 通知模块

### NotificationTemplate（通知模板）

| 字段 | 类型 | 说明 |
|------|------|------|
| Code | string | 模板编码（唯一） |
| Name | string | 模板名称 |
| Type | NotificationType | 通知类型 |
| Title | string | 标题模板（支持变量） |
| Content | string | 内容模板（支持变量） |
| Channels | string | 支持的渠道（JSON数组） |
| Variables | string | 变量说明（JSON） |
| IsActive | bool | 是否启用 |
| IsSystem | bool | 是否系统模板（不可删除） |
| PushTitle | string | 推送标题 |
| PushContent | string | 推送内容 |
| SMSTemplateID | string | 短信模板ID（预留） |
| WechatTemplateID | string | 微信模板ID（预留） |

### UserNotification（用户通知）

> 注意：与 social.go 中的简单 Notification 共存，此模型用于更复杂的通知场景

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID |
| TemplateID | *uint64 | 模板ID |
| Type | NotificationType | 通知类型 |
| Channel | NotificationChannel | 通知渠道 |
| Title | string | 标题 |
| Content | string | 内容 |
| Status | NotificationStatus | 状态 |
| ReadAt | *time.Time | 已读时间 |
| SentAt | *time.Time | 发送时间 |
| RelatedType | string | 关联类型：order/coupon/activity/vip |
| RelatedID | *uint64 | 关联ID |
| PushID | string | 推送ID（第三方返回） |
| FailureReason | string | 失败原因 |

### UserNotificationSetting（用户通知设置）

| 字段 | 类型 | 说明 |
|------|------|------|
| UserID | uint64 | 用户ID（唯一） |
| OrderStatusEnabled | bool | 订单状态通知（默认true） |
| VipExpireEnabled | bool | VIP到期提醒（默认true） |
| CouponExpireEnabled | bool | 优惠券过期提醒（默认true） |
| ActivityEnabled | bool | 活动提醒（默认true） |
| SystemEnabled | bool | 系统公告（默认true） |
| PromotionEnabled | bool | 营销推广（预留，默认true） |
| ChatEnabled | bool | 聊天消息（预留，默认true） |
| InAppEnabled | bool | 站内消息（默认true） |
| PushEnabled | bool | App推送（默认true） |
| SMSEnabled | bool | 短信（预留，默认false） |
| WechatEnabled | bool | 微信（预留，默认false） |
| EmailEnabled | bool | 邮件（预留，默认false） |
| DoNotDisturbEnabled | bool | 是否启用免打扰 |
| DoNotDisturbStart | string | 免打扰开始时间（HH:MM） |
| DoNotDisturbEnd | string | 免打扰结束时间（HH:MM） |

### NotificationConfig（通知系统配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ConfigKey | string | 配置键（唯一） |
| ConfigValue | string | 配置值 |
| Description | string | 描述 |

**配置键常量：**
- `vip_expire_days` - VIP到期提前提醒天数（JSON数组，如 [7,3,1]）
- `coupon_expire_days` - 优惠券过期提前提醒天数（JSON数组）
- `push_provider` - 推送服务商：jpush/getui/umeng
- `sms_provider` - 短信服务商（预留）

### NotificationSchedule（定时通知任务）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 任务名称 |
| Type | NotificationType | 通知类型 |
| TemplateID | uint64 | 模板ID |
| ScheduleAt | time.Time | 计划发送时间 |
| Status | string | pending/processing/completed/failed |
| TargetType | string | 目标类型：all/vip/specific |
| TargetIDs | string | 目标用户ID列表（JSON数组） |
| TotalCount | int | 总发送数 |
| SentCount | int | 已发送数 |
| FailedCount | int | 失败数 |
| StartedAt | *time.Time | 开始时间 |
| CompletedAt | *time.Time | 完成时间 |

---

## 通知业务流程

### 通知类型

| 类型 | 说明 | 触发时机 |
|------|------|----------|
| order_status | 订单状态变更 | 下单成功、接单、完成、取消等 |
| vip_expire | VIP到期提醒 | 提前7/3/1天（可配置） |
| coupon_expire | 优惠券过期提醒 | 提前7/3/1天（可配置） |
| activity_start | 活动开始提醒 | 活动开始时 |
| activity_end | 活动结束提醒 | 活动结束前 |
| system | 系统公告 | 管理员发布 |
| promotion | 营销推广（预留） | - |
| chat | 聊天消息（预留） | - |

### 通知渠道

| 渠道 | 说明 | 状态 |
|------|------|------|
| in_app | 站内消息（App内通知中心） | ✅ 已实现 |
| push | App推送通知 | ✅ 已实现 |
| sms | 短信 | 🔜 预留 |
| wechat | 微信模板消息 | 🔜 预留 |
| email | 邮件 | 🔜 预留 |

### 发送流程

```
【即时通知】
事件触发 → 检查用户设置 → 检查免打扰 → 发送通知
  ├── 站内消息：直接写入 UserNotification
  └── App推送：调用推送服务商API

【定时通知】
定时任务扫描 → 找到待发送任务 → 批量发送 → 更新统计

【VIP/优惠券到期提醒】
每日定时任务 → 查询即将到期的VIP/优惠券
  → 检查是否已发送过该天数的提醒
  → 发送通知
```

### 用户设置检查

```go
// 检查用户是否接收某类通知
func ShouldSendNotification(setting *UserNotificationSetting, notifType NotificationType, channel NotificationChannel) bool {
    // 1. 检查通知类型开关
    switch notifType {
    case NotificationTypeOrderStatus:
        if !setting.OrderStatusEnabled { return false }
    case NotificationTypeVipExpire:
        if !setting.VipExpireEnabled { return false }
    // ...
    }
    
    // 2. 检查渠道开关
    switch channel {
    case NotificationChannelInApp:
        if !setting.InAppEnabled { return false }
    case NotificationChannelPush:
        if !setting.PushEnabled { return false }
    // ...
    }
    
    // 3. 检查免打扰时段（仅对推送生效）
    if channel == NotificationChannelPush && setting.IsInDoNotDisturbPeriod() {
        return false
    }
    
    return true
}
```

### 模板变量

模板支持变量替换，格式：`{{.变量名}}`

```
【订单状态变更】
标题：您的订单 {{.OrderNo}} 状态已更新
内容：订单状态已变更为「{{.Status}}」，点击查看详情

【VIP到期提醒】
标题：VIP即将到期提醒
内容：您的 {{.VipLevel}} 会员将于 {{.ExpireDate}} 到期，续费可享专属优惠

【优惠券过期提醒】
标题：优惠券即将过期
内容：您有 {{.Count}} 张优惠券将于 {{.ExpireDate}} 过期，快去使用吧
```

---

## 与现有 Notification 的关系

`social.go` 中已有简单的 `Notification` 模型（表名 `notifications`），用于基础通知场景。

新的 `UserNotification` 模型（表名 `user_notifications`）提供更完整的功能：
- 支持多渠道
- 支持模板
- 支持关联数据
- 支持发送状态追踪

两者共存，可根据场景选择使用。
