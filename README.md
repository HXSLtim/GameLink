# 🎮 GameLink - 现代化游戏陪玩管理平台

[![Go Version](https://img.shields.io/badge/Go-1.25.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)](https://github.com/HXSLtim/GameLink/actions)
[![Coverage](https://img.shields.io/badge/Coverage-~80%25-green)](backend/LATEST_COVERAGE_REPORT.md)

**Go + React 全栈项目 | 智能订单分发 | 多角色管理 | 实时通讯**

---

## 🌟 项目简介

GameLink 是一个现代化的游戏陪玩管理平台，采用 Go 后端 + React 前端的架构，为游戏陪玩服务提供高效的订单分发、用户管理和陪玩师管理功能。

### 核心功能

- 🎯 **智能订单分发** - 自动匹配用户与陪玩师，支持抢单池和客服指派
- 👥 **多角色管理** - 用户/陪玩师/管理员权限体系 + RBAC 权限控制
- 💬 **实时通讯** - WebSocket 即时通讯，支持公共聊天室和订单群聊（不支持私聊）
- 💳 **完整支付** - 订单支付、退款、收益结算一体化
- 📊 **数据监控** - 实时订单状态、收益统计、系统指标
- 🔐 **安全认证** - JWT + RBAC 权限控制 + AES-256-CBC 通信加密

---

## 📊 项目状态

| 指标 | 当前值 | 目标 |
|------|--------|------|
| 后端模块完成度 | 100%（36/36 模块） | 100% ✅ |
| 前端完成度 | 75% | 100% |
| 测试覆盖率 | ~80% | 80%+ |
| CI/CD | ✅ 完善 | - |

### 模块实现状态

| 分类 | 完成 | 进行中 | 仅Model | 说明 |
|------|------|--------|---------|------|
| 核心模块 | 19 | 0 | 0 | user/order/player/chat/dispute 等 |
| 新增业务模块 | 4 | 0 | 0 | player-rank/order-timeout/user-block/vip |
| 营销模块 | 6 | 0 | 0 | vip/coupon/recharge/activity/team/referral |
| 辅助模块 | 7 | 0 | 0 | commission/ranking/routing-rule/operation-log 等 |

---

## 🚀 快速开始

### 🐳 方式一：Docker 部署（推荐）

**环境要求**: Docker 20.10+ 和 Docker Compose 2.0+

#### 开发环境
```powershell
# Windows PowerShell
.\scripts\docker-dev-start.ps1

# 或手动启动
docker-compose up -d
```

#### 生产环境部署
```powershell
# 1. 配置环境变量
Copy-Item .env.example .env
notepad .env  # 编辑配置

# 2. 启动服务（加密版，推荐）
.\scripts\deploy-production-encrypted.ps1
```

**访问地址**:
- 🌐 **前端应用**: http://localhost
- 🔌 **后端API**: http://localhost:8080
- 📚 **Swagger文档**: http://localhost:8080/swagger/index.html
- 👤 **默认管理员**: admin@gamelink.com / admin123456

### 💻 方式二：本地开发

**环境要求**: Go 1.25.5+, Node.js 18+, PostgreSQL 16+, Redis 7+

#### 后端服务
```bash
cd backend
go mod download
go run cmd/main.go
```

#### 前端应用
```bash
cd admin
npm install
npm run dev
```

---

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端应用       │    │   后端API服务    │    │   数据存储       │
│                │    │                │    │                │
│ • React 18     │◄──►│ • Go 1.25.5    │◄──►│ • PostgreSQL   │
│ • TypeScript   │    │ • Gin + GORM   │    │ • Redis        │
│ • Ant Design   │    │ • JWT Auth     │    │                │
│ • WebSocket    │    │ • Swagger API  │    │                │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25.5+, Gin, GORM, JWT, WebSocket, Wire |
| 前端 | React 18, TypeScript, Vite, Ant Design 6.0 |
| 数据库 | PostgreSQL 16+, Redis 7+ |
| 部署 | Docker, Nginx, GitHub Actions |
| 安全 | AES-256-CBC + SHA-256, JWT + RBAC |

### 后端分层架构

```
Handler → Service → Repository → Model
```

- **Handler**: HTTP 请求处理、参数验证、响应格式化
- **Service**: 业务逻辑、事务管理、跨模块协调
- **Repository**: 数据库操作、缓存、查询封装
- **Model**: 数据结构、数据库映射、验证规则

---

## 📁 项目结构

```
GameLink/
├── admin/                  # 管理后台前端 (React)
├── api/                    # Go 后端服务
│   ├── cmd/               # 应用入口
│   ├── internal/          # 内部模块
│   │   ├── handler/       # HTTP 处理器
│   │   ├── service/       # 业务逻辑层
│   │   ├── repository/    # 数据访问层
│   │   └── model/         # 数据模型
│   ├── pkg/               # 公共包
│   └── docs/              # API 文档
├── app/                    # 小程序 (Taro)
├── client/                 # 用户端前端 (待开发)
├── scripts/                # 部署脚本
├── docs/                   # 项目文档
└── .kiro/steering/         # Kiro steering 规则
```

---

## 🎯 功能特色

### 三端架构

| 端 | 主要功能 |
|------|----------|
| 用户端 | 首页、游戏列表、陪玩师浏览、订单创建、支付、评价 |
| 陪玩师端 | 工作台、订单管理、收益管理、服务管理、团队功能 |
| 管理后台 | 仪表盘、用户管理、订单监控、财务管理、系统设置 |

### 核心业务流程

```
用户下单 → 支付 → 等待接单 → 陪玩师接单 → 服务进行中 → 服务完成 → 用户评价
```

### 商业模式

- **平台抽成**: 15-25%（三维度计算：项目抽成 + 陪玩师个人抽成 + 上月排名调整）
- **系统统一定价**: ¥20-60+/小时（按段位，陪玩师不可自定义）
- **收入来源**: 服务佣金、认证费用、推广服务

---

## 📚 文档导航

### 开发文档

| 文档 | 说明 |
|------|------|
| [PROGRESS.md](PROGRESS.md) | 开发进度文档 |
| [docs/PRD.md](docs/PRD.md) | 产品需求文档 |
| [.kiro/steering/](/.kiro/steering/) | Steering 规则文档 |

### Steering 规则

| 文档 | 说明 |
|------|------|
| [01-product.md](.kiro/steering/01-product.md) | 产品概述 |
| [02-tech-stack.md](.kiro/steering/02-tech-stack.md) | 技术栈 |
| [03-project-structure.md](.kiro/steering/03-project-structure.md) | 项目结构 |
| [04-data-models.md](.kiro/steering/04-data-models.md) | 数据模型 |
| [05-testing-standard.md](.kiro/steering/05-testing-standard.md) | 测试规范 |
| [06-project-management.md](.kiro/steering/06-project-management.md) | 项目管理 |

---

## 🔧 常用命令

### 后端

```bash
cd api
go mod tidy          # 整理依赖
go run cmd/main.go   # 运行应用
make test            # 运行测试
make swagger         # 生成 API 文档
```

### 管理后台

```bash
cd admin
npm install          # 安装依赖
npm run dev          # 开发服务器 (localhost:5173)
npm run build        # 生产构建
```

### Docker

```powershell
# 查看所有命令
.\scripts\docker-manager.ps1 help

# 快速启动
.\scripts\docker-manager.ps1 start

# 查看服务状态
docker-compose ps
```

---

## 🛡️ 安全特性

| 特性 | 说明 |
|------|------|
| 通信加密 | AES-256-CBC + SHA-256 签名（生产环境强制启用） |
| 认证 | JWT Token + 刷新机制 |
| 权限控制 | RBAC 角色权限系统 |
| 数据保护 | 敏感数据加密存储 |

---

## 🤝 贡献指南

1. Fork 项目到你的 GitHub 账户
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交代码 (`git commit -m 'feat: add AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范

- 遵循 [Go 编码规范](https://golang.org/doc/effective_go.html)
- 遵循 [TypeScript 编码规范](https://www.typescriptlang.org/docs/)
- 添加必要的测试用例
- 更新相关文档
- 通过所有 CI 检查

---

## 📞 联系我们

- **项目负责人**: GameLink开发团队
- **技术支持**: a2778978136@163.com
- **项目仓库**: https://github.com/HXSLtim/GameLink.git

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个Star！**

**🚀 让我们一起构建更好的游戏陪玩生态！**

*最后更新: 2025-01-10 | 版本: v3.1*

</div>
