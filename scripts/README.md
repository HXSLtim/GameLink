# GameLink 部署脚本说明

本目录包含 GameLink 项目的各种部署和管理脚本。

## 📋 生产环境部署脚本

### 1. 标准部署（不含加密）

**脚本**: `deploy-production.ps1`

**用途**: 快速部署生产环境，不启用加密中间件。适合内网环境或已有其他安全措施的场景。

**使用方法**:
```powershell
# 完整部署
.\scripts\deploy-production.ps1

# 跳过构建（使用现有镜像）
.\scripts\deploy-production.ps1 -SkipBuild

# 跳过前端构建
.\scripts\deploy-production.ps1 -SkipFrontend

# 不拉取基础镜像
.\scripts\deploy-production.ps1 -NoPull
```

**特点**:
- ✅ 快速部署，步骤简单
- ✅ 不需要配置加密密钥
- ✅ 适合开发和测试环境
- ⚠️ 请求体数据未加密

---

### 2. 加密部署（推荐生产环境）

**脚本**: `deploy-production-encrypted.ps1`

**用途**: 部署启用 AES-256-CBC 加密的生产环境。保护前后端通信数据安全。

**使用方法**:
```powershell
# 完整部署（自动生成密钥）
.\scripts\deploy-production-encrypted.ps1

# 重新生成加密密钥
.\scripts\deploy-production-encrypted.ps1 -RegenerateKeys

# 跳过构建
.\scripts\deploy-production-encrypted.ps1 -SkipBuild

# 跳过前端
.\scripts\deploy-production-encrypted.ps1 -SkipFrontend
```

**特点**:
- ✅ AES-256-CBC 加密保护请求体
- ✅ SHA-256 签名验证请求完整性
- ✅ 自动生成和同步加密密钥
- ✅ 防止中间人攻击和数据篡改
- ⚠️ 需要安装 crypto-js 依赖

**加密说明**:
- 加密算法: AES-256-CBC
- 签名算法: SHA-256
- 加密范围: POST/PUT/PATCH 请求体
- 排除路径: /health, /ping, /auth/refresh

---

## 🛠️ 其他工具脚本

### Docker 管理工具

**脚本**: `docker-manager.ps1`

**用途**: 统一管理 Docker 服务的工具，提供 30+ 个常用命令。

**使用方法**:
```powershell
# 查看帮助
.\scripts\docker-manager.ps1 help

# 启动服务
.\scripts\docker-manager.ps1 start

# 停止服务
.\scripts\docker-manager.ps1 stop

# 查看日志
.\scripts\docker-manager.ps1 logs backend

# 健康检查
.\scripts\docker-manager.ps1 health
```

---

### 本地生产环境部署

**脚本**: `docker-prod-local-start.ps1`

**用途**: 在本地启动生产环境（使用不同端口避免冲突）。

**端口映射**:
- 前端: 8081 (HTTP), 8443 (HTTPS)
- 后端: 8082
- PostgreSQL: 5433
- Redis: 6380

**使用方法**:
```powershell
.\scripts\docker-prod-local-start.ps1
```

---

### 密钥同步工具

**脚本**: `sync-crypto-keys.ps1`

**用途**: 将后端 `.env` 中的加密密钥同步到前端 `.env.production`。

**使用方法**:
```powershell
.\scripts\sync-crypto-keys.ps1
```

---

### 健康检查

**脚本**: `docker-health-check.ps1`

**用途**: 检查所有服务的健康状态。

**使用方法**:
```powershell
.\scripts\docker-health-check.ps1
```

---

### 日志查看

**脚本**: `docker-logs.ps1`

**用途**: 查看服务日志，支持实时跟踪和过滤。

**使用方法**:
```powershell
# 查看后端日志
.\scripts\docker-logs.ps1 -Service backend

# 实时跟踪
.\scripts\docker-logs.ps1 -Service backend -Follow

# 查看最近 100 行
.\scripts\docker-logs.ps1 -Service backend -Tail 100
```

---

### 数据备份

**脚本**: `docker-backup.ps1`

**用途**: 备份数据库和 Redis 数据。

**使用方法**:
```powershell
# 完整备份
.\scripts\docker-backup.ps1

# 只备份数据库
.\scripts\docker-backup.ps1 -DatabaseOnly

# 只备份 Redis
.\scripts\docker-backup.ps1 -RedisOnly
```

---

### 数据恢复

**脚本**: `docker-restore.ps1`

**用途**: 从备份恢复数据。

