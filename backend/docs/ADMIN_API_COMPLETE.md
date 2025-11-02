# 🎯 GameLink 管理端API完整文档

## 📊 总览

**完成日期**: 2025-11-02  
**API版本**: v1  
**总端点数**: 50+  
**认证方式**: JWT Bearer Token  
**权限控制**: RBAC细粒度权限  

---

## 🎯 新增的管理端接口

### 1. 服务项目管理（Service Items）⭐ 核心

**统一管理护航服务和礼物**

```bash
# 创建服务项目（护航或礼物）
POST /api/v1/admin/service-items
{
  "itemCode": "ESCORT_GIFT_ROSE",
  "name": "高端玫瑰",
  "subCategory": "gift",          # solo/team/gift
  "basePriceCents": 10000,
  "serviceHours": 0,              # 礼物为0，护航>=1
  "commissionRate": 0.20
}

# 获取服务项目列表
GET /api/v1/admin/service-items?subCategory=gift&page=1&pageSize=20

# 获取服务项目详情
GET /api/v1/admin/service-items/:id

# 更新服务项目
PUT /api/v1/admin/service-items/:id
{
  "basePriceCents": 12000,
  "isActive": true
}

# 删除服务项目
DELETE /api/v1/admin/service-items/:id

# 批量更新状态
POST /api/v1/admin/service-items/batch-update-status
{
  "ids": [1, 2, 3],
  "isActive": false
}

# 批量更新价格
POST /api/v1/admin/service-items/batch-update-price
{
  "ids": [1, 2, 3],
  "basePriceCents": 15000
}
```

---

### 2. 抽成管理（Commission）⭐ 核心

**管理平台抽成规则和结算**

```bash
# 创建抽成规则
POST /api/v1/admin/commission/rules
{
  "name": "王者荣耀特殊抽成",
  "type": "special",
  "rate": 15,                     # 15%
  "gameId": 1
}

# 更新抽成规则
PUT /api/v1/admin/commission/rules/:id
{
  "rate": 18,
  "isActive": true
}

# 手动触发月度结算
POST /api/v1/admin/commission/settlements/trigger?month=2024-11

# 获取平台统计
GET /api/v1/admin/commission/stats?month=2024-11
Response:
{
  "month": "2024-11",
  "totalOrders": 156,
  "totalIncome": 1560000,        # 总收入156万分 = 15600元
  "totalCommission": 312000,     # 平台抽成31.2万分 = 3120元
  "totalPlayerIncome": 1248000   # 陪玩师收入124.8万分 = 12480元
}
```

---

### 3. 提现审核（Withdraw）⭐ 新增

**审核和管理陪玩师提现申请**

```bash
# 获取提现申请列表
GET /api/v1/admin/withdraws?status=pending&page=1&pageSize=20
Response:
{
  "withdraws": [
    {
      "id": 1,
      "playerId": 5,
      "amountCents": 50000,
      "method": "alipay",
      "accountInfo": "185****1234",
      "status": "pending",
      "createdAt": "2024-11-02T10:00:00Z"
    }
  ],
  "total": 15
}

# 获取提现详情
GET /api/v1/admin/withdraws/:id

# 批准提现
POST /api/v1/admin/withdraws/:id/approve
{
  "remark": "已确认账户信息，批准提现"
}

# 拒绝提现
POST /api/v1/admin/withdraws/:id/reject
{
  "reason": "账户信息不符"
}

# 完成提现（已打款）
POST /api/v1/admin/withdraws/:id/complete
```

**提现流程：**
```
pending → approve → complete
   ↓
pending → reject
```

---

### 4. Dashboard统计（Dashboard）⭐ 新增

**管理后台首页数据展示**

