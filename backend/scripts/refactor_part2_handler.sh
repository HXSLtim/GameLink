#!/bin/bash
# Part 2: Handler层结构整合脚本
# 执行前请确保Part 1已完成或创建独立分支

set -e

echo "🚀 Part 2: Handler层结构整合开始..."

cd "$(dirname "$0")/.."

# 备份
echo "📦 创建分支..."
git checkout -b refactor/part2-handler

echo "📁 创建新目录结构..."
mkdir -p internal/handler/admin
mkdir -p internal/handler/user
mkdir -p internal/handler/player

echo "📝 迁移Admin Handler..."
# 从 internal/admin/ 迁移
if [ -f "internal/admin/game_handler.go" ]; then
    git mv internal/admin/game_handler.go internal/handler/admin/game.go
fi
if [ -f "internal/admin/user_handler.go" ]; then
    git mv internal/admin/user_handler.go internal/handler/admin/user.go
fi
if [ -f "internal/admin/player_handler.go" ]; then
    git mv internal/admin/player_handler.go internal/handler/admin/player.go
fi
if [ -f "internal/admin/order_handler.go" ]; then
    git mv internal/admin/order_handler.go internal/handler/admin/order.go
fi
if [ -f "internal/admin/payment_handler.go" ]; then
    git mv internal/admin/payment_handler.go internal/handler/admin/payment.go
fi
if [ -f "internal/admin/review_handler.go" ]; then
    git mv internal/admin/review_handler.go internal/handler/admin/review.go
fi

# 从 internal/handler/ 迁移admin相关
if [ -f "internal/handler/admin_commission.go" ]; then
    git mv internal/handler/admin_commission.go internal/handler/admin/commission.go
fi
if [ -f "internal/handler/admin_service_item.go" ]; then
    git mv internal/handler/admin_service_item.go internal/handler/admin/item.go
fi
if [ -f "internal/handler/admin_dashboard.go" ]; then
    git mv internal/handler/admin_dashboard.go internal/handler/admin/dashboard.go
fi
if [ -f "internal/handler/admin_withdraw.go" ]; then
    git mv internal/handler/admin_withdraw.go internal/handler/admin/withdraw.go
fi
if [ -f "internal/handler/admin_stats.go" ]; then
    git mv internal/handler/admin_stats.go internal/handler/admin/stats.go
fi
if [ -f "internal/handler/admin_ranking_commission.go" ]; then
    git mv internal/handler/admin_ranking_commission.go internal/handler/admin/ranking.go
fi

echo "📝 迁移User Handler..."
if [ -f "internal/handler/user_order.go" ]; then
    git mv internal/handler/user_order.go internal/handler/user/order.go
fi
if [ -f "internal/handler/user_payment.go" ]; then
    git mv internal/handler/user_payment.go internal/handler/user/payment.go
fi
if [ -f "internal/handler/user_player.go" ]; then
    git mv internal/handler/user_player.go internal/handler/user/player.go
fi
if [ -f "internal/handler/user_review.go" ]; then
    git mv internal/handler/user_review.go internal/handler/user/review.go
fi
if [ -f "internal/handler/user_gift.go" ]; then
    git mv internal/handler/user_gift.go internal/handler/user/gift.go
fi

echo "📝 迁移Player Handler..."
if [ -f "internal/handler/player_profile.go" ]; then
    git mv internal/handler/player_profile.go internal/handler/player/profile.go
fi
if [ -f "internal/handler/player_order.go" ]; then
    git mv internal/handler/player_order.go internal/handler/player/order.go
fi
if [ -f "internal/handler/player_earnings.go" ]; then
    git mv internal/handler/player_earnings.go internal/handler/player/earnings.go
fi
if [ -f "internal/handler/player_commission.go" ]; then
    git mv internal/handler/player_commission.go internal/handler/player/commission.go
fi
if [ -f "internal/handler/player_gift.go" ]; then
    git mv internal/handler/player_gift.go internal/handler/player/gift.go
fi

echo "🗑️  删除旧admin目录..."
if [ -d "internal/admin" ]; then
    # 检查是否还有文件
    if [ -z "$(ls -A internal/admin)" ]; then
        rm -rf internal/admin
    else
        echo "⚠️  internal/admin/ 还有文件，请手动检查"
    fi
fi

echo "🔄 更新cmd/main.go导入路径..."
# 更新导入路径
sed -i 's/"gamelink\/internal\/admin"/"gamelink\/internal\/handler\/admin"/g' cmd/main.go

echo "✅ 编译测试..."
go build ./...

echo "✅ Part 2 完成！"
echo "⚠️  重要：需要手动检查和更新cmd/main.go中的路由注册！"
echo "📝 请检查修改，确认无误后提交："
echo "   git add ."
echo "   git commit -m 'refactor(handler): reorganize handlers into admin/user/player directories'"
echo "   git push origin refactor/part2-handler"

