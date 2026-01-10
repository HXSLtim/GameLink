# Zustand 最佳实践

本文档提供使用 Zustand 的最佳实践指南，帮助编写高性能、可维护的状态管理代码。

## 目录

1. [选择器模式](#1-选择器模式)
2. [状态设计原则](#2-状态设计原则)
3. [性能优化](#3-性能优化)
4. [错误处理](#4-错误处理)
5. [测试](#5-测试)
6. [常见陷阱](#6-常见陷阱)
7. [Service 层集成](#7-service-层集成)

---

## 1. 选择器模式

### 基础选择器

**推荐**：只订阅需要的状态切片

```typescript
// ✅ 好：只订阅需要的状态
const users = useUserStore(s => s.users);
const loading = useUserStore(s => s.loading);

// ❌ 不好：订阅整个 store（任何状态变化都会重渲染）
const { users, loading, error, filters, pagination } = useUserStore();
```

### 多字段选择

使用 `shallow` 比较避免不必要的重渲染：

```typescript
import { shallow } from 'zustand/shallow';

// ✅ 好：使用 shallow 比较多个字段
const { users, loading } = useUserStore(
  s => ({ users: s.users, loading: s.loading }),
  shallow
);

// ❌ 不好：每次都创建新对象
const { users, loading } = useUserStore(s => ({
  users: s.users,
  loading: s.loading
})); // 没有 shallow，每次都会重渲染
```

### 计算属性

对于需要计算的数据，使用 store 内部的 selector：

```typescript
// ✅ 好：在 store 内定义 selector
const getActiveUsers = useUserStore(s => s.getActiveUsers());

// ❌ 不好：在组件中计算
const activeUsers = useUserStore(s => s.users.filter(u => u.status === 'active'));
// 每次渲染都会重新计算
```

### 条件选择

只在需要时订阅状态：

```typescript
// ✅ 好：条件性订阅
const users = useUserStore(s => showModal ? s.users : []);

// ❌ 不好：总是订阅
const users = useUserStore(s => s.users);
```

---

## 2. 状态设计原则

### 适合全局状态的数据

**适合放 store 的**：
- API 响应缓存（用户列表、订单列表等）
- 跨页面共享状态（用户信息、权限）
- 全局 UI 状态（侧边栏展开状态、主题）
- WebSocket 连接状态
- 复杂的业务逻辑状态

**不适合放 store 的**：
- 组件本地 UI 状态（模态框显示、输入框内容）
- 表单临时数据
- URL 参数（使用 `react-router` 的 `useSearchParams`）
- 短暂的动画状态

### 示例对比

```typescript
// ✅ 适合全局状态
const { userInfo, isAuthenticated } = useAuthStore();
const { users, filters } = useUserStore();

// ❌ 不适合全局状态（应使用 useState）
const [modalVisible, setModalVisible] = useState(false);
const [formData, setFormData] = useState<FormData>({});
```

### 最小化状态

只存储最小化的数据，派生数据通过计算获得：

```typescript
// ✅ 好：只存储原始数据
interface State {
  users: User[];
  getUserById: (id: number) => User | undefined;
}

// ❌ 不好：存储冗余的派生数据
interface State {
  users: User[];
  activeUsers: User[];  // 可以从 users 派生
  adminUsers: User[];   // 可以从 users 派生
}
```

---

## 3. 性能优化

### 使用 React.memo

结合选择器使用 `React.memo` 避免不必要的重渲染：

```typescript
const UserCard = React.memo(({ userId }: { userId: number }) => {
  // 只在 userId 对应的用户变化时重渲染
  const user = useUserStore(s => s.getUserById(userId));

  if (!user) return null;
  return <div>{user.name}</div>;
});
```

### 批量更新

使用 `set` 的函数形式避免多次渲染：

```typescript
// ✅ 好：批量更新
set((state) => ({
  users: [...state.users, newUser],
  pagination: {
    ...state.pagination,
    total: state.pagination.total + 1
  }
}));

// ❌ 不好：多次调用 set
set({ users: [...state.users, newUser] });
set({ pagination: { ...state.pagination, total: state.pagination.total + 1 } });
```

### 持久化优化

限制持久化数据的大小：

```typescript
persist(
  (set, get) => ({
    users: [], // 可能有成千上万条数据
  }),
  {
    name: 'user-storage',
    partialize: (state) => ({
      // ✅ 只缓存最近 50 条
      users: state.users.slice(0, 50),
      filters: state.filters, // 保留筛选条件
    }),
  }
)
```

### 避免闭包陷阱

使用 `get()` 获取最新状态：

```typescript
// ✅ 好：使用 get() 获取最新状态
fetchUsers: async () => {
  set({ loading: true });
  try {
    const response = await api.getUsers();
    // 使用 get() 而不是闭包中的 state
    const currentFilters = get().filters;
    // ...
  } finally {
    set({ loading: false });
  }
}

// ❌ 不好：闭包中的旧状态
fetchUsers: async () => {
  const { filters } = get(); // 在 async 外捕获
  set({ loading: true });
  try {
    // filters 可能是旧值
    const response = await api.getUsers(filters);
  } finally {
    set({ loading: false });
  }
}
```

---

## 4. 错误处理

### 统一错误处理模式

在 store 中统一处理错误，组件只需关注 UI 反馈：

```typescript
// Store 实现
createUser: async (userData: Partial<User>) => {
  set({ loading: true, error: null });

  try {
    const response = await api.createUser(userData);
    set((state) => ({
      users: [response.data, ...state.users],
      loading: false,
    }));
  } catch (error) {
    const errorMessage = error.response?.data?.message || '创建用户失败';
    set({ error: errorMessage, loading: false });
    throw error; // 向上抛出，让组件决定是否显示提示
  }
}

// 组件使用
const { createUser, error } = useUserStore();

const handleSubmit = async () => {
  try {
    await createUser(formData);
    message.success('创建成功');
  } catch (err) {
    // 错误已在 store 中设置
    message.error(error || '创建失败');
  }
};
```

### 错误清理

在重试操作前清理旧错误：

```typescript
const { login, clearError } = useAuthStore();

const handleLogin = async () => {
  clearError(); // 清除旧错误
  try {
    await login(credentials);
  } catch (err) {
    // 显示错误
  }
};
```

### 错误边界

使用 React Error Boundary 捕获未预期的错误：

```typescript
<ErrorBoundary fallback={<ErrorPage />}>
  <UserList />
</ErrorBoundary>
```

---

## 5. 测试

### 单元测试 Store

使用 Vitest 测试 store 逻辑：

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { useUserStore } from './userStore';

describe('UserStore', () => {
  beforeEach(() => {
    // 每个测试前重置状态
    useUserStore.getState().reset();
  });

  it('should create user', async () => {
    const { createUser, users } = useUserStore.getState();

    await createUser({ name: '张三', email: 'test@example.com' });

    const updatedUsers = useUserStore.getState().users;
    expect(updatedUsers).toHaveLength(1);
    expect(updatedUsers[0].name).toBe('张三');
  });

  it('should filter users by status', () => {
    const state = useUserStore.getState();
    // 设置测试数据
    state.users = [
      { id: 1, name: 'User1', status: 'active' },
      { id: 2, name: 'User2', status: 'inactive' },
    ];

    const activeUsers = state.getActiveUsers();
    expect(activeUsers).toHaveLength(1);
    expect(activeUsers[0].status).toBe('active');
  });
});
```

### 测试选择器

```typescript
it('getUserById should return correct user', () => {
  const state = useUserStore.getState();
  state.users = [
    { id: 1, name: 'User1' },
    { id: 2, name: 'User2' },
  ];

  const user = state.getUserById(2);
  expect(user?.name).toBe('User2');
});
```

### 测试异步操作

```typescript
it('should handle API errors', async () => {
  // Mock API 调用
  vi.mock('@/api/admin', () => ({
    adminApi: {
      createUser: vi.fn().mockRejectedValue(new Error('Network error')),
    },
  }));

  const { createUser, error } = useUserStore.getState();

  await expect(createUser({})).rejects.toThrow('Network error');

  const updatedError = useUserStore.getState().error;
  expect(updatedError).toBeTruthy();
});
```

### 测试组件集成

```typescript
import { renderHook, act } from '@testing-library/react';
import { useUserStore } from './userStore';

it('should update users in component', async () => {
  const { result } = renderHook(() => useUserStore());

  await act(async () => {
    await result.current.fetchUsers();
  });

  expect(result.current.users).toBeDefined();
});
```

---

## 6. 常见陷阱

### 陷阱 1：订阅整个 store

```typescript
// ❌ 不好：任何状态变化都会重渲染
const store = useUserStore();

// ✅ 好：只订阅需要的状态
const users = useUserStore(s => s.users);
```

### 陷阱 2：在渲染中调用 getState()

```typescript
// ❌ 不好：违反 Hooks 规则
const users = useUserStore().getState().users;

// ✅ 好：使用选择器
const users = useUserStore(s => s.users);

// ✅ 或在事件处理中使用
const handleClick = () => {
  const users = useUserStore.getState().users;
};
```

### 陷阱 3：依赖闭包中的状态

```typescript
// ❌ 不好：闭包中的状态是旧值
useEffect(() => {
  const interval = setInterval(() => {
    console.log(users); // 永远是初始值
  }, 1000);
  return () => clearInterval(interval);
}, []); // 空依赖数组

// ✅ 好：使用 getState()
useEffect(() => {
  const interval = setInterval(() => {
    console.log(useUserStore.getState().users);
  }, 1000);
  return () => clearInterval(interval);
}, []);

// ✅ 或使用 subscribe()
useEffect(() => {
  const unsubscribe = useUserStore.subscribe(
    (state) => state.users,
    (users) => console.log(users)
  );
  return unsubscribe;
}, []);
```

### 陷阱 4：过度使用全局状态

```typescript
// ❌ 不好：模态框状态不需要全局
const { modalVisible, setModalVisible } = useModalStore();

// ✅ 好：本地状态足够
const [modalVisible, setModalVisible] = useState(false);
```

### 陷阱 5：持久化敏感信息

```typescript
// ❌ 不好：持久化 token 到 localStorage（安全风险）
persist(
  (set) => ({ token: null }),
  { storage: createJSONStorage(() => localStorage) }
)

// ✅ 好：使用 sessionStorage
persist(
  (set) => ({ token: null }),
  { storage: createJSONStorage(() => sessionStorage) }
)
```

---

## 性能检查清单

在提交代码前，检查以下项目：

- [ ] 是否只订阅了需要的状态切片？
- [ ] 是否使用了 `shallow` 比较多字段？
- [ ] 是否使用了 `React.memo` 优化列表渲染？
- [ ] 是否限制了持久化数据的大小？
- [ ] 是否避免了闭包陷阱？
- [ ] 异步操作中是否使用了 `get()` 获取最新状态？
- [ ] 是否正确清理了错误状态？
- [ ] 组件本地状态是否误用了全局状态？

---

## 参考资源

- [Zustand 官方文档](https://github.com/pmndrs/zustand)
- [Zustand 最佳实践](https://docs.pmnd.rs/zustand/guides/performance)
- [React 性能优化](https://react.dev/learn/render-and-commit)
- [项目 README.md](./README.md)
- [TypeScript 类型指南](./TYPE_GUIDE.md)
- [Service 层文档](../services/README.md)

---

## 7. Service 层集成

Phase 3 引入了 Service 层架构，将业务逻辑从 Store 中抽离。本节介绍如何正确集成 Service 层。

### 架构概览

```
UI Components → Zustand Stores → Domain Services → API Layer
                    ↓
              (UI State, Cache)    (Business Logic)
```

**职责分离**：
- **Store**: 管理 UI 状态、缓存、分页
- **Service**: 处理业务逻辑、数据验证、API 编排

### 导入 Service

```typescript
// ✅ 推荐：从统一入口导入
import { userService, orderService, playerService } from '@/services';

// ✅ 也可以：从具体模块导入
import { userService } from '@/services/domain/userService';

// ❌ 不推荐：直接导入类（除非需要自定义实例）
import { UserService } from '@/services/domain/userService';
```

### Store 中使用 Service

**推荐模式**：Store 调用 Service 方法，处理结果并更新状态

```typescript
// ✅ 好：Store 使用 Service 进行业务逻辑处理
import { userService } from '@/services';

export const useUserStore = create<UserState>()((set, get) => ({
  users: [],
  loading: false,
  error: null,

  // 使用 Service 进行数据验证
  createUser: async (userData: CreateUserDto) => {
    // 1. 先用 Service 验证数据
    const validation = userService.validateUserData(userData);
    if (!validation.valid) {
      set({ error: validation.errors.map(e => e.message).join(', ') });
      return;
    }

    set({ loading: true, error: null });
    try {
      // 2. 调用 Service 创建用户
      const result = await userService.createUser(userData);
      if (result.success && result.data) {
        set((state) => ({
          users: [result.data!, ...state.users],
          loading: false,
        }));
      } else {
        set({ error: result.error?.message || '创建失败', loading: false });
      }
    } catch (error) {
      set({ error: '创建用户失败', loading: false });
    }
  },

  // 使用 Service 进行批量操作
  batchUpdateStatus: async (userIds: number[], status: string) => {
    set({ loading: true, error: null });
    try {
      const result = await userService.batchUpdateStatus(userIds, status);
      if (result.success) {
        // 更新本地状态
        set((state) => ({
          users: state.users.map(user =>
            userIds.includes(user.id) ? { ...user, status } : user
          ),
          loading: false,
        }));
      } else {
        // 处理部分失败
        const failedIds = result.results
          .filter(r => !r.success)
          .map(r => r.index);
        set({
          error: `${result.failed} 个用户更新失败`,
          loading: false,
        });
      }
    } catch (error) {
      set({ error: '批量更新失败', loading: false });
    }
  },
}));
```

### Service 结果处理

Service 方法返回 `ServiceResult<T>` 或 `BatchResult<T>`，需要正确处理：

```typescript
import type { ServiceResult, BatchResult } from '@/services';

// 处理单个操作结果
const handleSingleResult = async () => {
  const result: ServiceResult<User> = await userService.getUserById(1);

  if (result.success && result.data) {
    // 成功：使用 result.data
    console.log('User:', result.data);
  } else {
    // 失败：使用 result.error
    console.error('Error:', result.error?.message);
    console.error('Code:', result.error?.code);
  }
};

// 处理批量操作结果
const handleBatchResult = async () => {
  const result: BatchResult<void> = await userService.batchDelete([1, 2, 3]);

  console.log(`Total: ${result.total}`);
  console.log(`Succeeded: ${result.succeeded}`);
  console.log(`Failed: ${result.failed}`);

  // 处理每个项目的结果
  result.results.forEach((item) => {
    if (item.success) {
      console.log(`Item ${item.index} succeeded`);
    } else {
      console.error(`Item ${item.index} failed:`, item.error?.message);
    }
  });
};
```

### 数据验证

**推荐**：在 Store 调用 API 前使用 Service 验证数据

```typescript
// ✅ 好：使用 Service 验证
const createUser = async (userData: CreateUserDto) => {
  // 验证用户数据
  const validation = userService.validateUserData(userData);
  if (!validation.valid) {
    set({ error: validation.errors.map(e => e.message).join(', ') });
    return;
  }

  // 验证通过，继续创建
  const result = await userService.createUser(userData);
  // ...
};

// ✅ 好：单独验证字段
const validateEmail = (email: string) => {
  return userService.validateEmail(email);
};

const validatePassword = (password: string) => {
  const result = userService.validatePassword(password);
  return result; // { valid: boolean, errors: string[] }
};
```

### 数据导出

使用 Service 进行数据导出：

```typescript
import { userService } from '@/services';

const exportUsers = () => {
  const users = useUserStore.getState().users;
  const exportData = userService.exportUsers(users);

  // exportData.headers: string[] - 表头
  // exportData.rows: string[][] - 数据行

  // 使用 xlsx 或其他库生成文件
  const worksheet = XLSX.utils.aoa_to_sheet([
    exportData.headers,
    ...exportData.rows,
  ]);
  // ...
};
```

### 数据导入

使用 ImportService 进行批量数据导入：

```typescript
import { importService, parseFile } from '@/services';

const handleImport = async (file: File) => {
  // 1. 解析文件
  const preview = await importService.parseFile(file, 'user');

  if (preview.structureErrors.length > 0) {
    // 显示结构错误
    console.error('Structure errors:', preview.structureErrors);
    return;
  }

  // 2. 显示预览（validRows, invalidRows）
  console.log('Valid rows:', preview.validRows.length);
  console.log('Invalid rows:', preview.invalidRows.length);

  // 3. 用户确认后执行导入
  const result = await importService.importUsers(preview.validRows);

  console.log(`Imported: ${result.importedCount}`);
  console.log(`Skipped: ${result.skippedCount}`);
};
```

### 订单操作

使用 OrderService 进行订单相关操作：

```typescript
import { orderService } from '@/services';

// 检查是否可以取消订单
const canCancel = orderService.canCancel(order);
if (!canCancel.allowed) {
  message.error(canCancel.reason);
  return;
}

// 计算退款金额
const refundCalc = orderService.calculateRefund(order, requestedAmount);
console.log('Refund amount:', refundCalc.refundAmount);
console.log('Platform fee:', refundCalc.platformFee);
console.log('Player amount:', refundCalc.playerAmount);

// 计算订单统计
const stats = orderService.computeStatistics(orders);
console.log('Total revenue:', stats.totalRevenue);
console.log('Completion rate:', stats.completionRate);
```

### 陪玩师操作

使用 PlayerService 进行陪玩师相关操作：

```typescript
import { playerService } from '@/services';

// 检查是否可以验证
const canVerify = playerService.canVerify(player, 'verified');
if (!canVerify.allowed) {
  message.error(canVerify.reason);
  return;
}

// 计算收益
const earnings = playerService.calculateEarnings(order, commissionRule);
console.log('Gross amount:', earnings.grossAmount);
console.log('Commission:', earnings.commissionAmount);
console.log('Net amount:', earnings.netAmount);

// 计算统计
const stats = playerService.computeStatistics(playerId, orders);
console.log('Total earnings:', stats.totalEarnings);
console.log('Average rating:', stats.averageRating);
```

### 错误处理

Service 层使用统一的错误格式：

```typescript
import { ServiceErrorCodes, ServiceException } from '@/services';

try {
  const result = await userService.createUser(userData);
  if (!result.success) {
    // 根据错误码处理
    switch (result.error?.code) {
      case ServiceErrorCodes.USER_EMAIL_EXISTS:
        message.error('邮箱已存在');
        break;
      case ServiceErrorCodes.USER_PHONE_EXISTS:
        message.error('手机号已存在');
        break;
      case ServiceErrorCodes.VALIDATION_ERROR:
        message.error(result.error.message);
        break;
      default:
        message.error('创建失败');
    }
  }
} catch (error) {
  if (error instanceof ServiceException) {
    console.error('Service error:', error.code, error.message);
  }
}
```

### 日志和性能监控

Service 层内置日志和性能监控：

```typescript
import { createServiceLogger, createPerformanceMonitor } from '@/services';

// 创建自定义 logger
const logger = createServiceLogger('MyService', {
  level: 'debug',
  enableConsole: true,
});

logger.info('Operation started', { userId: 1 });
logger.error('Operation failed', new Error('Something went wrong'));

// 性能监控
const monitor = createPerformanceMonitor(logger, {
  slowThresholdMs: 3000,
  maxMetrics: 1000,
});

const stopTimer = monitor.startTimer('fetchUsers');
// ... 执行操作
const duration = stopTimer();
console.log(`Operation took ${duration}ms`);
```

### Service 层检查清单

在使用 Service 层时，检查以下项目：

- [ ] 是否从统一入口 `@/services` 导入？
- [ ] 是否正确处理 `ServiceResult` 的 success/error？
- [ ] 是否在调用 API 前使用 Service 验证数据？
- [ ] 批量操作是否处理了部分失败的情况？
- [ ] 是否根据错误码提供用户友好的错误信息？
- [ ] 是否使用 Service 的计算方法而非在组件中计算？
- [ ] 导入操作是否显示了预览和确认步骤？
