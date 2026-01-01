/**
 * 团队成员管理页面
 * 全局成员列表，跨团队查看和管理成员
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    message,
    Avatar,
    Card,
    Row,
    Col,
    Statistic,
    Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    TeamOutlined,
    UserOutlined,
    CrownOutlined,
    SearchOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type SearchField } from '@/components';
import { teamApi, type TeamMember, type Team } from '@/api/team';
import dayjs from 'dayjs';

const { Text } = Typography;

/**
 * 状态映射
 */
const statusMap: Record<TeamMember['status'], { color: string; text: string }> = {
    active: { color: 'success', text: '正常' },
    left: { color: 'default', text: '已退出' },
    kicked: { color: 'error', text: '已移除' },
};

/**
 * 角色映射
 */
const roleMap: Record<TeamMember['role'], { color: string; text: string }> = {
    leader: { color: 'gold', text: '队长' },
    member: { color: 'blue', text: '成员' },
};

/**
 * 成员管理页面
 */
const MembersPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [members, setMembers] = useState<TeamMember[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 详情弹窗状态
    const [detailVisible, setDetailVisible] = useState(false);
    const [selectedMember, setSelectedMember] = useState<TeamMember | null>(null);
    const [memberTeam, setMemberTeam] = useState<Team | null>(null);

    /**
     * 加载成员数据
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
            const response = await teamApi.listMembers(queryParams);
            if (response.data.success) {
                setMembers(response.data.data?.items || []);
                setTotal(response.data.data?.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load members error:', error);
            message.error('加载成员列表失败');
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
     * 查看成员详情
     */
    const handleViewDetail = async (member: TeamMember) => {
        setSelectedMember(member);
        setDetailVisible(true);

        // 加载所属团队信息
        try {
            const response = await teamApi.getTeamDetail(member.teamId);
            if (response.data.success) {
                setMemberTeam(response.data.data);
            }
        } catch (error) {
            console.error('Load team error:', error);
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '成员名称/ID' },
        {
            name: 'role',
            label: '角色',
            type: 'select',
            options: [
                { label: '队长', value: 'leader' },
                { label: '成员', value: 'member' },
            ],
        },
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
    const columns: ColumnsType<TeamMember> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '成员信息',
            key: 'member',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size={40}
                        src={record.player?.avatar || undefined}
                        icon={<UserOutlined />}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>
                            {record.player?.nickname || `ID:${record.playerId}`}
                        </div>
                        {record.player?.rank && (
                            <Tag color="blue" style={{ fontSize: 11 }}>{record.player.rank}</Tag>
                        )}
                    </div>
                </Space>
            ),
        },
        {
            title: '所属团队',
            dataIndex: 'teamId',
            key: 'teamId',
            width: 180,
            render: (teamId: number) => (
                <Space>
                    <TeamOutlined />
                    <span>团队 #{teamId}</span>
                </Space>
            ),
        },
        {
            title: '角色',
            dataIndex: 'role',
            key: 'role',
            width: 100,
            render: (role: TeamMember['role']) => (
                <Tag color={roleMap[role].color} icon={role === 'leader' ? <CrownOutlined /> : undefined}>
                    {roleMap[role].text}
                </Tag>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: TeamMember['status']) => (
                <Tag color={statusMap[status].color}>
                    {statusMap[status].text}
                </Tag>
            ),
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 80,
            align: 'center',
        },
        {
            title: '统计数据',
            key: 'stats',
            width: 150,
            render: (_, record) => (
                <div style={{ fontSize: 12 }}>
                    <div>订单: {record.orderCount}</div>
                    <div>收益: ¥{(record.incomeCents / 100).toFixed(2)}</div>
                </div>
            ),
        },
        {
            title: '加入时间',
            dataIndex: 'joinedAt',
            key: 'joinedAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '退出时间',
            dataIndex: 'leftAt',
            key: 'leftAt',
            width: 180,
            render: date => date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 120,
            fixed: 'right',
            render: (_, record) => (
                <Button
                    type="link"
                    size="small"
                    icon={<SearchOutlined />}
                    onClick={() => handleViewDetail(record)}
                >
                    详情
                </Button>
            ),
        },
    ];

    return (
        <PageContainer
            title="成员管理"
            subTitle="全局团队成员列表"
        >
            <SearchTable
                columns={columns}
                dataSource={members}
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
                scroll={{ x: 1400 }}
            />

            {/* 成员详情弹窗 */}
            <Modal
                title="成员详情"
                open={detailVisible}
                onCancel={() => setDetailVisible(false)}
                footer={null}
                width={600}
            >
                {selectedMember && (
                    <>
                        {/* 成员基本信息 */}
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <div style={{ textAlign: 'center', marginBottom: 16 }}>
                                <Avatar
                                    size={80}
                                    src={selectedMember.player?.avatar || undefined}
                                    icon={<UserOutlined />}
                                />
                                <h3 style={{ marginTop: 12, marginBottom: 4 }}>
                                    {selectedMember.player?.nickname || `ID:${selectedMember.playerId}`}
                                </h3>
                                <Space>
                                    <Tag color={roleMap[selectedMember.role].color}>
                                        {selectedMember.role === 'leader' && <CrownOutlined />}
                                        {roleMap[selectedMember.role].text}
                                    </Tag>
                                    <Tag color={statusMap[selectedMember.status].color}>
                                        {statusMap[selectedMember.status].text}
                                    </Tag>
                                    {selectedMember.player?.rank && (
                                        <Tag color="blue">{selectedMember.player.rank}</Tag>
                                    )}
                                </Space>
                            </div>

                            <Row gutter={16}>
                                <Col span={12}>
                                    <Statistic
                                        title="订单数"
                                        value={selectedMember.orderCount}
                                        suffix="单"
                                    />
                                </Col>
                                <Col span={12}>
                                    <Statistic
                                        title="总收益"
                                        value={`¥${(selectedMember.incomeCents / 100).toFixed(2)}`}
                                    />
                                </Col>
                            </Row>
                        </Card>

                        {/* 所属团队信息 */}
                        {memberTeam && (
                            <Card
                                size="small"
                                title="所属团队"
                                style={{ marginBottom: 16 }}
                            >
                                <Space>
                                    <Avatar size={40} src={memberTeam.avatarUrl} icon={<TeamOutlined />} />
                                    <div>
                                        <div style={{ fontWeight: 500 }}>{memberTeam.name}</div>
                                        <Text type="secondary" style={{ fontSize: 12 }}>
                                            成员: {memberTeam.memberCount}/{memberTeam.maxMembers}
                                        </Text>
                                    </div>
                                </Space>
                            </Card>
                        )}

                        {/* 详细信息 */}
                        <Card size="small" title="详细信息">
                            <Row gutter={[16, 16]}>
                                <Col span={12}>
                                    <Text type="secondary">成员ID</Text>
                                    <div>{selectedMember.id}</div>
                                </Col>
                                <Col span={12}>
                                    <Text type="secondary">团队ID</Text>
                                    <div>{selectedMember.teamId}</div>
                                </Col>
                                <Col span={12}>
                                    <Text type="secondary">陪玩师ID</Text>
                                    <div>{selectedMember.playerId}</div>
                                </Col>
                                <Col span={12}>
                                    <Text type="secondary">排序</Text>
                                    <div>{selectedMember.sortOrder}</div>
                                </Col>
                                <Col span={12}>
                                    <Text type="secondary">加入时间</Text>
                                    <div>{selectedMember.joinedAt ? dayjs(selectedMember.joinedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}</div>
                                </Col>
                                <Col span={12}>
                                    <Text type="secondary">退出时间</Text>
                                    <div>{selectedMember.leftAt ? dayjs(selectedMember.leftAt).format('YYYY-MM-DD HH:mm:ss') : '-'}</div>
                                </Col>
                            </Row>
                        </Card>
                    </>
                )}
            </Modal>
        </PageContainer>
    );
};

export default MembersPage;
