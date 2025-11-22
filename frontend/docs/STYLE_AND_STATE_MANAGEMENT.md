# 样式系统和状态管理使用指南

## 📋 概述

本项目集成了现代化的样式系统和状态管理方案：
- **Tailwind CSS**：原子化样式系统，提供丰富的实用工具类
- **Zustand**：轻量级、高性能的状态管理库

## 🎨 Tailwind CSS 样式系统

### 配置

Tailwind CSS 配置文件：`tailwind.config.ts`

```typescript
// 主要配置
- 扫描路径: './src/**/*.{js,ts,jsx,tsx}'
- 主题颜色: primary (50-900)，gamelink (#1772f6)
- 字体: Inter 系统字体栈
```

### 使用方法

#### 基础使用

Tailwind CSS 提供原子化的实用类，直接在 JSX 中使用：

```tsx
// 按钮示例
<button className="btn btn-primary">
  主要按钮
</button>

// 卡片示例
<div className="card p-6">
  <h2 className="text-xl font-bold text-gray-800">标题</h2>
  <p className="text-gray-600 mt-2">内容描述</p>
</div>
```

#### 自定义组件

在全局样式文件 `src/styles/tailwind.css` 中定义了一些基础组件：

```css
/* 卡片组件 */
.card { @apply bg-white rounded-lg shadow-sm border border-gray-200; }

/* 按钮组件 */
.btn { @apply px-4 py-2 rounded-md font-medium transition-colors duration-200; }
.btn-primary { @apply btn bg-primary-600 text-white hover:bg-primary-700; }

/* 输入框组件 */
.input { @apply px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500; }
```

#### 响应式设计

Tailwind 提供响应式前缀：

```tsx
// 移动端单列，桌面端三列
<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
  <div>内容1</div>
  <div>内容2</div>
  <div>内容3</div>
</div>
```

响应式断点：
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px
- `2xl`: 1536px

#### 颜色系统

使用主题颜色：

```tsx
// 主色
<div className="bg-primary-500 text-white">主色调</div>

// GameLink 品牌色
<div className="bg-gamelink text-white">品牌色</div>

// 语义化颜色
<div className="text-green-600 bg-green-50">成功</div>
<div className="text-red-600 bg-red-50">错误</div>
<div className="text-yellow-600 bg-yellow-50">警告</div>
```

## 🗃️ Zustand 状态管理

### 核心概念

Zustand 是一个轻量级的状态管理库，特点：
- 零配置，简单易用
- 支持 TypeScript，类型安全
- 中间件支持（持久化、异步等）
- 性能优秀，按需重渲染
- 体积小巧（< 1KB）

### Store 结构

#### 1. 认证状态管理 (`useAuthStore.ts`)

```typescript
// 使用示例
import { useAuthStore } from '@/stores';

function LoginComponent() {
  const { token, user, isAuthenticated, setAuth, clearAuth } = useAuthStore();

  const handleLogin = async () => {
    // 登录成功后
    setAuth(token, user);
  };

  const handleLogout = () => {
    clearAuth();
  };

  return (
    <div>
      {isAuthenticated ? (
        <div>欢迎, {user?.username}</div>
      ) : (
        <button onClick={handleLogin}>登录</button>
      )}
    </div>
  );
}
```

**特性**：
- 使用 `persist` 中间件自动持久化到 localStorage
- 包含 token、user、isAuthenticated 状态
- 提供 setAuth、clearAuth、updateUser 操作

#### 2. 应用全局状态管理 (`useAppStore.ts`)

```typescript
// 使用示例
import { useAppStore, useNotification, useLoading } from '@/stores';
import type { NotificationType } from '@/stores';

function MyComponent() {
  const { isSidebarCollapsed, toggleSidebar } = useAppStore();

  return (
    <button onClick={() => toggleSidebar()}>
      {isSidebarCollapsed ? '展开' : '收起'}
    </button>
  );
}
```

##### 通知系统

```tsx
import { useNotification } from '@/stores';

function NotifyExample() {
  const showNotification = useNotification();

  const handleClick = () => {
    // 信息通知
    showNotification('info', '提示', '操作成功');

    // 成功通知
    showNotification('success', '成功', '数据已保存', 3000);

    // 警告通知
    showNotification('warning', '警告', '请检查输入');

    // 错误通知
    showNotification('error', '错误', '操作失败，请重试');
  };

  return <button onClick={handleClick}>显示通知</button>;
}
```

