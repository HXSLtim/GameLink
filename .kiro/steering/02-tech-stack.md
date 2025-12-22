# 技术栈

## 后端 (Go)

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25.5+ |
| Web 框架 | Gin | - |
| ORM | GORM | - |
| 认证 | JWT (golang-jwt/jwt/v5) | - |
| 数据库 | PostgreSQL | 16+ |
| 缓存 | Redis (go-redis/v9) | 7+ |
| WebSocket | gorilla/websocket | - |
| API 文档 | Swagger/OpenAPI (swaggo/swag) | - |
| 依赖注入 | Google Wire | - |
| 测试 | testify + mockery | - |
| 加密 | AES-256-CBC + SHA-256 | - |

## 前端 - 管理后台 (React)

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | React + TypeScript | 18.2+ / 5.2+ |
| 构建工具 | Vite | 7.2+ |
| UI 库 | Ant Design | 6.0 |
| 路由 | React Router | 7.9+ |
| HTTP 客户端 | Axios（带加密拦截器） | 1.13+ |
| 状态管理 | React Context API | - |
| 加密 | crypto-js (AES-256-CBC) | - |
| 样式 | Less | 4.2 |
| WebSocket | socket.io-client | 4.8+ |
| 测试 | Vitest + Testing Library + Playwright | 4.0+ |

## 基础设施

- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **CI/CD**: GitHub Actions

## CI/CD 流程

### CI 流程 (ci.yml)

| 功能 | 说明 |
|------|------|
| 变更检测 | 只运行受影响的任务（backend/admin） |
| 并发控制 | 自动取消重复运行 |
| 后端测试 | PostgreSQL + Redis 服务，race 检测 |
| 覆盖率检查 | 70% 以下警告 |
| 前端测试 | 类型检查 + Lint + 单元测试 |
| Docker 构建 | main/dev 分支自动构建镜像 |

### Security 流程 (security.yml)

| 功能 | 说明 |
|------|------|
| Go 安全扫描 | Gosec + govulncheck |
| NPM 审计 | 高危/严重漏洞检测 |
| Docker 扫描 | Trivy 镜像漏洞扫描 |
| 密钥检测 | Gitleaks 密钥泄露检测 |
| 定时运行 | 每周一凌晨自动运行 |

### Deploy 流程 (deploy.yml)

| 功能 | 说明 |
|------|------|
| 触发方式 | Tag 推送 / 手动触发 |
| 环境选择 | staging / production |
| 预检查 | 版本和环境自动判断 |
| 健康检查 | 部署后自动验证 |
| 回滚机制 | 失败时自动回滚 |

## 部署

### 推荐方式（部署脚本）

```powershell
# 生产环境部署（加密版，推荐）
.\scripts\deploy-production-encrypted.ps1

# 部署脚本参数
-SkipBuild        # 跳过 Docker 镜像构建
-SkipAdmin        # 跳过管理后台构建
-NoPull           # 不拉取基础镜像
-RegenerateKeys   # 重新生成加密密钥
```

### 环境变量 (.env)

```bash
# 数据库
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=<安全密码，不含特殊字符>
POSTGRES_DB=gamelink

# Redis
REDIS_PASSWORD=<安全密码>

# JWT
JWT_SECRET_KEY=<32字符以上>

# 加密（生产环境必须）
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=<32字符>
CRYPTO_IV=<16字符>

# 超级管理员
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=<包含大小写数字特殊字符，8位以上>
```

### 注意事项

- 生产环境必须启用加密中间件（CRYPTO_ENABLED=true）
- 缓存类型必须设置为 redis（CACHE_TYPE=redis）
- 超级管理员角色 slug 是 `superAdmin`（驼峰式）
- 后端健康检查路径是 `/api/v1/healthz`
- 数据库密码不能包含 URL 特殊字符（如 `%`）

## 常用命令

### 后端

```bash
cd backend
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

## 架构模式

### 后端分层

```
Handler → Service → Repository → Model
```

- **Handler**: HTTP 请求处理、参数验证、响应格式化
- **Service**: 业务逻辑、事务管理、跨模块协调
- **Repository**: 数据库操作、缓存、查询封装
- **Model**: 数据结构、数据库映射、验证规则

### 统一响应规范

Handler 层使用 `resp` 包和 `apierr` 包统一响应格式：

```go
// 位置: backend/internal/handler/resp/

// 成功响应
resp.OK(c, data)           // 200 + data
resp.Created(c, data)      // 201 + data
resp.Updated(c, data)      // 200 + "updated" + data
resp.Deleted(c)            // 200 + "deleted"
resp.List(c, items, pagination)  // 200 + 分页列表

// 错误响应（使用 apierr 包）
resp.Error(c, apierr.BadRequest("invalid input"))
resp.Error(c, apierr.NotFound("resource not found"))
resp.Error(c, apierr.InternalError("operation failed").WithDetails(err.Error()))

// Admin 包辅助函数（backend/internal/handler/admin/helpers.go）
respondSuccess(c, data)    // 通用成功
respondCreated(c, data)    // 创建成功
respondUpdated(c, data)    // 更新成功
respondDeleted(c)          // 删除成功
respondList(c, items, pagination)  // 分页列表
respondMsg(c, "message")   // 仅消息无数据
respondError(c, err)       // 统一错误处理
respondBadRequest(c, msg)  // 400 错误
ParseIDAndRespond(c, "id") // 解析ID并自动响应错误
ValidateAndRespond(c, &req) // 绑定JSON并自动响应错误
```

### 错误处理

使用 `fmt.Errorf("context: %w", err)` 进行错误包装，三层错误机制：
- Repository errors: 数据库级错误
- Service errors: 业务逻辑错误（带上下文）
- API errors: 标准化 HTTP 响应（使用 `apierr` 包）

```go
// apierr 包常用函数（backend/pkg/apierr/errors.go）
apierr.BadRequest(msg)      // 400
apierr.Unauthorized(msg)    // 401
apierr.Forbidden(msg)       // 403
apierr.NotFound(msg)        // 404
apierr.Conflict(msg)        // 409
apierr.InternalError(msg)   // 500

// 链式调用添加详情
apierr.BadRequest("validation failed").WithDetails(err.Error())

// 错误类型判断
apierr.IsNotFound(err)
apierr.IsValidationError(err)
```

### 前端请求层规范

Axios 拦截器已在响应层自动解析 `response.data.data`，页面组件直接使用返回值即可。

```typescript
// ✅ 正确用法 - 拦截器已解析
const resp = await api.getUsers();
const users = resp.data;  // 直接是业务数据

// ❌ 错误用法 - 多余的解析
const users = resp.data.data;  // 不需要！
```

API 响应结构：
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }  // ← 拦截器返回这一层
}
```
