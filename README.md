# GameLink

现代化游戏陪玩社交平台，连接游戏玩家与陪玩师，提供从下单、匹配、服务、评价到结算的完整闭环。

---

## 技术栈

| 层 | 技术 | 版本 |
|----|------|------|
| 后端 | Go + Gin + GORM | Go 1.24.5 |
| 管理后台 | React + Ant Design + TypeScript | React 19 / AntD 6 |
| 小程序/H5 | uni-app + Vue 3 + TypeScript | Vue 3.4 |
| 数据库 | PostgreSQL | 16+ |
| 缓存 | Redis | 7+ |
| 部署 | Docker + Nginx | - |

---

## 项目结构

```
GameLink/
├── api/                # Go 后端服务
│   ├── internal/
│   │   ├── handler/    # 路由处理（admin/user/player/public）
│   │   ├── service/    # 业务逻辑（57 个模块）
│   │   ├── repository/ # 数据访问（56 个模块）
│   │   ├── model/      # 数据模型（67 个）
│   │   ├── router/     # 路由注册
│   │   └── ws/         # WebSocket
│   └── pkg/            # 公共包（auth/config/db/scheduler）
├── admin/              # 管理后台（React 19）
│   └── src/
│       ├── pages/      # 40+ 功能模块
│       ├── components/ # 通用组件
│       ├── api/        # API 封装
│       └── router/     # 路由配置
├── app/                # 小程序/H5（uni-app）
│   └── src/
│       ├── pages/      # 28 个页面
│       ├── components/ # 133 个组件
│       ├── composables/# 38 个 Hook
│       ├── api/        # 14 个 API 模块
│       ├── store/      # 3 个 Store
│       ├── styles/     # 主题系统（CSS 变量）
│       └── types/      # 27 个类型定义
├── docs/               # 项目文档（PRD/进度/架构）
├── scripts/            # 部署脚本
└── docker-compose.yml  # Docker 编排
```

---

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 20+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose（可选）

### 1. 克隆项目

```bash
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

### 2. 环境配置

```bash
cp .env.example .env
# 编辑 .env，配置数据库、Redis、JWT 密钥等
```

### 3. Docker 一键启动（推荐）

```bash
# 启动基础服务（PostgreSQL + Redis）
docker compose up -d

# 启动开发环境（后端 + 管理后台）
docker compose -f docker-compose.dev.yml up -d
```

### 4. 手动启动

**后端：**

```bash
cd api
go mod download
go run main.go
# 服务运行在 http://localhost:8080
# Swagger 文档：http://localhost:8080/swagger/index.html
```

**管理后台：**

```bash
cd admin
npm install
npm run dev
# 访问 http://localhost:5173
```

**小程序/H5：**

```bash
cd app
npm install

# H5 开发
npm run dev:h5
# 访问 http://localhost:5174

# 微信小程序
npm run dev:mp-weixin
# 用微信开发者工具打开 dist/dev/mp-weixin
```

---

## 默认账号

| 角色 | 账号 | 密码 |
|------|------|------|
| 超级管理员 | admin@gamelink.com | Admin123456 |

启动后端后自动创建种子数据，包含演示用户和陪玩师。

---

## 核心功能

### 用户端（小程序/H5）

- 浏览陪玩师（筛选、排序、搜索）
- 下单 + 支付（余额/微信/支付宝）
- 实时聊天（WebSocket）
- 订单管理（状态流转、取消、评价）
- 钱包与充值
- 收藏、帮助中心、在线客服

### 陪玩师端

- 工作台（今日统计、快捷入口）
- 接单管理（可接订单、我的订单）
- 服务管理（上下架、定价）
- 收益与提现
- 实名认证 + 段位认证
- 在线/离线状态切换

### 管理后台

- 仪表盘（核心指标）
- 用户/陪玩师/订单管理
- 财务（结算、佣金、提现）
- 内容审核（聊天监控、敏感词）
- 营销（活动、优惠券、VIP、推荐）
- RBAC 权限系统
- 实时监控 + 数据分析 + KPI

---

## API 文档

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

---

## 测试

**后端测试：**

```bash
cd api
go test ./... -v -cover
# 159 个测试文件，覆盖率目标 70%+
```

**管理后台测试：**

```bash
cd admin
npm run test
# 88 个单元测试
```

---

## CI/CD

项目配置了 4 条 GitHub Actions 流水线：

| 流水线 | 说明 |
|--------|------|
| `ci.yml` | 主 CI（lint、test、build、Docker 镜像） |
| `deploy.yml` | 自动部署 |
| `pre-commit-check.yml` | 提交前检查 |
| `security.yml` | 安全扫描 |

---

## 文档

| 文档 | 路径 | 说明 |
|------|------|------|
| PRD | [docs/PRD.md](docs/PRD.md) | 产品需求文档 |
| 进度 | [docs/PROGRESS.md](docs/PROGRESS.md) | 项目进度与版本历史 |
| 组件文档 | [app/docs/COMPONENTS.md](app/docs/COMPONENTS.md) | 小程序组件职责 |
| 页面计划 | [app/docs/PAGE_COMPONENT_USAGE_PLAN.md](app/docs/PAGE_COMPONENT_USAGE_PLAN.md) | 页面组件使用清单 |
| 重构计划 | [app/docs/REFACTOR_PLAN.md](app/docs/REFACTOR_PLAN.md) | 前端重构记录 |

---

## License

Private - All Rights Reserved
