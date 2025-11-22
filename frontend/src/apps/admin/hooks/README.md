# Admin Hooks

## 📁 目录用途

该目录用于存放管理后台（Admin）专用的自定义 React Hooks。

## 🎯 设计理念

Hooks 用于提取组件之间的可复用逻辑，特别是涉及状态和副作用的逻辑。

## 📦 Hook 分类

### 数据获取 Hooks
- `useAdminUsers` - 获取管理员列表
- `useAdminPlayers` - 获取陪玩师列表
- `useAdminOrders` - 获取订单列表和统计
- `useAdminGames` - 获取游戏配置
- `useAdminRevenue` - 获取收益统计
- `useAdminPermissions` - 获取权限列表
- `useAdminRoles` - 获取角色列表
- `useAdminReviews` - 获取评价列表

### 业务逻辑 Hooks
- `useUserBan` - 用户封禁操作
- `usePlayerVerify` - 陪玩师审核
- `useOrderManage` - 订单管理操作
- `useGameConfig` - 游戏配置管理
- `useRevenueStats` - 收益统计计算
- `usePermissionAssign` - 权限分配
- `useRoleManage` - 角色管理
- `useContentReview` - 内容审核

### UI 交互 Hooks
- `useAdminSearch` - 后台搜索功能
- `useAdminFilter` - 筛选逻辑
- `useAdminSort` - 排序逻辑
- `useAdminPagination` - 分页逻辑
- `useAdminExport` - 数据导出
- `useBulkActions` - 批量操作

### Form Hooks
- `useUserForm` - 用户表单逻辑
- `usePlayerForm` - 陪玩师表单逻辑
- `useGameForm` - 游戏配置表单
- `useSystemConfigForm` - 系统设置表单

## 📋 命名规范

### Hook 命名规则
- 必须以 `use` 开头
- 使用 camelCase（小驼峰）
- 命名应具有描述性，例如：`useUserBan`, `useOrderStats`

```typescript
// ✅ 推荐
export const useAdminUsers = () => { ... }
export const useOrderStats = () => { ... }

// ❌ 避免
export const useData = () => { ... }
export const useAdminHook = () => { ... }
```

### 返回值规范

#### 数据获取 Hook
```typescript
interface UseDataResult<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

function useSomeData(): UseDataResult<DataType> {
  // ...
}
```

#### 操作 Hook
```typescript
interface UseOperationResult {
  execute: (...args: any[]) => Promise<any>;
  loading: boolean;
  error: Error | null;
}

function useSomeAction(): UseOperationResult {
  // ...
}
```

## 🎯 开发规范

### TypeScript 类型定义

每个 hook 都应该有明确的类型定义：

```typescript
// ✅ 推荐
interface UseAdminUsersOptions {
  page?: number;
  pageSize?: number;
  status?: 'active' | 'banned' | 'pending';
}

interface UseAdminUsersResult {
  users: User[];
  total: number;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
  changePage: (page: number) => void;
}

export function useAdminUsers(options: UseAdminUsersOptions = {}): UseAdminUsersResult {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // ... hook 逻辑

  return { users, total, loading, error, refetch, changePage };
}
```

### 错误处理

所有涉及异步操作的 hook 都应该有良好的错误处理：

```typescript
// ✅ 推荐
export function useAdminUsers() {
  const [error, setError] = useState<Error | null>(null);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getUsers();
      setUsers(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('获取用户失败'));
      // 显示错误通知
      showNotification('error', '错误', '获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  return { error, /* ... */ };
}
```

### 依赖管理

正确使用 useEffect 的依赖数组：

```typescript
// ✅ 正确
export function useAdminUser(userId: number) {
  useEffect(() => {
    fetchUser(userId);
  }, [userId]); // 明确依赖
}

// ❌ 避免
export function useAdminUser(userId: number) {
  useEffect(() => {
    fetchUser(userId);
  }, []); // userId 变化时不会重新获取
}
```

## 💡 最佳实践

### 1. 保持 Hook 纯粹

Hook 应该只处理逻辑，不包含 JSX：

```typescript
// ✅ 推荐 - 纯逻辑
export function useUserStatus(userId: number) {
  const [status, setStatus] = useState<'online' | 'offline' | 'busy'>('offline');

  useEffect(() => {
    const unsub = subscribeUserStatus(userId, setStatus);
    return unsub;
  }, [userId]);

  return status;
}

// 在组件中使用
function UserBadge({ userId }: { userId: number }) {
  const status = useUserStatus(userId);
  return <span className={`status-${status}`}>{status}</span>;
}

// ❌ 避免 - 包含 JSX
export function useUserStatus(userId: number) {
  const [status, setStatus] = useState('offline');

  // ...

  // 返回 JSX - 这使得 hook 不灵活
  return <span className={`status-${status}`}>{status}</span>;
}
```

### 2. 合理拆分 Hook

复杂的逻辑应该拆分为多个小 hook：

```typescript
// ✅ 推荐 - 拆分小的 hook
export function useUserList() { /* ... */ }
export function useUserSearch() { /* ... */ }
export function useUserFilter() { /* ... */ }

// 在组件中组合使用
function UserManagement() {
  const { users } = useUserList();
  const { searchTerm, setSearchTerm } = useUserSearch();
  const { filteredUsers } = useUserFilter(users, searchTerm);

  // ...
}

// ❌ 避免 - 一个巨大的 hook
export function useUserManagement() {
  // 包含列表、搜索、筛选、分页、排序所有逻辑
  // ... 几百行代码
}
```

