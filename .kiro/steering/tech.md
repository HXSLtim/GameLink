# Technology Stack

## Backend

- **Language**: Go 1.25.3+
- **Web Framework**: Gin (HTTP routing and middleware)
- **ORM**: GORM (database operations with auto-migration)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Database**: PostgreSQL 16+ (全流程使用)
- **Cache**: Redis 7+ (go-redis/v9，全流程使用)
- **WebSocket**: gorilla/websocket
- **API Documentation**: Swagger/OpenAPI (swaggo/swag)
- **Dependency Injection**: Google Wire
- **Testing**: Go testing + testify, mockery for mocks
- **Monitoring**: Prometheus client
- **Encryption**: AES-256-CBC + SHA-256 签名

## Frontend

- **Framework**: React 18.2+ with TypeScript 5.2+
- **Build Tool**: Vite 7.2+
- **UI Library**: Ant Design 6.0
- **Routing**: React Router 7.9+
- **HTTP Client**: Axios 1.13+ (带加密拦截器)
- **State Management**: React Context API
- **Encryption**: crypto-js (AES-256-CBC)
- **Styling**: Less 4.2
- **WebSocket**: socket.io-client 4.8+
- **Testing**: Vitest 4.0+, Testing Library, Playwright
- **Code Quality**: ESLint, Prettier

## Infrastructure

- **Containerization**: Docker + Docker Compose
- **Database**: PostgreSQL 16 (Alpine)
- **Cache**: Redis 7 (Alpine)
- **Reverse Proxy**: Nginx
- **CI/CD**: GitHub Actions (已配置)

## Deployment

### 部署脚本（推荐方式）

```powershell
# 生产环境部署（加密版，推荐）
.\scripts\deploy-production-encrypted.ps1

# 生产环境部署（标准版）
.\scripts\deploy-production.ps1

# 部署脚本参数
-SkipBuild        # 跳过 Docker 镜像构建
-SkipFrontend     # 跳过前端构建
-NoPull           # 不拉取基础镜像
-RegenerateKeys   # 重新生成加密密钥（仅加密版）
```

**部署脚本自动执行：**
1. 检查环境变量（.env 文件）
2. 生成/验证加密密钥（加密版）
3. 同步密钥到前端（加密版）
4. 安装前端依赖（包括 crypto-js）
5. 构建前端（npm run build）
6. 构建 Docker 镜像
7. 停止旧服务
8. 启动新服务 + 健康检查

**重要：修改代码后直接运行部署脚本即可，不需要手动构建！**

### 手动 Docker 操作（仅调试时使用）

```bash
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml down
docker logs gamelink-backend --tail=50
```

## Key Configuration

### 环境变量（.env）

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

### 部署注意事项

- 生产环境必须启用加密中间件（CRYPTO_ENABLED=true）
- 缓存类型必须设置为 redis（CACHE_TYPE=redis）
- 超级管理员角色 slug 是 `superAdmin`（驼峰式）
- 后端健康检查路径是 `/api/v1/healthz`
- 数据库密码不能包含 URL 特殊字符（如 `%`）
- 管理员密码必须包含大写、小写、数字和特殊字符，至少 8 个字符
- AES-256 加密需要 32 字节密钥，IV 需要 16 字节

## Common Commands

### Backend

```bash
cd backend
go mod tidy                      # 整理依赖
go run cmd/main.go               # 运行应用
make test                        # 运行测试
make lint                        # 代码检查
make swagger                     # 生成 API 文档
```

### Frontend

```bash
cd frontend
npm install                      # 安装依赖
npm run dev                      # 开发服务器 (localhost:5173)
npm run build                    # 生产构建
npm run lint                     # 代码检查
```

## Architecture

### Backend Layered Architecture

```
Handler → Service → Repository → Model
```

- **Handler**: HTTP 请求处理、参数验证、响应格式化
- **Service**: 业务逻辑、事务管理、跨模块协调
- **Repository**: 数据库操作、缓存、查询封装
- **Model**: 数据结构、数据库映射、验证规则

### Error Handling

三层错误机制：
- Repository errors: 数据库级错误
- Service errors: 业务逻辑错误（带上下文）
- API errors: 标准化 HTTP 响应（带错误码）

使用 `fmt.Errorf("context: %w", err)` 进行错误包装。

## Testing Strategy

- **Unit Tests**: 表驱动测试，mock 依赖
- **Integration Tests**: 真实数据库（PostgreSQL），测试 fixtures
- **Concurrency Tests**: Race detector，压力测试
- **Coverage Target**: 76.4%（当前），目标 80%+
