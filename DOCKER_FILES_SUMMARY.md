# Docker 环境部署文件清单

本文档列出了为 GameLink 项目创建的所有 Docker 相关文件及其用途。

## 📋 创建的文件列表

### 核心配置文件

#### 1. `docker-compose.yml` - 开发环境配置
**位置**: 项目根目录
**用途**: 开发环境的 Docker Compose 配置
**特点**:
- 使用 SQLite 数据库（无需外部数据库）
- 使用内存缓存（无需 Redis）
- 支持热更新和调试
- 简化的服务依赖

**启动命令**:
```powershell
docker-compose up -d
```

---

#### 2. `docker-compose.prod.yml` - 生产环境配置
**位置**: 项目根目录
**用途**: 生产环境的 Docker Compose 配置
**特点**:
- 完整的技术栈（PostgreSQL + Redis）
- 健康检查和自动重启
- 数据持久化
- 生产级安全配置

**启动命令**:
```powershell
# 需要先配置 .env 文件
docker-compose -f docker-compose.prod.yml up -d
```

---

### Dockerfile 文件

#### 3. `backend/Dockerfile` - 后端镜像（已更新）
**位置**: backend/
**用途**: 构建 Go 后端服务镜像
**更新内容**:
- 默认使用生产环境配置
- 添加健康检查依赖（wget）
- 创建日志目录
- 多阶段构建优化镜像大小

**构建特点**:
- 使用 Go 1.24 Alpine 构建
- 静态编译（CGO_ENABLED=0）
- 最终镜像基于 Alpine 3.19
- 支持构建缓存

---

#### 4. `frontend/Dockerfile` - 前端镜像（新建）
**位置**: frontend/
**用途**: 构建 React 前端应用镜像
**特点**:
- 多阶段构建（构建 + 运行）
- 使用 Nginx 提供静态文件服务
- 自动代理后端 API
- 优化的资源压缩和缓存

**构建产物**: Nginx + 构建后的静态文件

---

#### 5. `frontend/nginx.conf` - Nginx 配置（新建）
**位置**: frontend/
**用途**: 前端容器的 Nginx 服务配置
**功能**:
- 静态资源服务和缓存
- API 请求代理到后端
- Swagger 文档代理
- React Router 支持（SPA 路由）
- Gzip 压缩
- 健康检查端点

---

### 构建优化文件

#### 6. `backend/.dockerignore` - 后端构建排除（新建）
**位置**: backend/
**用途**: 优化后端 Docker 镜像构建
**排除内容**:
- 测试文件（*_test.go）
- 开发工具配置
- Git 文件
- 临时文件和数据库
- 文档和 CI/CD 配置

---

#### 7. `frontend/.dockerignore` - 前端构建排除（新建）
**位置**: frontend/
**用途**: 优化前端 Docker 镜像构建
**排除内容**:
- node_modules（会重新安装）
- 开发工具配置
- Git 文件
- 测试和覆盖率文件
- 已有的构建产物
- 环境变量文件

---

### 环境变量配置

#### 8. `.env.example` - 环境变量模板（新建）
**位置**: 项目根目录
**用途**: 生产环境环境变量配置示例
**包含配置**:
- 数据库配置（PostgreSQL）
- Redis 配置
- JWT 密钥
- 加密密钥和 IV
- 管理员账号配置
- 应用环境设置

**使用方法**:
```powershell
Copy-Item .env.example .env
# 然后编辑 .env 文件
```

---

#### 9. `.env.development` - 开发环境配置（新建）
**位置**: 项目根目录
**用途**: 开发环境的简化配置
**特点**: 最小化配置，使用 SQLite 和内存缓存

---

### 启动脚本

#### 10. `scripts/docker-dev-start.ps1` - 开发环境启动脚本（新建）
**位置**: scripts/
**用途**: 自动化开发环境启动
**功能**:
- 检查 Docker 环境
- 停止旧容器
- 构建镜像
- 启动服务
- 显示访问信息和常用命令

**使用方法**:
```powershell
.\scripts\docker-dev-start.ps1
```

---

#### 11. `scripts/docker-prod-start.ps1` - 生产环境启动脚本（新建）
**位置**: scripts/
**用途**: 自动化生产环境启动
**功能**:
- 检查 Docker 环境
- 验证 .env 文件和必需变量
- 安全确认提示
- 构建和启动服务
- 显示访问信息和安全提醒

**使用方法**:
```powershell
.\scripts\docker-prod-start.ps1
```

---

### 文档

#### 12. `DOCKER_DEPLOYMENT.md` - Docker 部署完整指南（新建）
**位置**: 项目根目录
**用途**: Docker 部署的详细文档
**内容**:
- 前置要求和环境验证
- 开发环境和生产环境部署步骤
- 密钥生成方法
- 常用命令参考
- 健康检查
- 故障排查
- 性能优化
- 数据备份与恢复
- 安全建议

