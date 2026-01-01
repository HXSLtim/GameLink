/**
 * Activity Rewards Management Page
 * View and manage all activity rewards across all activities
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    message,
    Typography,
    Row,
    Col,
    Statistic,
    Select,
    Progress,
    theme,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    TrophyOutlined,
    ReloadOutlined,
    GiftOutlined,
} from '@ant-design/icons';
import {
    activityApi,
    type Activity,
    type ActivityReward,
    calculateStockPercentage,
} from '@/api/activity';

const { Title, Text } = Typography;

const RewardsPage: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [activities, setActivities] = useState<Activity[]>([]);
    const [selectedActivityId, setSelectedActivityId] = useState<number | undefined>(undefined);
    const [rewards, setRewards] = useState<ActivityReward[]>([]);

    const loadActivities = useCallback(async () => {
        try {
            const res = await activityApi.getActivities({ page_size: 1000 });
            if (res.data?.success && res.data?.data) {
                setActivities(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load activities:', err);
        }
    }, []);

    const loadRewards = useCallback(async (activityId?: number) => {
        if (!activityId) {
            setRewards([]);
            return;
        }

        setLoading(true);
        try {
            const res = await activityApi.getActivityRewards(activityId);
            if (res.data?.success && res.data?.data) {
                setRewards(res.data.data);
            } else {
                message.error(res.data?.message || '加载失败');
            }
        } catch (err) {
            console.error('Failed to load rewards:', err);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadActivities();
    }, [loadActivities]);

    useEffect(() => {
        loadRewards(selectedActivityId);
    }, [selectedActivityId, loadRewards]);

    // Calculate stats
    const totalRewards = rewards.length;
    const totalStock = rewards.reduce((sum, r) => sum + (r.totalStock || 0), 0);
    const totalRemaining = rewards.reduce((sum, r) => sum + r.remainingStock, 0);
    const totalProbability = rewards.reduce((sum, r) => sum + r.probability, 0);

    const columns: ColumnsType<ActivityReward> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 70,
        },
        {
            title: '优惠券模板',
            key: 'coupon',
            width: 220,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{record.couponTemplate?.name || `ID: ${record.couponTemplateId}`}</Text>
                    {record.couponTemplate?.type && (
                        <Tag color="blue" style={{ fontSize: 12 }}>
                            {record.couponTemplate.type}
                        </Tag>
                    )}
                </Space>
            ),
        },
        {
            title: '发放数量',
            dataIndex: 'couponCount',
            key: 'couponCount',
            width: 100,
            render: (count) => (
                <Space>
                    <GiftOutlined />
                    <Text strong>{count}</Text>
                    <Text type="secondary">张/人</Text>
                </Space>
            ),
        },
        {
            title: '中奖概率',
            dataIndex: 'probability',
            key: 'probability',
            width: 130,
            render: (probability) => (
                <Space direction="vertical" size={0}>
                    <Text strong style={{ color: probability > 20 ? '#ff4d4f' : undefined }}>
                        {probability}%
                    </Text>
                    <Progress
                        percent={probability}
                        size="small"
                        status={probability > 50 ? 'exception' : 'normal'}
                        showInfo={false}
                    />
                </Space>
            ),
        },
        {
            title: '库存状态',
            key: 'stock',
            width: 200,
            render: (_, record) => {
                const percentage = record.totalStock > 0
                    ? calculateStockPercentage(record)
                    : 100;
                const isLowStock = percentage < 20;
                const isUnlimited = record.totalStock === 0;

                return (
                    <Space direction="vertical" size={0} style={{ width: '100%' }}>
                        <Text>
                            {record.remainingStock.toLocaleString()} / {' '}
                            {isUnlimited ? '无限制' : record.totalStock.toLocaleString()}
                        </Text>
                        {!isUnlimited && (
                            <Progress
                                percent={percentage}
                                size="small"
                                status={isLowStock ? 'exception' : 'normal'}
                                showInfo={false}
                            />
                        )}
                        {isLowStock && !isUnlimited && (
                            <Tag color="red" style={{ fontSize: 12 }}>库存紧张</Tag>
                        )}
                    </Space>
                );
            },
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 80,
            sorter: (a, b) => a.sortOrder - b.sortOrder,
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 160,
            render: (date) => new Date(date).toLocaleString('zh-CN'),
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 160,
            render: (date) => new Date(date).toLocaleString('zh-CN'),
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}>
                <TrophyOutlined /> 活动奖励管理
            </Title>

            {/* Activity Selector */}
            <Card style={{ marginBottom: 16 }}>
                <Row gutter={16} align="middle">
                    <Col>
                        <Space>
                            <Text strong>选择活动：</Text>
                            <Select
                                placeholder="请选择活动查看奖励"
                                style={{ width: 300 }}
                                value={selectedActivityId}
                                onChange={setSelectedActivityId}
                                options={activities.map(a => ({
                                    label: `${a.name} (${getActivityStatusLabel(a.status)})`,
                                    value: a.id,
                                }))}
                                showSearch
                                optionFilterProp="label"
                            />
                            <Button icon={<ReloadOutlined />} onClick={() => loadRewards(selectedActivityId)}>
                                刷新
                            </Button>
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* Statistics */}
            {selectedActivityId && (
                <Row gutter={16} style={{ marginBottom: 16 }}>
                    <Col xs={12} sm={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="奖励数量"
                                value={totalRewards}
                                prefix={<TrophyOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col xs={12} sm={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="总库存"
                                value={totalStock}
                                valueStyle={{ color: token.colorPrimary }}
                            />
                        </Card>
                    </Col>
                    <Col xs={12} sm={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="剩余库存"
                                value={totalRemaining}
                                valueStyle={{ color: totalRemaining < totalStock * 0.2 ? token.colorError : token.colorSuccess }}
                            />
                        </Card>
                    </Col>
                    <Col xs={12} sm={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="总概率"
                                value={totalProbability}
                                suffix="%"
                                valueStyle={{ color: totalProbability > 100 ? token.colorError : token.colorPrimary }}
                            />
                        </Card>
                    </Col>
                </Row>
            )}

            {/* Table */}
            <Card>
                <Table
                    columns={columns}
                    dataSource={rewards}
                    rowKey="id"
                    loading={loading}
                    pagination={{
                        pageSize: 20,
                        showSizeChanger: true,
                        showTotal: (t) => `共 ${t} 条`,
                    }}
                />
            </Card>
        </div>
    );
};

// Helper function for status label
function getActivityStatusLabel(status: string): string {
    const labels: Record<string, string> = {
        draft: '草稿',
        preheat: '预热中',
        active: '进行中',
        paused: '已暂停',
        ended: '已结束',
        canceled: '已取消',
    };
    return labels[status] || status;
}

export default RewardsPage;
