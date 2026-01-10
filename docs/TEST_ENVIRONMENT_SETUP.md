# 测试环境设置（本地 & CI）

本文档聚焦“如何把测试跑起来”，覆盖：
- `api/`：Go 单元测试、集成测试（PostgreSQL）
- `admin/`：Vitest 单元测试、Playwright E2E

> 约定：**集成测试默认会在数据库不可用时自动跳过**（`SkipIfNoTestDB`），要真正执行集成测试需要先准备好测试数据库并正确设置 `TEST_DB_*`。

---

## 1. 先决条件

- Docker Desktop / Docker Engine（本地集成测试推荐）
- Go（版本以 `api/go.mod` 与 CI 为准）
- Node.js（版本以 `.github/workflows/ci.yml` 为准，当前为 Node 20）

---

## 2. 后端（api/）测试

### 2.1 单元测试（不依赖外部服务）

```bash
cd api
make test
```

覆盖率：

```bash
cd api
make test-coverage
```

说明：
- 会包含 `api/internal/service/integration` 下的用例，但当测试数据库不可用时会 `skip`（不会让 `make test` 失败）。

### 2.2 集成测试（需要 PostgreSQL）

集成测试通过 `TEST_DB_*` 环境变量连接 PostgreSQL（见 `api/internal/service/integration/testdb.go` 的 `DefaultTestConfig()`）。

#### 方案 A：使用仓库自带 `docker-compose.test.yml`（推荐）

1）启动测试数据库：

```powershell
docker compose -f docker-compose.test.yml up -d
docker compose -f docker-compose.test.yml ps
```

2）运行集成测试（自动设置 `TEST_DB_*`）：

```bash
cd api
make test-integration-db
```

端口说明（本地默认）：
- PostgreSQL：宿主机 `5433` → 容器 `5432`
- Redis：宿主机 `6380` → 容器 `6379`（当前集成测试主要使用 PostgreSQL，但保留 Redis 便于后续用例扩展）

> 注意：`docker-compose.yml`（开发环境）也会占用 `5433/6380`，两者不能同时启动；如已启动开发环境，请先 `docker compose down` 或改端口后再启动测试环境。

#### 方案 B：手动设置 `TEST_DB_*`（适合自建数据库/CI 对齐）

PowerShell 示例（CI 同款端口通常是 `5432`；本地使用 `docker-compose.test.yml` 时端口是 `5433`）：

```powershell
$env:TEST_DB_HOST = "localhost"
$env:TEST_DB_PORT = "5432"
$env:TEST_DB_USER = "gamelink"
$env:TEST_DB_PASSWORD = "gamelink"
$env:TEST_DB_NAME = "gamelink_test"

cd api
go test ./internal/service/integration/... -v -count=1
```

### 2.3 （可选）初始化演示数据（用于联调/UI/E2E）

后端支持通过 `SEED_ENABLED=true` 启动时自动写入演示数据（订单/支付/优惠券/VIP/活动/通知等），便于联调与手工验证：

```powershell
$env:SEED_ENABLED = "true"
cd api
go run cmd/main.go
```

更多说明见 `docs/backend/database-seed.md`。

---

## 3. 管理后台（admin/）测试

### 3.1 单元测试（Vitest）

```bash
cd admin
npm ci
npm run test:run
```

### 3.2 E2E（Playwright）

Playwright 配置会自动启动管理后台开发服务器（见 `admin/playwright.config.ts` 的 `webServer.command = npm run dev`），因此你需要先确保：
- 后端 API 可访问：默认 `http://localhost:8080/api/v1`
- 有可用的管理员账号用于登录（建议使用 `SUPER_ADMIN_EMAIL/SUPER_ADMIN_PASSWORD` 对应的账号）

推荐启动后端方式（本地 Docker 开发环境）：

```powershell
docker compose up -d backend postgres redis
docker compose ps
```

设置 E2E 环境变量（建议显式设置，避免默认值与本地账号不一致）：

```powershell
$env:API_URL = "http://localhost:8080/api/v1"
$env:BASE_URL = "http://localhost:5173"
$env:TEST_ADMIN_USERNAME = "admin@gamelink.com"
$env:TEST_ADMIN_PASSWORD = "admin123456"
```

> 说明：若你启用了 `SEED_ENABLED=true`，也可以使用演示管理员账号登录（例如 `sysadmin@gamelink.com / Admin@123456`），以 seed 数据为准。

运行 E2E：

```bash
cd admin
npm run test:e2e:install
npm run test:e2e
```

---

## 4. 常见问题（FAQ）

### 4.1 集成测试全部被跳过（Skipped）

- 确认 PostgreSQL 已启动（`docker compose -f docker-compose.test.yml ps`）
- 确认 `TEST_DB_PORT` 与实际宿主机端口一致
  - 使用 `docker-compose.test.yml`：默认是 `5433`
  - CI services：通常是 `5432`

### 4.2 E2E 登录失败

- 确认 `TEST_ADMIN_USERNAME/TEST_ADMIN_PASSWORD` 与后端实际初始化的超管一致
  - `docker-compose.yml`（开发环境）默认：`admin@gamelink.com / admin123456`
  - `docker-compose.prod*.yml` 依赖 `.env` 中的 `SUPER_ADMIN_*`
  - 若启用 `SEED_ENABLED=true`，也可使用 seed 内置演示管理员账号（见 `docs/backend/database-seed.md`）
