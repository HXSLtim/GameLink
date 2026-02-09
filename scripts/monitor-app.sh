#!/bin/bash
# GameLink 应用监控脚本
# 用于监控应用服务状态和性能指标
# 使用方法：bash scripts/monitor-app.sh [staging|production]

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
    echo " GameLink 应用监控工具"
    echo "=================================="
    echo ""
}

# 检查环境参数
ENVIRONMENT=${1:-staging}

if [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "production" ]; then
    log_error "无效的环境参数: $ENVIRONMENT"
    echo "使用方法: bash $0 [staging|production]"
    exit 1
fi

show_header

log_info "环境: $ENVIRONMENT"
log_info "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 配置参数
case $ENVIRONMENT in
    staging)
        BACKEND_CONTAINER="gamelink-backend-staging"
        ADMIN_CONTAINER="gamelink-admin-staging"
        POSTGRES_CONTAINER="gamelink-postgres-staging"
        REDIS_CONTAINER="gamelink-redis-staging"
        BACKEND_URL="http://localhost:8080"
        ADMIN_URL="http://localhost"
        ;;
    production)
        BACKEND_CONTAINER="gamelink-backend"
        ADMIN_CONTAINER="gamelink-admin"
        POSTGRES_CONTAINER="gamelink-postgres"
        REDIS_CONTAINER="gamelink-redis"
        BACKEND_URL="http://localhost:8080"
        ADMIN_URL="http://localhost"
        ;;
esac

# 1. 容器状态检查
log_info "1. 容器状态检查"
echo ""

containers=(
    "$BACKEND_CONTAINER:后端 API"
    "$ADMIN_CONTAINER:管理后台"
    "$POSTGRES_CONTAINER:数据库"
    "$REDIS_CONTAINER:缓存"
)

for item in "${containers[@]}"; do
    IFS=':' read -r container service <<< "$item"

    if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        # 获取容器状态
        status=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null)
        health=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "none")

        if [ "$health" = "healthy" ] || [ "$health" = "none" ]; then
            log_success "$service ($container): 运行中 ($status)"
        else
            log_warning "$service ($container): $health"
        fi
    else
        log_error "$service ($container): 未运行"
    fi
done

echo ""

# 2. 资源使用情况
log_info "2. 资源使用情况"
echo ""

if command -v docker-stats &> /dev/null; then
    echo "容器资源使用："
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" \
        | grep -E "NAME|$BACKEND_CONTAINER|$ADMIN_CONTAINER|$POSTGRES_CONTAINER|$REDIS_CONTAINER"
else
    log_warning "docker-stats 不可用，跳过资源统计"
fi

echo ""

# 3. 后端 API 健康检查
log_info "3. 后端 API 健康检查"
echo ""

if curl -s -f "$BACKEND_URL/api/v1/healthz" > /dev/null 2>&1; then
    # 测量响应时间
    response_time=$(curl -o /dev/null -s -w '%{time_total}' "$BACKEND_URL/api/v1/healthz")

    # 转换为毫秒
    response_ms=$(echo "$response_time * 1000" | bc)

    if (( $(echo "$response_ms < 200" | bc -l) )); then
        log_success "后端 API 响应时间: ${response_ms}ms (优秀)"
    elif (( $(echo "$response_ms < 500" | bc -l) )); then
        log_warning "后端 API 响应时间: ${response_ms}ms (可接受)"
    else
        log_error "后端 API 响应时间: ${response_ms}ms (需要优化)"
    fi

    # 获取健康检查详情
    health_data=$(curl -s "$BACKEND_URL/api/v1/healthz")
    echo "健康检查响应: $health_data"
else
    log_error "后端 API 无响应"
fi

echo ""

# 4. 数据库连接检查
log_info "4. 数据库连接检查"
echo ""

if docker ps --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER}$"; then
    # 检查数据库连接
    if docker exec "$POSTGRES_CONTAINER" pg_isready -U gamelink > /dev/null 2>&1; then
        log_success "PostgreSQL: 连接正常"

        # 获取连接数
        connections=$(docker exec "$POSTGRES_CONTAINER" psql \
            -U gamelink -d gamelink -t -c "SELECT count(*) FROM pg_stat_activity;" 2>/dev/null | tr -d ' ')

        max_connections=$(docker exec "$POSTGRES_CONTAINER" psql \
            -U gamelink -d postgres -t -c "SHOW max_connections;" 2>/dev/null | tr -d ' ')

        echo "当前连接数: $connections / $max_connections"

        # 检查慢查询
        if docker exec "$POSTGRES_CONTAINER" psql \
            -U gamelink -d gamelink -c "SELECT EXISTS (SELECT 1 FROM pg_stat_statements WHERE mean_exec_time > 1000);" \
            > /dev/null 2>&1; then
            slow_queries=$(docker exec "$POSTGRES_CONTAINER" psql \
                -U gamelink -d gamelink -t -c "SELECT count(*) FROM pg_stat_statements WHERE mean_exec_time > 1000;" \
                2>/dev/null | tr -d ' ')

            if [ "$slow_queries" -gt 0 ]; then
                log_warning "发现 $slow_queries 个慢查询 (>1秒)"
            else
                log_success "无慢查询"
            fi
        fi
    else
        log_error "PostgreSQL: 连接失败"
    fi
