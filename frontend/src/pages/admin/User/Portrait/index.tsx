import React from 'react';
import { Card, Row, Col, Tag, Descriptions, Avatar, Input, Space } from 'antd';
import { PageContainer } from '@/components';
import { UserOutlined, SearchOutlined } from '@ant-design/icons';
import { RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar, ResponsiveContainer, Tooltip } from 'recharts';

const data = [
    { subject: '活跃度', A: 120, fullMark: 150 },
    { subject: '消费能力', A: 98, fullMark: 150 },
    { subject: '互动频率', A: 86, fullMark: 150 },
    { subject: '信誉分', A: 99, fullMark: 150 },
    { subject: '游戏时长', A: 85, fullMark: 150 },
    { subject: '技能水平', A: 65, fullMark: 150 },
];

const UserPortrait: React.FC = () => {
    return (
        <PageContainer title="用户画像分析" subTitle="分析用户群体特征、地域分布和偏好">
            <Card style={{ marginBottom: 24 }}>
                <Space>
                    <Input placeholder="输入用户ID/昵称搜索" style={{ width: 300 }} prefix={<SearchOutlined />} />
                </Space>
            </Card>

            <Row gutter={[24, 24]}>
                <Col span={8}>
                    <Card title="基础信息">
                        <div style={{ textAlign: 'center', marginBottom: 24 }}>
                            <Avatar size={100} icon={<UserOutlined />} style={{ backgroundColor: '#1890ff' }} />
                            <h2 style={{ marginTop: 16 }}>示例用户</h2>
                            <Tag color="blue">90后</Tag>
                            <Tag color="green">北京</Tag>
                            <Tag color="gold">高价值用户</Tag>
                        </div>
                        <Descriptions column={1} bordered size="small">
                            <Descriptions.Item label="注册时长">365天</Descriptions.Item>
                            <Descriptions.Item label="最后登录">2小时前</Descriptions.Item>
                            <Descriptions.Item label="偏好游戏">王者荣耀, LOL</Descriptions.Item>
                            <Descriptions.Item label="消费偏好">皮肤, 陪玩</Descriptions.Item>
                        </Descriptions>
                    </Card>
                </Col>
                <Col span={16}>
                    <Card title="能力模型">
                        <div style={{ height: 400 }}>
                            <ResponsiveContainer width="100%" height="100%">
                                <RadarChart cx="50%" cy="50%" outerRadius="80%" data={data}>
                                    <PolarGrid />
                                    <PolarAngleAxis dataKey="subject" />
                                    <PolarRadiusAxis />
                                    <Radar name="用户能力" dataKey="A" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
                                    <Tooltip />
                                </RadarChart>
                            </ResponsiveContainer>
                        </div>
                    </Card>
                </Col>
            </Row>
        </PageContainer>
    );
};

export default UserPortrait;
