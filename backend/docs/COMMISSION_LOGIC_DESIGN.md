# 💰 GameLink 抽成逻辑设计文档

## 🎯 核心原则

**三层抽成，取最低值（对陪玩师最优惠）**

```
实际抽成 = MIN(服务项目抽成, 陪玩师专属抽成, 排名抽成)
```

---

## 📊 三层抽成详解

### 第1层：服务项目抽成

**来源**: `service_items.commission_rate`  
**配置**: 管理员为每个服务项目单独设置  
**范围**: 0-100%  

**示例：**
```json
{
  "itemCode": "ESCORT_RANK_DIAMOND",
  "name": "钻石段位护航",
  "basePriceCents": 50000,
  "commissionRate": 0.20  // 20%抽成
}

{
  "itemCode": "ESCORT_GIFT_ROSE",
  "name": "玫瑰花",
  "basePriceCents": 10000,
  "commissionRate": 0.20  // 礼物固定20%，不享受排名优惠
}
```

---

### 第2层：陪玩师专属抽成

**来源**: `commission_rules WHERE player_id = ?`  
**配置**: 管理员为特定陪玩师设置优惠抽成  
**用途**: 奖励优质陪玩师  

**示例：**
```json
{
  "name": "金牌陪玩师优惠",
  "type": "special",
  "playerId": 5,
  "rate": 15,  // 该陪玩师所有订单只抽15%
  "isActive": true
}
```

**管理员操作：**
```bash
POST /admin/commission/rules
{
  "name": "张三专属优惠",
  "type": "special",
  "playerId": 5,
  "rate": 15
}
```

---

### 第3层：排名抽成（管理员自定义）

**来源**: `ranking_commission_configs`  
**配置**: 管理员自定义排名抽成规则  
**计算**: 基于上月排名（单量 OR 金额）  
**排除**: 礼物订单不计入排名，礼物订单不享受排名优惠  

#### 配置方式A：统一抽成

```json
{
  "name": "2024年11月单量前8名优惠",
  "rankingType": "order_count",  // 按单量排名
  "month": "2024-12",            // 应用于12月的订单
  "rulesJson": "[
    {
      \"rankStart\": 1,
      \"rankEnd\": 8,
      \"commissionRate\": 10     // 前8名统一10%抽成
    }
  ]",
  "isActive": true
}
```

**效果：** 上月单量前8名的陪玩师，本月所有订单只抽10%

#### 配置方式B：阶梯抽成

```json
{
  "name": "2024年11月金额排名阶梯优惠",
  "rankingType": "income",       // 按金额排名
  "month": "2024-12",
  "rulesJson": "[
    {
      \"rankStart\": 1,
      \"rankEnd\": 3,
      \"commissionRate\": 10     // 前3名10%抽成
    },
    {
      \"rankStart\": 4,
      \"rankEnd\": 8,
      \"commissionRate\": 12     // 4-8名12%抽成
    },
    {
      \"rankStart\": 9,
      \"rankEnd\": 20,
      \"commissionRate\": 15     // 9-20名15%抽成
    }
  ]",
  "isActive": true
}
```

**效果：** 根据陪玩师上月金额排名，给予不同的抽成优惠

---

## 🔢 抽成计算示例

### 示例1：普通护航订单

**条件：**
- 服务项目抽成：20%
- 陪玩师专属抽成：无
- 排名抽成：无（陪玩师上月未进前20）

**计算：**
```
候选抽成: [20%]
实际抽成 = MIN(20%) = 20%

订单金额: 50000分 (500元)
平台抽成: 10000分 (100元) = 500 × 20%
陪玩师收入: 40000分 (400元) = 500 × 80%
```

---

### 示例2：护航订单 + 陪玩师专属优惠

**条件：**
- 服务项目抽成：20%
- 陪玩师专属抽成：15%（金牌陪玩师）
- 排名抽成：无

