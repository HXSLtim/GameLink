/**
 * 权限审计日志页面
 * Requirements: 6.3 - 审计日志查询与筛选
 * Requirements: 6.5 - 审计日志导出
 * 
 * 功能：
 * - 分页显示审计日志
 * - 显示操作前后数据对比
 * - 按时间范围、操作类型、操作者筛选
 * - 导出 CSV 格式
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Tag,
    Space,
    Button,
    DatePicker,
    Select,
    Input,
    Typography,
    Tooltip,
    Modal,
    Descriptions,
    message,
    Row,
    Col,
    Spin,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    DownloadOutlined,
    SearchOutlined,
    ReloadOutlined,
    EyeOutlined,
    UserOutlined,
    ClockCircleOutlined,
    SwapOutlined,
    FileTextOutlined,
} from '@ant-design/icons';
import { motion } from 'framer-motion';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import { auditLogApi } from '@/api/permission';
import { logger } from '@/utils/logger';
import type {
    PermissionAuditLog,
    AuditAction,
    AuditTargetType,
    AuditLogQueryParams,
} from '@/types/permission';

const { Text, Paragraph } = Typography;
const { RangePicker } = DatePicker;

/**
 * 操作类型标签颜色映射
 */
const ACTION_COLORS: Record<AuditAction, string> = {
    permission_create: 'green',
    permission_update: 'blue',
    permission_delete: 'red',
    role_create: 'green',
    role_update: 'blue',
    role_delete: 'red',
    role_permission_assign: 'purple',
    user_role_assign: 'orange',
};

/**
 * 操作类型中文映射
 */
const ACTION_LABELS: Record<AuditAction, string> = {
    permission_create: '创建权限',
    permission_update: '更新权限',
    permission_delete: '删除权限',
    role_create: '创建角色',
    role_update: '更新角色',
    role_delete: '删除角色',
    role_permission_assign: '分配角色权限',
    user_role_assign: '分配用户角色',
};

/**
 * 目标类型中文映射
 */
const TARGET_TYPE_LABELS: Record<AuditTargetType, string> = {
    permission: '权限',
    role: '角色',
    user: '用户',
};

/**
 * 目标类型颜色映射
 */
const TARGET_TYPE_COLORS: Record<AuditTargetType, string> = {
    permission: 'cyan',
    role: 'geekblue',
    user: 'gold',
};

/**
 * 操作类型选项
 */
const ACTION_OPTIONS = Object.entries(ACTION_LABELS).map(([value, label]) => ({
    value,
    label,
}));

/**
 * 目标类型选项
 */
const TARGET_TYPE_OPTIONS = Object.entries(TARGET_TYPE_LABELS).map(([value, label]) => ({
    value,
    label,
}));

/**
 * 格式化 JSON 数据用于显示
 */
const formatJsonData = (jsonStr?: string): Record<string, unknown> | null => {
    if (!jsonStr) return null;
    try {
        return JSON.parse(jsonStr);
    } catch {
        return null;
    }
};

/**
 * 数据对比组件
 * 显示操作前后数据的差异
 */
const DataCompare: React.FC<{
    beforeData?: string;
    afterData?: string;
}> = ({ beforeData, afterData }) => {
    const before = formatJsonData(beforeData);
    const after = formatJsonData(afterData);

    if (!before && !after) {
        return <Text type="secondary">无数据变更</Text>;
    }

    return (
        <Row gutter={16}>
            <Col span={12}>
                <Card
                    size="small"
                    title={<Text type="secondary">变更前</Text>}
                    style={{ backgroundColor: 'rgba(255, 77, 79, 0.05)' }}
                >
                    {before ? (
                        <Paragraph>
                            <pre style={{ 
                                margin: 0, 
                                fontSize: 12, 
                                maxHeight: 300, 
                                overflow: 'auto',
                                whiteSpace: 'pre-wrap',
                                wordBreak: 'break-all',
                            }}>
                                {JSON.stringify(before, null, 2)}
                            </pre>
                        </Paragraph>
                    ) : (
                        <Text type="secondary">无数据（新建操作）</Text>
                    )}
                </Card>
            </Col>
            <Col span={12}>
                <Card
                    size="small"
                    title={<Text type="secondary">变更后</Text>}
                    style={{ backgroundColor: 'rgba(82, 196, 26, 0.05)' }}
                >
                    {after ? (
                        <Paragraph>
                            <pre style={{ 
                                margin: 0, 
                                fontSize: 12, 
                                maxHeight: 300, 
                                overflow: 'auto',
                                whiteSpace: 'pre-wrap',
                                wordBreak: 'break-all',
                            }}>
                                {JSON.stringify(after, null, 2)}
                            </pre>
                        </Paragraph>
                    ) : (
                        <Text type="secondary">无数据（删除操作）</Text>
                    )}
                </Card>
            </Col>
        </Row>
    );
};

