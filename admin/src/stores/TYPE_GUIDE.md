# TypeScript 类型指南

本文档提供 Zustand stores 的 TypeScript 使用指南，确保类型安全和良好的开发体验。

## 目录

1. [类型定义位置](#1-类型定义位置)
2. [Store 类型定义](#2-store-类型定义)
3. [类型导入](#3-类型导入)
4. [泛型使用](#4-泛型使用)
5. [类型推断](#5-类型推断)
6. [常见类型问题](#6-常见类型问题)

---

## 1. 类型定义位置

### 通用类型定义

所有共享类型定义在 `types/index.ts` 中：

```typescript
// admin/src/stores/types/index.ts

// ============================================
// Common Types
// ============================================

export interface UserInfo {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  avatar?: string;
  role: string;
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

export interface MenuItem {
  id: number;
  key: string;
  label: string;
  icon?: string;
  path?: string;
  parentId?: number | null;
  sort: number;
  permission?: string;
  children?: MenuItem[];
}

// User types
export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  status: 'active' | 'inactive' | 'blocked';
  role: string;
  createdAt: string;
}

// Order types
export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  playerIds: number[];
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled' | 'refunded';
  amount: number;
  duration: number;
  createdAt: string;
}

// Player types
export interface Player {
  id: number;
  userId: number;
  nickname: string;
  avatar: string;
  rank: number;
  level: number;
  pricePerHour: number;
  status: 'available' | 'busy' | 'offline';
  tags: string[];
}

// Chat types
export interface ChatMessage {
  id: number;
  orderId: number;
  senderId: number;
  senderType: 'user' | 'player' | 'admin';
  content: string;
  type: 'text' | 'image' | 'voice';
  createdAt: string;
  read: boolean;
}

export interface ChatRoom {
  id: number;
  orderId: number;
  participants: number[];
  lastMessage?: ChatMessage;
  unreadCount: number;
}
```

---

## 2. Store 类型定义

### State 接口

每个 store 定义自己的 State 接口：

```typescript
// admin/src/stores/modules/authStore.ts

interface AuthState {
  // ========== State ==========
  token: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;

  // ========== Actions ==========
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  setToken: (token: string) => void;
  setUserInfo: (user: UserInfo) => void;
  clearError: () => void;
  setLoading: (loading: boolean) => void;

  // ========== Selectors ==========
  isAdmin: () => boolean;
  hasPermission: (permission: string | string[], mode?: 'any' | 'all') => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasRole: (role: string | string[]) => boolean;
}
```

### Store 创建

使用泛型创建类型安全的 store：

```typescript
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 实现代码
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);
```

### 导出 State 类型

供外部使用的类型导出：

```typescript
// 导出完整的 State 类型
export type { AuthState };

// 导出便捷 hooks（自动推断返回类型）
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
```

---

## 3. 类型导入

### 导入类型

使用 `import type` 只导入类型：

```typescript
// ✅ 好：只导入类型
import type { UserInfo, LoginRequest, LoginResponse } from '@/stores/types';

// ❌ 不好：导入值（虽然是类型，但混淆意图）
import { UserInfo, LoginRequest } from '@/stores/types';
```

### 从 API 导入类型

复用 API 层的类型定义：

```typescript
import { adminApi, type User as ApiUser } from '@/api/admin';

interface UserState {
  users: ApiUser[];  // 使用 API 类型
  updateUser: (id: number, data: Partial<ApiUser>) => Promise<void>;
}
```

### 从 Store 导入类型

导入其他 store 的类型：

```typescript
import type { AuthState } from '@/stores/modules/authStore';

// 使用 AuthState 类型
const checkAuth = (auth: AuthState) => {
  return auth.isAuthenticated;
};
```

---

## 4. 泛型使用

### 函数泛型

在 selector 中使用泛型：

```typescript
interface UserState {
  // 泛型 selector
  getUsersByStatus: <T extends User['status']>(status: T) => User[];
}

// 实现
getUsersByStatus: <T extends User['status']>(status: T) => {
  return get().users.filter((u) => u.status === status);
}
```

### API 调用泛型

处理 API 响应：

```typescript
import type { ApiResponse } from '@/api/client';

interface UserState {
  fetchUsers: () => Promise<ApiResponse<User[]>>;
}
```

---

## 5. 类型推断

### 自动推断返回类型

Zustand 会自动推断选择器的返回类型：

```typescript
// 自动推断为 string | null
const token = useAuthStore((state) => state.token);

// 自动推断为 boolean
const isAdmin = useAuthStore((state) => state.isAdmin());

// 自动推断为 UserInfo | null
const userInfo = useAuthStore((state) => state.userInfo);
```

### 手动指定类型

需要时手动指定类型：

```typescript
interface UserFilters {
  status?: string;
  role?: string;
  keyword?: string;
}

const filters: UserFilters = useUserStore((state) => state.filters);
```

---

## 6. 常见类型问题

### 问题 1：可选链的类型推断

```typescript
// ❌ 类型可能是 undefined
const userName = useAuthStore((state) => state.userInfo?.name);

// ✅ 提供默认值
const userName = useAuthStore((state) => state.userInfo?.name ?? '');
```

### 问题 2：数组的类型保护

```typescript
// ❌ 可能是 undefined
const activeUsers = useUserStore((state) => state.getActiveUsers());

// ✅ 确保返回数组
interface UserState {
  getActiveUsers: () => User[];  // 总是返回数组
}

// 使用时
const activeUsers = useUserStore((state) => state.getActiveUsers());
// 类型保证是 User[]
```

### 问题 3：Partial 类型

更新数据时使用 `Partial`：

```typescript
interface UserState {
  updateUser: (id: number, data: Partial<User>) => Promise<void>;
}

// 调用时可以只提供部分字段
await updateUser(1, { name: '新名字' });
```

### 问题 4：枚举类型 vs 字符串字面量

```typescript
// ✅ 推荐：字符串字面量（更灵活）
interface User {
  status: 'active' | 'inactive' | 'blocked';
}

// ❌ 不推荐：枚举（除非有特殊需求）
enum UserStatus {
  Active = 'active',
  Inactive = 'inactive',
  Blocked = 'blocked',
}
```

### 问题 5：persist 的类型问题

使用 `createJSONStorage` 避免类型错误：

```typescript
import { createJSONStorage } from 'zustand/middleware';

persist(
  (set, get) => ({
    // ...
  }),
  {
    name: 'user-storage',
    storage: createJSONStorage(() => localStorage),  // ✅ 类型安全
    // storage: localStorage,  // ❌ 类型错误
  }
)
```

---

## 类型检查清单

在编写 store 代码时，确保：

- [ ] 所有状态字段都有明确的类型定义
- [ ] Actions 的参数和返回值都有类型
- [ ] Selectors 的返回类型正确
- [ ] 使用 `import type` 导入类型
- [ ] 复用 API 层的类型定义
- [ ] 可选字段正确使用 `?` 标记
- [ ] 数组类型不包含 `undefined`
- [ ] 使用 `Partial` 处理部分更新
- [ ] 导出 State 类型供外部使用
- [ ] 使用 TypeScript 严格模式

---

## 示例：完整类型定义

```typescript
// ============================================
// admin/src/stores/modules/userStore.ts
// ============================================

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { adminApi, type User as ApiUser } from '@/api/admin';

// 分页参数
interface Pagination {
  current: number;
  pageSize: number;
  total: number;
}

// 筛选条件
interface UserFilters {
  status?: string;
  role?: string;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}

// State 接口
interface UserState {
  // State
  users: ApiUser[];
  loading: boolean;
  error: string | null;
  pagination: Pagination;
  filters: UserFilters;

  // Actions
  fetchUsers: (page?: number, pageSize?: number) => Promise<void>;
  createUser: (userData: Partial<ApiUser>) => Promise<void>;
  updateUser: (id: number, userData: Partial<ApiUser>) => Promise<void>;
  deleteUser: (id: number) => Promise<void>;
  batchDeleteUsers: (userIds: number[]) => Promise<void>;
  updateUserStatus: (id: number, status: string) => Promise<void>;
  batchUpdateUserStatus: (userIds: number[], status: string) => Promise<void>;
  updateUserRole: (id: number, role: string) => Promise<void>;
  batchUpdateUserRole: (userIds: number[], role: string) => Promise<void>;
  setFilters: (newFilters: Partial<UserFilters>) => void;
  clearFilters: () => void;
  reset: () => void;

  // Selectors
  getUserById: (id: number) => ApiUser | undefined;
  getActiveUsers: () => ApiUser[];
  getUsersByRole: (role: string) => ApiUser[];
  getUsersByStatus: (status: string) => ApiUser[];
}

// 导出类型
export type { UserState };

// 创建 store
export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // 实现...
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({
        users: state.users.slice(0, 50),
        filters: state.filters,
      }),
    }
  )
);
```

---

## 参考资源

- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)
- [Zustand TypeScript 指南](https://github.com/pmndrs/zustand/blob/main/docs/guides/typescript.md)
- [项目 README.md](./README.md)
- [最佳实践指南](./BEST_PRACTICES.md)
