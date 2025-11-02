# Phase 1 快速开始指南

## 🚀 5分钟快速体验新功能

### 前置条件
- ✅ 后端服务正常运行
- ✅ 数据库已迁移
- ✅ 已有管理员账号和普通用户账号

---

## 步骤1: 启动应用（1分钟）

```bash
cd backend
go run ./cmd/main.go
```

**预期日志:**
```
created default commission rule: 20% (id=1)
Settlement scheduler started - will run on 1st of each month at 02:00
Server started on :8080
```

---

## 步骤2: 创建服务（1分钟）

### 2.1 管理员登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@gamelink.local",
    "password": "Admin@123456"
  }'

# 保存返回的token
export ADMIN_TOKEN="eyJhbGc..."
```

### 2.2 创建护航服务
```bash
# 创建段位护航
curl -X POST http://localhost:8080/api/v1/admin/services \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "gameId": 1,
    "name": "王者荣耀 - 王者段位护航",
    "description": "专业王者段位陪玩，快速上分",
    "type": "rank_escort",
    "pricePerHour": 8000,
    "minDuration": 1.0,
    "maxDuration": 10.0,
    "requiredRank": "王者",
    "commissionRate": 20,
    "sortOrder": 1,
    "icon": "👑"
  }'

# 保存返回的serviceId
export SERVICE_ID=1
```

### 2.3 创建礼物
```bash
# 创建玫瑰花
curl -X POST http://localhost:8080/api/v1/admin/gifts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "玫瑰花",
    "description": "表达感谢的玫瑰花",
    "icon": "🌹",
    "priceCents": 1000,
    "commissionRate": 20,
    "category": "flower",
    "sortOrder": 1
  }'

# 保存返回的giftId
export GIFT_ID=1
```

---

## 步骤3: 用户下单（1分钟）

### 3.1 用户登录
```bash
# 用普通用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "User@123456"
  }'

export USER_TOKEN="eyJhbGc..."
```

### 3.2 浏览服务
```bash
curl http://localhost:8080/api/v1/admin/services \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 3.3 创建订单（关联服务）
```bash
curl -X POST http://localhost:8080/api/v1/user/orders \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "playerId": 1,
    "gameId": 1,
    "serviceId": 1,
    "title": "王者段位护航4小时",
    "description": "需要从钻石上到王者",
    "durationHours": 4.0,
    "scheduledStart": "2024-11-20T14:00:00Z"
  }'

# 返回：订单价格 = 8000 × 4 = 32000分 (320元)
export ORDER_ID=1
```

---

## 步骤4: 完成订单（1分钟）

### 4.1 支付订单
```bash
curl -X POST http://localhost:8080/api/v1/user/payments \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": 1,
    "method": "wechat"
  }'
```

### 4.2 陪玩师接单
```bash
# 陪玩师登录
export PLAYER_TOKEN="..."

curl -X POST http://localhost:8080/api/v1/player/orders/1/accept \
  -H "Authorization: Bearer $PLAYER_TOKEN"
```

### 4.3 完成订单
```bash
# 用户确认完成
curl -X POST http://localhost:8080/api/v1/user/orders/1/complete \
  -H "Authorization: Bearer $USER_TOKEN"

# ✅ 系统自动记录抽成！
# CommissionRecord创建：
# - 订单总额: 320元
# - 平台抽成: 64元 (20%)
# - 陪玩师收入: 256元 (80%)
```

---

## 步骤5: 赠送礼物（1分钟）

```bash
# 浏览礼物
curl http://localhost:8080/api/v1/user/gifts \
  -H "Authorization: Bearer $USER_TOKEN"

# 赠送10朵玫瑰花
curl -X POST http://localhost:8080/api/v1/user/gifts/send \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "playerId": 1,
    "giftId": 1,
    "quantity": 10,
    "message": "感谢大神带我上分！",
    "isAnonymous": false,
    "orderId": 1
  }'

# ✅ 礼物记录创建：
# - 总价: 100元 (10元 × 10)
# - 平台抽成: 20元 (20%)
# - 陪玩师收入: 80元 (80%)
```

---

## 步骤6: 查看收入（1分钟）

### 6.1 陪玩师查看抽成记录
```bash
curl http://localhost:8080/api/v1/player/commission/summary?month=2024-11 \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Response:
{
  "monthlyIncome": 25600,      # 订单收入 256元
  "totalCommission": 6400,      # 平台抽成 64元
  "totalIncome": 25600,         # 累计收入
  "totalOrders": 1
}
```