**使用方法**:
```powershell
# 恢复最新备份
.\scripts\docker-restore.ps1

# 恢复指定备份
.\scripts\docker-restore.ps1 -BackupFile "backup-20231213-120000.tar.gz"
```

---

### 清理工具

**脚本**: `docker-clean.ps1`

**用途**: 清理 Docker 资源（容器、镜像、卷、网络）。

**使用方法**:
```powershell
# 清理所有资源
.\scripts\docker-clean.ps1 -All

# 只清理容器
.\scripts\docker-clean.ps1 -Containers

# 只清理镜像
.\scripts\docker-clean.ps1 -Images

# 清理未使用的卷
.\scripts\docker-clean.ps1 -Volumes
```

---

## 📚 部署流程对比

### 标准部署流程

```
1. 检查环境变量 (.env)
2. 安装前端依赖
3. 构建前端
4. 构建 Docker 镜像
5. 停止旧服务
6. 启动新服务
```

### 加密部署流程

```
1. 检查环境变量 (.env)
2. 检查/生成加密密钥
3. 同步密钥到前端
4. 安装前端依赖（包括 crypto-js）
5. 构建前端
6. 构建 Docker 镜像
7. 停止旧服务
8. 启动新服务
```

---

## 🔐 加密配置说明

### 环境变量要求

**后端** (`.env`):
```bash
CRYPTO_SECRET_KEY=<32字符密钥>  # AES-256 密钥
CRYPTO_IV=<16字符IV>            # 初始化向量
```

**前端** (`frontend/.env.production`):
```bash
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=<与后端相同>
VITE_CRYPTO_IV=<与后端相同>
VITE_CRYPTO_USE_SIGNATURE=true
```

### 密钥生成

加密部署脚本会自动生成符合要求的密钥：
- `CRYPTO_SECRET_KEY`: 32 字符（A-Z, a-z, 0-9）
- `CRYPTO_IV`: 16 字符（A-Z, a-z, 0-9）

也可以手动生成：
```powershell
# 生成 32 字符密钥
$chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
-join ((1..32) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })

# 生成 16 字符 IV
-join ((1..16) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })
```

---

## 🚨 常见问题

### Q1: 部署失败，提示找不到 .env 文件

**解决方案**:
```powershell
# 复制示例文件
Copy-Item .env.example .env

# 编辑配置
notepad .env
```

### Q2: 前端构建失败，提示 crypto-js 未找到

**解决方案**:
```powershell
cd frontend
npm install crypto-js @types/crypto-js
cd ..
```

### Q3: 加密中间件未启用

**解决方案**:
```powershell
# 检查后端日志
docker logs gamelink-backend | Select-String "crypto"

# 应该看到: "crypto middleware enabled"
# 如果没有，检查 docker-compose.prod.yml 中的环境变量
```

### Q4: 请求解密失败

**可能原因**:
1. 前后端密钥不一致
2. 前端环境变量未生效（需要重新构建）
3. 密钥长度不正确

**解决方案**:
```powershell
# 重新同步密钥并部署
.\scripts\deploy-production-encrypted.ps1 -RegenerateKeys
```

---

## 📖 相关文档

- **[加密配置指南](../CRYPTO_SETUP_GUIDE.md)** - 详细的加密配置说明
- **[部署总结](../DEPLOYMENT_SUMMARY.md)** - 完整的部署步骤
- **[最终部署状态](../FINAL_DEPLOYMENT_STATUS.md)** - 当前部署状态
- **[Docker 快速参考](../DOCKER_QUICK_REFERENCE.md)** - 常用 Docker 命令
- **[生产部署指南](../PRODUCTION_DEPLOYMENT_GUIDE.md)** - 生产环境部署详解

---

## 💡 最佳实践

1. **生产环境推荐使用加密部署**
   ```powershell
   .\scripts\deploy-production-encrypted.ps1
   ```

2. **定期备份数据**
   ```powershell
   .\scripts\docker-backup.ps1
   ```

3. **定期更换加密密钥**（建议 3-6 个月）
   ```powershell
   .\scripts\deploy-production-encrypted.ps1 -RegenerateKeys
   ```

4. **监控服务健康状态**
   ```powershell
   .\scripts\docker-health-check.ps1
   ```

5. **查看日志排查问题**
   ```powershell
   .\scripts\docker-logs.ps1 -Service backend -Follow
   ```

---

**更新时间**: 2025-12-13  
**版本**: 1.0.0
