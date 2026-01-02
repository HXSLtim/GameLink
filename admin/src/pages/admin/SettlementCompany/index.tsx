/**
 * 结算公司管理页面
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
    Card,
    Row,
    Col,
    Statistic,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    BankOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

import { logger } from '@/utils/logger';
interface SettlementCompany {
    id: number;
    name: string;
    creditCode: string;
    taxRegistrationNo: string;
    contactName: string;
    contactPhone: string;
    address: string;
    bankName: string;
    bankAccount: string;
    bankBranch: string;
    status: 'active' | 'inactive';
    playerCount: number;
    createdAt: string;
    updatedAt: string;
}

const statusMap: Record<string, { color: string; text: string }> = {
    active: { color: 'success', text: '启用' },
    inactive: { color: 'default', text: '禁用' },
};

const exportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '公司名称' },
    { key: 'creditCode', title: '统一社会信用代码' },
    { key: 'contactName', title: '联系人' },
    { key: 'contactPhone', title: '联系电话' },
    { key: 'bankName', title: '银行名称' },
    { key: 'playerCount', title: '陪玩师数' },
    { key: 'status', title: '状态', render: (v) => statusMap[v as string]?.text || String(v) },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

const SettlementCompanyPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [companies, setCompanies] = useState<SettlementCompany[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    const [editVisible, setEditVisible] = useState(false);
    const [currentCompany, setCurrentCompany] = useState<SettlementCompany | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get('/admin/settlement-companies', {
                params: { page: current, pageSize, ...searchParams },
            });
            if (response.data.success) {
                setCompanies(response.data.data || []);
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
        setCurrentCompany(null);
        form.resetFields();
        setEditVisible(true);
    };

    const handleEdit = (record: SettlementCompany) => {
        setCurrentCompany(record);
        form.setFieldsValue(record);
        setEditVisible(true);
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            if (currentCompany) {
                await apiClient.put(`/admin/settlement-companies/${currentCompany.id}`, values);
                message.success('更新成功');
            } else {
                await apiClient.post('/admin/settlement-companies', values);
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

    const handleToggleStatus = async (record: SettlementCompany) => {
        try {
            await apiClient.post(`/admin/settlement-companies/${record.id}/toggle`, {
                enabled: record.status !== 'active',
            });
            message.success('状态更新成功');
            loadData();
        } catch {
            message.error('操作失败');
        }
    };

    const handleExport = () => {
        exportToCSV(companies as unknown as Record<string, unknown>[], exportColumns, 'settlement_companies');
        message.success('导出成功');
    };

    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '公司名称/代码' },
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

    const columns: ColumnsType<SettlementCompany> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        { title: '公司名称', dataIndex: 'name', key: 'name', width: 200 },
        { title: '信用代码', dataIndex: 'creditCode', key: 'creditCode', width: 180 },
        { title: '联系人', dataIndex: 'contactName', key: 'contactName', width: 120 },
        { title: '联系电话', dataIndex: 'contactPhone', key: 'contactPhone', width: 140 },
        {
            title: '陪玩师数',
            dataIndex: 'playerCount',
            key: 'playerCount',
            width: 100,
            render: (count: number) => <Tag color="blue">{count || 0}</Tag>,
        },
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
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 120,
            render: (_, record) => (
                <Space size="small">
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        { text: '导出数据', icon: <DownloadOutlined />, needSelection: false, onClick: () => handleExport() },
    ];

    return (
        <PageContainer title="结算公司管理" subTitle="管理陪玩师结算公司">
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={8}>
                    <Card><Statistic title="公司总数" value={total} prefix={<BankOutlined />} /></Card>
                </Col>
                <Col span={8}>
                    <Card><Statistic title="启用公司" value={companies.filter(c => c.status === 'active').length} valueStyle={{ color: '#52c41a' }} /></Card>
                </Col>
                <Col span={8}>
                    <Card><Statistic title="关联陪玩师" value={companies.reduce((sum, c) => sum + (c.playerCount || 0), 0)} /></Card>
                </Col>
            </Row>

            <SearchTable
                columns={columns}
                dataSource={companies}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={true}
                createText="新增公司"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={{ current, pageSize, total, showSizeChanger: true, showTotal: t => `共 ${t} 条`, onChange: (p, s) => { setCurrent(p); setPageSize(s); } }}
                scroll={{ x: 1200 }}
            />

            <Modal title={currentCompany ? '编辑结算公司' : '新增结算公司'} open={editVisible} onOk={handleSave} onCancel={() => setEditVisible(false)} confirmLoading={submitting} width={600}>
                <Form form={form} layout="vertical">
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="name" label="公司名称" rules={[{ required: true, message: '请输入公司名称' }]}>
                                <Input placeholder="请输入公司名称" maxLength={200} />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="creditCode" label="统一社会信用代码" rules={[{ required: true, message: '请输入18位统一社会信用代码' }, { len: 18, message: '必须为18位' }]}>
                                <Input placeholder="请输入18位统一社会信用代码" maxLength={18} disabled={!!currentCompany} />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="contactName" label="联系人">
                                <Input placeholder="请输入联系人" maxLength={50} />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="contactPhone" label="联系电话">
                                <Input placeholder="请输入联系电话" maxLength={20} />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item name="address" label="地址">
                        <Input placeholder="请输入地址" maxLength={500} />
                    </Form.Item>
                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item name="bankName" label="银行名称">
                                <Input placeholder="请输入银行名称" maxLength={100} />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="bankAccount" label="银行账号">
                                <Input placeholder="请输入银行账号" maxLength={30} />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item name="bankBranch" label="开户支行">
                                <Input placeholder="请输入开户支行" maxLength={200} />
                            </Form.Item>
                        </Col>
                    </Row>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default SettlementCompanyPage;
