/**
 * 用户拉黑管理页面
 * Manages user blocking relationships
 *
 * Features:
 * - View and filter user blocks by status, user type, and user ID
 * - View detailed block information including reason and status
 * - Admin force unblock (cancel active blocks)
 * - Batch unblock operations
 * - Delete block records
 * - View block statistics
 */
import React, { useState, useCallback, useEffect, useMemo } from 'react';
import {
    Card,
    Row,
    Col,
    Statistic,
    Tag,
    Space,
    Button,
    Avatar,
    Modal,
    Form,
    Input,
    message,
    Popconfirm,
    Drawer,
    Descriptions,
    Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    StopOutlined,
    UnlockOutlined,
    DeleteOutlined,
    EyeOutlined,
    UserOutlined,
    ExclamationCircleOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { userBlockApi } from '@/api/userBlock';
import type {
    UserBlock,
    UserBlockQueryParams,
    UserBlockStats,
    BlockStatus,
} from '@/types/userBlock';
import {
    BLOCK_STATUS_TEXT,
    BLOCK_STATUS_COLOR,
    BLOCK_USER_TYPE_TEXT,
    BLOCK_USER_TYPE_COLOR,
} from '@/types/userBlock';
import { logger } from '@/utils/logger';
import dayjs from 'dayjs';

const { Text } = Typography;

/**
 * User Block Management Page
 */
const UserBlockPage: React.FC = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [blocks, setBlocks] = useState<UserBlock[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<UserBlockQueryParams>({});

    // Statistics
    const [stats, setStats] = useState<UserBlockStats | null>(null);

    // Modal states
    const [detailVisible, setDetailVisible] = useState(false);
    const [unblockModalVisible, setUnblockModalVisible] = useState(false);
    const [currentBlock, setCurrentBlock] = useState<UserBlock | null>(null);
    const [submitting, setSubmitting] = useState(false);

    // Selected rows for batch operations
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    /**
     * Load user block data
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: UserBlockQueryParams = {
                page: current,
                pageSize: pageSize,
                ...searchParams,
            };
            const response = await userBlockApi.getUserBlocks(params);
            const data = response.data?.data;
            const blockList = Array.isArray(data?.blocks) ? data.blocks : [];
            setBlocks(blockList);
            setTotal(data?.total || 0);
        } catch (error) {
            logger.error('Load user blocks error:', error);
            message.error('加载拉黑列表失败');
            setBlocks([]);
            setTotal(0);
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, message]);

    /**
     * Load block statistics
     */
    const loadStats = useCallback(async () => {
        try {
            const response = await userBlockApi.getUserBlockStats();
            setStats(response.data.data);
        } catch (error) {
            logger.error('Load stats error:', error);
        }
    }, []);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    /**
     * Handle search
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        const params: UserBlockQueryParams = {
            ...values,
        } as UserBlockQueryParams;
        setSearchParams(params);
        setCurrent(1);
    }, []);

    /**
     * View block detail
     */
    const handleViewDetail = useCallback((block: UserBlock) => {
        setCurrentBlock(block);
        setDetailVisible(true);
    }, []);

    /**
     * Open unblock modal
     */
    const handleOpenUnblock = (block: UserBlock) => {
        if (block.status !== 'active') {
            message.warning('只能解除生效中的拉黑');
            return;
        }
        setCurrentBlock(block);
        form.resetFields();
        setUnblockModalVisible(true);
    };

    /**
     * Admin unblock
     */
    const handleUnblock = async () => {
        if (!currentBlock) return;

        try {
            const values = await form.validateFields();
            setSubmitting(true);
            await userBlockApi.adminUnblock(currentBlock.id, {
                remark: values.remark,
            });
            message.success('解除拉黑成功');
            setUnblockModalVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Unblock error:', error);
            if (error && typeof error === 'object' && 'errorFields' in error) {
                // Form validation error
                return;
            }
            message.error('解除拉黑失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * Delete block record
     */
    const handleDelete = async (block: UserBlock) => {
        try {
            await userBlockApi.deleteUserBlock(block.id);
            message.success('删除成功');
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Delete error:', error);
            message.error('删除失败');
        }
    };

    /**
     * Batch unblock
     */
    const handleBatchUnblock = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要解除拉黑的记录');
            return;
        }

        Modal.confirm({
            title: '确认批量解除拉黑',
            icon: <ExclamationCircleOutlined />,
            content: `确定要批量解除选中的 ${selectedRowKeys.length} 条拉黑记录吗？`,
            okText: '确认',
            cancelText: '取消',
            onOk: async () => {
                try {
                    setSubmitting(true);
                    await userBlockApi.batchUnblock({
                        blockIds: selectedRowKeys as number[],
                        remark: '批量解除拉黑',
                    });
                    message.success('批量解除拉黑成功');
                    setSelectedRowKeys([]);
                    loadData();
                    loadStats();
                } catch (error) {
                    logger.error('Batch unblock error:', error);
                    message.error('批量解除拉黑失败');
                } finally {
                    setSubmitting(false);
                }
            },
        });
    };

    /**
     * Batch delete
     */
    const handleBatchDelete = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要删除的记录');
            return;
        }

        Modal.confirm({
            title: '确认批量删除',
            icon: <ExclamationCircleOutlined />,
            content: `确定要批量删除选中的 ${selectedRowKeys.length} 条拉黑记录吗？此操作不可恢复！`,
            okText: '确认',
            okButtonProps: { danger: true },
            cancelText: '取消',
            onOk: async () => {
                try {
                    setSubmitting(true);
                    await userBlockApi.batchDelete({
                        blockIds: selectedRowKeys as number[],
                    });
                    message.success('批量删除成功');
                    setSelectedRowKeys([]);
                    loadData();
                    loadStats();
                } catch (error) {
                    logger.error('Batch delete error:', error);
                    message.error('批量删除失败');
                } finally {
                    setSubmitting(false);
                }
            },
        });
    };

    /**
     * Search fields configuration
     */
    const searchFields: SearchField[] = useMemo(
        () => [
            {
                name: 'blockerId',
                label: '拉黑人ID',
                type: 'input',
                placeholder: '请输入拉黑人ID',
            },
            {
                name: 'blockedId',
                label: '被拉黑人ID',
                type: 'input',
                placeholder: '请输入被拉黑人ID',
            },
            {
                name: 'blockerType',
                label: '拉黑人类型',
                type: 'select',
                options: [
                    { label: '用户', value: 'user' },
                    { label: '陪玩师', value: 'player' },
                ],
            },
            {
                name: 'blockedType',
                label: '被拉黑人类型',
                type: 'select',
                options: [
                    { label: '用户', value: 'user' },
                    { label: '陪玩师', value: 'player' },
                ],
            },
            {
                name: 'status',
                label: '状态',
                type: 'select',
                options: [
                    { label: '生效中', value: 'active' },
                    { label: '已取消', value: 'canceled' },
                    { label: '管理员解除', value: 'admin_canceled' },
                ],
            },
        ],
        []
    );

    /**
     * Quick filter handlers
     */
    const handleQuickFilter = useCallback((filter: string) => {
        switch (filter) {
            case 'active':
                setSearchParams({ status: 'active' as BlockStatus });
                break;
            case 'all':
                setSearchParams({});
                break;
        }
        setCurrent(1);
    }, []);

    /**
     * Toolbar buttons
     */
    const toolbarButtons: ToolbarButton[] = useMemo(
        () => [
            {
                text: '生效中',
                icon: <StopOutlined />,
                needSelection: false,
                type: stats?.active ? 'primary' : 'default',
                onClick: () => handleQuickFilter('active'),
            },
            {
                text: '全部',
                needSelection: false,
                onClick: () => handleQuickFilter('all'),
            },
            {
                text: '批量解除',
                icon: <UnlockOutlined />,
                needSelection: true,
                onClick: () => handleBatchUnblock(),
            },
            {
                text: '批量删除',
                icon: <DeleteOutlined />,
                needSelection: true,
                danger: true,
                onClick: () => handleBatchDelete(),
            },
        ],
        [stats, handleQuickFilter]
    );

    /**
     * Table columns
     */
    const columns: ColumnsType<UserBlock> = useMemo(
        () => [
            {
                title: 'ID',
                dataIndex: 'id',
                key: 'id',
                width: 80,
            },
            {
                title: '拉黑人',
                key: 'blocker',
                width: 200,
                render: (_, record) => (
                    <Space>
                        <Avatar
                            src={record.blockerAvatar}
                            icon={<UserOutlined />}
                            size="small"
                        />
                        <div>
                            <div>{record.blockerName || `ID: ${record.blockerId}`}</div>
                            <Tag color={BLOCK_USER_TYPE_COLOR[record.blockerType]}>
                                {BLOCK_USER_TYPE_TEXT[record.blockerType]}
                            </Tag>
                        </div>
                    </Space>
                ),
            },
            {
                title: '被拉黑人',
                key: 'blocked',
                width: 200,
                render: (_, record) => (
                    <Space>
                        <Avatar
                            src={record.blockedAvatar}
                            icon={<UserOutlined />}
                            size="small"
                        />
                        <div>
                            <div>{record.blockedName || `ID: ${record.blockedId}`}</div>
                            <Tag color={BLOCK_USER_TYPE_COLOR[record.blockedType]}>
                                {BLOCK_USER_TYPE_TEXT[record.blockedType]}
                            </Tag>
                        </div>
                    </Space>
                ),
            },
            {
                title: '拉黑原因',
                dataIndex: 'reason',
                key: 'reason',
                width: 200,
                ellipsis: true,
                render: (text) => text || '-',
            },
            {
                title: '状态',
                dataIndex: 'status',
                key: 'status',
                width: 120,
                render: (status: BlockStatus) => (
                    <Tag color={BLOCK_STATUS_COLOR[status]}>
                        {BLOCK_STATUS_TEXT[status]}
                    </Tag>
                ),
            },
            {
                title: '创建时间',
                dataIndex: 'createdAt',
                key: 'createdAt',
                width: 180,
                render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
            },
            {
                title: '操作',
                key: 'action',
                width: 180,
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
                        {record.status === 'active' && (
                            <Button
                                type="link"
                                size="small"
                                icon={<UnlockOutlined />}
                                onClick={() => handleOpenUnblock(record)}
                            >
                                解除
                            </Button>
                        )}
                        <Popconfirm
                            title="确定要删除这条记录吗？"
                            onConfirm={() => handleDelete(record)}
                            okText="确认"
                            cancelText="取消"
                        >
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                            >
                                删除
                            </Button>
                        </Popconfirm>
                    </Space>
                ),
            },
        ],
        [handleViewDetail, handleDelete]
    );

    return (
        <PageContainer title="用户拉黑管理" subTitle="管理用户之间的拉黑关系">
            {/* Statistics Cards */}
            {stats && (
                <Row gutter={16} style={{ marginBottom: 16 }}>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="总计"
                                value={stats.total}
                                prefix={<StopOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="生效中"
                                value={stats.active}
                                valueStyle={{ color: '#ff4d4f' }}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="已取消"
                                value={stats.canceled}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="管理员解除"
                                value={stats.adminCanceled}
                                valueStyle={{ color: '#faad14' }}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="今日新增"
                                value={stats.todayCount}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="用户拉黑陪玩师"
                                value={stats.userBlocksPlayer}
                            />
                        </Card>
                    </Col>
                </Row>
            )}

            <SearchTable
                columns={columns}
                dataSource={blocks}
                loading={loading}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total) => `共 ${total} 条`,
                }}
                searchFields={searchFields}
                toolbarButtons={toolbarButtons}
                rowKey="id"
                onSearch={handleSearch}
                onRefresh={loadData}
                onChange={(pagination) => {
                    setCurrent(pagination.current ?? 1);
                    setPageSize(pagination.pageSize ?? 10);
                }}
                rowSelection={{
                    selectedRowKeys,
                    onChange: (keys) => setSelectedRowKeys(keys),
                }}
            />

            {/* Detail Drawer */}
            <Drawer
                title="拉黑详情"
                open={detailVisible}
                onClose={() => setDetailVisible(false)}
                width={600}
            >
                {currentBlock && (
                    <Descriptions bordered column={1}>
                        <Descriptions.Item label="记录ID">{currentBlock.id}</Descriptions.Item>
                        <Descriptions.Item label="拉黑人">
                            <Space>
                                <Avatar
                                    src={currentBlock.blockerAvatar}
                                    icon={<UserOutlined />}
                                />
                                <div>
                                    <div>{currentBlock.blockerName || '未知'}</div>
                                    <Text type="secondary">ID: {currentBlock.blockerId}</Text>
                                    <br />
                                    <Tag color={BLOCK_USER_TYPE_COLOR[currentBlock.blockerType]}>
                                        {BLOCK_USER_TYPE_TEXT[currentBlock.blockerType]}
                                    </Tag>
                                </div>
                            </Space>
                        </Descriptions.Item>
                        <Descriptions.Item label="被拉黑人">
                            <Space>
                                <Avatar
                                    src={currentBlock.blockedAvatar}
                                    icon={<UserOutlined />}
                                />
                                <div>
                                    <div>{currentBlock.blockedName || '未知'}</div>
                                    <Text type="secondary">ID: {currentBlock.blockedId}</Text>
                                    <br />
                                    <Tag color={BLOCK_USER_TYPE_COLOR[currentBlock.blockedType]}>
                                        {BLOCK_USER_TYPE_TEXT[currentBlock.blockedType]}
                                    </Tag>
                                </div>
                            </Space>
                        </Descriptions.Item>
                        <Descriptions.Item label="拉黑原因">
                            {currentBlock.reason || '-'}
                        </Descriptions.Item>
                        <Descriptions.Item label="状态">
                            <Tag color={BLOCK_STATUS_COLOR[currentBlock.status]}>
                                {BLOCK_STATUS_TEXT[currentBlock.status]}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="创建时间">
                            {dayjs(currentBlock.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Descriptions.Item>
                        <Descriptions.Item label="更新时间">
                            {dayjs(currentBlock.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Descriptions.Item>
                        {currentBlock.status !== 'active' && (
                            <>
                                <Descriptions.Item label="取消时间">
                                    {currentBlock.canceledAt
                                        ? dayjs(currentBlock.canceledAt).format('YYYY-MM-DD HH:mm:ss')
                                        : '-'}
                                </Descriptions.Item>
                                {currentBlock.status === 'admin_canceled' && (
                                    <>
                                        <Descriptions.Item label="操作管理员">
                                            {currentBlock.adminCanceledByName ||
                                                `ID: ${currentBlock.adminCanceledBy}`}
                                        </Descriptions.Item>
                                        <Descriptions.Item label="解除备注">
                                            {currentBlock.adminCanceledRemark || '-'}
                                        </Descriptions.Item>
                                    </>
                                )}
                            </>
                        )}
                    </Descriptions>
                )}
            </Drawer>

            {/* Unblock Modal */}
            <Modal
                title="解除拉黑"
                open={unblockModalVisible}
                onOk={handleUnblock}
                onCancel={() => setUnblockModalVisible(false)}
                confirmLoading={submitting}
                okText="确认解除"
                cancelText="取消"
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="remark"
                        label="解除备注"
                        rules={[{ max: 500, message: '备注最多500个字符' }]}
                    >
                        <Input.TextArea
                            rows={4}
                            placeholder="请输入解除拉黑的原因（可选）"
                            maxLength={500}
                            showCount
                        />
                    </Form.Item>
                </Form>
                {currentBlock && (
                    <div style={{ marginTop: 16, padding: 12, background: '#f5f5f5', borderRadius: 4 }}>
                        <Text strong>当前拉黑记录：</Text>
                        <br />
                        <Text type="secondary">
                            {currentBlock.blockerName || `ID: ${currentBlock.blockerId}`} →{' '}
                            {currentBlock.blockedName || `ID: ${currentBlock.blockedId}`}
                        </Text>
                        {currentBlock.reason && (
                            <>
                                <br />
                                <Text type="secondary">原因：{currentBlock.reason}</Text>
                            </>
                        )}
                    </div>
                )}
            </Modal>
        </PageContainer>
    );
};

export default UserBlockPage;
