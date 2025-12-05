import React from 'react';
import { Timeline, Card, Tag, Typography, Space, DatePicker, Button } from 'antd';
import { ClockCircleOutlined, UserOutlined, LoginOutlined, EditOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';

const { Text } = Typography;

const Audit: React.FC = () => {
    const activities = [
        {
            id: 1,
            user: 'Admin',
            action: '更新了系统设置',
            target: '平台名称',
            time: '2023-11-26 10:30:00',
            type: 'system',
            ip: '192.168.1.1'
        },
        {
            id: 2,
            user: 'Admin',
            action: '封禁了用户',
            target: 'Spammer#1234',
            time: '2023-11-26 09:15:22',
            type: 'security',
            ip: '192.168.1.1'
        },
        {
            id: 3,
            user: 'Moderator',
            action: '审核通过了游戏',
            target: '新游戏申请: Valorant',
            time: '2023-11-25 16:45:10',
            type: 'audit',
            ip: '10.0.0.5'
        },
        {
            id: 4,
            user: 'System',
            action: '自动备份数据库',
            target: 'backup_20231125.sql',
            time: '2023-11-25 00:00:00',
            type: 'system',
            ip: 'localhost'
        },
        {
            id: 5,
            user: 'Admin',
            action: '登录系统',
            target: '-',
            time: '2023-11-24 08:30:00',
            type: 'login',
            ip: '192.168.1.1'
        },
    ];

    const getIcon = (type: string) => {
        switch (type) {
            case 'system': return <ClockCircleOutlined style={{ fontSize: '16px' }} />;
            case 'login': return <LoginOutlined style={{ fontSize: '16px', color: '#3ba55c' }} />;
            case 'security': return <SafetyCertificateOutlined style={{ fontSize: '16px', color: '#ed4245' }} />;
            case 'audit': return <EditOutlined style={{ fontSize: '16px', color: '#faa61a' }} />;
            default: return <UserOutlined style={{ fontSize: '16px' }} />;
        }
    };

    const getColor = (type: string) => {
        switch (type) {
            case 'system': return 'blue';
            case 'login': return 'green';
            case 'security': return 'red';
            case 'audit': return 'gold';
            default: return 'gray';
        }
    };

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card bordered={false} styles={{ body: { padding: '24px' } }}>
                <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h2 style={{ margin: 0, fontSize: 20 }}>系统审计日志</h2>
                    <Space>
                        <DatePicker.RangePicker />
                        <Button type="primary" style={{ backgroundColor: '#5865F2' }}>查询</Button>
                    </Space>
                </div>

                <Timeline mode="left">
                    {activities.map(item => (
                        <Timeline.Item
                            key={item.id}
                            color={getColor(item.type)}
                            dot={getIcon(item.type)}
                            label={<span style={{ color: 'rgba(255,255,255,0.45)' }}>{item.time}</span>}
                        >
                            <Card
                                size="small"
                                bordered={false}
                                style={{
                                    backgroundColor: 'rgba(255,255,255,0.04)',
                                    marginBottom: 16,
                                    borderRadius: 8
                                }}
                            >
                                <Space orientation="vertical" size={4} style={{ width: '100%' }}>
                                    <Space>
                                        <Tag color="geekblue">{item.user}</Tag>
                                        <Text style={{ color: '#dcddde' }}>{item.action}</Text>
                                        <Tag>{item.target}</Tag>
                                    </Space>
                                    <div style={{ fontSize: 12, color: 'rgba(255,255,255,0.3)' }}>
                                        IP: {item.ip}
                                    </div>
                                </Space>
                            </Card>
                        </Timeline.Item>
                    ))}
                </Timeline>
            </Card>
        </motion.div>
    );
};

export default Audit;
