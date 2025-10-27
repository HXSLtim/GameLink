# GameLink 项目结构说明

## 📋 概述

本文档详细说明了 GameLink 陪玩管理平台的项目结构设计，包括后端Go微服务架构、前端应用结构和部署配置等。

## 🏗 整体项目结构

```
GameLink/
├── README.md                    # 项目主文档
├── LICENSE                      # 开源许可证
├── .gitignore                   # Git忽略文件
├── .github/                     # GitHub配置
│   ├── workflows/               # CI/CD工作流
│   │   ├── backend-ci.yml       # 后端CI
│   │   ├── frontend-ci.yml      # 前端CI
│   │   └── deploy.yml           # 部署工作流
│   ├── ISSUE_TEMPLATE/          # Issue模板
│   └── PULL_REQUEST_TEMPLATE.md # PR模板
├── docs/                        # 项目文档
│   ├── api/                     # API文档
│   ├── deployment/              # 部署文档
│   ├── development/             # 开发文档
│   └── architecture/            # 架构文档
├── backend/                     # Go后端服务
├── frontend/                    # 前端应用
├── deployments/                 # 部署配置
├── scripts/                     # 构建和部署脚本
├── tools/                       # 开发工具
├── configs/                     # 配置文件
└── tests/                       # 端到端测试
```

## 🔧 后端项目结构

