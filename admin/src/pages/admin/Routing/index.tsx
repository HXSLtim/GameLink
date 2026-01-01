/**
 * Routing Rule Management Page
 * Payment routing rule list and management
 */
import React, { useState, useCallback, useEffect } from 'react';
import {
    Card,
    Button,
    Space,
    Tag,
    Badge,
    Modal,
    message,
    Tooltip,
    Row,
    Col,
    Statistic,
    Popconfirm,
    Timeline,
    Descriptions,
    theme,
} from 'antd';
import {
    EditOutlined,
    DeleteOutlined,
    ExperimentOutlined,
    HistoryOutlined,
    DragOutlined,
    CheckCircleOutlined,
    StopOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton, type SearchField } from '@/components';
import { routingApi } from '@/api/routing';
import type {
    RoutingRule,
    RoutingRuleHistory,
    RuleStatus,
} from '@/api/routing';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';
import RoutingForm from './components/RoutingForm';

const statusMap: Record<RuleStatus, { color: string; text: string }> = {
    active: { color: 'success', text: '启用' },
    inactive: { color: 'default', text: '禁用' },
};

const exportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '规则名称' },
    { key: 'priority', title: '优先级' },
    { key: 'targetEntity.name', title: '目标主体', render: (v, r) => String((r as unknown as RoutingRule & { targetEntityName?: string }).targetEntityName || v) },
    { key: 'status', title: '状态', render: (v) => statusMap[v as RuleStatus]?.text || String(v) },
    { key: 'description', title: '描述' },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

/**
 * Render condition tags for display
 */
const renderConditions = (conditions: Array<{ field: string; operator: string; value: unknown }>) => {
    if (!conditions || conditions.length === 0) {
        return <Tag color="default">无条件</Tag>;
    }

    const fieldLabels: Record<string, string> = {
        game_type: '游戏',
        service_type: '服务',
        order_amount: '金额',
        region: '地区',
    };

    return (
        <Space size={4} wrap>
            {conditions.slice(0, 3).map((cond, index: number) => (
                <Tag key={index} color="blue" style={{ fontSize: 11 }}>
                    {fieldLabels[cond.field] || cond.field}
                    {cond.operator === 'eq' ? '=' :
                     cond.operator === 'neq' ? '!=' :
                     cond.operator === 'gt' ? '>' :
                     cond.operator === 'lt' ? '<' :
                     cond.operator === 'in' ? '∈' :
                     cond.operator === 'not_in' ? '∉' :
                     cond.operator === 'between' ? '~' : cond.operator}
                    {Array.isArray(cond.value)
                        ? cond.value.length > 1
                            ? `${cond.value[0]}...`
                            : cond.value[0]
                        : cond.value}
                </Tag>
            ))}
            {conditions.length > 3 && (
                <Tag color="default">+{conditions.length - 3}</Tag>
            )}
        </Space>
    );
};

/**
 * Routing Rule Management Page
 */
