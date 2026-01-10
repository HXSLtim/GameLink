#!/bin/bash
# GameLink 安全密钥生成工具
# 用于生成生产环境所需的各类安全密钥

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
WHITE='\033[0;37m'
NC='\033[0m' # No Color

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  GameLink 安全密钥生成工具${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 检查 OpenSSL 是否安装
if ! command -v openssl &> /dev/null; then
    echo -e "${RED}错误: 未找到 OpenSSL 命令${NC}"
    echo ""
    echo -e "${YELLOW}请安装 OpenSSL:${NC}"
    echo -e "${WHITE}  Ubuntu/Debian: sudo apt-get install openssl${NC}"
    echo -e "${WHITE}  CentOS/RHEL: sudo yum install openssl${NC}"
    echo -e "${WHITE}  macOS: brew install openssl${NC}"
    echo ""
    exit 1
fi

echo -e "${GREEN}生成的密钥如下（请妥善保存，不要泄露）：${NC}"
echo ""

# 生成 32 字节加密密钥 (AES-256-CBC)
SECRET_KEY=$(openssl rand -base64 32)
echo -e "${YELLOW}1. 加密密钥 (CRYPTO_SECRET_KEY) - 32字节:${NC}"
echo -e "   ${WHITE}$SECRET_KEY${NC}"
echo ""

# 生成 16 字节初始化向量
IV=$(openssl rand -base64 16)
echo -e "${YELLOW}2. 初始化向量 (CRYPTO_IV) - 16字节:${NC}"
echo -e "   ${WHITE}$IV${NC}"
echo ""

# 生成 32 字节 JWT 密钥
JWT_SECRET=$(openssl rand -base64 32)
echo -e "${YELLOW}3. JWT 密钥 (JWT_SECRET_KEY) - 32字节:${NC}"
echo -e "   ${WHITE}$JWT_SECRET${NC}"
echo ""

# 生成 24 字节超级管理员密码
ADMIN_PASSWORD=$(openssl rand -base64 24)
echo -e "${YELLOW}4. 超级管理员密码 (SUPER_ADMIN_PASSWORD) - 24字节:${NC}"
echo -e "   ${WHITE}$ADMIN_PASSWORD${NC}"
echo ""

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}使用方法：${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""
echo -e "${YELLOW}Linux/Mac 环境变量设置:${NC}"
echo -e "${WHITE}export CRYPTO_SECRET_KEY='$SECRET_KEY'${NC}"
echo -e "${WHITE}export CRYPTO_IV='$IV'${NC}"
echo -e "${WHITE}export JWT_SECRET_KEY='$JWT_SECRET'${NC}"
echo -e "${WHITE}export SUPER_ADMIN_PASSWORD='$ADMIN_PASSWORD'${NC}"
echo ""

echo -e "${YELLOW}Docker Compose 环境变量 (.env 文件):${NC}"
echo -e "${WHITE}CRYPTO_SECRET_KEY=$SECRET_KEY${NC}"
echo -e "${WHITE}CRYPTO_IV=$IV${NC}"
echo -e "${WHITE}JWT_SECRET_KEY=$JWT_SECRET${NC}"
echo -e "${WHITE}SUPER_ADMIN_PASSWORD=$ADMIN_PASSWORD${NC}"
echo ""

echo -e "${CYAN}========================================${NC}"
echo -e "${RED}重要提示：${NC}"
echo -e "${CYAN}========================================${NC}"
echo -e "${WHITE}1. 请将这些密钥保存到安全的位置（密码管理器）${NC}"
echo -e "${WHITE}2. 不要将密钥提交到 Git 仓库${NC}"
echo -e "${WHITE}3. 生产环境每次部署都应使用不同的密钥${NC}"
echo -e "${WHITE}4. 密钥泄露后立即重新生成并更新${NC}"
echo ""

# 可选：导出到 .env 文件
read -p "是否导出到 .env 文件? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    cat > .env <<EOF
# GameLink 生产环境密钥
# 警告: 请勿将此文件提交到版本控制系统

# 加密配置 (AES-256-CBC)
CRYPTO_ENABLED=true
CRYPTO_SECRET_KEY=$SECRET_KEY
CRYPTO_IV=$IV

# JWT 配置
JWT_SECRET_KEY=$JWT_SECRET

# 超级管理员配置
SUPER_ADMIN_EMAIL=admin@gamelink.com
SUPER_ADMIN_PASSWORD=$ADMIN_PASSWORD
SUPER_ADMIN_NAME=Super Admin

# 数据库配置 (请根据实际情况修改)
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=your_secure_db_password_here
POSTGRES_DB=gamelink

# Redis 配置
REDIS_PASSWORD=your_secure_redis_password_here
EOF

    echo ""
    echo -e "${GREEN}已导出到 .env${NC}"
    echo -e "${YELLOW}请编辑文件并填写数据库和 Redis 密码${NC}"
    echo -e "${YELLOW}然后将 .env 添加到 .gitignore${NC}"
fi
