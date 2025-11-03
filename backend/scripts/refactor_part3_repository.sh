#!/bin/bash
# Part 3: Repository层文件重命名脚本
# 执行前请确保代码已提交

set -e

echo "🚀 Part 3: Repository层文件重命名开始..."

cd "$(dirname "$0")/.."

# 备份
echo "📦 创建分支..."
git checkout -b refactor/part3-repository

echo "📝 重命名Repository文件..."

# User Repository
cd internal/repository/user
git mv user_gorm_repository.go repository.go 2>/dev/null || mv user_gorm_repository.go repository.go
git mv user_gorm_repository_test.go repository_test.go 2>/dev/null || mv user_gorm_repository_test.go repository_test.go
cd ../../..

# Player Repository
cd internal/repository/player
git mv player_gorm_repository.go repository.go 2>/dev/null || mv player_gorm_repository.go repository.go
git mv player_gorm_repository_test.go repository_test.go 2>/dev/null || mv player_gorm_repository_test.go repository_test.go
cd ../../..

# Game Repository
cd internal/repository/game
git mv game_gorm_repository.go repository.go 2>/dev/null || mv game_gorm_repository.go repository.go
git mv game_gorm_repository_test.go repository_test.go 2>/dev/null || mv game_gorm_repository_test.go repository_test.go
cd ../../..

# Order Repository
cd internal/repository/order
git mv order_gorm_repository.go repository.go 2>/dev/null || mv order_gorm_repository.go repository.go
git mv order_gorm_repository_test.go repository_test.go 2>/dev/null || mv order_gorm_repository_test.go repository_test.go
cd ../../..

# Payment Repository
cd internal/repository/payment
git mv payment_gorm_repository.go repository.go 2>/dev/null || mv payment_gorm_repository.go repository.go
git mv payment_gorm_repository_test.go repository_test.go 2>/dev/null || mv payment_gorm_repository_test.go repository_test.go
cd ../../..

# Review Repository
cd internal/repository/review
git mv review_gorm_repository.go repository.go 2>/dev/null || mv review_gorm_repository.go repository.go
git mv review_gorm_repository_test.go repository_test.go 2>/dev/null || mv review_gorm_repository_test.go repository_test.go
cd ../../..

# PlayerTag Repository
cd internal/repository/player_tag
git mv player_tag_gorm_repository.go repository.go 2>/dev/null || mv player_tag_gorm_repository.go repository.go
git mv player_tag_gorm_repository_test.go repository_test.go 2>/dev/null || mv player_tag_gorm_repository_test.go repository_test.go
cd ../../..

# Stats Repository
cd internal/repository/stats
git mv stats_gorm_repository.go repository.go 2>/dev/null || mv stats_gorm_repository.go repository.go
git mv stats_gorm_repository_test.go repository_test.go 2>/dev/null || mv stats_gorm_repository_test.go repository_test.go
cd ../../..

# Permission Repository
cd internal/repository/permission
git mv permission_gorm_repository.go repository.go 2>/dev/null || mv permission_gorm_repository.go repository.go
git mv permission_gorm_repository_test.go repository_test.go 2>/dev/null || mv permission_gorm_repository_test.go repository_test.go
cd ../../..

# Role Repository
cd internal/repository/role
git mv role_gorm_repository.go repository.go 2>/dev/null || mv role_gorm_repository.go repository.go
git mv role_gorm_repository_test.go repository_test.go 2>/dev/null || mv role_gorm_repository_test.go repository_test.go
cd ../../..

# OperationLog Repository (如果存在)
if [ -d "internal/repository/operation_log" ]; then
    cd internal/repository/operation_log
    git mv operation_log_gorm_repository.go repository.go 2>/dev/null || mv operation_log_gorm_repository.go repository.go
    git mv operation_log_gorm_repository_test.go repository_test.go 2>/dev/null || mv operation_log_gorm_repository_test.go repository_test.go
    cd ../../..
fi

echo "✅ 编译测试..."
go build ./...

echo "✅ 运行Repository测试..."
go test ./internal/repository/... -v

echo "✅ Part 3 完成！"
echo "📝 请检查修改，确认无误后提交："
echo "   git add ."
echo "   git commit -m 'refactor(repository): simplify repository filenames'"
echo "   git push origin refactor/part3-repository"


