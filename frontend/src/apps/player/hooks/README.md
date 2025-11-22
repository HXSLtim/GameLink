# Player Hooks

## 📁 目录用途

该目录用于存放陪玩师端（Player）专用的自定义 React Hooks。

## 🎯 陪玩师端特点

陪玩师端是一个**移动端优先**的应用，核心功能包括：
- 实时接单
- 收益管理
- 个人形象展示
- 技能认证
- 订单管理

因此，这些 Hooks 需要特别考虑：**移动网络环境、实时性、操作便捷性、性能优化**。

## 📦 Hook 分类

### 实时接单 Hooks
- `useOrderNotifications` - 实时订单推送监听
- `useOrderRingtone` - 订单提醒音控制
- `useOrderVibration` - 订单震动提醒（移动端）
- `useOrderCountdown` - 接单倒计时（30秒内必须响应）
- `useAutoAccept` - 自动接单配置
- `useOrderPriority` - 订单优先级计算

### 收益管理 Hooks
- `useEarnings` - 收益明细查询
- `useWithdraw` - 提现操作
- `useBankCards` - 银行卡管理
- `useCommissionRate` - 佣金比例查询
- `useRevenueStats` - 收益统计分析
- `useWithdrawHistory` - 提现历史查询

### 个人中心 Hooks
- `usePlayerProfile` - 个人资料管理
- `useGameStats` - 游戏数据展示（胜率、段位等）
- `useSkillTags` - 技能标签管理
- `useServiceGallery` - 服务相册管理
- `useVideoIntro` - 视频介绍管理
- `useAvailability` - 接单状态设置
- `useWorkingHours` - 工作时间配置

### 订单管理 Hooks
- `usePlayerOrders` - 陪玩师订单列表
- `useActiveOrder` - 当前进行中的订单
- `useOrderStatusUpdate` - 订单状态更新
- `useOrderChat` - 订单内聊天
- `useOrderReview` - 订单评价查看
- `useOrderHistory` - 历史订单查询

### 技能认证 Hooks
- `useGameCertification` - 游戏认证状态
- `useUploadCredentials` - 上传认证材料
- `useCertificationProgress` - 认证进度查询
- `useSkillVerification` - 技能审核状态

### 评价管理 Hooks
- `useReceivedReviews` - 收到的评价列表
- `useReplyReview` - 回复评价
- `useReportReview` - 举报恶意评价
- `useReviewStats` - 评价统计分析（好评率）

### 系统设置 Hooks
- `useMobileNotifications` - 移动端推送通知设置
- `useSoundSettings` - 声音设置（提醒音、震动）
- `usePrivacySettings` - 隐私设置
- `useAccountSecurity` - 账户安全设置

### 移动端专用 Hooks
- `usePushNotifications` - 推送通知处理
- `useBackgroundMode` - 应用后台运行优化
- `useNetworkStatus` - 网络状态监听
- `useBatteryStatus` - 电池状态监听（低电量优化）
- `useAppBadge` - 应用角标数量（未读订单）

## 📋 命名规范

### Hook 命名规则
- 必须以 `use` 开头
- 使用 camelCase（小驼峰）
- 命名应具有描述性，突出**陪玩师场景**

```typescript
// ✅ 推荐（突出陪玩师场景）
export const usePlayerOrders = () => { ... }
export const useOrderAcceptCountdown = () => { ... }
export const useEarningsToday = () => { ... }

// ❌ 避免（过于通用）
export const useOrders = () => { ... }
export const useCountdown = () => { ... }
export const useRevenue = () => { ... }
```

### 返回值规范

#### 实时数据 Hook（WebSocket）
```typescript
interface UseRealTimeDataResult<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  connected: boolean;  // WebSocket 连接状态
  reconnect: () => void;
}

function useRealTimeOrder(userId: number): UseRealTimeDataResult<Order> {
  // ...
}
```

#### 移动端 Hook（带平台判断）
```typescript
interface UseMobileFeatureResult {
  supported: boolean;  // 是否支持该功能
  enabled: boolean;    // 是否启用
  setEnabled: (enabled: boolean) => void;
  permission?: PermissionStatus; // 权限状态
}

function usePushNotifications(): UseMobileFeatureResult {
  // ...
}
```

#### 收益相关 Hook
```typescript
interface UseEarningsResult {
  earnings: EarningsData;
  loading: boolean;
  error: Error | null;
  refresh: () => Promise<void>;
  canWithdraw: boolean;
}

function usePlayerEarnings(): UseEarningsResult {
  // ...
}
```

## 🎯 开发规范

