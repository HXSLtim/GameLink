import React from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import { PageContainer } from '@/components';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { UserOutlined, ClockCircleOutlined, DollarOutlined } from '@ant-design/icons';

const data = [
    { name: 'Mon', uv: 4000, pv: 2400, amt: 2400 },
    { name: 'Tue', uv: 3000, pv: 1398, amt: 2210 },
    { name: 'Wed', uv: 2000, pv: 9800, amt: 2290 },
    { name: 'Thu', uv: 2780, pv: 3908, amt: 2000 },
    { name: 'Fri', uv: 1890, pv: 4800, amt: 2181 },
    { name: 'Sat', uv: 2390, pv: 3800, amt: 2500 },
    { name: 'Sun', uv: 3490, pv: 4300, amt: 2100 },
];

const pieData = [
    { name: 'Group A', value: 400 },
    { name: 'Group B', value: 300 },
    { name: 'Group C', value: 300 },
    { name: 'Group D', value: 200 },
];

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const UserBehavior: React.FC = () => {
    return (
        <PageContainer title="用户行为分析" subTitle="分析用户登录、使用习惯和消费行为">
            <Row gutter={[16, 16]}>
                <Col span={8}>
                    <Card>
                        <Statistic title="日活跃用户 (DAU)" value={1128} prefix={<UserOutlined />} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic title="平均在线时长" value="45m" prefix={<ClockCircleOutlined />} />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic title="人均消费" value={128} prefix={<DollarOutlined />} precision={2} />
                    </Card>
                </Col>
            </Row>

            <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
                <Col span={16}>
                    <Card title="流量趋势">
                        <div style={{ height: 300 }}>
                            <ResponsiveContainer width="100%" height="100%">
                                <LineChart data={data}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="name" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Line type="monotone" dataKey="pv" stroke="#8884d8" activeDot={{ r: 8 }} />
                                    <Line type="monotone" dataKey="uv" stroke="#82ca9d" />
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
                <Col span={8}>
                    <Card title="用户分布">
                        <div style={{ height: 300 }}>
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={pieData}
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={60}
                                        outerRadius={80}
                                        fill="#8884d8"
                                        paddingAngle={5}
                                        dataKey="value"
                                    >
                                        {pieData.map((entry, index) => (
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
        </PageContainer>
    );
};

export default UserBehavior;
