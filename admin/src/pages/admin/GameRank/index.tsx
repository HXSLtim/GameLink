/**
 * 段位管理页面
 * 管理游戏段位配置（如：青铜、白银、黄金等）
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Popconfirm,
    Form,
    Input,
    InputNumber,
    Switch,
    Select,
    ColorPicker,
    Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { adminApi, type GameRank, type Game } from '@/api/admin';
import dayjs from 'dayjs';
import { GAMELINK_PRIMARY } from '@/theme';
import { logger } from '@/utils/logger';

const { Text } = Typography;

const GameRankPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [ranks, setRanks] = useState<GameRank[]>([]);
    const [games, setGames] = useState<Game[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [modalVisible, setModalVisible] = useState(false);
    const [editingRank, setEditingRank] = useState<GameRank | null>(null);
    const [form] = Form.useForm();
    const [submitLoading, setSubmitLoading] = useState(false);

    // 批量操作
    const [selectedKeys, setSelectedKeys] = useState<React.Key[]>([]);

    // 加载游戏列表
    const loadGames = useCallback(async () => {
        try {
            const response = await adminApi.getGames({ page_size: 1000 });
            if (response.data.success) {
                setGames(response.data.data || []);
            }
        } catch (error) {
            logger.error('Load games error:', error);
        }
    }, []);

    // 加载段位数据
    const loadData = useCallback(async (params?: Record<string, unknown>) => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                pageSize,
                ...searchParams,
                ...params,
            };
            const response = await adminApi.getGameRanks(queryParams);
            if (response.data.success) {
                setRanks(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load game ranks error:', error);
            message.error('加载段位列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadGames();
    }, [loadGames]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleCreate = () => {
        setEditingRank(null);
        form.resetFields();
        form.setFieldsValue({ isActive: true, sortOrder: 0, level: 0, priceCents: 0 });
        setModalVisible(true);
    };

    const handleEdit = (record: GameRank) => {
        setEditingRank(record);
        form.setFieldsValue({
            ...record,
            color: record.color || GAMELINK_PRIMARY.base,
        });
        setModalVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await adminApi.deleteGameRank(id);
            message.success('删除成功');
            loadData();
        } catch (error) {
            logger.error('Delete game rank error:', error);
            message.error('删除失败');
        }
    };

    const handleSubmit = async () => {
        try {
            setSubmitLoading(true);
            const values = await form.validateFields();
            const data = {
                ...values,
                color: typeof values.color === 'string'
                    ? values.color
                    : values.color?.toHexString?.() || GAMELINK_PRIMARY.base,
            };

            if (editingRank) {
                await adminApi.updateGameRank(editingRank.id, data);
                message.success('更新成功');
            } else {
                await adminApi.createGameRank(data);
                message.success('创建成功');
            }
            setModalVisible(false);
            loadData();
        } catch (error) {
            logger.error('Submit game rank error:', error);
            message.error('操作失败');
        } finally {
            setSubmitLoading(false);
        }
    };

    const handleBatchDelete = async () => {
        if (selectedKeys.length === 0) {
            message.warning('请先选择要删除的段位');
            return;
        }
        try {
            await adminApi.batchDeleteGameRanks(selectedKeys.map(String));
            message.success(`批量删除 ${selectedKeys.length} 个段位成功`);
            setSelectedKeys([]);
            loadData();
        } catch (error) {
            logger.error('Batch delete error:', error);
            message.error('批量删除失败');
        }
    };

    const handleBatchStatus = async (isActive: boolean) => {
        if (selectedKeys.length === 0) {
            message.warning('请先选择段位');
            return;
        }
        try {
            await adminApi.batchUpdateGameRankStatus(selectedKeys.map(String), isActive);
            message.success(`批量${isActive ? '启用' : '禁用'} ${selectedKeys.length} 个段位成功`);
            setSelectedKeys([]);
            loadData();
        } catch (error) {
            logger.error('Batch status error:', error);
            message.error('批量操作失败');
        }
    };

    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '段位名称' },
        {
            name: 'gameId',
            label: '游戏',
            type: 'select',
            options: games.map(g => ({ label: g.name, value: g.id })),
        },
        {
            name: 'isActive',
            label: '状态',
            type: 'select',
            options: [
                { label: '启用', value: 'true' },
                { label: '禁用', value: 'false' },
            ],
        },
    ];

    const columns: ColumnsType<GameRank> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '段位名称',
            key: 'name',
            width: 180,
            render: (_, record) => (
                <Space>
                    {record.color && (
                        <span style={{
                            display: 'inline-block',
                            width: 14,
                            height: 14,
                            borderRadius: '50%',
                            backgroundColor: record.color,
                            border: '1px solid #d9d9d9',
                        }} />
                    )}
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.name}</div>
                        {record.description && (
                            <Text type="secondary" style={{ fontSize: 12 }}>
                                {record.description}
                            </Text>
                        )}
                    </div>
                </Space>
            ),
        },
        {
            title: '所属游戏',
            key: 'game',
            width: 120,
            render: (_, record) => (
                <Tag color="blue">{record.game?.name || '-'}</Tag>
            ),
        },
        {
            title: '等级',
            dataIndex: 'level',
            key: 'level',
            width: 80,
            render: level => <Tag>{level}</Tag>,
        },
        {
            title: '定价',
            dataIndex: 'priceCents',
            key: 'priceCents',
            width: 120,
            render: cents => cents ? (
                <Text type="danger" strong>¥{(cents / 100).toFixed(2)}/小时</Text>
            ) : '-',
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 80,
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: 80,
            render: isActive => (
                <Tag color={isActive ? 'success' : 'default'}>
                    {isActive ? '启用' : '禁用'}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 160,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 160, // 2个按钮 × 80px
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                    <Popconfirm title="确定删除该段位？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量启用',
            needSelection: true,
            onClick: () => handleBatchStatus(true),
        },
        {
            text: '批量禁用',
            needSelection: true,
            onClick: () => handleBatchStatus(false),
        },
        {
            text: '批量删除',
            icon: <DeleteOutlined />,
            needSelection: true,
            danger: true,
            onClick: handleBatchDelete,
        },
    ];

    return (
        <PageContainer title="段位管理" subTitle="管理游戏段位配置，设置不同段位的定价">
            <SearchTable
                columns={columns}
                dataSource={ranks}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate
                createText="新增段位"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                rowSelection={{
                    selectedRowKeys: selectedKeys,
                    onChange: setSelectedKeys,
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
                scroll={{ x: 1100 }}
            />

            <Modal
                title={editingRank ? '编辑段位' : '新增段位'}
                open={modalVisible}
                onOk={handleSubmit}
                onCancel={() => setModalVisible(false)}
                confirmLoading={submitLoading}
                width={560}
                destroyOnHidden
            >
                <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
                    <Form.Item name="gameId" label="所属游戏" rules={[{ required: true, message: '请选择游戏' }]}>
                        <Select placeholder="请选择游戏" disabled={!!editingRank} showSearch optionFilterProp="children">
                            {games.map(g => (
                                <Select.Option key={g.id} value={g.id}>{g.name}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item name="name" label="段位名称" rules={[{ required: true, message: '请输入段位名称' }]}>
                        <Input placeholder="如：青铜、白银、黄金" maxLength={64} />
                    </Form.Item>
                    <Space style={{ width: '100%' }} size={16}>
                        <Form.Item name="level" label="等级" tooltip="数字越大等级越高" style={{ width: 160 }}>
                            <InputNumber min={0} style={{ width: '100%' }} />
                        </Form.Item>
                        <Form.Item name="sortOrder" label="排序" tooltip="数字越小越靠前" style={{ width: 160 }}>
                            <InputNumber min={0} style={{ width: '100%' }} />
                        </Form.Item>
                        <Form.Item name="color" label="段位颜色" style={{ width: 120 }}>
                            <ColorPicker />
                        </Form.Item>
                    </Space>
                    <Form.Item name="priceCents" label="定价（分）" tooltip="该段位的小时定价，单位为分">
                        <InputNumber min={0} style={{ width: '100%' }} addonAfter="分/小时" placeholder="如：2000 表示 20元/小时" />
                    </Form.Item>
                    <Form.Item name="iconUrl" label="图标URL">
                        <Input placeholder="段位图标地址（可选）" />
                    </Form.Item>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={2} placeholder="段位描述（可选）" />
                    </Form.Item>
                    <Form.Item name="isActive" label="启用状态" valuePropName="checked">
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default GameRankPage;
