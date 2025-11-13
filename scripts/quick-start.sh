#!/bin/bash

# GameLink 快速启动脚本
# 一键搭建开发环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 项目信息
PROJECT_NAME="GameLink"
VERSION="v2.1.0"
BACKEND_PORT="8080"
FRONTEND_PORT="5173"

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')] ✅ $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%H:%M:%S')] ⚠️  $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%H:%M:%S')] ❌ $1${NC}"
}

info() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')] ℹ️  $1${NC}"
}

# 显示横幅
show_banner() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    🎮 GameLink 快速启动工具                  ║"
    echo "║                                                              ║"
    echo "║    现代化游戏陪玩管理平台 - Go + React 全栈项目               ║"
    echo "║                       $VERSION                          ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# 检查系统要求
check_requirements() {
    log "检查系统要求..."

    # 检查操作系统
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macOS"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="Linux"
    else
        warn "未检测到的操作系统: $OSTYPE"
    fi
    info "操作系统: $OS"

    # 检查必要命令
    local commands=("git" "curl" "wget")
    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            error "未找到命令: $cmd，请先安装"
            exit 1
        fi
    done

    # 检查 Docker
    if command -v docker &> /dev/null; then
        if docker info &> /dev/null; then
            log "Docker 已安装并运行"
        else
            error "Docker 未运行，请启动 Docker 服务"
            exit 1
        fi
    else
        warn "Docker 未安装，将使用本地部署模式"
        DEPLOY_MODE="local"
    fi

    # 检查 Docker Compose
    if command -v docker-compose &> /dev/null; then
        log "Docker Compose 已安装"
    else
        warn "Docker Compose 未安装，将使用本地部署模式"
        DEPLOY_MODE="local"
    fi

    # 设置部署模式
    if [ -z "$DEPLOY_MODE" ]; then
        DEPLOY_MODE="docker"
        log "将使用 Docker 部署模式"
    fi
}

# 检查端口占用
check_ports() {
    log "检查端口占用..."

    local ports=("$BACKEND_PORT" "$FRONTEND_PORT")
    for port in "${ports[@]}"; do
        if lsof -i :"$port" &> /dev/null; then
            warn "端口 $port 已被占用"
            read -p "是否继续？(y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    done
}

# 安装 Go
install_go() {
    if command -v go &> /dev/null; then
        local go_version=$(go version | awk '{print $3}')
        log "Go 已安装: $go_version"
        return
    fi

    log "安装 Go 1.25.3..."

    local go_version="1.25.3"
    local go_arch=""

    # 检测架构
    case $(uname -m) in
        x86_64) go_arch="amd64" ;;
        arm64) go_arch="arm64" ;;
        *)
            error "不支持的架构: $(uname -m)"
            exit 1
            ;;
    esac

    local os=""
    case "$OSTYPE" in
        darwin*) os="darwin" ;;
        linux*) os="linux" ;;
        *)
            error "不支持的操作系统: $OSTYPE"
            exit 1
            ;;
    esac

    local go_file="go${go_version}.${os}-${go_arch}.tar.gz"
    local go_url="https://golang.org/dl/${go_file}"

    # 下载并安装 Go
    cd /tmp
    wget "$go_url"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "$go_file"

    # 设置环境变量
    if ! grep -q 'export PATH=' ~/.zshrc 2>/dev/null && ! grep -q 'export PATH=' ~/.bashrc 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc 2>/dev/null || true
    fi

    export PATH=$PATH:/usr/local/go/bin

    log "Go 安装完成"
}

# 安装 Node.js
install_nodejs() {
    if command -v node &> /dev/null; then
        local node_version=$(node --version)
        if [[ "$node_version" =~ ^v1[8-9] || "$node_version" =~ ^v2[0-9] ]]; then
            log "Node.js 已安装: $node_version"
            return
        else
            warn "Node.js 版本过低: $node_version，需要 18+"
        fi
    fi

    log "安装 Node.js 18..."

    # 安装 nvm
    if [ ! -d "$HOME/.nvm" ]; then
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
        export NVM_DIR="$HOME/.nvm"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    fi

    # 安装 Node.js
    nvm install 18
    nvm use 18

    log "Node.js 安装完成"
}

