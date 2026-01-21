# GameLink Admin - 组件使用文档

本文档介绍 GameLink Admin 面板的核心组件使用方法。

---

## 目录

1. [PermissionGuard - 权限守卫](#permissionguard---权限守卫)
2. [SearchTable - 搜索表格](#searchtable---搜索表格)
3. [PageContainer - 页面容器](#pagecontainer---页面容器)
4. [自定义 Hooks](#自定义-hooks)

---

## PermissionGuard - 权限守卫

权限守卫组件，用于根据用户权限控制 UI 元素的显示/隐藏/禁用。

### 导入

```tsx
import { PermissionGuard, PermissionButton } from '@/components';
```

### 权限格式

```
模块.资源.操作
```

示例:
- `user.list.view` - 查看用户列表
- `order.edit` - 编辑订单
- `player.rank.update` - 更新陪练师等级
- `*` - 超级管理员（所有权限）

### PermissionGuard 组件

#### 基础用法

```tsx
// 有权限则显示，无权限则隐藏
<PermissionGuard permission="user.delete">
  <Button danger>删除用户</Button>
</PermissionGuard>
```

#### 多权限（任一满足）

```tsx
<PermissionGuard
  permission={['user.edit', 'user.delete']}
  mode="any"
>
  <Button>操作用户</Button>
</PermissionGuard>
```

#### 多权限（全部满足）

```tsx
<PermissionGuard
  permission={['user.view', 'user.edit']}
  mode="all"
>
  <Button>编辑用户</Button>
</PermissionGuard>
```

#### 带有 Fallback

```tsx
<PermissionGuard
  permission="admin.settings"
  fallback={<span>权限不足</span>}
>
  <Button>系统设置</Button>
</PermissionGuard>
```

#### 禁用模式

无权限时禁用按钮而非隐藏：

```tsx
<PermissionGuard permission="order.cancel" disabled>
  <Button>取消订单</Button>
</PermissionGuard>
```

#### 加载状态

```tsx
<PermissionGuard
  permission="user.view"
  loading={<Skeleton active />}
>
  <UserList />
</PermissionGuard>
```

### PermissionButton 组件

专门为按钮设计的权限包装器。

#### 默认行为（无权限隐藏）

```tsx
<PermissionButton permission="order.create">
  <Button type="primary">创建订单</Button>
</PermissionButton>
```

#### 禁用模式（无权限禁用）

```tsx
<PermissionButton
  permission="order.delete"
  disableOnNoPermission
>
  <Button danger>删除订单</Button>
</PermissionButton>
```

### API

#### PermissionGuard Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `permission` | `string \| string[]` | - | **必填**，所需权限码 |
| `mode` | `'any' \| 'all'` | `'any'` | 多权限检查模式 |
| `children` | `ReactNode` | - | **必填**，有权限时显示的内容 |
| `fallback` | `ReactNode` | `null` | 无权限时显示的内容 |
| `loading` | `ReactNode` | `null` | 权限加载中显示的内容 |
| `disabled` | `boolean` | `false` | 无权限时是否禁用子组件 |

#### PermissionButton Props

继承 `PermissionGuard` 所有属性，额外增加：

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `disableOnNoPermission` | `boolean` | `false` | 无权限时禁用按钮而非隐藏 |

---

## SearchTable - 搜索表格

封装表格搜索、分页、批量操作等通用功能的高级组件。

### 导入

```tsx
import { SearchTable } from '@/components';
import type { SearchField, ToolbarButton } from '@/components';
```

### 基础用法

```tsx
const columns = [
  {
    title: '用户名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: '操作',
    key: 'action',
    render: (_, record) => (
      <Space>
        <Button onClick={() => handleEdit(record)}>编辑</Button>
        <Button danger onClick={() => handleDelete(record.id)}>删除</Button>
      </Space>
    ),
  },
];

const searchFields: SearchField[] = [
  { name: 'keyword', label: '搜索', type: 'input', placeholder: '输入用户名或邮箱' },
  { name: 'status', label: '状态', type: 'select', options: [
    { label: '启用', value: 'active' },
    { label: '禁用', value: 'inactive' },
  ]},
];

<SearchTable
  columns={columns}
  dataSource={data}
  loading={loading}
  searchFields={searchFields}
  onSearch={handleSearch}
  onRefresh={fetchData}
  pagination={{
    total,
    current,
    pageSize,
    onChange: handlePageChange,
  }}
/>
```

### 高级用法

#### 带工具栏

```tsx
<SearchTable
  columns={columns}
  dataSource={data}
  cardTitle="用户列表"
  showCreate={true}
  createText="新增用户"
  createPermission="user.create"
  onCreate={handleCreate}
  showBatchDelete={true}
  batchDeletePermission="user.delete"
  onBatchDelete={handleBatchDelete}
  toolbarButtons={[
    {
      text: '导出',
      icon: <ExportOutlined />,
      onClick: handleExport,
    },
    {
      text: '批量启用',
      type: 'primary',
      permission: 'user.enable',
      needSelection: true,
      onClick: handleBatchEnable,
    },
  ]}
/>
```

#### 危险操作确认

```tsx
toolbarButtons={[
  {
    text: '批量删除',
    danger: true,
    needSelection: true,
    confirmText: '确定要删除选中的数据吗？此操作不可撤销！',
    permission: 'user.delete',
    onClick: handleBatchDelete,
  },
]}
```

### API

#### SearchTable Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `columns` | `ColumnType[]` | - | **必填**，表格列配置 |
| `dataSource` | `T[]` | - | 表格数据 |
| `loading` | `boolean` | `false` | 加载状态 |
| `searchFields` | `SearchField[]` | `[]` | 搜索字段配置 |
| `toolbarButtons` | `ToolbarButton[]` | `[]` | 工具栏按钮配置 |
| `showCreate` | `boolean` | `true` | 是否显示新增按钮 |
| `createText` | `string` | `'新增'` | 新增按钮文本 |
| `createPermission` | `string` | - | 新增权限码 |
| `onCreate` | `() => void` | - | 新增回调 |
| `showBatchDelete` | `boolean` | `false` | 是否显示批量删除 |
| `batchDeletePermission` | `string` | - | 批量删除权限码 |
| `onBatchDelete` | `(keys) => Promise<void>` | - | 批量删除回调 |
| `onSearch` | `(values) => void` | - | 搜索回调 |
| `onRefresh` | `() => void` | - | 刷新回调 |
| `cardTitle` | `ReactNode` | - | 卡片标题 |
| `defaultExpanded` | `boolean` | `true` | 搜索栏默认展开 |

继承 Ant Design `Table` 的所有其他属性。

#### SearchField 类型

```typescript
interface SearchField {
  name: string;           // 字段名
  label: string;          // 标签文本
  type: 'input' | 'select' | 'dateRange' | 'numberRange';  // 字段类型
  placeholder?: string;   // 占位符
  options?: Array<{       // 选项（select 类型）
    label: string;
    value: any;
  }>;
}
```

#### ToolbarButton 类型

```typescript
interface ToolbarButton {
  text: string;              // 按钮文本
  type?: 'primary' | 'default' | 'dashed' | 'link' | 'text';
  icon?: ReactNode;          // 图标
  danger?: boolean;          // 危险按钮
  permission?: string;       // 权限码
  onClick?: (keys?: React.Key[]) => void;
  needSelection?: boolean;   // 是否需要选中行
  confirmText?: string;      // 确认提示文本
  simpleAction?: boolean;    // 是否为简单批量操作（有选中时显示）
}
```

---

## PageContainer - 页面容器

提供统一的页面布局，包含标题、操作区、标签页等。

### 导入

```tsx
import { PageContainer } from '@/components';
```

### 基础用法

```tsx
<PageContainer title="用户管理">
  <UserList />
</PageContainer>
```

### 带操作按钮

```tsx
<PageContainer
  title="用户管理"
  extra={
    <Space>
      <Button>导出</Button>
      <Button type="primary">新增用户</Button>
    </Space>
  }
>
  <UserList />
</PageContainer>
```

### 带刷新按钮

```tsx
<PageContainer
  title="实时监控"
  showRefresh={true}
  onRefresh={fetchData}
  loading={loading}
>
  <MonitorDashboard />
</PageContainer>
```

### 带标签页

```tsx
<PageContainer
  title="订单详情"
  tabList={[
    { key: 'info', label: '基本信息' },
    { key: 'timeline', label: '订单时间线' },
    { key: 'payment', label: '支付信息' },
  ]}
  activeTab={activeTab}
  onTabChange={setActiveTab}
>
  {activeTab === 'info' && <OrderInfo />}
  {activeTab === 'timeline' && <OrderTimeline />}
  {activeTab === 'payment' && <PaymentInfo />}
</PageContainer>
```

### API

#### PageContainer Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `title` | `ReactNode` | - | 页面标题 |
| `subTitle` | `ReactNode` | - | 页面副标题 |
| `extra` | `ReactNode` | - | 额外操作区域 |
| `tabList` | `TabsProps['items']` | - | 标签页配置 |
| `activeTab` | `string` | - | 当前激活标签 |
| `onTabChange` | `(key: string) => void` | - | 标签切换回调 |
| `showRefresh` | `boolean` | `false` | 显示刷新按钮 |
| `onRefresh` | `() => void` | - | 刷新回调 |
| `loading` | `boolean` | `false` | 加载状态 |
| `children` | `ReactNode` | - | 页面内容 |
| `className` | `string` | - | 自定义类名 |

---

## 自定义 Hooks

### usePermission

权限检查 Hook。

```tsx
import { usePermission } from '@/hooks/usePermission';

function MyComponent() {
  const { hasPermission, loading } = usePermission('user.delete');

  if (loading) return <Spin />;
  if (!hasPermission) return <div>权限不足</div>;

  return <Button danger>删除</Button>;
}
```

#### API

```typescript
function usePermission(
  permission: string | string[],
  mode?: 'any' | 'all'
): {
  hasPermission: boolean;
  loading: boolean;
}
```

### useWebSocket

WebSocket 连接管理 Hook。

```tsx
import { useWebSocket } from '@/hooks/useWebSocket';

function MonitorPage() {
  const { data, connected, connect, disconnect } = useWebSocket({
    url: import.meta.env.VITE_WS_URL,
    enabled: true,
    onMessage: (msg) => console.log('收到消息:', msg),
  });

  return (
    <div>
      <div>状态: {connected ? '已连接' : '未连接'}</div>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
```

#### API

```typescript
interface UseWebSocketOptions {
  url: string;
  enabled?: boolean;
  reconnect?: boolean;
  reconnectInterval?: number;
  onMessage?: (data: any) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
}

function useWebSocket(options: UseWebSocketOptions): {
  data: any;
  connected: boolean;
  connect: () => void;
  disconnect: () => void;
  send: (data: any) => void;
}
```

---

## 常见使用模式

### 完整 CRUD 页面示例

```tsx
import { PageContainer, SearchTable } from '@/components';
import { PermissionButton } from '@/components';
import { usePermission } from '@/hooks/usePermission';

function UserManagement() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const { hasPermission } = usePermission('user.create');

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: '用户名', dataIndex: 'name' },
    { title: '邮箱', dataIndex: 'email' },
    {
      title: '操作',
      render: (_, record) => (
        <Space>
          <PermissionButton permission="user.edit">
            <Button onClick={() => handleEdit(record)}>编辑</Button>
          </PermissionButton>
          <PermissionButton permission="user.delete" disableOnNoPermission>
            <Button danger onClick={() => handleDelete(record.id)}>删除</Button>
          </PermissionButton>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer
      title="用户管理"
      extra={
        hasPermission && (
          <Button type="primary" onClick={handleCreate}>
            新增用户
          </Button>
        )
      }
    >
      <SearchTable
        columns={columns}
        dataSource={data}
        loading={loading}
        searchFields={[
          { name: 'keyword', label: '搜索', type: 'input' },
          { name: 'status', label: '状态', type: 'select', options: statusOptions },
        ]}
        onSearch={handleSearch}
        onRefresh={fetchData}
        showBatchDelete={true}
        batchDeletePermission="user.delete"
        onBatchDelete={handleBatchDelete}
        pagination={{
          total,
          current,
          pageSize,
          onChange: handlePageChange,
        }}
      />
    </PageContainer>
  );
}
```

---

## 更多组件

其他可用组件可从 `@/components` 导入：

- `Button`, `IconButton` - 统一按钮
- `Card`, `StatisticCard`, `ContentCard` - 卡片组件
- `StatCard` - 统计卡片
- `ThemeToggle` - 主题切换
- `PermissionTree` - 权限树选择器
- `ErrorBoundary`, `PageErrorBoundary` - 错误边界
- `AnimatedContainer`, `AnimatedListItem` - 动画容器
- `CollapsibleSection` - 可折叠区块
- `ImportModal`, `ImportHistoryTable` - 导入功能

查看组件源码获取更多使用示例。
