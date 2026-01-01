/**
 * 充值记录管理页面
 * 表格形式展示充值记录，支持查看详情和退款操作
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    message,
    Input,
    Select,
    DatePicker,
    Modal,
    Form,
    Descriptions,
    Avatar,
    Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    ReloadOutlined,
    EyeOutlined,
    SearchOutlined,
    RollbackOutlined,
    GiftOutlined,
} from '@ant-design/icons';
import { rechargeApi, type RechargeRecord, type RechargeRecordQueryParams, type RechargeStatus } from '@/api/recharge';
import { RECHARGE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import RefundModal from './components/RefundModal';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;
const { Text } = Typography;

interface RecordsProps {
    onStatsUpdate?: () => void;
}

/**
 * 状态选项
 */
const statusOptions = [
    { label: '全部', value: undefined },
    { label: '待支付', value: 'pending' },
    { label: '已支付', value: 'paid' },
    { label: '支付失败', value: 'failed' },
    { label: '已退款', value: 'refunded' },
    { label: '已过期', value: 'expired' },
];

const statusColorMap: Record<string, string> = {
    pending: 'orange',
    paid: 'success',
    failed: 'error',
    refunded: 'default',
    expired: 'default',
};

const statusTextMap: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    failed: '支付失败',
    refunded: '已退款',
    expired: '已过期',
};

/**
 * 充值记录管理页面
 */
