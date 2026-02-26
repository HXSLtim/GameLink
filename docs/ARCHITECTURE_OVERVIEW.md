# GameLink 项目架构总览

> 最后更新: 2025-02-09
> 文档版本: 1.0

## 项目简介

GameLink 是一个现代化的游戏陪玩社交平台，连接游戏玩家与陪玩师，提供从下单、匹配、服务、评价到结算的完整闭环。

### 核心功能
- 🔍 **陪玩师发现**: 浏览陪玩师档案、技能标签、实时状态
- 📝 **订单管理**: 创建护航服务订单、多人组队、礼物赠送
- 💬 **即时通讯**: 订单服务群、私聊、团队群、LFG快速匹配
- 💰 **财务系统**: 钱包充值、订单结算、提现管理、佣金分配
- 🎮 **内容社区**: 动态发布、游戏攻略、精彩时刻分享
- 👥 **团队管理**: 陪玩师团队、排班管理、收入统计

---

## 技术栈

| 层级 | 技术选型 | 版本 | 说明 |
|------|---------|------|------|
| **后端** | Go + Gin + GORM | 1.24.5 | RESTful API + WebSocket |
| **管理后台** | React + Ant Design + TypeScript | React 19 / AntD 6 | 后台管理系统 |
| **用户端 Web** | React + shadcn/ui + Tailwind CSS | React 19 | Web 前端 |
| **数据库** | PostgreSQL | 16+ | 主数据存储 |
| **缓存** | Redis | 7+ | 会话/热点数据缓存 |
| **对象存储** | OSS | - | 图片/文件存储 |
| **实时通讯** | TRTC | - | 腾讯云语音服务 |
| **部署** | Docker + Nginx | - | 容器化部署 |

---

## 项目结构

```
GameLink/
├── api/                          # Go 后端服务
│   ├── cmd/
│   │   └── main.go              # 应用入口
│   ├── internal/
│   │   ├── handler/             # HTTP 处理器
│   │   │   ├── admin/           # 管理员接口
│   │   │   ├── user/            # 用户接口
│   │   │   ├── player/          # 陪玩师接口
│   │   │   ├── public/          # 公开接口
│   │   │   ├── middleware/      # 中间件
│   │   │   └── ws/              # WebSocket
│   │   ├── service/             # 业务逻辑层 (57个模块)
│   │   ├── repository/          # 数据访问层 (56个模块)
│   │   ├── model/               # 数据模型 (67个)
│   │   └── router/              # 路由注册
│   └── pkg/                     # 公共包
│       ├── auth/                # JWT/认证
│       ├── config/              # 配置管理
│       ├── db/                  # 数据库/迁移/种子数据
│       └── scheduler/           # 定时任务
│
├── admin/                        # React 管理后台
│   └── src/
│       ├── pages/               # 40+ 功能模块
│       │   ├── admin/           # 系统管理
│       │   ├── adminChat/       # 聊天管理
│       │   ├── auth/            # 认证
│       │   └── biz/             # 业务管理
│       ├── components/          # 通用组件
│       ├── api/                 # API 封装
│       ├── router/              # 路由配置
│       └── utils/               # 工具函数
│
├── app/                          # 用户端 Web
│   └── src/
│       ├── features/            # 业务功能页面
│       ├── components/          # 组件（含 shadcn/ui）
│       ├── hooks/               # 业务 Hook
│       ├── services/            # API 服务
│       └── lib/                 # 工具函数
│
├── docs/                         # 项目文档
├── scripts/                      # 部署脚本
└── docker-compose.yml            # Docker 编排
```

---

## 数据库架构

### 核心数据表分类（80+张表）

#### 1. 用户与权限系统 (10张表)
- `users` - 用户主表（多角色、VIP、微信登录）
- `roles` - 角色定义（super_admin, admin, player, user）
- `permissions` - 权限定义
- `user_roles` - 用户角色关联
- `user_tags` / `user_tag_relations` - 用户标签
- `user_login_history` - 登录历史
- `user_behavior` - 行为追踪
- `permission_audit_logs` - 权限审计
- `menus` - 动态菜单

