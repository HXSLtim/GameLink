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
    Rate,
    Card,
    Row,
    Col,
    Statistic,
    Form,
    Input,
    Select,
    Typography,
    Divider,
    Image,
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
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { PLAYER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import dayjs from 'dayjs';

const { Text, Paragraph } = Typography;

/**
 * 陪玩师数据接口
 */
interface Player {
    id: number;
    userId: number;
    name: string;
    avatar: string;
    gender: 'male' | 'female';
    age: number;
    games: string[];
    introduction: string;
    voiceSample: string;
    idCardFront: string;
    idCardBack: string;
    status: 'pending' | 'approved' | 'rejected' | 'banned';
    rating: number;
    orderCount: number;
    totalEarnings: number;
    onlineStatus: 'online' | 'offline' | 'busy';
    createdAt: string;
    updatedAt: string;
    auditRemark?: string;
}

/**
 * 状态映射
 */
const statusMap = {
    pending: { color: 'gold', text: '待审核' },
    approved: { color: 'success', text: '已通过' },
    rejected: { color: 'error', text: '已拒绝' },
    banned: { color: 'default', text: '已封禁' },
};

/**
 * 在线状态映射
 */
const onlineStatusMap = {
    online: { color: 'success', text: '在线' },
    offline: { color: 'default', text: '离线' },
    busy: { color: 'warning', text: '忙碌' },
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

    // 弹窗状态
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [auditModalVisible, setAuditModalVisible] = useState(false);
    const [currentPlayer, setCurrentPlayer] = useState<Player | null>(null);
    const [auditForm] = Form.useForm();

    /**
     * 加载陪玩师数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        await new Promise(resolve => setTimeout(resolve, 500));

        const mockPlayers: Player[] = Array.from({ length: 25 }, (_, i) => ({
            id: i + 1,
            userId: 100 + i,
            name: `陪玩师${i + 1}`,
            avatar: '',
            gender: i % 2 === 0 ? 'female' : 'male',
            age: 18 + (i % 10),
            games: [['王者荣耀', '英雄联盟'], ['和平精英', '原神'], ['王者荣耀', '和平精英']][i % 3],
            introduction: `大家好，我是陪玩师${i + 1}，擅长多种游戏，期待与您一起游戏！`,
            voiceSample: '',
            idCardFront: 'https://via.placeholder.com/200x120',
            idCardBack: 'https://via.placeholder.com/200x120',
            status: ['pending', 'approved', 'approved', 'rejected', 'banned'][i % 5] as Player['status'],
            rating: 4 + Math.random(),
            orderCount: Math.floor(Math.random() * 200),
            totalEarnings: Math.floor(Math.random() * 50000),
            onlineStatus: ['online', 'offline', 'busy'][i % 3] as Player['onlineStatus'],
            createdAt: dayjs().subtract(i, 'day').format('YYYY-MM-DD HH:mm:ss'),
            updatedAt: dayjs().subtract(i, 'hour').format('YYYY-MM-DD HH:mm:ss'),
            auditRemark: i % 5 === 3 ? '身份信息不清晰，请重新提交' : undefined,
        }));

        const start = (current - 1) * pageSize;
        const end = start + pageSize;
        setPlayers(mockPlayers.slice(start, end));
        setTotal(mockPlayers.length);
        setLoading(false);
    }, [current, pageSize]);

    useEffect(() => {
        loadData();
    }, [loadData]);

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
        try {
            const values = await auditForm.validateFields();
            console.log('Audit:', { approved, ...values });
            message.success(approved ? '审核通过' : '审核拒绝');
            setAuditModalVisible(false);
            loadData();
        } catch {
            // 验证失败
        }
    };

    /**
     * 封禁/解封
     */
    const handleToggleBan = async (player: Player) => {
        const action = player.status === 'banned' ? '解封' : '封禁';
        message.success(`${action}成功`);
        loadData();
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
        {
            name: 'onlineStatus',
            label: '在线状态',
            type: 'select',
            options: Object.entries(onlineStatusMap).map(([key, val]) => ({ label: val.text, value: key })),
        },
        { name: 'dateRange', label: '申请时间', type: 'dateRange' },
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
                        src={record.avatar}
                        icon={<UserOutlined />}
                        style={{ backgroundColor: record.gender === 'female' ? '#eb2f96' : '#1890ff' }}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.name}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            {record.gender === 'female' ? '女' : '男'} · {record.age}岁
                        </Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '擅长游戏',
            dataIndex: 'games',
            key: 'games',
            width: 180,
            render: games => (
                <Space size={4} wrap>
                    {games.map((g: string) => <Tag key={g}>{g}</Tag>)}
                </Space>
            ),
        },
        {
            title: '评分',
            dataIndex: 'rating',
            key: 'rating',
            width: 120,
            render: rating => (
                <Space>
                    <StarOutlined style={{ color: '#faad14' }} />
                    <span>{rating.toFixed(1)}</span>
                </Space>
            ),
        },
        {
            title: '订单数',
            dataIndex: 'orderCount',
            key: 'orderCount',
            width: 100,
            render: count => `${count} 单`,
        },
        {
            title: '收益',
            dataIndex: 'totalEarnings',
            key: 'totalEarnings',
            width: 120,
            render: earnings => <Text strong>¥{earnings.toLocaleString()}</Text>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: status => <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>,
        },
        {
            title: '在线',
            dataIndex: 'onlineStatus',
            key: 'onlineStatus',
            width: 80,
            render: status => <Tag color={onlineStatusMap[status].color}>{onlineStatusMap[status].text}</Tag>,
        },
        {
            title: '申请时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
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
                    {record.status === 'pending' && (
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
                    {record.status === 'approved' && (
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
                    {record.status === 'banned' && (
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

    return (
        <PageContainer title="陪玩师管理" subTitle="管理平台陪玩师认证和审核">
            <SearchTable
                columns={columns}
                dataSource={players}
                rowKey="id"
                searchFields={searchFields}
                onSearch={() => loadData()}
                onRefresh={loadData}
                loading={loading}
                showCreate={false}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: total => `共 ${total} 条`,
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
                width={600}
            >
                {currentPlayer && (
                    <>
                        {/* 基本信息卡片 */}
                        <Card size="small">
                            <div style={{ textAlign: 'center', marginBottom: 16 }}>
                                <Avatar
                                    size={80}
                                    src={currentPlayer.avatar}
                                    icon={<UserOutlined />}
                                    style={{ backgroundColor: currentPlayer.gender === 'female' ? '#eb2f96' : '#1890ff' }}
                                />
                                <h2 style={{ marginTop: 12, marginBottom: 4 }}>{currentPlayer.name}</h2>
                                <Space>
                                    <Tag color={statusMap[currentPlayer.status].color}>{statusMap[currentPlayer.status].text}</Tag>
                                    <Tag color={onlineStatusMap[currentPlayer.onlineStatus].color}>{onlineStatusMap[currentPlayer.onlineStatus].text}</Tag>
                                </Space>
                            </div>

                            <Row gutter={16}>
                                <Col span={8}>
                                    <Statistic title="评分" value={currentPlayer.rating.toFixed(1)} prefix={<StarOutlined />} />
                                </Col>
                                <Col span={8}>
                                    <Statistic title="订单数" value={currentPlayer.orderCount} suffix="单" />
                                </Col>
                                <Col span={8}>
                                    <Statistic title="总收益" value={currentPlayer.totalEarnings} prefix="¥" />
                                </Col>
                            </Row>
                        </Card>

                        <Divider />

                        {/* 详细信息 */}
                        <Descriptions title="基本信息" column={2} size="small">
                            <Descriptions.Item label="ID">{currentPlayer.id}</Descriptions.Item>
                            <Descriptions.Item label="用户ID">{currentPlayer.userId}</Descriptions.Item>
                            <Descriptions.Item label="性别">{currentPlayer.gender === 'female' ? '女' : '男'}</Descriptions.Item>
                            <Descriptions.Item label="年龄">{currentPlayer.age}岁</Descriptions.Item>
                            <Descriptions.Item label="擅长游戏" span={2}>
                                <Space wrap>
                                    {currentPlayer.games.map(g => <Tag key={g}>{g}</Tag>)}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="自我介绍" span={2}>
                                <Paragraph ellipsis={{ rows: 3, expandable: true }}>
                                    {currentPlayer.introduction}
                                </Paragraph>
                            </Descriptions.Item>
                            <Descriptions.Item label="申请时间">{currentPlayer.createdAt}</Descriptions.Item>
                            <Descriptions.Item label="更新时间">{currentPlayer.updatedAt}</Descriptions.Item>
                            {currentPlayer.auditRemark && (
                                <Descriptions.Item label="审核备注" span={2}>
                                    <Text type="danger">{currentPlayer.auditRemark}</Text>
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <Divider />

                        {/* 身份证照片 */}
                        <Descriptions title="身份验证" column={2}>
                            <Descriptions.Item label="身份证正面">
                                <Image
                                    width={150}
                                    src={currentPlayer.idCardFront}
                                    fallback="https://via.placeholder.com/150x100"
                                />
                            </Descriptions.Item>
                            <Descriptions.Item label="身份证反面">
                                <Image
                                    width={150}
                                    src={currentPlayer.idCardBack}
                                    fallback="https://via.placeholder.com/150x100"
                                />
                            </Descriptions.Item>
                        </Descriptions>
                    </>
                )}
            </Drawer>

            {/* 审核弹窗 */}
            <Modal
                title="审核陪玩师申请"
                open={auditModalVisible}
                onCancel={() => setAuditModalVisible(false)}
                footer={
                    <Space>
                        <Button onClick={() => setAuditModalVisible(false)}>取消</Button>
                        <Button danger icon={<CloseOutlined />} onClick={() => handleAudit(false)}>
                            拒绝
                        </Button>
                        <Button type="primary" icon={<CheckOutlined />} onClick={() => handleAudit(true)}>
                            通过
                        </Button>
                    </Space>
                }
            >
                <Descriptions column={1} size="small">
                    <Descriptions.Item label="申请人">{currentPlayer?.name}</Descriptions.Item>
                    <Descriptions.Item label="擅长游戏">
                        {currentPlayer?.games.join('、')}
                    </Descriptions.Item>
                </Descriptions>

                <Divider />

                <Form form={auditForm} layout="vertical">
                    <Form.Item name="remark" label="审核备注">
                        <Input.TextArea rows={3} placeholder="请输入审核备注（选填）" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default PlayerPage;