### TypeScript 严格类型

陪玩师端业务复杂，必须有严格的类型定义：

```typescript
// ✅ 推荐
type OrderStatus = 'pending' | 'accepted' | 'in_progress' | 'completed' | 'cancelled';

type NotificationType = 'order' | 'review' | 'payment' | 'system';

interface PlayerOrder extends Order {
  earnings: number;           // 陪玩师收益
  commission: number;         // 平台佣金
  tip: number;                // 小费
  playerStatus: 'waiting' | 'ready' | 'serving';
}

export interface UsePlayerOrdersOptions {
  status?: OrderStatus[];
  dateRange?: [Date, Date];
  gameId?: number;
}

export interface UsePlayerOrdersResult {
  orders: PlayerOrder[];
  totalEarnings: number;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
  loadMore: () => Promise<void>;
}

export function usePlayerOrders(
  options: UsePlayerOrdersOptions = {}
): UsePlayerOrdersResult {
  const [orders, setOrders] = useState<PlayerOrder[]>([]);
  const [totalEarnings, setTotalEarnings] = useState(0);
  // ...
  return { orders, totalEarnings, /* ... */ };
}
```

### 错误处理（移动端友好）

移动端显示空间有限，错误提示要简洁：

```typescript
// ✅ 推荐 - 简洁的错误处理
export function usePlayerOrders() {
  const [error, setError] = useState<Error | null>(null);

  const fetchOrders = async () => {
    try {
      setError(null);
      const data = await api.player.getOrders();
      setOrders(data);
    } catch (err) {
      setError(err);
      // 移动端用简洁的提示
      showNotification('error', '订单获取失败', '请检查网络');
    }
  };

  return { error /* ... */ };
}

// ❌ 避免 - 复杂的长消息
showNotification('error', '订单获取失败', `
  原因: ${error.message}
  代码: ${error.code}
  建议: 请稍后重试或联系客服
`);
```

### 性能优化（移动端重要）

移动端性能更重要，必须优化：

```typescript
// ✅ 推荐 - 使用 useCallback 和 useMemo
export function usePlayerOrders(userId: number) {
  // 使用 useCallback 缓存函数
  const fetchOrders = useCallback(async () => {
    const data = await api.player.getOrders(userId);
    setOrders(data);
  }, [userId]); // 明确依赖

  // 使用 useMemo 缓存计算结果
  const totalEarnings = useMemo(() => {
    return orders.reduce((sum, order) => sum + order.earnings, 0);
  }, [orders]);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]); // 使用缓存的函数

  return { orders, totalEarnings, fetchOrders };
}

// ❌ 避免 - 每次渲染都创建新函数
export function usePlayerOrders(userId: number) {
  const fetchOrders = async () => {
    // 每次渲染都创建新函数，导致 useEffect 频繁触发
  };

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]); // 依赖变化频繁
}
```

### 实时数据优化（WebSocket）

陪玩师端大量使用实时数据，需要优化连接：

```typescript
// ✅ 推荐 - WebSocket 连接管理
export function useOrderNotifications(playerId: number) {
  const [connection, setConnection] = useState<WebSocket | null>(null);
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    // 创建连接
    const ws = new WebSocket(`wss://api.gamelink.com/players/${playerId}/orders`);
    setConnection(ws);

    ws.onmessage = (event) => {
      const newOrder = JSON.parse(event.data);
      setOrders(prev => [newOrder, ...prev]);

      // 移动端震动提醒
      if (navigator.vibrate) {
        navigator.vibrate([200, 100, 200]); // 震动 200ms, 停 100ms, 再震动 200ms
      }
    };

    // 清理函数：关闭连接
    return () => {
      ws.close();
    };
  }, [playerId]);

  // 断线重连逻辑
  const reconnect = useCallback(() => {
    if (connection?.readyState === WebSocket.CLOSED) {
      // 重连逻辑
      showNotification('warning', '连接已断开', '正在尝试重新连接');
      // ...
    }
  }, [connection]);

  useEffect(() => {
    const interval = setInterval(reconnect, 5000); // 每 5 秒检查一次
    return () => clearInterval(interval);
  }, [reconnect]);

  return { orders, connection, reconnect };
}
```

### 移动端专用 Hook

移动端有独特的 API 和限制：

```typescript
// ✅ 推荐 - 推送通知 Hook
export function usePushNotifications() {
  const [permission, setPermission] = useState<NotificationPermission | 'unsupported'>(
    'default'
  );
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    // 检查是否支持
    if (!('Notification' in window)) {
      setPermission('unsupported');
      return;
    }

    // 检查权限
    setPermission(Notification.permission);
  }, []);

  const requestPermission = useCallback(async () => {
    if (!('Notification' in window)) return false;

    const result = await Notification.requestPermission();
    setPermission(result);
    return result === 'granted';
  }, []);

  const showNotification = useCallback((title: string, body: string) => {
    if (!enabled || permission !== 'granted') return;

    new Notification(title, {
      body,
      icon: '/logo-192.png',
      badge: '/badge-72.png',
      vibrate: [200, 100, 200],
    });
  }, [enabled, permission]);

  const isSupported = permission !== 'unsupported';

  return {
    permission,
    enabled,
    setEnabled,
    isSupported,
    requestPermission,
    showNotification,
  };
}
```

## 📚 示例代码

### 实时订单 Hook（WebSocket + 本地通知）

```typescript
// useRealTimeOrders.ts
import { useState, useEffect, useCallback } from 'react';
import { api } from '@/api';
import type { Order } from '@/api';
import { useNotification } from '@/stores';

