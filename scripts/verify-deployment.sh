#!/bin/bash
# GameLink 部署验证工具
# 用途：验证生产环境部署是否成功
# 使用方法：bash scripts/verify-deployment.sh

set -e

echo "=================================="
echo " GameLink 部署验证工具"
echo "=================================="
echo ""

# 配置
BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"
FRONTEND_URL="${FRONTEND_URL:-http://localhost}"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
PASSED=0
FAILED=0
WARNINGS=0

# 测试函数
test_endpoint() {
    local name="$1"
    local url="$2"
    local expected_code="${3:-200}"

    echo -n "测试 $name... "

    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")

    if [ "$response" == "$expected_code" ]; then
        echo -e "${GREEN}✅ 通过${NC} (HTTP $response)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ 失败${NC} (HTTP $response, 期望 $expected_code)"
        ((FAILED++))
        return 1
    fi
}

test_json_endpoint() {
    local name="$1"
    local url="$2"
    local field="$3"

    echo -n "测试 $name... "

    response=$(curl -s "$url" 2>/dev/null)

    if echo "$response" | jq -e ".${field}" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ 通过${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ 失败${NC} (字段 $field 不存在)"
        ((FAILED++))
        return 1
    fi
}

echo "开始验证部署..."
echo ""

# 1. 后端健康检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 1. 后端健康检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "后端健康检查" "$BACKEND_URL/api/v1/healthz" 200
echo ""

# 2. 前端访问检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 2. 前端访问检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "前端首页" "$FRONTEND_URL/" 200
echo ""

# 3. API 端点检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 3. API 端点检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_json_endpoint "游戏列表API" "$BACKEND_URL/api/v1/games?page=1&page_size=10" "data"
test_endpoint "游戏列表(未授权)" "$BACKEND_URL/api/v1/players?page=1&page_size=10" 401
echo ""

# 4. 数据库连接检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 4. 数据库连接检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -n "检查数据库连接... "
db_check=$(curl -s "$BACKEND_URL/api/v1/healthz" | jq -r '.database' 2>/dev/null || echo "error")
if [ "$db_check" == "ok" ]; then
    echo -e "${GREEN}✅ 正常${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ 失败${NC}"
    ((FAILED++))
fi
echo ""

# 5. Redis 连接检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 5. Redis 连接检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -n "检查 Redis 连接... "
redis_check=$(curl -s "$BACKEND_URL/api/v1/healthz" | jq -r '.cache' 2>/dev/null || echo "error")
if [ "$redis_check" == "ok" ]; then
    echo -e "${GREEN}✅ 正常${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ 失败${NC}"
    ((FAILED++))
fi
echo ""

# 6. WebSocket 检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 6. WebSocket 连接检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -n "检查 WebSocket 端点... "
ws_check=$(curl -s -i -H "Connection: Upgrade" -H "Upgrade: websocket" "$BACKEND_URL/api/v1/ws" 2>/dev/null | grep -i "101 switching protocols" || echo "")
if [ -n "$ws_check" ]; then
    echo -e "${GREEN}✅ 正常${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  跳过${NC} (需要 WebSocket 客户端)"
    ((WARNINGS++))
fi
echo ""

# 7. 容器健康检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 7. 容器健康检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if command -v docker &> /dev/null; then
    echo "检查 Docker 容器状态..."

    # 检查后端容器
    echo -n "  后端容器... "
    if docker ps | grep -q "gamelink-backend"; then
        echo -e "${GREEN}✅ 运行中${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ 未运行${NC}"
        ((FAILED++))
    fi

    # 检查前端容器
    echo -n "  前端容器... "
    if docker ps | grep -q "gamelink-admin"; then
        echo -e "${GREEN}✅ 运行中${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  未找到${NC}"
        ((WARNINGS++))
    fi

    # 检查数据库容器
    echo -n "  数据库容器... "
    if docker ps | grep -q "gamelink-postgres"; then
        echo -e "${GREEN}✅ 运行中${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  未找到${NC}"
        ((WARNINGS++))
    fi

    # 检查 Redis 容器
    echo -n "  Redis 容器... "
    if docker ps | grep -q "gamelink-redis"; then
        echo -e "${GREEN}✅ 运行中${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  未找到${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠️  Docker 未安装，跳过容器检查${NC}"
    ((WARNINGS++))
fi
echo ""

# 8. 安全配置检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 8. 安全配置检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查加密是否启用
echo -n "检查加密配置... "
if [ -f ".env.production" ]; then
    source .env.production
    if [ "$CRYPTO_ENABLED" == "true" ]; then
        echo -e "${GREEN}✅ 已启用${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ 未启用${NC}"
        ((FAILED++))
    fi
else
    echo -e "${YELLOW}⚠️  配置文件未找到${NC}"
    ((WARNINGS++))
fi

# 检查 HTTPS（如果配置了）
if [[ "$BACKEND_URL" == https://* ]]; then
    echo -n "检查 HTTPS 配置... "
    if curl -s -k "$BACKEND_URL/api/v1/healthz" | grep -q "ok"; then
        echo -e "${GREEN}✅ 正常${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ 失败${NC}"
        ((FAILED++))
    fi
fi
echo ""

# 9. 日志检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 9. 日志检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if command -v docker &> /dev/null; then
    echo "检查后端日志中的错误..."
    error_count=$(docker logs gamelink-backend --tail 100 2>&1 | grep -i "error" | wc -l)

    if [ "$error_count" -eq 0 ]; then
        echo -e "${GREEN}✅ 未发现错误${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  发现 $error_count 个错误${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠️  Docker 未安装，跳过日志检查${NC}"
    ((WARNINGS++))
fi
echo ""

# 10. 性能检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " 10. 性能检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -n "API 响应时间... "
start_time=$(date +%s%N)
curl -s "$BACKEND_URL/api/v1/healthz" > /dev/null 2>&1
end_time=$(date +%s%N)
response_time=$(( (end_time - start_time) / 1000000 ))

if [ "$response_time" -lt 1000 ]; then
    echo -e "${GREEN}✅ ${response_time}ms${NC} (良好)"
    ((PASSED++))
elif [ "$response_time" -lt 3000 ]; then
    echo -e "${YELLOW}⚠️  ${response_time}ms${NC} (可接受)"
    ((WARNINGS++))
else
    echo -e "${RED}❌ ${response_time}ms${NC} (需要优化)"
    ((FAILED++))
fi
echo ""

# 总结
echo "=================================="
echo " 验证结果总结"
echo "=================================="
echo ""
echo -e "${GREEN}✅ 通过: $PASSED${NC}"
echo -e "${YELLOW}⚠️  警告: $WARNINGS${NC}"
echo -e "${RED}❌ 失败: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 部署验证通过！${NC}"
    echo ""
    echo "系统已准备就绪，可以开始使用。"
    exit 0
else
    echo -e "${RED}❌ 部署验证失败！${NC}"
    echo ""
    echo "请检查失败的项并修复问题。"
    exit 1
fi
