#!/bin/bash
# GameLink 生产环境密钥生成工具
# 用途：为生产环境生成所有必需的密钥和密码
# 使用方法：bash scripts/generate-production-keys.sh

set -e

echo "=================================="
echo " GameLink 生产环境密钥生成工具"
echo "=================================="
echo ""
echo "⚠️  警告：此脚本将生成生产环境密钥"
echo "请妥善保管生成的密钥，不要泄露给任何人！"
echo ""
read -p "是否继续？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "已取消操作"
    exit 0
fi

echo ""
echo "开始生成密钥..."
echo ""

# 生成数据库密码
echo "1. 生成数据库密码..."
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
echo "   数据库密码已生成（24字节）"

# 生成 Redis 密码
echo "2. 生成 Redis 密码..."
REDIS_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
echo "   Redis 密码已生成（24字节）"

# 生成 JWT 密钥
echo "3. 生成 JWT 密钥..."
JWT_SECRET_KEY=$(openssl rand -base64 32 | tr -d '\n')
echo "   JWT 密钥已生成（32字节）"

# 生成后端加密密钥（base64 格式）
echo "4. 生成后端加密密钥..."
BACKEND_CRYPTO_KEY=$(openssl rand -base64 32 | tr -d '\n')
echo "   后端 CRYPTO_SECRET_KEY 已生成（base64，32字节）"

# 生成后端 IV（base64 格式）
echo "5. 生成后端 IV..."
BACKEND_CRYPTO_IV=$(openssl rand -base64 16 | tr -d '\n')
echo "   后端 CRYPTO_IV 已生成（base64，16字节）"

# 生成前端加密密钥（原始字节格式）
echo "6. 生成前端加密密钥..."
FRONTEND_CRYPTO_KEY=$(openssl rand -base64 32 | base64 -d | xxd -p -c 32 | tr -d '\n')
echo "   前端 VITE_CRYPTO_SECRET_KEY 已生成（原始字节，32字节）"

# 生成前端 IV（原始字节格式）
echo "7. 生成前端 IV..."
FRONTEND_CRYPTO_IV=$(openssl rand -base64 16 | base64 -d | xxd -p -c 16 | tr -d '\n')
echo "   前端 VITE_CRYPTO_IV 已生成（原始字节，16字节）"

# 生成超级管理员密码
echo "8. 生成超级管理员密码..."
SUPER_ADMIN_PASSWORD=$(openssl rand -base64 24 | tr -d '\n' | head -c 16)
echo "   超级管理员密码已生成（16字符）"

echo ""
echo "=================================="
echo " 密钥生成完成！"
echo "=================================="
echo ""

# 创建 .env.production 文件
cat > .env.production << EOF
# GameLink 生产环境配置
# 生成时间: $(date)
# ⚠️  警告：请妥善保管此文件，不要提交到版本控制

# =============================================================================
# 应用环境
# =============================================================================
APP_ENV=production
GIN_MODE=release

# =============================================================================
# 数据库配置 (PostgreSQL)
# =============================================================================
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=${DB_PASSWORD}
POSTGRES_DB=gamelink
POSTGRES_PORT=5432

# 连接池配置
DB_MAX_CONNS=50
DB_MAX_IDLE=25

# =============================================================================
# 缓存配置 (Redis)
# =============================================================================
CACHE_TYPE=redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=${REDIS_PASSWORD}
REDIS_DB=0

# =============================================================================
# 安全: 加密配置
# =============================================================================
# ⚠️  生产环境必须启用加密
CRYPTO_ENABLED=true

# AES-256-CBC 加密密钥（base64 编码）
CRYPTO_SECRET_KEY=${BACKEND_CRYPTO_KEY}

# 初始化向量（base64 编码）
CRYPTO_IV=${BACKEND_CRYPTO_IV}

# 使用 SHA-256 签名
CRYPTO_USE_SIGNATURE=true

# =============================================================================
# 安全: JWT 认证
# =============================================================================
JWT_SECRET_KEY=${JWT_SECRET_KEY}
JWT_TOKEN_TTL_HOURS=24

# =============================================================================
# 超级管理员配置
# =============================================================================
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=${SUPER_ADMIN_PASSWORD}
SUPER_ADMIN_NAME=Super Admin

# =============================================================================
# 种子数据
# =============================================================================
# 生产环境不初始化种子数据
SEED_ENABLED=false

# =============================================================================
# 外部 API（根据需要配置）
# =============================================================================
# 微信支付
WECHAT_PAY_ENABLED=false

# 支付宝
ALIPAY_ENABLED=false

# 短信服务
SMS_ENABLED=false

