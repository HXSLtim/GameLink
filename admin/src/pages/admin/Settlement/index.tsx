/**
 * 结算公司管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Card,
    Row,
    Col,
    Statistic,
    Popconfirm,
    Tooltip,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    BankOutlined,
    UserOutlined,
    PlusOutlined,
    HistoryOutlined,
    DeleteOutlined,
    CheckOutlined,
    CloseOutlined,
    TeamOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { settlementApi } from '@/api/settlement';
import type {
    SettlementCompany,
    CompanyType,
    CompanyStatus,
} from '@/api/settlement';
import CompanyForm from './components/CompanyForm';
import dayjs from 'dayjs';

/**
 * 公司类型映射
 */
const companyTypeMap: Record<CompanyType, { text: string; color: string }> = {
    individual: { text: '个人', color: 'blue' },
    company: { text: '企业', color: 'green' },
};

/**
 * 公司状态映射
 */
const companyStatusMap: Record<CompanyStatus, { text: string; color: string }> = {
    active: { text: '启用', color: 'success' },
    suspended: { text: '停用', color: 'error' },
};

/**
 * 结算公司管理页面
 */
const SettlementPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [companies, setCompanies] = useState<SettlementCompany[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [formVisible, setFormVisible] = useState(false);
    const [historyDrawerVisible, setHistoryDrawerVisible] = useState(false);
    const [currentCompany, setCurrentCompany] = useState<SettlementCompany | null>(null);
    const [companyHistory, setCompanyHistory] = useState<any[]>([]);

    // 批量操作状态
    const [selectedCompanyIds, setSelectedCompanyIds] = useState<number[]>([]);
    const [batchStatusVisible, setBatchStatusVisible] = useState(false);
    const [batchDeleteVisible, setBatchDeleteVisible] = useState(false);
    const [batchTarget, setBatchTarget] = useState<'selected' | 'status' | 'all'>('selected');

    /**
     * 加载结算公司数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                pageSize,
                ...searchParams,
            };
            const response = await settlementApi.getSettlementCompanies(queryParams);
            if (response.data.success) {
                setCompanies(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load settlement companies error:', error);
            message.error('加载结算公司列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    /**
     * 打开新增弹窗
     */
    const handleAdd = () => {
        setCurrentCompany(null);
        setFormVisible(true);
    };

    /**
     * 打开编辑弹窗
     */
    const handleEdit = (company: SettlementCompany) => {
        setCurrentCompany(company);
        setFormVisible(true);
    };

    /**
     * 切换状态
     */
    const handleToggleStatus = async (company: SettlementCompany) => {
        try {
            await settlementApi.toggleSettlementCompanyStatus(company.id, {
                enabled: company.status === 'suspended',
            });
            message.success('状态更新成功');
            loadData();
        } catch (error) {
            console.error('Toggle status error:', error);
            message.error('状态更新失败');
        }
    };

    /**
     * 删除公司
     */
    const handleDelete = async (company: SettlementCompany) => {
        try {
            await settlementApi.deleteSettlementCompany(company.id);
            message.success('删除成功');
            loadData();
        } catch (error) {
            console.error('Delete company error:', error);
            message.error('删除失败');
        }
    };

    /**
     * 查看变更历史
     */
    const handleViewHistory = async (company: SettlementCompany) => {
        try {
            const response = await settlementApi.getSettlementCompanyHistory(company.id);
            if (response.data.success) {
                setCompanyHistory(response.data.data || []);
                setCurrentCompany(company);
                setHistoryDrawerVisible(true);
            }
        } catch (error) {
            console.error('Load history error:', error);
            message.error('加载历史记录失败');
        }
    };

    /**
     * 批量修改状态
     */
    const handleBatchStatus = (keys: React.Key[]) => {
        setSelectedCompanyIds(keys ? keys.map(k => Number(k)) : []);
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchStatusVisible(true);
    };

    const submitBatchStatus = async (isActive: boolean) => {
        try {
            let companyIds: number[] = [];

            if (batchTarget === 'selected') {
                companyIds = selectedCompanyIds;
            } else if (batchTarget === 'status') {
                const response = await settlementApi.getSettlementCompanies({
                    status: isActive ? 'suspended' : 'active',
                    pageSize: 1000,
                });
                if (response.data.success && response.data.data) {
                    companyIds = response.data.data.map((c: SettlementCompany) => c.id);
                }
            } else {
                const response = await settlementApi.getSettlementCompanies({ pageSize: 1000 });
                if (response.data.success && response.data.data) {
                    companyIds = response.data.data.map((c: SettlementCompany) => c.id);
                }
            }

            if (companyIds.length === 0) {
                message.warning('没有符合条件的公司');
                return;
            }

            const response = await settlementApi.batchUpdateCompanyStatus({
                companyIds,
                isActive,
            });

            if (response.data.success) {
                const result = response.data.data;
                message.success(`批量修改完成：成功 ${result.successCount}，失败 ${result.failedCount}`);
                setBatchStatusVisible(false);
                loadData();
            }
        } catch (error) {
            console.error('Batch status error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 批量删除
     */
    const handleBatchDelete = (keys: React.Key[]) => {
        setSelectedCompanyIds(keys ? keys.map(k => Number(k)) : []);
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchDeleteVisible(true);
    };

    const submitBatchDelete = async () => {
        try {
            let companyIds: number[] = [];

            if (batchTarget === 'selected') {
                companyIds = selectedCompanyIds;
            } else {
                const response = await settlementApi.getSettlementCompanies({ pageSize: 1000 });
                if (response.data.success && response.data.data) {
                    companyIds = response.data.data.map((c: SettlementCompany) => c.id);
                }
            }

            if (companyIds.length === 0) {
                message.warning('没有符合条件的公司');
                return;
            }

            const response = await settlementApi.batchDeleteCompanies({ companyIds });

            if (response.data.success) {
                const result = response.data.data;
                message.success(`批量删除完成：成功 ${result.successCount}，失败 ${result.failedCount}`);
                setBatchDeleteVisible(false);
                loadData();
            }
        } catch (error) {
            console.error('Batch delete error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '公司名称/税号' },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: Object.entries(companyStatusMap).map(([key, val]) => ({
                label: val.text,
                value: key,
            })),
        },
        {
            name: 'type',
            label: '类型',
            type: 'select',
            options: Object.entries(companyTypeMap).map(([key, val]) => ({
                label: val.text,
                value: key,
            })),
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<SettlementCompany> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '公司名称',
            dataIndex: 'name',
            key: 'name',
            width: 200,
            ellipsis: { showTitle: false },
            render: (name: string) => (
                <Tooltip title={name}>
                    <span>{name}</span>
                </Tooltip>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 100,
            render: (type: CompanyType) => (
                <Tag color={companyTypeMap[type].color}>
                    {companyTypeMap[type].text}
                </Tag>
            ),
        },
        {
            title: '税号',
            dataIndex: 'taxNumber',
            key: 'taxNumber',
            width: 180,
            ellipsis: true,
        },
        {
            title: '联系人',
            key: 'contact',
            width: 150,
            render: (_, record) => (
                <div>
                    <div>{record.contactPerson || '-'}</div>
                    {record.contactPhone && (
                        <div style={{ fontSize: 12, color: '#999' }}>
                            {record.contactPhone}
                        </div>
                    )}
                </div>
            ),
        },
        {
            title: '陪玩师数',
            dataIndex: 'playerCount',
            key: 'playerCount',
            width: 100,
            align: 'center',
            render: (count: number) => (
                <Tag icon={<UserOutlined />} color="blue">
                    {count || 0}
                </Tag>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: CompanyStatus) => (
                <Tag color={companyStatusMap[status].color}>
                    {companyStatusMap[status].text}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
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
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<HistoryOutlined />}
                        onClick={() => handleViewHistory(record)}
                    >
                        历史
                    </Button>
                    <Popconfirm
                        title={record.status === 'active' ? '确定要停用该公司吗？' : '确定要启用该公司吗？'}
                        onConfirm={() => handleToggleStatus(record)}
                    >
                        <Button type="link" size="small">
                            {record.status === 'active' ? '停用' : '启用'}
                        </Button>
                    </Popconfirm>
                    <Popconfirm
                        title="确定要删除该公司吗？"
                        description="删除后关联的陪玩师将变为未分配状态"
                        onConfirm={() => handleDelete(record)}
                        okText="确认"
                        cancelText="取消"
                        okButtonProps={{ danger: true }}
                    >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量启用',
            icon: <CheckOutlined />,
            needSelection: false,
            onClick: () => {
                handleBatchStatus([]);
                setTimeout(() => submitBatchStatus(true), 100);
            },
        },
        {
            text: '批量停用',
            icon: <CloseOutlined />,
            needSelection: false,
            onClick: () => {
                handleBatchStatus([]);
                setTimeout(() => submitBatchStatus(false), 100);
            },
        },
        {
            text: '批量删除',
            icon: <DeleteOutlined />,
            needSelection: false,
            danger: true,
            onClick: (keys) => handleBatchDelete(keys || []),
        },
    ];

    // 统计数据
    const stats = {
        total: total,
        active: companies.filter(c => c.status === 'active').length,
        totalPlayers: companies.reduce((sum, c) => sum + (c.playerCount || 0), 0),
    };

    return (
        <PageContainer
            title="结算公司管理"
            subTitle="管理陪玩师结算公司与归属关系"
            extra={
                <Space size="large">
                    <Statistic title="公司总数" value={stats.total} prefix={<BankOutlined />} />
                    <Statistic title="启用公司" value={stats.active} valueStyle={{ color: '#52c41a' }} />
                    <Statistic title="关联陪玩师" value={stats.totalPlayers} prefix={<TeamOutlined />} />
                </Space>
            }
        >
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
                onCreate={handleAdd}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1400 }}
            />

            {/* 新增/编辑弹窗 */}
            <CompanyForm
                open={formVisible}
                company={currentCompany}
                onOk={() => {
                    setFormVisible(false);
                    loadData();
                }}
                onCancel={() => setFormVisible(false)}
            />

            {/* 历史记录抽屉 */}
            <Modal
                title="变更历史"
                open={historyDrawerVisible}
                onCancel={() => setHistoryDrawerVisible(false)}
                footer={null}
                width={800}
            >
                {companyHistory.length > 0 ? (
                    <div style={{ maxHeight: 500, overflowY: 'auto' }}>
                        {companyHistory.map((item, index) => (
                            <Card
                                key={item.id || index}
                                size="small"
                                style={{ marginBottom: 8 }}
                            >
                                <Row gutter={16}>
                                    <Col span={6}>
                                        <Tag color="blue">{item.fieldName}</Tag>
                                    </Col>
                                    <Col span={7}>
                                        <del style={{ color: '#999' }}>
                                            {item.oldValue || '-'}
                                        </del>
                                    </Col>
                                    <Col span={1} style={{ textAlign: 'center' }}>
                                        →
                                    </Col>
                                    <Col span={7}>
                                        <span style={{ color: '#52c41a', fontWeight: 500 }}>
                                            {item.newValue || '-'}
                                        </span>
                                    </Col>
                                    <Col span={3} style={{ textAlign: 'right', fontSize: 12, color: '#999' }}>
                                        {dayjs(item.changedAt).format('MM-DD HH:mm')}
                                    </Col>
                                </Row>
                                {item.changedByAdmin && (
                                    <div style={{ fontSize: 12, color: '#999', marginTop: 4 }}>
                                        操作人: {item.changedByAdmin.name}
                                    </div>
                                )}
                            </Card>
                        ))}
                    </div>
                ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                        暂无变更记录
                    </div>
                )}
            </Modal>
        </PageContainer>
    );
};

export default SettlementPage;
