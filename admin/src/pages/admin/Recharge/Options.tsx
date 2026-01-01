/**
 * 充值选项管理页面
 * 表格形式展示充值档位，支持增删改查操作
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
    Switch,
    Typography,
    Image,
    InputNumber,
    Tooltip,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    ReloadOutlined,
    GiftOutlined,
    StarOutlined,
    DollarOutlined,
    CrownOutlined,
} from '@ant-design/icons';
import { rechargeApi, type RechargeOption, type RechargeOptionQueryParams } from '@/api/recharge';
import { RECHARGE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { MONEY, PAGINATION, LAYOUT, SIZES, TABLE, MODAL, COLORS, BUSINESS } from '@/constants/common';
import OptionForm from './components/OptionForm';
import dayjs from 'dayjs';

const { Text } = Typography;

interface OptionsProps {
    onStatsUpdate?: () => void;
}

/**
 * 状态选项
 */
const activeOptions = [
    { label: '全部', value: undefined },
    { label: '已启用', value: true },
    { label: '已禁用', value: false },
];

const recommendedOptions = [
    { label: '全部', value: undefined },
    { label: '推荐', value: true },
    { label: '普通', value: false },
];

/**
 * 充值选项管理页面
 */
const RechargeOptions: React.FC<OptionsProps> = ({ onStatsUpdate }) => {
    const [loading, setLoading] = useState(false);
    const [options, setOptions] = useState<RechargeOption[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<RechargeOptionQueryParams>({});
    const [formVisible, setFormVisible] = useState(false);
    const [currentOption, setCurrentOption] = useState<RechargeOption | null>(null);

    /**
     * 加载充值选项数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const queryParams: RechargeOptionQueryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
            };
            const response = await rechargeApi.getRechargeOptions(queryParams);
            if (response.data.success) {
                const data = response.data.data || [];
                setOptions(data);
                setTotal(data.length);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load recharge options error:', error);
            message.error('加载充值选项失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 新增
     */
    const handleAdd = () => {
        setCurrentOption(null);
        setFormVisible(true);
    };

    /**
     * 编辑
     */
    const handleEdit = (record: RechargeOption) => {
        setCurrentOption(record);
        setFormVisible(true);
    };

    /**
     * 删除
     */
    const handleDelete = async (record: RechargeOption) => {
        try {
            await rechargeApi.deleteRechargeOption(record.id);
            message.success('删除成功');
            loadData();
            onStatsUpdate?.();
        } catch (error) {
            console.error('Delete error:', error);
            message.error('删除失败');
        }
    };

    /**
     * 切换启用状态
     */
    const handleToggleActive = async (record: RechargeOption, isActive: boolean) => {
        try {
            await rechargeApi.toggleRechargeOptionStatus([record.id], isActive);
            message.success(isActive ? '已启用' : '已禁用');
            loadData();
            onStatsUpdate?.();
        } catch (error) {
            console.error('Toggle status error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values as RechargeOptionQueryParams);
        setCurrent(1);
    };

    /**
     * 重置搜索
     */
    const handleReset = () => {
        setSearchParams({});
        setCurrent(1);
    };

    /**
     * 表格列配置
     */
    const columns: ColumnsType<RechargeOption> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: TABLE.COLUMN_WIDTH.SMALL,
        },
        {
            title: '图标',
            dataIndex: 'iconUrl',
            key: 'iconUrl',
            width: TABLE.COLUMN_WIDTH.SMALL,
            render: (iconUrl: string) => (
                <Image
                    src={iconUrl}
                    alt="icon"
                    width={SIZES.AVATAR.MEDIUM}
                    height={SIZES.AVATAR.MEDIUM}
                    preview={false}
                    style={{ borderRadius: SIZES.IMAGE_BORDER_RADIUS }}
                    fallback="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mN8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="
                />
            ),
        },
        {
            title: '名称',
            dataIndex: 'name',
            key: 'name',
            width: TABLE.COLUMN_WIDTH.XXLARGE,
            render: (name: string, record) => (
                <div>
                    <div style={{ fontWeight: 500 }}>{name}</div>
                    {record.tag && (
                        <Tag color="gold" style={{ marginTop: 4 }}>
                            {record.tag}
                        </Tag>
                    )}
                </div>
            ),
        },
        {
            title: '充值金额',
            key: 'amount',
            width: TABLE.COLUMN_WIDTH.XXLARGE,
            render: (_, record) => (
                <div>
                    <div style={{ fontWeight: 500, color: COLORS.INFO }}>
                        ¥{(record.amountCents / MONEY.YUAN_TO_FEN).toFixed(BUSINESS.PRECISION.AMOUNT)}
                    </div>
                    {record.originalCents && record.originalCents > record.amountCents && (
                        <Text delete type="secondary" style={{ fontSize: SIZES.SECONDARY_FONT_SIZE }}>
                            ¥{(record.originalCents / MONEY.YUAN_TO_FEN).toFixed(BUSINESS.PRECISION.AMOUNT)}
                        </Text>
                    )}
                </div>
            ),
        },
        {
            title: '赠送金额',
            dataIndex: 'bonusCents',
            key: 'bonusCents',
            width: TABLE.COLUMN_WIDTH.XLARGE,
            render: (cents: number) => (
                <div style={{ color: COLORS.SUCCESS, fontWeight: 500 }}>
                    {cents > 0 ? (
                        <>
                            <GiftOutlined style={{ marginRight: 4 }} />
                            ¥{(cents / MONEY.YUAN_TO_FEN).toFixed(BUSINESS.PRECISION.AMOUNT)}
                        </>
                    ) : (
                        '-'
                    )}
                </div>
            ),
        },
        {
            title: '折扣',
            dataIndex: 'discountPercent',
            key: 'discountPercent',
            width: TABLE.COLUMN_WIDTH.LARGE,
            render: (percent?: number) => {
                if (percent && percent > 0) {
                    return <Tag color="red">{(percent * 100).toFixed(BUSINESS.PRECISION.PERCENT)}% OFF</Tag>;
                }
                return '-';
            },
        },
        {
            title: '优惠券',
            key: 'coupon',
            width: TABLE.COLUMN_WIDTH.XXLARGE,
            render: (_, record) => {
                if (record.couponCount > 0) {
                    return (
                        <Tooltip title={`优惠券模板ID: ${record.couponTemplateId}`}>
                            <Tag color="purple" icon={<GiftOutlined />}>
                                x{record.couponCount}
                            </Tag>
                        </Tooltip>
                    );
                }
                return '-';
            },
        },
        {
            title: 'VIP等级限制',
            dataIndex: 'minVipLevel',
            key: 'minVipLevel',
            width: TABLE.COLUMN_WIDTH.XLARGE,
            render: (level?: number) => {
                if (level !== undefined && level !== null) {
                    return (
                        <Tag color="blue" icon={<CrownOutlined />}>
                            Lv.{level}+
                        </Tag>
                    );
                }
                return <Tag color="default">不限</Tag>;
            },
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: TABLE.COLUMN_WIDTH.SMALL,
            render: (_: unknown, record: RechargeOption) => (
                <InputNumber
                    size="small"
                    value={record.sortOrder}
                    min={0}
                    style={{ width: TABLE.COLUMN_WIDTH.SMALL - 10 }}
                    onChange={async (value) => {
                        if (value !== null && value !== record.sortOrder) {
                            try {
                                await rechargeApi.updateRechargeOption(record.id, { ...record, sortOrder: value });
                                message.success('排序已更新');
                                loadData();
                            } catch {
                                message.error('更新失败');
                            }
                        }
                    }}
                />
            ),
        },
        {
            title: '推荐',
            dataIndex: 'isRecommended',
            key: 'isRecommended',
            width: TABLE.COLUMN_WIDTH.LARGE,
            render: (recommended: boolean) =>
                recommended ? (
                    <Tag color="gold" icon={<StarOutlined />}>
                        推荐
                    </Tag>
                ) : (
                    <Tag>普通</Tag>
                ),
        },
        {
            title: '状态',
            dataIndex: 'isActive',
            key: 'isActive',
            width: TABLE.COLUMN_WIDTH.MEDIUM,
            render: (isActive: boolean) => (
                <Tag color={isActive ? 'success' : 'default'}>
                    {isActive ? '启用' : '禁用'}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: TABLE.COLUMN_WIDTH.XXXLARGE,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: TABLE.COLUMN_WIDTH.ACTION,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <PermissionGuard permission={RECHARGE_PERMISSIONS.UPDATE_OPTION}>
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        >
                            编辑
                        </Button>
                    </PermissionGuard>
                    <PermissionGuard permission={RECHARGE_PERMISSIONS.UPDATE_OPTION}>
                        <Switch
                            size="small"
                            checked={record.isActive}
                            onChange={(checked) => handleToggleActive(record, checked)}
                            checkedChildren="启用"
                            unCheckedChildren="禁用"
                        />
                    </PermissionGuard>
                    <PermissionGuard permission={RECHARGE_PERMISSIONS.DELETE_OPTION}>
                        <Popconfirm
                            title="确定要删除此充值选项吗？"
                            onConfirm={() => handleDelete(record)}
                            okText="确定"
                            cancelText="取消"
                        >
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                            >
                                删除
                            </Button>
                        </Popconfirm>
                    </PermissionGuard>
                </Space>
            ),
        },
    ];

    return (
        <>
            {/* 操作栏 */}
            <Card style={{ marginBottom: LAYOUT.CARD_MARGIN }}>
                <Space wrap>
                    <Button icon={<ReloadOutlined />} onClick={loadData} loading={loading}>
                        刷新
                    </Button>
                    <PermissionGuard permission={RECHARGE_PERMISSIONS.CREATE_OPTION}>
                        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                            新增充值选项
                        </Button>
                    </PermissionGuard>
                </Space>
            </Card>

            {/* 表格 */}
            <Table
                columns={columns}
                dataSource={options}
                rowKey="id"
                loading={loading}
                scroll={{ x: TABLE.SCROLL_WIDTH.LARGE }}
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

            {/* 编辑表单弹窗 */}
            <OptionForm
                visible={formVisible}
                option={currentOption}
                onCancel={() => setFormVisible(false)}
                onSuccess={() => {
                    loadData();
                    onStatsUpdate?.();
                }}
            />
        </>
    );
};

export default RechargeOptions;
