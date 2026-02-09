#!/bin/bash

# GameLink 监控服务启动脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}GameLink 监控服务启动${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 检查必要的目录
echo "检查监控配置目录..."
if [ ! -d "monitoring/prometheus" ]; then
    echo -e "${YELLOW}警告: monitoring/prometheus 目录不存在${NC}"
    echo "请先创建监控配置文件"
    exit 1
fi

# 检查 Docker Compose 文件
if [ ! -f "docker-compose.monitoring.yml" ]; then
    echo -e "${YELLOW}警告: docker-compose.monitoring.yml 文件不存在${NC}"
    exit 1
fi

# 停止现有服务
echo "停止现有监控服务..."
docker-compose -f docker-compose.monitoring.yml down 2>/dev/null || true

# 拉取最新镜像
echo "拉取最新镜像..."
docker-compose -f docker-compose.monitoring.yml pull

# 启动监控服务
echo "启动监控服务..."
docker-compose -f docker-compose.monitoring.yml up -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "=== 服务状态 ==="
docker-compose -f docker-compose.monitoring.yml ps

# 验证服务可访问
echo ""
echo "=== 服务验证 ==="

echo -n "Prometheus (http://localhost:9090)... "
if curl -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 运行正常${NC}"
else
    echo -e "${YELLOW}⚠ 启动中...${NC}"
fi

echo -n "Grafana (http://localhost:3001)... "
if curl -s http://localhost:3001/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 运行正常${NC}"
else
    echo -e "${YELLOW}⚠ 启动中...${NC}"
fi

echo -n "cAdvisor (http://localhost:8081)... "
if curl -s http://localhost:8081/healthz > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 运行正常${NC}"
else
    echo -e "${YELLOW}⚠ 启动中...${NC}"
fi

# 显示访问信息
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}监控服务启动完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "访问地址："
echo "  - Prometheus:  http://localhost:9090"
echo "  - Grafana:      http://localhost:3001 (admin/admin123)"
echo "  - cAdvisor:     http://localhost:8081"
echo "  - Alertmanager: http://localhost:9093"
echo ""
echo "命令："
echo "  - 查看日志: docker-compose -f docker-compose.monitoring.yml logs -f"
echo "  - 停止服务: docker-compose -f docker-compose.monitoring.yml down"
echo "  - 重启服务: docker-compose -f docker-compose.monitoring.yml restart"
echo ""
