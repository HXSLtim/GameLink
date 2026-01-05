/**
 * Referral List Component
 * 推荐关系列表组件
 *
 * Displays all referral relationships with filtering and actions.
 */
import React, { useState, useCallback, useEffect, useMemo } from 'react';
import {
    Table,
    Tag,
    Space,
    Button,
    Avatar,
    Drawer,
    Descriptions,
    Divider,
    Modal,
    Form,
    Select,
    message,
    Card,
} from 'antd';
import {
    UserOutlined,
    EyeOutlined,
    CheckOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { referralApi } from '@/api/referral';
import type {
    Referral,
    ReferralQueryParams,
    ReferralStatus,
    ReferralType,
} from '@/api/referral';
import {
    getReferralTypeLabel,
    getReferralStatusLabel,
    getReferralStatusColor,
} from '@/api/referral';

interface ReferralListProps {
    onDataChange?: () => void;
}

const ReferralList: React.FC<ReferralListProps> = ({ onDataChange }) => {
    // State
    const [loading, setLoading] = useState(false);
    const [referrals, setReferrals] = useState<Referral[]>([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 20,
        total: 0,
    });

    // Filter states
    const [keyword, setKeyword] = useState('');
    const [typeFilter, setTypeFilter] = useState<ReferralType | undefined>();
    const [statusFilter, setStatusFilter] = useState<ReferralStatus | undefined>();

    // Detail drawer
    const [detailVisible, setDetailVisible] = useState(false);
    const [currentReferral, setCurrentReferral] = useState<Referral | null>(null);

    // Update status modal
    const [updateModalVisible, setUpdateModalVisible] = useState(false);
    const [updateLoading, setUpdateLoading] = useState(false);
    const [updateForm] = Form.useForm();

    // Batch selection
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [batchModalVisible, setBatchModalVisible] = useState(false);
    const [batchLoading, setBatchLoading] = useState(false);
    const [batchForm] = Form.useForm();

    /**
     * Load referrals
     */
    const fetchReferrals = useCallback(async () => {
        setLoading(true);
        try {
            const params: ReferralQueryParams = {
                page: pagination.current,
                page_size: pagination.pageSize,
                keyword: keyword || undefined,
                type: typeFilter,
                status: statusFilter,
            };

            const response = await referralApi.getReferrals(params);
            if (response.data.success) {
                setReferrals(response.data.data || []);
                const responsePagination = (response.data as { pagination?: { total: number } }).pagination;
                if (responsePagination) {
                    setPagination(prev => ({
                        ...prev,
                        total: responsePagination.total,
                    }));
                }
            }
        } catch {
            message.error('获取推荐列表失败');
        } finally {
            setLoading(false);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [pagination.current, pagination.pageSize, keyword, typeFilter, statusFilter]);

    useEffect(() => {
        fetchReferrals();
    }, [fetchReferrals]);

    /**
     * Handle search
     */
    const handleSearch = () => {
        setPagination(prev => ({ ...prev, current: 1 }));
    };

    /**
     * Handle reset
     */
    const handleReset = () => {
        setKeyword('');
        setTypeFilter(undefined);
        setStatusFilter(undefined);
        setPagination(prev => ({ ...prev, current: 1 }));
    };

    /**
     * Handle table change
     */
    const handleTableChange = (paginationConfig: TablePaginationConfig) => {
        setPagination(prev => ({
            ...prev,
            current: paginationConfig.current || 1,
            pageSize: paginationConfig.pageSize || 20,
        }));
    };

    /**
     * View detail
     */
    const handleViewDetail = useCallback((record: Referral) => {
        setCurrentReferral(record);
        setDetailVisible(true);
    }, []);

    /**
     * Open update modal
     */
    const openUpdateModal = useCallback((record: Referral) => {
        setCurrentReferral(record);
        updateForm.setFieldsValue({
            status: record.status,
        });
        setUpdateModalVisible(true);
    }, [updateForm]);

    /**
     * Handle update status
     */
    const handleUpdateStatus = async () => {
        if (!currentReferral) return;

        try {
            const values = await updateForm.validateFields();
            setUpdateLoading(true);

            const response = await referralApi.updateReferralStatus(
                currentReferral.id,
                { status: values.status }
            );

            if (response.data.success) {
                message.success('状态更新成功');
                setUpdateModalVisible(false);
                setCurrentReferral(null);
                fetchReferrals();
                onDataChange?.();
            }
        } catch {
            message.error('状态更新失败');
        } finally {
            setUpdateLoading(false);
        }
    };

    /**
     * Batch update status
     */
    const handleBatchUpdate = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要操作的推荐记录');
            return;
        }

        try {
            const values = await batchForm.validateFields();
            setBatchLoading(true);

            const response = await referralApi.batchUpdateReferralsStatus({
                ids: selectedRowKeys.map(key => Number(key)),
                status: values.status,
            });

            if (response.data.success) {
                const result = response.data.data as { successCount?: number; failedCount?: number };
                message.success(`成功更新 ${result?.successCount || selectedRowKeys.length} 条记录`);
                setBatchModalVisible(false);
                setSelectedRowKeys([]);
                fetchReferrals();
                onDataChange?.();
            }
        } catch {
            message.error('批量操作失败');
        } finally {
            setBatchLoading(false);
        }
    };

    /**
     * Table columns
     */
    const columns: ColumnsType<Referral> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '推荐人',
            key: 'referrer',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size="small"
                        src={record.referrer?.avatarUrl}
                        icon={<UserOutlined />}
                    />
                    <span>{record.referrer?.name || `用户${record.referrerId}`}</span>
                </Space>
            ),
        },
        {
            title: '被推荐人',
            key: 'referee',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size="small"
                        src={record.referee?.avatarUrl}
                        icon={<UserOutlined />}
                    />
                    <span>{record.referee?.name || `用户${record.refereeId}`}</span>
                </Space>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 100,
            render: (type: ReferralType) => (
                <Tag color={type === 'user' ? 'blue' : 'purple'}>
                    {getReferralTypeLabel(type)}
                </Tag>
            ),
        },
        {
            title: '邀请码',
            dataIndex: ['code', 'code'],
            key: 'code',
            width: 120,
            ellipsis: true,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: ReferralStatus) => (
                <Tag color={getReferralStatusColor(status)}>
                    {getReferralStatusLabel(status)}
                </Tag>
            ),
        },
        {
            title: '完成时间',
            dataIndex: 'completedAt',
            key: 'completedAt',
            width: 160,
            render: (time?: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 160,
            render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm:ss'),
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
                    {record.status === 'pending' && (
                        <Button
                            type="link"
                            size="small"
                            onClick={() => openUpdateModal(record)}
                        >
                            更新状态
                        </Button>
                    )}
                </Space>
            ),
        },
    ], [handleViewDetail, openUpdateModal]);

    /**
     * Row selection config
     */
    const rowSelection = {
        selectedRowKeys,
        onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
    };

    return (
        <>
            {/* Filter Bar */}
            <Card size="small" style={{ marginBottom: 16 }}>
                <Space wrap>
                    <Space.Compact>
                        <Select
                            placeholder="类型"
                            value={typeFilter}
                            onChange={setTypeFilter}
                            style={{ width: 100 }}
                            allowClear
                            options={[
                                { value: 'user', label: '用户推荐' },
                                { value: 'player', label: '陪玩师推荐' },
                            ]}
                        />
                        <Select
                            placeholder="状态"
                            value={statusFilter}
                            onChange={setStatusFilter}
                            style={{ width: 100 }}
                            allowClear
                            options={[
                                { value: 'pending', label: '待完成' },
                                { value: 'completed', label: '已完成' },
                                { value: 'canceled', label: '已取消' },
                            ]}
                        />
                    </Space.Compact>
                    <Button type="primary" onClick={handleSearch}>
                        搜索
                    </Button>
                    <Button onClick={handleReset}>
                        重置
                    </Button>
                    {selectedRowKeys.length > 0 && (
                        <>
                            <Button
                                onClick={() => setBatchModalVisible(true)}
                            >
                                批量更新 ({selectedRowKeys.length})
                            </Button>
                        </>
                    )}
                </Space>
            </Card>

            {/* Table */}
            <Table
                columns={columns}
                dataSource={referrals}
                rowKey="id"
                loading={loading}
                rowSelection={rowSelection}
                pagination={{
                    ...pagination,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total) => `共 ${total} 条`,
                }}
                onChange={handleTableChange}
                scroll={{ x: 1200 }}
            />

            {/* Detail Drawer */}
            <Drawer
                title="推荐关系详情"
                placement="right"
                size="large"
                onClose={() => setDetailVisible(false)}
                open={detailVisible}
            >
                {currentReferral && (
                    <>
                        <Descriptions column={2} bordered>
                            <Descriptions.Item label="推荐ID">{currentReferral.id}</Descriptions.Item>
                            <Descriptions.Item label="类型">
                                <Tag color={currentReferral.type === 'user' ? 'blue' : 'purple'}>
                                    {getReferralTypeLabel(currentReferral.type)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="状态" span={2}>
                                <Tag color={getReferralStatusColor(currentReferral.status)}>
                                    {getReferralStatusLabel(currentReferral.status)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="推荐人">
                                <Space>
                                    <Avatar
                                        src={currentReferral.referrer?.avatarUrl}
                                        icon={<UserOutlined />}
                                    />
                                    {currentReferral.referrer?.name || `用户${currentReferral.referrerId}`}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="被推荐人">
                                <Space>
                                    <Avatar
                                        src={currentReferral.referee?.avatarUrl}
                                        icon={<UserOutlined />}
                                    />
                                    {currentReferral.referee?.name || `用户${currentReferral.refereeId}`}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="邀请码" span={2}>
                                {currentReferral.code?.code || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间" span={2}>
                                {dayjs(currentReferral.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            {currentReferral.completedAt && (
                                <Descriptions.Item label="完成时间" span={2}>
                                    {dayjs(currentReferral.completedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                            )}
                            {currentReferral.cancelReason && (
                                <Descriptions.Item label="取消原因" span={2}>
                                    {currentReferral.cancelReason}
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider />
                        <Space>
                            {currentReferral.status === 'pending' && (
                                <Button
                                    type="primary"
                                    icon={<CheckOutlined />}
                                    onClick={() => {
                                        setDetailVisible(false);
                                        openUpdateModal(currentReferral);
                                    }}
                                >
                                    更新状态
                                </Button>
                            )}
                        </Space>
                    </>
                )}
            </Drawer>

            {/* Update Status Modal */}
            <Modal
                title="更新推荐状态"
                open={updateModalVisible}
                onOk={handleUpdateStatus}
                onCancel={() => {
                    setUpdateModalVisible(false);
                    setCurrentReferral(null);
                    updateForm.resetFields();
                }}
                confirmLoading={updateLoading}
            >
                <Form form={updateForm} layout="vertical">
                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{ required: true, message: '请选择状态' }]}
                    >
                        <Select
                            placeholder="请选择状态"
                            options={[
                                { value: 'pending', label: '待完成' },
                                { value: 'completed', label: '已完成' },
                                { value: 'canceled', label: '已取消' },
                            ]}
                        />
                    </Form.Item>
                </Form>
            </Modal>

            {/* Batch Update Modal */}
            <Modal
                title="批量更新推荐状态"
                open={batchModalVisible}
                onOk={handleBatchUpdate}
                onCancel={() => {
                    setBatchModalVisible(false);
                    batchForm.resetFields();
                }}
                confirmLoading={batchLoading}
            >
                <div style={{ marginBottom: 16 }}>
                    已选择 <strong>{selectedRowKeys.length}</strong> 条记录
                </div>
                <Form form={batchForm} layout="vertical">
                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{ required: true, message: '请选择状态' }]}
                    >
                        <Select
                            placeholder="请选择状态"
                            options={[
                                { value: 'pending', label: '待完成' },
                                { value: 'completed', label: '已完成' },
                                { value: 'canceled', label: '已取消' },
                            ]}
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
};

export default ReferralList;