const RoutingRulePage: React.FC = () => {
    const navigate = useNavigate();
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [rules, setRules] = useState<RoutingRule[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // Form modal state
    const [formVisible, setFormVisible] = useState(false);
    const [currentRule, setCurrentRule] = useState<RoutingRule | null>(null);

    // History modal state
    const [historyVisible, setHistoryVisible] = useState(false);
    const [histories, setHistories] = useState<RoutingRuleHistory[]>([]);
    const [historyLoading, setHistoryLoading] = useState(false);

    // Statistics
    const [stats, setStats] = useState({
        total: 0,
        active: 0,
        inactive: 0,
        defaultEntity: '',
    });

    // Load data
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await routingApi.getRoutingRules({
                page: current,
                page_size: pageSize,
                ...searchParams,
            });
            if (response.data?.data) {
                const data = response.data.data;
                setRules(data);
                setTotal(data.length);

                // Calculate stats from new data
                setStats({
                    total: data.length,
                    active: data.filter(r => r.status === 'active').length,
                    inactive: data.filter(r => r.status === 'inactive').length,
                    defaultEntity: '默认主体',
                });
            }
        } catch {
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleCreate = () => {
        setCurrentRule(null);
        setFormVisible(true);
    };

    const handleEdit = (record: RoutingRule) => {
        setCurrentRule(record);
        setFormVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await routingApi.deleteRoutingRule(id);
            message.success('删除成功');
            loadData();
        } catch {
            message.error('删除失败');
        }
    };

    const handleToggleStatus = async (record: RoutingRule) => {
        try {
            await routingApi.toggleRoutingRuleStatus(
                record.id,
                record.status !== 'active'
            );
            message.success('状态更新成功');
            loadData();
        } catch {
            message.error('操作失败');
        }
    };

    const handleViewHistory = async (record: RoutingRule) => {
        setHistoryVisible(true);
        setHistoryLoading(true);
        try {
            const response = await routingApi.getRoutingRuleHistory(record.id);
            if (response.data?.data) {
                setHistories(response.data.data);
            }
        } catch {
            message.error('加载历史失败');
        } finally {
            setHistoryLoading(false);
        }
    };

    const handleExport = () => {
        exportToCSV(rules as unknown as Record<string, unknown>[], exportColumns, 'routing_rules');
        message.success('导出成功');
    };

    const navigateToTest = () => {
        navigate('/admin/routing/test');
    };

    const searchFields: SearchField[] = [
        {
            name: 'keyword',
            label: '关键词',
            type: 'input',
            placeholder: '规则名称',
        },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: [
                { label: '启用', value: 'active' },
                { label: '禁用', value: 'inactive' },
            ],
        },
    ];

    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 70,
        },
        {
            title: '规则名称',
            dataIndex: 'name',
            key: 'name',
            width: 180,
            render: (text: string, record: RoutingRule) => (
                <Space direction="vertical" size={2}>
                    <span>{text}</span>
                    {record.description && (
                        <span style={{ fontSize: 11, color: '#999' }}>
                            {record.description}
                        </span>
                    )}
                </Space>
            ),
        },
        {
            title: '优先级',
            dataIndex: 'priority',
            key: 'priority',
            width: 80,
            sorter: (a: RoutingRule, b: RoutingRule) => a.priority - b.priority,
            render: (priority: number, record: RoutingRule) => (
                <Badge
                    count={priority}
                    style={{
                        backgroundColor: record.status === 'active' ? '#52c41a' : '#d9d9d9',
                    }}
                />
            ),
        },
        {
            title: '条件',
            dataIndex: 'conditions',
            key: 'conditions',
            width: 200,
            render: (conditions: Array<{ field: string; operator: string; value: unknown }>) => renderConditions(conditions),
        },
        {
            title: '目标主体',
            dataIndex: ['targetEntity', 'name'],
            key: 'targetEntity',
            width: 150,
            render: (name: string, record: RoutingRule) => {
                const entityName = name || (record as RoutingRule & { targetEntityName?: string }).targetEntityName || `主体 #${record.targetEntityId}`;
                return <Tag color="purple">{entityName}</Tag>;
            },
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 90,
            render: (status: RuleStatus, record: RoutingRule) => (
                <Tag
                    color={statusMap[status].color}
                    icon={status === 'active' ? <CheckCircleOutlined /> : <StopOutlined />}
                    style={{ cursor: 'pointer' }}
                    onClick={() => handleToggleStatus(record)}
                >
                    {statusMap[status].text}
                </Tag>
            ),
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 160,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 180,
            fixed: 'right' as const,
            render: (_: unknown, record: RoutingRule) => (
                <Space size="small">
                    <Tooltip title="编辑">
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        />
                    </Tooltip>
                    <Tooltip title="历史">
                        <Button
                            type="link"
                            size="small"
                            icon={<HistoryOutlined />}
                            onClick={() => handleViewHistory(record)}
                        />
                    </Tooltip>
                    <Popconfirm
                        title="确定要删除这条规则吗？"
                        onConfirm={() => handleDelete(record.id)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Tooltip title="删除">
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                            />
                        </Tooltip>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        {
            text: '测试规则',
            icon: <ExperimentOutlined />,
            needSelection: false,
            onClick: navigateToTest,
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: handleExport,
        },
    ];

    return (
        <PageContainer title="支付路由规则" subTitle="管理订单支付路由规则">
            {/* Statistics */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="规则总数"
                            value={stats.total}
                            prefix={<DragOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="启用规则"
                            value={stats.active}
                            valueStyle={{ color: token.colorSuccess }}
                            prefix={<CheckCircleOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="禁用规则"
                            value={stats.inactive}
                            valueStyle={{ color: token.colorTextSecondary }}
                            prefix={<StopOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="默认主体"
                            value={stats.defaultEntity}
                            valueStyle={{ fontSize: 16 }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* Table */}
            <SearchTable
                columns={columns}
                dataSource={rules}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={loadData}
                loading={loading}
                showCreate
                createText="新增规则"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: (t) => `共 ${t} 条`,
                    onChange: (p, s) => {
                        setCurrent(p);
                        setPageSize(s);
                    },
                }}
                scroll={{ x: 1200 }}
            />

            {/* Form Modal */}
            <RoutingForm
                visible={formVisible}
                rule={currentRule}
                onCancel={() => setFormVisible(false)}
                onOk={() => {
                    setFormVisible(false);
                    loadData();
                }}
            />

            {/* History Modal */}
            <Modal
                title="修改历史"
                open={historyVisible}
                onCancel={() => setHistoryVisible(false)}
                footer={null}
                width={700}
            >
                {historyLoading ? (
                    <div style={{ textAlign: 'center', padding: 20 }}>加载中...</div>
                ) : histories.length > 0 ? (
                    <Timeline
                        items={histories.map((h) => ({
                            children: (
                                <div>
                                    <Descriptions size="small" column={1}>
                                        <Descriptions.Item label="字段">
                                            <Tag color="blue">{h.fieldName}</Tag>
                                        </Descriptions.Item>
                                        <Descriptions.Item label="变更">
                                            <Space>
                                                <span style={{ color: '#ff4d4f' }}>{h.oldValue}</span>
                                                <span>-</span>
                                                <span style={{ color: '#52c41a' }}>{h.newValue}</span>
                                            </Space>
                                        </Descriptions.Item>
                                        <Descriptions.Item label="时间">
                                            {dayjs(h.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                                        </Descriptions.Item>
                                    </Descriptions>
                                </div>
                            ),
                        }))}
                    />
                ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                        暂无修改历史
                    </div>
                )}
            </Modal>
        </PageContainer>
    );
};

export default RoutingRulePage;
