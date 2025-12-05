# 陪玩师绩效与动态抽成模块设计文档

## 概述

陪玩师绩效与动态抽成模块是 GameLink 陪玩管理平台的核心激励系统,通过灵活的抽成规则配置和智能的最优抽成计算,实现对优秀陪玩师的激励。该模块支持三种抽成规则:全站排名抽成、个人签约抽成和服务项目抽成,并自动选择对陪玩师最有利的抽成比例。

**核心特性:**
- 多维度排名计算(订单量、评分、服务时长、满意度)
- 三种抽成规则灵活配置
- 最低抽成原则自动选择
- 抽成历史记录和影响分析
- 规则模拟和效果评估

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                      前端管理界面                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │排名抽成  │  │签约抽成  │  │项目抽成  │  │抽成历史│ │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │排名计算  │  │抽成模拟  │  │影响分析  │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
                        ↕ HTTP/REST API
┌─────────────────────────────────────────────────────────┐
│                      后端服务层                           │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  抽成规则服务     │  │  排名计算服务     │            │
│  │  - 规则管理       │  │  - 多维度评分     │            │
│  │  - 最低抽成选择   │  │  - 自动更新       │            │
│  └──────────────────┘  └──────────────────┘            │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  抽成计算服务     │  │  历史记录服务     │            │
│  │  - 订单抽成计算   │  │  - 记录查询       │            │
│  │  - 模拟计算       │  │  - 统计分析       │            │
│  └──────────────────┘  └──────────────────┘            │
└─────────────────────────────────────────────────────────┘
                        ↕ 数据访问
┌─────────────────────────────────────────────────────────┐
│                      数据存储层                           │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  PostgreSQL      │  │  Redis 缓存       │            │
│  │  - 抽成规则       │  │  - 排名数据       │            │
│  │  - 抽成历史       │  │  - 规则缓存       │            │
│  │  - 排名记录       │  │                   │            │
│  └──────────────────┘  └──────────────────┘            │
└─────────────────────────────────────────────────────────┘
```

## 组件和接口

### 前端组件

#### 1. RankingCommissionConfig (排名抽成配置)
```typescript
interface RankingCommissionConfigProps {
  onSave: (config: RankingRule[]) => Promise<void>;
}

interface RankingRule {
  rankStart: number;  // 排名起始
  rankEnd: number;    // 排名结束
  commissionRate: number; // 抽成比例
}
```

#### 2. ContractCommissionManager (签约抽成管理)
```typescript
interface ContractCommissionManagerProps {
  onCreateContract: (contract: ContractConfig) => Promise<void>;
}

interface ContractConfig {
  playerId: number;
  commissionRate: number;
  startDate: string;
  endDate: string;
  reason: string;
}
```

#### 3. ServiceCommissionConfig (服务项目抽成配置)
```typescript
interface ServiceCommissionConfigProps {
  services: ServiceItem[];
  onBatchUpdate: (updates: ServiceCommissionUpdate[]) => Promise<void>;
}

interface ServiceCommissionUpdate {
  serviceId: number;
  commissionRate: number;
}
```

#### 4. CommissionCalculator (抽成计算器)
```typescript
interface CommissionCalculatorProps {
  playerId: number;
  serviceId: number;
  orderAmount: number;
  onCalculate: () => Promise<CommissionResult>;
}

interface CommissionResult {
  rankingCommission: number;
  contractCommission: number;
  serviceCommission: number;
  finalCommission: number;
  selectedRule: 'ranking' | 'contract' | 'service';
}
```

### 后端 API 接口

#### 排名抽成 API
```go
// POST /api/v1/admin/commission/ranking-rules
type CreateRankingRuleRequest struct {
    RankStart      int     `json:"rankStart"`
    RankEnd        int     `json:"rankEnd"`
    CommissionRate float64 `json:"commissionRate"`
}

// GET /api/v1/admin/commission/ranking-rules
type RankingRuleResponse struct {
    ID             int     `json:"id"`
    RankStart      int     `json:"rankStart"`
    RankEnd        int     `json:"rankEnd"`
    CommissionRate float64 `json:"commissionRate"`
    IsActive       bool    `json:"isActive"`
}
```

#### 签约抽成 API
```go
// POST /api/v1/admin/commission/contracts
type CreateContractRequest struct {
    PlayerID       int     `json:"playerId"`
    CommissionRate float64 `json:"commissionRate"`
    StartDate      string  `json:"startDate"`
    EndDate        string  `json:"endDate"`
    Reason         string  `json:"reason"`
}