export interface UseRealTimeOrdersResult {
  orders: Order[];
  loading: boolean;
  error: Error | null;
  connection: WebSocket | null;
  connected: boolean;
  unreadCount: number;
  markAsRead: (orderId: number) => void;
}

export function useRealTimeOrders(playerId: number): UseRealTimeOrdersResult {
  const [orders, setOrders] = useState<Order[]>([]);
  const [connection, setConnection] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [unreadCount, setUnreadCount] = useState(0);
  const showNotification = useNotification();

  // 连接 WebSocket
  useEffect(() => {
    setLoading(true);

    const ws = new WebSocket(`wss://api.gamelink.com/ws/players/${playerId}/orders`);

    ws.onopen = () => {
      setConnection(ws);
      setConnected(true);
      setLoading(false);
      console.log('WebSocket 连接已建立');
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === 'NEW_ORDER') {
          const order: Order = data.order;

          // 添加到订单列表
          setOrders(prev => [order, ...prev]);

          // 增加未读数
          setUnreadCount(prev => prev + 1);

          // 显示通知
          showNotification('success', '新订单', `
            ${order.game?.name} - ¥${order.amount}
            ${order.user?.username}
          `);

          // 移动端本地通知
          if (navigator.vibrate) {
            navigator.vibrate([200, 100, 200]);
          }

          // 更新应用角标
          if ('setAppBadge' in navigator) {
            // @ts-ignore
            navigator.setAppBadge(unreadCount + 1);
          }
        } else if (data.type === 'ORDER_CANCELLED') {
          const orderId = data.orderId;

          // 从列表中移除
          setOrders(prev => prev.filter(o => o.id !== orderId));

          showNotification('info', '订单已取消', '用户取消了订单');
        }
      } catch (err) {
        console.error('解析消息失败:', err);
      }
    };

    ws.onerror = (event) => {
      console.error('WebSocket 错误:', event);
      setError(new Error('连接错误'));
      setConnected(false);
    };

    ws.onclose = () => {
      console.log('WebSocket 连接已关闭');
      setConnected(false);
      setConnection(null);
    };

    return () => {
      ws.close();
    };
  }, [playerId, showNotification, unreadCount]);

  const markAsRead = useCallback((orderId: number) => {
    setOrders(prev =>
      prev.map(order =>
        order.id === orderId ? { ...order, isRead: true } : order
      )
    );
    setUnreadCount(prev => Math.max(0, prev - 1));
  }, []);

  return {
    orders,
    loading,
    error,
    connection,
    connected,
    unreadCount,
    markAsRead,
  };
}
```

### 收益统计 Hook（移动端优化）

```typescript
// useEarningsStats.ts
import { useState, useEffect, useCallback, useMemo } from 'react';
import { api } from '@/api';
import type { EarningsData } from '@/api';

export interface UseEarningsStatsResult {
  today: number;
  thisWeek: number;
  thisMonth: number;
  availableForWithdraw: number;
  loading: boolean;
  error: Error | null;
  refresh: () => Promise<void>;
}

// 缓存 key
const EARNINGS_CACHE_KEY = 'player_earnings_cache';
const CACHE_DURATION = 5 * 60 * 1000; // 5 分钟

