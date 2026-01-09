/**
 * 分流规则管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Space,
    Button,
    Modal,
    Form,
    Input,
    message,
    Switch,
    Select,
    InputNumber,
    Descriptions,
    Timeline,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
    HistoryOutlined,
    DownloadOutlined,
    ExperimentOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

import { logger } from '@/utils/logger';
interface RoutingRule {
    id: number;
    name: string;
    description: string;
    priority: number;
    status: 'active' | 'inactive';
    targetEntityId: number;
    targetEntityName: string;
    conditions: Record<string, unknown>;
    createdAt: string;
    updatedAt: string;
}

interface RuleHistory {
    id: number;
    ruleId: number;
    action: string;
    changes: string;
    operatorId: number;
    operatorName: string;
    createdAt: string;
}

const statusMap: Record<string, { color: string; text: string }> = {
    active: { color: 'success', text: '启用' },
    inactive: { color: 'default', text: '禁用' },
};

const exportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '规则名称' },
    { key: 'priority', title: '优先级' },
    { key: 'targetEntityName', title: '目标主体' },
    { key: 'status', title: '状态', render: (v) => statusMap[v as string]?.text || String(v) },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];


const RoutingRulePage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [rules, setRules] = useState<RoutingRule[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    const [editVisible, setEditVisible] = useState(false);
    const [currentRule, setCurrentRule] = useState<RoutingRule | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    const [historyVisible, setHistoryVisible] = useState(false);
    const [histories, setHistories] = useState<RuleHistory[]>([]);
    const [historyLoading, setHistoryLoading] = useState(false);

    const [testVisible, setTestVisible] = useState(false);
    const [testForm] = Form.useForm();
    const [testResult, setTestResult] = useState<Record<string, unknown> | null>(null);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get('/admin/routing-rules', {
                params: { page: current, pageSize, ...searchParams },
            });
            if (response.data.success) {
                setRules(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            }
        } catch (error) {
            logger.error('Load error:', error);
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
        form.resetFields();
        setEditVisible(true);
    };

    const handleEdit = (record: RoutingRule) => {
        setCurrentRule(record);
        form.setFieldsValue(record);
        setEditVisible(true);
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            // Transform form data to API format with conditions array
            // Value needs to be JSON-encoded based on field type
            let conditionValue: unknown = values.conditionValue;
            if (values.conditionType === 'order_amount') {
                // For amount, try to parse as number
                conditionValue = parseFloat(values.conditionValue) || values.conditionValue;
            }
            const payload = {
                name: values.name,
                description: values.description || '',
                priority: values.priority,
                targetEntityId: values.targetEntityId,
                conditions: [{
                    field: values.conditionType,
                    operator: 'eq',
                    value: conditionValue,
                }],
            };
            if (currentRule) {
                await apiClient.put(`/admin/routing-rules/${currentRule.id}`, payload);
                message.success('更新成功');
            } else {
                await apiClient.post('/admin/routing-rules', payload);
                message.success('创建成功');
            }
            setEditVisible(false);
            loadData();
        } catch {
            message.error('保存失败');
        } finally {
            setSubmitting(false);
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await apiClient.delete(`/admin/routing-rules/${id}`);
            message.success('删除成功');
            loadData();
        } catch {
            message.error('删除失败');
        }
    };

    const handleToggleStatus = async (record: RoutingRule) => {
        try {
            await apiClient.post(`/admin/routing-rules/${record.id}/toggle`, {
                enabled: record.status !== 'active',
            });
            message.success('状态更新成功');
            loadData();
        } catch {
            message.error('操作失败');
        }
    };

    const handleViewHistory = async (record: RoutingRule) => {
        setHistoryLoading(true);
        setHistoryVisible(true);
        try {
            const response = await apiClient.get(`/admin/routing-rules/${record.id}/history`);
            if (response.data.success) {
                setHistories(response.data.data || []);
            }
        } catch {
            message.error('加载历史失败');
        } finally {
            setHistoryLoading(false);
        }
    };

    const handleTest = async () => {
        try {
            const values = await testForm.validateFields();
            const response = await apiClient.post('/admin/routing-rules/test', values);
            if (response.data.success) {
                setTestResult(response.data.data);
            }
        } catch {
            message.error('测试失败');
        }
    };

    const handleExport = () => {
        exportToCSV(rules as unknown as Record<string, unknown>[], exportColumns, 'routing_rules');
        message.success('导出成功');
    };

    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '规则名称' },
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

    const columns: ColumnsType<RoutingRule> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        { title: '规则名称', dataIndex: 'name', key: 'name', width: 200 },
        { title: '优先级', dataIndex: 'priority', key: 'priority', width: 80, sorter: (a, b) => a.priority - b.priority },
        { title: '目标主体', dataIndex: 'targetEntityName', key: 'targetEntityName', width: 150 },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string, record) => (
                <Switch
                    checked={status === 'active'}
                    onChange={() => handleToggleStatus(record)}
                    checkedChildren="启用"
                    unCheckedChildren="禁用"
                />
            ),
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 240, // 3个按钮 × 80px
            render: (_, record) => (
                <Space size={4}>
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
                    <Button type="link" size="small" icon={<HistoryOutlined />} onClick={() => handleViewHistory(record)}>历史</Button>
                    <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        { text: '测试规则', icon: <ExperimentOutlined />, needSelection: false, onClick: () => setTestVisible(true) },
        { text: '导出数据', icon: <DownloadOutlined />, needSelection: false, onClick: () => handleExport() },
    ];


    return (
        <PageContainer title="分流规则管理" subTitle="管理订单分流规则">
            <SearchTable
                columns={columns}
                dataSource={rules}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={true}
                createText="新增规则"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={{ current, pageSize, total, showSizeChanger: true, showTotal: t => `共 ${t} 条`, onChange: (p, s) => { setCurrent(p); setPageSize(s); } }}
                scroll={{ x: 1000 }}
            />

            <Modal title={currentRule ? '编辑分流规则' : '新增分流规则'} open={editVisible} onOk={handleSave} onCancel={() => setEditVisible(false)} confirmLoading={submitting} width={600}>
                <Form form={form} layout="vertical">
                    <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }]}>
                        <Input placeholder="请输入规则名称" />
                    </Form.Item>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea placeholder="请输入描述" rows={2} />
                    </Form.Item>
                    <Form.Item name="priority" label="优先级" rules={[{ required: true, message: '请输入优先级' }]}>
                        <InputNumber placeholder="数字越小优先级越高" min={1} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item name="targetEntityId" label="目标收款主体ID" rules={[{ required: true, message: '请输入目标主体ID' }]}>
                        <InputNumber placeholder="请输入目标收款主体ID" min={1} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item name="conditionType" label="条件类型" rules={[{ required: true, message: '请选择条件类型' }]}>
                        <Select placeholder="请选择条件类型">
                            <Select.Option value="order_amount">金额范围</Select.Option>
                            <Select.Option value="game_type">游戏类型</Select.Option>
                            <Select.Option value="region">地区</Select.Option>
                            <Select.Option value="service_type">服务类型</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="conditionValue" label="条件值" rules={[{ required: true, message: '请输入条件值' }]}>
                        <Input placeholder="如: 100-500 (金额范围), LOL (游戏), 华东 (地区)" />
                    </Form.Item>
                </Form>
            </Modal>

            <Modal title="修改历史" open={historyVisible} onCancel={() => setHistoryVisible(false)} footer={null} width={600}>
                {historyLoading ? (
                    <div style={{ textAlign: 'center', padding: 20 }}>加载中...</div>
                ) : (
                    <Timeline
                        items={histories.map(h => ({
                            children: (
                                <div>
                                    <div><strong>{h.action}</strong> - {h.operatorName}</div>
                                    <div style={{ color: '#999', fontSize: 12 }}>{dayjs(h.createdAt).format('YYYY-MM-DD HH:mm:ss')}</div>
                                    {h.changes && <div style={{ marginTop: 4, fontSize: 12 }}>{h.changes}</div>}
                                </div>
                            ),
                        }))}
                    />
                )}
            </Modal>

            <Modal title="测试分流规则" open={testVisible} onOk={handleTest} onCancel={() => { setTestVisible(false); setTestResult(null); }} width={500}>
                <Form form={testForm} layout="vertical">
                    <Form.Item name="amount" label="订单金额" rules={[{ required: true }]}>
                        <InputNumber placeholder="请输入订单金额" min={0} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item name="gameId" label="游戏ID">
                        <InputNumber placeholder="请输入游戏ID" min={1} style={{ width: '100%' }} />
                    </Form.Item>
                </Form>
                {testResult && (
                    <Descriptions title="测试结果" bordered size="small" column={1} style={{ marginTop: 16 }}>
                        <Descriptions.Item label="匹配规则">{testResult.matchedRuleName as string || '默认规则'}</Descriptions.Item>
                        <Descriptions.Item label="目标主体">{testResult.targetEntityName as string}</Descriptions.Item>
                    </Descriptions>
                )}
            </Modal>
        </PageContainer>
    );
};

export default RoutingRulePage;
