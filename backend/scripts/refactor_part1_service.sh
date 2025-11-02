#!/bin/bash
# Part 1: Service层文件重命名脚本
# 执行前请确保代码已提交

set -e

echo "🚀 Part 1: Service层文件重命名开始..."

cd "$(dirname "$0")/.."

# 备份当前状态
echo "📦 创建备份分支..."
git checkout -b refactor/part1-service-backup
git checkout -b refactor/part1-service

echo "📝 重命名Service文件..."

# Auth
cd internal/service/auth
git mv auth_service.go auth.go 2>/dev/null || mv auth_service.go auth.go
git mv auth_service_test.go auth_test.go 2>/dev/null || mv auth_service_test.go auth_test.go
cd ../../..

# Order
cd internal/service/order
git mv order_service.go order.go 2>/dev/null || mv order_service.go order.go
git mv order_service_test.go order_test.go 2>/dev/null || mv order_service_test.go order_test.go
cd ../../..

# Player
cd internal/service/player
git mv player_service.go player.go 2>/dev/null || mv player_service.go player.go
git mv player_service_test.go player_test.go 2>/dev/null || mv player_service_test.go player_test.go
cd ../../..

# Payment
cd internal/service/payment
git mv payment_service.go payment.go 2>/dev/null || mv payment_service.go payment.go
git mv payment_service_test.go payment_test.go 2>/dev/null || mv payment_service_test.go payment_test.go
cd ../../..

# Review
cd internal/service/review
git mv review_service.go review.go 2>/dev/null || mv review_service.go review.go
git mv review_service_test.go review_test.go 2>/dev/null || mv review_service_test.go review_test.go
cd ../../..

# Earnings
cd internal/service/earnings
git mv earnings_service.go earnings.go 2>/dev/null || mv earnings_service.go earnings.go
git mv earnings_service_test.go earnings_test.go 2>/dev/null || mv earnings_service_test.go earnings_test.go
cd ../../..

# Gift
cd internal/service/gift
git mv gift_service.go gift.go 2>/dev/null || mv gift_service.go gift.go
git mv gift_service_test.go gift_test.go 2>/dev/null || mv gift_service_test.go gift_test.go
cd ../../..

# ServiceItem → Item (包重命名)
echo "📦 重命名 serviceitem → item ..."
git mv internal/service/serviceitem internal/service/item 2>/dev/null || mv internal/service/serviceitem internal/service/item
cd internal/service/item
git mv service_item.go item.go 2>/dev/null || mv service_item.go item.go
git mv service_item_test.go item_test.go 2>/dev/null || mv service_item_test.go item_test.go
cd ../../..

# Commission
cd internal/service/commission
git mv commission_service.go commission.go 2>/dev/null || mv commission_service.go commission.go
git mv commission_service_test.go commission_test.go 2>/dev/null || mv commission_service_test.go commission_test.go
cd ../../..

# Ranking
cd internal/service/ranking
git mv ranking_service.go ranking.go 2>/dev/null || mv ranking_service.go ranking.go
git mv ranking_service_test.go ranking_test.go 2>/dev/null || mv ranking_service_test.go ranking_test.go
cd ../../..

echo "🔄 更新cmd/main.go中的导入路径..."
# serviceitem → item
sed -i 's/serviceitemservice/itemservice/g' cmd/main.go
sed -i 's/service\/serviceitem/service\/item/g' cmd/main.go

echo "✅ 编译测试..."
go build ./...

echo "✅ 运行测试..."
go test ./internal/service/... -v

echo "✅ Part 1 完成！"
echo "📝 请检查修改，确认无误后提交："
echo "   git add ."
echo "   git commit -m 'refactor(service): remove redundant _service suffix'"
echo "   git push origin refactor/part1-service"

