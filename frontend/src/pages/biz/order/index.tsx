import React, { useState } from 'react';
import { Table, Tag, Space, Card, Input, DatePicker, Button, Tooltip } from 'antd';
import { SearchOutlined, EyeOutlined, DownloadOutlined, FilterOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { motion } from 'framer-motion';

interface Order {
    id: string;
    orderNo: string;
    user: string;
    companion: string;
    game: string;
    amount: number;
    status: 'completed' | 'pending' | 'cancelled' | 'processing';
    createdAt: string;
}

const Orders: React.FC = () => {
    const [searchText, setSearchText] = useState('');

    const orders: Order[] = [
        { id: '1', orderNo: 'ORD-20231125-001', user: 'GamerOne', companion: 'ProPlayer', game: '英雄联盟', amount: 50.00, status: 'completed', createdAt: '2023-11-25 14:30' },
        { id: '2', orderNo: 'ORD-20231125-002', user: 'Newbie', companion: 'ProPlayer', game: '无畏契约', amount: 35.00, status: 'pending', createdAt: '2023-11-25 15:45' },
        { id: '3', orderNo: 'ORD-20231124-005', user: 'RichGuy', companion: 'StarGamer', game: '原神', amount: 200.00, status: 'processing', createdAt: '2023-11-24 09:12' },
        { id: '4', orderNo: 'ORD-20231124-001', user: 'Troll', companion: 'ProPlayer', game: 'Apex 英雄', amount: 45.00, status: 'cancelled', createdAt: '2023-11-24 08:00' },
    ];

    const columns: ColumnsType<Order> = [
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            render: text => <span style={{ fontFamily: 'monospace' }}>{text}</span>,
        },
        {
            title: '下单用户',
            dataIndex: 'user',
            key: 'user',
        },
        {
            title: '陪玩',
            dataIndex: 'companion',
            key: 'companion',
            render: text => <span style={{ color: '#5865F2' }}>@{text}</span>,
        },
        {
            title: '游戏',
            dataIndex: 'game',
            key: 'game',
        },
        {
            title: '金额',
            dataIndex: 'amount',
            key: 'amount',
            render: amount => <span style={{ fontWeight: 'bold' }}>¥{amount.toFixed(2)}</span>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: status => {
                let color = 'default';
                let text = '未知';
                switch (status) {
                    case 'completed': color = 'success'; text = '已完成'; break;
                    case 'pending': color = 'warning'; text = '待支付'; break;
                    case 'processing': color = 'processing'; text = '进行中'; break;
                    case 'cancelled': color = 'error'; text = '已取消'; break;
                }
                return <Tag color={color}>{text}</Tag>;
            },
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: text => <span style={{ color: 'rgba(255,255,255,0.45)' }}>{text}</span>,
        },
        {
            title: '操作',
            key: 'action',
            render: () => (
                <Space size="middle">
                    <Tooltip title="查看详情">
                        <Button type="text" icon={<EyeOutlined />} style={{ color: '#fff' }} />
                    </Tooltip>
                </Space>
            ),
        },
    ];

    const filteredOrders = orders.filter(order =>
        order.orderNo.toLowerCase().includes(searchText.toLowerCase()) ||
        order.user.toLowerCase().includes(searchText.toLowerCase())
    );

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card bordered={false} bodyStyle={{ padding: '24px' }}>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h2 style={{ margin: 0, fontSize: 20 }}>订单管理</h2>
                    <Space>
                        <DatePicker.RangePicker style={{ width: 250 }} />
                        <Input
                            placeholder="搜索订单号/用户..."
                            prefix={<SearchOutlined />}
                            onChange={e => setSearchText(e.target.value)}
                            style={{ width: 200 }}
                        />
                        <Button icon={<FilterOutlined />}>筛选</Button>
                        <Button type="primary" icon={<DownloadOutlined />} style={{ backgroundColor: '#5865F2' }}>导出</Button>
                    </Space>
                </div>
                <Table
                    columns={columns}
                    dataSource={filteredOrders}
                    rowKey="id"
                    pagination={{ pageSize: 10 }}
                />
            </Card>
        </motion.div>
    );
};

export default Orders;