// GET /api/v1/admin/commission/contracts
type ContractResponse struct {
    ID             int     `json:"id"`
    PlayerID       int     `json:"playerId"`
    PlayerName     string  `json:"playerName"`
    CommissionRate float64 `json:"commissionRate"`
    StartDate      string  `json:"startDate"`
    EndDate        string  `json:"endDate"`
    IsActive       bool    `json:"isActive"`
}
```

#### 抽成计算 API
```go
// POST /api/v1/admin/commission/calculate
type CalculateCommissionRequest struct {
    PlayerID  int     `json:"playerId"`
    ServiceID int     `json:"serviceId"`
    Amount    float64 `json:"amount"`
}

type CalculateCommissionResponse struct {
    RankingCommission  float64 `json:"rankingCommission"`
    ContractCommission float64 `json:"contractCommission"`
    ServiceCommission  float64 `json:"serviceCommission"`
    FinalCommission    float64 `json:"finalCommission"`
    SelectedRule       string  `json:"selectedRule"`
}
```

## 数据模型

### 排名抽成规则表 (ranking_commission_rules)
```sql
CREATE TABLE ranking_commission_rules (
    id SERIAL PRIMARY KEY,
    rank_start INTEGER NOT NULL,
    rank_end INTEGER NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (rank_start <= rank_end),
    CHECK (commission_rate >= 0 AND commission_rate <= 100)
);
```

### 签约抽成表 (contract_commissions)
```sql
CREATE TABLE contract_commissions (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES users(id),
    commission_rate DECIMAL(5,2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (start_date <= end_date),
    CHECK (commission_rate >= 0 AND commission_rate <= 100)
);
```

### 陪玩师排名表 (player_rankings)
```sql
CREATE TABLE player_rankings (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES users(id),
    rank_date DATE NOT NULL,
    ranking INTEGER NOT NULL,
    order_count INTEGER DEFAULT 0,
    avg_rating DECIMAL(3,2) DEFAULT 0,
    service_hours DECIMAL(10,2) DEFAULT 0,
    satisfaction_score DECIMAL(5,2) DEFAULT 0,
    total_score DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (player_id, rank_date)
);
```

### 抽成历史表 (commission_history)
```sql
CREATE TABLE commission_history (
    id BIGSERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    player_id INTEGER NOT NULL REFERENCES users(id),
    order_amount DECIMAL(10,2) NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL,
    commission_amount DECIMAL(10,2) NOT NULL,
    rule_type VARCHAR(20) NOT NULL,
    rule_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_player_id (player_id),
    INDEX idx_order_id (order_id),
    INDEX idx_created_at (created_at)
);
```

## 正确性属性

*属性是指在系统所有有效执行中都应该成立的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性是人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性反思

**冗余分析:**
1. 属性 1.4 和 5.1 都测试排名计算的多维度性,合并为一个属性
2. 属性 4.2 和 4.3 都测试最低抽成原则,合并为一个属性
3. 其他属性测试不同功能层面,保持独立

### 正确性属性列表

**属性 1: 排名规则修改立即生效**
*对于任意*排名规则修改操作,修改后创建的订单应该使用新的排名抽成比例进行计算
**验证需求: 1.3**

**属性 2: 多维度排名计算完整性**
*对于任意*陪玩师数据集,排名计算应该综合考虑订单量、评分、服务时长和用户满意度四个维度,且每个维度都应该对最终排名产生影响
**验证需求: 1.4, 5.1**

**属性 3: 签约协议记录完整性**
*对于任意*签约协议创建操作,系统应该记录陪玩师ID、抽成比例、生效时间和结束时间,且所有字段都应该非空
**验证需求: 2.2**

**属性 4: 签约到期自动停用**
*对于任意*签约协议,当当前日期超过结束日期时,系统应该自动停用该协议并恢复使用其他抽成规则
**验证需求: 2.3**

**属性 5: 签约修改历史记录**
*对于任意*个人抽成修改操作,系统应该记录修改前后的抽成比例、修改时间和修改原因
**验证需求: 2.4**

**属性 6: 服务项目抽成必填**
*对于任意*新创建的服务项目,系统应该要求设置默认抽成比例,且抽成比例应该在0-100之间
**验证需求: 3.2**

**属性 7: 项目抽成修改立即生效**
*对于任意*服务项目抽成修改,修改后创建的该项目订单应该使用新的抽成比例
**验证需求: 3.3**

**属性 8: 批量设置正确应用**
*对于任意*批量抽成设置操作,所有匹配筛选条件的服务项目都应该被更新为指定的抽成比例
**验证需求: 3.4**

**属性 9: 三种抽成全部计算**
*对于任意*订单创建操作,系统应该计算全站排名抽成、个人签约抽成和服务项目抽成三个值,且三个值都应该非负
**验证需求: 4.1**

**属性 10: 最低抽成原则**
*对于任意*三个抽成比例值,系统应该选择数值最小的一个作为最终抽成比例,且最终抽成应该小于或等于其他两个抽成
**验证需求: 4.2, 4.3**

**属性 11: 订单抽成不可变**
*对于任意*已创建的订单,其抽成比例应该在订单创建时确定并保持不变,即使后续抽成规则发生变化
**验证需求: 4.4**

**属性 12: 排名权重影响结果**
*对于任意*排名权重配置变更,变更后计算的排名应该反映新的权重分配,且权重更高的指标应该对排名产生更大影响
**验证需求: 5.2**

**属性 13: 排名每日自动更新**
*对于任意*日期,系统应该在该日期自动触发排名计算,且所有陪玩师的排名都应该被重新计算
**验证需求: 5.3**

**属性 14: 排名变更历史记录**
*对于任意*陪玩师排名变化,系统应该记录变更前排名、变更后排名和变更时间
**验证需求: 5.5**

**属性 15: 抽成历史筛选正确性**
*对于任意*抽成历史记录集合和筛选条件,筛选结果中的所有记录都应该满足指定的筛选条件(陪玩师/时间/类型)
**验证需求: 6.2**

**属性 16: 抽成统计计算准确性**
*对于任意*抽成历史记录集合,总抽成金额应该等于所有记录的抽成金额之和,平均抽成比例应该等于所有记录抽成比例的算术平均值
**验证需求: 6.3**

**属性 17: 订单结算记录抽成**
*对于任意*订单结算操作,系统应该在抽成历史表中创建一条记录,包含订单ID、抽成比例、抽成金额和规则类型
**验证需求: 6.5**

**属性 18: 规则自动启用**
*对于任意*抽成规则,当当前时间达到生效开始时间时,系统应该自动将规则状态设置为启用
**验证需求: 7.2**

**属性 19: 规则自动停用**
*对于任意*抽成规则,当当前时间达到结束时间时,系统应该自动将规则状态设置为停用
**验证需求: 7.3**

**属性 20: 多规则最低抽成选择**
*对于任意*同时生效的多个抽成规则,系统应该选择抽成比例最低的规则应用到订单
**验证需求: 7.4**

**属性 21: 模拟计算完整性**
*对于任意*抽成模拟请求,系统应该返回三种抽成规则的计算结果和最终选择的规则,且最终抽成应该是三者中的最小值
**验证需求: 8.2**

**属性 22: 批量模拟正确处理**
*对于任意*陪玩师列表,批量模拟应该为每个陪玩师计算抽成,且结果数量应该等于输入的陪玩师数量
**验证需求: 8.3**

**属性 23: 模拟结果信息完整**
*对于任意*模拟计算结果,应该包含陪玩师ID、抽成比例和预估抽成金额,且所有字段都应该非空
**验证需求: 8.4**

**属性 24: 影响分析效果计算**
*对于任意*抽成规则变更,影响分析应该计算变更前后的陪玩师活跃度和用户满意度差异
**验证需求: 9.2**

**属性 25: 规则对比差异展示**
*对于任意*两个不同的抽成规则,对比分析应该展示它们在订单量、收入和陪玩师满意度等指标上的差异
**验证需求: 9.3**

**属性 26: 受益陪玩师识别准确性**
*对于任意*抽成规则调整,系统应该识别出所有因调整而降低抽成比例的陪玩师,且识别结果应该准确无遗漏
**验证需求: 9.4**

**属性 27: 优先级规则选择**
*对于任意*规则集合,当最低抽成原则应用后,系统应该在抽成比例相同的规则中选择优先级更高的规则
**验证需求: 10.2, 10.3**

**属性 28: 优先级调整立即生效**
*对于任意*规则优先级调整,调整后的规则选择逻辑应该立即应用新的优先级顺序
**验证需求: 10.4**

## 错误处理

### 前端错误处理
1. API 请求失败 - 显示错误提示并提供重试
2. 数据验证失败 - 显示具体的验证错误信息
3. 权限不足 - 引导用户联系管理员

### 后端错误处理
1. 规则冲突 - 返回冲突详情并拒绝操作
2. 数据不一致 - 记录错误日志并回滚事务
3. 计算异常 - 使用默认抽成比例并记录异常

## 测试策略

### 单元测试
- 抽成计算逻辑测试
- 排名计算算法测试
- 规则选择逻辑测试

### 属性测试
使用 fast-check (前端) 和 gopter (后端),每个测试至少 100 次迭代,覆盖所有 28 个正确性属性。

### 集成测试
- 订单创建时抽成计算流程
- 规则修改后的抽成变化
- 排名更新触发抽成调整

### 端到端测试
- 完整的抽成规则配置流程
- 订单创建和结算流程
- 抽成历史查询和导出流程
