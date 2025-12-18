import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Avatar,
    Modal,
    Form,
    Input,
    Select,
    message,
    Popconfirm,
    Drawer,
    Descriptions,
    Divider,
    theme,
    Row,
    Col,
    Card,
    Statistic,
    Tabs,
    Radio
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    UserOutlined,
    EditOutlined,
    DeleteOutlined,
    LockOutlined,
    UnlockOutlined,
    EyeOutlined,
    TeamOutlined,
    CrownOutlined,
    SafetyOutlined,
    MailOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { exportToCSV, userExportColumns } from '@/utils/export';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { USER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type User, type CreateUserDto, type UpdateUserDto, type UserQueryParams, type UserStats, type ApiResponse, type LoginHistory, type AuditLog } from '@/api/admin';
import { Table } from 'antd';
import dayjs from 'dayjs';

/**
 * 角色映射
 */
const roleMap = {
    user: { color: 'blue', text: '普通用户' },
    player: { color: 'purple', text: '陪玩师' },
    admin: { color: 'red', text: '管理员' },
};

/**
 * 状态映射
 */
const statusMap = {
    active: { color: 'success', text: '正常' },
    banned: { color: 'error', text: '已封禁' },
    suspended: { color: 'warning', text: '已停用' },
};

/**
 * 用户管理页面
 */
