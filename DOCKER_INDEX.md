# GameLink Docker 文档索引

## 📚 快速导航

### 🚀 新手入门
- **[快速上手指南](DOCKER_GETTING_STARTED.md)** ⭐ 推荐首选
  - 5分钟快速开始
  - 常用操作示例
  - 故障排查指南
  - 最佳实践

### 📖 参考文档
- **[快速参考手册](DOCKER_QUICK_REFERENCE.md)**
  - 命令速查表
  - 访问地址
  - 常用操作
  - 故障排查

### 📘 完整指南
- **[完整部署指南](DOCKER_DEPLOYMENT.md)**
  - 详细部署步骤
  - 环境配置
  - 性能优化
  - 安全建议

### 🔧 工具说明
- **[工具集详解](README.docker.md)**
  - 所有脚本说明
  - 配置文件说明
  - 端口映射表
  - 架构对比

### 📊 项目总结
- **[完成总结](DOCKER_SETUP_COMPLETE.md)**
  - 完成内容清单
  - 功能特性
  - 使用方式
  - 学习路径

- **[部署总结](DOCKER_DEPLOYMENT_SUMMARY.md)**
  - 统计数据
  - 性能指标
  - 核心优势
  - 使用场景

## 🛠️ 工具脚本

### 统一管理工具（推荐）
```powershell
.\scripts\docker-manager.ps1 <command>
```

**常用命令**:
- `help` - 查看所有命令
- `start` - 启动本地生产环境
- `health` - 健康检查
- `logs-backend` - 查看后端日志
- `backup` - 备份数据
- `stop` - 停止服务

