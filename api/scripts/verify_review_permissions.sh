#!/bin/bash

# 验证评价管理路由是否都应用了权限中间件
# Verify that all review management routes have permission middleware applied

echo "=== 验证评价管理权限中间件 ==="
echo ""

# 检查评价管理路由
echo "1. 检查评价管理路由 (RegisterRoutes)..."
grep -n "RequirePermission.*reviews" backend/internal/handler/admin/router.go | head -20
echo ""

# 检查敏感词管理路由
echo "2. 检查敏感词管理路由 (RegisterSensitiveWordRoutes)..."
grep -n "RequirePermission.*sensitive-words" backend/internal/handler/admin/router.go | head -10
echo ""

# 检查评价统计路由
echo "3. 检查评价统计路由 (RegisterReviewStatsRoutes)..."
grep -n "RequirePermission.*reviews/stats\|RequirePermission.*reviews/trend\|RequirePermission.*reviews/top-players\|RequirePermission.*reviews/game-stats\|RequirePermission.*reviews/export" backend/internal/handler/admin/router.go
echo ""

# 检查评价展示设置路由
echo "4. 检查评价展示设置路由 (RegisterReviewSettingsRoutes)..."
grep -n "RequirePermission.*review-settings" backend/internal/handler/admin/router.go
echo ""

# 检查权限种子数据
echo "5. 检查权限种子数据..."
echo "运行权限种子数据测试..."
cd backend && go test -v ./pkg/db -run TestSeedReviewPermissions
echo ""

echo "=== 验证完成 ==="
