# API 模块文档

## 概述

本项目采用三层架构设计 API 模块，确保代码的可维护性和可扩展性。

### 三层架构

1. **Interface 层** (`src/api/interface/`)
   - 定义 TypeScript 接口和类型
   - 与后端 API 契约保持一致
   - 提供强类型支持

2. **Request 层** (`src/api/request/`)
   - HTTP 请求工具函数
   - 请求重试机制
   - 请求/响应拦截器

3. **Modules 层** (`src/api/modules/`)
   - 业务逻辑封装
   - 提供具体的 API 函数
   - 统一响应格式处理

## 目录结构

```
src/api/
├── interface/          # 类型定义层
│   ├── auth.ts        # 认证相关接口
│   ├── order.ts       # 订单相关接口
│   ├── game.ts        # 游戏相关接口
│   ├── payment.ts     # 支付相关接口
│   ├── review.ts      # 评价相关接口
│   └── notification.ts # 通知相关接口
│
├── modules/            # 业务模块层
│   ├── auth.ts        # 认证 API
│   ├── order.ts       # 订单 API
│   ├── player.ts      # 陪玩师 API
│   ├── game.ts        # 游戏 API
│   ├── payment.ts     # 支付 API
│   ├── review.ts      # 评价 API
│   └── notification.ts # 通知 API
│
├── request/            # 请求工具层
│   ├── http.ts        # HTTP 请求函数
│   └── retry.ts       # 重试机制
│
├── client.ts           # Axios 实例配置
├── index.ts            # 统一导出
└── API_USAGE_EXAMPLES.md # 使用示例
```

## 核心特性

### 1. 统一响应格式处理

所有 API 响应都经过统一处理，符合后端标准格式：

```typescript
// 成功响应
{
  success: true,
  data: T,
  message?: string
}

// 失败响应
{
  success: false,
  code: number,
  message: string,
  details?: any
}
```

### 2. 自动 Token 管理

请求拦截器自动添加 JWT Token：

```typescript
// 自动从 localStorage 获取 token
const token = storage.getItem<string>(STORAGE_KEYS.token);
if (token) {
  config.headers.Authorization = `Bearer ${token}`;
}
```

### 3. 智能错误处理

响应拦截器自动处理常见错误：

- **401 未授权**：清除 token 并重定向到登录页
- **403 禁止访问**：记录错误并抛出异常
- **404 资源不存在**：记录错误并抛出异常
- **500 服务器错误**：记录错误并抛出异常
- **网络错误**：自动重试机制

### 4. 请求重试机制

自动重试网络错误和 5xx 服务器错误：

```typescript
// 默认配置
maxRetries: 3
initialDelay: 1000ms
backoffFactor: 2 (指数退避)
maxDelay: 10000ms
```

### 5. 请求加密

支持请求/响应数据加密（通过 cryptoMiddleware）：

```typescript
// 请求加密
return cryptoMiddleware.requestInterceptor(config);

// 响应解密
const decryptedResponse = cryptoMiddleware.responseInterceptor(response);
```

## 使用方式

### 基础导入

```typescript
// 导入所有 API 模块
import api, { authApi, orderApi, playerApi } from '@/api';

// 导入类型
import type { User, Order, Player } from '@/api';
```

### 认证 API

```typescript
import { authApi } from '@/api';

// 登录
const { token, user } = await authApi.login({
  username: 'admin',
  password: 'password123'
});

// 获取当前用户
const user = await authApi.getCurrentUser();

// 登出
await authApi.logout();
```

### 订单 API

```typescript
import { orderApi } from '@/api';

// 创建订单
const order = await orderApi.createOrder({
  playerId: 1,
  gameId: 1,
  serviceType: '排位陪玩',
  duration: 2,
  amount: 10000
});

// 获取订单列表
const { list, total } = await orderApi.getOrders({
  status: 'pending',
  page: 1,
  pageSize: 10
});

// 获取订单详情
const order = await orderApi.getOrderById(1);

// 更新订单
const updatedOrder = await orderApi.updateOrder(1, {
  status: 'confirmed'
});

// 取消订单
await orderApi.cancelOrder(1, '临时有事');

// 完成订单
await orderApi.completeOrder(1);
```

### 陪玩师 API

