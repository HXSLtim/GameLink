#!/bin/bash
# GameLink 端口占用检查脚本
# 检查所有 GameLink 服务所需的端口是否被占用

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
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

show_header() {
    echo ""
    echo "=================================="
    echo " GameLink 端口占用检查工具"
    echo "=================================="
    echo ""
}

# 端口配置
declare -A ports=(
    ["5173"]="Admin 前端 (Vite)"
    ["8080"]="后端 API (Go)"
    ["5175"]="User 前端 (React/Vite)"
    ["5433"]="PostgreSQL 数据库"
    ["6380"]="Redis 缓存"
)

show_header

log_info "检查 GameLink 服务端口占用情况..."
echo ""

has_conflicts=false

for port in "${!ports[@]}"; do
    service="${ports[$port]}"

    # 检查端口是否被占用
    if netstat -ano 2>/dev/null | grep -q ":$port.*LISTENING"; then
        # 获取占用进程的 PID
        pid=$(netstat -ano 2>/dev/null | grep ":$port.*LISTENING" | awk '{print $5}' | head -n1)

        # 获取进程名称
        if [ -n "$pid" ]; then
            process=$(tasklist | grep "$pid" | awk '{print $1}' | head -n1)
            log_warning "$service (端口 $port): 被 $process (PID $pid) 占用"
        else
            log_warning "$service (端口 $port): 被占用"
        fi
        has_conflicts=true
    else
        log_success "$service (端口 $port): 可用"
    fi
done

echo ""

if [ "$has_conflicts" = true ]; then
    echo "=================================="
    log_error "发现端口冲突！"
    echo "=================================="
    echo ""
    log_info "解决方案："
    echo ""
    echo "1. 查看占用进程详情："
    echo "   tasklist | findstr <PID>"
    echo ""
    echo "2. 结束占用进程："
    echo "   taskkill /PID <PID> /F"
    echo ""
    echo "3. 或者修改端口配置："
    echo "   后端: .env 中的 BACKEND_PORT"
    echo "   Admin: admin/vite.config.ts 中的 server.port"
    echo "   User: app/vite.config.ts 中的 server.port"
    echo ""
    echo "4. 重启开发服务器"
    echo ""
    exit 1
else
    echo "=================================="
    log_success "所有端口可用！"
    echo "=================================="
    echo ""
    log_info "可以启动 GameLink 开发环境"
    echo ""
    exit 0
fi
