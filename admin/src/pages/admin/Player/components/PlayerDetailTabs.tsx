/**
 * 陪玩师详情 Tabs 组件
 * 将陪玩师详情分为多个 Tab：基本信息、订单历史、收益统计、认证信息
 */
import React, { useState, useEffect } from 'react';
import {
    Tabs,
    Descriptions,
    Avatar,
    Tag,
    Space,
    Card,
    Row,
    Col,
    Statistic,
    Table,
    Typography,
    Empty,
    Spin,
    Timeline,
    Badge,
} from 'antd';
import {
    UserOutlined,
    StarOutlined,
    ShoppingOutlined,
    DollarOutlined,
    SafetyOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
} from '@ant-design/icons';
import type { Player } from '@/api/admin';
import { adminApi } from '@/api/admin';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

const { Text, Paragraph } = Typography;

interface PlayerDetailTabsProps {
    player: Player;
}

/**
 * 状态映射
 */
const statusMap = {
    pending: { color: 'gold', text: '待审核' },
    verified: { color: 'success', text: '已通过' },
    rejected: { color: 'error', text: '已拒绝' },
};

const orderStatusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'gold', text: '待确认' },
    confirmed: { color: 'blue', text: '已确认' },
    in_progress: { color: 'processing', text: '进行中' },
    completed: { color: 'success', text: '已完成' },
    canceled: { color: 'default', text: '已取消' },
    refunded: { color: 'error', text: '已退款' },
};

/**
 * 基本信息 Tab
 */