```bash
# 总览统计
GET /api/v1/admin/dashboard/overview
Response:
{
  "totalUsers": 1500,          # 总用户数
  "totalPlayers": 85,          # 总陪玩师数
  "totalOrders": 3200,         # 总订单数
  "todayOrders": 45,           # 今日订单
  "todayRevenue": 456000,      # 今日收入（分）
  "monthRevenue": 5600000,     # 本月收入（分）
  "pendingWithdraws": 8,       # 待审核提现
  "activeServices": 25         # 活跃服务项目
}

# 最近订单
GET /api/v1/admin/dashboard/recent-orders?limit=10

# 最近提现
GET /api/v1/admin/dashboard/recent-withdraws?limit=10

# 月度收入趋势
GET /api/v1/admin/dashboard/monthly-revenue?months=12
Response:
{
  "revenue": [
    {
      "month": "2024-01",
      "totalRevenue": 4500000,
      "totalCommission": 900000,
      "totalOrders": 450
    },
    ...
  ]
}
```

---

### 5. 数据统计（Stats）⭐ 新增

**深度数据分析**

```bash
# 服务项目销售统计
GET /api/v1/admin/stats/service-items
Response:
{
  "items": [
    {
      "itemId": 1,
      "itemCode": "ESCORT_RANK_DIAMOND",
      "itemName": "钻石段位护航",
      "subCategory": "solo",
      "orderCount": 156,
      "totalRevenue": 780000
    },
    {
      "itemId": 2,
      "itemCode": "GIFT_ROSE",
      "itemName": "玫瑰花",
      "subCategory": "gift",
      "orderCount": 89,
      "totalRevenue": 890000
    }
  ]
}

# Top陪玩师排行
GET /api/v1/admin/stats/top-players?month=2024-11&limit=10

# 礼物销售统计
GET /api/v1/admin/stats/gift-stats
Response:
{
  "gifts": [
    {
      "giftId": 2,
      "giftName": "玫瑰花",
      "totalSent": 450,           # 赠送总数
      "totalRevenue": 4500000     # 总收入（分）
    }
  ]
}

# 按游戏统计收入
GET /api/v1/admin/stats/revenue-by-game
Response:
{
  "games": [
    {
      "gameId": 1,
      "orderCount": 856,
      "totalRevenue": 8560000
    }
  ]
}
```

---

## 📋 完整Admin API列表

### 用户管理

```
GET    /admin/users             # 用户列表
GET    /admin/users/:id         # 用户详情
POST   /admin/users             # 创建用户
PUT    /admin/users/:id         # 更新用户
DELETE /admin/users/:id         # 删除用户
PUT    /admin/users/:id/status  # 更新状态
PUT    /admin/users/:id/balance # 更新余额
```

### 陪玩师管理

```
GET    /admin/players              # 陪玩师列表
GET    /admin/players/:id          # 陪玩师详情
POST   /admin/players              # 创建陪玩师
PUT    /admin/players/:id          # 更新陪玩师
DELETE /admin/players/:id          # 删除陪玩师
PUT    /admin/players/:id/verify   # 审核通过
PUT    /admin/players/:id/reject   # 审核拒绝
```

### 游戏管理

```
GET    /admin/games        # 游戏列表
GET    /admin/games/:id    # 游戏详情
POST   /admin/games        # 创建游戏
PUT    /admin/games/:id    # 更新游戏
DELETE /admin/games/:id    # 删除游戏
```

### 订单管理

```
GET    /admin/orders                 # 订单列表
GET    /admin/orders/:id             # 订单详情
POST   /admin/orders                 # 创建订单
PUT    /admin/orders/:id             # 更新订单
DELETE /admin/orders/:id             # 删除订单
POST   /admin/orders/:id/assign      # 指派陪玩师
PUT    /admin/orders/:id/status      # 更新状态
POST   /admin/orders/:id/refund      # 退款
```

### 服务项目管理 ⭐ 新增

