#!/bin/bash
# GameLink 数据库恢复脚本
# 用途：从备份文件恢复 PostgreSQL 数据库
# 使用方法：bash scripts/restore-database.sh [staging|production] <backup_file>

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
    echo " GameLink 数据库恢复工具"
    echo "=================================="
    echo ""
}

# 检查参数
ENVIRONMENT=${1:-staging}
BACKUP_FILE=${2}

if [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "production" ]; then
    log_error "无效的环境参数: $ENVIRONMENT"
    echo "使用方法: bash $0 [staging|production] <backup_file>"
    exit 1
fi

if [ -z "$BACKUP_FILE" ]; then
    log_error "未指定备份文件"
    echo "使用方法: bash $0 [staging|production] <backup_file>"
    echo ""
    echo "可用的备份文件："
    case $ENVIRONMENT in
        staging)
            ls -lh ./backups/postgres/staging/*.sql.gz 2>/dev/null || echo "无备份文件"
            ;;
        production)
            ls -lh ./backups/postgres/production/*.sql.gz 2>/dev/null || echo "无备份文件"
            ;;
    esac
    exit 1
fi

# 检查备份文件是否存在
if [ ! -f "$BACKUP_FILE" ]; then
    log_error "备份文件不存在: $BACKUP_FILE"
    exit 1
fi

# 配置参数
case $ENVIRONMENT in
    staging)
        CONTAINER_NAME="gamelink-postgres-staging"
        DB_NAME="${POSTGRES_DB:-gamelink_staging}"
        ;;
    production)
        CONTAINER_NAME="gamelink-postgres"
        DB_NAME="${POSTGRES_DB:-gamelink}"
        ;;
esac

show_header

log_warning "⚠️  警告：此操作将覆盖现有数据库！"
log_info "环境: $ENVIRONMENT"
log_info "容器: $CONTAINER_NAME"
log_info "数据库: $DB_NAME"
log_info "备份文件: $BACKUP_FILE"
echo ""

# 确认操作
read -p "确认要恢复数据库？(yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    log_info "操作已取消"
    exit 0
fi

# 检查容器是否运行
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    log_error "容器 $CONTAINER_NAME 未运行"
    exit 1
fi

# 创建当前数据库的备份（安全措施）
log_info "创建当前数据库的备份..."
CURRENT_BACKUP="./backups/postgres/${ENVIRONMENT}/pre_restore_$(date +%Y%m%d_%H%M%S).sql.gz"
docker exec "$CONTAINER_NAME" pg_dump \
    -U "${POSTGRES_USER:-gamelink}" \
    -d "$DB_NAME" \
    --no-owner \
    --no-acl \
    2>&1 | gzip > "$CURRENT_BACKUP"

if [ $? -eq 0 ]; then
    log_success "当前数据库已备份到: $CURRENT_BACKUP"
else
    log_error "备份当前数据库失败，操作取消"
    exit 1
fi

# 删除现有数据库
log_info "删除现有数据库..."
docker exec "$CONTAINER_NAME" psql \
    -U "${POSTGRES_USER:-gamelink}" \
    -d postgres \
    -c "DROP DATABASE IF EXISTS $DB_NAME;"

if [ $? -eq 0 ]; then
    log_success "现有数据库已删除"
else
    log_error "删除数据库失败"
    exit 1
fi

# 创建新数据库
log_info "创建新数据库..."
docker exec "$CONTAINER_NAME" psql \
    -U "${POSTGRES_USER:-gamelink}" \
    -d postgres \
    -c "CREATE DATABASE $DB_NAME;"

if [ $? -eq 0 ]; then
    log_success "新数据库已创建"
else
    log_error "创建数据库失败"
    exit 1
fi

# 恢复数据
log_info "开始恢复数据..."
gunzip -c "$BACKUP_FILE" | docker exec -i "$CONTAINER_NAME" psql \
    -U "${POSTGRES_USER:-gamelink}" \
    -d "$DB_NAME" \
    --quiet

if [ $? -eq 0 ]; then
    log_success "数据恢复完成"
else
    log_error "数据恢复失败"
    log_info "尝试从备份恢复: $CURRENT_BACKUP"
    exit 1
fi

# 验证恢复结果
log_info "验证恢复结果..."
TABLE_COUNT=$(docker exec "$CONTAINER_NAME" psql \
    -U "${POSTGRES_USER:-gamelink}" \
    -d "$DB_NAME" \
    -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | tr -d ' ')

if [ -n "$TABLE_COUNT" ] && [ "$TABLE_COUNT" -gt 0 ]; then
    log_success "数据库已恢复，包含 $TABLE_COUNT 个表"
else
    log_warning "警告：恢复的数据库中没有表"
fi

echo ""
log_success "数据库恢复完成！"
