# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 📋 项目概述

GameLink是一个现代化的陪玩管理平台，专注于为游戏陪玩服务提供高效的订单分发、用户管理和打手管理功能。平台采用Go语言后端+React前端的架构，支持高并发、低延迟的业务场景。

### 🎯 核心目标
- **订单智能分发**: 基于算法的智能订单匹配系统
- **实时通信管理**: WebSocket实现的实时状态同步
- **多端用户支持**: 用户端、打手端、管理端三端协同
- **安全支付保障**: 多渠道支付集成和风控系统
- **高并发处理**: 支持万级用户同时在线

## 🏗 技术架构

### 系统架构概览
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

### 技术栈详情

#### 后端技术栈
- **运行环境**: Go 1.24.0
- **Web框架**: Gin + GORM
- **数据库**:
  - SQLite (开发环境)
  - PostgreSQL (生产环境)
  - Redis 7.0 (缓存)
- **认证**: JWT (golang-jwt/jwt/v5)
- **配置管理**: YAML配置文件
- **构建工具**: Go Modules + Makefile

#### 前端技术栈
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **UI组件**: Arco Design
- **状态管理**: React Context (AuthContext, ThemeContext)
- **路由**: React Router v6
- **HTTP客户端**: 原生Fetch API
- **样式**: Less + CSS-in-JS
- **测试**: Vitest + Testing Library

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
│   ├── package.json        # NPM依赖配置
│   ├── tsconfig.json       # TypeScript配置
│   ├── vite.config.ts      # Vite构建配置
│   └── .eslintrc.cjs       # ESLint配置
├── configs/                # 全局配置文件
│   ├── config.development.yaml  # 开发环境配置
│   └── config.production.yaml   # 生产环境配置
├── docs/                   # 项目文档
├── scripts/                # 构建脚本
├── .gitignore              # Git忽略文件
├── README.md               # 项目说明
├── CONTRIBUTING.md         # 贡献指南
├── AGENTS.md               # AI开发指南
└── optimization_guide.md   # 性能优化指南
```

## 🚀 常用开发命令

### 后端开发命令

在 `backend/` 目录下执行：

```powershell
# 安装依赖
make deps

# 代码检查
make lint

# 运行测试
make test

# 启动用户服务 (开发模式)
make run CMD=user-service

# 构建所有服务
make build

# 手动运行服务
go run ./cmd/user-service
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

# 运行测试(单次)
npm run test:run
```

## 🗄️ 数据库结构和模型

### 数据模型层次结构

#### 基础模型 (Base)
```go
type Base struct {
    ID        uint64         `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

#### 核心业务模型

1. **用户模型 (User)**
   ```go
   type User struct {
       Base
       Phone        string     `json:"phone,omitempty" gorm:"size:32;uniqueIndex"`
       Email        string     `json:"email,omitempty" gorm:"size:128;uniqueIndex"`
       PasswordHash string     `json:"-" gorm:"size:255"`
       Name         string     `json:"name" gorm:"size:64"`
       AvatarURL    string     `json:"avatar_url,omitempty" gorm:"size:255"`
       Role         Role       `json:"role" gorm:"size:32"`         // user/player/admin
       Status       UserStatus `json:"status" gorm:"size:32;index"`  // active/suspended/banned
       LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
   }
   ```

2. **订单模型 (Order)**
   ```go
   type Order struct {
       Base
       UserID         uint64      `json:"user_id" gorm:"index"`
       PlayerID       uint64      `json:"player_id,omitempty" gorm:"index"`
       GameID         uint64      `json:"game_id" gorm:"index"`
       Title          string      `json:"title" gorm:"size:128"`
       Description    string      `json:"description,omitempty" gorm:"type:text"`
       Status         OrderStatus `json:"status" gorm:"size:32;index"`  // pending/confirmed/in_progress/completed/canceled/refunded
       PriceCents     int64       `json:"price_cents"`
       Currency       Currency    `json:"currency,omitempty" gorm:"type:char(3)"`
       ScheduledStart *time.Time  `json:"scheduled_start,omitempty"`
       ScheduledEnd   *time.Time  `json:"scheduled_end,omitempty"`
       CancelReason   string      `json:"cancel_reason,omitempty" gorm:"type:text"`
   }
   ```

3. **其他核心模型**
   - **Game**: 游戏信息
   - **Player**: 打手信息
   - **Payment**: 支付记录
   - **Review**: 评价系统

### 数据库迁移

数据库迁移通过GORM AutoMigrate处理：
```go
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.Game{},
        &model.Player{},
        &model.PlayerGame{},
        &model.PlayerSkillTag{},
        &model.User{},
        &model.Order{},
        &model.Payment{},
        &model.Review{},
    )
}
```

## 🔌 API设计模式

### RESTful API规范

#### API基础路径
- 基础路径: `/api/v1`
- 管理端路径: `/api/v1/admin`

#### 统一响应格式
```typescript
interface SuccessResponse<T> {
  success: true;
  data: T;
  message?: string;
}