**计算：**
```
候选抽成: [20%, 15%]
实际抽成 = MIN(20%, 15%) = 15% ✅

订单金额: 50000分
平台抽成: 7500分 = 500 × 15%
陪玩师收入: 42500分 = 500 × 85%
```

---

### 示例3：护航订单 + 排名优惠（单量前5）

**条件：**
- 服务项目抽成：20%
- 陪玩师专属抽成：无
- 排名抽成：10%（上月单量第5名，前8名10%抽成）

**计算：**
```
候选抽成: [20%, 10%]
实际抽成 = MIN(20%, 10%) = 10% ✅

订单金额: 50000分
平台抽成: 5000分 = 500 × 10%
陪玩师收入: 45000分 = 500 × 90%
```

---

### 示例4：护航订单 + 多重优惠（取最低）

**条件：**
- 服务项目抽成：20%
- 陪玩师专属抽成：15%（金牌陪玩师）
- 排名抽成：12%（上月金额第7名，4-8名12%抽成）

**计算：**
```
候选抽成: [20%, 15%, 12%]
实际抽成 = MIN(20%, 15%, 12%) = 12% ✅

订单金额: 50000分
平台抽成: 6000分 = 500 × 12%
陪玩师收入: 44000分 = 500 × 88%
```

---

### 示例5：礼物订单（固定抽成）

**条件：**
- 服务项目抽成：20%（礼物固定）
- 陪玩师专属抽成：15%（即使有专属优惠）
- 排名抽成：10%（即使有排名）

**计算：**
```
礼物订单不参与排名优惠！
候选抽成: [20%, 15%]  // 不包括排名抽成
实际抽成 = MIN(20%, 15%) = 15% ✅

礼物金额: 30000分 (3个玫瑰 × 100元)
平台抽成: 4500分
陪玩师收入: 25500分
```

**注意：** 礼物可以享受陪玩师专属抽成，但不享受排名优惠

---

## 🏆 排名计算规则

### 排名统计范围

**包含：**
- ✅ 单人护航订单 (sub_category = 'solo')
- ✅ 团队护航订单 (sub_category = 'team')

**排除：**
- ❌ 礼物订单 (sub_category = 'gift')

**代码实现：**
```go
// 计算排名时过滤
for _, order := range orders {
    if order.IsGiftOrder() {
        continue  // 跳过礼物订单
    }
    // 统计单量和金额
}
```

---

### 排名维度

#### 1. 按单量排名

```sql
-- 统计上月订单数（排除礼物）
SELECT 
    cr.player_id,
    COUNT(*) as order_count
FROM commission_records cr
JOIN orders o ON cr.order_id = o.id
JOIN service_items si ON o.item_id = si.id
WHERE cr.settlement_month = '2024-11'
  AND si.sub_category IN ('solo', 'team')  -- 排除礼物
GROUP BY cr.player_id
ORDER BY order_count DESC
LIMIT 20
```

#### 2. 按金额排名

```sql
-- 统计上月收入（排除礼物）
SELECT 
    cr.player_id,
    SUM(cr.total_amount_cents) as total_income
FROM commission_records cr
JOIN orders o ON cr.order_id = o.id
JOIN service_items si ON o.item_id = si.id
WHERE cr.settlement_month = '2024-11'
  AND si.sub_category IN ('solo', 'team')  -- 排除礼物
GROUP BY cr.player_id
ORDER BY total_income DESC
LIMIT 20
```

---

## 📅 时间轴

### 月度循环

```
11月 1-30日
  ├── 陪玩师接单服务
  ├── 订单完成自动记录抽成
  └── 抽成使用11月的排名配置
      （基于10月的排名）

12月 1日凌晨 2:00
  ├── 自动结算11月数据
  ├── 计算11月排名
  │   ├── 单量排名（排除礼物）
  │   └── 金额排名（排除礼物）
  └── 生成 PlayerRanking 记录

12月 1-30日
  ├── 陪玩师接单服务
  ├── 订单完成自动记录抽成
  └── 抽成使用12月的排名配置
      （基于11月的排名）✅

1月 1日凌晨 2:00
  ├── 自动结算12月数据
  └── 计算12月排名
```

