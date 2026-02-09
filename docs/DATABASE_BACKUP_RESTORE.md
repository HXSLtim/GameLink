# GameLink 数据库备份恢复操作手册

**文档版本：** 1.0.0
**创建日期：** 2026-02-09
**维护人：** DevOps-Engineer, Database-Architect

---

## 概述

本文档说明 GameLink 项目的数据库备份和恢复操作流程，确保数据安全和灾难恢复能力。

## 备份策略

### 备份类型

| 备份类型 | 频率 | 保留时间 | 用途 |
|---------|------|---------|------|
| **完全备份** | 每日 | 7 天 | 日常数据保护 |
| **部署前备份** | 按需 | 30 天 | 部署前数据保护 |
| **重大变更前备份** | 按需 | 永久 | 重大变更前保护 |

### 备份文件命名规范

```
{数据库名}_{时间戳}.sql.gz

示例：
gamelink_staging_20260209_143022.sql.gz
gamelink_20260209_143022.sql.gz
pre_restore_20260209_143022.sql.gz  # 恢复前的自动备份
```

---

## 备份操作

### 自动化备份（推荐）

使用自动化脚本进行备份：

```bash
# 备份 Staging 环境
bash scripts/backup-database.sh staging

# 备份 Production 环境
bash scripts/backup-database.sh production
```

**脚本功能：**
- ✅ 自动检测容器状态
- ✅ 创建带时间戳的备份文件
- ✅ 压缩备份文件（gzip）
- ✅ 验证备份文件完整性
- ✅ 自动清理旧备份（保留最近 7 个）
- ✅ 显示备份文件大小

### 手动备份

如需手动备份，使用以下命令：

```bash
# 进入容器
docker exec -it gamelink-postgres bash

# 执行备份
pg_dump -U gamelink -d gamelink \
  --no-owner \
  --no-acl \
  --verbose \
  > /tmp/gamelink_backup.sql

# 退出容器并复制备份文件
docker cp gamelink-postgres:/tmp/gamelink_backup.sql ./backups/postgres/

# 压缩备份文件
gzip ./backups/postgres/gamelink_backup.sql
```

---

## 恢复操作

### 自动化恢复（推荐）

使用自动化脚本进行恢复：

```bash
# 恢复 Staging 环境
bash scripts/restore-database.sh staging ./backups/postgres/staging/gamelink_staging_20260209_143022.sql.gz

# 恢复 Production 环境
bash scripts/restore-database.sh production ./backups/postgres/production/gamelink_20260209_143022.sql.gz
```

**脚本功能：**
- ✅ 检查备份文件存在性
- ✅ 显示可用备份文件列表
- ✅ 交互式确认（防止误操作）
- ✅ 自动创建恢复前备份
- ✅ 删除并重建数据库
- ✅ 恢复数据并验证

### 恢复前检查清单

- [ ] 确认备份文件完整（使用 `gzip -t` 验证）
- [ ] 确认容器运行正常
- [ ] 确认有足够的磁盘空间
- [ ] 记录恢复前的表数量
- [ ] 通知相关人员（生产环境）

### 恢复流程

1. **查看可用备份**
   ```bash
   ls -lh ./backups/postgres/staging/
   ```

2. **验证备份文件**
   ```bash
   gzip -t ./backups/postgres/staging/gamelink_staging_20260209_143022.sql.gz
   ```

3. **执行恢复**
   ```bash
   bash scripts/restore-database.sh staging ./backups/postgres/staging/gamelink_staging_20260209_143022.sql.gz
   ```

4. **验证恢复结果**
   ```bash
   docker exec -it gamelink-postgres-staging psql \
     -U gamelink -d gamelink_staging \
     -c "\dt"  # 查看表列表
   ```

---

## 定时备份配置

### 使用 Cron 定时任务

**编辑 crontab：**
```bash
crontab -e
```

**添加定时任务：**
```bash
# 每日凌晨 2 点备份 Staging 环境
0 2 * * * cd /path/to/GameLink && bash scripts/backup-database.sh staging >> /var/log/gamelink-backup.log 2>&1

# 每日凌晨 3 点备份 Production 环境
0 3 * * * cd /path/to/GameLink && bash scripts/backup-database.sh production >> /var/log/gamelink-backup.log 2>&1
```

