/**
 * 订单详情 Tabs 组件
 * 将订单详情分为多个 Tab：基本信息、服务进度、支付信息、相关记录
 */
import React, { useState, useEffect } from 'react';
import {
    Tabs,
    Descriptions,
    Avatar,
    Tag,
    Space,
    Card,
    Typography,
    Timeline,
    Table,
    Empty,
    Spin,
    Row,
    Col,
    Statistic,
} from 'antd';
import {
    UserOutlined,
    ShoppingOutlined,
    DollarOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    FileTextOutlined,
    ExclamationCircleOutlined,
    CloseCircleOutlined,
} from '@ant-design/icons';
import type { Order } from '@/api/admin';
import { disputeApi } from '@/api/dispute';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

const { Text, Title } = Typography;

interface OrderDetailTabsProps {
    order: Order;
}

/**
 * 订单状态映射
 */
const statusMap: Record<string, { color: string; text: string; icon: React.ReactNode }> = {
    pending: { color: 'gold', text: '待确认', icon: <ClockCircleOutlined /> },
    confirmed: { color: 'blue', text: '已确认', icon: <CheckCircleOutlined /> },
    in_progress: { color: 'processing', text: '进行中', icon: <ClockCircleOutlined /> },
    completed: { color: 'success', text: '已完成', icon: <CheckCircleOutlined /> },
    canceled: { color: 'default', text: '已取消', icon: <CloseCircleOutlined /> },
    refunded: { color: 'error', text: '已退款', icon: <ExclamationCircleOutlined /> },
};

/**
 * 基本信息 Tab
 */
