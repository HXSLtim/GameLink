/**
 * 告警管理页面
 * 系统告警记录、告警规则配置
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    Modal,
    Form,
    Input,
    InputNumber,
    Select,
    Switch,
    message,
    Popconfirm,
    Typography,
    Row,
    Col,
    Statistic,
    Tabs,
    Badge,
    Alert,
    DatePicker,
    theme,
} from 'antd';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    BellOutlined,
    WarningOutlined,
    CloseCircleOutlined,
    ExclamationCircleOutlined,
    SettingOutlined,
    ReloadOutlined,
    CheckOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { monitorApi } from '@/api/monitor';
import type { Alert as AlertType } from '@/types/monitor';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

interface AlertRule {
    id: number;
    name: string;
    type: 'system' | 'business' | 'security';
    metric: string;
    condition: 'gt' | 'lt' | 'eq' | 'gte' | 'lte';
    threshold: number;
    duration: number;
    level: 'critical' | 'warning' | 'info';
    isActive: boolean;
    notifyChannels: string[];
    createdAt: string;
}

// 本地告警规则（后端暂未实现规则管理 API）
const defaultRules: AlertRule[] = [
    { id: 1, name: 'CPU使用率过高', type: 'system', metric: 'cpu_usage', condition: 'gt', threshold: 80, duration: 5, level: 'warning', isActive: true, notifyChannels: ['email', 'sms'], createdAt: '2024-01-01' },
    { id: 2, name: '内存使用率过高', type: 'system', metric: 'memory_usage', condition: 'gt', threshold: 90, duration: 3, level: 'critical', isActive: true, notifyChannels: ['email', 'sms', 'webhook'], createdAt: '2024-01-01' },
    { id: 3, name: '订单异常增长', type: 'business', metric: 'order_rate', condition: 'gt', threshold: 200, duration: 10, level: 'warning', isActive: true, notifyChannels: ['email'], createdAt: '2024-02-15' },
    { id: 4, name: '支付失败率过高', type: 'business', metric: 'payment_fail_rate', condition: 'gt', threshold: 5, duration: 5, level: 'critical', isActive: true, notifyChannels: ['email', 'sms'], createdAt: '2024-02-15' },
    { id: 5, name: '异常登录检测', type: 'security', metric: 'login_fail_count', condition: 'gt', threshold: 10, duration: 1, level: 'warning', isActive: true, notifyChannels: ['email', 'sms'], createdAt: '2024-03-01' },
];

const AdminAlert: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [rules, setRules] = useState<AlertRule[]>(defaultRules);
    const [alerts, setAlerts] = useState<AlertType[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [modalVisible, setModalVisible] = useState(false);
    const [editingRule, setEditingRule] = useState<AlertRule | null>(null);
    const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
    const [form] = Form.useForm();

    // 筛选条件
    const [filters, setFilters] = useState<{
        level?: string;
        type?: string;
        isRead?: boolean;
        dateRange?: [dayjs.Dayjs, dayjs.Dayjs];
    }>({});

    const loadAlerts = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = {
                page: current,
                page_size: pageSize,
            };
            if (filters.level) params.level = filters.level;
            if (filters.type) params.type = filters.type;
            if (filters.isRead !== undefined) params.is_read = filters.isRead;
            if (filters.dateRange) {
                params.date_from = filters.dateRange[0].format('YYYY-MM-DD');
                params.date_to = filters.dateRange[1].format('YYYY-MM-DD');
            }

            const res = await monitorApi.getAlerts(params);
            if (res.data?.success) {
                setAlerts(res.data.data || []);
                setTotal(res.data.pagination?.total || 0);
            } else {
                message.error(res.data?.message || '加载失败');
            }
        } catch (err) {
            logger.error('加载告警失败:', err);
            message.error('加载告警数据失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, filters]);

    useEffect(() => {
        loadAlerts();
    }, [loadAlerts]);

    const handleMarkRead = async (id: string) => {
        try {
            const res = await monitorApi.markAlertRead(id);
            if (res.data?.success) {
                message.success('已标记为已读');
                loadAlerts();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            logger.error('标记已读失败:', err);
            message.error('操作失败');
        }
    };

    const handleBatchMarkRead = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请选择要标记的告警');
            return;
        }
        try {
            const res = await monitorApi.markAlertsRead(selectedRowKeys.map(String));
            if (res.data?.success) {
                message.success(`已标记 ${selectedRowKeys.length} 条告警为已读`);
                setSelectedRowKeys([]);
                loadAlerts();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            logger.error('批量标记已读失败:', err);
            message.error('操作失败');
        }
    };

    // 规则管理（本地状态，后端暂未实现）
    const handleAddRule = () => {
        setEditingRule(null);
        form.resetFields();
        setModalVisible(true);
    };

    const handleEditRule = (record: AlertRule) => {
        setEditingRule(record);
        form.setFieldsValue(record);
        setModalVisible(true);
    };

    const handleDeleteRule = (id: number) => {
        setRules(rules.filter(r => r.id !== id));
        message.success('删除成功');
    };

    const handleSubmitRule = (values: Partial<AlertRule>) => {
        if (editingRule) {
            setRules(rules.map(r => r.id === editingRule.id ? { ...r, ...values } : r));
            message.success('更新成功');
        } else {
            const newRule: AlertRule = {
                ...values,
                id: Math.max(...rules.map(r => r.id)) + 1,
                createdAt: dayjs().format('YYYY-MM-DD'),
            } as AlertRule;
            setRules([...rules, newRule]);
            message.success('创建成功');
        }
        setModalVisible(false);
    };

    const handleRuleStatusChange = (id: number, isActive: boolean) => {
        setRules(rules.map(r => r.id === id ? { ...r, isActive } : r));
        message.success(isActive ? '已启用' : '已禁用');
    };

    const getLevelTag = (level: string) => {
        const config: Record<string, { color: string; icon: React.ReactNode; text: string }> = {
            high: { color: 'red', icon: <CloseCircleOutlined />, text: '严重' },
            critical: { color: 'red', icon: <CloseCircleOutlined />, text: '严重' },
            medium: { color: 'orange', icon: <ExclamationCircleOutlined />, text: '警告' },
            warning: { color: 'orange', icon: <ExclamationCircleOutlined />, text: '警告' },
            low: { color: 'blue', icon: <BellOutlined />, text: '信息' },
            info: { color: 'blue', icon: <BellOutlined />, text: '信息' },
        };
        const cfg = config[level] || config.info;
        return <Tag color={cfg.color} icon={cfg.icon}>{cfg.text}</Tag>;
    };

    const getTypeTag = (type: string) => {
        const config: Record<string, { color: string; text: string }> = {
            system: { color: 'purple', text: '系统' },
            business: { color: 'cyan', text: '业务' },
            security: { color: 'red', text: '安全' },
        };
        const cfg = config[type] || { color: 'default', text: type };
        return <Tag color={cfg.color}>{cfg.text}</Tag>;
    };

    const getConditionText = (condition: AlertRule['condition']) => {
        const map = { gt: '>', lt: '<', eq: '=', gte: '>=', lte: '<=' };
        return map[condition];
    };

    const alertColumns: ColumnsType<AlertType> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
        { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
        { title: '类型', dataIndex: 'type', key: 'type', width: 80, render: (type) => getTypeTag(type) },
        { title: '级别', dataIndex: 'level', key: 'level', width: 80, render: (level) => getLevelTag(level) },
        { title: '告警信息', dataIndex: 'message', key: 'message', ellipsis: true },
        { title: '来源', dataIndex: 'source', key: 'source', width: 100 },
        { title: '时间', dataIndex: 'createdAt', key: 'createdAt', width: 170, render: (t) => dayjs(t).format('YYYY-MM-DD HH:mm:ss') },
        {
            title: '状态',
            dataIndex: 'isRead',
            key: 'isRead',
            width: 80,
            render: (isRead) => isRead ? <Tag color="default">已读</Tag> : <Tag color="red">未读</Tag>,
        },
        {
            title: '操作',
            key: 'action',
            width: 100,
            render: (_, record) => (
                <Space>
                    {!record.isRead && (
                        <Button type="link" size="small" icon={<CheckOutlined />} onClick={() => handleMarkRead(String(record.id))}>
                            已读
                        </Button>
                    )}
                </Space>
            ),
        },
    ];

    const ruleColumns: ColumnsType<AlertRule> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
        { title: '规则名称', dataIndex: 'name', key: 'name' },
        { title: '类型', dataIndex: 'type', key: 'type', render: (type) => getTypeTag(type) },
        { title: '级别', dataIndex: 'level', key: 'level', render: (level) => getLevelTag(level) },
        {
            title: '条件',
            key: 'condition',
            render: (_, record) => (
                <Text code>{record.metric} {getConditionText(record.condition)} {record.threshold}</Text>
            ),
        },
        { title: '持续时间', dataIndex: 'duration', key: 'duration', render: (d) => `${d} 分钟` },
        {
            title: '通知渠道',
            dataIndex: 'notifyChannels',
            key: 'notifyChannels',
            render: (channels) => <Space>{channels.map((c: string) => <Tag key={c}>{c}</Tag>)}</Space>,
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            render: (isActive, record) => (
                <Switch
                    checked={isActive}
                    onChange={(checked) => handleRuleStatusChange(record.id, checked)}
                    checkedChildren="启用"
                    unCheckedChildren="禁用"
                />
            ),
        },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
                <Space>
                    <Button type="link" icon={<EditOutlined />} onClick={() => handleEditRule(record)}>编辑</Button>
                    <Popconfirm title="确定删除此规则？" onConfirm={() => handleDeleteRule(record.id)} okText="确定" cancelText="取消">
                        <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const unreadCount = alerts.filter(a => !a.isRead).length;

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}><BellOutlined /> 告警管理</Title>

            {/* 活跃告警提示 */}
            {unreadCount > 0 && (
                <Alert
                    type="warning"
                    showIcon
                    icon={<WarningOutlined />}
                    style={{ marginBottom: 16 }}
                    message={`当前有 ${unreadCount} 条未读告警`}
                />
            )}

            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}><Statistic title="告警总数" value={total} /></Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="未读告警"
                            value={unreadCount}
                            valueStyle={{ color: unreadCount > 0 ? token.colorError : token.colorSuccess }}
                            prefix={<Badge status={unreadCount > 0 ? 'error' : 'success'} />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}><Statistic title="告警规则" value={rules.length} /></Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="启用规则"
                            value={rules.filter(r => r.isActive).length}
                            valueStyle={{ color: token.colorSuccess }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* 告警记录和规则 */}
            <Card>
                <Tabs
                    defaultActiveKey="alerts"
                    items={[
                        {
                            key: 'alerts',
                            label: <><BellOutlined /> 告警记录 <Badge count={unreadCount} /></>,
                            children: (
                                <>
                                    {/* 筛选条件 */}
                                    <Row gutter={16} style={{ marginBottom: 16 }}>
                                        <Col flex="auto">
                                            <Space wrap>
                                                <Select
                                                    placeholder="告警级别"
                                                    allowClear
                                                    style={{ width: 120 }}
                                                    value={filters.level}
                                                    onChange={(v) => { setFilters({ ...filters, level: v }); setCurrent(1); }}
                                                    options={[
                                                        { value: 'high', label: '严重' },
                                                        { value: 'medium', label: '警告' },
                                                        { value: 'low', label: '信息' },
                                                    ]}
                                                />
                                                <Select
                                                    placeholder="告警类型"
                                                    allowClear
                                                    style={{ width: 120 }}
                                                    value={filters.type}
                                                    onChange={(v) => { setFilters({ ...filters, type: v }); setCurrent(1); }}
                                                    options={[
                                                        { value: 'system', label: '系统' },
                                                        { value: 'business', label: '业务' },
                                                        { value: 'security', label: '安全' },
                                                    ]}
                                                />
                                                <Select
                                                    placeholder="阅读状态"
                                                    allowClear
                                                    style={{ width: 120 }}
                                                    value={filters.isRead}
                                                    onChange={(v) => { setFilters({ ...filters, isRead: v }); setCurrent(1); }}
                                                    options={[
                                                        { value: false, label: '未读' },
                                                        { value: true, label: '已读' },
                                                    ]}
                                                />
                                                <RangePicker
                                                    value={filters.dateRange}
                                                    onChange={(dates) => {
                                                        setFilters({ ...filters, dateRange: dates as [dayjs.Dayjs, dayjs.Dayjs] });
                                                        setCurrent(1);
                                                    }}
                                                />
                                            </Space>
                                        </Col>
                                        <Col>
                                            <Space>
                                                <Button icon={<ReloadOutlined />} onClick={loadAlerts}>刷新</Button>
                                                <Button
                                                    icon={<CheckOutlined />}
                                                    onClick={handleBatchMarkRead}
                                                    disabled={selectedRowKeys.length === 0}
                                                >
                                                    批量已读
                                                </Button>
                                            </Space>
                                        </Col>
                                    </Row>
                                    <Table
                                        columns={alertColumns}
                                        dataSource={alerts}
                                        rowKey="id"
                                        loading={loading}
                                        rowSelection={{
                                            selectedRowKeys,
                                            onChange: setSelectedRowKeys,
                                        }}
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
                                </>
                            ),
                        },
                        {
                            key: 'rules',
                            label: <><SettingOutlined /> 告警规则</>,
                            children: (
                                <>
                                    <div style={{ marginBottom: 16 }}>
                                        <Button type="primary" icon={<PlusOutlined />} onClick={handleAddRule}>
                                            新增规则
                                        </Button>
                                        <Text type="secondary" style={{ marginLeft: 16 }}>
                                            注：告警规则目前为本地配置，后端 API 开发中
                                        </Text>
                                    </div>
                                    <Table columns={ruleColumns} dataSource={rules} rowKey="id" pagination={{ pageSize: 10 }} />
                                </>
                            ),
                        },
                    ]}
                />
            </Card>

            {/* 规则编辑弹窗 */}
            <Modal
                title={editingRule ? '编辑告警规则' : '新增告警规则'}
                open={modalVisible}
                onCancel={() => setModalVisible(false)}
                footer={null}
                width={600}
            >
                <Form form={form} layout="vertical" onFinish={handleSubmitRule}>
                    <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }]}>
                        <Input placeholder="如：CPU使用率过高" />
                    </Form.Item>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="type" label="告警类型" rules={[{ required: true, message: '请选择类型' }]}>
                                <Select
                                    placeholder="选择类型"
                                    options={[
                                        { value: 'system', label: '系统告警' },
                                        { value: 'business', label: '业务告警' },
                                        { value: 'security', label: '安全告警' },
                                    ]}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="level" label="告警级别" rules={[{ required: true, message: '请选择级别' }]}>
                                <Select
                                    placeholder="选择级别"
                                    options={[
                                        { value: 'critical', label: '严重' },
                                        { value: 'warning', label: '警告' },
                                        { value: 'info', label: '信息' },
                                    ]}
                                />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item name="metric" label="监控指标" rules={[{ required: true, message: '请输入指标' }]}>
                                <Input placeholder="如：cpu_usage" />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="condition" label="条件" rules={[{ required: true, message: '请选择条件' }]}>
                                <Select
                                    placeholder="选择"
                                    options={[
                                        { value: 'gt', label: '>' },
                                        { value: 'gte', label: '>=' },
                                        { value: 'lt', label: '<' },
                                        { value: 'lte', label: '<=' },
                                        { value: 'eq', label: '=' },
                                    ]}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="threshold" label="阈值" rules={[{ required: true, message: '请输入阈值' }]}>
                                <InputNumber style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="duration" label="持续时间（分钟）" rules={[{ required: true, message: '请输入持续时间' }]}>
                                <InputNumber min={1} style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="notifyChannels" label="通知渠道" rules={[{ required: true, message: '请选择通知渠道' }]}>
                                <Select
                                    mode="multiple"
                                    placeholder="选择通知渠道"
                                    options={[
                                        { value: 'email', label: '邮件' },
                                        { value: 'sms', label: '短信' },
                                        { value: 'webhook', label: 'Webhook' },
                                    ]}
                                />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item name="isActive" label="状态" valuePropName="checked" initialValue={true}>
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                    <Form.Item>
                        <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
                            <Button onClick={() => setModalVisible(false)}>取消</Button>
                            <Button type="primary" htmlType="submit">保存</Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default AdminAlert;
