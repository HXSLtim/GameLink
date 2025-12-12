# GameLink Docker 快速参考

## 🚀 快速启动

### 使用管理工具（推荐）
```powershell
# 查看所有命令
.\scripts\docker-manager.ps1 help

# 启动本地生产环境
.\scripts\docker-manager.ps1 start

# 启动开发环境
.\scripts\docker-manager.ps1 start-dev

# 健康检查
.\scripts\docker-manager.ps1 health

# 查看日志
.\scripts\docker-manager.ps1 logs-backend
```

### 使用独立脚本
```powershell
# 开发环境（SQLite）
.\scripts\docker-dev-start.ps1

# 本地生产测试（PostgreSQL + Redis）
.\scripts\docker-prod-local-start.ps1

# 生产环境
.\scripts\docker-prod-start.ps1
```

### 访问地址
- 开发环境: http://localhost:8080
- 本地生产: http://localhost:8081
- 生产环境: http://localhost:8080

## 📊 服务管理

### 查看状态
```powershell
# 健康检查
.\scripts\docker-health-check.ps1 -Environment prod-local

# 查看容器
docker-compose -f docker-compose.prod.local.yml ps

# 查看资源使用
docker stats
```

### 查看日志
```powershell
# 查看所有日志
.\scripts\docker-logs.ps1 -Environment prod-local

# 查看后端日志（持续跟踪）
.\scripts\docker-logs.ps1 -Service backend -Follow

# 查看最近100行
.\scripts\docker-logs.ps1 -Service backend -Lines 100
```

### 重启服务
```powershell
# 重启所有服务
docker-compose -f docker-compose.prod.local.yml restart

# 重启后端
docker-compose -f docker-compose.prod.local.yml restart backend
```

### 停止服务
```powershell
# 停止所有服务
docker-compose -f docker-compose.prod.local.yml down

# 停止并删除数据
docker-compose -f docker-compose.prod.local.yml down -v
```

## 💾 数据管理

### 备份
```powershell
# 备份所有数据
.\scripts\docker-backup.ps1 -Environment prod-local

# 指定备份目录
.\scripts\docker-backup.ps1 -Environment prod-local -BackupDir "D:\backups"
```

### 恢复
```powershell
# 从备份恢复
.\scripts\docker-restore.ps1 -BackupFile ".\backups\20250113_120000.zip" -Environment prod-local
```

### 手动数据库操作
```powershell
# 进入 PostgreSQL
docker exec -it gamelink-postgres psql -U gamelink -d gamelink

# 导出数据库
docker exec gamelink-postgres pg_dump -U gamelink gamelink > backup.sql

# 导入数据库
docker exec -i gamelink-postgres psql -U gamelink gamelink < backup.sql
```

### Redis 操作
```powershell
# 进入 Redis CLI
docker exec -it gamelink-redis redis-cli -a redis123

# 查看所有键
docker exec gamelink-redis redis-cli -a redis123 KEYS "*"

# 清空 Redis
docker exec gamelink-redis redis-cli -a redis123 FLUSHALL
```

## 🧹 清理

### 软清理（仅停止）
```powershell
.\scripts\docker-clean.ps1 -Level soft -Environment prod-local
```

### 中等清理（删除容器和镜像）
```powershell
.\scripts\docker-clean.ps1 -Level medium -Environment prod-local
```

### 完全清理（删除所有数据）
```powershell
.\scripts\docker-clean.ps1 -Level hard -Environment prod-local
```

## 🔧 故障排查

### 服务无法启动
```powershell
# 查看详细日志
docker-compose -f docker-compose.prod.local.yml logs backend

# 检查配置
docker exec gamelink-backend cat /app/configs/config.production.yaml

# 检查环境变量
docker exec gamelink-backend env
```

### 数据库连接失败
```powershell
# 检查 PostgreSQL 状态
docker exec gamelink-postgres pg_isready -U gamelink

# 查看数据库日志
docker-compose -f docker-compose.prod.local.yml logs postgres

# 测试连接
docker exec gamelink-postgres psql -U gamelink -d gamelink -c "SELECT 1"
```