通知参数：
- `type`: 通知类型 (`info` | `success` | `warning` | `error`)
- `title`: 标题
- `message`: 内容
- `duration`: 持续时间（毫秒，默认 4000）

##### 加载状态

```tsx
import { useLoading } from '@/stores';

function LoadingExample() {
  const { showLoading, hideLoading, isLoading, message } = useLoading();

  const handleAsyncOperation = async () => {
    showLoading('正在加载数据...');

    try {
      await fetchData();
    } finally {
      hideLoading();
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="loading-spinner mr-3" />
        <span>{message || '加载中...'}</span>
      </div>
    );
  }

  return <button onClick={handleAsyncOperation}>执行操作</button>;
}
```

##### 数据缓存

缓存用户的筛选条件：

```tsx
import { useAppStore } from '@/stores';
import type { OrderStatus } from '@/api';

function OrderList() {
  const { cachedFilters, updateCachedFilters, clearCachedFilters } = useAppStore();

  const handleFilterChange = (status: OrderStatus) => {
    updateCachedFilters({ orderStatus: status });
  };

  return (
    <div>
      当前筛选: {cachedFilters.orderStatus}
      <button onClick={clearCachedFilters}>清除筛选</button>
    </div>
  );
}
```

### 最佳实践

#### 1. 按功能拆分 Store

不要将所有状态放在一个 Store 中，按功能模块拆分：

```typescript
// ✅ 推荐
- useAuthStore.ts      // 认证相关
- useAppStore.ts       // 应用通用
- useGameStore.ts      // 游戏相关（未来）
- useChatStore.ts      // 聊天相关（未来）

// ❌ 避免
- useGlobalStore.ts    // 所有状态混在一起
```

#### 2. 使用选择性订阅

只订阅需要的状态，避免不必要的重渲染：

```typescript
// ✅ 推荐：只订阅需要的状态
const { token } = useAuthStore();

// ❌ 避免：订阅整个状态
const store = useAuthStore(); // 任何状态变化都会触发重渲染
```

#### 3. 使用快捷 Hooks

对于常用操作，提供快捷 Hooks：

```typescript
// ✅ 推荐：使用快捷 hook
const showNotification = useNotification();
const { showLoading, hideLoading } = useLoading();

// ❌ 避免：每次都从 store 中解构
const showNotification = useAppStore((state) => state.showNotification);
```

#### 4. 合理使用持久化

只在需要时持久化：

```typescript
// 认证状态需要持久化
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({ /* ... */ }),
    { name: 'auth-storage' }
  )
);

// UI 状态不需要持久化
export const useAppStore = create<AppState>((set, get) => ({
  isSidebarCollapsed: false, // 刷新页面后重置
  // ...
}));
```

#### 5. 类型安全

始终定义接口/类型，享受 TypeScript 的类型检查：

```typescript
interface UserState {
  // 类型化状态
  user: User | null;
  permissions: Record<string, boolean>;

  // 类型化操作
  updateUser: (user: Partial<User>) => void;
  checkPermission: (key: string) => boolean;
}
```

## 🔧 开发工具

### Prettier 插件

Tailwind CSS 类名自动排序（已配置）：

```json
{
  "plugins": ["prettier-plugin-tailwindcss"],
  "tailwindConfig": "./tailwind.config.ts"
}
```

### VS Code 扩展

推荐安装：
- **Tailwind CSS IntelliSense**：智能提示、自动完成
- **Headwind**：自动排序 Tailwind 类名

### 开发服务器代理

在 `vite.config.ts` 中配置：

```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8000', // 后端 API 地址
      changeOrigin: true,
    },
  },
}
```

## 📝 总结

- **Tailwind CSS**：原子化样式，快速构建响应式界面
- **Zustand**：轻量状态管理，支持持久化和 TypeScript
- 两者结合：高效开发，维护简单，性能优秀

更多详情参考：
- [Tailwind CSS 官方文档](https://tailwindcss.com/docs)
- [Zustand 官方文档](https://zustand-demo.pmnd.rs/)
