# Project Structure

## Repository Layout

```
GameLink/
├── admin/                # 管理后台前端 (React)
├── app/                  # 小程序 (Taro)
├── backend/              # Go 后端服务
├── client/               # 用户端+陪玩师端前端 (待开发)
├── scripts/              # 部署和工具脚本
├── docs/                 # 项目文档
├── .kiro/                # Kiro 配置和 steering 规则
└── [root files]          # README, LICENSE, docker-compose 等
```

## Backend Structure (`backend/`)

```
backend/
├── cmd/main.go                    # 应用入口
├── internal/                      # 私有应用代码
│   ├── handler/                   # HTTP 处理器
│   │   ├── admin/                 # 管理端接口
│   │   ├── middleware/            # 中间件（auth, crypto, CORS）
│   │   ├── player/                # 陪玩师接口
│   │   └── user/                  # 用户接口
│   ├── service/                   # 业务逻辑层
│   │   ├── admin/                 # 管理服务
│   │   ├── auth/                  # 认证服务
│   │   ├── order/                 # 订单服务
│   │   └── player/                # 陪玩师服务
│   ├── repository/                # 数据访问层
│   │   ├── admin/                 # 管理数据访问
│   │   ├── interfaces/            # 仓库接口
│   │   └── mocks/                 # Mock 实现
│   ├── model/                     # 数据模型
│   │   ├── user.go                # 用户模型
│   │   ├── order.go               # 订单模型
│   │   ├── role.go                # RBAC 角色
│   │   └── permission.go          # RBAC 权限
│   ├── router/                    # 路由定义
│   ├── integration/               # 集成测试
│   └── ws/                        # WebSocket 处理
├── pkg/                           # 公共/可复用包
│   ├── auth/                      # JWT 工具
│   ├── cache/                     # 缓存（memory/redis）
│   ├── config/                    # 配置管理
│   ├── db/                        # 数据库工具（迁移）
│   └── logging/                   # 日志工具
├── configs/                       # 配置文件
│   ├── config.development.yaml
│   └── config.production.yaml
├── docs/                          # Swagger 文档
├── Makefile                       # 构建命令
└── Dockerfile                     # 容器镜像
```

## Admin Structure (`admin/`)

```
admin/
├── src/
│   ├── api/                       # API 客户端
│   │   ├── auth.ts                # 认证 API
│   │   ├── admin.ts               # 管理 API
│   │   └── client.ts              # Axios 客户端（带加密）
│   ├── components/                # 可复用组件
│   │   ├── common/                # 通用 UI 组件
│   │   └── layout/                # 布局组件
│   ├── pages/                     # 页面组件
│   │   ├── admin/                 # 管理端页面
│   │   ├── auth/                  # 认证页面
│   │   └── sys/                   # 系统页面
│   ├── layouts/                   # 布局包装器
│   │   └── AdminLayout/           # 管理端布局
│   ├── router/                    # 路由配置
│   │   ├── index.tsx              # 路由入口
│   │   ├── routes.tsx             # 静态路由
│   │   └── componentMap.tsx       # 组件映射
│   ├── context/                   # React Context
│   │   ├── AuthContext.tsx        # 认证上下文
│   │   ├── AdminContext.tsx       # 管理上下文
│   │   └── ThemeContext.tsx       # 主题上下文
│   ├── utils/                     # 工具函数
│   │   ├── crypto.ts              # 加密工具
│   │   ├── dynamicRoutes.tsx      # 动态路由生成
│   │   └── menuPermission.ts      # 菜单权限
│   ├── types/                     # TypeScript 类型
│   ├── App.tsx                    # 根组件
│   └── main.tsx                   # 应用入口
├── public/                        # 静态文件
├── dist/                          # 构建输出
├── package.json
├── vite.config.ts                 # Vite 配置（代码分割）
└── Dockerfile
```

## Scripts Structure (`scripts/`)

```
scripts/
├── deploy-production.ps1          # 标准部署脚本
├── deploy-production-encrypted.ps1 # 加密版部署脚本（推荐）
├── sync-crypto-keys.ps1           # 加密密钥同步
└── README.md                      # 脚本说明
```

## Naming Conventions

### Backend (Go)

- **Files**: `snake_case.go`
- **Packages**: lowercase, single word
- **Types**: PascalCase
- **Functions**: PascalCase (exported), camelCase (private)
- **Variables**: camelCase
- **Test files**: `*_test.go`

### Admin (TypeScript/React)

- **Components**: PascalCase (`UserProfile.tsx`)
- **Utilities**: camelCase (`formatDate.ts`)
- **Types/Interfaces**: PascalCase
- **Constants**: UPPER_SNAKE_CASE
- **CSS classes**: kebab-case

## Key Architectural Patterns

### Backend

1. **Dependency Injection**: Google Wire
2. **Repository Pattern**: 数据访问抽象
3. **Service Layer**: 业务逻辑分离
4. **Middleware Chain**: 请求处理管道
5. **Error Wrapping**: 上下文感知错误传播

### Admin

1. **Dynamic Routing**: 基于后端菜单的动态路由
2. **Component Mapping**: 组件名到组件的映射
3. **Context API**: 全局状态管理
4. **Lazy Loading**: 路由级代码分割
5. **Encrypted Communication**: API 请求自动加密

## Import Organization

### Go

```go
import (
    // 标准库
    "context"
    "fmt"

    // 第三方包
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // 内部包
    "gamelink/internal/model"
    "gamelink/pkg/auth"
)
```

### TypeScript

```typescript
// React 和第三方
import React from 'react';
import { Button } from 'antd';

// 内部模块
import { User } from '@/types/models';
import { useAuth } from '@/hooks/useAuth';
```