/**
 * 审计日志详情弹窗
 */
const AuditLogDetailModal: React.FC<{
    visible: boolean;
    log: PermissionAuditLog | null;
    onClose: () => void;
}> = ({ visible, log, onClose }) => {
    if (!log) return null;

    return (
        <Modal
            title={
                <Space>
                    <FileTextOutlined />
                    审计日志详情
                </Space>
            }
            open={visible}
            onCancel={onClose}
            footer={[
                <Button key="close" onClick={onClose}>
                    关闭
                </Button>,
            ]}
            width={800}
        >
            <Descriptions bordered column={2} size="small" style={{ marginBottom: 16 }}>
                <Descriptions.Item label="日志ID">{log.id}</Descriptions.Item>
                <Descriptions.Item label="操作时间">
                    {dayjs(log.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
                <Descriptions.Item label="操作者">
                    <Space>
                        <UserOutlined />
                        {log.operatorName} (ID: {log.operatorId})
                    </Space>
                </Descriptions.Item>
                <Descriptions.Item label="操作类型">
                    <Tag color={ACTION_COLORS[log.action]}>
                        {ACTION_LABELS[log.action] || log.action}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="目标类型">
                    <Tag color={TARGET_TYPE_COLORS[log.targetType]}>
                        {TARGET_TYPE_LABELS[log.targetType] || log.targetType}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="目标">
                    {log.targetName} (ID: {log.targetId})
                </Descriptions.Item>
                <Descriptions.Item label="IP 地址" span={2}>
                    {log.ipAddress || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="请求ID" span={2}>
                    <Text copyable={{ text: log.requestId || '' }}>
                        {log.requestId || '-'}
                    </Text>
                </Descriptions.Item>
                <Descriptions.Item label="User Agent" span={2}>
                    <Tooltip title={log.userAgent}>
                        <Text ellipsis style={{ maxWidth: 500 }}>
                            {log.userAgent || '-'}
                        </Text>
                    </Tooltip>
                </Descriptions.Item>
            </Descriptions>

            <Card title="数据变更对比" size="small">
                <DataCompare beforeData={log.beforeData} afterData={log.afterData} />
            </Card>
        </Modal>
    );
};

/**
 * 权限审计日志页面组件
 */
const AuditLogPage: React.FC = () => {
    // 数据状态
    const [loading, setLoading] = useState(false);
    const [logs, setLogs] = useState<PermissionAuditLog[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 筛选状态
    const [dateRange, setDateRange] = useState<[Dayjs | null, Dayjs | null] | null>(null);
    const [actionFilter, setActionFilter] = useState<AuditAction | undefined>();
    const [targetTypeFilter, setTargetTypeFilter] = useState<AuditTargetType | undefined>();
    const [operatorKeyword, setOperatorKeyword] = useState<string>('');

    // 详情弹窗状态
    const [detailVisible, setDetailVisible] = useState(false);
    const [selectedLog, setSelectedLog] = useState<PermissionAuditLog | null>(null);

    // 导出状态
    const [exporting, setExporting] = useState(false);

    /**
     * 构建查询参数
     */
    const buildQueryParams = useCallback((): AuditLogQueryParams => {
        const params: AuditLogQueryParams = {
            page: current,
            page_size: pageSize,
        };

        if (dateRange && dateRange[0] && dateRange[1]) {
            params.date_from = dateRange[0].format('YYYY-MM-DD');
            params.date_to = dateRange[1].format('YYYY-MM-DD');
        }

        if (actionFilter) {
            params.action = actionFilter;
        }

        if (targetTypeFilter) {
            params.target_type = targetTypeFilter;
        }

        return params;
    }, [current, pageSize, dateRange, actionFilter, targetTypeFilter]);

    /**
     * 加载审计日志数据
     * Requirements: 6.3 - 审计日志查询
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params = buildQueryParams();
            const res = await auditLogApi.list(params);
            if (res.data.success && res.data.data) {
                // Backend returns data as array with pagination object
                const data = res.data.data;
                const pagination = (res.data as { pagination?: { total?: number } }).pagination;
                if (Array.isArray(data)) {
                    setLogs(data);
                    setTotal(pagination?.total || data.length);
                } else {
                    const { items, totalCount } = data as { items: PermissionAuditLog[]; totalCount: number };
                    setLogs(items || []);
                    setTotal(totalCount || 0);
                }
            }
        } catch (error) {
            logger.error('加载审计日志失败:', error);
            message.error('加载审计日志失败');
        } finally {
            setLoading(false);
        }
    }, [buildQueryParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索处理
     */
    const handleSearch = () => {
        setCurrent(1);
        loadData();
    };

    /**
     * 重置筛选条件
     */
    const handleReset = () => {
        setDateRange(null);
        setActionFilter(undefined);
        setTargetTypeFilter(undefined);
        setOperatorKeyword('');
        setCurrent(1);
    };

    /**
     * 查看详情
     */
    const handleViewDetail = (log: PermissionAuditLog) => {
        setSelectedLog(log);
        setDetailVisible(true);
    };

    /**
     * 导出审计日志
     * Requirements: 6.5 - 导出 CSV 格式
     */
    const handleExport = async () => {
        setExporting(true);
        try {
            const params = {
                date_from: dateRange?.[0]?.format('YYYY-MM-DD'),
                date_to: dateRange?.[1]?.format('YYYY-MM-DD'),
                action: actionFilter,
                target_type: targetTypeFilter,
                format: 'csv' as const,
            };

            const response = await auditLogApi.export(params);
            
            // 创建下载链接
            const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' });
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = `audit_logs_${dayjs().format('YYYYMMDD_HHmmss')}.csv`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);

            message.success('导出成功');
        } catch (error) {
            logger.error('导出审计日志失败:', error);
            message.error('导出失败');
        } finally {
            setExporting(false);
        }
    };

    /**
     * 表格列配置
     */
    const columns: ColumnsType<PermissionAuditLog> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 70,
        },
        {
            title: '操作时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 170,
            render: (time: string) => (
                <Space>
                    <ClockCircleOutlined />
                    {dayjs(time).format('YYYY-MM-DD HH:mm:ss')}
                </Space>
            ),
        },
        {
            title: '操作者',
            dataIndex: 'operatorName',
            key: 'operatorName',
            width: 120,
            render: (name: string, record) => (
                <Tooltip title={`ID: ${record.operatorId}`}>
                    <Space>
                        <UserOutlined />
                        {name}
                    </Space>
                </Tooltip>
            ),
        },
        {
            title: '操作类型',
            dataIndex: 'action',
            key: 'action',
            width: 130,
            render: (action: AuditAction) => (
                <Tag color={ACTION_COLORS[action]}>
                    {ACTION_LABELS[action] || action}
                </Tag>
            ),
        },
        {
            title: '目标类型',
            dataIndex: 'targetType',
            key: 'targetType',
            width: 100,
            render: (type: AuditTargetType) => (
                <Tag color={TARGET_TYPE_COLORS[type]}>
                    {TARGET_TYPE_LABELS[type] || type}
                </Tag>
            ),
        },
        {
            title: '目标',
            key: 'target',
            width: 180,
            ellipsis: true,
            render: (_, record) => (
                <Tooltip title={`ID: ${record.targetId}`}>
                    <Text>{record.targetName}</Text>
                </Tooltip>
            ),
        },
        {
            title: '数据变更',
            key: 'dataChange',
            width: 100,
            align: 'center',
            render: (_, record) => {
                const hasBefore = !!record.beforeData;
                const hasAfter = !!record.afterData;
                
                if (!hasBefore && !hasAfter) {
                    return <Text type="secondary">-</Text>;
                }
                
                return (
                    <Tooltip title="点击查看详情">
                        <Button
                            type="link"
                            size="small"
                            icon={<SwapOutlined />}
                            onClick={() => handleViewDetail(record)}
                        >
                            {hasBefore && hasAfter ? '变更' : hasBefore ? '删除' : '新建'}
                        </Button>
                    </Tooltip>
                );
            },
        },
        {
            title: 'IP 地址',
            dataIndex: 'ipAddress',
            key: 'ipAddress',
            width: 130,
            render: (ip: string) => ip || '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 80,
            fixed: 'right',
            render: (_, record) => (
                <Button
                    type="link"
                    size="small"
                    icon={<EyeOutlined />}
                    onClick={() => handleViewDetail(record)}
                >
                    详情
                </Button>
            ),
        },
    ];

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card bordered={false}>
                {/* 页面标题 */}
                <div style={{ marginBottom: 24 }}>
                    <h2 style={{ margin: 0, fontSize: 20 }}>权限审计日志</h2>
                    <Text type="secondary">
                        查看权限、角色、用户角色分配等操作的审计记录
                    </Text>
                </div>

                {/* 筛选区域 - Requirements: 6.3 */}
                <Card
                    size="small"
                    style={{ marginBottom: 16, backgroundColor: 'rgba(255,255,255,0.02)' }}
                >
                    <Row gutter={[16, 16]} align="middle">
                        <Col>
                            <Space>
                                <Text>时间范围：</Text>
                                <RangePicker
                                    value={dateRange}
                                    onChange={(dates) => setDateRange(dates)}
                                    style={{ width: 260 }}
                                    placeholder={['开始日期', '结束日期']}
                                />
                            </Space>
                        </Col>
                        <Col>
                            <Space>
                                <Text>操作类型：</Text>
                                <Select
                                    value={actionFilter}
                                    onChange={setActionFilter}
                                    placeholder="全部"
                                    allowClear
                                    style={{ width: 150 }}
                                    options={ACTION_OPTIONS}
                                />
                            </Space>
                        </Col>
                        <Col>
                            <Space>
                                <Text>目标类型：</Text>
                                <Select
                                    value={targetTypeFilter}
                                    onChange={setTargetTypeFilter}
                                    placeholder="全部"
                                    allowClear
                                    style={{ width: 120 }}
                                    options={TARGET_TYPE_OPTIONS}
                                />
                            </Space>
                        </Col>
                        <Col>
                            <Space>
                                <Text>操作者：</Text>
                                <Input
                                    value={operatorKeyword}
                                    onChange={(e) => setOperatorKeyword(e.target.value)}
                                    placeholder="搜索操作者"
                                    style={{ width: 150 }}
                                    allowClear
                                />
                            </Space>
                        </Col>
                        <Col flex="auto" style={{ textAlign: 'right' }}>
                            <Space>
                                <Button
                                    icon={<SearchOutlined />}
                                    type="primary"
                                    onClick={handleSearch}
                                >
                                    搜索
                                </Button>
                                <Button
                                    icon={<ReloadOutlined />}
                                    onClick={handleReset}
                                >
                                    重置
                                </Button>
                                <Button
                                    icon={<DownloadOutlined />}
                                    onClick={handleExport}
                                    loading={exporting}
                                >
                                    导出 CSV
                                </Button>
                            </Space>
                        </Col>
                    </Row>
                </Card>

                {/* 数据表格 */}
                <Spin spinning={loading}>
                    <Table
                        columns={columns}
                        dataSource={logs}
                        rowKey="id"
                        pagination={{
                            current,
                            pageSize,
                            total,
                            showSizeChanger: true,
                            showQuickJumper: true,
                            showTotal: (t) => `共 ${t} 条记录`,
                            onChange: (page, size) => {
                                setCurrent(page);
                                setPageSize(size);
                            },
                        }}
                        scroll={{ x: 1200 }}
                        size="middle"
                    />
                </Spin>
            </Card>

            {/* 详情弹窗 */}
            <AuditLogDetailModal
                visible={detailVisible}
                log={selectedLog}
                onClose={() => {
                    setDetailVisible(false);
                    setSelectedLog(null);
                }}
            />
        </motion.div>
    );
};

export default AuditLogPage;
