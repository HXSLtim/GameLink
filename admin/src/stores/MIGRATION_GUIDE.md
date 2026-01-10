# 从 Context API 迁移到 Zustand

本文档说明如何从旧的 Context API 状态管理迁移到新的 Zustand stores。

## 迁移优势

| 特性 | Context API | Zustand |
|------|-------------|---------|
| 重复请求 | 每页独立请求 | 全局缓存，只请求一次 |
| 组件通信 | Props drilling | 直接使用 store |
| 代码复用 | 难以复用 | 自定义 hooks |
| 性能 | 容易重渲染 | 选择器优化 |
| 测试 | 难以测试 | 纯函数，易测试 |
| 样板代码 | 需要创建 Context | 零配置 |

## 迁移步骤

### 1. 移除 Context Provider

**之前** (Context API):
```tsx
import { AdminProvider } from '@/contexts/AdminContext';

function App() {
  return (
    <AdminProvider>
      <MainApp />
    </AdminProvider>
  );
}
```

**之后** (Zustand):
```tsx
// 不需要 Provider，直接使用
function App() {
  return <MainApp />;
}
```

### 2. 替换 Hook 使用

**之前**:
```tsx
import { useAdmin } from '@/contexts/AdminContext';

const Component = () => {
  const { menus, hasPermission, userInfo } = useAdmin();
  // ...
};
```

**之后**:
```tsx
import { useMenuStore, useAuthStore, useUserStore } from '@/stores';

const Component = () => {
  const { menus } = useMenuStore();
  const { hasPermission } = useAuthStore();
  const { userInfo } = useUserStore();
  // ...
};
```

### 3. 数据获取模式

**之前**: 每个 page 的 useEffect 中独立调用 API
```tsx
const UserProfile = () => {
  const { userInfo, fetchUserInfo } = useAdmin();

  useEffect(() => {
    fetchUserInfo(); // 每个组件都调用一次
  }, []);

  return <div>{userInfo?.name}</div>;
};
```

**之后**: 使用 stores 的缓存数据，只在需要时刷新
```tsx
const UserProfile = () => {
  const { userInfo } = useUserStore(); // 自动从缓存读取

  return <div>{userInfo?.name}</div>;
};

// 如果需要刷新，显式调用
const OtherComponent = () => {
  const { fetchUserInfo } = useUserStore();

  const handleRefresh = () => {
    fetchUserInfo(); // 手动刷新
  };

  return <Button onClick={handleRefresh}>刷新</Button>;
};
```

## Store 使用示例

### Auth Store

```tsx
import { useAuthStore } from '@/stores/modules/authStore';

// 登录
const { login, isLoggedIn } = useAuthStore();
await login({ phone: '13800138000', password: 'password' });

// 权限检查
const { hasPermission } = useAuthStore();
if (hasPermission('admin.users.delete')) {
  // 有权限
}

// 登出
const { logout } = useAuthStore();
logout();
```

### User Store

```tsx
import { useUserStore } from '@/stores/modules/userStore';

// 获取用户列表
const { users, fetchUsers } = useUserStore();
await fetchUsers({ page: 1, pageSize: 10 });

// 创建用户
const { createUser } = useUserStore();
await createUser({ name: '张三', phone: '13800138000' });

// 更新用户
const { updateUser } = useUserStore();
await updateUser(1, { name: '李四' });

// 删除用户
const { deleteUser } = useUserStore();
await deleteUser(1);
```

### Menu Store

```tsx
import { useMenuStore } from '@/stores/modules/menuStore';

// 获取菜单树
const { menuTree, fetchMenuTree } = useMenuStore();
await fetchMenuTree();

// 检查菜单权限
const { hasMenuPermission } = useMenuStore();
if (hasMenuPermission('/admin/users')) {
  // 有菜单权限
}
```

### Player Store

```tsx
import { usePlayerStore } from '@/stores/modules/playerStore';

// 获取陪玩师列表
const { players, fetchPlayers } = usePlayerStore();
await fetchPlayers({ page: 1, pageSize: 10 });

// 审核陪玩师
const { approvePlayer, rejectPlayer } = usePlayerStore();
await approvePlayer(1);
await rejectPlayer(1, '资料不全');
```

## 完整示例