### 独立脚本
| 脚本 | 功能 | 文档 |
|------|------|------|
| `docker-manager.ps1` | 统一管理工具 | [README.docker.md](README.docker.md#统一管理工具推荐) |
| `docker-prod-local-start.ps1` | 本地生产环境启动 | [README.docker.md](README.docker.md#docker-prod-local-startps1) |
| `docker-dev-start.ps1` | 开发环境启动 | [README.docker.md](README.docker.md#docker-dev-startps1) |
| `docker-prod-start.ps1` | 生产环境启动 | [README.docker.md](README.docker.md#docker-prod-startps1) |
| `docker-health-check.ps1` | 健康检查 | [README.docker.md](README.docker.md#docker-health-checkps1) |
| `docker-logs.ps1` | 日志查看 | [README.docker.md](README.docker.md#docker-logsps1) |
| `docker-backup.ps1` | 数据备份 | [README.docker.md](README.docker.md#docker-backupps1) |
| `docker-restore.ps1` | 数据恢复 | [README.docker.md](README.docker.md#docker-restoreps1) |
| `docker-clean.ps1` | 清理工具 | [README.docker.md](README.docker.md#docker-cleanps1) |

## 📄 配置文件

### Docker Compose
| 文件 | 环境 | 说明 |
|------|------|------|
| `docker-compose.yml` | 开发 | SQLite + 内存缓存 |
| `docker-compose.prod.local.yml` | 本地生产 | PostgreSQL + Redis |
| `docker-compose.prod.yml` | 生产 | PostgreSQL + Redis |
| `docker-compose.backend-only.yml` | 后端 | 仅后端服务 |

### 环境变量
| 文件 | 用途 |
|------|------|
| `.env.example` | 环境变量模板 |
| `.env.development` | 开发环境配置 |
| `.env.production.local` | 本地生产配置 |
| `.env` | 生产环境配置（需创建） |

### 其他配置
| 文件 | 说明 |
|------|------|
| `backend/Dockerfile` | 后端镜像构建 |
| `frontend/Dockerfile` | 前端镜像构建 |
| `frontend/nginx.conf` | Nginx 配置 |
| `backend/scripts/sql/01_init.sql` | 数据库初始化 |
| `Makefile.docker` | Make 快捷命令 |

## 🎯 使用场景

### 场景 1: 第一次使用
1. 阅读 [快速上手指南](DOCKER_GETTING_STARTED.md)
2. 运行 `.\scripts\docker-manager.ps1 start`
3. 访问 http://localhost:8081

### 场景 2: 日常开发
```powershell
# 启动
.\scripts\docker-manager.ps1 start

# 查看日志
.\scripts\docker-manager.ps1 logs-backend

# 停止
.\scripts\docker-manager.ps1 stop
```

### 场景 3: 测试验证
```powershell
# 清理并启动
.\scripts\docker-manager.ps1 start-clean

# 健康检查
.\scripts\docker-manager.ps1 health

# 备份测试数据
.\scripts\docker-manager.ps1 backup
```

### 场景 4: 故障排查
```powershell
# 快速检查
.\scripts\docker-manager.ps1 quick-check

# 查看日志
.\scripts\docker-manager.ps1 logs-backend

# 重启服务
.\scripts\docker-manager.ps1 restart
```

### 场景 5: 数据管理
```powershell
# 备份
.\scripts\docker-manager.ps1 backup

# 恢复
.\scripts\docker-manager.ps1 restore .\backups\20250113_120000.zip

# 进入数据库
.\scripts\docker-manager.ps1 db-shell
```

## 📊 三种环境对比

| 特性 | 开发环境 | 本地生产 | 生产环境 |
|------|----------|----------|----------|
| 数据库 | SQLite | PostgreSQL 16 | PostgreSQL 16 |
| 缓存 | 内存 | Redis 7 | Redis 7 |
| 后端端口 | 8080 | 8081 | 8080 |
| 数据库端口 | - | 5433 | 5432 |
| Redis 端口 | - | 6380 | 6379 |
| 启动命令 | `start-dev` | `start` | `prod` |
| 适用场景 | 快速开发 | 本地测试 | 实际部署 |

## 🔍 快速查找

### 我想...

#### 启动服务
- 开发环境: `.\scripts\docker-manager.ps1 start-dev`
- 本地生产: `.\scripts\docker-manager.ps1 start`
- 生产环境: `.\scripts\docker-prod-start.ps1`

#### 查看状态
- 健康检查: `.\scripts\docker-manager.ps1 health`
- 容器状态: `.\scripts\docker-manager.ps1 ps`
- 资源使用: `.\scripts\docker-manager.ps1 stats`

#### 查看日志
- 所有日志: `.\scripts\docker-manager.ps1 logs`
- 后端日志: `.\scripts\docker-manager.ps1 logs-backend`
- 数据库日志: `.\scripts\docker-manager.ps1 logs-db`

#### 管理数据
- 备份数据: `.\scripts\docker-manager.ps1 backup`
- 恢复数据: `.\scripts\docker-manager.ps1 restore <file>`
- 进入数据库: `.\scripts\docker-manager.ps1 db-shell`

#### 清理环境
- 停止容器: `.\scripts\docker-manager.ps1 stop`
- 软清理: `.\scripts\docker-manager.ps1 clean`
- 完全清理: `.\scripts\docker-manager.ps1 clean-hard`

#### 更新服务
- 更新后端: `.\scripts\docker-manager.ps1 update-backend`
- 更新前端: `.\scripts\docker-manager.ps1 update-frontend`
- 重新构建: `.\scripts\docker-manager.ps1 rebuild`

## 📖 学习路径

### Level 1: 入门（第1天）
1. ✅ 阅读 [快速上手指南](DOCKER_GETTING_STARTED.md)
2. ✅ 启动本地生产环境
3. ✅ 访问应用并测试
4. ✅ 查看日志和状态

**目标**: 能够启动和停止服务

### Level 2: 熟练（第2-3天）
1. ✅ 阅读 [快速参考手册](DOCKER_QUICK_REFERENCE.md)
2. ✅ 掌握所有基本命令
3. ✅ 学习数据库操作
4. ✅ 练习备份恢复

**目标**: 能够独立管理环境

### Level 3: 精通（第4-7天）
1. ✅ 阅读 [完整部署指南](DOCKER_DEPLOYMENT.md)
2. ✅ 学习故障排查
3. ✅ 了解性能优化
4. ✅ 掌握安全配置

**目标**: 能够处理复杂问题

### Level 4: 专家（持续学习）
1. ✅ 阅读 [工具集详解](README.docker.md)
2. ✅ 自定义配置
3. ✅ 优化工作流
4. ✅ 贡献改进

**目标**: 能够优化和扩展

## 🆘 获取帮助

### 查看帮助
```powershell
# 查看所有命令
.\scripts\docker-manager.ps1 help

# 查看脚本帮助
Get-Help .\scripts\docker-prod-local-start.ps1

# 查看文档列表
.\scripts\docker-manager.ps1 docs
```

### 故障排查
1. 查看 [快速参考手册 - 故障排查](DOCKER_QUICK_REFERENCE.md#-故障排查)
2. 查看 [完整部署指南 - 故障排查](DOCKER_DEPLOYMENT.md#故障排查)
3. 查看 [快速上手指南 - 故障排查](DOCKER_GETTING_STARTED.md#-故障排查)

### 常见问题
- [快速上手指南 - 常见问题](DOCKER_GETTING_STARTED.md#-常见问题)
- [完整部署指南 - 常见问题](DOCKER_DEPLOYMENT.md#常见问题)

### 联系支持
- 查看日志: `.\scripts\docker-manager.ps1 logs-backend`
- 健康检查: `.\scripts\docker-manager.ps1 health`
- 提交 Issue 到项目仓库

## 📈 统计信息

### 文档统计
- 文档数量: 6 个
- 总行数: ~2,500 行
- 覆盖场景: 20+ 个

### 工具统计
- 脚本数量: 7 个
- 命令数量: 77+ 个
- 代码行数: ~2,500 行

### 配置统计
- Docker Compose: 4 个
- 环境变量: 4 个
- 其他配置: 4 个

## 🎯 快速链接

### 文档
- [快速上手](DOCKER_GETTING_STARTED.md)
- [快速参考](DOCKER_QUICK_REFERENCE.md)
- [完整指南](DOCKER_DEPLOYMENT.md)
- [工具说明](README.docker.md)
- [完成总结](DOCKER_SETUP_COMPLETE.md)
- [部署总结](DOCKER_DEPLOYMENT_SUMMARY.md)

### 主要文档
- [项目 README](README.md)
- [产品概述](PRODUCT_OVERVIEW.md)
- [技术架构](TECHNICAL_ARCHITECTURE.md)
- [后端文档](backend/README.md)
- [前端文档](frontend/README.md)

### 外部资源
- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Redis 文档](https://redis.io/documentation)

## 💡 提示

### 命令别名
在 PowerShell 配置文件中添加：
```powershell
Set-Alias -Name dm -Value "C:\path\to\GameLink\scripts\docker-manager.ps1"
```

然后可以使用：
```powershell
dm start
dm health
dm logs-backend
```

### 快捷键
推荐在 Windows Terminal 中配置快捷键快速执行常用命令。

### 自动补全
PowerShell 支持 Tab 自动补全，输入命令时可以按 Tab 键补全。

## 🎉 开始使用

```powershell
# 1. 查看帮助
.\scripts\docker-manager.ps1 help

# 2. 启动服务
.\scripts\docker-manager.ps1 start

# 3. 访问应用
# http://localhost:8081

# 4. 查看状态
.\scripts\docker-manager.ps1 health
```

---

**最后更新**: 2025-12-13  
**版本**: 1.0.0  
**维护者**: GameLink Team

**祝您使用愉快！** 🚀
