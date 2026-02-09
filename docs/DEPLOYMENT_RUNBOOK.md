# GameLink 部署运行手册

**文档版本：** 1.0.0
**创建日期：** 2026-02-09
**维护人：** DevOps-Engineer
**适用环境：** Staging / Production

---

## 目录

1. [部署前检查](#部署前检查)
2. [Staging 环境部署](#staging-环境部署)
3. [生产环境部署](#生产环境部署)
4. [部署验证](#部署验证)
5. [回滚流程](#回滚流程)
6. [监控和告警](#监控和告警)
7. [应急预案](#应急预案)

---

## 部署前检查

### 检查清单

#### 1. 环境准备

- [ ] 服务器资源充足（CPU、内存、磁盘）
- [ ] Docker 和 Docker Compose 已安装
- [ ] 网络连接正常
- [ ] SSL/TLS 证书已准备（生产环境）

#### 2. 配置文件

- [ ] `.env.production` 已创建并配置
- [ ] `admin/.env.production` 已创建并配置
- [ ] `docker-compose.prod.yml` 已准备
- [ ] `deploy/nginx-production.conf` 已准备

#### 3. 数据库

- [ ] 数据库备份已完成
- [ ] 迁移脚本已准备
- [ ] 回滚计划已确认

#### 4. 代码

- [ ] 代码已合并到主分支
- [ ] Docker 镜像已构建
- [ ] 版本标签已创建

#### 5. 依赖服务

- [ ] PostgreSQL 已准备
- [ ] Redis 已准备
- [ ] 对象存储已配置（如需要）

### 运行检查脚本

```bash
# 运行部署前检查
bash scripts/pre-deployment-check.sh

# 运行端口检查
bash scripts/check-ports.sh

# 运行加密密钥验证
bash scripts/verify-crypto-keys.sh
```

---

## Staging 环境部署

### 目的

在类生产环境中测试新功能和配置，确保生产部署顺利。

### 部署步骤

#### 1. 生成 Staging 配置

```bash
# 生成 Staging 环境配置
bash scripts/generate-production-keys.sh

# 重命名为 Staging 配置
mv .env.production .env.staging
mv admin/.env.production admin/.env.staging
```

#### 2. 部署服务

```bash
# 使用自动化部署脚本
bash scripts/deploy.sh staging

# 或手动部署
docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d
```

#### 3. 执行数据库迁移

```bash
# 连接到后端容器
docker exec -it gamelink-backend-staging bash

# 运行迁移
go run cmd/main.go migrate

# 退出容器
exit
```

#### 4. 验证部署

```bash
# 运行部署验证脚本
bash scripts/verify-deployment.sh
```

### Staging 环境测试

#### 功能测试

- [ ] 用户认证流程
- [ ] 陪玩师浏览和详情
- [ ] 订单创建和支付（测试模式）
- [ ] 聊天功能
- [ ] 钱包和充值
- [ ] 管理后台所有功能

#### 性能测试

- [ ] API 响应时间 < 200ms (P95)
- [ ] 页面加载时间 < 2s
- [ ] 数据库查询优化

#### 安全测试

- [ ] 加密通讯正常
- [ ] SQL 注入防护
- [ ] XSS 防护
- [ ] CSRF 防护

---

## 生产环境部署

### 部署前准备

#### 1. 通知团队

```markdown
@team 生产环境部署将在 30 分钟后开始

预计影响时间：10-15 分钟
影响范围：短暂的服务中断

请知悉相关方并做好准备。
```

#### 2. 备份数据

```bash
# 备份数据库
bash scripts/backup-database.sh production

# 验证备份文件
ls -lh ./backups/postgres/production/
```

#### 3. 创建部署标签

```bash
git tag -a v1.0.0 -m "Production release v1.0.0"
git push origin v1.0.0
```

### 部署步骤

#### 1. 使用自动化部署（推荐）

```bash
# 一键部署
bash scripts/deploy.sh production
```

#### 2. 手动部署

**步骤 1：拉取最新代码**

```bash
git fetch origin
git checkout main
git pull origin main
```

**步骤 2：构建 Docker 镜像**

```bash
# 构建后端镜像
cd api
docker build -t gamelink-backend:v1.0.0 .
cd ..

# 构建前端镜像
cd admin
docker build -t gamelink-admin:v1.0.0 .
cd ..
```

**步骤 3：部署服务**

```bash
# 停止现有服务
docker-compose -f docker-compose.prod.yml --env-file .env.production down

# 启动新服务
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
```

**步骤 4：等待服务就绪**

```bash
# 检查服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看后端日志
docker logs -f gamelink-backend
```

**步骤 5：执行数据库迁移**

```bash
# 如果有新的迁移
docker exec -it gamelink-backend go run cmd/main.go migrate
```

**步骤 6：验证部署**

```bash
# 运行验证脚本
bash scripts/verify-deployment.sh production
```

---

## 部署验证

### 健康检查

#### 后端 API

```bash
# 健康检查端点
curl http://localhost:8080/api/v1/healthz

# 预期响应
{
  "status": "healthy",
  "timestamp": "2026-02-09T10:00:00Z"
}
```

#### 前端

```bash
# 检查 Admin
curl -I http://localhost/

# 检查响应头
# HTTP/1.1 200 OK
# Content-Type: text/html
```

#### 数据库

```bash
# 检查 PostgreSQL
docker exec gamelink-postgres pg_isready -U gamelink

# 检查数据库连接
docker exec gamelink-postgres psql \
  -U gamelink -d gamelink -c "SELECT 1;"
```

#### Redis

```bash
# 检查 Redis
docker exec gamelink-redis redis-cli ping

# 预期响应
# PONG
```

### 功能验证清单

- [ ] 用户可以登录
- [ ] 陪玩师列表可以浏览
- [ ] 订单可以创建
- [ ] 聊天功能正常
- [ ] 管理后台可以访问
- [ ] WebSocket 连接正常
- [ ] 支付回调正常（测试环境）

### 性能验证

```bash
# API 性能测试
ab -n 1000 -c 10 http://localhost:8080/api/v1/healthz

# 检查响应时间
# 应该 < 100ms (P95)
```

---

## 回滚流程

### 回滚触发条件

- 部署后发现严重 bug
- 性能严重下降
- 数据库迁移失败
- 安全漏洞

### 回滚步骤

#### 方案 A：使用备份回滚（推荐）

```bash
# 停止当前服务
docker-compose -f docker-compose.prod.yml down

# 恢复数据库备份
bash scripts/restore-database.sh production \
  ./backups/postgres/production/gamelink_<timestamp>.sql.gz

# 恢复到上一个版本
git checkout v0.9.0

# 重新构建和部署
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
```

#### 方案 B：使用 Docker 镜像回滚

```bash
# 停止当前服务
docker-compose -f docker-compose.prod.yml down

# 修改 docker-compose.prod.yml 使用旧版本镜像
# image: gamelink-backend:v0.9.0

# 启动服务
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
```

### 回滚验证

- [ ] 服务正常运行
- [ ] 核心功能可用
- [ ] 数据完整性验证
- [ ] 性能指标正常

---

## 监控和告警

### 关键指标

#### 服务可用性

- API 响应时间
- 错误率
- 服务正常运行时间

#### 系统资源

- CPU 使用率 < 70%
- 内存使用率 < 80%
- 磁盘使用率 < 80%
- 网络流量

#### 数据库

- 连接数
- 查询性能
- 慢查询日志
- 复制延迟（如果使用主从）

#### 应用日志

- 错误日志
- 访问日志
- 慢请求日志

### 监控工具

#### 实时监控

```bash
# 查看容器资源使用
docker stats

# 查看服务日志
docker-compose -f docker-compose.prod.yml logs -f

# 查看特定服务日志
docker logs -f gamelink-backend --tail 100
```

#### 日志分析

```bash
# 查看错误日志
docker logs gamelink-backend 2>&1 | grep ERROR

# 查看访问日志
docker logs nginx 2>&1 | grep POST
```

---

## 应急预案

### 场景 1：服务无法启动

**症状：**
```bash
docker-compose ps
# 显示 Exit 1 或 Restarting
```

**诊断：**
```bash
# 查看日志
docker logs gamelink-backend

# 检查配置
cat .env.production
```

**解决方案：**
1. 检查环境变量配置
2. 检查数据库连接
3. 检查端口占用
4. 查看详细错误日志

### 场景 2：数据库连接失败

**症状：**
API 返回 500 错误，日志显示数据库连接错误

**诊断：**
```bash
# 检查数据库状态
docker exec gamelink-postgres pg_isready -U gamelink

# 检查连接数
docker exec gamelink-postgres psql \
  -U gamelink -d gamelink -c "SELECT count(*) FROM pg_stat_activity;"
```

**解决方案：**
1. 检查数据库是否正常运行
2. 检查连接池配置
3. 检查网络连接
4. 必要时重启数据库

### 场景 3：磁盘空间不足

**症状：**
服务写入失败，日志显示磁盘空间不足

**诊断：**
```bash
# 检查磁盘使用
df -h

# 检查 Docker 占用
docker system df
```

**解决方案：**
1. 清理未使用的 Docker 资源
```bash
docker system prune -a
```

2. 清理旧日志
```bash
docker exec gamelink-backend sh -c "rm /app/logs/*.log"
```

3. 扩展磁盘空间

### 场景 4：内存溢出

**症状：**
服务被杀死，日志显示 OOM

**诊断：**
```bash
# 检查容器内存限制
docker inspect gamelink-backend | grep Memory

# 查看内存使用
docker stats --no-stream
```

**解决方案：**
1. 增加内存限制
2. 优化应用内存使用
3. 检查内存泄漏

### 场景 5：性能下降

**症状：**
API 响应时间变长，用户体验变差

**诊断：**
```bash
# 检查慢查询
docker exec gamelink-postgres psql \
  -U gamelink -d gamelink \
  -c "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# 检查缓存命中率
docker exec gamelink-redis redis-cli info stats
```

**解决方案：**
1. 优化慢查询
2. 增加缓存
3. 扩展服务实例

---

## 部署后检查清单

### 立即检查（部署后 5 分钟）

- [ ] 所有容器状态正常
- [ ] 健康检查通过
- [ ] 无错误日志
- [ ] API 响应正常

### 短期检查（部署后 30 分钟）

- [ ] 监控指标正常
- [ ] 用户反馈正常
- [ ] 核心功能可用
- [ ] 性能指标达标

### 长期检查（部署后 24 小时）

- [ ] 系统稳定运行
- [ ] 无重大问题报告
- [ ] 资源使用正常
- [ ] 备份正常执行

---

## 相关文档

- **部署检查清单：** `docs/DEPLOYMENT_CHECKLIST.md`
- **安全加固指南：** `docs/SECURITY_HARDENING.md`
- **数据库备份恢复：** `docs/DATABASE_BACKUP_RESTORE.md`
- **端口配置指南：** `docs/PORT_CONFIGURATION_GUIDE.md`

---

## 联系人

| 角色 | 负责人 | 职责 |
|------|--------|------|
| **DevOps-Engineer** | DevOps-Engineer | 部署执行、监控 |
| **Backend-Lead** | Backend-Lead | 后端问题处理 |
| **Database-Architect** | Database-Architect | 数据库问题处理 |
| **Team Lead** | Team-Lead | 决策、协调 |

---

**更新历史：**

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-02-09 | 1.0.0 | 初始版本 | DevOps-Engineer |
