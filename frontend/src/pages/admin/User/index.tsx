/**
 * 用户管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
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
import dayjs from 'dayjs';

/**
 * 用户数据接口
 */
interface User {
    id: number;
    username: string;
    email: string;
    phone: string;
    avatar: string;
    role: 'user' | 'player' | 'admin';
    status: 'active' | 'banned' | 'pending';
    balance: number;
    createdAt: string;
    lastLoginAt: string;
}

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
    pending: { color: 'warning', text: '待审核' },
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
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [form] = Form.useForm();

    /**
     * 加载用户数据
     */
    const loadData = useCallback(async () => {
        // ... (keep existing loadData)
        setLoading(true);
        // 模拟API调用
        await new Promise(resolve => setTimeout(resolve, 500));

        const mockUsers: User[] = Array.from({ length: 50 }, (_, i) => ({
            id: i + 1,
            username: `user${i + 1}`,
            email: `user${i + 1}@example.com`,
            phone: `138${String(i).padStart(8, '0')}`,
            avatar: '',
            role: ['user', 'player', 'admin'][i % 3] as User['role'],
            status: ['active', 'banned', 'pending'][i % 3] as User['status'],
            balance: Math.floor(Math.random() * 10000),
            createdAt: dayjs().subtract(i, 'day').format('YYYY-MM-DD HH:mm:ss'),
            lastLoginAt: dayjs().subtract(Math.floor(Math.random() * 24), 'hour').format('YYYY-MM-DD HH:mm:ss'),
        }));

        // 模拟分页
        const start = (current - 1) * pageSize;
        const end = start + pageSize;
        setUsers(mockUsers.slice(start, end));
        setTotal(mockUsers.length);
        setLoading(false);
    }, [current, pageSize]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 搜索
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setSearchParams(values);
        setCurrent(1);
        loadData();
    };

    /**
     * 编辑用户
     */
    const handleEdit = (user: User) => {
        setCurrentUser(user);
        form.setFieldsValue(user);
        setEditModalVisible(true);
    };

    /**
     * 查看详情
     */
    const handleViewDetail = (user: User) => {
        setCurrentUser(user);
        setDetailDrawerVisible(true);
    };

    /**
     * 保存编辑
     */
    const handleSaveEdit = async () => {
        try {
            const values = await form.validateFields();
            console.log('Save:', values);
            message.success('保存成功');
            setEditModalVisible(false);
            loadData();
        } catch {
            // 验证失败
        }
    };

    /**
     * 封禁/解封用户
     */
    const handleToggleBan = async (user: User) => {
        const action = user.status === 'banned' ? '解封' : '封禁';
        message.success(`${action}成功`);
        loadData();
    };

    /**
     * 删除用户
     */
    const handleDelete = async (user: User) => {
        message.success(`删除用户 ${user.username} 成功`);
        loadData();
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
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
                { label: '待审核', value: 'pending' },
            ],
        },
        { name: 'dateRange', label: '注册时间', type: 'dateRange' },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<User> = [
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
                        src={record.avatar}
                        icon={<UserOutlined />}
                        style={{ backgroundColor: token.colorPrimary }}
                    />
                    <div>
                        <div style={{ fontWeight: 500 }}>{record.username}</div>
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
            render: role => <Tag color={roleMap[role].color}>{roleMap[role].text}</Tag>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: status => <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>,
        },
        {
            title: '余额',
            dataIndex: 'balance',
            key: 'balance',
            width: 100,
            render: balance => <span>¥{balance.toFixed(2)}</span>,
        },
        {
            title: '注册时间',
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
                            title="确定要删除该用户吗？此操作不可恢复。"
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
    ];

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
                onBatchDelete={async (keys) => {
                    console.log('批量删除:', keys);
                    await new Promise(resolve => setTimeout(resolve, 500));
                    loadData();
                }}
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
                scroll={{ x: 1200 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentUser ? '编辑用户' : '新增用户'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="username"
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
                        rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' }]}
                    >
                        <Input placeholder="请输入手机号" />
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
                                { label: '待审核', value: 'pending' },
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
                                src={currentUser.avatar}
                                icon={<UserOutlined />}
                                style={{ backgroundColor: token.colorPrimary }}
                            />
                            <h2 style={{ marginTop: 16, marginBottom: 4 }}>{currentUser.username}</h2>
                            <Tag color={roleMap[currentUser.role].color}>{roleMap[currentUser.role].text}</Tag>
                            <Tag color={statusMap[currentUser.status].color}>{statusMap[currentUser.status].text}</Tag>
                        </div>


                        <Divider />

                        <Descriptions column={1} labelStyle={{ width: 100 }}>
                            <Descriptions.Item label="用户ID">{currentUser.id}</Descriptions.Item>
                            <Descriptions.Item label="邮箱">{currentUser.email}</Descriptions.Item>
                            <Descriptions.Item label="手机号">{currentUser.phone}</Descriptions.Item>
                            <Descriptions.Item label="余额">¥{currentUser.balance.toFixed(2)}</Descriptions.Item>
                            <Descriptions.Item label="注册时间">{currentUser.createdAt}</Descriptions.Item>
                            <Descriptions.Item label="最后登录">{currentUser.lastLoginAt}</Descriptions.Item>
                        </Descriptions>
                    </>
                )}
            </Drawer>
        </PageContainer>
    );
};

export default UserPage;
