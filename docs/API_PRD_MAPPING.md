# API 与 PRD 功能点映射表

> **文档版本**: v1.0
> **创建日期**: 2025-01-05
> **用途**: 追踪 API 端点与 PRD 功能需求的对应关系

---

## 概述

本文档将 `docs/PRD.md` 第4章的功能需求映射到实际的 API 端点，确保每个功能点都有对应的 API 支持。

**映射标记**:
- ✅ 已实现并验证
- ⚠️ 已实现但需测试
- ❌ 未实现

---

## 用户端 API (User端)

### 1. 首页浏览功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 游戏列表 | `/api/v1/user/games` | GET | GameService | ✅ | 获取所有游戏 |
| 陪玩师列表 | `/api/v1/user/players` | GET | PlayerService | ✅ | 支持分页和筛选 |
| 搜索功能 | `/api/v1/user/players/search` | GET | PlayerService | ✅ | 按昵称/ID搜索 |
| 排行榜 | `/api/v1/user/rankings` | GET | RankingService | ✅ | 人气/订单/好评榜 |

### 2. 订单功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 创建订单 | `/api/v1/user/orders` | POST | OrderService | ✅ | solo/team/gift订单 |
| 支付 | `/api/v1/user/payments` | POST | PaymentService | ✅ | 支持组合支付 |
| 订单列表 | `/api/v1/user/orders` | GET | OrderService | ✅ | 我的订单，支持状态筛选 |
| 取消订单 | `/api/v1/user/orders/:id/cancel` | POST | OrderService | ✅ | 仅pending状态可取消 |
| 退款申请 | `/api/v1/user/orders/:id/refund` | POST | RefundService | ✅ | 订单完成7天内可申请 |

### 3. 评价功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 提交评价 | `/api/v1/user/reviews` | POST | ReviewService | ✅ | 订单完成后可评价 |
| 一键评价 | `/api/v1/user/reviews/quick` | POST | ReviewService | ✅ | 快速5星好评 |
| 单独评价 | `/api/v1/user/reviews/:id` | PUT | ReviewService | ✅ | 分维度评分(态度/技术/沟通) |
| 修改评价 | `/api/v1/user/reviews/:id` | PATCH | ReviewService | ✅ | 仅限评价者本人 |

### 4. 钱包功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 余额查询 | `/api/v1/user/wallet` | GET | WalletService | ✅ | 查询可用和冻结余额 |
| 充值 | `/api/v1/user/recharges` | POST | RechargeService | ✅ | 第三方支付充值 |
| 交易记录 | `/api/v1/user/wallet/transactions` | GET | WalletService | ✅ | 分页查询交易记录 |

### 5. 聊天功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 公共聊天室 | `/api/v1/user/chats/public/:gameId` | GET | ChatService | ✅ | 获取公共聊天室消息 |
| 发送消息 | `/api/v1/user/chats/:roomId/messages` | POST | ChatService | ✅ | WebSocket或HTTP |
| 订单群聊 | `/api/v1/user/chats/order/:orderId` | GET | ChatService | ✅ | 订单聊天室信息 |

### 6. 营销功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| VIP状态 | `/api/v1/user/vip` | GET | VipService | ✅ | 查询VIP等级和权益 |
| 我的优惠券 | `/api/v1/user/coupons` | GET | CouponService | ✅ | 可用优惠券列表 |
| 领取优惠券 | `/api/v1/user/coupons/:id/claim` | POST | CouponService | ✅ | 领取优惠券 |
| 活动列表 | `/api/v1/user/activities` | GET | ActivityService | ✅ | 可参与活动 |
| 参与活动 | `/api/v1/user/activities/:id/participate` | POST | ActivityService | ✅ | 参与活动领取奖励 |
| 推荐信息 | `/api/v1/user/referral` | GET | ReferralService | ✅ | 我的推荐码和收益 |

