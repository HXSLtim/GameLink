/**
 * 实名审核页面
 * 审核陪玩师的实名认证申请
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Popconfirm,
    Drawer,
    Descriptions,
    Avatar,
    Image,
    Form,
    Input,
    Card,
    Statistic,
    Row,
    Col,
    Tabs,
    Typography,
    Divider,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    CheckOutlined,
    CloseOutlined,
    UserOutlined,
    DeleteOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    IdcardOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { adminApi, type PlayerCertification } from '@/api/admin';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

const { Text } = Typography;

const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'gold', text: '待审核' },
    verified: { color: 'success', text: '已通过' },
    rejected: { color: 'error', text: '已拒绝' },
};

const PlayerCertificationPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [certifications, setCertifications] = useState<PlayerCertification[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});
    const [activeTab, setActiveTab] = useState<string>('all');

    // 统计数据
    const [stats, setStats] = useState<Record<string, number>>({});

    // 弹窗状态
    const [detailVisible, setDetailVisible] = useState(false);
    const [verifyVisible, setVerifyVisible] = useState(false);
    const [currentCert, setCurrentCert] = useState<PlayerCertification | null>(null);
    const [verifyForm] = Form.useForm();
    const [verifyLoading, setVerifyLoading] = useState(false);

    // 加载统计数据
    const loadStats = useCallback(async () => {
        try {
            const response = await adminApi.getPlayerCertificationStats();
            if (response.data.success) {
                setStats(response.data.data || {});
            }
        } catch (error) {
            logger.error('Load stats error:', error);
        }
    }, []);

    // 加载数据
    const loadData = useCallback(async (params?: Record<string, unknown>) => {
        setLoading(true);
        try {
            const queryParams: Record<string, unknown> = {
                page: current,
                pageSize,
                ...searchParams,
                ...params,
            };

            if (activeTab !== 'all') {
                queryParams.status = activeTab;
            }

            const response = await adminApi.getPlayerCertifications(queryParams);
            if (response.data.success) {
                setCertifications(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load certifications error:', error);
            message.error('加载实名认证列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, activeTab]);

    useEffect(() => {
        loadStats();
    }, [loadStats]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleViewDetail = (cert: PlayerCertification) => {
        setCurrentCert(cert);
        setDetailVisible(true);
    };

    const handleOpenVerify = (cert: PlayerCertification) => {
        setCurrentCert(cert);
        verifyForm.resetFields();
        setVerifyVisible(true);
    };

    const handleVerify = async (approved: boolean) => {
        if (!currentCert) return;
        try {
            setVerifyLoading(true);
            const values = await verifyForm.validateFields();
            await adminApi.verifyPlayerCertification(currentCert.id, {
                status: approved ? 'verified' : 'rejected',
                rejectReason: values.rejectReason,
            });
            message.success(approved ? '审核通过' : '审核拒绝');
            setVerifyVisible(false);
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Verify error:', error);
            message.error('审核操作失败');
        } finally {
            setVerifyLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await adminApi.deletePlayerCertification(id);
            message.success('删除成功');
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Delete error:', error);
            message.error('删除失败');
        }
    };

    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '真实姓名' },
        { name: 'playerId', label: '陪玩师ID', type: 'input', placeholder: '陪玩师ID' },
    ];

    const columns: ColumnsType<PlayerCertification> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '陪玩师',
            key: 'player',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Avatar size={40} src={record.player?.user?.avatarUrl} icon={<UserOutlined />} />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.player?.nickname || '-'}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>ID: {record.playerId}</Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '真实姓名',
            dataIndex: 'realName',
            key: 'realName',
            width: 120,
            render: name => name ? (
                <Space>
                    <IdcardOutlined />
                    <span>{name}</span>
                </Space>
            ) : '-',
        },
        {
            title: '身份证照片',
            key: 'idCard',
            width: 140,
            render: (_, record) => (
                <Image.PreviewGroup>
                    <Space>
                        {record.idCardFrontUrl && (
                            <Image width={50} height={35} src={record.idCardFrontUrl} style={{ objectFit: 'cover', borderRadius: 4 }} />
                        )}
                        {record.idCardBackUrl && (
                            <Image width={50} height={35} src={record.idCardBackUrl} style={{ objectFit: 'cover', borderRadius: 4 }} />
                        )}
                        {!record.idCardFrontUrl && !record.idCardBackUrl && <Text type="secondary">-</Text>}
                    </Space>
                </Image.PreviewGroup>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: status => (
                <Tag color={statusMap[status]?.color || 'default'}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '申请时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 160,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-',
        },
        {
            title: '审核时间',
            dataIndex: 'verifiedAt',
            key: 'verifiedAt',
            width: 160,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 240, // 3个按钮 × 80px
            fixed: 'right',
            render: (_, record) => (
                <Space size={4}>
                    <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record)}>
                        详情
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<CheckOutlined />}
                        onClick={() => handleOpenVerify(record)}
                        disabled={record.status !== 'pending'}
                    >
                        审核
                    </Button>
                    <Popconfirm title="确定删除该记录？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const tabItems = [
        { key: 'all', label: '全部' },
        { key: 'pending', label: `待审核 (${stats.pending || 0})` },
        { key: 'verified', label: `已通过 (${stats.verified || 0})` },
        { key: 'rejected', label: `已拒绝 (${stats.rejected || 0})` },
    ];

    return (
        <PageContainer title="实名审核" subTitle="审核陪玩师的实名认证申请">
            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card size="small">
                        <Statistic
                            title="待审核"
                            value={stats.pending || 0}
                            prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card size="small">
                        <Statistic
                            title="已通过"
                            value={stats.verified || 0}
                            prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card size="small">
                        <Statistic
                            title="已拒绝"
                            value={stats.rejected || 0}
                            prefix={<CloseCircleOutlined style={{ color: '#ff4d4f' }} />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card size="small">
                        <Statistic
                            title="总计"
                            value={(stats.pending || 0) + (stats.verified || 0) + (stats.rejected || 0)}
                        />
                    </Card>
                </Col>
            </Row>

            <Card>
                <Tabs
                    activeKey={activeTab}
                    onChange={key => { setActiveTab(key); setCurrent(1); }}
                    items={tabItems}
                    style={{ marginBottom: 16 }}
                />
                <SearchTable
                    columns={columns}
                    dataSource={certifications}
                    rowKey="id"
                    searchFields={searchFields}
                    onSearch={handleSearch}
                    onRefresh={() => loadData()}
                    loading={loading}
                    showCreate={false}
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
                    scroll={{ x: 1100 }}
                />
            </Card>

            {/* 详情抽屉 */}
            <Drawer
                title="实名认证详情"
                open={detailVisible}
                onClose={() => setDetailVisible(false)}
                size="large"
            >
                {currentCert && (
                    <>
                        {/* 陪玩师信息卡片 */}
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space align="start">
                                <Avatar size={64} src={currentCert.player?.user?.avatarUrl} icon={<UserOutlined />} />
                                <div>
                                    <div style={{ fontSize: 16, fontWeight: 500 }}>
                                        {currentCert.player?.nickname || '-'}
                                    </div>
                                    <Text type="secondary">ID: {currentCert.playerId}</Text>
                                    <div style={{ marginTop: 8 }}>
                                        <Tag icon={<IdcardOutlined />} color="blue">
                                            {currentCert.realName || '未填写'}
                                        </Tag>
                                    </div>
                                </div>
                            </Space>
                        </Card>

                        <Descriptions column={1} bordered size="small">
                            <Descriptions.Item label="真实姓名">{currentCert.realName || '-'}</Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={statusMap[currentCert.status]?.color}>
                                    {statusMap[currentCert.status]?.text}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="申请时间">
                                {currentCert.createdAt ? dayjs(currentCert.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                            </Descriptions.Item>
                            {currentCert.verifiedAt && (
                                <Descriptions.Item label="审核时间">
                                    {dayjs(currentCert.verifiedAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                            )}
                            {currentCert.rejectReason && (
                                <Descriptions.Item label="拒绝原因">
                                    <Text type="danger">{currentCert.rejectReason}</Text>
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider style={{ marginTop: 24 }}>身份证照片</Divider>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Card size="small" title="正面" type="inner">
                                    {currentCert.idCardFrontUrl ? (
                                        <Image width="100%" src={currentCert.idCardFrontUrl} style={{ borderRadius: 4 }} />
                                    ) : (
                                        <div style={{ color: '#999', padding: 40, textAlign: 'center', background: '#f5f5f5', borderRadius: 4 }}>
                                            暂无照片
                                        </div>
                                    )}
                                </Card>
                            </Col>
                            <Col span={12}>
                                <Card size="small" title="背面" type="inner">
                                    {currentCert.idCardBackUrl ? (
                                        <Image width="100%" src={currentCert.idCardBackUrl} style={{ borderRadius: 4 }} />
                                    ) : (
                                        <div style={{ color: '#999', padding: 40, textAlign: 'center', background: '#f5f5f5', borderRadius: 4 }}>
                                            暂无照片
                                        </div>
                                    )}
                                </Card>
                            </Col>
                        </Row>

                        {currentCert.photoUrl && (
                            <>
                                <Divider style={{ marginTop: 24 }}>个人照片</Divider>
                                <Image width={200} src={currentCert.photoUrl} style={{ borderRadius: 4 }} />
                            </>
                        )}

                        {currentCert.status === 'pending' && (
                            <div style={{ marginTop: 24, textAlign: 'center' }}>
                                <Button
                                    type="primary"
                                    icon={<CheckOutlined />}
                                    onClick={() => {
                                        setDetailVisible(false);
                                        handleOpenVerify(currentCert);
                                    }}
                                >
                                    审核
                                </Button>
                            </div>
                        )}
                    </>
                )}
            </Drawer>

            {/* 审核弹窗 */}
            <Modal
                title="审核实名认证"
                open={verifyVisible}
                onCancel={() => setVerifyVisible(false)}
                width={700}
                footer={
                    <Space>
                        <Button onClick={() => setVerifyVisible(false)}>取消</Button>
                        <Button danger icon={<CloseOutlined />} loading={verifyLoading} onClick={() => handleVerify(false)}>
                            拒绝
                        </Button>
                        <Button type="primary" icon={<CheckOutlined />} loading={verifyLoading} onClick={() => handleVerify(true)}>
                            通过
                        </Button>
                    </Space>
                }
            >
                {currentCert && (
                    <>
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space align="start">
                                <Avatar size={48} src={currentCert.player?.user?.avatarUrl} icon={<UserOutlined />} />
                                <div>
                                    <div style={{ fontWeight: 500 }}>{currentCert.player?.nickname || '-'}</div>
                                    <Tag icon={<IdcardOutlined />} style={{ marginTop: 4 }}>
                                        {currentCert.realName || '未填写'}
                                    </Tag>
                                </div>
                            </Space>
                        </Card>

                        <Row gutter={16} style={{ marginBottom: 16 }}>
                            <Col span={12}>
                                <Text strong style={{ display: 'block', marginBottom: 8 }}>身份证正面</Text>
                                {currentCert.idCardFrontUrl ? (
                                    <Image width="100%" src={currentCert.idCardFrontUrl} style={{ borderRadius: 4 }} />
                                ) : (
                                    <div style={{ color: '#999', padding: 30, textAlign: 'center', background: '#f5f5f5', borderRadius: 4 }}>
                                        暂无照片
                                    </div>
                                )}
                            </Col>
                            <Col span={12}>
                                <Text strong style={{ display: 'block', marginBottom: 8 }}>身份证背面</Text>
                                {currentCert.idCardBackUrl ? (
                                    <Image width="100%" src={currentCert.idCardBackUrl} style={{ borderRadius: 4 }} />
                                ) : (
                                    <div style={{ color: '#999', padding: 30, textAlign: 'center', background: '#f5f5f5', borderRadius: 4 }}>
                                        暂无照片
                                    </div>
                                )}
                            </Col>
                        </Row>

                        <Form form={verifyForm} layout="vertical">
                            <Form.Item name="rejectReason" label="拒绝原因（拒绝时必填）">
                                <Input.TextArea rows={3} placeholder="请输入拒绝原因，如：照片模糊、信息不匹配等" />
                            </Form.Item>
                        </Form>
                    </>
                )}
            </Modal>
        </PageContainer>
    );
};

export default PlayerCertificationPage;
