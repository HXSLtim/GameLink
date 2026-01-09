/**
 * DisputeList Component
 * Displays the list of disputes with filters and actions
 * 
 * Features:
 * - Inline quick assignment via Select dropdown
 * - Quick filter buttons for common scenarios
 * - SLA status indicators
 */
import React, { useState, useEffect, useCallback } from 'react';
import { Tag, Space, Button, Avatar, Typography, Select, App } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    UserSwitchOutlined,
    CheckCircleOutlined,
    RollbackOutlined,
    ClockCircleOutlined,
    UserOutlined,
    ThunderboltOutlined,
} from '@ant-design/icons';
import { SearchTable, type ToolbarButton, type SearchField } from '@/components';
import type { Dispute } from '@/types/dispute';
import {
    DISPUTE_STATUS_LABELS,
    DISPUTE_STATUS_COLORS,
    DISPUTE_TYPE_LABELS,
} from '@/types/dispute';
import { adminApi } from '@/api/admin';
import { disputeApi } from '@/api/dispute';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

export interface DisputeListProps {
    /** Dispute data */
    disputes: Dispute[];
    /** Loading state */
    loading: boolean;
    /** Total count for pagination */
    total: number;
    /** Current page */
    current: number;
    /** Page size */
    pageSize: number;
    /** Search field configurations */
    searchFields: SearchField[];
    /** Toolbar button configurations */
    toolbarButtons: ToolbarButton[];
    /** Search callback */
    onSearch: (values: Record<string, unknown>) => void;
    /** Refresh callback */
    onRefresh: () => void;
    /** Page change callback */
    onPageChange: (page: number, size: number) => void;
    /** View detail callback */
    onViewDetail: (dispute: Dispute) => void;
    /** Assign callback */
    onAssign: (dispute: Dispute) => void;
    /** Resolve callback */
    onResolve: (dispute: Dispute) => void;
    /** Rollback callback */
    onRollback: (dispute: Dispute) => void;
}

interface CSUser {
    id: number;
    name: string;
}

/**
 * DisputeList Component
 */
