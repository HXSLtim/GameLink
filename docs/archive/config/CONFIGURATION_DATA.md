# GameLink 系统配置数据汇总

> 本文档汇总了 GameLink 平台各模块的可量化配置数据，用于系统配置、业务规则调整和运维参考。

**文档版本**: v1.0
**最后更新**: 2025-01-10
**适用范围**: 全系统 36 个模块

---

## 目录

1. [认证与权限](#认证与权限)
2. [订单与支付](#订单与支付)
3. [佣金与结算](#佣金与结算)
4. [争议处理](#争议处理)
5. [营销系统](#营销系统)
6. [用户与陪玩师](#用户与陪玩师)
7. [系统限制](#系统限制)
8. [缓存与性能](#缓存与性能)

---

## 认证与权限

### JWT Token 配置

| 配置项 | 常量名 | 默认值 | 说明 |
|--------|--------|--------|------|
| Token 有效期 | `DefaultTokenDuration` | 24 小时 | JWT Token 默认有效期 |
| 自动刷新窗口 | `TokenAutoRefreshWindow` | 15 分钟 | 剩余时间小于此值时自动刷新 |
| 刷新建议窗口 | `TokenRefreshRecommendationWindow` | 1 小时 | 提醒前端刷新的时间窗口 |
| 最小刷新阈值 | `TokenMinRefreshThreshold` | 30 秒 | Token 刷新的最小阈值 |
| 最大刷新窗口 | `DefaultMaxRefreshWindow` | 7 天 | Token 签发后超过此时间无法刷新 |

**代码位置**: [`api/pkg/auth/constants.go`](../api/pkg/auth/constants.go)

### 角色继承配置

| 配置项 | 常量名 | 默认值 | 说明 |
|--------|--------|--------|------|
| 最大继承深度 | `MaxRoleInheritanceDepth` | 5 层 | 角色继承的最大层级限制 |

**代码位置**: [`api/internal/model/role.go`](../api/internal/model/role.go)

### 预定义角色

| 角色标识 | 角色名称 | 说明 |
|----------|----------|------|
| `superAdmin` | 超级管理员 | 拥有所有权限 |
| `admin` | 管理员/店长 | 平台管理员 |
| `finance` | 财务 | 财务操作权限 |
| `customerService` | 客服 | 客服操作权限 |
| `player` | 陪玩师 | 接单角色 |
| `user` | 普通用户 | 下单用户 |

---

## 订单与支付

### 订单超时配置

| 配置键 | 默认值 | 单位 | 说明 |
|--------|--------|------|------|
| `payment_timeout_minutes` | 30 | 分钟 | 支付超时时间 |
| `order_accept_timeout_minutes` | 30 | 分钟 | 接单超时时间 |
| `auto_cancel_enabled` | true | - | 是否启用自动取消 |
| `auto_refund_enabled` | true | - | 是否启用自动退款 |
| `auto_assign_service_enabled` | true | - | 接单后是否自动分配客服 |

**代码位置**: [`api/internal/model/orderTimeout.go`](../api/internal/model/orderTimeout.go)

### 订单金额配置

| 配置项 | 范围 | 单位 | 说明 |
|--------|------|------|------|
| 最低订单金额 | 20 | 元 | 系统控制的最低服务价格 |
| 最高订单金额 | 60+ | 元 | 根据陪玩师等级动态调整 |
| 订单金额精度 | 0.01 | 元 | 金额精确到分 |

### 陪玩师等级价格

| 等级 | 每小时价格范围 |
|------|---------------|
| 新手陪玩师 | ¥20-30/小时 |
| 中级陪玩师 | ¥30-45/小时 |
| 高级陪玩师 | ¥45-60/小时 |
| 顶级陪玩师 | ¥60+/小时 |

---

## 佣金与结算

### 默认佣金规则

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 平台抽成比例 | 20% | 默认平台佣金 |
| 陪玩师所得比例 | 80% | 陪玩师实际收入比例 |

**代码位置**: [`api/pkg/db/migrate.go`](../api/pkg/db/migrate.go#L559-L572)

### 佣金计算三层优先级

1. **第一优先级**: 陪玩师个人抽成比例 (`CommissionRule.PlayerID`)
2. **第二优先级**: 服务项目抽成比例 (`ServiceItem.CommissionRate`)
3. **第三优先级**: 月度排行榜折扣 (`RankingCommissionConfig`)

### 月度排行榜佣金折扣

| 排名区间 | 佣金折扣 | 实际抽成比例 |
|----------|----------|--------------|
| #1-3 | 5% | 15% |
| #4-10 | 3% | 17% |
| #11-50 | 1% | 19% |
| #51+ | 0% | 20% |

### 收入结算规则

| 配置项 | 数值 | 说明 |
|--------|------|------|
| T+7 冻结期 | 7 天 | 订单完成后收入冻结时间 |
| 争议期 | 7 天 | 订单完成后可发起争议的时间 |
| 最低提现金额 | - | 待配置 |

---

## 争议处理

### 争议 SLA

| 配置项 | 时间限制 | 说明 |
|--------|----------|------|
| 客服响应 SLA | 30 分钟 | 客服必须响应的时间 |
| 争议发起期限 | 订单完成后 7 天 | 用户/陪玩师可发起争议的时间窗口 |

### 争议类型

| 类型标识 | 名称 | 说明 |
|----------|------|------|
| `service_quality` | 服务质量问题 | 服务质量不达标 |
| `bad_attitude` | 态度问题 | 陪玩师态度恶劣 |
| `incomplete_service` | 未完成服务 | 服务未按约定完成 |
| `late_arrival` | 迟到 | 陪玩师迟到 |
| `early_leave` | 早退 | 陪玩师提前结束 |

### 争议处理结果

| 结果类型 | 说明 |
|----------|------|
| `refund` | 全额退款 |
| `partial` | 部分退款 |
| `reassign` | 重新指派陪玩师 |
| `reject` | 驳回争议 |

**代码位置**: [`api/internal/model/dispute.go`](../api/internal/model/dispute.go)

---

## 营销系统

### VIP 等级配置

| 等级 | 订单折扣 | 说明 |
|------|----------|------|
| 普通 | 1.0 (无折扣) | 无折扣 |
| VIP1 | 0.98 (98折) | 2% 折扣 |
| VIP2 | 0.95 (95折) | 5% 折扣 |
| VIP3 | 0.90 (90折) | 10% 折扣 |

**代码位置**: [`api/internal/model/vip.go`](../api/internal/model/vip.go)

### 优惠券配置

#### 优惠券类型

| 类型 | 标识 | 说明 |
|------|------|------|
| 满减券 | `deduct` | 直接减免金额 |
| 折扣券 | `discount` | 按比例打折 |

#### 优惠券来源

| 来源 | 标识 | 说明 |
|------|------|------|
| 系统发放 | `system` | 平台主动发放 |
| 活动领取 | `activity` | 用户参与活动获得 |
| 充值赠送 | `recharge` | 充值奖励 |
| 推荐奖励 | `referral` | 推荐新用户奖励 |

#### 优惠券范围

| 范围 | 标识 | 说明 |
|------|------|------|
| 全部订单 | `all` | 所有订单可用 |
| 指定游戏 | `game` | 特定游戏订单可用 |
| 指定服务 | `service` | 特定服务项目可用 |

**代码位置**: [`api/internal/model/coupon.go`](../api/internal/model/coupon.go)

### 充值档位配置

充值档位是动态配置的，典型配置示例：

| 档位 | 充值金额 | 赠送金额 | 折扣比例 |
|------|----------|----------|----------|
| 档位1 | ¥30 | ¥0 | 无折扣 |
| 档位2 | ¥68 | ¥12 | 约 8.3 折 |
| 档位3 | ¥128 | ¥32 | 约 7.8 折 |
| 档位4 | ¥328 | ¥92 | 约 7.2 折 |
| 档位5 | ¥648 | ¥200 | 约 6.9 折 |

**数据模型**:
- `AmountCents`: 充值金额（分）
- `BonusCents`: 赠送金额（分）
- `OriginalCents`: 原价（用于显示划线价）
- `DiscountPercent`: 折扣百分比（显示用）

**代码位置**: [`api/internal/model/recharge.go`](../api/internal/model/recharge.go)

### 推荐奖励配置

| 奖励类型 | 标识 | 说明 |
|----------|------|------|
| 现金 | `cash` | 直接到钱包 |
| 优惠券 | `coupon` | 发放优惠券 |
| 积分 | `points` | 预留功能 |

**代码位置**: [`api/internal/model/referral.go`](../api/internal/model/referral.go)

### 活动类型

| 活动类型 | 标识 | 说明 |
|----------|------|------|
| 优惠券发放 | `coupon` | 发放优惠券活动 |
| 限时折扣 | `discount` | 预留功能 |
| 赠品活动 | `gift` | 预留功能 |

**代码位置**: [`api/internal/model/activity.go`](../api/internal/model/activity.go)

---

## 用户与陪玩师

### 陪玩师等级配置

| 等级 | 说明 |
|------|------|
| Bronze | 青铜陪玩师 |
| Silver | 白银陪玩师 |
| Gold | 黄金陪玩师 |
| Platinum | 铂金陪玩师 |
| Diamond | 钻石陪玩师 |
| Master | 大师陪玩师 |

### 订单类型

| 类型 | 标识 | 所需陪玩师数 | 支付流程 |
|------|------|-------------|----------|
| 单人陪玩 | `solo` | 1 | 标准流程 |
| 多人组队 | `team` | 2+ | 需匹配所有位置后开始 |
| 礼物订单 | `gift` | 1 | 即时完成，无 T+7 |

### 订单状态

| 状态 | 标识 | 说明 |
|------|------|------|
| 待支付 | `pending` | 等待用户支付 |
| 待接单 | `waiting` | 等待陪玩师接单 |
| 服务中 | `in_progress` | 陪玩中 |
| 待确认 | `pending_confirmation` | 等待用户确认完成 |
| 已完成 | `completed` | 订单完成 |
| 已取消 | `canceled` | 订单取消 |
| 争议中 | `disputed` | 争议处理中 |

---

## 系统限制

### 批量操作限制

| 操作类型 | 最大数量 | 说明 |
|----------|----------|------|
| 批量更新佣金 | 100 个服务项 | 单次最多更新 100 个 |
| 批量删除 | 可变 | 根据具体模块限制 |

### 数据库字段长度限制

| 字段类型 | 最大长度 | 说明 |
|----------|----------|------|
| 角色名称 | 128 字符 | `RoleModel.Name` |
| 角色描述 | 255 字符 | `RoleModel.Description` |
| 角色标识 | 64 字符 | `RoleModel.Slug` |
| 优惠券名称 | 128 字符 | `Coupon.Name` |
| 配置值 | 255 字符 | `OrderTimeoutConfig.ConfigValue` |
| 配置描述 | 255 字符 | `OrderTimeoutConfig.Description` |

---

## 缓存与性能

### 缓存配置

| 缓存类型 | TTL | 说明 |
|----------|-----|------|
| 服务项目缓存 | 1 小时 | 热点数据缓存 |
| VIP 配置缓存 | 1 小时 | 热点数据缓存 |
| 佣金规则缓存 | 1 小时 | 热点数据缓存 |
| 权限缓存 | 5 分钟 | 用户权限缓存 |

### WebSocket 配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 心跳间隔 | - | 待配置 |
| 连接超时 | - | 待配置 |
| 消息队列大小 | - | 待配置 |

---

## 定时任务

### 结算任务

| 配置项 | 时间 | 说明 |
|--------|------|------|
| 每月结算时间 | 每月 1 日 02:00 | 自动执行月度结算 |
| 分布式锁 | 启用 | 防止多实例重复执行 |

### 聊天记录保留

| 配置项 | 时间 | 说明 |
|--------|------|------|
| 执行时间 | 每天 03:15 | 清理过期聊天记录 |
| 分布式锁 | 启用 | 防止多实例重复执行 |

**代码位置**: [`api/pkg/scheduler/`](../api/pkg/scheduler/)

---

## 通知系统

### 通知类型

| 类型 | 标识 | 说明 |
|------|------|------|
| VIP 到期 | `vip_expire` | VIP 会员到期提醒 |
| 优惠券过期 | `coupon_expire` | 优惠券即将过期提醒 |
| 活动开始 | `activity_start` | 活动开始提醒 |

**代码位置**: [`api/internal/model/notification.go`](../api/internal/model/notification.go)

---

## 监控指标

### 系统状态监控

| 指标 | 说明 |
|------|------|
| CPU 使用率 | 系统整体 CPU 使用率 |
| 内存使用率 | 系统内存使用率 |
| Goroutines 数量 | Go 协程数量 |
| 数据库连接数 | 活跃/空闲/最大连接数 |
| 在线用户数 | 当前在线用户统计 |
| 每秒请求数 | 系统吞吐量 |
| 系统运行时间 | 服务启动时长 |

### 健康状态

| 状态 | 触发条件 |
|------|----------|
| `healthy` | 所有指标正常 |
| `degraded` | CPU 或 内存 > 70%，或 DB 连接 > 60% |
| `critical` | CPU 或 内存 > 90%，或 DB 连接 > 80% |

**代码位置**: [`api/internal/service/monitor/realtime.go`](../api/internal/service/monitor/realtime.go)

---

## API 配置

### 路由配置

| 路由前缀 | 说明 |
|----------|------|
| `/api/v1/admin` | 管理后台 API |
| `/api/v1/user` | 用户端 API |
| `/api/v1/player` | 陪玩师端 API |
| `/api/v1/public` | 公开 API |
| `/metrics` | Prometheus 监控指标 |
| `/healthz` | 健康检查端点 |

### 端点统计

| 类型 | 数量 |
|------|------|
| 总 API 端点数 | 435 |
| 管理端 API | ~300 |
| 用户端 API | ~80 |
| 陪玩师端 API | ~50 |

---

## 数据库索引

### 外键索引

已为以下表添加外键索引（共 12 个）:

1. `orders.user_id`
2. `orders.player_id`
3. `orders.game_id`
4. `payments.order_id`
5. `reviews.order_id`
6. `chat_groups.order_id`
7. `disputes.order_id`
8. `user_blocks.blocked_user_id`
9. `coupons.user_id`
10. `recharge_records.user_id`
11. `order_service_assignments.order_id`
12. `order_service_assignments.service_user_id`

---

## 配置管理建议

### 可配置化优先级

| 优先级 | 配置项 | 当前状态 | 建议 |
|--------|--------|----------|------|
| P0 | JWT Token 有效期 | 硬编码 | 建议移至配置文件 |
| P0 | 订单超时时间 | 数据库配置 | ✅ 已可配置 |
| P0 | 佣金比例 | 数据库配置 | ✅ 已可配置 |
| P1 | VIP 等级折扣 | 数据库配置 | ✅ 已可配置 |
| P1 | 充值档位 | 数据库配置 | ✅ 已可配置 |
| P2 | 批量操作限制 | 硬编码 | 建议移至配置文件 |
| P2 | 缓存 TTL | 硬编码 | 建议移至配置文件 |

### 配置文件位置

- **开发环境**: `api/configs/config.development.yaml`
- **生产环境**: `api/configs/config.production.yaml`
- **默认配置**: `api/pkg/config/defaults.go`

---

## 变更日志

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2025-01-10 | v1.0 | 初始版本，汇总 36 个模块的配置数据 |

---

## 附录

### 相关文档

- [数据模型定义](../.kiro/steering/04-data-models.md)
- [项目结构说明](../.kiro/steering/03-project-structure.md)
- [测试标准](../.kiro/steering/05-testing-standard.md)
- [项目管理规则](../.kiro/06-project-management.md)

### 代码位置索引

- **认证常量**: `api/pkg/auth/constants.go`
- **角色模型**: `api/internal/model/role.go`
- **订单超时**: `api/internal/model/orderTimeout.go`
- **争议模型**: `api/internal/model/dispute.go`
- **优惠券**: `api/internal/model/coupon.go`
- **VIP**: `api/internal/model/vip.go`
- **充值**: `api/internal/model/recharge.go`
- **推荐**: `api/internal/model/referral.go`
- **活动**: `api/internal/model/activity.go`
- **通知**: `api/internal/model/notification.go`
- **监控服务**: `api/internal/service/monitor/realtime.go`
- **定时任务**: `api/pkg/scheduler/`

---

**维护说明**: 本文档应随着系统配置的变更及时更新。所有新增配置或修改现有配置都应在变更日志中记录。
