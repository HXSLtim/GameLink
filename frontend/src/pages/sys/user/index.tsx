import React, { useEffect, useState } from 'react';
import { Table, Tag, Space, Button, Avatar, Input, Card } from 'antd';
import { SearchOutlined, EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { motion } from 'framer-motion';
import { adminApi } from '@/api/admin';
import type { User } from '@/api/admin';
import { PermissionGuard } from '@/components/PermissionGuard';
import { usePermissions } from '@/hooks/usePermission';
import { USER_PERMISSIONS } from '@/constants/permissions';

const Users: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [searchText, setSearchText] = useState('');

    // 批量检查权限
    const permissions = usePermissions({
        canCreate: USER_PERMISSIONS.CREATE,
        canEdit: USER_PERMISSIONS.UPDATE,
        canDelete: USER_PERMISSIONS.DELETE,
    });

    useEffect(() => {
        // Mock data
        setUsers([
            { id: '1', username: 'GamerOne', email: 'gamer1@example.com', role: 'USER', status: 'active', createdAt: '2023-01-01' },
            { id: '2', username: 'ProPlayer', email: 'pro@example.com', role: 'COMPANION', status: 'active', createdAt: '2023-01-05' },
            { id: '3', username: 'BannedUser', email: 'bad@example.com', role: 'USER', status: 'banned', createdAt: '2023-02-10' },
            { id: '4', username: 'AdminUser', email: 'admin@gamelink.com', role: 'ADMIN', status: 'active', createdAt: '2022-12-01' },
            { id: '5', username: 'Newbie', email: 'new@example.com', role: 'USER', status: 'pending', createdAt: '2023-03-15' },
        ]);
        // Silence unused warning
        console.log(adminApi, loading, setLoading);
    }, []);

    const columns: ColumnsType<User> = [
        {
            title: 'User',
            dataIndex: 'username',
            key: 'username',
            render: (text, record) => (
                <Space>
                    <Avatar style={{ backgroundColor: '#5865F2' }} icon={<UserOutlined />} />
                    <div>
                        <div style={{ fontWeight: 500 }}>{text}</div>
                        <div style={{ fontSize: 12, color: 'rgba(255,255,255,0.45)' }}>{record.email}</div>
                    </div>
                </Space>
            ),
        },
        {
            title: 'Role',
            dataIndex: 'role',
            key: 'role',
            render: (role) => {
                let color = 'default';
                if (role === 'ADMIN') color = 'magenta';
                if (role === 'COMPANION') color = 'purple';
                if (role === 'USER') color = 'blue';
                return <Tag color={color}>{role}</Tag>;
            },
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status) => {
                let color = 'default';
                if (status === 'active') color = 'success';
                if (status === 'banned') color = 'error';
                if (status === 'pending') color = 'warning';
                return <Tag color={color}>{status.toUpperCase()}</Tag>;
            },
        },
        {
            title: 'Joined Date',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: (text) => <span style={{ color: 'rgba(255,255,255,0.65)' }}>{text}</span>,
        },
        {
            title: 'Actions',
            key: 'actions',
            render: () => (
                <Space size="small">
                    {/* 使用 PermissionGuard 组件控制按钮显示 */}
                    <PermissionGuard permission={USER_PERMISSIONS.UPDATE}>
                        <Button type="text" icon={<EditOutlined />} style={{ color: '#5865F2' }} />
                    </PermissionGuard>
                    <PermissionGuard permission={USER_PERMISSIONS.DELETE}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                    </PermissionGuard>
                </Space>
            ),
        },
    ];

    const filteredUsers = users.filter(user =>
        user.username.toLowerCase().includes(searchText.toLowerCase()) ||
        user.email.toLowerCase().includes(searchText.toLowerCase())
    );

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
        >
            <Card variant="borderless" bodyStyle={{ padding: '24px' }}>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h2 style={{ margin: 0, fontSize: 20 }}>User Management</h2>
                    <Space>
                        <Input
                            placeholder="Search users..."
                            prefix={<SearchOutlined />}
                            onChange={e => setSearchText(e.target.value)}
                            style={{ width: 250 }}
                        />
                        {/* 使用 usePermissions 批量检查结果控制按钮显示 */}
                        {permissions.canCreate && (
                            <Button type="primary" style={{ backgroundColor: '#5865F2' }}>Add User</Button>
                        )}
                    </Space>
                </div>
                <Table
                    columns={columns}
                    dataSource={filteredUsers}
                    rowKey="id"
                    pagination={{ pageSize: 10 }}
                />
            </Card>
        </motion.div>
    );
};

export default Users;
