import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, message } from 'antd';
import { PageContainer } from '@/components';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { UserOutlined, ClockCircleOutlined, DollarOutlined } from '@ant-design/icons';
import { adminApi } from '@/api/admin';
import type { UserBehaviorStats, UserDistribution, TrendData } from '@/api/admin';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const UserBehavior: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [stats, setStats] = useState<UserBehaviorStats | null>(null);
    const [trendData, setTrendData] = useState<TrendData[]>([]);
    const [distribution, setDistribution] = useState<UserDistribution | null>(null);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [statsRes, trendRes, distributionRes] = await Promise.all([
                adminApi.getUserBehaviorStats(),
                adminApi.getUserActivityTrend({ days: 7 }),
                adminApi.getUserDistribution()
            ]);

            if (statsRes.data.success) {
                setStats(statsRes.data.data);
            }
            if (trendRes.data.success) {
                setTrendData(trendRes.data.data || []);
            }
            if (distributionRes.data.success) {
                setDistribution(distributionRes.data.data);
            }
        } catch (error) {
            message.error('获取用户行为数据失败');
            console.error('Failed to fetch user behavior data:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <PageContainer title="用户行为分析" subTitle="分析用户登录、使用习惯和消费行为">
            <Spin spinning={loading}>
                <Row gutter={[16, 16]}>
                    <Col span={8}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="日活跃用户 (DAU)"
                                value={stats?.dau || 0}
                                prefix={<UserOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="平均在线时长"
                                value={stats?.avgOnlineTime || '0m'}
                                prefix={<ClockCircleOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="人均消费"
                                value={stats?.avgConsumption || 0}
                                prefix={<DollarOutlined />}
                                precision={2}
                            />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
                    <Col span={16}>
                        <Card title="用户活动趋势">
                            <div style={{ height: 300 }}>
                                <ResponsiveContainer width="100%" height="100%">
                                    <LineChart data={trendData}>
                                        <CartesianGrid strokeDasharray="3 3" />
                                        <XAxis dataKey="date" />
                                        <YAxis />
                                        <Tooltip />
                                        <Legend />
                                        <Line
                                            type="monotone"
                                            dataKey="value"
                                            stroke="#8884d8"
                                            activeDot={{ r: 8 }}
                                            name="活跃用户数"
                                        />
                                    </LineChart>
                                </ResponsiveContainer>
                            </div>
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card title="用户地域分布">
                            <div style={{ height: 300 }}>
                                <ResponsiveContainer width="100%" height="100%">
                                    <PieChart>
                                        <Pie
                                            data={distribution?.byRegion || []}
                                            cx="50%"
                                            cy="50%"
                                            innerRadius={60}
                                            outerRadius={80}
                                            fill="#8884d8"
                                            paddingAngle={5}
                                            dataKey="value"
                                            label={(entry) => entry.name}
                                        >
                                            {(distribution?.byRegion || []).map((_entry, index) => (
                                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                            ))}
                                        </Pie>
                                        <Tooltip />
                                    </PieChart>
                                </ResponsiveContainer>
                            </div>
                        </Card>
                    </Col>
                </Row>
            </Spin>
        </PageContainer>
    );
};

export default UserBehavior;
