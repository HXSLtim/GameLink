# API 使用示例

本指南介绍如何使用项目中的 API 模块进行后端交互。

## 目录结构

```
src/api/
├── interface/          # 类型定义层（TypeScript 接口）
├── modules/            # 业务模块层（API 函数）
├── request/            # 请求工具层（HTTP 客户端）
├── client.ts           # Axios 实例配置
└── index.ts            # 统一导出
```

## 基本使用

### 1. 导入 API 模块

```typescript
// 导入整个 API 模块
import api, { authApi, orderApi, playerApi } from '@/api';

// 或按需导入特定模块
import { authApi } from '@/api';
import { orderApi } from '@/api';
```

### 2. 认证相关 API

#### 用户登录

```typescript
import { authApi } from '@/api';
import { useAuthStore } from '@/stores/useAuthStore';

const handleLogin = async () => {
  try {
    const { token, user } = await authApi.login({
      username: 'admin',
      password: 'password123',
    });
    
    // 保存 token 和用户信息
    useAuthStore.getState().setAuth(token, user);
    
    console.log('登录成功', user);
  } catch (error) {
    console.error('登录失败', error);
  }
};
```

#### 获取当前用户信息

```typescript
import { authApi } from '@/api';

const fetchCurrentUser = async () => {
  try {
    const user = await authApi.getCurrentUser();
    console.log('当前用户', user);
  } catch (error) {
    console.error('获取用户信息失败', error);
  }
};
```

#### 用户注册

```typescript
import { authApi } from '@/api';

const handleRegister = async () => {
  try {
    const result = await authApi.register({
      username: 'newuser',
      email: 'user@example.com',
      password: 'password123',
      confirmPassword: 'password123',
    });
    
    console.log('注册成功', result);
  } catch (error) {
    console.error('注册失败', error);
  }
};
```

### 3. 订单相关 API

#### 创建订单

```typescript
import { orderApi } from '@/api';

const createNewOrder = async () => {
  try {
    const order = await orderApi.createOrder({
      playerId: 1,
      gameId: 1,
      serviceType: '排位陪玩',
      duration: 2,
      amount: 10000, // 单位为分，即 100 元
      notes: '希望陪玩师能够带我上王者',
      requirements: '需要语音沟通',
    });
    
    console.log('订单创建成功', order);
  } catch (error) {
    console.error('创建订单失败', error);
  }
};
```

#### 获取订单列表

```typescript
import { orderApi } from '@/api';

const fetchOrders = async () => {
  try {
    const { list, total } = await orderApi.getOrders({
      status: 'pending',
      page: 1,
      pageSize: 10,
    });
    
    console.log('订单列表', list);
    console.log('总数', total);
  } catch (error) {
    console.error('获取订单列表失败', error);
  }
};
```

#### 获取订单详情

```typescript
import { orderApi } from '@/api';

const fetchOrderDetail = async (orderId: number) => {
  try {
    const order = await orderApi.getOrderById(orderId);
    console.log('订单详情', order);
  } catch (error) {
    console.error('获取订单详情失败', error);
  }
};
```

#### 取消订单

```typescript
import { orderApi } from '@/api';

const cancelOrder = async (orderId: number) => {
  try {
    await orderApi.cancelOrder(orderId, '临时有事，需要取消');
    console.log('订单已取消');
  } catch (error) {
    console.error('取消订单失败', error);
  }
};
```

### 4. 陪玩师相关 API

#### 获取陪玩师列表

```typescript
import { playerApi } from '@/api';

const fetchPlayers = async () => {
  try {
    const { list, total } = await playerApi.getPlayers({
      gameId: 1,
      status: 'available',
      minPrice: 50,
      maxPrice: 200,
      sortBy: 'rating',
      sortOrder: 'desc',
      page: 1,
      pageSize: 20,
    });
    
    console.log('陪玩师列表', list);
    console.log('总数', total);
  } catch (error) {
    console.error('获取陪玩师列表失败', error);
  }
};
```