```
backend/
├── go.mod                       # Go模块定义
├── go.sum                       # 依赖校验和
├── Makefile                     # 构建脚本
├── .golangci.yml               # 代码检查配置
├── cmd/                         # 应用入口点
│   ├── user-service/            # 用户服务
│   │   └── main.go
│   ├── order-service/           # 订单服务
│   │   └── main.go
│   ├── payment-service/         # 支付服务
│   │   └── main.go
│   ├── notification-service/    # 通知服务
│   │   └── main.go
│   ├── game-service/            # 游戏服务
│   │   └── main.go
│   ├── analytics-service/       # 统计分析服务
│   │   └── main.go
│   └── gateway/                 # API网关
│       └── main.go
├── internal/                    # 私有应用代码
│   ├── config/                  # 配置管理
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── redis.go
│   │   └── jwt.go
│   ├── handler/                 # HTTP处理器
│   │   ├── user/
│   │   │   ├── user_handler.go
│   │   │   ├── auth_handler.go
│   │   │   └── profile_handler.go
│   │   ├── order/
│   │   │   ├── order_handler.go
│   │   │   └── order_status_handler.go
│   │   ├── payment/
│   │   │   └── payment_handler.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       ├── cors.go
│   │       ├── logging.go
│   │       ├── rate_limit.go
│   │       └── recovery.go
│   ├── service/                 # 业务逻辑层
│   │   ├── user/
│   │   │   ├── user_service.go
│   │   │   ├── auth_service.go
│   │   │   └── profile_service.go
│   │   ├── order/
│   │   │   ├── order_service.go
│   │   │   ├── order_dispatcher.go
│   │   │   └── order_matcher.go
│   │   ├── payment/
│   │   │   ├── payment_service.go
│   │   │   ├── wechat_payment.go
│   │   │   └── alipay_payment.go
│   │   └── notification/
│   │       ├── notification_service.go
│   │       ├── websocket_hub.go
│   │       └── push_service.go
│   ├── repository/              # 数据访问层
│   │   ├── user/
│   │   │   ├── user_repository.go
│   │   │   └── user_repository_test.go
│   │   ├── order/
│   │   │   ├── order_repository.go
│   │   │   └── order_repository_test.go
│   │   ├── payment/
│   │   │   └── payment_repository.go
│   │   └── cache/
│   │       ├── redis_cache.go
│   │       └── local_cache.go
│   ├── model/                   # 数据模型
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   ├── user_request.go
│   │   │   └── user_response.go
│   │   ├── order/
│   │   │   ├── order.go
│   │   │   ├── order_request.go
│   │   │   └── order_response.go
│   │   ├── payment/
│   │   │   ├── payment.go
│   │   │   └── transaction.go
│   │   └── common/
│   │       ├── base_model.go
│   │       ├── response.go
│   │       └── error.go
│   ├── domain/                  # 领域模型
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   ├── order/
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   └── payment/
│   │       ├── entity.go
│   │       ├── service.go
│   │       └── repository.go
│   └── utils/                   # 工具函数
│       ├── validator.go
│       ├── password.go
│       ├── phone.go
│       ├── id_generator.go
│       ├── time.go
│       └── response.go
├── pkg/                         # 可被外部应用使用的库代码
│   ├── database/                # 数据库连接
│   │   ├── mysql.go
│   │   ├── redis.go
│   │   └── mongodb.go
│   ├── logger/                  # 日志工具
│   │   ├── logger.go
│   │   ├── zap.go
│   │   └── context.go
│   ├── cache/                   # 缓存封装
│   │   ├── redis_cache.go
│   │   ├── interface.go
│   │   └── local_cache.go
│   ├── auth/                    # 认证工具
│   │   ├── jwt.go
│   │   ├── oauth.go
│   │   └── password.go
│   ├── payment/                 # 支付工具
│   │   ├── wechat/
│   │   │   ├── client.go
│   │   │   └── types.go
│   │   └── alipay/
│   │       ├── client.go
│   │       └── types.go
│   ├── notification/            # 通知工具
│   │   ├── sms.go
│   │   ├── email.go
│   │   └── push.go
│   ├── storage/                 # 文件存储
│   │   ├── oss.go
│   │   ├── s3.go
│   │   └── interface.go
│   ├── middleware/              # 公共中间件
│   │   ├── tracing.go
│   │   ├── metrics.go
│   │   └── health_check.go
│   └── errors/                  # 错误处理
│       ├── errors.go
│       ├── codes.go
│       └── handler.go
├── api/                         # API定义和文档
│   ├── proto/                   # Protocol Buffers定义
│   │   ├── user/
│   │   │   └── user.proto
│   │   └── order/
│   │       └── order.proto
│   ├── openapi/                 # OpenAPI规范
│   │   ├── user-service.yaml
│   │   ├── order-service.yaml
│   │   └── payment-service.yaml
│   └── graphql/                 # GraphQL定义（可选）
│       ├── schema.graphql
│       └── resolvers/
├── migrations/                  # 数据库迁移文件
│   ├── 001_create_users_table.sql
│   ├── 002_create_orders_table.sql
│   ├── 003_create_payments_table.sql
│   └── 004_create_indexes.sql
├── configs/                     # 配置文件
│   ├── dev.yaml                 # 开发环境配置
│   ├── staging.yaml             # 测试环境配置
│   ├── prod.yaml                # 生产环境配置
│   └── local.yaml               # 本地环境配置
├── scripts/                     # 脚本文件
│   ├── build.sh                 # 构建脚本
│   ├── deploy.sh                # 部署脚本
│   ├── migrate.sh               # 数据库迁移脚本
│   └── seed.sh                  # 数据种子脚本
└── tests/                       # 测试文件
    ├── integration/             # 集成测试
    ├── e2e/                     # 端到端测试
    ├── performance/             # 性能测试
    └── fixtures/                # 测试数据
```

## 🎨 前端项目结构

