# GameLink 新成员入职指南

> **版本**: v2.0
> **更新日期**: 2026-02-09
> **适用对象**: 新加入团队的开发工程师、产品经理、测试工程师

---

## 📚 目录

1. [欢迎加入 GameLink](#欢迎加入-gamelink)
2. [第一天清单](#第一天清单)
3. [项目概览](#项目概览)
4. [开发环境搭建](#开发环境搭建)
5. [代码仓库与分支管理](#代码仓库与分支管理)
6. [核心业务流程](#核心业务流程)
7. [开发工作流程](#开发工作流程)
8. [常用命令与工具](#常用命令与工具)
9. [问题排查与求助](#问题排查与求助)
10. [学习资源](#学习资源)

---

## 欢迎加入 GameLink

你好！欢迎加入 GameLink 团队。这是一份为你准备的新成员入职指南，帮助你快速了解项目、搭建开发环境并开始工作。

### 1.1 我们的产品

**GameLink** 是一款现代化游戏陪玩社交平台，连接游戏玩家与陪玩师，提供从下单、匹配、服务、评价到结算的完整商业闭环。

**核心价值**：
- 用户可以快速找到优质陪玩师，享受专业游戏服务
- 陪玩师可以接单赚钱，展示技能，建立个人品牌
- 平台通过交易佣金、VIP会员、增值服务实现商业价值

### 1.2 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                      客户端层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ 小程序/H5    │  │ 管理后台     │  │ 未来: APP       │   │
│  │ uni-app      │  │ React 19     │  │ uni-app 编译    │   │
│  │ Vue 3.4 + TS │  │ AntD 6 + TS │  │                 │   │
│  └──────┬───────┘  └──────┬───────┘  └─────────────────┘   │
│         │                 │                                  │
│         └─────────────────┴────────────────┐                │
│                   │ HTTPS / WSS           │                │
├───────────────────┼────────────────────────┼────────────────┤
│          Nginx 反向代理 + 负载均衡         │                │
├───────────────────┼────────────────────────┼────────────────┤
│                   │                        │                │
│  ┌────────────────┴────────────────────────┴────────┐     │
│  │                  API 层 (Go 1.24)                 │     │
│  │  Handler → Service → Repository → Model          │     │
│  │  WebSocket、Scheduler、Middleware                 │     │
│  └───────────────────────────────────────────────────┘     │
│                         │                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    数据层                             │   │
│  │  PostgreSQL 16+ (主存储) + Redis 7+ (缓存)         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 1.3 团队成员

| 角色 | 职责 |
|-----|------|
| **Backend-Lead** | 后端架构设计、API 开发 |
| **Frontend-Lead** | 管理后台开发、React 技术栈 |
| **Mobile-Lead** | 小程序/H5 开发、Vue 技术栈 |
| **Database-Architect** | 数据库设计、性能优化 |
| **DevOps-Engineer** | 部署、CI/CD、运维 |
| **Product-Manager** | 产品规划、需求分析 |

---

## 第一天清单

### ✅ 必做事项（第1天）

**上午（9:00-12:00）**

1. **环境准备**（30分钟）
   - [ ] 安装 Go 1.24+
   - [ ] 安装 Node.js 20+
   - [ ] 安装 Docker & Docker Compose
   - [ ] 安装 Git
   - [ ] 安装 IDE（推荐：VS Code / GoLand / WebStorm）

2. **代码克隆**（15分钟）
   ```bash
   git clone https://github.com/your-org/gamelink.git
   cd gamelink
   ```

3. **阅读文档**（1小时15分钟）
   - [ ] 阅读 `README.md` - 项目简介
   - [ ] 阅读 `PROJECT_OVERVIEW.md` - 项目概览
   - [ ] 阅读 `docs/PRD_COMPREHENSIVE.md` - 产品需求（可选）

**下午（14:00-18:00）**

4. **开发环境搭建**（2小时）
   - [ ] 启动数据库服务（Docker）
   - [ ] 配置环境变量（.env）
   - [ ] 启动后端服务
   - [ ] 启动管理后台
   - [ ] 启动小程序/H5

5. **熟悉项目**（2小时）
   - [ ] 运行后端测试
   - [ ] 访问管理后台，熟悉界面
   - [ ] 浏览代码仓库结构
   - [ ] 阅读 API 文档（Swagger）

**完成标志**

- [ ] 后端服务运行在 `http://localhost:8080`
- [ ] 管理后台运行在 `http://localhost:5173`
- [ ] H5 应用运行在 `http://localhost:5174`
- [ ] 可以使用超级管理员账号登录管理后台

---

## 项目概览

### 3.1 项目结构

```
GameLink/
├── api/                # Go 后端服务
│   ├── cmd/            # 应用入口
│   ├── internal/       # 内部代码
│   │   ├── handler/    # HTTP 处理器
│   │   ├── service/    # 业务逻辑（57个模块）
│   │   ├── repository/ # 数据访问（56个模块）
│   │   ├── model/      # 数据模型（67个）
│   │   └── router/     # 路由注册
│   ├── pkg/            # 公共包
│   ├── ws/             # WebSocket
│   └── go.mod
├── admin/              # 管理后台（React 19）
│   ├── src/
│   │   ├── pages/      # 40+ 页面
│   │   ├── components/ # 通用组件
│   │   ├── api/        # API 封装
│   │   └── router/     # 路由配置
│   └── package.json
├── app/                # 小程序/H5（uni-app）
│   ├── src/
│   │   ├── pages/      # 28 页面
│   │   ├── components/ # 133 组件
│   │   ├── composables/# 38 Hook
│   │   ├── api/        # 14 API 模块
│   │   └── store/      # Pinia Store
│   └── package.json
├── docs/               # 项目文档
│   ├── PRD.md
│   ├── PROGRESS.md
│   └── ...
├── scripts/            # 部署脚本
├── docker-compose.yml  # Docker 编排
├── .env.example        # 环境变量模板
└── README.md           # 项目说明
```

### 3.2 技术栈速查

#### 后端

| 技术 | 版本 | 用途 |
|-----|------|------|
| Go | 1.24+ | 开发语言 |
| Gin | 最新 | Web 框架 |
| GORM | 最新 | ORM 框架 |
| PostgreSQL | 16+ | 数据库 |
| Redis | 7+ | 缓存 |
| Gorilla WebSocket | - | 实时通讯 |

#### 管理后台

| 技术 | 版本 | 用途 |
|-----|------|------|
| React | 19 | UI 框架 |
| TypeScript | 5.9 | 类型系统 |
| Ant Design | 6 | 组件库 |
| Vite | 5 | 构建工具 |
| TanStack Query | - | 数据请求 |

#### 小程序/H5

| 技术 | 版本 | 用途 |
|-----|------|------|
| Vue | 3.4+ | UI 框架 |
| uni-app | 最新 | 跨平台框架 |
| TypeScript | 5+ | 类型系统 |
| Pinia | - | 状态管理 |

### 3.3 数据库概览

- **数据库类型**: PostgreSQL 16+
- **总表数**: 80+ 张
- **核心表**: users, players, orders, payments, reviews, chat_messages, wallets, coupons, activities, 等

---

## 开发环境搭建

### 4.1 安装必需软件

#### Go 安装

```bash
# macOS
brew install go

# Windows
# 下载安装包: https://go.dev/dl/

# Linux
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz

# 验证安装
go version
```

#### Node.js 安装

```bash
# 使用 nvm（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
nvm use 20

# 验证安装
node -v
npm -v
```

#### Docker 安装

```bash
# macOS
brew install docker docker-compose

# Windows
# 下载 Docker Desktop: https://www.docker.com/products/docker-desktop

# Linux
curl -fsSL https://get.docker.com | sh
```

### 4.2 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，设置以下关键配置：
# - 数据库配置
# - Redis 配置
# - JWT 密钥
# - 超级管理员账号
```

**开发环境最小配置**：

```env
# 应用环境
APP_ENV=development
GIN_MODE=debug

# 数据库
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=gamelink123
POSTGRES_DB=gamelink
POSTGRES_PORT=5433

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PORT=6380

# JWT（开发环境）
JWT_SECRET_KEY=dev-jwt-secret-key-minimum-32-chars
JWT_TOKEN_TTL_HOURS=24

# 超级管理员
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=Admin123456
SUPER_ADMIN_NAME=Super Admin

# 种子数据
SEED_ENABLED=true
```

### 4.3 启动数据库服务

```bash
# 启动 PostgreSQL 和 Redis
docker compose up -d

# 检查服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

### 4.4 启动后端服务

```bash
cd api

# 下载依赖
go mod download

# 运行数据库迁移
go run cmd/main.go migrate

# 初始化种子数据
go run cmd/main.go seed

# 启动服务
go run cmd/main.go

# 服务运行在 http://localhost:8080
# Swagger 文档: http://localhost:8080/swagger/index.html
```

### 4.5 启动管理后台

```bash
cd admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

### 4.6 启动小程序/H5

```bash
cd app

# 安装依赖
npm install

# H5 开发
npm run dev:h5
# 访问 http://localhost:5174

# 微信小程序
npm run dev:mp-weixin
# 用微信开发者工具打开 dist/dev/mp-weixin
```

### 4.7 验证安装

**后端验证**

```bash
# 访问 Swagger 文档
open http://localhost:8080/swagger/index.html

# 测试 API
curl http://localhost:8080/api/v1/public/health
```

**管理后台验证**

```bash
# 访问管理后台
open http://localhost:5173

# 使用超级管理员登录
账号: admin@gamelink.com
密码: Admin123456
```

**小程序/H5 验证**

```bash
# 访问 H5
open http://localhost:5174

# 应该能看到首页
```

---

## 代码仓库与分支管理

### 5.1 代码仓库

- **主仓库**: `https://github.com/your-org/gamelink.git`
- **类型**: Git 仓库
- **可见性**: Private

### 5.2 分支策略

```
main (生产环境)
  ↑
  dev (开发环境)
  ↑
  feature/* (功能分支)
  fix/* (修复分支)
  hotfix/* (紧急修复)
```

### 5.3 分支规范

| 分支类型 | 命名规范 | 说明 | 合并目标 |
|---------|---------|------|---------|
| 主分支 | `main` | 生产环境代码 | 受保护 |
| 开发分支 | `dev` | 开发环境代码 | 功能分支合并到此 |
| 功能分支 | `feature/功能名称` | 新功能开发 | `dev` |
| 修复分支 | `fix/问题描述` | Bug 修复 | `dev` |
| 紧急修复 | `hotfix/问题描述` | 生产环境紧急修复 | `main` 和 `dev` |

### 5.4 Git 工作流

#### 开始新功能

```bash
# 1. 切换到 dev 分支并更新
git checkout dev
git pull origin dev

# 2. 创建功能分支
git checkout -b feature/user-management

# 3. 开发并提交
git add .
git commit -m "feat(user): add user management page"

# 4. 推送到远程
git push origin feature/user-management

# 5. 创建 Pull Request 到 dev
```

#### 提交规范

遵循 Conventional Commits 规范：

```bash
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式（不影响功能）
refactor: 重构
test: 测试
chore: 构建/工具

# 示例
feat(user): add user profile page
fix(order): resolve payment timeout issue
docs(readme): update setup instructions
```

### 5.5 代码审查

所有代码变更都需要经过 Pull Request 和代码审查：

1. 创建 Pull Request
2. 填写 PR 模板（功能描述、测试情况、 breaking changes）
3. 请求团队成员审查
4. 至少 1 人批准后方可合并
5. 解决审查意见
6. 合并到目标分支

---

## 核心业务流程

### 6.1 用户端核心流程

#### 订单流程

```
用户浏览陪玩师
    ↓
选择服务、时长
    ↓
创建订单
    ↓
支付（余额/微信/支付宝）
    ↓
待接单 → 陪玩师接单
    ↓
待确认 → 双方确认开始
    ↓
服务中 → 服务完成
    ↓
待评价 → 用户评价
    ↓
已完成 → 订单归档
```

#### 支付流程

```
选择支付方式
    ↓
创建支付订单
    ↓
调用支付接口
    ↓
用户完成支付
    ↓
接收支付回调
    ↓
验证签名
    ↓
更新订单状态
    ↓
分账（平台佣金 + 陪玩师收入）
```

### 6.2 陪玩师端核心流程

#### 接单流程

```
查看可接订单
    ↓
接单/拒绝
    ↓
确认开始服务
    ↓
提供服务
    ↓
完成服务
    ↓
等待用户评价
    ↓
获得收入
```

#### 提现流程

```
查看收益
    ↓
申请提现
    ↓
管理员审核
    ↓
审核通过 → 打款
    ↓
提现完成
```

### 6.3 管理后台核心流程

#### 陪玩师认证流程

```
陪玩师提交认证
    ↓
上传资料（身份证/段位截图）
    ↓
管理员审核
    ↓
审核通过/拒绝
    ↓
通知陪玩师
```

#### 争议处理流程

```
用户/陪玩师发起争议
    ↓
保存聊天快照
    ↓
分配客服（独立客服）
    ↓
调解处理
    ↓
做出决定（退款/驳回）
    ↓
执行结果
    ↓
关闭争议
```

---

## 开发工作流程

### 7.1 开发流程

1. **需求分析**（产品经理）
   - 编写需求文档
   - 评审需求
   - 确定优先级

2. **技术设计**（技术负责人）
   - 编写技术设计文档
   - 评审技术方案
   - 确定实现方案

3. **开发任务**（开发工程师）
   - 创建功能分支
   - 编写代码
   - 编写单元测试
   - 自测通过

4. **代码审查**（团队成员）
   - 提交 Pull Request
   - 代码审查
   - 修改反馈
   - 审查通过

5. **测试验证**（测试工程师）
   - 功能测试
   - 集成测试
   - 回归测试
   - Bug 修复

6. **部署上线**（DevOps）
   - 合并到 dev
   - 测试环境验证
   - 合并到 main
   - 生产环境部署

### 7.2 代码规范

#### Go 代码规范

```go
// 1. 包命名：小写单词
package user

// 2. 导出函数：大写开头
func GetUserByID(id int64) (*User, error) {
    // ...
}

// 3. 错误处理
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// 4. 结构体
type User struct {
    ID       uint64 `json:"id" gorm:"primaryKey"`
    Name     string `json:"name" gorm:"size:64"`
    Email    string `json:"email" gorm:"size:128"`
}
```

#### TypeScript 代码规范

```typescript
// 1. 接口定义：PascalCase
interface User {
  id: number;
  name: string;
  email: string;
}

// 2. 组件：PascalCase
function UserProfile({ user }: { user: User }) {
  return <div>{user.name}</div>;
}

// 3. 常量：UPPER_CASE
const API_BASE_URL = 'http://localhost:8080';
```

#### Vue 代码规范

```vue
<!-- 1. 组件名：PascalCase -->
<template>
  <div class="user-profile">
    {{ userName }}
  </div>
</template>

<script setup lang="ts">
// 2. Composition API
import { ref, computed } from 'vue';

const userName = ref('John');
</script>

<style scoped>
/* 3. 样式：scoped */
.user-profile {
  color: #333;
}
</style>
```

### 7.3 测试规范

#### 单元测试

```go
// 文件命名：xxx_test.go
// 测试函数：TestXxx
func TestGetUserByID(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    // Act
    user, err := service.GetUserByID(1)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John", user.Name)
}
```

#### 前端测试

```typescript
// 文件命名：xxx.test.tsx
import { render, screen } from '@testing-library/react';
import UserProfile from './UserProfile';

test('renders user name', () => {
  render(<UserProfile user={{ name: 'John' }} />);
  expect(screen.getByText('John')).toBeInTheDocument();
});
```

---

## 常用命令与工具

### 8.1 后端常用命令

```bash
# 运行服务
go run cmd/main.go

# 运行测试
go test ./... -v -cover

# 运行特定测试
go test ./internal/service/user -v

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...

# 依赖管理
go mod tidy
go mod download

# 数据库迁移
go run cmd/main.go migrate
go run cmd/main.go seed
```

### 8.2 前端常用命令

#### 管理后台

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 运行测试
npm run test

# 代码检查
npm run lint

# 代码格式化
npm run format
```

#### 小程序/H5

```bash
# 安装依赖
npm install

# H5 开发
npm run dev:h5

# H5 构建
npm run build:h5

# 微信小程序
npm run dev:mp-weixin
npm run build:mp-weixin
```

### 8.3 Docker 常用命令

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose down

# 查看日志
docker compose logs -f

# 查看服务状态
docker compose ps

# 重启服务
docker compose restart

# 进入容器
docker compose exec postgres bash
docker compose exec redis sh
```

### 8.4 Git 常用命令

```bash
# 查看状态
git status

# 查看分支
git branch -a

# 切换分支
git checkout dev

# 创建并切换分支
git checkout -b feature/xxx

# 提交代码
git add .
git commit -m "feat: xxx"

# 推送代码
git push origin feature/xxx

# 拉取更新
git pull origin dev

# 合并分支
git merge feature/xxx
```

---

## 问题排查与求助

### 9.1 常见问题

#### 后端问题

**Q: 后端启动失败，提示数据库连接错误？**

A: 检查以下几点：
1. PostgreSQL 和 Redis 是否已启动：`docker compose ps`
2. `.env` 中的数据库配置是否正确
3. 数据库端口是否被占用：`lsof -i:5433`

**Q: 测试失败？**

A:
1. 确保数据库已迁移：`go run cmd/main.go migrate`
2. 确保种子数据已初始化：`go run cmd/main.go seed`
3. 检查环境变量配置

**Q: Swagger 文档无法访问？**

A:
1. 确保后端服务已启动
2. 检查端口 8080 是否被占用
3. 尝试访问 `http://localhost:8080/swagger/index.html`

#### 前端问题

**Q: 管理后台无法启动？**

A:
1. 删除 `node_modules` 和 `package-lock.json`
2. 重新安装依赖：`npm install`
3. 检查 Node.js 版本：`node -v`（应该是 20+）

**Q: 管理后台无法登录？**

A:
1. 确保后端服务已启动
2. 检查 `admin/src/api/client.ts` 中的 API 地址
3. 打开浏览器控制台查看错误信息

**Q: 小程序/H5 无法访问接口？**

A:
1. 检查后端是否启动
2. 检查 `app/src/api/request.ts` 中的 baseURL
3. 检查跨域配置（后端 CORS 中间件）

#### Docker 问题

**Q: Docker 容器无法启动？**

A:
1. 检查 Docker 是否运行：`docker ps`
2. 查看容器日志：`docker compose logs`
3. 尝试重建容器：`docker compose down && docker compose up -d`

### 9.2 获取帮助

#### 技术问题

1. **查看文档**
   - 项目 README：`README.md`
   - 项目概览：`PROJECT_OVERVIEW.md`
   - API 文档：`http://localhost:8080/swagger/index.html`

2. **搜索代码**
   - 使用 VS Code 的全局搜索功能
   - 使用 `grep` 或 `rg` 命令行工具

3. **询问团队**
   - 在团队群里提问
   - @ 相关的技术负责人
   - 创建 Issue 描述问题

#### 求助模板

```
【问题描述】
简洁描述遇到的问题

【复现步骤】
1. 执行了什么操作
2. 预期结果是什么
3. 实际结果是什么

【环境信息】
- 操作系统：
- Go/Node 版本：
- 错误日志：

【已尝试的解决方案】
1.
2.
```

---

## 学习资源

### 10.1 技术文档

#### Go 相关

- [Go 官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)

#### 前端相关

- [React 官方文档](https://react.dev/)
- [TypeScript 手册](https://www.typescriptlang.org/docs/)
- [Ant Design 文档](https://ant.design/components/overview-cn/)
- [Vue 3 文档](https://cn.vuejs.org/)
- [uni-app 文档](https://uniapp.dcloud.net.cn/)

#### 数据库相关

- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Redis 文档](https://redis.io/docs/)
- [SQL 教程](https://www.w3schools.com/sql/)

### 10.2 项目文档

| 文档 | 路径 | 说明 |
|-----|------|------|
| 项目 README | `/README.md` | 项目简介 |
| 项目概览 | `/PROJECT_OVERVIEW.md` | 综合概览 |
| 产品需求 | `/docs/PRD_COMPREHENSIVE.md` | PRD v2.0 |
| 项目进度 | `/docs/PROGRESS.md` | 版本历史 |
| 组件文档 | `/app/docs/COMPONENTS.md` | 小程序组件 |
| 重构计划 | `/app/docs/REFACTOR_PLAN.md` | 前端重构 |

### 10.3 推荐学习路径

#### 后端工程师

**第1周：环境搭建与熟悉**
- 安装开发环境
- 运行项目
- 阅读 API 文档
- 熟悉代码结构

**第2-3周：核心业务理解**
- 用户/陪玩师/订单模块
- 支付流程
- 聊天系统
- 权限系统

**第4周：实践开发**
- 修复简单 Bug
- 添加小功能
- 参与代码审查

#### 前端工程师（管理后台）

**第1周：环境搭建与熟悉**
- 安装开发环境
- 运行项目
- 熟悉 React + Ant Design
- 熟悉项目结构

**第2-3周：核心模块理解**
- 路由与权限系统
- 用户/陪玩师管理
- 订单管理
- 实时通信

**第4周：实践开发**
- 修复简单 Bug
- 添加小功能
- 参与代码审查

#### 前端工程师（小程序/H5）

**第1周：环境搭建与熟悉**
- 安装开发环境
- 运行项目
- 熟悉 Vue 3 + uni-app
- 熟悉项目结构

**第2-3周：核心模块理解**
- 页面路由
- 用户端功能（浏览、下单、支付）
- 聊天功能
- 个人中心

**第4周：实践开发**
- 修复简单 Bug
- 添加小功能
- 参与代码审查

### 10.4 社区资源

- [Go 语言中文网](https://studygolang.com/)
- [React 中文网](https://react.docschina.org/)
- [Vue 中文网](https://cn.vuejs.org/)
- [SegmentFault](https://segmentfault.com/)
- [Stack Overflow](https://stackoverflow.com/)

---

## 附录

### A. 默认账号

| 角色 | 账号 | 密码 | 说明 |
|-----|------|------|------|
| 超级管理员 | admin@gamelink.com | Admin123456 | 管理后台登录 |

### B. 核心端口

| 服务 | 端口 | 说明 |
|-----|------|------|
| 后端 API | 8080 | Go API 服务 |
| 管理后台 | 5173 | React 开发服务器 |
| H5 应用 | 5174 | Vue 开发服务器 |
| PostgreSQL | 5433 | 数据库 |
| Redis | 6380 | 缓存 |

### C. 有用的链接

- API 文档：`http://localhost:8080/swagger/index.html`
- 管理后台：`http://localhost:5173`
- H5 应用：`http://localhost:5174`

### D. 团队沟通

- **日常沟通**：团队群
- **代码审查**：GitHub Pull Request
- **文档协作**：GitHub Wiki
- **问题追踪**：GitHub Issues

---

**祝你入职顺利！** 🎉

如果在任何时候遇到问题，请不要犹豫，随时向团队求助。我们都在这里支持你！

**重要提示**：

1. **不要害怕提问** - 没有问题是愚蠢的
2. **先搜索再提问** - 很多问题可能已经有答案
3. **描述清楚问题** - 提供足够的上下文和错误信息
4. **保持耐心** - 学习新东西需要时间

欢迎加入 GameLink 团队！让我们一起打造优秀的游戏陪玩平台！🚀
