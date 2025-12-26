#!/bin/bash
# 批量发送通知API测试脚本

# 配置
API_URL="http://localhost:8080/api/v1/admin/users/batch/notification"
# 请替换为你的实际token
TOKEN="YOUR_ADMIN_TOKEN_HERE"

echo "========================================"
echo "批量发送通知API测试脚本"
echo "========================================"
echo ""

# 测试1: 指定用户列表模式
echo "测试1: 为指定用户发送通知"
echo "目标: 用户ID 1,2,3"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "users",
    "userIds": [1, 2, 3],
    "title": "测试-系统通知",
    "content": "这是一条测试系统通知消息",
    "type": "system"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试2: 按角色筛选模式 - 单个角色
echo "测试2: 为所有player角色发送通知"
echo "角色: player"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["player"],
    "title": "测试-陪玩师培训通知",
    "content": "本周五举办陪玩师技能培训，请准时参加",
    "type": "activity"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试3: 按角色筛选模式 - 多个角色
echo "测试3: 为user和player角色发送通知"
echo "角色: user, player"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "role",
    "roles": ["user", "player"],
    "title": "测试-营销活动",
    "content": "新年优惠活动来袭，下单立减！",
    "type": "marketing"
  }' | jq .

echo ""
echo "========================================"
echo ""

# 测试4: 全体用户模式
echo "测试4: 为全体用户发送通知"
curl -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target": "all",
    "title": "测试-平台公告",
    "content": "平台即将进行系统维护，请提前做好准备",
    "type": "system"
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
    "title": "测试-错误场景",
    "content": "这是一个错误的请求",
    "type": "system"
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
    "title": "测试-错误场景",
    "content": "这是一个错误的请求",
    "type": "system"
  }' | jq .

echo ""
echo "========================================"
echo ""

echo "测试完成！"
echo ""
echo "注意事项："
echo "1. 请确保已替换脚本中的TOKEN为有效的管理员token"
echo "2. 请确保数据库中存在测试用户"
echo "3. 检查用户通知表确认通知是否正确创建"
echo "4. 查看操作日志确认记录是否正确"
echo ""
echo "查询用户通知的SQL示例："
echo "SELECT * FROM notification_events ORDER BY created_at DESC LIMIT 10;"
