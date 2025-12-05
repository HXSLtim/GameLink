# 设计文档 - 佣金与结算模块

## 概述

佣金与结算模块是 GameLink 平台的核心财务结算系统，负责管理平台抽成规则、陪玩师收益结算、提现审核和打款操作。该模块采用前后端分离架构，确保结算数据的准确性和及时性，维护平台与陪玩师之间的财务透明度。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层 (React)                        │
├─────────────────────────────────────────────────────────────┤
│  佣金规则  │  结算管理  │  提现审核  │  打款管理  │  统计报表  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                      API 层 (Go Backend)                     │
├─────────────────────────────────────────────────────────────┤
│  佣金API  │  结算API  │  提现API  │  打款API  │  统计API      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层                              │
├─────────────────────────────────────────────────────────────┤
│  佣金服务  │  结算服务  │  提现服务  │  打款服务  │  统计服务  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite  │  Redis (缓存)  │  MongoDB (日志)       │
│  (结算数据)         │                │                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      外部服务层                              │
├─────────────────────────────────────────────────────────────┤
│  支付宝API  │  微信支付API  │  银行接口  │  定时任务         │
└─────────────────────────────────────────────────────────────┘
```

## 组件和接口

### 前端组件

#### 1. CommissionRuleList 组件
```typescript
interface CommissionRuleListProps {
  onAdd: () => void;
  onEdit: (ruleId: number) => void;
  onDelete: (ruleId: number) => void;
  onToggleStatus: (ruleId: number) => void;
}
```

#### 2. SettlementManager 组件
```typescript
interface SettlementManagerProps {
  onTrigger: (startDate: Date, endDate: Date) => Promise<void>;
  onViewDetail: (settlementId: number) => void;
}
```

#### 3. WithdrawalReview 组件
```typescript
interface WithdrawalReviewProps {
  onApprove: (withdrawalId: number) => Promise<void>;
  onReject: (withdrawalId: number, reason: string) => Promise<void>;
}
```

### 后端接口

```go
type CommissionService interface {
    CreateRule(ctx context.Context, req *CreateRuleRequest) (*CommissionRule, error)
    UpdateRule(ctx context.Context, ruleID int64, req *UpdateRuleRequest) error
    DeleteRule(ctx context.Context, ruleID int64) error
    ListRules(ctx context.Context) ([]*CommissionRule, error)
    CalculateCommission(ctx context.Context, orderAmount float64, gameID int64) (float64, error)
}

type SettlementService interface {
    TriggerSettlement(ctx context.Context, req *SettlementRequest) (*SettlementResult, error)
    ListSettlements(ctx context.Context, req *ListSettlementsRequest) (*ListSettlementsResponse, error)
    GetSettlement(ctx context.Context, settlementID int64) (*Settlement, error)
}

type WithdrawalService interface {
    ListWithdrawals(ctx context.Context, req *ListWithdrawalsRequest) (*ListWithdrawalsResponse, error)
    ApproveWithdrawal(ctx context.Context, withdrawalID int64) error
    RejectWithdrawal(ctx context.Context, withdrawalID int64, reason string) error
    ProcessPayout(ctx context.Context, withdrawalIDs []int64) error
}
```

## 数据模型

### CommissionRule (佣金规则)
```typescript
interface CommissionRule {
  id: number;
  name: string;
  gameID?: number;
  minAmount: number;
  maxAmount?: number;
  rate: number;              // 抽成比例 (0-100)
  description?: string;
  isActive: boolean;
  priority: number;
  createdAt: Date;
  updatedAt: Date;
}
```

### Settlement (结算记录)
```typescript
interface Settlement {
  id: number;
  playerID: number;
  startDate: Date;
  endDate: Date;
  totalAmount: number;       // 总金额
  commission: number;        // 平台佣金
  netAmount: number;         // 净收入
  orderCount: number;        // 订单数量
  status: SettlementStatus;
  createdAt: Date;
  processedAt?: Date;
}

enum SettlementStatus {
  PENDING = 'pending',
  PROCESSED = 'processed',
  PAID = 'paid'
}
```

### Withdrawal (提现记录)
```typescript
interface Withdrawal {
  id: number;
  playerID: number;
  amount: number;
  method: WithdrawalMethod;
  accountInfo: string;
  status: WithdrawalStatus;
  rejectReason?: string;
  createdAt: Date;
  approvedAt?: Date;
  paidAt?: Date;
}

enum WithdrawalMethod {
  ALIPAY = 'alipay',
  WECHAT = 'wechat',
  BANK_CARD = 'bank_card'
}

enum WithdrawalStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  PROCESSING = 'processing',
  COMPLETED = 'completed',
  FAILED = 'failed'
}
```

## 正确性属性

### 属性 1：佣金比例范围约束
*对于任何*佣金规则，抽成比例必须在0到100之间（包含0和100）
**验证：需求 1.3**

### 属性 2：金额范围合理性
*对于任何*佣金规则，如果设置了最高金额，则最高金额必须大于最低金额
**验证：需求 1.3**

### 属性 3：佣金计算准确性
*对于任何*订单，计算的佣金金额必须不超过订单金额
**验证：需求 9.3**

### 属性 4：净收入计算一致性
*对于任何*结算记录，净收入必须等于总金额减去佣金
**验证：需求 9.4**

### 属性 5：结算金额汇总正确性
*对于任何*结算操作，总结算金额必须等于所有订单净收入之和
**验证：需求 9.5**

### 属性 6：提现金额约束
*对于任何*提现申请，提现金额必须大于0且不超过陪玩师可用余额
**验证：需求 6.5**

### 属性 7：结算状态转换合法性
*对于任何*结算状态更新，状态转换必须遵循合法路径：pending → processed → paid
**验证：需求 8.5**

### 属性 8：提现状态转换合法性
*对于任何*提现状态更新，状态转换必须遵循合法路径：pending → approved → processing → completed 或 pending → rejected
**验证：需求 6.3, 6.4, 7.4, 7.5**

### 属性 9：权限验证一致性
*对于任何*需要权限的操作，系统必须先验证用户权限，验证失败时必须返回403错误且不执行任何业务逻辑
**验证：需求 10.1, 10.2, 10.3, 10.4**

### 属性 10：批量打款原子性
*对于任何*批量打款操作，要么所有提现都成功打款，要么所有提现都保持原状态
**验证：需求 7.2**

## 错误处理

### 错误分类
- **400 Bad Request**: 请求参数错误（佣金比例超出范围、金额范围不合理等）
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 资源不存在
- **500 Internal Server Error**: 服务器内部错误
- **502 Bad Gateway**: 第三方支付服务错误

## 测试策略

### 单元测试
使用 Vitest + Testing Library 进行组件和工具函数测试

### 属性测试
使用 fast-check 验证系统的通用正确性属性，每个属性测试运行至少100次迭代

### 集成测试
测试完整结算流程、提现审核流程、打款流程

### E2E测试
使用 Playwright 进行端到端测试
