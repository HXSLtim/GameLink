# GameLink Docker 部署指南

本文档提供 GameLink 项目的 Docker 容器化部署完整指南。

## 📋 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [开发环境部署](#开发环境部署)
- [生产环境部署](#生产环境部署)
- [常用命令](#常用命令)
- [健康检查](#健康检查)
- [故障排查](#故障排查)
- [性能优化](#性能优化)

## 前置要求

### 必需软件

- **Docker**: 20.10.0 或更高版本
- **Docker Compose**: 2.0.0 或更高版本

### 验证安装

```powershell
# 检查 Docker 版本
docker --version

# 检查 Docker Compose 版本
docker-compose --version

# 验证 Docker 服务运行状态
docker info
```

## 快速开始

### 1. 克隆项目

```powershell
git clone <repository-url>
cd GameLink
```

### 2. 开发环境快速启动

```powershell
# 启动开发环境（使用 SQLite + 内存缓存）
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 3. 访问应用

- **前端**: http://localhost
- **后端API**: http://localhost:8080
- **Swagger文档**: http://localhost:8080/swagger/index.html
- **默认管理员账号**:
  - 邮箱: admin@gamelink.com
  - 密码: 123456

## 开发环境部署

### 架构说明

开发环境使用简化配置，便于快速开发和调试：

- **数据库**: SQLite（文件存储）
- **缓存**: 内存缓存
- **配置**: `config.development.yaml`

### 启动服务

```powershell
# 启动所有服务
docker-compose up -d

# 仅启动后端
docker-compose up -d backend

# 仅启动前端
docker-compose up -d frontend
```

### 查看日志

```powershell
# 查看所有服务日志
docker-compose logs -f

# 查看后端日志
docker-compose logs -f backend

# 查看前端日志
docker-compose logs -f frontend

# 查看最近100行日志
docker-compose logs --tail=100 -f
```

### 停止服务

```powershell
# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v

# 停止并删除镜像
docker-compose down --rmi all
```

## 生产环境部署

### 架构说明

生产环境使用完整的技术栈：

- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **配置**: `config.production.yaml`

### 1. 配置环境变量

```powershell
# 复制环境变量模板
Copy-Item .env.example .env

# 编辑 .env 文件，设置生产环境配置
notepad .env
```

**重要配置项**:

```env
# PostgreSQL 配置
POSTGRES_DB=gamelink
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=<强密码>

# Redis 配置
REDIS_PASSWORD=<强密码>

# JWT 密钥（必须修改）
JWT_SECRET_KEY=<至少32字符的随机字符串>

# 加密配置
CRYPTO_SECRET_KEY=<24字符密钥>
CRYPTO_IV=<16字符IV>

# 管理员配置
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=<强密码>
```

### 2. 生成密钥

#### Windows PowerShell

```powershell
# 生成 JWT 密钥（32字节）
$jwt = [Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
Write-Host "JWT_SECRET_KEY=$jwt"

# 生成加密密钥（24字符）
$crypto = -join ((65..90) + (97..122) + (48..57) | Get-Random -Count 24 | ForEach-Object {[char]$_})
Write-Host "CRYPTO_SECRET_KEY=$crypto"

# 生成加密IV（16字符）
$iv = -join ((65..90) + (97..122) + (48..57) | Get-Random -Count 16 | ForEach-Object {[char]$_})
Write-Host "CRYPTO_IV=$iv"
```

#### Linux/macOS

```bash
# 生成 JWT 密钥
echo "JWT_SECRET_KEY=$(openssl rand -base64 32)"

# 生成加密密钥（24字符）
echo "CRYPTO_SECRET_KEY=$(openssl rand -base64 24 | cut -c1-24)"

# 生成加密IV（16字符）
echo "CRYPTO_IV=$(openssl rand -base64 16 | cut -c1-16)"
```

### 3. 启动生产环境

```powershell
# 构建并启动所有服务
docker-compose -f docker-compose.prod.yml up -d --build

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 4. 数据库初始化验证

```powershell
# 进入后端容器
docker exec -it gamelink-backend sh

# 查看日志确认数据库初始化
docker-compose -f docker-compose.prod.yml logs backend | grep -i "migration\|seed"
```

## 本地生产环境测试

### 快速启动

使用本地生产环境配置（PostgreSQL + Redis）进行测试：

```powershell
# 启动本地生产环境
.\scripts\docker-prod-local-start.ps1

# 清理数据重新开始
.\scripts\docker-prod-local-start.ps1 -Clean

# 跳过镜像构建（使用现有镜像）
.\scripts\docker-prod-local-start.ps1 -NoBuild
```

### 访问地址

- **后端API**: http://localhost:8081
- **Swagger文档**: http://localhost:8081/swagger/index.html
- **PostgreSQL**: localhost:5433
- **Redis**: localhost:6380

端口已调整避免与本地开发服务冲突。

## 常用命令

### 快捷脚本

```powershell
# 健康检查
.\scripts\docker-health-check.ps1 -Environment prod-local

# 查看日志
.\scripts\docker-logs.ps1 -Environment prod-local -Service backend -Follow

# 数据备份
.\scripts\docker-backup.ps1 -Environment prod-local

# 数据恢复
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod-local

# 清理环境
.\scripts\docker-clean.ps1 -Level soft -Environment prod-local
```

### 容器管理

```powershell
# 查看运行中的容器
docker-compose ps

# 重启服务
docker-compose restart

# 重启特定服务
docker-compose restart backend

# 查看容器资源使用情况
docker stats
```

### 镜像管理

```powershell
# 查看镜像
docker images

# 重新构建镜像
docker-compose build

# 强制重新构建（不使用缓存）
docker-compose build --no-cache

# 删除未使用的镜像
docker image prune -a
```

### 数据卷管理

```powershell
# 查看数据卷
docker volume ls

# 查看数据卷详情
docker volume inspect gamelink-postgres-data

# 备份数据卷
docker run --rm -v gamelink-postgres-data:/data -v ${PWD}:/backup alpine tar czf /backup/postgres-backup.tar.gz /data
```

### 进入容器

```powershell
# 进入后端容器
docker exec -it gamelink-backend sh

# 进入前端容器
docker exec -it gamelink-frontend sh

# 进入数据库容器
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

# 进入 Redis 容器
docker exec -it gamelink-redis redis-cli -a <REDIS_PASSWORD>
```

## 健康检查

### 自动健康检查

所有服务都配置了健康检查端点：

- **后端**: `GET /api/v1/health`
- **前端**: `GET /health`
- **PostgreSQL**: `pg_isready`
- **Redis**: `redis-cli ping`

### 手动健康检查

```powershell
# 检查后端健康
curl http://localhost:8080/api/v1/health

# 检查前端健康
curl http://localhost/health

# 查看所有服务健康状态
docker-compose ps
```

## 故障排查

### 服务无法启动

```powershell
# 查看详细日志
docker-compose logs backend

# 检查配置文件
docker exec gamelink-backend cat /app/configs/config.production.yaml

# 检查环境变量
docker exec gamelink-backend env | grep APP_ENV
```

### 数据库连接失败

```powershell
# 检查 PostgreSQL 是否就绪
docker exec gamelink-postgres pg_isready -U gamelink

# 查看数据库日志
docker-compose logs postgres

# 测试数据库连接
docker exec gamelink-postgres psql -U gamelink -d gamelink -c "SELECT 1"
```

### Redis 连接失败

```powershell
# 检查 Redis 状态
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> ping

# 查看 Redis 日志
docker-compose logs redis

# 检查 Redis 配置
docker exec gamelink-redis redis-cli -a <REDIS_PASSWORD> CONFIG GET requirepass
```

### 端口冲突

```powershell
# 查看端口占用情况
netstat -ano | findstr :8080
netstat -ano | findstr :80
netstat -ano | findstr :5432

# 修改 docker-compose.yml 中的端口映射
# 例如: "8081:8080" 将主机端口改为 8081
```

### 前端无法访问后端

1. 检查 nginx 配置中的后端地址
2. 确认后端容器名称为 `backend`
3. 检查网络配置

```powershell
# 查看网络
docker network ls

# 检查容器网络连接
docker network inspect gamelink-network
```

## 性能优化

### 构建优化

```powershell
# 使用构建缓存加速构建
docker-compose build

# 多阶段构建已默认启用
# 参考 Dockerfile 中的 build 和 runtime 阶段
```

### 资源限制

在 `docker-compose.prod.yml` 中添加资源限制：

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

### 日志管理

配置日志轮转：

```yaml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 数据备份与恢复

### 自动备份（推荐）

使用提供的备份脚本进行完整备份：

```powershell
# 备份本地生产环境
.\scripts\docker-backup.ps1 -Environment prod-local

# 备份生产环境
.\scripts\docker-backup.ps1 -Environment prod

# 指定备份目录
.\scripts\docker-backup.ps1 -Environment prod-local -BackupDir "D:\backups"
```

备份内容包括：
- PostgreSQL 数据库完整导出
- Redis RDB 快照
- 环境配置文件
- 自动压缩为 .zip 文件

### 数据恢复

```powershell
# 从备份恢复
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod-local

# 恢复到生产环境
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod
```

### 手动备份（高级）

#### PostgreSQL 备份

```powershell
# 备份数据库
docker exec gamelink-postgres pg_dump -U gamelink gamelink > backup.sql

# 恢复数据库
docker exec -i gamelink-postgres psql -U gamelink gamelink < backup.sql

# 备份到压缩文件
docker exec gamelink-postgres pg_dump -U gamelink gamelink | gzip > backup.sql.gz
```

#### Redis 备份

```powershell
# 创建 RDB 快照
docker exec gamelink-redis redis-cli -a redis123 SAVE

# 复制备份文件
docker cp gamelink-redis:/data/dump.rdb ./redis-backup.rdb

# 恢复 Redis 数据
docker cp ./redis-backup.rdb gamelink-redis:/data/dump.rdb
docker-compose restart redis
```

## 更新部署

### 滚动更新

```powershell
# 拉取最新代码
git pull

# 重新构建并启动（零停机更新）
docker-compose -f docker-compose.prod.yml up -d --build --no-deps backend
docker-compose -f docker-compose.prod.yml up -d --build --no-deps frontend
```

### 完全重启

```powershell
# 停止所有服务
docker-compose -f docker-compose.prod.yml down

# 重新构建并启动
docker-compose -f docker-compose.prod.yml up -d --build
```

## 安全建议

1. **修改默认密码**: 生产环境必须修改所有默认密码
2. **使用强密钥**: JWT 和加密密钥必须使用随机生成的强密钥
3. **限制访问**: 使用防火墙限制数据库和 Redis 端口访问
4. **HTTPS**: 生产环境使用 HTTPS（配置 SSL 证书）
5. **定期备份**: 建立定期数据备份策略
6. **监控日志**: 定期检查应用日志和容器日志
7. **更新镜像**: 定期更新基础镜像以修复安全漏洞

## 监控和日志

### 集成 Prometheus 和 Grafana（可选）

参考 `docker-compose.monitoring.yml`（需要创建）配置监控服务。

### 日志聚合（可选）

使用 ELK Stack 或 Loki 进行日志聚合和分析。

## 脚本工具说明

### docker-prod-local-start.ps1

本地生产环境启动脚本，使用完整的生产技术栈（PostgreSQL + Redis）。

**参数**:
- `-Clean`: 清理所有数据重新开始
- `-NoBuild`: 跳过镜像构建

**特点**:
- 自动检查 Docker 环境
- 检测端口冲突
- 等待服务就绪
- 显示详细的访问信息

### docker-health-check.ps1

健康检查工具，检查所有服务状态。

**参数**:
- `-Environment`: dev, prod-local 或 prod

**检查项**:
- 容器运行状态
- HTTP 服务健康端点
- 数据库连接
- Redis 连接
- 网络配置
- 资源使用情况

### docker-backup.ps1

数据备份工具，自动备份数据库和配置。

**参数**:
- `-Environment`: prod-local 或 prod
- `-BackupDir`: 备份目录路径

**备份内容**:
- PostgreSQL 完整导出
- Redis RDB 快照
- 环境配置文件
- 自动压缩打包

### docker-restore.ps1

数据恢复工具，从备份文件恢复数据。

**参数**:
- `-BackupFile`: 备份文件路径（.zip）
- `-Environment`: prod-local 或 prod

**安全特性**:
- 需要确认操作
- 自动解压备份
- 验证容器状态

### docker-logs.ps1

日志查看工具，方便查看和过滤日志。

**参数**:
- `-Environment`: dev, prod-local 或 prod
- `-Service`: backend, frontend, postgres, redis 或 all
- `-Follow`: 持续跟踪日志
- `-Lines`: 显示行数（默认50）

### docker-clean.ps1

清理工具，清理容器、镜像和数据。

**参数**:
- `-Level`: soft, medium 或 hard
- `-Environment`: dev, prod-local, prod 或 all

**清理级别**:
- `soft`: 仅停止容器
- `medium`: 删除容器和镜像
- `hard`: 删除所有包括数据（需要确认）

## 常见问题

### Q: 如何切换环境？

A: 使用不同的 docker-compose 文件和启动脚本：
- 开发环境: `docker-compose up -d` 或 `.\scripts\docker-dev-start.ps1`
- 本地生产测试: `.\scripts\docker-prod-local-start.ps1`
- 生产环境: `.\scripts\docker-prod-start.ps1`

### Q: 端口冲突怎么办？

A: 本地生产环境使用不同端口避免冲突：
- 后端: 8081（开发环境用8080）
- PostgreSQL: 5433（标准5432）
- Redis: 6380（标准6379）

可以在 `docker-compose.prod.local.yml` 中修改端口映射。

### Q: 如何查看服务状态？

A: 使用健康检查脚本：
```powershell
.\scripts\docker-health-check.ps1 -Environment prod-local
```

### Q: 数据丢失了怎么办？

A: 使用备份恢复：
```powershell
# 查看备份文件
ls .\backups\

# 恢复数据
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip"
```

### Q: 如何清理测试数据？

A: 使用清理脚本：
```powershell
# 软清理（仅停止）
.\scripts\docker-clean.ps1 -Level soft

# 完全清理（删除数据）
.\scripts\docker-clean.ps1 -Level hard
```

### Q: 如何扩展服务？

A: 使用 `docker-compose up -d --scale backend=3` 扩展服务实例。

### Q: 如何查看容器内文件？

A: 使用 `docker exec -it <container> sh` 进入容器，或使用 `docker cp` 复制文件。

### Q: 如何连接数据库？

A: 使用数据库客户端连接：
```
主机: localhost
端口: 5433 (本地生产) 或 5432 (生产)
数据库: gamelink
用户名: gamelink
密码: gamelink123 (本地测试) 或环境变量中的密码
```

## 支持与反馈

如遇到问题，请：

1. 查看日志：`docker-compose logs -f`
2. 检查服务状态：`docker-compose ps`
3. 查阅本文档的故障排查部分
4. 提交 Issue 到项目仓库

---

**最后更新**: 2025-12-13
**版本**: 1.0.0