```
frontend/
├── user-app/                    # 用户端应用
│   ├── public/                  # 静态资源
│   │   ├── index.html
│   │   ├── favicon.ico
│   │   └── manifest.json
│   ├── src/
│   │   ├── components/          # 可复用组件
│   │   │   ├── common/          # 通用组件
│   │   │   │   ├── Header/
│   │   │   │   ├── Footer/
│   │   │   │   ├── Loading/
│   │   │   │   └── Modal/
│   │   │   ├── forms/           # 表单组件
│   │   │   │   ├── LoginForm/
│   │   │   │   ├── RegisterForm/
│   │   │   │   └── OrderForm/
│   │   │   └── business/        # 业务组件
│   │   │       ├── GameSelector/
│   │   │       ├── OrderCard/
│   │   │       └── PlayerProfile/
│   │   ├── pages/               # 页面组件
│   │   │   ├── Home/
│   │   │   ├── Login/
│   │   │   ├── Register/
│   │   │   ├── Orders/
│   │   │   ├── Profile/
│   │   │   └── Wallet/
│   │   ├── hooks/               # 自定义Hooks
│   │   │   ├── useAuth.ts
│   │   │   ├── useOrder.ts
│   │   │   └── useWebSocket.ts
│   │   ├── store/               # 状态管理
│   │   │   ├── authStore.ts
│   │   │   ├── orderStore.ts
│   │   │   └── userStore.ts
│   │   ├── services/            # API服务
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   ├── order.ts
│   │   │   └── user.ts
│   │   ├── utils/               # 工具函数
│   │   │   ├── request.ts
│   │   │   ├── storage.ts
│   │   │   └── validation.ts
│   │   ├── types/               # TypeScript类型定义
│   │   │   ├── api.ts
│   │   │   ├── user.ts
│   │   │   └── order.ts
│   │   ├── styles/              # 样式文件
│   │   │   ├── globals.css
│   │   │   ├── variables.css
│   │   │   └── components.css
│   │   ├── App.tsx
│   │   ├── index.tsx
│   │   └── vite-env.d.ts
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── .eslintrc.js
├── player-app/                  # 打手端应用
│   └── [类似user-app结构]
├── admin-app/                   # 管理端应用
│   └── [类似user-app结构]
├── shared/                      # 共享代码
│   ├── components/              # 共享组件
│   ├── types/                   # 共享类型
│   ├── utils/                   # 共享工具
│   └── constants/               # 共享常量
├── build/                       # 构建输出
└── docs/                        # 前端文档
```

## 🚀 部署配置结构

```
deployments/
├── docker/                      # Docker配置
│   ├── backend/
│   │   ├── user-service/
│   │   │   └── Dockerfile
│   │   ├── order-service/
│   │   │   └── Dockerfile
│   │   └── ...
│   ├── frontend/
│   │   ├── user-app/
│   │   │   └── Dockerfile
│   │   └── ...
│   └── docker-compose.yml
├── kubernetes/                  # Kubernetes配置
│   ├── namespaces/
│   │   ├── gamelink-dev.yaml
│   │   ├── gamelink-staging.yaml
│   │   └── gamelink-prod.yaml
│   ├── configmaps/
│   │   ├── backend-config.yaml
│   │   └── frontend-config.yaml
│   ├── secrets/
│   │   ├── db-credentials.yaml
│   │   └── jwt-secret.yaml
│   ├── deployments/
│   │   ├── user-service.yaml
│   │   ├── order-service.yaml
│   │   ├── payment-service.yaml
│   │   └── ...
│   ├── services/
│   │   ├── user-service.yaml
│   │   ├── order-service.yaml
│   │   └── ...
│   ├── ingress/
│   │   ├── api-ingress.yaml
│   │   └── web-ingress.yaml
│   └── hpa/
│       ├── user-service-hpa.yaml
│       └── order-service-hpa.yaml
├── helm/                        # Helm Charts
│   ├── gamelink/
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   ├── values-dev.yaml
│   │   ├── values-staging.yaml
│   │   ├── values-prod.yaml
│   │   └── templates/
│   │       ├── deployment.yaml
│   │       ├── service.yaml
│   │       ├── configmap.yaml
│   │       └── ingress.yaml
│   └── dependencies/
│       ├── mysql/
│       ├── redis/
│       └── mongodb
├── terraform/                   # 基础设施即代码
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── modules/
│   │   ├── vpc/
│   │   ├── rds/
│   │   └── eks/
│   └── environments/
│       ├── dev/
│       ├── staging/
│       └── prod/
└── ansible/                     # 配置管理
    ├── playbooks/
    ├── roles/
    └── inventory/
```

## 🛠 开发工具结构

```
tools/
├── code-generator/              # 代码生成工具
│   ├── api-generator/
│   ├── model-generator/
│   └── service-generator/
├── migration-tool/              # 数据库迁移工具
├── performance-test/            # 性能测试工具
│   ├── load-test/
│   └── stress-test/
├── monitoring/                  # 监控工具
│   ├── prometheus-config/
│   ├── grafana-dashboards/
│   └── alertmanager-config/
└── scripts/                     # 辅助脚本
    ├── setup-dev-env.sh
    ├── clean-docker.sh
    └── backup-data.sh
```

