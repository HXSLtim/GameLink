# 设计文档 - 支付与财务模块

## 概述

支付与财务模块是 GameLink 平台的核心财务管理系统，负责处理所有与资金相关的业务流程，包括支付管理、退款处理、财务对账、收款分流和提现分流等功能。该模块支持多公司主体运营，实现收款和工资发放的自动分流，满足税务合规要求。模块采用前后端分离架构，前端基于 React + TypeScript + Ant Design 构建，后端通过 RESTful API 提供服务。模块设计遵循高内聚低耦合原则，确保系统的可维护性和可扩展性。

## 架构

### 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              前端层 (React)                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  支付记录  │  退款管理  │  财务对账  │  报表中心  │  分流管理  │  税务报表  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────────────────────┐
│                            API 网关层 (Go)                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│  路由管理  │  权限验证  │  请求限流  │  日志记录  │  错误处理  │  多租户隔离 │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ↓
┌─────────────────────────────────────────────────────────────────────────────┐
│                            业务逻辑层 (Go)                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│  支付服务  │  退款服务  │  对账服务  │  报表服务  │  分流服务  │  税务服务  │
│            │            │            │            │            │            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        分流引擎 (Routing Engine)                     │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │  收款分流器  │  提现分流器  │  规则匹配器  │  主体选择器  │  容错处理  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ↓
┌─────────────────────────────────────────────────────────────────────────────┐
│                            数据访问层 (Go)                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│  支付仓储  │  流水仓储  │  对账仓储  │  主体仓储  │  规则仓储  │  税务仓储  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ↓
┌─────────────────────────────────────────────────────────────────────────────┐
│                            数据存储层                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  PostgreSQL/SQLite (主数据)  │  Redis (缓存/规则)  │  MongoDB (日志/审计)    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ↓
┌─────────────────────────────────────────────────────────────────────────────┐
│                            外部服务层                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  支付宝 API  │  微信支付 API  │  银行转账接口  │  税务系统接口  │  消息队列  │
│  (多商户)    │  (多商户)      │  (多账户)      │               │            │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 分流架构详图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              收款分流流程                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   用户支付请求                                                               │
│        │                                                                    │
│        ↓                                                                    │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────────┐    │
│   │ 订单信息    │───→│ 规则匹配器  │───→│ 收款主体选择                │    │
│   │ (游戏/服务) │    │             │    │ ┌─────────┐ ┌─────────┐    │    │
│   └─────────────┘    └─────────────┘    │ │ 公司A   │ │ 公司B   │    │    │
│                                         │ │ 商户号A │ │ 商户号B │    │    │
│                                         │ └────┬────┘ └────┬────┘    │    │
│                                         └──────┼───────────┼─────────┘    │
│                                                ↓           ↓              │
│                                         ┌─────────────────────────┐       │
│                                         │ 第三方支付平台          │       │
│                                         │ (支付宝/微信/银行)      │       │
│                                         └─────────────────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              提现分流流程                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   陪玩师提现申请                                                             │
│        │                                                                    │
│        ↓                                                                    │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────────┐    │
│   │ 陪玩师信息  │───→│ 结算公司    │───→│ 工资发放主体选择            │    │
│   │ (所属公司)  │    │ 分配查询    │    │ ┌─────────┐ ┌─────────┐    │    │
│   └─────────────┘    └─────────────┘    │ │ 结算公司A│ │ 结算公司B│    │    │
│                                         │ │ 银行账户A│ │ 银行账户B│    │    │
│                                         │ └────┬────┘ └────┬────┘    │    │
│                                         └──────┼───────────┼─────────┘    │
│                                                ↓           ↓              │
│                                         ┌─────────────────────────┐       │
│                                         │ 银行转账系统            │       │
│                                         │ (工资发放/个税代扣)     │       │
│                                         └─────────────────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘
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

#### 4. SettlementCompanyList 组件
**职责**：结算公司列表管理

**Props**：
```typescript
interface SettlementCompanyListProps {
  onAdd: () => void;
  onEdit: (companyId: number) => void;
  onToggleStatus: (companyId: number, enabled: boolean) => Promise<void>;
}
```

#### 5. CollectionEntityList 组件
**职责**：收款主体列表管理

**Props**：
```typescript
interface CollectionEntityListProps {
  onAdd: () => void;
  onEdit: (entityId: number) => void;
  onConfigureChannel: (entityId: number) => void;
}
```