### Redis 连接失败
```powershell
# 检查 Redis 状态
docker exec gamelink-redis redis-cli -a redis123 PING

# 查看 Redis 日志
docker-compose -f docker-compose.prod.local.yml logs redis
```

### 端口被占用
```powershell
# 查看端口占用
netstat -ano | findstr :8081
netstat -ano | findstr :5433
netstat -ano | findstr :6380

# 修改端口（编辑 docker-compose.prod.local.yml）
# 例如: "8082:8080" 将主机端口改为 8082
```

### 容器内调试
```powershell
# 进入后端容器
docker exec -it gamelink-backend sh

# 进入前端容器
docker exec -it gamelink-frontend sh

# 查看容器文件
docker exec gamelink-backend ls -la /app
```

## 📍 访问地址

### 开发环境
- 前端: http://localhost
- 后端: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html

### 本地生产环境
- 后端: http://localhost:8081
- Swagger: http://localhost:8081/swagger/index.html
- PostgreSQL: localhost:5433
- Redis: localhost:6380

### 默认账号
- 邮箱: admin@gamelink.com
- 密码: admin123456 (本地测试)

## 🔑 数据库连接信息

### 本地生产环境
```
主机: localhost
端口: 5433
数据库: gamelink
用户名: gamelink
密码: gamelink123
```

### Redis 连接
```
主机: localhost
端口: 6380
密码: redis123
```

## 📝 常用 Docker 命令

### 容器操作
```powershell
# 查看所有容器
docker ps -a

# 查看 GameLink 容器
docker ps --filter "name=gamelink"

# 删除停止的容器
docker container prune

# 强制删除容器
docker rm -f gamelink-backend
```

### 镜像操作
```powershell
# 查看镜像
docker images

# 删除镜像
docker rmi gamelink-backend

# 清理未使用的镜像
docker image prune -a
```

### 数据卷操作
```powershell
# 查看数据卷
docker volume ls

# 查看 GameLink 数据卷
docker volume ls --filter "name=gamelink"

# 查看数据卷详情
docker volume inspect gamelink-postgres-data

# 删除数据卷
docker volume rm gamelink-postgres-data

# 清理未使用的数据卷
docker volume prune
```

### 网络操作
```powershell
# 查看网络
docker network ls

# 查看网络详情
docker network inspect gamelink-network

# 删除网络
docker network rm gamelink-network
```

## 🔄 更新部署

### 更新代码
```powershell
# 拉取最新代码
git pull

# 重新构建并启动
.\scripts\docker-prod-local-start.ps1

# 或手动执行
docker-compose -f docker-compose.prod.local.yml build
docker-compose -f docker-compose.prod.local.yml up -d
```

### 仅更新后端
```powershell
docker-compose -f docker-compose.prod.local.yml build backend
docker-compose -f docker-compose.prod.local.yml up -d --no-deps backend
```

### 仅更新前端
```powershell
docker-compose -f docker-compose.prod.local.yml build frontend
docker-compose -f docker-compose.prod.local.yml up -d --no-deps frontend
```

## 🎯 最佳实践

1. **定期备份**: 每天自动备份数据
   ```powershell
   # 创建定时任务
   .\scripts\docker-backup.ps1 -Environment prod
   ```

2. **监控日志**: 定期检查错误日志
   ```powershell
   .\scripts\docker-logs.ps1 -Service backend | Select-String "ERROR"
   ```

3. **健康检查**: 每小时检查服务状态
   ```powershell
   .\scripts\docker-health-check.ps1 -Environment prod-local
   ```

4. **资源监控**: 监控容器资源使用
   ```powershell
   docker stats --no-stream
   ```

5. **清理空间**: 定期清理未使用的资源
   ```powershell
   docker system prune -a --volumes
   ```

## 📚 相关文档

- [完整部署指南](DOCKER_DEPLOYMENT.md)
- [项目结构](docs/guides/project-structure.md)
- [技术栈说明](.kiro/steering/tech.md)
- [产品概述](PRODUCT_OVERVIEW.md)

---

**提示**: 所有脚本都支持 `-?` 参数查看帮助信息
```powershell
.\scripts\docker-prod-local-start.ps1 -?
```