## 📊 监控和日志结构

```
monitoring/
├── prometheus/
│   ├── prometheus.yml
│   ├── rules/
│   │   ├── api.yml
│   │   ├── business.yml
│   │   └── infrastructure.yml
│   └── targets/
├── grafana/
│   ├── dashboards/
│   │   ├── api-performance.json
│   │   ├── business-metrics.json
│   │   └── system-overview.json
│   └── provisioning/
│       ├── dashboards/
│       └── datasources/
├── alertmanager/
│   └── alertmanager.yml
├── loki/
│   └── loki.yml
└── jaeger/
    └── jaeger.yml
```

## 🔍 测试结构

```
tests/
├── unit/                        # 单元测试
│   ├── user-service/
│   ├── order-service/
│   └── payment-service/
├── integration/                 # 集成测试
│   ├── api-integration/
│   ├── database-integration/
│   └── cache-integration/
├── e2e/                         # 端到端测试
│   ├── user-journey/
│   ├── order-flow/
│   └── payment-flow/
├── performance/                 # 性能测试
│   ├── load-tests/
│   ├── stress-tests/
│   └── benchmark-tests/
├── security/                    # 安全测试
│   ├── penetration/
│   └── vulnerability/
└── fixtures/                    # 测试数据
    ├── users.json
    ├── orders.json
    └── payments.json
```

## 📝 文档结构

```
docs/
├── api/                         # API文档
│   ├── user-service.md
│   ├── order-service.md
│   ├── payment-service.md
│   └── websocket-api.md
├── architecture/                # 架构文档
│   ├── system-overview.md
│   ├── microservices.md
│   ├── database-design.md
│   └── security-design.md
├── deployment/                  # 部署文档
│   ├── local-setup.md
│   ├── docker-deployment.md
│   ├── kubernetes-deployment.md
│   └── production-deployment.md
├── development/                 # 开发文档
│   ├── getting-started.md
│   ├── coding-standards.md
│   ├── testing-guide.md
│   └── contribution-guide.md
└── operations/                  # 运维文档
    ├── monitoring.md
    ├── troubleshooting.md
    ├── backup-restore.md
    └── security-best-practices.md
```

## 🎯 关键设计原则

### 1. 分层架构
- **Handler层**: 处理HTTP请求，参数验证，响应格式化
- **Service层**: 业务逻辑处理，事务管理
- **Repository层**: 数据访问，缓存管理
- **Model层**: 数据模型定义

### 2. 依赖注入
- 使用Go的接口和依赖注入
- 便于测试和模块解耦
- 支持配置驱动的服务发现

### 3. 配置管理
- 环境隔离（dev/staging/prod）
- 配置热更新支持
- 敏感信息加密存储

### 4. 错误处理

## 🆕 前端管理端初始化（React + Arco Design）

- 位置：`frontend/`
- 技术栈：Vite + React + TypeScript + Arco Design
- 开发命令：
  - `npm run dev` 本地开发（已配置 `/api` 代理到 `http://localhost:8080`）
  - `npm run build` 产物构建
  - `npm run lint` 代码检查（ESLint + Prettier）
  - `npm run test` 单元测试（Vitest + Testing Library）
- 入口：`frontend/index.html`，应用入口 `frontend/src/main.tsx`
- 基础页面：Dashboard（总览）、Login（登录占位）、Settings（占位）
- 样式：引入 `@arco-design/web-react/dist/css/arco.css` 并提供全局样式 `src/styles/global.css`

- 统一的错误码和错误消息
- 结构化的错误响应
- 完善的错误日志记录

### 5. 可观测性
- 结构化日志记录
- 分布式链路追踪
- 业务指标监控
- 健康检查机制

这个项目结构设计遵循了Go语言的最佳实践，支持微服务架构，具有良好的可扩展性和可维护性。每个目录和文件都有明确的职责，便于团队协作和代码管理。
