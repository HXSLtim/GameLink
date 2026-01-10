# 技术栈

## 后端 (Go)

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25.3+ |
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
| 框架 | React + TypeScript | 19.x / 5.9+ |
| 构建工具 | Vite | 7.2+ |
| UI 库 | Ant Design | 6.0 |
| 路由 | React Router | 7.9+ |
| HTTP 客户端 | Axios（带加密拦截器） | 1.13+ |
| 状态管理 | React Context API | - |
| 加密 | crypto-js (AES-256-CBC) | - |
| 样式 | Less | 4.4 |
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
| 变更检测 | 只运行受影响的任务（api/admin） |
| 并发控制 | 自动取消重复运行 |
| 后端测试 | PostgreSQL + Redis 服务，race 检测 |
| 覆盖率检查 | **70% 以下构建失败**（质量门强制执行） |
| 代码检查 | **Linter 失败会阻塞构建**（质量门强制执行） |
| 前端测试 | 类型检查 + Lint + 单元测试（全部强制执行） |
| Docker 构建 | main/dev 分支自动构建镜像 |

**质量门策略**:
- 后端/前端 Linter 失败 → **构建失败**
- 测试失败 → **构建失败**
- 覆盖率 < 70% → **构建失败**
- 覆盖率上传失败 → **不阻塞构建**（非关键步骤）

### Security 流程 (security.yml)

| 功能 | 说明 |
|------|------|
| Go 安全扫描 | Gosec（high 级别以上失败）+ govulncheck（有漏洞即失败） |
| NPM 审计 | **高危/严重漏洞检测并失败构建** |
| Docker 扫描 | Trivy 镜像漏洞扫描（CRITICAL/HIGH 级别失败） |
| 密钥检测 | Gitleaks 密钥泄露检测（发现即失败） |
| 定时运行 | 每周一凌晨自动运行 |

**安全门策略**:
- Go 已知漏洞 → **构建失败**
- NPM高危/严重漏洞 → **构建失败**
- Docker 镜像高危漏洞 → **构建失败**
- 密钥泄露 → **构建失败**
- SARIF 上传失败 → **不阻塞构建**（非关键步骤）

### Deploy 流程 (deploy.yml)

| 功能 | 说明 |
|------|------|
| 触发方式 | Tag 推送 / 手动触发 |
| 环境选择 | staging / production |
| 预检查 | 版本和环境自动判断 |
| 健康检查 | HTTP 请求验证（后端 `/api/v1/healthz` + 前端 `/`） |
| 重试机制 | 5 次尝试，每次间隔 10 秒 |
| 超时时间 | 每次 curl 请求 10 秒超时 |
| 回滚机制 | 失败时自动回滚 |

健康检查实现：
```bash
# 后端健康检查端点
GET /api/v1/healthz  # 返回 200 + {"status": "ok"}

# 健康检查逻辑（5 次重试）
for i in {1..5}; do
  curl -f -s --max-time 10 $BACKEND_URL/api/v1/healthz && exit 0
  sleep 10
done
exit 1
```

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

> **重要**: 生产环境必须配置以下安全密钥，否则应用将拒绝启动。详见 [安全配置指南](../../docs/SECURITY_CONFIG.md)

```bash
# 数据库
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=<安全密码，不含特殊字符>
POSTGRES_DB=gamelink

# Redis
REDIS_PASSWORD=<安全密码>

# JWT（必须，32+字符）
# 生成命令: openssl rand -base64 32
JWT_SECRET_KEY=<生成的32+字符密钥>

# 加密（生产环境必须）
CRYPTO_ENABLED=true
# 密钥长度要求：必须是 16/24/32 字节
# 生成命令: openssl rand -base64 32 (密钥) 或 openssl rand -base64 16 (IV)
CRYPTO_SECRET_KEY=<生成的32字节密钥>
CRYPTO_IV=<生成的16字节IV>

# 超级管理员（必须，8+字符，包含大小写、数字、特殊符号）
# 生成命令: openssl rand -base64 24
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=<生成的强密码>
SUPER_ADMIN_NAME=Super Admin
```

**快速生成所有密钥**:
```bash
# Windows
.\scripts\generate-secrets.ps1

# Linux/Mac
./scripts/generate-secrets.sh
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
cd api
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
// 位置: api/internal/handler/resp/

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

// Admin 包辅助函数（api/internal/handler/admin/helpers.go）
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
// apierr 包常用函数（api/pkg/apierr/errors.go）
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

Axios 拦截器返回完整的 `AxiosResponse`，调用方需要手动解析响应数据。

```typescript
// API 响应结构
interface ApiResponse<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;  // 业务数据
}

// 正确用法 - 需要访问 response.data.data
const response = await api.getUsers();
const users = response.data.data;  // response.data 是 ApiResponse，.data 是业务数据

// 错误处理时可以访问
if (!response.data.success) {
  console.error(response.data.message);
}
```

> **注意**：`useCrud` hook 内部已处理响应解析，使用该 hook 时无需手动解析。
