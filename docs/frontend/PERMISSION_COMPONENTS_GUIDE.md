# 前端权限组件使用指南

本文档介绍 GameLink 平台前端权限控制组件的使用方法，包括 `PermissionGuard` 组件、`usePermission` Hook 等。

## 概述

前端权限控制系统提供以下核心功能：

1. **PermissionGuard 组件** - 声明式权限控制，包裹需要权限保护的元素
2. **PermissionButton 组件** - 带权限控制的按钮组件
3. **usePermission Hook** - 单个权限检查
4. **usePermissions Hook** - 批量权限检查
5. **usePermissionChecker Hook** - 动态权限检查函数
6. **withPermission HOC** - 高阶组件方式添加权限控制

## PermissionGuard 组件

### 基本用法

```tsx
import { PermissionGuard } from '@/components/PermissionGuard';

// 无权限时隐藏按钮
<PermissionGuard permission="admin.users.create">
    <Button type="primary">创建用户</Button>
</PermissionGuard>
```

### 多权限检查

```tsx
// 任一权限满足即可（默认模式）
<PermissionGuard permission={['admin.users.update', 'admin.users.delete']} mode="any">
    <Button>操作</Button>
</PermissionGuard>

// 全部权限都需要满足
<PermissionGuard permission={['admin.users.read', 'admin.users.update']} mode="all">
    <Button>编辑</Button>
</PermissionGuard>
```

### 禁用模式

无权限时显示禁用状态的按钮，而不是隐藏：

```tsx
<PermissionGuard 
    permission="admin.users.delete" 
    disabled 
    tooltip="您没有删除权限"
>
    <Button danger>删除</Button>
</PermissionGuard>
```

### 自定义加载状态

```tsx
import { Spin } from 'antd';

<PermissionGuard 
    permission="admin.users.create" 
    loading={<Spin size="small" />}
>
    <Button>创建</Button>
</PermissionGuard>
```

### 自定义回退内容

```tsx
<PermissionGuard 
    permission="admin.users.create" 
    fallback={<span className="text-gray-400">无权限</span>}
>
    <Button>创建</Button>
</PermissionGuard>
```

### Props 说明

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| permission | `string \| string[]` | - | 需要检查的权限码 |
| mode | `'any' \| 'all'` | `'any'` | 权限检查模式 |
| children | `ReactNode` | - | 子元素 |
| fallback | `ReactNode` | `null` | 无权限时的回退内容 |
| loading | `ReactNode` | `null` | 加载中时显示的内容 |
| disabled | `boolean` | `false` | 是否使用禁用模式 |
| tooltip | `string` | `'您没有此操作的权限'` | 禁用模式下的提示文本 |

## PermissionButton 组件

带权限控制的按钮组件，简化常见场景：

```tsx
import { PermissionButton } from '@/components/PermissionGuard';

// 无权限时显示禁用按钮（默认）
<PermissionButton permission="admin.users.delete" tooltip="没有删除权限">
    删除
</PermissionButton>

// 无权限时隐藏按钮
<PermissionButton permission="admin.users.delete" hideOnNoPermission>
    删除
</PermissionButton>

// 支持所有 Ant Design Button 属性
<PermissionButton 
    permission="admin.users.create" 
    type="primary" 
    icon={<PlusOutlined />}
>
    创建用户
</PermissionButton>
```

### Props 说明

继承所有 Ant Design `ButtonProps`，额外属性：

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| permission | `string \| string[]` | - | 需要检查的权限码 |
| mode | `'any' \| 'all'` | `'any'` | 权限检查模式 |
| tooltip | `string` | `'您没有此操作的权限'` | 无权限时的提示文本 |
| hideOnNoPermission | `boolean` | `false` | 是否在无权限时隐藏 |

## usePermission Hook

### 基本用法

```tsx
import { usePermission } from '@/hooks/usePermission';

function MyComponent() {
    const { hasPermission, loading } = usePermission('admin.games.create');

    if (loading) {
        return <Spin />;
    }

    return hasPermission ? <Button>创建游戏</Button> : null;
}
```

### 多权限检查

```tsx
// 任一权限满足
const { hasPermission } = usePermission(
    ['admin.games.create', 'admin.games.update']
);

// 全部权限满足
const { hasPermission } = usePermission(
    ['admin.games.create', 'admin.games.update'], 
    'all'
);
```

### 返回值

| 属性 | 类型 | 说明 |
|------|------|------|
| hasPermission | `boolean` | 是否拥有权限 |
| loading | `boolean` | 权限检查是否正在加载 |

## usePermissions Hook

批量检查多个权限，返回每个权限的检查结果：

```tsx
import { usePermissions } from '@/hooks/usePermission';

function MyComponent() {
    const permissions = usePermissions({
        canCreate: 'admin.games.create',
        canEdit: 'admin.games.update',
        canDelete: 'admin.games.delete'
    });

    if (permissions.loading) {
        return <Spin />;
    }

    return (
        <Space>
            {permissions.canCreate && <Button type="primary">创建</Button>}
            {permissions.canEdit && <Button>编辑</Button>}
            {permissions.canDelete && <Button danger>删除</Button>}
        </Space>
    );
}
```

### 返回值

返回一个对象，包含：
- 每个传入 key 对应的 boolean 值
- `loading`: 权限检查是否正在加载

