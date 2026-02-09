#!/bin/bash
# GameLink 快速回滚脚本
# 用于快速回滚到上一个稳定版本
# 使用方法：bash scripts/rollback.sh [staging|production]

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
    echo " GameLink 快速回滚工具"
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

log_warning "⚠️  警告：此操作将回滚到上一个版本！"
log_info "环境: $ENVIRONMENT"
echo ""

# 确认操作
read -p "确认要执行回滚操作？(yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    log_info "操作已取消"
    exit 0
fi

echo ""

# 配置参数
case $ENVIRONMENT in
    staging)
        COMPOSE_FILE="docker-compose.staging.yml"
        ENV_FILE=".env.staging"
        BACKUP_DIR="./backups/postgres/staging"
        ;;
    production)
        COMPOSE_FILE="docker-compose.prod.yml"
        ENV_FILE=".env.production"
        BACKUP_DIR="./backups/postgres/production"
        ;;
esac

# 1. 创建回滚前备份
log_info "步骤 1: 创建回滚前备份..."
echo ""

backup_timestamp=$(date +%Y%m%d_%H%M%S)
pre_rollback_backup="$BACKUP_DIR/pre_rollback_${backup_timestamp}.sql.gz"

# 备份数据库
log_info "备份数据库..."
if docker ps --format '{{.Names}}' | grep -q "gamelink-postgres${ENVIRONMENT:+-$ENVIRONMENT}"; then
    bash scripts/backup-database.sh "$ENVIRONMENT"
    log_success "数据库备份完成"
else
    log_warning "数据库容器未运行，跳过备份"
fi

echo ""

# 2. 获取当前版本信息
log_info "步骤 2: 获取当前版本信息..."
echo ""

current_commit=$(git rev-parse --short HEAD)
current_branch=$(git branch --show-current)
current_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "无标签")

echo "当前分支: $current_branch"
echo "当前提交: $current_commit"
echo "当前标签: $current_tag"
echo ""

# 3. 查找可用的回滚点
log_info "步骤 3: 查找可用的回滚点..."
echo ""

# 查找最近的 Git 标签
previous_tag=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")

if [ -z "$previous_tag" ]; then
    log_warning "未找到上一个 Git 标签"
    echo "可用的回滚选项："
    echo "1. 回滚到上一个提交: $(git rev-parse --short HEAD^)"
    echo "2. 手动指定版本"
    echo ""
    read -p "选择回滚方式 (1/2): " rollback_choice

    case $rollback_choice in
        1)
            rollback_target=$(git rev-parse --short HEAD^)
            rollback_type="commit"
            ;;
        2)
            read -p "输入要回滚的版本 (提交哈希或标签): " rollback_target
            rollback_type="manual"
            ;;
        *)
            log_error "无效的选择"
            exit 1
            ;;
    esac
else
    echo "找到上一个标签: $previous_tag"
    echo "标签提交: $(git rev-parse --short $previous_tag)"
    echo ""
    read -p "是否回滚到标签 $previous_tag？(yes/no): " use_tag

    if [ "$use_tag" = "yes" ]; then
        rollback_target=$previous_tag
        rollback_type="tag"
    else
        rollback_target=$(git rev-parse --short HEAD^)
        rollback_type="commit"
        log_info "将回滚到上一个提交: $rollback_target"
    fi
fi

echo ""
log_info "回滚目标: $rollback_target ($rollback_type)"
echo ""

# 4. 停止当前服务
log_info "步骤 4: 停止当前服务..."
echo ""

docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down

log_success "服务已停止"
echo ""

# 5. 回滚代码
log_info "步骤 5: 回滚代码到目标版本..."
echo ""

git checkout "$rollback_target"

if [ $? -eq 0 ]; then
    log_success "代码已回滚到 $rollback_target"
else
    log_error "代码回滚失败"
    exit 1
fi

echo ""

# 6. 重新构建和部署
log_info "步骤 6: 重新构建和部署服务..."
echo ""

# 构建镜像
log_info "构建 Docker 镜像..."
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build

if [ $? -eq 0 ]; then
    log_success "镜像构建完成"
else
    log_error "镜像构建失败"
    exit 1
fi

# 启动服务
log_info "启动服务..."
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d

if [ $? -eq 0 ]; then
    log_success "服务已启动"
else
    log_error "服务启动失败"
    exit 1
fi

echo ""

# 7. 等待服务就绪
log_info "步骤 7: 等待服务就绪..."
echo ""

sleep 10

# 检查容器状态
log_info "检查容器状态..."
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

echo ""

# 8. 验证服务
log_info "步骤 8: 验证服务..."
echo ""

# 运行验证脚本
if [ -f "scripts/verify-deployment.sh" ]; then
    bash scripts/verify-deployment.sh "$ENVIRONMENT"
else
    log_warning "验证脚本不存在，跳过自动验证"
fi

echo ""

# 9. 创建回滚记录
log_info "步骤 9: 记录回滚信息..."
echo ""

rollback_log="./logs/rollback_${ENVIRONMENT}_${backup_timestamp}.log"
mkdir -p "$(dirname "$rollback_log")"

cat > "$rollback_log" << EOF
回滚记录
===========================================
时间: $(date '+%Y-%m-%d %H:%M:%S')
环境: $ENVIRONMENT
回滚类型: $rollback_type
回滚目标: $rollback_target
回滚前版本: $current_commit
备份文件: $pre_rollback_backup
===========================================
EOF

log_success "回滚记录已保存: $rollback_log"
echo ""

# 10. 总结
echo "=================================="
echo " 回滚完成"
echo "=================================="
echo ""
echo "环境: $ENVIRONMENT"
echo "回滚到: $rollback_target"
echo "回滚前版本: $current_commit"
echo "备份位置: $pre_rollback_backup"
echo ""
log_success "回滚操作成功完成！"
echo ""
echo "下一步操作："
echo "1. 验证所有功能正常"
echo "2. 检查应用日志: bash scripts/collect-logs.sh $ENVIRONMENT all analyze"
echo "3. 如有问题，可以从备份恢复: bash scripts/restore-database.sh $ENVIRONMENT <backup_file>"
echo ""