### 6.2 查看礼物收入
```bash
curl http://localhost:8080/api/v1/player/gifts/stats \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Response:
{
  "totalReceived": 10,          # 收到10朵玫瑰
  "totalIncome": 8000,          # 礼物收入 80元
  "totalCount": 1               # 1条礼物记录
}
```

### 6.3 查看抽成记录
```bash
curl http://localhost:8080/api/v1/player/commission/records \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Response:
{
  "records": [
    {
      "id": 1,
      "orderId": 1,
      "totalAmountCents": 32000,
      "commissionRate": 20,
      "commissionCents": 6400,
      "playerIncomeCents": 25600,
      "settlementStatus": "pending",
      "settlementMonth": "2024-11"
    }
  ],
  "total": 1
}
```

---

## 步骤7: 管理员查看统计（1分钟）

### 7.1 平台统计
```bash
curl "http://localhost:8080/api/v1/admin/commission/stats?month=2024-11" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Response:
{
  "month": "2024-11",
  "totalOrders": 1,
  "totalIncome": 32000,         # 订单总额 320元
  "totalCommission": 6400,      # 平台抽成 64元
  "totalPlayerIncome": 25600    # 陪玩师收入 256元
}
```

### 7.2 手动触发结算（测试）
```bash
curl -X POST "http://localhost:8080/api/v1/admin/commission/settlements/trigger?month=2024-11" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# ✅ 结算成功！
# - 所有pending记录变为settled
# - 创建MonthlySettlement记录
```

### 7.3 陪玩师查看结算记录
```bash
curl http://localhost:8080/api/v1/player/commission/settlements \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Response:
{
  "settlements": [
    {
      "id": 1,
      "settlementMonth": "2024-11",
      "totalOrderCount": 1,
      "totalAmountCents": 32000,
      "totalCommissionCents": 6400,
      "totalIncomeCents": 25600,
      "bonusCents": 0,
      "finalIncomeCents": 25600,
      "status": "pending"
    }
  ],
  "total": 1
}
```

---

## ✅ 验证清单

完成上述步骤后，验证以下功能：

### 抽成机制
- [x] 订单完成自动记录抽成 ✅
- [x] 抽成比例正确（20%）✅
- [x] 陪玩师可查看抽成记录 ✅
- [x] 月度结算正常运行 ✅
- [x] 管理员可查看平台统计 ✅

### 服务分类
- [x] 可以创建6种服务类型 ✅
- [x] 订单可以关联服务 ✅
- [x] 价格从服务获取 ✅
- [x] 时长范围验证 ✅
- [x] 批量操作正常 ✅

### 礼物系统
- [x] 可以创建礼物 ✅
- [x] 用户可以赠送礼物 ✅
- [x] 礼物收入计算正确 ✅
- [x] 陪玩师可查看收到的礼物 ✅
- [x] 礼物统计正常 ✅

---

## 🎯 核心数据验证

### 数学验证

#### 订单抽成
```
订单金额: 320元
抽成比例: 20%
------------------
平台抽成: 320 × 20% = 64元 ✅
陪玩师收入: 320 - 64 = 256元 ✅
```

#### 礼物抽成
```
礼物单价: 10元
数量: 10朵
总价: 100元
抽成比例: 20%
------------------
平台抽成: 100 × 20% = 20元 ✅
陪玩师收入: 100 - 20 = 80元 ✅
```

#### 月度结算
```
订单收入: 256元
礼物收入: 80元
------------------
总收入: 256 + 80 = 336元 ✅
```

---

## 🐛 常见问题

### Q1: 定时任务何时运行？
**A**: 每月1号凌晨2点自动运行，也可以手动触发测试。

### Q2: 如何修改抽成比例？
**A**: 
```bash
# 创建特殊抽成规则
POST /admin/commission/rules
{
  "name": "VIP陪玩师15%抽成",
  "type": "special",
  "rate": 15,
  "playerId": 1
}
```

### Q3: 礼物需要支付吗？
**A**: 当前版本礼物赠送会创建记录，但实际支付逻辑需要集成（TODO标记）。

### Q4: 订单如何关联服务？
**A**: 创建订单时传入`serviceId`参数即可。如果不传，使用陪玩师时薪计算（向后兼容）。

---

## 📊 数据库验证

### 检查新表
```sql
-- 连接数据库
sqlite3 var/dev.db

-- 查看所有表
.tables

-- 应该看到以下新表:
-- commission_rules
-- commission_records
-- monthly_settlements
-- services
-- gifts
-- gift_records
```

### 查看默认规则
```sql
-- 查看默认抽成规则
SELECT * FROM commission_rules WHERE type = 'default';

-- 应该看到:
-- id=1, name="默认抽成规则", rate=20, is_active=1
```

