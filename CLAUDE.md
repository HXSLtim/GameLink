# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GameLink是一个现代化的陪玩管理平台，采用Go语言后端+React前端的架构。项目专注于为游戏陪玩服务提供高效的订单分发、用户管理和打手管理功能。

## 常用开发命令

### 后端开发命令

在 `backend/` 目录下执行：

```powershell
# 安装依赖
make deps
go mod tidy

# 代码检查
make lint
golangci-lint run --timeout=5m

# 运行测试
make test
go test ./...

# 启动用户服务 (开发模式)
make run CMD=user-service
go run ./cmd/user-service

# 构建所有服务
make build
go build ./cmd/...

# 生成Swagger文档
make swagger
```

### 前端开发命令

在 `frontend/` 目录下执行：

```powershell
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
npm run build:analyze  # 带分析报告

# 预览构建结果
npm run preview

# 代码检查
npm run lint

# 代码格式化
npm run format

# 类型检查
npm run typecheck

# 运行测试
npm run test
npm run test:run
npm run test:coverage
```

## 项目架构

### 后端架构 (Go 1.25.3)

**技术栈:**
- Web框架: Gin + GORM
- 数据库: SQLite (开发) / PostgreSQL (生产) / Redis
- 认证: JWT (golang-jwt/jwt/v5)
- 文档: Swagger (swaggo)
- 测试: testify + golang/mock

**目录结构:**
```
backend/
├── cmd/user-service/          # 应用入口点
├── internal/                  # 内部包
│   ├── admin/                 # 管理端处理器
│   ├── handler/               # HTTP处理器
│   ├── service/               # 业务逻辑层
│   ├── repository/            # 数据访问层
│   ├── model/                 # 数据模型
│   ├── auth/                  # JWT认证
│   ├── cache/                 # 缓存层
│   ├── config/                # 配置管理
│   └── middleware/            # 中间件
├── configs/                   # 配置文件
├── docs/swagger/              # API文档
└── scripts/sql/               # SQL脚本
```

### 前端架构 (React 18 + TypeScript)

**技术栈:**
- 框架: React 18 + TypeScript
- 构建工具: Vite
- 路由: React Router v6
- HTTP客户端: axios
- 样式: Less
- 测试: Vitest + Testing Library
- 代码检查: ESLint + Prettier

**目录结构:**
```
frontend/
├── src/
│   ├── api/                   # API调用层
│   ├── components/            # 可复用组件
│   ├── pages/                 # 页面组件
│   ├── layouts/               # 布局组件
│   ├── types/                 # TypeScript类型
│   └── utils/                 # 工具函数
├── public/                    # 静态资源
└── docs/                      # 前端文档
```

## 核心业务模型

### 数据模型 (GORM)
- **User**: 用户基础信息 (角色: user/player/admin)
- **Player**: 打手认证信息
- **Game**: 游戏配置
- **Order**: 订单管理 (状态: pending/confirmed/in_progress/completed/canceled)
- **Payment**: 支付记录
- **Review**: 评价系统

### API设计规范
- 基础路径: `/api/v1`
- 管理端路径: `/api/v1/admin`
- 认证方式: Bearer Token (JWT)
- 响应格式: 统一JSON格式
- 命名规范: camelCase

## 开发环境配置

### 环境要求
- Go: 1.25.3+
- Node.js: 18+
- PowerShell: Windows 11环境
- Git: 版本控制

### 快速启动

1. **后端服务**
```powershell
cd backend
make deps
make run CMD=user-service
```

2. **前端应用**
```powershell
cd frontend
npm install
npm run dev
```

3. **访问地址**
- 前端: http://localhost:5173
- 后端API: http://localhost:8080
- Swagger文档: http://localhost:8080/swagger/index.html

## 测试策略

### 后端测试
- 单元测试: `*_test.go` 文件
- 测试命令: `make test` 或 `go test ./...`
- 覆盖率要求: 80%+
- Mock: 使用golang/mock

