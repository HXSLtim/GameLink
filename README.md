# 🎮 GameLink - 陪玩管理平台

<div align="center">

![GameLink Logo](https://img.shields.io/badge/GameLink-陪玩平台-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**高性能陪玩订单分发和用户管理平台**

[功能特性](#-功能特性) • [技术架构](#-技术架构) • [快速开始](#-快速开始) • [部署指南](#-部署指南) • [API文档](#-api文档)

</div>

## 📋 项目概述

GameLink是一个现代化的陪玩管理平台，专注于为游戏陪玩服务提供高效的订单分发、用户管理和打手管理功能。平台采用Go语言微服务架构，支持高并发、低延迟的业务场景。

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
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 7.0+

### 本地开发环境搭建

1. **克隆项目**
```bash
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

2. **后端服务启动**
```bash
# 安装依赖
cd backend
go mod download

# 启动基础服务 (MySQL, Redis, MongoDB)
docker-compose up -d

# 运行数据库迁移
go run cmd/migrate/main.go

# 启动用户服务
go run cmd/user-service/main.go

# 启动订单服务
go run cmd/order-service/main.go

# 启动API网关
go run cmd/gateway/main.go
```

3. **前端应用启动**
```bash
# 用户端应用
cd frontend/user-app
npm install
npm run dev

# 管理端应用
cd frontend/admin-app
npm install
npm run dev
```

4. **验证安装**
```bash
# 检查服务状态
curl http://localhost:8080/health

# 查看API文档
open http://localhost:8080/swagger/index.html
```

## 📁 项目结构

```
GameLink/
├── backend/                 # Go后端服务
│   ├── cmd/                # 应用入口
│   │   ├── user-service/   # 用户服务
│   │   ├── order-service/  # 订单服务
│   │   ├── payment-service/# 支付服务
│   │   └── gateway/        # API网关
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置管理
│   │   ├── handler/        # HTTP处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── repository/     # 数据访问层
│   │   ├── model/          # 数据模型
│   │   └── middleware/     # 中间件
│   ├── pkg/                # 公共包
│   │   ├── database/       # 数据库连接
│   │   ├── cache/          # 缓存封装
│   │   ├── logger/         # 日志工具
│   │   └── utils/          # 工具函数
│   ├── api/                # API定义
│   ├── docs/               # 文档
│   ├── scripts/            # 脚本文件
│   └── configs/            # 配置文件
├── frontend/               # 前端应用
│   ├── user-app/           # 用户端应用
│   ├── player-app/         # 打手端应用
│   └── admin-app/          # 管理端应用
├── deployments/            # 部署配置
│   ├── docker/             # Docker配置
│   ├── k8s/                # Kubernetes配置
│   └── helm/               # Helm Charts
├── docs/                   # 项目文档
├── scripts/                # 构建脚本
└── tools/                  # 开发工具
```

## 🔧 开发指南

### 代码规范
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `golangci-lint` 进行代码检查
- 函数和方法必须有注释
- 单元测试覆盖率 > 80%

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
- Swagger UI: http://localhost:8080/swagger/index.html
- API文档: [docs/api.md](docs/api.md)

### 主要API端点
```
用户服务:
POST   /api/v1/auth/register     # 用户注册
POST   /api/v1/auth/login        # 用户登录
GET    /api/v1/users/profile     # 获取用户信息

订单服务:
POST   /api/v1/orders            # 创建订单
GET    /api/v1/orders            # 获取订单列表
PUT    /api/v1/orders/{id}       # 更新订单状态

支付服务:
POST   /api/v1/payments/wechat   # 微信支付
POST   /api/v1/payments/alipay   # 支付宝支付
```

## 🧪 测试

### 运行测试
```bash
# 单元测试
go test ./...

# 集成测试
go test -tags=integration ./...

# 性能测试
go test -bench=. ./...

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试环境
- 单元测试覆盖率要求: 80%+
- 集成测试环境: 独立的测试数据库
- 性能测试: 模拟真实负载场景

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