#### 获取陪玩师详情

```typescript
import { playerApi } from '@/api';

const fetchPlayerDetail = async (playerId: number) => {
  try {
    const player = await playerApi.getPlayerById(playerId);
    console.log('陪玩师详情', player);
  } catch (error) {
    console.error('获取陪玩师详情失败', error);
  }
};
```

#### 获取推荐陪玩师

```typescript
import { playerApi } from '@/api';

const fetchRecommendedPlayers = async () => {
  try {
    const players = await playerApi.getRecommendedPlayers(10);
    console.log('推荐陪玩师', players);
  } catch (error) {
    console.error('获取推荐陪玩师失败', error);
  }
};
```

### 5. 游戏相关 API

#### 获取游戏列表

```typescript
import { gameApi } from '@/api';

const fetchGames = async () => {
  try {
    const { list } = await gameApi.getGames({
      category: 'moba',
      status: 'active',
      page: 1,
      pageSize: 50,
    });
    
    console.log('游戏列表', list);
  } catch (error) {
    console.error('获取游戏列表失败', error);
  }
};
```

#### 获取热门游戏

```typescript
import { gameApi } from '@/api';

const fetchPopularGames = async () => {
  try {
    const games = await gameApi.getPopularGames(10);
    console.log('热门游戏', games);
  } catch (error) {
    console.error('获取热门游戏失败', error);
  }
};
```

### 6. 支付相关 API

#### 创建支付

```typescript
import { paymentApi, PaymentMethod } from '@/api';

const createPayment = async (orderId: number) => {
  try {
    const payment = await paymentApi.createPayment({
      orderId,
      method: PaymentMethod.WECHAT,
    });
    
    // 跳转到支付页面或显示二维码
    if (payment.paymentUrl) {
      window.location.href = payment.paymentUrl;
    } else if (payment.qrCode) {
      console.log('支付二维码', payment.qrCode);
    }
  } catch (error) {
    console.error('创建支付失败', error);
  }
};
```

#### 查询支付状态

```typescript
import { paymentApi } from '@/api';

const checkPaymentStatus = async (paymentId: number) => {
  try {
    const status = await paymentApi.queryPaymentStatus(paymentId);
    console.log('支付状态', status);
  } catch (error) {
    console.error('查询支付状态失败', error);
  }
};
```

### 7. 评价相关 API

#### 创建评价

```typescript
import { reviewApi } from '@/api';

const createReview = async (orderId: number) => {
  try {
    const review = await reviewApi.createReview({
      orderId,
      revieweeId: 1,
      rating: 5,
      content: '非常棒的陪玩体验，技术很好，态度也很棒！',
      tags: ['技术好', '态度好', '沟通顺畅'],
    });
    
    console.log('评价创建成功', review);
  } catch (error) {
    console.error('创建评价失败', error);
  }
};
```

#### 获取评价列表

```typescript
import { reviewApi } from '@/api';

const fetchReviews = async () => {
  try {
    const { list, averageRating } = await reviewApi.getReviews({
      revieweeId: 1,
      page: 1,
      pageSize: 10,
    });
    
    console.log('评价列表', list);
    console.log('平均评分', averageRating);
  } catch (error) {
    console.error('获取评价列表失败', error);
  }
};
```

### 8. 通知相关 API

#### 获取通知列表

```typescript
import { notificationApi } from '@/api';

const fetchNotifications = async () => {
  try {
    const { list, unreadCount } = await notificationApi.getNotifications({
      status: 'unread',
      page: 1,
      pageSize: 20,
    });
    
    console.log('通知列表', list);
    console.log('未读数量', unreadCount);
  } catch (error) {
    console.error('获取通知列表失败', error);
  }
};
```

#### 标记通知为已读

