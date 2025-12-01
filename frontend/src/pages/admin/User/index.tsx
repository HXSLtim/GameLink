/**
 * 用户管理页面
 */
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
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    UserOutlined,
    EditOutlined,
    DeleteOutlined,
    LockOutlined,
    UnlockOutlined,
    EyeOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { USER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type User, type CreateUserDto, type UpdateUserDto, type UserQueryParams } from '@/api/admin';
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

    // 搜索参数
    const [searchParams, setSearchParams] = useState<UserQueryParams>({});

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [form] = Form.useForm();

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
                setTotal(response.data.pagination?.total || 0);
            } else {
                message.error(response.data.message || '获取用户列表失败');
            }
        } catch (error: any) {
            console.error('加载用户列表失败:', error);
            message.error(error.response?.data?.message || '获取用户列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索处理
     */
    const handleSearch = useCallback((values: any) => {
        const params: UserQueryParams = {};

        if (values.keyword) {
            params.keyword = values.keyword;
        }
        if (values.role) {
            params.role = Array.isArray(values.role) ? values.role : [values.role];
        }
        if (values.status) {
            params.status = Array.isArray(values.status) ? values.status : [values.status];
        }
        if (values.dateRange && values.dateRange.length === 2) {
            params.date_from = values.dateRange[0].format('YYYY-MM-DD');
            params.date_to = values.dateRange[1].format('YYYY-MM-DD');
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
        } catch (error: any) {
            console.error('保存用户失败:', error);
            if (error.errorFields) {
                // 表单验证失败
                return;
            }
            message.error(error.response?.data?.message || '保存用户失败');
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
        } catch (error: any) {
            console.error('更新用户状态失败:', error);
            message.error(error.response?.data?.message || '更新用户状态失败');
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
        } catch (error: any) {
            console.error('删除用户失败:', error);
            message.error(error.response?.data?.message || '删除用户失败');
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
        } catch (error: any) {
            console.error('批量删除用户失败:', error);
            message.error(error.response?.data?.message || '批量删除失败');
        }
    }, [loadData]);

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = useMemo(() => [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '用户名/邮箱/手机号' },
        {
            name: 'role',
            label: '角色',
            type: 'select',
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
            width: 80,
        },
        {
            title: '用户信息',
            key: 'userInfo',
            width: 250,
            render: (_, record) => (
                <Space>
                    <Avatar
                        src={record.avatarUrl}
                        icon={<UserOutlined />}
                        style={{ backgroundColor: token.colorPrimary }}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.name}</div>
                        <div style={{ fontSize: 12, color: token.colorTextSecondary }}>{record.email}</div>
                    </div>
                </Space>
            ),
        },
        {
            title: '手机号',
            dataIndex: 'phone',
            key: 'phone',
            width: 130,
        },
        {
            title: '角色',
            dataIndex: 'role',
            key: 'role',
            width: 100,
            render: (role: User['role']) => <Tag color={roleMap[role].color}>{roleMap[role].text}</Tag>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: User['status']) => <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>,
        },
        {
            title: '注册时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (createdAt: string) => dayjs(createdAt).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '最后登录',
            dataIndex: 'lastLoginAt',
            key: 'lastLoginAt',
            width: 180,
            render: (lastLoginAt?: string) => lastLoginAt ? dayjs(lastLoginAt).format('YYYY-MM-DD HH:mm:ss') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 250,
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

    return (
        <PageContainer title="用户管理" subTitle="管理平台所有注册用户">
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

            {/* 编辑弹窗 */}
            <Modal
                title={currentUser ? '编辑用户' : '新增用户'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={600}
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
                width={500}
            >
                {currentUser && (
                    <>
                        <div style={{ textAlign: 'center', marginBottom: 24 }}>
                            <Avatar
                                size={80}
                                src={currentUser.avatarUrl}
                                icon={<UserOutlined />}
                                style={{ backgroundColor: token.colorPrimary }}
                            />
                            <h2 style={{ marginTop: 16, marginBottom: 4 }}>{currentUser.name}</h2>
                            <Tag color={roleMap[currentUser.role].color}>{roleMap[currentUser.role].text}</Tag>
                            <Tag color={statusMap[currentUser.status].color}>{statusMap[currentUser.status].text}</Tag>
                        </div>

                        <Divider />

                        <Descriptions column={1} labelStyle={{ width: 100 }}>
                            <Descriptions.Item label="用户ID">{currentUser.id}</Descriptions.Item>
                            <Descriptions.Item label="邮箱">{currentUser.email}</Descriptions.Item>
                            <Descriptions.Item label="手机号">{currentUser.phone}</Descriptions.Item>
                            <Descriptions.Item label="注册时间">
                                {dayjs(currentUser.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                            </Descriptions.Item>
                            <Descriptions.Item label="最后登录">
                                {currentUser.lastLoginAt ? dayjs(currentUser.lastLoginAt).format('YYYY-MM-DD HH:mm:ss') : '从未登录'}
                            </Descriptions.Item>
                        </Descriptions>
                    </>
                )}
            </Drawer>
        </PageContainer>
    );
};

export default UserPage;
