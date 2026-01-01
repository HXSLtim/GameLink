/**
 * Coupon Template Management Page
 * Manage coupon templates - create, edit, delete, enable/disable, issue to users
 */
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    message,
    Popconfirm,
    Typography,
    Row,
    Col,
    Statistic,
    Switch,
    Modal,
    Input,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    GiftOutlined,
    ReloadOutlined,
    SendOutlined,
    EyeOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import {
    couponApi,
    type CouponTemplate,
    type CouponStats,
    type CreateTemplateDto,
    getCouponTypeLabel,
    getCouponScopeLabel,
    getCouponSourceLabel,
    centsToYuan,
    parseJsonArray,
} from '@/api/coupon';
import { TemplateForm, IssueModal } from './components';

const { Title, Text } = Typography;

const CouponPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [templates, setTemplates] = useState<CouponTemplate[]>([]);
    const [stats, setStats] = useState<CouponStats | null>(null);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // Modal states
    const [formVisible, setFormVisible] = useState(false);
    const [formLoading, setFormLoading] = useState(false);
    const [editingTemplate, setEditingTemplate] = useState<CouponTemplate | null>(null);

    // Issue modal states
    const [issueVisible, setIssueVisible] = useState(false);
    const [issueTemplateId, setIssueTemplateId] = useState<number | undefined>(undefined);

    // Detail modal
    const [detailVisible, setDetailVisible] = useState(false);
    const [detailTemplate, setDetailTemplate] = useState<CouponTemplate | null>(null);

    // Filters
    const [keyword, setKeyword] = useState('');
    const [typeFilter, setTypeFilter] = useState<string>('');
    const [sourceFilter, setSourceFilter] = useState<string>('');
    const [activeFilter, setActiveFilter] = useState<boolean | undefined>(undefined);

    const loadStats = useCallback(async () => {
        try {
            const res = await couponApi.getCouponStats();
            if (res.data?.success && res.data?.data) {
                setStats(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load coupon stats:', err);
        }
    }, []);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = {
                page: current,
                page_size: pageSize,
            };
            if (keyword) params.keyword = keyword;
            if (typeFilter) params.type = typeFilter;
            if (sourceFilter) params.source = sourceFilter;
            if (activeFilter !== undefined) params.isActive = activeFilter;

            const res = await couponApi.getTemplates(params);
            if (res.data?.success && res.data?.data) {
                setTemplates(res.data.data);
                setTotal(res.data.pagination?.total || res.data.data.length || 0);
            } else {
                message.error(res.data?.message || '加载失败');
            }
        } catch (err) {
            console.error('Failed to load templates:', err);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, keyword, typeFilter, sourceFilter, activeFilter]);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    const handleAdd = () => {
        setEditingTemplate(null);
        setFormVisible(true);
    };

    const handleEdit = (record: CouponTemplate) => {
        setEditingTemplate(record);
        setFormVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            const res = await couponApi.deleteTemplate(id);
            if (res.data?.success) {
                message.success('删除成功');
                loadData();
                loadStats();
            } else {
                message.error(res.data?.message || '删除失败');
            }
        } catch (err) {
            console.error('Delete failed:', err);
            message.error('删除失败');
        }
    };

    const handleSubmit = async (values: CreateTemplateDto) => {
        setFormLoading(true);
        try {
            let res;
            if (editingTemplate) {
                res = await couponApi.updateTemplate(editingTemplate.id, values);
            } else {
                res = await couponApi.createTemplate(values);
            }

            if (res.data?.success) {
                message.success(editingTemplate ? '更新成功' : '创建成功');
                setFormVisible(false);
                loadData();
                loadStats();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            console.error('Submit failed:', err);
            message.error('操作失败');
        } finally {
            setFormLoading(false);
        }
    };

    const handleToggleActive = async (id: number, isActive: boolean) => {
        try {
            const res = await couponApi.toggleCouponTemplate(id, isActive);
            if (res.data?.success) {
                message.success(isActive ? '已启用' : '已禁用');
                loadData();
                loadStats();
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            console.error('Toggle active failed:', err);
            message.error('操作失败');
        }
    };

    const handleIssue = (templateId?: number) => {
        setIssueTemplateId(templateId);
        setIssueVisible(true);
    };

    const handleIssueSuccess = () => {
        setIssueVisible(false);
        loadData();
        loadStats();
    };

    const handleViewDetail = (template: CouponTemplate) => {
        setDetailTemplate(template);
        setDetailVisible(true);
    };

    const columns: ColumnsType<CouponTemplate> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '优惠券名称',
            dataIndex: 'name',
            key: 'name',
            width: 150,
            render: (name, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{name}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        {getCouponTypeLabel(record.type)}
                    </Text>
                </Space>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 80,
            render: (type) => <Tag color="blue">{getCouponTypeLabel(type)}</Tag>,
        },
        {
            title: '来源',
            dataIndex: 'source',
            key: 'source',
            width: 90,
            render: (source) => <Tag>{getCouponSourceLabel(source)}</Tag>,
        },
        {
            title: '优惠内容',
            key: 'discount',
            width: 150,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    {record.type === 'deduct' ? (
                        <Text type="danger" strong>
                            减 ¥{centsToYuan(record.deductAmountCents)}
                        </Text>
                    ) : (
                        <Text type="danger" strong>
                            {(record.discountRate * 10).toFixed(1)} 折
                        </Text>
                    )}
                    {record.minAmountCents > 0 && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            满 ¥{centsToYuan(record.minAmountCents)} 可用
                        </Text>
                    )}
                </Space>
            ),
        },
        {
            title: '适用范围',
            dataIndex: 'scope',
            key: 'scope',
            width: 90,
            render: (scope) => <Tag color="purple">{getCouponScopeLabel(scope)}</Tag>,
        },
        {
            title: '有效期',
            key: 'validity',
            width: 120,
            render: (_, record) => (
                <Text style={{ fontSize: 12 }}>
                    {record.validityType === 'days'
                        ? `${record.validityDays} 天`
                        : record.fixedExpireAt
                          ? dayjs(record.fixedExpireAt).format('YYYY-MM-DD')
                          : '-'}
                </Text>
            ),
        },
        {
            title: '已领/总数',
            key: 'claimed',
            width: 90,
            render: (_, record) => (
                <Text>
                    {record.claimedCount} / {record.totalCount}
                </Text>
            ),
        },
        {
            title: '每人限领',
            dataIndex: 'perUserLimit',
            key: 'perUserLimit',
            width: 80,
            render: (limit) => (limit === 0 ? <Text type="secondary">不限</Text> : limit),
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: 90,
            render: (isActive, record) => (
                <Switch
                    checked={isActive}
                    onChange={(checked) => handleToggleActive(record.id, checked)}
                    checkedChildren="启用"
                    unCheckedChildren="禁用"
                />
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 150,
            render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm'),
        },
        {
            title: '操作',
            key: 'action',
            width: 220,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small" wrap>
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<SendOutlined />}
                        onClick={() => handleIssue(record.id)}
                    >
                        发放
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定删除此优惠券模板？"
                        onConfirm={() => handleDelete(record.id)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}>
                <GiftOutlined /> 优惠券模板管理
            </Title>

            {/* Statistics */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="模板总数"
                            value={stats?.totalTemplates || 0}
                            prefix={<GiftOutlined />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="启用模板"
                            value={stats?.activeTemplates || 0}
                            valueStyle={{ color: '#3f8600' }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已发放"
                            value={stats?.totalCoupons || 0}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已使用"
                            value={stats?.usedCoupons || 0}
                            valueStyle={{ color: '#1890ff' }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* Filters and Actions */}
            <Card style={{ marginBottom: 16 }}>
                <Row gutter={16} align="middle">
                    <Col flex="auto">
                        <Space wrap>
                            <Input
                                placeholder="搜索名称"
                                style={{ width: 150 }}
                                value={keyword}
                                onChange={(e) => setKeyword(e.target.value)}
                                onPressEnter={() => {
                                    setCurrent(1);
                                    loadData();
                                }}
                            />
                            <Input
                                placeholder="类型 (deduct/discount)"
                                style={{ width: 150 }}
                                value={typeFilter}
                                onChange={(e) => setTypeFilter(e.target.value)}
                                onPressEnter={() => {
                                    setCurrent(1);
                                    loadData();
                                }}
                            />
                            <Input
                                placeholder="来源"
                                style={{ width: 150 }}
                                value={sourceFilter}
                                onChange={(e) => setSourceFilter(e.target.value)}
                                onPressEnter={() => {
                                    setCurrent(1);
                                    loadData();
                                }}
                            />
                            <Switch
                                checked={activeFilter === true}
                                onChange={(checked) => {
                                    setActiveFilter(checked ? true : undefined);
                                    setCurrent(1);
                                }}
                                checkedChildren="启用"
                                unCheckedChildren="全部"
                            />
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <Button icon={<ReloadOutlined />} onClick={loadData}>
                                刷新
                            </Button>
                            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                                新建模板
                            </Button>
                            <Button icon={<SendOutlined />} onClick={() => handleIssue(undefined)}>
                                发放优惠券
                            </Button>
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* Table */}
            <Card>
                <Table
                    columns={columns}
                    dataSource={templates}
                    rowKey="id"
                    loading={loading}
                    scroll={{ x: 1400 }}
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
            </Card>

            {/* Template Form Modal */}
            <TemplateForm
                visible={formVisible}
                editing={!!editingTemplate}
                initialValues={editingTemplate || undefined}
                loading={formLoading}
                onCancel={() => setFormVisible(false)}
                onSubmit={handleSubmit}
            />

            {/* Issue Modal */}
            <IssueModal
                visible={issueVisible}
                templateId={issueTemplateId}
                onCancel={() => setIssueVisible(false)}
                onSuccess={handleIssueSuccess}
            />

            {/* Detail Modal */}
            <Modal
                title="优惠券模板详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={700}
            >
                {detailTemplate && (
                    <Card type="inner">
                        <Row gutter={[16, 16]}>
                            <Col span={12}>
                                <Text strong>名称：</Text>
                                <Text>{detailTemplate.name}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>类型：</Text>
                                <Tag color="blue">{getCouponTypeLabel(detailTemplate.type)}</Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>来源：</Text>
                                <Tag>{getCouponSourceLabel(detailTemplate.source)}</Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>适用范围：</Text>
                                <Tag color="purple">{getCouponScopeLabel(detailTemplate.scope)}</Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>最低消费：</Text>
                                <Text>¥{centsToYuan(detailTemplate.minAmountCents)}</Text>
                            </Col>
                            {detailTemplate.type === 'deduct' && (
                                <Col span={12}>
                                    <Text strong>减免金额：</Text>
                                    <Text type="danger" strong>
                                        ¥{centsToYuan(detailTemplate.deductAmountCents)}
                                    </Text>
                                </Col>
                            )}
                            {detailTemplate.type === 'discount' && (
                                <>
                                    <Col span={12}>
                                        <Text strong>折扣率：</Text>
                                        <Text type="danger" strong>
                                            {(detailTemplate.discountRate * 10).toFixed(1)} 折
                                        </Text>
                                    </Col>
                                    <Col span={12}>
                                        <Text strong>最大优惠：</Text>
                                        <Text>¥{centsToYuan(detailTemplate.maxDiscountCents)}</Text>
                                    </Col>
                                </>
                            )}
                            <Col span={12}>
                                <Text strong>有效期类型：</Text>
                                <Text>{detailTemplate.validityType === 'days' ? '相对天数' : '固定时间'}</Text>
                            </Col>
                            {detailTemplate.validityType === 'days' && (
                                <Col span={12}>
                                    <Text strong>有效天数：</Text>
                                    <Text>{detailTemplate.validityDays} 天</Text>
                                </Col>
                            )}
                            {detailTemplate.validityType === 'fixed' && (
                                <Col span={12}>
                                    <Text strong>过期时间：</Text>
                                    <Text>{detailTemplate.fixedExpireAt || '-'}</Text>
                                </Col>
                            )}
                            <Col span={12}>
                                <Text strong>发放总数：</Text>
                                <Text>{detailTemplate.totalCount}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>已领取：</Text>
                                <Text>{detailTemplate.claimedCount}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>每人限领：</Text>
                                <Text>{detailTemplate.perUserLimit === 0 ? '不限制' : detailTemplate.perUserLimit}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>领取链接：</Text>
                                <Text>{detailTemplate.claimLink || '-'}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>状态：</Text>
                                <Tag color={detailTemplate.isActive ? 'green' : 'red'}>
                                    {detailTemplate.isActive ? '启用' : '禁用'}
                                </Tag>
                            </Col>
                            <Col span={24}>
                                <Text strong>描述：</Text>
                                <Text>{detailTemplate.description || '-'}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>创建时间：</Text>
                                <Text>{dayjs(detailTemplate.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>更新时间：</Text>
                                <Text>{dayjs(detailTemplate.updatedAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
                            </Col>
                        </Row>
                    </Card>
                )}
            </Modal>
        </div>
    );
};

export default CouponPage;