const BasicInfoTab: React.FC<{ order: Order }> = ({ order }) => (
    <>
        {/* 状态卡片 */}
        <Card size="small" style={{ marginBottom: 16 }}>
            <div style={{ textAlign: 'center' }}>
                <Tag color={statusMap[order.status]?.color} style={{ fontSize: 16, padding: '4px 16px' }}>
                    {statusMap[order.status]?.icon} {statusMap[order.status]?.text}
                </Tag>
                <Title level={2} style={{ margin: '16px 0 0' }}>
                    ¥{(order.totalPriceCents / 100).toFixed(2)}
                </Title>
            </div>
        </Card>

        {/* 基本信息 */}
        <Descriptions title="订单信息" column={2} size="small" bordered>
            <Descriptions.Item label="订单号" span={2}>
                <Text copyable>{order.orderNo}</Text>
            </Descriptions.Item>
            <Descriptions.Item label="游戏">{order.game?.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="标题">{order.title || '-'}</Descriptions.Item>
            <Descriptions.Item label="金额">
                ¥{(order.totalPriceCents / 100).toFixed(2)} {order.currency}
            </Descriptions.Item>
            <Descriptions.Item label="状态">
                <Tag color={statusMap[order.status]?.color}>
                    {statusMap[order.status]?.text}
                </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="预约开始">
                {order.scheduledStart ? dayjs(order.scheduledStart).format('YYYY-MM-DD HH:mm') : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="预约结束">
                {order.scheduledEnd ? dayjs(order.scheduledEnd).format('YYYY-MM-DD HH:mm') : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
                {order.createdAt ? dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="完成时间">
                {order.completedAt ? dayjs(order.completedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
            </Descriptions.Item>
            {order.description && (
                <Descriptions.Item label="描述" span={2}>{order.description}</Descriptions.Item>
            )}
            {order.cancelReason && (
                <Descriptions.Item label="取消原因" span={2}>
                    <Text type="danger">{order.cancelReason}</Text>
                </Descriptions.Item>
            )}
        </Descriptions>

        {/* 用户信息 */}
        <Descriptions title="用户信息" column={2} size="small" style={{ marginTop: 16 }}>
            <Descriptions.Item label="用户">
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={order.user?.avatarUrl} />
                    {order.user?.name || `用户${order.userId}`}
                </Space>
            </Descriptions.Item>
            <Descriptions.Item label="用户ID">{order.userId}</Descriptions.Item>
        </Descriptions>

        <Descriptions title="陪玩师信息" column={2} size="small">
            <Descriptions.Item label="陪玩师">
                <Space>
                    <Avatar
                        size="small"
                        icon={<UserOutlined />}
                        src={order.player?.user?.avatarUrl}
                        style={{ backgroundColor: '#722ed1' }}
                    />
                    {order.player?.nickname || (order.playerId ? `陪玩师${order.playerId}` : '未分配')}
                </Space>
            </Descriptions.Item>
            <Descriptions.Item label="陪玩师ID">{order.playerId || '-'}</Descriptions.Item>
        </Descriptions>
    </>
);


/**
 * 服务进度 Tab
 */
const ProgressTab: React.FC<{ order: Order }> = ({ order }) => {
    const getTimelineItems = () => {
        const items = [
            {
                color: 'green',
                dot: <CheckCircleOutlined />,
                children: (
                    <div>
                        <div>订单创建</div>
                        <Text type="secondary">
                            {order.createdAt ? dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                        </Text>
                    </div>
                ),
            },
        ];

        if (['confirmed', 'in_progress', 'completed'].includes(order.status)) {
            items.push({
                color: 'green',
                dot: <CheckCircleOutlined />,
                children: (
                    <div>
                        <div>订单确认</div>
                        <Text type="secondary">陪玩师已接单</Text>
                    </div>
                ),
            });
        } else if (order.status === 'pending') {
            items.push({
                color: 'gray',
                dot: <ClockCircleOutlined />,
                children: <div>等待确认</div>,
            });
        }

        if (['in_progress', 'completed'].includes(order.status)) {
            items.push({
                color: 'green',
                dot: <CheckCircleOutlined />,
                children: (
                    <div>
                        <div>开始服务</div>
                        <Text type="secondary">
                            {order.scheduledStart ? dayjs(order.scheduledStart).format('YYYY-MM-DD HH:mm') : '服务进行中'}
                        </Text>
                    </div>
                ),
            });
        } else if (order.status === 'confirmed') {
            items.push({
                color: 'gray',
                dot: <ClockCircleOutlined />,
                children: <div>等待服务开始</div>,
            });
        }

        if (order.status === 'completed') {
            items.push({
                color: 'green',
                dot: <CheckCircleOutlined />,
                children: (
                    <div>
                        <div>服务完成</div>
                        <Text type="secondary">
                            {order.completedAt ? dayjs(order.completedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                        </Text>
                    </div>
                ),
            });
        } else if (order.status === 'canceled') {
            items.push({
                color: 'red',
                dot: <CloseCircleOutlined />,
                children: (
                    <div>
                        <div>订单取消</div>
                        {order.cancelReason && <Text type="danger">{order.cancelReason}</Text>}
                    </div>
                ),
            });
        } else if (order.status === 'refunded') {
            items.push({
                color: 'red',
                dot: <ExclamationCircleOutlined />,
                children: <div>已退款</div>,
            });
        }

        return items;
    };

    return (
        <Card size="small">
            <Timeline items={getTimelineItems()} />
        </Card>
    );
};

/**
 * 支付信息 Tab
 */
const PaymentTab: React.FC<{ order: Order }> = ({ order }) => {
    return (
        <>
            <Row gutter={[16, 16]}>
                <Col span={12}>
                    <Card size="small">
                        <Statistic
                            title="订单金额"
                            value={(order.totalPriceCents / 100).toFixed(2)}
                            prefix="¥"
                        />
                    </Card>
                </Col>
                <Col span={12}>
                    <Card size="small">
                        <Statistic
                            title="支付状态"
                            value={order.status === 'completed' || order.status === 'in_progress' ? '已支付' :
                                   order.status === 'refunded' ? '已退款' : '待支付'}
                        />
                    </Card>
                </Col>
            </Row>

            <Descriptions column={1} size="small" bordered style={{ marginTop: 16 }}>
                <Descriptions.Item label="订单金额">
                    ¥{(order.totalPriceCents / 100).toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="货币类型">
                    {order.currency || 'CNY'}
                </Descriptions.Item>
                <Descriptions.Item label="支付方式">
                    余额支付
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">
                    {order.createdAt ? dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                </Descriptions.Item>
            </Descriptions>
        </>
    );
};


/**
 * 相关记录 Tab - 显示纠纷、评价等关联记录
 */
const RelatedRecordsTab: React.FC<{ order: Order }> = ({ order }) => {
    const [loading, setLoading] = useState(false);
    const [disputes, setDisputes] = useState<Array<{
        id: number;
        status: string;
        reason: string;
        createdAt: string;
    }>>([]);

    useEffect(() => {
        const loadRelatedData = async () => {
            setLoading(true);
            try {
                // 尝试加载相关纠纷
                const response = await disputeApi.getDisputes({ orderNo: order.orderNo });
                if (response.data.success && response.data.data?.disputes) {
                    setDisputes(response.data.data.disputes);
                }
            } catch (error) {
                logger.error('Load related records error:', error);
            } finally {
                setLoading(false);
            }
        };
        loadRelatedData();
    }, [order.orderNo]);

    const disputeColumns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 80,
            render: (status: string) => {
                const statusColors: Record<string, string> = {
                    pending: 'gold',
                    assigned: 'blue',
                    mediating: 'processing',
                    resolved: 'success',
                    rejected: 'error',
                };
                return <Tag color={statusColors[status] || 'default'}>{status}</Tag>;
            },
        },
        {
            title: '原因',
            dataIndex: 'reason',
            key: 'reason',
            ellipsis: true,
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 120,
            render: (date: string) => date ? dayjs(date).format('MM-DD HH:mm') : '-',
        },
    ];

    if (loading) {
        return <Spin style={{ display: 'block', textAlign: 'center', padding: 40 }} />;
    }

    return (
        <>
            <Card title="纠纷记录" size="small" style={{ marginBottom: 16 }}>
                {disputes.length > 0 ? (
                    <Table
                        columns={disputeColumns}
                        dataSource={disputes}
                        rowKey="id"
                        size="small"
                        pagination={false}
                    />
                ) : (
                    <Empty description="暂无纠纷记录" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                )}
            </Card>

            {/* 评价信息 - 如果有的话 */}
            <Card title="评价信息" size="small">
                <Empty description="暂无评价" image={Empty.PRESENTED_IMAGE_SIMPLE} />
            </Card>
        </>
    );
};

/**
 * 订单详情 Tabs 主组件
 */
const OrderDetailTabs: React.FC<OrderDetailTabsProps> = ({ order }) => {
    const tabItems = [
        {
            key: 'basic',
            label: (
                <span>
                    <ShoppingOutlined />
                    基本信息
                </span>
            ),
            children: <BasicInfoTab order={order} />,
        },
        {
            key: 'progress',
            label: (
                <span>
                    <ClockCircleOutlined />
                    服务进度
                </span>
            ),
            children: <ProgressTab order={order} />,
        },
        {
            key: 'payment',
            label: (
                <span>
                    <DollarOutlined />
                    支付信息
                </span>
            ),
            children: <PaymentTab order={order} />,
        },
        {
            key: 'records',
            label: (
                <span>
                    <FileTextOutlined />
                    相关记录
                </span>
            ),
            children: <RelatedRecordsTab order={order} />,
        },
    ];

    return <Tabs items={tabItems} />;
};

export default OrderDetailTabs;