#### 2. 陪玩师体系 (8张表)
- `players` - 陪玩师档案
- `player_games` - 游戏关联
- `player_skill_tags` - 技能标签
- `game_ranks` - 段位配置
- `player_rank_records` - 段位认证
- `player_certifications` - 实名认证
- `player_schedules` - 排班表
- `player_presence` - 在线状态

#### 3. 订单系统 (7张表)
- `orders` - 订单主表（护航/礼物/团队单）
- `order_groups` - 主订单（多时段拆分）
- `order_items` - 订单明细（多人槽位）
- `order_players` - 订单陪玩师关联
- `payments` - 支付记录（钱包/微信/组合）
- `order_disputes` - 争议处理
- `order_timeout_logs` - 超时日志

#### 4. 游戏与服务 (6张表)
- `games` / `game_categories` - 游戏配置
- `service_items` - 服务项目（护航/礼物）
- `coupons` / `coupon_templates` - 优惠券
- `recharge_options` / `recharge_records` - 充值

#### 5. 评价与社交 (10张表)
- `reviews` - 评价主表
- `review_reports` / `review_replies` - 评价管理
- `favorites` - 收藏
- `user_blocks` - 拉黑
- `chat_groups` - 聊天群组（7种类型）
- `chat_group_members` - 群成员
- `chat_messages` - 聊天消息
- `chat_reports` - 聊天举报

#### 6. 内容社区 (4张表)
- `feeds` / `feed_images` - 动态
- `feed_reports` - 举报
- `content_categories` - 分类

#### 7. 财务系统 (8张表)
- `wallets` - 钱包
- `withdraws` - 提现
- `commission_rules` / `commission_records` - 抽成
- `monthly_settlements` - 月度结算
- `collection_entities` - 收款主体
- `settlement_companies` - 结算公司
- `refund_records` - 退款

#### 8. 团队与活动 (8张表)
- `teams` / `team_members` / `team_invites` - 团队管理
- `activities` / `activity_rewards` / `activity_participations` - 活动
- `referrals` / `referral_codes` / `referral_rewards` - 推荐
- `vip_levels` / `vip_configs` - VIP系统

#### 9. 通知与监控 (9张表)
- `notification_templates` / `user_notifications` / `user_notification_settings` - 通知
- `alerts` / `kpi_targets` - 监控
- `operation_logs` - 操作日志
- `sensitive_words` - 敏感词
- `order_timeout_configs` / `order_service_assignments` - 超时配置

#### 10. 统计与LFG (7张表)
- `user_statistics` / `player_statistics` / `service_item_statistics` - 统计
- `lfg_requests` - 快速组队
- `tag_thresholds` - 标签阈值

### 核心数据关系

```
用户 (User) 1:1 陪玩师档案 (Player)
    ├── 1:N 订单 (Order) → 1:N 支付 (Payment)
    ├── 1:1 钱包 (Wallet)
    ├── 1:N 收藏 (Favorite)
    └── N:M 角色 (UserRole + Role)

游戏 (Game) 1:N 服务项目 (ServiceItem)
陪玩师 (Player) 1:N 团队成员 (TeamMember)

订单 (Order) N:1 主订单 (OrderGroup)
    ├── 1:N 订单明细 (OrderItem)
    ├── 1:N 订单陪玩师 (OrderPlayer)
    └── 1:1 聊天群组 (ChatGroup - order类型)
```

---

## 前端架构

### 用户端 Web 技术栈

```
app/
├── src/
│   ├── features/           # 页面功能
│   ├── components/         # 组件（含 shadcn/ui）
│   ├── hooks/              # 业务 Hook
│   ├── services/           # API 模块
│   ├── lib/                # 工具函数
│   └── index.css           # Tailwind + 主题变量
```

#### 组件分层

1. **基础组件 (shadcn/ui)**: 基于 Radix + Tailwind 的原子组件
   - `Button`, `Input`, `Dialog`, `Sheet` 等

2. **模式组件**: 可复用的UI模式
   - `NavBar`, `SearchBar`, `FilterPanel`, `InfiniteList` 等

3. **业务组件**: 领域特定的复合组件
   - `PlayerCard`, `OrderCard`, `ChatMessageBubble`, `WalletBalanceCard` 等