const UserPage: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [users, setUsers] = useState<User[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 统计数据
    const [stats, setStats] = useState<UserStats | null>(null);
    const [statsLoading, setStatsLoading] = useState(false);

    // 搜索参数
    const [searchParams, setSearchParams] = useState<UserQueryParams>({});

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [form] = Form.useForm();

    // Batch Operation States
    const [batchRoleVisible, setBatchRoleVisible] = useState(false);
    const [batchStatusVisible, setBatchStatusVisible] = useState(false);
    const [batchNotificationVisible, setBatchNotificationVisible] = useState(false);
    const [batchPointsVisible, setBatchPointsVisible] = useState(false);
    const [selectedUserIds, setSelectedUserIds] = useState<number[]>([]);
    const [batchForm] = Form.useForm();
    const [notificationForm] = Form.useForm();
    const [statusForm] = Form.useForm();
    const [pointsForm] = Form.useForm();

    // Login History State
    const [loginHistory, setLoginHistory] = useState<LoginHistory[]>([]);
    const [loginHistoryLoading, setLoginHistoryLoading] = useState(false);

    // Operation Logs State
    const [operationLogs, setOperationLogs] = useState<AuditLog[]>([]);
    const [operationLogsLoading, setOperationLogsLoading] = useState(false);

    const fetchLoginHistory = async (userId: number) => {
        setLoginHistoryLoading(true);
        try {
            const res = await adminApi.getUserLogs(userId, { page: 1, page_size: 10, type: 'login' }) as unknown as ApiResponse<AuditLog[]>;
            if (res.success) {
                const history: LoginHistory[] = (res.data || []).map(log => ({
                    id: log.id,
                    ip: log.ip,
                    location: log.location,
                    device: log.device,
                    loginAt: log.createdAt,
                    status: log.action.includes('fail') ? 'failed' : 'success'
                }));
                setLoginHistory(history);
            }
        } catch (err) {
            console.error('Failed to fetch login history', err);
        } finally {
            setLoginHistoryLoading(false);
        }
    };

    const fetchOperationLogs = async (userId: number) => {
        setOperationLogsLoading(true);
        try {
            const res = await adminApi.getUserLogs(userId, { page: 1, page_size: 10 }) as unknown as ApiResponse<AuditLog[]>;
            if (res.success) {
                setOperationLogs(res.data || []);
            }
        } catch (err) {
            console.error('Failed to fetch operation logs', err);
        } finally {
            setOperationLogsLoading(false);
        }
    };

    const handleTabChange = (key: string) => {
        if (!currentUser) return;
        if (key === '2') {
            fetchLoginHistory(currentUser.id);
        } else if (key === '3') {
            fetchOperationLogs(currentUser.id);
        }
    };

    /**
     * 加载统计数据
     */
    const loadStats = useCallback(async () => {
        try {
            setStatsLoading(true);
            const response = await adminApi.getUserStats();
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (err) {
            console.error('加载统计数据失败:', err);
            // 静默失败，不影响主要功能
        } finally {
            setStatsLoading(false);
        }
    }, []);

    /**
     * 加载用户数据
     */
    const loadData = useCallback(async () => {
        try {
            setLoading(true);
            const params: UserQueryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
            };

            const response = await adminApi.getUsers(params);

            if (response.data.success) {
                setUsers(response.data.data || []);
                setTotal((response.data as { pagination?: { total?: number } }).pagination?.total || 0);
            } else {
                message.error(response.data.message || '获取用户列表失败');
            }
        } catch (error: unknown) {
            console.error('加载用户列表失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '获取用户列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
        loadStats();
    }, [loadData, loadStats]);

    /**
     * 搜索处理
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        const params: UserQueryParams = {};

        const keyword = values.keyword as string | undefined;
        if (keyword?.trim()) {
            params.keyword = keyword.trim();
        }
        if (values.role) {
            const roles = Array.isArray(values.role) ? values.role : [values.role];
            if (roles.length > 0) {
                params.role = roles as string[];
            }
        }
        if (values.status) {
            const statuses = Array.isArray(values.status) ? values.status : [values.status];
            if (statuses.length > 0) {
                params.status = statuses as string[];
            }
        }
        const dateRange = values.dateRange as { format: (f: string) => string }[] | undefined;
        if (dateRange?.length === 2) {
            params.date_from = dateRange[0].format('YYYY-MM-DD');
            params.date_to = dateRange[1].format('YYYY-MM-DD');
        }

        setSearchParams(params);
        setCurrent(1); // 重置到第一页
    }, []);

    /**
     * 编辑用户
     */
    const handleEdit = useCallback((user: User) => {
        setCurrentUser(user);
        form.setFieldsValue({
            name: user.name,
            email: user.email,
            phone: user.phone,
            avatarUrl: user.avatarUrl,
            role: user.role,
            status: user.status,
        });
        setEditModalVisible(true);
    }, [form]);

    /**
     * 查看详情
     */
    const handleViewDetail = useCallback((user: User) => {
        setCurrentUser(user);
        setDetailDrawerVisible(true);
    }, []);

    /**
     * 保存编辑
     */
    const handleSaveEdit = useCallback(async () => {
        try {
            const values = await form.validateFields();

            if (currentUser) {
                // 更新用户
                const updateData: UpdateUserDto = {
                    name: values.name,
                    email: values.email,
                    phone: values.phone,
                    avatarUrl: values.avatarUrl,
                    role: values.role,
                    status: values.status,
                };

                // 如果输入了新密码,则包含密码
                if (values.password && values.password.trim()) {
                    updateData.password = values.password.trim();
                }

                const response = await adminApi.updateUser(currentUser.id, updateData);

                if (response.data.success) {
                    message.success('更新用户成功');
                    setEditModalVisible(false);
                    loadData();
                } else {
                    message.error(response.data.message || '更新用户失败');
                }
            } else {
                // 创建用户
                const createData: CreateUserDto = {
                    name: values.name,
                    email: values.email,
                    phone: values.phone,
                    password: values.password,
                    avatarUrl: values.avatarUrl,
                    role: values.role,
                    status: values.status,
                };

                const response = await adminApi.createUser(createData);

                if (response.data.success) {
                    message.success('创建用户成功');
                    setEditModalVisible(false);
                    loadData();
                } else {
                    message.error(response.data.message || '创建用户失败');
                }
            }
        } catch (error: unknown) {
            console.error('保存用户失败:', error);
            const err = error as { errorFields?: unknown; response?: { data?: { message?: string } } };
            if (err.errorFields) {
                // 表单验证失败
                return;
            }
            message.error(err.response?.data?.message || '保存用户失败');
        }
    }, [form, currentUser, loadData]);

    /**
     * 封禁/解封用户
     */
    const handleToggleBan = useCallback(async (user: User) => {
        try {
            const newStatus = user.status === 'banned' ? 'active' : 'banned';
            const action = newStatus === 'banned' ? '封禁' : '解封';

            const response = await adminApi.updateUserStatus(user.id, newStatus);

            if (response.data.success) {
                message.success(`${action}用户成功`);
                loadData();
            } else {
                message.error(response.data.message || `${action}用户失败`);
            }
        } catch (error: unknown) {
            console.error('更新用户状态失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '更新用户状态失败');
        }
    }, [loadData]);

    /**
     * 删除用户
     */
    const handleDelete = useCallback(async (user: User) => {
        try {
            const response = await adminApi.deleteUser(user.id);

            if (response.data.success) {
                message.success(`删除用户 ${user.name} 成功`);
                loadData();
            } else {
                message.error(response.data.message || '删除用户失败');
            }
        } catch (error: unknown) {
            console.error('删除用户失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '删除用户失败');
        }
    }, [loadData]);

    /**
     * 批量删除
     */
    const handleBatchDelete = useCallback(async (selectedRowKeys: React.Key[]) => {
        try {
            const ids = selectedRowKeys.map(key => Number(key));
            const response = await adminApi.batchDeleteUsers(ids);

            if (response.data.success) {
                message.success(`成功删除 ${ids.length} 个用户`);
                loadData();
            } else {
                message.error(response.data.message || '批量删除失败');
            }
        } catch (error: unknown) {
            console.error('批量删除用户失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '批量删除失败');
        }
    }, [loadData]);

    const handleBatchModifyRole = (keys: React.Key[]) => {
        if (!keys || keys.length === 0) return;
        setSelectedUserIds(keys.map(k => Number(k)));
        batchForm.resetFields();
        setBatchRoleVisible(true);
    };

    const handleBatchModifyStatus = (keys: React.Key[]) => {
        if (!keys || keys.length === 0) return;
        setSelectedUserIds(keys.map(k => Number(k)));
        statusForm.resetFields();
        setBatchStatusVisible(true);
    };

    const handleBatchSendNotification = (keys: React.Key[]) => {
        setSelectedUserIds(keys ? keys.map(k => Number(k)) : []);
        notificationForm.resetFields();
        notificationForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'users' : 'all',
            type: 'system'
        });
        setBatchNotificationVisible(true);
    };

    const submitBatchRole = async () => {
        try {
            const values = await batchForm.validateFields();
            const res = await adminApi.batchUpdateUserRole({
                userIds: selectedUserIds,
                role: values.role
            }) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success('批量修改角色成功');
                setBatchRoleVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    const submitBatchStatus = async () => {
        try {
            const values = await statusForm.validateFields();
            const res = await adminApi.batchUpdateUserStatus({
                userIds: selectedUserIds,
                status: values.status
            }) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success('批量修改状态成功');
                setBatchStatusVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    const submitBatchNotification = async () => {
        try {
            const values = await notificationForm.validateFields();
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const payload: any = {
                target: values.target,
                title: values.title,
                content: values.content,
                type: values.type
            };

            if (values.target === 'users') {
                payload.userIds = selectedUserIds;
            } else if (values.target === 'role') {
                payload.roles = values.roles;
            }

            const res = await adminApi.batchSendNotification(payload) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success('批量发送通知成功');
                setBatchNotificationVisible(false);
            }
        } catch {
            message.error('操作失败');
        }
    };

    const handleBatchAddPoints = (keys: React.Key[]) => {
        setSelectedUserIds(keys ? keys.map(k => Number(k)) : []);
        pointsForm.resetFields();
        pointsForm.setFieldsValue({
            target: (keys && keys.length > 0) ? 'users' : 'all',
            type: 'activity'
        });
        setBatchPointsVisible(true);
    };

    const submitBatchPoints = async () => {
        try {
            const values = await pointsForm.validateFields();
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const payload: any = {
                target: values.target,
                cents: Number(values.points),
                reason: values.reason,
                type: values.type
            };

            if (values.target === 'users') {
                payload.userIds = selectedUserIds;
            } else if (values.target === 'role') {
                payload.roles = values.roles;
            }

            const res = await adminApi.batchAddUserPoints(payload) as unknown as ApiResponse<void>;

            if (res.success) {
                message.success('批量增加积分成功');
                setBatchPointsVisible(false);
                loadData();
            }
        } catch {
            message.error('操作失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = useMemo(() => [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '用户名/邮箱/手机号' },
        {
            name: 'role',
            label: '角色',
            type: 'select',
            mode: 'multiple',
            options: [
                { label: '普通用户', value: 'user' },
                { label: '陪玩师', value: 'player' },
                { label: '管理员', value: 'admin' },
            ],
        },
        {
            name: 'status',
            label: '状态',
            type: 'select',
            mode: 'multiple',
            options: [
                { label: '正常', value: 'active' },
                { label: '已封禁', value: 'banned' },
                { label: '已停用', value: 'suspended' },
            ],
        },
        { name: 'dateRange', label: '注册时间', type: 'dateRange' },
    ], []);

    /**
     * 表格列配置
     */
    const columns: ColumnsType<User> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 60,
        },
        {
            title: '用户信息',
            key: 'userInfo',
            width: 180,
            ellipsis: true,
            render: (_, record) => (
                <Space>
                    <Avatar
                        src={record.avatarUrl || undefined}
                        icon={<UserOutlined />}
                        style={{ backgroundColor: token.colorPrimary }}
                    />
                    <div style={{ overflow: 'hidden' }}>
                        <div style={{ fontWeight: 500, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{record.name}</div>
                        <div style={{ fontSize: 12, color: token.colorTextSecondary, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{record.email}</div>
                    </div>
                </Space>
            ),
        },
        {
            title: '手机号',
            dataIndex: 'phone',
            key: 'phone',
            width: 115,
        },
        {
            title: '角色',
            dataIndex: 'role',
            key: 'role',
            width: 80,
            render: (role: User['role']) => <Tag color={roleMap[role].color}>{roleMap[role].text}</Tag>,
        },
        {
            title: '等级',
            dataIndex: 'level',
            key: 'level',
            width: 60,
            render: (level: number) => <Tag color="gold">Lv.{level || 0}</Tag>,
        },
        {
            title: '积分',
            dataIndex: ['wallet', 'balanceCents'],
            key: 'points',
            width: 70,
            render: (_: unknown, record: User) => {
                const points = (record as User & { wallet?: { balanceCents?: number } }).wallet?.balanceCents ?? 0;
                return <span style={{ fontWeight: 500, color: token.colorPrimary }}>{points}</span>;
            },
        },
        {
            title: '标签',
            dataIndex: 'tags',
            key: 'tags',
            width: 100,
            ellipsis: true,
            render: (tags: string[]) => (
                <Space size={[0, 4]} wrap>
                    {tags?.slice(0, 2).map(tag => <Tag key={tag}>{tag}</Tag>)}
                    {tags?.length > 2 && <Tag>+{tags.length - 2}</Tag>}
                </Space>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 70,
            render: (status: User['status']) => <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>,
        },
        {
            title: '注册时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 145,
            render: (createdAt: string) => dayjs(createdAt).format('YYYY-MM-DD HH:mm'),
        },
        {
            title: '最后登录',
            dataIndex: 'lastLoginAt',
            key: 'lastLoginAt',
            width: 145,
            render: (lastLoginAt?: string) => lastLoginAt ? dayjs(lastLoginAt).format('YYYY-MM-DD HH:mm') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 240,
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
                    <PermissionGuard permission={USER_PERMISSIONS.UPDATE}>
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        >
                            编辑
                        </Button>
                    </PermissionGuard>
                    <PermissionGuard permission={USER_PERMISSIONS.STATUS}>
                        <Popconfirm
                            title={`确定要${record.status === 'banned' ? '解封' : '封禁'}该用户吗？`}
                            onConfirm={() => handleToggleBan(record)}
                        >
                            <Button
                                type="link"
                                size="small"
                                danger={record.status !== 'banned'}
                                icon={record.status === 'banned' ? <UnlockOutlined /> : <LockOutlined />}
                            >
                                {record.status === 'banned' ? '解封' : '封禁'}
                            </Button>
                        </Popconfirm>
                    </PermissionGuard>
                    <PermissionGuard permission={USER_PERMISSIONS.DELETE}>
                        <Popconfirm
                            title="确定要删除该用户吗?此操作不可恢复。"
                            onConfirm={() => handleDelete(record)}
                        >
                            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                                删除
                            </Button>
                        </Popconfirm>
                    </PermissionGuard>
                </Space>
            ),
        },
    ], [token, handleViewDetail, handleEdit, handleToggleBan, handleDelete]);

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = [
        {
            text: '批量修改角色',
            icon: <EditOutlined />,
            needSelection: true,
            onClick: (keys) => handleBatchModifyRole(keys || []),
            permission: USER_PERMISSIONS.UPDATE,
        },
        {
            text: '批量修改状态',
            icon: <SafetyOutlined />,
            needSelection: true,
            onClick: (keys) => handleBatchModifyStatus(keys || []),
            permission: USER_PERMISSIONS.UPDATE,
        },
        {
            text: '批量发送通知',
            icon: <MailOutlined />,
            needSelection: false,
            onClick: (keys) => handleBatchSendNotification(keys || []),
            permission: USER_PERMISSIONS.UPDATE,
        },
        {
            text: '批量增加积分',
            icon: <CrownOutlined />,
            needSelection: false,
            onClick: (keys) => handleBatchAddPoints(keys || []),
            permission: USER_PERMISSIONS.UPDATE,
        },
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
            permission: USER_PERMISSIONS.LIST,
        },
    ];

    /**
     * 导出用户数据
     */
    const handleExport = useCallback(async () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            // 获取所有符合条件的数据
            const response = await adminApi.getUsers({ ...searchParams, page_size: 10000 });
            if (response.data.success && response.data.data) {
                exportToCSV(response.data.data as unknown as Record<string, unknown>[], userExportColumns, 'users');
                message.success({ content: '导出成功', key: 'export' });
            } else {
                message.error({ content: '导出失败', key: 'export' });
            }
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    }, [searchParams]);

    return (
        <PageContainer title="用户管理" subTitle="管理平台所有注册用户">
            {/* 统计卡片 */}
            <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading}>
                        <Statistic
                            title="用户总数"
                            value={stats?.total || 0}
                            prefix={<TeamOutlined />}
                            styles={{ content: { color: token.colorPrimary } }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading}>
                        <Statistic
                            title="陪玩师"
                            value={stats?.byRole.player || 0}
                            prefix={<CrownOutlined />}
                            styles={{ content: { color: '#722ed1' } }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading}>
                        <Statistic
                            title="正常用户"
                            value={stats?.byStatus.active || 0}
                            prefix={<SafetyOutlined />}
                            styles={{ content: { color: '#52c41a' } }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading}>
                        <Statistic
                            title="最近注册"
                            value={stats?.recentRegistrations || 0}
                            prefix={<UserOutlined />}
                            styles={{ content: { color: '#52c41a' } }}
                        />
                    </Card>
                </Col>
            </Row>

            <SearchTable
                columns={columns}
                dataSource={users}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={loadData}
                loading={loading}
                showCreate={true}
                createText="新增用户"
                createPermission={USER_PERMISSIONS.CREATE}
                onCreate={() => {
                    setCurrentUser(null);
                    form.resetFields();
                    setEditModalVisible(true);
                }}
                showBatchDelete={true}
                batchDeletePermission={USER_PERMISSIONS.DELETE}
                onBatchDelete={handleBatchDelete}
                toolbarButtons={toolbarButtons}
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
                scroll={{ x: 'max-content' }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentUser ? '编辑用户' : '新增用户'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={600}
                style={{ maxWidth: 'calc(100vw - 32px)', top: 20 }}
                okText="保存"
                cancelText="取消"
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="用户名"
                        rules={[{ required: true, message: '请输入用户名' }]}
                    >
                        <Input placeholder="请输入用户名" />
                    </Form.Item>
                    <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                            { required: true, message: '请输入邮箱' },
                            { type: 'email', message: '请输入正确的邮箱格式' },
                        ]}
                    >
                        <Input placeholder="请输入邮箱" />
                    </Form.Item>
                    <Form.Item
                        name="phone"
                        label="手机号"
                        rules={[
                            { required: true, message: '请输入手机号' },
                            { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' },
                        ]}
                    >
                        <Input placeholder="请输入手机号" />
                    </Form.Item>
                    <Form.Item
                        name="password"
                        label={currentUser ? '新密码(留空则不修改)' : '密码'}
                        rules={currentUser ? [] : [
                            { required: true, message: '请输入密码' },
                            { min: 6, message: '密码至少6位' },
                        ]}
                    >
                        <Input.Password placeholder={currentUser ? '留空则不修改密码' : '请输入密码(至少6位)'} />
                    </Form.Item>
                    <Form.Item name="avatarUrl" label="头像URL">
                        <Input placeholder="请输入头像URL(可选)" />
                    </Form.Item>
                    <Form.Item name="role" label="角色" rules={[{ required: true, message: '请选择角色' }]}>
                        <Select
                            placeholder="请选择角色"
                            options={[
                                { label: '普通用户', value: 'user' },
                                { label: '陪玩师', value: 'player' },
                                { label: '管理员', value: 'admin' },
                            ]}
                        />
                    </Form.Item>
                    <Form.Item name="status" label="状态" rules={[{ required: true, message: '请选择状态' }]}>
                        <Select
                            placeholder="请选择状态"
                            options={[
                                { label: '正常', value: 'active' },
                                { label: '已封禁', value: 'banned' },
                                { label: '已停用', value: 'suspended' },
                            ]}
                        />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 详情抽屉 */}
            <Drawer
                title="用户详情"
                open={detailDrawerVisible}
                onClose={() => setDetailDrawerVisible(false)}
                size="large"
                style={{ maxWidth: '100%' }}
            >
                {currentUser && (
                    <Tabs defaultActiveKey="1" onChange={handleTabChange}>
                        <Tabs.TabPane tab="基本信息" key="1">
                            <div style={{ textAlign: 'center', marginBottom: 24 }}>
                                <Avatar
                                    size={80}
                                    src={currentUser.avatarUrl || undefined}
                                    icon={<UserOutlined />}
                                    style={{ backgroundColor: token.colorPrimary }}
                                />
                                <h2 style={{ marginTop: 16, marginBottom: 4 }}>{currentUser.name}</h2>
                                <Tag color={roleMap[currentUser.role].color}>{roleMap[currentUser.role].text}</Tag>
                                <Tag color={statusMap[currentUser.status].color}>{statusMap[currentUser.status].text}</Tag>
                            </div>

                            <Divider />

                            <Descriptions column={2} labelStyle={{ width: 100 }}>
                                <Descriptions.Item label="用户ID">{currentUser.id}</Descriptions.Item>
                                <Descriptions.Item label="邮箱">{currentUser.email}</Descriptions.Item>
                                <Descriptions.Item label="手机号">{currentUser.phone}</Descriptions.Item>
                                <Descriptions.Item label="等级">Lv.{currentUser.level || 0}</Descriptions.Item>
                                <Descriptions.Item label="标签">
                                    {currentUser.tags?.map(tag => <Tag key={tag}>{tag}</Tag>)}
                                </Descriptions.Item>
                                <Descriptions.Item label="注册时间">
                                    {dayjs(currentUser.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                                </Descriptions.Item>
                                <Descriptions.Item label="最后登录">
                                    {currentUser.lastLoginAt ? dayjs(currentUser.lastLoginAt).format('YYYY-MM-DD HH:mm:ss') : '从未登录'}
                                </Descriptions.Item>
                            </Descriptions>
                        </Tabs.TabPane>
                        <Tabs.TabPane tab="登录历史" key="2">
                            <Table
                                dataSource={loginHistory}
                                loading={loginHistoryLoading}
                                rowKey="id"
                                pagination={false}
                                columns={[
                                    { title: '时间', dataIndex: 'loginAt', render: (t) => dayjs(t).format('YYYY-MM-DD HH:mm:ss') },
                                    { title: 'IP', dataIndex: 'ip' },
                                    { title: '地点', dataIndex: 'location' },
                                    { title: '设备', dataIndex: 'device' },
                                    { title: '状态', dataIndex: 'status', render: (s) => <Tag color={s === 'success' ? 'success' : 'error'}>{s === 'success' ? '成功' : '失败'}</Tag> },
                                ]}
                            />
                        </Tabs.TabPane>
                        <Tabs.TabPane tab="操作日志" key="3">
                            <Table
                                dataSource={operationLogs}
                                loading={operationLogsLoading}
                                rowKey="id"
                                pagination={false}
                                columns={[
                                    { title: '时间', dataIndex: 'createdAt', render: (t) => dayjs(t).format('YYYY-MM-DD HH:mm:ss') },
                                    { title: '操作', dataIndex: 'action' },
                                    { title: '详情', dataIndex: 'details' },
                                    { title: 'IP', dataIndex: 'ip' },
                                ]}
                            />
                        </Tabs.TabPane>
                    </Tabs>
                )}
            </Drawer>

            {/* Batch Role Modal */}
            <Modal
                title="批量修改角色"
                open={batchRoleVisible}
                onOk={submitBatchRole}
                onCancel={() => setBatchRoleVisible(false)}
            >
                <Form form={batchForm} layout="vertical">
                    <Form.Item name="role" label="选择角色" rules={[{ required: true }]}>
                        <Select>
                            <Select.Option value="user">普通用户</Select.Option>
                            <Select.Option value="player">陪玩师</Select.Option>
                            <Select.Option value="admin">管理员</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            {/* Batch Status Modal */}
            <Modal
                title="批量修改状态"
                open={batchStatusVisible}
                onOk={submitBatchStatus}
                onCancel={() => setBatchStatusVisible(false)}
            >
                <Form form={statusForm} layout="vertical">
                    <Form.Item name="status" label="选择状态" rules={[{ required: true }]}>
                        <Select>
                            <Select.Option value="active">正常</Select.Option>
                            <Select.Option value="banned">封禁</Select.Option>
                            <Select.Option value="suspended">停用</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            {/* Batch Notification Modal */}
            <Modal
                title="批量发送通知"
                open={batchNotificationVisible}
                onOk={submitBatchNotification}
                onCancel={() => setBatchNotificationVisible(false)}
            >
                <Form form={notificationForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group>
                            <Radio value="users" disabled={selectedUserIds.length === 0}>
                                指定用户 {selectedUserIds.length > 0 ? `(${selectedUserIds.length})` : ''}
                            </Radio>
                            <Radio value="role">指定角色</Radio>
                            <Radio value="all">全体用户</Radio>
                        </Radio.Group>
                    </Form.Item>

                    <Form.Item
                        noStyle
                        shouldUpdate={(prev, current) => prev.target !== current.target}
                    >
                        {({ getFieldValue }) =>
                            getFieldValue('target') === 'role' ? (
                                <Form.Item name="roles" label="选择角色" rules={[{ required: true, message: '请选择角色' }]}>
                                    <Select mode="multiple" placeholder="请选择角色">
                                        <Select.Option value="user">普通用户</Select.Option>
                                        <Select.Option value="player">陪玩师</Select.Option>
                                        <Select.Option value="admin">管理员</Select.Option>
                                    </Select>
                                </Form.Item>
                            ) : null
                        }
                    </Form.Item>

                    <Form.Item name="title" label="标题" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="type" label="类型" initialValue="system">
                        <Radio.Group>
                            <Radio value="system">系统通知</Radio>
                            <Radio value="marketing">营销通知</Radio>
                            <Radio value="personal">个人通知</Radio>
                            <Radio value="activity">活动通知</Radio>
                        </Radio.Group>
                    </Form.Item>
                    <Form.Item name="content" label="内容" rules={[{ required: true }]}>
                        <Input.TextArea rows={4} />
                    </Form.Item>
                </Form>
            </Modal>

            {/* Batch Points Modal */}
            <Modal
                title="批量增加积分"
                open={batchPointsVisible}
                onOk={submitBatchPoints}
                onCancel={() => setBatchPointsVisible(false)}
            >
                <Form form={pointsForm} layout="vertical">
                    <Form.Item name="target" label="目标对象" rules={[{ required: true }]}>
                        <Radio.Group>
                            <Radio value="users" disabled={selectedUserIds.length === 0}>
                                指定用户 {selectedUserIds.length > 0 ? `(${selectedUserIds.length})` : ''}
                            </Radio>
                            <Radio value="role">指定角色</Radio>
                            <Radio value="all">全体用户</Radio>
                        </Radio.Group>
                    </Form.Item>

                    <Form.Item
                        noStyle
                        shouldUpdate={(prev, current) => prev.target !== current.target}
                    >
                        {({ getFieldValue }) =>
                            getFieldValue('target') === 'role' ? (
                                <Form.Item name="roles" label="选择角色" rules={[{ required: true, message: '请选择角色' }]}>
                                    <Select mode="multiple" placeholder="请选择角色">
                                        <Select.Option value="user">普通用户</Select.Option>
                                        <Select.Option value="player">陪玩师</Select.Option>
                                        <Select.Option value="admin">管理员</Select.Option>
                                    </Select>
                                </Form.Item>
                            ) : null
                        }
                    </Form.Item>

                    <Form.Item name="points" label="积分数量 (分)" rules={[{ required: true, message: '请输入积分数量' }]}>
                        <Input type="number" placeholder="请输入积分数量 (100分 = 1元)" />
                    </Form.Item>
                    <Form.Item name="type" label="积分类型" rules={[{ required: true, message: '请选择积分类型' }]}>
                        <Select placeholder="请选择积分类型">
                            <Select.Option value="admin">管理赠送</Select.Option>
                            <Select.Option value="activity">活动奖励</Select.Option>
                            <Select.Option value="compensation">补偿</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="reason" label="变动原因" rules={[{ required: true, message: '请输入变动原因' }]}>
                        <Input placeholder="请输入变动原因" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default UserPage;
