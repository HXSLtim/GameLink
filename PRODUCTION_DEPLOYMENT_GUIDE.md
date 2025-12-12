# GameLink 生产环境部署指南

## ✅ 已完成的准备工作

1. ✅ 生成了安全的生产环境配置文件 `.env`
2. ✅ 后端 Docker 镜像构建成功
3. ✅ PostgreSQL 和 Redis 容器已启动

## 🔑 生产环境密钥

已为您生成安全的密钥并保存在 `.env` 文件中：

- **数据库密码**: 已设置强密码
- **Redis 密码**: 已设置强密码  
- **JWT 密钥**: 已生成安全密钥
- **加密密钥**: 已生成 24 字符密钥
- **管理员密码**: 已生成强密码

**重要**: 请妥善保管 `.env` 文件，不要提交到版本控制系统！

## 🚀 完成部署步骤

### 步骤 1: 重启 Docker Desktop

Docker Desktop 服务需要重启：

1. 打开 Docker Desktop 应用
2. 点击右上角设置图标
3. 选择 "Restart Docker Desktop"
4. 等待 Docker 完全启动（图标变绿）

### 步骤 2: 启动生产环境

```powershell
# 方式 1: 使用部署脚本（推荐）
.\scripts\deploy-production.ps1

# 方式 2: 手动启动
docker-compose -f docker-compose.prod.yml up -d postgres redis backend
```

### 步骤 3: 等待服务就绪

```powershell
# 等待约 30 秒，然后检查服务状态
Start-Sleep -Seconds 30

# 查看容器状态
docker ps --filter "name=gamelink"

# 查看后端日志
docker logs gamelink-backend --tail=50
```

### 步骤 4: 健康检查

```powershell
# 使用健康检查脚本
.\scripts\docker-health-check.ps1 -Environment prod

# 或手动检查
curl http://localhost:8080/api/v1/health
```

### 步骤 5: 访问应用

- **后端 API**: http://localhost:8080
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **管理员账号**: 
  - 邮箱: `admin@gamelink.com`
  - 密码: 查看 `.env` 文件中的 `SUPER_ADMIN_PASSWORD`

## 📊 服务架构

```
生产环境部署架构:
┌─────────────────────────────────────┐
│         GameLink Backend            │
│         (Port 8080)                 │
└──────────┬──────────────────────────┘
           │
           ├─────────────┐
           │             │
┌──────────▼──────┐  ┌──▼──────────┐
│   PostgreSQL    │  │    Redis    │
│   (Port 5432)   │  │ (Port 6379) │
└─────────────────┘  └─────────────┘
```

## 🔧 故障排查

### 问题 1: 后端容器一直重启

**原因**: 可能是数据库连接失败或配置错误

**解决方案**:
```powershell
# 1. 查看后端日志
docker logs gamelink-backend --tail=100

# 2. 检查数据库是否就绪
docker exec gamelink-postgres pg_isready -U gamelink

# 3. 检查环境变量
docker exec gamelink-backend env | Select-String "DB_DSN"

# 4. 重启后端
docker-compose -f docker-compose.prod.yml restart backend
```

### 问题 2: 数据库连接失败

**解决方案**:
```powershell
# 1. 检查 PostgreSQL 日志
docker logs gamelink-postgres --tail=50

# 2. 测试数据库连接
docker exec gamelink-postgres psql -U gamelink -d gamelink -c "SELECT 1"

# 3. 重启数据库
docker-compose -f docker-compose.prod.yml restart postgres
```

### 问题 3: Redis 连接失败

**解决方案**:
```powershell
# 1. 检查 Redis 状态
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> PING

# 2. 查看 Redis 日志
docker logs gamelink-redis --tail=50

# 3. 重启 Redis
docker-compose -f docker-compose.prod.yml restart redis
```

## 📝 常用管理命令

### 查看服务状态
```powershell
# 查看所有容器
docker ps --filter "name=gamelink"

# 查看资源使用
docker stats --no-stream --filter "name=gamelink"

# 健康检查
.\scripts\docker-health-check.ps1 -Environment prod
```

