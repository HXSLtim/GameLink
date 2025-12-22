# 项目结构

## 仓库布局

```
GameLink/
├── admin/                # 管理后台前端 (React)
├── app/                  # 小程序 (Taro)
├── backend/              # Go 后端服务
├── client/               # 用户端+陪玩师端前端 (待开发)
├── scripts/              # 部署和工具脚本
├── docs/                 # 项目文档
├── .kiro/steering/       # Kiro steering 规则
└── [root files]          # README, LICENSE, docker-compose 等
```

## 后端结构 (`backend/`)

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
│   ├── repository/                # 数据访问层
│   ├── model/                     # 数据模型（详见 04-data-models.md）
│   ├── router/                    # 路由定义
│   └── ws/                        # WebSocket 处理
├── pkg/                           # 公共/可复用包
│   ├── auth/                      # JWT 工具
│   ├── cache/                     # 缓存（memory/redis）
│   ├── config/                    # 配置管理
│   └── db/                        # 数据库工具
├── configs/                       # 配置文件
└── docs/                          # Swagger 文档
```

## 管理后台结构 (`admin/`)

```
admin/
├── src/
│   ├── api/                       # API 客户端
│   │   ├── auth.ts                # 认证 API
│   │   ├── admin.ts               # 管理 API
│   │   └── client.ts              # Axios 客户端（带加密）
│   ├── components/                # 可复用组件
│   ├── pages/                     # 页面组件
│   ├── layouts/                   # 布局包装器
│   ├── router/                    # 路由配置
│   ├── context/                   # React Context
│   ├── utils/                     # 工具函数
│   └── types/                     # TypeScript 类型
├── public/                        # 静态文件
└── dist/                          # 构建输出
```

## 命名规范

### 后端 (Go)

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件 | camelCase（小驼峰） | `userService.go`, `routingRule.go` |
| 包名 | lowercase | `handler` |
| 类型 | PascalCase | `UserService` |
| 导出函数 | PascalCase | `CreateUser` |
| 私有函数 | camelCase | `validateInput` |
| 变量 | camelCase | `userID` |
| 测试文件 | *_test.go | `userService_test.go` |

> ⚠️ **注意**: Go 文件使用 camelCase（小驼峰），不使用 snake_case

### 前端 (TypeScript/React)

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | PascalCase | `UserProfile.tsx` |
| 工具函数 | camelCase | `formatDate.ts` |
| 类型/接口 | PascalCase | `UserResponse` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL` |
| CSS 类 | kebab-case | `user-profile` |

## Import 组织

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
