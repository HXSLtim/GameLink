/**
 * 角色管理页面
 * Requirements: 2.1, 2.4, 2.5
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    message,
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
import type { Role } from '@/api/admin';
import type { ApiResponse } from '@/api/admin';

const { Text } = Typography;

/**
 * 角色管理页面
 * Requirements: 2.1, 2.4, 2.5
 */
const RolePage: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [roles, setRoles] = useState<Role[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [currentRole, setCurrentRole] = useState<Role | null>(null);
    const [form] = Form.useForm();

    /**
     * 加载角色数据
     */
    const loadData = useCallback(async (params: any = {}) => {
        setLoading(true);
        try {
            const res = await adminApi.getRoles({
                page: current,
                page_size: pageSize,
                ...params
            }) as unknown as ApiResponse<{ items: Role[], totalCount: number }>;
            if (res.success) {
                const { items, totalCount } = res.data;
                setRoles(items);
                setTotal(totalCount);
            }
        } catch (error) {
            console.error(error);
            message.error('加载角色列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    /**
     * 编辑角色
     */
    const handleEdit = (role: Role) => {
        setCurrentRole(role);
        form.setFieldsValue(role);
        setEditModalVisible(true);
    };

    /**
     * 配置权限 - 跳转到权限配置页面
     * Requirements: 2.4 - 查看角色已有权限
     * Requirements: 2.5 - 系统角色显示特殊提示
     */
    const handleConfigPermission = (role: Role) => {
        navigate(`/admin/sys/role/${role.id}/permissions`);
    };

    /**
     * 保存编辑
     */
    const handleSaveEdit = async () => {
        try {
            const values = await form.validateFields();
            if (currentRole) {
                await adminApi.updateRole(currentRole.id, values);
                message.success('更新成功');
            } else {
                await adminApi.createRole(values);
                message.success('创建成功');
            }
            setEditModalVisible(false);
            loadData();
        } catch (error) {
            console.error(error);
            message.error('保存失败');
        }
    };

    /**
     * 删除角色
     */
    const handleDelete = async (role: Role) => {
        if (role.isSystem) {
            message.error('系统角色不可删除');
            return;
        }
        try {
            await adminApi.deleteRole(role.id);
            message.success(`删除角色 ${role.name} 成功`);
            loadData();
        } catch (error) {
            console.error(error);
            message.error('删除失败');
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '角色名称/编码' },
    ];

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
            render: (users: any[]) => <Tag>{users?.length || 0} 人</Tag>,
        },
        {
            title: '权限数',
            dataIndex: 'permissions',
            key: 'permissions',
            width: 100,
            render: (permissions: any[]) => <Tag color="purple">{permissions?.length || 0} 项</Tag>,
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
                onSearch={(values) => {
                    setCurrent(1);
                    loadData(values);
                }}
                onRefresh={loadData}
                loading={loading}
                showCreate={true}
                createText="新增角色"
                createPermission={ROLE_PERMISSIONS.CREATE}
                onCreate={() => {
                    setCurrentRole(null);
                    form.resetFields();
                    setEditModalVisible(true);
                }}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: total => `共 ${total} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
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
