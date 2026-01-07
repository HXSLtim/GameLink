/**
 * 陪玩师管理页面
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
    Card,
    Form,
    Input,
    Typography,
    Select,
    Radio,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    CheckOutlined,
    CloseOutlined,
    UserOutlined,
    LockOutlined,
    UnlockOutlined,
    StarOutlined,
    SafetyOutlined,
    DeleteOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { exportToCSV, playerExportColumns } from '@/utils/export';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { PLAYER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type Player, type ApiResponse } from '@/api/admin';
import dayjs from 'dayjs';
import PlayerDetailTabs from './components/PlayerDetailTabs';

import { logger } from '@/utils/logger';
const { Text, Paragraph } = Typography;

/**
 * 状态映射
 */
const statusMap = {
    pending: { color: 'gold', text: '待审核' },
    verified: { color: 'success', text: '已通过' },
    rejected: { color: 'error', text: '已拒绝' },
};

/**
 * 陪玩师管理页面
 */
const PlayerPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [players, setPlayers] = useState<Player[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [auditModalVisible, setAuditModalVisible] = useState(false);
    const [currentPlayer, setCurrentPlayer] = useState<Player | null>(null);
    const [auditForm] = Form.useForm();

    // 批量操作状态
    const [batchStatusVisible, setBatchStatusVisible] = useState(false);
    const [batchDeleteVisible, setBatchDeleteVisible] = useState(false);
    const [selectedPlayerIds, setSelectedPlayerIds] = useState<number[]>([]);
    const [batchTarget, setBatchTarget] = useState<'selected' | 'status' | 'all'>('selected');
    const [batchForm] = Form.useForm();

    /**
     * 加载陪玩师数据
     */
    const loadData = useCallback(async (params?: Record<string, unknown>) => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
                ...params,
            };
            const response = await adminApi.getPlayers(queryParams);
            if (response.data.success) {
                setPlayers(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load players error:', error);
            message.error('加载陪玩师列表失败');
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
     * 查看详情
     */
    const handleViewDetail = (player: Player) => {
        setCurrentPlayer(player);
        setDetailDrawerVisible(true);
    };

    /**
     * 打开审核弹窗
     */
    const handleOpenAudit = (player: Player) => {
        setCurrentPlayer(player);
        auditForm.resetFields();
        setAuditModalVisible(true);
    };

    /**
     * 执行审核
     */
    const handleAudit = async (approved: boolean) => {
        if (!currentPlayer) return;
        try {
            const values = await auditForm.validateFields();
            const status = approved ? 'verified' : 'rejected';
            // 传递审核备注到后端保存
            await adminApi.updatePlayerVerification(currentPlayer.id, status, values.remark);
            message.success(approved ? '审核通过' : '审核拒绝');
            setAuditModalVisible(false);
            loadData();
        } catch (error) {
            logger.error('Audit error:', error);
            message.error('审核操作失败');
        }
    };

    /**
     * 封禁/解封
     */
    const handleToggleBan = async (player: Player) => {
        try {
            const newStatus = player.verificationStatus === 'rejected' ? 'verified' : 'rejected';
            await adminApi.updatePlayerVerification(player.id, newStatus);
            const action = player.verificationStatus === 'rejected' ? '解封' : '封禁';
            message.success(`${action}成功`);
            loadData();
        } catch (error) {
            logger.error('Toggle ban error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 批量修改状态
     */
    const handleBatchStatus = (keys: React.Key[]) => {
        setSelectedPlayerIds(keys ? keys.map(k => Number(k)) : []);
        batchForm.resetFields();
        batchForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'selected' : 'all',
        });
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchStatusVisible(true);
    };

    const submitBatchStatus = async () => {
        try {
            const values = await batchForm.validateFields();
            let playerIds: number[] = [];

            if (values.target === 'selected') {
                playerIds = selectedPlayerIds;
            } else if (values.target === 'status') {
                // 按状态筛选 - 需要获取符合条件的陪玩师ID
                const response = await adminApi.getPlayers({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    playerIds = response.data.data.map((p: Player) => p.id);
                }
            } else {
                // 全部 - 获取所有陪玩师ID
                const response = await adminApi.getPlayers({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    playerIds = response.data.data.map((p: Player) => p.id);
                }
            }

            if (playerIds.length === 0) {
                message.warning('没有符合条件的陪玩师');
                return;
            }

            const res = await adminApi.batchUpdatePlayerStatus({
                playerIds,
                status: values.status,
            }) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success(`批量修改 ${playerIds.length} 个陪玩师状态成功`);
                setBatchStatusVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * 批量删除
     */
    const handleBatchDelete = (keys: React.Key[]) => {
        setSelectedPlayerIds(keys ? keys.map(k => Number(k)) : []);
        batchForm.resetFields();
        batchForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'selected' : 'all',
        });
        setBatchTarget((keys && keys.length > 0) ? 'selected' : 'all');
        setBatchDeleteVisible(true);
    };

    const submitBatchDelete = async () => {
        try {
            const values = await batchForm.validateFields();
            let playerIds: number[] = [];

            if (values.target === 'selected') {
                playerIds = selectedPlayerIds;
            } else if (values.target === 'status') {
                const response = await adminApi.getPlayers({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    playerIds = response.data.data.map((p: Player) => p.id);
                }
            } else {
                const response = await adminApi.getPlayers({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    playerIds = response.data.data.map((p: Player) => p.id);
                }
            }

            if (playerIds.length === 0) {
                message.warning('没有符合条件的陪玩师');
                return;
            }

            const res = await adminApi.batchDeletePlayers(playerIds) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success(`批量删除 ${playerIds.length} 个陪玩师成功`);
                setBatchDeleteVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '名称/ID' },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: Object.entries(statusMap).map(([key, val]) => ({ label: val.text, value: key })),
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Player> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '陪玩师',
            key: 'player',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size={40}
                        src={record.user?.avatarUrl || undefined}
                        icon={<UserOutlined />}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.nickname || record.user?.name || '-'}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            用户ID: {record.userId}
                        </Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '主游戏',
            key: 'mainGame',
            width: 120,
            render: (_, record) => record.mainGame?.name || '-',
        },
        {
            title: '时薪',
            dataIndex: 'hourlyRateCents',
            key: 'hourlyRateCents',
            width: 100,
            render: cents => cents ? `¥${(cents / 100).toFixed(2)}` : '-',
        },
        {
            title: '评分',
            key: 'rating',
            width: 120,
            render: (_, record) => (
                <Space>
                    <StarOutlined style={{ color: '#faad14' }} />
                    <span>{record.ratingAverage?.toFixed(1) || '0.0'}</span>
                    <Text type="secondary">({record.ratingCount || 0})</Text>
                </Space>
            ),
        },
        {
            title: '技能标签',
            dataIndex: 'skillTags',
            key: 'skillTags',
            width: 180,
            render: tags => (
                <Space size={4} wrap>
                    {(tags || []).slice(0, 3).map((tag: string) => <Tag key={tag}>{tag}</Tag>)}
                    {(tags || []).length > 3 && <Tag>+{tags.length - 3}</Tag>}
                </Space>
            ),
        },
        {
            title: '状态',
            dataIndex: 'verificationStatus',
            key: 'verificationStatus',
            width: 100,
            render: (status: Player['verificationStatus']) => (
                <Tag color={statusMap[status]?.color || 'default'}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
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
                        icon={<EyeOutlined />}
                        onClick={() => handleViewDetail(record)}
                    >
                        详情
                    </Button>
                    {record.verificationStatus === 'pending' && (
                        <PermissionGuard permission={PLAYER_PERMISSIONS.AUDIT}>
                            <Button
                                type="link"
                                size="small"
                                icon={<CheckOutlined />}
                                onClick={() => handleOpenAudit(record)}
                            >
                                审核
                            </Button>
                        </PermissionGuard>
                    )}
                    {record.verificationStatus === 'verified' && (
                        <PermissionGuard permission={PLAYER_PERMISSIONS.UPDATE}>
                            <Popconfirm
                                title="确定要封禁该陪玩师吗？"
                                onConfirm={() => handleToggleBan(record)}
                            >
                                <Button type="link" size="small" danger icon={<LockOutlined />}>
                                    封禁
                                </Button>
                            </Popconfirm>
                        </PermissionGuard>
                    )}
                    {record.verificationStatus === 'rejected' && (
                        <PermissionGuard permission={PLAYER_PERMISSIONS.UPDATE}>
                            <Popconfirm
                                title="确定要解封该陪玩师吗？"
                                onConfirm={() => handleToggleBan(record)}
                            >
                                <Button type="link" size="small" icon={<UnlockOutlined />}>
                                    解封
                                </Button>
                            </Popconfirm>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    /**
     * 导出陪玩师数据
     */
    const handleExport = async () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            const response = await adminApi.getPlayers({ ...searchParams, page_size: 10000 });
            if (response.data.success && response.data.data) {
                exportToCSV(response.data.data as unknown as Record<string, unknown>[], playerExportColumns, 'players');
                message.success({ content: '导出成功', key: 'export' });
            } else {
                message.error({ content: '导出失败', key: 'export' });
            }
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    };

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量修改状态',
            icon: <SafetyOutlined />,
            needSelection: false,
            onClick: (keys) => handleBatchStatus(keys || []),
            permission: PLAYER_PERMISSIONS.UPDATE,
        },
        {
            text: '批量删除',
            icon: <DeleteOutlined />,
            needSelection: false,
            danger: true,
            onClick: (keys) => handleBatchDelete(keys || []),
            permission: PLAYER_PERMISSIONS.DELETE,
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
            permission: PLAYER_PERMISSIONS.LIST,
        },
    ];

    return (
        <PageContainer title="陪玩师管理" subTitle="管理平台陪玩师认证和审核">
            <SearchTable
                columns={columns}
                dataSource={players}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={false}
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

            {/* 详情抽屉 */}
            <Drawer
                title="陪玩师详情"
                open={detailDrawerVisible}
                onClose={() => setDetailDrawerVisible(false)}
                size="large"
            >
                {currentPlayer && <PlayerDetailTabs player={currentPlayer} />}
            </Drawer>

            {/* 审核弹窗 */}
            <Modal
                title="审核陪玩师申请"
                open={auditModalVisible}
                onCancel={() => setAuditModalVisible(false)}
                width={600}
                footer={
                    <Space>
                        <Button onClick={() => setAuditModalVisible(false)}>取消</Button>
                        <PermissionGuard permission={PLAYER_PERMISSIONS.AUDIT}>
                            <Button danger icon={<CloseOutlined />} onClick={() => handleAudit(false)}>
                                拒绝
                            </Button>
                        </PermissionGuard>
                        <PermissionGuard permission={PLAYER_PERMISSIONS.AUDIT}>
                            <Button type="primary" icon={<CheckOutlined />} onClick={() => handleAudit(true)}>
                                通过
                            </Button>
                        </PermissionGuard>
                    </Space>
                }
            >
                {currentPlayer && (
                    <>
                        {/* 申请人基本信息 */}
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space align="start">
                                <Avatar
                                    size={64}
                                    src={currentPlayer.user?.avatarUrl || undefined}
                                    icon={<UserOutlined />}
                                />
                                <div>
                                    <div style={{ fontSize: 16, fontWeight: 500 }}>
                                        {currentPlayer.nickname || currentPlayer.user?.name || '-'}
                                    </div>
                                    <Text type="secondary">用户ID: {currentPlayer.userId}</Text>
                                    <div style={{ marginTop: 8 }}>
                                        {currentPlayer.mainGame && (
                                            <Tag color="blue">{currentPlayer.mainGame.name}</Tag>
                                        )}
                                        {currentPlayer.rank && <Tag>{currentPlayer.rank}</Tag>}
                                    </div>
                                </div>
                            </Space>
                        </Card>

                        {/* 详细信息 */}
                        <Descriptions column={1} size="small" bordered style={{ marginBottom: 16 }}>
                            <Descriptions.Item label="个人简介">
                                <Paragraph ellipsis={{ rows: 3, expandable: true }}>
                                    {currentPlayer.bio || '暂无介绍'}
                                </Paragraph>
                            </Descriptions.Item>
                            <Descriptions.Item label="期望时薪">
                                {currentPlayer.hourlyRateCents ? `¥${(currentPlayer.hourlyRateCents / 100).toFixed(2)}` : '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label="申请时间">
                                {currentPlayer.createdAt ? dayjs(currentPlayer.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                            </Descriptions.Item>
                        </Descriptions>

                        <Divider />

                        <Form form={auditForm} layout="vertical">
                            <Form.Item name="remark" label="审核备注">
                                <Input.TextArea rows={3} placeholder="请输入审核备注（选填，拒绝时建议填写原因）" />
                            </Form.Item>
                        </Form>
                    </>
                )}
            </Modal>

            {/* 批量修改状态弹窗 */}
            <Modal
                title="批量修改陪玩师状态"
                open={batchStatusVisible}
                onOk={submitBatchStatus}
                onCancel={() => setBatchStatusVisible(false)}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedPlayerIds.length === 0}>
                                选中的陪玩师 {selectedPlayerIds.length > 0 ? `(${selectedPlayerIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部陪玩师</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true, message: '请选择筛选状态' }]}>
                            <Select placeholder="请选择要筛选的状态">
                                <Select.Option value="pending">待审核</Select.Option>
                                <Select.Option value="verified">已通过</Select.Option>
                                <Select.Option value="rejected">已拒绝</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <Form.Item name="status" label="修改为" rules={[{ required: true, message: '请选择目标状态' }]}>
                        <Select placeholder="请选择目标状态">
                            <Select.Option value="pending">待审核</Select.Option>
                            <Select.Option value="verified">已通过</Select.Option>
                            <Select.Option value="rejected">已拒绝</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 批量删除弹窗 */}
            <Modal
                title="批量删除陪玩师"
                open={batchDeleteVisible}
                onOk={submitBatchDelete}
                onCancel={() => setBatchDeleteVisible(false)}
                okText="确认删除"
                okButtonProps={{ danger: true }}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedPlayerIds.length === 0}>
                                选中的陪玩师 {selectedPlayerIds.length > 0 ? `(${selectedPlayerIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部陪玩师</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true, message: '请选择筛选状态' }]}>
                            <Select placeholder="请选择要筛选的状态">
                                <Select.Option value="pending">待审核</Select.Option>
                                <Select.Option value="verified">已通过</Select.Option>
                                <Select.Option value="rejected">已拒绝</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <div style={{ color: '#ff4d4f', marginTop: 16 }}>
                        ⚠️ 警告：此操作不可恢复，请谨慎操作！
                    </div>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default PlayerPage;