else
    log_warning "PostgreSQL 容器未运行"
fi

echo ""

# 5. Redis 缓存检查
log_info "5. Redis 缓存检查"
echo ""

if docker ps --format '{{.Names}}' | grep -q "^${REDIS_CONTAINER}$"; then
    # 检查 Redis 连接
    if docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" ping > /dev/null 2>&1; then
        log_success "Redis: 连接正常"

        # 获取内存使用
        memory_info=$(docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" INFO memory 2>/dev/null | grep used_memory_human | cut -d: -f2 | tr -d '\r')
        echo "内存使用: $memory_info"

        # 获取键数量
        key_count=$(docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" DBSIZE 2>/dev/null | tr -d '\r')
        echo "键数量: $key_count"

        # 检查缓存命中率
        if docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" INFO stats 2>/dev/null | grep -q keyspace; then
            hits=$(docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" INFO stats 2>/dev/null | grep keyspace_hits | cut -d: -f2 | tr -d '\r')
            misses=$(docker exec "$REDIS_CONTAINER" redis-cli -a "${REDIS_PASSWORD:-}" INFO stats 2>/dev/null | grep keyspace_misses | cut -d: -f2 | tr -d '\r')

            if [ -n "$hits" ] && [ -n "$misses" ] && [ "$((hits + misses))" -gt 0 ]; then
                hit_rate=$(echo "scale=2; $hits * 100 / ($hits + $misses)" | bc)
                echo "缓存命中率: ${hit_rate}%"
            fi
        fi
    else
        log_error "Redis: 连接失败"
    fi
else
    log_warning "Redis 容器未运行"
fi

echo ""

# 6. 日志错误检查
log_info "6. 最近的错误日志"
echo ""

# 后端错误日志
backend_errors=$(docker logs --tail 100 "$BACKEND_CONTAINER" 2>&1 | grep -i "error" | wc -l)
if [ "$backend_errors" -gt 0 ]; then
    log_warning "后端最近 100 行日志中发现 $backend_errors 个错误"
    echo "最近的错误："
    docker logs --tail 50 "$BACKEND_CONTAINER" 2>&1 | grep -i "error" | tail -n 5
else
    log_success "后端无错误日志"
fi

echo ""

# 7. 磁盘空间检查
log_info "7. 磁盘空间检查"
echo ""

disk_usage=$(df -h / | tail -n 1 | awk '{print $5}' | tr -d '%')
echo "根分区使用率: ${disk_usage}%"

if [ "$disk_usage" -gt 80 ]; then
    log_warning "磁盘空间不足 (${disk_usage}% 已使用)"
elif [ "$disk_usage" -gt 90 ]; then
    log_error "磁盘空间严重不足 (${disk_usage}% 已使用)"
else
    log_success "磁盘空间充足"
fi

# Docker 卷使用
docker_usage=$(docker system df -v 2>/dev/null | grep "VOLUME" | awk '{print $4}' | head -n 1)
if [ -n "$docker_usage" ]; then
    echo "Docker 卷使用: $docker_usage"
fi

echo ""

# 8. 总结
echo "=================================="
echo " 监控总结"
echo "=================================="
echo ""
echo "环境: $ENVIRONMENT"
echo "检查时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo "建议："

# 根据检查结果给出建议
if [ -n "$response_ms" ] && (( $(echo "$response_ms > 500" | bc -l) )); then
    echo "- ⚠️  后端 API 响应时间偏慢，建议优化"
fi

if [ -n "$connections" ] && [ "$connections" -gt 50 ]; then
    echo "- ⚠️  数据库连接数较高，建议检查连接池配置"
fi

if [ "$disk_usage" -gt 80 ]; then
    echo "- ⚠️  磁盘空间不足，建议清理旧日志或备份数据"
fi

if [ "$backend_errors" -gt 10 ]; then
    echo "- ⚠️  后端错误日志较多，建议检查应用日志"
fi

echo ""
echo "如需详细信息，请查看各服务日志："
echo "- 后端: docker logs -f $BACKEND_CONTAINER"
echo "- 数据库: docker logs -f $POSTGRES_CONTAINER"
echo "- Redis: docker logs -f $REDIS_CONTAINER"
echo ""