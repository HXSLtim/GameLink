#!/bin/bash
# GameLink 部署前检查脚本
# 在部署到 Staging 或 Production 之前验证所有配置

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 统计变量
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNING_CHECKS=0
BACKEND_ENV_FILE=".env.production"
ADMIN_ENV_FILE="admin/.env.production"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
    WARNING_CHECKS=$((WARNING_CHECKS + 1))
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
}

show_header() {
    echo ""
    echo "=================================="
    echo " GameLink 部署前检查工具"
    echo "=================================="
    echo ""
}

resolve_env_file() {
    local target_var="$1"
    local primary="$2"
    local fallback="$3"

    if [ -f "$primary" ]; then
        log_success "$primary 文件存在"
        printf -v "$target_var" "%s" "$primary"
        return 0
    fi

    if [ -f "$fallback" ]; then
        log_warning "$primary 文件不存在，使用 $fallback 进行预检查（仅示例配置）"
        printf -v "$target_var" "%s" "$fallback"
        return 0
    fi

    log_error "$primary 文件不存在，且未找到备用文件 $fallback"
    printf -v "$target_var" "%s" ""
    return 1
}

# 检查环境变量文件
check_env_files() {
    log_info "检查环境变量文件..."

    local env_missing=false

    # 检查并解析后端环境文件
    resolve_env_file BACKEND_ENV_FILE ".env.production" ".env.example"
    if [ -n "$BACKEND_ENV_FILE" ]; then
        # 检查必需的变量
        if grep -q "^CRYPTO_ENABLED=true" "$BACKEND_ENV_FILE"; then
            log_success "  生产环境加密已启用"
        else
            log_warning "  生产环境加密未启用（CRYPTO_ENABLED != true）"
        fi

        if grep -q "^SEED_ENABLED=false" "$BACKEND_ENV_FILE"; then
            log_success "  生产环境种子数据已禁用"
        else
            log_warning "  生产环境种子数据未禁用（SEED_ENABLED != false）"
        fi
    else
        env_missing=true
    fi

    # 检查并解析管理后台环境文件
    resolve_env_file ADMIN_ENV_FILE "admin/.env.production" "admin/.env.example"
    if [ -z "$ADMIN_ENV_FILE" ]; then
        env_missing=true
    fi

    if [ "$env_missing" = true ]; then
        log_warning "建议运行: bash scripts/generate-production-keys.sh"
    fi
}

