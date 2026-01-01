/**
 * Payment Records Page
 * Display all payment records in the system
 */
import React, { useState, useEffect } from 'react';
import { Card, Table, Tag, Space, Button, Input, DatePicker, Select, message, Statistic, Row, Col } from 'antd';
import { SearchOutlined, ReloadOutlined, DownloadOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

interface PaymentRecord {
    id: number;
    orderNo: string;
    userId: number;
    userName?: string;
    amountCents: number;
    status: 'pending' | 'success' | 'failed' | 'refunded';
    paymentMethod: string;
    transactionId?: string;
    createdAt: string;
    updatedAt: string;
}

const PaymentRecords: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<PaymentRecord[]>([]);
    const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
    const [filters, setFilters] = useState<any>({});

    const fetchRecords = async (page = 1, pageSize = 10) => {
        setLoading(true);
        try {
            // Mock data for now since API might not be fully implemented
            const mockData: PaymentRecord[] = [
                {
                    id: 1,
                    orderNo: 'ESC20251231113618283671',
                    userId: 4,
                    userName: '用户 #4',
                    amountCents: 21900,
                    status: 'success',
                    paymentMethod: 'wechat',
                    transactionId: 'TX1234567890',
                    createdAt: '2025-12-31 11:36:18',
                    updatedAt: '2025-12-31 11:36:18',
                },
                {
                    id: 2,
                    orderNo: 'ESC20251231113618476417',
                    userId: 2,
                    userName: '用户 #2',
                    amountCents: 16900,
                    status: 'refunded',
                    paymentMethod: 'alipay',
                    transactionId: 'TX1234567891',
                    createdAt: '2025-12-31 11:36:18',
                    updatedAt: '2025-12-31 11:36:18',
                },
            ];

            setData(mockData);
            setPagination({ current: page, pageSize, total: mockData.length });
        } catch (error) {
            message.error('获取支付记录失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRecords();
    }, []);

    const getStatusTag = (status: string) => {
        const statusMap: Record<string, { color: string; text: string }> = {
            pending: { color: 'default', text: '待支付' },
            success: { color: 'success', text: '已支付' },
            failed: { color: 'error', text: '支付失败' },
            refunded: { color: 'warning', text: '已退款' },
        };
        const { color, text } = statusMap[status] || { color: 'default', text: status };
        return <Tag color={color}>{text}</Tag>;
    };

    const getPaymentMethodTag = (method: string) => {
        const methodMap: Record<string, { text: string }> = {
            wechat: { text: '微信支付' },
            alipay: { text: '支付宝' },
            balance: { text: '余额支付' },
        };
        const { text } = methodMap[method] || { text: method };
        return <span>{text}</span>;
    };

    const columns: ColumnsType<PaymentRecord> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            width: 200,
            render: (text: string) => (
                <Space>
                    <span>{text}</span>
                    <Button type="link" size="small" icon={<SearchOutlined />}>查看</Button>
                </Space>
            ),
        },
        {
            title: '用户',
            dataIndex: 'userName',
            key: 'userName',
            width: 120,
        },
        {
            title: '金额',
            dataIndex: 'amountCents',
            key: 'amountCents',
            width: 120,
            render: (amount: number) => `¥${(amount / 100).toFixed(2)}`,
        },
        {
            title: '支付方式',
            dataIndex: 'paymentMethod',
            key: 'paymentMethod',
            width: 120,
            render: (method: string) => getPaymentMethodTag(method),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => getStatusTag(status),
        },
        {
            title: '交易流水号',
            dataIndex: 'transactionId',
            key: 'transactionId',
            width: 180,
            render: (text: string) => text || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
        },
    ];

    return (
        <div style={{ padding: '24px' }}>
            <Row gutter={16} style={{ marginBottom: '24px' }}>
                <Col span={6}>
                    <Statistic title="总支付记录" value={data.length} />
                </Col>
                <Col span={6}>
                    <Statistic title="今日支付" value="0" />
                </Col>
                <Col span={6}>
                    <Statistic title="今日金额" value="¥0.00" />
                </Col>
                <Col span={6}>
                    <Statistic title="成功率" value="100%" />
                </Col>
            </Row>

            <Card
                title="支付记录"
                extra={
                    <Space>
                        <Button icon={<ReloadOutlined />} onClick={() => fetchRecords(pagination.current, pagination.pageSize)}>
                            刷新
                        </Button>
                        <Button icon={<DownloadOutlined />}>
                            导出
                        </Button>
                    </Space>
                }
            >
                <div style={{ marginBottom: '16px' }}>
                    <Space>
                        <Input placeholder="搜索订单号" prefix={<SearchOutlined />} style={{ width: 200 }} />
                        <Select placeholder="支付方式" style={{ width: 120 }} allowClear>
                            <Select.Option value="wechat">微信支付</Select.Option>
                            <Select.Option value="alipay">支付宝</Select.Option>
                            <Select.Option value="balance">余额支付</Select.Option>
                        </Select>
                        <Select placeholder="状态" style={{ width: 120 }} allowClear>
                            <Select.Option value="pending">待支付</Select.Option>
                            <Select.Option value="success">已支付</Select.Option>
                            <Select.Option value="failed">支付失败</Select.Option>
                            <Select.Option value="refunded">已退款</Select.Option>
                        </Select>
                        <RangePicker style={{ width: 280 }} />
                        <Button type="primary" icon={<SearchOutlined />}>搜索</Button>
                        <Button>重置</Button>
                    </Space>
                </div>

                <Table
                    columns={columns}
                    dataSource={data}
                    loading={loading}
                    rowKey="id"
                    pagination={{
                        current: pagination.current,
                        pageSize: pagination.pageSize,
                        total: pagination.total,
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total) => `共 ${total} 条`,
                        onChange: (page, pageSize) => fetchRecords(page, pageSize),
                    }}
                />
            </Card>
        </div>
    );
};

export default PaymentRecords;
