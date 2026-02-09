#!/bin/bash
# GameLink 支付集成环境准备脚本
# 用途：为支付集成准备所需的证书和配置
# 使用方法：bash scripts/setup-payment.sh

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_header() {
    echo ""
    echo "=================================="
    echo " GameLink 支付集成环境准备工具"
    echo "=================================="
    echo ""
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v openssl &> /dev/null; then
        log_error "OpenSSL 未安装，请先安装 OpenSSL"
        exit 1
    fi

    log_success "依赖检查通过"
}

# 创建证书目录
create_cert_directories() {
    log_info "创建证书目录..."

    mkdir -p certs/wechat
    mkdir -p certs/alipay

    log_success "证书目录创建完成"
}

# 生成测试用 RSA 密钥对（支付宝）
generate_alipay_test_keys() {
    log_info "生成支付宝测试用 RSA 密钥对..."

    # 生成应用私钥
    openssl genrsa -out certs/alipay/app_private_key.pem 2048

    # 提取公钥
    openssl rsa -in certs/alipay/app_private_key.pem -pubout -out certs/alipay/app_public_key.pem

    # 验证密钥
    if [ -f "certs/alipay/app_private_key.pem" ] && [ -f "certs/alipay/app_public_key.pem" ]; then
        log_success "支付宝测试密钥生成完成"
    else
        log_error "支付宝测试密钥生成失败"
        exit 1
    fi
}

# 显示配置指南
show_config_guide() {
    echo ""
    echo "=================================="
    echo " 支付集成配置指南"
    echo "=================================="
    echo ""
    echo "1. 微信支付配置"
    echo "   1.1 登录微信商户平台"
    echo "   1.2 下载商户证书（apiclient_cert.p12）"
    echo "   1.3 设置 API 密钥"
    echo "   1.4 配置回调地址："
    echo "       生产: https://api.gamelink.com/api/v1/payments/wechat/notify"
    echo "       测试: https://staging.gamelink.com/api/v1/payments/wechat/notify"
    echo "   1.5 将证书文件复制到：certs/wechat/"
    echo ""
    echo "2. 支付宝配置"
    echo "   2.1 登录支付宝开放平台"
    echo "   2.2 创建应用并获取 APP_ID"
    echo "   2.3 生成 RSA 密钥对（或使用工具生成）"
    echo "   2.4 上传应用公钥到支付宝"
    echo "   2.5 下载支付宝公钥"
    echo "   2.6 配置回调地址："
    echo "       生产: https://api.gamelink.com/api/v1/payments/alipay/notify"
    echo "       测试: https://staging.gamelink.com/api/v1/payments/alipay/notify"
    echo "   2.7 将密钥文件复制到：certs/alipay/"
    echo ""
    echo "3. 环境变量配置"
    echo "   3.1 更新 .env.production 文件："
    echo "       cp .env.example .env.production"
    echo "       nano .env.production"
    echo ""
    echo "   3.2 填入支付相关配置："
    echo "       WECHAT_PAY_ENABLED=true"
    echo "       WECHAT_PAY_APP_ID=<你的AppID>"
    echo "       WECHAT_PAY_MCH_ID=<你的商户号>"
    echo "       WECHAT_PAY_API_KEY=<你的API密钥>"
    echo "       WECHAT_PAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/wechat/notify"
    echo ""
    echo "       ALIPAY_ENABLED=true"
    echo "       ALIPAY_APP_ID=<你的应用ID>"
    echo "       ALIPAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/alipay/notify"
    echo ""
    echo "4. 域名和 HTTPS 配置"
    echo "   4.1 确保域名已备案"
    echo "   4.2 申请 SSL/TLS 证书"
    echo "   4.3 配置 Nginx SSL 证书"
    echo ""
    echo "5. Docker 卷挂载配置"
    echo "   5.1 确保 docker-compose.yml 包含："
    echo "       volumes:"
    echo "         - ./certs/wechat:/app/certs/wechat:ro"
    echo "         - ./certs/alipay:/app/certs/alipay:ro"
    echo ""
}

# 创建环境变量模板
create_env_template() {
    log_info "创建支付配置模板..."

    cat > .env.payment.example << EOF
# GameLink 支付集成配置模板
# 复制此文件为 .env.payment 并填入实际值

# =============================================================================
# 微信支付配置
# =============================================================================
WECHAT_PAY_ENABLED=false

# 微信公众号/小程序 AppID
WECHAT_PAY_APP_ID=wx####################

# 微信商户号
WECHAT_PAY_MCH_ID=1####################

# API 密钥（32字节，在微信商户平台设置）
WECHAT_PAY_API_KEY=##############################

# 商户证书路径（Docker 容器内路径）
WECHAT_PAY_CERT_PATH=/app/certs/wechat/apiclient_cert.p12
WECHAT_PAY_KEY_PATH=/app/certs/wechat/apiclient_key.pem

# 证书序列号
WECHAT_PAY_CERT_SERIAL_NO=######################

# 回调地址（必须是 HTTPS）
WECHAT_PAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/wechat/notify

# =============================================================================
# 支付宝配置
# =============================================================================
ALIPAY_ENABLED=false

# 支付宝应用 ID
ALIPAY_APP_ID=####################

# 支付宝网关地址
ALIPAY_GATEWAY=https://openapi.alipay.com/gateway.do

# 应用私钥路径
ALIPAY_PRIVATE_KEY_PATH=/app/certs/alipay/app_private_key.pem

# 支付宝公钥路径
ALIPAY_PUBLIC_KEY_PATH=/app/certs/alipay/app_public_key.pem

# 回调地址（必须是 HTTPS）
ALIPAY_NOTIFY_URL=https://api.gamelink.com/api/v1/payments/alipay/notify

# =============================================================================
# 回调 URL 说明
# =============================================================================
#
# 微信支付和支付宝都要求回调地址必须是 HTTPS
#
# 生产环境回调地址：
# - 微信: https://api.gamelink.com/api/v1/payments/wechat/notify
# - 支付宝: https://api.gamelink.com/api/v1/payments/alipay/notify
#
# Staging 环境回调地址：
# - 微信: https://staging.gamelink.com/api/v1/payments/wechat/notify
# - 支付宝: https://staging.gamelink.com/api/v1/payments/alipay/notify
#
# 开发环境：
# - 可以使用 ngrok 等内网穿透工具获取 HTTPS 地址
# - 或者暂时禁用真实支付，使用 Mock 模式
#
# =============================================================================
EOF

    log_success "支付配置模板创建完成"
}

# 主流程
main() {
    show_header

    # 检查环境
    check_dependencies

    # 创建目录
    create_cert_directories

    # 生成测试密钥
    log_info "生成测试用密钥..."
    generate_alipay_test_keys

    # 创建环境模板
    create_env_template

    # 显示配置指南
    show_config_guide

    echo ""
    echo "=================================="
    echo " 下一步操作"
    echo "=================================="
    echo ""
    echo "1. 查看支付配置文档："
    echo "   cat docs/PAYMENT_WEBHOOK_CONFIG.md"
    echo ""
    echo "2. 按照上述配置指南完成支付平台配置"
    echo ""
    echo "3. 生成生产环境配置："
    echo "   cp .env.payment.example .env.payment"
    echo "   nano .env.payment"
    echo ""
    echo "4. 将实际证书文件复制到对应目录："
    echo "   微信: certs/wechat/"
    echo "   支付宝: certs/alipay/"
    echo ""
    echo "5. 测试配置："
    echo "   bash scripts/verify-crypto-keys.sh"
    echo ""

    log_success "支付集成环境准备完成！"
}

# 运行主流程
main "$@"
