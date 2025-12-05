#!/bin/bash
# 批量增加积分API测试脚本

# 配置
API_URL="http://localhost:8080/api/v1/admin/users/batch/points"
# 请替换为你的实际token
TOKEN="YOUR_ADMIN_TOKEN_HERE"

echo "========================================"
echo "批量增加积分API测试脚本"
echo "========================================"
echo ""

# 测试1: 指定用户列表模式
echo "测试1: 为指定用户增加积分"
echo "目标: 用户ID 1,2,3"
echo "积分: 500分 (5元)"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "users",
    "userIds": [1, 2, 3],
    "cents": 500,
    "reason": "测试-指定用户奖励",
    "type": "admin"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试2: 按角色筛选模式 - 单个角色
echo "测试2: 为所有player角色增加积分"
echo "角色: player"
echo "积分: 1000分 (10元)"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["player"],
    "cents": 1000,
    "reason": "测试-陪玩师月度奖励",
    "type": "admin"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试3: 按角色筛选模式 - 多个角色
echo "测试3: 为user和player角色增加积分"
echo "角色: user, player"
echo "积分: 200分 (2元)"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["user", "player"],
    "cents": 200,
    "reason": "测试-系统升级补偿",
    "type": "compensation"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试4: 全体用户模式
echo "测试4: 为全体用户增加积分"
echo "积分: 100分 (1元)"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "all",
    "cents": 100,
    "reason": "测试-平台活动奖励",
    "type": "activity"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试5: 错误场景 - target=users但未提供userIds
echo "测试5: 错误场景 - 缺少userIds"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "users",
    "cents": 100,
    "reason": "测试-错误场景",
    "type": "admin"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试6: 错误场景 - target=role但未提供roles
echo "测试6: 错误场景 - 缺少roles"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "cents": 100,
    "reason": "测试-错误场景",
    "type": "admin"
  }' | jq .

echo ""
echo "========================================"
echo ""

echo "测试完成！"
echo ""
echo "注意事项："
echo "1. 请确保已替换脚本中的TOKEN为有效的管理员token"
echo "2. 请确保数据库中存在测试用户"
echo "3. 检查用户钱包余额确认积分是否正确增加"
echo "4. 查看操作日志确认记录是否正确"
