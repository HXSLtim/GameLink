# GameLink 快速参考手册

**版本**: 1.0
**更新日期**: 2026-02-09
**适用场景**: 日常开发查阅

---

## 目录

1. [常用命令](#1-常用命令)
2. [端口与服务](#2-端口与服务)
3. [默认账号](#3-默认账号)
4. [快速定位文件](#4-快速定位文件)
5. [常见代码模式](#5-常见代码模式)
6. [调试技巧](#6-调试技巧)
7. [故障排查](#7-故障排查)

---

## 1. 常用命令

### 后端开发

```bash
# 启动后端
cd api
go run cmd/main.go

# 运行测试
go test ./... -v

# 运行特定测试
go test ./internal/service/order -v

# 查看覆盖率
go test ./... -cover

# 数据库迁移
go run cmd/main.go migrate

# 填充种子数据
go run cmd/main.go seed

# 代码检查
golangci-lint run

# 格式化代码
go fmt ./...

# 查看依赖
go mod graph

# 更新依赖
go get -u ./...
```

### 管理后台开发

```bash
# 启动开发服务器
cd admin
npm run dev

# 生产构建
npm run build

# 预览构建
npm run preview

# 运行测试
npm run test

# 运行测试 UI
npm run test:ui

# 代码检查
npm run lint

# 代码格式化
npm run format

# 类型检查
npm run type-check
```

### 用户端 Web 开发

```bash
# 启动开发
cd app
npm run dev

# 生产构建
npm run build

# 预览构建
npm run preview
```

### Docker 命令

```bash
# 启动所有服务
docker-compose up -d

# 启动特定服务
docker-compose up -d postgres redis

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f api

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 进入容器
docker-compose exec postgres bash
docker-compose exec api sh

# 查看服务状态
docker-compose ps
```

### Git 命令

```bash
# 创建功能分支
git checkout -b feature/your-feature

# 提交代码
git add .
git commit -m "feat: add feature"

# 推送到远程
git push origin feature/your-feature

# 更新 dev 分支
git checkout dev
git pull origin dev

# 合并到 dev
git merge feature/your-feature

# 删除分支
git branch -d feature/your-feature
```

---

## 2. 端口与服务

### 开发环境端口

| 服务 | 端口 | 用途 |
|------|------|------|
| Go API | 8080 | 后端 API |
| Admin Dev | 5173 | 管理后台开发服务器 |
| User Dev | 5175 | 用户端 Web 开发服务器 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |
| Swagger | 8080/swagger | API 文档 |

### 生产环境端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 80, 443 | 反向代理 |
| Go API | 8080 | 后端服务 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |

### 访问地址

```
开发环境:
- 后端 API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- 管理后台: http://localhost:5173
- H5 应用: http://localhost:5174

生产环境:
- 主应用: https://www.gamelink.com
- 管理后台: https://admin.gamelink.com
- API: https://api.gamelink.com
```

---

## 3. 默认账号

### 管理员账号

```
账号: admin@gamelink.com
密码: Admin123456
权限: 超级管理员 (*)
```

### 测试账号

```
普通用户:
- 账号: user@test.com
- 密码: Test123456

陪玩师:
- 账号: player@test.com
- 密码: Test123456
```

### 创建新账号

```bash
# 方式 1: 使用种子数据
go run cmd/main.go seed

# 方式 2: 注册接口
POST /api/v1/auth/register
{
  "email": "newuser@test.com",
  "password": "Password123",
  "name": "New User",
  "role": "user"
}

# 方式 3: 直接数据库插入
INSERT INTO users (email, password_hash, name, role)
VALUES ('new@test.com', '$2a$10$...', 'New User', 'user');
```

---

## 4. 快速定位文件

### 后端文件定位

```bash
# 用户相关
api/internal/handler/user/          # 用户接口
api/internal/service/user/          # 用户服务
api/internal/repository/user/       # 用户数据
api/internal/model/user.go          # 用户模型

# 订单相关
api/internal/handler/user/order.go  # 订单接口
api/internal/service/order/         # 订单服务
api/internal/repository/order/      # 订单数据
api/internal/model/order.go         # 订单模型

# 聊天相关
api/internal/handler/user/chat.go   # 聊天接口
api/internal/service/chat/          # 聊天服务
api/internal/repository/chat/       # 聊天数据
api/internal/model/chat.go          # 聊天模型
ws/                                  # WebSocket 实现

# 认证相关
api/internal/handler/auth.go        # 认证接口
api/pkg/auth/jwt.go                 # JWT 处理
api/pkg/auth/wechat.go              # 微信认证
```

### 管理后台文件定位

```bash
# 页面组件
admin/src/pages/admin/              # 管理端页面
admin/src/pages/admin/Dashboard/    # 仪表盘
admin/src/pages/admin/User/         # 用户管理
admin/src/pages/admin/Order/        # 订单管理
admin/src/pages/admin/Player/       # 陪玩师管理

# 通用组件
admin/src/components/common/         # 通用组件
admin/src/components/PermissionGuard/ # 权限组件
admin/src/components/SearchTable/    # 搜索表格

# 路由配置
admin/src/router/index.tsx           # 路由入口
admin/src/router/routes.tsx          # 路由配置
admin/src/router/componentMap.tsx    # 组件映射

# API 封装
admin/src/api/admin.ts               # 管理端 API
admin/src/api/client.ts              # Axios 客户端

# 权限系统
admin/src/context/AdminContext.tsx   # 管理员上下文
admin/src/utils/menuPermission.ts    # 菜单权限
admin/src/config/adminRoutes.ts      # 路由配置
```

### 小程序文件定位

```bash
# 页面
app/src/pages/auth/login/            # 登录页
app/src/pages/index/                 # 首页
app/src/pages/player/list/           # 陪玩师列表
app/src/pages/order/create/          # 创建订单
app/src/pages/message/chat/          # 聊天

# 组件
app/src/components/gl/               # 基础组件
app/src/components/PlayerCard.vue    # 陪玩师卡片
app/src/components/OrderCard.vue     # 订单卡片

# 业务逻辑
app/src/composables/usePlayerList.ts # 陪玩师列表
app/src/composables/useOrderCreate.ts # 创建订单
app/src/composables/useChatRoom.ts   # 聊天室

# API
app/src/api/player.ts               # 陪玩师 API
app/src/api/order.ts                # 订单 API
app/src/api/request.ts              # 请求封装
```

---

## 5. 常见代码模式

### 后端：创建新 API

```go
// 1. 定义模型 (api/internal/model/myresource.go)
type MyResource struct {
    Base
    Name string `json:"name" gorm:"size:128;not null"`
    // ...
}

// 2. 定义 Repository 接口 (api/internal/repository/interfaces.go)
type MyResourceRepository interface {
    Create(ctx context.Context, res *MyResource) error
    FindByID(ctx context.Context, id uint64) (*MyResource, error)
}

// 3. 实现 Repository (api/internal/repository/myresource.go)
type myResourceRepository struct {
    db *gorm.DB
}

func NewMyResourceRepository(db *gorm.DB) MyResourceRepository {
    return &myResourceRepository{db: db}
}

func (r *myResourceRepository) Create(ctx context.Context, res *MyResource) error {
    return r.db.WithContext(ctx).Create(res).Error
}

// 4. 定义 Service (api/internal/service/myresource.go)
type MyResourceService struct {
    repo MyResourceRepository
}

func (s *MyResourceService) Create(ctx context.Context, req *CreateRequest) (*MyResource, error) {
    // 业务逻辑
    res := &MyResource{Name: req.Name}
    if err := s.repo.Create(ctx, res); err != nil {
        return nil, err
    }
    return res, nil
}

// 5. 定义 Handler (api/internal/handler/admin/myresource.go)
func CreateMyResource(svc *MyResourceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            resp.Error(c, err)
            return
        }

        res, err := svc.Create(c.Request.Context(), &req)
        if err != nil {
            resp.Error(c, err)
            return
        }

        resp.Success(c, res)
    }
}

// 6. 注册路由 (api/internal/router/admin.go)
func RegisterAdminRoutes(r *gin.RouterGroup, h *Handlers) {
    r.POST("/myresources", h.CreateMyResource)
}
```

### 前端：创建新页面 (React)

```typescript
// 1. 创建页面组件 (admin/src/pages/admin/MyPage/index.tsx)
import { PageContainer, SearchTable } from '@/components';
import { usePermission } from '@/hooks/usePermission';

const MyPage: React.FC = () => {
  const { hasPermission } = usePermission();

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '名称', dataIndex: 'name' },
  ];

  return (
    <PageContainer title="我的页面">
      <SearchTable
        columns={columns}
        dataSource={data}
        loading={loading}
      />
    </PageContainer>
  );
};

export default MyPage;

// 2. 添加路由 (admin/src/router/routes.tsx)
{
  path: 'mypage',
  element: <LazyLoad><MyPage /></LazyLoad>,
  meta: { title: '我的页面', permission: 'admin.mypage.list' }
}

// 3. 添加权限 (admin/src/config/adminRoutes.ts)
{
  name: '我的页面',
  path: '/admin/mypage',
  component: 'MyPage',
  permission: 'admin.mypage.list'
}
```

### 前端：创建新页面 (Vue)

```vue
<!-- 1. 创建页面组件 (app/src/pages/mypage/index.vue) -->
<template>
  <view class="mypage">
    <gl-button @click="handleClick">点击</gl-button>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';

const data = ref([]);

onLoad(() => {
  fetchData();
});

const fetchData = async () => {
  // 获取数据
};

const handleClick = () => {
  // 处理点击
};
</script>

<style lang="scss" scoped>
.mypage {
  padding: 20rpx;
}
</style>

<!-- 2. 配置路由 (app/src/pages.json) -->
{
  "pages": [
    {
      "path": "pages/mypage/index",
      "style": {
        "navigationBarTitleText": "我的页面"
      }
    }
  ]
}
```

### 权限检查

```typescript
// React 管理后台
import { PermissionGuard } from '@/components';

<PermissionGuard permission="admin.users.delete">
  <Button danger>删除</Button>
</PermissionGuard>

// Vue 小程序
import { useAdmin } from '@/composables/useAdmin';

const { hasPermission } = useAdmin();

if (hasPermission('admin.users.delete')) {
  // 有权限
}
```

### WebSocket 使用

```typescript
// 管理后台
import { wsManager } from '@/utils/websocket';

// 连接
wsManager.connect();

// 监听消息
wsManager.on('ORDER_UPDATE', (data) => {
  console.log('订单更新:', data);
});

// 发送消息
wsManager.send('CHAT_MESSAGE', {
  roomId: 123,
  content: 'Hello',
});

// 断开连接
wsManager.disconnect();

// 小程序
import { useWebSocket } from '@/composables/useWebSocket';

const { connected, send, onMessage } = useWebSocket();

onMessage('ORDER_UPDATE', (data) => {
  console.log('订单更新:', data);
});
```

---

## 6. 调试技巧

### 后端调试

```go
// 使用日志
logger.Info("Processing order", "orderID", order.ID)

// 使用 Delve 调试器
dlv debug cmd/main.go

// 设置断点
// 在 VS Code 中点击行号左侧

// 查看 SQL 日志
db.LogMode(true)

// 打印请求
c.Request.URL.Path
c.Request.Method
```

### 前端调试

```typescript
// React DevTools
// 安装 React Developer Tools 浏览器扩展

// console.log
console.log('Data:', data);

// debugger
debugger; // 代码暂停

// Redux DevTools (如果使用)
// 查看状态变化

// Network 面板
// 查看 API 请求

// Vue DevTools
// 安装 Vue DevTools 浏览器扩展
```

### 小程序调试

```bash
# 1. 打开微信开发者工具
# 2. 开启调试模式
# 3. 使用 console.log
# 4. 查看 Network 面板
# 5. 使用 vconsole
```

### 数据库调试

```bash
# 连接数据库
psql -h localhost -U gamelink -d gamelink

# 查看表
\dt

# 查看表结构
\d users

# 执行查询
SELECT * FROM users LIMIT 10;

# 查看索引
\di

# 执行计划
EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = 1;
```

### Redis 调试

```bash
# 连接 Redis
redis-cli

# 查看所有键
KEYS *

# 查看键值
GET key

# 设置键值
SET key value

# 查看键类型
TYPE key

# 查看缓存信息
INFO
```

---

## 7. 故障排查

### 后端问题

**问题**: 端口被占用
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

**问题**: 数据库连接失败
```bash
# 检查 PostgreSQL 是否运行
# Windows
sc query postgresql-x64-16

# Linux
sudo systemctl status postgresql

# 检查连接配置
# .env 文件
DB_HOST=localhost
DB_PORT=5432
DB_USER=gamelink
DB_PASSWORD=your_password
DB_NAME=gamelink
```

**问题**: Redis 连接失败
```bash
# 检查 Redis 是否运行
redis-cli ping

# 应返回: PONG

# 检查连接配置
# .env 文件
REDIS_HOST=localhost
REDIS_PORT=6379
```

**问题**: Go 依赖错误
```bash
go mod download
go mod tidy
```

### 前端问题

**问题**: npm install 失败
```bash
# 清除缓存
npm cache clean --force

# 删除 node_modules
rm -rf node_modules
rm package-lock.json

# 重新安装
npm install
```

**问题**: Vite 开发服务器启动失败
```bash
# 检查端口占用
netstat -ano | findstr :5173

# 更改端口
# vite.config.ts
server: {
  port: 5174,
}
```

**问题**: API 请求失败
```bash
# 检查后端是否启动
curl http://localhost:8080/health

# 检查代理配置
# vite.config.ts
server: {
  proxy: {
    '/api/v1': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

**问题**: TypeScript 类型错误
```bash
# 重新生成类型
npm run type-check

# 或重启 TypeScript Server
# VS Code: Cmd+Shift+P -> "TypeScript: Restart TS Server"
```

### 小程序问题

**问题**: 微信开发者工具报错
```bash
# 1. 检查 AppID
# manifest.json

# 2. 检查服务器域名配置
# 微信公众平台 -> 开发 -> 开发设置 -> 服务器域名

# 3. 清除缓存
# 微信开发者工具 -> 清除缓存
```

**问题**: H5 页面白屏
```bash
# 1. 检查 console 错误
# 2. 检查 API 请求
# 3. 检查是否缺少 polyfill
```

### Docker 问题

**问题**: 容器无法启动
```bash
# 查看日志
docker-compose logs -f

# 重建容器
docker-compose down
docker-compose up -d --build
```

**问题**: 数据库数据丢失
```bash
# 检查数据卷
docker volume ls

# 恢复数据
docker volume inspect gamelink_postgres_data
```

---

## 8. 常用工具

### API 测试

```bash
# Swagger
http://localhost:8080/swagger/index.html

# Postman
# 导入 API collection

# curl
curl -X GET http://localhost:8080/api/v1/health
```

### 数据库工具

```bash
# psql 命令行
psql -h localhost -U gamelink -d gamelink

# DBeaver (GUI)
# 下载: https://dbeaver.io/

# DataGrip (JetBrains)
# 付费工具
```

### Redis 工具

```bash
# redis-cli 命令行
redis-cli

# Redis Insight (GUI)
# 下载: https://redis.com/redis-enterprise/redis-insight/
```

### 日志查看

```bash
# 实时查看日志
tail -f api/logs/app.log

# 查看错误日志
tail -f api/logs/error.log

# 搜索日志
grep "ERROR" api/logs/app.log
```

---

## 9. 快捷键

### VS Code

```bash
# 命令面板
Ctrl+Shift+P (Windows/Linux)
Cmd+Shift+P (Mac)

# 快速打开文件
Ctrl+P (Windows/Linux)
Cmd+P (Mac)

# 切换终端
Ctrl+` (Windows/Linux)
Cmd+` (Mac)

# 格式化代码
Shift+Alt+F (Windows/Linux)
Shift+Option+F (Mac)

# 多光标选择
Alt+Click (Windows/Linux)
Option+Click (Mac)
```

### GoLand

```bash
# 查找文件
Ctrl+Shift+N (Windows/Linux)
Cmd+Shift+O (Mac)

# 查找操作
Ctrl+Shift+A (Windows/Linux)
Cmd+Shift+A (Mac)

# 运行
Shift+F10 (Windows/Linux)
Ctrl+R (Mac)

# 调试
Shift+F9 (Windows/Linux)
Ctrl+D (Mac)
```

### Chrome DevTools

```bash
# 打开开发者工具
F12 (Windows/Linux)
Cmd+Option+I (Mac)

# 元素检查
Ctrl+Shift+C (Windows/Linux)
Cmd+Shift+C (Mac)

# 控制台
Ctrl+Shift+J (Windows/Linux)
Cmd+Option+J (Mac)

# 网络
Ctrl+Shift+I -> Network (Windows/Linux)
Cmd+Option+I -> Network (Mac)
```

---

## 10. 常见错误

### 后端错误

```
Error: dial tcp 127.0.0.1:5432: connect: connection refused
解决: 检查 PostgreSQL 是否启动

Error: NOCONN: connection refused
解决: 检查 Redis 是否启动

Error: unauthorized
解决: 检查 JWT Token 是否有效

Error: record not found
解决: 检查数据库中是否存在该记录
```

### 前端错误

```
Error: Module not found: Can't resolve '@/xxx'
解决: 检查路径别名配置

Error: Network Error
解决: 检查后端服务是否启动

Error: 403 Forbidden
解决: 检查权限配置

Error: 401 Unauthorized
解决: 检查 Token 是否过期
```

### 小程序错误

```
Error: request:fail url not in domain list
解决: 在微信公众平台配置服务器域名

Error: webView is not defined
解决: 检查是否在 H5 环境中使用小程序 API

Error: module 'miniprogram_npm/xxx' is not defined
解决: 重新构建 npm
```

---

**提示**: 将此文档加入浏览器书签，方便随时查阅！

**更新**: 当发现新的常用模式或解决方案时，请及时更新此文档。