---

## 💳 提现模式

### 提现时机

```
每周: ✅ 陪玩师可随时申请提现
      - 不限制次数
      - 最低金额100元

月末: ✅ 自动结算所有已完成订单
      - 每月1号凌晨2点
      - 生成月度结算报表
```

### 可提现余额计算

```
可提现余额 = 累计收入 - 累计提现 - 待处理提现 - 待结算余额

其中：
- 累计收入 = SUM(已完成订单的 player_income_cents)
- 待结算余额 = SUM(进行中订单的 player_income_cents)
```

---

## 🎯 管理员配置示例

### 配置1：按单量排名，前8名10%抽成

```bash
POST /admin/ranking-commission/configs
{
  "name": "2024年12月单量前8名优惠",
  "rankingType": "order_count",
  "month": "2024-12",
  "rulesJson": "[{\"rankStart\":1,\"rankEnd\":8,\"commissionRate\":10}]",
  "description": "上月单量前8名的陪玩师，本月享受10%抽成"
}
```

### 配置2：按金额排名，阶梯抽成

```bash
POST /admin/ranking-commission/configs
{
  "name": "2024年12月金额排名阶梯优惠",
  "rankingType": "income",
  "month": "2024-12",
  "rulesJson": "[
    {\"rankStart\":1,\"rankEnd\":3,\"commissionRate\":10},
    {\"rankStart\":4,\"rankEnd\":8,\"commissionRate\":12},
    {\"rankStart\":9,\"rankEnd\":15,\"commissionRate\":15}
  ]",
  "description": "前3名10%，4-8名12%，9-15名15%"
}
```

### 配置3：双维度并行

**可以同时配置：**
- 单量排名配置（例如：前8名10%）
- 金额排名配置（例如：前5名10%）

**效果：** 如果陪玩师同时满足两个排名，取抽成更低的那个

**示例：**
```
陪玩师A:
- 单量排名: 第5名 → 10%抽成
- 金额排名: 第2名 → 10%抽成
→ 使用10%抽成

陪玩师B:
- 单量排名: 第15名 → 无优惠
- 金额排名: 第6名 → 12%抽成
→ 使用12%抽成
```

---

## 🎁 礼物订单特殊规则

### 礼物不参与排名

```
计算11月排名时：
✅ 统计: 护航订单（solo + team）
❌ 排除: 礼物订单（gift）

例如：
陪玩师A:
- 护航订单: 50单，收入50万分
- 礼物订单: 30单，收入30万分
→ 排名只看护航: 50单，50万分
```

### 礼物抽成计算

```
礼物订单抽成：
1. 服务项目抽成（礼物的commission_rate）
2. 陪玩师专属抽成（如果有）
3. 不考虑排名抽成 ❌

实际抽成 = MIN(服务项目抽成, 陪玩师专属抽成)
```

**示例：**
```
礼物订单（玫瑰花）:
- 服务项目抽成: 20%
- 陪玩师专属: 15%（金牌陪玩师）
- 排名优惠: 10%（上月第1名）❌ 不适用
→ 实际抽成 = MIN(20%, 15%) = 15%
```

---

## 📋 数据模型

### ranking_commission_configs（排名抽成配置）

```sql
CREATE TABLE ranking_commission_configs (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    ranking_type VARCHAR(32) NOT NULL,  -- 'income' or 'order_count'
    period VARCHAR(32) NOT NULL,         -- 'monthly'
    month VARCHAR(7) NOT NULL,           -- 'YYYY-MM' 应用于哪个月
    rules_json TEXT NOT NULL,            -- JSON序列化的规则
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    
    INDEX idx_type_month (ranking_type, month),
    INDEX idx_month_active (month, is_active)
);
```