```
GET    /admin/service-items                    # 服务列表（护航+礼物）
GET    /admin/service-items/:id                # 服务详情
POST   /admin/service-items                    # 创建服务
PUT    /admin/service-items/:id                # 更新服务
DELETE /admin/service-items/:id                # 删除服务
POST   /admin/service-items/batch-update-status # 批量启用/禁用
POST   /admin/service-items/batch-update-price  # 批量调价
```

### 抽成管理 ⭐ 新增

```
POST /admin/commission/rules                   # 创建抽成规则
PUT  /admin/commission/rules/:id               # 更新抽成规则
POST /admin/commission/settlements/trigger     # 手动触发结算
GET  /admin/commission/stats                   # 平台统计
```

### 提现审核 ⭐ 新增

```
GET  /admin/withdraws             # 提现列表
GET  /admin/withdraws/:id         # 提现详情
POST /admin/withdraws/:id/approve # 批准提现
POST /admin/withdraws/:id/reject  # 拒绝提现
POST /admin/withdraws/:id/complete # 完成提现（已打款）
```

### Dashboard ⭐ 新增

```
GET /admin/dashboard/overview         # 总览统计
GET /admin/dashboard/recent-orders    # 最近订单
GET /admin/dashboard/recent-withdraws # 最近提现
GET /admin/dashboard/monthly-revenue  # 月度收入趋势
```

### 数据统计 ⭐ 新增

```
GET /admin/stats/service-items   # 服务项目销售统计
GET /admin/stats/top-players     # Top陪玩师排行
GET /admin/stats/gift-stats      # 礼物销售统计
GET /admin/stats/revenue-by-game # 按游戏统计收入
```

### 支付管理

```
GET /admin/payments        # 支付列表
GET /admin/payments/:id    # 支付详情
```

### 评价管理

```
GET    /admin/reviews       # 评价列表
GET    /admin/reviews/:id   # 评价详情
DELETE /admin/reviews/:id   # 删除评价
```

### RBAC权限管理

```
GET    /admin/roles              # 角色列表
POST   /admin/roles              # 创建角色
PUT    /admin/roles/:id          # 更新角色
DELETE /admin/roles/:id          # 删除角色
PUT    /admin/roles/:id/permissions # 分配权限

GET    /admin/permissions        # 权限列表
POST   /admin/permissions        # 创建权限
PUT    /admin/permissions/:id    # 更新权限
DELETE /admin/permissions/:id    # 删除权限
```

### 系统管理

```
GET /admin/system/health    # 系统健康检查
GET /admin/system/info      # 系统信息
GET /admin/system/stats     # 系统统计
```

---

## 🔥 核心业务流程

### 服务项目配置流程

```
1. 管理员登录
POST /api/v1/auth/login

2. 创建护航服务
POST /admin/service-items
{
  "subCategory": "solo",
  "gameId": 1,
  "basePriceCents": 50000,
  "serviceHours": 1
}

3. 创建礼物
POST /admin/service-items
{
  "subCategory": "gift",
  "basePriceCents": 10000,
  "serviceHours": 0
}

4. 批量调价
POST /admin/service-items/batch-update-price
{
  "ids": [1, 2, 3],
  "basePriceCents": 12000
}

5. 查看销售统计
GET /admin/stats/service-items
```

---

### 提现审核流程

```
1. 查看待审核提现
GET /admin/withdraws?status=pending

2. 查看提现详情
GET /admin/withdraws/123

3. 批准提现
POST /admin/withdraws/123/approve
{
  "remark": "账户信息已确认"
}

4. 打款后标记完成
POST /admin/withdraws/123/complete

状态变化：
pending → approved → completed
```

---

### 月度结算流程

```
# 自动结算（每月1号凌晨2点）
Cron自动执行

# 手动补偿结算
POST /admin/commission/settlements/trigger?month=2024-10

# 查看结算统计
GET /admin/commission/stats?month=2024-11

# 查看月度收入趋势
GET /admin/dashboard/monthly-revenue?months=12
```

---

## 📊 Dashboard数据展示建议

### 首页卡片