## useHasPermission Hook

简化版本，直接返回 boolean：

```tsx
import { useHasPermission } from '@/hooks/usePermission';

function MyComponent() {
    const canCreate = useHasPermission('admin.games.create');

    return canCreate && <Button>创建</Button>;
}
```

## usePermissionChecker Hook

返回一个检查函数，可以动态检查权限：

```tsx
import { usePermissionChecker } from '@/hooks/usePermission';

function MyComponent() {
    const checkPermission = usePermissionChecker();

    const handleClick = () => {
        if (checkPermission('admin.games.delete')) {
            // 执行删除操作
            deleteGame();
        } else {
            message.error('没有删除权限');
        }
    };

    // 支持批量检查
    const canOperate = checkPermission(
        ['admin.games.update', 'admin.games.delete'], 
        'any'
    );

    return <Button onClick={handleClick}>删除</Button>;
}
```

### 带加载状态的版本

```tsx
import { usePermissionCheckerWithLoading } from '@/hooks/usePermission';

function MyComponent() {
    const { check, loading } = usePermissionCheckerWithLoading();

    if (loading) {
        return <Spin />;
    }

    const handleClick = () => {
        if (check('admin.games.delete')) {
            deleteGame();
        }
    };

    return <Button onClick={handleClick}>删除</Button>;
}
```

## withPermission 高阶组件

为组件添加权限控制：

```tsx
import { withPermission } from '@/components/PermissionGuard';
import { Button } from 'antd';

// 创建带权限控制的按钮
const ProtectedButton = withPermission(Button, 'admin.users.create');

// 使用
<ProtectedButton type="primary">创建用户</ProtectedButton>

// 带选项
const ProtectedButton = withPermission(Button, 'admin.users.create', {
    mode: 'all',
    fallback: DisabledButton,
    loading: <Spin size="small" />
});
```

## AdminContext

权限数据通过 `AdminContext` 提供：

```tsx
import { useAdmin } from '@/context/AdminContext';

function MyComponent() {
    const { 
        permissions,    // 权限码数组
        loading,        // 加载状态
        hasPermission,  // 权限检查函数
        isSuperAdmin,   // 是否超级管理员
        refreshMenus    // 刷新菜单和权限
    } = useAdmin();

    // 超级管理员拥有所有权限
    if (isSuperAdmin) {
        return <AdminPanel />;
    }

    // 使用 hasPermission 函数
    if (hasPermission(['admin.users.read'], 'any')) {
        return <UserList />;
    }

    return <NoPermission />;
}
```

## 最佳实践

### 1. 优先使用声明式组件

```tsx
// ✅ 推荐：使用 PermissionGuard
<PermissionGuard permission="admin.users.create">
    <Button>创建</Button>
</PermissionGuard>

// ❌ 不推荐：手动检查
const { hasPermission } = usePermission('admin.users.create');
return hasPermission ? <Button>创建</Button> : null;
```

### 2. 批量检查使用 usePermissions

```tsx
// ✅ 推荐：一次性检查多个权限
const perms = usePermissions({
    canCreate: 'admin.users.create',
    canEdit: 'admin.users.update',
    canDelete: 'admin.users.delete'
});

// ❌ 不推荐：多次调用 usePermission
const { hasPermission: canCreate } = usePermission('admin.users.create');
const { hasPermission: canEdit } = usePermission('admin.users.update');
const { hasPermission: canDelete } = usePermission('admin.users.delete');
```

### 3. 处理加载状态

```tsx
// ✅ 推荐：显示加载状态，避免闪烁
<PermissionGuard 
    permission="admin.users.create" 
    loading={<Button disabled loading>创建</Button>}
>
    <Button>创建</Button>
</PermissionGuard>
```

### 4. 使用禁用模式提供更好的用户体验

```tsx
// ✅ 推荐：无权限时显示禁用按钮，让用户知道功能存在
<PermissionGuard 
    permission="admin.users.delete" 
    disabled 
    tooltip="需要删除权限"
>
    <Button danger>删除</Button>
</PermissionGuard>
```

### 5. 动态检查使用 usePermissionChecker

```tsx
// ✅ 推荐：事件处理中使用 checker
const checkPermission = usePermissionChecker();

const handleAction = (action: string) => {
    const permissionCode = `admin.users.${action}`;
    if (checkPermission(permissionCode)) {
        performAction(action);
    } else {
        message.error('没有权限');
    }
};
```

## 权限码命名规范

权限码采用三段式命名：`{module}.{resource}.{action}`

常用权限码示例：
- `admin.users.read` - 查看用户
- `admin.users.create` - 创建用户
- `admin.users.update` - 更新用户
- `admin.users.delete` - 删除用户
- `admin.orders.read` - 查看订单
- `admin.orders.refund` - 退款订单

完整权限码列表请参考 [权限码列表文档](../backend/PERMISSION_CODES.md)。

## 注意事项

1. **超级管理员**：拥有 `*` 权限的用户自动通过所有权限检查
2. **加载状态**：权限数据加载期间，默认不显示受保护的元素
3. **缓存**：权限数据会被缓存，权限变更后需要刷新页面或调用 `refreshMenus()`
4. **性能**：Hook 内部使用 `useMemo` 优化，避免不必要的重渲染
