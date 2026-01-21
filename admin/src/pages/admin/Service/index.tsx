/**
 * 服务项目管理页面
 * 管理陪玩服务项目、定价、状态
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    Switch,
    App,
    Popconfirm,
    Typography,
    
    Row,
    Col,
    Statistic,
} from 'antd';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    AppstoreOutlined,
    SearchOutlined,
    ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { adminApi } from '@/api/admin';

import { logger } from '@/utils/logger';
const { Title, Text } = Typography;
const { TextArea } = Input;

interface ServiceItem {
    id: number;
    itemCode: string;
    name: string;
    description: string;
    category: string;
    subCategory: string;
    gameId?: number;
    basePriceCents: number;
    serviceHours: number;
    commissionRate: number;
    minUsers: number;
    maxPlayers: number;
    tags: string;
    iconUrl: string;
    isActive: boolean;
    sortOrder: number;
    createdAt: string;
    updatedAt: string;
}

interface Game {
    id: number;
    name: string;
}

const subCategoryMap: Record<string, string> = {
    solo: '单人护航',
    team: '团队护航',
    gift: '礼物',
};

const AdminService: React.FC = () => {
    const { message } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [services, setServices] = useState<ServiceItem[]>([]);
    const [games, setGames] = useState<Game[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [modalVisible, setModalVisible] = useState(false);
    const [editingService, setEditingService] = useState<ServiceItem | null>(null);
    const [submitting, setSubmitting] = useState(false);
    const [form] = Form.useForm();
    const [filters, setFilters] = useState({ gameId: undefined as number | undefined, keyword: '', subCategory: undefined as string | undefined });

    const loadGames = useCallback(async () => {
        try {
            const res = await adminApi.getGames({ page_size: 100 });
            if (res.data?.success && res.data?.data) {
                const data = res.data.data;
                const gameList = Array.isArray(data) ? data : [];
                setGames(gameList.map((g: { id: number; name: string }) => ({ id: g.id, name: g.name })));
            } else {
                setGames([]);
            }
        } catch (err) {
            logger.error('加载游戏列表失败:', err);
            setGames([]);
        }
    }, []);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = {
                page: current,
                page_size: pageSize,
            };
            if (filters.gameId) params.game_id = filters.gameId;
            if (filters.subCategory) params.sub_category = filters.subCategory;
            if (filters.keyword) params.keyword = filters.keyword;

            const res = await adminApi.getServiceItems(params);
            if (res.data?.success) {
                const data = res.data.data;
                // 处理不同的响应格式，确保返回数组
                let items: ServiceItem[] = [];
                if (Array.isArray(data)) {
                    items = data;
                } else if (data && typeof data === 'object') {
                    items = Array.isArray((data as { items?: ServiceItem[] }).items) 
                        ? (data as { items: ServiceItem[] }).items 
                        : [];
                }
                setServices(items);
                setTotal(res.data.pagination?.total || items.length || 0);
            } else {
                message.error(res.data?.message || '加载失败');
                setServices([]);
                setTotal(0);
            }
        } catch (err) {
            logger.error('加载服务项目失败:', err);
            message.error('加载数据失败');
            setServices([]);
            setTotal(0);
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, filters, message]);

    useEffect(() => {
        loadGames();
    }, [loadGames]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleAdd = () => {
        setEditingService(null);
        form.resetFields();
        form.setFieldsValue({
            category: 'escort',
            subCategory: 'solo',
            isActive: true,
            sortOrder: 0,
            commissionRate: 0.2,
            minUsers: 1,
            maxPlayers: 1,
            serviceHours: 1,
        });
        setModalVisible(true);
    };

    const handleEdit = (record: ServiceItem) => {
        setEditingService(record);
        form.setFieldsValue({
            ...record,
            basePriceCents: record.basePriceCents / 100,
            tags: record.tags ? JSON.parse(record.tags) : [],
        });
        setModalVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            const res = await adminApi.deleteServiceItem(id);
            if (res.data?.success) {
                message.success('删除成功');
                loadData();
            } else {
                message.error(res.data?.message || '删除失败');
            }
        } catch (err) {
            logger.error('删除失败:', err);
            message.error('删除失败');
        }
    };

    const handleSubmit = async (values: Record<string, unknown>) => {
        setSubmitting(true);
        try {
            const data = {
                itemCode: values.itemCode as string,
                name: values.name as string,
                description: values.description as string,
                subCategory: values.subCategory as 'solo' | 'team' | 'gift',
                gameId: values.gameId as number | undefined,
                basePriceCents: Math.round((values.basePriceCents as number) * 100),
                serviceHours: values.serviceHours as number,
                commissionRate: values.commissionRate as number,
                sortOrder: values.sortOrder as number,
                isActive: values.isActive as boolean,
                tags: JSON.stringify(values.tags || []),
            };

            let res;
            if (editingService) {
                res = await adminApi.updateServiceItem(editingService.id, data);
            } else {
                res = await adminApi.createServiceItem(data);
            }

            if (res.data?.success) {
                message.success(editingService ? '更新成功' : '创建成功');
                setModalVisible(false);
                loadData();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            logger.error('保存失败:', err);
            message.error('操作失败');
        } finally {
            setSubmitting(false);
        }
    };

    const handleStatusChange = async (id: number, isActive: boolean) => {
        try {
            const res = await adminApi.updateServiceItem(id, { isActive });
            if (res.data?.success) {
                message.success(isActive ? '已启用' : '已禁用');
                loadData();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            logger.error('更新状态失败:', err);
            message.error('操作失败');
        }
    };

    const columns: ColumnsType<ServiceItem> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
        {
            title: '服务名称',
            dataIndex: 'name',
            key: 'name',
            render: (name, record) => (
                <Space orientation="vertical" size={0}>
                    <Text strong>{name}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>{record.itemCode}</Text>
                </Space>
            ),
        },
        {
            title: '类型',
            dataIndex: 'subCategory',
            key: 'subCategory',
            width: 100,
            render: (sub) => <Tag color="blue">{subCategoryMap[sub] || sub}</Tag>,
        },
        {
            title: '价格',
            dataIndex: 'basePriceCents',
            key: 'basePriceCents',
            width: 100,
            render: (cents) => <Text strong style={{ color: '#f5222d' }}>¥{(cents / 100).toFixed(2)}</Text>,
        },
        {
            title: '服务时长',
            dataIndex: 'serviceHours',
            key: 'serviceHours',
            width: 90,
            render: (hours) => hours > 0 ? `${hours}小时` : '-',
        },
        {
            title: '抽成比例',
            dataIndex: 'commissionRate',
            key: 'commissionRate',
            width: 90,
            render: (rate) => `${(rate * 100).toFixed(0)}%`,
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: 90,
            render: (isActive, record) => (
                <Switch
                    checked={isActive}
                    onChange={(checked) => handleStatusChange(record.id, checked)}
                    checkedChildren="启用"
                    unCheckedChildren="禁用"
                />
            ),
        },
        { title: '排序', dataIndex: 'sortOrder', key: 'sortOrder', width: 60 },
        {
            title: '操作',
            key: 'action',
            width: 160, // 2个按钮 × 80px
            render: (_, record) => (
                <Space size={4}>
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                    <Popconfirm title="确定删除此服务项目？" onConfirm={() => handleDelete(record.id)} okText="确定" cancelText="取消">
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const stats = {
        total: total,
        active: services.filter((s) => s.isActive).length,
    };

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}>
                <AppstoreOutlined /> 服务项目管理
            </Title>

            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic title="服务项目总数" value={stats.total} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic title="已启用" value={stats.active} />
                    </Card>
                </Col>
            </Row>

            {/* 筛选和操作 */}
            <Card style={{ marginBottom: 16 }}>
                <Row gutter={16} align="middle">
                    <Col flex="auto">
                        <Space wrap>
                            <Select
                                placeholder="选择游戏"
                                allowClear
                                style={{ width: 150 }}
                                value={filters.gameId}
                                onChange={(v) => {
                                    setFilters({ ...filters, gameId: v });
                                    setCurrent(1);
                                }}
                                options={games.map((g) => ({ value: g.id, label: g.name }))}
                            />
                            <Select
                                placeholder="服务类型"
                                allowClear
                                style={{ width: 120 }}
                                value={filters.subCategory}
                                onChange={(v) => {
                                    setFilters({ ...filters, subCategory: v });
                                    setCurrent(1);
                                }}
                                options={[
                                    { value: 'solo', label: '单人护航' },
                                    { value: 'team', label: '团队护航' },
                                    { value: 'gift', label: '礼物' },
                                ]}
                            />
                            <Input
                                placeholder="搜索服务名称"
                                prefix={<SearchOutlined />}
                                style={{ width: 200 }}
                                value={filters.keyword}
                                onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
                                onPressEnter={() => {
                                    setCurrent(1);
                                    loadData();
                                }}
                            />
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <Button icon={<ReloadOutlined />} onClick={loadData}>
                                刷新
                            </Button>
                            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                                新增服务
                            </Button>
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* 服务列表 */}
            <Card>
                <Table
                    columns={columns}
                    dataSource={services}
                    rowKey="id"
                    loading={loading}
                    pagination={{
                        current,
                        pageSize,
                        total,
                        showSizeChanger: true,
                        showTotal: (t) => `共 ${t} 条`,
                        onChange: (page, size) => {
                            setCurrent(page);
                            setPageSize(size);
                        },
                    }}
                />
            </Card>

            {/* 编辑弹窗 */}
            <Modal
                title={editingService ? '编辑服务项目' : '新增服务项目'}
                open={modalVisible}
                onCancel={() => setModalVisible(false)}
                footer={null}
                width={650}
            >
                <Form form={form} layout="vertical" onFinish={handleSubmit}>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="itemCode" label="服务编码" rules={[{ required: true, message: '请输入服务编码' }]}>
                                <Input placeholder="如：ESCORT_SOLO_001" disabled={!!editingService} />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="name" label="服务名称" rules={[{ required: true, message: '请输入服务名称' }]}>
                                <Input placeholder="如：上分陪玩" />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="subCategory" label="服务类型" rules={[{ required: true, message: '请选择服务类型' }]}>
                                <Select
                                    placeholder="选择类型"
                                    options={[
                                        { value: 'solo', label: '单人护航' },
                                        { value: 'team', label: '团队护航' },
                                        { value: 'gift', label: '礼物' },
                                    ]}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="gameId" label="关联游戏">
                                <Select placeholder="选择游戏（可选）" allowClear options={games.map((g) => ({ value: g.id, label: g.name }))} />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item name="basePriceCents" label="基础价格(元)" rules={[{ required: true, message: '请输入价格' }]}>
                                <InputNumber min={0} precision={2} style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="serviceHours" label="服务时长(小时)">
                                <InputNumber min={0} style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="commissionRate" label="抽成比例">
                                <InputNumber min={0} max={1} step={0.01} style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item name="description" label="服务描述">
                        <TextArea rows={3} placeholder="描述服务内容和特点" />
                    </Form.Item>
                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item name="sortOrder" label="排序">
                                <InputNumber min={0} style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="isActive" label="状态" valuePropName="checked">
                                <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item>
                        <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
                            <Button onClick={() => setModalVisible(false)}>取消</Button>
                            <Button type="primary" htmlType="submit" loading={submitting}>
                                保存
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default AdminService;
