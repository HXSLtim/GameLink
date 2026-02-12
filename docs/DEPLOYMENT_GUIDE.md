# GameLink 部署手册

> **版本**: v2.0
> **最后更新**: 2026-02-11
> **适用环境**: 开发/测试/生产

---

## 目录

1. [环境规划](#1-环境规划)
2. [开发环境部署](#2-开发环境部署)
3. [测试环境部署](#3-测试环境部署)
4. [生产环境部署](#4-生产环境部署)
5. [数据库迁移](#5-数据库迁移)
6. [回滚方案](#6-回滚方案)
7. [监控告警](#7-监控告警)
8. [常见问题](#8-常见问题)

---

## 1. 环境规划

### 1.1 环境配置对比

| 配置项 | 开发环境 | 测试环境 | 生产环境 |
|--------|----------|----------|----------|
| APP_ENV | development | staging | production |
| GIN_MODE | debug | debug | release |
| CRYPTO_ENABLED | false | true | true |
| SEED_ENABLED | true | false | false |
| 日志级别 | DEBUG | INFO | WARN |
| 数据库 | 本地Docker | 云数据库 | 云数据库主从 |
| Redis | 本地Docker | 云缓存 | 云缓存集群 |

### 1.2 端口规划

| 服务 | 开发环境 | 生产环境 |
|------|----------|----------|
| 后端 API | 8081 | 8080 |
| 管理后台 | 5173 | 80/443 (Nginx) |
| 小程序 H5 | 5174 | 80/443 (Nginx) |
| PostgreSQL | 5432 | 内网 |
| Redis | 6379 | 内网 |
| WebSocket | 8081/ws | 8080/ws |

---

## 2. 开发环境部署

### 2.1 前置要求

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose
- Git

### 2.2 启动基础服务

```bash
# 启动 PostgreSQL 和 Redis
docker compose up -d postgres redis

# 验证服务
docker compose ps
```

### 2.3 配置环境变量

```bash
# 复制配置模板
cp .env.example .env

# 编辑配置（开发环境默认值即可）
# 重点配置：
# - 数据库连接
# - Redis 连接
# - JWT 密钥（开发环境可用简单值）
```

### 2.4 启动后端

```bash
cd api

# 安装依赖
go mod download

# 运行迁移（如有）
go run cmd/migrate/main.go

# 启动服务
go run cmd/main.go

# 验证
curl http://localhost:8081/health
```

### 2.5 启动管理后台

```bash
cd admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

### 2.6 启动小程序/H5

```bash
cd app

# 安装依赖
npm install

# H5 开发
npm run dev:h5
# 访问 http://localhost:5174

# 微信小程序
npm run dev:mp-weixin
# 用微信开发者工具打开 dist/dev/mp-weixin
```

---

## 3. 测试环境部署

### 3.1 服务器准备

```bash
# 系统要求
# - Ubuntu 22.04+ / CentOS 8+
# - 2核4GB 起步
# - 50GB 存储

# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 3.2 配置环境变量

```bash
# 生产级配置
cat > .env.staging << EOF
# 应用配置
APP_ENV=staging
GIN_MODE=debug
PORT=8080

# 数据库（使用云数据库）
DB_HOST=your-staging-db.xxx.rds.aliyuncs.com
DB_PORT=5432
DB_NAME=gamelink_staging
DB_USER=gamelink
DB_PASSWORD=your_strong_password

# Redis
REDIS_HOST=your-staging-redis.xxx.redis.rds.aliyuncs.com
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# 加密配置
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=your_32_byte_secret_key_here
CRYPTO_IV=your_16_byte_iv_here

# JWT
JWT_SECRET_KEY=your_jwt_secret_key_min_32_chars
JWT_EXPIRATION=24h

# 第三方服务（测试环境配置）
WECHAT_APP_ID=your_wechat_app_id
WECHAT_APP_SECRET=your_wechat_app_secret
# ... 其他配置
EOF
```

### 3.3 使用 Docker Compose 部署

```bash
# 构建镜像
docker compose -f docker-compose.staging.yml build

# 启动服务
docker compose -f docker-compose.staging.yml up -d

# 查看日志
docker compose -f docker-compose.staging.yml logs -f
```

### 3.4 验证部署

```bash
# 健康检查
curl https://staging-api.yourdomain.com/health

# API 文档访问
# https://staging-api.yourdomain.com/swagger/index.html

# 管理后台访问
# https://staging-admin.yourdomain.com

# H5 访问
# https://staging-h5.yourdomain.com
```

---

## 4. 生产环境部署

### 4.1 部署前检查清单

#### 4.1.1 配置检查

- [ ] `APP_ENV=production`
- [ ] `GIN_MODE=release`
- [ ] `CRYPTO_ENABLED=true`
- [ ] `CRYPTO_SECRET_KEY` 已设置（32字节随机值）
- [ ] `CRYPTO_IV` 已设置（16字节随机值）
- [ ] `JWT_SECRET_KEY` 已设置（强随机值）
- [ ] `SUPER_ADMIN_EMAIL` 已设置
- [ ] `SUPER_ADMIN_PASSWORD` 已设置（强密码）
- [ ] `SEED_ENABLED=false`
- [ ] 数据库密码已修改为强密码
- [ ] Redis 密码已设置

#### 4.1.2 第三方服务配置

- [ ] 微信支付：商户号、API密钥、证书
- [ ] 支付宝：应用ID、私钥、公钥
- [ ] 微信小程序：AppID、AppSecret
- [ ] 短信服务：AccessKey、签名
- [ ] OSS：Bucket、AccessKey

#### 4.1.3 安全检查

- [ ] 数据库只允许应用服务器访问
- [ ] Redis 设置密码且只允许内网访问
- [ ] API 服务配置 CORS 白名单
- [ ] 启用 HTTPS（SSL 证书）
- [ ] 配置防火墙规则
- [ ] 敏感信息不记录在日志中

### 4.2 数据库部署

```sql
-- 创建生产数据库
CREATE DATABASE gamelink_prod
  ENCODING 'UTF8'
  LC_COLLATE='en_US.UTF-8'
  LC_CTYPE='en_US.UTF-8';

-- 创建用户（强密码）
CREATE USER gamelink_prod WITH PASSWORD 'your_very_strong_password';

-- 授权
GRANT ALL PRIVILEGES ON DATABASE gamelink_prod TO gamelink_prod;

-- 启用必要扩展
\c gamelink_prod
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### 4.3 Redis 配置

```conf
# /etc/redis/redis.conf

# 绑定内网IP
bind 127.0.0.1 10.0.0.1

# 设置密码
requirepass your_redis_password

# 禁用危险命令
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG ""
```

### 4.4 Nginx 配置

```nginx
# /etc/nginx/sites-available/gamelink-api.conf

# 后端 API
upstream api_backend {
    server 127.0.0.1:8080;
    keepalive 64;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL 证书
    ssl_certificate /etc/ssl/certs/yourdomain.crt;
    ssl_certificate_key /etc/ssl/private/yourdomain.key;
    ssl_protocols TLSv1.2 TLSv1.3;

    # API 代理
    location /api/ {
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket
    location /ws {
        proxy_pass http://api_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }

    # Swagger 文档（生产环境可选择关闭）
    location /swagger {
        proxy_pass http://api_backend;
    }
}

# 管理后台静态文件
server {
    listen 443 ssl http2;
    server_name admin.yourdomain.com;

    root /var/www/gamelink/admin;
    index index.html;

    # SPA 路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
}

# H5 静态文件
server {
    listen 443 ssl http2;
    server_name h5.yourdomain.com;

    root /var/www/gamelink/h5;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

### 4.5 应用部署

```bash
# 1. 拉取代码
cd /opt/gamelink
git pull origin main

# 2. 构建前端
cd admin
npm ci
npm run build
cp -r dist/* /var/www/gamelink/admin/

cd ../app
npm ci
npm run build:h5
cp -r/dist/dev/h5/* /var/www/gamelink/h5/

# 3. 构建后端
cd ../api
go build -o gamelink-api cmd/main.go

# 4. 重启服务（使用 systemd）
sudo systemctl restart gamelink-api

# 5. 检查状态
sudo systemctl status gamelink-api
```

### 4.6 Systemd 配置

```ini
# /etc/systemd/system/gamelink-api.service

[Unit]
Description=GameLink API Service
After=network.target postgresql.service

[Service]
Type=simple
User=gamelink
Group=gamelink
WorkingDirectory=/opt/gamelink/api
ExecStart=/opt/gamelink/api/gamelink-api
Restart=always
RestartSec=5

# 环境变量文件
EnvironmentFile=/opt/gamelink/.env.production

# 安全配置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/gamelink /tmp

[Install]
WantedBy=multi-user.target
```

---

## 5. 数据库迁移

### 5.1 创建迁移文件

```bash
cd api
go run cmd/migrate/main.go create add_user_timezone
```

### 5.2 执行迁移

```bash
# 开发环境
go run cmd/migrate/main.go up

# 生产环境（先备份！）
pg_dump gamelink_prod > backup_$(date +%Y%m%d_%H%M%S).sql
go run cmd/migrate/main.go up
```

### 5.3 回滚迁移

```bash
go run cmd/migrate/main.go down
```

---

## 6. 回滚方案

### 6.1 应用回滚

```bash
# 1. 停止服务
sudo systemctl stop gamelink-api

# 2. 切换到上一版本
cd /opt/gamelink
git reset --hard HEAD~1

# 3. 重新构建
cd api
go build -o gamelink-api cmd/main.go

# 4. 重启服务
sudo systemctl start gamelink-api

# 5. 验证
curl https://api.yourdomain.com/health
```

### 6.2 数据库回滚

```bash
# 1. 停止应用
sudo systemctl stop gamelink-api

# 2. 恢复备份
psql -U gamelink_prod -d gamelink_prod < backup_20260211_120000.sql

# 3. 回滚迁移
cd /opt/gamelink/api
go run cmd/migrate/main.go down

# 4. 启动应用
sudo systemctl start gamelink-api
```

### 6.3 快速回滚脚本

```bash
#!/bin/bash
# rollback.sh - 快速回滚脚本

BACKUP_DIR="/opt/backups/gamelink"
DB_NAME="gamelink_prod"

echo "开始回滚..."

# 1. 停止服务
systemctl stop gamelink-api

# 2. 恢复最新备份
LATEST_BACKUP=$(ls -t $BACKUP_DIR/*.sql | head -1)
echo "恢复备份: $LATEST_BACKUP"
psql -U gamelink_prod -d $DB_NAME < $LATEST_BACKUP

# 3. 回滚代码
cd /opt/gamelink
git reset --hard HEAD~1

# 4. 重新构建
cd api
go build -o gamelink-api cmd/main.go

# 5. 启动服务
systemctl start gamelink-api

echo "回滚完成！"
```

---

## 7. 监控告警

### 7.1 Prometheus 配置

```yaml
# prometheus.yml

scrape_configs:
  - job_name: 'gamelink-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s

  - job_name: 'postgres'
    static_configs:
      - targets: ['localhost:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['localhost:9121']
```

### 7.2 告警规则

```yaml
# alert_rules.yml

groups:
  - name: gamelink_alerts
    rules:
      # API 可用性
      - alert: APIDown
        expr: up{job="gamelink-api"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "API 服务不可用"

      # 高错误率
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API 错误率过高"

      # 数据库连接
      - alert: DatabaseConnectionHigh
        expr: pg_stat_database_numbackends{datname="gamelink_prod"} > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "数据库连接数过高"
```

### 7.3 Grafana 仪表盘

- 系统概览：CPU、内存、磁盘、网络
- API 指标：请求量、响应时间、错误率
- 数据库：连接数、查询性能、慢查询
- Redis：内存使用、命中率、连接数
- 业务指标：订单量、支付成功率、在线用户

---

## 8. 常见问题

### 8.1 启动失败

**问题**: 后端启动失败，提示数据库连接错误

**排查**:
```bash
# 检查数据库是否启动
systemctl status postgresql

# 测试连接
psql -h localhost -U gamelink -d gamelink_dev

# 检查防火墙
sudo ufw status
```

### 8.2 端口冲突

**问题**: 端口已被占用

**解决**:
```bash
# 查找占用进程
sudo lsof -i :8080

# 杀死进程
sudo kill -9 <PID>
```

### 8.3 权限问题

**问题**: 无法写入日志或上传文件

**解决**:
```bash
# 创建目录并设置权限
sudo mkdir -p /var/log/gamelink
sudo chown gamelink:gamelink /var/log/gamelink
sudo chmod 755 /var/log/gamelink
```

### 8.4 SSL 证书

**问题**: HTTPS 证书过期

**解决**:
```bash
# 使用 Let's Encrypt 自动续期
sudo certbot renew
sudo systemctl reload nginx
```

### 8.5 性能问题

**问题**: 响应慢

**排查**:
```bash
# 查看 API 日志
journalctl -u gamelink-api -f

# 查看数据库慢查询
sudo tail -f /var/log/postgresql/postgresql-slow.log

# 查看系统资源
htop
iostat -x 1
```

---

## 附录

### A. 环境变量完整列表

```bash
# 应用配置
APP_ENV=production
GIN_MODE=release
PORT=8080
ALLOWED_ORIGINS=https://admin.yourdomain.com,https://h5.yourdomain.com

# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_NAME=gamelink_prod
DB_USER=gamelink
DB_PASSWORD=your_password
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_DB=0

# JWT
JWT_SECRET_KEY=your_32_byte_secret_key_here
JWT_EXPIRATION=24h

# 加密
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=your_32_byte_secret_key_here
CRYPTO_IV=your_16_byte_iv_here

# 超级管理员
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=your_strong_password

# 种子数据
SEED_ENABLED=false

# 第三方服务
WECHAT_APP_ID=
WECHAT_APP_SECRET=
WECHAT_PAY_MCH_ID=
WECHAT_PAY_API_KEY=
ALIPAY_APP_ID=
ALIPAY_PRIVATE_KEY=
# ... 其他配置
```

### B. 部署检查清单

- [ ] 代码已合并到 main 分支
- [ ] 所有测试通过
- [ ] 环境变量已配置
- [ ] 数据库已迁移
- [ ] SSL 证书有效
- [ ] 监控告警已配置
- [ ] 备份已创建
- [ ] 回滚方案已准备

---

**文档维护**: 产品经理 + DevOps
**更新频率**: 每次部署后更新
