/**
 * User Coupon Management Page
 * View and manage user coupons - status, usage details, etc.
 */
import React, { useState, useEffect, useCallback, useMemo } from 'react';
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
    Input,
    Select,
    Modal,
    Descriptions,
    theme,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    GiftOutlined,
    ReloadOutlined,
    EyeOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import {
    couponApi,
    type CouponWithTemplate,
    type CouponState,
    getCouponTypeLabel,
    getCouponScopeLabel,
    getCouponSourceLabel,
    getCouponStateLabel,
    getCouponStateColor,
    centsToYuan,
} from '@/api/coupon';
import { adminApi } from '@/api/admin';

import { logger } from '@/utils/logger';
const { Title, Text } = Typography;

const UserCouponPage: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [coupons, setCoupons] = useState<CouponWithTemplate[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // Detail modal
    const [detailVisible, setDetailVisible] = useState(false);
    const [detailCoupon, setDetailCoupon] = useState<CouponWithTemplate | null>(null);

    // Filters
    const [userIdFilter, setUserIdFilter] = useState<number | undefined>(undefined);
    const [stateFilter, setStateFilter] = useState<string>('');
    const [typeFilter, setTypeFilter] = useState<string>('');
    const [userSearchKeyword, setUserSearchKeyword] = useState('');
    const [users, setUsers] = useState<Array<{ id: number; name: string; phone: string }>>([]);

    // Load users for filter
    useEffect(() => {
        const loadUsers = async () => {
            try {
                const res = await adminApi.getUsers({ page_size: 100 });
                if (res.data?.success && res.data?.data) {
                    setUsers(res.data.data);
                }
            } catch (err) {
                logger.error('Failed to load users:', err);
            }
        };
        loadUsers();
    }, []);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: Record<string, unknown> = {
                page: current,
                page_size: pageSize,
            };
            if (userIdFilter) params.userId = userIdFilter;
            if (stateFilter) params.state = stateFilter;
            if (typeFilter) params.type = typeFilter;

            const res = await couponApi.getCoupons(params);
            if (res.data?.success && res.data?.data) {
                setCoupons(res.data.data);
                setTotal(res.data.pagination?.total || res.data.data.length || 0);
            } else {
                message.error(res.data?.message || '加载失败');
            }
        } catch (err) {
            logger.error('Failed to load coupons:', err);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, userIdFilter, stateFilter, typeFilter]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleViewDetail = (coupon: CouponWithTemplate) => {
        setDetailCoupon(coupon);
        setDetailVisible(true);
    };

    const columns: ColumnsType<CouponWithTemplate> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '用户',
            key: 'user',
            width: 150,
            render: (_, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{record.user?.name || `用户${record.userId}`}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        ID: {record.userId}
                    </Text>
                </Space>
            ),
        },
        {
            title: '优惠券',
            dataIndex: 'name',
            key: 'name',
            width: 150,
            render: (name, record) => (
                <Space direction="vertical" size={0}>
                    <Text strong>{name}</Text>
                    <Space size={4}>
                        <Tag color="blue" style={{ margin: 0 }}>
                            {getCouponTypeLabel(record.type)}
                        </Tag>
                        <Tag style={{ margin: 0 }}>
                            {getCouponSourceLabel(record.source)}
                        </Tag>
                    </Space>
                </Space>
            ),
        },
        {
            title: '优惠内容',
            key: 'discount',
            width: 120,
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
                            满 ¥{centsToYuan(record.minAmountCents)}
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
            title: '状态',
            dataIndex: 'state',
            key: 'state',
            width: 90,
            render: (state) => (
                <Tag color={getCouponStateColor(state as CouponState)}>
                    {getCouponStateLabel(state as CouponState)}
                </Tag>
            ),
        },
        {
            title: '领取时间',
            dataIndex: 'claimedAt',
            key: 'claimedAt',
            width: 150,
            render: (date) => (date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'),
        },
        {
            title: '过期时间',
            dataIndex: 'expireAt',
            key: 'expireAt',
            width: 150,
            render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm'),
        },
        {
            title: '使用时间',
            dataIndex: 'usedAt',
            key: 'usedAt',
            width: 150,
            render: (date) => (date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'),
        },
        {
            title: '优惠金额',
            dataIndex: 'discountCents',
            key: 'discountCents',
            width: 100,
            render: (cents) => (
                <Text strong style={{ color: cents > 0 ? '#f5222d' : undefined }}>
                    {cents > 0 ? `¥${centsToYuan(cents)}` : '-'}
                </Text>
            ),
        },
        {
            title: '操作',
            key: 'action',
            width: 80,
            fixed: 'right',
            render: (_, record) => (
                <Space size={4}>
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                </Space>
            ),
        },
    ];

    // Calculate statistics from current data
    const stats = useMemo(() => {
        const available = coupons.filter((c) => c.state === 'available').length;
        const used = coupons.filter((c) => c.state === 'used').length;
        const expired = coupons.filter((c) => c.state === 'expired').length;
        const locked = coupons.filter((c) => c.state === 'locked').length;
        return { available, used, expired, locked };
    }, [coupons]);

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}>
                <GiftOutlined /> 用户优惠券管理
            </Title>

            {/* Statistics */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="可用"
                            value={stats.available}
                            valueStyle={{ color: token.colorSuccess }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已使用"
                            value={stats.used}
                            valueStyle={{ color: token.colorPrimary }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已锁定"
                            value={stats.locked}
                            valueStyle={{ color: token.colorWarning }}
                        />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已过期"
                            value={stats.expired}
                            valueStyle={{ color: token.colorTextSecondary }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* Filters */}
            <Card style={{ marginBottom: 16 }}>
                <Row gutter={16} align="middle">
                    <Col flex="auto">
                        <Space wrap>
                            <Input
                                placeholder="搜索用户"
                                style={{ width: 150 }}
                                value={userSearchKeyword}
                                onChange={(e) => setUserSearchKeyword(e.target.value)}
                                onPressEnter={() => {
                                    const user = users.find(
                                        (u) =>
                                            u.name.includes(userSearchKeyword) ||
                                            u.phone.includes(userSearchKeyword)
                                    );
                                    if (user) {
                                        setUserIdFilter(user.id);
                                    } else {
                                        message.warning('未找到该用户');
                                    }
                                }}
                            />
                            <Select
                                placeholder="选择用户"
                                allowClear
                                showSearch
                                style={{ width: 200 }}
                                value={userIdFilter}
                                onChange={(v) => setUserIdFilter(v)}
                                filterOption={(input, option) =>
                                    (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                }
                                options={users.map((u) => ({
                                    label: `${u.name} (${u.phone})`,
                                    value: u.id,
                                }))}
                            />
                            <Select
                                placeholder="状态"
                                allowClear
                                style={{ width: 120 }}
                                value={stateFilter}
                                onChange={(v) => setStateFilter(v || '')}
                            >
                                <Select.Option value="available">可用</Select.Option>
                                <Select.Option value="locked">已锁定</Select.Option>
                                <Select.Option value="used">已使用</Select.Option>
                                <Select.Option value="expired">已过期</Select.Option>
                            </Select>
                            <Select
                                placeholder="类型"
                                allowClear
                                style={{ width: 120 }}
                                value={typeFilter}
                                onChange={(v) => setTypeFilter(v || '')}
                            >
                                <Select.Option value="deduct">满减券</Select.Option>
                                <Select.Option value="discount">折扣券</Select.Option>
                            </Select>
                        </Space>
                    </Col>
                    <Col>
                        <Button icon={<ReloadOutlined />} onClick={loadData}>
                            刷新
                        </Button>
                    </Col>
                </Row>
            </Card>

            {/* Table */}
            <Card>
                <Table
                    columns={columns}
                    dataSource={coupons}
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

            {/* Detail Modal */}
            <Modal
                title="优惠券详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={700}
            >
                {detailCoupon && (
                    <Card type="inner">
                        <Descriptions column={2} bordered size="small">
                            <Descriptions.Item label="优惠券ID" span={2}>
                                {detailCoupon.id}
                            </Descriptions.Item>
                            <Descriptions.Item label="优惠券名称" span={2}>
                                {detailCoupon.name}
                            </Descriptions.Item>
                            <Descriptions.Item label="用户ID">
                                {detailCoupon.userId}
                            </Descriptions.Item>
                            <Descriptions.Item label="用户名">
                                {detailCoupon.user?.name || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="模板ID">
                                {detailCoupon.templateId}
                            </Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={getCouponStateColor(detailCoupon.state as CouponState)}>
                                    {getCouponStateLabel(detailCoupon.state as CouponState)}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="类型">
                                <Tag color="blue">{getCouponTypeLabel(detailCoupon.type)}</Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="来源">
                                <Tag>{getCouponSourceLabel(detailCoupon.source)}</Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="适用范围">
                                <Tag color="purple">{getCouponScopeLabel(detailCoupon.scope)}</Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="最低消费">
                                ¥{centsToYuan(detailCoupon.minAmountCents)}
                            </Descriptions.Item>
                            {detailCoupon.type === 'deduct' && (
                                <Descriptions.Item label="减免金额" span={2}>
                                    <Text type="danger" strong>
                                        ¥{centsToYuan(detailCoupon.deductAmountCents)}
                                    </Text>
                                </Descriptions.Item>
                            )}
                            {detailCoupon.type === 'discount' && (
                                <>
                                    <Descriptions.Item label="折扣率">
                                        {(detailCoupon.discountRate * 10).toFixed(1)} 折
                                    </Descriptions.Item>
                                    <Descriptions.Item label="最大优惠">
                                        ¥{centsToYuan(detailCoupon.maxDiscountCents)}
                                    </Descriptions.Item>
                                </>
                            )}
                            <Descriptions.Item label="领取时间">
                                {detailCoupon.claimedAt
                                    ? dayjs(detailCoupon.claimedAt).format('YYYY-MM-DD HH:mm:ss')
                                    : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="过期时间">
                                {dayjs(detailCoupon.expireAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            <Descriptions.Item label="使用时间">
                                {detailCoupon.usedAt
                                    ? dayjs(detailCoupon.usedAt).format('YYYY-MM-DD HH:mm:ss')
                                    : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="锁定订单">
                                {detailCoupon.lockedByOrderId || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="使用订单">
                                {detailCoupon.usedOrderId || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="实际优惠金额" span={2}>
                                <Text strong style={{ color: detailCoupon.discountCents > 0 ? '#f5222d' : undefined }}>
                                    {detailCoupon.discountCents > 0
                                        ? `¥${centsToYuan(detailCoupon.discountCents)}`
                                        : '-'}
                                </Text>
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间" span={2}>
                                {dayjs(detailCoupon.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                        </Descriptions>
                    </Card>
                )}
            </Modal>
        </div>
    );
};

export default UserCouponPage;