```tsx
<Row gutter={16}>
  <Col span={6}>
    <Card>
      <Statistic title="总用户数" value={stats.totalUsers} />
    </Card>
  </Col>
  <Col span={6}>
    <Card>
      <Statistic title="总陪玩师" value={stats.totalPlayers} />
    </Card>
  </Col>
  <Col span={6}>
    <Card>
      <Statistic title="今日订单" value={stats.todayOrders} />
    </Card>
  </Col>
  <Col span={6}>
    <Card>
      <Statistic 
        title="今日收入" 
        value={stats.todayRevenue / 100} 
        prefix="¥"
      />
    </Card>
  </Col>
</Row>
```

### 月度收入趋势图

```tsx
<Card title="月度收入趋势">
  <Line
    data={{
      labels: revenue.map(r => r.month),
      datasets: [
        {
          label: '总收入',
          data: revenue.map(r => r.totalRevenue / 100)
        },
        {
          label: '平台抽成',
          data: revenue.map(r => r.totalCommission / 100)
        }
      ]
    }}
  />
</Card>
```

### 待处理事项

```tsx
<Card title="待处理">
  <List>
    <List.Item>
      <Badge count={stats.pendingWithdraws} />
      <span>待审核提现</span>
      <Button onClick={() => navigate('/admin/withdraws')}>
        去处理
      </Button>
    </List.Item>
  </List>
</Card>
```

---

## 🔐 权限控制

### API权限要求

```
所有 /admin/* 接口：
✅ 需要JWT认证
✅ 需要admin或super_admin角色
✅ 部分接口需要细粒度权限

示例：
GET /admin/withdraws
  需要: 认证 + admin角色 + "withdraw:read" 权限
  
POST /admin/withdraws/:id/approve
  需要: 认证 + admin角色 + "withdraw:approve" 权限
```

---

## 📈 数据统计维度

### 1. 服务项目维度

```
- 各服务项目的销售量
- 各服务项目的收入
- 护航 vs 礼物收入对比
- 热门服务Top10
```

### 2. 游戏维度

```
- 各游戏的订单量
- 各游戏的收入
- 游戏活跃度排名
```

### 3. 时间维度

```
- 今日/本周/本月数据
- 同比环比增长
- 月度收入趋势
- 高峰时段分析
```

### 4. 陪玩师维度

```
- Top陪玩师排行
- 新增陪玩师趋势
- 陪玩师活跃度
- 提现统计
```

---

## 🎯 使用场景

### 场景1: 配置新游戏的服务

```
1. 创建游戏
POST /admin/games

2. 为游戏创建段位护航服务
POST /admin/service-items
{
  "subCategory": "solo",
  "gameId": 新游戏ID,
  "rankLevel": "钻石",
  "basePriceCents": 50000,
  "serviceHours": 1
}

3. 为游戏创建团队服务
POST /admin/service-items
{
  "subCategory": "team",
  "gameId": 新游戏ID,
  "basePriceCents": 80000,
  "serviceHours": 2
}

4. 激活服务
PUT /admin/service-items/:id
{
  "isActive": true
}
```

### 场景2: 礼物营销活动

```
1. 创建节日礼物
POST /admin/service-items
{
  "itemCode": "GIFT_VALENTINE_ROSE",
  "name": "情人节玫瑰",
  "subCategory": "gift",
  "basePriceCents": 9900,    # 促销价
  "commissionRate": 0.15      # 特殊抽成
}

2. 查看销售情况
GET /admin/stats/gift-stats

3. 活动结束后调价
POST /admin/service-items/batch-update-price
{
  "ids": [礼物ID],
  "basePriceCents": 10000
}
```

### 场景3: 财务结算

```
1. 查看本月收入
GET /admin/commission/stats?month=2024-11

2. 月初手动触发结算
POST /admin/commission/settlements/trigger?month=2024-10

3. 审核提现申请
GET /admin/withdraws?status=pending
POST /admin/withdraws/:id/approve

4. 打款后标记完成
POST /admin/withdraws/:id/complete
```

