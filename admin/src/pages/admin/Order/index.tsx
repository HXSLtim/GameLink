/**
 * 订单管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Popconfirm,
    Drawer,
    Descriptions,
    Timeline,
    Card,
    Typography,
    Divider,
    Avatar,
    Form,
    Input,
    InputNumber,
    Select,
    Radio,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    CloseCircleOutlined,
    DollarOutlined,
    UserOutlined,
    CheckCircleOutlined,
    ClockCircleOutlined,
    ExclamationCircleOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { ORDER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type Order, type ApiResponse } from '@/api/admin';
import { exportToCSV, orderExportColumns } from '@/utils/export';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
const { Text, Title } = Typography;

/**
 * 订单状态映射
 */
const statusMap = {
    pending: { color: 'gold', text: '待确认', icon: <ClockCircleOutlined /> },
    confirmed: { color: 'blue', text: '已确认', icon: <CheckCircleOutlined /> },
    in_progress: { color: 'processing', text: '进行中', icon: <ClockCircleOutlined /> },
    completed: { color: 'success', text: '已完成', icon: <CheckCircleOutlined /> },
    cancelled: { color: 'default', text: '已取消', icon: <CloseCircleOutlined /> },
    refunded: { color: 'error', text: '已退款', icon: <ExclamationCircleOutlined /> },
};

/**
 * 订单管理页面
 */
const OrderPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<Order[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [refundModalVisible, setRefundModalVisible] = useState(false);
    const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
    const [refundForm] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    // 批量操作状态
    const [batchCancelVisible, setBatchCancelVisible] = useState(false);
    const [batchCompleteVisible, setBatchCompleteVisible] = useState(false);
    const [selectedOrderIds, setSelectedOrderIds] = useState<number[]>([]);
    const [batchTarget, setBatchTarget] = useState<'selected' | 'status' | 'all'>('selected');
    const [batchForm] = Form.useForm();

    /**
     * 加载订单数据
     */
    const loadData = useCallback(async (params?: Record<string, unknown>) => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
                ...params,
            };
            const response = await adminApi.getOrders(queryParams);
            if (response.data.success) {
                setOrders(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load orders error:', error);
            message.error('加载订单列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    /**
     * 查看详情
     */
    const handleViewDetail = async (order: Order) => {
        try {
            const response = await adminApi.getOrder(order.id);
            if (response.data.success) {
                setCurrentOrder(response.data.data);
            } else {
                setCurrentOrder(order);
            }
        } catch {
            setCurrentOrder(order);
        }
        setDetailDrawerVisible(true);
    };

    /**
     * 取消订单
     */
    const handleCancel = async (order: Order) => {
        try {
            await adminApi.cancelOrder(order.id);
            message.success(`订单 ${order.orderNo} 已取消`);
            loadData();
        } catch (error) {
            logger.error('Cancel order error:', error);
            message.error('取消订单失败');
        }
    };

    /**
     * 打开退款弹窗
     */
    const handleOpenRefund = (order: Order) => {
        setCurrentOrder(order);
        refundForm.setFieldsValue({
            amount: order.totalPriceCents / 100,
            reason: '',
        });
        setRefundModalVisible(true);
    };

    /**
     * 执行退款
     */
    const handleRefund = async () => {
        if (!currentOrder) return;
        try {
            const values = await refundForm.validateFields();
            setSubmitting(true);
            await adminApi.refundOrder(currentOrder.id, {
                reason: values.reason,
                amount_cents: Math.round(values.amount * 100),
            });
            message.success('退款成功');
            setRefundModalVisible(false);
            loadData();
        } catch (error) {
            logger.error('Refund error:', error);
            message.error('退款失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * 批量取消订单
     */
    const handleBatchCancel = (keys: React.Key[]) => {
        setSelectedOrderIds(keys ? keys.map(k => Number(k)) : []);
        batchForm.resetFields();
        batchForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'selected' : 'all',
        });
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchCancelVisible(true);
    };

    const submitBatchCancel = async () => {
        try {
            const values = await batchForm.validateFields();
            let orderIds: number[] = [];

            if (values.target === 'selected') {
                orderIds = selectedOrderIds;
            } else if (values.target === 'status') {
                const response = await adminApi.getOrders({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    orderIds = response.data.data.map((o: Order) => o.id);
                }
            } else {
                // 全部可取消的订单（pending, confirmed）
                const response = await adminApi.getOrders({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    orderIds = response.data.data
                        .filter((o: Order) => ['pending', 'confirmed'].includes(o.status))
                        .map((o: Order) => o.id);
                }
            }

            if (orderIds.length === 0) {
                message.warning('没有符合条件的订单');
                return;
            }

            const res = await adminApi.batchCancelOrders(orderIds, values.reason) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success(`批量取消 ${orderIds.length} 个订单成功`);
                setBatchCancelVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * 批量完成订单
     */
    const handleBatchComplete = (keys: React.Key[]) => {
        setSelectedOrderIds(keys ? keys.map(k => Number(k)) : []);
        batchForm.resetFields();
        batchForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'selected' : 'all',
        });
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchCompleteVisible(true);
    };

    const submitBatchComplete = async () => {
        try {
            const values = await batchForm.validateFields();
            let orderIds: number[] = [];

            if (values.target === 'selected') {
                orderIds = selectedOrderIds;
            } else if (values.target === 'status') {
                const response = await adminApi.getOrders({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    orderIds = response.data.data.map((o: Order) => o.id);
                }
            } else {
                // 全部可完成的订单（in_progress）
                const response = await adminApi.getOrders({ status: 'in_progress', page_size: 1000 });
                if (response.data.success && response.data.data) {
                    orderIds = response.data.data.map((o: Order) => o.id);
                }
            }

            if (orderIds.length === 0) {
                message.warning('没有符合条件的订单');
                return;
            }

            const res = await adminApi.batchCompleteOrders(orderIds) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success(`批量完成 ${orderIds.length} 个订单成功`);
                setBatchCompleteVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'orderNo', label: '订单号', type: 'input', placeholder: '请输入订单号' },
        {
            name: 'status',
            label: '订单状态',
            type: 'select',
            options: Object.entries(statusMap).map(([key, val]) => ({ label: val.text, value: key })),
        },
        { name: 'dateRange', label: '创建时间', type: 'dateRange' },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Order> = [
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            width: 180,
            render: text => <Text copyable={{ text }}>{text}</Text>,
        },
        {
            title: '用户',
            key: 'user',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={record.user?.avatarUrl || undefined} />
                    <span>{record.user?.name || `用户${record.userId}`}</span>
                </Space>
            ),
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size="small"
                        icon={<UserOutlined />}
                        src={record.player?.user?.avatarUrl || undefined}
                        style={{ backgroundColor: '#722ed1' }}
                    />
                    <span>{record.player?.nickname || (record.playerId ? `陪玩师${record.playerId}` : '-')}</span>
                </Space>
            ),
        },
        {
            title: '游戏',
            key: 'game',
            width: 120,
            render: (_, record) => record.game?.name || '-',
        },
        {
            title: '标题',
            dataIndex: 'title',
            key: 'title',
            width: 150,
            ellipsis: true,
            render: title => title || '-',
        },
        {
            title: '金额',
            dataIndex: 'totalPriceCents',
            key: 'totalPriceCents',
            width: 100,
            render: cents => <Text strong style={{ color: '#f5222d' }}>¥{(cents / 100).toFixed(2)}</Text>,
        },
        {
            title: '订单状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: Order['status']) => (
                <Tag color={statusMap[status]?.color || 'default'} icon={statusMap[status]?.icon}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                    {['pending', 'confirmed'].includes(record.status) && (
                        <PermissionGuard permission={ORDER_PERMISSIONS.CANCEL}>
                            <Popconfirm
                                title="确定要取消该订单吗？"
                                onConfirm={() => handleCancel(record)}
                            >
                                <Button type="link" size="small" danger icon={<CloseCircleOutlined />}>
                                    取消
                                </Button>
                            </Popconfirm>
                        </PermissionGuard>
                    )}
                    {!['cancelled', 'refunded'].includes(record.status) && record.totalPriceCents > 0 && (
                        <PermissionGuard permission={ORDER_PERMISSIONS.REFUND}>
                            <Button
                                type="link"
                                size="small"
                                icon={<DollarOutlined />}
                                onClick={() => handleOpenRefund(record)}
                            >
                                退款
                            </Button>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    /**
     * 导出订单数据
     */
    const handleExport = async () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            const response = await adminApi.getOrders({ ...searchParams, page_size: 10000 });
            if (response.data.success && response.data.data) {
                exportToCSV(response.data.data as unknown as Record<string, unknown>[], orderExportColumns, 'orders');
                message.success({ content: '导出成功', key: 'export' });
            } else {
                message.error({ content: '导出失败', key: 'export' });
            }
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    };

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量取消',
            icon: <CloseCircleOutlined />,
            needSelection: false,
            danger: true,
            onClick: (keys) => handleBatchCancel(keys || []),
            permission: ORDER_PERMISSIONS.CANCEL,
        },
        {
            text: '批量完成',
            icon: <CheckCircleOutlined />,
            needSelection: false,
            onClick: (keys) => handleBatchComplete(keys || []),
            permission: ORDER_PERMISSIONS.UPDATE,
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
            permission: ORDER_PERMISSIONS.LIST,
        },
    ];

    return (
        <PageContainer title="订单管理" subTitle="管理平台所有订单">
            <SearchTable
                columns={columns}
                dataSource={orders}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={false}
                toolbarButtons={toolbarButtons}
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
                scroll={{ x: 1500 }}
            />

            {/* 订单详情抽屉 */}
            <Drawer
                title="订单详情"
                open={detailDrawerVisible}
                onClose={() => setDetailDrawerVisible(false)}
                width={600}
            >
                {currentOrder && (
                    <>
                        {/* 状态卡片 */}
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <div style={{ textAlign: 'center' }}>
                                <Tag color={statusMap[currentOrder.status]?.color} style={{ fontSize: 16, padding: '4px 16px' }}>
                                    {statusMap[currentOrder.status]?.icon} {statusMap[currentOrder.status]?.text}
                                </Tag>
                                <Title level={2} style={{ margin: '16px 0 0' }}>
                                    ¥{(currentOrder.totalPriceCents / 100).toFixed(2)}
                                </Title>
                            </div>
                        </Card>

                        {/* 基本信息 */}
                        <Descriptions title="订单信息" column={2} size="small" bordered>
                            <Descriptions.Item label="订单号" span={2}>
                                <Text copyable>{currentOrder.orderNo}</Text>
                            </Descriptions.Item>
                            <Descriptions.Item label="游戏">{currentOrder.game?.name || '-'}</Descriptions.Item>
                            <Descriptions.Item label="标题">{currentOrder.title || '-'}</Descriptions.Item>
                            <Descriptions.Item label="金额">
                                ¥{(currentOrder.totalPriceCents / 100).toFixed(2)} {currentOrder.currency}
                            </Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={statusMap[currentOrder.status]?.color}>
                                    {statusMap[currentOrder.status]?.text}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="预约开始">
                                {currentOrder.scheduledStart ? dayjs(currentOrder.scheduledStart).format('YYYY-MM-DD HH:mm') : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="预约结束">
                                {currentOrder.scheduledEnd ? dayjs(currentOrder.scheduledEnd).format('YYYY-MM-DD HH:mm') : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间">
                                {currentOrder.createdAt ? dayjs(currentOrder.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="完成时间">
                                {currentOrder.completedAt ? dayjs(currentOrder.completedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                            </Descriptions.Item>
                            {currentOrder.description && (
                                <Descriptions.Item label="描述" span={2}>{currentOrder.description}</Descriptions.Item>
                            )}
                            {currentOrder.cancelReason && (
                                <Descriptions.Item label="取消原因" span={2}>
                                    <Text type="danger">{currentOrder.cancelReason}</Text>
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider />

                        {/* 用户信息 */}
                        <Descriptions title="用户信息" column={2} size="small">
                            <Descriptions.Item label="用户">
                                <Space>
                                    <Avatar size="small" icon={<UserOutlined />} src={currentOrder.user?.avatarUrl} />
                                    {currentOrder.user?.name || `用户${currentOrder.userId}`}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="用户ID">{currentOrder.userId}</Descriptions.Item>
                        </Descriptions>

                        <Descriptions title="陪玩师信息" column={2} size="small">
                            <Descriptions.Item label="陪玩师">
                                <Space>
                                    <Avatar
                                        size="small"
                                        icon={<UserOutlined />}
                                        src={currentOrder.player?.user?.avatarUrl}
                                        style={{ backgroundColor: '#722ed1' }}
                                    />
                                    {currentOrder.player?.nickname || (currentOrder.playerId ? `陪玩师${currentOrder.playerId}` : '未分配')}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="陪玩师ID">{currentOrder.playerId || '-'}</Descriptions.Item>
                        </Descriptions>

                        <Divider />

                        {/* 订单进度 */}
                        <Title level={5}>订单进度</Title>
                        <Timeline
                            items={[
                                {
                                    color: 'green',
                                    children: `${currentOrder.createdAt ? dayjs(currentOrder.createdAt).format('YYYY-MM-DD HH:mm:ss') : ''} 订单创建`,
                                },
                                {
                                    color: ['confirmed', 'in_progress', 'completed'].includes(currentOrder.status) ? 'green' : 'gray',
                                    children: '订单确认',
                                },
                                {
                                    color: ['in_progress', 'completed'].includes(currentOrder.status) ? 'green' : 'gray',
                                    children: '开始服务',
                                },
                                {
                                    color: currentOrder.status === 'completed' ? 'green' : 'gray',
                                    children: currentOrder.completedAt
                                        ? `${dayjs(currentOrder.completedAt).format('YYYY-MM-DD HH:mm:ss')} 服务完成`
                                        : '服务完成',
                                },
                            ]}
                        />
                    </>
                )}
            </Drawer>

            {/* 退款弹窗 */}
            <Modal
                title="订单退款"
                open={refundModalVisible}
                onOk={handleRefund}
                onCancel={() => setRefundModalVisible(false)}
                confirmLoading={submitting}
                width={550}
            >
                {currentOrder && (
                    <>
                        {/* 订单信息概览 */}
                        <Card size="small" style={{ marginBottom: 16, backgroundColor: '#fafafa' }}>
                            <Descriptions column={2} size="small">
                                <Descriptions.Item label="订单号" span={2}>
                                    <Text copyable>{currentOrder.orderNo}</Text>
                                </Descriptions.Item>
                                <Descriptions.Item label="用户">
                                    <Space>
                                        <Avatar size="small" icon={<UserOutlined />} src={currentOrder.user?.avatarUrl} />
                                        {currentOrder.user?.name || `用户${currentOrder.userId}`}
                                    </Space>
                                </Descriptions.Item>
                                <Descriptions.Item label="陪玩师">
                                    <Space>
                                        <Avatar
                                            size="small"
                                            icon={<UserOutlined />}
                                            src={currentOrder.player?.user?.avatarUrl}
                                            style={{ backgroundColor: '#722ed1' }}
                                        />
                                        {currentOrder.player?.nickname || '-'}
                                    </Space>
                                </Descriptions.Item>
                                <Descriptions.Item label="游戏">{currentOrder.game?.name || '-'}</Descriptions.Item>
                                <Descriptions.Item label="订单金额">
                                    <Text strong style={{ color: '#f5222d' }}>
                                        ¥{(currentOrder.totalPriceCents / 100).toFixed(2)}
                                    </Text>
                                </Descriptions.Item>
                                <Descriptions.Item label="订单状态">
                                    <Tag color={statusMap[currentOrder.status]?.color}>
                                        {statusMap[currentOrder.status]?.text}
                                    </Tag>
                                </Descriptions.Item>
                            </Descriptions>
                        </Card>

                        <Divider style={{ margin: '12px 0' }} />

                        <Form form={refundForm} layout="vertical">
                            <Form.Item
                                name="amount"
                                label="退款金额"
                                rules={[
                                    { required: true, message: '请输入退款金额' },
                                    {
                                        type: 'number',
                                        max: currentOrder.totalPriceCents / 100,
                                        message: `退款金额不能超过 ¥${(currentOrder.totalPriceCents / 100).toFixed(2)}`,
                                    },
                                ]}
                                extra={`最大可退款金额: ¥${(currentOrder.totalPriceCents / 100).toFixed(2)}`}
                            >
                                <InputNumber
                                    min={0.01}
                                    max={currentOrder.totalPriceCents / 100}
                                    precision={2}
                                    prefix="¥"
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>
                            <Form.Item
                                name="reason"
                                label="退款原因"
                                rules={[{ required: true, message: '请输入退款原因' }]}
                            >
                                <Input.TextArea rows={3} placeholder="请输入退款原因" />
                            </Form.Item>
                        </Form>
                    </>
                )}
            </Modal>

            {/* 批量取消订单弹窗 */}
            <Modal
                title="批量取消订单"
                open={batchCancelVisible}
                onOk={submitBatchCancel}
                onCancel={() => setBatchCancelVisible(false)}
                okText="确认取消"
                okButtonProps={{ danger: true }}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedOrderIds.length === 0}>
                                选中的订单 {selectedOrderIds.length > 0 ? `(${selectedOrderIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部可取消订单</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true, message: '请选择筛选状态' }]}>
                            <Select placeholder="请选择要筛选的状态">
                                <Select.Option value="pending">待确认</Select.Option>
                                <Select.Option value="confirmed">已确认</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <Form.Item name="reason" label="取消原因">
                        <Input.TextArea rows={3} placeholder="请输入取消原因（选填）" />
                    </Form.Item>

                    <div style={{ color: '#ff4d4f', marginTop: 16 }}>
                        ⚠️ 警告：此操作将取消所有符合条件的订单，请谨慎操作！
                    </div>
                </Form>
            </Modal>

            {/* 批量完成订单弹窗 */}
            <Modal
                title="批量完成订单"
                open={batchCompleteVisible}
                onOk={submitBatchComplete}
                onCancel={() => setBatchCompleteVisible(false)}
                okText="确认完成"
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedOrderIds.length === 0}>
                                选中的订单 {selectedOrderIds.length > 0 ? `(${selectedOrderIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部进行中订单</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true, message: '请选择筛选状态' }]}>
                            <Select placeholder="请选择要筛选的状态">
                                <Select.Option value="in_progress">进行中</Select.Option>
                                <Select.Option value="confirmed">已确认</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <div style={{ color: '#1890ff', marginTop: 16 }}>
                        ℹ️ 提示：此操作将完成所有符合条件的订单
                    </div>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default OrderPage;