#### 6. RoutingRuleConfig 组件
**职责**：收款分流规则配置

**Props**：
```typescript
interface RoutingRuleConfigProps {
  onSave: (rule: RoutingRule) => Promise<void>;
  onTest: (rule: RoutingRule, testOrder: TestOrder) => Promise<RoutingResult>;
}
```

#### 7. WithdrawalRoutingStats 组件
**职责**：提现分流统计报表

**Props**：
```typescript
interface WithdrawalRoutingStatsProps {
  dateRange: [Date, Date];
  settlementCompanyId?: number;
  onExport: (format: 'excel' | 'pdf') => void;
}
```

#### 8. TaxReportGenerator 组件
**职责**：税务报表生成

**Props**：
```typescript
interface TaxReportGeneratorProps {
  reportType: 'income_tax' | 'vat';
  entityId: number;
  period: { year: number; month: number };
  onGenerate: () => Promise<TaxReport>;
  onExport: (format: string) => void;
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

#### 2. SettlementCompanyService 接口
```go
type SettlementCompanyService interface {
    CreateCompany(ctx context.Context, req *CreateSettlementCompanyRequest) (*SettlementCompany, error)
    UpdateCompany(ctx context.Context, companyID int64, req *UpdateSettlementCompanyRequest) error
    GetCompany(ctx context.Context, companyID int64) (*SettlementCompany, error)
    ListCompanies(ctx context.Context, req *ListCompaniesRequest) (*ListCompaniesResponse, error)
    ToggleCompanyStatus(ctx context.Context, companyID int64, enabled bool) error
    AssignPlayerToCompany(ctx context.Context, playerID int64, companyID int64, reason string) error
    BatchAssignPlayers(ctx context.Context, playerIDs []int64, companyID int64, reason string) error
    GetPlayerAssignmentHistory(ctx context.Context, playerID int64) ([]PlayerCompanyAssignment, error)
}
```

#### 3. CollectionEntityService 接口
```go
type CollectionEntityService interface {
    CreateEntity(ctx context.Context, req *CreateCollectionEntityRequest) (*CollectionEntity, error)
    UpdateEntity(ctx context.Context, entityID int64, req *UpdateCollectionEntityRequest) error
    GetEntity(ctx context.Context, entityID int64) (*CollectionEntity, error)
    ListEntities(ctx context.Context, req *ListEntitiesRequest) (*ListEntitiesResponse, error)
    ToggleEntityStatus(ctx context.Context, entityID int64, enabled bool) error
    ConfigurePaymentChannel(ctx context.Context, entityID int64, channel *PaymentChannelConfig) error
}
```

#### 4. RoutingService 接口
```go
type RoutingService interface {
    // 收款分流
    CreateRoutingRule(ctx context.Context, req *CreateRoutingRuleRequest) (*RoutingRule, error)
    UpdateRoutingRule(ctx context.Context, ruleID int64, req *UpdateRoutingRuleRequest) error
    ListRoutingRules(ctx context.Context, req *ListRoutingRulesRequest) (*ListRoutingRulesResponse, error)
    ToggleRuleStatus(ctx context.Context, ruleID int64, enabled bool) error
    SetDefaultEntity(ctx context.Context, entityID int64) error
    MatchCollectionEntity(ctx context.Context, order *Order) (*CollectionEntity, error)
    
    // 提现分流
    RouteWithdrawal(ctx context.Context, withdrawal *Withdrawal) (*SettlementCompany, error)
    GetWithdrawalsByCompany(ctx context.Context, companyID int64, dateRange DateRange) ([]Withdrawal, error)
}
```

#### 5. TaxReportService 接口
```go
type TaxReportService interface {
    GenerateIncomeTaxReport(ctx context.Context, companyID int64, period Period) (*IncomeTaxReport, error)
    GenerateVATReport(ctx context.Context, entityID int64, period Period) (*VATReport, error)
    ExportTaxReport(ctx context.Context, reportID int64, format string) ([]byte, error)
    ListTaxReports(ctx context.Context, req *ListTaxReportsRequest) (*ListTaxReportsResponse, error)
}
```

#### 6. MultiEntityReconciliationService 接口
```go
type MultiEntityReconciliationService interface {
    ReconcileByCollectionEntity(ctx context.Context, entityID int64, dateRange DateRange) (*ReconciliationResult, error)
    ReconcileBySettlementCompany(ctx context.Context, companyID int64, dateRange DateRange) (*ReconciliationResult, error)
    GetReconciliationDifferences(ctx context.Context, reconciliationID int64) ([]ReconciliationDifference, error)
    MarkDifferenceResolved(ctx context.Context, differenceID int64, resolution string) error
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
  collectionEntityID: number;      // 收款主体ID
  merchantNo: string;              // 实际使用的商户号
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

### SettlementCompany (结算公司)
```typescript
interface SettlementCompany {
  id: number;
  name: string;                    // 公司名称
  creditCode: string;              // 统一社会信用代码
  taxRegistrationNo: string;       // 税务登记号
  bankName: string;                // 开户银行
  bankAccount: string;             // 银行账号
  bankBranch: string;              // 开户支行
  contactName: string;             // 联系人
  contactPhone: string;            // 联系电话
  address: string;                 // 公司地址
  status: CompanyStatus;           // 状态
  playerCount: number;             // 关联陪玩师数量
  createdAt: Date;
  updatedAt: Date;
}

enum CompanyStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive'
}
```

### PlayerCompanyAssignment (陪玩师公司分配)
```typescript
interface PlayerCompanyAssignment {
  id: number;
  playerID: number;
  settlementCompanyID: number;
  effectiveDate: Date;             // 生效日期
  endDate?: Date;                  // 结束日期
  reason: string;                  // 分配原因
  assignedBy: number;              // 分配操作人
  createdAt: Date;
}
```

### CollectionEntity (收款主体)
```typescript
interface CollectionEntity {
  id: number;
  name: string;                    // 公司名称
  creditCode: string;              // 统一社会信用代码
  taxRegistrationNo: string;       // 税务登记号
  status: EntityStatus;            // 状态
  isDefault: boolean;              // 是否默认收款主体
  paymentChannels: PaymentChannelConfig[];  // 支付渠道配置
  totalCollection: number;         // 累计收款金额
  createdAt: Date;
  updatedAt: Date;
}

interface PaymentChannelConfig {
  channel: PaymentMethod;          // 支付渠道
  merchantNo: string;              // 商户号
  merchantKey: string;             // 商户密钥 (加密存储)
  callbackUrl: string;             // 回调地址
  enabled: boolean;                // 是否启用
}

enum EntityStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive'
}
```

### RoutingRule (分流规则)
```typescript
interface RoutingRule {
  id: number;
  name: string;                    // 规则名称
  priority: number;                // 优先级 (数字越小优先级越高)
  conditions: RoutingCondition[];  // 匹配条件
  targetEntityID: number;          // 目标收款主体ID
  status: RuleStatus;              // 状态
  createdAt: Date;
  updatedAt: Date;
}

interface RoutingCondition {
  field: ConditionField;           // 条件字段
  operator: ConditionOperator;     // 操作符
  value: string | number | string[];  // 条件值
}

enum ConditionField {
  GAME_TYPE = 'game_type',
  SERVICE_TYPE = 'service_type',
  ORDER_AMOUNT = 'order_amount',
  REGION = 'region'
}

enum ConditionOperator {
  EQUALS = 'eq',
  NOT_EQUALS = 'neq',
  IN = 'in',
  NOT_IN = 'not_in',
  GREATER_THAN = 'gt',
  LESS_THAN = 'lt',
  BETWEEN = 'between'
}

enum RuleStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive'
}
```

### Withdrawal (提现记录) - 扩展
```typescript
interface Withdrawal {
  id: number;
  playerID: number;
  amount: number;
  status: WithdrawalStatus;
  settlementCompanyID: number;     // 结算公司ID
  settlementCompanyName: string;   // 结算公司名称
  paymentBankAccount: string;      // 付款银行账户
  bankTransactionNo?: string;      // 银行流水号
  taxDeducted: number;             // 代扣个税金额
  actualAmount: number;            // 实际到账金额
  approvedBy?: number;             // 审批人
  approvedAt?: Date;               // 审批时间
  paidAt?: Date;                   // 付款时间
  createdAt: Date;
  updatedAt: Date;
}

enum WithdrawalStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  PAID = 'paid',
  FAILED = 'failed'
}
```

### TaxReport (税务报表)
```typescript
interface TaxReport {
  id: number;
  reportType: TaxReportType;
  entityID: number;                // 主体ID (结算公司或收款主体)
  entityType: EntityType;          // 主体类型
  period: Period;                  // 报表周期
  totalAmount: number;             // 总金额
  taxableAmount: number;           // 应税金额
  taxAmount: number;               // 税额
  details: TaxReportDetail[];      // 明细
  status: ReportStatus;
  generatedAt: Date;
  exportedAt?: Date;
}

interface Period {
  year: number;
  month: number;
}

enum TaxReportType {
  INCOME_TAX = 'income_tax',       // 个人所得税
  VAT = 'vat'                      // 增值税
}

enum EntityType {
  SETTLEMENT_COMPANY = 'settlement_company',
  COLLECTION_ENTITY = 'collection_entity'
}

enum ReportStatus {
  DRAFT = 'draft',
  FINALIZED = 'finalized',
  EXPORTED = 'exported'
}
```

## 正确性属性

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的正式陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性 1：支付记录完整性
*对于任何*支付记录列表请求，返回的所有支付记录必须包含完整的必填字段（订单ID、用户ID、支付方式、金额、状态、收款主体ID、创建时间），且字段值不为空或null
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

### 属性 11：统一社会信用代码格式验证
*对于任何*结算公司或收款主体的创建/更新操作，统一社会信用代码必须符合18位标准格式，且在系统内唯一
**验证：需求 11.2, 15.2**

### 属性 12：陪玩师结算公司分配唯一性
*对于任何*时间点，每个陪玩师只能有一个生效的结算公司分配，新分配生效时旧分配自动结束
**验证：需求 12.1, 12.4**

### 属性 13：提现分流一致性
*对于任何*陪玩师提现申请，系统分配的结算公司必须与该陪玩师当前生效的结算公司分配一致
**验证：需求 13.1**

### 属性 14：提现分流统计准确性
*对于任何*结算公司的提现统计，统计的提现总额必须等于该公司所有已完成提现记录金额之和
**验证：需求 14.1**

### 属性 15：收款分流规则优先级
*对于任何*支付请求，当存在多条匹配的分流规则时，系统必须选择优先级最高（数值最小）的规则确定收款主体
**验证：需求 16.2**

### 属性 16：收款分流默认主体回退
*对于任何*支付请求，当没有规则匹配时，系统必须使用标记为默认的收款主体处理支付
**验证：需求 16.3, 17.4**

### 属性 17：收款分流记录完整性
*对于任何*成功的支付记录，必须包含实际收款主体ID和对应的商户号信息
**验证：需求 17.3**

### 属性 18：收款分流统计准确性
*对于任何*收款主体的收款统计，统计的收款总额必须等于该主体所有已支付记录金额之和
**验证：需求 18.1**

### 属性 19：多主体对账隔离性
*对于任何*按主体进行的对账操作，对账结果只包含该主体的交易记录，不包含其他主体的数据
**验证：需求 19.2, 19.3**

### 属性 20：个税计算一致性
*对于任何*个税申报报表，每个陪玩师的应扣税额必须根据其收入金额按照个税计算公式正确计算
**验证：需求 20.2**

## 营收报表设计

### 报表类型

#### 1. 平台总营收报表
展示平台整体营收情况，包括：
- 总收入、总支出、净利润
- 按时间维度（日/周/月/年）的收入趋势
- 按支付方式的收入分布
- 按游戏类型的收入分布
- 按服务类型的收入分布

#### 2. 收款主体营收报表
按收款主体分类展示营收：
- 各收款主体的收款总额、退款总额、净收入
- 各主体的收款笔数和平均订单金额
- 各主体的收款趋势对比图
- 各主体的收款占比饼图

#### 3. 结算公司支出报表
按结算公司分类展示支出：
- 各结算公司的工资发放总额
- 各公司的发放笔数和平均发放金额
- 各公司的发放趋势对比图
- 各公司的支出占比饼图

#### 4. 利润分析报表
- 毛利润 = 总收入 - 陪玩师分成
- 净利润 = 毛利润 - 运营成本 - 税费
- 各收款主体的利润贡献
- 各游戏/服务类型的利润率分析

### 报表数据模型

```typescript
interface RevenueReport {
  id: number;
  reportType: RevenueReportType;
  period: ReportPeriod;
  generatedAt: Date;
  
  // 汇总数据
  summary: RevenueSummary;
  
  // 按主体分类数据
  byCollectionEntity: EntityRevenue[];
  bySettlementCompany: CompanyExpense[];
  
  // 趋势数据
  trendData: TrendDataPoint[];
  
  // 分布数据
  byPaymentMethod: DistributionItem[];
  byGameType: DistributionItem[];
  byServiceType: DistributionItem[];
}

interface RevenueSummary {
  totalRevenue: number;           // 总收入
  totalRefund: number;            // 总退款
  netRevenue: number;             // 净收入
  totalPayout: number;            // 总支出（工资发放）
  grossProfit: number;            // 毛利润
  platformFee: number;            // 平台手续费收入
  transactionCount: number;       // 交易笔数
  averageOrderValue: number;      // 平均客单价
}

interface EntityRevenue {
  entityId: number;
  entityName: string;
  revenue: number;
  refund: number;
  netRevenue: number;
  transactionCount: number;
  percentage: number;             // 占比
}

interface CompanyExpense {
  companyId: number;
  companyName: string;
  totalPayout: number;
  payoutCount: number;
  taxDeducted: number;            // 代扣税额
  percentage: number;             // 占比
}

interface TrendDataPoint {
  date: string;
  revenue: number;
  refund: number;
  payout: number;
  profit: number;
}

interface DistributionItem {
  name: string;
  value: number;
  percentage: number;
}

enum RevenueReportType {
  DAILY = 'daily',
  WEEKLY = 'weekly',
  MONTHLY = 'monthly',
  QUARTERLY = 'quarterly',
  YEARLY = 'yearly',
  CUSTOM = 'custom'
}

interface ReportPeriod {
  startDate: Date;
  endDate: Date;
  type: RevenueReportType;
}
```

### 报表服务接口

```go
type RevenueReportService interface {
    // 生成营收报表
    GenerateRevenueReport(ctx context.Context, req *GenerateReportRequest) (*RevenueReport, error)
    
    // 获取平台营收汇总
    GetRevenueSummary(ctx context.Context, period *ReportPeriod) (*RevenueSummary, error)
    
    // 按收款主体获取营收
    GetRevenueByEntity(ctx context.Context, period *ReportPeriod) ([]EntityRevenue, error)
    
    // 按结算公司获取支出
    GetExpenseByCompany(ctx context.Context, period *ReportPeriod) ([]CompanyExpense, error)
    
    // 获取营收趋势
    GetRevenueTrend(ctx context.Context, period *ReportPeriod, granularity string) ([]TrendDataPoint, error)
    
    // 获取收入分布
    GetRevenueDistribution(ctx context.Context, period *ReportPeriod, dimension string) ([]DistributionItem, error)
    
    // 导出报表
    ExportRevenueReport(ctx context.Context, reportID int64, format string) ([]byte, error)
    
    // 获取历史报表列表
    ListRevenueReports(ctx context.Context, req *ListReportsRequest) (*ListReportsResponse, error)
}
```

### 前端报表组件

```typescript
// 营收报表页面组件
interface RevenueReportPageProps {
  defaultPeriod?: ReportPeriod;
}

// 营收汇总卡片
interface RevenueSummaryCardsProps {
  summary: RevenueSummary;
  loading?: boolean;
  comparePeriod?: RevenueSummary;  // 对比周期数据
}

// 主体营收对比图
interface EntityRevenueChartProps {
  data: EntityRevenue[];
  chartType: 'bar' | 'pie';
  loading?: boolean;
}

// 营收趋势图
interface RevenueTrendChartProps {
  data: TrendDataPoint[];
  metrics: ('revenue' | 'refund' | 'payout' | 'profit')[];
  loading?: boolean;
}

// 收入分布图
interface RevenueDistributionChartProps {
  data: DistributionItem[];
  dimension: 'paymentMethod' | 'gameType' | 'serviceType';
  loading?: boolean;
}
```

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
- 前端：使用 Vitest + Testing Library 进行组件和工具函数测试
- 后端：使用 Go testing + testify 进行服务层和仓储层测试

### 属性测试
- 前端：使用 fast-check 验证系统的通用正确性属性，每个属性测试运行至少100次迭代
- 后端：使用 gopter 或 rapid 进行属性测试
- 每个属性测试必须标注对应的正确性属性编号，格式：`**Feature: payment-finance-module, Property {number}: {property_text}**`

### 集成测试
- 测试 API 端点完整流程
- 数据库事务测试
- 第三方服务集成（使用 mock）
- 权限控制测试
- 分流规则匹配测试
- 多主体对账测试

### E2E测试
使用 Playwright 进行端到端测试，覆盖：
- 结算公司管理流程
- 收款主体管理流程
- 分流规则配置流程
- 提现分流审批流程
- 税务报表生成流程
