# 🚀 GameLink 部署和运维指南

本文档提供 GameLink 项目的完整部署方案，包括开发、测试和生产环境的配置。

---

## 📋 目录

- [部署概览](#部署概览)
- [环境要求](#环境要求)
- [开发环境部署](#开发环境部署)
- [Docker 部署](#docker-部署)
- [生产环境部署](#生产环境部署)
- [CI/CD 配置](#cicd-配置)
- [监控和日志](#监控和日志)
- [运维操作](#运维操作)
- [安全配置](#安全配置)
- [性能优化](#性能优化)
- [故障排查](#故障排查)

---

## 🎯 部署概览

### 部署架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   Web Server    │    │   Database      │
│                │    │                │    │                │
│ • Nginx        │◄──►│ • Go Services   │◄──►│ • MySQL         │
│ • SSL/TLS      │    │ • React App     │    │ • Redis         │
│ • HTTP/2       │    │ • Static Files  │    │ • Backups       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CDN           │    │   Monitoring    │    │   Storage       │
│                │    │                │    │                │
│ • Static Assets │    │ • Prometheus    │    │ • File Storage  │
│ • Cache         │    │ • Grafana       │    │ • Image Upload  │
│ • DDoS Protect  │    │ • Alerts        │    │ • CDN           │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 服务组件

**后端服务**
- `user-service`: 用户管理服务 (端口 8080)
- `order-service`: 订单管理服务 (端口 8081)
- `payment-service`: 支付服务 (端口 8082)
- `notification-service`: 通知服务 (端口 8083)

**前端应用**
- `Web App`: React 单页应用
- `Admin Panel`: 管理后台
- `Mobile PWA`: 移动端渐进式应用

**基础设施**
- `MySQL`: 主数据库
- `Redis`: 缓存和会话存储
- `Nginx`: 反向代理和负载均衡
- `Prometheus`: 监控系统
- `Grafana`: 可视化面板

---

## 🔧 环境要求

### 最低配置
- **CPU**: 2 cores
- **内存**: 4GB RAM
- **存储**: 20GB SSD
- **网络**: 10 Mbps

### 推荐配置
- **CPU**: 4 cores
- **内存**: 8GB RAM
- **存储**: 50GB SSD
- **网络**: 100 Mbps

### 生产环境配置
- **CPU**: 8 cores
- **内存**: 16GB RAM
- **存储**: 100GB SSD
- **网络**: 1 Gbps

### 软件依赖

**操作系统**
- Ubuntu 20.04 LTS / CentOS 8 / RHEL 8
- Windows Server 2019+ (可选)
- macOS (开发环境)

**运行时环境**
- Go 1.25.3+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- Nginx 1.18+

**容器环境**
- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.24+ (生产环境)

---

## 🏠 开发环境部署

### 本地快速启动

#### 1. 环境准备
```bash
# 安装必需工具
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.12.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

#### 2. 克隆项目
```bash
git clone https://github.com/your-org/GameLink.git
cd GameLink
```

#### 3. 配置环境变量
```bash
# 复制环境配置文件
cp .env.example .env
cp docker-compose.example.yml docker-compose.yml

# 编辑配置文件
nano .env
```

**开发环境配置 (`.env`)**
```env
# 应用配置
APP_ENV=development
APP_VERSION=2.1.0
DEBUG=true

# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_NAME=gamelink_dev
DB_USER=gamelink
DB_PASSWORD=dev_password_123
DB_ROOT_PASSWORD=root_password_123

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=dev_jwt_secret_key_change_in_production
JWT_EXPIRE_HOURS=24

# 服务端口
API_PORT=8080
WEB_PORT=5173

# 文件上传
UPLOAD_MAX_SIZE=10485760
UPLOAD_PATH=./uploads

# 日志配置
LOG_LEVEL=debug
LOG_FILE=./logs/app.log

# 外部服务 (可选)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

#### 4. 启动服务
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

#### 5. 初始化数据库
```bash
# 运行数据库迁移
docker-compose exec api make migrate

# 插入测试数据
docker-compose exec api make seed
```

#### 6. 验证部署
访问以下地址验证部署是否成功：
- 前端应用: http://localhost:5173
- 后端API: http://localhost:8080/health
- API文档: http://localhost:8080/swagger/index.html
- 管理后台: http://localhost:5173/admin

### 单独服务部署

#### 后端服务
```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量
export DB_HOST=localhost
export DB_PASSWORD=your_password
export JWT_SECRET=your_jwt_secret

# 运行数据库迁移
make migrate

# 启动服务
make run CMD=user-service
```

#### 前端应用
```bash
cd frontend

# 安装依赖
npm install

# 配置环境变量
echo "VITE_API_BASE_URL=http://localhost:8080/api/v1" > .env.local

# 启动开发服务器
npm run dev
```

---

## 🐳 Docker 部署

### Docker Compose 配置

**`docker-compose.yml`**
```yaml
version: '3.8'

services:
  # MySQL 数据库
  mysql:
    image: mysql:8.0
    container_name: gamelink_mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/sql:/docker-entrypoint-initdb.d
    networks:
      - gamelink_network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: gamelink_redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - gamelink_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      timeout: 3s
      retries: 5

  # 后端 API 服务
  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: gamelink_api
    restart: unless-stopped
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=${JWT_SECRET}
      - APP_ENV=${APP_ENV}
    ports:
      - "${API_PORT}:8080"
    volumes:
      - ./uploads:/app/uploads
      - ./logs:/app/logs
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - gamelink_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      timeout: 10s
      retries: 5

  # 前端应用
  web:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: gamelink_web
    restart: unless-stopped
    environment:
      - VITE_API_BASE_URL=http://localhost:${API_PORT}/api/v1
    ports:
      - "${WEB_PORT}:80"
    depends_on:
      - api
    networks:
      - gamelink_network

  # Nginx 反向代理
  nginx:
    image: nginx:alpine
    container_name: gamelink_nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - ./uploads:/var/www/uploads
    depends_on:
      - api
      - web
    networks:
      - gamelink_network

volumes:
  mysql_data:
    driver: local
  redis_data:
    driver: local

networks:
  gamelink_network:
    driver: bridge
```

### 后端 Dockerfile

**`backend/Dockerfile`**
```dockerfile
# 构建阶段
FROM golang:1.25.3-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/user-service

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates curl

# 创建用户和目录
RUN adduser -D -s /bin/sh appuser
RUN mkdir -p /app/uploads /app/logs
RUN chown -R appuser:appuser /app

# 复制构建的二进制文件
COPY --from=builder /app/main /app/
COPY --from=builder /app/configs /app/configs

# 切换到非 root 用户
USER appuser

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"]
```

### 前端 Dockerfile

**`frontend/Dockerfile`**
```dockerfile
# 构建阶段
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制构建结果
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
```

### Nginx 配置

**`nginx/nginx.conf`**
```nginx
events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log;

    # 基本设置
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # 文件上传大小限制
    client_max_body_size 10M;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;

    # 上游服务器
    upstream api_backend {
        server api:8080;
    }

    # HTTP 服务器
    server {
        listen 80;
        server_name localhost;

        # 重定向到 HTTPS (生产环境)
        # return 301 https://$server_name$request_uri;

        # 前端静态文件
        location / {
            proxy_pass http://web;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # API 接口
        location /api/ {
            proxy_pass http://api_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # WebSocket 支持
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }

        # WebSocket 连接
        location /ws/ {
            proxy_pass http://api_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # 文件上传
        location /uploads/ {
            alias /var/www/uploads/;
            expires 1y;
            add_header Cache-Control "public, immutable";
        }

        # 健康检查
        location /health {
            proxy_pass http://api_backend;
            access_log off;
        }
    }

    # HTTPS 服务器 (生产环境)
    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        # SSL 证书
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        # SSL 设置
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;
        ssl_session_cache shared:SSL:10m;

        # HSTS
        add_header Strict-Transport-Security "max-age=63072000" always;

        # 其他配置与 HTTP 相同...
    }
}
```

---

## 🏭 生产环境部署

### 服务器准备

#### 1. 系统初始化
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git vim htop

# 配置防火墙
sudo ufw enable
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 创建应用用户
sudo useradd -m -s /bin/bash gamelink
sudo usermod -aG sudo gamelink
```

#### 2. Docker 安装
```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 添加用户到 docker 组
sudo usermod -aG docker gamelink
```

#### 3. 目录结构
```bash
# 创建应用目录
sudo mkdir -p /opt/gamelink
sudo mkdir -p /opt/gamelink/app
sudo mkdir -p /opt/gamelink/data/mysql
sudo mkdir -p /opt/gamelink/data/redis
sudo mkdir -p /opt/gamelink/logs
sudo mkdir -p /opt/gamelink/uploads
sudo mkdir -p /opt/gamelink/ssl

# 设置权限
sudo chown -R gamelink:gamelink /opt/gamelink
```

### 生产环境配置

#### 1. 环境变量配置
```bash
# 生产环境配置 (.env.prod)
APP_ENV=production
APP_VERSION=2.1.0
DEBUG=false

# 数据库配置 (使用强密码)
DB_HOST=mysql
DB_PORT=3306
DB_NAME=gamelink_prod
DB_USER=gamelink_prod
DB_PASSWORD=your_super_secure_password_here
DB_ROOT_PASSWORD=your_super_secure_root_password_here

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password_here

# JWT 配置 (使用强密钥)
JWT_SECRET=your_super_secure_jwt_secret_key_here_min_32_chars
JWT_EXPIRE_HOURS=24

# 服务端口
API_PORT=8080
WEB_PORT=5173

# 文件上传
UPLOAD_MAX_SIZE=10485760
UPLOAD_PATH=/opt/gamelink/uploads

# 日志配置
LOG_LEVEL=info
LOG_FILE=/opt/gamelink/logs/app.log

# SSL 配置
SSL_CERT_PATH=/opt/gamelink/ssl/cert.pem
SSL_KEY_PATH=/opt/gamelink/ssl/key.pem

# 邮件配置
SMTP_HOST=smtp.your-domain.com
SMTP_PORT=587
SMTP_USER=noreply@your-domain.com
SMTP_PASSWORD=your_smtp_password

# 监控配置
PROMETHEUS_URL=http://prometheus:9090
GRAFANA_URL=http://grafana:3000
```

#### 2. 生产环境 Docker Compose
**`docker-compose.prod.yml`**
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: gamelink_mysql_prod
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    volumes:
      - /opt/gamelink/data/mysql:/var/lib/mysql
      - ./scripts/sql:/docker-entrypoint-initdb.d
    networks:
      - gamelink_network
    command: --default-authentication-plugin=mysql_native_password

  redis:
    image: redis:7-alpine
    container_name: gamelink_redis_prod
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - /opt/gamelink/data/redis:/data
    networks:
      - gamelink_network

  api:
    image: gamelink/api:latest
    container_name: gamelink_api_prod
    restart: always
    environment:
      - DB_HOST=mysql
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=redis
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - APP_ENV=production
    volumes:
      - /opt/gamelink/uploads:/app/uploads
      - /opt/gamelink/logs:/app/logs
    depends_on:
      - mysql
      - redis
    networks:
      - gamelink_network
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '1.0'
          memory: 1G

  web:
    image: gamelink/web:latest
    container_name: gamelink_web_prod
    restart: always
    networks:
      - gamelink_network
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G

  nginx:
    image: nginx:alpine
    container_name: gamelink_nginx_prod
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.prod.conf:/etc/nginx/nginx.conf
      - /opt/gamelink/ssl:/etc/nginx/ssl
      - /opt/gamelink/uploads:/var/www/uploads
      - /opt/gamelink/logs/nginx:/var/log/nginx
    depends_on:
      - api
      - web
    networks:
      - gamelink_network

  # 监控服务
  prometheus:
    image: prom/prometheus:latest
    container_name: gamelink_prometheus
    restart: always
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - gamelink_network

  grafana:
    image: grafana/grafana:latest
    container_name: gamelink_grafana
    restart: always
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin123
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana:/etc/grafana/provisioning
    networks:
      - gamelink_network

volumes:
  prometheus_data:
  grafana_data:

networks:
  gamelink_network:
    driver: bridge
```

### 部署脚本

**`scripts/deploy.sh`**
```bash
#!/bin/bash

set -e

# 配置变量
APP_NAME="gamelink"
DEPLOY_USER="gamelink"
DEPLOY_PATH="/opt/gamelink"
BACKUP_PATH="/opt/backups"
LOG_FILE="/var/log/deploy.log"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >> $LOG_FILE
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1" >> $LOG_FILE
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1" >> $LOG_FILE
}

# 检查权限
check_permissions() {
    if [ "$EUID" -ne 0 ]; then
        error "请使用 root 权限运行此脚本"
        exit 1
    fi
}

# 备份数据
backup_data() {
    log "开始备份数据..."

    BACKUP_DIR="$BACKUP_PATH/$(date +%Y%m%d_%H%M%S)"
    mkdir -p $BACKUP_DIR

    # 备份数据库
    docker exec gamelink_mysql_prod mysqldump \
        -u root -p${DB_ROOT_PASSWORD} \
        ${DB_NAME} > $BACKUP_DIR/database.sql

    # 备份上传文件
    cp -r /opt/gamelink/uploads $BACKUP_DIR/

    log "数据备份完成: $BACKUP_DIR"
}

# 构建镜像
build_images() {
    log "开始构建 Docker 镜像..."

    cd $DEPLOY_PATH

    # 构建后端镜像
    docker build -t gamelink/api:latest ./backend/

    # 构建前端镜像
    docker build -t gamelink/web:latest ./frontend/

    log "Docker 镜像构建完成"
}

# 数据库迁移
migrate_database() {
    log "开始数据库迁移..."

    docker exec gamelink_api_prod make migrate

    log "数据库迁移完成"
}

# 更新服务
update_services() {
    log "开始更新服务..."

    cd $DEPLOY_PATH

    # 拉取最新镜像
    docker-compose -f docker-compose.prod.yml pull

    # 重启服务
    docker-compose -f docker-compose.prod.yml up -d

    # 等待服务启动
    sleep 30

    log "服务更新完成"
}

# 健康检查
health_check() {
    log "开始健康检查..."

    # 检查 API 服务
    if curl -f http://localhost/api/v1/health > /dev/null 2>&1; then
        log "API 服务健康检查通过"
    else
        error "API 服务健康检查失败"
        exit 1
    fi

    # 检查前端服务
    if curl -f http://localhost/ > /dev/null 2>&1; then
        log "前端服务健康检查通过"
    else
        error "前端服务健康检查失败"
        exit 1
    fi

    log "所有服务健康检查通过"
}

# 清理旧镜像
cleanup() {
    log "开始清理旧镜像..."

    # 删除未使用的镜像
    docker image prune -f

    # 删除旧版本镜像 (保留最近3个版本)
    docker images --format "table {{.Repository}}:{{.Tag}}" | \
        grep gamelink | tail -n +4 | \
        awk '{print $1}' | xargs -r docker rmi

    log "清理完成"
}

# 主函数
main() {
    log "开始部署 $APP_NAME..."

    check_permissions
    backup_data
    build_images
    migrate_database
    update_services
    health_check
    cleanup

    log "部署完成！"
    log "访问地址: http://your-domain.com"
}

# 执行主函数
main "$@"
```

### 监控配置

#### 1. Prometheus 配置
**`monitoring/prometheus.yml`**
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'gamelink-api'
    static_configs:
      - targets: ['api:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx:80']
    metrics_path: '/metrics'

  - job_name: 'mysql'
    static_configs:
      - targets: ['mysql:3306']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

#### 2. Grafana 仪表盘
```json
{
  "dashboard": {
    "title": "GameLink 监控面板",
    "panels": [
      {
        "title": "API 请求量",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "响应时间",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "数据库连接数",
        "type": "singlestat",
        "targets": [
          {
            "expr": "mysql_global_status_threads_connected"
          }
        ]
      }
    ]
  }
}
```

---

## 🔄 CI/CD 配置

### GitHub Actions

**`.github/workflows/deploy.yml`**
```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: gamelink_test
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s --health-timeout=5s --health-retries=3

      redis:
        image: redis:7
        ports:
          - 6379:6379
        options: --health-cmd="redis-cli ping" --health-interval=10s --health-timeout=5s --health-retries=3

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.25.3

    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: 18
        cache: 'npm'
        cache-dependency-path: frontend/package-lock.json

    - name: Install dependencies
      run: |
        cd backend && go mod download
        cd ../frontend && npm ci

    - name: Run tests
      run: |
        cd backend && go test ./... -v -cover
        cd ../frontend && npm test

    - name: Run linting
      run: |
        cd backend && golangci-lint run
        cd ../frontend && npm run lint

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}

    - name: Build and push API image
      uses: docker/build-push-action@v4
      with:
        context: ./backend
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/api:latest
        labels: ${{ steps.meta.outputs.labels }}

    - name: Build and push Web image
      uses: docker/build-push-action@v4
      with:
        context: ./frontend
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/web:latest
        labels: ${{ steps.meta.outputs.labels }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
    - uses: actions/checkout@v3

    - name: Deploy to production
      uses: appleboy/ssh-action@v0.1.5
      with:
        host: ${{ secrets.HOST }}
        username: ${{ secrets.USERNAME }}
        key: ${{ secrets.SSH_KEY }}
        script: |
          cd /opt/gamelink
          git pull origin main
          ./scripts/deploy.sh

    - name: Notify deployment
      uses: 8398a7/action-slack@v3
      with:
        status: ${{ job.status }}
        channel: '#deployments'
        text: '部署完成: ${{ github.sha }}'
      env:
        SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

---

## 📊 监控和日志

### 日志管理

#### 1. 应用日志配置
```go
// logger/logger.go
package logger

import (
    "io"
    "os"
    "path/filepath"
    "time"

    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func Init(logLevel string, logFile string) error {
    Logger = logrus.New()

    // 设置日志级别
    level, err := logrus.ParseLevel(logLevel)
    if err != nil {
        return err
    }
    Logger.SetLevel(level)

    // 设置日志格式
    Logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })

    // 设置日志输出
    if logFile != "" {
        // 确保日志目录存在
        if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
            return err
        }

        // 日志轮转配置
        output := &lumberjack.Logger{
            Filename:   logFile,
            MaxSize:    100, // MB
            MaxBackups: 3,
            MaxAge:     28, // days
            Compress:   true,
        }

        Logger.SetOutput(output)
    } else {
        Logger.SetOutput(os.Stdout)
    }

    return nil
}
```

#### 2. 日志收集配置
```yaml
# docker-compose.logging.yml
version: '3.8'

services:
  # ELK Stack for log management
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    networks:
      - monitoring

  logstash:
    image: docker.elastic.co/logstash/logstash:8.5.0
    container_name: logstash
    volumes:
      - ./monitoring/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
      - /opt/gamelink/logs:/logs
    ports:
      - "5044:5044"
    depends_on:
      - elasticsearch
    networks:
      - monitoring

  kibana:
    image: docker.elastic.co/kibana/kibana:8.5.0
    container_name: kibana
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    networks:
      - monitoring

  # Filebeat for log shipping
  filebeat:
    image: docker.elastic.co/beats/filebeat:8.5.0
    container_name: filebeat
    user: root
    volumes:
      - ./monitoring/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /opt/gamelink/logs:/logs:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    depends_on:
      - logstash
    networks:
      - monitoring

volumes:
  elasticsearch_data:

networks:
  monitoring:
    driver: bridge
```

### 监控指标

#### 1. 应用指标
```go
// metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP 请求计数器
    HttpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    // HTTP 请求持续时间
    HttpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    // 数据库连接池
    DatabaseConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Number of active database connections",
        },
    )

    // 订单计数器
    OrdersTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orders_total",
            Help: "Total number of orders",
        },
        []string{"status"},
    )

    // 用户在线数
    UsersOnline = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "users_online",
            Help: "Number of currently online users",
        },
    )
)
```

#### 2. 告警规则
```yaml
# monitoring/rules/gamelink.yml
groups:
  - name: gamelink.rules
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "高错误率告警"
          description: "5xx 错误率超过 10%"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "响应时间过长"
          description: "95% 请求响应时间超过 1 秒"

      - alert: DatabaseConnectionsHigh
        expr: database_connections_active > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "数据库连接数过高"
          description: "活跃数据库连接数: {{ $value }}"

      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "服务不可用"
          description: "{{ $labels.instance }} 服务已停止"
```

---

## 🔧 运维操作

### 日常维护

#### 1. 数据库维护
```bash
# 数据库备份脚本
#!/bin/bash
# backup.sh

BACKUP_DIR="/opt/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="gamelink_prod"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
docker exec gamelink_mysql_prod mysqldump \
    -u root -p$DB_ROOT_PASSWORD \
    --single-transaction \
    --routines \
    --triggers \
    $DB_NAME | gzip > $BACKUP_DIR/gamelink_$DATE.sql.gz

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "数据库备份完成: $BACKUP_DIR/gamelink_$DATE.sql.gz"
```

```bash
# 数据库优化脚本
#!/bin/bash
# optimize.sh

docker exec gamelink_mysql_prod mysql -u root -p$DB_ROOT_PASSWORD -e "
    OPTIMIZE TABLE users;
    OPTIMIZE TABLE orders;
    OPTIMIZE TABLE payments;
    ANALYZE TABLE users;
    ANALYZE TABLE orders;
    ANALYZE TABLE payments;
"

echo "数据库优化完成"
```

#### 2. 日志轮转
```bash
# /etc/logrotate.d/gamelink
/opt/gamelink/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 gamelink gamelink
    postrotate
        docker kill -s USR1 gamelink_api_prod
    endscript
}
```

#### 3. 系统监控脚本
```bash
#!/bin/bash
# monitor.sh

# 检查服务状态
check_service() {
    local service=$1
    local url=$2

    if curl -f $url > /dev/null 2>&1; then
        echo "✅ $service 服务正常"
    else
        echo "❌ $service 服务异常"
        # 发送告警
        send_alert "$service 服务异常"
    fi
}

# 检查磁盘空间
check_disk_space() {
    local usage=$(df /opt/gamelink | awk 'NR==2 {print $5}' | sed 's/%//')

    if [ $usage -gt 80 ]; then
        echo "⚠️  磁盘使用率过高: ${usage}%"
        send_alert "磁盘使用率过高: ${usage}%"
    else
        echo "✅ 磁盘使用率正常: ${usage}%"
    fi
}

# 发送告警
send_alert() {
    local message=$1
    # 发送邮件、Slack 或其他通知
    echo "$message" | mail -s "GameLink 告警" admin@your-domain.com
}

# 主检查函数
main() {
    echo "=== 系统监控检查 $(date) ==="

    check_service "API" "http://localhost/api/v1/health"
    check_service "前端" "http://localhost/"
    check_disk_space

    echo "检查完成"
}

main
```

### 扩容操作

#### 1. 水平扩容
```bash
# 扩容 API 服务实例
docker-compose -f docker-compose.prod.yml up -d --scale api=3

# 负载均衡配置更新
# 更新 nginx 配置以包含新的后端实例
```

#### 2. 垂直扩容
```bash
# 更新 docker-compose.yml 中的资源限制
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '4.0'
          memory: 4G
        reservations:
          cpus: '2.0'
          memory: 2G
```

### 故障恢复

#### 1. 服务重启
```bash
#!/bin/bash
# restart_service.sh

SERVICE_NAME=$1

if [ -z "$SERVICE_NAME" ]; then
    echo "Usage: $0 <service_name>"
    exit 1
fi

echo "重启服务: $SERVICE_NAME"

# 优雅关闭
docker-compose -f docker-compose.prod.yml stop $SERVICE_NAME

# 等待服务完全停止
sleep 10

# 启动服务
docker-compose -f docker-compose.prod.yml start $SERVICE_NAME

# 健康检查
sleep 30
if curl -f http://localhost/api/v1/health > /dev/null 2>&1; then
    echo "✅ 服务重启成功"
else
    echo "❌ 服务重启失败，需要手动检查"
fi
```

#### 2. 数据恢复
```bash
#!/bin/bash
# restore.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

echo "从备份恢复数据库: $BACKUP_FILE"

# 停止 API 服务
docker-compose -f docker-compose.prod.yml stop api

# 恢复数据库
gunzip -c $BACKUP_FILE | docker exec -i gamelink_mysql_prod mysql -u root -p$DB_ROOT_PASSWORD gamelink_prod

# 重启服务
docker-compose -f docker-compose.prod.yml start api

echo "数据库恢复完成"
```

---

## 🔒 安全配置

### 1. SSL/TLS 配置
```bash
# 生成自签名证书 (开发环境)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /opt/gamelink/ssl/key.pem \
    -out /opt/gamelink/ssl/cert.pem

# 使用 Let's Encrypt (生产环境)
certbot certonly --webroot -w /var/www/html -d your-domain.com
```

### 2. 防火墙配置
```bash
# UFW 防火墙规则
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 3. 安全头配置
```nginx
# Nginx 安全头配置
add_header X-Frame-Options DENY;
add_header X-Content-Type-Options nosniff;
add_header X-XSS-Protection "1; mode=block";
add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'";
```

### 4. 数据库安全
```sql
-- 创建专用数据库用户
CREATE USER 'gamelink_web'@'%' IDENTIFIED BY 'strong_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON gamelink.* TO 'gamelink_web'@'%';

-- 创建只读用户 (用于报表)
CREATE USER 'gamelink_readonly'@'%' IDENTIFIED BY 'readonly_password';
GRANT SELECT ON gamelink.* TO 'gamelink_readonly'@'%';

-- 禁用不必要的功能
FLUSH PRIVILEGES;
```

---

## ⚡ 性能优化

### 1. 数据库优化
```sql
-- 添加索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- 配置优化 (my.cnf)
[mysqld]
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
max_connections = 200
query_cache_size = 64M
```

### 2. Redis 优化
```conf
# redis.conf
maxmemory 1gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
```

### 3. Nginx 优化
```nginx
# nginx.conf 优化配置
worker_processes auto;
worker_connections 2048;

keepalive_timeout 30;
keepalive_requests 100;

# 开启 gzip
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css application/json application/javascript;

# 缓存配置
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

---

## 🔍 故障排查

### 常见问题诊断

#### 1. 服务无法启动
```bash
# 检查容器状态
docker-compose -f docker-compose.prod.yml ps

# 查看容器日志
docker-compose -f docker-compose.prod.yml logs api

# 检查端口占用
netstat -tlnp | grep :8080
```

#### 2. 数据库连接问题
```bash
# 检查数据库服务
docker exec gamelink_mysql_prod mysql -u root -p -e "SHOW PROCESSLIST;"

# 检查连接数
docker exec gamelink_mysql_prod mysql -u root -p -e "SHOW STATUS LIKE 'Threads_connected';"

# 测试连接
docker exec gamelink_mysql_prod mysql -u gamelink -p -e "SELECT 1;"
```

#### 3. 性能问题
```bash
# 检查系统资源
top
htop
iotop

# 检查磁盘使用
df -h
du -sh /opt/gamelink/*

# 检查网络连接
netstat -an | grep :8080
```

### 监控指标分析

#### 关键指标
- **CPU 使用率**: < 80%
- **内存使用率**: < 85%
- **磁盘使用率**: < 80%
- **网络延迟**: < 100ms
- **API 响应时间**: < 500ms (95th percentile)
- **错误率**: < 1%

#### 告警阈值
```yaml
# 告警配置示例
alerts:
  - name: high_cpu
    condition: cpu_usage > 80
    duration: 5m
    severity: warning

  - name: high_memory
    condition: memory_usage > 85
    duration: 5m
    severity: warning

  - name: api_error_rate
    condition: error_rate > 5
    duration: 2m
    severity: critical
```

---

## 📞 技术支持

### 联系方式
- **技术团队**: dev-team@gamelink.com
- **运维团队**: ops-team@gamelink.com
- **紧急联系**: +86-xxx-xxxx-xxxx

### 文档资源
- [API 文档](./API.md)
- [架构设计](./ARCHITECTURE.md)
- [开发指南](./DEVELOPMENT.md)
- [故障排查](./TROUBLESHOOTING.md)

---

*最后更新: 2025-11-13*