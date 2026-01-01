/**
 * 陪玩师端订单列表页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Tag,
    Button,
    Space,
    Avatar,
    Typography,
    Tabs,
    Empty,
    App,
} from 'antd';
import {
    UserOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    PlayCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Text, Title } = Typography;

interface Order {
    id: number;
    orderNo: string;
    user: {
        id: number;
        nickname: string;
        avatar: string;
    };
    game: string;
    service: string;
    price: number;
    duration: number;
    totalAmount: number;
    commission: number;
    status: string;
    createdAt: string;
    startTime?: string;
    endTime?: string;
}

const statusConfig: Record<string, { color: string; text: string }> = {
    paid: { color: 'gold', text: '待接单' },
    accepted: { color: 'cyan', text: '已接单' },
    in_progress: { color: 'processing', text: '进行中' },
    completed: { color: 'green', text: '已完成' },
    cancelled: { color: 'default', text: '已取消' },
    refunded: { color: 'red', text: '已退款' },
};

const PlayerOrders: React.FC = () => {
    const { message, modal } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<Order[]>([]);
    const [activeTab, setActiveTab] = useState('all');

    const loadOrders = useCallback(async () => {
        setLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 500));
            
            const mockOrders: Order[] = [
                {
                    id: 1,
                    orderNo: 'ORD202412160001',
                    user: { id: 1, nickname: '快乐玩家', avatar: '' },
                    game: '王者荣耀',
                    service: '上分陪玩',
                    price: 50,
                    duration: 2,
                    totalAmount: 100,
                    commission: 85,
                    status: 'in_progress',
                    createdAt: '2024-12-16 10:30:00',
                    startTime: '2024-12-16 11:00:00',
                },
                {
                    id: 2,
                    orderNo: 'ORD202412150002',
                    user: { id: 2, nickname: '游戏达人', avatar: '' },
                    game: '英雄联盟',
                    service: '娱乐陪玩',
                    price: 80,
                    duration: 1,
                    totalAmount: 80,
                    commission: 68,
                    status: 'completed',
                    createdAt: '2024-12-15 14:00:00',
                    startTime: '2024-12-15 15:00:00',
                    endTime: '2024-12-15 16:00:00',
                },
                {
                    id: 3,
                    orderNo: 'ORD202412160003',
                    user: { id: 3, nickname: '新手小白', avatar: '' },
                    game: '和平精英',
                    service: '陪聊',
                    price: 45,
                    duration: 1,
                    totalAmount: 45,
                    commission: 38,
                    status: 'paid',
                    createdAt: '2024-12-16 15:00:00',
                },
            ];
            setOrders(mockOrders);
        } catch {
            message.error('加载订单失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadOrders();
    }, [loadOrders]);

    const handleAccept = (order: Order) => {
        modal.confirm({
            title: '确认接单',
            content: `确定接受订单 ${order.orderNo} 吗？`,
            onOk: async () => {
                message.success('接单成功');
                loadOrders();
            },
        });
    };

    const handleStart = (_order: Order) => {
        modal.confirm({
            title: '开始服务',
            content: '确定开始服务吗？开始后将计时。',
            onOk: async () => {
                message.success('服务已开始');
                loadOrders();
            },
        });
    };

    const handleComplete = (_order: Order) => {
        modal.confirm({
            title: '完成服务',
            content: '确定完成服务吗？',
            onOk: async () => {
                message.success('服务已完成');
                loadOrders();
            },
        });
    };

    const columns: ColumnsType<Order> = [
        {
            title: '订单信息',
            key: 'info',
            render: (_, record) => (
                <Space>
                    <Avatar src={record.user.avatar} icon={<UserOutlined />} />
                    <div>
                        <div><Text strong>{record.user.nickname}</Text></div>
                        <Text type="secondary" style={{ fontSize: 12 }}>{record.orderNo}</Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '服务内容',
            key: 'service',
            render: (_, record) => (
                <div>
                    <div><Tag color="blue">{record.game}</Tag></div>
                    <Text>{record.service}</Text>
                </div>
            ),
        },
        {
            title: '时长',
            dataIndex: 'duration',
            key: 'duration',
            render: (duration: number) => `${duration}小时`,
        },
        {
            title: '订单金额',
            dataIndex: 'totalAmount',
            key: 'totalAmount',
            render: (amount: number) => <Text>¥{amount}</Text>,
        },
        {
            title: '我的收益',
            dataIndex: 'commission',
            key: 'commission',
            render: (commission: number) => (
                <Text type="success" strong>¥{commission}</Text>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => {
                const config = statusConfig[status] || { color: 'default', text: status };
                return <Tag color={config.color}>{config.text}</Tag>;
            },
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
        },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => {
                const actions = [];
                if (record.status === 'paid') {
                    actions.push(
                        <Button key="accept" type="primary" size="small" onClick={() => handleAccept(record)}>
                            接单
                        </Button>
                    );
                }
                if (record.status === 'accepted') {
                    actions.push(
                        <Button key="start" type="primary" size="small" onClick={() => handleStart(record)}>
                            开始服务
                        </Button>
                    );
                }
                if (record.status === 'in_progress') {
                    actions.push(
                        <Button key="complete" type="primary" size="small" onClick={() => handleComplete(record)}>
                            完成服务
                        </Button>
                    );
                }
                actions.push(
                    <Button key="detail" type="link" size="small">
                        详情
                    </Button>
                );
                return <Space>{actions}</Space>;
            },
        },
    ];

    const filteredOrders = activeTab === 'all' 
        ? orders 
        : orders.filter(o => o.status === activeTab);

    const tabItems = [
        { key: 'all', label: '全部', icon: <ClockCircleOutlined /> },
        { key: 'paid', label: '待接单', icon: <ClockCircleOutlined /> },
        { key: 'in_progress', label: '进行中', icon: <PlayCircleOutlined /> },
        { key: 'completed', label: '已完成', icon: <CheckCircleOutlined /> },
        { key: 'cancelled', label: '已取消', icon: <CloseCircleOutlined /> },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}>我的订单</Title>
            
            <Card>
                <Tabs
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    items={tabItems.map(item => ({
                        key: item.key,
                        label: (
                            <span>
                                {item.icon} {item.label}
                            </span>
                        ),
                    }))}
                />
                
                <Table
                    columns={columns}
                    dataSource={filteredOrders}
                    rowKey="id"
                    loading={loading}
                    locale={{ emptyText: <Empty description="暂无订单" /> }}
                    pagination={{ pageSize: 10 }}
                    scroll={{ x: 'max-content' }}
                />
            </Card>
        </div>
    );
};

export default PlayerOrders;
