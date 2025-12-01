/**
 * 角色管理页面
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
    Tree,
    Card,
    Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { TreeDataNode } from 'antd';
import {
    EditOutlined,
    DeleteOutlined,
    SafetyCertificateOutlined,
    SettingOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { ROLE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi } from '@/api/admin';

const { Text } = Typography;

/**
 * 角色数据接口
 */
interface Role {
    id: number;
    name: string;
    code: string;
    description: string;
    isSystem: boolean;
    userCount: number;
    permissionCount: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * 角色管理页面
 */
const RolePage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [roles, setRoles] = useState<Role[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 弹窗状态
    const [editModalVisible, setEditModalVisible] = useState(false);
    const [permModalVisible, setPermModalVisible] = useState(false);
    const [currentRole, setCurrentRole] = useState<Role | null>(null);
    const [form] = Form.useForm();

    // 权限树
    const [permissionTree, setPermissionTree] = useState<TreeDataNode[]>([]);
    const [checkedKeys, setCheckedKeys] = useState<React.Key[]>([]);
    const [permLoading, setPermLoading] = useState(false);

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
            }) as any;
            const { items, totalCount } = res.data || { items: [], totalCount: 0 };
            setRoles(items);
            setTotal(totalCount);
        } catch (error) {
            console.error(error);
            message.error('加载角色列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize]);

    /**
     * 加载权限树
     */
    const loadPermissionTree = useCallback(async () => {
        setPermLoading(true);
        try {
            // 这里假设有一个获取所有权限树的接口，如果没有，可能需要从菜单接口转换
            // 暂时使用模拟数据，或者调用 getMenus 获取菜单作为权限树
            // 实际项目中应该有 adminApi.getPermissionTree()
            // 实际项目中应该有 adminApi.getPermissionTree()
            const res = await adminApi.getMenus() as any;
            const menus = res.data || [];

            const convertToTree = (items: any[]): TreeDataNode[] => {
                return items.map(item => ({
                    title: item.name,
                    key: item.permission || `menu_${item.id}`,
                    children: item.children ? convertToTree(item.children) : undefined
                }));
            };

            setPermissionTree(convertToTree(menus));
        } catch (error) {
            console.error(error);
        } finally {
            setPermLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
        loadPermissionTree();
    }, [loadData, loadPermissionTree]);

    /**
     * 编辑角色
     */
    const handleEdit = (role: Role) => {
        setCurrentRole(role);
        form.setFieldsValue(role);
        setEditModalVisible(true);
    };

    /**
     * 配置权限
     */
    const handleConfigPermission = async (role: Role) => {
        setCurrentRole(role);
        setPermModalVisible(true);
        setPermLoading(true);
        try {
            // 获取角色当前权限
            // const res = await adminApi.getRolePermissions(role.id);
            // setCheckedKeys(res.data.data);
            // 暂时置空或模拟
            setCheckedKeys([]);
        } catch (error) {
            message.error('获取角色权限失败');
        } finally {
            setPermLoading(false);
        }
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
     * 保存权限配置
     */
    const handleSavePermission = async () => {
        if (!currentRole) return;
        try {
            // @ts-ignore
            await adminApi.assignRolePermissions(currentRole.id, checkedKeys);
            message.success('权限配置成功');
            setPermModalVisible(false);
            loadData();
        } catch (error) {
            console.error(error);
            message.error('权限配置失败');
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
            width: 150,
            render: (text, record) => (
                <Space>
                    <SafetyCertificateOutlined style={{ color: record.isSystem ? '#1890ff' : '#52c41a' }} />
                    <span style={{ fontWeight: 500 }}>{text}</span>
                    {record.isSystem && <Tag color="blue">系统</Tag>}
                </Space>
            ),
        },
        {
            title: '角色编码',
            dataIndex: 'code',
            key: 'code',
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
            dataIndex: 'userCount',
            key: 'userCount',
            width: 100,
            render: count => <Tag>{count || 0} 人</Tag>,
        },
        {
            title: '权限数',
            dataIndex: 'permissionCount',
            key: 'permissionCount',
            width: 100,
            render: count => <Tag color="purple">{count || 0} 项</Tag>,
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
                        name="code"
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

            {/* 权限配置弹窗 */}
            <Modal
                title={`配置权限 - ${currentRole?.name}`}
                open={permModalVisible}
                onOk={handleSavePermission}
                onCancel={() => setPermModalVisible(false)}
                width={600}
            >
                <Card loading={permLoading}>
                    <Tree
                        checkable
                        defaultExpandAll
                        checkedKeys={checkedKeys}
                        onCheck={(checked) => setCheckedKeys(checked as React.Key[])}
                        treeData={permissionTree}
                    />
                </Card>
            </Modal>
        </PageContainer>
    );
};

export default RolePage;
