#!/bin/bash
# GameLink 服务器部署脚本（1Panel 环境）
# 前提：已在 1Panel 应用商店安装 PostgreSQL 和 Redis
# 在服务器上运行: bash deploy.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "🚀 GameLink 部署开始..."
echo "   项目目录: $PROJECT_DIR"

# 检查配置文件
if [ ! -f ".env.production" ]; then
    echo "❌ 未找到 .env.production 配置文件"
    echo "   请复制 deploy/.env.production.template 为 .env.production 并填入配置"
    echo ""
    echo "   cp deploy/.env.production.template .env.production"
    echo "   nano .env.production"
    exit 1
fi

# 加载配置
source .env.production

# 检查必要配置
check_config() {
    echo "📋 检查配置..."
    local missing=0
    
    [ -z "$POSTGRES_PASSWORD" ] && echo "   ❌ 缺少 POSTGRES_PASSWORD" && missing=1
    [ -z "$REDIS_PASSWORD" ] && echo "   ❌ 缺少 REDIS_PASSWORD" && missing=1
    [ -z "$SUPER_ADMIN_PASSWORD" ] && echo "   ❌ 缺少 SUPER_ADMIN_PASSWORD" && missing=1
    
    if [ $missing -eq 1 ]; then
        echo "   请编辑 .env.production 填入必要配置"
        exit 1
    fi
    echo "   ✅ 配置检查通过"
}

# 生成密钥
generate_keys() {
    echo "🔐 检查/生成加密密钥..."
    local updated=0
    
    if [ -z "$CRYPTO_SECRET_KEY" ] || [ "$CRYPTO_SECRET_KEY" == "你的32字节加密密钥" ]; then
        NEW_SECRET=$(openssl rand -base64 32 | tr -d '\n' | head -c 32)
        sed -i "s|CRYPTO_SECRET_KEY=.*|CRYPTO_SECRET_KEY=$NEW_SECRET|" .env.production
        echo "   ✅ 已生成 CRYPTO_SECRET_KEY"
        updated=1
    fi
    
    if [ -z "$CRYPTO_IV" ] || [ "$CRYPTO_IV" == "你的16字节IV" ]; then
        NEW_IV=$(openssl rand -base64 16 | tr -d '\n' | head -c 16)
        sed -i "s|CRYPTO_IV=.*|CRYPTO_IV=$NEW_IV|" .env.production
        echo "   ✅ 已生成 CRYPTO_IV"
        updated=1
    fi
    
    if [ -z "$JWT_SECRET_KEY" ] || [ "$JWT_SECRET_KEY" == "你的JWT密钥至少32个字符" ]; then
        NEW_JWT=$(openssl rand -base64 32 | tr -d '\n')
        sed -i "s|JWT_SECRET_KEY=.*|JWT_SECRET_KEY=$NEW_JWT|" .env.production
        echo "   ✅ 已生成 JWT_SECRET_KEY"
        updated=1
    fi
    
    # 重新加载配置
    [ $updated -eq 1 ] && source .env.production
}

# 测试数据库连接
test_db_connection() {
    echo "🔍 测试数据库连接..."
    
    # 测试 PostgreSQL
    if command -v psql &> /dev/null; then
        if PGPASSWORD="$POSTGRES_PASSWORD" psql -h "${DB_HOST:-127.0.0.1}" -p "${DB_PORT:-5432}" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1" &> /dev/null; then
            echo "   ✅ PostgreSQL 连接正常"
        else
            echo "   ⚠️  PostgreSQL 连接失败，请检查配置"
        fi
    else
        echo "   ⏭️  跳过 PostgreSQL 测试（psql 未安装）"
    fi
    
    # 测试 Redis
    if command -v redis-cli &> /dev/null; then
        if redis-cli -h "${REDIS_HOST:-127.0.0.1}" -p "${REDIS_PORT:-6379}" -a "$REDIS_PASSWORD" ping &> /dev/null; then
            echo "   ✅ Redis 连接正常"
        else
            echo "   ⚠️  Redis 连接失败，请检查配置"
        fi
    else
        echo "   ⏭️  跳过 Redis 测试（redis-cli 未安装）"
    fi
}

# 构建管理后台前端
build_admin() {
    echo "📦 构建管理后台前端..."
    
    # 创建前端环境变量
    cat > admin/.env.production << EOF
VITE_API_BASE_URL=/api/v1
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=$CRYPTO_SECRET_KEY
VITE_CRYPTO_IV=$CRYPTO_IV
VITE_CRYPTO_USE_SIGNATURE=true
EOF
    
    cd admin
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        echo "   ⚠️  Node.js 未安装"
        echo "   请安装 Node.js 20+，或在本地构建后上传 dist 目录"
        cd ..
        return 1
    fi
    
    echo "   安装依赖..."
    npm install --legacy-peer-deps 2>/dev/null || npm install
    
    echo "   构建..."
    npm run build
    
    cd ..
    echo "   ✅ 前端构建完成: admin/dist/"
}

# 启动后端服务
start_backend() {
    echo "🐳 启动后端服务..."
    
    # 停止旧容器
    docker compose -f deploy/docker-compose.server.yml --env-file .env.production down 2>/dev/null || true
    
    # 构建并启动
    docker compose -f deploy/docker-compose.server.yml --env-file .env.production up -d --build
    
    echo "   ⏳ 等待服务启动..."
    sleep 10
    
    # 健康检查
    for i in {1..6}; do
        if curl -s http://127.0.0.1:8080/api/v1/healthz > /dev/null 2>&1; then
            echo "   ✅ 后端服务正常"
            return 0
        fi
        echo "   等待中... ($i/6)"
        sleep 5
    done
    
    echo "   ⚠️  后端服务可能还在启动中"
    echo "   查看日志: docker logs gamelink-backend"
}

# 显示状态
show_status() {
    echo ""
    echo "═══════════════════════════════════════════════════════"
    echo "📊 服务状态:"
    docker ps --filter "name=gamelink" --format "table {{.Names}}\t{{.Status}}"
    
    echo ""
    echo "🎉 部署完成！"
    echo ""
    echo "📝 在 1Panel 中配置网站:"
    echo "   1. 网站 → 创建网站 → 静态网站"
    echo "   2. 根目录: $PROJECT_DIR/admin/dist"
    echo "   3. 添加反向代理:"
    echo "      - 代理目录: /api"
    echo "      - 目标URL: http://127.0.0.1:8080"
    echo ""
    echo "🔑 管理员登录:"
    echo "   邮箱: $SUPER_ADMIN_EMAIL"
    echo "   密码: (查看 .env.production)"
    echo ""
    echo "🔧 常用命令:"
    echo "   查看日志: docker logs -f gamelink-backend"
    echo "   重启服务: docker restart gamelink-backend"
    echo "   停止服务: docker stop gamelink-backend"
    echo "═══════════════════════════════════════════════════════"
}

# 主流程
check_config
generate_keys
test_db_connection
build_admin || echo "   ⚠️  前端构建跳过，请手动构建或上传"
start_backend
show_status