export const DisputeList: React.FC<DisputeListProps> = ({
    disputes,
    loading,
    total,
    current,
    pageSize,
    searchFields,
    toolbarButtons,
    onSearch,
    onRefresh,
    onPageChange,
    onViewDetail,
    onAssign,
    onResolve,
    onRollback,
}) => {
    const { message } = App.useApp();
    const [csUsers, setCsUsers] = useState<CSUser[]>([]);
    const [assigningId, setAssigningId] = useState<number | null>(null);

    // 加载客服列表
    const loadCSUsers = useCallback(async () => {
        try {
            const response = await adminApi.getUsers({ role: ['admin'], page_size: 100 });
            const users = response.data?.data || [];
            setCsUsers(Array.isArray(users) ? users.map(u => ({ id: u.id, name: u.name })) : []);
        } catch (error) {
            logger.error('Load CS users error:', error);
        }
    }, []);

    useEffect(() => {
        loadCSUsers();
    }, [loadCSUsers]);

    // 行内快速分配
    const handleQuickAssign = async (disputeId: number, csId: number) => {
        setAssigningId(disputeId);
        try {
            await disputeApi.assignDispute(disputeId, { assignedServiceId: csId });
            message.success('分配成功');
            onRefresh();
        } catch (error) {
            logger.error('Quick assign error:', error);
            message.error('分配失败');
        } finally {
            setAssigningId(null);
        }
    };
    /**
     * Table columns configuration
     */
    const columns: ColumnsType<Dispute> = [
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
            width: 180,
            render: (text) => <Typography.Text copyable>{text || '-'}</Typography.Text>,
        },
        {
            title: '发起人',
            key: 'initiator',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} />
                    <span>{record.initiatorName || `ID: ${record.initiatorId}`}</span>
                    <Tag color={record.initiatorType === 'user' ? 'blue' : 'purple'}>
                        {record.initiatorType === 'user' ? '用户' : '陪玩师'}
                    </Tag>
                </Space>
            ),
        },
        {
            title: '纠纷类型',
            dataIndex: 'type',
            key: 'type',
            width: 120,
            render: (type) => DISPUTE_TYPE_LABELS[type as keyof typeof DISPUTE_TYPE_LABELS] || type,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status) => (
                <Tag color={DISPUTE_STATUS_COLORS[status as keyof typeof DISPUTE_STATUS_COLORS] as string}>
                    {DISPUTE_STATUS_LABELS[status as keyof typeof DISPUTE_STATUS_LABELS]}
                </Tag>
            ),
        },
        {
            title: '处理客服',
            key: 'assignedService',
            width: 160,
            render: (_, record) => {
                // 待处理状态：显示快速分配下拉框
                if (record.status === 'pending') {
                    return (
                        <Select
                            placeholder="快速分配"
                            size="small"
                            style={{ width: 140 }}
                            loading={assigningId === record.id}
                            disabled={assigningId !== null}
                            onChange={(csId) => handleQuickAssign(record.id, csId)}
                            options={csUsers.map(u => ({ label: u.name, value: u.id }))}
                            suffixIcon={<ThunderboltOutlined style={{ color: '#faad14' }} />}
                        />
                    );
                }
                // 已分配状态：显示客服信息
                if (record.assignedServiceName) {
                    return (
                        <Space direction="vertical" size={0}>
                            <span>{record.assignedServiceName}</span>
                            {record.originalServiceName && (
                                <span style={{ fontSize: '12px', color: '#999' }}>
                                    原: {record.originalServiceName}
                                </span>
                            )}
                        </Space>
                    );
                }
                return '-';
            },
        },
        {
            title: 'SLA状态',
            key: 'sla',
            width: 100,
            render: (_, record) => {
                if (['resolved', 'rejected', 'canceled'].includes(record.status)) {
                    return <Tag color="default">已完成</Tag>;
                }
                if (record.slaBreached) {
                    return <Tag color="error" icon={<ClockCircleOutlined />}>已超时</Tag>;
                }
                if (record.slaDeadline) {
                    const deadline = dayjs(record.slaDeadline);
                    const now = dayjs();
                    const diff = deadline.diff(now, 'minute');
                    if (diff < 10) {
                        return <Tag color="warning">剩余{diff}分钟</Tag>;
                    }
                    return <Tag color="success">正常</Tag>;
                }
                return <Tag color="default">-</Tag>;
            },
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
            width: 320,
            fixed: 'right',
            render: (_, record) => (
                <Space size={4}>
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => onViewDetail(record)}
                    >
                        详情
                    </Button>
                    {record.status === 'pending' && (
                        <Button
                            type="link"
                            size="small"
                            icon={<UserSwitchOutlined />}
                            onClick={() => onAssign(record)}
                        >
                            分配
                        </Button>
                    )}
                    {['assigned', 'mediating'].includes(record.status) && (
                        <>
                            <Button
                                type="link"
                                size="small"
                                icon={<CheckCircleOutlined />}
                                style={{ color: '#52c41a' }}
                                onClick={() => onResolve(record)}
                            >
                                解决
                            </Button>
                            <Button
                                type="link"
                                size="small"
                                icon={<RollbackOutlined />}
                                onClick={() => onRollback(record)}
                            >
                                回滚
                            </Button>
                        </>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <SearchTable
            columns={columns}
            dataSource={disputes}
            rowKey="id"
            searchFields={searchFields}
            onSearch={onSearch}
            onRefresh={onRefresh}
            loading={loading}
            showCreate={false}
            toolbarButtons={toolbarButtons}
            pagination={{
                current,
                pageSize,
                total,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (t) => `共 ${t} 条`,
                onChange: onPageChange,
            }}
            scroll={{ x: 1400 }}
        />
    );
};

export default DisputeList;
