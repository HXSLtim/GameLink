#!/bin/bash
# GameLink 日志收集和分析脚本
# 用于收集、分析和管理应用日志
# 使用方法：bash scripts/collect-logs.sh [staging|production] [service]

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
    echo " GameLink 日志收集工具"
    echo "=================================="
    echo ""
}

# 显示帮助信息
show_help() {
    echo "使用方法: bash $0 [environment] [service] [action]"
    echo ""
    echo "参数："
    echo "  environment    环境 (staging|production)"
    echo "  service        服务 (backend|admin|postgres|redis|all)"
    echo "  action         动作 (view|analyze|export|clean)"
    echo ""
    echo "示例："
    echo "  bash $0 staging backend view        # 查看后端日志"
    echo "  bash $0 production all analyze      # 分析所有服务日志"
    echo "  bash $0 staging backend export       # 导出后端日志"
    echo "  bash $0 staging postgres clean       # 清理数据库日志"
    echo ""
}

# 检查参数
if [ $# -lt 2 ]; then
    show_help
    exit 1
fi

ENVIRONMENT=$1
SERVICE=$2
ACTION=${3:-view}

# 验证环境参数
if [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "production" ]; then
    log_error "无效的环境参数: $ENVIRONMENT"
    show_help
    exit 1
fi

show_header

log_info "环境: $ENVIRONMENT"
log_info "服务: $SERVICE"
log_info "动作: $ACTION"
echo ""

# 配置参数
case $ENVIRONMENT in
    staging)
        BACKEND_CONTAINER="gamelink-backend-staging"
        ADMIN_CONTAINER="gamelink-admin-staging"
        POSTGRES_CONTAINER="gamelink-postgres-staging"
        REDIS_CONTAINER="gamelink-redis-staging"
        LOG_DIR="./logs/staging"
        ;;
    production)
        BACKEND_CONTAINER="gamelink-backend"
        ADMIN_CONTAINER="gamelink-admin"
        POSTGRES_CONTAINER="gamelink-postgres"
        REDIS_CONTAINER="gamelink-redis"
        LOG_DIR="./logs/production"
        ;;
esac

# 创建日志目录
mkdir -p "$LOG_DIR"

# 获取容器列表
get_containers() {
    case $SERVICE in
        backend)
            echo "$BACKEND_CONTAINER"
            ;;
        admin)
            echo "$ADMIN_CONTAINER"
            ;;
        postgres)
            echo "$POSTGRES_CONTAINER"
            ;;
        redis)
            echo "$REDIS_CONTAINER"
            ;;
        all)
            echo "$BACKEND_CONTAINER" "$ADMIN_CONTAINER" "$POSTGRES_CONTAINER" "$REDIS_CONTAINER"
            ;;
        *)
            log_error "无效的服务参数: $SERVICE"
            show_help
            exit 1
            ;;
    esac
}

# 查看日志
view_logs() {
    local container=$1

    log_info "查看 $container 日志（最近 100 行）..."
    echo ""

    docker logs --tail 100 --follow "$container" 2>&1
}

# 分析日志
analyze_logs() {
    local container=$1
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local log_file="$LOG_DIR/${container}_${timestamp}.log"

    log_info "分析 $container 日志..."

    # 导出日志到文件
    docker logs --tail 1000 "$container" 2>&1 > "$log_file"

    # 统计信息
    local total_lines=$(wc -l < "$log_file")
    local error_count=$(grep -i "error" "$log_file" | wc -l)
    local warning_count=$(grep -i "warning\|warn" "$log_file" | wc -l)
    local fatal_count=$(grep -i "fatal\|panic" "$log_file" | wc -l)

    echo ""
    echo "=================================="
    echo " $container 日志分析"
    echo "=================================="
    echo ""
    echo "总行数: $total_lines"
    echo "错误数: $error_count"
    echo "警告数: $warning_count"
    echo "致命错误: $fatal_count"
    echo ""

    # 如果有错误，显示最近的错误
    if [ "$error_count" -gt 0 ]; then
        echo "最近的错误："
        echo "-----------------------------------"
        grep -i "error" "$log_file" | tail -n 10
        echo ""
    fi

    # 如果有致命错误，显示最近的致命错误
    if [ "$fatal_count" -gt 0 ]; then
        log_error "发现 $fatal_count 个致命错误！"
        echo "最近的致命错误："
        echo "-----------------------------------"
        grep -i "fatal\|panic" "$log_file" | tail -n 10
        echo ""
    fi

    log_success "日志已保存到: $log_file"
}

# 导出日志
export_logs() {
    local container=$1
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local log_file="$LOG_DIR/${container}_${timestamp}.log"
    local archive_file="$LOG_DIR/${container}_${timestamp}.tar.gz"

    log_info "导出 $container 日志..."

    # 导出完整日志
    docker logs "$container" 2>&1 > "$log_file"

    # 压缩日志
    tar -czf "$archive_file" -C "$LOG_DIR" "$(basename "$log_file")"
    rm "$log_file"

    # 显示文件大小
    local size=$(du -h "$archive_file" | cut -f1)

    log_success "日志已导出并压缩: $archive_file ($size)"
}

# 清理日志
clean_logs() {
    local container=$1

    log_info "清理 $container 容器内日志..."

    # 对于不同的容器类型，使用不同的清理方法
    case $container in
        *postgres*)
            # PostgreSQL 日志轮转
            docker exec "$container" sh -c "rm -f /var/log/postgresql/*.log*" 2>/dev/null || true
            log_success "PostgreSQL 日志已清理"
            ;;
        *redis*)
            # Redis 不需要特殊处理
            log_success "Redis 日志已清理（容器重启后自动清理）"
            ;;
        *backend*|*admin*)
            # 应用日志清理
            docker exec "$container" sh -c "rm -f /app/logs/*.log" 2>/dev/null || true
            log_success "应用日志已清理"
            ;;
    esac
}

# 执行动作
containers=$(get_containers)

for container in $containers; do
    # 检查容器是否存在
    if ! docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
        log_warning "容器 $container 不存在，跳过"
        continue
    fi

    # 检查容器是否运行
    if ! docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        log_warning "容器 $container 未运行，跳过"
        continue
    fi

    echo ""
    echo "=========================================="
    echo " 处理容器: $container"
    echo "=========================================="

    case $ACTION in
        view)
            view_logs "$container"
            ;;
        analyze)
            analyze_logs "$container"
            ;;
        export)
            export_logs "$container"
            ;;
        clean)
            clean_logs "$container"
            ;;
        *)
            log_error "无效的动作参数: $ACTION"
            show_help
            exit 1
            ;;
    esac
done

echo ""
log_success "日志处理完成！"
echo ""
echo "日志目录: $LOG_DIR"
