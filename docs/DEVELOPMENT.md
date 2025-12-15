# GameLink 开发文档

## 1. 技术架构

### 1.1 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| **后端** | Go + Gin + GORM | Go 1.25.3+ |
| **前端** | React + TypeScript + Vite | React 18.2+, Vite 7.2+ |
| **UI** | Ant Design | 6.0 |
| **数据库** | PostgreSQL | 16+ |
| **缓存** | Redis | 7+ |
| **容器** | Docker + Docker Compose | - |

### 1.2 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (React)                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │  Pages  │  │ Context │  │  API    │  │ Crypto (AES-256)│ │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ HTTPS + AES-256-CBC
┌─────────────────────────────────────────────────────────────┐
│                        Backend (Go/Gin)                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │ Handler │→ │ Service │→ │  Repo   │→ │     Model       │ │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Middleware: Auth(JWT) | Crypto | CORS | Permission     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │PostgreSQL│   │  Redis   │   │WebSocket │
        └──────────┘   └──────────┘   └──────────┘
```

---

## 2. 模块分布

### 2.1 后端模块 (`backend/internal/`)

#### Handler 层 (HTTP 处理器)
```
handler/
├── admin/           # 管理端 API (35+ handlers)
│   ├── analytics.go      # 运营分析
│   ├── commission.go     # 佣金管理
│   ├── content.go        # 内容管理
│   ├── dashboard.go      # 仪表盘
│   ├── game.go           # 游戏管理
│   ├── kpi.go            # KPI 指标
│   ├── menu.go           # 菜单管理
│   ├── monitor.go        # 实时监控
│   ├── order.go          # 订单管理
│   ├── permission.go     # 权限管理
│   ├── player.go         # 陪玩师管理
│   ├── review.go         # 评价管理
│   ├── role.go           # 角色管理
│   ├── user.go           # 用户管理
│   ├── withdraw.go       # 提现管理
│   └── ...
├── middleware/      # 中间件
│   ├── auth.go           # JWT 认证
│   ├── crypto.go         # 加密解密
│   ├── cors.go           # 跨域处理
│   └── permission.go     # 权限校验
├── player/          # 陪玩师端 API
├── user/            # 用户端 API
└── notification/    # 通知 API
```

#### Service 层 (业务逻辑)
```
service/
├── admin/           # 管理服务
├── analytics/       # 数据分析
├── assignment/      # 订单分配
├── audit/           # 审计日志
├── auth/            # 认证服务
├── chat/            # 聊天服务
├── commission/      # 佣金计算
├── content/         # 内容服务
├── earnings/        # 收益服务
├── feed/            # 动态服务
├── gift/            # 礼物服务
├── item/            # 服务项目
├── kpi/             # KPI 服务
├── menu/            # 菜单服务
├── monitor/         # 监控服务
├── notification/    # 通知服务
├── order/           # 订单服务
├── payment/         # 支付服务
├── permission/      # 权限服务
├── player/          # 陪玩师服务
├── ranking/         # 排行榜
├── review/          # 评价服务
├── role/            # 角色服务
├── sensitiveword/   # 敏感词
├── stats/           # 统计服务
├── team/            # 团队服务
├── user/            # 用户服务
├── wallet/          # 钱包服务
└── withdraw/        # 提现服务
```

#### Repository 层 (数据访问)
```
repository/
├── admin/           # 管理数据访问
├── chat/            # 聊天数据
├── commission/      # 佣金数据
├── content/         # 内容数据
├── dispute/         # 纠纷数据
├── feed/            # 动态数据
├── game/            # 游戏数据
├── menu/            # 菜单数据
├── notification/    # 通知数据
├── order/           # 订单数据
├── payment/         # 支付数据
├── permission/      # 权限数据
├── player/          # 陪玩师数据
├── review/          # 评价数据
├── user/            # 用户数据
├── wallet/          # 钱包数据
├── withdraw/        # 提现数据
├── interfaces/      # 接口定义
└── mocks/           # Mock 实现
```

#### Model 层 (数据模型)
```
model/
├── user.go          # 用户模型
├── player.go        # 陪玩师模型
├── order.go         # 订单模型
├── payment.go       # 支付模型
├── review.go        # 评价模型
├── wallet.go        # 钱包模型
├── withdraw.go      # 提现模型
├── game.go          # 游戏模型
├── serviceItem.go   # 服务项目
├── role.go          # 角色模型
├── permission.go    # 权限模型
├── menu.go          # 菜单模型
├── chat.go          # 聊天模型
├── commission.go    # 佣金模型
├── dispute.go       # 纠纷模型
└── ...
```

### 2.2 前端模块 (`frontend/src/`)

#### 页面结构
```
pages/
├── admin/           # 管理后台页面
│   ├── Commission/       # 佣金管理
│   ├── Content/          # 内容管理
│   ├── Dashboard/        # 仪表盘
│   ├── Game/             # 游戏管理
│   ├── Monitor/          # 监控中心
│   ├── Notifications/    # 通知管理
│   ├── Order/            # 订单管理
│   ├── Permission/       # 权限管理
│   ├── Player/           # 陪玩师管理
│   ├── Review/           # 评价管理
│   ├── Role/             # 角色管理
│   ├── User/             # 用户管理
│   └── Withdraw/         # 提现管理
├── auth/            # 认证页面
│   ├── Login.tsx
│   └── Register.tsx
├── biz/             # 业务页面
│   └── service/
└── sys/             # 系统页面
    ├── log/
    ├── menu/
    ├── permission/
    ├── setting/
    └── user-role/