### 7. 其他功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 用户拉黑 | `/api/v1/user/blocks` | POST | UserBlockService | ✅ | 拉黑陪玩师 |
| 取消拉黑 | `/api/v1/user/blocks/:id` | DELETE | UserBlockService | ✅ | 取消拉黑 |
| 动态列表 | `/api/v1/user/feeds` | GET | FeedService | ✅ | 平台动态 |

---

## 陪玩师端 API (Player端)

### 1. 入驻认证功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 入驻申请 | `/api/v1/player/certifications` | POST | PlayerCertificationService | ✅ | 提交入驻申请 |
| 实名认证 | `/api/v1/player/certifications/identity` | POST | PlayerCertificationService | ✅ | 提交实名信息 |
| 段位认证 | `/api/v1/player/certifications/rank` | POST | PlayerCertificationService | ✅ | 上传段位截图 |

### 2. 订单管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 待接订单 | `/api/v1/player/orders/pending` | GET | OrderService | ✅ | 抢单池订单 |
| 接单 | `/api/v1/player/orders/:id/accept` | POST | OrderService | ✅ | 接受订单 |
| 我的订单 | `/api/v1/player/orders` | GET | OrderService | ✅ | 陪玩师订单列表 |
| 完成订单 | `/api/v1/player/orders/:id/complete` | POST | OrderService | ✅ | 标记服务完成 |

### 3. 收益管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 收益统计 | `/api/v1/player/earnings/stats` | GET | StatisticsService | ✅ | 收益概览和趋势 |
| 提现申请 | `/api/v1/player/withdrawals` | POST | WithdrawService | ✅ | 申请提现 |
| 提现记录 | `/api/v1/player/withdrawals` | GET | WithdrawService | ✅ | 提现历史记录 |

### 4. 状态管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 在线状态 | `/api/v1/player/status/online` | POST | PlayerService | ✅ | 设置在线状态 |
| 离线状态 | `/api/v1/player/status/offline` | POST | PlayerService | ✅ | 设置离线状态 |
| 接单开关 | `/api/v1/player/status/toggle-accepting` | POST | PlayerService | ✅ | 开启/关闭接单 |

### 5. 团队功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 创建团队 | `/api/v1/player/teams` | POST | TeamService | ✅ | 创建团队 |
| 团队信息 | `/api/v1/player/teams/:id` | GET | TeamService | ✅ | 查询团队详情 |
| 团队接单 | `/api/v1/player/teams/:id/orders` | POST | TeamService | ✅ | 团队接单 |
| 收入分配 | `/api/v1/player/teams/:id/distribute` | POST | TeamService | ✅ | 分配收入给成员 |

---

## 管理后台 API (Admin端)

### 1. 仪表盘功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 数据概览 | `/api/v1/admin/dashboard/overview` | GET | DashboardService | ✅ | 核心指标统计 |
| 趋势图表 | `/api/v1/admin/dashboard/trends` | GET | StatisticsService | ✅ | 趋势数据 |
| 实时监控 | `/api/v1/admin/dashboard/realtime` | GET | MonitorService | ✅ | WebSocket实时推送 |

### 2. 用户管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 用户列表 | `/api/v1/admin/users` | GET | UserService | ✅ | 支持分页和筛选 |
| 用户详情 | `/api/v1/admin/users/:id` | GET | UserService | ✅ | 用户完整信息 |
| 封禁/解封 | `/api/v1/admin/users/:id/status` | PUT | UserService | ✅ | 修改用户状态 |
| 用户标签 | `/api/v1/admin/users/:id/tags` | POST | UserService | ✅ | 添加用户标签 |
| 用户拉黑 | `/api/v1/admin/user-blocks` | GET | UserBlockService | ✅ | 查询拉黑记录 |

