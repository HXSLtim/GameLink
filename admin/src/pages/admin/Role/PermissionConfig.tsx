/**
 * 角色权限配置页面
 * Requirements: 2.1, 2.2, 2.4, 2.5
 *
 * 功能：
 * - 以树形结构展示所有权限（按分组懒加载）
 * - 支持父子节点联动选择
 * - 高亮已分配权限
 * - 系统角色显示特殊提示
 */
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import {
    Card,
    Button,
    Space,
    Spin,
    Alert,
    Typography,
    Descriptions,
    Tag,
    Divider,
    Modal,
    App,
} from 'antd';
import {
    SaveOutlined,
    RollbackOutlined,
    SafetyCertificateOutlined,
    LockOutlined,
    ExclamationCircleOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import { PageContainer } from '@/components';
import { PermissionTree } from '@/components/PermissionTree';
import { roleApi } from '@/api/permission';
import { ROLE_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import type { Role } from '@/types/permission';

import { logger } from '@/utils/logger';
const { Text } = Typography;

/**
 * 角色权限配置页面组件
 */
const RolePermissionConfig: React.FC = () => {
    const { message, modal } = App.useApp();
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const roleId = id ? parseInt(id, 10) : null;

    // 状态
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [role, setRole] = useState<Role | null>(null);
    const [checkedKeys, setCheckedKeys] = useState<number[]>([]);
    const [originalCheckedKeys, setOriginalCheckedKeys] = useState<number[]>([]);

    /**
     * 加载角色信息
     */
    const loadRole = useCallback(async () => {
        if (!roleId) return;

        try {
            const res = await roleApi.get(roleId);
            if (res.data.success && res.data.data) {
                setRole(res.data.data);
            }
        } catch (error) {
            logger.error('Failed to load role:', error);
            message.error('加载角色信息失败');
        }
    }, [roleId, message]);

    /**
     * 加载角色已有权限
     */
    const loadRolePermissions = useCallback(async () => {
        if (!roleId) return;

        try {
            const res = await roleApi.getPermissions(roleId);
            if (res.data.success && res.data.data) {
                const permissionIds = res.data.data;
                setCheckedKeys(permissionIds);
                setOriginalCheckedKeys(permissionIds);
            }
        } catch (error) {
            logger.error('Failed to load role permissions:', error);
            message.error('加载角色权限失败');
        }
    }, [roleId, message]);

    /**
     * 初始化加载
     */
    useEffect(() => {
        const loadData = async () => {
            setLoading(true);
            await Promise.all([loadRole(), loadRolePermissions()]);
            setLoading(false);
        };

        if (roleId) {
            loadData();
        }
    }, [roleId, loadRole, loadRolePermissions]);

    /**
     * 处理权限选择变化
     * Requirements: 2.2 - 支持批量选择（选中父节点自动选中所有子节点）
     */
    const handleCheck = useCallback((keys: number[]) => {
        setCheckedKeys(keys);
    }, []);

    /**
     * 检查是否有未保存的更改
     */
    const hasChanges = useMemo(() => {
        if (checkedKeys.length !== originalCheckedKeys.length) return true;
        const sortedCurrent = [...checkedKeys].sort((a, b) => a - b);
        const sortedOriginal = [...originalCheckedKeys].sort((a, b) => a - b);
        return !sortedCurrent.every((key, index) => key === sortedOriginal[index]);
    }, [checkedKeys, originalCheckedKeys]);

    /**
     * 保存权限配置
     * Requirements: 2.3 - 保存角色权限并立即生效
     */
    const handleSave = useCallback(async () => {
        if (!roleId) return;

        // 如果是超级管理员角色，显示提示
        if (role?.slug === 'superAdmin') {
            Modal.info({
                title: '提示',
                content: '超级管理员默认拥有所有权限，无需单独配置。',
                icon: <ExclamationCircleOutlined />,
            });
            return;
        }

        setSaving(true);
        try {
            await roleApi.batchAssignPermissions(roleId, { permissionIds: checkedKeys });
            message.success('权限配置保存成功');
            setOriginalCheckedKeys(checkedKeys);
        } catch (error: unknown) {
            logger.error('Failed to save permissions:', error);
            const err = error as { response?: { data?: { message?: string } } };
            if (err.response?.data?.message) {
                message.error(err.response.data.message);
            } else {
                message.error('保存权限配置失败');
            }
        } finally {
            setSaving(false);
        }
    }, [roleId, role, checkedKeys, message]);

    /**
     * 重置为原始状态
     */
    const handleReset = useCallback(() => {
        setCheckedKeys(originalCheckedKeys);
        message.info('已重置为原始配置');
    }, [originalCheckedKeys, message]);

    /**
     * 返回角色列表
     */
    const handleBack = useCallback(() => {
        if (hasChanges) {
            modal.confirm({
                title: '确认离开',
                content: '您有未保存的更改，确定要离开吗？',
                icon: <ExclamationCircleOutlined />,
                okText: '离开',
                cancelText: '取消',
                onOk: () => navigate('/admin/sys/role'),
            });
        } else {
            navigate('/admin/sys/role');
        }
    }, [hasChanges, navigate, modal]);

    /**
     * 判断是否为系统角色（超级管理员）
     * Requirements: 2.5 - 系统角色显示特殊提示
     */
    const isSuperAdmin = role?.slug === 'superAdmin';

    if (loading) {
        return (
            <PageContainer title="角色权限配置">
                <div style={{ textAlign: 'center', padding: '100px 0' }}>
                    <Spin size="large" tip="加载中..." />
                </div>
            </PageContainer>
        );
    }

    if (!role) {
        return (
            <PageContainer title="角色权限配置">
                <Alert
                    type="error"
                    message="角色不存在"
                    description="未找到指定的角色信息，请返回角色列表重试。"
                    showIcon
                    action={
                        <Button onClick={() => navigate('/admin/sys/role')}>
                            返回角色列表
                        </Button>
                    }
                />
            </PageContainer>
        );
    }

    return (
        <PageContainer
            title={`角色权限配置 - ${role.name}`}
            subTitle="配置角色的访问权限"
            extra={
                <Space>
                    <Button icon={<RollbackOutlined />} onClick={handleBack}>
                        返回
                    </Button>
                    <PermissionGuard permission={ROLE_PERMISSIONS.ASSIGN_PERMISSIONS}>
                        <Button
                            onClick={handleReset}
                            disabled={!hasChanges || saving || isSuperAdmin}
                        >
                            重置
                        </Button>
                        <Button
                            type="primary"
                            icon={<SaveOutlined />}
                            onClick={handleSave}
                            loading={saving}
                            disabled={!hasChanges || isSuperAdmin}
                        >
                            保存配置
                        </Button>
                    </PermissionGuard>
                </Space>
            }
        >
            {/* 角色信息卡片 */}
            <Card style={{ marginBottom: 16 }}>
                <Descriptions
                    title={
                        <Space>
                            <SafetyCertificateOutlined
                                style={{ color: role.isSystem ? '#1890ff' : '#52c41a' }}
                            />
                            <span>角色信息</span>
                            {role.isSystem && <Tag color="blue">系统角色</Tag>}
                        </Space>
                    }
                    column={3}
                >
                    <Descriptions.Item label="角色名称">{role.name}</Descriptions.Item>
                    <Descriptions.Item label="角色编码">
                        <Text code>{role.slug}</Text>
                    </Descriptions.Item>
                    <Descriptions.Item label="优先级">{role.priority}</Descriptions.Item>
                    <Descriptions.Item label="描述" span={3}>
                        {role.description || '-'}
                    </Descriptions.Item>
                </Descriptions>
            </Card>

            {/* 超级管理员提示 - Requirements: 2.5 */}
            {isSuperAdmin && (
                <Alert
                    type="info"
                    message="超级管理员权限说明"
                    description={
                        <div>
                            <p style={{ margin: '0 0 8px 0' }}>
                                <LockOutlined style={{ marginRight: 8 }} />
                                超级管理员（superAdmin）默认拥有系统所有权限，无需单独配置权限。
                            </p>
                            <ul style={{ margin: 0, paddingLeft: 20 }}>
                                <li>系统会自动跳过权限检查，直接放行所有操作</li>
                                <li>权限配置面板仅供查看，无法进行修改</li>
                                <li>如需限制超级管理员权限，请创建新的自定义角色</li>
                            </ul>
                        </div>
                    }
                    showIcon
                    style={{ marginBottom: 16 }}
                />
            )}

            {/* 系统角色提示（非超级管理员的系统角色） - Requirements: 2.5 */}
            {role.isSystem && !isSuperAdmin && (
                <Alert
                    type="warning"
                    message="系统角色提示"
                    description={
                        <span>
                            <SafetyCertificateOutlined style={{ marginRight: 8 }} />
                            这是一个系统预设角色，修改其权限可能影响系统正常运行。请谨慎操作。
                        </span>
                    }
                    showIcon
                    style={{ marginBottom: 16 }}
                />
            )}

            <Card
                title={
                    <Space>
                        <span>权限配置</span>
                        {hasChanges && (
                            <Tag color="warning">有未保存的更改</Tag>
                        )}
                    </Space>
                }
            >
                <PermissionTree
                    checkedKeys={checkedKeys}
                    onCheck={handleCheck}
                    loading={loading}
                    disabled={isSuperAdmin}
                    showSearch
                    showSelectAll
                    height={500}
                    virtual
                    isSystemRole={isSuperAdmin}
                    lazyLoadByGroup
                />

                <Divider />

                {/* 底部操作栏 */}
                <div style={{ textAlign: 'right' }}>
                    <Space>
                        <Text type="secondary">
                            已选择 {checkedKeys.length} 项权限
                        </Text>
                        <PermissionGuard permission={ROLE_PERMISSIONS.ASSIGN_PERMISSIONS}>
                            <Button
                                onClick={handleReset}
                                disabled={!hasChanges || saving || isSuperAdmin}
                            >
                                重置
                            </Button>
                            <Button
                                type="primary"
                                icon={<SaveOutlined />}
                                onClick={handleSave}
                                loading={saving}
                                disabled={!hasChanges || isSuperAdmin}
                            >
                                保存配置
                            </Button>
                        </PermissionGuard>
                    </Space>
                </div>
            </Card>
        </PageContainer>
    );
};

export default RolePermissionConfig;
