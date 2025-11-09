#!/bin/bash

# GameLink 改进项目结构快速搭建脚本
# 使用方法: bash scripts/setup-improvement-structure.sh

set -e

echo "🚀 GameLink 改进项目结构搭建开始..."
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "📂 项目根目录: $PROJECT_ROOT"
echo ""

# ============================================
# 第一部分: 后端数据模型文件
# ============================================

echo "${BLUE}📊 第一步: 创建后端数据模型文件...${NC}"

# 创建新的数据模型文件
MODELS=(
    "dispute"
    "ticket"
    "notification"
    "chat"
    "favorite"
    "tag"
)

for model in "${MODELS[@]}"; do
    file="backend/internal/model/${model}.go"
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

echo ""

# ============================================
# 第二部分: Repository 层
# ============================================

echo "${BLUE}📚 第二步: 创建 Repository 层文件...${NC}"

REPOS=(
    "dispute"
    "ticket"
    "notification"
    "chat"
    "favorite"
    "tag"
)

for repo in "${REPOS[@]}"; do
    dir="backend/internal/repository/${repo}"
    
    # 创建目录
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo "${GREEN}✓${NC} 创建目录: $dir"
    fi
    
    # 创建文件
    files=("repository.go" "repository_test.go")
    for file in "${files[@]}"; do
        full_path="$dir/$file"
        if [ ! -f "$full_path" ]; then
            touch "$full_path"
            echo "${GREEN}✓${NC} 创建: $full_path"
        else
            echo "${YELLOW}⚠${NC}  已存在: $full_path"
        fi
    done
done

echo ""

# ============================================
# 第三部分: Service 层
# ============================================

echo "${BLUE}💼 第三步: 创建 Service 层文件...${NC}"

SERVICES=(
    "dispute"
    "ticket"
    "notification"
    "chat"
    "favorite"
    "upload"
)

for service in "${SERVICES[@]}"; do
    dir="backend/internal/service/${service}"
    
    # 创建目录
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo "${GREEN}✓${NC} 创建目录: $dir"
    fi
    
    # 创建文件
    files=("service.go" "service_test.go")
    for file in "${files[@]}"; do
        full_path="$dir/$file"
        if [ ! -f "$full_path" ]; then
            touch "$full_path"
            echo "${GREEN}✓${NC} 创建: $full_path"
        else
            echo "${YELLOW}⚠${NC}  已存在: $full_path"
        fi
    done
done

# 创建支付服务文件
PAYMENT_FILES=(
    "backend/internal/service/payment/alipay.go"
    "backend/internal/service/payment/wechat.go"
)

for file in "${PAYMENT_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

# 创建聊天Hub
if [ ! -f "backend/internal/service/chat/hub.go" ]; then
    touch "backend/internal/service/chat/hub.go"
    echo "${GREEN}✓${NC} 创建: backend/internal/service/chat/hub.go"
fi

echo ""

# ============================================
# 第四部分: Handler 层
# ============================================

echo "${BLUE}🎯 第四步: 创建 Handler 层文件...${NC}"

# User Handler
USER_HANDLERS=(
    "dispute"
    "ticket"
    "notification"
    "favorite"
)

for handler in "${USER_HANDLERS[@]}"; do
    file="backend/internal/handler/user/${handler}.go"
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

# Player Handler
PLAYER_HANDLERS=(
    "online"
)

for handler in "${PLAYER_HANDLERS[@]}"; do
    file="backend/internal/handler/player/${handler}.go"
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

# WebSocket Handler
mkdir -p "backend/internal/handler/websocket"
WEBSOCKET_FILES=(
    "chat.go"
    "notification.go"
)

for file in "${WEBSOCKET_FILES[@]}"; do
    full_path="backend/internal/handler/websocket/$file"
    if [ ! -f "$full_path" ]; then
        touch "$full_path"
        echo "${GREEN}✓${NC} 创建: $full_path"
    else
        echo "${YELLOW}⚠${NC}  已存在: $full_path"
    fi
done

# Upload Handler
mkdir -p "backend/internal/handler/upload"
if [ ! -f "backend/internal/handler/upload/upload.go" ]; then
    touch "backend/internal/handler/upload/upload.go"
    echo "${GREEN}✓${NC} 创建: backend/internal/handler/upload/upload.go"
fi

echo ""

# ============================================
# 第五部分: 调度器和中间件
# ============================================

echo "${BLUE}⏰ 第五步: 创建调度器和中间件文件...${NC}"

# 调度器
mkdir -p "backend/internal/scheduler"
SCHEDULER_FILES=(
    "order_scheduler.go"
    "settlement_scheduler.go"
)

for file in "${SCHEDULER_FILES[@]}"; do
    full_path="backend/internal/scheduler/$file"
    if [ ! -f "$full_path" ]; then
        touch "$full_path"
        echo "${GREEN}✓${NC} 创建: $full_path"
    else
        echo "${YELLOW}⚠${NC}  已存在: $full_path"
    fi
done

# Prometheus中间件
if [ ! -f "backend/internal/middleware/prometheus.go" ]; then
    touch "backend/internal/middleware/prometheus.go"
    echo "${GREEN}✓${NC} 创建: backend/internal/middleware/prometheus.go"
fi

echo ""

# ============================================
# 第六部分: 前端用户端页面
# ============================================

echo "${BLUE}👥 第六步: 创建前端用户端页面...${NC}"

