# GameLink 依赖和配置文档

本文档详细记录 GameLink 项目的所有关键依赖、配置和环境变量，确保团队新成员快速上手，避免配置问题。

## 目录

- [环境变量配置](#环境变量配置)
- [外部服务依赖](#外部服务依赖)
- [构建和部署配置](#构建和部署配置)
- [测试设置和覆盖率](#测试设置和覆盖率)
- [CI/CD 管道配置](#cicd-管道配置)
- [安全考虑和敏感配置](#安全考虑和敏感配置)

---

## 环境变量配置

### 1. 后端环境变量 (`.env.example`)

#### 应用环境
```bash
# 应用运行环境
APP_ENV=development  # 选项: development, staging, production

# Gin 框架模式
GIN_MODE=debug       # 选项: debug, release, test
                     # 生产环境必须使用 release

# 服务端口
BACKEND_PORT=8080
```

#### 数据库配置 (PostgreSQL 16)
```bash
# PostgreSQL 连接
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=gamelink123  # 生产环境使用强密码
POSTGRES_DB=gamelink
POSTGRES_PORT=5433             # 容器内使用 5432

# 连接池配置
DB_MAX_CONNS=25                # 最大连接数
DB_MAX_IDLE=5                  # 最大空闲连接数

# 读写分离（可选，高可用配置）
# DB_READER_DSNS=host=postgres-replica port=5432 ...
```

#### 缓存配置 (Redis 7)
```bash
# 缓存类型
CACHE_TYPE=redis               # 选项: redis, memory

# Redis 连接
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=                # 生产环境必须设置密码
REDIS_DB=0
REDIS_PORT=6380                # 容器内使用 6379
```

#### 安全配置
```bash
# 加密配置（重要！）
CRYPTO_ENABLED=false           # 开发环境: false, 生产环境: true

# AES-256-CBC 加密密钥（CRYPTO_ENABLED=true 时必填）
# 生成命令: openssl rand -base64 32
CRYPTO_SECRET_KEY=             # 32字节 (base64 编码后约 44 字符)

# 初始化向量（CRYPTO_ENABLED=true 时必填）
# 生成命令: openssl rand -base64 16
CRYPTO_IV=                     # 16字节 (base64 编码后约 24 字符)

# SHA-256 签名
CRYPTO_USE_SIGNATURE=false
```

**加密密钥生成示例：**
```bash
# 生成 SECRET_KEY (32 字节)
openssl rand -base64 32 | tr -d '\n'

# 生成 IV (16 字节)
openssl rand -base64 16 | tr -d '\n'
```

#### JWT 认证
```bash
# JWT 密钥（生产环境必填，建议 32+ 字符）
# 生成命令: openssl rand -base64 32
JWT_SECRET_KEY=dev-jwt-secret-key-minimum-32-chars

# Token 有效期
JWT_TOKEN_TTL_HOURS=24
```

#### 超级管理员配置
```bash
# 超级管理员账号（系统初始化时创建）
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=Admin123456  # 至少 8 字符，包含大小写字母、数字
SUPER_ADMIN_NAME=Super Admin
```

#### 种子数据配置
```bash
# 是否初始化演示数据
SEED_ENABLED=true               # 生产环境必须为 false
```

#### 外部 API 配置

**微信小程序（可选）**
```bash
WECHAT_APPID=                   # 微信小程序 AppID
WECHAT_SECRET=                  # 微信小程序 Secret
# 获取方式: 微信公众平台 -> 开发 -> 开发管理 -> 开发设置
```

**微信支付（可选）**
```bash
WECHAT_PAY_APP_ID=
WECHAT_PAY_MCH_ID=              # 商户号
WECHAT_PAY_API_KEY=             # API 密钥
WECHAT_PAY_API_CERT_PATH=       # 证书路径
WECHAT_PAY_NOTIFY_URL=          # 回调地址
WECHAT_PAY_ENABLED=false
```

**支付宝（可选）**
```bash
ALIPAY_APP_ID=
ALIPAY_PRIVATE_KEY_PATH=        # 商户私钥路径
ALIPAY_PUBLIC_KEY_PATH=         # 支付宝公钥路径
ALIPAY_NOTIFY_URL=
ALIPAY_ENABLED=false
```

**短信服务（可选）**
```bash
SMS_PROVIDER=                   # 选项: aliyun, tencent
SMS_ACCESS_KEY=
SMS_SECRET_KEY=
SMS_SIGN_NAME=                  # 短信签名
SMS_ENABLED=false
```

**对象存储（可选）**
```bash
OSS_PROVIDER=                   # 选项: aliyun, qcloud, minio
OSS_ENDPOINT=                   # 服务端点
OSS_ACCESS_KEY=
OSS_SECRET_KEY=
OSS_BUCKET=                     # 存储桶名称
OSS_REGION=                     # 区域
OSS_ENABLED=false
OSS_SIGNED_URL_TTL_SECONDS=3600 # 签名 URL 有效期（秒）
```

---

### 2. 管理后台环境变量 (`admin/.env.example`)

#### API 配置
```bash
# 后端 API 基础 URL（不带尾部斜杠）
VITE_API_BASE_URL=http://localhost:8080
```

#### 安全配置
```bash
# 客户端加密（必须与后端 CRYPTO_ENABLED 匹配）
VITE_CRYPTO_ENABLED=false

# AES-256-CBC 加密密钥（VITE_CRYPTO_ENABLED=true 时必填）
# 注意：前端使用原始字节，不是 base64 编码
# 生成命令:
#   openssl rand -base64 32 | base64 -d | xxd -p -c 32
VITE_CRYPTO_SECRET_KEY=

# 初始化向量（16 字节原始字节）
# 生成命令:
#   openssl rand -base64 16 | base64 -d | xxd -p -c 16
VITE_CRYPTO_IV=

# SHA-256 签名（必须与后端匹配）
VITE_CRYPTO_USE_SIGNATURE=true
```

**重要提示：**
- 前端加密密钥会暴露在浏览器中，应与后端密钥保持一致
- 开发环境可以禁用加密 (`VITE_CRYPTO_ENABLED=false`)
- 生产环境如果后端启用加密，前端也必须启用

#### 应用配置
```bash
# 应用信息
VITE_APP_TITLE=GameLink Admin Panel
VITE_APP_VERSION=1.0.0

# 调试模式
VITE_DEBUG=false
```

#### 功能开关
```bash
# WebSocket 实时更新
VITE_ENABLE_WEBSOCKET=true
VITE_WEBSOCKET_RECONNECT_ATTEMPTS=5
VITE_WEBSOCKET_RECONNECT_INTERVAL=3000
```

#### UI 配置
```bash
# 分页配置
VITE_DEFAULT_PAGE_SIZE=20
VITE_MAX_PAGE_SIZE=100

# 日期格式
VITE_DATE_FORMAT=YYYY-MM-DD HH:mm:ss

# 时区
VITE_TIMEZONE=Asia/Shanghai
```

---

### 3. 移动端环境变量 (`app/.env.example`)

```bash
# API 基础路径
VITE_API_BASE_URL=http://localhost:8081/api/v1

# WebSocket 地址
VITE_WS_URL=ws://localhost:8081/api/v1/ws

# 环境标识
VITE_ENV=development
```

---

## 外部服务依赖

### 核心依赖（必需）

#### 1. PostgreSQL 16
**用途：** 主数据库

**连接配置：**
```bash
host=localhost
port=5432
user=gamelink
password=***
database=gamelink
sslmode=disable
```

**性能调优（生产环境）：**
```bash
max_connections=200
shared_buffers=256MB
effective_cache_size=768MB
maintenance_work_mem=64MB
```

**健康检查：**
```bash
pg_isready -U gamelink -d gamelink
```

#### 2. Redis 7
**用途：** 缓存、会话存储、WebSocket Pub/Sub

**连接配置：**
```bash
host=localhost
port=6379
password=***       # 生产环境必须设置
db=0
```

**性能调优（生产环境）：**
```bash
maxmemory=512MB
maxmemory-policy=allkeys-lru
appendonly=yes
appendfsync=everysec
```

**健康检查：**
```bash
redis-cli -a *** ping
```

### 可选依赖

#### 3. 微信小程序
**用途：** 小程序登录、支付

**必需配置：**
- `WECHAT_APPID`: 小程序 AppID
- `WECHAT_SECRET`: 小程序 Secret

**获取方式：**
1. 登录微信公众平台
2. 开发 → 开发管理 → 开发设置
3. 复制 AppID 和生成 Secret

#### 4. 微信支付
**用途：** 在线支付

**必需配置：**
- `WECHAT_PAY_APP_ID`: 应用 ID
- `WECHAT_PAY_MCH_ID`: 商户号
- `WECHAT_PAY_API_KEY`: API 密钥
- `WECHAT_PAY_API_CERT_PATH`: 证书路径

**安全注意事项：**
- 证书文件权限应为 `600`
- API 密钥应定期轮换
- 回调 URL 必须是 HTTPS

#### 5. 支付宝
**用途：** 在线支付

**必需配置：**
- `ALIPAY_APP_ID`: 应用 ID
- `ALIPAY_PRIVATE_KEY_PATH`: 商户私钥文件路径
- `ALIPAY_PUBLIC_KEY_PATH`: 支付宝公钥文件路径

**密钥生成：**
```bash
# 生成商户私钥
openssl genrsa -out app_private_key.pem 2048

# 提取公钥
openssl rsa -in app_private_key.pem -pubout -out app_public_key.pem
```

#### 6. 短信服务
**用途：** 发送验证码、通知短信

**支持提供商：**
- 阿里云 (`aliyun`)
- 腾讯云 (`tencent`)

**必需配置：**
- `SMS_PROVIDER`: 提供商
- `SMS_ACCESS_KEY`: 访问密钥
- `SMS_SECRET_KEY`: 秘密密钥
- `SMS_SIGN_NAME`: 短信签名

#### 7. 对象存储 (OSS)
**用途：** 文件上传、图片存储

**支持提供商：**
- 阿里云 OSS (`aliyun`)
- 腾讯云 COS (`qcloud`)
- MinIO (`minio`)

**必需配置：**
- `OSS_PROVIDER`: 提供商
- `OSS_ENDPOINT`: 服务端点
- `OSS_ACCESS_KEY`: 访问密钥
- `OSS_SECRET_KEY`: 秘密密钥
- `OSS_BUCKET`: 存储桶名称
- `OSS_REGION`: 区域

---

## 构建和部署配置

### 1. 后端构建配置

#### Go 模块 (`api/go.mod`)
```go
module gamelink

go 1.24.5
```

**主要依赖：**
- `gin-gonic/gin v1.10.0` - Web 框架
- `gorm.io/gorm v1.30.0` - ORM
- `gorm.io/driver/postgres v1.5.9` - PostgreSQL 驱动
- `redis/go-redis/v9 v9.6.3` - Redis 客户端
- `golang-jwt/jwt/v5 v5.2.3` - JWT 认证
- `gorilla/websocket v1.5.3` - WebSocket
- `swaggo/gin-swagger v1.6.1` - API 文档
- `prometheus/client_golang v1.18.0` - 监控指标

#### Dockerfile 特性
- **多阶段构建：** Builder + Runtime
- **基础镜像：** `golang:1.25.5-alpine` → `alpine:3.20`
- **安全特性：**
  - 非 root 用户运行（gamelink:1000）
  - 静态编译（CGO_ENABLED=0）
  - 最小化镜像
- **优化：**
  - BuildKit 缓存加速
  - 版本信息注入
  - 健康检查

#### 环境差异

| 配置项 | 开发环境 | 生产环境 | 高可用 |
|--------|---------|---------|--------|
| CRYPTO_ENABLED | false | true | true |
| SEED_ENABLED | true | false | false |
| GIN_MODE | debug | release | release |
| 数据库实例 | 单实例 | 单实例 | 主从复制 |
| 后端实例 | 1 | 1 | 2+ |
| 资源限制 | 无 | 有 | 有 |

### 2. 管理后台构建配置

#### 包管理 (`admin/package.json`)
```json
{
  "name": "frontend",
  "type": "module",
  "dependencies": {
    "react": "^19.2.0",
    "antd": "^6.0.0",
    "react-router-dom": "^7.9.6",
    "axios": "^1.13.2",
    "socket.io-client": "^4.8.1",
    "zustand": "^5.0.9"
  }
}
```

**构建脚本：**
```bash
npm run dev          # 开发服务器
npm run build        # 生产构建
npm run build:analyze # Bundle 分析
npm run lint         # 代码检查
npm run test         # 单元测试
npm run test:e2e     # E2E 测试
```

#### Vite 配置特性

**性能优化：**
- 代码分割策略（React、AntD、图表、工具库分离）
- Gzip + Brotli 压缩
- PWA 支持（Service Worker 缓存）
- 源码映射（开发环境启用）

**安全头：**
```typescript
'Content-Security-Policy': "default-src 'self'..."
'X-Frame-Options': 'DENY'
'X-Content-Type-Options': 'nosniff'
'Strict-Transport-Security': 'max-age=31536000'
```

**PWA 缓存策略：**
- API 请求：NetworkFirst（24小时）
- 静态资源：CacheFirst（30天）
- CSS/JS：StaleWhileRevalidate（7天）

### 3. 移动端构建配置

#### UniApp 支持
- 微信小程序
- 支付宝小程序
- H5
- App (iOS/Android)

**构建脚本：**
```bash
npm run dev:mp-weixin    # 微信小程序开发
npm run build:mp-weixin  # 微信小程序构建
npm run dev:h5           # H5 开发
npm run build:h5         # H5 构建
```

---

## 测试设置和覆盖率

### 1. 后端测试

#### 测试框架
- **测试运行器：** Go 内置 testing
- **断言库：** `testify v1.11.1`
- **Mock 工具：** `golang/mock v1.6.0`
- **Redis Mock：** `miniredis v2.35.0`

#### 测试配置
```bash
# 运行所有测试
go test -v ./...

# 运行测试并生成覆盖率
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 查看覆盖率
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

#### 覆盖率要求
- **最低覆盖率：** 70%
- **CI 门禁：** 低于 70% 将失败
- **关键模块：** 80%+

#### 集成测试
```bash
# 设置测试数据库
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=gamelink
export TEST_DB_PASSWORD=testpassword
export TEST_DB_NAME=gamelink_test

# 运行集成测试
go test -v -tags=integration ./...
```

### 2. 管理后台测试

#### 测试框架
- **单元测试：** Vitest 4.0
- **E2E 测试：** Playwright 1.57
- **测试工具库：** Testing Library

#### 测试配置 (`vitest`)

```typescript
test: {
  globals: true,
  environment: 'jsdom',
  setupFiles: ['./src/test/setup.ts'],
  coverage: {
    reporter: ['text', 'json', 'html'],
    exclude: ['node_modules/', 'src/test/']
  }
}
```

#### 测试命令
```bash
npm run test           # 单元测试（交互模式）
npm run test:ui        # 测试 UI 界面
npm run test:run       # 单元测试（CI 模式）
npm run test:e2e       # E2E 测试
npm run test:e2e:ui    # E2E 测试 UI
```

#### Playwright 配置
- **浏览器：** Chromium
- **超时：** 30秒
- **重试：** 2次
- **并行执行：** 是

### 3. 代码质量检查

#### Go Linters (`.golangci.yml`)

**启用的检查器：**
- `govet` - Go 静态分析
- `gofmt` - 代码格式化
- `goimports` - import 排序
- `gocyclo` - 圈复杂度检查（阈值 15）
- `errcheck` - 错误检查
- `dupl` - 重复代码检测（阈值 150）

**运行：**
```bash
golangci-lint run --timeout=5m
```

#### Node.js Linters

**ESLint 配置：**
- `eslint` 9.39
- `typescript-eslint` 8.46
- `eslint-plugin-react` 7.0
- `prettier` 3.6

**运行：**
```bash
npm run lint
```

---

## CI/CD 管道配置

### 1. CI Pipeline (`.github/workflows/ci.yml`)

#### 触发条件
- Push 到 `main`, `dev`, `develop` 分支
- Pull Request 到上述分支

#### 工作流程

**阶段 1：变更检测**
- 检测 `api/`, `admin/`, `client/` 变更
- 跳过未变更模块的构建

**阶段 2：测试和构建**
```
Backend (Go 1.25.5)
  ├── golangci-lint
  ├── go test (race + coverage)
  ├── 检查覆盖率 ≥ 70%
  └── 构建

Admin (Node.js 20)
  ├── npm ci
  ├── tsc --noEmit (类型检查)
  ├── eslint
  ├── npm test
  └── npm build

Client (Node.js 20)
  ├── npm ci
  ├── tsc --noEmit
  ├── eslint
  ├── npm test
  └── npm build
```

**阶段 3：Docker 构建**
- 仅在 `main`/`dev` 分支
- 所有测试通过后
- 使用 GitHub Actions 缓存

#### 质量门禁

| 检查项 | 要求 | 失败处理 |
|--------|------|---------|
| golangci-lint | 0 错误 | 阻断构建 |
| 测试覆盖率 | ≥ 70% | 阻断构建 |
| TypeScript | 0 错误 | 阻断构建 |
| ESLint | 0 错误 | 阻断构建 |
| 单元测试 | 全部通过 | 阻断构建 |

### 2. Deploy Pipeline (`.github/workflows/deploy.yml`)

#### 触发条件
- Tag: `v*`
- 手动触发

#### 部署流程

**预检查：**
- 确定版本号
- 确定目标环境（staging/production）

**测试（可选跳过）：**
- 后端测试（short 模式）
- 前端测试

**构建和推送：**
- 构建 Docker 镜像
- 推送到 GHCR
- 标签策略：`v1.0.0`, `v1.0`, `latest`

**部署：**
- Staging: `https://staging.gamelink.example.com`
- Production: `https://gamelink.example.com`

**健康检查：**
- 后端健康检查：`/api/v1/healthz`
- 前端健康检查：`/`
- 重试 5 次，间隔 10 秒

**失败回滚：**
- 自动回滚到上一版本

### 3. Security Pipeline (`.github/workflows/security.yml`)

#### 触发条件
- Push 到 `main`, `dev`
- Pull Request 到 `main`
- 定时：每周一凌晨
- 手动触发

#### 安全检查

**Go 安全扫描：**
- `govulncheck` - 已知漏洞扫描
- 检测标准：任何漏洞 → 失败

**Node.js 审计：**
- `npm audit --audit-level=high`
- Critical 级别 → 失败
- High 级别 → 警告

**Docker 镜像扫描：**
- `Trivy` 漏洞扫描
- 仅检查 CRITICAL 级别
- CRITICAL → 失败，HIGH → 警告

**密钥泄露检测：**
- `Gitleaks` 扫描
- 检测标准：任何密钥 → 失败

---

## 安全考虑和敏感配置

### 1. 敏感信息管理

#### 永远不要提交的内容
```
.env                    # 环境变量
*.key                   # 私钥文件
*.pem                   # 证书文件
*.cert                  # 证书文件
passwords.txt           # 密码文件
credentials.json        # 凭证文件
```

#### .gitignore 配置
```gitignore
# 环境变量
.env
.env.local
.env.*.local

# 密钥和证书
*.key
*.pem
*.cert
*.crt
secrets/

# 日志文件
*.log
logs/
```

### 2. 密钥生成和存储

#### 生成密钥

**JWT 密钥：**
```bash
openssl rand -base64 32
```

**加密密钥：**
```bash
# SECRET_KEY (32 字节)
openssl rand -base64 32 | tr -d '\n'

# IV (16 字节)
openssl rand -base64 16 | tr -d '\n'
```

**生产环境密钥存储：**
- Kubernetes Secrets
- AWS Secrets Manager
- Azure Key Vault
- HashiCorp Vault

### 3. 加密配置

#### 前后端加密一致性

**配置检查清单：**
- [ ] 后端 `CRYPTO_ENABLED` = 前端 `VITE_CRYPTO_ENABLED`
- [ ] 后端 `CRYPTO_SECRET_KEY` = 前端 `VITE_CRYPTO_SECRET_KEY`
- [ ] 后端 `CRYPTO_IV` = 前端 `VITE_CRYPTO_IV`
- [ ] 后端 `CRYPTO_USE_SIGNATURE` = 前端 `VITE_CRYPTO_USE_SIGNATURE`

#### 加密流程
1. 前端使用 AES-256-CBC 加密请求数据
2. 添加 HMAC-SHA256 签名
3. 后端验证签名
4. 后端解密数据
5. 处理业务逻辑
6. 后端加密响应数据
7. 前端解密响应

### 4. CORS 和安全头

#### CORS 配置
```go
// 允许的源
AllowOrigins: []string{
    "https://gamelink.example.com",
    "https://admin.gamelink.example.com",
}

// 允许的方法
AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

// 允许的头
AllowHeaders: []string{
    "Origin",
    "Content-Type",
    "Authorization",
    "X-Signature",
}
```

#### 安全头
```http
Content-Security-Policy: default-src 'self'; ...
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Strict-Transport-Security: max-age=31536000
X-XSS-Protection: 1; mode=block
```

### 5. 数据库安全

#### 连接安全
- 生产环境使用 SSL/TLS
- 强密码策略
- 限制数据库访问 IP
- 定期轮换密码

#### 备份策略
```bash
# 每日备份
0 2 * * * pg_dump -U gamelink gamelink | gzip > backup_$(date +\%Y\%m\%d).sql.gz

# 保留 30 天
find /backup -name "backup_*.sql.gz" -mtime +30 -delete
```

### 6. Redis 安全

#### 访问控制
```bash
# 设置密码
requirepass your_strong_password

# 禁用危险命令
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG ""
```

#### 网络安全
- 绑定到本地回环地址
- 使用防火墙限制访问
- 定期轮换密码

### 7. 容器安全

#### Docker 最佳实践
- 使用非 root 用户运行
- 最小化镜像（Alpine）
- 定期更新基础镜像
- 扫描镜像漏洞
- 限制容器资源

#### 安全扫描
```bash
# Trivy 扫描
trivy image gamelink-backend:latest

# 漏洞报告
trivy image --format json --output report.json gamelink-backend:latest
```

---

## 附录：快速配置检查清单

### 新环境设置检查清单

**数据库：**
- [ ] PostgreSQL 16 已安装并运行
- [ ] 创建数据库和用户
- [ ] 设置强密码
- [ ] 配置连接池参数
- [ ] 测试连接

**缓存：**
- [ ] Redis 7 已安装并运行
- [ ] 设置密码
- [ ] 配置内存策略
- [ ] 测试连接

**后端：**
- [ ] 创建 `.env` 文件
- [ ] 配置数据库连接
- [ ] 配置 Redis 连接
- [ ] 生成并设置 JWT 密钥
- [ ] 根据环境设置 CRYPTO_ENABLED
- [ ] 设置超级管理员账号
- [ ] 设置 SEED_ENABLED（生产=false）

**前端：**
- [ ] 创建 `.env` 文件
- [ ] 配置 API URL
- [ ] 配置加密密钥（与后端匹配）
- [ ] 测试 API 连接

**外部服务（可选）：**
- [ ] 配置微信登录
- [ ] 配置微信支付
- [ ] 配置支付宝
- [ ] 配置短信服务
- [ ] 配置对象存储

### 常见问题排查

**问题：数据库连接失败**
```bash
# 检查 PostgreSQL 状态
systemctl status postgresql

# 测试连接
psql -h localhost -U gamelink -d gamelink

# 检查防火墙
sudo ufw status
```

**问题：Redis 连接失败**
```bash
# 检查 Redis 状态
systemctl status redis

# 测试连接
redis-cli -a password ping

# 检查配置
redis-cli CONFIG GET requirepass
```

**问题：前端加密失败**
```bash
# 检查前后端加密配置是否一致
# 后端
echo $CRYPTO_SECRET_KEY
echo $CRYPTO_IV

# 前端
echo $VITE_CRYPTO_SECRET_KEY
echo $VITE_CRYPTO_IV

# 确保完全匹配
```

**问题：Docker 容器无法访问宿主机服务**
```bash
# 使用 host.docker.internal (Mac/Windows)
# 或使用 --network host (Linux)

# 或使用宿主机 IP
hostname -I  # 查看 IP 地址
```

---

## 联系和支持

如有配置问题，请：
1. 查阅本文档相关章节
2. 检查 `.env.example` 文件注释
3. 查看 GitHub Issues
4. 联系 DevOps 团队

**文档维护：** DevOps 团队
**最后更新：** 2026-02-09
**版本：** 1.0.0