export function useEarningsStats(playerId: number): UseEarningsStatsResult {
  const [today, setToday] = useState(0);
  const [thisWeek, setThisWeek] = useState(0);
  const [thisMonth, setThisMonth] = useState(0);
  const [availableForWithdraw, setAvailableForWithdraw] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const loadFromCache = useCallback(() => {
    try {
      const cached = localStorage.getItem(EARNINGS_CACHE_KEY);
      if (cached) {
        const { data, timestamp } = JSON.parse(cached);
        if (Date.now() - timestamp < CACHE_DURATION) {
          setToday(data.today);
          setThisWeek(data.thisWeek);
          setThisMonth(data.thisMonth);
          setAvailableForWithdraw(data.availableForWithdraw);
          return true;
        }
      }
    } catch (err) {
      console.error('加载缓存失败:', err);
    }
    return false;
  }, []);

  const saveToCache = useCallback((data: EarningsData) => {
    try {
      const cache = {
        data,
        timestamp: Date.now(),
      };
      localStorage.setItem(EARNINGS_CACHE_KEY, JSON.stringify(cache));
    } catch (err) {
      console.error('保存缓存失败:', err);
    }
  }, []);

  const fetchEarnings = useCallback(async () => {
    setLoading(true);
    try {
      setError(null);

      // 优先从缓存加载
      const cached = loadFromCache();
      if (cached) {
        // 后台静默更新
        try {
          const data = await api.player.getEarnings(playerId);
          saveToCache(data);
          // 更新数据（不显示 loading）
          setToday(data.today);
          setThisWeek(data.thisWeek);
          setThisMonth(data.thisMonth);
          setAvailableForWithdraw(data.availableForWithdraw);
        } catch (error) {
          console.error('静默更新失败:', error);
        }
      } else {
        // 无缓存，从 API 获取
        const data = await api.player.getEarnings(playerId);
        saveToCache(data);
        setToday(data.today);
        setThisWeek(data.thisWeek);
        setThisMonth(data.thisMonth);
        setAvailableForWithdraw(data.availableForWithdraw);
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error('获取收益数据失败');
      setError(error);

      // 移动端：不显示复杂错误，简单提示
      showNotification('error', '数据加载失败', '请下拉刷新重试');
    } finally {
      setLoading(false);
    }
  }, [playerId, loadFromCache, saveToCache]);

  useEffect(() => {
    fetchEarnings();
  }, [fetchEarnings]);

  // 每分钟自动刷新一次
  useEffect(() => {
    const interval = setInterval(() => {
      // 只在非 loading 状态下刷新
      if (!loading) {
        fetchEarnings().catch(console.error);
      }
    }, 60 * 1000);

    return () => clearInterval(interval);
  }, [fetchEarnings, loading]);

  const refresh = useCallback(async () => {
    // 清除缓存，强制刷新
    localStorage.removeItem(EARNINGS_CACHE_KEY);
    return fetchEarnings();
  }, [fetchEarnings]);

  return {
    today,
    thisWeek,
    thisMonth,
    availableForWithdraw,
    loading,
    error,
    refresh,
  };
}
```

## 🎯 移动端优化最佳实践

### 1. 减少网络请求
- **缓存**: 使用 localStorage 缓存数据
- **防抖**: 搜索等输入延迟请求
- **分页**: 大量数据分页加载
- **预加载**: 提前加载可能需要的数据

### 2. 优化渲染性能
- **虚拟滚动**: 长列表使用虚拟滚动（react-window）
- **图片优化**: 使用 WebP 格式，懒加载
- **懒加载**: 路由懒加载、组件懒加载
- **代码分割**: 按需加载代码

### 3. 提升用户体验
- **离线支持**: 缓存关键数据，离线可查看
- **骨架屏**: 加载时显示骨架屏
- **进度指示**: 所有异步操作显示进度
- **快速反馈**: 点击后立即有反馈（即使还在加载）

### 4. 电池和网络优化
- **减少动画**: 复杂动画消耗电量
- **降低频率**: 后台数据更新降低频率
- **按需加载**: 只在需要时加载数据
- **压缩数据**: 使用 Gzip/Brotli 压缩

## 📝 注意事项

1. **移动端测试**: 必须在真实手机上测试，不只是浏览器模拟器
2. **网络兼容性**: 支持 4G/5G，兼容弱网环境
3. **浏览器兼容**: 兼容主流移动端浏览器（Safari, Chrome）
4. **PWA 支持**: 考虑添加 PWA 支持，提升体验
5. **错误处理**: 移动端错误处理要简洁，不要显示复杂堆栈
6. **日志收集**: 移动端错误日志要自动上报

---

**最后更新**: 2025-11-22
**维护者**: GameLink 前端团队
**特点**: 🎯 移动端优先 | ⚡ 实时更新 | 🎮 游戏陪玩场景
**相关文档**: [全局 Hooks 目录](../../shared/hooks/)
