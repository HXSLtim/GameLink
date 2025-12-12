# GameLink Docker 部署工具集

本文档介绍 GameLink 项目的 Docker 部署工具和脚本。

## 📚 文档索引

- **[DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)** - 完整的 Docker 部署指南
- **[DOCKER_QUICK_REFERENCE.md](DOCKER_QUICK_REFERENCE.md)** - 快速参考手册
- **本文档** - 工具和脚本说明

## 🎯 三种部署环境

### 1. 开发环境 (Development)
- **配置文件**: `docker-compose.yml`
- **数据库**: SQLite（文件存储）
- **缓存**: 内存缓存
- **端口**: 8080
- **用途**: 快速开发和调试

### 2. 本地生产环境 (Production Local)
- **配置文件**: `docker-compose.prod.local.yml`
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **端口**: 8081, 5433, 6380
- **用途**: 本地测试生产环境配置

### 3. 生产环境 (Production)
- **配置文件**: `docker-compose.prod.yml`
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **端口**: 8080, 5432, 6379
- **用途**: 实际生产部署

## 🛠️ 工具脚本

### 统一管理工具（推荐）

#### docker-manager.ps1
统一的 Docker 管理命令入口。

```powershell
# 查看所有命令
.\scripts\docker-manager.ps1 help

# 环境管理
.\scripts\docker-manager.ps1 start           # 启动本地生产环境
.\scripts\docker-manager.ps1 start-dev       # 启动开发环境
.\scripts\docker-manager.ps1 stop            # 停止服务
.\scripts\docker-manager.ps1 restart         # 重启服务

# 监控
.\scripts\docker-manager.ps1 health          # 健康检查
.\scripts\docker-manager.ps1 logs-backend    # 查看后端日志
.\scripts\docker-manager.ps1 ps              # 查看容器状态
.\scripts\docker-manager.ps1 stats           # 查看资源使用

# 数据管理
.\scripts\docker-manager.ps1 backup          # 备份数据
.\scripts\docker-manager.ps1 db-shell        # 进入数据库
.\scripts\docker-manager.ps1 redis-shell     # 进入 Redis

# 快捷操作
.\scripts\docker-manager.ps1 quick-start     # 快速启动
.\scripts\docker-manager.ps1 quick-check     # 快速检查
```

**优势**:
- 统一的命令入口
- 简洁的命令名称
- 自动处理参数
- 彩色输出提示

### 启动脚本

#### docker-dev-start.ps1
开发环境启动脚本。

```powershell
.\scripts\docker-dev-start.ps1
```

**功能**:
- 自动检查 Docker 环境
- 停止现有容器
- 构建镜像
- 启动服务
- 显示访问信息

#### docker-prod-local-start.ps1
本地生产环境启动脚本。

```powershell
# 正常启动
.\scripts\docker-prod-local-start.ps1

# 清理数据重新开始
.\scripts\docker-prod-local-start.ps1 -Clean

# 跳过镜像构建
.\scripts\docker-prod-local-start.ps1 -NoBuild
```

**功能**:
- 检查 Docker 环境
- 检查端口冲突
- 可选清理旧数据
- 构建镜像（可跳过）
- 启动服务
- 等待服务就绪
- 健康检查

**特点**:
- 使用不同端口避免冲突（8081, 5433, 6380）
- 完整的生产技术栈（PostgreSQL + Redis）
- 适合本地测试生产环境

#### docker-prod-start.ps1
生产环境启动脚本。

```powershell
.\scripts\docker-prod-start.ps1
```

**功能**:
- 检查环境变量配置
- 安全确认
- 启动生产环境

### 监控脚本

#### docker-health-check.ps1
健康检查工具。

```powershell
# 检查本地生产环境
.\scripts\docker-health-check.ps1 -Environment prod-local

# 检查开发环境
.\scripts\docker-health-check.ps1 -Environment dev

# 检查生产环境
.\scripts\docker-health-check.ps1 -Environment prod
```

**检查项**:
- 容器运行状态
- HTTP 服务健康端点
- 数据库连接和大小
- Redis 连接和统计
- 网络配置
- 数据卷状态
- 资源使用情况

#### docker-logs.ps1
日志查看工具。

```powershell
# 查看所有日志
.\scripts\docker-logs.ps1 -Environment prod-local

# 查看后端日志（持续跟踪）
.\scripts\docker-logs.ps1 -Service backend -Follow

# 查看最近100行
.\scripts\docker-logs.ps1 -Service backend -Lines 100

# 查看特定服务
.\scripts\docker-logs.ps1 -Service postgres
```

**参数**:
- `-Environment`: dev, prod-local, prod
- `-Service`: backend, frontend, postgres, redis, all
- `-Follow`: 持续跟踪日志
- `-Lines`: 显示行数（默认50）