### 前端测试
- 测试框架: Vitest
- 测试命令: `npm run test`
- 覆盖率: `npm run test:coverage`
- 组件测试: Testing Library

## 代码规范

### Go代码规范
- 遵循Go官方代码规范
- 使用golangci-lint进行检查
- 导出函数必须有JSDoc风格注释
- 命名: 包名小写，函数名大写开头，变量名小写开头

### TypeScript代码规范
- 严格TypeScript模式
- 函数必须有参数和返回值类型
- 使用interface定义对象类型
- 避免使用any类型

### 提交规范
使用Conventional Commits格式:
```
<type>(<scope>): <subject>

feat(user): add user registration feature
fix(order): resolve order status update issue
docs(api): update payment API documentation
```

## 常见问题排查

### 后端启动失败
1. 检查Go版本: `go version` (需要1.25.3+)
2. 安装依赖: `go mod download`
3. 检查配置文件路径
4. 查看端口占用情况

### 前端构建失败
1. 清除依赖重新安装: `rm -rf node_modules && npm install`
2. 检查Node.js版本: `node --version` (需要18+)
3. 检查TypeScript配置
4. 查看具体错误信息

### 数据库连接问题
1. 确认数据库配置正确
2. 检查数据库服务状态
3. 验证连接字符串格式

## Claude Code 角色和职责

### 🎯 我的角色
作为**项目测试和质量管理负责人**，我负责：
- 测试策略制定和执行
- 代码质量评估和审查
- 业务流程完整性验证
- 用户体验问题发现和改进建议
- 项目风险管理

### 🔍 质量检查要点
在开发过程中，我会重点关注：

#### 代码质量检查
1. **规范性检查**
   - 命名是否清晰、一致
   - 函数/类/组件职责是否单一
   - 是否遵循项目编码规范

2. **逻辑性检查**
   - 业务逻辑是否清晰合理
   - 是否存在冗余代码
   - 错误处理是否完善

3. **可维护性检查**
   - 代码是否易于理解和修改
   - 注释是否充分
   - 是否存在硬编码

4. **性能和安全检查**
   - 是否存在性能瓶颈
   - 安全性检查
   - 边界情况处理

#### 业务流程验证
1. **用户端体验**
   - 注册登录流程是否顺畅
   - 订单创建和支付流程是否完整
   - 用户界面是否友好

2. **管理端功能**
   - 用户管理是否完善
   - 订单监控是否实时
   - 数据统计是否准确

3. **陪玩师端功能**
   - 订单接收和处理流程
   - 收益统计和提现功能
   - 个人信息管理

#### 问题反馈规范
发现问题后，我会按以下格式提出：
```
🔍 **问题发现**: [问题描述]
📍 **位置**: [文件路径:行号]
💡 **建议**: [改进建议]
⚠️ **优先级**: [高/中/低]
🎯 **影响范围**: [功能模块/用户体验/系统稳定性]
```

### 📊 当前项目状态评估
**质量评级**: ⚠️ **需要改进** (存在阻塞性问题)

**主要问题**:
- 测试系统完全崩溃 (0% 可运行性)
- 数据模型不一致导致编译错误
- 前端组件测试失败率高
- 业务流程完整性不足

**改进重点**:
1. 立即修复测试系统和数据模型问题
2. 完善订单和支付业务流程
3. 提升用户界面和交互体验
4. 加强API安全性和一致性

### 🚀 开发阶段质量检查清单
- **开发前**: 需求分析和设计评审
- **开发中**: 代码质量和进度检查
- **测试后**: 功能完整性和性能验证
- **发布前**: 整体质量和风险评估

## 重要提醒

- 项目使用中文响应和注释
- 前后端API命名统一使用camelCase
- 所有新功能需要包含测试用例
- 遵循现有的错误处理模式
- 代码审查前请运行lint检查
- 查看docs/目录获取详细文档
- **质量优先**: 功能实现必须以保证代码质量为前提