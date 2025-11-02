# 统一架构快速开始指南

## 🚀 5分钟快速上手

### 1️⃣ 启动应用

```bash
cd backend
go run ./cmd/main.go
```

**预期输出：**
```
created default commission rule: 20% (id=1)
Settlement scheduler started - will run on 1st of each month at 02:00
Server started on :8080
```

---

### 2️⃣ 创建服务项目

#### 创建护航服务（管理员）

```bash
curl -X POST http://localhost:8080/api/v1/admin/service-items \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "itemCode": "ESCORT_RANK_DIAMOND_LOL",
    "name": "英雄联盟钻石段位护航",
    "description": "专业钻石陪玩师，带你上分",
    "subCategory": "solo",
    "gameId": 1,
    "rankLevel": "钻石",
    "basePriceCents": 50000,
    "serviceHours": 1,
    "commissionRate": 0.20,
    "minUsers": 1,
    "maxPlayers": 1,
    "tags": "[\"专业\", \"上分\", \"钻石\"]",
    "iconUrl": "/icons/diamond.png",
    "sortOrder": 1
  }'
```

#### 创建礼物（管理员）

```bash
curl -X POST http://localhost:8080/api/v1/admin/service-items \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "itemCode": "ESCORT_GIFT_ROSE_PREMIUM",
    "name": "高端玫瑰",
    "description": "送给陪玩师表达感谢",
    "subCategory": "gift",
    "basePriceCents": 10000,
    "serviceHours": 0,
    "commissionRate": 0.20,
    "tags": "[\"礼物\", \"浪漫\"]",
    "iconUrl": "/icons/rose.png",
    "sortOrder": 1
  }'
```

---

### 3️⃣ 用户使用流程

#### 浏览礼物

```bash
curl http://localhost:8080/api/v1/user/gifts \
  -H "Authorization: Bearer {user_token}"
```

#### 赠送礼物

```bash
curl -X POST http://localhost:8080/api/v1/user/gifts/send \
  -H "Authorization: Bearer {user_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "playerId": 5,
    "giftItemId": 2,
    "quantity": 1,
    "message": "谢谢你的陪伴！",
    "isAnonymous": false
  }'
```

---

### 4️⃣ 陪玩师查看收入

#### 查看收到的礼物

```bash
curl http://localhost:8080/api/v1/player/gifts/received \
  -H "Authorization: Bearer {player_token}"
```

#### 查看礼物统计

```bash
curl http://localhost:8080/api/v1/player/gifts/stats \
  -H "Authorization: Bearer {player_token}"
```

#### 查看抽成汇总

```bash
curl "http://localhost:8080/api/v1/player/commission/summary?month=2024-11" \
  -H "Authorization: Bearer {player_token}"
```

---

## 💡 核心概念

### ServiceItem = 一切可购买的服务

| 类型 | sub_category | service_hours | 示例 |
|------|-------------|---------------|------|
| 单人护航 | solo | >= 1 | 段位护航、技能提升 |
| 团队护航 | team | >= 1 | 五排上分、战术训练 |
| 礼物 | gift | 0 | 玫瑰、巧克力、特效 |

### Order = 统一的订单

| 订单类型 | 特征字段 | 流程 |
|---------|---------|------|
| 护航订单 | PlayerID, GameID, ScheduledStart | 需要陪玩师接单和服务 |
| 礼物订单 | RecipientPlayerID, GiftMessage, IsAnonymous | 支付后立即送达 |

### Commission = 统一的抽成

| 来源 | 计算 | 结算 |
|------|------|------|
| 护航订单 | 总价 × 20% | 月度结算 |
| 礼物订单 | 总价 × 20% | 月度结算 |

**完全一致的逻辑！**

---

## 🎨 前端对接建议

### 服务项目展示

```tsx
// 护航服务卡片
<ServiceCard 
  item={serviceItem}
  type="escort"
  onClick={() => createOrder(serviceItem)}
/>

// 礼物卡片
<GiftCard 
  item={serviceItem}
  type="gift"
  onClick={() => sendGift(serviceItem)}
/>

// 它们都是 ServiceItem，只是展示方式不同
```

### 陪玩师收入统计

```tsx
const IncomeStats = () => {
  const { data } = usePlayerIncome();
  
  return (
    <div>
      <Statistic title="总收入" value={data.totalIncome} />
      <Statistic title="护航收入" value={data.escortIncome} />
      <Statistic title="礼物收入" value={data.giftIncome} />
    </div>
  );
};
```

---

## 📋 管理后台功能

### 服务项目管理

