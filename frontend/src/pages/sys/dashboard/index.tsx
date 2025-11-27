import React from 'react';
import { Card, Row, Col, Statistic, List, Avatar, Tag } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined, UserOutlined, TrophyOutlined, ShoppingOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';

const Dashboard: React.FC = () => {
    const containerVariants = {
        hidden: { opacity: 0 },
        visible: {
            opacity: 1,
            transition: {
                staggerChildren: 0.1
            }
        }
    };

    const itemVariants = {
        hidden: { y: 20, opacity: 0 },
        visible: {
            y: 0,
            opacity: 1
        }
    };

    return (
        <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
        >
            <h2 style={{ color: '#fff', marginBottom: 24 }}>仪表盘概览</h2>

            <Row gutter={[16, 16]}>
                <Col xs={24} sm={12} lg={6}>
                    <motion.div variants={itemVariants}>
                        <Card bordered={false} hoverable>
                            <Statistic
                                title="总用户数"
                                value={12345}
                                precision={0}
                                styles={{ content: { color: '#5865F2' } }}
                                prefix={<UserOutlined />}
                                suffix="人"
                            />
                            <div style={{ color: '#3ba55c', marginTop: 8, fontSize: 12 }}>
                                <ArrowUpOutlined /> 较上周增长 5.2%
                            </div>
                        </Card>
                    </motion.div>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <motion.div variants={itemVariants}>
                        <Card bordered={false} hoverable>
                            <Statistic
                                title="活跃游戏"
                                value={48}
                                styles={{ content: { color: '#3ba55c' } }}
                                prefix={<TrophyOutlined />}
                                suffix="款"
                            />
                            <div style={{ color: '#3ba55c', marginTop: 8, fontSize: 12 }}>
                                <ArrowUpOutlined /> 新增 2 款
                            </div>
                        </Card>
                    </motion.div>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <motion.div variants={itemVariants}>
                        <Card bordered={false} hoverable>
                            <Statistic
                                title="今日营收"
                                value={1240}
                                precision={2}
                                styles={{ content: { color: '#faa61a' } }}
                                prefix="¥"
                            />
                            <div style={{ color: '#3ba55c', marginTop: 8, fontSize: 12 }}>
                                <ArrowUpOutlined /> 增长 12%
                            </div>
                        </Card>
                    </motion.div>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <motion.div variants={itemVariants}>
                        <Card bordered={false} hoverable>
                            <Statistic
                                title="待处理订单"
                                value={15}
                                styles={{ content: { color: '#ed4245' } }}
                                prefix={<ShoppingOutlined />}
                                suffix="单"
                            />
                            <div style={{ color: '#ed4245', marginTop: 8, fontSize: 12 }}>
                                <ArrowDownOutlined /> 较昨日减少 3 单
                            </div>
                        </Card>
                    </motion.div>
                </Col>
            </Row>

            <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
                <Col xs={24} lg={16}>
                    <motion.div variants={itemVariants}>
                        <Card title="活动趋势" bordered={false} style={{ height: 400 }}>
                            <div style={{
                                height: 300, display: 'flex', alignItems: 'center', justifyContent: 'center',
                                color: 'rgba(255,255,255,0.3)', border: '2px dashed rgba(255,255,255,0.1)',
                                borderRadius: 8
                            }}>
                                图表可视化占位符
                            </div>
                        </Card>
                    </motion.div>
                </Col>
                <Col xs={24} lg={8}>
                    <motion.div variants={itemVariants}>
                        <Card title="近期工单" bordered={false} style={{ height: 400, overflow: 'auto' }}>
                            <List
                                itemLayout="horizontal"
                                dataSource={[
                                    { title: '登录问题', user: 'User#1234', time: '2分钟前', status: '待处理' },
                                    { title: '支付失败', user: 'Gamer#9999', time: '15分钟前', status: '紧急' },
                                    { title: 'Bug 反馈', user: 'Dev#0000', time: '1小时前', status: '已解决' },
                                    { title: '账号找回', user: 'Lost#1111', time: '2小时前', status: '待处理' },
                                ]}
                                renderItem={(item) => (
                                    <List.Item>
                                        <List.Item.Meta
                                            avatar={<Avatar style={{ backgroundColor: '#5865F2' }}>{item.user[0]}</Avatar>}
                                            title={<span style={{ color: '#fff' }}>{item.title}</span>}
                                            description={
                                                <span style={{ color: 'rgba(255,255,255,0.5)', fontSize: 12 }}>
                                                    <UserOutlined style={{ marginRight: 4 }} />{item.user}
                                                    <ClockCircleOutlined style={{ marginLeft: 8, marginRight: 4 }} />{item.time}
                                                </span>
                                            }
                                        />
                                        <Tag color={item.status === '紧急' ? 'red' : item.status === '已解决' ? 'green' : 'blue'}>
                                            {item.status}
                                        </Tag>
                                    </List.Item>
                                )}
                            />
                        </Card>
                    </motion.div>
                </Col>
            </Row>
        </motion.div>
    );
};

export default Dashboard;
