# 订单拆分与转单功能设计

## 业务背景

用户下单 N 小时陪玩服务时，可能出现以下场景：
- 陪玩 A 打了 1 小时后无法继续
- 需要将剩余时间转给陪玩 B
- 用户端仍然看到一个完整订单

## 数据模型

### OrderGroup（主订单 - 用户视角）

```
order_groups
├── id                  主键
├── group_no            主订单号（用户看到的）
├── user_id             下单用户
├── game_id             游戏ID
├── item_id             服务项目ID
├── original_player_id  原始陪玩师
├── total_price_cents   总价
├── total_hours         总时长
├── completed_hours     已完成时长
├── status              主订单状态
└── ...
```

### Order（子订单 - 每小时一个）

新增字段：
```
orders
├── group_id            关联主订单ID
├── hour_index          第几小时 (1, 2, 3...)
├── is_sub_order        是否为子订单
├── can_transfer        是否可转单
├── transfer_from       转单来源订单ID
├── transfer_to         转单目标订单ID
└── transfer_note       转单备注
```

## 订单创建流程

```
用户下单 3 小时
    │
    ▼
┌─────────────────────────────────────┐
│  CreateOrder(duration=3)            │
│  检测到 duration > 1 小时            │
│  调用 createOrderWithSplit()        │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  1. 创建 OrderGroup (主订单)         │
│     - group_no: G202601060001       │
│     - total_hours: 3                │
│     - status: pending               │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  2. 创建 3 个 Order (子订单)         │
│     - order_1: hour_index=1         │
│     - order_2: hour_index=2         │
│     - order_3: hour_index=3         │
│     - 每个关联 group_id             │
└─────────────────────────────────────┘
    │
    ▼
返回 { orderId: group.ID, isSplit: true }
```

## 转单流程

```
陪玩 A 完成第 1 小时，无法继续
    │
    ▼
┌─────────────────────────────────────┐
│  TransferSubOrder(subOrderId=2,     │
│                   newPlayerId=B)    │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  1. 验证原订单可转单                 │
│  2. 验证新陪玩师有效                 │
│  3. 创建新子订单（复制信息，换陪玩）  │
│  4. 原订单标记为已取消               │
│  5. 更新主订单状态                   │
└─────────────────────────────────────┘
```

## API 接口

### 用户端

#### 创建订单
```
POST /api/v1/user/orders
{
  "playerId": 1,
  "gameId": 1,
  "durationHours": 3,
  "scheduledStart": "2026-01-06T10:00:00Z",
  "title": "王者荣耀陪玩"
}

Response:
{
  "orderId": 100,        // 主订单ID
  "groupNo": "G202601060001",
  "priceCents": 15000,
  "totalHours": 3,
  "isSplit": true,
  "subOrderCount": 3
}
```

#### 获取主订单列表
```
GET /api/v1/user/order-groups?page=1&pageSize=10
```

#### 获取主订单详情
```
GET /api/v1/user/order-groups/:id
```

#### 获取子订单列表
```
GET /api/v1/user/order-groups/:id/sub-orders
```

#### 转单
```
POST /api/v1/user/order-groups/:subOrderId/transfer
{
  "newPlayerId": 2,
  "transferNote": "陪玩A有事，转给陪玩B"
}

Response:
{
  "success": true,
  "newSubOrderId": 105
}
```

#### 批量转单
```
POST /api/v1/user/order-groups/batch-transfer
{
  "subOrderIds": [102, 103],
  "newPlayerId": 2,
  "transferNote": "批量转单"
}
```

### 管理端

#### 获取主订单列表
```
GET /api/v1/admin/order-groups?userId=1&status=in_progress&page=1
```

#### 获取主订单详情
```
GET /api/v1/admin/order-groups/:id
```

#### 管理员转单
```
POST /api/v1/admin/order-groups/:subOrderId/transfer
```

#### 管理员批量转单
```
POST /api/v1/admin/order-groups/batch-transfer
```

## 状态流转

### 主订单状态
- `pending` - 待支付
- `paid` - 已支付，待接单
- `in_progress` - 进行中
- `completed` - 已完成
- `canceled` - 已取消
- `partial` - 部分完成（有转单）

### 子订单状态
保持原有 Order 状态不变

## 用户端展示

用户只看到 OrderGroup：
- 显示总价、总时长
- 显示进度（已完成 X/Y 小时）
- 如有转单，显示服务陪玩师列表

## 向后兼容

- 时长 ≤ 1 小时的订单不拆分，保持原有逻辑
- `orderGroups` 仓储未注入时，不拆分
- 原有 Order 相关接口保持兼容


## 测试

### 单元测试

测试文件位于 `api/internal/service/order/` 目录：

| 文件 | 覆盖内容 |
|------|----------|
| `creation_test.go` | 订单拆分构建、创建逻辑 |
| `transfer_test.go` | 转单、批量转单、可转单查询 |
| `orderService_test.go` | 基础订单服务 |

运行单元测试：
```bash
cd api && go test ./internal/service/order/... -v
```

### 集成测试

集成测试文件：`api/internal/service/order/integration_test.go`

测试用例：
- `TestIntegration_CreateOrderWithSplit` - 订单拆分创建完整流程
- `TestIntegration_CreateOrderNoSplit` - 短时长订单不拆分
- `TestIntegration_TransferSubOrder` - 转单完整流程
- `TestIntegration_BatchTransferSubOrders` - 批量转单
- `TestIntegration_GetTransferableSubOrders` - 获取可转单子订单
- `TestIntegration_OrderGroupStatusUpdate` - 主订单状态更新

运行集成测试（需要测试数据库）：
```bash
cd api && make test-integration-db
# 或手动指定环境变量
TEST_DB_HOST=localhost TEST_DB_PORT=5433 TEST_DB_NAME=gamelink_test \
  go test ./internal/service/order/... -tags=integration -v
```

### 测试覆盖率

当前覆盖率：**79.6%**（超过 70% 要求）

生成覆盖率报告：
```bash
cd api && make test-coverage
```