```
1. 统一的服务项目列表
   ├── 筛选：所有 | 护航 | 礼物
   ├── 搜索：按名称/编码
   └── 操作：编辑 | 启用/禁用 | 删除

2. 创建服务项目
   ├── 选择类型：solo | team | gift
   ├── 设置价格和抽成
   └── 上传图标

3. 批量操作
   ├── 批量调价
   └── 批量启用/禁用
```

### 财务管理

```
1. 月度结算
   ├── 查看月度统计
   ├── 手动触发结算
   └── 导出结算报表

2. 抽成规则
   ├── 默认规则：20%
   ├── 特殊规则：游戏/陪玩师/类型
   └── 规则优先级管理
```

---

## 🔍 数据查询示例

### 查询所有礼物

```go
items, total, _ := serviceItemRepo.List(ctx, ServiceItemListOptions{
    SubCategory: &model.SubCategoryGift,
    IsActive: boolPtr(true),
    Page: 1,
    PageSize: 20,
})
```

### 查询某游戏的护航服务

```go
gameID := uint64(1)
items, total, _ := serviceItemRepo.List(ctx, ServiceItemListOptions{
    GameID: &gameID,
    SubCategory: &model.SubCategorySolo,
    IsActive: boolPtr(true),
})
```

### 查询陪玩师收到的礼物订单

```go
orders, _ := orderRepo.List(ctx, OrderListOptions{
    PlayerID: &playerID,
    // 在 orders 中有 RecipientPlayerID 的就是礼物订单
})

// 过滤礼物订单
for _, order := range orders {
    if order.IsGiftOrder() {
        // 这是礼物订单
    }
}
```

---

## 🎯 下一步开发建议

### Phase 1: 完善核心功能（1周）

```
✅ 统一架构已完成
□ 创建初始服务项目数据
□ 前端对接API
□ 端到端测试
□ 性能优化
```

### Phase 2: 功能增强（2周）

```
□ OrderService集成ServiceItem（从表获取价格）
□ 通知系统（礼物送达通知）
□ 礼物特效展示
□ 陪玩师动态功能
```

### Phase 3: 运营功能（2周）

```
□ 排名激励系统
□ 数据分析报表
□ 用户行为分析
□ 推荐系统
```

---

## 🎁 示例：完整的礼物赠送流程

### Step 1: 用户浏览礼物

```
GET /api/v1/user/gifts

Response:
{
  "items": [
    {
      "id": 2,
      "itemCode": "ESCORT_GIFT_ROSE_PREMIUM",
      "name": "高端玫瑰",
      "subCategory": "gift",
      "basePriceCents": 10000,
      "serviceHours": 0,
      "commissionRate": 0.20,
      "iconUrl": "/icons/rose.png"
    }
  ]
}
```

### Step 2: 用户赠送

```
POST /api/v1/user/gifts/send
{
  "playerId": 5,
  "giftItemId": 2,
  "quantity": 3,
  "message": "感谢陪伴！"
}

系统自动：
1. 创建 Order
   - OrderNo: GIFT20241102...
   - ItemID: 2
   - RecipientPlayerID: 5
   - TotalPriceCents: 30000
   - CommissionCents: 6000
   - PlayerIncomeCents: 24000

2. 立即送达
   - Status: completed
   - DeliveredAt: now

3. 记录抽成
   - CommissionRecord 自动创建
```

### Step 3: 陪玩师查看

```
GET /api/v1/player/gifts/received

Response:
{
  "gifts": [
    {
      "orderId": 1001,
      "giftName": "高端玫瑰",
      "quantity": 3,
      "totalPrice": 30000,
      "income": 24000,      // 已扣除20%抽成
      "message": "感谢陪伴！",
      "deliveredAt": "2024-11-02T15:00:00Z"
    }
  ]
}
```

### Step 4: 月度结算

```
每月1号凌晨2点自动：
1. 汇总该陪玩师所有订单的抽成记录
2. 创建 MonthlySettlement
   - TotalIncomeCents: 包含护航+礼物
3. 可提现余额更新
```

---

## ✨ 核心价值

### 对用户 👤
- ✅ 购买护航服务提升游戏体验
- ✅ 赠送礼物表达感谢
- ✅ 统一的订单管理

### 对陪玩师 🎮
- ✅ 护航服务收入
- ✅ 礼物额外收入
- ✅ 透明的收入明细
- ✅ 自动月度结算

### 对平台 💰
- ✅ 护航订单20%抽成
- ✅ 礼物订单20%抽成
- ✅ 自动化财务管理
- ✅ 灵活的抽成规则

---

**准备好开始业务了！** 🎉🚀