# 设置环境变量
setup_environment() {
    log "设置环境变量..."

    # 复制环境配置文件
    if [ ! -f .env ]; then
        cp .env.example .env
        log "已创建 .env 配置文件"
    fi

    if [ "$DEPLOY_MODE" = "docker" ]; then
        if [ ! -f docker-compose.yml ]; then
            cp docker-compose.example.yml docker-compose.yml
            log "已创建 docker-compose.yml 配置文件"
        fi
    fi

    # 生成随机密钥
    if ! grep -q "JWT_SECRET=" .env || grep -q "change_me" .env; then
        local jwt_secret=$(openssl rand -hex 32)
        sed -i.bak "s/JWT_SECRET=.*/JWT_SECRET=$jwt_secret/" .env
        log "已生成新的 JWT 密钥"
    fi
}

# Docker 部署
deploy_docker() {
    log "使用 Docker 部署..."

    # 构建 Docker 镜像
    info "构建 Docker 镜像..."
    docker-compose build

    # 启动服务
    info "启动服务..."
    docker-compose up -d

    # 等待服务启动
    log "等待服务启动..."
    sleep 30

    # 运行数据库迁移
    info "运行数据库迁移..."
    docker-compose exec -T api make migrate || true

    log "Docker 部署完成"
}

# 本地部署
deploy_local() {
    log "使用本地模式部署..."

    # 启动数据库服务 (如果需要)
    if [ "$SKIP_DATABASE" != "true" ]; then
        start_database_services
    fi

    # 构建后端
    info "构建后端服务..."
    cd backend
    go mod download
    go build -o bin/user-service ./cmd/user-service

    # 构建前端
    info "构建前端应用..."
    cd ../frontend
    npm install

    cd ..

    # 启动服务
    info "启动服务..."
    start_local_services

    log "本地部署完成"
}

# 启动数据库服务
start_database_services() {
    log "启动数据库服务..."

    # MySQL
    if command -v brew &> /dev/null; then
        brew services start mysql 2>/dev/null || true
    elif command -v systemctl &> /dev/null; then
        sudo systemctl start mysql 2>/dev/null || true
    fi

    # Redis
    if command -v brew &> /dev/null; then
        brew services start redis 2>/dev/null || true
    elif command -v systemctl &> /dev/null; then
        sudo systemctl start redis 2>/dev/null || true
    fi
}

# 启动本地服务
start_local_services() {
    log "启动本地服务..."

    # 启动后端服务
    cd backend
    nohup ./bin/user-service > ../logs/api.log 2>&1 &
    local api_pid=$!
    echo $api_pid > ../logs/api.pid

    # 启动前端服务
    cd ../frontend
    nohup npm run dev > ../logs/frontend.log 2>&1 &
    local frontend_pid=$!
    echo $frontend_pid > ../logs/frontend.pid

    cd ..

    # 保存 PID 到文件
    echo $api_pid > logs/api.pid
    echo $frontend_pid > logs/frontend.pid

    log "服务已启动"
    log "后端 PID: $api_pid"
    log "前端 PID: $frontend_pid"
}

# 验证部署
verify_deployment() {
    log "验证部署..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:$BACKEND_PORT/health &> /dev/null; then
            log "后端服务验证成功"
            break
        fi

        if [ $attempt -eq $max_attempts ]; then
            error "后端服务验证失败"
            return 1
        fi

        info "等待后端服务启动... ($attempt/$max_attempts)"
        sleep 5
        ((attempt++))
    done

    # 验证前端服务
    if curl -f http://localhost:$FRONTEND_PORT &> /dev/null; then
        log "前端服务验证成功"
    else
        warn "前端服务验证失败，请手动检查"
    fi
}

# 显示访问信息
show_access_info() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    🎉 部署成功！                          ║"
    echo "║                                                              ║"
    echo "║    访问地址:                                                   ║"
    echo "║    🌐 前端应用: http://localhost:$FRONTEND_PORT                ║"
    echo "║    🔌 后端API: http://localhost:$BACKEND_PORT                 ║"
    echo "║    📚 API文档: http://localhost:$BACKEND_PORT/swagger         ║"
    echo "║                                                              ║"
    echo "║    管理命令:                                                   ║"
    echo "║    📋 查看状态: ./scripts/status.sh                          ║"
    echo "║    🛑 停止服务: ./scripts/stop.sh                            ║"
    echo "║    🔄 重启服务: ./scripts/restart.sh                         ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# 创建管理脚本
