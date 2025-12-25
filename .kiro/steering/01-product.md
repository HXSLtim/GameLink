# 产品概述

> GameLink 游戏陪玩服务平台

## 核心功能

- **智能订单分配**: 用户与陪玩师自动匹配，支持抢单池和客服分配两种模式
- **多角色系统**: 三层架构（用户/陪玩师/管理员）+ RBAC 权限控制
- **实时通信**: 基于 WebSocket 的即时消息，支持公共聊天室和订单群聊（不支持私聊）
- **支付结算**: 完整支付流程，包括订单、退款和收益结算
- **数据监控**: 实时订单状态、收益统计和系统指标

## 用户角色

| 角色 | 说明 | 主要功能 |
|------|------|----------|
| 用户 (User) | 普通玩家 | 浏览陪玩师、创建订单、支付、评价 |
| 陪玩师 (Player) | 服务提供者 | 接单、管理服务、收益追踪、团队功能 |
| 管理员 (Admin) | 平台运营 | 仪表盘、用户管理、订单监控、财务管理、系统设置 |

## 商业模式

- 平台抽成: 15-25%（三维度计算：项目抽成 + 陪玩师个人抽成 + 上月排名调整）
- 系统统一定价: ¥20-60+/小时（按段位，陪玩师不可自定义）
- 收入来源: 服务佣金、认证费用、推广服务

### 抽成计算三维度

| 维度 | 说明 | 优先级 |
|------|------|--------|
| 项目抽成 | 服务项目级别的基础抽成（默认 20%） | 基础 |
| 陪玩师个人抽成 | 特定陪玩师的专属抽成比例 | 覆盖项目抽成 |
| 上月排名抽成 | 根据上月排名的阶梯式抽成减免（激励机制） | 叠加减免 |

> 详细抽成计算规则见 `04-data-models.md` 的"抽成计算规则"章节

## 项目进度

| 指标 | 当前值 | 目标 |
|------|--------|------|
| 测试覆盖率 | ~80% | 80%+ |
| 后端模块完成度 | 74%（26/35 模块） | 100% |
| 前端完成度 | 75% | 100% |
| PRD 功能覆盖 | 82%（45/55 功能点） | 100% |
| CI/CD | ✅ 完善 | - |

### 模块实现状态

| 分类 | 完成 | 进行中 | 仅Model | 说明 |
|------|------|--------|---------|------|
| 核心模块 | 19 | 0 | 0 | user/order/player/chat/dispute 等 |
| 新增业务模块 | 0 | 0 | 3 | player-rank/order-timeout/user-block |
| 营销模块 | 0 | 1 | 5 | vip/coupon/recharge/activity/team/referral |
| 辅助模块 | 7 | 0 | 0 | commission/ranking/routing-rule/operation-log 等 |

> 详细模块状态见 `06-project-management.md`

### 最近完成

- ✅ PRD 功能需求对照完善（82% 功能覆盖率）
- ✅ 全量代码审查完成（所有模块无编译错误）
- ✅ 项目进度更新（74% 完成度，26/35 模块）
- ✅ operation-log 模块确认完成（admin service 包含）
- ✅ 全量代码审查修复（eventHooks/assignment/user dispute handler 字段与模型同步）
- ✅ 争议模块代码修复（disputeService.go 和 dispute handler 字段与模型同步）
- ✅ 模块进度核实（核心模块全部完成，整体完成度提升至 71%）
- ✅ 代码与文档一致性同步（User/OrderDispute/ChatGroup/Commission/OrderItem/OrderPlayer/GameRank 模型更新）
- ✅ 业务文档完善（游戏/服务项目/敏感词/内容管理/统计模块）
- ✅ 用户/认证模块业务流程文档（注册/登录/密码/封禁/Token/RBAC）
- ✅ 项目管理 steering 规则（AI 辅助维护文档）
- ✅ 陪玩师等级/认证 Model（GameRank、PlayerRankRecord、PlayerCertification）
- ✅ 订单超时处理 Model（OrderTimeoutConfig、OrderTimeoutLog、OrderServiceAssignment）
- ✅ 用户拉黑 Model（UserBlock）
- ✅ 统一响应处理（resp 包，96% 覆盖率）
- ✅ 后端代码瘦身（删除 ~4900 行冗余代码）

### 待实现功能

| 优先级 | 模块 | 说明 | 当前状态 |
|--------|------|------|----------|
| P1 | player-rank | 陪玩师等级/认证系统 | 🟡 Model完成，需实现 Repo/Service/Handler |
| P1 | order-timeout | 订单超时处理 | 🟡 Model完成，需实现 Repo/Service/Handler |
| P1 | user-block | 用户拉黑功能 | 🟡 Model完成，需实现 Repo/Service/Handler |
| P2 | vip | VIP会员系统 | 🟡 Model完成 |
| P2 | coupon | 优惠券系统 | 🟡 Model完成 |
| P2 | recharge | 充值系统 | 🟡 Model完成 |
| P2 | activity | 活动系统 | 🟡 Model完成 |
| P2 | team | 团队系统 | 🔄 Service完成，需实现 Repo/Handler |
| P3 | referral | 推荐/邀请系统 | 🟡 Model完成（预留） |

- **通信加密**: AES-256-CBC + SHA-256 签名（生产环境强制启用）
- **认证**: JWT Token + 刷新机制
- **权限控制**: RBAC 角色权限系统
- **数据保护**: 敏感数据加密存储


## 安全特性

- **通信加密**: AES-256-CBC + SHA-256 签名（生产环境强制启用）
- **认证**: JWT Token + 刷新机制
- **权限控制**: RBAC 角色权限系统
- **数据保护**: 敏感数据加密存储
