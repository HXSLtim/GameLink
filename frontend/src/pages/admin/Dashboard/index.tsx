/**
 * 管理后台仪表盘
 */
import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Table, Tag, Avatar, Typography, Space, message, Select, theme } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    UserOutlined,
    ShoppingCartOutlined,
    TeamOutlined,
    DollarOutlined,
    ClockCircleOutlined,
    RiseOutlined,
} from '@ant-design/icons';
import { PageContainer, StatCard } from '@/components';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { adminApi } from '@/api/admin';
import type { DashboardStats, TrendData, TopPlayer } from '@/api/admin';

const { Text } = Typography;
const { Option } = Select;

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'gold', text: '待确认' },
    confirmed: { color: 'blue', text: '已确认' },
    in_progress: { color: 'processing', text: '进行中' },
    completed: { color: 'success', text: '已完成' },
    canceled: { color: 'default', text: '已取消' },
    refunded: { color: 'error', text: '已退款' },
};

const paymentStatusMap: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    refunded: '已退款',
    failed: '支付失败',
};

const Dashboard: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(true);
    const [days, setDays] = useState(7);
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [revenueTrend, setRevenueTrend] = useState<TrendData[]>([]);
    const [userGrowth, setUserGrowth] = useState<TrendData[]>([]);
    const [recentOrders, setRecentOrders] = useState<any[]>([]);
    const [topPlayers, setTopPlayers] = useState<TopPlayer[]>([]);

    const loadData = async () => {
        setLoading(true);
        try {
            const [dashboardRes, revenueRes, userGrowthRes, ordersRes, topPlayersRes] = await Promise.all([
                adminApi.getDashboardStats(),
                adminApi.getRevenueTrend({ days }),
                adminApi.getUserGrowth({ days }),
                adminApi.getOrders({ page: 1, page_size: 5 }),
                adminApi.getTopPlayers({ limit: 5 })
            ]);

            // @ts-ignore
            setStats(dashboardRes.data?.data || dashboardRes.data || {});
            // @ts-ignore
            setRevenueTrend(revenueRes.data?.data || revenueRes.data || []);
            // @ts-ignore
            setUserGrowth(userGrowthRes.data?.data || userGrowthRes.data || []);

            // @ts-ignore
            const ordersData = ordersRes.data?.data || ordersRes.data || [];
            // Handle case where API returns array directly or { list: [] }
            setRecentOrders(Array.isArray(ordersData) ? ordersData : (ordersData.list || []));

            // @ts-ignore
            setTopPlayers(topPlayersRes.data?.data || topPlayersRes.data || []);

        } catch (error) {
            console.error(error);
            message.error('加载仪表盘数据失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadData();
    }, [days]);

    const orderColumns: ColumnsType<any> = [
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            render: text => <Text copyable={{ text }}>{text}</Text>,
        },
        {
            title: '用户',
            dataIndex: 'userId',
            key: 'userId',
            render: (userId) => (
                <Space>
                    <Avatar style={{ backgroundColor: token.colorPrimary }} icon={<UserOutlined />} />
                    <Text>用户 #{userId}</Text>
                </Space>
            ),
        },
        {
            title: '金额',
            dataIndex: 'totalPriceCents',
            key: 'totalPriceCents',
            render: price => <Text strong>¥{(price / 100).toFixed(2)}</Text>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: status => (
                <Tag color={statusMap[status]?.color || 'default'}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: text => <Text type="secondary" style={{ fontSize: 12 }}>{new Date(text).toLocaleString()}</Text>,
        },
    ];

    const orderStatusData = stats?.ordersByStatus
        ? Object.entries(stats.ordersByStatus).map(([name, value]) => ({
            name: statusMap[name]?.text || name,
            value
        }))
        : [];

    const paymentStatusData = stats?.paymentsByStatus
        ? Object.entries(stats.paymentsByStatus).map(([name, value]) => ({
            name: paymentStatusMap[name] || name,
            value
        }))
        : [];

    return (
        <PageContainer
            title="仪表盘"
            subTitle="系统性能概览"
            extra={
                <Select value={days} onChange={setDays} style={{ width: 120 }}>
                    <Option value={7}>近7天</Option>
                    <Option value={30}>近30天</Option>
                    <Option value={90}>近90天</Option>
                </Select>
            }
        >
            {/* Stats Cards */}
            <Row gutter={[16, 16]}>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="总订单数"
                        value={stats?.totalOrders || 0}
                        icon={<ShoppingCartOutlined />}
                        iconBgColor={token.colorSuccess}
                        loading={loading}
                        suffix=""
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="交易总额"
                        value={(stats?.totalPaidAmountCents || 0) / 100}
                        icon={<DollarOutlined />}
                        iconBgColor={token.colorWarning}
                        loading={loading}
                        prefix="¥"
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="总用户数"
                        value={stats?.totalUsers || 0}
                        icon={<UserOutlined />}
                        iconBgColor={token.colorPrimary}
                        loading={loading}
                    />
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <StatCard
                        title="总陪玩数"
                        value={stats?.totalPlayers || 0}
                        icon={<TeamOutlined />}
                        iconBgColor="#722ed1"
                        loading={loading}
                    />
                </Col>
            </Row>

            {/* Status Distribution Charts */}
            <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
                <Col xs={24} lg={12}>
                    <Card title="订单状态分布" loading={loading} style={{ border: 'none' }}>
                        <div style={{ height: 300, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                            <ResponsiveContainer width="100%" height="100%" debounce={300} minWidth={0} minHeight={0}>
                                <PieChart>
                                    <Pie
                                        data={orderStatusData}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        label={({ name, percent }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
                                        outerRadius={100}
                                        fill="#8884d8"
                                        dataKey="value"
                                    >
                                        {orderStatusData.map((_, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }} itemStyle={{ color: token.colorText }} />
                                    <Legend />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
                <Col xs={24} lg={12}>
                    <Card title="支付状态分布" loading={loading} style={{ border: 'none' }}>
                        <div style={{ height: 300, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                            <ResponsiveContainer width="100%" height="100%" debounce={300} minWidth={0} minHeight={0}>
                                <PieChart>
                                    <Pie
                                        data={paymentStatusData}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        label={({ name, percent }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
                                        outerRadius={100}
                                        fill="#82ca9d"
                                        dataKey="value"
                                    >
                                        {paymentStatusData.map((_, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }} itemStyle={{ color: token.colorText }} />
                                    <Legend />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
            </Row>

            {/* Charts */}
            <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
                <Col xs={24} lg={12}>
                    <Card title="收入趋势 (近7天)" loading={loading} style={{ border: 'none' }}>
                        <div style={{ height: 300, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                            <ResponsiveContainer width="100%" height="100%" debounce={300} minWidth={0} minHeight={0}>
                                <LineChart data={revenueTrend}>
                                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                                    <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                                    <YAxis stroke={token.colorTextSecondary} />
                                    <Tooltip contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }} itemStyle={{ color: token.colorText }} />
                                    <Legend />
                                    <Line type="monotone" dataKey="value" name="收入" stroke={token.colorWarning} activeDot={{ r: 8 }} />
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
                <Col xs={24} lg={12}>
                    <Card title="用户增长 (近7天)" loading={loading} style={{ border: 'none' }}>
                        <div style={{ height: 300, width: '100%', minWidth: 0, overflow: 'hidden' }}>
                            <ResponsiveContainer width="100%" height="100%" debounce={300} minWidth={0} minHeight={0}>
                                <LineChart data={userGrowth}>
                                    <CartesianGrid strokeDasharray="3 3" stroke={token.colorSplit} />
                                    <XAxis dataKey="date" stroke={token.colorTextSecondary} />
                                    <YAxis stroke={token.colorTextSecondary} />
                                    <Tooltip contentStyle={{ backgroundColor: token.colorBgElevated, borderColor: token.colorBorder }} itemStyle={{ color: token.colorText }} />
                                    <Legend />
                                    <Line type="monotone" dataKey="value" name="新增用户" stroke={token.colorPrimary} activeDot={{ r: 8 }} />
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
            </Row>

            {/* Recent Orders & Top Players */}
            <Row gutter={[16, 16]} style={{ marginTop: 16, display: 'flex', flexWrap: 'wrap', alignItems: 'stretch' }}>
                <Col xs={24} xl={16} style={{ display: 'flex', flexDirection: 'column' }}>
                    <Card
                        title={
                            <Space>
                                <ClockCircleOutlined />
                                <span>最新订单</span>
                            </Space>
                        }
                        extra={<a href="/admin/biz/order">查看全部</a>}
                        loading={loading}
                        style={{ flex: 1, width: '100%', border: 'none' }}
                    >
                        <Table
                            columns={orderColumns}
                            dataSource={recentOrders}
                            rowKey="id"
                            pagination={false}
                            size="small"
                        />
                    </Card>
                </Col>
                <Col xs={24} xl={8} style={{ display: 'flex', flexDirection: 'column' }}>
                    <Card
                        title={
                            <Space>
                                <RiseOutlined />
                                <span>热门陪玩</span>
                            </Space>
                        }
                        loading={loading}
                        style={{ flex: 1, width: '100%', border: 'none' }}
                    >
                        <Table
                            columns={[
                                {
                                    title: '排名',
                                    key: 'rank',
                                    width: 60,
                                    render: (_: any, __: any, index: number) => <Tag color="gold">#{index + 1}</Tag>,
                                },
                                {
                                    title: '陪玩',
                                    key: 'player',
                                    render: (_: any, record: TopPlayer, index: number) => (
                                        <Space>
                                            <Avatar
                                                style={{ backgroundColor: `hsl(${index * 60}, 70%, 50%)` }}
                                                icon={<UserOutlined />}
                                            />
                                            <Text>{record.nickname}</Text>
                                        </Space>
                                    ),
                                },
                                {
                                    title: '评分',
                                    key: 'rating',
                                    render: (_: any, record: TopPlayer) => (
                                        <Space>
                                            <Text strong>{record.ratingAverage}</Text>
                                            <Text type="secondary">({record.ratingCount})</Text>
                                        </Space>
                                    ),
                                },
                            ]}
                            dataSource={topPlayers}
                            rowKey="playerId"
                            pagination={false}
                            size="small"
                        />
                    </Card>
                </Col>
            </Row>
        </PageContainer>
    );
};

export default Dashboard;