### RulesJSON 格式

```json
[
  {
    "rankStart": 1,      // 排名开始
    "rankEnd": 3,        // 排名结束
    "commissionRate": 10 // 抽成比例（%）
  },
  {
    "rankStart": 4,
    "rankEnd": 8,
    "commissionRate": 12
  }
]
```

---

## 🔄 完整业务流程

### 11月运营流程

```
11月 1日
  ├── 管理员配置11月排名抽成规则
  │   POST /admin/ranking-commission/configs
  │   {
  │     "rankingType": "order_count",
  │     "month": "2024-11",
  │     "rulesJson": "[{\"rankStart\":1,\"rankEnd\":8,\"commissionRate\":10}]"
  │   }
  │
  └── 系统自动计算10月排名
      ├── 统计10月订单（排除礼物）
      ├── 计算单量排名
      └── 计算金额排名

11月 2-30日
  ├── 用户下单
  ├── 陪玩师服务
  ├── 订单完成
  └── 自动计算抽成
      ├── 查询陪玩师10月排名
      ├── 应用11月排名规则
      └── 取最低抽成

12月 1日 凌晨2:00
  ├── 自动结算11月数据
  │   └── 所有已完成订单的抽成记录
  │
  └── 自动计算11月排名
      ├── 统计11月订单（排除礼物）
      ├── 生成 PlayerRanking 记录
      └── 用于12月的抽成计算
```

---

## 🎯 管理后台配置界面建议

### 排名抽成配置页面

```tsx
<Form>
  <Form.Item label="配置名称">
    <Input placeholder="例如：2024年12月单量前8名优惠" />
  </Form.Item>

  <Form.Item label="排名维度">
    <Radio.Group>
      <Radio value="order_count">按单量排名</Radio>
      <Radio value="income">按金额排名</Radio>
    </Radio.Group>
  </Form.Item>

  <Form.Item label="应用月份">
    <DatePicker.MonthPicker />
  </Form.Item>

  <Form.Item label="抽成规则">
    <RankingRulesEditor />
  </Form.Item>
</Form>

// RankingRulesEditor 组件
<div>
  <Button onClick={addRule}>添加规则</Button>
  
  {rules.map((rule, index) => (
    <Row key={index}>
      <Col>
        排名 
        <InputNumber value={rule.rankStart} /> 
        - 
        <InputNumber value={rule.rankEnd} />
        名
      </Col>
      <Col>
        抽成 
        <InputNumber value={rule.commissionRate} min={0} max={100} />
        %
      </Col>
      <Col>
        <Button onClick={() => removeRule(index)}>删除</Button>
      </Col>
    </Row>
  ))}
</div>
```

---

## 📊 排名展示

### 陪玩师查看排名

```bash
GET /player/rankings?month=2024-11

Response:
{
  "rankings": [
    {
      "rankingType": "order_count",
      "period": "monthly",
      "month": "2024-11",
      "rank": 5,
      "orderCount": 156,
      "享受优惠": "本月订单10%抽成"
    },
    {
      "rankingType": "income",
      "period": "monthly",
      "month": "2024-11",
      "rank": 8,
      "totalIncome": 560000,
      "享受优惠": "本月订单12%抽成"
    }
  ]
}
```

---

## 🎯 API设计

### 管理端API

```bash
# 创建排名抽成配置
POST /admin/ranking-commission/configs
{
  "name": "12月单量前8名优惠",
  "rankingType": "order_count",
  "month": "2024-12",
  "rulesJson": "[{\"rankStart\":1,\"rankEnd\":8,\"commissionRate\":10}]"
}

# 获取配置列表
GET /admin/ranking-commission/configs?month=2024-12

# 更新配置
PUT /admin/ranking-commission/configs/:id

# 删除配置
DELETE /admin/ranking-commission/configs/:id

# 查看排名
GET /admin/rankings?month=2024-11&type=order_count

# 手动触发排名计算
POST /admin/rankings/calculate?month=2024-11
```