create_management_scripts() {
    log "创建管理脚本..."

    mkdir -p scripts logs

    # 状态检查脚本
    cat > scripts/status.sh << 'EOF'
#!/bin/bash

echo "=== GameLink 服务状态 ==="

if [ -f docker-compose.yml ]; then
    echo "Docker 模式:"
    docker-compose ps
else
    echo "本地模式:"

    # 检查后端服务
    if [ -f logs/api.pid ]; then
        local api_pid=$(cat logs/api.pid)
        if ps -p $api_pid > /dev/null; then
            echo "✅ 后端服务运行中 (PID: $api_pid)"
        else
            echo "❌ 后端服务未运行"
        fi
    else
        echo "❌ 后端服务未启动"
    fi

    # 检查前端服务
    if [ -f logs/frontend.pid ]; then
        local frontend_pid=$(cat logs/frontend.pid)
        if ps -p $frontend_pid > /dev/null; then
            echo "✅ 前端服务运行中 (PID: $frontend_pid)"
        else
            echo "❌ 前端服务未运行"
        fi
    else
        echo "❌ 前端服务未启动"
    fi
fi

echo ""
echo "端口检查:"
lsof -i :8080 2>/dev/null || echo "8080 端口未占用"
lsof -i :5173 2>/dev/null || echo "5173 端口未占用"
EOF

    # 停止服务脚本
    cat > scripts/stop.sh << 'EOF'
#!/bin/bash

echo "停止 GameLink 服务..."

if [ -f docker-compose.yml ]; then
    echo "停止 Docker 服务..."
    docker-compose down
else
    echo "停止本地服务..."

    # 停止后端服务
    if [ -f logs/api.pid ]; then
        local api_pid=$(cat logs/api.pid)
        if ps -p $api_pid > /dev/null; then
            kill $api_pid
            echo "已停止后端服务 (PID: $api_pid)"
        fi
        rm -f logs/api.pid
    fi

    # 停止前端服务
    if [ -f logs/frontend.pid ]; then
        local frontend_pid=$(cat logs/frontend.pid)
        if ps -p $frontend_pid > /dev/null; then
            kill $frontend_pid
            echo "已停止前端服务 (PID: $frontend_pid)"
        fi
        rm -f logs/frontend.pid
    fi
fi

echo "服务已停止"
EOF

    # 重启服务脚本
    cat > scripts/restart.sh << 'EOF'
#!/bin/bash

echo "重启 GameLink 服务..."

./scripts/stop.sh
sleep 5

if [ -f docker-compose.yml ]; then
    docker-compose up -d
else
    cd backend && nohup ./bin/user-service > ../logs/api.log 2>&1 & echo $! > ../logs/api.pid &
    cd ../frontend && nohup npm run dev > ../logs/frontend.log 2>&1 & echo $! > ../logs/frontend.pid &
    cd ..
fi

echo "服务重启完成"
EOF

    # 设置执行权限
    chmod +x scripts/*.sh

    log "管理脚本创建完成"
}

# 主函数
main() {
    show_banner

    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --mode)
                DEPLOY_MODE="$2"
                shift 2
                ;;
            --skip-database)
                SKIP_DATABASE="true"
                shift
                ;;
            --help)
                echo "用法: $0 [选项]"
                echo "选项:"
                echo "  --mode MODE      部署模式 (docker|local)"
                echo "  --skip-database 跳过数据库启动"
                echo "  --help          显示帮助信息"
                exit 0
                ;;
            *)
                error "未知参数: $1"
                exit 1
                ;;
        esac
    done

    log "开始部署 $PROJECT_NAME..."

    check_requirements
    check_ports

    if [ "$DEPLOY_MODE" = "local" ]; then
        install_go
        install_nodejs
    fi

    setup_environment
    create_management_scripts

    # 创建日志目录
    mkdir -p logs

    if [ "$DEPLOY_MODE" = "docker" ]; then
        deploy_docker
    else
        deploy_local
    fi

    verify_deployment
    show_access_info

    log "快速启动完成！"
}

# 捕获中断信号
trap 'error "部署被中断"; exit 1' INT TERM

# 运行主函数
main "$@"