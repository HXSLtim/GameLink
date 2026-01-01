/**
 * Referral Codes Management Component
 * 邀请码管理组件
 *
 * Manages invitation codes with CRUD operations.
 */
import React, { useState, useCallback, useEffect, useMemo } from 'react';
import {
    Table,
    Tag,
    Space,
    Button,
    Avatar,
    message,
    Popconfirm,
    Card,
    Input,
    Select,
    Progress,
    Drawer,
    Descriptions,
    Divider,
} from 'antd';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    EyeOutlined,
    CheckOutlined,
    CloseOutlined,
    UserOutlined,
    CopyOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { referralApi } from '@/api/referral';
import type {
    ReferralCode,
    ReferralCodeQueryParams,
    ReferralType,
    CreateReferralCodeDto,
    UpdateReferralCodeDto,
} from '@/api/referral';
import {
    getReferralTypeLabel,
    isCodeExpired,
    isCodeFullyUsed,
    getCodeUsagePercent,
} from '@/api/referral';
import CodeForm from './components/CodeForm';

interface CodesProps {
    onDataChange?: () => void;
}

const Codes: React.FC<CodesProps> = ({ onDataChange }) => {
    // State
    const [loading, setLoading] = useState(false);
    const [codes, setCodes] = useState<ReferralCode[]>([]);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 20,
        total: 0,
    });

    // Filter states
    const [keyword, setKeyword] = useState('');
    const [typeFilter, setTypeFilter] = useState<ReferralType | undefined>();
    const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive' | 'expired' | 'full'>('all');

    // Create/Edit modal
    const [formModalVisible, setFormModalVisible] = useState(false);
    const [editingCode, setEditingCode] = useState<ReferralCode | null>(null);

    // Detail drawer
    const [detailVisible, setDetailVisible] = useState(false);
    const [currentCode, setCurrentCode] = useState<ReferralCode | null>(null);

    // Batch selection
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

    /**
     * Load codes
     */
    const fetchCodes = useCallback(async () => {
        setLoading(true);
        try {
            const params: ReferralCodeQueryParams = {
                page: pagination.current,
                page_size: pagination.pageSize,
                keyword: keyword || undefined,
                type: typeFilter,
            };

            // Handle status filter
            if (statusFilter === 'active') {
                params.isActive = true;
            } else if (statusFilter === 'inactive') {
                params.isActive = false;
            }

            const response = await referralApi.getReferralCodes(params);
            if (response.data.success) {
                let filteredData = response.data.data || [];

                // Client-side filtering for expired and full
                if (statusFilter === 'expired') {
                    filteredData = filteredData.filter(code => isCodeExpired(code));
                } else if (statusFilter === 'full') {
                    filteredData = filteredData.filter(code => isCodeFullyUsed(code));
                }

                setCodes(filteredData);
                const responsePagination = (response.data as { pagination?: { total: number } }).pagination;
                if (responsePagination) {
                    setPagination(prev => ({
                        ...prev,
                        total: statusFilter === 'expired' || statusFilter === 'full'
                            ? filteredData.length
                            : responsePagination.total,
                    }));
                }
            }
        } catch {
            message.error('获取邀请码列表失败');
        } finally {
            setLoading(false);
        }
    }, [pagination.current, pagination.pageSize, keyword, typeFilter, statusFilter]);

    useEffect(() => {
        fetchCodes();
    }, [fetchCodes]);

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
        setStatusFilter('all');
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
     * Open create modal
     */
    const openCreateModal = () => {
        setEditingCode(null);
        setFormModalVisible(true);
    };

    /**
     * Open edit modal
     */
    const openEditModal = (code: ReferralCode) => {
        setEditingCode(code);
        setFormModalVisible(true);
    };

    /**
     * Handle create/update code
     */
    const handleSaveCode = async (values: CreateReferralCodeDto | UpdateReferralCodeDto) => {
        try {
            if (editingCode) {
                // Update
                const response = await referralApi.updateReferralCode(
                    editingCode.id,
                    values as UpdateReferralCodeDto
                );
                if (response.data.success) {
                    message.success('更新成功');
                }
            } else {
                // Create
                const response = await referralApi.createReferralCode(
                    values as CreateReferralCodeDto
                );
                if (response.data.success) {
                    message.success('创建成功');
                }
            }
            setFormModalVisible(false);
            setEditingCode(null);
            fetchCodes();
            onDataChange?.();
        } catch {
            message.error(editingCode ? '更新失败' : '创建失败');
        }
    };

    /**
     * Handle delete code
     */
    const handleDelete = async (id: number) => {
        try {
            const response = await referralApi.deleteReferralCode(id);
            if (response.data.success) {
                message.success('删除成功');
                fetchCodes();
                onDataChange?.();
            }
        } catch {
            message.error('删除失败');
        }
    };

    /**
     * Handle toggle active status
     */
    const handleToggleActive = async (code: ReferralCode) => {
        try {
            const response = await referralApi.updateReferralCode(code.id, {
                isActive: !code.isActive,
            });
            if (response.data.success) {
                message.success(code.isActive ? '已禁用' : '已启用');
                fetchCodes();
                onDataChange?.();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * View detail
     */
    const handleViewDetail = useCallback((code: ReferralCode) => {
        setCurrentCode(code);
        setDetailVisible(true);
    }, []);

    /**
     * Copy code to clipboard
     */
    const handleCopyCode = (code: string) => {
        navigator.clipboard.writeText(code);
        message.success('已复制到剪贴板');
    };

    /**
     * Batch toggle status
     */
    const handleBatchToggleStatus = async (isActive: boolean) => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要操作的邀请码');
            return;
        }

        try {
            const response = await referralApi.batchUpdateCodesStatus({
                ids: selectedRowKeys.map(key => Number(key)),
                isActive,
            });

            if (response.data.success) {
                const result = response.data.data as { successCount?: number };
                message.success(`成功更新 ${result?.successCount || selectedRowKeys.length} 个邀请码`);
                setSelectedRowKeys([]);
                fetchCodes();
                onDataChange?.();
            }
        } catch {
            message.error('批量操作失败');
        }
    };

    /**
     * Batch delete
     */
    const handleBatchDelete = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要删除的邀请码');
            return;
        }

        try {
            const response = await referralApi.batchDeleteCodes({
                ids: selectedRowKeys.map(key => Number(key)),
            });

            if (response.data.success) {
                const result = response.data.data as { successCount?: number };
                message.success(`成功删除 ${result?.successCount || selectedRowKeys.length} 个邀请码`);
                setSelectedRowKeys([]);
                fetchCodes();
                onDataChange?.();
            }
        } catch {
            message.error('批量删除失败');
        }
    };

    /**
     * Table columns
     */
    const columns: ColumnsType<ReferralCode> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '邀请码',
            dataIndex: 'code',
            key: 'code',
            width: 150,
            render: (code: string, _record) => (
                <Space>
                    <Tag color="blue" style={{ fontFamily: 'monospace', fontSize: 14 }}>
                        {code}
                    </Tag>
                    <Button
                        type="text"
                        size="small"
                        icon={<CopyOutlined />}
                        onClick={() => handleCopyCode(code)}
                    />
                </Space>
            ),
        },
        {
            title: '拥有者',
            key: 'owner',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size="small"
                        src={record.owner?.avatarUrl}
                        icon={<UserOutlined />}
                    />
                    <span>{record.owner?.name || `用户${record.ownerId}`}</span>
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
            title: '使用情况',
            key: 'usage',
            width: 150,
            render: (_, record) => (
                <div>
                    <Progress
                        percent={getCodeUsagePercent(record)}
                        size="small"
                        status={isCodeFullyUsed(record) ? 'exception' : 'active'}
                    />
                    <div style={{ fontSize: 12, marginTop: 4 }}>
                        {record.usedCount} / {record.maxUses}
                    </div>
                </div>
            ),
        },
        {
            title: '有效期至',
            dataIndex: 'expiresAt',
            key: 'expiresAt',
            width: 140,
            render: (time: string, record) => {
                const expired = isCodeExpired(record);
                return (
                    <span style={{ color: expired ? '#ff4d4f' : undefined }}>
                        {dayjs(time).format('YYYY-MM-DD HH:mm')}
                    </span>
                );
            },
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: 100,
            render: (isActive: boolean, record) => {
                const expired = isCodeExpired(record);
                const full = isCodeFullyUsed(record);
                let color = 'success';
                let text = '启用';

                if (!isActive) {
                    color = 'default';
                    text = '禁用';
                } else if (expired) {
                    color = 'error';
                    text = '已过期';
                } else if (full) {
                    color = 'warning';
                    text = '已用完';
                }

                return <Tag color={color}>{text}</Tag>;
            },
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
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => openEditModal(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={record.isActive ? <CloseOutlined /> : <CheckOutlined />}
                        onClick={() => handleToggleActive(record)}
                    >
                        {record.isActive ? '禁用' : '启用'}
                    </Button>
                    <Popconfirm
                        title="确定要删除这个邀请码吗?"
                        onConfirm={() => handleDelete(record.id)}
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
    ], [handleViewDetail]);

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
                    <Input
                        placeholder="搜索邀请码/拥有者"
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
                            { value: 'user', label: '用户推荐' },
                            { value: 'player', label: '陪玩师推荐' },
                        ]}
                    />
                    <Select
                        placeholder="状态"
                        value={statusFilter}
                        onChange={setStatusFilter}
                        style={{ width: 100 }}
                        options={[
                            { value: 'all', label: '全部' },
                            { value: 'active', label: '启用' },
                            { value: 'inactive', label: '禁用' },
                            { value: 'expired', label: '已过期' },
                            { value: 'full', label: '已用完' },
                        ]}
                    />
                    <Button type="primary" onClick={handleSearch}>
                        搜索
                    </Button>
                    <Button onClick={handleReset}>
                        重置
                    </Button>
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={openCreateModal}
                    >
                        创建邀请码
                    </Button>
                </Space>
            </Card>

            {/* Table */}
            <Table
                columns={columns}
                dataSource={codes}
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
                scroll={{ x: 1300 }}
                title={() => selectedRowKeys.length > 0 && (
                    <Space style={{ marginBottom: 8 }}>
                        <span>已选择 {selectedRowKeys.length} 项</span>
                        <Button
                            size="small"
                            icon={<CheckOutlined />}
                            onClick={() => handleBatchToggleStatus(true)}
                        >
                            批量启用
                        </Button>
                        <Button
                            size="small"
                            icon={<CloseOutlined />}
                            onClick={() => handleBatchToggleStatus(false)}
                        >
                            批量禁用
                        </Button>
                        <Popconfirm
                            title={`确定要删除选中的 ${selectedRowKeys.length} 个邀请码吗？`}
                            onConfirm={handleBatchDelete}
                        >
                            <Button
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                            >
                                批量删除
                            </Button>
                        </Popconfirm>
                    </Space>
                )}
            />

            {/* Create/Edit Modal */}
            <CodeForm
                visible={formModalVisible}
                code={editingCode}
                onSave={handleSaveCode}
                onCancel={() => {
                    setFormModalVisible(false);
                    setEditingCode(null);
                }}
            />

            {/* Detail Drawer */}
            <Drawer
                title="邀请码详情"
                placement="right"
                size="large"
                onClose={() => setDetailVisible(false)}
                open={detailVisible}
            >
                {currentCode && (
                    <>
                        <div style={{ textAlign: 'center', marginBottom: 24 }}>
                            <Tag
                                color="blue"
                                style={{
                                    fontFamily: 'monospace',
                                    fontSize: 24,
                                    padding: '8px 24px',
                                }}
                            >
                                {currentCode.code}
                            </Tag>
                            <Button
                                type="link"
                                icon={<CopyOutlined />}
                                onClick={() => handleCopyCode(currentCode.code)}
                            >
                                复制
                            </Button>
                        </div>

                        <Descriptions column={2} bordered>
                            <Descriptions.Item label="邀请码ID">{currentCode.id}</Descriptions.Item>
                            <Descriptions.Item label="类型">
                                <Tag color={currentCode.type === 'user' ? 'blue' : 'purple'}>
                                    {getReferralTypeLabel(currentCode.type)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="拥有者">
                                <Space>
                                    <Avatar
                                        src={currentCode.owner?.avatarUrl}
                                        icon={<UserOutlined />}
                                    />
                                    {currentCode.owner?.name || `用户${currentCode.ownerId}`}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={currentCode.isActive ? 'success' : 'default'}>
                                    {currentCode.isActive ? '启用' : '禁用'}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="使用次数">
                                {currentCode.usedCount} / {currentCode.maxUses}
                            </Descriptions.Item>
                            <Descriptions.Item label="使用率">
                                {getCodeUsagePercent(currentCode)}%
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间" span={2}>
                                {dayjs(currentCode.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            <Descriptions.Item label="有效期至" span={2}>
                                <span style={{ color: isCodeExpired(currentCode) ? '#ff4d4f' : undefined }}>
                                    {dayjs(currentCode.expiresAt).format('YYYY-MM-DD HH:mm:ss')}
                                </span>
                                {isCodeExpired(currentCode) && (
                                    <Tag color="error" style={{ marginLeft: 8 }}>已过期</Tag>
                                )}
                            </Descriptions.Item>
                        </Descriptions>

                        <Divider />
                        <Space>
                            <Button
                                type="primary"
                                icon={<EditOutlined />}
                                onClick={() => {
                                    setDetailVisible(false);
                                    openEditModal(currentCode);
                                }}
                            >
                                编辑
                            </Button>
                            <Button
                                icon={currentCode.isActive ? <CloseOutlined /> : <CheckOutlined />}
                                onClick={() => {
                                    handleToggleActive(currentCode);
                                    setDetailVisible(false);
                                }}
                            >
                                {currentCode.isActive ? '禁用' : '启用'}
                            </Button>
                        </Space>
                    </>
                )}
            </Drawer>
        </>
    );
};

export default Codes;
