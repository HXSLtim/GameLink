/**
 * 服务项目管理页面
 */
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    Switch,
    Image,
    App,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    CheckCircleOutlined,
    StopOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { adminApi, type ServiceItem } from '@/api/admin';

/**
 * 服务分类映射
 */
const categoryMap = {
    solo: '单人护航',
    team: '团队护航',
    gift: '礼物',
};

/**
 * 服务项目管理页面
 */
const ServiceItemList: React.FC = () => {
    const { message } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<ServiceItem[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    // 编辑弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [currentItem, setCurrentItem] = useState<ServiceItem | null>(null);
    const [form] = Form.useForm();

    /**
     * 加载数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const res = await adminApi.getServiceItems({
                ...searchParams,
                page: current,
                page_size: pageSize,
                game_id: searchParams.gameId as number | undefined,
            });

            if (res.data.success) {
                const items = res.data.data;
                setData(Array.isArray(items) ? items : []);
                setTotal((res.data as { pagination?: { total?: number } }).pagination?.total || 0);
            } else {
                message.error(res.data.message || '加载数据失败');
                setData([]);
                setTotal(0);
            }
        } catch (error) {
            console.error(error);
            message.error('加载数据失败');
            setData([]);
            setTotal(0);
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, message]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        const params: Record<string, unknown> = {};
        if (values.keyword) params.keyword = values.keyword;
        if (values.gameId) params.gameId = values.gameId;
        if (values.category) params.category = values.category;
        if (values.status) params.status = values.status;
        setSearchParams(params);
        setCurrent(1);
    }, []);

    /**
     * 打开编辑弹窗
     */
    const handleEdit = useCallback((item: ServiceItem) => {
        setCurrentItem(item);
        form.setFieldsValue({
            name: item.name,
            itemCode: item.itemCode,
            subCategory: item.subCategory,
            basePriceCents: item.basePriceCents ? item.basePriceCents / 100 : undefined,
            commissionRate: item.commissionRate,
            sortOrder: item.sortOrder,
            iconUrl: item.iconUrl,
            description: item.description,
            isActive: item.isActive ?? true,
        });
        setEditModalVisible(true);
    }, [form]);

    /**
     * 打开新增弹窗
     */
    const handleCreate = useCallback(() => {
        setCurrentItem(null);
        form.resetFields();
        form.setFieldsValue({ isActive: true, sortOrder: 0 });
        setEditModalVisible(true);
    }, [form]);

    /**
     * 保存
     */
    const handleSave = useCallback(async () => {
        try {
            const values = await form.validateFields();
            const payload = {
                ...values,
                basePriceCents: values.basePriceCents ? Math.round(values.basePriceCents * 100) : 0,
            };

            if (currentItem) {
                // 更新
                const res = await adminApi.updateServiceItem(currentItem.id, payload);
                if (res.data.success) {
                    message.success('更新成功');
                    setEditModalVisible(false);
                    loadData();
                } else {
                    message.error(res.data.message || '更新失败');
                }
            } else {
                // 新增
                const res = await adminApi.createServiceItem(payload);
                if (res.data.success) {
                    message.success('创建成功');
                    setEditModalVisible(false);
                    loadData();
                } else {
                    message.error(res.data.message || '创建失败');
                }
            }
        } catch (error) {
            console.error('Save failed:', error);
            message.error('保存失败');
        }
    }, [currentItem, form, loadData]);

    /**
     * 批量启用/禁用
     */
    const handleBatchStatus = useCallback(async (keys: React.Key[], status: 'active' | 'inactive') => {
        if (keys.length === 0) {
            message.warning('请选择要操作的数据');
            return;
        }
        try {
            await adminApi.batchUpdateServiceItemStatus(keys as number[], status);
            message.success(`批量${status === 'active' ? '启用' : '禁用'}成功`);
            setSelectedRowKeys([]);
            loadData();
        } catch {
            message.error('操作失败');
        }
    }, [loadData]);

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = useMemo(() => [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '搜索服务名称' },
        {
            name: 'gameId',
            label: '游戏',
            type: 'select',
            options: [
                { label: '王者荣耀', value: 1 },
                { label: '英雄联盟', value: 2 },
            ],
        },
        {
            name: 'category',
            label: '服务分类',
            type: 'select',
            options: [
                { label: '上分', value: 'rank' },
                { label: '陪玩', value: 'rush' },
                { label: '教学', value: 'teach' },
                { label: '娱乐', value: 'entertain' },
            ],
        },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: [
                { label: '启用', value: 'active' },
                { label: '禁用', value: 'inactive' },
            ],
        },
    ], []);

    /**
     * 表格列配置
     */
    const columns: ColumnsType<ServiceItem> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            width: 80,
        },
        {
            title: '服务名称',
            dataIndex: 'name',
            render: (text: string, record: ServiceItem) => (
                <Space>
                    {record.iconUrl && <Image src={record.iconUrl} alt="" width={24} height={24} style={{ borderRadius: 4 }} />}
                    <span style={{ fontWeight: 500 }}>{text}</span>
                </Space>
            ),
        },
        {
            title: '服务编码',
            dataIndex: 'itemCode',
            width: 150,
        },
        {
            title: '分类',
            dataIndex: 'subCategory',
            width: 100,
            render: (cat: string) => (
                <Tag color="blue">{categoryMap[cat as keyof typeof categoryMap] || cat}</Tag>
            ),
        },
        {
            title: '价格',
            dataIndex: 'basePriceCents',
            width: 100,
            render: (cents: number | undefined | null) => (
                <span style={{ color: '#faa61a', fontWeight: 'bold' }}>
                    ¥{((cents ?? 0) / 100).toFixed(2)}
                </span>
            ),
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            width: 100,
            render: (isActive: boolean) => (
                <Tag color={isActive ? 'success' : 'default'}>
                    {isActive ? '已启用' : '已禁用'}
                </Tag>
            ),
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            width: 80,
        },
        {
            title: '操作',
            key: 'action',
            width: 180,
            fixed: 'right',
            render: (_: unknown, record: ServiceItem) => (
                <Space size="small">
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={record.isActive ? <StopOutlined /> : <CheckCircleOutlined />}
                        style={{ color: record.isActive ? '#ff4d4f' : '#52c41a' }}
                        onClick={() => handleBatchStatus([record.id], record.isActive ? 'inactive' : 'active')}
                    >
                        {record.isActive ? '禁用' : '启用'}
                    </Button>
                </Space>
            ),
        },
    ], [handleEdit, handleBatchStatus]);

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量启用',
            icon: <CheckCircleOutlined />,
            needSelection: true,
            onClick: (keys) => handleBatchStatus(keys || [], 'active'),
            simpleAction: true,
        },
        {
            text: '批量禁用',
            icon: <StopOutlined />,
            needSelection: true,
            onClick: (keys) => handleBatchStatus(keys || [], 'inactive'),
            simpleAction: true,
            danger: true,
        },
    ];

    return (
        <PageContainer title="服务项目管理" subTitle="管理平台服务项目">
            <SearchTable
                columns={columns}
                dataSource={data}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={loadData}
                loading={loading}
                showCreate={true}
                createText="新建服务"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                rowSelection={{
                    selectedRowKeys,
                    onChange: setSelectedRowKeys,
                }}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1200 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentItem ? '编辑服务' : '新建服务'}
                open={editModalVisible}
                onOk={handleSave}
                onCancel={() => setEditModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="服务名称"
                        rules={[{ required: true, message: '请输入服务名称' }]}
                    >
                        <Input placeholder="请输入服务名称" />
                    </Form.Item>
                    <Form.Item
                        name="itemCode"
                        label="服务编码"
                        rules={[{ required: true, message: '请输入服务编码' }]}
                    >
                        <Input placeholder="请输入服务编码" />
                    </Form.Item>
                    <Form.Item
                        name="subCategory"
                        label="服务分类"
                        rules={[{ required: true, message: '请选择服务分类' }]}
                    >
                        <Select placeholder="请选择服务分类">
                            <Select.Option value="solo">单人护航</Select.Option>
                            <Select.Option value="team">团队护航</Select.Option>
                            <Select.Option value="gift">礼物</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        name="basePriceCents"
                        label="基础价格（元）"
                        rules={[{ required: true, message: '请输入基础价格' }]}
                    >
                        <InputNumber
                            min={0}
                            precision={2}
                            placeholder="请输入基础价格"
                            style={{ width: '100%' }}
                            prefix="¥"
                        />
                    </Form.Item>
                    <Form.Item
                        name="commissionRate"
                        label="佣金比例（%）"
                        rules={[{ required: true, message: '请输入佣金比例' }]}
                    >
                        <InputNumber
                            min={0}
                            max={100}
                            precision={2}
                            placeholder="请输入佣金比例"
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                    <Form.Item
                        name="sortOrder"
                        label="排序"
                        initialValue={0}
                    >
                        <InputNumber
                            min={0}
                            placeholder="请输入排序值"
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                    <Form.Item
                        name="iconUrl"
                        label="图标URL"
                    >
                        <Input placeholder="请输入图标URL" />
                    </Form.Item>
                    <Form.Item
                        name="description"
                        label="描述"
                    >
                        <Input.TextArea rows={3} placeholder="请输入服务描述" />
                    </Form.Item>
                    <Form.Item
                        name="isActive"
                        label="状态"
                        valuePropName="checked"
                        initialValue={true}
                    >
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default ServiceItemList;
