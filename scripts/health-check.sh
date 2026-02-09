#!/bin/bash

# GameLink 健康检查脚本
# 用于检查所有服务的运行状态

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 状态统计
TOTAL=0
PASSED=0
FAILED=0
WARNINGS=0

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    ((WARNINGS++))
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    ((FAILED++))
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((PASSED++))
}

((TOTAL++))

# 分隔线
separator() {
    echo "========================================="
}

# 检查端口是否监听
check_port() {
    local port=$1
    local name=$2
    ((TOTAL++))

    if netstat -ano | grep "LISTENING" | grep ":${port}" > /dev/null 2>&1; then
        log_success "${name} (端口 ${port}) 正在监听"
        return 0
    else
        log_error "${name} (端口 ${port}) 未监听"
        return 1
    fi
}

# 检查 Docker 容器状态
check_container() {
    local container=$1
    local name=$2
    ((TOTAL++))

    local status=$(docker inspect -f '{{.State.Status}}' ${container} 2>/dev/null || echo "not_found")
    local health=$(docker inspect -f '{{.State.Health.Status}}' ${container} 2>/dev/null || echo "none")

    if [ "${status}" == "running" ]; then
        if [ "${health}" == "healthy" ] || [ "${health}" == "none" ]; then
            log_success "${name} 容器运行中 (${status})"
            return 0
        else
            log_warn "${name} 容器运行中但健康检查异常: ${health}"
            return 1
        fi
    else
        log_error "${name} 容器未运行: ${status}"
        return 1
    fi
}

# 检查 HTTP 端点
check_http() {
    local url=$1
    local name=$2
    ((TOTAL++))

    local response=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 ${url} 2>/dev/null || echo "000")

    if [ "${response}" == "200" ] || [ "${response}" == "304" ]; then
        log_success "${name} HTTP 状态正常 (${response})"
        return 0
    else
        log_error "${name} HTTP 状态异常: ${response}"
        return 1
    fi
}

# 检查磁盘使用
check_disk() {
    ((TOTAL++))

    local usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')

    if [ ${usage} -lt 70 ]; then
        log_success "磁盘使用率正常: ${usage}%"
        return 0
    elif [ ${usage} -lt 85 ]; then
        log_warn "磁盘使用率较高: ${usage}%"
        return 1
    else
        log_error "磁盘空间不足: ${usage}%"
        return 1
    fi
}

# 检查内存使用
check_memory() {
    ((TOTAL++))

    # 使用 free 命令（Linux）或 vmstat
    if command -v free &> /dev/null; then
        local mem_percent=$(free | grep Mem | awk '{printf("%.0f", $3/$2 * 100.0)}')

        if [ ${mem_percent} -lt 70 ]; then
            log_success "内存使用率正常: ${mem_percent}%"
            return 0
        elif [ ${mem_percent} -lt 90 ]; then
            log_warn "内存使用率较高: ${mem_percent}%"
            return 1
        else
            log_error "内存不足: ${mem_percent}%"
            return 1
        fi
    else
        log_warn "无法检查内存使用率（free 命令不可用）"
        return 1
    fi
}

# 主检查函数
main() {
    separator
    echo "GameLink 服务健康检查"
    echo "检查时间: $(date '+%Y-%m-%d %H:%M:%S')"
    separator
    echo ""

    # Docker 容器检查
    echo "=== Docker 容器检查 ==="
    check_container "gamelink-postgres" "PostgreSQL 数据库"
    check_container "gamelink-redis" "Redis 缓存"
    check_container "gamelink-admin" "Admin 前端"
    echo ""

    # 端口监听检查
    echo "=== 服务端口检查 ==="
    check_port "5433" "PostgreSQL 数据库"
    check_port "6380" "Redis 缓存"
    check_port "8080" "后端 API"
    check_port "5173" "Admin 前端"
    check_port "3000" "App 前端"
    echo ""

    # HTTP 健康检查
    echo "=== HTTP 端点检查 ==="
    check_http "http://localhost:8080/health" "后端 API"
    check_http "http://localhost:5173" "Admin 前端"
    echo ""

    # 系统资源检查
    echo "=== 系统资源检查 ==="
    check_disk
    check_memory
    echo ""

    # 容器资源使用
    echo "=== 容器资源使用 ==="
    ((TOTAL++))
    if docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" 2>/dev/null | grep -v "CONTAINER"; then
        log_success "容器资源使用获取成功"
        docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
    else
        log_error "无法获取容器资源使用"
    fi
    echo ""

    # 汇总
    separator
    echo "=== 检查汇总 ==="
    echo "总检查项: ${TOTAL}"
    echo -e "${GREEN}通过: ${PASSED}${NC}"
    echo -e "${YELLOW}警告: ${WARNINGS}${NC}"
    echo -e "${RED}失败: ${FAILED}${NC}"
    separator

    # 返回状态码
    if [ ${FAILED} -gt 0 ]; then
        echo -e "${RED}✗ 健康检查失败${NC}"
        exit 1
    else
        echo -e "${GREEN}✓ 健康检查通过${NC}"
        exit 0
    fi
}

# 执行主函数
main
