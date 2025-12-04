import React, { useState } from 'react';
import { Card, List, Button, Tag, Space, Input, Modal, Form, Select } from 'antd';
import { PlusOutlined, SearchOutlined, TrophyOutlined, FireOutlined, ThunderboltOutlined, CrownOutlined } from '@ant-design/icons';
import { motion } from 'framer-motion';

interface Game {
    id: number;
    name: string;
    icon: React.ReactNode;
    category: string;
    status: 'active' | 'maintenance';
    players: number;
}

const Games: React.FC = () => {
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [searchText, setSearchText] = useState('');

    const games: Game[] = [
        { id: 1, name: '英雄联盟', icon: <TrophyOutlined />, category: 'MOBA', status: 'active', players: 12050 },
        { id: 2, name: '无畏契约', icon: <FireOutlined />, category: 'FPS', status: 'active', players: 8500 },
        { id: 3, name: 'Apex 英雄', icon: <ThunderboltOutlined />, category: '大逃杀', status: 'maintenance', players: 0 },
        { id: 4, name: '原神', icon: <CrownOutlined />, category: 'RPG', status: 'active', players: 15000 },
    ];

    const filteredGames = games.filter(game => game.name.includes(searchText));

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2 style={{ color: '#fff', margin: 0 }}>游戏库管理</h2>
                <Space>
                    <Input
                        placeholder="搜索游戏..."
                        prefix={<SearchOutlined />}
                        onChange={e => setSearchText(e.target.value)}
                        style={{ width: 250 }}
                    />
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)} style={{ backgroundColor: '#3ba55c', borderColor: '#3ba55c' }}>
                        添加游戏
                    </Button>
                </Space>
            </div>

            <List
                grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
                dataSource={filteredGames}
                renderItem={item => (
                    <List.Item>
                        <Card
                            hoverable
                            variant="borderless"
                            actions={[
                                <Button type="text" key="edit">编辑</Button>,
                                <Button type="text" danger key="delete">下架</Button>
                            ]}
                        >
                            <Card.Meta
                                avatar={
                                    <div style={{
                                        width: 48, height: 48, borderRadius: 8,
                                        backgroundColor: '#5865F2', display: 'flex',
                                        alignItems: 'center', justifyContent: 'center',
                                        fontSize: 24, color: '#fff'
                                    }}>
                                        {item.icon}
                                    </div>
                                }
                                title={<span style={{ color: '#fff' }}>{item.name}</span>}
                                description={
                                    <Space orientation="vertical" size={0}>
                                        <Tag color="blue">{item.category}</Tag>
                                        <div style={{ marginTop: 8, fontSize: 12 }}>
                                            <span style={{
                                                display: 'inline-block', width: 8, height: 8,
                                                borderRadius: '50%', backgroundColor: item.status === 'active' ? '#3ba55c' : '#faa61a',
                                                marginRight: 6
                                            }} />
                                            {item.status === 'active' ? '运营中' : '维护中'}
                                        </div>
                                    </Space>
                                }
                            />
                        </Card>
                    </List.Item>
                )}
            />

            <Modal
                title="添加新游戏"
                open={isModalVisible}
                onOk={() => setIsModalVisible(false)}
                onCancel={() => setIsModalVisible(false)}
            >
                <Form layout="vertical">
                    <Form.Item label="游戏名称" name="name" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item label="分类" name="category" rules={[{ required: true }]}>
                        <Select>
                            <Select.Option value="moba">MOBA</Select.Option>
                            <Select.Option value="fps">FPS</Select.Option>
                            <Select.Option value="rpg">RPG</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </motion.div>
    );
};

export default Games;
