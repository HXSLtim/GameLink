/**
 * Payment Records Page
 * Display all payment records in the system
 */
import React, { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import { Card, Tag, Space, Statistic, Row, Col, App, Button, Drawer, Table, Typography } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { PageContainer, SearchTable, type ToolbarButton, type SearchField } from '@/components';
import { 
    adminApi, 
    type Payment, 
    type PaymentStatus, 
    type PaymentMethod,
    type PaymentQueryParams,
    type RefundRecord,
    type RefundStatus,
} from '@/api/admin';

const PaymentRecords: React.FC = () => {
    const { message } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<Payment[]>([]);
    const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
    
    const searchParamsRef = useRef<Record<string, unknown>>({});

    const [refundDrawerOpen, setRefundDrawerOpen] = useState(false);
    const [refundLoading, setRefundLoading] = useState(false);
    const [refundRecords, setRefundRecords] = useState<RefundRecord[]>([]);
    const [refundPayment, setRefundPayment] = useState<Payment | null>(null);

    // 统计数据
    const [stats, setStats] = useState({
        total: 0,
        todayCount: 0,
        todayAmount: 0,
        successRate: 100,
    });

    const fetchRecords = useCallback(async (page = 1, pageSize = 10, params?: Record<string, unknown>) => {
        setLoading(true);
        try {
            const mergedParams = params ?? searchParamsRef.current;
            const keyword = typeof mergedParams.keyword === 'string' ? mergedParams.keyword.trim() : '';
            const method = mergedParams.method as PaymentMethod | undefined;
            const status = mergedParams.status as PaymentStatus | undefined;
            const dateRange = mergedParams.dateRange as [dayjs.Dayjs, dayjs.Dayjs] | undefined;

            const queryParams: PaymentQueryParams = {
                page,
                pageSize,
                keyword: keyword || undefined,
                status,
                method,
                dateFrom: dateRange?.[0]?.format('YYYY-MM-DD'),
                dateTo: dateRange?.[1]?.format('YYYY-MM-DD'),
            };

            const response = await adminApi.getPayments(queryParams);
            if (response.data.success) {
                const payments = response.data.data || [];
                setData(payments);
                
                // 从响应中获取分页信息
                const paginationData = response.data.pagination;
                setPagination({ 
                    current: paginationData?.page || page, 
                    pageSize: paginationData?.page_size || pageSize, 
                    total: paginationData?.total || payments.length,
                });

                // 计算统计数据
                const today = dayjs().format('YYYY-MM-DD');
                const todayPayments = payments.filter(p => p.createdAt?.startsWith(today));
                const successPayments = payments.filter(p => p.status === 'paid');
                
                setStats({
                    total: paginationData?.total || payments.length,
                    todayCount: todayPayments.length,
                    todayAmount: todayPayments.reduce((sum, p) => sum + (p.amountCents || 0), 0),
                    successRate: payments.length > 0 
                        ? Math.round((successPayments.length / payments.length) * 100) 
                        : 100,
                });
            }
        } catch {
            message.error('获取支付记录失败');
        } finally {
            setLoading(false);
        }
    }, [message]);

    useEffect(() => {
        fetchRecords();
    }, [fetchRecords]);

    const handleSearch = (values: Record<string, unknown>) => {
        searchParamsRef.current = values;
        fetchRecords(1, pagination.pageSize, values);
    };

    const handleExport = useCallback(() => {
        message.info('导出功能开发中');
    }, [message]);

    const getStatusTag = (status: PaymentStatus) => {
        const statusMap: Record<PaymentStatus, { color: string; text: string }> = {
            pending: { color: 'default', text: '待支付' },
            paid: { color: 'success', text: '已支付' },
            failed: { color: 'error', text: '支付失败' },
            refunded: { color: 'warning', text: '已退款' },
            canceled: { color: 'default', text: '已取消' },
        };
        const { color, text } = statusMap[status] || { color: 'default', text: status };
        return <Tag color={color}>{text}</Tag>;
    };

    const getPaymentMethodTag = (method: PaymentMethod) => {
        const methodMap: Record<PaymentMethod, { text: string }> = {
            wechat: { text: '微信支付' },
            alipay: { text: '支付宝' },
            balance: { text: '余额支付' },
            bank: { text: '银行卡' },
        };
        const { text } = methodMap[method] || { text: method };
        return <span>{text}</span>;
    };

    const getRefundStatusTag = (status: RefundStatus) => {
        const statusMap: Record<RefundStatus, { color: string; text: string }> = {
            pending: { color: 'default', text: '处理中' },
            processed: { color: 'success', text: '已完成' },
            failed: { color: 'error', text: '失败' },
        };
        const { color, text } = statusMap[status] || { color: 'default', text: status };
        return <Tag color={color}>{text}</Tag>;
    };

    const totalRefundedCents = useMemo(
        () => refundRecords.reduce((sum, r) => sum + (r.amountCents || 0), 0),
        [refundRecords],
    );

    const openRefundDrawer = useCallback(async (payment: Payment) => {
        setRefundPayment(payment);
        setRefundDrawerOpen(true);
        setRefundLoading(true);
        try {
            const res = await adminApi.getPaymentRefunds(payment.id);
            if (res.data.success) {
                setRefundRecords(res.data.data || []);
            } else {
                setRefundRecords([]);
            }
        } catch {
            setRefundRecords([]);
            message.error('获取退款记录失败');
        } finally {
            setRefundLoading(false);
        }
    }, [message]);

    const closeRefundDrawer = useCallback(() => {
        setRefundDrawerOpen(false);
        setRefundRecords([]);
        setRefundPayment(null);
        setRefundLoading(false);
    }, []);

    const columns: ColumnsType<Payment> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '订单',
            key: 'order',
            width: 200,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <span>订单ID: {record.orderId}</span>
                    {record.order?.orderNo && (
                        <span style={{ fontSize: 12, color: '#999' }}>{record.order.orderNo}</span>
                    )}
                </Space>
            ),
        },
        {
            title: '用户',
            key: 'user',
            width: 120,
            render: (_, record) => record.user?.name || `用户 #${record.userId}`,
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
            dataIndex: 'method',
            key: 'method',
            width: 120,
            render: (method: PaymentMethod) => getPaymentMethodTag(method),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: PaymentStatus) => getStatusTag(status),
        },
        {
            title: '交易流水号',
            dataIndex: 'providerTradeNo',
            key: 'providerTradeNo',
            width: 180,
            render: (text: string) => text || '-',
        },
        {
            title: '支付时间',
            dataIndex: 'paidAt',
            key: 'paidAt',
            width: 180,
            render: (text: string) => text || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
        },
        {
            title: '操作',
            key: 'actions',
            width: 120,
            render: (_, record) => {
                const hasRefund = (record.refundedAmountCents || 0) > 0 || record.status === 'refunded';
                return (
                    <Space size={8}>
                        <Button
                            type="link"
                            size="small"
                            disabled={!hasRefund}
                            onClick={() => openRefundDrawer(record)}
                        >
                            退款记录
                        </Button>
                    </Space>
                );
            },
        },
    ];

    const refundColumns: ColumnsType<RefundRecord> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '金额',
            dataIndex: 'amountCents',
            key: 'amountCents',
            width: 120,
            render: (amount: number) => `¥${(amount / 100).toFixed(2)}`,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: RefundStatus) => getRefundStatusTag(status),
        },
        {
            title: '原因',
            dataIndex: 'reason',
            key: 'reason',
            width: 220,
            render: (text?: string) => text || '-',
        },
        {
            title: '流水号',
            dataIndex: 'providerTradeNo',
            key: 'providerTradeNo',
            width: 220,
            render: (text?: string) => text || '-',
        },
        {
            title: '退款时间',
            dataIndex: 'refundedAt',
            key: 'refundedAt',
            width: 180,
            render: (text?: string) => text || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
        },
    ];

    const searchFields: SearchField[] = useMemo(() => ([
        {
            name: 'keyword',
            label: '关键词',
            type: 'input',
            placeholder: '搜索订单号/流水号',
        },
        {
            name: 'method',
            label: '支付方式',
            type: 'select',
            options: [
                { label: '微信支付', value: 'wechat' },
                { label: '支付宝', value: 'alipay' },
                { label: '余额支付', value: 'balance' },
                { label: '银行卡', value: 'bank' },
            ],
        },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: [
                { label: '待支付', value: 'pending' },
                { label: '已支付', value: 'paid' },
                { label: '支付失败', value: 'failed' },
                { label: '已退款', value: 'refunded' },
                { label: '已取消', value: 'canceled' },
            ],
        },
        {
            name: 'dateRange',
            label: '支付时间',
            type: 'dateRange',
        },
    ]), []);

    const toolbarButtons: ToolbarButton[] = useMemo(() => ([
        {
            text: '导出',
            icon: <DownloadOutlined />,
            onClick: handleExport,
        },
    ]), [handleExport]);

    return (
        <PageContainer title="支付记录" subTitle="查看平台支付流水与统计">
            <Row gutter={16} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card
                        style={{
                            border: 'none',
                            borderRadius: 12,
                            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.03)',
                            transition: 'all 0.2s ease-in-out'
                        }}
                        bodyStyle={{ padding: '24px' }}
                    >
                        <Statistic
                            title="总支付记录"
                            value={stats.total}
                            valueStyle={{ fontWeight: 600 }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card
                        style={{
                            border: 'none',
                            borderRadius: 12,
                            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.03)',
                            transition: 'all 0.2s ease-in-out'
                        }}
                        bodyStyle={{ padding: '24px' }}
                    >
                        <Statistic
                            title="今日支付"
                            value={stats.todayCount}
                            suffix="笔"
                            valueStyle={{ fontWeight: 600 }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card
                        style={{
                            border: 'none',
                            borderRadius: 12,
                            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.03)',
                            transition: 'all 0.2s ease-in-out'
                        }}
                        bodyStyle={{ padding: '24px' }}
                    >
                        <Statistic
                            title="今日金额"
                            value={stats.todayAmount / 100}
                            precision={2}
                            prefix="¥"
                            valueStyle={{ fontWeight: 600, color: '#52c41a' }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card
                        style={{
                            border: 'none',
                            borderRadius: 12,
                            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.03)',
                            transition: 'all 0.2s ease-in-out'
                        }}
                        bodyStyle={{ padding: '24px' }}
                    >
                        <Statistic
                            title="成功率"
                            value={stats.successRate}
                            suffix="%"
                            valueStyle={{
                                fontWeight: 600,
                                color: stats.successRate >= 90 ? '#52c41a' : stats.successRate >= 70 ? '#faad14' : '#ff4d4f'
                            }}
                        />
                    </Card>
                </Col>
            </Row>
            <SearchTable<Payment>
                columns={columns}
                dataSource={data}
                loading={loading}
                searchFields={searchFields}
                toolbarButtons={toolbarButtons}
                showCreate={false}
                onSearch={handleSearch}
                onRefresh={() => fetchRecords(pagination.current, pagination.pageSize)}
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
            <Drawer
                title="退款记录"
                open={refundDrawerOpen}
                width={920}
                onClose={closeRefundDrawer}
                destroyOnClose
            >
                <Space direction="vertical" size={16} style={{ width: '100%' }}>
                    <Card size="small">
                        <Space direction="vertical" size={4}>
                            <Typography.Text type="secondary">
                                支付ID：{refundPayment?.id ?? '-'}
                            </Typography.Text>
                            <Typography.Text type="secondary">
                                订单ID：{refundPayment?.orderId ?? '-'}
                            </Typography.Text>
                            <Typography.Text>
                                已退款合计：
                                <Typography.Text strong style={{ marginLeft: 8 }}>
                                    ¥{(totalRefundedCents / 100).toFixed(2)}
                                </Typography.Text>
                            </Typography.Text>
                        </Space>
                    </Card>
                    <Table<RefundRecord>
                        columns={refundColumns}
                        dataSource={refundRecords}
                        loading={refundLoading}
                        rowKey="id"
                        pagination={false}
                        locale={{ emptyText: '暂无退款记录' }}
                        scroll={{ x: 1000 }}
                    />
                </Space>
            </Drawer>
        </PageContainer>
    );
};

export default PaymentRecords;
