# 客服管理与计费模块设计文档

## 概述

客服管理与计费模块是 GameLink 陪玩管理平台的客服支持系统,提供灵活的薪酬计算方案。该模块支持两种计费模式:单量模式(按工单数量计费)和抽成模式(按订单金额抽成),并提供完整的工单统计、绩效评估和薪酬结算功能。

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                      前端管理界面                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │计费配置  │  │工单统计  │  │绩效评估  │  │薪酬报表│ │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↕ HTTP/REST API
┌─────────────────────────────────────────────────────────┐
│                      后端服务层                           │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  计费管理服务     │  │  工单统计服务     │            │
│  │  - 模式配置       │  │  - 数量统计       │            │
│  │  - 单价设置       │  │  - 金额统计       │            │
│  └──────────────────┘  └──────────────────┘            │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  薪酬计算服务     │  │  绩效评估服务     │            │
│  │  - 单量计算       │  │  - 指标统计       │            │
│  │  - 抽成计算       │  │  - 满意度计算     │            │
│  │  - 保底处理       │  │  - 奖金计算       │            │
│  └──────────────────┘  └──────────────────┘            │
└─────────────────────────────────────────────────────────┘
```

## 组件和接口

### 前端组件

```typescript
// 计费模式配置
interface BillingModeConfigProps {
  staffId: number;
  currentMode: 'quantity' | 'commission';
  onSave: (config: BillingConfig) => Promise<void>;
}

interface BillingConfig {
  mode: 'quantity' | 'commission';
  quantityRate?: number;  // 单量模式单价
  commissionRate?: number; // 抽成模式比例
  minimumSalary?: number;  // 保底薪酬
}

// 工单类型单价配置
interface TicketPriceConfigProps {
  ticketTypes: TicketType[];
  onUpdate: (updates: TicketPriceUpdate[]) => Promise<void>;
}

interface TicketPriceUpdate {
  ticketTypeId: number;
  price: number;
}
```

### 后端 API 接口

```go
// 计费配置 API
type BillingConfigRequest struct {
    StaffID        int     `json:"staffId"`
    Mode           string  `json:"mode"` // quantity/commission
    QuantityRate   float64 `json:"quantityRate,omitempty"`
    CommissionRate float64 `json:"commissionRate,omitempty"`
    MinimumSalary  float64 `json:"minimumSalary,omitempty"`
}

