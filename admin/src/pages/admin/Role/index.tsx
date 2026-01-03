/**
 * 角色管理页面
 * Requirements: 2.1, 2.4, 2.5
 */
import React, { useState, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    Popconfirm,
    Typography,
    Tooltip,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
    SafetyCertificateOutlined,
    SettingOutlined,
    LockOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { ROLE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi } from '@/api/admin';
import type { Role, CreateRoleDto, UpdateRoleDto } from '@/api/admin';
import { useCrud } from '@/hooks';

const { Text } = Typography;

/**
 * 角色管理页面
 */
const RolePage: React.FC = () => {
    const navigate = useNavigate();

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [currentRole, setCurrentRole] = useState<Role | null>(null);
    const [form] = Form.useForm();

    /**
     * 使用 CRUD Hook 管理角色数据
     */
    const {
        data: roles,
        loading,
        pagination,
        fetchAll,
        create: createRole,
        update: updateRole,
        remove: deleteRole,
        setSearchParams,
    } = useCrud<Role, CreateRoleDto, UpdateRoleDto>({
        api: {
            getAll: adminApi.getRoles as any,
            create: adminApi.createRole as any,
            update: adminApi.updateRole as any,
            remove: adminApi.deleteRole as any,
        },
        messages: {
            fetchError: '加载角色列表失败',
            createSuccess: '创建角色成功',
            updateSuccess: '更新角色成功',
            deleteSuccess: '删除角色成功',
        },
        initialPagination: {
            pageSize: 10,
        },
        paginationExtractor: (response) => {
            // Extract total from nested pagination object
            const res = response as { data?: { pagination?: { total?: number } } };
            return res.data?.pagination?.total;
        },
        dataTransformer: (rawData) => {
            // Handle different response formats
            if (Array.isArray(rawData)) {
                return rawData as Role[];
            }
            const data = rawData as { items?: Role[]; totalCount?: number };
            return data.items || [];
        },
    });

    /**
     * 配置权限 - 跳转到权限配置页面
     * Requirements: 2.4 - 查看角色已有权限
     * Requirements: 2.5 - 系统角色显示特殊提示
     */
    const handleConfigPermission = useCallback((role: Role) => {
        navigate(`/admin/sys/role/${role.id}/permissions`);
    }, [navigate]);

    /**
     * 编辑角色
     */
    const handleEdit = useCallback((role: Role) => {
        setCurrentRole(role);
        form.setFieldsValue(role);
        setEditModalVisible(true);
    }, [form]);

    /**
     * 新增角色
     */
    const handleCreate = useCallback(() => {
        setCurrentRole(null);
        form.resetFields();
        setEditModalVisible(true);
    }, [form]);

    /**
     * 保存编辑
     */
    const handleSaveEdit = useCallback(async () => {
        try {
            const values = await form.validateFields();

            if (currentRole) {
                await updateRole(currentRole.id, values);
            } else {
                await createRole(values);
            }

            setEditModalVisible(false);
        } catch (err) {
            // Form validation error or API error (handled by hook)
            if (err && typeof err === 'object' && 'errorFields' in err) {
                // Form validation error from Ant Design
                console.error('Form validation error:', err);
            } else {
                console.error('Save error:', err);
            }
        }
    }, [currentRole, form, updateRole, createRole]);

    /**
     * 删除角色
     */
    const handleDelete = useCallback(async (role: Role) => {
        if (role.isSystem) {
            Modal.error({
                title: '操作失败',
                content: '系统角色不可删除',
            });
            return;
        }

        await deleteRole(role.id, {
            confirmMessage: `确定要删除角色 "${role.name}" 吗？`,
        });
    }, [deleteRole]);

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '角色名称/编码' },
    ];

    /**
     * 搜索处理
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        setSearchParams(values);
    }, [setSearchParams]);

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Role> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '角色名称',
            dataIndex: 'name',
            key: 'name',
            width: 180,
            render: (text, record) => (
                <Space>
                    <SafetyCertificateOutlined style={{ color: record.isSystem ? '#1890ff' : '#52c41a' }} />
                    <span style={{ fontWeight: 500 }}>{text}</span>
                    {record.isSystem && (
                        <Tooltip title="系统角色不可删除">
                            <Tag color="blue" icon={<LockOutlined />}>系统</Tag>
                        </Tooltip>
                    )}
                    {record.slug === 'superAdmin' && (
                        <Tooltip title="超级管理员拥有所有权限">
                            <Tag color="gold">超管</Tag>
                        </Tooltip>
                    )}
                </Space>
            ),
        },
        {
            title: '角色编码',
            dataIndex: 'slug',
            key: 'slug',
            width: 120,
            render: text => <Text code>{text}</Text>,
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            width: 200,
            ellipsis: true,
        },
        {
            title: '用户数',
            dataIndex: 'users',
            key: 'users',
            width: 100,
            render: (users: unknown[]) => <Tag>{users?.length || 0} 人</Tag>,
        },
        {
            title: '权限数',
            dataIndex: 'permissions',
            key: 'permissions',
            width: 100,
            render: (permissions: unknown[]) => <Tag color="purple">{permissions?.length || 0} 项</Tag>,
        },
        {
            title: '更新时间',
            dataIndex: 'updatedAt',
            key: 'updatedAt',
            width: 180,
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <PermissionGuard permission={ROLE_PERMISSIONS.ASSIGN_PERMISSIONS}>
                        <Button
                            type="link"
                            size="small"
                            icon={<SettingOutlined />}
                            onClick={() => handleConfigPermission(record)}
                        >
                            权限
                        </Button>
                    </PermissionGuard>
                    <PermissionGuard permission={ROLE_PERMISSIONS.UPDATE}>
                        <Button
                            type="link"
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => handleEdit(record)}
                        >
                            编辑
                        </Button>
                    </PermissionGuard>
                    {!record.isSystem && (
                        <PermissionGuard permission={ROLE_PERMISSIONS.DELETE}>
                            <Popconfirm
                                title="确定要删除该角色吗？"
                                onConfirm={() => handleDelete(record)}
                            >
                                <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                                    删除
                                </Button>
                            </Popconfirm>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <PageContainer title="角色管理" subTitle="管理系统角色和权限分配">
            <SearchTable
                columns={columns}
                dataSource={roles}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => fetchAll()}
                loading={loading}
                showCreate={true}
                createText="新增角色"
                createPermission={ROLE_PERMISSIONS.CREATE}
                onCreate={handleCreate}
                pagination={pagination}
                scroll={{ x: 1100 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentRole ? '编辑角色' : '新增角色'}
                open={editModalVisible}
                onOk={handleSaveEdit}
                onCancel={() => setEditModalVisible(false)}
                width={500}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="角色名称"
                        rules={[{ required: true, message: '请输入角色名称' }]}
                    >
                        <Input placeholder="请输入角色名称" />
                    </Form.Item>
                    <Form.Item
                        name="slug"
                        label="角色编码"
                        rules={[
                            { required: true, message: '请输入角色编码' },
                            { pattern: /^[a-z_]+$/, message: '只能输入小写字母和下划线' },
                        ]}
                    >
                        <Input
                            placeholder="请输入角色编码，如 operator"
                            disabled={!!currentRole?.isSystem}
                        />
                    </Form.Item>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea placeholder="请输入角色描述" rows={3} />
                    </Form.Item>
                </Form>
            </Modal>

        </PageContainer>
    );
};

export default RolePage;