### 3. 使用自定义 Hook 封装重复逻辑

```typescript
// ✅ 推荐 - 封装重复的 CRUD 逻辑
export function useCRUD<T>(resource: string) {
  const [items, setItems] = useState<T[]>([]);
  const [loading, setLoading] = useState(false);

  const create = async (data: Partial<T>) => { /* ... */ };
  const update = async (id: number, data: Partial<T>) => { /* ... */ };
  const remove = async (id: number) => { /* ... */ };
  const fetchAll = async () => { /* ... */ };

  return { items, loading, create, update, remove, fetchAll };
}

// 使用
const userCRUD = useCRUD<User>('users');
const orderCRUD = useCRUD<Order>('orders');
```

### 4. 处理异步操作

```typescript
// ✅ 推荐 - 提供加载状态和错误信息
export function useAsyncOperation<T extends (...args: any[]) => Promise<any>>(
  asyncFunction: T
) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const execute = async (...args: Parameters<T>): Promise<ReturnType<T> | null> => {
    try {
      setLoading(true);
      setError(null);
      const result = await asyncFunction(...args);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err : new Error('操作失败'));
      showNotification('error', '错误', '操作失败');
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { execute, loading, error };
}

// 使用
const { execute: banUser, loading: banning } = useAsyncOperation(api.banUser);
```

### 5. 使用 useCallback 优化性能

```typescript
// ✅ 推荐 - 使用 useCallback
export function useAdminSearch() {
  const [searchTerm, setSearchTerm] = useState('');

  const debouncedSearch = useCallback(
    debounce((value: string) => {
      performSearch(value);
    }, 300),
    []
  );

  const handleSearch = useCallback((value: string) => {
    setSearchTerm(value);
    debouncedSearch(value);
  }, [debouncedSearch]);

  return { searchTerm, handleSearch };
}
```

## 📚 示例代码

### 基础数据获取 Hook

```typescript
// useAdminUsers.ts
import { useState, useEffect, useCallback } from 'react';
import { api } from '@/api';
import type { User, PaginationParams } from '@/api';

export interface UseAdminUsersOptions extends PaginationParams {
  status?: 'active' | 'banned' | 'pending';
  search?: string;
}

export interface UseAdminUsersResult {
  users: User[];
  total: number;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
  changePage: (page: number) => void;
  changePageSize: (pageSize: number) => void;
}

export function useAdminUsers(
  options: UseAdminUsersOptions = {}
): UseAdminUsersResult {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await api.admin.getUsers(options);

      setUsers(response.data);
      setTotal(response.total);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('获取用户列表失败');
      setError(error);
      // 显示错误通知
      showNotification('error', '错误', error.message);
    } finally {
      setLoading(false);
    }
  }, [options.page, options.pageSize, options.status, options.search]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const changePage = useCallback((page: number) => {
    // 处理页码变化
  }, []);

  const changePageSize = useCallback((pageSize: number) => {
    // 处理每页数量变化
  }, []);

  const refetch = useCallback(() => {
    return fetchUsers();
  }, [fetchUsers]);

  return {
    users,
    total,
    loading,
    error,
    refetch,
    changePage,
    changePageSize,
  };
}
```

### 业务操作 Hook

```typescript
// useUserBan.ts
import { useState, useCallback } from 'react';
import { api } from '@/api';
import { useNotification } from '@/stores';

export interface UseUserBanResult {
  banUser: (userId: number, reason: string) => Promise<boolean>;
  unbanUser: (userId: number) => Promise<boolean>;
  loading: boolean;
}

export function useUserBan(): UseUserBanResult {
  const [loading, setLoading] = useState(false);
  const showNotification = useNotification();

  const banUser = useCallback(async (userId: number, reason: string) => {
    setLoading(true);
    try {
      await api.admin.banUser(userId, reason);
      showNotification('success', '操作成功', '用户已封禁');
      return true;
    } catch (error) {
      const message = error instanceof Error ? error.message : '封禁用户失败';
      showNotification('error', '操作失败', message);
      return false;
    } finally {
      setLoading(false);
    }
  }, [showNotification]);

  const unbanUser = useCallback(async (userId: number) => {
    setLoading(true);
    try {
      await api.admin.unbanUser(userId);
      showNotification('success', '操作成功', '用户已解封');
      return true;
    } catch (error) {
      const message = error instanceof Error ? error.message : '解封用户失败';
      showNotification('error', '操作失败', message);
      return false;
    } finally {
      setLoading(false);
    }
  }, [showNotification]);

  return { banUser, unbanUser, loading };
}
```

## 📝 注意事项

1. **避免直接导出其它 Hook** - 如果需要组合，在组件中进行
2. **保持 Hook 独立** - 一个 Hook 不依赖另一个 Hook
3. **提供良好的 TypeScript 支持** - 完整的类型定义
4. **测试 Hook** - 使用 @testing-library/react-hooks 测试
5. **文档注释** - 复杂 Hook 需要说明使用场景

---

**最后更新**: 2025-11-22
**维护者**: GameLink 前端团队
**相关文档**: [全局 Hooks 目录](../../shared/hooks/)
