# Docker 环境配置说明

## 环境概览

| 环境 | 配置文件 | 用途 |
|------|----------|------|
| 开发环境 | `docker-compose.yml` | 本地开发（仅数据库服务） |
| 生产环境 | `docker-compose.prod.yml` | 生产部署（完整服务栈） |

## 环境差异对比

| 配置项 | 开发环境 | 生产环境 |
|--------|----------|----------|
| 服务 | PostgreSQL + Redis | PostgreSQL + Redis + Backend + Admin |
| `CRYPTO_ENABLED` | false | true |
| `SEED_ENABLED` | true | false |
| `GIN_MODE` | debug | release |
| PostgreSQL 端口 | 5433 | 5432 |
| Redis 端口 | 6380 | 6379 |
| 性能调优 | 无 | 有 |
| 资源限制 | 无 | 有 |

## 开发环境

开发环境仅部署数据库服务，后端和前端在本地运行。

### 启动

```bash
# 启动数据库服务
docker-compose up -d

# 查看状态
docker-compose ps
```

### 本地运行应用

```bash
# 后端 (新终端)
cd api
go run cmd/main.go

# 前端 (新终端)
cd admin
npm run dev
```

### 停止

```bash
docker-compose down

# 清理数据卷
docker-compose down -v
```

## 生产环境

生产环境部署完整服务栈。

### 准备

```bash
# 创建生产环境配置
cp .env.example .env.production

# 编辑配置，设置强密码和加密密钥
# 必须设置: POSTGRES_PASSWORD, REDIS_PASSWORD, JWT_SECRET_KEY, 
#          CRYPTO_SECRET_KEY, CRYPTO_IV, SUPER_ADMIN_PASSWORD
```

### 启动

```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

### 查看状态

```bash
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f backend
```

### 停止

```bash
docker-compose -f docker-compose.prod.yml down
```

## 端口映射

| 服务 | 开发环境 | 生产环境 |
|------|----------|----------|
| PostgreSQL | 5433 | 5432 |
| Redis | 6380 | 6379 |
| Backend | 本地 8080 | 8080 |
| Admin | 本地 5173 | 80/443 |

## 生成密钥

```bash
# JWT 密钥
openssl rand -base64 32

# 加密密钥 (32字节)
openssl rand -base64 32

# 加密 IV (16字节)
openssl rand -base64 16

# 管理员密码
openssl rand -base64 24
```

## 健康检查

- **PostgreSQL**: `pg_isready` 命令
- **Redis**: `redis-cli ping` 命令
- **Backend**: `GET /api/v1/healthz`
- **Admin**: `GET /health`

## 故障排查

```bash
# 查看容器日志
docker-compose logs -f postgres
docker-compose logs -f redis

# 进入数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

# 进入 Redis
docker exec -it gamelink-redis redis-cli
```
