# 按钮级权限控制使用指南

## 概述

本项目提供了完整的按钮级权限控制方案，包括：

1. **PermissionGuard 组件** - 声明式权限控制
2. **usePermission Hook** - 编程式权限检查
3. **权限工具函数** - 非 React 环境下的权限检查
4. **权限常量** - 统一的权限码定义

## 快速开始

### 1. 基础用法 - PermissionGuard 组件

```tsx
import { PermissionGuard } from '@/components/PermissionGuard';
import { GAME_PERMISSIONS } from '@/constants/permissions';

// 有权限则显示按钮
<PermissionGuard permission={GAME_PERMISSIONS.CREATE}>
    <Button type="primary">创建游戏</Button>
</PermissionGuard>

// 无权限时显示提示
<PermissionGuard
    permission={GAME_PERMISSIONS.DELETE}
    fallback={<span className="text-gray-400">无删除权限</span>}
>
    <Button danger>删除</Button>
</PermissionGuard>

// 无权限时禁用按钮（而非隐藏）
<PermissionGuard permission={GAME_PERMISSIONS.UPDATE} disabled>
    <Button>编辑</Button>
</PermissionGuard>
```

### 2. 多权限检查

```tsx
// 任一权限满足即可（默认模式）
<PermissionGuard permission={[GAME_PERMISSIONS.CREATE, GAME_PERMISSIONS.UPDATE]}>
    <Button>操作</Button>
</PermissionGuard>

// 全部权限都需要
<PermissionGuard
    permission={[GAME_PERMISSIONS.READ, GAME_PERMISSIONS.UPDATE]}
    mode="all"
>
    <Button>编辑详情</Button>
</PermissionGuard>
```

### 3. 使用 Hook - usePermission

```tsx
import { usePermission, useHasPermission, usePermissions } from '@/hooks/usePermission';
import { GAME_PERMISSIONS } from '@/constants/permissions';

function GamePage() {
    // 方式1：获取完整检查结果
    const { hasPermission, loading } = usePermission(GAME_PERMISSIONS.CREATE);

    // 方式2：简化版本，直接返回 boolean
    const canDelete = useHasPermission(GAME_PERMISSIONS.DELETE);

    // 方式3：批量检查多个权限
    const permissions = usePermissions({
        canCreate: GAME_PERMISSIONS.CREATE,
        canEdit: GAME_PERMISSIONS.UPDATE,
        canDelete: GAME_PERMISSIONS.DELETE,
    });

    return (
        <div>
            {hasPermission && <Button type="primary">创建</Button>}
            {canDelete && <Button danger>删除</Button>}

            {permissions.canEdit && <Button>编辑</Button>}
        </div>
    );
}
```

### 4. 使用 usePermissionChecker 进行动态检查

```tsx
import { usePermissionChecker } from '@/hooks/usePermission';

function ActionComponent() {
    const checkPermission = usePermissionChecker();

    const handleAction = (action: string) => {
        if (action === 'delete') {
            if (!checkPermission(GAME_PERMISSIONS.DELETE)) {
                message.error('您没有删除权限');
                return;
            }
            // 执行删除操作
        }
    };

    return <Button onClick={() => handleAction('delete')}>删除</Button>;
}
```

### 5. 使用 AdminContext

```tsx
import { useAdmin } from '@/context/AdminContext';

function MyComponent() {
    const { hasPermission, hasAllPermissions, hasAnyPermission, isSuperAdmin } = useAdmin();

    // 检查单个权限
    if (hasPermission('admin.games.create')) {
        // 可以创建
    }

    // 检查是否为超级管理员
    if (isSuperAdmin) {
        // 显示所有功能
    }

    // 检查多个权限（全部满足）
    if (hasAllPermissions(['admin.games.read', 'admin.games.update'])) {
        // 可以查看和编辑
    }

    return null;
}
```

### 6. 非 React 环境中使用

```tsx
import { hasPermission, filterActionsByPermission } from '@/utils/permission';

// 在普通函数中检查权限
function doSomething() {
    if (hasPermission('admin.games.delete')) {
        // 执行删除
    }
}

// 过滤操作列表
const actions = [
    { key: 'edit', label: '编辑', permission: 'admin.games.update' },
    { key: 'delete', label: '删除', permission: 'admin.games.delete' },
];
const allowedActions = filterActionsByPermission(actions);
```

### 7. 表格操作列权限控制

```tsx
import { PermissionGuard } from '@/components/PermissionGuard';
import { GAME_PERMISSIONS } from '@/constants/permissions';

const columns = [
    // ...其他列
    {
        title: '操作',
        key: 'action',
        render: (_, record) => (
            <Space>
                <PermissionGuard permission={GAME_PERMISSIONS.UPDATE}>
                    <Button type="link" onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                </PermissionGuard>
                <PermissionGuard permission={GAME_PERMISSIONS.DELETE}>
                    <Button type="link" danger onClick={() => handleDelete(record)}>
                        删除
                    </Button>
                </PermissionGuard>
            </Space>
        ),
    },
];
```

### 8. 使用 withPermission HOC

```tsx
import { withPermission } from '@/components/PermissionGuard';
import { Button } from 'antd';

// 创建带权限控制的按钮组件
const CreateButton = withPermission(Button, GAME_PERMISSIONS.CREATE);
const DeleteButton = withPermission(Button, GAME_PERMISSIONS.DELETE, {
    disabled: true, // 无权限时禁用而非隐藏
});

function MyPage() {
    return (
        <div>
            <CreateButton type="primary">创建游戏</CreateButton>
            <DeleteButton danger>删除游戏</DeleteButton>
        </div>
    );
}
```

## 权限码规范

权限码格式：`{module}.{resource}.{action}`

常用操作类型：
- `list` - 列表查询
- `read` - 详情查询
- `create` - 创建
- `update` - 更新
- `delete` - 删除

示例：
- `admin.games.list` - 游戏列表
- `admin.games.create` - 创建游戏
- `admin.users.update` - 更新用户

## 注意事项

1. **权限加载时机**：权限数据在 `AdminProvider` 初始化时自动加载，登录成功后会自动刷新。

2. **超级管理员**：权限码包含 `*` 的用户为超级管理员，自动拥有所有权限。

3. **默认行为**：
   - 权限加载中时，默认返回 `false`（无权限）
   - 无权限时默认隐藏元素（可通过 `fallback` 或 `disabled` 属性改变行为）

4. **性能优化**：
   - 权限检查使用 `useMemo` 优化，避免不必要的重新计算
   - 权限数据会同步到 `permissionStore`，支持非 React 环境访问

5. **TypeScript 支持**：
   - 所有 Hook 和组件都有完整的类型定义
   - 使用权限常量可以获得类型提示