const RechargeRecords: React.FC<RecordsProps> = ({ onStatsUpdate }) => {
    const [loading, setLoading] = useState(false);
    const [records, setRecords] = useState<RechargeRecord[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<RechargeRecordQueryParams>({});
    const [detailVisible, setDetailVisible] = useState(false);
    const [currentRecord, setCurrentRecord] = useState<RechargeRecord | null>(null);
    const [refundVisible, setRefundVisible] = useState(false);
    const [refundingRecord, setRefundingRecord] = useState<RechargeRecord | null>(null);

    /**
     * 加载充值记录数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const queryParams: RechargeRecordQueryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
            };
            const response = await rechargeApi.getRechargeOrders(queryParams);
            if (response.data.success) {
                const data = response.data.data || [];
                setRecords(data);
                setTotal(data.length);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load recharge records error:', error);
            message.error('加载充值记录失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 查看详情
     */
    const handleViewDetail = (record: RechargeRecord) => {
        setCurrentRecord(record);
        setDetailVisible(true);
    };

    /**
     * 退款
     */
    const handleRefund = (record: RechargeRecord) => {
        setRefundingRecord(record);
        setRefundVisible(true);
    };

    /**
     * 退款成功回调
     */
    const handleRefundSuccess = () => {
        loadData();
        onStatsUpdate?.();
        setRefundVisible(false);
        setRefundingRecord(null);
    };

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        const params: RechargeRecordQueryParams = { ...searchParams };

        if (values.orderNo) {
            params.orderNo = values.orderNo as string;
        } else {
            delete params.orderNo;
        }

        if (values.userId) {
            params.userId = Number(values.userId);
        } else {
            delete params.userId;
        }

        if (values.status) {
            params.status = values.status as RechargeStatus;
        } else {
            delete params.status;
        }

        if (values.paymentChannel) {
            params.paymentChannel = values.paymentChannel as string;
        } else {
            delete params.paymentChannel;
        }

        if (values.dateRange) {
            const range = values.dateRange as [dayjs.Dayjs, dayjs.Dayjs];
            params.startTime = range[0].startOf('day').toISOString();
            params.endTime = range[1].endOf('day').toISOString();
        } else {
            delete params.startTime;
            delete params.endTime;
        }

        setSearchParams(params);
        setCurrent(1);
    };

    /**
     * 重置搜索
     */
    const handleReset = () => {
        setSearchParams({});
        setCurrent(1);
    };

    /**
     * 表格列配置
     */
    const columns: ColumnsType<RechargeRecord> = [
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
            render: (orderNo: string) => (
                <Text copyable style={{ fontFamily: 'monospace' }}>
                    {orderNo}
                </Text>
            ),
        },
        {
            title: '用户',
            key: 'user',
            width: 150,
            render: (_, record) => (
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <Avatar
                        src={record.user?.avatarUrl}
                        size={32}
                        style={{ marginRight: 8 }}
                    >
                        {record.user?.name?.[0] || 'U'}
                    </Avatar>
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.user?.name || `用户${record.userId}`}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            ID: {record.userId}
                        </Text>
                    </div>
                </div>
            ),
        },
        {
            title: '充值档位',
            key: 'option',
            width: 150,
            render: (_, record) => (
                <div>
                    <div style={{ fontWeight: 500 }}>{record.option?.name || '-'}</div>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        ¥{(record.amountCents / 100).toFixed(2)}
                    </Text>
                </div>
            ),
        },
        {
            title: '充值金额',
            key: 'amount',
            width: 120,
            render: (_, record) => (
                <div>
                    <div style={{ fontWeight: 500, color: '#1890ff' }}>
                        ¥{(record.amountCents / 100).toFixed(2)}
                    </div>
                    {record.bonusCents > 0 && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            <GiftOutlined /> 赠送 ¥{(record.bonusCents / 100).toFixed(2)}
                        </Text>
                    )}
                </div>
            ),
        },
        {
            title: '到账金额',
            dataIndex: 'totalCents',
            key: 'totalCents',
            width: 100,
            render: (cents: number) => (
                <div style={{ fontWeight: 500, color: '#52c41a' }}>
                    ¥{(cents / 100).toFixed(2)}
                </div>
            ),
        },
        {
            title: '支付方式',
            dataIndex: 'paymentChannel',
            key: 'paymentChannel',
            width: 100,
            render: (channel?: string) => channel || '-',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => (
                <Tag color={statusColorMap[status]}>{statusTextMap[status]}</Tag>
            ),
        },
        {
            title: '支付时间',
            dataIndex: 'paidAt',
            key: 'paidAt',
            width: 180,
            render: (date?: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
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
            width: 150,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                    {record.status === 'paid' && (
                        <PermissionGuard permission={RECHARGE_PERMISSIONS.REFUND}>
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<RollbackOutlined />}
                                onClick={() => handleRefund(record)}
                            >
                                退款
                            </Button>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <>
            {/* 搜索栏 */}
            <Card style={{ marginBottom: 16 }}>
                <Form layout="inline" onFinish={handleSearch}>
                    <Form.Item name="orderNo" label="订单号">
                        <Input placeholder="请输入订单号" allowClear style={{ width: 200 }} />
                    </Form.Item>
                    <Form.Item name="userId" label="用户ID">
                        <Input placeholder="请输入用户ID" allowClear style={{ width: 120 }} />
                    </Form.Item>
                    <Form.Item name="status" label="状态">
                        <Select
                            placeholder="请选择状态"
                            allowClear
                            style={{ width: 120 }}
                            options={statusOptions}
                        />
                    </Form.Item>
                    <Form.Item name="paymentChannel" label="支付方式">
                        <Input placeholder="支付方式" allowClear style={{ width: 120 }} />
                    </Form.Item>
                    <Form.Item name="dateRange" label="时间范围">
                        <RangePicker style={{ width: 240 }} />
                    </Form.Item>
                    <Form.Item>
                        <Space>
                            <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                                搜索
                            </Button>
                            <Button onClick={handleReset}>重置</Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Card>

            {/* 操作栏 */}
            <Card style={{ marginBottom: 16 }}>
                <Space>
                    <Button icon={<ReloadOutlined />} onClick={loadData} loading={loading}>
                        刷新
                    </Button>
                </Space>
            </Card>

            {/* 表格 */}
            <Table
                columns={columns}
                dataSource={records}
                rowKey="id"
                loading={loading}
                scroll={{ x: 1400 }}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: (t) => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
            />

            {/* 详情弹窗 */}
            <Modal
                title="充值详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={700}
            >
                {currentRecord && (
                    <Descriptions column={2} bordered size="small">
                        <Descriptions.Item label="记录ID">{currentRecord.id}</Descriptions.Item>
                        <Descriptions.Item label="订单号">
                            <Text copyable>{currentRecord.orderNo}</Text>
                        </Descriptions.Item>
                        <Descriptions.Item label="用户">
                            <Space>
                                <Avatar
                                    src={currentRecord.user?.avatarUrl}
                                    size={24}
                                >
                                    {currentRecord.user?.name?.[0] || 'U'}
                                </Avatar>
                                {currentRecord.user?.name || `用户${currentRecord.userId}`}
                            </Space>
                        </Descriptions.Item>
                        <Descriptions.Item label="状态">
                            <Tag color={statusColorMap[currentRecord.status]}>
                                {statusTextMap[currentRecord.status]}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="充值档位">
                            {currentRecord.option?.name || '-'}
                        </Descriptions.Item>
                        <Descriptions.Item label="充值金额">
                            <span style={{ color: '#1890ff', fontWeight: 500 }}>
                                ¥{(currentRecord.amountCents / 100).toFixed(2)}
                            </span>
                        </Descriptions.Item>
                        <Descriptions.Item label="赠送金额">
                            {currentRecord.bonusCents > 0 ? (
                                <span style={{ color: '#52c41a', fontWeight: 500 }}>
                                    ¥{(currentRecord.bonusCents / 100).toFixed(2)}
                                </span>
                            ) : '-'}
                        </Descriptions.Item>
                        <Descriptions.Item label="到账金额">
                            <span style={{ color: '#52c41a', fontWeight: 500 }}>
                                ¥{(currentRecord.totalCents / 100).toFixed(2)}
                            </span>
                        </Descriptions.Item>
                        <Descriptions.Item label="支付方式">
                            {currentRecord.paymentChannel || '-'}
                        </Descriptions.Item>
                        <Descriptions.Item label="支付单号">
                            {currentRecord.paymentNo ? (
                                <Text copyable>{currentRecord.paymentNo}</Text>
                            ) : '-'}
                        </Descriptions.Item>
                        <Descriptions.Item label="创建时间" span={2}>
                            {dayjs(currentRecord.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Descriptions.Item>
                        {currentRecord.paidAt && (
                            <Descriptions.Item label="支付时间" span={2}>
                                {dayjs(currentRecord.paidAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                        )}
                        {currentRecord.refundedAt && (
                            <>
                                <Descriptions.Item label="退款时间">
                                    {dayjs(currentRecord.refundedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                                <Descriptions.Item label="退款原因">
                                    {currentRecord.refundReason || '-'}
                                </Descriptions.Item>
                            </>
                        )}
                    </Descriptions>
                )}
            </Modal>

            {/* 退款弹窗 */}
            {refundingRecord && (
                <RefundModal
                    visible={refundVisible}
                    record={refundingRecord}
                    onCancel={() => {
                        setRefundVisible(false);
                        setRefundingRecord(null);
                    }}
                    onSuccess={handleRefundSuccess}
                />
            )}
        </>
    );
};

export default RechargeRecords;