### 备份监控

创建备份监控脚本 `scripts/check-backup-age.sh`：

```bash
#!/bin/bash
# 检查备份文件是否过旧（超过 48 小时）

BACKUP_DIR="./backups/postgres/production"
MAX_AGE_HOURS=48

if [ -z "$(find "$BACKUP_DIR" -name "*.sql.gz" -mtime -2)" ]; then
    echo "警告：没有最近 48 小时内的备份文件！"
    exit 1
fi
```

---

## 灾难恢复流程

### 场景 1：数据误删除

**恢复步骤：**
1. 立即停止应用写入
2. 找到误删除发生前的备份
3. 恢复备份
4. 应用重放日志（如果启用了 WAL）

### 场景 2：数据库损坏

**恢复步骤：**
1. 停止数据库容器
2. 备份损坏的数据目录
3. 删除并重建容器
4. 恢复最近的备份

### 场景 3：完全灾难恢复

**恢复步骤：**
1. 准备新的服务器
2. 安装 Docker 和 Docker Compose
3. 部署 GameLink 应用
4. 恢复最近的数据库备份
5. 恢复 Redis 数据（如果需要）
6. 验证所有功能正常

---

## 备份存储策略

### 本地存储

**默认位置：**
```
./backups/postgres/
├── staging/
│   └── *.sql.gz
└── production/
    └── *.sql.gz
```

### 远程存储（推荐）

**使用 rsync 同步到远程服务器：**
```bash
rsync -avz --delete \
  ./backups/postgres/ \
  user@remote-server:/backups/gamelink/postgres/
```

**使用云存储（AWS S3/阿里云 OSS）：**
```bash
# 安装 AWS CLI
pip install awscli

# 配置 AWS 凭证
aws configure

# 同步到 S3
aws s3 sync ./backups/postgres/ s3://gamelink-backups/postgres/
```

---

## 最佳实践

### 1. 定期测试恢复

**建议频率：** 每月至少一次

**测试流程：**
1. 在 Staging 环境恢复最新备份
2. 验证数据完整性
3. 运行关键功能测试
4. 记录恢复时间和遇到的问题

### 2. 加密敏感备份

**使用 GPG 加密：**
```bash
# 加密备份文件
gpg --symmetric --cipher-algo AES256 gamelink_backup.sql.gz

# 解密备份文件
gpg --decrypt gamelink_backup.sql.gz.gpg > gamelink_backup.sql.gz
```

### 3. 备份验证

**定期验证备份文件：**
```bash
# 验证 gzip 压缩文件
gzip -t *.sql.gz

# 验证备份内容（不解压）
gunzip -c gamelink_backup.sql.gz | head -n 50
```

### 4. 监控备份大小

**设置告警：**
- 正常备份大小范围
- 异常增长检测（可能指示数据问题）

---

## 常见问题

### Q1: 备份文件损坏怎么办？

**解决方案：**
1. 检查是否有其他可用备份
2. 尝试使用 `gzip -v` 获取详细信息
3. 如果所有备份都损坏，考虑使用 WAL 日志恢复

### Q2: 恢复时磁盘空间不足

**解决方案：**
1. 清理旧备份文件
2. 清理 Docker 未使用的资源：`docker system prune -a`
3. 扩展磁盘空间

### Q3: 恢复后数据不一致

**解决方案：**
1. 检查备份是否在写入期间创建的
2. 考虑使用更早的备份
3. 检查应用程序日志

---

## 相关文档

- **部署检查清单：** `docs/DEPLOYMENT_CHECKLIST.md`
- **安全加固指南：** `docs/SECURITY_HARDENING.md`
- **数据库架构：** Database-Architect 提供的架构文档

---

## 联系人

| 角色 | 负责人 | 职责 |
|------|--------|------|
| **DevOps-Engineer** | DevOps-Engineer | 备份恢复脚本、自动化 |
| **Database-Architect** | Database-Architect | 数据库架构、性能优化 |
| **Backend-Lead** | Backend-Lead | 数据模型、迁移脚本 |

---

**更新历史：**

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-02-09 | 1.0.0 | 初始版本 | DevOps-Engineer |