### 数据管理脚本

#### docker-backup.ps1
数据备份工具。

```powershell
# 备份本地生产环境
.\scripts\docker-backup.ps1 -Environment prod-local

# 备份生产环境
.\scripts\docker-backup.ps1 -Environment prod

# 指定备份目录
.\scripts\docker-backup.ps1 -Environment prod-local -BackupDir "D:\backups"
```

**备份内容**:
- PostgreSQL 完整数据库导出（.sql）
- Redis RDB 快照（.rdb）
- 环境配置文件
- 备份信息（JSON）
- 自动压缩为 .zip 文件

**输出**:
```
backups/
└── 20250113_120000.zip
    ├── postgres_backup.sql
    ├── redis_dump.rdb
    ├── env_backup.txt
    ├── config.yaml
    └── backup_info.json
```

#### docker-restore.ps1
数据恢复工具。

```powershell
# 从备份恢复
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod-local

# 恢复到生产环境
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod
```

**功能**:
- 检查备份文件
- 安全确认（需要输入 yes）
- 解压备份
- 恢复 PostgreSQL（删除并重建数据库）
- 恢复 Redis（重启容器）
- 自动清理临时文件

**注意**: 此操作会覆盖现有数据，请谨慎使用！

### 清理脚本

#### docker-clean.ps1
清理工具。

```powershell
# 软清理（仅停止容器）
.\scripts\docker-clean.ps1 -Level soft -Environment prod-local

# 中等清理（删除容器和镜像）
.\scripts\docker-clean.ps1 -Level medium -Environment prod-local

# 完全清理（删除所有包括数据）
.\scripts\docker-clean.ps1 -Level hard -Environment prod-local

# 清理所有环境
.\scripts\docker-clean.ps1 -Level medium -Environment all
```

**清理级别**:
- `soft`: 仅停止容器，保留所有数据
- `medium`: 删除容器和镜像，保留数据卷
- `hard`: 删除所有容器、镜像和数据卷（需要输入 DELETE 确认）

**参数**:
- `-Level`: soft, medium, hard
- `-Environment`: dev, prod-local, prod, all

## 📋 Makefile 快捷命令

使用 `Makefile.docker` 简化操作：

```powershell
# 查看所有命令
make -f Makefile.docker help

# 启动环境
make -f Makefile.docker dev              # 开发环境
make -f Makefile.docker prod-local       # 本地生产环境
make -f Makefile.docker prod-local-clean # 清理并启动

# 监控
make -f Makefile.docker health           # 健康检查
make -f Makefile.docker logs-backend     # 查看后端日志
make -f Makefile.docker ps               # 查看容器状态
make -f Makefile.docker stats            # 查看资源使用

# 数据管理
make -f Makefile.docker backup           # 备份数据
make -f Makefile.docker db-shell         # 进入数据库
make -f Makefile.docker redis-shell      # 进入 Redis

# 清理
make -f Makefile.docker clean-soft       # 软清理
make -f Makefile.docker clean-medium     # 中等清理
make -f Makefile.docker prune            # 清理未使用资源

# 更新
make -f Makefile.docker update-backend   # 更新后端
make -f Makefile.docker update-frontend  # 更新前端

# 快捷操作
make -f Makefile.docker quick-start      # 快速启动
make -f Makefile.docker quick-check      # 快速检查
```

## 🔧 配置文件

### Docker Compose 文件

| 文件 | 用途 | 环境 |
|------|------|------|
| `docker-compose.yml` | 开发环境 | SQLite + 内存缓存 |
| `docker-compose.prod.local.yml` | 本地生产测试 | PostgreSQL + Redis |
| `docker-compose.prod.yml` | 生产环境 | PostgreSQL + Redis |
| `docker-compose.backend-only.yml` | 仅后端服务 | 用于前端独立开发 |

### 环境变量文件

| 文件 | 用途 |
|------|------|
| `.env.example` | 环境变量模板 |
| `.env.development` | 开发环境配置 |
| `.env.production.local` | 本地生产测试配置 |
| `.env` | 生产环境配置（需创建） |

### Dockerfile

| 文件 | 说明 |
|------|------|
| `backend/Dockerfile` | 后端多阶段构建 |
| `frontend/Dockerfile` | 前端多阶段构建 + Nginx |
| `frontend/nginx.conf` | Nginx 配置（反向代理） |

### 数据库初始化

| 文件 | 说明 |
|------|------|
| `backend/scripts/sql/01_init.sql` | PostgreSQL 初始化脚本 |

## 📊 端口映射

### 开发环境
| 服务 | 容器端口 | 主机端口 |
|------|----------|----------|
| 后端 | 8080 | 8080 |
| 前端 | 80 | 80 |

