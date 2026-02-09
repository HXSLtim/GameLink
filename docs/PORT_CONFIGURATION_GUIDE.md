# GameLink 端口配置指南

**文档版本：** 1.0.0
**创建日期：** 2026-02-09
**维护人：** DevOps-Engineer

---

## 端口分配规范

### 开发环境

| 服务 | 端口 | 配置文件 | 说明 |
|------|------|---------|------|
| **后端 API** | 8080 | `.env` (BACKEND_PORT) | Go API 服务器 |
| **Admin 前端** | 5173 | `admin/vite.config.ts` (server.port) | Vite 开发服务器 |
| **App 前端** | 3000 | `app/vite.config.ts` (server.port) | UniApp/Vite 开发服务器 |
| **PostgreSQL** | 5433 | `.env` (POSTGRES_PORT) | 数据库 |
| **Redis** | 6380 | `.env` (REDIS_PORT) | 缓存 |

### 生产环境

| 服务 | 端口 | 暴露端口 | 说明 |
|------|------|---------|------|
| **Nginx** | 80/443 | 80, 443 | HTTP/HTTPS 入口 |
| **后端 API** | 8080 | 容器内部 | Docker 网络内部 |
| **Admin 静态文件** | - | 80 (Nginx) | Nginx 托管 |
| **PostgreSQL** | 5432 | 容器内部 | Docker 网络内部 |
| **Redis** | 6379 | 容器内部 | Docker 网络内部 |

---

## 开发环境启动顺序

### 1. 启动基础设施（Docker）

```bash
# 启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis
```

### 2. 启动后端 API

```bash
cd api
go run cmd/main.go
# 后端运行在 http://localhost:8080
```

### 3. 启动 Admin 前端

```bash
cd admin
npm run dev
# Admin 运行在 http://localhost:5173
# API 请求通过代理转发到 http://localhost:8080
```

### 4. 启动 App 前端

```bash
cd app
npm run dev:h5
# App 运行在 http://localhost:3000
# API 请求需要配置 VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## 端口冲突解决

### 检查端口占用

**Windows:**
```bash
netstat -ano | findstr ":8080"
netstat -ano | findstr ":5173"
netstat -ano | findstr ":3000"
```

**Linux/Mac:**
```bash
lsof -i :8080
lsof -i :5173
lsof -i :3000
```

### 解决冲突

**如果端口被占用：**

1. **查找占用进程**
   ```bash
   # Windows
   netstat -ano | findstr ":8080"
   tasklist | findstr <PID>

   # Linux/Mac
   lsof -i :8080
   ps aux | grep <PID>
   ```

2. **结束进程**
   ```bash
   # Windows
   taskkill /PID <PID> /F

   # Linux/Mac
   kill -9 <PID>
   ```

3. **或者修改端口配置**
   - 后端: `.env` 中的 `BACKEND_PORT`
   - Admin: `admin/vite.config.ts` 中的 `server.port`
   - App: `app/vite.config.ts` 中的 `server.port`

---

## BUG-001 解决方案

### 问题

移动端前端（App）占用 8080 端口，导致后端 API 无法启动。

### 根本原因

App 的 Vite 配置中虽然指定了 `port: 3000`，但实际运行时使用了 8080。可能原因：
- 配置未生效
- 有其他配置覆盖
- 缓存问题

### 解决方案

**方案 A：强制 App 使用 3000 端口（推荐）**

修改 `app/vite.config.ts`:
```typescript
server: {
  host: '0.0.0.0',
  port: 3000,
  strictPort: true,  // 添加此行，端口被占用时报错
},
```

重启 App 开发服务器：
```bash
cd app
# 结束占用 8080 的进程
# 重新启动
npm run dev:h5
```

**方案 B：后端改用其他端口（不推荐）**

需要修改多处配置，影响较大。

---

## API 代理配置

### Admin 代理

`admin/vite.config.ts`:
```typescript
server: {
  proxy: {
    '/api/v1': {
      target: 'http://localhost:8080',  // 后端 API
      changeOrigin: true,
      ws: true,  // 支持 WebSocket
    },
  },
}
```

### App API 配置

`app/.env.development`:
```bash
# 留空则使用默认配置
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/api/v1/ws
```

---

## Docker Compose 端口映射

### 开发环境

`docker-compose.yml`:
```yaml
services:
  postgres:
    ports:
      - "5433:5432"  # 宿主机:容器

  redis:
    ports:
      - "6380:6379"
```

### 生产环境

`docker-compose.prod.yml`:
```yaml
services:
  postgres:
    # 不暴露端口，仅容器内部访问
    # ports: []

  redis:
    # 不暴露端口，仅容器内部访问

  backend:
    # 不暴露端口，通过 Nginx 代理
    # ports: []
```

---

## 验证配置

### 1. 检查所有服务运行状态

```bash
# 检查后端
curl http://localhost:8080/api/v1/healthz

# 检查 Admin
curl http://localhost:5173

# 检查 App
curl http://localhost:3000

# 检查 PostgreSQL
docker exec gamelink-postgres pg_isready -U gamelink

# 检查 Redis
docker exec gamelink-redis redis-cli ping
```

### 2. 检查端口占用

```bash
# 运行端口检查脚本
bash scripts/check-ports.sh
```

---

## 常见问题

### Q1: 修改配置后端口没有生效

**解决方案：**
1. 完全结束开发服务器（Ctrl+C）
2. 检查端口是否被占用
3. 清理缓存：`rm -rf node_modules/.vite`
4. 重新启动

### Q2: 多个开发环境端口冲突

**解决方案：**
- 每个开发者使用不同的端口范围
- 或使用 Docker Compose 隔离环境

### Q3: 生产环境端口配置

**解决方案：**
- 生产环境使用 Nginx 统一入口
- 所有服务仅容器内部通信
- 不对外暴露非必要端口

---

## 相关文档

- **部署检查清单：** `docs/DEPLOYMENT_CHECKLIST.md`
- **依赖配置文档：** `docs/DEPENDENCIES_AND_CONFIG.md`
- **安全加固指南：** `docs/SECURITY_HARDENING.md`

---

**更新历史：**

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-02-09 | 1.0.0 | 初始版本，BUG-001 解决方案 | DevOps-Engineer |