# 对象存储
OSS_ENABLED=false
EOF

echo "✅ 已创建 .env.production 文件"
echo ""

# 创建前端 .env.production 文件
mkdir -p admin
cat > admin/.env.production << EOF
# GameLink Admin Panel - 生产环境配置
# 生成时间: $(date)

# =============================================================================
# API Configuration
# =============================================================================
VITE_API_BASE_URL=/api/v1

# =============================================================================
# Security: Encryption Configuration
# =============================================================================
# ⚠️  必须与后端加密配置保持一致
VITE_CRYPTO_ENABLED=true
VITE_CRYPTO_SECRET_KEY=${FRONTEND_CRYPTO_KEY}
VITE_CRYPTO_IV=${FRONTEND_CRYPTO_IV}
VITE_CRYPTO_USE_SIGNATURE=true

# =============================================================================
# Application Configuration
# =============================================================================
VITE_APP_TITLE=GameLink Admin Panel
VITE_APP_VERSION=1.0.0
VITE_DEBUG=false

# =============================================================================
# Feature Flags
# =============================================================================
VITE_ENABLE_WEBSOCKET=true
VITE_WEBSOCKET_RECONNECT_ATTEMPTS=5
VITE_WEBSOCKET_RECONNECT_INTERVAL=3000

# =============================================================================
# UI Configuration
# =============================================================================
VITE_DEFAULT_PAGE_SIZE=20
VITE_MAX_PAGE_SIZE=100
VITE_DATE_FORMAT=YYYY-MM-DD HH:mm:ss
VITE_TIMEZONE=Asia/Shanghai
EOF

echo "✅ 已创建 admin/.env.production 文件"
echo ""

# 创建密钥备份文件
BACKUP_DIR="backups/keys/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

cat > "$BACKUP_DIR/keys.txt" << EOF
# GameLink 生产环境密钥备份
# 生成时间: $(date)
# ⚠️  请妥善保管此文件，建议加密存储

# 数据库密码
DB_PASSWORD=${DB_PASSWORD}

# Redis 密码
REDIS_PASSWORD=${REDIS_PASSWORD}

# JWT 密钥
JWT_SECRET_KEY=${JWT_SECRET_KEY}

# 超级管理员密码
SUPER_ADMIN_PASSWORD=${SUPER_ADMIN_PASSWORD}

# 后端加密密钥（base64）
BACKEND_CRYPTO_KEY=${BACKEND_CRYPTO_KEY}
BACKEND_CRYPTO_IV=${BACKEND_CRYPTO_IV}

# 前端加密密钥（原始字节）
FRONTEND_CRYPTO_KEY=${FRONTEND_CRYPTO_KEY}
FRONTEND_CRYPTO_IV=${FRONTEND_CRYPTO_IV}
EOF

echo "✅ 已创建密钥备份: $BACKUP_DIR/keys.txt"
echo ""

# 显示密钥摘要
echo "=================================="
echo " 密钥摘要（请妥善保管）"
echo "=================================="
echo ""
echo "数据库密码: ${DB_PASSWORD:0:8}..."
echo "Redis 密码: ${REDIS_PASSWORD:0:8}..."
echo "JWT 密钥: ${JWT_SECRET_KEY:0:8}..."
echo "后端加密密钥: ${BACKEND_CRYPTO_KEY:0:8}..."
echo "前端加密密钥: ${FRONTEND_CRYPTO_KEY:0:8}..."
echo "管理员密码: ${SUPER_ADMIN_PASSWORD:0:8}..."
echo ""

# 安全提示
echo "=================================="
echo " ⚠️  安全提示"
echo "=================================="
echo ""
echo "1. 请立即将密钥备份文件复制到安全位置"
echo "2. 建议使用密码管理器或密钥管理服务存储"
echo "3. 不要将 .env.production 文件提交到 Git"
echo "4. 定期轮换密钥（建议每季度）"
echo "5. 如果密钥泄露，请立即重新生成并更新"
echo ""

echo "=================================="
echo " 下一步操作"
echo "=================================="
echo ""
echo "1. 查看生成的配置文件："
echo "   cat .env.production"
echo ""
echo "2. 查看前端配置文件："
echo "   cat admin/.env.production"
echo ""
echo "3. 验证密钥一致性："
echo "   ./scripts/verify-crypto-keys.sh"
echo ""
echo "4. 部署到 Staging 环境："
echo "   docker-compose -f docker-compose.staging.yml --env-file .env.staging up -d"
echo ""
echo "5. 验证部署："
echo "   ./scripts/verify-deployment.sh"
echo ""

echo "✅ 密钥生成完成！"
echo ""
