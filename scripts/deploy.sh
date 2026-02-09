#!/bin/bash
# GameLink 部署脚本
# 用途：自动化部署到 Staging 或 Production 环境
# 使用方法：bash scripts/deploy.sh [staging|production]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
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

# 显示标题
show_header() {
    echo ""
    echo "=================================="
    echo " GameLink 部署工具"
    echo "=================================="
    echo ""
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_success "依赖检查通过"
}

# 生成密钥
generate_keys() {
    log_info "检查密钥配置..."

    if [ ! -f ".env.${ENV}" ]; then
        log_warning "未找到 .env.${ENV}，开始生成密钥..."

        if [ ! -f "scripts/generate-production-keys.sh" ]; then
            log_error "密钥生成脚本不存在"
            exit 1
        fi

        bash scripts/generate-production-keys.sh

        if [ ! -f ".env.${ENV}" ]; then
            log_error "密钥生成失败"
            exit 1
        fi

        log_success "密钥生成完成"
    else
        log_success "密钥配置已存在，跳过生成"
    fi
}

# 验证配置
validate_config() {
    log_info "验证配置..."

    # 加载配置
    if [ -f ".env.${ENV}" ]; then
        source .env.${ENV}
    else
        log_error "配置文件 .env.${ENV} 不存在"
        exit 1
    fi

    # 检查必需的配置项
    if [ -z "$POSTGRES_PASSWORD" ]; then
        log_error "POSTGRES_PASSWORD 未配置"
        exit 1
    fi

    if [ -z "$JWT_SECRET_KEY" ]; then
        log_error "JWT_SECRET_KEY 未配置"
        exit 1
    fi

    if [ "$CRYPTO_ENABLED" == "true" ]; then
        if [ -z "$CRYPTO_SECRET_KEY" ] || [ -z "$CRYPTO_IV" ]; then
            log_error "加密已启用但密钥未配置"
            exit 1
        fi
    fi

    log_success "配置验证通过"
}

# 构建应用
build_app() {
    log_info "构建应用..."

    # 构建管理后台
    if [ -d "admin" ]; then
        log_info "构建管理后台..."
        cd admin

        if [ ! -d "node_modules" ]; then
            log_info "安装依赖..."
            npm install --legacy-peer-deps
        fi

        log_info "构建前端..."
        if [ "$ENV" == "production" ]; then
            npm run build
        else
            npm run build
        fi

        cd ..
        log_success "管理后台构建完成"
    fi

    # 构建后端
    if [ -d "api" ]; then
        log_info "构建后端..."
        cd api

        # 这里假设后端是 Go 应用
        if [ -f "go.mod" ]; then
            go build -o bin/gamelink cmd/main.go
        fi

        cd ..
        log_success "后端构建完成"
    fi
}

# 部署数据库
deploy_database() {
    log_info "部署数据库..."

    # 启动数据库
    if [ "$ENV" == "staging" ]; then
        docker-compose -f docker-compose.${ENV}.yml up -d postgres redis
    else
        docker-compose -f docker-compose.${ENV}.yml up -d postgres redis
    fi

    # 等待数据库就绪
    log_info "等待数据库启动..."
    sleep 10

    # 执行迁移
    log_info "执行数据库迁移..."
    if [ -f "api/bin/gamelink" ]; then
        docker-compose -f docker-compose.${ENV}.yml exec -T backend /app/server migrate || true
    fi

    log_success "数据库部署完成"
}

# 部署应用
deploy_app() {
    log_info "部署应用..."

    # 启动所有服务
    docker-compose -f docker-compose.${ENV}.yml --env-file .env.${ENV} up -d --build

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 15

    log_success "应用部署完成"
}

# 验证部署
verify_deployment() {
    log_info "验证部署..."

    if [ -f "scripts/verify-deployment.sh" ]; then
        bash scripts/verify-deployment.sh
    else
        log_warning "验证脚本不存在，跳过验证"
    fi
}

# 显示状态
show_status() {
    echo ""
    echo "=================================="
    echo " 部署状态"
    echo "=================================="
    echo ""

    docker-compose -f docker-compose.${ENV}.yml --env-file .env.${ENV} ps

    echo ""
    echo "=================================="
    echo " 访问地址"
    echo "=================================="
    echo ""

    if [ "$ENV" == "staging" ]; then
        echo "前端: https://staging.gamelink.com"
        echo "后端: https://staging.gamelink.com/api/v1"
        echo "管理后台: https://admin.staging.gamelink.com"
    else
        echo "前端: https://gamelink.com"
        echo "后端: https://gamelink.com/api/v1"
        echo "管理后台: https://admin.gamelink.com"
    fi

    echo ""
    echo "=================================="
    echo " 管理员账号"
    echo "=================================="
    echo ""
    echo "邮箱: $SUPER_ADMIN_EMAIL"
    echo "密码: (查看 .env.${ENV})"
    echo ""
}

# 主流程
main() {
    show_header

    # 检查环境参数
    if [ -z "$1" ]; then
        log_error "请指定部署环境: staging 或 production"
        echo "使用方法: bash $0 [staging|production]"
        exit 1
    fi

    ENV=$1

    if [ "$ENV" != "staging" ] && [ "$ENV" != "production" ]; then
        log_error "无效的环境参数: $ENV"
        echo "支持的环境: staging, production"
        exit 1
    fi

    log_info "部署到环境: $ENV"
    echo ""

    # 确认部署
    if [ "$ENV" == "production" ]; then
        log_warning "警告：即将部署到生产环境！"
        read -p "确认继续？(yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            log_info "已取消部署"
            exit 0
        fi
    fi

    # 执行部署步骤
    check_dependencies
    generate_keys
    validate_config
    build_app
    deploy_database
    deploy_app
    verify_deployment
    show_status

    log_success "部署完成！"
}

# 运行主流程
main "$@"
