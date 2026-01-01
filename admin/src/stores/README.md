# GameLink Admin - Zustand 状态管理

## 概述

本项目使用 [Zustand](https://github.com/pmndrs/zustand) 进行全局状态管理，替代原有的 Context API 方案。Zustand 是一个轻量级、高性能的状态管理库，具有以下优势：

- 轻量简洁 (~1KB)
- 无需 Provider 包裹
- 完善的 TypeScript 支持
- 内置 DevTools 集成
- 支持状态持久化 (persist middleware)

## 架构设计

### 文件结构

```
admin/src/stores/
├── index.ts              # 统一导出入口
├── types/
│   └── index.ts         # 通用类型定义
├── modules/
│   ├── authStore.ts     # 认证和权限管理
│   ├── userStore.ts     # 用户数据管理
│   ├── menuStore.ts     # 菜单管理
│   ├── orderStore.ts    # 订单数据管理
│   ├── playerStore.ts   # 陪玩师管理
│   └── chatStore.ts     # 聊天管理
├── README.md            # 本文档
├── BEST_PRACTICES.md    # 最佳实践指南
└── TYPE_GUIDE.md        # TypeScript 类型指南
```

### Stores 列表

| Store | 用途 | 主要功能 |
|-------|------|---------|
| **authStore** | 认证和权限 | 登录/登出、Token 管理、权限检查、角色判断 |
| **userStore** | 用户管理 | 用户列表 CRUD、分页、筛选、批量操作 |
| **menuStore** | 菜单管理 | 菜单树缓存、权限过滤、面包屑生成 |
| **orderStore** | 订单管理 | 订单列表 CRUD、状态更新、批量操作、退款 |
| **playerStore** | 陪玩师管理 | 陪玩师列表、状态管理、等级管理 |
| **chatStore** | 聊天管理 | 消息列表、WebSocket 连接、未读数统计 |

## 快速开始

### 安装

```bash
npm install zustand
```

### 基础用法

#### 1. 在组件中使用 Store

```typescript
import { useAuthStore } from '@/stores/modules/authStore';

function LoginPage() {
  // 解构需要的状态和方法
  const { login, loading, error } = useAuthStore();

  const handleLogin = async () => {
    try {
      await login({ username, password });
      message.success('登录成功');
    } catch (err) {
      message.error(error || '登录失败');
    }
  };

  return <Button onClick={handleLogin} loading={loading}>登录</Button>;
}
```

#### 2. 使用选择器优化性能

```typescript
// 推荐：只订阅需要的状态
const users = useUserStore(s => s.users);
const loading = useUserStore(s => s.loading);

// 推荐：使用 shallow 比较多个字段
import { shallow } from 'zustand/shallow';

const { users, loading } = useUserStore(
  s => ({ users: s.users, loading: s.loading }),
  shallow
);

// 不推荐：订阅整个 store（任何变化都会重渲染）
const { users, loading, error, filters, pagination } = useUserStore();
```

#### 3. 使用便捷 hooks

每个 store 都导出了便捷的 hooks：

```typescript
import {
  useUsers,
  useUserLoading,
  useUserError,
  useUserPagination,
} from '@/stores/modules/userStore';

function UserList() {
  const users = useUsers();
  const loading = useUserLoading();
  const pagination = useUserPagination();

  // ...
}
```

## 核心 API

### authStore - 认证和权限

```typescript
import { useAuthStore } from '@/stores/modules/authStore';

// 状态
const { token, userInfo, isAuthenticated, loading, error } = useAuthStore();

// 操作
const { login, logout, setToken, clearError } = useAuthStore();

// 选择器 (权限检查)
const isAdmin = useAuthStore(s => s.isAdmin());
const hasPermission = useAuthStore(s => s.hasPermission('user:delete'));
const hasAllPermissions = useAuthStore(s => s.hasAllPermissions(['user:read', 'user:write']));
const hasAnyPermission = useAuthStore(s => s.hasAnyPermission(['user:delete', 'user:update']));
const hasRole = useAuthStore(s => s.hasRole('admin'));
```

**权限检查示例**：

```typescript
// 检查单个权限
if (hasPermission('user:delete')) {
  // 显示删除按钮
}

// 检查多个权限 (满足任意一个)
if (hasAnyPermission(['user:delete', 'user:update'])) {
  // 显示操作按钮
}

// 检查多个权限 (满足全部)
if (hasAllPermissions(['user:read', 'order:read'])) {
  // 显示管理面板
}

// 检查角色
if (hasRole(['admin', 'superAdmin'])) {
  // 显示管理员功能
}
```

### userStore - 用户管理

```typescript
import { useUserStore } from '@/stores/modules/userStore';

// 状态
const { users, loading, error, pagination, filters } = useUserStore();

// 操作
const { fetchUsers, createUser, updateUser, deleteUser } = useUserStore();
const { batchDeleteUsers, updateUserStatus, batchUpdateUserStatus } = useUserStore();
const { updateUserRole, batchUpdateUserRole, setFilters, clearFilters } = useUserStore();

// 选择器
const getUserById = useUserStore(s => s.getUserById(userId));
const activeUsers = useUserStore(s => s.getActiveUsers());
const admins = useUserStore(s => s.getUsersByRole('admin'));
```

**分页和筛选示例**：

```typescript
// 设置筛选条件
useUserStore.getState().setFilters({
  status: 'active',
  role: 'player',
  keyword: '张三'
});

// 获取用户列表（会应用筛选条件）
useUserStore.getState().fetchUsers(1, 20);
```

### orderStore - 订单管理

```typescript
import { useOrderStore } from '@/stores/modules/orderStore';

// 状态
const { orders, selectedOrderIds, loading, error, pagination, filters } = useOrderStore();

// 操作
const { fetchOrders, cancelOrder, refundOrder } = useOrderStore();
const { batchCancelOrders, batchCompleteOrders } = useOrderStore();
const { setSelectedOrders, setFilters } = useOrderStore();

// 选择器
const pendingOrders = useOrderStore(s => s.getOrdersByStatus('pending'));
const userOrders = useOrderStore(s => s.getOrdersByUser(userId));
const pendingCount = useOrderStore(s => s.getPendingOrdersCount());
```

### menuStore - 菜单管理

```typescript
import { useMenuStore } from '@/stores/modules/menuStore';

// 状态
const { menus, loading, flatMenus, breadcrumb } = useMenuStore();

// 操作
const { fetchMenus, generateBreadcrumb } = useMenuStore();
```

### playerStore - 陪玩师管理

```typescript
import { usePlayerStore } from '@/stores/modules/playerStore';

// 状态
const { players, loading, pagination, filters } = usePlayerStore();

// 操作
const { fetchPlayers, updatePlayerStatus, setFilters } = usePlayerStore();
```

### chatStore - 聊天管理

```typescript
import { useChatStore } from '@/stores/modules/chatStore';

// 状态
const { rooms, messages, currentRoomId, unreadCount } = useChatStore();

// 操作
const { fetchRooms, fetchMessages, sendMessage, connectWebSocket } = useChatStore();
```

## 状态持久化

所有 stores 使用 `persist` middleware 进行状态持久化：

```typescript
import { persist } from 'zustand/middleware';

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // store implementation
    }),
    {
      name: 'user-storage',  // localStorage key
      partialize: (state) => ({
        // 只持久化部分状态
        users: state.users.slice(0, 50),  // 限制缓存大小
        filters: state.filters,
      }),
    }
  )
);
```

**注意事项**：
- `authStore` 使用 `sessionStorage`（关闭标签页即清除，安全最佳实践）
- 其他 stores 使用 `localStorage`（保留筛选条件和缓存数据）
- 登出时会自动清理所有 stores 的持久化数据

## 开发工具

### 启用 DevTools

在开发环境中，Zustand 会自动集成 Redux DevTools：

```typescript
// 开发环境自动启用
if (import.meta.env.DEV) {
  // DevTools 已在 store 配置中启用
}
```

在浏览器 DevTools 中查看状态变化：
1. 打开 Redux DevTools 扩展
2. 选择 "Zustand" 状态
3. 查看状态变化历史

## 常见问题

### 1. 如何在组件外使用 store？

```typescript
// 使用 getState() 获取当前状态
const currentUser = useAuthStore.getState().userInfo;

// 使用 setState() 更新状态
useUserStore.getState().setFilters({ status: 'active' });

// 调用 actions
useAuthStore.getState().logout();
```

### 2. 如何监听状态变化？

```typescript
// 使用 subscribe() 监听状态变化
const unsubscribe = useUserStore.subscribe(
  (state) => state.users,
  (users) => {
    console.log('Users updated:', users);
  }
);

// 取消订阅
unsubscribe();
```

### 3. 如何重置 store？

```typescript
// 每个都有 reset() 方法
useUserStore.getState().reset();

// authStore 的 logout 会重置所有 stores
await useAuthStore.getState().logout();
```

### 4. 如何处理异步错误？

```typescript
const { createUser, error } = useUserStore();

try {
  await createUser({ name: '张三', email: 'test@example.com' });
  message.success('创建成功');
} catch (err) {
  // 错误已在 store 中设置，可直接使用
  message.error(error || '创建失败');
}
```

## 相关文档

- [BEST_PRACTICES.md](./BEST_PRACTICES.md) - 最佳实践指南
- [TYPE_GUIDE.md](./TYPE_GUIDE.md) - TypeScript 类型指南
- [Zustand 官方文档](https://github.com/pmndrs/zustand)

## 迁移指南

### 从 Context API 迁移

**旧代码 (Context API)**：
```typescript
const { users, loading } = useUserContext();
```

**新代码 (Zustand)**：
```typescript
import { useUserStore } from '@/stores/modules/userStore';

const users = useUserStore(s => s.users);
const loading = useUserStore(s => s.loading);
```

**优势**：
- 无需 Provider 包裹
- 减少不必要的重渲染
- 更简洁的 API
- 内置持久化支持

## 贡献指南

添加新 store 时，请遵循以下规范：

1. 在 `modules/` 目录创建文件
2. 在 `types/index.ts` 添加类型定义
3. 在 `index.ts` 导出 store
4. 编写 JSDoc 注释
5. 导出便捷 hooks
6. 在 authStore 的 logout 中清理新 store
7. 更新本文档的 Stores 列表
