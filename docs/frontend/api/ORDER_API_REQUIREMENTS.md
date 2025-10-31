# 订单管理接口需求文档

## 📊 接口现状

### ✅ 已实现接口（10个）

| 方法               | 端点                                   | 说明         | 状态 |
| ------------------ | -------------------------------------- | ------------ | ---- |
| `getList()`        | `GET /api/v1/admin/orders`             | 获取订单列表 | ✅   |
| `getDetail(id)`    | `GET /api/v1/admin/orders/:id`         | 获取订单详情 | ✅   |
| `create(data)`     | `POST /api/v1/admin/orders`            | 创建订单     | ✅   |
| `update(id, data)` | `PUT /api/v1/admin/orders/:id`         | 更新订单信息 | ✅   |
| `assign(id, data)` | `POST /api/v1/admin/orders/:id/assign` | 分配陪玩师   | ✅   |
| `review(id, data)` | `POST /api/v1/admin/orders/:id/review` | 审核订单     | ✅   |
| `cancel(id, data)` | `POST /api/v1/admin/orders/:id/cancel` | 取消订单     | ✅   |
| `delete(id)`       | `DELETE /api/v1/admin/orders/:id`      | 删除订单     | ✅   |
| `getLogs(id)`      | `GET /api/v1/admin/orders/:id/logs`    | 获取操作日志 | ✅   |
| `getStatistics()`  | `GET /api/v1/admin/stats/orders`       | 获取订单统计 | ✅   |

---

## 🔧 建议补充的接口

### 1. 订单状态流转接口

#### 1.1 确认订单

```typescript
confirm: (id: number, data?: { note?: string }): Promise<Order>
```

**端点**: `POST /api/v1/admin/orders/:id/confirm`

**说明**: 管理员确认订单，状态从 `pending` → `confirmed`

**请求体**:

```typescript
{
  note?: string;  // 备注信息
}
```

**使用场景**:

- 管理员审核用户提交的订单
- 确认订单信息无误后进入下一流程

---

#### 1.2 开始服务

```typescript
start: (id: number, data?: { note?: string }): Promise<Order>
```

**端点**: `POST /api/v1/admin/orders/:id/start`

**说明**: 陪玩师开始服务，状态从 `confirmed` → `in_progress`

**请求体**:

```typescript
{
  note?: string;  // 开始服务备注
}
```

**使用场景**:

- 陪玩师接受订单并开始服务
- 记录实际开始时间

---

#### 1.3 完成订单

```typescript
complete: (id: number, data?: { note?: string }): Promise<Order>
```

**端点**: `POST /api/v1/admin/orders/:id/complete`

**说明**: 完成订单服务，状态从 `in_progress` → `completed`

**请求体**:

```typescript
{
  note?: string;  // 完成备注
}
```

**使用场景**:

- 服务完成后标记订单完成
- 触发评价、结算等后续流程

---

### 2. 财务相关接口

#### 2.1 退款处理

```typescript
refund: (id: number, data: RefundRequest): Promise<Order>
```

**端点**: `POST /api/v1/admin/orders/:id/refund`

**说明**: 处理订单退款

**请求体**:

```typescript
interface RefundRequest {
  reason: string; // 退款原因（必填）
  amount_cents?: number; // 退款金额（分），不填则全额退款
  note?: string; // 备注信息
}
```

**响应**:

```typescript
{
  id: number;
  status: 'refunded';
  refund_amount_cents: number;
  refund_reason: string;
  refunded_at: string;
}
```

**使用场景**:

- 用户申请退款
- 订单异常需要退款
- 部分退款或全额退款

---

### 3. 批量操作接口

#### 3.1 批量分配订单

```typescript
batchAssign: (data: BatchAssignRequest): Promise<{ success: number; failed: number }>
```

**端点**: `POST /api/v1/admin/orders/batch/assign`

**请求体**:

```typescript
interface BatchAssignRequest {
  order_ids: number[]; // 订单ID数组
  player_id: number; // 陪玩师ID
  note?: string; // 备注
}
```

**响应**:

```typescript
{
  success: number; // 成功数量
  failed: number; // 失败数量
  results: Array<{
    order_id: number;
    success: boolean;
    error?: string;
  }>;
}
```

---

#### 3.2 批量审核订单

```typescript
batchReview: (data: BatchReviewRequest): Promise<{ success: number; failed: number }>
```

**端点**: `POST /api/v1/admin/orders/batch/review`

**请求体**:

```typescript
interface BatchReviewRequest {
  order_ids: number[];
  result: 'approved' | 'rejected';
  reason?: string;
}
```

---

#### 3.3 批量取消订单

```typescript
batchCancel: (data: BatchCancelRequest): Promise<{ success: number; failed: number }>
```

**端点**: `POST /api/v1/admin/orders/batch/cancel`

**请求体**:

```typescript
interface BatchCancelRequest {
  order_ids: number[];
  reason: string;
}
```

---

### 4. 数据导出接口

#### 4.1 导出订单列表

```typescript
exportList: (params: OrderListQuery & { format?: 'csv' | 'excel' }): Promise<Blob>
```

**端点**: `GET /api/v1/admin/orders/export`

**Query 参数**: 与 `getList` 相同，额外增加 `format` 参数

**响应**: 文件流（CSV 或 Excel）

**使用场景**:

- 导出财务报表
- 导出运营数据
- 数据备份

---

### 5. 高级查询接口

#### 5.1 获取订单时间线

```typescript
getTimeline: (id: number): Promise<OrderTimeline[]>
```

**端点**: `GET /api/v1/admin/orders/:id/timeline`

**响应**:

```typescript
interface OrderTimeline {
  id: number;
  order_id: number;
  event_type: 'status_change' | 'action' | 'system' | 'note';
  title: string;
  description?: string;
  operator?: string;
  operator_role?: string;
  status_before?: string;
  status_after?: string;
  metadata?: Record<string, any>;
  created_at: string;
}
```

**使用场景**:

- 查看订单完整历史
- 问题追溯和调查
- 数据审计

---

#### 5.2 获取审核记录

```typescript
getReviews: (id: number): Promise<OrderReview[]>
```

**端点**: `GET /api/v1/admin/orders/:id/reviews`

**响应**:

```typescript
interface OrderReview {
  id: number;
  order_id: number;
  reviewer_id: number;
  reviewer_name: string;
  result: 'approved' | 'rejected';
  reason?: string;
  comment?: string;
  created_at: string;
}
```

---

### 6. 关联数据接口

#### 6.1 获取支付记录

```typescript
getPayments: (id: number): Promise<Payment[]>
```

**端点**: `GET /api/v1/admin/orders/:id/payments`

**响应**:

```typescript
interface Payment {
  id: number;
  order_id: number;
  amount_cents: number;
  currency: string;
  payment_method: string;
  payment_status: 'pending' | 'success' | 'failed';
  transaction_id?: string;
  paid_at?: string;
  created_at: string;
}
```

---

#### 6.2 获取退款记录

```typescript
getRefunds: (id: number): Promise<Refund[]>
```

**端点**: `GET /api/v1/admin/orders/:id/refunds`

**响应**:

```typescript
interface Refund {
  id: number;
  order_id: number;
  amount_cents: number;
  reason: string;
  status: 'pending' | 'processing' | 'success' | 'failed';
  refund_method: string;
  refunded_at?: string;
  created_at: string;
}
```

---

## 📝 实现建议

### 优先级分类

#### 🔴 高优先级（核心业务流程）

1. `confirm()` - 确认订单
2. `start()` - 开始服务
3. `complete()` - 完成订单
4. `refund()` - 退款处理
5. `getTimeline()` - 订单时间线

#### 🟡 中优先级（提升效率）

1. `batchAssign()` - 批量分配
2. `batchReview()` - 批量审核
3. `batchCancel()` - 批量取消
4. `exportList()` - 数据导出

#### 🟢 低优先级（优化体验）

1. `getReviews()` - 审核记录（如果 `getDetail` 已包含）
2. `getPayments()` - 支付记录（如果有独立的支付管理模块）
3. `getRefunds()` - 退款记录（如果有独立的退款管理模块）

---

## 🎯 接口实现示例

### 示例 1: 添加状态流转接口

```typescript
// src/services/api/order.ts

export const orderApi = {
  // ... 现有接口 ...

  /**
   * 确认订单
   */
  confirm: (id: number, data?: { note?: string }): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/confirm`, data);
  },

  /**
   * 开始服务
   */
  start: (id: number, data?: { note?: string }): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/start`, data);
  },

  /**
   * 完成订单
   */
  complete: (id: number, data?: { note?: string }): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/complete`, data);
  },

  /**
   * 退款处理
   */
  refund: (id: number, data: RefundRequest): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/refund`, data);
  },
};
```

### 示例 2: 添加类型定义

```typescript
// src/types/order.ts

/**
 * 退款请求
 */
export interface RefundRequest {
  reason: string;
  amount_cents?: number;
  note?: string;
}

/**
 * 批量分配请求
 */
export interface BatchAssignRequest {
  order_ids: number[];
  player_id: number;
  note?: string;
}

/**
 * 批量审核请求
 */
export interface BatchReviewRequest {
  order_ids: number[];
  result: 'approved' | 'rejected';
  reason?: string;
}

/**
 * 批量取消请求
 */
export interface BatchCancelRequest {
  order_ids: number[];
  reason: string;
}

/**
 * 批量操作响应
 */
export interface BatchOperationResponse {
  success: number;
  failed: number;
  results: Array<{
    order_id: number;
    success: boolean;
    error?: string;
  }>;
}

/**
 * 订单时间线事件
 */
export interface OrderTimeline {
  id: number;
  order_id: number;
  event_type: 'status_change' | 'action' | 'system' | 'note';
  title: string;
  description?: string;
  operator?: string;
  operator_role?: string;
  status_before?: string;
  status_after?: string;
  metadata?: Record<string, any>;
  created_at: string;
}
```

---

## 🔍 与现有接口的关系

### 状态流转接口 vs 更新接口

- `update()`: 修改订单基本信息（标题、描述、价格等）
- `confirm()`, `start()`, `complete()`: 修改订单状态，附带业务逻辑（如时间记录、通知等）

### 单个操作 vs 批量操作

- 单个操作: `assign()`, `review()`, `cancel()`
- 批量操作: `batchAssign()`, `batchReview()`, `batchCancel()`
- 批量操作适用于需要同时处理多个订单的场景

### 日志 vs 时间线

- `getLogs()`: 获取操作日志（管理员操作记录）
- `getTimeline()`: 获取完整时间线（包括状态变更、系统事件、用户操作等）

---

## ✅ 总结

### 当前状态

- ✅ 基础 CRUD 操作完整
- ✅ 订单分配、审核、取消功能完整
- ✅ 统计和日志查询完整

### 建议补充

- ❌ 状态流转接口（confirm, start, complete）
- ❌ 退款处理接口
- ❌ 批量操作接口
- ❌ 数据导出接口
- ❌ 订单时间线接口

### 实施建议

1. 优先实现**高优先级**接口（核心业务流程）
2. 根据实际使用频率逐步补充**中优先级**接口
3. 在功能稳定后考虑**低优先级**接口

---

**文档版本**: v1.0  
**创建时间**: 2025-10-28  
**维护者**: GameLink 开发团队