### 示例 1: Admin 用户列表页

参见: `admin/src/pages/admin/Users/stores-example.tsx`

```tsx
import { useEffect } from 'react';
import { Table, message } from 'antd';
import { useUserStore } from '@/stores/modules/userStore';
import { useAuthStore } from '@/stores/modules/authStore';

const UsersList = () => {
  const { users, loading, pagination, fetchUsers, deleteUser } = useUserStore();
  const { hasPermission } = useAuthStore();

  useEffect(() => {
    fetchUsers();
  }, []);

  if (!hasPermission('admin.users.read')) {
    return <div>无权限访问</div>;
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteUser(id);
      message.success('删除成功');
      await fetchUsers();
    } catch (error) {
      message.error('删除失败');
    }
  };

  return (
    <Table
      dataSource={users}
      loading={loading}
      rowKey="id"
      pagination={pagination}
    />
  );
};
```

### 示例 2: App 用户个人中心

参见: `app/src/pages/profile/stores-example.tsx`

```tsx
import { useEffect } from 'react';
import { View, Text, Button } from '@tarojs/components';
import Taro from '@tarojs/taro';
import { useUserStore } from '@/stores/modules/userStore';
import { useAuthStore } from '@/stores/modules/authStore';

const Profile = () => {
  const { userInfo, loading, updateProfile } = useUserStore();
  const { isLoggedIn } = useAuthStore();

  useEffect(() => {
    if (isLoggedIn()) {
      // 如果缓存中没有数据，会自动请求
      if (!userInfo) {
        // fetchUserInfo(); // 可选：显式刷新
      }
    }
  }, []);

  const handleUpdate = async () => {
    try {
      await updateProfile({ name: '新昵称' });
      Taro.showToast({ title: '更新成功', icon: 'success' });
    } catch (error) {
      Taro.showToast({ title: '更新失败', icon: 'error' });
    }
  };

  return (
    <View>
      <Text>昵称: {userInfo?.name}</Text>
      <Button onClick={handleUpdate}>更新</Button>
    </View>
  );
};
```

## 选择器优化

Zustand 支持选择器，只在特定状态变化时重渲染：

```tsx
// ✅ 好：只订阅 users，loading 变化不会重渲染
const users = useUserStore(state => state.users);

// ❌ 差：订阅整个 store，任何状态变化都重渲染
const { users, loading, error } = useUserStore();

// ✅ 好：使用浅比较选择器
import { shallow } from 'zustand/shallow';
const { users, loading } = useUserStore(
  state => ({ users: state.users, loading: state.loading }),
  shallow
);
```

## 开发工具

Zustand 提供了 Redux DevTools 集成：

```tsx
// store 中已配置
devTools: process.env.NODE_ENV === 'development',

// 在浏览器 Redux DevTools 中查看状态变化
```

## 常见问题

### Q: 如何处理组件间状态同步？

A: Zustand 是全局状态，自动同步：
```tsx
// 组件 A
const ComponentA = () => {
  const { setUser } = useUserStore();
  return <button onClick={() => setUser({ name: 'Alice' })}>设置</button>;
};

// 组件 B (自动获取更新)
const ComponentB = () => {
  const { user } = useUserStore();
  return <div>{user?.name}</div>;
};
```

### Q: 如何重置状态？

A: 使用 store 的 reset 方法：
```tsx
const { reset } = useUserStore();
reset(); // 重置为初始状态
```

### Q: 如何处理异步错误？

A: store 自动设置 error 状态：
```tsx
const { error, fetchUsers } = useUserStore();

try {
  await fetchUsers();
} catch (err) {
  // error 状态已自动设置
}

if (error) {
  message.error(error);
}
```

## 迁移检查清单

- [ ] 移除所有 Context Provider
- [ ] 替换 useAdmin() 为具体 store hooks
- [ ] 移除重复的 API 调用
- [ ] 使用选择器优化性能
- [ ] 测试权限检查
- [ ] 测试数据缓存
- [ ] 验证错误处理
- [ ] 删除旧的 Context 文件

## 相关文档

- [Zustand 官方文档](https://github.com/pmndrs/zustand)
- [Store 架构设计](../stores/modules/README.md)
- [API 文档](../../api/docs/)
