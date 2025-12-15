/**
 * 用户角色分配页面
 * Requirements: 9.1, 9.2, 9.3, 10.3
 * 
 * 功能：
 * - 显示用户列表及其已分配的RBAC角色
 * - 支持为单个用户分配角色
 * - 支持批量为多个用户分配角色
 * - 查看用户有效权限（合并后的完整权限列表）
 */
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
    Tag,
    Space,
    Button,
    Avatar,
    Modal,
    message,
    Drawer,
    Descriptions,
    Divider,
    theme,
    Checkbox,
    List,
    Typography,
    Alert,
    Spin,
    Result,
    Tooltip,
    Badge,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    UserOutlined,
    SafetyCertificateOutlined,
    TeamOutlined,
    EyeOutlined,
    CheckCircleOutlined,
    InfoCircleOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField, ToolbarButton } from '@/components';
import { ROLE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { adminApi, type User, type UserQueryParams, type ApiResponse } from '@/api/admin';
import { userRoleApi, roleApi } from '@/api/permission';
import type { Role, UserEffectivePermissions } from '@/types/permission';
import dayjs from 'dayjs';

const { Text, Title } = Typography;

/**
 * 用户状态映射
 */
const statusMap = {
    active: { color: 'success', text: '正常' },
    banned: { color: 'error', text: '已封禁' },
    suspended: { color: 'warning', text: '已停用' },
};

/**
 * 扩展用户类型，包含RBAC角色
 */
interface UserWithRoles extends User {
    rbacRoles?: Role[];
}

/**
 * 用户角色分配页面
 */
const UserRolePage: React.FC = () => {
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [users, setUsers] = useState<UserWithRoles[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchParams, setSearchParams] = useState<UserQueryParams>({});

    // 角色列表
    const [roles, setRoles] = useState<Role[]>([]);
    const [rolesLoading, setRolesLoading] = useState(false);

    // 角色分配弹窗
    const [assignModalVisible, setAssignModalVisible] = useState(false);
    const [currentUser, setCurrentUser] = useState<UserWithRoles | null>(null);
    const [selectedRoleIds, setSelectedRoleIds] = useState<number[]>([]);
    const [assignLoading, setAssignLoading] = useState(false);

    // 批量角色分配
    const [batchAssignVisible, setBatchAssignVisible] = useState(false);
    const [selectedUserIds, setSelectedUserIds] = useState<number[]>([]);
    const [batchSelectedRoleIds, setBatchSelectedRoleIds] = useState<number[]>([]);
    const [batchAssignLoading, setBatchAssignLoading] = useState(false);

    // 有效权限查看
    const [permissionDrawerVisible, setPermissionDrawerVisible] = useState(false);
    const [effectivePermissions, setEffectivePermissions] = useState<UserEffectivePermissions | null>(null);
    const [permissionsLoading, setPermissionsLoading] = useState(false);

    /**
     * 加载角色列表
     */
    const loadRoles = useCallback(async () => {
        setRolesLoading(true);
        try {
            const res = await roleApi.list({ page: 1, page_size: 100 });
            if (res.data?.success) {
                setRoles(res.data.data?.items || []);
            }
        } catch (error) {
            console.error('加载角色列表失败:', error);
        } finally {
            setRolesLoading(false);
        }
    }, []);

    /**
     * 加载用户数据（包含RBAC角色）
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const params: UserQueryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
            };

            const response = await adminApi.getUsers(params) as unknown as ApiResponse<User[]>;

            if (response.success) {
                const userList = response.data || [];
                
                // 并行获取每个用户的RBAC角色
                const usersWithRoles = await Promise.all(
                    userList.map(async (user) => {
                        try {
                            const rolesRes = await userRoleApi.getUserRoles(user.id);
                            return {
                                ...user,
                                rbacRoles: rolesRes.data?.success ? rolesRes.data.data : [],
                            };
                        } catch {
                            return { ...user, rbacRoles: [] };
                        }
                    })
                );

                setUsers(usersWithRoles);
                setTotal(response.pagination?.total || 0);
            } else {
                message.error(response.message || '获取用户列表失败');
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
        loadRoles();
    }, [loadData, loadRoles]);

    /**
     * 搜索处理
     */
    const handleSearch = useCallback((values: Record<string, unknown>) => {
        const params: UserQueryParams = {};
        const keyword = values.keyword as string | undefined;
        if (keyword?.trim()) {
            params.keyword = keyword.trim();
        }
        if (values.status) {
            const statuses = Array.isArray(values.status) ? values.status : [values.status];
            if (statuses.length > 0) {
                params.status = statuses as string[];
            }
        }
        setSearchParams(params);
        setCurrent(1);
    }, []);

    /**
     * 打开角色分配弹窗
     */
    const handleOpenAssignModal = useCallback((user: UserWithRoles) => {
        setCurrentUser(user);
        setSelectedRoleIds(user.rbacRoles?.map(r => r.id) || []);
        setAssignModalVisible(true);
    }, []);

    /**
     * 提交角色分配
     */
    const handleAssignRoles = useCallback(async () => {
        if (!currentUser) return;

        setAssignLoading(true);
        try {
            const res = await userRoleApi.assignRoles(currentUser.id, { roleIds: selectedRoleIds });
            if (res.data?.success) {
                message.success('角色分配成功');
                setAssignModalVisible(false);
                loadData();
            } else {
                message.error('角色分配失败');
            }
        } catch (error: unknown) {
            console.error('角色分配失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '角色分配失败');
        } finally {
            setAssignLoading(false);
        }
    }, [currentUser, selectedRoleIds, loadData]);

    /**
     * 打开批量角色分配弹窗
     */
    const handleOpenBatchAssign = useCallback((keys?: React.Key[]) => {
        if (!keys || keys.length === 0) {
            message.warning('请先选择用户');
            return;
        }
        setSelectedUserIds(keys.map(k => Number(k)));
        setBatchSelectedRoleIds([]);
        setBatchAssignVisible(true);
    }, []);

    /**
     * 提交批量角色分配
     */
    const handleBatchAssignRoles = useCallback(async () => {
        if (selectedUserIds.length === 0 || batchSelectedRoleIds.length === 0) {
            message.warning('请选择用户和角色');
            return;
        }

        setBatchAssignLoading(true);
        try {
            const res = await userRoleApi.batchAssignRoles({
                userIds: selectedUserIds,
                roleIds: batchSelectedRoleIds,
            });

            if (res.data?.success) {
                const result = res.data.data;
                // Backend returns successCount/failedCount
                const successCount = result?.successCount ?? 0;
                const failedCount = result?.failedCount ?? 0;
                if (failedCount > 0) {
                    message.warning(`成功: ${successCount}, 失败: ${failedCount}`);
                } else {
                    message.success(`成功为 ${successCount || selectedUserIds.length} 个用户分配角色`);
                }
                setBatchAssignVisible(false);
                loadData();
            } else {
                message.error('批量角色分配失败');
            }
        } catch (error: unknown) {
            console.error('批量角色分配失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '批量角色分配失败');
        } finally {
            setBatchAssignLoading(false);
        }
    }, [selectedUserIds, batchSelectedRoleIds, loadData]);

    /**
     * 查看用户有效权限
     */
    const handleViewPermissions = useCallback(async (user: UserWithRoles) => {
        setCurrentUser(user);
        setPermissionDrawerVisible(true);
        setPermissionsLoading(true);

        try {
            const res = await userRoleApi.getUserPermissions(user.id);
            if (res.data?.success) {
                setEffectivePermissions(res.data.data);
            } else {
                message.error('获取用户权限失败');
            }
        } catch (error: unknown) {
            console.error('获取用户权限失败:', error);
            const err = error as { response?: { data?: { message?: string } } };
            message.error(err.response?.data?.message || '获取用户权限失败');
        } finally {
            setPermissionsLoading(false);
        }
    }, []);

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = useMemo(() => [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '用户名/邮箱/手机号' },
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
    ], []);

    /**
     * 表格列配置
     */
    const columns: ColumnsType<UserWithRoles> = useMemo(() => [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '用户信息',
            key: 'userInfo',
            width: 220,
            render: (_, record) => (
                <Space>
                    <Avatar
                        src={record.avatarUrl || undefined}
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
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: User['status']) => (
                <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>
            ),
        },
        {
            title: '已分配角色',
            key: 'rbacRoles',
            width: 300,
            render: (_, record) => {
                const rbacRoles = record.rbacRoles || [];
                if (rbacRoles.length === 0) {
                    return <Text type="secondary">未分配角色</Text>;
                }
                return (
                    <Space size={[0, 4]} wrap>
                        {rbacRoles.map(role => (
                            <Tag
                                key={role.id}
                                color={role.isSystem ? 'blue' : 'green'}
                                icon={<SafetyCertificateOutlined />}
                            >
                                {role.name}
                                {role.slug === 'superAdmin' && (
                                    <Badge count="超管" style={{ marginLeft: 4, backgroundColor: '#faad14' }} />
                                )}
                            </Tag>
                        ))}
                    </Space>
                );
            },
        },
        {
            title: '注册时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (createdAt: string) => dayjs(createdAt).format('YYYY-MM-DD HH:mm'),
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <PermissionGuard permission={ROLE_PERMISSIONS.ASSIGN_USER}>
                        <Button
                            type="link"
                            size="small"
                            icon={<TeamOutlined />}
                            onClick={() => handleOpenAssignModal(record)}
                        >
                            分配角色
                        </Button>
                    </PermissionGuard>
                    <Button
                        type="link"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewPermissions(record)}
                    >
                        查看权限
                    </Button>
                </Space>
            ),
        },
    ], [token, handleOpenAssignModal, handleViewPermissions]);

    /**
     * 工具栏按钮
     */
    const toolbarButtons: ToolbarButton[] = useMemo(() => [
        {
            text: '批量分配角色',
            icon: <TeamOutlined />,
            needSelection: true,
            onClick: handleOpenBatchAssign,
            permission: ROLE_PERMISSIONS.ASSIGN_USER,
        },
    ], [handleOpenBatchAssign]);

    /**
     * 渲染角色选择列表
     */
    const renderRoleCheckboxList = (
        selectedIds: number[],
        onChange: (ids: number[]) => void
    ) => (
        <Checkbox.Group
            value={selectedIds}
            onChange={(values) => onChange(values as number[])}
            style={{ width: '100%' }}
        >
            <List
                loading={rolesLoading}
                dataSource={roles}
                renderItem={(role) => (
                    <List.Item>
                        <Checkbox value={role.id} style={{ width: '100%' }}>
                            <Space>
                                <SafetyCertificateOutlined
                                    style={{ color: role.isSystem ? token.colorPrimary : token.colorSuccess }}
                                />
                                <span style={{ fontWeight: 500 }}>{role.name}</span>
                                <Text type="secondary" code>{role.slug}</Text>
                                {role.isSystem && <Tag color="blue">系统</Tag>}
                                {role.slug === 'superAdmin' && (
                                    <Tooltip title="超级管理员拥有所有权限">
                                        <Tag color="gold">超管</Tag>
                                    </Tooltip>
                                )}
                            </Space>
                            {role.description && (
                                <div style={{ marginLeft: 24, color: token.colorTextSecondary, fontSize: 12 }}>
                                    {role.description}
                                </div>
                            )}
                        </Checkbox>
                    </List.Item>
                )}
            />
        </Checkbox.Group>
    );

    /**
     * 渲染有效权限列表
     */
    const renderEffectivePermissions = () => {
        if (permissionsLoading) {
            return <Spin tip="加载中..." />;
        }

        if (!effectivePermissions) {
            return <Result status="warning" title="无法获取权限信息" />;
        }

        const { permissions, roles: userRoles, isSuperAdmin } = effectivePermissions;

        return (
            <div>
                {/* 超级管理员提示 */}
                {isSuperAdmin && (
                    <Alert
                        message="超级管理员"
                        description="该用户是超级管理员，拥有系统所有权限"
                        type="warning"
                        showIcon
                        icon={<SafetyCertificateOutlined />}
                        style={{ marginBottom: 16 }}
                    />
                )}

                {/* 角色来源 */}
                <Title level={5}>
                    <TeamOutlined style={{ marginRight: 8 }} />
                    角色来源
                </Title>
                <div style={{ marginBottom: 16 }}>
                    {userRoles.length === 0 ? (
                        <Text type="secondary">未分配任何角色</Text>
                    ) : (
                        <Space wrap>
                            {userRoles.map(role => (
                                <Tag
                                    key={role.id}
                                    color={role.isSystem ? 'blue' : 'green'}
                                    icon={<SafetyCertificateOutlined />}
                                >
                                    {role.name}
                                </Tag>
                            ))}
                        </Space>
                    )}
                </div>

                <Divider />

                {/* 有效权限列表 */}
                <Title level={5}>
                    <CheckCircleOutlined style={{ marginRight: 8 }} />
                    有效权限 ({permissions.length})
                </Title>
                {isSuperAdmin ? (
                    <Alert
                        message="拥有所有权限 (*)"
                        type="info"
                        showIcon
                    />
                ) : permissions.length === 0 ? (
                    <Text type="secondary">暂无权限</Text>
                ) : (
                    <div style={{ maxHeight: 400, overflow: 'auto' }}>
                        <List
                            size="small"
                            dataSource={permissions}
                            renderItem={(perm) => (
                                <List.Item>
                                    <Text code>{perm}</Text>
                                </List.Item>
                            )}
                        />
                    </div>
                )}
            </div>
        );
    };

    return (
        <PageContainer title="用户角色分配" subTitle="管理用户的RBAC角色分配">
            <SearchTable
                columns={columns}
                dataSource={users}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={loadData}
                loading={loading}
                toolbarButtons={toolbarButtons}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (t) => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1100 }}
            />

            {/* 单用户角色分配弹窗 */}
            <Modal
                title={
                    <Space>
                        <TeamOutlined />
                        <span>分配角色 - {currentUser?.name}</span>
                    </Space>
                }
                open={assignModalVisible}
                onOk={handleAssignRoles}
                onCancel={() => setAssignModalVisible(false)}
                confirmLoading={assignLoading}
                width={600}
                okText="保存"
                cancelText="取消"
            >
                {currentUser && (
                    <>
                        <Descriptions column={2} style={{ marginBottom: 16 }}>
                            <Descriptions.Item label="用户ID">{currentUser.id}</Descriptions.Item>
                            <Descriptions.Item label="用户名">{currentUser.name}</Descriptions.Item>
                            <Descriptions.Item label="邮箱">{currentUser.email}</Descriptions.Item>
                            <Descriptions.Item label="状态">
                                <Tag color={statusMap[currentUser.status].color}>
                                    {statusMap[currentUser.status].text}
                                </Tag>
                            </Descriptions.Item>
                        </Descriptions>
                        <Divider />
                        <Title level={5}>选择角色</Title>
                        {renderRoleCheckboxList(selectedRoleIds, setSelectedRoleIds)}
                    </>
                )}
            </Modal>

            {/* 批量角色分配弹窗 */}
            <Modal
                title={
                    <Space>
                        <TeamOutlined />
                        <span>批量分配角色</span>
                    </Space>
                }
                open={batchAssignVisible}
                onOk={handleBatchAssignRoles}
                onCancel={() => setBatchAssignVisible(false)}
                confirmLoading={batchAssignLoading}
                width={600}
                okText="确认分配"
                cancelText="取消"
            >
                <Alert
                    message={`已选择 ${selectedUserIds.length} 个用户`}
                    description="分配的角色将覆盖这些用户现有的角色配置"
                    type="info"
                    showIcon
                    icon={<InfoCircleOutlined />}
                    style={{ marginBottom: 16 }}
                />
                <Title level={5}>选择要分配的角色</Title>
                {renderRoleCheckboxList(batchSelectedRoleIds, setBatchSelectedRoleIds)}
            </Modal>

            {/* 有效权限查看抽屉 */}
            <Drawer
                title={
                    <Space>
                        <EyeOutlined />
                        <span>用户权限详情 - {currentUser?.name}</span>
                    </Space>
                }
                open={permissionDrawerVisible}
                onClose={() => {
                    setPermissionDrawerVisible(false);
                    setEffectivePermissions(null);
                }}
                width={500}
            >
                {renderEffectivePermissions()}
            </Drawer>
        </PageContainer>
    );
};

export default UserRolePage;
