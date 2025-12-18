/**
 * 纠纷管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    Select,
    message,
    Drawer,
    Descriptions,
    Timeline,
    Card,
    Row,
    Col,
    Statistic,
    InputNumber,
    Avatar,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    UserSwitchOutlined,
    CheckCircleOutlined,
    RollbackOutlined,
    ExclamationCircleOutlined,
    UserOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

/**
 * 纠纷数据接口
 */
interface Dispute {
    id: number;
    orderId: number;
    orderNo: string;
    userId: number;
    playerId: number;
    reason: string;
    description: string;
    status: 'pending' | 'assigned' | 'processing' | 'resolved' | 'closed';
    assignedToUserId?: number;
    assignedAt?: string;
    resolution?: string;
    resolutionAmount?: number;
    resolutionNotes?: string;
    resolvedAt?: string;
    createdAt: string;
    updatedAt: string;
    user?: { id: number; name: string; avatarUrl?: string };
    player?: { id: number; nickname: string; user?: { avatarUrl?: string } };
    assignedTo?: { id: number; name: string };
    order?: { id: number; orderNo: string; totalPriceCents: number };
}

/**
 * 状态映射
 */
const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'orange', text: '待处理' },
    assigned: { color: 'blue', text: '已分配' },
    processing: { color: 'processing', text: '处理中' },
    resolved: { color: 'success', text: '已解决' },
    closed: { color: 'default', text: '已关闭' },
};

/**
 * 解决方案映射
 */
const resolutionMap: Record<string, string> = {
    refund: '全额退款',
    partial: '部分退款',
    reassign: '重新分配',
    reject: '驳回申诉',
};

/**
 * 导出列配置
 */
const disputeExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'orderNo', title: '订单号' },
    { key: 'user.name', title: '用户' },
    { key: 'player.nickname', title: '陪玩师' },
    { key: 'reason', title: '纠纷原因' },
    { key: 'status', title: '状态', render: (v) => statusMap[v as string]?.text || String(v) },
    { key: 'assignedTo.name', title: '处理人' },
    { key: 'resolution', title: '解决方案', render: (v) => resolutionMap[v as string] || String(v) },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
    { key: 'resolvedAt', title: '解决时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

const DisputePage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [disputes, setDisputes] = useState<Dispute[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 统计数据
    const [stats, setStats] = useState({ pending: 0, processing: 0, resolved: 0, total: 0 });

    // 弹窗状态
    const [detailVisible, setDetailVisible] = useState(false);
    const [assignVisible, setAssignVisible] = useState(false);
    const [resolveVisible, setResolveVisible] = useState(false);
    const [currentDispute, setCurrentDispute] = useState<Dispute | null>(null);
    const [assignForm] = Form.useForm();
    const [resolveForm] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    // 客服列表
    const [staffList, setStaffList] = useState<{ id: number; name: string }[]>([]);

    /**
     * 加载纠纷数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params = {
                page: current,
                pageSize: pageSize,
                ...searchParams,
            };
            const response = await apiClient.get('/admin/disputes', { params });
            if (response.data.success) {
                setDisputes(response.data.data?.disputes || []);
                setTotal(response.data.data?.total || 0);
            }
        } catch (error) {
            console.error('Load disputes error:', error);
            message.error('加载纠纷列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    /**
     * 加载统计数据
     */
    const loadStats = useCallback(async () => {
        try {
            const response = await apiClient.get('/admin/disputes/stats');
            if (response.data.success) {
                const data = response.data.data || {};
                setStats({
                    pending: data.pending || 0,
                    processing: (data.assigned || 0) + (data.mediating || 0),
                    resolved: data.resolved || 0,
                    total: data.total || 0,
                });
            }
        } catch (error) {
            console.error('Load stats error:', error);
        }
    }, []);

    /**
     * 加载客服列表
     */
    const loadStaffList = useCallback(async () => {
        try {
            const response = await apiClient.get('/admin/users', { params: { role: 'admin', page_size: 100 } });
            if (response.data.success) {
                setStaffList(response.data.data?.map((u: { id: number; name: string }) => ({ id: u.id, name: u.name })) || []);
            }
        } catch {
            // 静默失败
        }
    }, []);

    useEffect(() => {
        loadData();
        loadStats();
        loadStaffList();
    }, [loadData, loadStats, loadStaffList]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleViewDetail = (record: Dispute) => {
        setCurrentDispute(record);
        setDetailVisible(true);
    };

    const handleOpenAssign = (record: Dispute) => {
        setCurrentDispute(record);
        assignForm.resetFields();
        setAssignVisible(true);
    };

    const handleAssign = async () => {
        if (!currentDispute) return;
        try {
            const values = await assignForm.validateFields();
            setSubmitting(true);
            await apiClient.post(`/admin/disputes/${currentDispute.id}/assign`, {
                assignedToUserId: values.assignedToUserId,
                source: 'manual',
            });
            message.success('分配成功');
            setAssignVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            console.error('Assign error:', error);
            message.error('分配失败');
        } finally {
            setSubmitting(false);
        }
    };

    const handleOpenResolve = (record: Dispute) => {
        setCurrentDispute(record);
        resolveForm.resetFields();
        resolveForm.setFieldsValue({
            resolutionAmount: record.order?.totalPriceCents ? record.order.totalPriceCents / 100 : 0,
        });
        setResolveVisible(true);
    };

    const handleResolve = async () => {
        if (!currentDispute) return;
        try {
            const values = await resolveForm.validateFields();
            setSubmitting(true);
            await apiClient.post(`/admin/disputes/${currentDispute.id}/resolve`, {
                resolution: values.resolution,
                resolutionAmount: Math.round((values.resolutionAmount || 0) * 100),
                resolutionNotes: values.resolutionNotes,
            });
            message.success('处理成功');
            setResolveVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            console.error('Resolve error:', error);
            message.error('处理失败');
        } finally {
            setSubmitting(false);
        }
    };

    const handleRollback = async (record: Dispute) => {
        Modal.confirm({
            title: '确认回滚',
            content: '确定要回滚此纠纷的分配吗？',
            onOk: async () => {
                try {
                    await apiClient.post(`/admin/disputes/${record.id}/rollback`, {
                        rollbackReason: '管理员手动回滚',
                    });
                    message.success('回滚成功');
                    loadData();
                } catch {
                    message.error('回滚失败');
                }
            },
        });
    };

    const handleExport = async () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            exportToCSV(disputes as unknown as Record<string, unknown>[], disputeExportColumns, 'disputes');
            message.success({ content: '导出成功', key: 'export' });
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    };

    const searchFields: SearchField[] = [
        { name: 'orderNo', label: '订单号', type: 'input', placeholder: '请输入订单号' },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: Object.entries(statusMap).map(([k, v]) => ({ label: v.text, value: k })),
        },
    ];

    const columns: ColumnsType<Dispute> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            width: 180,
            render: (_, record) => record.order?.orderNo || record.orderNo || '-',
        },
        {
            title: '用户',
            key: 'user',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={record.user?.avatarUrl} />
                    <span>{record.user?.name || `用户${record.userId}`}</span>
                </Space>
            ),
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={record.player?.user?.avatarUrl} style={{ backgroundColor: '#722ed1' }} />
                    <span>{record.player?.nickname || `陪玩师${record.playerId}`}</span>
                </Space>
            ),
        },
        {
            title: '纠纷原因',
            dataIndex: 'reason',
            key: 'reason',
            width: 150,
            ellipsis: true,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => (
                <Tag color={statusMap[status]?.color}>{statusMap[status]?.text || status}</Tag>
            ),
        },
        {
            title: '处理人',
            key: 'assignedTo',
            width: 120,
            render: (_, record) => record.assignedTo?.name || '-',
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
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record)}>
                        详情
                    </Button>
                    {record.status === 'pending' && (
                        <Button type="link" size="small" icon={<UserSwitchOutlined />} onClick={() => handleOpenAssign(record)}>
                            分配
                        </Button>
                    )}
                    {['assigned', 'processing'].includes(record.status) && (
                        <>
                            <Button type="link" size="small" icon={<CheckCircleOutlined />} onClick={() => handleOpenResolve(record)}>
                                处理
                            </Button>
                            <Button type="link" size="small" icon={<RollbackOutlined />} onClick={() => handleRollback(record)}>
                                回滚
                            </Button>
                        </>
                    )}
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
        },
    ];

    return (
        <PageContainer title="纠纷管理" subTitle="处理用户与陪玩师之间的订单纠纷">
            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={6}>
                    <Card>
                        <Statistic title="待处理" value={stats.pending} valueStyle={{ color: '#faad14' }} prefix={<ExclamationCircleOutlined />} />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic title="处理中" value={stats.processing} valueStyle={{ color: '#1890ff' }} />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic title="已解决" value={stats.resolved} valueStyle={{ color: '#52c41a' }} />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic title="总计" value={stats.total} />
                    </Card>
                </Col>
            </Row>

            <SearchTable
                columns={columns}
                dataSource={disputes}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: t => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1400 }}
            />

            {/* 详情抽屉 */}
            <Drawer title="纠纷详情" open={detailVisible} onClose={() => setDetailVisible(false)} width={600}>
                {currentDispute && (
                    <>
                        <Descriptions column={2} bordered size="small">
                            <Descriptions.Item label="纠纷ID">{currentDispute.id}</Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={statusMap[currentDispute.status]?.color}>{statusMap[currentDispute.status]?.text}</Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="订单号">{currentDispute.order?.orderNo || '-'}</Descriptions.Item>
                            <Descriptions.Item label="订单金额">
                                ¥{((currentDispute.order?.totalPriceCents || 0) / 100).toFixed(2)}
                            </Descriptions.Item>
                            <Descriptions.Item label="用户">{currentDispute.user?.name || '-'}</Descriptions.Item>
                            <Descriptions.Item label="陪玩师">{currentDispute.player?.nickname || '-'}</Descriptions.Item>
                            <Descriptions.Item label="纠纷原因" span={2}>{currentDispute.reason}</Descriptions.Item>
                            <Descriptions.Item label="详细描述" span={2}>{currentDispute.description || '-'}</Descriptions.Item>
                            <Descriptions.Item label="处理人">{currentDispute.assignedTo?.name || '-'}</Descriptions.Item>
                            <Descriptions.Item label="分配时间">
                                {currentDispute.assignedAt ? dayjs(currentDispute.assignedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                            </Descriptions.Item>
                            {currentDispute.resolution && (
                                <>
                                    <Descriptions.Item label="解决方案">{resolutionMap[currentDispute.resolution] || currentDispute.resolution}</Descriptions.Item>
                                    <Descriptions.Item label="退款金额">
                                        ¥{((currentDispute.resolutionAmount || 0) / 100).toFixed(2)}
                                    </Descriptions.Item>
                                    <Descriptions.Item label="处理备注" span={2}>{currentDispute.resolutionNotes || '-'}</Descriptions.Item>
                                    <Descriptions.Item label="解决时间" span={2}>
                                        {currentDispute.resolvedAt ? dayjs(currentDispute.resolvedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                                    </Descriptions.Item>
                                </>
                            )}
                            <Descriptions.Item label="创建时间">{dayjs(currentDispute.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
                            <Descriptions.Item label="更新时间">{dayjs(currentDispute.updatedAt).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
                        </Descriptions>

                        <Card title="处理进度" size="small" style={{ marginTop: 16 }}>
                            <Timeline
                                items={[
                                    { color: 'green', children: `${dayjs(currentDispute.createdAt).format('MM-DD HH:mm')} 纠纷创建` },
                                    ...(currentDispute.assignedAt ? [{ color: 'blue', children: `${dayjs(currentDispute.assignedAt).format('MM-DD HH:mm')} 分配给 ${currentDispute.assignedTo?.name || '客服'}` }] : []),
                                    ...(currentDispute.resolvedAt ? [{ color: 'green', children: `${dayjs(currentDispute.resolvedAt).format('MM-DD HH:mm')} 处理完成` }] : []),
                                ]}
                            />
                        </Card>
                    </>
                )}
            </Drawer>

            {/* 分配弹窗 */}
            <Modal title="分配纠纷" open={assignVisible} onOk={handleAssign} onCancel={() => setAssignVisible(false)} confirmLoading={submitting}>
                <Form form={assignForm} layout="vertical">
                    <Form.Item name="assignedToUserId" label="分配给" rules={[{ required: true, message: '请选择处理人' }]}>
                        <Select placeholder="请选择客服人员">
                            {staffList.map(s => (
                                <Select.Option key={s.id} value={s.id}>{s.name}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 处理弹窗 */}
            <Modal title="处理纠纷" open={resolveVisible} onOk={handleResolve} onCancel={() => setResolveVisible(false)} confirmLoading={submitting} width={500}>
                <Form form={resolveForm} layout="vertical">
                    <Form.Item name="resolution" label="解决方案" rules={[{ required: true, message: '请选择解决方案' }]}>
                        <Select placeholder="请选择解决方案">
                            <Select.Option value="refund">全额退款</Select.Option>
                            <Select.Option value="partial">部分退款</Select.Option>
                            <Select.Option value="reassign">重新分配陪玩师</Select.Option>
                            <Select.Option value="reject">驳回申诉</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="resolutionAmount" label="退款金额（元）">
                        <InputNumber min={0} precision={2} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item name="resolutionNotes" label="处理备注" rules={[{ required: true, message: '请输入处理备注' }]}>
                        <Input.TextArea rows={3} placeholder="请输入处理备注" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default DisputePage;