```

#### 核心模块
```
src/
├── api/             # API 客户端
│   ├── client.ts         # Axios 实例 (带加密)
│   ├── auth.ts           # 认证 API
│   └── admin.ts          # 管理 API
├── components/      # 可复用组件
│   ├── common/           # 通用组件
│   └── layout/           # 布局组件
├── context/         # React Context
│   ├── AuthContext.tsx   # 认证上下文
│   ├── AdminContext.tsx  # 管理上下文
│   └── ThemeContext.tsx  # 主题上下文
├── utils/           # 工具函数
│   ├── crypto.ts         # 加密工具
│   ├── dynamicRoutes.tsx # 动态路由
│   └── menuPermission.ts # 菜单权限
├── constants/       # 常量定义
│   └── permissions.ts    # 权限码
└── types/           # TypeScript 类型
```

---

## 3. 权限系统

### 3.1 权限码格式
```
{module}.{resource}.{action}

示例：
- admin.users.list      # 用户列表
- admin.orders.refund   # 订单退款
- content.feed.approve  # 动态审核
```

### 3.2 权限模块分布

| 模块 | 权限数量 | 说明 |
|------|---------|------|
| 用户管理 | 6 | list, read, create, update, delete, status |
| 游戏管理 | 5 | list, read, create, update, delete |
| 订单管理 | 7 | list, read, create, update, delete, cancel, refund |
| 陪玩师管理 | 6 | list, read, create, update, delete, audit |
| 角色管理 | 7 | list, read, create, update, delete, permissions, assign-user |
| 权限管理 | 6 | list, read, create, update, delete, groups |
| 菜单管理 | 5 | list, read, create, update, delete |
| 服务项目 | 6 | list, read, create, update, delete, batch-status |
| 佣金管理 | 5 | list, read, create, update, settle |
| 提现管理 | 5 | list, read, approve, reject, complete |
| 内容管理 | 15+ | feed, chat, report, category 相关 |
| 评价管理 | 10+ | reviews, review-reports, sensitive-words 相关 |
| 监控中心 | 10+ | monitor, analytics, kpi 相关 |

---

## 4. 开发指南

### 4.1 环境准备

```bash
# 后端依赖
cd backend
go mod tidy

# 前端依赖
cd frontend
npm install
```

### 4.2 本地开发

```bash
# 启动后端 (需要 PostgreSQL + Redis)
cd backend
go run cmd/main.go

# 启动前端
cd frontend
npm run dev
```

### 4.3 Docker 部署

```powershell
# 生产环境部署（推荐）
.\scripts\deploy-production-encrypted.ps1

# 标准部署
.\scripts\deploy-production.ps1
```

### 4.4 常用命令

```bash
# 后端
make test           # 运行测试
make lint           # 代码检查
make swagger        # 生成 API 文档

# 前端
npm run build       # 生产构建
npm run lint        # 代码检查
npm run test        # 运行测试
```

---

## 5. API 规范

### 5.1 请求格式
```json
{
  "data": "加密后的请求体 (生产环境)",
  "sign": "SHA-256 签名"
}
```

### 5.2 响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 5.3 错误码
| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

---

## 6. 数据库设计

### 6.1 核心表

| 表名 | 说明 | 关联 |
|------|------|------|
| users | 用户表 | → players, orders, reviews |
| players | 陪玩师表 | → users, orders |
| orders | 订单表 | → users, players, payments |
| payments | 支付表 | → orders, users |
| reviews | 评价表 | → orders, users, players |
| wallets | 钱包表 | → users |
| withdraws | 提现表 | → users, wallets |
| games | 游戏表 | → service_items |
| service_items | 服务项目表 | → games |
| roles | 角色表 | → permissions, menus |
| permissions | 权限表 | → roles |
| menus | 菜单表 | → roles |

### 6.2 索引策略
- 主键索引：所有表 ID
- 外键索引：user_id, player_id, order_id 等
- 状态索引：status 字段
- 时间索引：created_at, updated_at

---

## 7. 测试策略

### 7.1 测试类型
- 单元测试：表驱动测试，mock 依赖
- 集成测试：真实数据库，测试 fixtures
- 并发测试：Race detector，压力测试

### 7.2 测试覆盖率
- 当前：76.4%
- 目标：80%+

### 7.3 集成测试文件
```
integration/
├── auth_integration_test.go
├── order_*_integration_test.go
├── payment_*_integration_test.go
├── review_*_integration_test.go
├── admin_*_integration_test.go
├── commission_integration_test.go
├── withdraw_integration_test.go
└── ...
```

---

## 8. 安全措施

### 8.1 通信安全
- AES-256-CBC 加密请求/响应体
- SHA-256 签名验证
- HTTPS 传输

### 8.2 认证安全
- JWT Token + 刷新机制
- Token 过期时间：2 小时
- 刷新 Token 过期：7 天

### 8.3 权限安全
- RBAC 角色权限控制
- 接口级权限校验
- 菜单级权限控制

### 8.4 数据安全
- 密码 bcrypt 加密存储
- 敏感数据脱敏
- SQL 注入防护