---

## 📊 报表导出（建议实现）

### 财务报表

```
# 月度财务报表
GET /admin/reports/monthly-finance?month=2024-11
导出：
- 总收入
- 平台抽成
- 陪玩师收入
- 提现金额
- 净利润

# 服务项目销售报表
GET /admin/reports/service-sales?startDate=2024-11-01&endDate=2024-11-30

# 陪玩师收入报表
GET /admin/reports/player-income?month=2024-11
```

---

## 🔍 查询参数说明

### 通用参数

```
page=1           # 页码（从1开始）
pageSize=20      # 每页数量（最大100）
sortBy=created   # 排序字段
sortOrder=desc   # 排序方向（asc/desc）
```

### 日期参数

```
month=2024-11         # 月份（YYYY-MM）
date=2024-11-02       # 日期（YYYY-MM-DD）
startDate=2024-11-01  # 开始日期
endDate=2024-11-30    # 结束日期
```

### 状态参数

```
status=pending        # 订单/提现状态
isActive=true         # 是否激活
```

---

## ✨ 新增API总结

### 核心业务API（7个路由组）

```
✅ Service Items   - 7个端点（统一管理）
✅ Commission      - 4个端点（抽成管理）
✅ Withdraw        - 5个端点（提现审核）
✅ Dashboard       - 4个端点（总览统计）
✅ Stats           - 4个端点（数据分析）
```

### 新增端点数

```
服务项目: 7个
抽成管理: 4个
提现审核: 5个
Dashboard: 4个
统计分析: 4个
----------
总计: 24个新端点
```

---

## 🚀 部署后管理员工作流

### 每日工作

```
1. 查看Dashboard
   GET /admin/dashboard/overview
   
2. 处理待审核提现
   GET /admin/withdraws?status=pending
   POST /admin/withdraws/:id/approve
   
3. 查看异常订单
   GET /admin/orders?status=in_progress
   
4. 查看今日数据
   - 今日订单数
   - 今日收入
   - 活跃陪玩师数
```

### 每周工作

```
1. 查看服务销售情况
   GET /admin/stats/service-items
   
2. 调整服务价格
   POST /admin/service-items/batch-update-price
   
3. 分析热门游戏
   GET /admin/stats/revenue-by-game
   
4. 查看Top陪玩师
   GET /admin/stats/top-players
```

### 每月工作

```
1. 查看月度收入
   GET /admin/commission/stats?month=上月
   
2. 检查月度结算
   GET /admin/dashboard/monthly-revenue
   
3. 生成财务报表
   导出月度数据
   
4. 规划下月运营
   基于数据调整策略
```

---

## 🎯 总结

### ✅ 完整的管理功能

- **用户管理** - 完整CRUD
- **陪玩师管理** - 包含审核
- **游戏管理** - 完整CRUD
- **订单管理** - 完整生命周期
- **服务项目管理** - 统一管理（护航+礼物）⭐
- **抽成管理** - 规则配置和结算 ⭐
- **提现审核** - 完整审批流程 ⭐
- **Dashboard** - 数据总览 ⭐
- **统计分析** - 多维度分析 ⭐
- **RBAC权限** - 细粒度控制
- **系统管理** - 健康检查

### 📊 新增价值

**运营效率提升：**
- ✅ 统一的服务项目管理（不再分散）
- ✅ 自动化财务结算（无需人工）
- ✅ 可视化数据Dashboard（一目了然）
- ✅ 便捷的提现审核（快速处理）

**决策支持：**
- ✅ 服务销售统计（了解热门）
- ✅ 收入趋势分析（把握走势）
- ✅ 游戏维度统计（优化资源）
- ✅ 陪玩师排行（激励优质）

---

**管理端API已完善！可以支持完整的平台运营！** 🎉✨