#### Hooks 设计

业务 Hook 封装核心逻辑：
- `usePlayerList` - 陪玩师列表
- `useOrderCreate` - 订单创建
- `useChatRoom` - 聊天室
- `useWallet` - 钱包管理
- `useWebSocket` - WebSocket连接

### React 管理后台

```
admin/
├── src/
│   ├── pages/              # 40+功能模块
│   ├── components/         # 通用组件
│   ├── api/                # API封装
│   ├── router/             # 路由配置
│   └── utils/              # 工具函数
```

---

## 后端架构

### 分层设计

```
Handler (HTTP/WebSocket)
    ↓
Service (业务逻辑)
    ↓
Repository (数据访问)
    ↓
Model (数据模型)
```

### 核心模块

- **Handler**: 57个路由处理器
  - `admin/` - 管理员接口
  - `user/` - 用户接口
  - `player/` - 陪玩师接口
  - `public/` - 公开接口

- **Service**: 57个业务服务
  - 订单、支付、聊天、推荐、结算等

- **Repository**: 56个数据仓储
  - 封装数据库操作，支持测试Mock

- **Model**: 67个数据模型
  - 用户、订单、陪玩师、游戏等

### 中间件

- `jwtAuth` - JWT认证
- `permissionSync` - 权限同步
- `redisRateLimiter` - 限流
- `signature` - 签名验证
- `cors` - 跨域

---

## 关键设计模式

### 1. 订单拆分机制
- `OrderGroup` - 用户视角的主订单
- `Order` - 按时段拆分的子订单
- 支持转单 (`transfer_from` → `transfer_to`)

### 2. 多陪玩师支持
- `OrderItem` - 订单槽位 (slot 1, 2, 3...)
- `OrderPlayer` - 槽位与陪玩师关联
- 团队服务项目 (`sub_category='team'`)

### 3. 混合支付
- 钱包支付 (`PaymentMethodWallet`)
- 第三方支付（微信）
- 组合支付 (`PaymentMethodCombined`)

### 4. 聊天群组（7种类型）
- `public` - 公开群
- `order` - 订单服务群
- `private` - 私聊 (1v1)
- `team` - 团队群
- `lfg` - 快速匹配群
- `custom` - 自定义群

### 5. 状态管理
- 用户状态 (active/suspended/banned)
- 订单状态 (pending/confirmed/in_progress/completed/canceled/refunded)
- 聊天审核 (pending/approved/rejected/deleted)

---

## 部署架构

### Docker Compose

```yaml
services:
  postgres:    # PostgreSQL 16
  redis:       # Redis 7
  api:         # Go后端 (8080)
  admin:       # React管理后台 (3000)
  nginx:       # 反向代理
```

### 环境配置

- `.env.example` - 环境变量模板
- 支持开发/生产环境配置
- 数据库连接、Redis、JWT密钥等

---

## 开发指南

### 快速开始

1. **克隆项目**
```bash
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

2. **环境配置**
```bash
cp .env.example .env
# 编辑 .env 配置数据库、Redis等
```

3. **启动服务**
```bash
# Docker方式
docker compose up -d

# 手动启动
cd api && go run main.go          # 后端
cd admin && npm run dev           # 管理后台
cd app && npm run dev              # 用户端 Web
```

### 代码规范

- **Go**: 遵循 `gofmt` 和 `golangci-lint`
- **React/TypeScript**: ESLint + Prettier
- **Tailwind/shadcn**: 组件优先 + utility-first 样式规范

### 测试

```bash
# 后端测试
cd api && go test ./...

# 前端测试
cd admin && npm test
```

---

## 文档索引

- [数据库架构详细说明](./DATABASE.md)
- [API接口文档](./API.md)
- [前端组件库文档](./COMPONENTS.md)
- [部署指南](./DEPLOYMENT.md)
- [开发规范](./CONTRIBUTING.md)

---

## 联系方式

- 项目负责人: [HXSL](mailto:a2778978136@163.com)
- 技术支持: [GitHub Issues](https://github.com/your-org/gamelink/issues)