const BasicInfoTab: React.FC<{ player: Player }> = ({ player }) => (
    <>
        {/* 头像和基本信息卡片 */}
        <Card size="small" style={{ marginBottom: 16 }}>
            <div style={{ textAlign: 'center', marginBottom: 16 }}>
                <Avatar
                    size={80}
                    src={player.user?.avatarUrl || undefined}
                    icon={<UserOutlined />}
                />
                <h2 style={{ marginTop: 12, marginBottom: 4 }}>
                    {player.nickname || player.user?.name || '-'}
                </h2>
                <Tag color={statusMap[player.verificationStatus]?.color}>
                    {statusMap[player.verificationStatus]?.text}
                </Tag>
            </div>

            <Row gutter={16}>
                <Col span={8}>
                    <Statistic
                        title="评分"
                        value={player.ratingAverage?.toFixed(1) || '0.0'}
                        prefix={<StarOutlined />}
                    />
                </Col>
                <Col span={8}>
                    <Statistic title="评价数" value={player.ratingCount || 0} suffix="条" />
                </Col>
                <Col span={8}>
                    <Statistic
                        title="时薪"
                        value={player.hourlyRateCents ? (player.hourlyRateCents / 100).toFixed(2) : 0}
                        prefix="¥"
                    />
                </Col>
            </Row>
        </Card>

        {/* 详细信息 */}
        <Descriptions title="详细信息" column={2} size="small" bordered>
            <Descriptions.Item label="ID">{player.id}</Descriptions.Item>
            <Descriptions.Item label="用户ID">{player.userId}</Descriptions.Item>
            <Descriptions.Item label="昵称">{player.nickname || '-'}</Descriptions.Item>
            <Descriptions.Item label="段位">{player.rank || '-'}</Descriptions.Item>
            <Descriptions.Item label="主游戏">{player.mainGame?.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="时薪">
                {player.hourlyRateCents ? `¥${(player.hourlyRateCents / 100).toFixed(2)}` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="技能标签" span={2}>
                <Space wrap>
                    {(player.skillTags || []).map(tag => <Tag key={tag}>{tag}</Tag>)}
                    {(!player.skillTags || player.skillTags.length === 0) && '-'}
                </Space>
            </Descriptions.Item>
            <Descriptions.Item label="个人简介" span={2}>
                <Paragraph ellipsis={{ rows: 3, expandable: true }}>
                    {player.bio || '暂无介绍'}
                </Paragraph>
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
                {player.createdAt ? dayjs(player.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
                {player.updatedAt ? dayjs(player.updatedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
            </Descriptions.Item>
        </Descriptions>
    </>
);


/**
 * 订单历史 Tab
 */
const OrderHistoryTab: React.FC<{ player: Player }> = ({ player }) => {
    const [loading, setLoading] = useState(false);
    const [orders, setOrders] = useState<Array<{
        id: number;
        orderNo: string;
        status: string;
        totalPriceCents: number;
        createdAt: string;
        user?: { name?: string };
        game?: { name?: string };
    }>>([]);

    useEffect(() => {
        const loadOrders = async () => {
            setLoading(true);
            try {
                const response = await adminApi.getOrders({ playerId: player.id, page_size: 20 });
                if (response.data.success) {
                    setOrders(response.data.data || []);
                }
            } catch (error) {
                logger.error('Load player orders error:', error);
            } finally {
                setLoading(false);
            }
        };
        loadOrders();
    }, [player.id]);

    const columns = [
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            width: 150,
            render: (text: string) => <Text copyable={{ text }}>{text.slice(-8)}</Text>,
        },
        {
            title: '用户',
            key: 'user',
            width: 100,
            render: (_: unknown, record: typeof orders[0]) => record.user?.name || '-',
        },
        {
            title: '游戏',
            key: 'game',
            width: 100,
            render: (_: unknown, record: typeof orders[0]) => record.game?.name || '-',
        },
        {
            title: '金额',
            dataIndex: 'totalPriceCents',
            key: 'totalPriceCents',
            width: 100,
            render: (cents: number) => <Text type="success">¥{(cents / 100).toFixed(2)}</Text>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => (
                <Tag color={orderStatusMap[status]?.color || 'default'}>
                    {orderStatusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 120,
            render: (date: string) => date ? dayjs(date).format('MM-DD HH:mm') : '-',
        },
    ];

    if (loading) {
        return <Spin style={{ display: 'block', textAlign: 'center', padding: 40 }} />;
    }

    if (orders.length === 0) {
        return <Empty description="暂无订单记录" />;
    }

    return (
        <Table
            columns={columns}
            dataSource={orders}
            rowKey="id"
            size="small"
            pagination={{ pageSize: 10, showSizeChanger: false }}
        />
    );
};

/**
 * 收益统计 Tab
 */
const EarningsTab: React.FC<{ player: Player }> = ({ player }) => {
    const [loading, setLoading] = useState(false);
    const [stats, setStats] = useState<{
        totalEarnings: number;
        monthlyEarnings: number;
        completedOrders: number;
        avgOrderValue: number;
    } | null>(null);

    useEffect(() => {
        const loadStats = async () => {
            setLoading(true);
            try {
                // 获取订单统计
                const response = await adminApi.getOrders({ 
                    playerId: player.id, 
                    status: 'completed',
                    page_size: 1000 
                });
                if (response.data.success && response.data.data) {
                    const orders = response.data.data;
                    const totalEarnings = orders.reduce((sum: number, o: { totalPriceCents: number }) => 
                        sum + (o.totalPriceCents || 0), 0);
                    
                    // 本月订单
                    const startOfMonth = dayjs().startOf('month');
                    const monthlyOrders = orders.filter((o: { createdAt: string }) => 
                        dayjs(o.createdAt).isAfter(startOfMonth));
                    const monthlyEarnings = monthlyOrders.reduce((sum: number, o: { totalPriceCents: number }) => 
                        sum + (o.totalPriceCents || 0), 0);

                    setStats({
                        totalEarnings,
                        monthlyEarnings,
                        completedOrders: orders.length,
                        avgOrderValue: orders.length > 0 ? totalEarnings / orders.length : 0,
                    });
                }
            } catch (error) {
                logger.error('Load player stats error:', error);
            } finally {
                setLoading(false);
            }
        };
        loadStats();
    }, [player.id]);

    if (loading) {
        return <Spin style={{ display: 'block', textAlign: 'center', padding: 40 }} />;
    }

    return (
        <Row gutter={[16, 16]}>
            <Col span={12}>
                <Card size="small">
                    <Statistic
                        title="总收益"
                        value={stats ? (stats.totalEarnings / 100).toFixed(2) : 0}
                        prefix={<DollarOutlined />}
                        suffix="元"
                        />
                </Card>
            </Col>
            <Col span={12}>
                <Card size="small">
                    <Statistic
                        title="本月收益"
                        value={stats ? (stats.monthlyEarnings / 100).toFixed(2) : 0}
                        prefix={<DollarOutlined />}
                        suffix="元"
                        />
                </Card>
            </Col>
            <Col span={12}>
                <Card size="small">
                    <Statistic
                        title="完成订单数"
                        value={stats?.completedOrders || 0}
                        prefix={<ShoppingOutlined />}
                        suffix="单"
                    />
                </Card>
            </Col>
            <Col span={12}>
                <Card size="small">
                    <Statistic
                        title="平均订单金额"
                        value={stats ? (stats.avgOrderValue / 100).toFixed(2) : 0}
                        prefix="¥"
                    />
                </Card>
            </Col>
        </Row>
    );
};


/**
 * 认证信息 Tab
 */
const CertificationTab: React.FC<{ player: Player }> = ({ player }) => {
    return (
        <>
            {/* 认证状态 */}
            <Card size="small" style={{ marginBottom: 16 }}>
                <div style={{ textAlign: 'center' }}>
                    <Badge
                        status={
                            player.verificationStatus === 'verified' ? 'success' :
                            player.verificationStatus === 'pending' ? 'processing' : 'error'
                        }
                        text={
                            <span style={{ fontSize: 18 }}>
                                {statusMap[player.verificationStatus]?.text || player.verificationStatus}
                            </span>
                        }
                    />
                </div>
            </Card>

            {/* 审核时间线 */}
            <Timeline
                items={[
                    {
                        color: 'green',
                        dot: <CheckCircleOutlined />,
                        children: (
                            <div>
                                <div>提交申请</div>
                                <Text type="secondary">
                                    {player.createdAt ? dayjs(player.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                                </Text>
                            </div>
                        ),
                    },
                    ...(player.verifiedAt ? [{
                        color: player.verificationStatus === 'verified' ? 'green' : 'red',
                        dot: player.verificationStatus === 'verified' ? 
                            <CheckCircleOutlined /> : <ClockCircleOutlined />,
                        children: (
                            <div>
                                <div>{player.verificationStatus === 'verified' ? '审核通过' : '审核拒绝'}</div>
                                <Text type="secondary">
                                    {dayjs(player.verifiedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Text>
                                {player.verifiedBy && (
                                    <div><Text type="secondary">审核人ID: {player.verifiedBy}</Text></div>
                                )}
                            </div>
                        ),
                    }] : []),
                ]}
            />

            {/* 审核详情 */}
            {player.verifiedAt && (
                <Descriptions column={1} size="small" bordered style={{ marginTop: 16 }}>
                    <Descriptions.Item label="审核时间">
                        {dayjs(player.verifiedAt).format('YYYY-MM-DD HH:mm:ss')}
                    </Descriptions.Item>
                    <Descriptions.Item label="审核人ID">
                        {player.verifiedBy || '-'}
                    </Descriptions.Item>
                    {player.verificationStatus === 'rejected' && player.rejectReason && (
                        <Descriptions.Item label="拒绝原因">
                            <Text type="danger">{player.rejectReason}</Text>
                        </Descriptions.Item>
                    )}
                    {player.verifyRemark && (
                        <Descriptions.Item label="审核备注">
                            {player.verifyRemark}
                        </Descriptions.Item>
                    )}
                </Descriptions>
            )}
        </>
    );
};

/**
 * 陪玩师详情 Tabs 主组件
 */
const PlayerDetailTabs: React.FC<PlayerDetailTabsProps> = ({ player }) => {
    const tabItems = [
        {
            key: 'basic',
            label: (
                <span>
                    <UserOutlined />
                    基本信息
                </span>
            ),
            children: <BasicInfoTab player={player} />,
        },
        {
            key: 'orders',
            label: (
                <span>
                    <ShoppingOutlined />
                    订单历史
                </span>
            ),
            children: <OrderHistoryTab player={player} />,
        },
        {
            key: 'earnings',
            label: (
                <span>
                    <DollarOutlined />
                    收益统计
                </span>
            ),
            children: <EarningsTab player={player} />,
        },
        {
            key: 'certification',
            label: (
                <span>
                    <SafetyOutlined />
                    认证信息
                </span>
            ),
            children: <CertificationTab player={player} />,
        },
    ];

    return <Tabs items={tabItems} />;
};

export default PlayerDetailTabs;