USER_PAGES=(
    "Home"
    "GameList"
    "PlayerList"
    "PlayerDetail"
    "OrderCreate"
    "MyOrders"
    "Profile"
)

for page in "${USER_PAGES[@]}"; do
    dir="frontend/src/pages/UserPortal/${page}"
    
    # 创建目录
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo "${GREEN}✓${NC} 创建目录: $dir"
    fi
    
    # 创建文件
    files=("index.tsx" "${page}.module.less")
    for file in "${files[@]}"; do
        full_path="$dir/$file"
        if [ ! -f "$full_path" ]; then
            touch "$full_path"
            echo "${GREEN}✓${NC} 创建: $full_path"
        else
            echo "${YELLOW}⚠${NC}  已存在: $full_path"
        fi
    done
done

echo ""

# ============================================
# 第七部分: 前端陪玩师端页面
# ============================================

echo "${BLUE}🎮 第七步: 创建前端陪玩师端页面...${NC}"

PLAYER_PAGES=(
    "Dashboard"
    "Orders"
    "Earnings"
    "Services"
    "Profile"
    "Reviews"
    "Schedule"
)

for page in "${PLAYER_PAGES[@]}"; do
    dir="frontend/src/pages/PlayerPortal/${page}"
    
    # 创建目录
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo "${GREEN}✓${NC} 创建目录: $dir"
    fi
    
    # 创建文件
    files=("index.tsx" "${page}.module.less")
    for file in "${files[@]}"; do
        full_path="$dir/$file"
        if [ ! -f "$full_path" ]; then
            touch "$full_path"
            echo "${GREEN}✓${NC} 创建: $full_path"
        else
            echo "${YELLOW}⚠${NC}  已存在: $full_path"
        fi
    done
done

echo ""

# ============================================
# 第八部分: 前端通用组件
# ============================================

echo "${BLUE}🧩 第八步: 创建前端通用组件...${NC}"

COMPONENTS=(
    "GameCard"
    "PlayerCard"
    "OrderStatusBadge"
    "ChatWindow"
    "DisputeModal"
    "TicketModal"
    "NotificationBell"
    "FavoriteButton"
)

for component in "${COMPONENTS[@]}"; do
    dir="frontend/src/components/${component}"
    
    # 创建目录
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo "${GREEN}✓${NC} 创建目录: $dir"
    fi
    
    # 创建文件
    files=("index.ts" "${component}.tsx" "${component}.module.less")
    for file in "${files[@]}"; do
        full_path="$dir/$file"
        if [ ! -f "$full_path" ]; then
            touch "$full_path"
            echo "${GREEN}✓${NC} 创建: $full_path"
        else
            echo "${YELLOW}⚠${NC}  已存在: $full_path"
        fi
    done
done

echo ""

# ============================================
# 第九部分: 前端服务层
# ============================================

echo "${BLUE}🔧 第九步: 创建前端服务层文件...${NC}"

API_FILES=(
    "dispute"
    "ticket"
    "notification"
    "favorite"
    "chat"
    "earnings"
)

for api in "${API_FILES[@]}"; do
    file="frontend/src/services/api/${api}.ts"
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

# WebSocket 服务
mkdir -p "frontend/src/services/websocket"
if [ ! -f "frontend/src/services/websocket/chat.ts" ]; then
    touch "frontend/src/services/websocket/chat.ts"
    echo "${GREEN}✓${NC} 创建: frontend/src/services/websocket/chat.ts"
fi

echo ""

# ============================================
# 第十部分: 前端类型定义
# ============================================

echo "${BLUE}📝 第十步: 创建前端类型定义文件...${NC}"

TYPE_FILES=(
    "dispute"
    "ticket"
    "notification"
    "favorite"
    "chat"
    "player"
)

for type in "${TYPE_FILES[@]}"; do
    file="frontend/src/types/${type}.ts"
    if [ ! -f "$file" ]; then
        touch "$file"
        echo "${GREEN}✓${NC} 创建: $file"
    else
        echo "${YELLOW}⚠${NC}  已存在: $file"
    fi
done

echo ""

# ============================================
# 完成
# ============================================

echo ""
echo "${GREEN}✅ 项目结构搭建完成!${NC}"
echo ""
echo "📊 统计信息:"
echo "  - 后端模型文件: 6个"
echo "  - Repository层: 6个目录, 12个文件"
echo "  - Service层: 6个目录, 12+个文件"
echo "  - Handler层: 10+个文件"
echo "  - 前端用户端页面: 7个目录, 14个文件"
echo "  - 前端陪玩师端页面: 7个目录, 14个文件"
echo "  - 前端组件: 8个目录, 24个文件"
echo "  - 前端服务层: 7个文件"
echo "  - 前端类型定义: 6个文件"
echo ""
echo "📖 下一步:"
echo "  1. 查看详细开发计划: cat GAMELINK_IMPROVEMENT_PLAN.md"
echo "  2. 查看快速摘要: cat IMPROVEMENT_SUMMARY.md"
echo "  3. 开始实现数据模型: cd backend/internal/model"
echo "  4. 运行数据库迁移: cd backend && go run cmd/server/main.go migrate up"
echo ""
echo "🎯 第一周任务:"
echo "  - Day 1-2: 实现6个新数据模型"
echo "  - Day 3-4: 实现Repository层"
echo "  - Day 5-7: 实现Service层"
echo ""
echo "💡 提示: 所有文件已创建为空文件,请根据GAMELINK_IMPROVEMENT_PLAN.md中的代码模板填充内容"
echo ""

