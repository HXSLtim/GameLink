# 配置说明（configs）

本目录包含不同环境的服务配置文件：

- `config.development.yaml` — 开发环境配置（默认）
- `config.production.yaml` — 生产环境配置

服务启动时会根据环境变量 `APP_ENV` 加载 `configs/config.<env>.yaml`，并由环境变量进行覆盖，最终生成运行配置。

## 配置结构

YAML 文件可包含以下字段（与 `internal/config/env.go` 对应）：

```yaml
server:
  port: "8080"
  enable_swagger: true

database:
  type: sqlite
  dsn: "file:./var/dev.db?mode=rwc&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"

cache:
  type: memory
  redis:
    addr: "127.0.0.1:6379"
    password: ""
    db: 0

crypto:
  enabled: false
  secret_key: "GameLink2025SecretKey!@#"
  iv: "GameLink2025IV!!!"
  methods: ["POST", "PUT", "PATCH"]
  exclude_paths: ["/api/v1/health", "/api/v1/ping", "/api/v1/auth/refresh"]
  use_signature: true

auth:
  jwt_secret: ""
  token_ttl_hours: 24

seed:
  enabled: false
```

## 环境变量覆盖

加载顺序为：默认值 → 文件配置 → 环境变量覆盖。可用的环境变量如下：

- `APP_ENV` — 选择配置文件（development/production），默认 `development`
- `SERVICE_PORT` — 覆盖 `server.port`
- `ENABLE_SWAGGER` — 覆盖 `server.enable_swagger`（true/false）
- `DB_TYPE` — 覆盖 `database.type`（sqlite/postgres/mysql/sqlserver）
- `DB_DSN` — 覆盖 `database.dsn`
- `CACHE_TYPE` — 覆盖 `cache.type`（memory/redis）
- `REDIS_ADDR` — 覆盖 `cache.redis.addr`
- `REDIS_PASSWORD` — 覆盖 `cache.redis.password`
- `REDIS_DB` — 覆盖 `cache.redis.db`（整数）
- `CRYPTO_ENABLED` — 启用/禁用加密（true/false）
- `CRYPTO_SECRET_KEY` — 加密密钥（长度必须为 16/24/32 字节）
- `CRYPTO_IV` — 初始化向量（至少 16 字节）
- `CRYPTO_METHODS` — 逗号分隔的 HTTP 方法列表（例如：`POST,PUT,PATCH`）
- `CRYPTO_EXCLUDE_PATHS` — 逗号分隔的排除路径（例如：`/api/v1/health,/api/v1/ping`）
- `CRYPTO_USE_SIGNATURE` — 是否启用签名（true/false）
- `JWT_SECRET_KEY` — JWT 密钥（生产环境必须提供）
- `JWT_TOKEN_TTL_HOURS` — Token 有效期小时数（整数）
- `SEED_ENABLED` — 是否注入演示数据（true/false）

## 校验与默认值

- 在生产环境下：
  - `DB_DSN` 必须提供，否则会报错。
  - 开启加密时：`CRYPTO_SECRET_KEY` 长度必须为 16/24/32，`CRYPTO_IV` 至少 16 字节，`CRYPTO_METHODS` 不能为空。
  - 未显式提供 `JWT_SECRET_KEY` 时，不会自动填充（请通过配置或环境变量提供）。
- 在开发环境下：
  - 若 `DB_DSN` 为空，会根据 `DB_TYPE` 自动填充示例 DSN（日志可见）。
  - 若 `JWT_SECRET_KEY` 为空，会降级使用开发兜底密钥（仅用于本地调试）。

## 使用建议

- 建议以文件配置为主，环境变量用于容器/CI 下的动态覆盖。
- 不要将敏感配置（如 `JWT_SECRET_KEY`、数据库密码）写入仓库，生产环境请通过环境变量或安全的配置注入。
- 如启用 `crypto`，确保前后端密钥/IV/签名策略一致，并将不需要加密的健康检查等路径加入 `exclude_paths`。

## 示例启动

- 本地开发（默认 development）：直接运行服务会读取 `config.development.yaml`，并按需由环境变量覆盖。
- 生产环境：设置 `APP_ENV=production` 并提供必要的环境变量，例如：

```powershell
$env:APP_ENV = "production"
$env:DB_TYPE = "postgres"
$env:DB_DSN = "postgres://user:password@db:5432/gamelink?sslmode=disable"
$env:JWT_SECRET_KEY = "<your-secret>"
# 启动服务
# go run ./cmd/user-service
```

## 联系方式

如需沟通或反馈，请通过邮箱：a2778978136@163.com