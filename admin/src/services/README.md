# Service Layer - 服务层

## 概述

Service 层是 Phase 3 引入的业务逻辑层，负责将复杂的业务逻辑从 UI 组件和 Zustand Store 中抽离，提升代码可维护性和可测试性。

## 架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        UI Components                             │
│  (Pages, Modals, Forms)                                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Zustand Stores                              │
│  (UI State, Cache, Pagination)                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Domain Services                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐    │
│  │UserService│  │OrderService│ │PlayerService│ │ImportService│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘    │
│                                                                  │
│  - Business Logic                                                │
│  - Data Validation                                               │
│  - Data Transformation                                           │
│  - Error Handling                                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer                                 │
│  (adminApi, userApi, etc.)                                      │
└─────────────────────────────────────────────────────────────────┘
```

## 目录结构

```
admin/src/services/
├── index.ts                     # 统一导出入口
├── init.ts                      # 应用初始化服务
├── README.md                    # 本文档
├── domain/                      # 领域服务
│   ├── base.ts                  # 基础服务类
│   ├── userService.ts           # 用户领域服务
│   ├── orderService.ts          # 订单领域服务
│   ├── playerService.ts         # 陪玩师领域服务
│   └── index.ts                 # 统一导出
├── import/                      # 数据导入服务
│   ├── importService.ts         # 导入核心服务
│   ├── parsers/                 # 文件解析器
│   ├── validators/              # 数据验证器
│   ├── templates/               # 导入模板
│   ├── history/                 # 导入历史
│   ├── transaction/             # 事务管理
│   └── index.ts
└── utils/                       # 服务工具
    ├── serviceError.ts          # 统一错误类
    ├── serviceResult.ts         # 统一结果类
    ├── logger.ts                # 日志服务
    ├── performance.ts           # 性能监控
    ├── concurrency.ts           # 并发控制
    └── index.ts
```

## 快速开始

### 导入服务

```typescript
// 推荐：从统一入口导入
import {
  userService,
  orderService,
  playerService,
  importService,
} from '@/services';

// 导入类型
import type {
  ServiceResult,
  BatchResult,
  UserValidationResult,
} from '@/services';
```

### 基本用法

```typescript
import { userService } from '@/services';

// 获取用户
const result = await userService.getUserById(1);
if (result.success) {
  console.log('User:', result.data);
} else {
  console.error('Error:', result.error?.message);
}

// 验证数据
const validation = userService.validateUserData({
  email: 'test@example.com',
  phone: '13800138000',
});
if (!validation.valid) {
  console.error('Validation errors:', validation.errors);
}

// 批量操作
const batchResult = await userService.batchUpdateStatus([1, 2, 3], 'active');
console.log(`Succeeded: ${batchResult.succeeded}, Failed: ${batchResult.failed}`);
```

## 领域服务

### UserService - 用户服务

```typescript
import { userService } from '@/services';

// CRUD 操作
await userService.getUsers(params);
await userService.getUserById(id);
await userService.createUser(data);
await userService.updateUser(id, data);
await userService.deleteUser(id);

// 状态和角色
await userService.updateUserStatus(id, status);
await userService.updateUserRole(id, role);

// 批量操作
await userService.batchUpdateStatus(userIds, status);
await userService.batchUpdateRole(userIds, role);
await userService.batchDelete(userIds);

// 验证
userService.validateUserData(data);
userService.validateEmail(email);
userService.validatePhone(phone);
userService.validatePassword(password);

// 导出
userService.exportUsers(users);
```

### OrderService - 订单服务

```typescript
import { orderService } from '@/services';

// 查询
await orderService.getOrders(params);
await orderService.getOrderById(id);

// 操作
await orderService.cancelOrder(id, reason);
await orderService.refundOrder(id, data);

// 批量操作
await orderService.batchCancel(orderIds, reason);
await orderService.batchComplete(orderIds);

// 验证和计算
orderService.canCancel(order);
orderService.canRefund(order);
orderService.calculateRefund(order, amount);

