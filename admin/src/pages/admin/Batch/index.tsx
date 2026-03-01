/**
 * 用户批量操作页面
 * 支持批量启用/禁用、发送通知、添加积分、修改角色等操作
 */
import React, { useState, useCallback, useMemo } from 'react';
import {
    Card,
    Form,
    Input,
    Select,
    Button,
    message,
    Space,
    Typography,
    Alert,
    Radio,
    InputNumber,
    Row,
    Col,
} from 'antd';
import {
    SendOutlined,
    PlusOutlined,
    TeamOutlined,
    LockOutlined,
    UnlockOutlined,
    CheckCircleOutlined,
    UserOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@/components';
import { SearchFilters } from '@/components/common/SearchFilters';
import BatchActions from '@/components/common/BatchActions';
import { adminApi } from '@/api/admin';
import type { BatchRoleDto, BatchStatusDto, BatchPointsDto, BatchNotificationDto } from '@/api/admin';
import { USER_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { logger } from '@/utils/logger';

const { Text } = Typography;
const { TextArea } = Input;

/**
 * 角色选项
 */
const ROLE_OPTIONS = [
    { label: '普通用户', value: 'user' },
    { label: '陪玩师', value: 'player' },
    { label: '管理员', value: 'admin' },
];

/**
 * 状态选项
 */
const STATUS_OPTIONS = [
    { label: '正常', value: 'active' },
    { label: '禁用', value: 'inactive' },
    { label: '锁定', value: 'locked' },
];

/**
 * 积分类型选项
 */
const POINTS_TYPE_OPTIONS = [
    { label: '活动奖励', value: 'activity' },
    { label: '系统赠送', value: 'system' },
    { label: '注册奖励', value: 'register' },
    { label: '签到奖励', value: 'checkin' },
    { label: '任务奖励', value: 'task' },
    { label: '补偿', value: 'compensation' },
];

/**
 * 通知类型选项
 */
const NOTIFICATION_TYPE_OPTIONS = [
    { label: '系统通知', value: 'system' },
    { label: '营销通知', value: 'marketing' },
    { label: '个人通知', value: 'personal' },
    { label: '活动通知', value: 'activity' },
];

/**
 * 批量操作表单值类型
 */
interface BatchNotificationFormValues {
    target: 'users' | 'role' | 'all';
    roles?: string[];
    title: string;
    content: string;
    type: BatchNotificationDto['type'];
}

interface BatchPointsFormValues {
    target: 'users' | 'role' | 'all';
    roles?: string[];
    cents: number;
    reason: string;
    type: string;
}

/**
 * 批量操作页面
 */
const BatchPage: React.FC = () => {
    // 表单实例
    const [notificationForm] = Form.useForm<BatchNotificationFormValues>();
    const [pointsForm] = Form.useForm<BatchPointsFormValues>();
    const [roleForm] = Form.useForm<{ role: string }>();
    const [statusForm] = Form.useForm<{ status: string }>();

    // 状态
    const [selectedUserIds, setSelectedUserIds] = useState<number[]>([]);
    const [filterValues, setFilterValues] = useState<Record<string, unknown>>({});
    const [loading, setLoading] = useState(false);

    /**
     * 筛选器配置
     */
    const filters = useMemo(() => [
        {
            key: 'keyword',
            type: 'input' as const,
            placeholder: '用户名/邮箱/手机号',
            allowClear: true,
            width: 200,
        },
        {
            key: 'role',
            type: 'select' as const,
            placeholder: '角色',
            options: [
                { label: '普通用户', value: 'user' },
                { label: '陪玩师', value: 'player' },
                { label: '管理员', value: 'admin' },
            ],
            allowClear: true,
            width: 120,
        },
        {
            key: 'status',
            type: 'select' as const,
            placeholder: '状态',
            options: [
                { label: '正常', value: 'active' },
                { label: '禁用', value: 'inactive' },
                { label: '锁定', value: 'locked' },
            ],
            allowClear: true,
            width: 120,
        },
    ], []);

    /**
     * 处理筛选变化
     */
    const handleFilterChange = useCallback((key: string, value: unknown) => {
        setFilterValues(prev => ({ ...prev, [key]: value }));
    }, []);

    /**
     * 查询用户（获取选中用户数量）
     */
    const handleQueryUsers = useCallback(async () => {
        try {
            setLoading(true);
            const params: Record<string, unknown> = {
                page: 1,
                page_size: 1, // 只获取数量
                ...filterValues,
            };

            const response = await adminApi.getUsers(params) as unknown as { data: { data: { total: number } } };
            const total = response.data.data?.total || 0;

            message.info(`符合条件的用户数量: ${total}`);
            return total;
        } catch (error) {
            logger.error('Failed to query users', error);
            message.error('查询用户失败');
            return 0;
        } finally {
            setLoading(false);
        }
    }, [filterValues]);

    /**
     * 重置筛选
     */
    const handleResetFilters = useCallback(() => {
        setFilterValues({});
        notificationForm.resetFields();
        pointsForm.resetFields();
        roleForm.resetFields();
        statusForm.resetFields();
    }, [notificationForm, pointsForm, roleForm, statusForm]);

    /**
     * 批量发送通知
     */
    const handleSendNotification = useCallback(async () => {
        try {
            const values = await notificationForm.validateFields();

            if (!values.title || !values.content) {
                message.warning('请填写通知标题和内容');
                return;
            }

            setLoading(true);

            const payload: BatchNotificationDto = {
                target: values.target,
                title: values.title,
                content: values.content,
                type: values.type,
            };

            if (values.target === 'users' && selectedUserIds.length > 0) {
                payload.userIds = selectedUserIds;
            } else if (values.target === 'role' && values.roles) {
                payload.roles = values.roles;
            }

            const response = await adminApi.batchSendNotification(payload) as unknown as { success: boolean };

            if (response.success) {
                message.success('批量发送通知成功');
                notificationForm.resetFields();
            } else {
                message.error('发送通知失败');
            }
        } catch (error) {
            logger.error('Failed to send notification', error);
            message.error('发送通知失败');
        } finally {
            setLoading(false);
        }
    }, [notificationForm, selectedUserIds]);

    /**
     * 批量添加积分
     */
    const handleAddPoints = useCallback(async () => {
        try {
            const values = await pointsForm.validateFields();

            if (!values.cents || values.cents <= 0) {
                message.warning('请输入有效的积分数量');
                return;
            }

            if (!values.reason) {
                message.warning('请输入积分变动原因');
                return;
            }

            setLoading(true);

            const payload: BatchPointsDto = {
                target: values.target,
                cents: values.cents,
                reason: values.reason,
                type: values.type,
            };

            if (values.target === 'users' && selectedUserIds.length > 0) {
                payload.userIds = selectedUserIds;
            } else if (values.target === 'role' && values.roles) {
                payload.roles = values.roles;
            }

            const response = await adminApi.batchAddUserPoints(payload) as unknown as { success: boolean };

            if (response.success) {
                message.success('批量添加积分成功');
                pointsForm.resetFields();
            } else {
                message.error('添加积分失败');
            }
        } catch (error) {
            logger.error('Failed to add points', error);
            message.error('添加积分失败');
        } finally {
            setLoading(false);
        }
    }, [pointsForm, selectedUserIds]);

    /**
     * 批量修改角色
     */
    const handleUpdateRole = useCallback(async () => {
        try {
            const values = await roleForm.validateFields();

            if (!values.role) {
                message.warning('请选择角色');
                return;
            }

            if (selectedUserIds.length === 0) {
                message.warning('请先选择用户或使用筛选条件');
                return;
            }

            setLoading(true);

            const payload: BatchRoleDto = {
                userIds: selectedUserIds,
                role: values.role,
            };

            const response = await adminApi.batchUpdateUserRole(payload) as unknown as { success: boolean };

            if (response.success) {
                message.success('批量修改角色成功');
                roleForm.resetFields();
            } else {
                message.error('修改角色失败');
            }
        } catch (error) {
            logger.error('Failed to update role', error);
            message.error('修改角色失败');
        } finally {
            setLoading(false);
        }
    }, [roleForm, selectedUserIds]);

    /**
     * 批量修改状态
     */
    const handleUpdateStatus = useCallback(async () => {
        try {
            const values = await statusForm.validateFields();

            if (!values.status) {
                message.warning('请选择状态');
                return;
            }

            if (selectedUserIds.length === 0) {
                message.warning('请先选择用户或使用筛选条件');
                return;
            }

            setLoading(true);

            const payload: BatchStatusDto = {
                userIds: selectedUserIds,
                status: values.status,
            };

            const response = await adminApi.batchUpdateUserStatus(payload) as unknown as { success: boolean };

            if (response.success) {
                message.success('批量修改状态成功');
                statusForm.resetFields();
            } else {
                message.error('修改状态失败');
            }
        } catch (error) {
            logger.error('Failed to update status', error);
            message.error('修改状态失败');
        } finally {
            setLoading(false);
        }
    }, [statusForm, selectedUserIds]);

    /**
     * 批量操作配置
     */
    const batchActions = useMemo(() => [
        {
            key: 'enable',
            label: '批量启用',
            icon: <UnlockOutlined />,
            type: 'default' as const,
            onConfirm: async () => {
                if (selectedUserIds.length === 0) {
                    message.warning('请先选择用户');
                    return;
                }
                const payload: BatchStatusDto = {
                    userIds: selectedUserIds,
                    status: 'active',
                };
                await adminApi.batchUpdateUserStatus(payload);
                message.success('批量启用成功');
            },
        },
        {
            key: 'disable',
            label: '批量禁用',
            icon: <LockOutlined />,
            danger: true,
            onConfirm: async () => {
                if (selectedUserIds.length === 0) {
                    message.warning('请先选择用户');
                    return;
                }
                const payload: BatchStatusDto = {
                    userIds: selectedUserIds,
                    status: 'inactive',
                };
                await adminApi.batchUpdateUserStatus(payload);
                message.success('批量禁用成功');
            },
        },
        {
            key: 'delete',
            label: '批量删除',
            icon: <TeamOutlined />,
            danger: true,
            onConfirm: async () => {
                if (selectedUserIds.length === 0) {
                    message.warning('请先选择用户');
                    return;
                }
                await adminApi.batchDeleteUsers(selectedUserIds);
                message.success('批量删除成功');
                setSelectedUserIds([]);
            },
        },
    ], [selectedUserIds]);

    return (
        <PageContainer
            title="用户批量操作"
            subTitle="支持批量启用/禁用、发送通知、添加积分、修改角色等操作"
        >
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
                {/* 筛选区域 */}
                <SearchFilters
                    filters={filters}
                    onFilterChange={handleFilterChange}
                    filterValues={filterValues}
                    showQueryButtons
                    onQuery={handleQueryUsers}
                    onReset={handleResetFilters}
                    loading={loading}
                />

                {/* 已选用户提示 */}
                {selectedUserIds.length > 0 && (
                    <Alert
                        message={
                            <Space>
                                <CheckCircleOutlined />
                                <Text>已选择 <Text strong>{selectedUserIds.length}</Text> 个用户进行操作</Text>
                            </Space>
                        }
                        type="info"
                        showIcon
                        closable
                        onClose={() => setSelectedUserIds([])}
                    />
                )}

                {/* 快速操作 */}
                <Card title="快速操作" extra={<Text type="secondary">对已选用户或筛选结果进行操作</Text>}>
                    <BatchActions
                        selectedCount={selectedUserIds.length}
                        actions={batchActions}
                        selectedRowKeys={selectedUserIds}
                    />
                </Card>

                {/* 批量发送通知 */}
                <PermissionGuard permission={USER_PERMISSIONS.UPDATE}>
                    <Card
                        title={
                            <Space>
                                <SendOutlined />
                                <span>批量发送通知</span>
                            </Space>
                        }
                    >
                        <Form
                            form={notificationForm}
                            layout="vertical"
                            initialValues={{
                                target: 'all',
                                type: 'system',
                            }}
                        >
                            <Row gutter={16}>
                                <Col span={8}>
                                    <Form.Item
                                        name="target"
                                        label="发送目标"
                                        rules={[{ required: true, message: '请选择发送目标' }]}
                                    >
                                        <Radio.Group>
                                            <Radio value="all">全部用户</Radio>
                                            <Radio value="role">按角色</Radio>
                                            <Radio value="users">指定用户</Radio>
                                        </Radio.Group>
                                    </Form.Item>
                                </Col>
                                <Col span={8}>
                                    <Form.Item noStyle shouldUpdate={(prev, curr) => prev.target !== curr.target}>
                                        {({ getFieldValue }) =>
                                            getFieldValue('target') === 'role' ? (
                                                <Form.Item
                                                    name="roles"
                                                    label="选择角色"
                                                    rules={[{ required: true, message: '请选择角色' }]}
                                                >
                                                    <Select mode="multiple" placeholder="请选择角色" options={ROLE_OPTIONS} />
                                                </Form.Item>
                                            ) : null
                                        }
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item
                                name="type"
                                label="通知类型"
                                rules={[{ required: true, message: '请选择通知类型' }]}
                            >
                                <Select options={NOTIFICATION_TYPE_OPTIONS} />
                            </Form.Item>

                            <Form.Item
                                name="title"
                                label="通知标题"
                                rules={[{ required: true, message: '请输入通知标题' }]}
                            >
                                <Input placeholder="请输入通知标题" maxLength={100} showCount />
                            </Form.Item>

                            <Form.Item
                                name="content"
                                label="通知内容"
                                rules={[{ required: true, message: '请输入通知内容' }]}
                            >
                                <TextArea
                                    placeholder="请输入通知内容"
                                    rows={4}
                                    maxLength={500}
                                    showCount
                                />
                            </Form.Item>

                            <Form.Item>
                                <Button
                                    type="primary"
                                    icon={<SendOutlined />}
                                    onClick={handleSendNotification}
                                    loading={loading}
                                >
                                    发送通知
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </PermissionGuard>

                {/* 批量添加积分 */}
                <PermissionGuard permission={USER_PERMISSIONS.UPDATE}>
                    <Card
                        title={
                            <Space>
                                <PlusOutlined />
                                <span>批量添加积分</span>
                            </Space>
                        }
                    >
                        <Form
                            form={pointsForm}
                            layout="vertical"
                            initialValues={{
                                target: 'all',
                                type: 'activity',
                            }}
                        >
                            <Row gutter={16}>
                                <Col span={8}>
                                    <Form.Item
                                        name="target"
                                        label="添加目标"
                                        rules={[{ required: true, message: '请选择添加目标' }]}
                                    >
                                        <Radio.Group>
                                            <Radio value="all">全部用户</Radio>
                                            <Radio value="role">按角色</Radio>
                                            <Radio value="users">指定用户</Radio>
                                        </Radio.Group>
                                    </Form.Item>
                                </Col>
                                <Col span={8}>
                                    <Form.Item noStyle shouldUpdate={(prev, curr) => prev.target !== curr.target}>
                                        {({ getFieldValue }) =>
                                            getFieldValue('target') === 'role' ? (
                                                <Form.Item
                                                    name="roles"
                                                    label="选择角色"
                                                    rules={[{ required: true, message: '请选择角色' }]}
                                                >
                                                    <Select mode="multiple" placeholder="请选择角色" options={ROLE_OPTIONS} />
                                                </Form.Item>
                                            ) : null
                                        }
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Row gutter={16}>
                                <Col span={8}>
                                    <Form.Item
                                        name="type"
                                        label="积分类型"
                                        rules={[{ required: true, message: '请选择积分类型' }]}
                                    >
                                        <Select options={POINTS_TYPE_OPTIONS} />
                                    </Form.Item>
                                </Col>
                                <Col span={8}>
                                    <Form.Item
                                        name="cents"
                                        label="积分数量（分）"
                                        rules={[
                                            { required: true, message: '请输入积分数量' },
                                            { type: 'number', min: 1, message: '积分数量必须大于0' },
                                        ]}
                                        tooltip="100分 = 1元"
                                    >
                                        <InputNumber
                                            style={{ width: '100%' }}
                                            placeholder="请输入积分数量"
                                            min={1}
                                            precision={0}
                                        />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item
                                name="reason"
                                label="变动原因"
                                rules={[{ required: true, message: '请输入积分变动原因' }]}
                            >
                                <Input placeholder="请输入积分变动原因" maxLength={200} />
                            </Form.Item>

                            <Form.Item>
                                <Button
                                    type="primary"
                                    icon={<PlusOutlined />}
                                    onClick={handleAddPoints}
                                    loading={loading}
                                >
                                    添加积分
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </PermissionGuard>

                {/* 批量修改角色 */}
                <PermissionGuard permission={USER_PERMISSIONS.UPDATE}>
                    <Card
                        title={
                            <Space>
                                <UserOutlined />
                                <span>批量修改角色</span>
                            </Space>
                        }
                    >
                        <Alert
                            message="此操作将修改选中用户的角色"
                            description="请先通过筛选条件选择用户，或使用指定用户ID"
                            type="warning"
                            showIcon
                            style={{ marginBottom: 16 }}
                        />
                        <Form
                            form={roleForm}
                            layout="vertical"
                        >
                            <Form.Item
                                name="role"
                                label="选择角色"
                                rules={[{ required: true, message: '请选择角色' }]}
                            >
                                <Select placeholder="请选择角色" options={ROLE_OPTIONS} />
                            </Form.Item>

                            <Form.Item>
                                <Button
                                    type="primary"
                                    icon={<TeamOutlined />}
                                    onClick={handleUpdateRole}
                                    loading={loading}
                                >
                                    修改角色
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </PermissionGuard>

                {/* 批量修改状态 */}
                <PermissionGuard permission={USER_PERMISSIONS.STATUS}>
                    <Card
                        title={
                            <Space>
                                <LockOutlined />
                                <span>批量修改状态</span>
                            </Space>
                        }
                    >
                        <Alert
                            message="此操作将修改选中用户的状态"
                            description="请先通过筛选条件选择用户，或使用指定用户ID"
                            type="warning"
                            showIcon
                            style={{ marginBottom: 16 }}
                        />
                        <Form
                            form={statusForm}
                            layout="vertical"
                        >
                            <Form.Item
                                name="status"
                                label="选择状态"
                                rules={[{ required: true, message: '请选择状态' }]}
                            >
                                <Select placeholder="请选择状态" options={STATUS_OPTIONS} />
                            </Form.Item>

                            <Form.Item>
                                <Space>
                                    <Button
                                        type="primary"
                                        icon={<UnlockOutlined />}
                                        onClick={() => statusForm.setFieldsValue({ status: 'active' })}
                                    >
                                        启用
                                    </Button>
                                    <Button
                                        danger
                                        icon={<LockOutlined />}
                                        onClick={() => statusForm.setFieldsValue({ status: 'inactive' })}
                                    >
                                        禁用
                                    </Button>
                                    <Button
                                        type="default"
                                        onClick={handleUpdateStatus}
                                        loading={loading}
                                    >
                                        确认修改
                                    </Button>
                                </Space>
                            </Form.Item>
                        </Form>
                    </Card>
                </PermissionGuard>
            </Space>
        </PageContainer>
    );
};

export default BatchPage;