```typescript
import { playerApi } from '@/api';

// 获取陪玩师列表
const { list, total } = await playerApi.getPlayers({
  gameId: 1,
  status: 'available',
  page: 1,
  pageSize: 20
});

// 获取陪玩师详情
const player = await playerApi.getPlayerById(1);

// 获取推荐陪玩师
const players = await playerApi.getRecommendedPlayers(10);
```

### 游戏 API

```typescript
import { gameApi } from '@/api';

// 获取游戏列表
const { list } = await gameApi.getGames({
  category: 'moba',
  page: 1,
  pageSize: 50
});

// 获取游戏详情
const game = await gameApi.getGameById(1);

// 获取热门游戏
const games = await gameApi.getPopularGames(10);

// 获取游戏服务类型
const services = await gameApi.getGameServices(1);
```

### 支付 API

```typescript
import { paymentApi, PaymentMethod } from '@/api';

// 创建支付
const payment = await paymentApi.createPayment({
  orderId: 1,
  method: PaymentMethod.WECHAT
});

// 获取支付详情
const payment = await paymentApi.getPaymentById(1);

// 查询支付状态
const status = await paymentApi.queryPaymentStatus(1);

// 申请退款
const refund = await paymentApi.requestRefund({
  paymentId: 1,
  amount: 5000,
  reason: '服务不满意'
});
```

### 评价 API

```typescript
import { reviewApi } from '@/api';

// 创建评价
const review = await reviewApi.createReview({
  orderId: 1,
  revieweeId: 1,
  rating: 5,
  content: '非常棒的体验！',
  tags: ['技术好', '态度好']
});

// 获取评价列表
const { list, averageRating } = await reviewApi.getReviews({
  revieweeId: 1,
  page: 1,
  pageSize: 10
});

// 获取评价统计
const stats = await reviewApi.getReviewStats({
  playerId: 1
});
```

### 通知 API

```typescript
import { notificationApi } from '@/api';

// 获取通知列表
const { list, unreadCount } = await notificationApi.getNotifications({
  status: 'unread',
  page: 1,
  pageSize: 20
});

// 标记通知为已读
await notificationApi.markNotificationAsRead(1);

// 获取未读通知数量
const { total } = await notificationApi.getUnreadNotificationCount();
```

## 错误处理

### API 错误类

```typescript
import { ApiError } from '@/shared/types/api';

try {
  const user = await authApi.getCurrentUser();
} catch (error) {
  if (error instanceof ApiError) {
    console.error('错误码:', error.code);
    console.error('错误信息:', error.message);
    console.error('用户友好信息:', error.getFriendlyMessage());
    
    if (error.isAuthError()) {
      // 认证错误，重定向到登录页
      window.location.href = '/login';
    } else if (error.isNetworkError()) {
      // 网络错误，显示网络错误提示
      alert('网络连接失败，请检查网络设置');
    } else if (error.isServerError()) {
      // 服务器错误，显示服务器错误提示
      alert('服务器暂时不可用，请稍后重试');
    }
  }
}
```

## 与 React Hooks 集成

### useApi Hook

```typescript
import { useApi } from '@/shared/hooks/useApi';
import { authApi } from '@/api';

const LoginForm = () => {
  const { data, loading, error, execute } = useApi(authApi.login);
  
  const handleLogin = async () => {
    const result = await execute({
      username: 'admin',
      password: 'password123'
    });
    
    if (result) {
      console.log('登录成功', result);
    }
  };
  
  return (
    <form onSubmit={handleLogin}>
      {/* 表单字段 */}
      <button type="submit" disabled={loading}>
        {loading ? '登录中...' : '登录'}
      </button>
      {error && <div>错误: {error.message}</div>}
    </form>
  );
};
```

### useApiQuery Hook

```typescript
import { useApiQuery } from '@/shared/hooks/useApi';
import { playerApi } from '@/api';

const PlayerDetail = ({ playerId }: { playerId: number }) => {
  const { data: player, loading, error, refresh } = useApiQuery(
    playerApi.getPlayerById,
    playerId
  );
  
  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error.message}</div>;
  if (!player) return <div>陪玩师不存在</div>;
  
  return (
    <div>
      <h1>{player.user.username}</h1>
      <button onClick={refresh}>刷新</button>
    </div>
  );
};
```

## 配置说明

### 环境变量