// 薪酬计算 API
type SalaryCalculationResponse struct {
    StaffID         int     `json:"staffId"`
    Period          string  `json:"period"`
    TicketCount     int     `json:"ticketCount"`
    OrderAmount     float64 `json:"orderAmount"`
    ActualSalary    float64 `json:"actualSalary"`
    MinimumSalary   float64 `json:"minimumSalary"`
    FinalSalary     float64 `json:"finalSalary"`
    BonusAmount     float64 `json:"bonusAmount"`
    IsMinimumApplied bool   `json:"isMinimumApplied"`
}
```

## 数据模型

### 客服计费配置表 (staff_billing_configs)
```sql
CREATE TABLE staff_billing_configs (
    id SERIAL PRIMARY KEY,
    staff_id INTEGER NOT NULL REFERENCES users(id),
    billing_mode VARCHAR(20) NOT NULL,
    quantity_rate DECIMAL(10,2),
    commission_rate DECIMAL(5,2),
    minimum_salary DECIMAL(10,2),
    is_active BOOLEAN DEFAULT TRUE,
    effective_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 工单类型单价表 (ticket_type_prices)
```sql
CREATE TABLE ticket_type_prices (
    id SERIAL PRIMARY KEY,
    ticket_type VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 工单处理记录表 (ticket_records)
```sql
CREATE TABLE ticket_records (
    id BIGSERIAL PRIMARY KEY,
    ticket_id VARCHAR(50) UNIQUE NOT NULL,
    staff_id INTEGER NOT NULL REFERENCES users(id),
    ticket_type VARCHAR(50) NOT NULL,
    order_id INTEGER REFERENCES orders(id),
    order_amount DECIMAL(10,2),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_minutes INTEGER,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 薪酬记录表 (salary_records)
```sql
CREATE TABLE salary_records (
    id BIGSERIAL PRIMARY KEY,
    staff_id INTEGER NOT NULL REFERENCES users(id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    billing_mode VARCHAR(20) NOT NULL,
    ticket_count INTEGER DEFAULT 0,
    order_amount DECIMAL(10,2) DEFAULT 0,
    actual_salary DECIMAL(10,2) NOT NULL,
    minimum_salary DECIMAL(10,2),
    bonus_amount DECIMAL(10,2) DEFAULT 0,
    final_salary DECIMAL(10,2) NOT NULL,
    is_minimum_applied BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 正确性属性

*属性是指在系统所有有效执行中都应该成立的特征或行为——本质上是关于系统应该做什么的形式化陈述。*

### 属性反思
合并冗余属性:属性 8.3 和 8.4 都测试保底逻辑,合并为一个属性。

### 正确性属性列表

**属性 1: 单量模式单价必填**
*对于任意*选择单量模式的客服配置,单价字段应该非空且大于零
**验证需求: 1.2**

**属性 2: 抽成模式比例必填**
*对于任意*选择抽成模式的客服配置,抽成比例字段应该非空且在0-100之间
**验证需求: 1.3**

**属性 3: 计费模式修改历史记录**
*对于任意*计费模式修改操作,系统应该记录修改前模式、修改后模式和生效时间
**验证需求: 1.4**

**属性 4: 工单类型单价必填**
*对于任意*新创建的工单类型,单价字段应该非空且大于零
**验证需求: 2.2**

**属性 5: 工单单价自动匹配**
*对于任意*工单处理操作,系统应该根据工单类型自动应用对应的单价,且单价应该与配置表中的单价一致
**验证需求: 2.3**

**属性 6: 单价修改立即生效**
*对于任意*工单单价修改,修改后处理的工单应该使用新单价
**验证需求: 2.4**

**属性 7: 工单处理记录完整性**
*对于任意*工单完成操作,系统应该记录处理时间和工单类型,且两个字段都应该非空
**验证需求: 3.1**

**属性 8: 工单数量统计准确性**
*对于任意*时间周期和客服,工单数量统计应该等于该客服在该周期内完成的工单总数
**验证需求: 3.2**

**属性 9: 单量薪酬计算准确性**
*对于任意*工单集合,单量薪酬应该等于所有工单单价的总和
**验证需求: 3.3**

**属性 10: 工单统计筛选正确性**
*对于任意*工单记录集合和筛选条件,筛选结果中的所有记录都应该满足指定的筛选条件(客服/时间/类型)
**验证需求: 3.4**

**属性 11: 订单工单金额记录**
*对于任意*订单相关工单,系统应该记录关联的订单金额,且金额应该大于零
**验证需求: 4.1**

**属性 12: 订单金额统计准确性**
*对于任意*时间周期和客服,订单金额统计应该等于该客服在该周期内处理的所有订单金额之和
**验证需求: 4.2**

**属性 13: 抽成薪酬计算准确性**
*对于任意*订单总金额和抽成比例,抽成薪酬应该等于订单总金额乘以抽成比例
**验证需求: 4.3**

**属性 14: 无订单工单排除**
*对于任意*未关联订单的工单,该工单不应该被计入抽成薪酬计算
**验证需求: 4.4**

**属性 15: 薪酬报表内容完整性**
*对于任意*薪酬报表生成,报表应该包含工单数量、订单金额和应付薪酬三个字段,且所有字段都应该非空
**验证需求: 5.2**

**属性 16: 单量模式报表信息**
*对于任意*使用单量模式的客服,薪酬报表应该展示各类型工单数量和对应单价
**验证需求: 5.3**

**属性 17: 抽成模式报表信息**
*对于任意*使用抽成模式的客服,薪酬报表应该展示订单总金额和抽成比例
**验证需求: 5.4**

**属性 18: 满意度计算准确性**
*对于任意*用户评价集合,满意度应该等于所有评价分数的算术平均值
**验证需求: 6.2**

**属性 19: 客服绩效排序正确性**
*对于任意*客服集合和排序指标,排序结果应该按照指定指标降序排列
**验证需求: 6.4**

**属性 20: 低绩效标记准确性**
*对于任意*客服绩效数据,如果绩效低于设定阈值,该客服应该被标记为需要关注对象
**验证需求: 6.5**

**属性 21: 登录打卡记录**
*对于任意*客服登录事件,系统应该记录登录时间作为上班打卡
**验证需求: 7.2**

**属性 22: 退出打卡记录**
*对于任意*客服退出事件,系统应该记录退出时间作为下班打卡
**验证需求: 7.3**

**属性 23: 工作时长计算准确性**
*对于任意*客服的打卡记录,实际工作时长应该等于所有上班打卡和下班打卡时间差的总和
**验证需求: 7.4**

**属性 24: 工作时长不足检测**
*对于任意*客服的工作时长,如果实际工作时长小于应工作时长,系统应该标记为异常
**验证需求: 7.5**

**属性 25: 薪酬保底逻辑**
*对于任意*薪酬计算,如果实际薪酬低于保底薪酬,最终薪酬应该等于保底薪酬;如果实际薪酬高于保底薪酬,最终薪酬应该等于实际薪酬
**验证需求: 8.2, 8.3, 8.4**

**属性 26: 薪酬历史筛选正确性**
*对于任意*薪酬记录集合和筛选条件,筛选结果中的所有记录都应该满足指定的筛选条件(客服/时间/模式)
**验证需求: 9.2**

**属性 27: 薪酬统计计算准确性**
*对于任意*薪酬记录集合,总薪酬支出应该等于所有记录最终薪酬的总和,平均薪酬应该等于总薪酬除以记录数量
**验证需求: 9.3**

**属性 28: 薪酬结算历史记录**
*对于任意*薪酬结算操作,系统应该在薪酬历史表中创建一条记录,包含周期、薪酬明细和结算状态
**验证需求: 9.4**

**属性 29: 绩效奖金自动计算**
*对于任意*满足奖金条件的客服,系统应该自动计算奖金金额并加入最终薪酬
**验证需求: 10.2**

**属性 30: 奖金规则修改立即生效**
*对于任意*奖金规则修改,修改后的薪酬计算应该使用新的奖金规则
**验证需求: 10.4**

## 错误处理

1. 计费模式配置冲突 - 验证单量和抽成参数的互斥性
2. 工单单价为负 - 拒绝并提示错误
3. 薪酬计算异常 - 记录错误日志并使用保底薪酬
4. 数据不一致 - 回滚事务并通知管理员

## 测试策略

### 单元测试
- 薪酬计算逻辑测试
- 保底薪酬处理测试
- 工单统计计算测试

### 属性测试
使用 fast-check (前端) 和 gopter (后端),每个测试至少 100 次迭代,覆盖所有 30 个正确性属性。

### 集成测试
- 工单处理到薪酬计算流程
- 计费模式切换后的薪酬变化
- 薪酬结算和历史记录流程