```typescript
import { notificationApi } from '@/api';

const markAsRead = async (notificationId: number) => {
  try {
    await notificationApi.markNotificationAsRead(notificationId);
    console.log('通知已标记为已读');
  } catch (error) {
    console.error('标记通知失败', error);
  }
};
```

## 错误处理

所有 API 函数都使用统一的错误处理机制，可以通过 try-catch 捕获错误：

```typescript
import { authApi } from '@/api';

const handleLogin = async () => {
  try {
    const result = await authApi.login(credentials);
  } catch (error) {
    // error 是 Error 对象或 AxiosError 对象
    if (error.response) {
      // 服务器响应错误
      console.error('服务器错误:', error.response.data);
      console.error('状态码:', error.response.status);
    } else if (error.request) {
      // 请求发送但没有收到响应
      console.error('网络错误:', error.message);
    } else {
      // 其他错误
      console.error('错误:', error.message);
    }
  }
};
```

## 在 React 组件中使用

### 使用 useEffect 获取数据

```typescript
import { useState, useEffect } from 'react';
import { playerApi } from '@/api';
import type { Player } from '@/api';

const PlayerList = () => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayers = async () => {
      try {
        setLoading(true);
        const { list } = await playerApi.getPlayers({
          page: 1,
          pageSize: 20,
        });
        setPlayers(list);
      } catch (err) {
        setError(err instanceof Error ? err.message : '获取数据失败');
      } finally {
        setLoading(false);
      }
    };

    fetchPlayers();
  }, []);

  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error}</div>;

  return (
    <ul>
      {players.map((player) => (
        <li key={player.id}>{player.user.username}</li>
      ))}
    </ul>
  );
};
```

### 使用自定义 Hook 封装

```typescript
import { useState, useEffect } from 'react';
import { playerApi } from '@/api';
import type { Player } from '@/api';

const usePlayers = (params?: Parameters<typeof playerApi.getPlayers>[0]) => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    const fetchPlayers = async () => {
      try {
        setLoading(true);
        const response = await playerApi.getPlayers(params);
        setPlayers(response.players);
        setTotal(response.total);
      } catch (err) {
        setError(err instanceof Error ? err.message : '获取数据失败');
      } finally {
        setLoading(false);
      }
    };

    fetchPlayers();
  }, [params]);

  return { players, loading, error, total };
};

// 使用 Hook
const PlayerList = () => {
  const { players, loading, error, total } = usePlayers({
    page: 1,
    pageSize: 20,
  });

  // ... 渲染逻辑
};
```

## 类型导入

可以从 API 模块导入类型：

```typescript
import type { User, Order, Player, Payment, Review, Notification } from '@/api';

// 或使用具体的路径
import type { User } from '@/api/interface/auth';
import type { Order } from '@/api/interface/order';
import type { Player } from '@/api/interface/player';
```

## 最佳实践

1. **错误处理**：始终使用 try-catch 捕获 API 错误
2. **加载状态**：为异步操作提供加载状态
3. **类型安全**：使用 TypeScript 接口确保类型安全
4. **代码复用**：将 API 调用封装到自定义 Hook 中
5. **请求取消**：在组件卸载时取消未完成的请求（Axios 自动处理）
6. **缓存**：对于不经常变化的数据，考虑添加缓存机制

## 常见问题

### 如何处理 401 未授权错误？

项目已配置自动处理 401 错误：
- 如果是认证相关接口（登录、获取用户信息等），会清除 token 并重定向到登录页
- 如果是其他接口，会保留 token 并抛出错误，由业务代码处理

### 如何添加新的 API 模块？

1. 在 `src/api/interface/` 中创建接口定义文件
2. 在 `src/api/modules/` 中创建业务模块文件
3. 在 `src/api/index.ts` 中导出模块

### 如何修改请求超时时间？

修改 `src/api/client.ts` 中的 `timeout` 配置：

```typescript
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 20000, // 修改为 20 秒
  // ...
});
```