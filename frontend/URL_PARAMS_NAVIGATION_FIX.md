# URL参数导航修复

## 问题描述

仪表盘中的快捷入口和订单状态卡片点击后会导航到对应页面（如 `/orders?status=pending`），但目标页面没有读取和应用这些URL参数，导致：
1. 筛选状态没有自动更新
2. 不能正确获取对应状态的数据

## 解决方案

### 1. 创建URL参数工具函数

新建 `src/utils/urlParams.ts`，提供通用的URL参数处理功能：

```typescript
export function mergeUrlParams<T extends Record<string, any>>(
  searchParams: URLSearchParams,
  initialParams: T,
  paramKeys: string[],
): T
```

该函数能够：
- 从URL查询参数中提取指定的参数
- 自动处理数字类型转换
- 与初始参数合并

### 2. 增强 `useListPage` Hook

更新 `src/hooks/useListPage.ts`，新增以下功能：

**新增配置项：**
- `urlParamKeys?: string[]` - 指定需要从URL读取的参数键

**实现细节：**
1. 组件初始化时从URL参数中读取并合并到初始参数
2. 监听URL参数变化，自动更新查询参数
3. URL参数变化时自动重置到第一页

**使用示例：**

```typescript
const { ... } = useListPage<User, UserListQuery>({
  initialParams: { ... },
  fetchData: async (params) => { ... },
  urlParamKeys: ['role', 'status'], // 从URL读取这些参数
});
```

### 3. 修复订单列表页面

更新 `src/pages/Orders/OrderList.tsx`：

1. 导入 `useSearchParams` hook
2. 创建 `getInitialParams` 函数从URL读取status参数
3. 添加 `useEffect` 监听URL参数变化
4. 在数据加载依赖中添加 `queryParams.status`

**修改代码：**

```typescript
const [searchParams] = useSearchParams();

// 从URL读取初始状态
const getInitialParams = (): OrderListQuery => {
  const statusFromUrl = searchParams.get('status') as OrderStatus | null;
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: statusFromUrl || undefined,
  };
};

const [queryParams, setQueryParams] = useState<OrderListQuery>(getInitialParams());

// 监听URL参数变化
useEffect(() => {
  const statusFromUrl = searchParams.get('status') as OrderStatus | null;
  setQueryParams((prev) => ({
    ...prev,
    status: statusFromUrl || undefined,
    page: 1,
  }));
}, [searchParams]);
```

### 4. 更新用户列表页面

更新 `src/pages/Users/UserList.tsx`，使用新增的 `urlParamKeys` 配置：

```typescript
const { ... } = useListPage<User, UserListQuery>({
  initialParams: { ... },
  fetchData: async (params) => { ... },
  urlParamKeys: ['role', 'status'],
});
```

## 功能验证

### 测试场景

1. **仪表盘订单状态卡片点击**
   - 点击"待处理"卡片 → 导航到 `/orders?status=pending`
   - 订单列表应自动显示状态筛选为"待处理"
   - 表格数据应只显示待处理订单

2. **仪表盘快捷入口点击**
   - 点击"待处理订单"快捷入口 → 导航到 `/orders?status=pending`
   - 点击"进行中订单"快捷入口 → 导航到 `/orders?status=in_progress`
   - 每次都应正确应用状态筛选

3. **URL直接访问**
   - 直接访问 `/orders?status=completed`
   - 页面应加载时就应用"已完成"状态筛选

4. **URL参数变化**
   - 在订单列表页面，手动修改URL中的status参数
   - 页面应自动更新筛选并重新加载数据

## 技术亮点

1. **通用性** - `urlParams.ts` 工具函数可复用于任何需要URL参数的场景
2. **Hook增强** - `useListPage` hook 优雅地集成URL参数支持，保持向后兼容
3. **自动化** - URL参数变化时自动更新，无需手动干预
4. **类型安全** - 完全支持TypeScript类型推断

## 未来改进

1. **支持更多页面** - 可以为GameList、PlayerList等页面添加URL参数支持
2. **URL同步** - 考虑在筛选条件变化时同步更新URL参数
3. **参数验证** - 添加URL参数的验证和错误处理
4. **历史记录** - 利用浏览器历史记录实现前进/后退功能