### 查看索引
```sql
-- 查看commission相关索引
SELECT name FROM sqlite_master 
WHERE type='index' AND name LIKE 'idx_commission%';
```

---

## 🎓 完整示例脚本

```bash
#!/bin/bash

# 配置
HOST="http://localhost:8080"
ADMIN_EMAIL="admin@gamelink.local"
ADMIN_PASSWORD="Admin@123456"

echo "=== Phase 1 功能测试 ==="

# 1. 管理员登录
echo "\n[1/7] 管理员登录..."
ADMIN_TOKEN=$(curl -s -X POST "$HOST/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" \
  | jq -r '.data.token')
echo "✅ Token: ${ADMIN_TOKEN:0:20}..."

# 2. 创建服务
echo "\n[2/7] 创建护航服务..."
SERVICE_RESULT=$(curl -s -X POST "$HOST/api/v1/admin/services" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "gameId": 1,
    "name": "王者荣耀 - 王者段位护航",
    "type": "rank_escort",
    "pricePerHour": 8000,
    "minDuration": 1.0,
    "maxDuration": 10.0,
    "requiredRank": "王者",
    "commissionRate": 20
  }')
SERVICE_ID=$(echo $SERVICE_RESULT | jq -r '.data.id')
echo "✅ 服务ID: $SERVICE_ID"

# 3. 创建礼物
echo "\n[3/7] 创建礼物..."
GIFT_RESULT=$(curl -s -X POST "$HOST/api/v1/admin/gifts" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "玫瑰花",
    "icon": "🌹",
    "priceCents": 1000,
    "commissionRate": 20,
    "category": "flower"
  }')
GIFT_ID=$(echo $GIFT_RESULT | jq -r '.data.id')
echo "✅ 礼物ID: $GIFT_ID"

# 4. 查看服务列表
echo "\n[4/7] 查看服务列表..."
curl -s "$HOST/api/v1/admin/services" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data'

# 5. 查看礼物列表
echo "\n[5/7] 查看礼物列表..."
curl -s "$HOST/api/v1/user/gifts" \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.data'

# 6. 查看平台统计
echo "\n[6/7] 查看平台统计..."
curl -s "$HOST/api/v1/admin/commission/stats?month=$(date +%Y-%m)" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data'

# 7. 测试手动结算
echo "\n[7/7] 测试手动结算..."
curl -s -X POST "$HOST/api/v1/admin/commission/settlements/trigger?month=$(date +%Y-%m)" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'

echo "\n=== ✅ 所有测试完成！ ==="
```

---

## 📖 进阶操作

### 1. 创建特殊抽成规则
```bash
# VIP陪玩师15%抽成
curl -X POST "$HOST/api/v1/admin/commission/rules" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "VIP陪玩师优惠",
    "type": "special",
    "rate": 15,
    "playerId": 1
  }'

# 王者荣耀游戏特殊抽成
curl -X POST "$HOST/api/v1/admin/commission/rules" \
  -d '{
    "name": "王者荣耀特惠",
    "type": "special",
    "rate": 18,
    "gameId": 1
  }'
```

### 2. 批量操作服务
```bash
# 批量禁用服务
curl -X POST "$HOST/api/v1/admin/services/batch-update-status" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "ids": [1, 2, 3],
    "isActive": false
  }'

# 批量调价
curl -X POST "$HOST/api/v1/admin/services/batch-update-price" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "ids": [1, 2],
    "pricePerHour": 10000
  }'
```

### 3. 查看详细收入
```bash
# 陪玩师查看抽成记录
curl "$HOST/api/v1/player/commission/records?page=1&pageSize=20" \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# 查看收到的礼物
curl "$HOST/api/v1/player/gifts/received?page=1" \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# 查看月度结算
curl "$HOST/api/v1/player/commission/settlements" \
  -H "Authorization: Bearer $PLAYER_TOKEN"
```

---

## 🎯 下一步

### Phase 2: 排名激励系统
- [ ] 收入排名
- [ ] 订单量排名
- [ ] 服务质量排名
- [ ] 自动奖金发放

### Phase 3: 社交功能
- [ ] 关注系统
- [ ] 通知系统
- [ ] 动态发布
- [ ] 私信功能

---

## ✨ 总结

恭喜！您已经体验了Phase 1的所有核心功能：

- ✅ **抽成机制** - 自动计算和记录
- ✅ **服务分类** - 6种服务类型
- ✅ **礼物系统** - 情感化互动
- ✅ **月度结算** - 自动化处理

**GameLink平台现在具备完整的商业化能力！** 🎉

---

**文档版本**: 1.0  
**最后更新**: 2025-11-02  
**适用版本**: Phase 1完整版

