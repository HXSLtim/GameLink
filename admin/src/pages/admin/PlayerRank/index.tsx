/**
 * 段位审核页面
 * 审核陪玩师的段位认证申请
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
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { adminApi, type PlayerRankRecord, type Game } from '@/api/admin';
import dayjs from 'dayjs';
import { logger } from '@/utils/logger';

const { Text } = Typography;

const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'gold', text: '待审核' },
    verified: { color: 'success', text: '已通过' },
    rejected: { color: 'error', text: '已拒绝' },
    revoked: { color: 'default', text: '已撤销' },
    expired: { color: 'default', text: '已过期' },
};

const PlayerRankPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [records, setRecords] = useState<PlayerRankRecord[]>([]);
    const [games, setGames] = useState<Game[]>([]);
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
    const [currentRecord, setCurrentRecord] = useState<PlayerRankRecord | null>(null);
    const [verifyForm] = Form.useForm();
    const [verifyLoading, setVerifyLoading] = useState(false);

    // 加载游戏列表
    const loadGames = useCallback(async () => {
        try {
            const response = await adminApi.getGames({ page_size: 1000 });
            if (response.data.success) {
                setGames(response.data.data || []);
            }
        } catch (error) {
            logger.error('Load games error:', error);
        }
    }, []);

    // 加载统计数据
    const loadStats = useCallback(async () => {
        try {
            const response = await adminApi.getPlayerRankStats();
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

            const response = await adminApi.getPlayerRanks(queryParams);
            if (response.data.success) {
                setRecords(response.data.data || []);
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            logger.error('Load player ranks error:', error);
            message.error('加载段位认证列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams, activeTab]);

    useEffect(() => {
        loadGames();
        loadStats();
    }, [loadGames, loadStats]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    const handleViewDetail = (record: PlayerRankRecord) => {
        setCurrentRecord(record);
        setDetailVisible(true);
    };

    const handleOpenVerify = (record: PlayerRankRecord) => {
        setCurrentRecord(record);
        verifyForm.resetFields();
        setVerifyVisible(true);
    };

    const handleVerify = async (approved: boolean) => {
        if (!currentRecord) return;
        try {
            setVerifyLoading(true);
            const values = await verifyForm.validateFields();
            await adminApi.verifyPlayerRank(currentRecord.id, {
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
            await adminApi.deletePlayerRank(id);
            message.success('删除成功');
            loadData();
            loadStats();
        } catch (error) {
            logger.error('Delete error:', error);
            message.error('删除失败');
        }
    };

    const searchFields: SearchField[] = [
        { name: 'playerId', label: '陪玩师ID', type: 'input', placeholder: '陪玩师ID' },
        {
            name: 'gameId',
            label: '游戏',
            type: 'select',
            options: games.map(g => ({ label: g.name, value: g.id })),
        },
    ];

    const columns: ColumnsType<PlayerRankRecord> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '陪玩师',
            key: 'player',
            width: 240,
            render: (_, record) => (
                <Space>
                    <Avatar size={40} src={record.player?.user?.avatarUrl} icon={<UserOutlined />} />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.player?.nickname || '-'}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            账号: {record.player?.user?.name || '-'} (ID: {record.player?.userId || record.playerId})
                        </Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '游戏',
            key: 'game',
            width: 120,
            render: (_, record) => <Tag color="blue">{record.game?.name || '-'}</Tag>,
        },
        {
            title: '申请段位',
            key: 'rank',
            width: 140,
            render: (_, record) => (
                <Space>
                    {record.rank?.color && (
                        <span style={{
                            display: 'inline-block',
                            width: 12,
                            height: 12,
                            borderRadius: '50%',
                            backgroundColor: record.rank.color,
                            border: '1px solid #d9d9d9',
                        }} />
                    )}
                    <span style={{ fontWeight: 500 }}>{record.rank?.name || '-'}</span>
                </Space>
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
            width: 240,  // 3个按钮 × 80px
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

    const parseScreenshots = (urls?: string): string[] => {
        if (!urls) return [];
        try {
            return JSON.parse(urls);
        } catch {
            return urls.split(',').filter(Boolean);
        }
    };

    return (
        <PageContainer title="段位审核" subTitle="审核陪玩师的段位认证申请">
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
                    dataSource={records}
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
                title="段位认证详情"
                open={detailVisible}
                onClose={() => setDetailVisible(false)}
                size="large"
            >
                {currentRecord && (
                    <>
                        {/* 陪玩师信息卡片 */}
                        <Card size="small" style={{ marginBottom: 16 }} title="陪玩师信息">
                            <Space align="start">
                                <Avatar size={64} src={currentRecord.player?.user?.avatarUrl} icon={<UserOutlined />} />
                                <div style={{ flex: 1 }}>
                                    <div style={{ fontSize: 16, fontWeight: 500, marginBottom: 8 }}>
                                        {currentRecord.player?.nickname || '-'}
                                    </div>
                                    <Descriptions column={2} size="small">
                                        <Descriptions.Item label="用户ID">{currentRecord.player?.userId || '-'}</Descriptions.Item>
                                        <Descriptions.Item label="陪玩师ID">{currentRecord.playerId}</Descriptions.Item>
                                        <Descriptions.Item label="账号名">{currentRecord.player?.user?.name || '-'}</Descriptions.Item>
                                        <Descriptions.Item label="邮箱">{currentRecord.player?.user?.email || '-'}</Descriptions.Item>
                                        <Descriptions.Item label="手机">{currentRecord.player?.user?.phone || '-'}</Descriptions.Item>
                                        <Descriptions.Item label="账号状态">
                                            <Tag color={currentRecord.player?.user?.status === 'active' ? 'success' : 'error'}>
                                                {currentRecord.player?.user?.status === 'active' ? '正常' : currentRecord.player?.user?.status || '-'}
                                            </Tag>
                                        </Descriptions.Item>
                                    </Descriptions>
                                </div>
                            </Space>
                        </Card>

                        {/* 申请信息 */}
                        <Card size="small" style={{ marginBottom: 16 }} title="申请信息">
                            <Descriptions column={1} size="small">
                                <Descriptions.Item label="申请游戏">
                                    <Tag color="blue">{currentRecord.game?.name || '-'}</Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="申请段位">
                                    <Space>
                                        {currentRecord.rank?.color && (
                                            <span style={{
                                                display: 'inline-block',
                                                width: 12,
                                                height: 12,
                                                borderRadius: '50%',
                                                backgroundColor: currentRecord.rank.color,
                                            }} />
                                        )}
                                        <Tag>{currentRecord.rank?.name || '-'}</Tag>
                                    </Space>
                                </Descriptions.Item>
                                <Descriptions.Item label="审核状态">
                                    <Tag color={statusMap[currentRecord.status]?.color}>
                                        {statusMap[currentRecord.status]?.text}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="申请时间">
                                    {currentRecord.createdAt ? dayjs(currentRecord.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                                </Descriptions.Item>
                                {currentRecord.verifiedAt && (
                                    <Descriptions.Item label="审核时间">
                                        {dayjs(currentRecord.verifiedAt).format('YYYY-MM-DD HH:mm:ss')}
                                    </Descriptions.Item>
                                )}
                                {currentRecord.rejectReason && (
                                    <Descriptions.Item label="拒绝原因">
                                        <Text type="danger">{currentRecord.rejectReason}</Text>
                                    </Descriptions.Item>
                                )}
                                {currentRecord.remark && (
                                    <Descriptions.Item label="备注">{currentRecord.remark}</Descriptions.Item>
                                )}
                            </Descriptions>
                        </Card>

                        <Divider style={{ marginTop: 24 }}>段位截图</Divider>
                        <Image.PreviewGroup>
                            <Space wrap>
                                {parseScreenshots(currentRecord.screenshotUrls).map((url, idx) => (
                                    <Image key={idx} width={150} src={url} style={{ borderRadius: 4 }} />
                                ))}
                            </Space>
                        </Image.PreviewGroup>
                        {parseScreenshots(currentRecord.screenshotUrls).length === 0 && (
                            <Text type="secondary">暂无截图</Text>
                        )}

                        {currentRecord.status === 'pending' && (
                            <div style={{ marginTop: 24, textAlign: 'center' }}>
                                <Button
                                    type="primary"
                                    icon={<CheckOutlined />}
                                    onClick={() => {
                                        setDetailVisible(false);
                                        handleOpenVerify(currentRecord);
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
                title="审核段位认证"
                open={verifyVisible}
                onCancel={() => setVerifyVisible(false)}
                width={640}
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
                {currentRecord && (
                    <>
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space align="start">
                                <Avatar size={48} src={currentRecord.player?.user?.avatarUrl} icon={<UserOutlined />} />
                                <div>
                                    <div style={{ fontWeight: 500 }}>{currentRecord.player?.nickname || '-'}</div>
                                    <Space style={{ marginTop: 4 }}>
                                        <Tag color="blue">{currentRecord.game?.name}</Tag>
                                        <Tag>{currentRecord.rank?.name}</Tag>
                                    </Space>
                                </div>
                            </Space>
                        </Card>

                        <div style={{ marginBottom: 16 }}>
                            <Text strong style={{ display: 'block', marginBottom: 8 }}>段位截图</Text>
                            <Image.PreviewGroup>
                                <Space wrap>
                                    {parseScreenshots(currentRecord.screenshotUrls).map((url, idx) => (
                                        <Image key={idx} width={120} src={url} style={{ borderRadius: 4 }} />
                                    ))}
                                </Space>
                            </Image.PreviewGroup>
                        </div>

                        <Form form={verifyForm} layout="vertical">
                            <Form.Item name="rejectReason" label="拒绝原因（拒绝时必填）">
                                <Input.TextArea rows={3} placeholder="请输入拒绝原因" />
                            </Form.Item>
                        </Form>
                    </>
                )}
            </Modal>
        </PageContainer>
    );
};

export default PlayerRankPage;
