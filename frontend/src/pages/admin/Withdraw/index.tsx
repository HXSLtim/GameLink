/**
 * 提现管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    message,
    Descriptions,
    Statistic,
    Card,
    Row,
    Col,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    CheckOutlined,
    CloseOutlined,
    DollarOutlined,
    EyeOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { WITHDRAW_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type Withdraw, type WithdrawQueryParams } from '@/api/admin';
import dayjs from 'dayjs';

/**
 * 状态选项
 */
const statusOptions = [
    { label: '待审核', value: 'pending' },
    { label: '已批准', value: 'approved' },
    { label: '已拒绝', value: 'rejected' },
    { label: '已完成', value: 'completed' },
];

const statusColorMap: Record<string, string> = {
    pending: 'orange',
    approved: 'blue',
    rejected: 'red',
    completed: 'green',
};

const statusTextMap: Record<string, string> = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝',
    completed: '已完成',
};

/**
 * 提现管理页面
 */
const WithdrawPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [withdraws, setWithdraws] = useState<Withdraw[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<WithdrawQueryParams>({});

    // 统计数据
    const [stats, setStats] = useState({ pending: 0, approved: 0, completed: 0, totalAmount: 0 });

    // 弹窗状态
    const [detailVisible, setDetailVisible] = useState(false);
    const [currentWithdraw, setCurrentWithdraw] = useState<Withdraw | null>(null);
    const [rejectVisible, setRejectVisible] = useState(false);
    const [rejectForm] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    /**
     * 加载提现数据
     */
    const loadData = useCallback(async (params?: WithdrawQueryParams) => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
                ...params,
            };
            // API client 响应拦截器已返回 response.data，直接访问 success
            const response = await adminApi.getWithdraws(queryParams) as unknown as {
                success: boolean;
                message?: string;
                data?: { withdraws: Withdraw[]; total: number };
            };
            if (response.success) {
                setWithdraws(response.data?.withdraws || []);
                setTotal(response.data?.total || 0);
            } else {
                message.error(response.message || '加载失败');
            }
        } catch (error) {
            console.error('Load withdraws error:', error);
            message.error('加载提现列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    /**
     * 加载统计数据
     */
    const loadStats = useCallback(async () => {
        try {
            // 获取各状态数量
            type WithdrawResponse = { success: boolean; data?: { total: number } };
            const [pendingRes, approvedRes, completedRes] = await Promise.all([
                adminApi.getWithdraws({ status: 'pending', page_size: 1 }) as unknown as WithdrawResponse,
                adminApi.getWithdraws({ status: 'approved', page_size: 1 }) as unknown as WithdrawResponse,
                adminApi.getWithdraws({ status: 'completed', page_size: 1 }) as unknown as WithdrawResponse,
            ]);
            setStats({
                pending: pendingRes.data?.total || 0,
                approved: approvedRes.data?.total || 0,
                completed: completedRes.data?.total || 0,
                totalAmount: 0,
            });
        } catch (error) {
            console.error('Load stats error:', error);
        }
    }, []);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values as WithdrawQueryParams);
        setCurrent(1);
    };

    /**
     * 查看详情
     */
    const handleViewDetail = (record: Withdraw) => {
        setCurrentWithdraw(record);
        setDetailVisible(true);
    };

    /**
     * 批准提现
     */
    const handleApprove = async (record: Withdraw) => {
        Modal.confirm({
            title: '确认批准',
            content: `确定要批准该提现申请吗？金额：¥${(record.amountCents / 100).toFixed(2)}`,
            onOk: async () => {
                try {
                    setSubmitting(true);
                    await adminApi.approveWithdraw(record.id);
                    message.success('批准成功');
                    loadData();
                    loadStats();
                } catch (error) {
                    console.error('Approve error:', error);
                    message.error('批准失败');
                } finally {
                    setSubmitting(false);
                }
            },
        });
    };

    /**
     * 拒绝提现
     */
    const handleReject = (record: Withdraw) => {
        setCurrentWithdraw(record);
        rejectForm.resetFields();
        setRejectVisible(true);
    };

    const submitReject = async () => {
        try {
            const values = await rejectForm.validateFields();
            setSubmitting(true);
            await adminApi.rejectWithdraw(currentWithdraw!.id, { reason: values.reason });
            message.success('已拒绝');
            setRejectVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            console.error('Reject error:', error);
            message.error('操作失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * 完成提现（打款）
     */
    const handleComplete = async (record: Withdraw) => {
        Modal.confirm({
            title: '确认完成',
            content: `确定已完成打款吗？金额：¥${(record.amountCents / 100).toFixed(2)}`,
            onOk: async () => {
                try {
                    setSubmitting(true);
                    await adminApi.completeWithdraw(record.id);
                    message.success('已标记为完成');
                    loadData();
                    loadStats();
                } catch (error) {
                    console.error('Complete error:', error);
                    message.error('操作失败');
                } finally {
                    setSubmitting(false);
                }
            },
        });
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'status', label: '状态', type: 'select', options: statusOptions },
        { name: 'playerId', label: '陪玩师ID', type: 'input', placeholder: '请输入陪玩师ID' },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Withdraw> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 150,
            render: (_, record) => record.player?.nickname || `ID: ${record.playerId}`,
        },
        {
            title: '提现金额',
            dataIndex: 'amountCents',
            key: 'amountCents',
            width: 120,
            render: (cents: number) => (
                <span style={{ color: '#f5222d', fontWeight: 500 }}>
                    ¥{(cents / 100).toFixed(2)}
                </span>
            ),
        },
        {
            title: '银行信息',
            key: 'bank',
            width: 200,
            render: (_, record) => (
                <div>
                    <div>{record.bankName || '-'}</div>
                    <div style={{ color: '#999', fontSize: 12 }}>
                        {record.bankAccount ? `****${record.bankAccount.slice(-4)}` : '-'}
                    </div>
                </div>
            ),
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
            title: '申请时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '处理时间',
            dataIndex: 'processedAt',
            key: 'processedAt',
            width: 180,
            render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
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
                    {record.status === 'pending' && (
                        <>
                            <PermissionGuard permission={WITHDRAW_PERMISSIONS.APPROVE}>
                                <Button
                                    type="link"
                                    size="small"
                                    icon={<CheckOutlined />}
                                    onClick={() => handleApprove(record)}
                                    loading={submitting}
                                >
                                    批准
                                </Button>
                            </PermissionGuard>
                            <PermissionGuard permission={WITHDRAW_PERMISSIONS.REJECT}>
                                <Button
                                    type="link"
                                    size="small"
                                    danger
                                    icon={<CloseOutlined />}
                                    onClick={() => handleReject(record)}
                                >
                                    拒绝
                                </Button>
                            </PermissionGuard>
                        </>
                    )}
                    {record.status === 'approved' && (
                        <PermissionGuard permission={WITHDRAW_PERMISSIONS.APPROVE}>
                            <Button
                                type="link"
                                size="small"
                                icon={<DollarOutlined />}
                                onClick={() => handleComplete(record)}
                                loading={submitting}
                            >
                                完成打款
                            </Button>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <PageContainer title="提现管理" subTitle="管理陪玩师提现申请">
            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="待审核"
                            value={stats.pending}
                            valueStyle={{ color: '#faad14' }}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="已批准待打款"
                            value={stats.approved}
                            valueStyle={{ color: '#1890ff' }}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="已完成"
                            value={stats.completed}
                            valueStyle={{ color: '#52c41a' }}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="本月提现总额"
                            value={stats.totalAmount / 100}
                            precision={2}
                            prefix="¥"
                        />
                    </Card>
                </Col>
            </Row>

            <SearchTable
                columns={columns}
                dataSource={withdraws}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1200 }}
            />

            {/* 详情弹窗 */}
            <Modal
                title="提现详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={600}
            >
                {currentWithdraw && (
                    <Descriptions column={2} bordered size="small">
                        <Descriptions.Item label="提现ID">{currentWithdraw.id}</Descriptions.Item>
                        <Descriptions.Item label="状态">
                            <Tag color={statusColorMap[currentWithdraw.status]}>
                                {statusTextMap[currentWithdraw.status]}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="陪玩师">
                            {currentWithdraw.player?.nickname || `ID: ${currentWithdraw.playerId}`}
                        </Descriptions.Item>
                        <Descriptions.Item label="提现金额">
                            <span style={{ color: '#f5222d', fontWeight: 500 }}>
                                ¥{(currentWithdraw.amountCents / 100).toFixed(2)}
                            </span>
                        </Descriptions.Item>
                        <Descriptions.Item label="银行名称">{currentWithdraw.bankName || '-'}</Descriptions.Item>
                        <Descriptions.Item label="银行账号">{currentWithdraw.bankAccount || '-'}</Descriptions.Item>
                        <Descriptions.Item label="账户名">{currentWithdraw.accountName || '-'}</Descriptions.Item>
                        <Descriptions.Item label="申请时间">
                            {dayjs(currentWithdraw.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Descriptions.Item>
                        <Descriptions.Item label="用户备注" span={2}>
                            {currentWithdraw.remark || '-'}
                        </Descriptions.Item>
                        {currentWithdraw.status === 'rejected' && (
                            <Descriptions.Item label="拒绝原因" span={2}>
                                <span style={{ color: '#f5222d' }}>{currentWithdraw.rejectReason}</span>
                            </Descriptions.Item>
                        )}
                        {currentWithdraw.processedAt && (
                            <>
                                <Descriptions.Item label="处理时间">
                                    {dayjs(currentWithdraw.processedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                                <Descriptions.Item label="管理员备注">
                                    {currentWithdraw.adminRemark || '-'}
                                </Descriptions.Item>
                            </>
                        )}
                        {currentWithdraw.completedAt && (
                            <Descriptions.Item label="完成时间" span={2}>
                                {dayjs(currentWithdraw.completedAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                        )}
                    </Descriptions>
                )}
            </Modal>

            {/* 拒绝弹窗 */}
            <Modal
                title="拒绝提现"
                open={rejectVisible}
                onOk={submitReject}
                onCancel={() => setRejectVisible(false)}
                confirmLoading={submitting}
            >
                <Form form={rejectForm} layout="vertical">
                    <Form.Item
                        name="reason"
                        label="拒绝原因"
                        rules={[{ required: true, message: '请输入拒绝原因' }]}
                    >
                        <Input.TextArea rows={4} placeholder="请输入拒绝原因" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default WithdrawPage;