---

#### 13. `README.md` - 项目首页（已更新）
**位置**: 项目根目录
**更新内容**: 在"快速开始"部分添加 Docker 部署方式
**新增**:
- Docker 部署推荐标识
- 开发环境和生产环境启动命令
- Docker 部署文档链接

---

## 🚀 快速使用指南

### 开发环境（最简单）

```powershell
# 1. 启动服务
.\scripts\docker-dev-start.ps1

# 2. 访问应用
# 前端: http://localhost
# 后端: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
```

### 生产环境

```powershell
# 1. 配置环境变量
Copy-Item .env.example .env
notepad .env

# 2. 生成密钥（可选）
# 使用 PowerShell 生成 JWT 密钥
$jwt = [Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
Write-Host "JWT_SECRET_KEY=$jwt"

# 3. 启动服务
.\scripts\docker-prod-start.ps1

# 4. 访问应用
# 前端: http://localhost
# 后端: http://localhost:8080
```

---

## 📂 文件结构概览

```
GameLink/
├── docker-compose.yml              # 开发环境配置
├── docker-compose.prod.yml         # 生产环境配置
├── .env.example                    # 环境变量模板
├── .env.development                # 开发环境变量
├── DOCKER_DEPLOYMENT.md            # Docker 部署文档
├── DOCKER_FILES_SUMMARY.md         # 本文档
├── README.md                       # 项目首页（已更新）
│
├── backend/
│   ├── Dockerfile                  # 后端镜像配置（已更新）
│   └── .dockerignore               # 后端构建排除
│
├── frontend/
│   ├── Dockerfile                  # 前端镜像配置
│   ├── nginx.conf                  # Nginx 配置
│   └── .dockerignore               # 前端构建排除
│
└── scripts/
    ├── docker-dev-start.ps1        # 开发环境启动脚本
    └── docker-prod-start.ps1       # 生产环境启动脚本
```

---

## 🔍 技术细节

### Docker 网络
- **网络名称**: gamelink-network
- **类型**: bridge
- **容器通信**: 通过容器名称（如 `backend`, `postgres`, `redis`）

### 数据持久化

**开发环境**:
- `backend-data`: SQLite 数据库文件

**生产环境**:
- `postgres-data`: PostgreSQL 数据
- `redis-data`: Redis 持久化数据
- `backend-logs`: 后端日志

### 健康检查

所有服务都配置了健康检查：
- **后端**: HTTP GET /api/v1/health
- **前端**: HTTP GET /health
- **PostgreSQL**: pg_isready
- **Redis**: redis-cli ping

### 端口映射

**开发环境**:
- 80 → 前端（Nginx）
- 8080 → 后端（Go API）

**生产环境**:
- 80 → 前端（Nginx）
- 443 → 前端（Nginx，预留 HTTPS）
- 8080 → 后端（Go API）
- 5432 → PostgreSQL（仅调试时暴露）
- 6379 → Redis（仅调试时暴露）

---

## ⚠️ 重要提醒

### 开发环境
- ✅ 使用 SQLite，无需配置数据库
- ✅ 使用内存缓存，无需 Redis
- ✅ 开启 Swagger 文档
- ✅ 自动创建测试数据

### 生产环境
- ⚠️ **必须**配置 .env 文件
- ⚠️ **必须**修改所有默认密码
- ⚠️ **必须**使用强密钥（JWT、加密）
- ⚠️ **建议**启用 HTTPS
- ⚠️ **建议**定期备份数据库
- ⚠️ 关闭 Swagger 文档（默认已关闭）

---

## 📊 资源要求

### 最小配置
- **CPU**: 2 核
- **内存**: 4GB
- **磁盘**: 10GB

### 推荐配置
- **CPU**: 4 核
- **内存**: 8GB
- **磁盘**: 20GB（含日志和数据库）

---

## 🛠️ 常见问题

### Q: 如何查看日志？
```powershell
# 开发环境
docker-compose logs -f

# 生产环境
docker-compose -f docker-compose.prod.yml logs -f
```

### Q: 如何重启服务？
```powershell
# 开发环境
docker-compose restart

# 生产环境
docker-compose -f docker-compose.prod.yml restart
```

### Q: 如何停止服务？
```powershell
# 开发环境（保留数据）
docker-compose down

# 生产环境（保留数据）
docker-compose -f docker-compose.prod.yml down

# 删除所有数据
docker-compose down -v
```

### Q: 如何更新代码？
```powershell
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动
docker-compose up -d --build
```

---

## 📝 下一步计划

- [ ] 添加 HTTPS 支持（Let's Encrypt）
- [ ] 集成监控（Prometheus + Grafana）
- [ ] 添加日志聚合（ELK Stack）
- [ ] 配置自动备份脚本
- [ ] CI/CD 集成

---

**创建日期**: 2025-12-13
**最后更新**: 2025-12-13
**版本**: 1.0.0
