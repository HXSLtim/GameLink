# GameLink 项目概览

> **版本**: v2.0
> **更新日期**: 2026-02-09
> **项目类型**: 游戏陪玩社交平台
> **技术栈**: Go + React + shadcn/ui + Tailwind CSS + PostgreSQL + Redis

---

## 📋 目录

1. [项目简介](#1-项目简介)
2. [技术架构总览](#2-技术架构总览)
3. [数据库设计](#3-数据库设计)
4. [功能模块](#4-功能模块)
5. [开发和部署指南](#5-开发和部署指南)
6. [团队和联系方式](#6-团队和联系方式)

---

## 1. 项目简介

### 1.1 产品定位

**GameLink** 是一款现代化游戏陪玩社交平台，通过技术手段连接有陪玩需求的用户与提供专业陪玩服务的陪玩师，打造从**发现、下单、匹配、服务、评价到结算**的完整商业闭环。

### 1.2 核心价值

| 用户角色 | 核心价值 |
|---------|---------|
| **用户** | 快速找到优质陪玩师、安全支付、服务保障 |
| **陪玩师** | 接单赚钱、管理服务、查看收益、技能变现 |
| **平台** | 交易佣金、VIP会员、增值服务 |

### 1.3 商业模式

- **平台佣金**：15%-25% 交易抽成
- **VIP会员**：等级权益、永久折扣、月度券
- **增值服务**：充值赠送、礼物打赏
- **数据服务**：行业报告、API开放（未来）

### 1.4 项目规模

```
┌──────────────────────────────────────────────────────┐
│                   项目规模统计                         │
├──────────────────────────────────────────────────────┤
│  代码行数    ~150,000+ 行                             │
│  数据表      80+ 张                                   │
│  API 接口    200+ 个                                  │
│  页面        40+ 个（管理端）                         │
│              用户端 Web（React）                        │
│  组件        100+ 个                                  │
│  测试用例    250+ 个                                  │
└──────────────────────────────────────────────────────┘
```

---

## 2. 技术架构总览

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      客户端层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ 用户端 Web   │  │ 管理后台     │  │ 未来: APP       │   │
│  │ React 19     │  │ React 19     │  │ 待评估技术栈    │   │
│  │ shadcn + TW  │  │ AntD 6 + TS │  │                 │   │
│  └──────┬───────┘  └──────┬───────┘  └─────────────────┘   │
│         │                 │                                  │
│         └─────────────────┴────────────────┐                │
│                   │ HTTPS / WSS           │                │
├───────────────────┼────────────────────────┼────────────────┤
│          Nginx 反向代理 + 负载均衡         │                │
├───────────────────┼────────────────────────┼────────────────┤
│                   │                        │                │
│  ┌────────────────┴────────────────────────┴────────┐     │
│  │                  API 层 (Go 1.24)                 │     │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │     │
│  │  │ Handler    │→ │ Service    │→ │ Repository │  │     │
│  │  │ (路由)     │  │ (业务逻辑) │  │ (数据访问) │  │     │
│  │  └────────────┘  └────────────┘  └────────────┘  │     │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │     │
│  │  │ WebSocket  │  │ Scheduler  │  │ Middleware │  │     │
│  │  │ (实时通讯) │  │ (定时任务) │  │ (中间件)   │  │     │
│  │  └────────────┘  └────────────┘  └────────────┘  │     │
│  └───────────────────────────────────────────────────┘     │
│                         │                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    数据层                             │   │
│  │  ┌──────────────┐  ┌──────────────┐                │   │
│  │  │ PostgreSQL   │  │ Redis        │                │   │
│  │  │ 主存储       │  │ 缓存/会话    │                │   │
│  │  │ 16+ 版本     │  │ 7+ 版本      │                │   │
│  │  └──────────────┘  └──────────────┘                │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   外部服务                           │   │
│  │  微信支付 / 支付宝 / 腾讯云TRTC / 对象存储          │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 后端技术栈

| 层级 | 技术 | 版本 | 说明 |
|-----|------|------|------|
| **语言** | Go | 1.24+ | 高性能、并发友好 |
| **框架** | Gin | 最新 | HTTP 路由框架 |
| **ORM** | GORM | 最新 | 数据库 ORM |
| **数据库** | PostgreSQL | 16+ | 主数据存储 |
| **缓存** | Redis | 7+ | 缓存、会话、PubSub |
| **WebSocket** | Gorilla WebSocket | - | 实时通讯 |
| **认证** | JWT + HMAC-SHA256 | - | 身份验证 + 请求签名 |
| **加密** | AES-256-CBC | - | 请求体加密（生产环境） |
| **文档** | Swagger | - | API 文档自动生成 |

### 2.3 前端技术栈

#### 管理后台

| 技术 | 版本 | 说明 |
|-----|------|------|
| React | 19 | UI 框架 |
| TypeScript | 5.9 | 类型系统 |
| Ant Design | 6 | UI 组件库 |
| React Router | 6 | 路由管理 |
| Vite | 5 | 构建工具 |
| TanStack Query | - | 数据请求 |

#### 用户端 Web

| 技术 | 版本 | 说明 |
|-----|------|------|
| React | 19 | UI 框架 |
| TypeScript | 5.9+ | 类型系统 |
| shadcn/ui | 最新 | UI 组件方案 |
| Tailwind CSS | 4.x | 样式系统 |

### 2.4 部署技术栈

| 技术 | 说明 |
|-----|------|
| Docker | 容器化 |
| Docker Compose | 本地开发编排 |
| Nginx | 反向代理 + 静态资源 |
| GitHub Actions | CI/CD |

---

## 3. 数据库设计

### 3.1 数据库概览

GameLink 使用 PostgreSQL 16+ 作为主数据库，设计了 **80+ 张表**，覆盖用户、订单、支付、聊天、营销等核心业务。

### 3.2 核心数据表分类

#### 用户与陪玩师体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `users` | 用户基础信息 | phone, email, role, vip_level_id |
| `players` | 陪玩师扩展信息 | user_id, rating_average, verification_status |
| `wallets` | 用户钱包 | user_id, balance_cents, frozen_cents |
| `player_certifications` | 陪玩师实名认证 | player_id, id_card, status |
| `player_rank_records` | 陪玩师段位认证 | player_id, game_id, rank_image |
| `player_services` | 陪玩师服务 | player_id, game_id, hourly_rate_cents |
| `player_schedules` | 陪玩师排班 | player_id, weekly_schedule |
| `player_presence` | 陪玩师在线状态 | player_id, online_status |

#### 订单与支付体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `orders` | 订单主表 | user_id, player_id, status, total_price_cents |
| `order_groups` | 主订单（订单拆分） | user_id, status, total_amount |
| `order_items` | 订单明细 | order_id, player_id, price_cents |
| `payments` | 支付记录 | order_id, method, amount_cents, status |
| `refund_records` | 退款记录 | payment_id, amount_cents, status |
| `order_disputes` | 订单争议 | order_id, initiator_id, resolution |

#### 服务与商品体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `service_items` | 服务项统一管理 | item_code, category, base_price_cents |
| `games` | 游戏列表 | key, name, category_id |
| `game_categories` | 游戏分类 | name, icon_url |
| `game_ranks` | 游戏段位配置 | game_id, rank_name |

#### 聊天与通讯体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `chat_groups` | 聊天房间 | group_name, group_type, related_order_id |
| `chat_group_members` | 房间成员 | group_id, user_id, role |
| `chat_messages` | 聊天消息 | group_id, sender_id, content |
| `chat_snapshots` | 争议聊天快照 | dispute_id, messages |

#### 营销与会员体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `vip_levels` | VIP等级配置 | slug, exp_required, order_discount |
| `coupons` | 用户优惠券 | user_id, template_id, state |
| `coupon_templates` | 优惠券模板 | name, type, deduct_amount_cents |
| `activities` | 营销活动 | name, type, status, start_at |
| `referral_codes` | 推荐邀请码 | user_id, code, use_count |
| `referrals` | 推荐关系 | referrer_id, referee_id, status |

#### 财务与结算体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `commission_rules` | 佣金规则 | type, rate, game_id |
| `commission_records` | 佣金记录 | player_id, order_id, commission_cents |
| `monthly_settlements` | 月度结算 | player_id, settlement_month |
| `withdraws` | 提现记录 | user_id, amount_cents, status |
| `settlement_companies` | 结算公司 | name, bank_account |
| `routing_rules` | 路由规则 | condition_json, target_entity_id |

#### 内容与审核体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `reviews` | 订单评价 | order_id, user_id, score |
| `sensitive_words` | 敏感词库 | word, category |
| `feeds` | 动态内容 | user_id, content, status |
| `content_categories` | 内容分类 | name, icon_url |

#### 权限与管理体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `permissions` | 权限定义 | code, resource, action |
| `role_models` | 角色定义 | slug, name, is_system |
| `role_permissions` | 角色权限关联 | role_id, permission_id |
| `user_roles` | 用户角色关联 | user_id, role_id |
| `menus` | 前端菜单配置 | name, path, parent_id |

#### 团队与社交体系

| 表名 | 说明 | 关键字段 |
|-----|------|---------|
| `teams` | 陪玩师团队 | name, leader_id, max_members |
| `team_members` | 团队成员 | team_id, user_id, role |
| `lfg_requests` | 组队请求 | user_id, game_id |
| `favorites` | 收藏记录 | user_id, target_type |

### 3.3 数据库索引优化

项目使用了 **PostgreSQL 覆盖索引** 优化查询性能：

```sql
-- 订单列表查询优化
CREATE INDEX idx_orders_user_status_created_covering ON orders
  (user_id, status, created_at DESC)
  INCLUDE (id, player_id, total_price_cents);

-- 支付记录查询优化
CREATE INDEX idx_payments_user_status_created_covering ON payments
  (user_id, status, created_at DESC)
  INCLUDE (id, amount_cents, payment_method);

-- 聊天消息查询优化
CREATE INDEX idx_chat_messages_group_sent_covering ON chat_messages
  (group_id, created_at DESC)
  INCLUDE (id, content, sender_id);
```

---

## 4. 功能模块

### 4.1 用户端（Web）

#### 核心功能

| 模块 | 功能 | 说明 |
|-----|------|------|
| **首页** | 热门推荐、游戏分类、陪玩师列表 | 首页入口 |
| **陪玩师浏览** | 筛选、排序、搜索、详情 | 发现服务 |
| **下单流程** | 创建订单、支付、订单管理 | 核心交易 |
| **聊天通讯** | 订单房间、实时消息 | 沟通协调 |
| **评价系统** | 星级评分、文字评价 | 质量反馈 |
| **钱包充值** | 余额、充值、交易记录 | 资金管理 |
| **个人中心** | 资料、收藏、设置 | 账户管理 |

#### 页面清单

用户端页面代码位于 `app/src/features/` 与 `app/src/components/`，采用 React Router 进行路由组织。

### 4.2 陪玩师端

#### 核心功能

| 模块 | 功能 | 说明 |
|-----|------|------|
| **工作台** | 今日统计、快捷入口、收益概览 | 数据看板 |
| **接单管理** | 订单大厅、我的订单、接单/拒绝 | 订单处理 |
| **服务管理** | 添加服务、上下架、定价 | 服务配置 |
| **收益提现** | 收益统计、提现申请、提现记录 | 收益管理 |
| **认证管理** | 实名认证、段位认证 | 准入控制 |
| **在线状态** | 在线/离线/忙碌切换 | 状态管理 |

### 4.3 管理后台

#### 核心功能

| 模块 | 功能 | 说明 |
|-----|------|------|
| **仪表盘** | 核心指标、实时监控、快捷入口 | 运营总览 |
| **用户管理** | 用户列表、详情、标签、封禁 | 用户运营 |
| **陪玩师管理** | 列表、认证审核、状态管理 | 供给管理 |
| **订单管理** | 列表、详情、状态流转、争议 | 交易管理 |
| **财务管理** | 结算、佣金、提现、充值 | 资金管理 |
| **内容审核** | 聊天监控、敏感词、举报 | 内容安全 |
| **营销管理** | 活动、优惠券、VIP、推荐 | 增长工具 |
| **系统管理** | 角色、权限、菜单、路由 | 系统配置 |
| **监控分析** | 实时监控、数据分析、KPI | 数据洞察 |

#### 页面清单

```
admin/src/pages/
├── Dashboard/             # 仪表盘
├── User/                  # 用户管理
├── Player/               # 陪玩师管理
├── PlayerCertification/  # 陪玩师认证
├── Order/                # 订单管理
├── Settlement/           # 结算管理
├── Commission/           # 佣金管理
├── Withdraw/             # 提现管理
├── PaymentRecords/       # 支付记录
├── Recharge/             # 充值管理
├── Content/              # 内容审核
├── Review/               # 评价管理
├── Activity/             # 活动管理
├── Coupon/               # 优惠券管理
├── VIP/                  # VIP管理
├── Referral/             # 推荐管理
├── Role/                 # 角色管理
├── Permission/           # 权限管理
├── Team/                 # 团队管理
├── Game/                 # 游戏管理
├── Service/              # 服务项管理
├── Dispute/              # 争议处理
├── RoutingRule/          # 路由规则
└── Monitor/              # 监控
```

---

## 5. 开发和部署指南

### 5.1 环境要求

| 工具 | 版本要求 |
|-----|---------|
| Go | 1.24+ |
| Node.js | 20+ |
| PostgreSQL | 16+ |
| Redis | 7+ |
| Docker & Docker Compose | 最新版（可选） |

### 5.2 快速开始

#### 1. 克隆项目

```bash
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

#### 2. 环境配置

```bash
cp .env.example .env
# 编辑 .env，配置数据库、Redis、JWT 密钥等
```

#### 3. Docker 启动（推荐）

```bash
# 启动基础服务（PostgreSQL + Redis）
docker compose up -d

# 查看服务状态
docker compose ps
```

#### 4. 启动后端

```bash
cd api
go mod download
go run cmd/main.go

# 服务运行在 http://localhost:8080
# Swagger 文档：http://localhost:8080/swagger/index.html
```

#### 5. 启动管理后台

```bash
cd admin
npm install
npm run dev

# 访问 http://localhost:5173
```

#### 6. 启动用户端 Web

```bash
cd app
npm install

# 启动开发
npm run dev
# 访问 http://localhost:5175
```

### 5.3 默认账号

| 角色 | 账号 | 密码 |
|------|------|------|
| 超级管理员 | admin@gamelink.com | Admin123456 |

> 启动后端后自动创建种子数据，包含演示用户和陪玩师。

### 5.4 API 文档

后端启动后访问 Swagger：

```
http://localhost:8080/swagger/index.html
```

API 基础路径：`/api/v1`

主要路由分组：

| 前缀 | 说明 | 认证 |
|------|------|------|
| `/api/v1/auth` | 登录、注册、刷新 Token | 部分需要 |
| `/api/v1/public` | 公开接口（陪玩师列表、游戏等） | 不需要 |
| `/api/v1/user` | 用户端接口（订单、钱包、聊天等） | 需要 |
| `/api/v1/player` | 陪玩师端接口（接单、服务、收益等） | 需要 |
| `/api/v1/admin` | 管理端接口（全部功能） | 需要 + RBAC |

### 5.5 测试

#### 后端测试

```bash
cd api
go test ./... -v -cover
# 159 个测试文件，覆盖率目标 70%+
```

#### 管理后台测试

```bash
cd admin
npm run test
# 88 个单元测试
```

### 5.6 CI/CD

项目配置了 4 条 GitHub Actions 流水线：

| 流水线 | 说明 |
|--------|------|
| `ci.yml` | 主 CI（lint、test、build、Docker 镜像） |
| `deploy.yml` | 自动部署 |
| `pre-commit-check.yml` | 提交前检查 |
| `security.yml` | 安全扫描 |

### 5.7 部署到生产环境

#### 环境变量检查清单

生产环境部署前，请确保以下配置正确：

- [ ] `APP_ENV=production`
- [ ] `GIN_MODE=release`
- [ ] `CRYPTO_ENABLED=true`
- [ ] `CRYPTO_SECRET_KEY` 已设置（32字节）
- [ ] `CRYPTO_IV` 已设置（16字节）
- [ ] `JWT_SECRET_KEY` 已设置（32+字符）
- [ ] `SUPER_ADMIN_EMAIL` 已设置
- [ ] `SUPER_ADMIN_PASSWORD` 已设置（强密码）
- [ ] `SEED_ENABLED=false`
- [ ] 数据库密码已修改
- [ ] Redis 密码已设置
- [ ] 微信支付/支付宝配置已填写（如需要）

---

## 6. 团队和联系方式

### 6.1 项目团队

| 角色 | 职责 |
|-----|------|
| **Backend-Lead** | 后端架构设计、API 开发 |
| **Frontend-Lead** | 管理后台开发、React 技术栈 |
| **Frontend-Lead (Web)** | 用户端 Web 开发、React 技术栈 |
| **Database-Architect** | 数据库设计、性能优化 |
| **DevOps-Engineer** | 部署、CI/CD、运维 |
| **Product-Manager** | 产品规划、需求分析 |

### 6.2 相关文档

| 文档 | 路径 | 说明 |
|-----|------|------|
| 项目 README | `/README.md` | 项目简介 |
| 产品需求文档 | `/docs/PRD.md` | PRD v1.0 |
| 完整PRD | `/docs/PRD_COMPREHENSIVE.md` | PRD v2.0（详细版） |
| 项目进度 | `/docs/PROGRESS.md` | 版本历史 |
| 用户端 Web 实现 | `/app/src` | React + shadcn/ui + Tailwind |

### 6.3 代码规范

- **Go**: 遵循 `gofmt` 和 `Effective Go` 规范
- **TypeScript**: 遵循 ESLint 配置
- **Vue 3**: 遵循 Vue 3 风格指南
- **Commit**: 遵循 Conventional Commits 规范

### 6.4 分支策略

```
main (生产)
  ↑
  dev (开发)
  ↑
  feature/* (功能分支)
  bugfix/* (修复分支)
```

---

## 附录

### A. 常见问题

**Q: 后端启动失败，提示数据库连接错误？**

A: 请检查 PostgreSQL 和 Redis 是否已启动，检查 `.env` 中的数据库配置是否正确。

**Q: 管理后台无法登录？**

A: 确保后端已启动，检查 `admin/src/api` 中的 API 地址配置。

**Q: 用户端 Web 无法访问接口？**

A: 检查 `app/src/services` 中的 API 配置，确保后端已启动并允许跨域。

### B. 性能优化建议

1. **数据库**：已使用覆盖索引优化，定期 VACUUM
2. **缓存**：热点数据使用 Redis 缓存
3. **CDN**：静态资源建议使用 CDN 加速
4. **负载均衡**：生产环境建议使用 Nginx 负载均衡

### C. 安全建议

1. **生产环境必须启用加密**：`CRYPTO_ENABLED=true`
2. **定期更新依赖**：`go get -u all` / `npm update`
3. **使用强密码**：数据库、Redis、JWT 密钥
4. **启用 HTTPS**：生产环境必须使用 SSL/TLS
5. **定期备份**：数据库定期备份

---

**文档结束**

> 本文档整合了技术团队所有成员的分析报告，提供了 GameLink 项目的全面概览。如有疑问，请联系相关负责人。