```env
# API 基础 URL
VITE_API_BASE_URL=http://localhost:8080

# 开发环境 Mock 配置
VITE_DEV_MOCK_USERNAME=admin
VITE_DEV_MOCK_PASSWORD=admin123
VITE_DEV_MOCK_TOKEN=dev-token

# 功能开关
VITE_ENABLE_MOCK=false
VITE_ENABLE_ANALYZER=false
```

### Axios 配置

```typescript
// src/api/client.ts
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000, // 10秒超时
  headers: {
    'Content-Type': 'application/json',
  },
});
```

### 重试配置

```typescript
// src/api/request/retry.ts
const DEFAULT_RETRY_OPTIONS: Required<RetryOptions> = {
  maxRetries: 3,
  initialDelay: 1000,
  backoffFactor: 2,
  maxDelay: 10000,
  shouldRetry: (error: unknown, attempt: number) => {
    // 网络错误
    if (error instanceof Error && error.message.includes('Failed to fetch')) {
      return true;
    }
    // HTTP 5xx 错误
    if (typeof error === 'object' && error !== null && 'status' in error) {
      const status = (error as { status: number }).status;
      return status >= 500 && status < 600;
    }
    return false;
  },
};
```

## 最佳实践

### 1. 错误处理

```typescript
// ✅ 推荐：使用 try-catch 捕获错误
try {
  const order = await orderApi.createOrder(data);
} catch (error) {
  if (error instanceof ApiError) {
    // 处理 API 错误
    console.error(error.getFriendlyMessage());
  }
}

// ✅ 推荐：使用 useApi Hook
const { data, loading, error, execute } = useApi(orderApi.createOrder);

// ❌ 避免：不使用错误处理
const order = await orderApi.createOrder(data); // 可能抛出未捕获的错误
```

### 2. 加载状态

```typescript
// ✅ 推荐：显示加载状态
const { loading, execute } = useApi(authApi.login);

<button disabled={loading}>
  {loading ? '登录中...' : '登录'}
</button>

// ❌ 避免：不提供加载反馈
<button onClick={handleLogin}>登录</button> // 用户不知道正在处理
```

### 3. 类型安全

```typescript
// ✅ 推荐：使用 TypeScript 接口
import type { CreateOrderRequest } from '@/api';

const data: CreateOrderRequest = {
  playerId: 1,
  gameId: 1,
  serviceType: '排位陪玩',
  duration: 2,
  amount: 10000
};

// ❌ 避免：使用 any 类型
const data: any = { ... }; // 失去类型检查
```

### 4. 代码复用

```typescript
// ✅ 推荐：封装自定义 Hook
const usePlayers = (params) => {
  const { data, loading, error } = useApiQuery(playerApi.getPlayers, params);
  return { players: data?.list || [], loading, error };
};

// ❌ 避免：重复编写相同逻辑
// 在每个组件中重复调用 playerApi.getPlayers
```

### 5. 请求取消

```typescript
// ✅ 推荐：使用 AbortController
const controller = new AbortController();
const response = await fetch(url, { signal: controller.signal });

// 组件卸载时取消请求
useEffect(() => {
  return () => controller.abort();
}, []);

// ❌ 避免：不取消未完成的请求
// 可能导致内存泄漏和状态更新错误
```

## 常见问题

### Q: 如何添加新的 API 模块？

A:
1. 在 `src/api/interface/` 创建接口定义文件
2. 在 `src/api/modules/` 创建业务模块文件
3. 在 `src/api/index.ts` 中导出模块

### Q: 如何修改请求超时时间？

A: 修改 `src/api/client.ts` 中的 `timeout` 配置：

```typescript
export const apiClient = axios.create({
  timeout: 20000, // 20秒
  // ...
});
```

### Q: 如何禁用请求重试？

A: 在调用时传入 `maxRetries: 0`：

```typescript
const response = await retryAsync(fetch, {
  maxRetries: 0
});
```

### Q: 如何调试 API 请求？

A: 使用浏览器开发者工具的 Network 面板，或添加日志：

```typescript
apiClient.interceptors.request.use((config) => {
  console.log('Request:', config.method, config.url, config.data);
  return config;
});

apiClient.interceptors.response.use((response) => {
  console.log('Response:', response.config.url, response.data);
  return response;
});
```

## 相关文档

- [API 使用示例](./API_USAGE_EXAMPLES.md)
- [后端 API 文档](https://github.com/HXSLtim/GameLink/blob/main/docs/API.md)