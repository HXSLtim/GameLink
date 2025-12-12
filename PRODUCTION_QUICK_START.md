# 🚀 GameLink 生产环境快速启动

## ✅ 当前状态

- ✅ 生产环境配置已生成（`.env` 文件）
- ✅ 后端 Docker 镜像已构建
- ✅ PostgreSQL 和 Redis 已配置
- ⏸️ 等待 Docker Desktop 重启

## 🎯 立即开始（3步）

### 1️⃣ 重启 Docker Desktop

打开 Docker Desktop → 右上角设置 → Restart Docker Desktop

### 2️⃣ 运行部署脚本

```powershell
.\scripts\deploy-production.ps1
```

### 3️⃣ 访问应用

- **API**: http://localhost:8080
- **文档**: http://localhost:8080/swagger/index.html
- **账号**: admin@gamelink.com
- **密码**: 查看 `.env` 文件中的 `SUPER_ADMIN_PASSWORD`

## 📋 常用命令

```powershell
# 查看状态
docker ps --filter "name=gamelink"

# 查看日志
docker logs gamelink-backend -f

# 健康检查
.\scripts\docker-health-check.ps1 -Environment prod

# 重启服务
docker-compose -f docker-compose.prod.yml restart

# 停止服务
docker-compose -f docker-compose.prod.yml down

# 备份数据
.\scripts\docker-backup.ps1 -Environment prod
```

## 🔧 故障排查

### 后端一直重启？
```powershell
# 查看日志
docker logs gamelink-backend --tail=100

# 检查数据库
docker exec gamelink-postgres pg_isready -U gamelink

# 重启后端
docker-compose -f docker-compose.prod.yml restart backend
```

### 无法访问 API？
```powershell
# 检查端口
netstat -ano | findstr :8080

# 检查防火墙
# 确保 8080 端口未被阻止
```

## 📖 详细文档

- **[完整部署指南](PRODUCTION_DEPLOYMENT_GUIDE.md)** - 详细步骤和故障排查
- **[Docker 文档索引](DOCKER_INDEX.md)** - 所有 Docker 相关文档

## 🔑 重要信息

### 生产环境密钥位置
所有密钥都在 `.env` 文件中：
- 数据库密码: `POSTGRES_PASSWORD`
- Redis 密码: `REDIS_PASSWORD`
- JWT 密钥: `JWT_SECRET_KEY`
- 管理员密码: `SUPER_ADMIN_PASSWORD`

### 数据库连接信息
```
主机: localhost
端口: 5432
数据库: gamelink
用户名: gamelink
密码: 见 .env 文件
```

### 服务端口
- 后端: 8080
- PostgreSQL: 5432
- Redis: 6379

## ⚠️ 注意事项

1. **不要提交 `.env` 文件到 Git**
2. **定期备份数据库**（使用 `.\scripts\docker-backup.ps1 -Environment prod`）
3. **监控服务日志**（使用 `docker logs gamelink-backend -f`）
4. **前端 TypeScript 错误已修复** ✅ 可以直接部署

## 🎉 部署成功后

1. 访问 Swagger 文档测试 API
2. 使用管理员账号登录
3. 配置系统设置
4. 开始使用 GameLink！

---

**需要帮助？** 查看 [PRODUCTION_DEPLOYMENT_GUIDE.md](PRODUCTION_DEPLOYMENT_GUIDE.md)
