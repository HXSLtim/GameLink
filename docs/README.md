# 🎮 GameLink - 陪玩管理平台

<div align="center">

![GameLink Logo](https://img.shields.io/badge/GameLink-陪玩平台-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?style=for-the-badge&logo=typescript)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**高性能陪玩订单分发和用户管理平台**

[功能特性](#-功能特性) • [技术架构](#-技术架构) • [快速开始](#-快速开始) • [部署指南](#-部署指南) • [API文档](#-api文档)

</div>

## 📋 项目概述

GameLink是一个现代化的陪玩管理平台，专注于为游戏陪玩服务提供高效的订单分发、用户管理和打手管理功能。平台采用Go语言后端+React前端的架构，支持高并发、低延迟的业务场景。目前已完成**camelCase命名规范统一**，前后端API接口完全一致。

### 🎯 核心目标
- **订单智能分发**: 基于算法的智能订单匹配系统
- **实时通信管理**: WebSocket实现的实时状态同步
- **多端用户支持**: 用户端、打手端、管理端三端协同
- **安全支付保障**: 多渠道支付集成和风控系统
- **高并发处理**: 支持万级用户同时在线

## ✨ 功能特性

### 👤 用户端功能
- 🔐 **安全认证**: 手机验证、第三方OAuth登录
- 🎮 **游戏管理**: 多游戏支持、段位选择、偏好设置
- 📝 **订单系统**: 订单创建、支付、跟踪、评价
- 💰 **钱包管理**: 充值、提现、积分兑换、优惠券
- 📊 **数据统计**: 个人数据、订单历史、消费分析

### 🎯 打手端功能
- ⭐ **认证体系**: 实名认证、技能认证、段位验证
- 📋 **订单大厅**: 智能订单推荐、批量接单、日程管理
- 💼 **收益管理**: 收益统计、预测分析、快速提现
- 📅 **智能日程**: 冲突检测、疲劳提醒、时间优化
- 🏆 **等级系统**: 技能等级、信用评级、成就体系

### 🛠 管理端功能
- 👥 **用户管理**: 用户审核、权限管理、黑名单
- 📈 **订单监控**: 实时监控、异常报警、数据分析
- 💳 **财务管理**: 资金流水、自动对账、风险控制
- 🎨 **系统配置**: 游戏配置、价格策略、规则设置
- 📊 **运营分析**: 用户分析、收益报表、趋势预测

## 🏗 技术架构

### 系统架构图
```
┌─────────────────────────────────────────────────────────┐
│                    前端应用层                           │
│    React Web + React Native + Admin Panel              │
└─────────────────────┬───────────────────────────────────┘
                      │ HTTPS + WebSocket
┌─────────────────────┴───────────────────────────────────┐
│                  API 网关层                             │
│         Kong Gateway + Custom Go Middleware            │
└─────────────────────┬───────────────────────────────────┘
                      │ Service Mesh (Istio)
┌─────────────────────┴───────────────────────────────────┐
│                  微服务层 (Go)                          │
│  用户服务 │ 订单服务 │ 支付服务 │ 通知服务 │ 游戏服务    │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────┴───────────────────────────────────┐
│                  数据存储层                             │
│   MySQL Cluster + Redis Cluster + MongoDB + MinIO      │
└─────────────────────────────────────────────────────────┘
```

### 技术栈

#### 后端技术栈
- **运行环境**: Go 1.21+
- **Web框架**: Gin + GORM
- **微服务**: Go-zero / go-kit
- **数据库**: MySQL 8.0 + Redis 7.0 + MongoDB
- **消息队列**: Redis + Asynq / Kafka
- **实时通信**: Gorilla WebSocket
- **服务发现**: Consul / etcd
- **监控追踪**: Prometheus + Jaeger + Zap
- **容器化**: Docker + Kubernetes

#### 前端技术栈
- **框架**: React 18 + TypeScript
- **状态管理**: Zustand + React Query
- **UI组件**: Ant Design + Tailwind CSS
- **构建工具**: Vite
- **实时通信**: Socket.io-client

## 🚀 快速开始

### 环境要求
- Go 1.24+
- Node.js 18+
- PowerShell (Windows 11)
- Git
- Docker & Docker Compose (可选)
- MySQL 8.0+ (生产环境)
- Redis 7.0+ (生产环境)

### 本地开发环境搭建

1. **克隆项目**
```bash
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

2. **后端服务启动**
```powershell
# 进入后端目录
cd backend

# 安装Go依赖
make deps

# 启动用户服务 (开发模式)
make run CMD=user-service

# 或者手动运行
go run ./cmd/user-service
```

3. **前端应用启动**
```powershell
# 进入前端目录
cd frontend

# 安装NPM依赖
npm install

# 启动开发服务器
npm run dev
```

4. **验证安装**
```powershell
# 检查后端API
curl http://localhost:8080/healthz

# 查看Swagger文档
# 浏览器访问: http://localhost:8080/swagger
```

## 📁 项目结构

```
GameLink/
├── backend/                 # Go后端服务
│   ├── cmd/                # 应用入口
│   │   └── user-service/   # 用户服务主程序
│   ├── internal/           # 内部包
│   │   ├── admin/          # 管理端处理器
│   │   ├── auth/           # 认证模块
│   │   ├── cache/          # 缓存层
│   │   ├── config/         # 配置管理
│   │   ├── db/             # 数据库连接
│   │   ├── handler/        # HTTP处理器
│   │   ├── model/          # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── service/        # 业务逻辑层
│   │   └── handler/middleware/ # 中间件
│   ├── scripts/            # 脚本文件
│   │   └── sql/            # SQL迁移脚本
│   ├── configs/            # 配置文件
│   ├── docs/               # 后端文档
│   ├── go.mod              # Go模块定义
│   └── Makefile            # 构建脚本
├── frontend/               # 前端应用
│   ├── src/                # 源代码
│   │   ├── api/            # API调用层
│   │   ├── components/     # 可复用组件
│   │   ├── contexts/       # React Context
│   │   ├── layouts/        # 布局组件
│   │   ├── pages/          # 页面组件
│   │   ├── services/       # 业务服务层
│   │   ├── types/          # TypeScript类型定义
│   │   └── utils/          # 工具函数
│   ├── public/             # 静态资源
│   ├── docs/               # 前端文档
│   ├── package.json        # NPM依赖配置
│   ├── tsconfig.json       # TypeScript配置
│   ├── vite.config.ts      # Vite构建配置
│   └── .eslintrc.cjs       # ESLint配置
├── configs/                # 全局配置文件
│   ├── config.development.yaml  # 开发环境配置
│   └── config.production.yaml   # 生产环境配置
├── docs/                   # 项目文档
│   ├── CAMELCASE_MIGRATION_REPORT.md # 迁移报告
│   ├── go-coding-standards.md     # Go编码规范
│   └── api-design-standards.md    # API设计规范
├── scripts/                # 构建脚本
├── .gitignore              # Git忽略文件
├── README.md               # 项目说明
├── CONTRIBUTING.md         # 贡献指南
├── AGENTS.md               # AI开发指南
├── CLAUDE.md               # Claude开发配置
└── optimization_guide.md   # 性能优化指南
```

## 🔧 开发指南

### 代码规范

#### Go代码规范
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `golangci-lint` 进行代码检查
- 所有导出的函数、类型、常量必须有注释
- 使用 JSDoc 风格的注释

#### TypeScript代码规范
- 使用严格的 TypeScript 配置
- 所有函数必须有参数和返回值类型
- 使用 interface 定义对象类型
- 避免使用 any 类型

#### 命名规范
- **API 接口**: 统一使用 camelCase 命名
- **数据库字段**: 保持 snake_case (GORM标签处理)
- **常量**: UpperCamelCase 或 SCREAMING_SNAKE_CASE
- **变量**: lowerCamelCase

### 测试要求
- 单元测试覆盖率 > 80%
- 前端组件测试覆盖
- API 集成测试
- E2E 测试覆盖核心流程

### 提交规范
```bash
# 提交格式
<type>(<scope>): <subject>

# 示例
feat(user): add user registration feature
fix(order): resolve order status update issue
docs(api): update payment API documentation
```

### 分支管理
- `main`: 主分支，用于生产环境
- `develop`: 开发分支，用于集成测试
- `feature/*`: 功能分支
- `hotfix/*`: 热修复分支

## 📊 性能指标

### 系统性能
- **并发用户**: 10,000+ 同时在线
- **订单处理**: 1,000+ TPS
- **API响应**: P99 < 100ms
- **系统可用性**: 99.95%+

### 资源占用
- **内存使用**: 50-100MB (单个服务)
- **CPU占用**: < 50% (正常负载)
- **启动时间**: < 1秒 (容器启动)

## 🚀 部署指南

### Docker部署
```bash
# 构建镜像
docker build -t gamelink/user-service:v1.0.0 .

# 运行容器
docker-compose up -d
```

### Kubernetes部署
```bash
# 部署到K8s集群
kubectl apply -f deployments/k8s/

# 检查部署状态
kubectl get pods -n gamelink
```

详细部署指南请参考 [部署文档](docs/deployment.md)

## 📚 API文档

### 在线文档
- **Swagger UI**: http://localhost:8080/swagger
- **API 文档**: http://localhost:8080/swagger.json
- **根路径**: http://localhost:8080/ (显示所有端点)

### API 规范
- **命名规范**: 统一 camelCase (已迁移完成)
- **认证方式**: Bearer Token (JWT)
- **响应格式**: 统一 JSON 格式
- **错误处理**: 标准化错误码和消息

### 主要API端点
```
认证相关:
POST   /api/v1/auth/login        # 用户登录
GET    /api/v1/auth/me           # 获取当前用户信息
POST   /api/v1/auth/logout       # 用户登出

管理端 - 用户管理:
GET    /api/v1/admin/users       # 获取用户列表
POST   /api/v1/admin/users       # 创建用户
GET    /api/v1/admin/users/{id}  # 获取用户详情
PUT    /api/v1/admin/users/{id}  # 更新用户信息
DELETE /api/v1/admin/users/{id}  # 删除用户

管理端 - 订单管理:
GET    /api/v1/admin/orders      # 获取订单列表
GET    /api/v1/admin/orders/{id} # 获取订单详情
PUT    /api/v1/admin/orders/{id} # 更新订单状态
DELETE /api/v1/admin/orders/{id} # 删除订单

管理端 - 游戏管理:
GET    /api/v1/admin/games       # 获取游戏列表
POST   /api/v1/admin/games       # 创建游戏
GET    /api/v1/admin/games/{id}  # 获取游戏详情
PUT    /api/v1/admin/games/{id}  # 更新游戏信息
DELETE /api/v1/admin/games/{id}  # 删除游戏
```

### 📖 详细文档
- [CamelCase 迁移报告](docs/CAMELCASE_MIGRATION_REPORT.md)
- [Go 编码规范](docs/go-coding-standards.md)
- [API 设计规范](docs/api-design-standards.md)

## 🧪 测试

### 后端测试
```powershell
# 在 backend/ 目录下执行

# 运行所有测试
make test

# 运行特定测试
go test ./internal/service

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试
```powershell
# 在 frontend/ 目录下执行

# 运行测试
npm run test

# 运行测试并生成覆盖率
npm run test:run

# 监听模式
npm run test -- --watch
```

### 测试环境
- **单元测试覆盖率要求**: 80%+
- **集成测试环境**: 独立的测试数据库
- **性能测试**: 模拟真实负载场景
- **E2E 测试**: 核心业务流程覆盖

## 📈 监控和日志

### 监控指标
- 应用性能指标 (APM)
- 基础设施监控
- 业务指标监控
- 错误率和延迟监控

### 日志管理
- 结构化日志 (JSON格式)
- 日志级别: DEBUG, INFO, WARN, ERROR
- 日志聚合: ELK Stack
- 日志保留: 30天

## 🤝 贡献指南

我们欢迎所有形式的贡献！请阅读 [贡献指南](CONTRIBUTING.md) 了解详细信息。

### 贡献流程
1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request
5. 代码审查
6. 合并代码

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 📞 联系我们

- **项目主页**: https://github.com/your-org/gamelink
- **问题反馈**: https://github.com/your-org/gamelink/issues
- **邮箱**: dev@gamelink.com
- **文档**: https://docs.gamelink.com

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者和社区成员！

---

<div align="center">

**[⬆ 回到顶部](#-gamelink---陪玩管理平台)**

Made with ❤️ by GameLink Team

</div>