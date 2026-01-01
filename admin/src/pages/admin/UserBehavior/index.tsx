/**
 * 用户行为分析页面
 */
import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, Select, App } from 'antd';
import {
    UserOutlined,
    ClockCircleOutlined,
    DollarOutlined,
    RiseOutlined,
    TeamOutlined,
    EnvironmentOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@/components';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

interface BehaviorStats {
    dau: number;
    mau: number;
    avgOnlineTime: number;
    avgSpending: number;
    newUsers: number;
    activeRate: number;
}

interface TrendData {
    date: string;
    dau: number;
    newUsers: number;
    orders: number;
}

interface DistributionData {
    regions: { name: string; count: number }[];
    ageGroups: { range: string; count: number }[];
    devices: { type: string; count: number }[];
}

const UserBehaviorPage: React.FC = () => {
    const { message } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [stats, setStats] = useState<BehaviorStats | null>(null);
    const [trend, setTrend] = useState<TrendData[]>([]);
    const [distribution, setDistribution] = useState<DistributionData | null>(null);
    const [days, setDays] = useState(7);

    const loadStats = async () => {
        try {
            const response = await apiClient.get('/admin/users/behavior/stats');
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (error) {
            console.error('Load stats error:', error);
        }
    };

    const loadTrend = async () => {
        try {
            const response = await apiClient.get('/admin/users/behavior/trend', { params: { days } });
            if (response.data.success) {
                setTrend(response.data.data || []);
            }
        } catch (error) {
            console.error('Load trend error:', error);
        }
    };

    const loadDistribution = async () => {
        try {
            const response = await apiClient.get('/admin/users/behavior/distribution');
            if (response.data.success) {
                setDistribution(response.data.data);
            }
        } catch (error) {
            console.error('Load distribution error:', error);
        }
    };

    const loadData = async () => {
        setLoading(true);
        try {
            await Promise.all([loadStats(), loadTrend(), loadDistribution()]);
        } catch {
            message.error('加载数据失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        loadTrend();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [days]);


    return (
        <PageContainer title="用户行为分析" subTitle="用户活跃度、消费习惯和分布统计">
            <Spin spinning={loading}>
                {/* 核心指标 */}
                <Row gutter={16} style={{ marginBottom: 24 }}>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="日活用户(DAU)" value={stats?.dau || 0} prefix={<UserOutlined />} />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="月活用户(MAU)" value={stats?.mau || 0} prefix={<TeamOutlined />} />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="平均在线时长" value={stats?.avgOnlineTime || 0} suffix="分钟" prefix={<ClockCircleOutlined />} />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="人均消费" value={stats?.avgSpending || 0} precision={2} prefix={<DollarOutlined />} suffix="元" />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="新增用户" value={stats?.newUsers || 0} prefix={<RiseOutlined />} />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic title="活跃率" value={stats?.activeRate || 0} precision={1} suffix="%" />
                        </Card>
                    </Col>
                </Row>

                {/* 趋势图 */}
                <Card title="用户活动趋势" extra={
                    <Select value={days} onChange={setDays} style={{ width: 120 }}>
                        <Select.Option value={7}>最近7天</Select.Option>
                        <Select.Option value={14}>最近14天</Select.Option>
                        <Select.Option value={30}>最近30天</Select.Option>
                    </Select>
                } style={{ marginBottom: 24 }}>
                    <div style={{ height: 300, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                        {trend.length > 0 ? (
                            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                <thead>
                                    <tr style={{ borderBottom: '1px solid #f0f0f0' }}>
                                        <th style={{ padding: 8, textAlign: 'left' }}>日期</th>
                                        <th style={{ padding: 8, textAlign: 'right' }}>DAU</th>
                                        <th style={{ padding: 8, textAlign: 'right' }}>新增用户</th>
                                        <th style={{ padding: 8, textAlign: 'right' }}>订单数</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {trend.map((item, index) => (
                                        <tr key={index} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                            <td style={{ padding: 8 }}>{dayjs(item.date).format('MM-DD')}</td>
                                            <td style={{ padding: 8, textAlign: 'right' }}>{item.dau}</td>
                                            <td style={{ padding: 8, textAlign: 'right' }}>{item.newUsers}</td>
                                            <td style={{ padding: 8, textAlign: 'right' }}>{item.orders}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        ) : (
                            <span style={{ color: '#999' }}>暂无数据</span>
                        )}
                    </div>
                </Card>

                {/* 用户分布 */}
                <Row gutter={16}>
                    <Col span={8}>
                        <Card title={<><EnvironmentOutlined /> 地域分布</>}>
                            <div style={{ height: 200 }}>
                                {distribution?.regions?.length ? (
                                    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                                        {distribution.regions.slice(0, 5).map((r, i) => (
                                            <li key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                                                <span>{r.name}</span>
                                                <span style={{ color: '#1890ff' }}>{r.count}</span>
                                            </li>
                                        ))}
                                    </ul>
                                ) : <div style={{ textAlign: 'center', color: '#999', paddingTop: 80 }}>暂无数据</div>}
                            </div>
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card title="年龄分布">
                            <div style={{ height: 200 }}>
                                {distribution?.ageGroups?.length ? (
                                    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                                        {distribution.ageGroups.map((a, i) => (
                                            <li key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                                                <span>{a.range}</span>
                                                <span style={{ color: '#52c41a' }}>{a.count}</span>
                                            </li>
                                        ))}
                                    </ul>
                                ) : <div style={{ textAlign: 'center', color: '#999', paddingTop: 80 }}>暂无数据</div>}
                            </div>
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card title="设备分布">
                            <div style={{ height: 200 }}>
                                {distribution?.devices?.length ? (
                                    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                                        {distribution.devices.map((d, i) => (
                                            <li key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                                                <span>{d.type}</span>
                                                <span style={{ color: '#faad14' }}>{d.count}</span>
                                            </li>
                                        ))}
                                    </ul>
                                ) : <div style={{ textAlign: 'center', color: '#999', paddingTop: 80 }}>暂无数据</div>}
                            </div>
                        </Card>
                    </Col>
                </Row>
            </Spin>
        </PageContainer>
    );
};

export default UserBehaviorPage;
