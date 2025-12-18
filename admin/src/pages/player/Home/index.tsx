/**
 * 陪玩师端首页
 * 展示接单状态、收益概览、待处理订单
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Statistic,
    Switch,
    Tag,
    Button,
    List,
    Avatar,
    Typography,
    Space,
    message,
    Badge,
    Progress,
} from 'antd';
import {
    DollarOutlined,
    ShoppingCartOutlined,
    StarOutlined,
    UserOutlined,
    BellOutlined,
    CheckCircleOutlined,
} from '@ant-design/icons';

const { Title, Text } = Typography;

interface DashboardStats {
    todayEarnings: number;
    monthEarnings: number;
    totalOrders: number;
    completedOrders: number;
    rating: number;
    pendingOrders: number;
}

interface PendingOrder {
    id: number;
    orderNo: string;
    user: { nickname: string; avatar: string };
    game: string;
    service: string;
    price: number;
    duration: number;
    createdAt: string;
}

const PlayerHome: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [isOnline, setIsOnline] = useState(true);
    const [stats, setStats] = useState<DashboardStats>({
        todayEarnings: 0,
        monthEarnings: 0,
        totalOrders: 0,
        completedOrders: 0,
        rating: 0,
        pendingOrders: 0,
    });
    const [pendingOrders, setPendingOrders] = useState<PendingOrder[]>([]);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 500));
            
            setStats({
                todayEarnings: 320,
                monthEarnings: 8650,
                totalOrders: 156,
                completedOrders: 148,
                rating: 4.9,
                pendingOrders: 3,
            });

            setPendingOrders([
                {
                    id: 1,
                    orderNo: 'ORD202412160010',
                    user: { nickname: '快乐玩家', avatar: '' },
                    game: '王者荣耀',
                    service: '上分陪玩',
                    price: 50,
                    duration: 2,
                    createdAt: '2024-12-16 14:30:00',
                },
                {
                    id: 2,
                    orderNo: 'ORD202412160011',
                    user: { nickname: '游戏达人', avatar: '' },
                    game: '英雄联盟',
                    service: '娱乐陪玩',
                    price: 60,
                    duration: 1,
                    createdAt: '2024-12-16 14:45:00',
                },
            ]);
        } catch {
            message.error('加载数据失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleStatusChange = (checked: boolean) => {
        setIsOnline(checked);
        message.success(checked ? '已上线，可以接单了' : '已下线');
    };

    const handleAcceptOrder = async (order: PendingOrder) => {
        message.success(`已接单：${order.orderNo}`);
        loadData();
    };

    const handleRejectOrder = async (order: PendingOrder) => {
        message.info(`已拒绝订单：${order.orderNo}`);
        loadData();
    };

    const completionRate = stats.totalOrders > 0 
        ? Math.round((stats.completedOrders / stats.totalOrders) * 100) 
        : 0;

    return (
        <div style={{ padding: 24 }}>
            {/* 状态切换 */}
            <Card style={{ marginBottom: 16 }}>
                <Row justify="space-between" align="middle">
                    <Col>
                        <Space size="large">
                            <Badge status={isOnline ? 'success' : 'default'} />
                            <Text strong style={{ fontSize: 16 }}>
                                {isOnline ? '在线接单中' : '已下线'}
                            </Text>
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <Text>接单状态：</Text>
                            <Switch
                                checked={isOnline}
                                onChange={handleStatusChange}
                                checkedChildren="在线"
                                unCheckedChildren="离线"
                            />
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* 数据概览 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic
                            title="今日收益"
                            value={stats.todayEarnings}
                            prefix={<DollarOutlined />}
                            suffix="元"
                            valueStyle={{ color: '#3f8600' }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic
                            title="本月收益"
                            value={stats.monthEarnings}
                            prefix={<DollarOutlined />}
                            suffix="元"
                            valueStyle={{ color: '#1890ff' }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic
                            title="总订单数"
                            value={stats.totalOrders}
                            prefix={<ShoppingCartOutlined />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic
                            title="综合评分"
                            value={stats.rating}
                            prefix={<StarOutlined />}
                            precision={1}
                            valueStyle={{ color: '#faad14' }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* 完成率 */}
            <Card style={{ marginBottom: 16 }} loading={loading}>
                <Row gutter={24} align="middle">
                    <Col span={12}>
                        <Title level={5}>订单完成率</Title>
                        <Progress
                            percent={completionRate}
                            status={completionRate >= 90 ? 'success' : 'normal'}
                            strokeColor={{
                                '0%': '#108ee9',
                                '100%': '#87d068',
                            }}
                        />
                    </Col>
                    <Col span={12}>
                        <Space size="large">
                            <Statistic title="已完成" value={stats.completedOrders} suffix="单" />
                            <Statistic title="总订单" value={stats.totalOrders} suffix="单" />
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* 待处理订单 */}
            <Card
                title={
                    <Space>
                        <BellOutlined />
                        <span>待接订单</span>
                        <Badge count={stats.pendingOrders} />
                    </Space>
                }
                loading={loading}
            >
                {pendingOrders.length === 0 ? (
                    <div style={{ textAlign: 'center', padding: 40 }}>
                        <CheckCircleOutlined style={{ fontSize: 48, color: '#52c41a' }} />
                        <Title level={5} style={{ marginTop: 16 }}>暂无待接订单</Title>
                        <Text type="secondary">保持在线状态，等待新订单</Text>
                    </div>
                ) : (
                    <List
                        dataSource={pendingOrders}
                        renderItem={order => (
                            <List.Item
                                actions={[
                                    <Button
                                        key="accept"
                                        type="primary"
                                        onClick={() => handleAcceptOrder(order)}
                                    >
                                        接单
                                    </Button>,
                                    <Button
                                        key="reject"
                                        onClick={() => handleRejectOrder(order)}
                                    >
                                        拒绝
                                    </Button>,
                                ]}
                            >
                                <List.Item.Meta
                                    avatar={<Avatar src={order.user.avatar} icon={<UserOutlined />} />}
                                    title={
                                        <Space>
                                            <Text strong>{order.user.nickname}</Text>
                                            <Tag color="blue">{order.game}</Tag>
                                        </Space>
                                    }
                                    description={
                                        <Space direction="vertical" size={0}>
                                            <Text>{order.service} · {order.duration}小时</Text>
                                            <Text type="secondary">{order.createdAt}</Text>
                                        </Space>
                                    }
                                />
                                <div style={{ textAlign: 'right' }}>
                                    <Text type="danger" strong style={{ fontSize: 18 }}>
                                        ¥{order.price * order.duration}
                                    </Text>
                                </div>
                            </List.Item>
                        )}
                    />
                )}
            </Card>
        </div>
    );
};

export default PlayerHome;
