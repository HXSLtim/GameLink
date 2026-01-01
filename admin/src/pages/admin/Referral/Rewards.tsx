/**
 * Referral Rewards Management Component
 * 推荐奖励管理组件
 *
 * Manages reward distribution for referrals.
 */
import React, { useState, useCallback, useEffect, useMemo } from 'react';
import {
    Table,
    Tag,
    Space,
    Button,
    Avatar,
    Modal,
    message,
    Card,
    Select,
    Drawer,
    Descriptions,
    Divider,
    Input,
} from 'antd';
import {
    DollarOutlined,
    GiftOutlined,
    CheckOutlined,
    CloseOutlined,
    EyeOutlined,
    UserOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { referralApi, type ReferralRewardQueryParams, type ReferralReward, type ReferralRewardStatus, type RewardType } from '@/api/referral';
import {
    getRewardTypeLabel,
    getRewardStatusLabel,
    getRewardStatusColor,
    centsToYuan,
} from '@/api/referral';
import IssueModal from './components/IssueModal';

interface RewardsProps {
    onDataChange?: () => void;
}

const Rewards: React.FC<RewardsProps> = ({ onDataChange }) => {
    // State
    const [loading, setLoading] = useState(false);
    const [rewards, setRewards] = useState<ReferralReward[]>([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 20,
        total: 0,
    });

    // Filter states
    const [keyword, setKeyword] = useState('');
    const [typeFilter, setTypeFilter] = useState<RewardType | undefined>();
    const [statusFilter, setStatusFilter] = useState<ReferralRewardStatus | undefined>();

    // Detail drawer
    const [detailVisible, setDetailVisible] = useState(false);
    const [currentReward, setCurrentReward] = useState<ReferralReward | null>(null);

    // Issue modal
    const [issueModalVisible, setIssueModalVisible] = useState(false);
    const [issueReward, setIssueReward] = useState<ReferralReward | null>(null);
    const [issueAction, setIssueAction] = useState<'issue' | 'fail'>('issue');

    // Batch selection
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    /**
     * Load rewards
     */
    const fetchRewards = useCallback(async () => {
        setLoading(true);
        try {
            const params: ReferralRewardQueryParams = {
                page: pagination.current,
                page_size: pagination.pageSize,
                keyword: keyword || undefined,
                type: typeFilter,
                status: statusFilter,
            };

            const response = await referralApi.getReferralRewards(params);
            if (response.data.success) {
                setRewards(response.data.data || []);
                const responsePagination = (response.data as { pagination?: { total: number } }).pagination;
                if (responsePagination) {
                    setPagination(prev => ({
                        ...prev,
                        total: responsePagination.total,
                    }));
                }
            }
        } catch (err) {
            message.error('获取奖励列表失败');
            console.error('Failed to fetch rewards:', err);
        } finally {
            setLoading(false);
        }
    }, [pagination.current, pagination.pageSize, keyword, typeFilter, statusFilter]);

    useEffect(() => {
        fetchRewards();
    }, [fetchRewards]);

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
    const handleViewDetail = useCallback((reward: ReferralReward) => {
        setCurrentReward(reward);
        setDetailVisible(true);
    }, []);

    /**
     * Open issue modal
     */
    const openIssueModal = (reward: ReferralReward, action: 'issue' | 'fail') => {
        setIssueReward(reward);
        setIssueAction(action);
        setIssueModalVisible(true);
    };

    /**
     * Handle issue reward
     */
    const handleIssueReward = async (reward: ReferralReward) => {
        try {
            const response = await referralApi.issueReferralReward(reward.id);
            if (response.data.success) {
                message.success('奖励发放成功');
                setIssueModalVisible(false);
                fetchRewards();
                onDataChange?.();
            }
        } catch (err) {
            message.error('奖励发放失败');
        }
    };

    /**
     * Handle fail reward
     */
    const handleFailReward = async (reward: ReferralReward, reason: string) => {
        try {
            const response = await referralApi.failReferralReward(reward.id, { reason });
            if (response.data.success) {
                message.success('已标记为发放失败');
                setIssueModalVisible(false);
                fetchRewards();
                onDataChange?.();
            }
        } catch (err) {
            message.error('操作失败');
        }
    };

    /**
     * Table columns
     */
    const columns: ColumnsType<ReferralReward> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '推荐关系ID',
            dataIndex: 'referralId',
            key: 'referralId',
            width: 100,
        },
        {
            title: '接收用户',
            key: 'user',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size="small"
                        src={record.user?.avatarUrl}
                        icon={<UserOutlined />}
                    />
                    <span>{record.user?.name || `用户${record.userId}`}</span>
                </Space>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 100,
            render: (type: RewardType) => (
                <Tag color={type === 'referrer' ? 'blue' : 'green'}>
                    {getRewardTypeLabel(type)}
                </Tag>
            ),
        },
        {
            title: '金额',
            dataIndex: 'amountCents',
            key: 'amountCents',
            width: 100,
            render: (amount: number) => (
                <span style={{ fontWeight: 500, color: '#faad14' }}>
                    ¥{centsToYuan(amount)}
                </span>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: ReferralRewardStatus) => (
                <Tag color={getRewardStatusColor(status)}>
                    {getRewardStatusLabel(status)}
                </Tag>
            ),
        },
        {
            title: '发放时间',
            dataIndex: 'issuedAt',
            key: 'issuedAt',
            width: 140,
            render: (time?: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-',
        },
        {
            title: '失败原因',
            dataIndex: 'failureReason',
            key: 'failureReason',
            width: 150,
            ellipsis: true,
            render: (reason?: string) => reason || '-',
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 140,
            render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm'),
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
                            <Button
                                type="link"
                                size="small"
                                icon={<CheckOutlined />}
                                onClick={() => openIssueModal(record, 'issue')}
                            >
                                发放
                            </Button>
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<CloseOutlined />}
                                onClick={() => openIssueModal(record, 'fail')}
                            >
                                失败
                            </Button>
                        </>
                    )}
                </Space>
            ),
        },
    ], [handleViewDetail]);

    /**
     * Row selection config
     */
    const rowSelection = {
        selectedRowKeys,
        onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
        getCheckboxProps: (record: ReferralReward) => ({
            disabled: record.status !== 'pending',
        }),
    };

    /**
     * Batch issue pending rewards
     */
    const handleBatchIssue = async () => {
        const pendingRewards = selectedRowKeys
            .map(key => rewards.find(r => r.id === key))
            .filter((r): r is ReferralReward => r !== undefined && r.status === 'pending');

        if (pendingRewards.length === 0) {
            message.warning('请选择待发放的奖励');
            return;
        }

        Modal.confirm({
            title: '批量发放奖励',
            content: `确定要发放 ${pendingRewards.length} 个奖励吗？`,
            onOk: async () => {
                let successCount = 0;
                let failCount = 0;

                for (const reward of pendingRewards) {
                    try {
                        await referralApi.issueReferralReward(reward.id);
                        successCount++;
                    } catch {
                        failCount++;
                    }
                }

                if (successCount > 0) {
                    message.success(`成功发放 ${successCount} 个奖励`);
                }
                if (failCount > 0) {
                    message.warning(`${failCount} 个奖励发放失败`);
                }

                setSelectedRowKeys([]);
                fetchRewards();
                onDataChange?.();
            },
        });
    };

    return (
        <>
            {/* Filter Bar */}
            <Card size="small" style={{ marginBottom: 16 }}>
                <Space wrap>
                    <Input
                        placeholder="搜索用户ID/推荐ID"
                        value={keyword}
                        onChange={e => setKeyword(e.target.value)}
                        style={{ width: 150 }}
                        allowClear
                    />
                    <Select
                        placeholder="类型"
                        value={typeFilter}
                        onChange={setTypeFilter}
                        style={{ width: 100 }}
                        allowClear
                        options={[
                            { value: 'referrer', label: '推荐人奖励' },
                            { value: 'referee', label: '被推荐人奖励' },
                        ]}
                    />
                    <Select
                        placeholder="状态"
                        value={statusFilter}
                        onChange={setStatusFilter}
                        style={{ width: 100 }}
                        allowClear
                        options={[
                            { value: 'pending', label: '待发放' },
                            { value: 'issued', label: '已发放' },
                            { value: 'failed', label: '发放失败' },
                        ]}
                    />
                    <Button type="primary" onClick={handleSearch}>
                        搜索
                    </Button>
                    <Button onClick={handleReset}>
                        重置
                    </Button>
                    {selectedRowKeys.length > 0 && (
                        <Button
                            type="primary"
                            icon={<CheckOutlined />}
                            onClick={handleBatchIssue}
                            style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
                        >
                            批量发放 ({selectedRowKeys.filter(key => {
                                const reward = rewards.find(r => r.id === key);
                                return reward?.status === 'pending';
                            }).length})
                        </Button>
                    )}
                </Space>
            </Card>

            {/* Table */}
            <Table
                columns={columns}
                dataSource={rewards}
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
                scroll={{ x: 1400 }}
            />

            {/* Detail Drawer */}
            <Drawer
                title="奖励详情"
                placement="right"
                size="large"
                onClose={() => setDetailVisible(false)}
                open={detailVisible}
            >
                {currentReward && (
                    <>
                        <div style={{ textAlign: 'center', marginBottom: 24 }}>
                            <div style={{ fontSize: 32, fontWeight: 'bold', color: '#faad14' }}>
                                <DollarOutlined style={{ marginRight: 8 }} />
                                ¥{centsToYuan(currentReward.amountCents)}
                            </div>
                            <Tag color={getRewardStatusColor(currentReward.status)} style={{ fontSize: 14, marginTop: 8 }}>
                                {getRewardStatusLabel(currentReward.status)}
                            </Tag>
                        </div>

                        <Descriptions column={2} bordered>
                            <Descriptions.Item label="奖励ID">{currentReward.id}</Descriptions.Item>
                            <Descriptions.Item label="推荐关系ID">{currentReward.referralId}</Descriptions.Item>
                            <Descriptions.Item label="类型">
                                <Tag color={currentReward.type === 'referrer' ? 'blue' : 'green'}>
                                    {getRewardTypeLabel(currentReward.type)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={getRewardStatusColor(currentReward.status)}>
                                    {getRewardStatusLabel(currentReward.status)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="接收用户">
                                <Space>
                                    <Avatar
                                        src={currentReward.user?.avatarUrl}
                                        icon={<UserOutlined />}
                                    />
                                    {currentReward.user?.name || `用户${currentReward.userId}`}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="金额">
                                <span style={{ fontWeight: 500, color: '#faad14' }}>
                                    ¥{centsToYuan(currentReward.amountCents)}
                                </span>
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间" span={2}>
                                {dayjs(currentReward.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            {currentReward.issuedAt && (
                                <Descriptions.Item label="发放时间" span={2}>
                                    {dayjs(currentReward.issuedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                            )}
                            {currentReward.failedAt && (
                                <Descriptions.Item label="失败时间" span={2}>
                                    {dayjs(currentReward.failedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                            )}
                            {currentReward.failureReason && (
                                <Descriptions.Item label="失败原因" span={2}>
                                    {currentReward.failureReason}
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider />
                        {currentReward.status === 'pending' && (
                            <Space>
                                <Button
                                    type="primary"
                                    icon={<CheckOutlined />}
                                    onClick={() => {
                                        setDetailVisible(false);
                                        openIssueModal(currentReward, 'issue');
                                    }}
                                >
                                    立即发放
                                </Button>
                                <Button
                                    danger
                                    icon={<CloseOutlined />}
                                    onClick={() => {
                                        setDetailVisible(false);
                                        openIssueModal(currentReward, 'fail');
                                    }}
                                >
                                    标记失败
                                </Button>
                            </Space>
                        )}
                    </>
                )}
            </Drawer>

            {/* Issue/Fail Modal */}
            <IssueModal
                visible={issueModalVisible}
                reward={issueReward}
                action={issueAction}
                onIssue={handleIssueReward}
                onFail={handleFailReward}
                onCancel={() => {
                    setIssueModalVisible(false);
                    setIssueReward(null);
                }}
            />
        </>
    );
};

export default Rewards;