interface ErrorResponse {
  success: false;
  code: number;
  message: string;
  details?: any;
}
```

#### 主要API端点

**认证相关**
```
POST   /api/v1/auth/login     # 用户登录
GET    /api/v1/auth/me        # 获取当前用户信息
POST   /api/v1/auth/logout    # 用户登出
```

**管理端 - 用户管理**
```
GET    /api/v1/admin/users           # 获取用户列表
POST   /api/v1/admin/users           # 创建用户
GET    /api/v1/admin/users/:id       # 获取用户详情
PUT    /api/v1/admin/users/:id       # 更新用户信息
DELETE /api/v1/admin/users/:id       # 删除用户
```

**管理端 - 订单管理**
```
GET    /api/v1/admin/orders          # 获取订单列表
GET    /api/v1/admin/orders/:id      # 获取订单详情
PUT    /api/v1/admin/orders/:id      # 更新订单状态
DELETE /api/v1/admin/orders/:id      # 删除订单
```

**管理端 - 游戏管理**
```
GET    /api/v1/admin/games           # 获取游戏列表
POST   /api/v1/admin/games           # 创建游戏
GET    /api/v1/admin/games/:id       # 获取游戏详情
PUT    /api/v1/admin/games/:id       # 更新游戏信息
DELETE /api/v1/admin/games/:id       # 删除游戏
```

#### 分页查询
所有列表API支持分页参数：
```
page: number          # 页码，从1开始
page_size: number     # 每页数量，默认20
sort_by: string       # 排序字段
sort_order: 'asc' | 'desc'  # 排序方向
keyword: string       # 搜索关键词
```

### 认证和授权

#### JWT认证
```go
// 中间件实现
func AdminAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"success": false, "code": 401, "message": "Authorization header required"})
            c.Abort()
            return
        }
        // 验证JWT token逻辑...
    }
}
```

#### 权限控制
- **用户角色**: user, player, admin
- **管理端权限**: 需要admin角色 + JWT认证
- **API限流**: 管理端API启用速率限制

## ⚙️ 开发环境配置

### 环境要求
- **Go**: 1.24.0+
- **Node.js**: 18+
- **PowerShell**: Windows 11环境
- **Git**: 版本控制

### 本地开发环境搭建

#### 1. 克隆项目
```powershell
git clone https://github.com/your-org/gamelink.git
cd gamelink
```

#### 2. 后端环境配置
```powershell
cd backend

# 安装Go依赖
go mod download

# 启动开发服务器
make run CMD=user-service
```

#### 3. 前端环境配置
```powershell
cd frontend

# 安装NPM依赖
npm install

# 启动开发服务器
npm run dev
```

### 配置文件说明

#### 后端配置 (configs/config.development.yaml)
```yaml
server:
  port: "8080"
  enable_swagger: true

database:
  type: "sqlite"
  dsn: "file:./var/dev.db?mode=rwc&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"

cache:
  type: "memory"
```

#### 前端配置 (src/config.ts)
```typescript
export const API_BASE = '/api/v1';
export const STORAGE_KEYS = {
  token: 'gamelink_token',
};
```

### 开发工具配置

#### VSCode配置推荐
- **Go扩展**: Go团队官方扩展
- **TypeScript扩展**: Microsoft TypeScript
- **ESLint扩展**: ESLint
- **Prettier扩展**: Prettier

#### Git配置
项目已配置完整的.gitignore，包括：
- Go构建产物
- Node.js依赖
- IDE配置文件
- 环境变量文件
- 数据库文件

## 🚀 部署流程

### 开发环境部署
```powershell
# 启动后端服务
cd backend
make run CMD=user-service

# 启动前端服务 (新终端)
cd frontend
npm run dev

# 访问应用
# 前端: http://localhost:5173
# 后端API: http://localhost:8080
# Swagger文档: http://localhost:8080/swagger/index.html
```

### 生产环境部署

#### 构建阶段
```powershell
# 构建后端
cd backend
make build

# 构建前端
cd frontend
npm run build
```

#### 配置管理
- 使用 `config.production.yaml` 生产配置
- 环境变量注入敏感信息
- 数据库连接使用PostgreSQL
- 缓存使用Redis

## 📋 代码规范和工具配置

### Go代码规范

#### 代码格式化
```powershell
# 格式化代码
go fmt ./...
goimports -w .

