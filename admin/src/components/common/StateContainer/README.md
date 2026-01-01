# 状态组件使用指南

## EmptyState - 空状态组件

用于展示无数据、无搜索结果等场景。

### 基础用法

```typescript
import { EmptyState } from '@/components/common/StateContainer';

<EmptyState type="no-data" />
<EmptyState
    type="no-search"
    actionText="清除搜索"
    onAction={() => setSearchParams({})}
/>
```

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type | `'no-data' \| 'no-search' \| 'no-permission' \| 'error' \| 'custom'` | `'no-data'` | 空状态类型 |
| title | `string` | - | 标题文本（覆盖预设值） |
| description | `string` | - | 描述文本（覆盖预设值） |
| actionText | `string` | - | 操作按钮文本 |
| onAction | `() => void` | - | 操作按钮回调 |
| showImage | `boolean` | `true` | 是否显示图片 |

### 预设类型

- **no-data**: 暂无数据
- **no-search**: 未找到相关结果
- **no-permission**: 暂无访问权限
- **error**: 加载失败
- **custom**: 自定义（需提供 title 和 description）

## LoadingState - 加载状态组件

提供统一的骨架屏加载效果。

### 基础用法

```typescript
import { LoadingState } from '@/components/common/StateContainer';

<LoadingState loading={true} rows={4} />
<LoadingState
    loading={true}
    card
    title="数据统计"
    rows={3}
/>
```

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| loading | `boolean` | `true` | 是否加载中 |
| card | `boolean` | `true` | 是否显示卡片包裹 |
| title | `string` | - | 卡片标题 |
| rows | `number` | `3` | 骨架屏行数 |
| children | `ReactNode` | - | 非加载状态下显示的内容 |

## StateContainer - 状态容器

统一管理加载、空状态、错误状态的容器组件。

### 基础用法

```typescript
import { StateContainer } from '@/components/common/StateContainer';

<StateContainer
    loading={loading}
    data={list}
    error={error}
    emptyType="no-data"
    emptyActionText="创建新数据"
    onEmptyAction={() => setVisible(true)}
>
    <Table dataSource={list} />
</StateContainer>
```

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| loading | `boolean` | `false` | 是否加载中 |
| data | `any[] \| any` | - | 数据数组或单个数据对象 |
| error | `string \| null` | - | 错误信息 |
| emptyType | `'no-data' \| 'no-search' \| 'no-permission'` | `'no-data'` | 空状态类型 |
| emptyTitle | `string` | - | 空状态标题 |
| emptyDescription | `string` | - | 空状态描述 |
| emptyActionText | `string` | - | 空状态操作按钮文本 |
| onEmptyAction | `() => void` | - | 空状态操作回调（也用于错误重试） |
| loadingConfig | `object` | - | 加载状态配置（card, title, rows） |
| children | `ReactNode` | - | 正常状态下的子内容 |

### 状态判断逻辑

1. **错误状态优先**: 如果 `error` 存在，显示错误空状态
2. **加载状态其次**: 如果 `loading` 为 true，显示骨架屏
3. **空状态判断**: 如果 `data` 为空数组或空值，显示空状态
4. **正常状态**: 否则显示 `children`

### 完整示例

```typescript
import { StateContainer } from '@/components/common/StateContainer';
import { Table, Button } from 'antd';

function UserList() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [users, setUsers] = useState<User[]>([]);

    const loadData = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await userApi.getUsers();
            setUsers(response.data);
        } catch (err) {
            setError('加载失败，请重试');
        } finally {
            setLoading(false);
        }
    };

    return (
        <StateContainer
            loading={loading}
            data={users}
            error={error}
            emptyType="no-data"
            emptyTitle="还没有用户"
            emptyDescription="创建第一个用户开始使用"
            emptyActionText="创建用户"
            onEmptyAction={() => setModalVisible(true)}
            loadingConfig={{
                card: false,
                rows: 5
            }}
        >
            <Table dataSource={users} columns={columns} />
        </StateContainer>
    );
}
```

## 最佳实践

1. **始终使用 StateContainer**: 对于需要处理加载、空状态、错误状态的列表页面，使用 StateContainer 可以大幅简化代码
2. **提供友好的空状态操作**: 在空状态时提供明确的操作指引（如"创建第一条数据"）
3. **错误状态允许重试**: 错误状态应提供"重新加载"按钮
4. **保持文案一致**: 使用预设类型时，如需自定义文案，尽量保持风格统一
5. **合理配置加载状态**: 根据列表内容调整骨架屏行数，提升用户体验
