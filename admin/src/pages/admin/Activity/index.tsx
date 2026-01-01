/**
 * Activity Management Page
 * Manage marketing activities - create, edit, delete, publish/unpublish, rewards management
 */
import React, { useState, useEffect, useCallback } from 'react';
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
    Input,
    Select,
    Switch,
    Modal,
    Divider,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    GiftOutlined,
    ReloadOutlined,
    EyeOutlined,
    TrophyOutlined,
    CheckCircleOutlined,
    StopOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import {
    activityApi,
    type Activity,
    type ActivityStats,
    type AllActivityStats,
    type CreateActivityDto,
    type CreateRewardDto,
    type ActivityReward,
    getActivityTypeLabel,
    getActivityStatusLabel,
    getActivityStatusColor,
    canEditActivity,
    canDeleteActivity,
    calculateStockPercentage,
} from '@/api/activity';
import { ActivityForm, RewardForm } from './components';

const { Title, Text } = Typography;

const ActivityPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [activities, setActivities] = useState<Activity[]>([]);
    const [stats, setStats] = useState<AllActivityStats | null>(null);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // Form modal states
    const [formVisible, setFormVisible] = useState(false);
    const [formLoading, setFormLoading] = useState(false);
    const [editingActivity, setEditingActivity] = useState<Activity | null>(null);

    // Detail modal
    const [detailVisible, setDetailVisible] = useState(false);
    const [detailActivity, setDetailActivity] = useState<Activity | null>(null);
    const [detailStats, setDetailStats] = useState<ActivityStats | null>(null);

    // Rewards modal
    const [rewardsVisible, setRewardsVisible] = useState(false);
    const [rewardsActivity, setRewardsActivity] = useState<Activity | null>(null);
    const [rewards, setRewards] = useState<ActivityReward[]>([]);
    const [rewardFormVisible, setRewardFormVisible] = useState(false);
    const [editingReward, setEditingReward] = useState<ActivityReward | null>(null);
    const [rewardLoading, setRewardLoading] = useState(false);

    // Filters
    const [keyword, setKeyword] = useState('');
    const [typeFilter, setTypeFilter] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string>('');
    const [visibleFilter, setVisibleFilter] = useState<boolean | undefined>(undefined);

    const loadStats = useCallback(async () => {
        try {
            const res = await activityApi.getAllActivityStats();
            if (res.data?.success && res.data?.data) {
                setStats(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load activity stats:', err);
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
            if (statusFilter) params.status = statusFilter;
            if (visibleFilter !== undefined) params.isVisible = visibleFilter;

            const res = await activityApi.getActivities(params);
            if (res.data?.success && res.data?.data) {
                setActivities(res.data.data);
                setTotal(res.data.pagination?.total || res.data.data.length || 0);
            } else {
                message.error(res.data?.message || '加载失败');
            }
        } catch (err) {
            console.error('Failed to load activities:', err);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, keyword, typeFilter, statusFilter, visibleFilter]);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    const handleAdd = () => {
        setEditingActivity(null);
        setFormVisible(true);
    };

    const handleEdit = (record: Activity) => {
        if (!canEditActivity(record)) {
            message.warning('只有草稿或预热状态的活动可以编辑');
            return;
        }
        setEditingActivity(record);
        setFormVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            const res = await activityApi.deleteActivity(id);
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

    const handleSubmit = async (values: CreateActivityDto) => {
        setFormLoading(true);
        try {
            let res;
            if (editingActivity) {
                res = await activityApi.updateActivity(editingActivity.id, values);
            } else {
                res = await activityApi.createActivity(values);
            }

            if (res.data?.success) {
                message.success(editingActivity ? '更新成功' : '创建成功');
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

    const handlePublish = async (activity: Activity) => {
        try {
            const res = await activityApi.publishActivity(activity.id);
            if (res.data?.success) {
                message.success('发布成功');
                loadData();
            } else {
                message.error(res.data?.message || '发布失败');
            }
        } catch (err) {
            console.error('Publish failed:', err);
            message.error('发布失败');
        }
    };

    const handleUnpublish = async (activity: Activity) => {
        try {
            const res = await activityApi.unpublishActivity(activity.id);
            if (res.data?.success) {
                message.success('下架成功');
                loadData();
            } else {
                message.error(res.data?.message || '下架失败');
            }
        } catch (err) {
            console.error('Unpublish failed:', err);
            message.error('下架失败');
        }
    };

    const handleViewDetail = async (activity: Activity) => {
        setDetailActivity(activity);
        setDetailVisible(true);

        // Load activity stats
        try {
            const res = await activityApi.getActivityStats(activity.id);
            if (res.data?.success && res.data?.data) {
                setDetailStats(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load activity stats:', err);
        }
    };

    const handleManageRewards = async (activity: Activity) => {
        setRewardsActivity(activity);
        setRewardsVisible(true);

        // Load rewards
        try {
            const res = await activityApi.getActivityRewards(activity.id);
            if (res.data?.success && res.data?.data) {
                setRewards(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load rewards:', err);
            message.error('加载奖励失败');
        }
    };

    const handleAddReward = () => {
        setEditingReward(null);
        setRewardFormVisible(true);
    };

    const handleEditReward = (reward: ActivityReward) => {
        setEditingReward(reward);
        setRewardFormVisible(true);
    };

    const handleDeleteReward = async (rewardId: number) => {
        try {
            const res = await activityApi.deleteReward(rewardId);
            if (res.data?.success) {
                message.success('删除奖励成功');
                // Reload rewards
                if (rewardsActivity) {
                    const res = await activityApi.getActivityRewards(rewardsActivity.id);
                    if (res.data?.success && res.data?.data) {
                        setRewards(res.data.data);
                    }
                }
            } else {
                message.error(res.data?.message || '删除失败');
            }
        } catch (err) {
            console.error('Delete reward failed:', err);
            message.error('删除失败');
        }
    };

    const handleSubmitReward = async (values: CreateRewardDto) => {
        setRewardLoading(true);
        try {
            let res;
            if (editingReward) {
                res = await activityApi.updateReward(editingReward.id, values);
            } else {
                res = await activityApi.createReward(values);
            }

            if (res.data?.success) {
                message.success(editingReward ? '更新成功' : '添加成功');
                setRewardFormVisible(false);
                // Reload rewards
                if (rewardsActivity) {
                    const res = await activityApi.getActivityRewards(rewardsActivity.id);
                    if (res.data?.success && res.data?.data) {
                        setRewards(res.data.data);
                    }
                }
            } else {
                message.error(res.data?.message || '操作失败');
            }
        } catch (err) {
            console.error('Submit reward failed:', err);
            message.error('操作失败');
        } finally {
            setRewardLoading(false);
        }
    };

    const columns: ColumnsType<Activity> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '活动名称',
            dataIndex: 'name',
            key: 'name',
            width: 180,
            render: (name, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{name}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        {getActivityTypeLabel(record.type)}
                    </Text>
                </Space>
            ),
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 100,
            render: (type) => <Tag color="blue">{getActivityTypeLabel(type)}</Tag>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 90,
            render: (status) => (
                <Tag color={getActivityStatusColor(status)}>
                    {getActivityStatusLabel(status)}
                </Tag>
            ),
        },
        {
            title: '时间',
            key: 'time',
            width: 180,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text style={{ fontSize: 12 }}>
                        开始：{dayjs(record.startAt).format('MM-DD HH:mm')}
                    </Text>
                    <Text style={{ fontSize: 12 }}>
                        结束：{dayjs(record.endAt).format('MM-DD HH:mm')}
                    </Text>
                </Space>
            ),
        },
        {
            title: '参与人数',
            key: 'participants',
            width: 100,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text>{record.totalParticipants.toLocaleString()}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        今日：{record.todayParticipants}
                    </Text>
                </Space>
            ),
        },
        {
            title: '已领取',
            dataIndex: 'totalClaimed',
            key: 'totalClaimed',
            width: 80,
            render: (count) => <Text strong>{count.toLocaleString()}</Text>,
        },
        {
            title: '显示',
            dataIndex: 'isVisible',
            key: 'isVisible',
            width: 70,
            render: (isVisible) => (
                <Tag color={isVisible ? 'green' : 'default'}>
                    {isVisible ? '显示' : '隐藏'}
                </Tag>
            ),
        },
        {
            title: 'VIP叠加',
            dataIndex: 'allowVipStack',
            key: 'allowVipStack',
            width: 80,
            render: (allow) => (
                <Tag color={allow ? 'purple' : 'default'}>
                    {allow ? '允许' : '不允许'}
                </Tag>
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
            width: 280,
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
                        icon={<TrophyOutlined />}
                        onClick={() => handleManageRewards(record)}
                    >
                        奖励
                    </Button>
                    {record.isVisible ? (
                        <Button
                            type="link"
                            size="small"
                            icon={<StopOutlined />}
                            onClick={() => handleUnpublish(record)}
                        >
                            下架
                        </Button>
                    ) : (
                        <Button
                            type="link"
                            size="small"
                            icon={<CheckCircleOutlined />}
                            onClick={() => handlePublish(record)}
                        >
                            发布
                        </Button>
                    )}
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                        disabled={!canEditActivity(record)}
                    >
                        编辑
                    </Button>
                    {canDeleteActivity(record) && (
                        <Popconfirm
                            title="确定删除此活动？"
                            onConfirm={() => handleDelete(record.id)}
                            okText="确定"
                            cancelText="取消"
                        >
                            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                                删除
                            </Button>
                        </Popconfirm>
                    )}
                </Space>
            ),
        },
    ];

    const rewardsColumns: ColumnsType<ActivityReward> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '优惠券模板',
            key: 'coupon',
            width: 200,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text>{record.couponTemplate?.name || `ID: ${record.couponTemplateId}`}</Text>
                    {record.couponTemplate?.type && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            类型：{record.couponTemplate.type}
                        </Text>
                    )}
                </Space>
            ),
        },
        {
            title: '发放数量',
            dataIndex: 'couponCount',
            key: 'couponCount',
            width: 90,
        },
        {
            title: '中奖概率',
            dataIndex: 'probability',
            key: 'probability',
            width: 90,
            render: (probability) => `${probability}%`,
        },
        {
            title: '库存',
            key: 'stock',
            width: 150,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text>
                        {record.remainingStock} / {record.totalStock || '无限制'}
                    </Text>
                    {record.totalStock > 0 && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            剩余：{calculateStockPercentage(record)}%
                        </Text>
                    )}
                </Space>
            ),
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 70,
        },
        {
            title: '操作',
            key: 'action',
            width: 140,
            render: (_, record) => (
                <Space size="small">
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEditReward(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定删除此奖励？"
                        onConfirm={() => handleDeleteReward(record.id)}
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
                <GiftOutlined /> 活动管理
            </Title>

            {/* Statistics */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card>
                        <Statistic
                            title="活动总数"
                            value={stats?.totalActivities || 0}
                            prefix={<GiftOutlined />}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card>
                        <Statistic
                            title="进行中"
                            value={stats?.activeActivities || 0}
                            valueStyle={{ color: '#3f8600' }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card>
                        <Statistic
                            title="草稿"
                            value={stats?.draftActivities || 0}
                            valueStyle={{ color: '#faad14' }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card>
                        <Statistic
                            title="总参与"
                            value={stats?.totalParticipants || 0}
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
                            <Select
                                placeholder="活动类型"
                                style={{ width: 120 }}
                                value={typeFilter || undefined}
                                onChange={(value) => {
                                    setTypeFilter(value || '');
                                    setCurrent(1);
                                }}
                                allowClear
                            >
                                <Select.Option value="coupon">优惠券发放</Select.Option>
                                <Select.Option value="discount">限时折扣</Select.Option>
                                <Select.Option value="gift">赠品活动</Select.Option>
                            </Select>
                            <Select
                                placeholder="活动状态"
                                style={{ width: 120 }}
                                value={statusFilter || undefined}
                                onChange={(value) => {
                                    setStatusFilter(value || '');
                                    setCurrent(1);
                                }}
                                allowClear
                            >
                                <Select.Option value="draft">草稿</Select.Option>
                                <Select.Option value="preheat">预热中</Select.Option>
                                <Select.Option value="active">进行中</Select.Option>
                                <Select.Option value="paused">已暂停</Select.Option>
                                <Select.Option value="ended">已结束</Select.Option>
                                <Select.Option value="canceled">已取消</Select.Option>
                            </Select>
                            <Switch
                                checked={visibleFilter === true}
                                onChange={(checked) => {
                                    setVisibleFilter(checked ? true : undefined);
                                    setCurrent(1);
                                }}
                                checkedChildren="显示"
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
                                新建活动
                            </Button>
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* Table */}
            <Card>
                <Table
                    columns={columns}
                    dataSource={activities}
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

            {/* Activity Form Modal */}
            <ActivityForm
                visible={formVisible}
                editing={!!editingActivity}
                initialValues={editingActivity}
                loading={formLoading}
                onCancel={() => setFormVisible(false)}
                onSubmit={handleSubmit}
            />

            {/* Detail Modal */}
            <Modal
                title="活动详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={800}
            >
                {detailActivity && (
                    <Card type="inner">
                        <Row gutter={[16, 16]}>
                            <Col span={12}>
                                <Text strong>活动名称：</Text>
                                <Text>{detailActivity.name}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>活动类型：</Text>
                                <Tag color="blue">{getActivityTypeLabel(detailActivity.type)}</Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>活动状态：</Text>
                                <Tag color={getActivityStatusColor(detailActivity.status)}>
                                    {getActivityStatusLabel(detailActivity.status)}
                                </Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>是否显示：</Text>
                                <Tag color={detailActivity.isVisible ? 'green' : 'default'}>
                                    {detailActivity.isVisible ? '显示' : '隐藏'}
                                </Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>预热时间：</Text>
                                <Text>{detailActivity.preheatAt ? dayjs(detailActivity.preheatAt).format('YYYY-MM-DD HH:mm:ss') : '-'}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>开始时间：</Text>
                                <Text>{dayjs(detailActivity.startAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>结束时间：</Text>
                                <Text>{dayjs(detailActivity.endAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>VIP叠加：</Text>
                                <Tag color={detailActivity.allowVipStack ? 'purple' : 'default'}>
                                    {detailActivity.allowVipStack ? '允许' : '不允许'}
                                </Tag>
                            </Col>
                            <Col span={12}>
                                <Text strong>总参与人数：</Text>
                                <Text>{detailActivity.totalParticipants}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>今日参与：</Text>
                                <Text>{detailActivity.todayParticipants}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>已领取奖励：</Text>
                                <Text>{detailActivity.totalClaimed}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>总限制：</Text>
                                <Text>{detailActivity.totalLimit === 0 ? '无限制' : detailActivity.totalLimit}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>每日限制：</Text>
                                <Text>{detailActivity.dailyLimit === 0 ? '无限制' : detailActivity.dailyLimit}</Text>
                            </Col>
                            <Col span={12}>
                                <Text strong>每人限制：</Text>
                                <Text>{detailActivity.perUserLimit === 0 ? '无限制' : detailActivity.perUserLimit}</Text>
                            </Col>
                            <Col span={24}>
                                <Text strong>活动描述：</Text>
                                <div>{detailActivity.description || '-'}</div>
                            </Col>
                            <Col span={24}>
                                <Text strong>活动规则：</Text>
                                <div style={{ whiteSpace: 'pre-wrap' }}>{detailActivity.rules || '-'}</div>
                            </Col>
                            {detailStats && (
                                <>
                                    <Col span={12}>
                                        <Text strong>统计参与人数：</Text>
                                        <Text>{detailStats.totalParticipants}</Text>
                                    </Col>
                                    <Col span={12}>
                                        <Text strong>统计今日参与：</Text>
                                        <Text>{detailStats.todayParticipants}</Text>
                                    </Col>
                                    <Col span={12}>
                                        <Text strong>统计已领取：</Text>
                                        <Text>{detailStats.totalClaimed}</Text>
                                    </Col>
                                </>
                            )}
                        </Row>
                    </Card>
                )}
            </Modal>

            {/* Rewards Management Modal */}
            <Modal
                title={`奖励管理 - ${rewardsActivity?.name || ''}`}
                open={rewardsVisible}
                onCancel={() => setRewardsVisible(false)}
                footer={null}
                width={900}
            >
                <div style={{ marginBottom: 16 }}>
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleAddReward}>
                        添加奖励
                    </Button>
                </div>
                <Table
                    columns={rewardsColumns}
                    dataSource={rewards}
                    rowKey="id"
                    pagination={false}
                    size="small"
                />
            </Modal>

            {/* Reward Form Modal */}
            {rewardsActivity && (
                <RewardForm
                    visible={rewardFormVisible}
                    editing={!!editingReward}
                    activityId={rewardsActivity.id}
                    initialValues={editingReward}
                    loading={rewardLoading}
                    onCancel={() => setRewardFormVisible(false)}
                    onSubmit={handleSubmitReward}
                />
            )}
        </div>
    );
};

export default ActivityPage;
