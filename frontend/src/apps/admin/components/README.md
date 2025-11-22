# Admin Components

## 📁 目录用途

该目录用于存放管理后台（Admin）专用的业务组件。

## 📦 组件分类

### 页面级组件
- 订单管理组件
- 用户管理组件
- 陪玩师管理组件
- 游戏配置组件
- 数据统计组件
- 权限管理组件
- 财务管理组件
- 系统配置组件

### 通用业务组件
- 数据表格（带分页、筛选、排序）
- 表单组件（增删改查）
- 详情弹窗/抽屉
- 批量操作组件
- 数据导出组件
- 图表统计组件
- 文件上传组件
- 日志查看组件

## 📋 命名规范

### 文件夹命名
```
ComponentName/
├── ComponentName.tsx       // 主组件
├── ComponentName.types.ts  // 类型定义
├── index.ts                // 导出
└── components/             // 子组件（如有）
    └── SubComponent.tsx
```

### 组件命名
- 使用 PascalCase（大驼峰）
- 以组件功能命名，例如：`OrderTable`, `UserForm`
- 避免使用通用名称，如：`List`, `Item`

## 🎯 开发规范

### TypeScript
- 所有组件使用 TypeScript
- 定义 Props 接口
- 使用明确的返回值类型

```typescript
// ✅ 推荐
interface OrderTableProps {
  orders: Order[];
  loading?: boolean;
  onViewDetail?: (id: number) => void;
}

export const OrderTable: React.FC<OrderTableProps> = ({ orders, loading }) => {
  return <div>...</div>;
};

// ❌ 避免
export const OrderTable = (props: any) => {
  return <div>...</div>;
};
```

### 组件设计原则
1. **单一职责**: 一个组件只做一件事
2. **可复用性**: 提高组件复用率
3. **可组合性**: 支持组件组合
4. **类型安全**: 使用 TypeScript 严格类型
5. **性能优化**: 使用 React.memo, useMemo, useCallback

### 样式规范
- 优先使用 Tailwind CSS 原子类
- 自定义组件样式放在 `src/styles/tailwind.css`
- 避免内联样式
- 使用响应式类名（md:, lg:）

```tsx
// ✅ 推荐
<div className="card p-6">
  <button className="btn btn-primary">按钮</button>
</div>

// ❌ 避免
<div style={{ padding: '24px', background: 'white' }}>
  <button style={{ background: 'blue', color: 'white' }}>按钮</button>
</div>
```

## 🔗 相关资源

- [全局组件目录](../../shared/components/) - 通用组件（跨应用）
- [管理页面目录](../pages/) - 页面级组件
- [Tailwind CSS 使用指南](../../../docs/STYLE_AND_STATE_MANAGEMENT.md)

## 📚 示例代码

### 创建新组件

```bash
# 在 admin/components 目录下
cd src/apps/admin/components

# 创建组件目录
mkdir UserForm
cd UserForm

# 创建组件文件
touch UserForm.tsx UserForm.types.ts index.ts
```

### 组件结构示例

```typescript
// UserForm.types.ts
export interface UserFormProps {
  initialData?: User;
  onSubmit: (data: UserFormData) => Promise<void>;
  onCancel?: () => void;
}

export interface UserFormData {
  username: string;
  email: string;
  role: 'admin' | 'user' | 'player';
}
```

```typescript
// UserForm.tsx
import React, { useState } from 'react';
import type { UserFormProps, UserFormData } from './UserForm.types';

export const UserForm: React.FC<UserFormProps> = ({
  initialData,
  onSubmit,
  onCancel
}) => {
  const [formData, setFormData] = useState<UserFormData>({
    username: initialData?.username || '',
    email: initialData?.email || '',
    role: initialData?.role || 'user',
  });

  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await onSubmit(formData);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card p-6 max-w-2xl">
      <h2 className="text-xl font-bold mb-4">
        {initialData ? '编辑用户' : '创建用户'}
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* 表单字段 */}
        <div>
          <label className="block text-sm font-medium mb-1">用户名</label>
          <input
            type="text"
            value={formData.username}
            onChange={(e) => setFormData({...formData, username: e.target.value})}
            className="input w-full"
            required
          />
        </div>

        {/* 提交按钮 */}
        <div className="flex justify-end gap-2 pt-4">
          {onCancel && (
            <button type="button" onClick={onCancel} className="btn btn-secondary">
              取消
            </button>
          )}
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? '提交中...' : '提交'}
          </button>
        </div>
      </form>
    </div>
  );
};
```

```typescript
// index.ts
export { UserForm } from './UserForm';
export type { UserFormProps, UserFormData } from './UserForm.types';
```

## 📝 注意事项

1. **避免复制全局组件** - 先检查 [shared/components](../../shared/components/) 是否已有通用组件
2. **保持组件纯函数** - 避免副作用，易于测试
3. **提供默认值** - 提高组件健壮性
4. **处理边界情况** - 空状态、加载状态、错误状态
5. **添加注释** - 复杂逻辑需要说明

## 🎨 设计参考

- **Arco Design** - 企业级设计系统
- **Ant Design** - React UI 库
- **Shadcn/ui** - 可复用的组件库

---

**最后更新**: 2025-11-22
**维护者**: GameLink 前端团队