// 统计
orderService.computeStatistics(orders);
orderService.computeTrend(orders, days);
```

### PlayerService - 陪玩师服务

```typescript
import { playerService } from '@/services';

// CRUD
await playerService.getPlayers(params);
await playerService.getPlayerById(id);
await playerService.createPlayer(data);
await playerService.updatePlayer(id, data);
await playerService.deletePlayer(id);

// 验证
await playerService.verifyPlayer(id, status, remark);
playerService.canVerify(player, newStatus);

// 技能标签
await playerService.updateSkillTags(id, tags);
playerService.validateSkillTags(tags);

// 批量操作
await playerService.batchUpdateStatus(playerIds, status);
await playerService.batchDelete(playerIds);

// 收益计算
playerService.calculateEarnings(order, commissionRule);
playerService.computeStatistics(playerId, orders);
```

## 数据导入服务

### ImportService - 导入服务

```typescript
import { importService, parseFile } from '@/services';

// 获取模板
const template = importService.getTemplate('user');
const blob = importService.downloadTemplate('user');

// 解析文件
const preview = await importService.parseFile(file, 'user');

// 执行导入
const result = await importService.importUsers(preview.validRows);
const result = await importService.importPlayers(preview.validRows);
const result = await importService.importGames(preview.validRows);

// 导入历史
const history = await importHistoryService.getHistory(params);
const details = await importHistoryService.getDetails(id);
```

## 服务工具

### 错误处理

```typescript
import { ServiceException, ServiceErrorCodes } from '@/services';

// 错误码
ServiceErrorCodes.VALIDATION_ERROR
ServiceErrorCodes.NOT_FOUND
ServiceErrorCodes.USER_EMAIL_EXISTS
ServiceErrorCodes.ORDER_CANNOT_CANCEL

// 抛出错误
throw new ServiceException(
  ServiceErrorCodes.VALIDATION_ERROR,
  '验证失败',
  { field: 'email' }
);
```

### 日志服务

```typescript
import { createServiceLogger } from '@/services';

const logger = createServiceLogger('MyService');
logger.debug('Debug message', { data });
logger.info('Info message');
logger.warn('Warning message');
logger.error('Error message', error);
```

### 性能监控

```typescript
import { createPerformanceMonitor } from '@/services';

const monitor = createPerformanceMonitor(logger);
const stopTimer = monitor.startTimer('operation');
// ... 执行操作
const duration = stopTimer();
```

### 并发控制

```typescript
import { ConcurrencyController, chunkArray } from '@/services';

const controller = new ConcurrencyController({ maxConcurrent: 5 });

// 防止重复提交
await controller.withDeduplication('key', async () => {
  // 操作
});

// 并发处理
const results = await controller.processWithConcurrency(items, async (item) => {
  // 处理每个项目
});

// 分块处理
const chunks = chunkArray(items, 50);
```

## 类型定义

### ServiceResult

```typescript
interface ServiceResult<T> {
  success: boolean;
  data?: T;
  error?: ServiceError;
}
```

### BatchResult

```typescript
interface BatchResult<T> {
  success: boolean;
  total: number;
  succeeded: number;
  failed: number;
  results: Array<{
    index: number;
    success: boolean;
    data?: T;
    error?: ServiceError;
  }>;
}
```

### ServiceError

```typescript
interface ServiceError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
  originalError?: Error;
}
```

## 测试

Service 层的测试文件与源文件放在同一目录：

```
domain/
├── userService.ts
├── userService.test.ts
├── orderService.ts
├── orderService.test.ts
└── ...
```

运行测试：

```bash
cd admin
npm run test -- --run src/services
```

## 相关文档

- [Store 最佳实践](../stores/BEST_PRACTICES.md) - 包含 Service 层集成指南
- [设计文档](../../../.kiro/specs/admin-phase3-improvements/design.md) - 详细设计规范
- [需求文档](../../../.kiro/specs/admin-phase3-improvements/requirements.md) - 需求规范
