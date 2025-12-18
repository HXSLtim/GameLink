/**
 * 排行榜抽成配置管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    message,
    Switch,
    Select,
    InputNumber,
    Card,
    Row,
    Col,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
    PlusOutlined,
    DownloadOutlined,
    TrophyOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

interface RankingRule {
    rankStart: number;
    rankEnd: number;
    commissionRate: number;
}

interface RankingCommissionConfig {
    id: number;
    name: string;
    rankingType: 'income' | 'order_count';
    period: string;
    month: string;
    rules: RankingRule[];
    description: string;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
}

const rankingTypeMap: Record<string, { color: string; text: string }> = {
    income: { color: 'gold', text: '收入排行' },
    order_count: { color: 'blue', text: '订单量排行' },
};

const exportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '配置名称' },
    { key: 'rankingType', title: '排行类型', render: (v) => rankingTypeMap[v as string]?.text || String(v) },
    { key: 'month', title: '月份' },
    { key: 'isActive', title: '状态', render: (v) => v ? '启用' : '禁用' },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];


const RankingCommissionPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [configs, setConfigs] = useState<RankingCommissionConfig[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    const [editVisible, setEditVisible] = useState(false);
    const [currentConfig, setCurrentConfig] = useState<RankingCommissionConfig | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get('/admin/ranking-commission/configs', {
                params: { page: current, pageSize, ...searchParams },
            });
            if (response.data.success) {
                setConfigs(response.data.data?.configs || []);
                setTotal(response.data.data?.total || 0);
            }
        } catch (error) {
            console.error('Load error:', error);
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
        setCurrentConfig(null);
        form.resetFields();
        form.setFieldsValue({ rules: [{ rankStart: 1, rankEnd: 10, commissionRate: 5 }] });
        setEditVisible(true);
    };

    const handleEdit = (record: RankingCommissionConfig) => {
        setCurrentConfig(record);
        form.setFieldsValue({
            ...record,
            rules: record.rules || [{ rankStart: 1, rankEnd: 10, commissionRate: 5 }],
        });
        setEditVisible(true);
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            if (currentConfig) {
                await apiClient.put(`/admin/ranking-commission/configs/${currentConfig.id}`, values);
                message.success('更新成功');
            } else {
                await apiClient.post('/admin/ranking-commission/configs', values);
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
            await apiClient.delete(`/admin/ranking-commission/configs/${id}`);
            message.success('删除成功');
            loadData();
        } catch {
            message.error('删除失败');
        }
    };

    const handleExport = () => {
        exportToCSV(configs as unknown as Record<string, unknown>[], exportColumns, 'ranking_commission');
        message.success('导出成功');
    };

    const searchFields: SearchField[] = [
        { name: 'month', label: '月份', type: 'input', placeholder: 'YYYY-MM' },
        {
            name: 'rankingType',
            label: '排行类型',
            type: 'select',
            options: [
                { label: '收入排行', value: 'income' },
                { label: '订单量排行', value: 'order_count' },
            ],
        },
    ];

    const columns: ColumnsType<RankingCommissionConfig> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        { title: '配置名称', dataIndex: 'name', key: 'name', width: 200 },
        {
            title: '排行类型',
            dataIndex: 'rankingType',
            key: 'rankingType',
            width: 120,
            render: (type: string) => {
                const info = rankingTypeMap[type];
                return info ? <Tag color={info.color}>{info.text}</Tag> : type;
            },
        },
        { title: '月份', dataIndex: 'month', key: 'month', width: 100 },
        {
            title: '规则数',
            dataIndex: 'rules',
            key: 'rules',
            width: 80,
            render: (rules: RankingRule[]) => <Tag color="blue">{rules?.length || 0}</Tag>,
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: 100,
            render: (active: boolean) => (
                <Tag color={active ? 'success' : 'default'}>{active ? '启用' : '禁用'}</Tag>
            ),
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
            width: 150,
            render: (_, record) => (
                <Space size="small">
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                    <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
                        删除
                    </Button>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        { text: '导出数据', icon: <DownloadOutlined />, needSelection: false, onClick: () => handleExport() },
    ];


    return (
        <PageContainer title="排行榜抽成配置" subTitle="管理陪玩师排行榜抽成规则">
            <SearchTable
                columns={columns}
                dataSource={configs}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={true}
                createText="新增配置"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={{ current, pageSize, total, showSizeChanger: true, showTotal: t => `共 ${t} 条`, onChange: (p, s) => { setCurrent(p); setPageSize(s); } }}
                scroll={{ x: 1000 }}
            />

            <Modal title={currentConfig ? '编辑抽成配置' : '新增抽成配置'} open={editVisible} onOk={handleSave} onCancel={() => setEditVisible(false)} confirmLoading={submitting} width={700}>
                <Form form={form} layout="vertical">
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="name" label="配置名称" rules={[{ required: true, message: '请输入配置名称' }]}>
                                <Input placeholder="请输入配置名称" />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="rankingType" label="排行类型" rules={[{ required: true, message: '请选择排行类型' }]}>
                                <Select placeholder="请选择排行类型" disabled={!!currentConfig}>
                                    <Select.Option value="income">收入排行</Select.Option>
                                    <Select.Option value="order_count">订单量排行</Select.Option>
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="month" label="月份" rules={[{ required: true, message: '请输入月份' }]}>
                                <Input placeholder="YYYY-MM" disabled={!!currentConfig} />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="isActive" label="状态" valuePropName="checked">
                                <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea placeholder="请输入描述" rows={2} />
                    </Form.Item>
                    <Card title={<><TrophyOutlined /> 抽成规则</>} size="small">
                        <Form.List name="rules">
                            {(fields, { add, remove }) => (
                                <>
                                    {fields.map(({ key, name, ...restField }) => (
                                        <Row key={key} gutter={8} style={{ marginBottom: 8 }}>
                                            <Col span={6}>
                                                <Form.Item {...restField} name={[name, 'rankStart']} rules={[{ required: true, message: '必填' }]} noStyle>
                                                    <InputNumber placeholder="起始排名" min={1} style={{ width: '100%' }} />
                                                </Form.Item>
                                            </Col>
                                            <Col span={6}>
                                                <Form.Item {...restField} name={[name, 'rankEnd']} rules={[{ required: true, message: '必填' }]} noStyle>
                                                    <InputNumber placeholder="结束排名" min={1} style={{ width: '100%' }} />
                                                </Form.Item>
                                            </Col>
                                            <Col span={8}>
                                                <Form.Item {...restField} name={[name, 'commissionRate']} rules={[{ required: true, message: '必填' }]} noStyle>
                                                    <InputNumber placeholder="抽成比例(%)" min={0} max={100} style={{ width: '100%' }} />
                                                </Form.Item>
                                            </Col>
                                            <Col span={4}>
                                                <Button type="link" danger onClick={() => remove(name)}>删除</Button>
                                            </Col>
                                        </Row>
                                    ))}
                                    <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>添加规则</Button>
                                </>
                            )}
                        </Form.List>
                    </Card>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default RankingCommissionPage;