### 查看日志
```powershell
# 查看后端日志
docker logs gamelink-backend -f

# 查看数据库日志
docker logs gamelink-postgres -f

# 查看 Redis 日志
docker logs gamelink-redis -f
```

### 重启服务
```powershell
# 重启所有服务
docker-compose -f docker-compose.prod.yml restart

# 重启后端
docker-compose -f docker-compose.prod.yml restart backend

# 重启数据库
docker-compose -f docker-compose.prod.yml restart postgres
```

### 停止服务
```powershell
# 停止所有服务
docker-compose -f docker-compose.prod.yml down

# 停止但保留数据
docker-compose -f docker-compose.prod.yml stop
```

## 💾 数据备份

### 自动备份
```powershell
# 使用备份脚本
.\scripts\docker-backup.ps1 -Environment prod
```

### 手动备份
```powershell
# 备份 PostgreSQL
docker exec gamelink-postgres pg_dump -U gamelink gamelink > backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql

# 备份 Redis
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> SAVE
docker cp gamelink-redis:/data/dump.rdb redis_backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').rdb
```

## 🔐 安全建议

### 1. 修改默认端口（可选）

编辑 `docker-compose.prod.yml`：
```yaml
services:
  backend:
    ports:
      - "8080:8080"  # 改为其他端口，如 "8888:8080"
```

### 2. 配置防火墙

```powershell
# 仅允许特定 IP 访问
# 使用 Windows 防火墙或云服务商的安全组配置
```

### 3. 启用 HTTPS（推荐）

使用 Nginx 反向代理并配置 SSL 证书：
```yaml
# 添加 Nginx 服务到 docker-compose.prod.yml
nginx:
  image: nginx:alpine
  ports:
    - "443:443"
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
    - ./ssl:/etc/nginx/ssl
```

### 4. 定期更新

```powershell
# 拉取最新代码
git pull

# 重新构建
docker-compose -f docker-compose.prod.yml build

# 重启服务
docker-compose -f docker-compose.prod.yml up -d
```

### 5. 监控和日志

```powershell
# 定期检查日志
docker logs gamelink-backend --tail=100 | Select-String "ERROR"

# 监控资源使用
docker stats --no-stream
```

## 📈 性能优化

### 1. 数据库优化

```sql
-- 连接到数据库
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

-- 创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
```

### 2. Redis 优化

```powershell
# 配置 Redis 最大内存
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> CONFIG SET maxmemory 256mb
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> CONFIG SET maxmemory-policy allkeys-lru
```

### 3. 容器资源限制

编辑 `docker-compose.prod.yml`：
```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

## 🎯 下一步

### 前端部署 ✅ 已修复

前端 TypeScript 错误已全部修复！现在可以部署：

```powershell
# 构建前端镜像
docker-compose -f docker-compose.prod.yml build frontend

# 启动前端服务
docker-compose -f docker-compose.prod.yml up -d frontend

# 或一次性部署所有服务
docker-compose -f docker-compose.prod.yml up -d --build
```

**修复详情**: 查看 [FRONTEND_BUILD_FIXED.md](FRONTEND_BUILD_FIXED.md)

### 监控系统（可选）

添加 Prometheus 和 Grafana 监控：
```powershell
# 使用监控配置
docker-compose -f docker-compose.prod.yml -f docker-compose.monitoring.yml up -d
```

### CI/CD 集成（可选）

配置 GitHub Actions 自动部署：
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy
        run: |
          docker-compose -f docker-compose.prod.yml build
          docker-compose -f docker-compose.prod.yml up -d
```

## 📞 获取帮助

1. 查看日志: `docker logs gamelink-backend`
2. 健康检查: `.\scripts\docker-health-check.ps1 -Environment prod`
3. 查阅文档: [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)
4. 提交 Issue 到项目仓库

---

**部署日期**: 2025-12-13  
**版本**: 1.0.0  
**状态**: 后端已构建，等待 Docker Desktop 重启后完成部署