### 本地生产环境
| 服务 | 容器端口 | 主机端口 |
|------|----------|----------|
| 后端 | 8080 | 8081 |
| PostgreSQL | 5432 | 5433 |
| Redis | 6379 | 6380 |

### 生产环境
| 服务 | 容器端口 | 主机端口 |
|------|----------|----------|
| 后端 | 8080 | 8080 |
| 前端 | 80 | 80 |
| 前端 (HTTPS) | 443 | 443 |
| PostgreSQL | 5432 | 5432 |
| Redis | 6379 | 6379 |

## 🚀 快速开始

### 第一次使用

1. **安装 Docker Desktop**
   ```powershell
   # 验证安装
   docker --version
   docker-compose --version
   ```

2. **启动开发环境**
   ```powershell
   .\scripts\docker-dev-start.ps1
   ```

3. **访问应用**
   - 前端: http://localhost
   - 后端: http://localhost:8080
   - Swagger: http://localhost:8080/swagger/index.html

### 测试生产环境

1. **启动本地生产环境**
   ```powershell
   .\scripts\docker-prod-local-start.ps1
   ```

2. **健康检查**
   ```powershell
   .\scripts\docker-health-check.ps1 -Environment prod-local
   ```

3. **查看日志**
   ```powershell
   .\scripts\docker-logs.ps1 -Service backend -Follow
   ```

4. **备份数据**
   ```powershell
   .\scripts\docker-backup.ps1 -Environment prod-local
   ```

## 🔍 故障排查

### 常见问题

1. **端口被占用**
   ```powershell
   # 查看端口占用
   netstat -ano | findstr :8081
   
   # 修改 docker-compose.prod.local.yml 中的端口
   ```

2. **容器无法启动**
   ```powershell
   # 查看日志
   .\scripts\docker-logs.ps1 -Service backend
   
   # 检查配置
   docker exec gamelink-backend cat /app/configs/config.production.yaml
   ```

3. **数据库连接失败**
   ```powershell
   # 检查数据库状态
   docker exec gamelink-postgres pg_isready -U gamelink
   
   # 查看数据库日志
   .\scripts\docker-logs.ps1 -Service postgres
   ```

4. **Redis 连接失败**
   ```powershell
   # 测试 Redis
   docker exec gamelink-redis redis-cli -a redis123 PING
   ```

### 调试技巧

```powershell
# 进入容器
docker exec -it gamelink-backend sh

# 查看环境变量
docker exec gamelink-backend env

# 查看文件
docker exec gamelink-backend ls -la /app

# 查看网络
docker network inspect gamelink-network

# 查看数据卷
docker volume inspect gamelink-postgres-data
```

## 📈 最佳实践

### 1. 定期备份
```powershell
# 每天备份
.\scripts\docker-backup.ps1 -Environment prod

# 保留最近7天的备份
Get-ChildItem .\backups\*.zip | 
    Where-Object {$_.LastWriteTime -lt (Get-Date).AddDays(-7)} | 
    Remove-Item
```

### 2. 监控健康
```powershell
# 定期健康检查
.\scripts\docker-health-check.ps1 -Environment prod-local

# 监控资源使用
docker stats --no-stream
```

### 3. 日志管理
```powershell
# 查看错误日志
.\scripts\docker-logs.ps1 -Service backend | Select-String "ERROR"

# 导出日志
.\scripts\docker-logs.ps1 -Service backend -Lines 1000 > backend.log
```

### 4. 更新部署
```powershell
# 拉取最新代码
git pull

# 更新后端
make -f Makefile.docker update-backend

# 验证更新
make -f Makefile.docker health
```

### 5. 清理维护
```powershell
# 定期清理未使用资源
make -f Makefile.docker prune

# 查看磁盘使用
docker system df
```

## 🔐 安全建议

1. **修改默认密码**: 生产环境必须修改所有默认密码
2. **使用强密钥**: JWT 和加密密钥使用随机生成
3. **限制访问**: 使用防火墙限制端口访问
4. **启用 HTTPS**: 生产环境配置 SSL 证书
5. **定期备份**: 建立自动备份策略
6. **监控日志**: 定期检查安全日志
7. **更新镜像**: 定期更新基础镜像

## 📚 相关资源

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [PostgreSQL Docker 镜像](https://hub.docker.com/_/postgres)
- [Redis Docker 镜像](https://hub.docker.com/_/redis)
- [Nginx Docker 镜像](https://hub.docker.com/_/nginx)

## 🤝 贡献

如果您发现问题或有改进建议，请：
1. 查看现有文档
2. 运行健康检查
3. 提交 Issue 或 Pull Request

---

**最后更新**: 2025-12-13  
**维护者**: GameLink Team
