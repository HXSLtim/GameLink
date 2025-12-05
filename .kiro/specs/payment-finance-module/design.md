# 设计文档 - 支付与财务模块

## 概述

支付与财务模块是 GameLink 平台的核心财务管理系统，负责处理所有与资金相关的业务流程。该模块采用前后端分离架构，前端基于 React + TypeScript + Ant Design 构建，后端通过 RESTful API 提供服务。模块设计遵循高内聚低耦合原则，确保系统的可维护性和可扩展性。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层 (React)                        │
├─────────────────────────────────────────────────────────────┤
│  支付记录页面  │  退款管理页面  │  财务对账页面  │  报表页面  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                      API 网关层 (Go)                         │
├─────────────────────────────────────────────────────────────┤
│  路由管理  │  权限验证  │  请求限流  │  日志记录  │  错误处理  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层 (Go)                         │
├─────────────────────────────────────────────────────────────┤
│  支付服务  │  退款服务  │  对账服务  │  报表服务  │  配置服务  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据访问层 (Go)                         │
├─────────────────────────────────────────────────────────────┤
│  支付仓储  │  流水仓储  │  对账仓储  │  配置仓储  │  日志仓储  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite  │  Redis (缓存)  │  MongoDB (日志)       │
│  (主数据)           │                │                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      外部服务层                              │
├─────────────────────────────────────────────────────────────┤
│  支付宝 API  │  微信支付 API  │  银行接口  │  消息队列       │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈

**前端**：
- React 19.2.0
- TypeScript
- Ant Design 6.0.0 (UI组件库)
- Axios 1.13.2
- Recharts 3.5.0

**后端**：
- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- Redis
- PostgreSQL (生产环境业务数据)
- SQLite (测试环境业务数据)

## 组件和接口

### 前端组件

#### 1. PaymentList 组件
**职责**：显示支付记录列表，支持搜索、筛选和分页

**Props**：
```typescript
interface PaymentListProps {
  onViewDetail: (paymentId: number) => void;
  onRefund: (paymentId: number) => void;
  onExport: () => void;
}
```

#### 2. RefundModal 组件
**职责**：处理退款操作的模态框

**Props**：
```typescript
interface RefundModalProps {
  visible: boolean;
  payment: Payment;
  onConfirm: (amount: number, reason: string) => Promise<void>;
  onCancel: () => void;
}
```

#### 3. ReconciliationPanel 组件
**职责**：财务对账面板，展示对账结果

**Props**：
```typescript
interface ReconciliationPanelProps {
  dateRange: [Date, Date];
  onReconcile: (dateRange: [Date, Date]) => Promise<ReconciliationResult>;
}
```

### 后端接口

#### 1. PaymentService 接口
```go
type PaymentService interface {
    ListPayments(ctx context.Context, req *ListPaymentsRequest) (*ListPaymentsResponse, error)
    GetPayment(ctx context.Context, paymentID int64) (*Payment, error)
    CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*Payment, error)
    UpdatePaymentStatus(ctx context.Context, paymentID int64, status PaymentStatus) error
    ProcessRefund(ctx context.Context, req *RefundRequest) (*RefundResult, error)
    CapturePayment(ctx context.Context, paymentID int64) error
}
```

## 数据模型

### Payment (支付记录)
```typescript
interface Payment {
  id: number;
  orderID: number;
  userID: number;
  method: PaymentMethod;
  amount: number;
  currency: string;
  status: PaymentStatus;
  providerTradeNo: string;
  refundedAmount: number;
  createdAt: Date;
  updatedAt: Date;
  paidAt?: Date;
}

enum PaymentMethod {
  ALIPAY = 'alipay',
  WECHAT = 'wechat',
  BANK_CARD = 'bank_card',
  BALANCE = 'balance'
}

enum PaymentStatus {
  PENDING = 'pending',
  PAID = 'paid',
  REFUNDED = 'refunded',
  FAILED = 'failed'
}
```

## 正确性属性

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的正式陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性 1：支付记录完整性
*对于任何*支付记录列表请求，返回的所有支付记录必须包含完整的必填字段（订单ID、用户ID、支付方式、金额、状态、创建时间），且字段值不为空或null
**验证：需求 1.1**

### 属性 2：退款金额约束
*对于任何*退款操作，退款金额必须大于0且不超过原支付金额减去已退款金额
**验证：需求 2.1, 9.1, 9.2**

### 属性 3：支付状态转换合法性
*对于任何*支付状态更新操作，状态转换必须遵循合法路径：pending → paid → refunded 或 pending → failed，不允许其他转换
**验证：需求 2.3, 2.4, 8.2**

### 属性 4：对账数据一致性
*对于任何*对账操作，计算的总收入必须等于所有已支付记录的金额之和减去所有退款记录的金额之和
**验证：需求 3.2**

### 属性 5：分页数据连续性
*对于任何*分页查询，当前页的最后一条记录的创建时间必须大于或等于下一页第一条记录的创建时间（降序排列时）
**验证：需求 1.4, 4.3**

### 属性 6：权限验证一致性
*对于任何*需要权限的操作，系统必须先验证用户权限，验证失败时必须返回403错误且不执行任何业务逻辑
**验证：需求 10.1, 10.2, 10.3, 10.4**

### 属性 7：回调签名验证
*对于任何*第三方支付回调，系统必须验证签名有效性，签名无效时必须拒绝请求且不更新任何数据
**验证：需求 8.1, 8.3**

### 属性 8：操作日志完整性
*对于任何*支付相关操作（创建、更新、退款、确认），系统必须记录操作日志，包含操作类型、操作人、操作时间、操作前状态和操作后状态
**验证：需求 7.1**

### 属性 9：导出数据一致性
*对于任何*数据导出操作，导出文件中的数据必须与当前筛选条件下的查询结果完全一致
**验证：需求 1.5, 4.5**

### 属性 10：退款幂等性
*对于任何*相同的退款请求（相同支付ID和退款金额），多次提交只应执行一次退款操作
**验证：需求 2.2**

## 错误处理

### 错误分类

#### 1. 客户端错误 (4xx)
- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 资源不存在

#### 2. 服务端错误 (5xx)
- **500 Internal Server Error**: 服务器内部错误
- **502 Bad Gateway**: 第三方服务错误
- **503 Service Unavailable**: 服务不可用

## 测试策略

### 单元测试
使用 Vitest + Testing Library 进行组件和工具函数测试

### 属性测试
使用 fast-check 验证系统的通用正确性属性，每个属性测试运行至少100次迭代

### 集成测试
测试 API 端点完整流程、数据库事务、第三方服务集成、权限控制

### E2E测试
使用 Playwright 进行端到端测试
