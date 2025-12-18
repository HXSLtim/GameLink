/**
 * 提现分流统计页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Card,
    Row,
    Col,
    Statistic,
    DatePicker,
    message,
    Spin,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    BankOutlined,
    DollarOutlined,
    DownloadOutlined,
    FileTextOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

interface WithdrawByCompany {
    id: number;
    playerId: number;
    playerName: string;
    amount: number;
    status: string;
    settlementCompanyId: number;
    settlementCompanyName: string;
    createdAt: string;
}

interface RoutingStats {
    totalAmount: number;
    totalCount: number;
    completedAmount: number;
    completedCount: number;
    pendingAmount: number;
    pendingCount: number;
    byCompany: { companyId: number; companyName: string; amount: number; count: number }[];
}

const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'processing', text: '待审核' },
    approved: { color: 'warning', text: '已批准' },
    rejected: { color: 'error', text: '已拒绝' },
    completed: { color: 'success', text: '已完成' },
    failed: { color: 'error', text: '失败' },
};

const exportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'playerName', title: '陪玩师' },
    { key: 'amount', title: '金额' },
    { key: 'settlementCompanyName', title: '结算公司' },
    { key: 'status', title: '状态', render: (v) => statusMap[v as string]?.text || String(v) },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];


const WithdrawRoutingPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [withdraws, setWithdraws] = useState<WithdrawByCompany[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    const [statsLoading, setStatsLoading] = useState(false);
    const [stats, setStats] = useState<RoutingStats | null>(null);
    const [dateRange, setDateRange] = useState<[dayjs.Dayjs | null, dayjs.Dayjs | null]>([null, null]);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = { page: current, pageSize, ...searchParams };
            if (dateRange[0]) params.dateFrom = dateRange[0].format('YYYY-MM-DD');
            if (dateRange[1]) params.dateTo = dateRange[1].format('YYYY-MM-DD');

            const response = await apiClient.get('/admin/withdrawals/by-company', { params });
            if (response.data.success) {
                setWithdraws(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            }
        } catch (error) {
            console.error('Load error:', error);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, dateRange]);

    const loadStats = useCallback(async () => {
        setStatsLoading(true);
        try {
            const params: Record<string, unknown> = {};
            if (dateRange[0]) params.dateFrom = dateRange[0].format('YYYY-MM-DD');
            if (dateRange[1]) params.dateTo = dateRange[1].format('YYYY-MM-DD');

            const response = await apiClient.get('/admin/withdrawals/routing-stats', { params });
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (error) {
            console.error('Load stats error:', error);
        } finally {
            setStatsLoading(false);
        }
    }, [dateRange]);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleExport = () => {
        exportToCSV(withdraws as unknown as Record<string, unknown>[], exportColumns, 'withdraw_routing');
        message.success('导出成功');
    };

    const handleGenerateReport = async (reportType: string) => {
        try {
            const now = dayjs();
            const params: Record<string, unknown> = { reportType, year: now.year() };
            if (reportType === 'monthly') params.month = now.month() + 1;
            if (reportType === 'quarterly') params.quarter = Math.ceil((now.month() + 1) / 3);

            const response = await apiClient.get('/admin/withdrawals/routing-report', { params });
            if (response.data.success) {
                message.success('报表生成成功');
                // 可以下载或展示报表
            }
        } catch {
            message.error('生成报表失败');
        }
    };

    const searchFields: SearchField[] = [
        { name: 'settlementCompanyId', label: '结算公司ID', type: 'input', placeholder: '公司ID' },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: [
                { label: '待审核', value: 'pending' },
                { label: '已批准', value: 'approved' },
                { label: '已完成', value: 'completed' },
                { label: '已拒绝', value: 'rejected' },
            ],
        },
    ];

    const columns: ColumnsType<WithdrawByCompany> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        { title: '陪玩师', dataIndex: 'playerName', key: 'playerName', width: 120 },
        {
            title: '金额',
            dataIndex: 'amount',
            key: 'amount',
            width: 120,
            render: (amount: number) => <span style={{ color: '#f5222d' }}>¥{amount.toFixed(2)}</span>,
        },
        { title: '结算公司', dataIndex: 'settlementCompanyName', key: 'settlementCompanyName', width: 150 },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => {
                const info = statusMap[status];
                return info ? <Tag color={info.color}>{info.text}</Tag> : status;
            },
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        { text: '月报', icon: <FileTextOutlined />, needSelection: false, onClick: () => handleGenerateReport('monthly') },
        { text: '导出数据', icon: <DownloadOutlined />, needSelection: false, onClick: () => handleExport() },
    ];


    return (
        <PageContainer title="提现分流统计" subTitle="按结算公司统计提现数据">
            <Spin spinning={statsLoading}>
                <Row gutter={16} style={{ marginBottom: 16 }}>
                    <Col span={6}>
                        <Card>
                            <Statistic title="总提现金额" value={stats?.totalAmount || 0} precision={2} prefix={<DollarOutlined />} suffix="元" />
                        </Card>
                    </Col>
                    <Col span={6}>
                        <Card>
                            <Statistic title="总提现笔数" value={stats?.totalCount || 0} prefix={<BankOutlined />} />
                        </Card>
                    </Col>
                    <Col span={6}>
                        <Card>
                            <Statistic title="已完成金额" value={stats?.completedAmount || 0} precision={2} suffix="元" />
                        </Card>
                    </Col>
                    <Col span={6}>
                        <Card>
                            <Statistic title="待处理金额" value={stats?.pendingAmount || 0} precision={2} suffix="元" />
                        </Card>
                    </Col>
                </Row>

                {/* 按公司统计 */}
                {stats?.byCompany && stats.byCompany.length > 0 && (
                    <Card title="按结算公司统计" style={{ marginBottom: 16 }}>
                        <Row gutter={16}>
                            {stats.byCompany.map((company, index) => (
                                <Col span={6} key={index}>
                                    <Card size="small">
                                        <Statistic
                                            title={company.companyName}
                                            value={company.amount}
                                            precision={2}
                                            suffix={`元 (${company.count}笔)`}
                                        />
                                    </Card>
                                </Col>
                            ))}
                        </Row>
                    </Card>
                )}
            </Spin>

            <Card
                title="提现明细"
                extra={
                    <Space>
                        <DatePicker.RangePicker
                            value={dateRange}
                            onChange={(dates) => setDateRange(dates as [dayjs.Dayjs | null, dayjs.Dayjs | null])}
                            allowClear
                        />
                    </Space>
                }
            >
                <SearchTable
                    columns={columns}
                    dataSource={withdraws}
                    rowKey="id"
                    searchFields={searchFields}
                    onSearch={handleSearch}
                    onRefresh={() => { loadData(); loadStats(); }}
                    loading={loading}
                    showCreate={false}
                    toolbarButtons={toolbarButtons}
                    pagination={{ current, pageSize, total, showSizeChanger: true, showTotal: t => `共 ${t} 条`, onChange: (p, s) => { setCurrent(p); setPageSize(s); } }}
                    scroll={{ x: 800 }}
                />
            </Card>
        </PageContainer>
    );
};

export default WithdrawRoutingPage;