### 陪玩师端API

```bash
# 查看我的排名
GET /player/rankings?month=2024-11

# 查看排名抽成说明
GET /player/rankings/commission-info
```

---

## 💡 抽成计算算法

### 伪代码

```go
func CalculateCommission(order) {
    candidates := []

    // 1. 服务项目抽成
    serviceItem := GetServiceItem(order.ItemID)
    candidates.append({
        source: "服务项目",
        rate: serviceItem.CommissionRate * 100
    })

    // 2. 陪玩师专属抽成
    if playerRule := GetPlayerRule(order.PlayerID) {
        candidates.append({
            source: "陪玩师专属",
            rate: playerRule.Rate
        })
    }

    // 3. 排名抽成（礼物订单跳过）
    if !order.IsGiftOrder() {
        lastMonth := GetLastMonth()
        playerRanking := GetPlayerRanking(order.PlayerID, lastMonth)
        
        if playerRanking {
            rankingConfig := GetRankingConfig(month)
            rules := ParseJSON(rankingConfig.RulesJSON)
            
            for rule in rules {
                if playerRanking.Rank >= rule.RankStart 
                   && playerRanking.Rank <= rule.RankEnd {
                    candidates.append({
                        source: "排名优惠",
                        rate: rule.CommissionRate
                    })
                    break
                }
            }
        }
    }

    // 取最低抽成
    actualRate := MIN(candidates)
    
    // 计算金额
    commission = order.TotalPrice * actualRate / 100
    playerIncome = order.TotalPrice - commission
    
    return {commission, playerIncome, actualRate}
}
```

---

## ✅ 验证规则

### 配置验证

```
1. RankStart >= 1
2. RankEnd >= RankStart
3. CommissionRate 在 0-100之间
4. 规则范围不能重叠
5. RulesJSON 格式正确
```

### 示例：无效配置

```json
// ❌ 范围重叠
[
  {"rankStart": 1, "rankEnd": 5, "commissionRate": 10},
  {"rankStart": 3, "rankEnd": 8, "commissionRate": 12}  // 3-5重叠
]

// ❌ 抽成超范围
[
  {"rankStart": 1, "rankEnd": 3, "commissionRate": 150}  // 超过100%
]

// ✅ 正确配置
[
  {"rankStart": 1, "rankEnd": 3, "commissionRate": 10},
  {"rankStart": 4, "rankEnd": 8, "commissionRate": 12},
  {"rankStart": 9, "rankEnd": 15, "commissionRate": 15}
]
```

---

## 📊 数据示例

### 陪玩师月度数据

```
陪玩师ID: 5
11月数据:
├── 护航订单: 50单，收入50万分 → 计入排名
├── 礼物订单: 15单，收入15万分 → 不计入排名
└── 总收入: 65万分

排名结果:
├── 单量排名: 第5名（基于50单）
└── 金额排名: 第8名（基于50万分）

12月抽成:
├── 单量前8名 → 10%抽成
├── 金额4-8名 → 12%抽成
└── 实际使用: MIN(10%, 12%) = 10% ✅
```

---

## 🎊 总结

### 抽成计算优先级

```
1. 计算三种抽成
   ├── 服务项目抽成（必有）
   ├── 陪玩师专属抽成（可选）
   └── 排名抽成（可选，礼物订单不适用）

2. 取最低值
   → 对陪玩师最优惠

3. 礼物特殊处理
   → 不计入排名统计
   → 不享受排名优惠
   → 可享受陪玩师专属优惠
```

### 排名规则特点

```
✅ 管理员完全控制
✅ 支持统一抽成
✅ 支持阶梯抽成
✅ 可配置多个维度
✅ JSON序列化存储
✅ 灵活扩展
```

---

**这样设计是否符合您的需求？** 🎯


