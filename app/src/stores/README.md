# GameLink Mini-Program - Zustand 状态管理

## 概述

本项目使用 [Zustand](https://github.com/pmndrs/zustand) 进行全局状态管理，针对 Taro 小程序环境进行了优化。Zustand 是一个轻量级、高性能的状态管理库，非常适合小程序开发。

## 架构设计

### 文件结构

```
app/src/stores/
├── index.ts              # 统一导出入口
├── types/
│   └── index.ts         # 通用类型定义
├── modules/
│   ├── authStore.ts     # 认证和授权（手机号登录、微信登录）
│   ├── userStore.ts     # 用户数据管理
│   ├── playerStore.ts   # 陪玩师管理
│   ├── orderStore.ts    # 订单管理
│   └── chatStore.ts     # 聊天管理
└── README.md            # 本文档
```

### Stores 列表

| Store | 用途 | 主要功能 |
|-------|------|---------|
| **authStore** | 认证和授权 | 手机验证码登录、微信登录、Token 管理 |
| **userStore** | 用户管理 | 用户信息、余额、VIP 状态 |
| **playerStore** | 陪玩师管理 | 陪玩师列表、筛选、收藏 |
| **orderStore** | 订单管理 | 订单创建、状态跟踪、历史订单 |
| **chatStore** | 聊天管理 | 消息列表、WebSocket 连接 |

## 快速开始

### 安装

```bash
npm install zustand
```

### 基础用法

#### 1. 在页面/组件中使用 Store

```typescript
import { useAuthStore } from '@/stores/modules/authStore';
import Taro from '@tarojs/taro';

function LoginPage() {
  const { loginWithCode, loading } = useAuthStore();

  const handleLogin = async () => {
    try {
      await loginWithCode(phone, code);
      Taro.showToast({ title: '登录成功', icon: 'success' });
    } catch (err) {
      Taro.showToast({ title: '登录失败', icon: 'none' });
    }
  };

  return <View onClick={handleLogin}>登录</View>;
}
```

#### 2. 使用选择器优化性能

```typescript
// 推荐：只订阅需要的状态
const userInfo = useAuthStore(s => s.userInfo);
const loading = useAuthStore(s => s.loading);

// 推荐：使用 shallow 比较多个字段
import { shallow } from 'zustand/shallow';

const { userInfo, token } = useAuthStore(
  s => ({ userInfo: s.userInfo, token: s.token }),
  shallow
);
```

## 核心 API

### authStore - 认证和授权

```typescript
import { useAuthStore } from '@/stores/modules/authStore';

// 状态
const { token, userInfo, isAuthenticated, loading } = useAuthStore();

// 操作
const { sendCode, loginWithCode, wechatLogin, logout } = useAuthStore();

// 发送验证码
await sendCode('13800138000');

// 验证码登录
await loginWithCode('13800138000', '123456');

// 微信登录
await wechatLogin();

// 退出登录
logout();
```

### userStore - 用户管理

```typescript
import { useUserStore } from '@/stores/modules/userStore';

// 状态
const { userInfo, balance, vipStatus } = useUserStore();

// 操作
const { fetchUserInfo, updateProfile } = useUserStore();

// 获取用户信息
await fetchUserInfo();

// 更新用户资料
await updateProfile({ nickname: '新昵称', avatar: 'xxx' });
```

### playerStore - 陪玩师管理

```typescript
import { usePlayerStore } from '@/stores/modules/playerStore';

// 状态
const { players, loading, filters } = usePlayerStore();

// 操作
const { fetchPlayers, setFilters, toggleFavorite } = usePlayerStore();

// 获取陪玩师列表
await fetchPlayers({ game: '王者荣耀', rank: '王者' });

// 设置筛选条件
setFilters({ game: '王者荣耀', minPrice: 30 });

// 收藏/取消收藏
await toggleFavorite(playerId);
```

### orderStore - 订单管理

```typescript
import { useOrderStore } from '@/stores/modules/orderStore';

// 状态
const { orders, currentOrder, loading } = useOrderStore();

// 操作
const { createOrder, fetchOrders, cancelOrder } = useOrderStore();

// 创建订单
await createOrder({
  playerId: 123,
  game: '王者荣耀',
  duration: 2,
});

// 获取订单列表
await fetchOrders({ status: 'in_progress' });

// 取消订单
await cancelOrder(orderId, '不需要了');
```

### chatStore - 聊天管理

```typescript
import { useChatStore } from '@/stores/modules/chatStore';

// 状态
const { messages, currentRoom, unreadCount } = useChatStore();

// 操作
const { fetchMessages, sendMessage, connectWebSocket } = useChatStore();

// 获取消息列表
await fetchMessages(orderId);

// 发送消息
await sendMessage({ orderId, content: '你好' });

// 连接 WebSocket
connectWebSocket();
```

## Taro 适配

### Storage 适配

小程序环境使用 `Taro.getStorageSync` 等 API，需要自定义 storage adapter：

```typescript
const taroStorage = {
  getItem: (name: string): string | null => {
    try {
      const value = Taro.getStorageSync(name);
      return value || null;
    } catch (e) {
      console.error('Taro getStorageSync error:', e);
      return null;
    }
  },
  setItem: (name: string, value: string): void => {
    try {
      Taro.setStorageSync(name, value);
    } catch (e) {
      console.error('Taro setStorageSync error:', e);
    }
  },
  removeItem: (name: string): void => {
    try {
      Taro.removeStorageSync(name);
    } catch (e) {
      console.error('Taro removeStorageSync error:', e);
    }
  },
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // ...
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => taroStorage),
    }
  )
);
```

### Toast 提示

使用 `Taro.showToast` 显示操作反馈：

```typescript
try {
  await someAction();
  Taro.showToast({
    title: '操作成功',
    icon: 'success',
    duration: 2000,
  });
} catch (error) {
  Taro.showToast({
    title: error.message || '操作失败',
    icon: 'none',
    duration: 2000,
  });
}
```

## 性能优化

### 选择器模式

只订阅需要的状态切片，避免不必要的重渲染：

```typescript
// ✅ 好：只订阅需要的状态
const userInfo = useAuthStore(s => s.userInfo);

// ❌ 不好：订阅整个 store
const { userInfo, token, isAuthenticated, loading } = useAuthStore();
```

### Shallow 比较

使用 `shallow` 比较多字段：

```typescript
import { shallow } from 'zustand/shallow';

const { userInfo, token } = useAuthStore(
  s => ({ userInfo: s.userInfo, token: s.token }),
  shallow
);
```

## 小程序特有注意事项

### 1. 存储限制

小程序 storage 有大小限制（通常 10MB），需要控制持久化数据量：

```typescript
persist(
  (set, get) => ({
    players: [], // 可能有大量数据
  }),
  {
    name: 'player-storage',
    partialize: (state) => ({
      // 只缓存必要数据
      favorites: state.favorites,
      // 不缓存整个列表
    }),
  }
)
```

### 2. 网络请求

小程序需要配置合法域名，确保 API 请求正常：

```typescript
// app.config.ts
export default {
  // ...
  permission: {
    'scope.userLocation': {
      desc: '你的位置信息将用于小程序位置接口的效果展示',
    },
  },
};
```

### 3. 微信登录

完整的微信登录流程：

```typescript
wechatLogin: async () => {
  set({ loading: true });

  try {
    // 1. 获取微信登录 code
    const loginRes = await Taro.login();
    console.log('wx.login code:', loginRes.code);

    // 2. 发送 code 到后端换取 session
    const response = await api.post<LoginResponse>('/auth/wechat-login', {
      code: loginRes.code,
    });

    // 3. 保存 token 和用户信息
    set({
      token: response.data.token,
      userInfo: response.data.user,
      isAuthenticated: true,
      loading: false,
    });
  } catch (error) {
    set({ loading: false });
    throw error;
  }
}
```

## 与 Admin Store 的差异

| 特性 | Admin (Web) | App (Mini-Program) |
|------|-------------|-------------------|
| Storage API | localStorage/sessionStorage | Taro.getStorageSync |
| Toast | Ant Design message | Taro.showToast |
| 路由 | react-router | Taro.navigateTo |
| 认证方式 | 用户名密码 | 手机验证码/微信登录 |
| 权限管理 | 基于 permission | 基于角色和 VIP |

## 相关文档

- [Admin Store 文档](../../../admin/src/stores/README.md)
- [Zustand 官方文档](https://github.com/pmndrs/zustand)
- [Taro 官方文档](https://taro-docs.jd.com/)

## 最佳实践

1. **只订阅需要的状态**：使用选择器避免不必要的重渲染
2. **控制持久化数据量**：小程序存储空间有限
3. **使用 Taro API**：Toast、Loading 等使用 Taro 提供的 API
4. **处理网络异常**：小程序网络不稳定，做好错误处理
5. **微信登录集成**：遵循微信登录流程规范

## 贡献指南

添加新 store 时，请遵循以下规范：

1. 在 `modules/` 目录创建文件
2. 在 `types/index.ts` 添加类型定义
3. 在 `index.ts` 导出 store
4. 编写 JSDoc 注释
5. 使用 Taro Storage 适配器
6. 在 authStore 的 logout 中清理新 store
7. 更新本文档的 Stores 列表
