/**
 * 用户端订单列表页面
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
    Modal,
    Rate,
    Input,
    App,
} from 'antd';
import {
    UserOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    ExclamationCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Text, Title } = Typography;
const { TextArea } = Input;

interface Order {
    id: number;
    orderNo: string;
    player: {
        id: number;
        nickname: string;
        avatar: string;
    };
    game: string;
    service: string;
    price: number;
    duration: number;
    totalAmount: number;
    status: string;
    createdAt: string;
    startTime?: string;
    endTime?: string;
}

const statusConfig: Record<string, { color: string; text: string }> = {
    pending: { color: 'gold', text: '待支付' },
    paid: { color: 'blue', text: '待接单' },
    accepted: { color: 'cyan', text: '已接单' },
    in_progress: { color: 'processing', text: '进行中' },
    completed: { color: 'green', text: '已完成' },
    canceled: { color: 'default', text: '已取消' },
    refunded: { color: 'red', text: '已退款' },
};

const UserOrders: React.FC = () => {
    const { message, modal } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<Order[]>([]);
    const [activeTab, setActiveTab] = useState('all');
    const [reviewModal, setReviewModal] = useState(false);
    const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
    const [rating, setRating] = useState(5);
    const [reviewContent, setReviewContent] = useState('');

    const loadOrders = useCallback(async () => {
        setLoading(true);
        try {
            // TODO: 替换为真实 API
            await new Promise(resolve => setTimeout(resolve, 500));
            
            const mockOrders: Order[] = [
                {
                    id: 1,
                    orderNo: 'ORD202412160001',
                    player: { id: 1, nickname: '小甜甜', avatar: '' },
                    game: '王者荣耀',
                    service: '上分陪玩',
                    price: 50,
                    duration: 2,
                    totalAmount: 100,
                    status: 'in_progress',
                    createdAt: '2024-12-16 10:30:00',
                    startTime: '2024-12-16 11:00:00',
                },
                {
                    id: 2,
                    orderNo: 'ORD202412150002',
                    player: { id: 2, nickname: '大神带飞', avatar: '' },
                    game: '英雄联盟',
                    service: '娱乐陪玩',
                    price: 80,
                    duration: 1,
                    totalAmount: 80,
                    status: 'completed',
                    createdAt: '2024-12-15 14:00:00',
                    startTime: '2024-12-15 15:00:00',
                    endTime: '2024-12-15 16:00:00',
                },
                {
                    id: 3,
                    orderNo: 'ORD202412140003',
                    player: { id: 3, nickname: '温柔学姐', avatar: '' },
                    game: '和平精英',
                    service: '陪聊',
                    price: 45,
                    duration: 1,
                    totalAmount: 45,
                    status: 'pending',
                    createdAt: '2024-12-14 20:00:00',
                },
            ];
            setOrders(mockOrders);
        } catch {
            message.error('加载订单失败');
        } finally {
            setLoading(false);
        }
    }, [message]);

    useEffect(() => {
        loadOrders();
    }, [loadOrders]);

    const handleCancel = (order: Order) => {
        modal.confirm({
            title: '确认取消订单',
            content: `确定要取消订单 ${order.orderNo} 吗？`,
            onOk: async () => {
                message.success('订单已取消');
                loadOrders();
            },
        });
    };

    const handlePay = (order: Order) => {
        modal.confirm({
            title: '确认支付',
            content: `订单金额：¥${order.totalAmount}`,
            onOk: async () => {
                message.success('支付成功');
                loadOrders();
            },
        });
    };

    const handleReview = (order: Order) => {
        setSelectedOrder(order);
        setRating(5);
        setReviewContent('');
        setReviewModal(true);
    };

    const submitReview = async () => {
        if (!reviewContent.trim()) {
            message.warning('请输入评价内容');
            return;
        }
        message.success('评价提交成功');
        setReviewModal(false);
        loadOrders();
    };

    const columns: ColumnsType<Order> = [
        {
            title: '订单信息',
            key: 'info',
            render: (_, record) => (
                <Space>
                    <Avatar src={record.player.avatar} icon={<UserOutlined />} />
                    <div>
                        <div><Text strong>{record.player.nickname}</Text></div>
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
            title: '时长/金额',
            key: 'amount',
            render: (_, record) => (
                <div>
                    <div><Text>{record.duration}小时</Text></div>
                    <Text type="danger" strong>¥{record.totalAmount}</Text>
                </div>
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
                if (record.status === 'pending') {
                    actions.push(
                        <Button key="pay" type="primary" size="small" onClick={() => handlePay(record)}>
                            支付
                        </Button>,
                        <Button key="cancel" size="small" onClick={() => handleCancel(record)}>
                            取消
                        </Button>
                    );
                }
                if (record.status === 'completed') {
                    actions.push(
                        <Button key="review" type="link" size="small" onClick={() => handleReview(record)}>
                            评价
                        </Button>
                    );
                }
                if (record.status === 'in_progress') {
                    actions.push(
                        <Button key="detail" type="link" size="small">
                            查看详情
                        </Button>
                    );
                }
                return <Space>{actions}</Space>;
            },
        },
    ];

    const filteredOrders = activeTab === 'all' 
        ? orders 
        : orders.filter(o => o.status === activeTab);

    const tabItems = [
        { key: 'all', label: '全部', icon: <ClockCircleOutlined /> },
        { key: 'pending', label: '待支付', icon: <ExclamationCircleOutlined /> },
        { key: 'in_progress', label: '进行中', icon: <ClockCircleOutlined /> },
        { key: 'completed', label: '已完成', icon: <CheckCircleOutlined /> },
        { key: 'canceled', label: '已取消', icon: <CloseCircleOutlined /> },
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
                />
            </Card>

            <Modal
                title="评价订单"
                open={reviewModal}
                onOk={submitReview}
                onCancel={() => setReviewModal(false)}
                okText="提交评价"
            >
                {selectedOrder && (
                    <div>
                        <div style={{ marginBottom: 16 }}>
                            <Text>陪玩师：{selectedOrder.player.nickname}</Text>
                        </div>
                        <div style={{ marginBottom: 16 }}>
                            <Text>评分：</Text>
                            <Rate value={rating} onChange={setRating} />
                        </div>
                        <div>
                            <Text>评价内容：</Text>
                            <TextArea
                                rows={4}
                                value={reviewContent}
                                onChange={e => setReviewContent(e.target.value)}
                                placeholder="请输入您的评价..."
                                maxLength={500}
                                showCount
                            />
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default UserOrders;
