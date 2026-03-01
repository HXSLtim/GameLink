/**
 * 对账管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    message,
    Card,
    Statistic,
    Popconfirm,
    Tooltip,
    Drawer,
    Descriptions,
    Table,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    CheckOutlined,
    CloseOutlined,
    EyeOutlined,
    SyncOutlined,
    FileTextOutlined,
    ExclamationCircleOutlined,
    ClockCircleOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { AmountDisplay } from '@/components/AmountDisplay';
import { reconciliationApi } from '@/api/reconciliation';
import type {
    Reconciliation,
    ReconciliationType,
    ReconciliationStatus,
    ReconciliationDetail,
    ReconciliationListParams,
} from '@/api/reconciliation';
import {
    RECONCILIATION_STATUS_TEXT,
    RECONCILIATION_STATUS_COLOR,
    RECONCILIATION_TYPE_TEXT,
    RECONCILIATION_TYPE_COLOR,
} from '@/types/reconciliation';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';

/**
 * 对账类型映射（兼容旧代码，建议使用 RECONCILIATION_TYPE_*.tsx）
 */
const reconciliationTypeMap: Record<ReconciliationType, { text: string; color: string }> = {
    payment: { text: '支付对账', color: 'blue' },
    internal: { text: '内部对账', color: 'green' },
    bank: { text: '银行对账', color: 'purple' },
    manual: { text: '手工对账', color: 'orange' },
};

/**
 * 对账状态映射
 */
const reconciliationStatusMap: Record<ReconciliationStatus, { text: string; color: string; icon: React.ReactNode }> = {
    pending: { text: '待对账', color: 'default', icon: <ClockCircleOutlined /> },
    progress: { text: '对账中', color: 'processing', icon: <SyncOutlined spin /> },
    success: { text: '对账成功', color: 'success', icon: <CheckOutlined /> },
    failed: { text: '对账失败', color: 'error', icon: <CloseOutlined /> },
    exception: { text: '异常', color: 'warning', icon: <ExclamationCircleOutlined /> },
};

/**
 * 对账管理页面
 */
const ReconciliationPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [reconciliations, setReconciliations] = useState<Reconciliation[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<ReconciliationListParams>({});

    // 弹窗/抽屉状态
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [currentReconciliation, setCurrentReconciliation] = useState<Reconciliation | null>(null);
    const [reconciliationDetails, setReconciliationDetails] = useState<ReconciliationDetail[]>([]);

    /**
     * 加载对账单数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                pageSize,
                ...searchParams,
            };
            const response = await reconciliationApi.getReconciliations(queryParams);
            if (response.data.success) {
                setReconciliations(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load reconciliations error:', error);
            message.error('加载对账单列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        const params: ReconciliationListParams = {};
        if (values.type) params.type = values.type as ReconciliationType;
        if (values.status) params.status = values.status as ReconciliationStatus;
        if (values.dateRange) {
            const [start, end] = values.dateRange as [dayjs.Dayjs, dayjs.Dayjs];
            params.date_from = start.format('YYYY-MM-DD');
            params.date_to = end.format('YYYY-MM-DD');
        }
        setSearchParams(params);
        setCurrent(1);
    };

    /**
     * 查看对账详情
     */
    const handleViewDetail = async (record: Reconciliation) => {
        try {
            const response = await reconciliationApi.getReconciliationDetail(record.id);
            if (response.data.success) {
                setCurrentReconciliation(response.data.data);
                setReconciliationDetails(response.data.data.details || []);
                setDetailDrawerVisible(true);
            } else {
                message.error(response.data.message || '加载详情失败');
            }
        } catch (error) {
            logger.error('Load reconciliation detail error:', error);
            message.error('加载对账详情失败');
        }
    };

    /**
     * 执行对账
     */
    const handleExecute = async (record: Reconciliation, targetStatus: ReconciliationStatus) => {
        try {
            const response = await reconciliationApi.executeReconciliation(record.id, { status: targetStatus });
            if (response.data.success) {
                message.success('操作成功');
                loadData();
            } else {
                message.error(response.data.message || '操作失败');
            }
        } catch (error) {
            logger.error('Execute reconciliation error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        {
            name: 'type',
            label: '对账类型',
            type: 'select',
            options: Object.entries(reconciliationTypeMap).map(([key, val]) => ({
                label: val.text,
                value: key,
            })),
        },
        {
            name: 'status',
            label: '对账状态',
            type: 'select',
            options: Object.entries(reconciliationStatusMap).map(([key, val]) => ({
                label: val.text,
                value: key,
            })),
        },
        {
            name: 'dateRange',
            label: '对账日期',
            type: 'dateRange',
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Reconciliation> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '对账单号',
            dataIndex: 'reconciliationNo',
            key: 'reconciliationNo',
            width: 180,
            render: (no: string) => (
                <Tooltip title={no}>
                    <span style={{ fontFamily: 'monospace' }}>{no}</span>
                </Tooltip>
            ),
        },
        {
            title: '对账类型',
            dataIndex: 'type',
            key: 'type',
            width: 120,
            render: (type: ReconciliationType) => (
                <Tag color={reconciliationTypeMap[type].color}>
                    {reconciliationTypeMap[type].text}
                </Tag>
            ),
        },
        {
            title: '对账状态',
            dataIndex: 'status',
            key: 'status',
            width: 120,
            render: (status: ReconciliationStatus) => (
                <Tag
                    color={reconciliationStatusMap[status].color}
                    icon={reconciliationStatusMap[status].icon}
                >
                    {reconciliationStatusMap[status].text}
                </Tag>
            ),
        },
        {
            title: '对账期间',
            key: 'period',
            width: 200,
            render: (_, record) => (
                <div>
                    <div style={{ fontSize: 12, color: '#999' }}>起: {dayjs(record.periodStart).format('YYYY-MM-DD')}</div>
                    <div style={{ fontSize: 12, color: '#999' }}>止: {dayjs(record.periodEnd).format('YYYY-MM-DD')}</div>
                </div>
            ),
        },
        {
            title: '记录统计',
            key: 'records',
            width: 120,
            align: 'center',
            render: (_, record) => (
                <div>
                    <div>总计: {record.totalRecords || 0}</div>
                    <div style={{ fontSize: 12, color: '#52c41a' }}>匹配: {record.matchedRecords || 0}</div>
                </div>
            ),
        },
        {
            title: '差异金额',
            dataIndex: 'differenceAmount',
            key: 'differenceAmount',
            width: 120,
            align: 'right',
            render: (amount: number) => (
                <AmountDisplay
                    value={amount}
                    type={amount > 0 ? 'expense' : amount < 0 ? 'income' : 'default'}
                    showCurrency={false}
                    showSign
                />
            ),
        },
        {
            title: '摘要',
            dataIndex: 'abstract',
            key: 'abstract',
            width: 200,
            ellipsis: { showTitle: false },
            render: (text: string) => (
                <Tooltip title={text}>
                    <span>{text || '-'}</span>
                </Tooltip>
            ),
        },
        {
            title: '对账日期',
            dataIndex: 'reconciliationDate',
            key: 'reconciliationDate',
            width: 120,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size={4}>
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                    {record.status === 'pending' && (
                        <Popconfirm
                            title="确认开始对账？"
                            onConfirm={() => handleExecute(record, 'progress')}
                        >
                            <Button type="link" size="small" icon={<SyncOutlined />}>
                                开始
                            </Button>
                        </Popconfirm>
                    )}
                    {record.status === 'exception' && (
                        <Popconfirm
                            title="确认重新处理？"
                            description="这将重新执行对账流程"
                            onConfirm={() => handleExecute(record, 'progress')}
                        >
                            <Button type="link" size="small" icon={<SyncOutlined />}>
                                重试
                            </Button>
                        </Popconfirm>
                    )}
                </Space>
            ),
        },
    ];

    /**
     * 明细表格列配置
     */
    const detailColumns: ColumnsType<ReconciliationDetail> = [
        {
            title: '行号',
            dataIndex: 'lineNo',
            key: 'lineNo',
            width: 80,
        },
        {
            title: '外部单号',
            key: 'external',
            width: 200,
            render: (_, record) => (
                <div>
                    <div style={{ fontSize: 12, color: '#999' }}>{record.externalType}</div>
                    <div style={{ fontFamily: 'monospace' }}>{record.externalNo}</div>
                    <div style={{ fontSize: 12, color: '#999' }}>
                        {dayjs(record.externalDate).format('YYYY-MM-DD')}
                    </div>
                </div>
            ),
        },
        {
            title: '外部金额',
            dataIndex: 'externalAmount',
            key: 'externalAmount',
            width: 120,
            align: 'right',
            render: (amount: number) => <AmountDisplay value={amount} />,
        },
        {
            title: '内部单号',
            key: 'internal',
            width: 200,
            render: (_, record) => (
                <div>
                    <div style={{ fontSize: 12, color: '#999' }}>{record.internalType}</div>
                    <div style={{ fontFamily: 'monospace' }}>{record.internalNo}</div>
                    <div style={{ fontSize: 12, color: '#999' }}>
                        {dayjs(record.internalDate).format('YYYY-MM-DD')}
                    </div>
                </div>
            ),
        },
        {
            title: '内部金额',
            dataIndex: 'internalAmount',
            key: 'internalAmount',
            width: 120,
            align: 'right',
            render: (amount: number) => <AmountDisplay value={amount} />,
        },
        {
            title: '差异',
            dataIndex: 'differenceAmount',
            key: 'differenceAmount',
            width: 100,
            align: 'right',
            render: (amount: number) => (
                <AmountDisplay
                    value={amount}
                    type={amount > 0 ? 'expense' : amount < 0 ? 'income' : 'default'}
                    showCurrency={false}
                    showSign
                    size="default"
                />
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => {
                const colorMap: Record<string, string> = {
                    matched: 'success',
                    unmatched: 'error',
                    pending: 'default',
                };
                const textMap: Record<string, string> = {
                    matched: '已匹配',
                    unmatched: '未匹配',
                    pending: '待处理',
                };
                return (
                    <Tag color={colorMap[status] || 'default'}>
                        {textMap[status] || status}
                    </Tag>
                );
            },
        },
        {
            title: '备注',
            dataIndex: 'remark',
            key: 'remark',
            ellipsis: true,
            render: (remark: string) => remark || '-',
        },
    ];

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [];

    // 统计数据
    const stats = {
        total: total,
        pending: reconciliations.filter(r => r.status === 'pending').length,
        success: reconciliations.filter(r => r.status === 'success').length,
        failed: reconciliations.filter(r => r.status === 'failed').length,
        exception: reconciliations.filter(r => r.status === 'exception').length,
    };

    return (
        <PageContainer
            title="对账管理"
            subTitle="管理与第三方支付、银行等对账记录"
            extra={
                <Space size="large">
                    <Statistic
                        title="总记录"
                        value={stats.total}
                        prefix={<FileTextOutlined />}
                        valueStyle={{ fontWeight: 600 }}
                    />
                    <Statistic
                        title="待对账"
                        value={stats.pending}
                        valueStyle={{ fontWeight: 600, color: '#faad14' }}
                    />
                    <Statistic
                        title="对账成功"
                        value={stats.success}
                        valueStyle={{ fontWeight: 600, color: '#52c41a' }}
                    />
                    <Statistic
                        title="异常"
                        value={stats.exception}
                        valueStyle={{ fontWeight: 600, color: '#ff4d4f' }}
                    />
                </Space>
            }
        >
            <SearchTable
                columns={columns}
                dataSource={reconciliations}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={false}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1600 }}
            />

            {/* 详情抽屉 */}
            <Drawer
                title={
                    <span style={{ fontSize: '16px', fontWeight: 600 }}>
                        对账单详情 - {currentReconciliation?.reconciliationNo}
                    </span>
                }
                placement="right"
                width={800}
                open={detailDrawerVisible}
                onClose={() => setDetailDrawerVisible(false)}
                styles={{
                    body: { padding: '24px' }
                }}
            >
                {currentReconciliation && (
                    <div>
                        <Card
                            size="small"
                            title="基本信息"
                            style={{ marginBottom: 16 }}
                            styles={{
                                body: { padding: '16px' }
                            }}
                        >
                            <Descriptions column={2} size="small">
                                <Descriptions.Item label="对账单号">
                                    {currentReconciliation.reconciliationNo}
                                </Descriptions.Item>
                                <Descriptions.Item label="对账类型">
                                    <Tag color={reconciliationTypeMap[currentReconciliation.type].color}>
                                        {reconciliationTypeMap[currentReconciliation.type].text}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="对账状态">
                                    <Tag
                                        color={reconciliationStatusMap[currentReconciliation.status].color}
                                        icon={reconciliationStatusMap[currentReconciliation.status].icon}
                                    >
                                        {reconciliationStatusMap[currentReconciliation.status].text}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="对账日期">
                                    {dayjs(currentReconciliation.reconciliationDate).format('YYYY-MM-DD')}
                                </Descriptions.Item>
                                <Descriptions.Item label="对账期间">
                                    {dayjs(currentReconciliation.periodStart).format('YYYY-MM-DD')} ~ {dayjs(currentReconciliation.periodEnd).format('YYYY-MM-DD')}
                                </Descriptions.Item>
                                <Descriptions.Item label="差异金额">
                                    <span
                                        style={{
                                            color: currentReconciliation.differenceAmount > 0 ? '#ff4d4f' : currentReconciliation.differenceAmount < 0 ? '#52c41a' : undefined,
                                            fontFamily: 'monospace',
                                            fontWeight: 600
                                        }}
                                    >
                                        ¥{(currentReconciliation.differenceAmount / 100).toFixed(2)}
                                    </span>
                                </Descriptions.Item>
                                <Descriptions.Item label="总记录数" span={2}>
                                    {currentReconciliation.totalRecords} 条 (匹配: {currentReconciliation.matchedRecords})
                                </Descriptions.Item>
                                <Descriptions.Item label="摘要" span={2}>
                                    {currentReconciliation.abstract || '-'}
                                </Descriptions.Item>
                                <Descriptions.Item label="创建时间">
                                    {dayjs(currentReconciliation.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                                {currentReconciliation.processedAt && (
                                    <Descriptions.Item label="处理时间">
                                        {dayjs(currentReconciliation.processedAt).format('YYYY-MM-DD HH:mm:ss')}
                                    </Descriptions.Item>
                                )}
                            </Descriptions>
                        </Card>

                        <Card
                            size="small"
                            title="对账明细"
                            styles={{
                                body: { padding: '16px' }
                            }}
                        >
                            <Table
                                columns={detailColumns}
                                dataSource={reconciliationDetails}
                                rowKey="id"
                                size="small"
                                pagination={{
                                    pageSize: 10,
                                    showSizeChanger: true,
                                    showTotal: t => `共 ${t} 条`,
                                }}
                                scroll={{ x: 1000 }}
                            />
                        </Card>
                    </div>
                )}
            </Drawer>
        </PageContainer>
    );
};

export default ReconciliationPage;
