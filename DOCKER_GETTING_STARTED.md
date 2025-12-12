# GameLink Docker 快速上手指南

## 🎯 5 分钟快速开始

### 第一步：检查环境

确保已安装 Docker Desktop：

```powershell
# 检查 Docker 版本
docker --version
# 应该显示: Docker version 20.10.0 或更高

# 检查 Docker Compose 版本
docker-compose --version
# 应该显示: Docker Compose version 2.0.0 或更高
```

如果未安装，请访问 [Docker Desktop 官网](https://www.docker.com/products/docker-desktop) 下载安装。

### 第二步：启动服务

使用统一管理工具（最简单）：

```powershell
# 启动本地生产环境
.\scripts\docker-manager.ps1 start
```

或使用独立脚本：

```powershell
.\scripts\docker-prod-local-start.ps1
```

等待约 30 秒，服务将自动启动并进行健康检查。

### 第三步：访问应用

打开浏览器访问：

- **后端 API**: http://localhost:8081
- **Swagger 文档**: http://localhost:8081/swagger/index.html

默认管理员账号：
- 邮箱: `admin@gamelink.com`
- 密码: `admin123456`

### 第四步：验证服务

```powershell
# 健康检查
.\scripts\docker-manager.ps1 health

# 查看容器状态
.\scripts\docker-manager.ps1 ps

# 查看后端日志
.\scripts\docker-manager.ps1 logs-backend
```

## 🎓 常用操作

### 查看帮助

```powershell
# 查看所有可用命令
.\scripts\docker-manager.ps1 help
```

### 启动和停止

```powershell
# 启动服务
.\scripts\docker-manager.ps1 start

# 停止服务
.\scripts\docker-manager.ps1 stop

# 重启服务
.\scripts\docker-manager.ps1 restart

# 清理数据并重新启动
.\scripts\docker-manager.ps1 start-clean
```

### 监控服务

```powershell
# 健康检查
.\scripts\docker-manager.ps1 health

# 查看容器状态
.\scripts\docker-manager.ps1 ps

# 查看资源使用
.\scripts\docker-manager.ps1 stats

# 快速检查（健康+状态+资源）
.\scripts\docker-manager.ps1 quick-check
```

### 查看日志

```powershell
# 查看所有日志
.\scripts\docker-manager.ps1 logs

# 查看后端日志（持续跟踪）
.\scripts\docker-manager.ps1 logs-backend

# 查看数据库日志
.\scripts\docker-manager.ps1 logs-db

# 查看 Redis 日志
.\scripts\docker-manager.ps1 logs-redis
```

### 数据管理

```powershell
# 备份数据
.\scripts\docker-manager.ps1 backup

# 恢复数据
.\scripts\docker-manager.ps1 restore .\backups\20250113_120000.zip

# 进入数据库
.\scripts\docker-manager.ps1 db-shell

# 进入 Redis
.\scripts\docker-manager.ps1 redis-shell
```

### 清理环境

```powershell
# 软清理（仅停止容器）
.\scripts\docker-manager.ps1 clean

# 中等清理（删除容器和镜像）
.\scripts\docker-manager.ps1 clean-medium

# 完全清理（删除所有数据）
.\scripts\docker-manager.ps1 clean-hard
```

## 🔧 数据库操作

### 连接数据库

使用数据库客户端（如 DBeaver、pgAdmin）：

```
主机: localhost
端口: 5433
数据库: gamelink
用户名: gamelink
密码: gamelink123
```

或使用命令行：

```powershell
# 进入 PostgreSQL Shell
.\scripts\docker-manager.ps1 db-shell

# 在 Shell 中执行 SQL
SELECT * FROM users LIMIT 10;
```

### 常用 SQL 操作

```sql
-- 查看所有表
\dt

-- 查看表结构
\d users

-- 查询用户数量
SELECT COUNT(*) FROM users;

-- 查看最近的订单
SELECT * FROM orders ORDER BY created_at DESC LIMIT 10;

-- 退出
\q
```

## 🔍 Redis 操作

### 连接 Redis

```powershell
# 进入 Redis CLI
.\scripts\docker-manager.ps1 redis-shell
```

### 常用 Redis 命令

```redis
# 查看所有键
KEYS *

# 获取键值
GET user:1

# 查看 Redis 信息
INFO

# 查看内存使用
INFO memory

# 退出
exit
```

## 🐛 故障排查

### 服务无法启动

```powershell
# 1. 查看日志
.\scripts\docker-manager.ps1 logs-backend

# 2. 检查端口占用
netstat -ano | findstr :8081

# 3. 重新启动
.\scripts\docker-manager.ps1 stop
.\scripts\docker-manager.ps1 start
```

### 数据库连接失败

```powershell
# 1. 检查数据库状态
docker exec gamelink-postgres pg_isready -U gamelink

# 2. 查看数据库日志
.\scripts\docker-manager.ps1 logs-db

# 3. 重启数据库
docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml restart postgres
```

### Redis 连接失败

```powershell
# 1. 测试 Redis
docker exec gamelink-redis redis-cli -a redis123 PING

# 2. 查看 Redis 日志
.\scripts\docker-manager.ps1 logs-redis

# 3. 重启 Redis
docker-compose --env-file .env.production.local -f docker-compose.prod.local.yml restart redis
```

### 端口被占用

如果端口被占用，可以修改 `docker-compose.prod.local.yml` 中的端口映射：

```yaml
ports:
  - "8082:8080"  # 将主机端口改为 8082
```

## 📚 进阶操作

### 更新代码

```powershell
# 1. 拉取最新代码
git pull

# 2. 重新构建
.\scripts\docker-manager.ps1 build

# 3. 重启服务
.\scripts\docker-manager.ps1 restart

# 4. 验证更新
.\scripts\docker-manager.ps1 health
```

### 仅更新后端

```powershell
.\scripts\docker-manager.ps1 update-backend
```

### 仅更新前端

```powershell
.\scripts\docker-manager.ps1 update-frontend
```

### 定期备份

建议每天备份一次：

```powershell
# 备份数据
.\scripts\docker-manager.ps1 backup

# 查看备份文件
ls .\backups\

# 清理旧备份（保留最近7天）
Get-ChildItem .\backups\*.zip | 
    Where-Object {$_.LastWriteTime -lt (Get-Date).AddDays(-7)} | 
    Remove-Item
```

## 🎯 最佳实践

### 开发流程

1. **启动环境**
   ```powershell
   .\scripts\docker-manager.ps1 start
   ```

2. **开发和测试**
   - 修改代码
   - 测试功能

3. **查看日志**
   ```powershell
   .\scripts\docker-manager.ps1 logs-backend
   ```

4. **定期备份**
   ```powershell
   .\scripts\docker-manager.ps1 backup
   ```

5. **结束工作**
   ```powershell
   .\scripts\docker-manager.ps1 stop
   ```

### 测试流程

1. **清理环境**
   ```powershell
   .\scripts\docker-manager.ps1 start-clean
   ```

2. **运行测试**
   - 执行测试用例
   - 验证功能

3. **检查状态**
   ```powershell
   .\scripts\docker-manager.ps1 quick-check
   ```

4. **备份测试数据**
   ```powershell
   .\scripts\docker-manager.ps1 backup
   ```

## 📖 相关文档

- **[DOCKER_QUICK_REFERENCE.md](DOCKER_QUICK_REFERENCE.md)** - 命令速查表
- **[DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)** - 完整部署指南
- **[README.docker.md](README.docker.md)** - 工具详细说明
- **[DOCKER_SETUP_COMPLETE.md](DOCKER_SETUP_COMPLETE.md)** - 完成总结

## 💡 提示

### 命令别名

可以在 PowerShell 配置文件中添加别名：

```powershell
# 编辑配置文件
notepad $PROFILE

# 添加别名
Set-Alias -Name dm -Value "C:\path\to\GameLink\scripts\docker-manager.ps1"

# 使用别名
dm start
dm health
dm logs-backend
```

### 快捷键

在 Windows Terminal 中可以配置快捷键：

```json
{
    "command": {
        "action": "sendInput",
        "input": ".\\scripts\\docker-manager.ps1 start\r"
    },
    "keys": "ctrl+shift+d"
}
```

## ❓ 常见问题

### Q: 如何切换到开发环境？

A: 使用 `.\scripts\docker-manager.ps1 start-dev` 或 `.\scripts\docker-dev-start.ps1`

### Q: 如何查看所有可用命令？

A: 使用 `.\scripts\docker-manager.ps1 help`

### Q: 数据会丢失吗？

A: 数据存储在 Docker 数据卷中，除非使用 `clean-hard` 命令，否则不会丢失。

### Q: 如何恢复到初始状态？

A: 使用 `.\scripts\docker-manager.ps1 start-clean`

### Q: 端口冲突怎么办？

A: 本地生产环境使用不同端口（8081, 5433, 6380）避免冲突。

## 🆘 获取帮助

如果遇到问题：

1. 查看日志：`.\scripts\docker-manager.ps1 logs-backend`
2. 健康检查：`.\scripts\docker-manager.ps1 health`
3. 查阅文档：`.\scripts\docker-manager.ps1 docs`
4. 提交 Issue 到项目仓库

---

**祝您使用愉快！** 🎉

如有问题，请参考完整文档或联系开发团队。
