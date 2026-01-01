/**
 * 团队管理页面
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
    Row,
    Col,
    Statistic,
    Form,
    Select,
    Radio,
    InputNumber,
    List,
    Typography,
    Input,
    Divider,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EyeOutlined,
    EditOutlined,
    DeleteOutlined,
    TeamOutlined,
    UserOutlined,
    CrownOutlined,
    PlusOutlined,
    ReloadOutlined,
    CheckOutlined,
    CloseOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { teamApi, type Team, type TeamStats, type TeamMember, type BatchOperationResponse } from '@/api/team';
import TeamForm from './components/TeamForm';
import MemberCard from './components/MemberCard';
import dayjs from 'dayjs';

const { Text, Title } = Typography;

/**
 * 状态映射
 */
const statusMap: Record<Team['status'], { color: string; text: string }> = {
    active: { color: 'success', text: '活跃' },
    busy: { color: 'processing', text: '接单中' },
    inactive: { color: 'default', text: '不活跃' },
};

/**
 * 收益分配类型映射
 */
const shareTypeMap: Record<'equal' | 'custom', { text: string }> = {
    equal: { text: '平均分配' },
    custom: { text: '自定义' },
};

/**
 * 团队管理页面
 */
const TeamPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [teams, setTeams] = useState<Team[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});
    const [stats, setStats] = useState<TeamStats | null>(null);

    // 弹窗状态
    const [formVisible, setFormVisible] = useState(false);
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [currentTeam, setCurrentTeam] = useState<Team | null>(null);

    // 批量操作状态
    const [batchStatusVisible, setBatchStatusVisible] = useState(false);
    const [batchDeleteVisible, setBatchDeleteVisible] = useState(false);
    const [selectedTeamIds, setSelectedTeamIds] = useState<number[]>([]);
    const [batchTarget, setBatchTarget] = useState<'selected' | 'status' | 'all'>('selected');
    const [batchForm] = Form.useForm();

    // 成员管理状态
    const [addMemberVisible, setAddMemberVisible] = useState(false);
    const [transferVisible, setTransferVisible] = useState(false);
    const [memberForm] = Form.useForm();

    /**
     * 加载团队数据
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
            const response = await teamApi.getTeams(queryParams);
            if (response.data.success) {
                setTeams(response.data.data?.items || []);
                setTotal(response.data.data?.pagination?.total || 0);
            } else {
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load teams error:', error);
            message.error('加载团队列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    /**
     * 加载统计数据
     */
    const loadStats = useCallback(async () => {
        try {
            const response = await teamApi.getTeamStats();
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (error) {
            console.error('Load stats error:', error);
        }
    }, []);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
    };

    /**
     * 打开新增弹窗
     */
    const handleAdd = () => {
        setCurrentTeam(null);
        setFormVisible(true);
    };

    /**
     * 打开编辑弹窗
     */
    const handleEdit = (team: Team) => {
        setCurrentTeam(team);
        setFormVisible(true);
    };

    /**
     * 查看详情
     */
    const handleViewDetail = async (team: Team) => {
        setCurrentTeam(team);
        setDetailDrawerVisible(true);
    };

    /**
     * 删除团队
     */
    const handleDelete = async (team: Team) => {
        try {
            await teamApi.deleteTeam(team.id);
            message.success('删除成功');
            loadData();
        } catch (error) {
            console.error('Delete team error:', error);
            message.error('删除失败');
        }
    };

    /**
     * 批量修改状态
     */
    const handleBatchStatus = (keys: React.Key[]) => {
        setSelectedTeamIds(keys ? keys.map(k => Number(k)) : []);
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
            let teamIds: number[] = [];

            if (values.target === 'selected') {
                teamIds = selectedTeamIds;
            } else if (values.target === 'status') {
                const response = await teamApi.getTeams({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    teamIds = response.data.data.items.map((t: Team) => t.id);
                }
            } else {
                const response = await teamApi.getTeams({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    teamIds = response.data.data.items.map((t: Team) => t.id);
                }
            }

            if (teamIds.length === 0) {
                message.warning('没有符合条件的团队');
                return;
            }

            const response = await teamApi.batchUpdateTeamsStatus({
                team_ids: teamIds,
                status: values.status,
            });

            if (response.data.success) {
                const result = response.data.data as BatchOperationResponse;
                message.success(`批量修改完成：成功 ${result.successCount}，失败 ${result.failedCount}`);
                setBatchStatusVisible(false);
                loadData();
                loadStats();
            }
        } catch (error) {
            console.error('Batch status error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 批量删除
     */
    const handleBatchDelete = (keys: React.Key[]) => {
        setSelectedTeamIds(keys ? keys.map(k => Number(k)) : []);
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
            let teamIds: number[] = [];

            if (values.target === 'selected') {
                teamIds = selectedTeamIds;
            } else if (values.target === 'status') {
                const response = await teamApi.getTeams({ status: values.filterStatus, page_size: 1000 });
                if (response.data.success && response.data.data) {
                    teamIds = response.data.data.items.map((t: Team) => t.id);
                }
            } else {
                const response = await teamApi.getTeams({ page_size: 1000 });
                if (response.data.success && response.data.data) {
                    teamIds = response.data.data.items.map((t: Team) => t.id);
                }
            }

            if (teamIds.length === 0) {
                message.warning('没有符合条件的团队');
                return;
            }

            const response = await teamApi.batchDeleteTeams({ team_ids: teamIds });

            if (response.data.success) {
                const result = response.data.data as BatchOperationResponse;
                message.success(`批量删除完成：成功 ${result.successCount}，失败 ${result.failedCount}`);
                setBatchDeleteVisible(false);
                loadData();
                loadStats();
            }
        } catch (error) {
            console.error('Batch delete error:', error);
            message.error('操作失败');
        }
    };

    /**
     * 成员操作
     */
    const handleAddMember = () => {
        memberForm.resetFields();
        setAddMemberVisible(true);
    };

    const submitAddMember = async () => {
        if (!currentTeam) return;
        try {
            const values = await memberForm.validateFields();
            await teamApi.addTeamMember(currentTeam.id, { playerId: values.playerId });
            message.success('添加成员成功');
            setAddMemberVisible(false);
            // 重新加载团队详情
            const response = await teamApi.getTeamDetail(currentTeam.id);
            if (response.data.success) {
                setCurrentTeam(response.data.data);
            }
        } catch (error) {
            console.error('Add member error:', error);
            message.error('添加成员失败');
        }
    };

    const handleRemoveMember = async (member: TeamMember) => {
        if (!currentTeam) return;
        try {
            await teamApi.removeTeamMember(currentTeam.id, member.playerId);
            message.success('移除成员成功');
            // 重新加载团队详情
            const response = await teamApi.getTeamDetail(currentTeam.id);
            if (response.data.success) {
                setCurrentTeam(response.data.data);
            }
        } catch (error) {
            console.error('Remove member error:', error);
            message.error('移除成员失败');
        }
    };

    const handleTransferLeader = (member: TeamMember) => {
        setCurrentTeam(prev => ({ ...prev!, pendingTransfer: member }));
        memberForm.resetFields();
        setTransferVisible(true);
    };

    const submitTransferLeader = async () => {
        if (!currentTeam || !currentTeam.pendingTransfer) return;
        try {
            await teamApi.transferCaptain(currentTeam.id, {
                newLeaderPlayerId: currentTeam.pendingTransfer.playerId,
            });
            message.success('转让队长成功');
            setTransferVisible(false);
            // 重新加载团队详情
            const response = await teamApi.getTeamDetail(currentTeam.id);
            if (response.data.success) {
                setCurrentTeam(response.data.data);
            }
        } catch (error) {
            console.error('Transfer leader error:', error);
            message.error('转让队长失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '团队名称/ID' },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            options: Object.entries(statusMap).map(([key, val]) => ({ label: val.text, value: key })),
        },
        {
            name: 'minMember',
            label: '最少成员数',
            type: 'input-number',
            placeholder: '最少成员数',
            props: { min: 1, max: 50, style: { width: '100%' } },
        },
        {
            name: 'maxMember',
            label: '最多成员数',
            type: 'input-number',
            placeholder: '最多成员数',
            props: { min: 1, max: 50, style: { width: '100%' } },
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Team> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '团队信息',
            key: 'team',
            width: 250,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size={48}
                        src={record.avatarUrl || undefined}
                        icon={<TeamOutlined />}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.name || '-'}</div>
                        {record.description && (
                            <Text type="secondary" style={{ fontSize: 12 }} ellipsis>
                                {record.description}
                            </Text>
                        )}
                    </div>
                </Space>
            ),
        },
        {
            title: '队长',
            key: 'leader',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Avatar
                        size={32}
                        src={record.leader?.avatar || undefined}
                        icon={<UserOutlined />}
                    />
                    <div>
                        <div>{record.leader?.nickname || `ID:${record.leaderId}`}</div>
                        {record.leader?.rank && (
                            <Tag color="blue" style={{ fontSize: 11 }}>{record.leader.rank}</Tag>
                        )}
                    </div>
                </Space>
            ),
        },
        {
            title: '成员/上限',
            key: 'members',
            width: 120,
            render: (_, record) => (
                <Space>
                    <TeamOutlined />
                    <span>{record.memberCount}/{record.maxMembers}</span>
                </Space>
            ),
        },
        {
            title: '收益分配',
            dataIndex: 'incomeShareType',
            key: 'incomeShareType',
            width: 120,
            render: (type: Team['incomeShareType'], record) => (
                <div>
                    <div>{shareTypeMap[type].text}</div>
                    {record.leaderBonusRate > 0 && (
                        <Text type="secondary" style={{ fontSize: 11 }}>
                            队长+{record.leaderBonusRate}%
                        </Text>
                    )}
                </div>
            ),
        },
        {
            title: '统计数据',
            key: 'stats',
            width: 150,
            render: (_, record) => (
                <div style={{ fontSize: 12 }}>
                    <div>订单: {record.totalOrderCount}</div>
                    <div>收益: ¥{(record.totalIncomeCents / 100).toFixed(2)}</div>
                </div>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: Team['status']) => (
                <Tag color={statusMap[status].color}>
                    {statusMap[status].text}
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
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定要删除该团队吗？"
                        onConfirm={() => handleDelete(record)}
                    >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量修改状态',
            icon: <CheckOutlined />,
            needSelection: false,
            onClick: (keys) => handleBatchStatus(keys || []),
        },
        {
            text: '批量删除',
            icon: <DeleteOutlined />,
            needSelection: false,
            danger: true,
            onClick: (keys) => handleBatchDelete(keys || []),
        },
    ];

    return (
        <PageContainer
            title="团队管理"
            subTitle="管理游戏陪玩团队"
            extra={
                stats && (
                    <Space size="large">
                        <Statistic title="总团队数" value={stats.totalTeams} />
                        <Statistic title="活跃团队" value={stats.activeTeams} valueStyle={{ color: '#52c41a' }} />
                        <Statistic title="总成员数" value={stats.totalMembers} />
                    </Space>
                )
            }
        >
            <SearchTable
                columns={columns}
                dataSource={teams}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => {
                    loadData();
                    loadStats();
                }}
                loading={loading}
                showCreate={true}
                createText="创建团队"
                onCreate={handleAdd}
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

            {/* 新增/编辑弹窗 */}
            <TeamForm
                open={formVisible}
                team={currentTeam}
                onOk={() => {
                    setFormVisible(false);
                    loadData();
                    loadStats();
                }}
                onCancel={() => setFormVisible(false)}
            />

            {/* 详情抽屉 */}
            <Drawer
                title="团队详情"
                open={detailDrawerVisible}
                onClose={() => setDetailDrawerVisible(false)}
                size="large"
                width={720}
            >
                {currentTeam && (
                    <>
                        {/* 团队基本信息 */}
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <div style={{ textAlign: 'center', marginBottom: 16 }}>
                                <Avatar
                                    size={80}
                                    src={currentTeam.avatarUrl || undefined}
                                    icon={<TeamOutlined />}
                                />
                                <Title level={4} style={{ marginTop: 12, marginBottom: 4 }}>
                                    {currentTeam.name}
                                </Title>
                                <Space>
                                    <Tag color={statusMap[currentTeam.status].color}>
                                        {statusMap[currentTeam.status].text}
                                    </Tag>
                                    <Tag>{shareTypeMap[currentTeam.incomeShareType].text}</Tag>
                                </Space>
                            </div>

                            <Row gutter={16}>
                                <Col span={6}>
                                    <Statistic title="成员数" value={`${currentTeam.memberCount}/${currentTeam.maxMembers}`} />
                                </Col>
                                <Col span={6}>
                                    <Statistic title="订单数" value={currentTeam.totalOrderCount} suffix="单" />
                                </Col>
                                <Col span={6}>
                                    <Statistic
                                        title="总收益"
                                        value={`¥${(currentTeam.totalIncomeCents / 100).toFixed(2)}`}
                                    />
                                </Col>
                                <Col span={6}>
                                    {currentTeam.leaderBonusRate > 0 && (
                                        <Statistic title="队长加成" value={`${currentTeam.leaderBonusRate}%`} />
                                    )}
                                </Col>
                            </Row>
                        </Card>

                        <Divider />

                        {/* 详细信息 */}
                        <Descriptions title="团队信息" column={2} size="small">
                            <Descriptions.Item label="团队ID">{currentTeam.id}</Descriptions.Item>
                            <Descriptions.Item label="当前状态">
                                <Tag color={statusMap[currentTeam.status].color}>
                                    {statusMap[currentTeam.status].text}
                                </Tag>
                            </Descriptions.Item>
                            <Descriptions.Item label="队长">
                                <Space>
                                    <Avatar
                                        size={24}
                                        src={currentTeam.leader?.avatar || undefined}
                                        icon={<UserOutlined />}
                                    />
                                    <span>{currentTeam.leader?.nickname || `ID:${currentTeam.leaderId}`}</span>
                                    {currentTeam.leader?.rank && (
                                        <Tag color="blue">{currentTeam.leader.rank}</Tag>
                                    )}
                                </Space>
                            </Descriptions.Item>
                            <Descriptions.Item label="成员上限">{currentTeam.maxMembers}人</Descriptions.Item>
                            <Descriptions.Item label="收益分配" span={2}>
                                {shareTypeMap[currentTeam.incomeShareType].text}
                                {currentTeam.leaderBonusRate > 0 && `（队长额外+${currentTeam.leaderBonusRate}%）`}
                            </Descriptions.Item>
                            <Descriptions.Item label="团队描述" span={2}>
                                <Text ellipsis>{currentTeam.description || '暂无描述'}</Text>
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间">
                                {dayjs(currentTeam.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            <Descriptions.Item label="更新时间">
                                {dayjs(currentTeam.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                        </Descriptions>

                        <Divider />

                        {/* 成员列表 */}
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                            <Title level={5} style={{ margin: 0 }}>团队成员</Title>
                            <Button
                                type="primary"
                                size="small"
                                icon={<PlusOutlined />}
                                onClick={handleAddMember}
                            >
                                添加成员
                            </Button>
                        </div>

                        <div style={{ maxHeight: 400, overflowY: 'auto' }}>
                            {currentTeam.members && currentTeam.members.length > 0 ? (
                                currentTeam.members.map(member => (
                                    <MemberCard
                                        key={member.id}
                                        member={member}
                                        onRemove={handleRemoveMember}
                                        onTransfer={handleTransferLeader}
                                    />
                                ))
                            ) : (
                                <div style={{ textAlign: 'center', padding: 24, color: '#999' }}>
                                    暂无成员
                                </div>
                            )}
                        </div>
                    </>
                )}
            </Drawer>

            {/* 批量修改状态弹窗 */}
            <Modal
                title="批量修改团队状态"
                open={batchStatusVisible}
                onOk={submitBatchStatus}
                onCancel={() => setBatchStatusVisible(false)}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedTeamIds.length === 0}>
                                选中的团队 {selectedTeamIds.length > 0 ? `(${selectedTeamIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部团队</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true }]}>
                            <Select placeholder="请选择筛选状态">
                                <Select.Option value="active">活跃</Select.Option>
                                <Select.Option value="busy">接单中</Select.Option>
                                <Select.Option value="inactive">不活跃</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <Form.Item name="status" label="修改为" rules={[{ required: true }]}>
                        <Select placeholder="请选择目标状态">
                            <Select.Option value="active">活跃</Select.Option>
                            <Select.Option value="busy">接单中</Select.Option>
                            <Select.Option value="inactive">不活跃</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 批量删除弹窗 */}
            <Modal
                title="批量删除团队"
                open={batchDeleteVisible}
                onOk={submitBatchDelete}
                onCancel={() => setBatchDeleteVisible(false)}
                okText="确认删除"
                okButtonProps={{ danger: true }}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group onChange={(e) => setBatchTarget(e.target.value)}>
                            <Radio value="selected" disabled={selectedTeamIds.length === 0}>
                                选中的团队 {selectedTeamIds.length > 0 ? `(${selectedTeamIds.length})` : ''}
                            </Radio>
                            <Radio value="status">按状态筛选</Radio>
                            <Radio value="all">全部团队</Radio>
                        </Radio.Group>
                    </Form.Item>

                    {batchTarget === 'status' && (
                        <Form.Item name="filterStatus" label="筛选状态" rules={[{ required: true }]}>
                            <Select placeholder="请选择筛选状态">
                                <Select.Option value="active">活跃</Select.Option>
                                <Select.Option value="busy">接单中</Select.Option>
                                <Select.Option value="inactive">不活跃</Select.Option>
                            </Select>
                        </Form.Item>
                    )}

                    <div style={{ color: '#ff4d4f', marginTop: 16 }}>
                        警告：此操作不可恢复，请谨慎操作！
                    </div>
                </Form>
            </Modal>

            {/* 添加成员弹窗 */}
            <Modal
                title="添加团队成员"
                open={addMemberVisible}
                onOk={submitAddMember}
                onCancel={() => setAddMemberVisible(false)}
            >
                <Form form={memberForm} layout="vertical">
                    <Form.Item
                        name="playerId"
                        label="陪玩师ID"
                        rules={[
                            { required: true, message: '请输入陪玩师ID' },
                            { type: 'number', min: 1, message: '请输入有效的陪玩师ID' },
                        ]}
                    >
                        <InputNumber
                            placeholder="请输入要添加的陪玩师ID"
                            min={1}
                            style={{ width: '100%' }}
                        />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 转让队长弹窗 */}
            <Modal
                title="转让队长"
                open={transferVisible}
                onOk={submitTransferLeader}
                onCancel={() => setTransferVisible(false)}
            >
                {currentTeam && currentTeam.pendingTransfer && (
                    <div>
                        <p>确定要将队长转让给以下成员吗？</p>
                        <Card size="small" style={{ marginBottom: 16 }}>
                            <Space>
                                <Avatar
                                    size={40}
                                    src={currentTeam.pendingTransfer.player?.avatar || undefined}
                                    icon={<UserOutlined />}
                                />
                                <div>
                                    <div>{currentTeam.pendingTransfer.player?.nickname || `ID:${currentTeam.pendingTransfer.playerId}`}</div>
                                    <Text type="secondary" style={{ fontSize: 12 }}>
                                        当前角色：{currentTeam.pendingTransfer.role === 'leader' ? '队长' : '成员'}
                                    </Text>
                                </div>
                            </Space>
                        </Card>
                        <p style={{ color: '#ff4d4f' }}>
                            转让后您将成为普通成员，此操作不可撤销！
                        </p>
                    </div>
                )}
            </Modal>
        </PageContainer>
    );
};

export default TeamPage;
