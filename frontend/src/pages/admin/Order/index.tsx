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
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { ORDER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import dayjs from 'dayjs';

const { Text, Title } = Typography;

/**
 * 订单数据接口
 */
interface Order {
    id: number;
    orderNo: string;
    userId: number;
    userName: string;
    userAvatar: string;
    playerId: number;
    playerName: string;
    playerAvatar: string;
    gameName: string;
    serviceName: string;
    amount: number;
    duration: number;
    status: 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'canceled' | 'refunded';
    paymentStatus: 'unpaid' | 'paid' | 'refunded';
    remark: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * 订单状态映射
 */
const statusMap = {
    pending: { color: 'gold', text: '待确认', icon: <ClockCircleOutlined /> },
    confirmed: { color: 'blue', text: '已确认', icon: <CheckCircleOutlined /> },
    in_progress: { color: 'processing', text: '进行中', icon: <ClockCircleOutlined /> },
    completed: { color: 'success', text: '已完成', icon: <CheckCircleOutlined /> },
    canceled: { color: 'default', text: '已取消', icon: <CloseCircleOutlined /> },
    refunded: { color: 'error', text: '已退款', icon: <ExclamationCircleOutlined /> },
};

/**
 * 支付状态映射
 */
const paymentStatusMap = {
    unpaid: { color: 'warning', text: '未支付' },
    paid: { color: 'success', text: '已支付' },
    refunded: { color: 'error', text: '已退款' },
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

    // 弹窗状态
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [refundModalVisible, setRefundModalVisible] = useState(false);
    const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
    const [refundForm] = Form.useForm();

    /**
     * 加载订单数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        await new Promise(resolve => setTimeout(resolve, 500));

        const mockOrders: Order[] = Array.from({ length: 30 }, (_, i) => ({
            id: i + 1,
            orderNo: `ORD${dayjs().format('YYYYMMDD')}${String(i + 1).padStart(4, '0')}`,
            userId: 100 + i,
            userName: `用户${i + 1}`,
            userAvatar: '',
            playerId: 200 + i,
            playerName: `陪玩师${(i % 10) + 1}`,
            playerAvatar: '',
            gameName: ['王者荣耀', '英雄联盟', '和平精英', '原神'][i % 4],
            serviceName: ['上分代练', '娱乐陪玩', '技术指导', '组队开黑'][i % 4],
            amount: [50, 80, 100, 150, 200][i % 5],
            duration: [1, 2, 3][i % 3],
            status: ['pending', 'confirmed', 'in_progress', 'completed', 'canceled', 'refunded'][i % 6] as Order['status'],
            paymentStatus: ['unpaid', 'paid', 'paid', 'paid', 'paid', 'refunded'][i % 6] as Order['paymentStatus'],
            remark: i % 3 === 0 ? '请准时上号' : '',
            createdAt: dayjs().subtract(i, 'hour').format('YYYY-MM-DD HH:mm:ss'),
            updatedAt: dayjs().subtract(i, 'hour').add(30, 'minute').format('YYYY-MM-DD HH:mm:ss'),
        }));

        const start = (current - 1) * pageSize;
        const end = start + pageSize;
        setOrders(mockOrders.slice(start, end));
        setTotal(mockOrders.length);
        setLoading(false);
    }, [current, pageSize]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 查看详情
     */
    const handleViewDetail = (order: Order) => {
        setCurrentOrder(order);
        setDetailDrawerVisible(true);
    };

    /**
     * 取消订单
     */
    const handleCancel = async (order: Order) => {
        message.success(`订单 ${order.orderNo} 已取消`);
        loadData();
    };

    /**
     * 打开退款弹窗
     */
    const handleOpenRefund = (order: Order) => {
        setCurrentOrder(order);
        refundForm.setFieldsValue({
            amount: order.amount,
            reason: '',
        });
        setRefundModalVisible(true);
    };

    /**
     * 执行退款
     */
    const handleRefund = async () => {
        try {
            const values = await refundForm.validateFields();
            console.log('Refund:', values);
            message.success('退款成功');
            setRefundModalVisible(false);
            loadData();
        } catch {
            // 验证失败
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'orderNo', label: '订单号', type: 'input', placeholder: '请输入订单号' },
        { name: 'userName', label: '用户', type: 'input', placeholder: '用户名/手机号' },
        { name: 'playerName', label: '陪玩师', type: 'input', placeholder: '陪玩师名称' },
        {
            name: 'status',
            label: '订单状态',
            type: 'select',
            options: Object.entries(statusMap).map(([key, val]) => ({ label: val.text, value: key })),
        },
        {
            name: 'paymentStatus',
            label: '支付状态',
            type: 'select',
            options: Object.entries(paymentStatusMap).map(([key, val]) => ({ label: val.text, value: key })),
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
                    <Avatar size="small" icon={<UserOutlined />} src={record.userAvatar || undefined} />
                    <span>{record.userName}</span>
                </Space>
            ),
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={record.playerAvatar || undefined} style={{ backgroundColor: '#722ed1' }} />
                    <span>{record.playerName}</span>
                </Space>
            ),
        },
        {
            title: '游戏/服务',
            key: 'service',
            width: 180,
            render: (_, record) => (
                <div>
                    <div>{record.gameName}</div>
                    <Text type="secondary" style={{ fontSize: 12 }}>{record.serviceName}</Text>
                </div>
            ),
        },
        {
            title: '金额',
            dataIndex: 'amount',
            key: 'amount',
            width: 100,
            render: amount => <Text strong style={{ color: '#f5222d' }}>¥{amount}</Text>,
        },
        {
            title: '时长',
            dataIndex: 'duration',
            key: 'duration',
            width: 80,
            render: duration => `${duration}小时`,
        },
        {
            title: '订单状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: status => (
                <Tag color={statusMap[status].color} icon={statusMap[status].icon}>
                    {statusMap[status].text}
                </Tag>
            ),
        },
        {
            title: '支付状态',
            dataIndex: 'paymentStatus',
            key: 'paymentStatus',
            width: 100,
            render: status => (
                <Tag color={paymentStatusMap[status].color}>{paymentStatusMap[status].text}</Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
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
                    {record.paymentStatus === 'paid' && !['canceled', 'refunded'].includes(record.status) && (
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

    return (
        <PageContainer title="订单管理" subTitle="管理平台所有订单">
            <SearchTable
                columns={columns}
                dataSource={orders}
                rowKey="id"
                searchFields={searchFields}
                onSearch={() => loadData()}
                onRefresh={loadData}
                loading={loading}
                showCreate={false}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: total => `共 ${total} 条`,
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
                                <Tag color={statusMap[currentOrder.status].color} style={{ fontSize: 16, padding: '4px 16px' }}>
                                    {statusMap[currentOrder.status].icon} {statusMap[currentOrder.status].text}
                                </Tag>
                                <Title level={2} style={{ margin: '16px 0 0' }}>¥{currentOrder.amount}</Title>
                            </div>
                        </Card>

                        {/* 基本信息 */}
                        <Descriptions title="订单信息" column={2} size="small" bordered>
                            <Descriptions.Item label="订单号" span={2}>
                                <Text copyable>{currentOrder.orderNo}</Text>
                            </Descriptions.Item>
                            <Descriptions.Item label="游戏">{currentOrder.gameName}</Descriptions.Item>
                            <Descriptions.Item label="服务">{currentOrder.serviceName}</Descriptions.Item>
                            <Descriptions.Item label="时长">{currentOrder.duration}小时</Descriptions.Item>
                            <Descriptions.Item label="支付状态">
                                <Tag color={paymentStatusMap[currentOrder.paymentStatus].color}>
                                    {paymentStatusMap[currentOrder.paymentStatus].text}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间">{currentOrder.createdAt}</Descriptions.Item>
                            <Descriptions.Item label="更新时间">{currentOrder.updatedAt}</Descriptions.Item>
                            {currentOrder.remark && (
                                <Descriptions.Item label="备注" span={2}>{currentOrder.remark}</Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider />

                        {/* 用户信息 */}
                        <Descriptions title="用户信息" column={2} size="small">
                            <Descriptions.Item label="用户">
                                <Space>
                                    <Avatar size="small" icon={<UserOutlined />} />
                                    {currentOrder.userName}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="用户ID">{currentOrder.userId}</Descriptions.Item>
                        </Descriptions>

                        <Descriptions title="陪玩师信息" column={2} size="small">
                            <Descriptions.Item label="陪玩师">
                                <Space>
                                    <Avatar size="small" icon={<UserOutlined />} style={{ backgroundColor: '#722ed1' }} />
                                    {currentOrder.playerName}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="陪玩师ID">{currentOrder.playerId}</Descriptions.Item>
                        </Descriptions>

                        <Divider />

                        {/* 订单进度 */}
                        <Title level={5}>订单进度</Title>
                        <Timeline
                            items={[
                                { color: 'green', children: `${currentOrder.createdAt} 订单创建` },
                                { color: currentOrder.paymentStatus !== 'unpaid' ? 'green' : 'gray', children: '用户支付' },
                                { color: ['confirmed', 'in_progress', 'completed'].includes(currentOrder.status) ? 'green' : 'gray', children: '陪玩师确认' },
                                { color: ['in_progress', 'completed'].includes(currentOrder.status) ? 'green' : 'gray', children: '开始服务' },
                                { color: currentOrder.status === 'completed' ? 'green' : 'gray', children: '服务完成' },
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
            >
                <Form form={refundForm} layout="vertical">
                    <Form.Item label="订单号">
                        <Text>{currentOrder?.orderNo}</Text>
                    </Form.Item>
                    <Form.Item
                        name="amount"
                        label="退款金额"
                        rules={[
                            { required: true, message: '请输入退款金额' },
                            { type: 'number', max: currentOrder?.amount, message: `退款金额不能超过 ¥${currentOrder?.amount}` },
                        ]}
                    >
                        <InputNumber
                            min={0.01}
                            max={currentOrder?.amount}
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
            </Modal>
        </PageContainer>
    );
};

export default OrderPage;