# 检查加密配置一致性
check_crypto_consistency() {
    log_info "检查前后端加密配置一致性..."

    if [ -z "$BACKEND_ENV_FILE" ] || [ -z "$ADMIN_ENV_FILE" ]; then
        log_warning "跳过加密配置检查（缺少 .env.production 文件）"
        return
    fi

    # 提取后端加密配置
    local backend_enabled=$(grep "^CRYPTO_ENABLED=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
    local backend_key=$(grep "^CRYPTO_SECRET_KEY=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
    local backend_iv=$(grep "^CRYPTO_IV=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
    local backend_sig=$(grep "^CRYPTO_USE_SIGNATURE=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)

    # 提取前端加密配置
    local frontend_enabled=$(grep "^VITE_CRYPTO_ENABLED=" "$ADMIN_ENV_FILE" | cut -d'=' -f2)
    local frontend_key=$(grep "^VITE_CRYPTO_SECRET_KEY=" "$ADMIN_ENV_FILE" | cut -d'=' -f2)
    local frontend_iv=$(grep "^VITE_CRYPTO_IV=" "$ADMIN_ENV_FILE" | cut -d'=' -f2)
    local frontend_sig=$(grep "^VITE_CRYPTO_USE_SIGNATURE=" "$ADMIN_ENV_FILE" | cut -d'=' -f2)

    # 检查启用状态
    if [ "$backend_enabled" = "true" ] && [ "$frontend_enabled" = "true" ]; then
        log_success "加密启用状态一致：前后端均为 true"
    elif [ "$backend_enabled" = "false" ] && [ "$frontend_enabled" = "false" ]; then
        log_success "加密启用状态一致：前后端均为 false"
        log_info "前后端加密均关闭，跳过签名和密钥长度检查"
        return
    else
        log_error "加密启用状态不一致（后端: $backend_enabled, 前端: $frontend_enabled）"
        return
    fi

    # 仅在启用加密时检查签名与密钥
    if [ "$backend_sig" = "$frontend_sig" ]; then
        log_success "签名配置一致：$backend_sig"
    else
        log_error "签名配置不一致（后端: $backend_sig, 前端: $frontend_sig）"
    fi

    # 检查密钥长度（后端 base64）
    if [ ${#backend_key} -eq 44 ]; then
        log_success "后端密钥长度正确（44字符，base64编码的32字节）"
    else
        log_warning "后端密钥长度异常：${#backend_key} 字符（期望44字符）"
    fi

    # 检查IV长度（后端 base64）
    if [ ${#backend_iv} -eq 24 ]; then
        log_success "后端IV长度正确（24字符，base64编码的16字节）"
    else
        log_warning "后端IV长度异常：${#backend_iv} 字符（期望24字符）"
    fi

    # 检查密钥长度（前端原始字节）
    if [ ${#frontend_key} -eq 32 ]; then
        log_success "前端密钥长度正确（32字符，原始字节）"
    else
        log_warning "前端密钥长度异常：${#frontend_key} 字符（期望32字符）"
    fi

    # 检查IV长度（前端原始字节）
    if [ ${#frontend_iv} -eq 16 ]; then
        log_success "前端IV长度正确（16字符，原始字节）"
    else
        log_warning "前端IV长度异常：${#frontend_iv} 字符（期望16字符）"
    fi
}

# 检查支付配置
check_payment_config() {
    log_info "检查支付配置..."

    if [ -z "$BACKEND_ENV_FILE" ]; then
        log_warning "跳过支付配置检查（缺少 .env.production 文件）"
        return
    fi

    # 检查微信支付
    local wechat_enabled=$(grep "^WECHAT_PAY_ENABLED=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
    if [ "$wechat_enabled" = "true" ]; then
        log_info "微信支付已启用"

        local wechat_app_id=$(grep "^WECHAT_PAY_APP_ID=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
        local wechat_mch_id=$(grep "^WECHAT_PAY_MCH_ID=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
        local wechat_api_key=$(grep "^WECHAT_PAY_API_KEY=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)

        if [ -n "$wechat_app_id" ] && [ ${#wechat_app_id} -gt 10 ]; then
            log_success "  微信AppID已配置"
        else
            log_warning "  微信AppID未配置或格式异常"
        fi

        if [ -n "$wechat_mch_id" ] && [ ${#wechat_mch_id} -gt 5 ]; then
            log_success "  微信商户号已配置"
        else
            log_warning "  微信商户号未配置或格式异常"
        fi

        if [ -n "$wechat_api_key" ] && [ ${#wechat_api_key} -eq 32 ]; then
            log_success "  微信API密钥已配置"
        else
            log_warning "  微信API密钥未配置或长度异常（应为32字符）"
        fi

        # 检查证书文件
        if [ -f "certs/wechat/apiclient_cert.p12" ]; then
            log_success "  微信商户证书存在"
        else
            log_error "  微信商户证书缺失：certs/wechat/apiclient_cert.p12"
        fi

        if [ -f "certs/wechat/apiclient_key.pem" ]; then
            log_success "  微信商户私钥存在"
        else
            log_warning "  微信商户私钥缺失：certs/wechat/apiclient_key.pem"
        fi
    else
        log_info "微信支付未启用"
    fi

    # 检查支付宝
    local alipay_enabled=$(grep "^ALIPAY_ENABLED=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)
    if [ "$alipay_enabled" = "true" ]; then
        log_info "支付宝已启用"

        local alipay_app_id=$(grep "^ALIPAY_APP_ID=" "$BACKEND_ENV_FILE" | cut -d'=' -f2)

        if [ -n "$alipay_app_id" ] && [ ${#alipay_app_id} -gt 10 ]; then
            log_success "  支付宝应用ID已配置"
        else
            log_warning "  支付宝应用ID未配置或格式异常"
        fi

        # 检查密钥文件
        if [ -f "certs/alipay/app_private_key.pem" ]; then
            log_success "  支付宝应用私钥存在"
        else
            log_error "  支付宝应用私钥缺失：certs/alipay/app_private_key.pem"
        fi

        if [ -f "certs/alipay/alipay_public_key.pem" ]; then
            log_success "  支付宝公钥存在"
        else
            log_warning "  支付宝公钥缺失：certs/alipay/alipay_public_key.pem"
        fi
    else
        log_info "支付宝未启用"
    fi
}

# 检查 Docker 配置
check_docker_config() {
    log_info "检查 Docker 配置..."

    # 检查 docker-compose 文件
    if [ -f "docker-compose.prod.yml" ]; then
        log_success "docker-compose.prod.yml 存在"

        # 检查是否包含支付证书挂载
        if grep -q "./certs/wechat:/app/certs/wechat:ro" docker-compose.prod.yml; then
            log_success "  微信支付证书挂载已配置"
        else
            log_warning "  微信支付证书挂载未配置"
        fi

        if grep -q "./certs/alipay:/app/certs/alipay:ro" docker-compose.prod.yml; then
            log_success "  支付宝证书挂载已配置"
        else
            log_warning "  支付宝证书挂载未配置"
        fi
    else
        log_error "docker-compose.prod.yml 不存在"
    fi

    if [ -f "docker-compose.staging.yml" ]; then
        log_success "docker-compose.staging.yml 存在"
    else
        log_warning "docker-compose.staging.yml 不存在"
    fi
}

# 检查证书目录权限
check_cert_permissions() {
    log_info "检查证书目录权限..."

    if [ ! -d "certs" ]; then
        log_warning "certs 目录不存在"
        return
    fi

    # 检查微信证书目录
    if [ -d "certs/wechat" ]; then
        log_info "微信证书目录存在"

        # 检查证书文件权限（应该是 400 或 600）
        if [ -f "certs/wechat/apiclient_cert.p12" ]; then
            local perms=$(stat -c %a certs/wechat/apiclient_cert.p12 2>/dev/null || stat -f %A certs/wechat/apiclient_cert.p12 2>/dev/null)
            if [ "$perms" = "400" ] || [ "$perms" = "600" ]; then
                log_success "  微信证书权限正确 ($perms)"
            else
                log_warning "  微信证书权限建议设为 400 或 600（当前: $perms）"
            fi
        fi
    else
        log_info "微信证书目录不存在"
    fi

    # 检查支付宝证书目录
    if [ -d "certs/alipay" ]; then
        log_info "支付宝证书目录存在"

        if [ -f "certs/alipay/app_private_key.pem" ]; then
            local perms=$(stat -c %a certs/alipay/app_private_key.pem 2>/dev/null || stat -f %A certs/alipay/app_private_key.pem 2>/dev/null)
            if [ "$perms" = "400" ] || [ "$perms" = "600" ]; then
                log_success "  支付宝私钥权限正确 ($perms)"
            else
                log_warning "  支付宝私钥权限建议设为 400 或 600（当前: $perms）"
            fi
        fi
    else
        log_info "支付宝证书目录不存在"
    fi
}

# 检查 Nginx 配置
check_nginx_config() {
    log_info "检查 Nginx 配置..."

    if [ -f "deploy/nginx-production.conf" ]; then
        log_success "Nginx 生产配置文件存在"

        # 检查是否包含支付回调路由
        if grep -q "/api/" deploy/nginx-production.conf; then
            log_success "  API 路由已配置"
        else
            log_warning "  API 路由未找到"
        fi

        # 检查 SSL 配置
        if grep -q "ssl_certificate" deploy/nginx-production.conf; then
            log_success "  SSL 证书配置存在"
        else
            log_warning "  SSL 证书配置未找到（需要手动配置）"
        fi
    else
        log_warning "Nginx 生产配置文件不存在"
    fi
}

# 检查部署脚本
check_deployment_scripts() {
    log_info "检查部署脚本..."

    local scripts=(
        "scripts/deploy.sh"
        "scripts/verify-deployment.sh"
        "scripts/verify-crypto-keys.sh"
        "scripts/setup-payment.sh"
        "scripts/generate-production-keys.sh"
    )

    for script in "${scripts[@]}"; do
        if [ -f "$script" ]; then
            if [ -x "$script" ]; then
                log_success "$script 存在且可执行"
            else
                log_warning "$script 存在但不可执行（运行: chmod +x $script）"
            fi
        else
            log_warning "$script 不存在"
        fi
    done
}

# 检查文档完整性
check_documentation() {
    log_info "检查文档完整性..."

    local docs=(
        "docs/DEPLOYMENT_CHECKLIST.md"
        "docs/SECURITY_HARDENING.md"
        "docs/PAYMENT_WEBHOOK_CONFIG.md"
        "docs/ENCRYPTION_VERIFICATION_REPORT.md"
        "docs/DEPENDENCIES_AND_CONFIG.md"
    )

    for doc in "${docs[@]}"; do
        if [ -f "$doc" ]; then
            log_success "$doc 存在"
        else
            log_warning "$doc 不存在"
        fi
    done
}

# 显示总结报告
show_summary() {
    echo ""
    echo "=================================="
    echo " 检查总结"
    echo "=================================="
    echo ""
    echo -e "总检查项: ${BLUE}$TOTAL_CHECKS${NC}"
    echo -e "${GREEN}通过: $PASSED_CHECKS${NC}"
    echo -e "${YELLOW}警告: $WARNING_CHECKS${NC}"
    echo -e "${RED}失败: $FAILED_CHECKS${NC}"
    echo ""

    if [ $FAILED_CHECKS -eq 0 ] && [ $WARNING_CHECKS -eq 0 ]; then
        echo -e "${GREEN}✓ 所有检查通过，可以开始部署！${NC}"
        return 0
    elif [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${YELLOW}⚠ 有 $WARNING_CHECKS 个警告，建议修复后再部署${NC}"
        return 0
    else
        echo -e "${RED}✗ 有 $FAILED_CHECKS 个错误必须修复才能部署${NC}"
        return 1
    fi
}

# 主流程
main() {
    show_header

    # 执行所有检查
    check_env_files
    check_crypto_consistency
    check_payment_config
    check_docker_config
    check_cert_permissions
    check_nginx_config
    check_deployment_scripts
    check_documentation

    # 显示总结
    show_summary
}

# 运行主流程
main "$@"