# 代码检查
make lint  # 使用golangci-lint
```

#### 命名规范
- **包名**: 小写，简短，有意义
- **常量**: UpperCamelCase 或 SCREAMING_SNAKE_CASE
- **变量**: lowerCamelCase
- **函数**: UpperCamelCase (导出) 或 lowerCamelCase (私有)

#### 注释规范
- 所有导出的函数、类型、常量必须有注释
- 使用JSDoc风格的注释
- 注释应说明函数的用途、参数、返回值

### TypeScript代码规范

#### ESLint配置
```javascript
// .eslintrc.cjs
module.exports = {
  extends: [
    '@typescript-eslint/recommended',
    'plugin:react/recommended',
    'prettier'
  ],
  rules: {
    // 自定义规则
  }
};
```

#### Prettier配置
```json
{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 100,
  "tabWidth": 2
}
```

#### TypeScript规范
- 使用严格的TypeScript配置
- 所有函数必须有参数和返回值类型
- 使用interface定义对象类型
- 避免使用any类型

### 提交规范

#### Conventional Commits
```
<type>(<scope>): <subject>

feat(user): add user registration feature
fix(order): resolve order status update issue
docs(api): update payment API documentation
```

#### 提交类型
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建工具、依赖更新等

## 🧪 测试策略

### 后端测试

#### 单元测试
```powershell
# 运行所有测试
make test

# 运行特定测试
go test ./internal/service

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### 测试文件组织
- 单元测试: `*_test.go`
- 集成测试: 使用 `-tags=integration`
- 测试数据: `tests/fixtures/`

### 前端测试

#### 测试命令
```powershell
# 运行测试
npm run test

# 运行测试并生成覆盖率
npm run test:run

# 监听模式
npm run test -- --watch
```

#### 测试类型
- **单元测试**: 组件测试
- **集成测试**: API集成测试
- **E2E测试**: 端到端测试

## 📊 性能优化指南

### 后端优化
- 参考 `optimization_guide.md`
- 启用Gin Release模式
- 优化数据库查询
- 使用缓存策略
- 批量操作优化

### 前端优化
- 使用Vite的构建优化
- 代码分割和懒加载
- 图片优化和压缩
- 缓存策略配置

## 🔍 监控和调试

### 日志管理
- 结构化日志输出
- 不同环境的日志级别
- 错误追踪和报警

### 性能监控
- API响应时间监控
- 数据库查询性能
- 内存使用情况

## 🤝 开发流程

### 分支管理
- `main`: 主分支，生产环境
- `develop`: 开发分支
- `feature/*`: 功能分支
- `hotfix/*`: 热修复分支

### 代码审查
- 所有代码需要PR审查
- 自动化测试必须通过
- 代码质量检查必须通过

### 发布流程
1. 代码合并到main分支
2. 自动化构建和测试
3. 部署到测试环境
4. 人工验证
5. 部署到生产环境

## 📚 文档资源

- **项目README**: `README.md`
- **贡献指南**: `CONTRIBUTING.md`
- **AI开发指南**: `AGENTS.md`
- **性能优化**: `optimization_guide.md`
- **API文档**: Swagger UI (`/swagger/index.html`)

## 🆘 故障排除

### 常见问题

#### 后端启动失败
- 检查Go版本是否为1.24.0+
- 确认依赖是否正确安装: `go mod download`
- 检查配置文件路径是否正确

#### 前端构建失败
- 清除node_modules重新安装: `rm -rf node_modules && npm install`
- 检查Node.js版本是否为18+
- 检查TypeScript配置

#### 数据库连接失败
- 确认数据库配置正确
- 检查数据库服务是否启动
- 验证连接字符串格式

### 获取帮助
- 查看项目Issues页面
- 联系开发团队
- 查看相关文档

## 🎯 Claude Code 工作规范

### 代码质量检查流程
作为产品经理、测试工程师和文档撰写人员，在每个开发阶段都需要进行代码整洁度检查：

#### 检查要点
1. **代码规范性**
   - 命名是否清晰、一致
   - 函数/类/组件职责是否单一
   - 是否遵循项目的编码规范

2. **代码逻辑性**
   - 业务逻辑是否清晰
   - 是否存在冗余代码
   - 错误处理是否完善

3. **代码可维护性**
   - 是否易于理解和修改
   - 注释是否充分
   - 是否存在硬编码

4. **性能和安全**
   - 是否存在性能问题
   - 安全性检查
   - 边界情况处理

#### 问题提出规范
发现问题后，请按以下格式提出：
```
🔍 **问题发现**: [问题描述]
📍 **位置**: [文件路径:行号]
💡 **建议**: [改进建议]
⚠️ **优先级**: [高/中/低]
```

#### 检查阶段
- **开发前**: 代码结构设计检查
- **开发中**: 代码质量和规范检查
- **开发后**: 功能完整性和测试覆盖检查
- **发布前**: 整体代码质量评估

---

**注意**: 本文档会随着项目的发展持续更新，请定期查看最新版本。