### 3. 陪玩师管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 陪玩师列表 | `/api/v1/admin/players` | GET | PlayerService | ✅ | 支持分页和筛选 |
| 入驻审核 | `/api/v1/admin/player-certifications/:id/approve` | POST | PlayerCertificationService | ✅ | 审核入驻申请 |
| 段位审核 | `/api/v1/admin/player-rank-records/:id/approve` | POST | PlayerRankService | ✅ | 审核段位认证 |
| 实名审核 | `/api/v1/admin/player-certifications/identity/:id/approve` | POST | PlayerCertificationService | ✅ | 审核实名认证 |

### 4. 订单管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 订单列表 | `/api/v1/admin/orders` | GET | OrderService | ✅ | 支持分页和筛选 |
| 订单详情 | `/api/v1/admin/orders/:id` | GET | OrderService | ✅ | 订单完整信息 |
| 退款审核 | `/api/v1/admin/refunds/:id/approve` | POST | RefundService | ✅ | 审核退款申请 |
| 争议处理 | `/api/v1/admin/disputes/:id` | GET/PUT | DisputeService | ✅ | 查询/处理争议 |
| 订单超时 | `/api/v1/admin/order-timeout-logs` | GET | OrderTimeoutService | ✅ | 超时日志查询 |

### 5. 财务管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 提现审核 | `/api/v1/admin/withdrawals/:id/approve` | POST | WithdrawService | ✅ | 审核提现申请 |
| 佣金配置 | `/api/v1/admin/commission-rules` | CRUD | CommissionService | ✅ | 佣金规则配置 |
| 结算管理 | `/api/v1/admin/settlement-companies` | CRUD | SettlementCompanyService | ✅ | 结算主体管理 |

### 6. 内容管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 游戏管理 | `/api/v1/admin/games` | CRUD | GameService | ✅ | 游戏CRUD |
| 段位管理 | `/api/v1/admin/game-ranks` | CRUD | GameRankService | ✅ | 段位CRUD |
| 服务项目 | `/api/v1/admin/service-items` | CRUD | ServiceItemService | ✅ | 服务项配置 |
| 敏感词 | `/api/v1/admin/sensitive-words` | CRUD | SensitiveWordService | ✅ | 敏感词管理 |

### 7. 系统管理功能

| PRD 功能点 | API 端点 | 方法 | Handler | 状态 | 说明 |
|-----------|---------|------|---------|------|------|
| 角色管理 | `/api/v1/admin/roles` | CRUD | RoleService | ✅ | 角色CRUD |
| 权限管理 | `/api/v1/admin/permissions` | CRUD | PermissionService | ✅ | 权限树管理 |
| 菜单管理 | `/api/v1/admin/menus` | CRUD | MenuService | ✅ | 菜单配置 |
| 操作日志 | `/api/v1/admin/operation-logs` | GET | AuditService | ✅ | 操作日志查询 |

---

## API 完整性统计

### 按模块统计

| 模块 | API 数量 | 覆盖功能点 | 完成率 |
|------|---------|-----------|--------|
| 用户端核心功能 | 20+ | 17 | 100% |
| 陪玩师端核心功能 | 15+ | 14 | 100% |
| 管理后台功能 | 30+ | 24 | 100% |
| **总计** | **65+** | **55** | **100%** ✅ |

### 待补充 API (前端缺失功能)

| 功能模块 | 缺失 API | 状态 | 优先级 |
|---------|---------|------|--------|
| 陪玩师端 - 实名/段位认证审核页面 | 需补充查询和审核UI | ⚠️ | P1 |
| 管理后台 - 段位管理页面 | 已有API，缺前端页面 | ⚠️ | P1 |
| 陪玩师端 - 团队管理页面 | 已有API，缺前端页面 | ⚠️ | P2 |

---

## 维护说明

### 更新频率
- 每次 API 变更后更新此文档
- 每月与 PRD.md 功能点进行对账

### 责任人
- 后端开发：确保 API 实现与文档一致
- 产品经理：确保 PRD 功能点有对应 API
- 测试：验证 API 与 PRD 功能点一致性

---

*文档维护：GameLink 开发团队*
*最后更新：2025-01-05*
