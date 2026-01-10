/**
 * 纠纷管理页面
 * Manages order disputes between users and players
 *
 * Features:
 * - View and filter disputes by status, order number, initiator type
 * - Assign disputes to customer service representatives (dual-CS mechanism)
 * - Resolve disputes with resolution type (refund, partial, reassign, reject)
 * - View detailed dispute information including SLA status
 * - Rollback dispute assignments
 */
import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
    Drawer,
    Card,
    Row,
    Col,
    Statistic,
    Alert,
    App,
    theme,
    Switch,
    Tooltip,
    Space,
} from 'antd';
import {
    ExclamationCircleOutlined,
    DownloadOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    BellOutlined,
    SoundOutlined,
} from '@ant-design/icons';
import { PageContainer, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { disputeApi } from '@/api/dispute';
import type {
    Dispute,
    DisputeStatus,
    DisputeInitiatorType,
    DisputeStats,
    DisputeResolution,
} from '@/types/dispute';
import { DISPUTE_STATUS_LABELS } from '@/types/dispute';
import { DisputeList } from './components/DisputeList';
import { DisputeDetail } from './components/DisputeDetail';
import { ResolveModal } from './components/ResolveModal';
import { AssignModal } from './components/AssignModal';
import { useNotification } from '@/hooks';

import { logger } from '@/utils/logger';
/**
 * Dispute Query Parameters Interface
 */
interface DisputeQueryParams {
    page?: number;
    pageSize?: number;
    status?: DisputeStatus;
    orderNo?: string;
    initiatorType?: DisputeInitiatorType;
}

/**
 * Dispute Management Page
 */
const DisputePage: React.FC = () => {
    const { message, modal } = App.useApp();
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [disputes, setDisputes] = useState<Dispute[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<DisputeQueryParams>({});

    // Statistics
    const [stats, setStats] = useState<DisputeStats | null>(null);

    // Ref for quick filter callback (to avoid stale closure in useNotification)
    const quickFilterRef = useRef<(filter: string) => void>(() => {});

    // Notification hook for SLA alerts
    const {
        isEnabled: notificationEnabled,
        requestPermission,
        preferences: notificationPrefs,
        updatePreferences: updateNotificationPrefs,
        testSound,
    } = useNotification({
        enableSLAMonitoring: true,
        checkInterval: 60000, // Check every minute
        getSLABreachCount: () => stats?.slaBreached || 0,
        onSLAAlert: () => {
            // Focus on SLA breached disputes
            quickFilterRef.current('sla_breached');
        },
    });

    // Modal states
    const [detailVisible, setDetailVisible] = useState(false);
    const [assignVisible, setAssignVisible] = useState(false);
    const [resolveVisible, setResolveVisible] = useState(false);
    const [currentDispute, setCurrentDispute] = useState<Dispute | null>(null);
    const [submitting, setSubmitting] = useState(false);

    /**
     * Load dispute data
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: DisputeQueryParams = {
                page: current,
                pageSize: pageSize,
                ...searchParams,
            };
            const response = await disputeApi.getDisputes(params);
            const data = response.data?.data;
            // 确保 disputes 是数组
            const disputeList = Array.isArray(data?.disputes) ? data.disputes : [];
            setDisputes(disputeList);
            setTotal(data?.total || 0);
        } catch (error) {
            logger.error('Load disputes error:', error);
            message.error('加载纠纷列表失败');
            // 确保错误时也设置为空数组
            setDisputes([]);
            setTotal(0);
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, message]);

    /**
     * Load dispute statistics
     */
    const loadStats = useCallback(async () => {
        try {
            const response = await disputeApi.getDisputeStats();
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
    const handleSearch = (values: Record<string, unknown>) => {
        const params: DisputeQueryParams = {
            ...values,
        } as DisputeQueryParams;
        setSearchParams(params);
        setCurrent(1);
    };

    /**
     * View dispute detail
     */
    const handleViewDetail = async (dispute: Dispute) => {
        setCurrentDispute(dispute);
        setDetailVisible(true);
    };

    /**
     * Open assign modal
     */
    const handleOpenAssign = (dispute: Dispute) => {
        setCurrentDispute(dispute);
        setAssignVisible(true);
    };

    /**
     * Assign dispute to CS
     */
    const handleAssign = async (assignedServiceId: number, originalServiceId?: number) => {
        if (!currentDispute) return;
        try {
            setSubmitting(true);
            await disputeApi.assignDispute(currentDispute.id, {
                assignedServiceId,
                originalServiceId,
            });
            message.success('分配成功');
            setAssignVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Assign error:', error);
            message.error('分配失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * Rollback assignment
     */
    const handleRollback = async (dispute: Dispute) => {
        modal.confirm({
            title: '确认回滚',
            content: '确定要回滚此纠纷的分配吗？回滚后纠纷将变为待处理状态。',
            okText: '确认',
            cancelText: '取消',
            okButtonProps: { danger: true },
            onOk: async () => {
                try {
                    await disputeApi.rollbackAssignment(dispute.id, {
                        rollbackReason: '管理员手动回滚',
                    });
                    message.success('回滚成功');
                    loadData();
                    loadStats();
                } catch (error) {
                    logger.error('Rollback error:', error);
                    message.error('回滚失败');
                }
            },
        });
    };

    /**
     * Open resolve modal
     */
    const handleOpenResolve = (dispute: Dispute) => {
        setCurrentDispute(dispute);
        setResolveVisible(true);
    };

    /**
     * Resolve dispute
     */
    const handleResolve = async (resolution: string, resolveRemark: string) => {
        if (!currentDispute) return;
        try {
            setSubmitting(true);
            await disputeApi.resolveDispute(currentDispute.id, {
                resolution: resolution as DisputeResolution,
                resolveRemark,
            });
            message.success('处理成功');
            setResolveVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Resolve error:', error);
            message.error('处理失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * Export disputes data
     */
    const handleExport = async () => {
        message.info('导出功能开发中');
    };

    /**
     * Search fields configuration
     */
    const searchFields: SearchField[] = [
        {
            name: 'orderNo',
            label: '订单号',
            type: 'input',
            placeholder: '请输入订单号',
        },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: Object.entries(DISPUTE_STATUS_LABELS).map(([key, val]) => ({
                label: val,
                value: key,
            })),
        },
        {
            name: 'initiatorType',
            label: '发起人类型',
            type: 'select',
            options: [
                { label: '用户', value: 'user' },
                { label: '陪玩师', value: 'player' },
            ],
        },
    ];

    /**
     * Quick filter handlers
     */
    const handleQuickFilter = useCallback((filter: string) => {
        switch (filter) {
            case 'pending':
                setSearchParams({ status: 'pending' as DisputeStatus });
                break;
            case 'sla_breached':
                // 筛选 SLA 超时的纠纷（需要后端支持，这里先用前端筛选）
                setSearchParams({ status: 'pending' as DisputeStatus });
                break;
            case 'today':
                // 今日纠纷筛选（清空筛选条件，后续可添加日期筛选）
                setSearchParams({});
                break;
            case 'all':
                setSearchParams({});
                break;
        }
        setCurrent(1);
    }, []);

    // Update ref when handleQuickFilter changes
    useEffect(() => {
        quickFilterRef.current = handleQuickFilter;
    }, [handleQuickFilter]);

    /**
     * Toolbar buttons with quick filters
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '待处理',
            icon: <ExclamationCircleOutlined />,
            needSelection: false,
            type: stats?.pending ? 'primary' : 'default',
            onClick: () => handleQuickFilter('pending'),
        },
        {
            text: 'SLA超时',
            icon: <ClockCircleOutlined />,
            needSelection: false,
            danger: (stats?.slaBreached || 0) > 0,
            onClick: () => handleQuickFilter('sla_breached'),
        },
        {
            text: '全部',
            needSelection: false,
            onClick: () => handleQuickFilter('all'),
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
        },
    ];

    return (
        <PageContainer title="纠纷管理" subTitle="处理用户与陪玩师之间的订单纠纷">
            {/* Statistics Cards */}
            {stats && (
                <Row gutter={16} style={{ marginBottom: 16 }}>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="待处理"
                                value={stats.pending}
                                valueStyle={{ color: token.colorWarning }}
                                prefix={<ExclamationCircleOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="已指派"
                                value={stats.assigned}
                                valueStyle={{ color: token.colorPrimary }}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="调解中"
                                value={stats.mediating}
                                valueStyle={{ color: token.colorInfo }}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="已解决"
                                value={stats.resolved}
                                valueStyle={{ color: token.colorSuccess }}
                                prefix={<CheckCircleOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="已驳回"
                                value={stats.rejected}
                                valueStyle={{ color: token.colorError }}
                            />
                        </Card>
                    </Col>
                    <Col span={4}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="SLA超时"
                                value={stats.slaBreached}
                                valueStyle={{ color: token.colorError }}
                                prefix={<ClockCircleOutlined />}
                            />
                        </Card>
                    </Col>
                </Row>
            )}

            {/* SLA Breached Warning */}
            {stats && stats.slaBreached > 0 && (
                <Alert
                    message="警告"
                    description={`有 ${stats.slaBreached} 个纠纷已超过30分钟SLA，请尽快处理！`}
                    type="warning"
                    showIcon
                    style={{ marginBottom: 16 }}
                    action={
                        !notificationEnabled && (
                            <Tooltip title="开启浏览器通知，及时收到 SLA 超时提醒">
                                <Switch
                                    checkedChildren={<BellOutlined />}
                                    unCheckedChildren={<BellOutlined />}
                                    checked={false}
                                    onChange={() => requestPermission()}
                                />
                            </Tooltip>
                        )
                    }
                />
            )}

            {/* Notification Settings */}
            <Card size="small" style={{ marginBottom: 16 }}>
                <Space>
                    <span>通知设置：</span>
                    <Tooltip title="浏览器通知">
                        <Switch
                            checkedChildren={<BellOutlined />}
                            unCheckedChildren={<BellOutlined />}
                            checked={notificationPrefs.browserNotificationEnabled && notificationEnabled}
                            onChange={(checked) => {
                                if (checked && !notificationEnabled) {
                                    requestPermission();
                                }
                                updateNotificationPrefs({ browserNotificationEnabled: checked });
                            }}
                        />
                    </Tooltip>
                    <Tooltip title="声音提醒">
                        <Switch
                            checkedChildren={<SoundOutlined />}
                            unCheckedChildren={<SoundOutlined />}
                            checked={notificationPrefs.soundEnabled}
                            onChange={(checked) => updateNotificationPrefs({ soundEnabled: checked })}
                        />
                    </Tooltip>
                    <Tooltip title="测试声音">
                        <a onClick={() => testSound('warning')}>测试</a>
                    </Tooltip>
                </Space>
            </Card>

            <DisputeList
                disputes={disputes}
                loading={loading}
                total={total}
                current={current}
                pageSize={pageSize}
                searchFields={searchFields}
                toolbarButtons={toolbarButtons}
                onSearch={handleSearch}
                onRefresh={loadData}
                onPageChange={(page, size) => {
                    setCurrent(page);
                    setPageSize(size);
                }}
                onViewDetail={handleViewDetail}
                onAssign={handleOpenAssign}
                onResolve={handleOpenResolve}
                onRollback={handleRollback}
            />

            {/* Detail Drawer */}
            <Drawer
                title="纠纷详情"
                open={detailVisible}
                onClose={() => setDetailVisible(false)}
                width={700}
            >
                {currentDispute && <DisputeDetail dispute={currentDispute} />}
            </Drawer>

            {/* Assign Modal */}
            <AssignModal
                open={assignVisible}
                dispute={currentDispute}
                onConfirm={handleAssign}
                onCancel={() => setAssignVisible(false)}
                confirmLoading={submitting}
            />

            {/* Resolve Modal */}
            <ResolveModal
                open={resolveVisible}
                dispute={currentDispute}
                onConfirm={handleResolve}
                onCancel={() => setResolveVisible(false)}
                confirmLoading={submitting}
            />
        </PageContainer>
    );
};

export default DisputePage;
