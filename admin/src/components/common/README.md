# 统一 UI 组件使用指南

## SearchFilters 组件

统一的标准搜索/筛选区域组件，替代各页面不一致的筛选布局。

### 基础用法

```tsx
import { useState } from 'react';
import { Button } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import SearchFilters, { type FilterItem } from '@/components/common/SearchFilters';

const MyPage: React.FC = () => {
    const [filters, setFilters] = useState({
        keyword: '',
        status: undefined,
        dateRange: undefined,
    });

    const filterConfigs: FilterItem[] = [
        {
            type: 'input',
            key: 'keyword',
            placeholder: '搜索名称',
            width: 200,
        },
        {
            type: 'select',
            key: 'status',
            placeholder: '选择状态',
            options: [
                { label: '启用', value: 'active' },
                { label: '禁用', value: 'inactive' },
            ],
        },
        {
            type: 'rangePicker',
            key: 'dateRange',
        },
    ];

    return (
        <SearchFilters
            filters={filterConfigs}
            filterValues={filters}
            onFilterChange={(key, value) => setFilters({ ...filters, [key]: value })}
            actions={
                <Button type="primary" icon={<PlusOutlined />}>
                    新建
                </Button>
            }
        />
    );
};
```

### Segmented 模式（类似 VIP 页面）

```tsx
const filterConfigs: FilterItem[] = [
    {
        type: 'segmented',
        key: 'viewMode',
        segmentedOptions: [
            { label: '全部', value: 'all' },
            { label: '已启用', value: 'active' },
            { label: '已禁用', value: 'inactive' },
        ],
    },
    {
        type: 'input',
        key: 'keyword',
        placeholder: '搜索等级名称',
        width: 300,
    },
];
```

### 带查询/重置按钮

```tsx
<SearchFilters
    filters={filterConfigs}
    filterValues={filters}
    onFilterChange={(key, value) => setFilters({ ...filters, [key]: value })}
    showQueryButtons={true}
    onQuery={() => fetchData()}
    onReset={() => setFilters({ keyword: '', status: undefined })}
    loading={loading}
/>
```

---

## BatchActions 组件

统一的批量操作组件，支持行内按钮和弹窗两种模式。

### 基础用法（行内模式）

```tsx
import BatchActions, { type BatchAction } from '@/components/common/BatchActions';

const MyPage: React.FC = () => {
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    const batchActions: BatchAction[] = [
        {
            key: 'enable',
            label: '批量启用',
            icon: <CheckCircleOutlined />,
            mode: 'inline',
            onConfirm: async (keys) => {
                await updateStatus(keys, 'active');
                message.success('启用成功');
            },
        },
        {
            key: 'disable',
            label: '批量禁用',
            icon: <StopOutlined />,
            type: 'danger',
            mode: 'inline',
            onConfirm: async (keys) => {
                await updateStatus(keys, 'inactive');
                message.success('禁用成功');
            },
        },
    ];

    return (
        <>
            <Table
                rowSelection={{
                    selectedRowKeys,
                    onChange: setSelectedRowKeys,
                }}
                ...
            />

            <BatchActions
                selectedCount={selectedRowKeys.length}
                actions={batchActions}
                selectedRowKeys={selectedRowKeys}
                onActionComplete={() => {
                    setSelectedRowKeys([]);
                    fetchData();
                }}
            />
        </>
    );
};
```

### 弹窗模式（需要额外输入）

```tsx
const [form] = Form.useForm();

const batchActions: BatchAction[] = [
    {
        key: 'assignRole',
        label: '分配角色',
        mode: 'modal',
        modalTitle: '批量分配角色',
        modalContent: (
            <Form form={form}>
                <Form.Item name="roleId" label="选择角色" rules={[{ required: true }]}>
                    <Select options={roleOptions} />
                </Form.Item>
            </Form>
        ),
        onConfirm: async (keys) => {
            const values = await form.validateFields();
            await assignRoles(keys, values.roleId);
            message.success('分配成功');
        },
    },
];
```

---

## 迁移指南

### Service 页面迁移

**迁移前**:
```tsx
<Row gutter={[16, 16]} justify="space-between">
    <Col>
        <Space>
            <Input placeholder="搜索服务名称" style={{ width: 200 }} />
            <Select style={{ width: 150 }} />
            {/* 更多筛选器... */}
        </Space>
    </Col>
    <Col>
        <Space>
            <Button type="primary">新建服务</Button>
        </Space>
    </Col>
</Row>
```

**迁移后**:
```tsx
<SearchFilters
    filters={[
        { type: 'input', key: 'keyword', placeholder: '搜索服务名称', width: 200 },
        { type: 'select', key: 'gameId', placeholder: '选择游戏', width: 150, options: gameOptions },
        // ...
    ]}
    filterValues={filters}
    onFilterChange={(key, value) => setFilters({ ...filters, [key]: value })}
    actions={<Button type="primary" icon={<PlusOutlined />}>新建服务</Button>}
    showQueryButtons
    onQuery={() => fetchData(1)}
    onReset={() => setFilters({ keyword: '', gameId: undefined })}
/>
```

### Alert 页面迁移

**迁移前**:
```tsx
<Row gutter={16}>
    <Col flex="auto">
        <Space wrap>
            <Select style={{ width: 120 }} />
            <Select style={{ width: 120 }} />
        </Space>
    </Col>
    <Col>
        <Space>
            <Button>刷新</Button>
            <Button>批量已读</Button>
        </Space>
    </Col>
</Row>
```

**迁移后**:
```tsx
<SearchFilters
    card={false}  // Alert 页面不需要 Card
    filters={[
        { type: 'select', key: 'level', placeholder: '告警级别', width: 120, options: levelOptions },
        { type: 'select', key: 'type', placeholder: '告警类型', width: 120, options: typeOptions },
    ]}
    filterValues={filters}
    onFilterChange={(key, value) => setFilters({ ...filters, [key]: value })}
    actions={
        <Space>
            <Button icon={<ReloadOutlined />} onClick={loadAlerts}>刷新</Button>
            <Button icon={<CheckOutlined />} disabled={selectedRowKeys.length === 0}>
                批量已读
            </Button>
        </Space>
    }
/>
```

---

## 优势总结

| 维度 | 迁移前 | 迁移后 |
|------|--------|--------|
| **代码一致性** | 3 种不同布局模式 | 统一组件 |
| **维护成本** | 修改需同步多个页面 | 只修改组件 |
| **新页面开发** | 复制粘贴现有代码 | 配置化使用 |
| **类型安全** | 部分页面缺少类型 | 完整 TypeScript 支持 |
| **响应式** | 不一致 | 统一处理 |
