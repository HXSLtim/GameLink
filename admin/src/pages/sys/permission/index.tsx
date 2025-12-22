/**
 * 权限管理页面
 * Requirements: 1.1 - 权限定义与管理
 * 
 * 功能：
 * - 分页、搜索、筛选权限列表
 * - 按分组展示权限
 * - 创建、编辑、删除权限
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Tag,
    Space,
    Button,
    Modal,
    Form,
    Input,
    Select,
    message,
    Typography,
    Tooltip,
    Badge,
    InputNumber,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    EditOutlined,
    DeleteOutlined,
    LockOutlined,
    ApiOutlined,
    ExclamationCircleOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable } from '@/components';
import type { SearchField } from '@/components';
import { PERMISSION_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { permissionApi } from '@/api/permission';
import type { Permission, HTTPMethod, CreatePermissionDto, UpdatePermissionDto } from '@/types/permission';

const { Text } = Typography;

/**
 * HTTP 方法颜色映射
 */
const METHOD_COLORS: Record<HTTPMethod, string> = {
    GET: 'green',
    POST: 'blue',
    PUT: 'orange',
    PATCH: 'purple',
    DELETE: 'red',
    '*': 'default',
};

/**
 * 权限码格式验证正则
 * 格式：module.resource.action（三段式）
 * Requirements: 1.3 - 权限码格式验证
 */
const PERMISSION_CODE_PATTERN = /^[a-z][a-z0-9]*\.[a-z][a-z0-9-]*\.[a-z][a-z0-9-]*$/;

/**
 * 权限管理页面组件
 */
const PermissionPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [permissions, setPermissions] = useState<Permission[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [groups, setGroups] = useState<string[]>([]);
    const [searchParams, setSearchParams] = useState<Record<string, unknown>>({});

    // 弹窗状态
    const [modalVisible, setModalVisible] = useState(false);
    const [currentPermission, setCurrentPermission] = useState<Permission | null>(null);
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [form] = Form.useForm();

    // 删除确认状态
    const [deleteConfirmVisible, setDeleteConfirmVisible] = useState(false);
    const [permissionToDelete, setPermissionToDelete] = useState<Permission | null>(null);
    const [referenceCount, setReferenceCount] = useState(0);

    /**
     * 加载权限数据
     */
    const loadData = useCallback(async (params: Record<string, unknown> = {}) => {
        setLoading(true);
        try {
            const queryParams = {
                page: current,
                page_size: pageSize,
                ...searchParams,
                ...params,
            };
            const res = await permissionApi.list(queryParams);
            if (res.data.success && res.data.data) {
                // Backend returns data as array with pagination object
                const data = res.data.data;
                const pagination = (res.data as { pagination?: { total?: number } }).pagination;
                if (Array.isArray(data)) {
                    setPermissions(data);
                    setTotal(pagination?.total || data.length);
                } else {
                    // Fallback for PaginatedList structure
                    const { items, totalCount } = data as { items: Permission[]; totalCount: number };
                    setPermissions(items || []);
                    setTotal(totalCount || 0);
                }
            }
        } catch (error) {
            console.error(error);
            message.error('加载权限列表失败');
        } finally {
            setLoading(false);
        }
    }, [current, pageSize, searchParams]);

    /**
     * 加载权限分组
     */
    const loadGroups = useCallback(async () => {
        try {
            const res = await permissionApi.getGroups();
            if (res.data.success && res.data.data) {
                setGroups(res.data.data);
            }
        } catch (error) {
            console.error(error);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    useEffect(() => {
        loadGroups();
    }, [loadGroups]);

    /**
     * 搜索处理
     */
    const handleSearch = (values: Record<string, unknown>) => {
        setCurrent(1);
        setSearchParams(values);
    };

    /**
     * 打开新增弹窗
     */
    const handleCreate = () => {
        setCurrentPermission(null);
        form.resetFields();
        form.setFieldsValue({ method: 'GET', sortOrder: 0 });
        setModalVisible(true);
    };

    /**
     * 打开编辑弹窗
     * Requirements: 1.4 - 编辑时权限码不可修改
     */
    const handleEdit = (permission: Permission) => {
        setCurrentPermission(permission);
        form.setFieldsValue({
            method: permission.method,
            path: permission.path,
            code: permission.code,
            group: permission.group,
            description: permission.description,
            sortOrder: permission.sortOrder,
        });
        setModalVisible(true);
    };

    /**
     * 保存权限
     * Requirements: 1.2 - 创建权限
     * Requirements: 1.4 - 编辑权限（权限码不可修改）
     */
    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setConfirmLoading(true);

            // 处理 group 字段（Select tags 模式返回数组）
            const groupValue = Array.isArray(values.group) ? values.group[0] : values.group;

            if (currentPermission) {
                // 编辑模式 - 权限码不可修改
                const updateData: UpdatePermissionDto = {
                    method: values.method,
                    path: values.path,
                    group: groupValue,
                    description: values.description,
                    sortOrder: values.sortOrder,
                };
                await permissionApi.update(currentPermission.id, updateData);
                message.success('更新权限成功');
            } else {
                // 新增模式
                const createData: CreatePermissionDto = {
                    method: values.method,
                    path: values.path,
                    code: values.code,
                    group: groupValue,
                    description: values.description,
                    sortOrder: values.sortOrder,
                };
                await permissionApi.create(createData);
                message.success('创建权限成功');
            }

            setModalVisible(false);
            loadData();
            loadGroups();
        } catch (error: unknown) {
            console.error(error);
            const err = error as { response?: { data?: { message?: string } }; errorFields?: unknown };
            if (err.response?.data?.message) {
                message.error(err.response.data.message);
            } else if (err.errorFields) {
                // Form validation error
                return;
            } else {
                message.error('保存失败');
            }
        } finally {
            setConfirmLoading(false);
        }
    };

    /**
     * 删除权限前检查
     * Requirements: 1.5 - 检查权限是否被角色引用，系统权限不可删除
     */
    const handleDeleteCheck = async (permission: Permission) => {
        if (permission.isSystem) {
            message.error('系统权限不可删除');
            return;
        }
        // 设置待删除权限并显示确认弹窗
        setPermissionToDelete(permission);
        setReferenceCount(0); // 实际应从后端获取引用数量
        setDeleteConfirmVisible(true);
    };

    /**
     * 确认删除权限
     */
    const handleDelete = async () => {
        if (!permissionToDelete) return;
        
        try {
            await permissionApi.delete(permissionToDelete.id);
            message.success(`删除权限 ${permissionToDelete.code} 成功`);
            setDeleteConfirmVisible(false);
            setPermissionToDelete(null);
            loadData();
        } catch (error: unknown) {
            console.error(error);
            const err = error as { response?: { data?: { message?: string } } };
            if (err.response?.data?.message) {
                message.error(err.response.data.message);
            } else {
                message.error('删除失败');
            }
        }
    };

    /**
     * 搜索字段配置
     */
    const searchFields: SearchField[] = [
        {
            name: 'keyword',
            label: '关键词',
            type: 'input',
            placeholder: '权限码/描述',
        },
        {
            name: 'group',
            label: '分组',
            type: 'select',
            placeholder: '选择分组',
            options: groups.map(g => ({ label: g, value: g })),
        },
        {
            name: 'is_system',
            label: '类型',
            type: 'select',
            placeholder: '全部',
            options: [
                { label: '系统权限', value: 'true' },
                { label: '自定义权限', value: 'false' },
            ],
        },
    ];

    /**
     * 表格列配置
     */
    const columns: ColumnsType<Permission> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 70,
        },
        {
            title: '权限码',
            dataIndex: 'code',
            key: 'code',
            width: 220,
            render: (code: string, record) => (
                <Space>
                    {record.isSystem && (
                        <Tooltip title="系统权限">
                            <LockOutlined style={{ color: '#1890ff' }} />
                        </Tooltip>
                    )}
                    <Text code copyable={{ text: code }}>
                        {code}
                    </Text>
                </Space>
            ),
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            width: 200,
            ellipsis: true,
        },
        {
            title: '分组',
            dataIndex: 'group',
            key: 'group',
            width: 120,
            render: (group: string) => (
                <Tag color="processing">{group}</Tag>
            ),
        },
        {
            title: 'HTTP 方法',
            dataIndex: 'method',
            key: 'method',
            width: 100,
            render: (method: HTTPMethod) => (
                <Tag color={METHOD_COLORS[method]}>{method}</Tag>
            ),
        },
        {
            title: 'API 路径',
            dataIndex: 'path',
            key: 'path',
            width: 250,
            ellipsis: true,
            render: (path: string) => (
                <Tooltip title={path}>
                    <Space>
                        <ApiOutlined />
                        <Text type="secondary">{path}</Text>
                    </Space>
                </Tooltip>
            ),
        },
        {
            title: '排序',
            dataIndex: 'sortOrder',
            key: 'sortOrder',
            width: 80,
            align: 'center',
        },
        {
            title: '状态',
            key: 'status',
            width: 100,
            render: (_, record) => (
                record.isSystem ? (
                    <Badge status="processing" text="系统" />
                ) : (
                    <Badge status="success" text="自定义" />
                )
            ),
        },
        {
            title: '操作',
            key: 'action',
            width: 150,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <PermissionGuard permission={PERMISSION_PERMISSIONS.UPDATE}>
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
                        <PermissionGuard permission={PERMISSION_PERMISSIONS.DELETE}>
                            <Button
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                                onClick={() => handleDeleteCheck(record)}
                            >
                                删除
                            </Button>
                        </PermissionGuard>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <PageContainer title="权限管理" subTitle="管理系统权限定义">
            <SearchTable
                columns={columns}
                dataSource={permissions}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={true}
                createText="新增权限"
                createPermission={PERMISSION_PERMISSIONS.CREATE}
                onCreate={handleCreate}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: (t) => `共 ${t} 条`,
                    onChange: (page, size) => {
                        setCurrent(page);
                        setPageSize(size);
                    },
                }}
                scroll={{ x: 1300 }}
            />

            {/* 新增/编辑弹窗 */}
            <Modal
                title={currentPermission ? '编辑权限' : '新增权限'}
                open={modalVisible}
                onOk={handleSave}
                onCancel={() => setModalVisible(false)}
                confirmLoading={confirmLoading}
                width={600}
                destroyOnHidden
            >
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{ method: 'GET', sortOrder: 0 }}
                >
                    <Form.Item
                        name="code"
                        label="权限码"
                        rules={[
                            { required: true, message: '请输入权限码' },
                            {
                                pattern: PERMISSION_CODE_PATTERN,
                                message: '权限码格式无效，应为 module.resource.action（如 admin.users.create）',
                            },
                        ]}
                        tooltip="格式：module.resource.action，创建后不可修改"
                        extra={currentPermission ? '权限码创建后不可修改' : '示例：admin.users.create'}
                    >
                        <Input
                            placeholder="请输入权限码，如 admin.users.create"
                            disabled={!!currentPermission}
                        />
                    </Form.Item>

                    <Form.Item
                        name="description"
                        label="描述"
                        rules={[{ required: true, message: '请输入描述' }]}
                    >
                        <Input placeholder="请输入权限描述" />
                    </Form.Item>

                    <Form.Item
                        name="group"
                        label="分组"
                        rules={[{ required: true, message: '请选择或输入分组' }]}
                    >
                        <Select
                            placeholder="选择或输入分组"
                            mode="tags"
                            maxCount={1}
                            options={groups.map(g => ({ label: g, value: g }))}
                        />
                    </Form.Item>

                    <Space style={{ width: '100%' }} size="large">
                        <Form.Item
                            name="method"
                            label="HTTP 方法"
                            rules={[{ required: true, message: '请选择 HTTP 方法' }]}
                            style={{ width: 150 }}
                        >
                            <Select
                                options={[
                                    { label: 'GET', value: 'GET' },
                                    { label: 'POST', value: 'POST' },
                                    { label: 'PUT', value: 'PUT' },
                                    { label: 'PATCH', value: 'PATCH' },
                                    { label: 'DELETE', value: 'DELETE' },
                                    { label: '* (全部)', value: '*' },
                                ]}
                            />
                        </Form.Item>

                        <Form.Item
                            name="path"
                            label="API 路径"
                            rules={[{ required: true, message: '请输入 API 路径' }]}
                            style={{ flex: 1 }}
                        >
                            <Input placeholder="请输入 API 路径，如 /api/admin/users" />
                        </Form.Item>
                    </Space>

                    <Form.Item
                        name="sortOrder"
                        label="排序"
                        tooltip="数值越小越靠前"
                    >
                        <InputNumber placeholder="排序值" style={{ width: 150 }} />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 删除确认弹窗 - Requirements: 1.5 */}
            <Modal
                title={
                    <Space>
                        <ExclamationCircleOutlined style={{ color: '#faad14' }} />
                        确认删除权限
                    </Space>
                }
                open={deleteConfirmVisible}
                onOk={handleDelete}
                onCancel={() => {
                    setDeleteConfirmVisible(false);
                    setPermissionToDelete(null);
                }}
                okText="确认删除"
                okButtonProps={{ danger: true }}
            >
                <p>确定要删除权限 <Text strong>{permissionToDelete?.code}</Text> 吗？</p>
                {referenceCount > 0 && (
                    <p style={{ color: '#faad14' }}>
                        <ExclamationCircleOutlined /> 该权限被 {referenceCount} 个角色引用，删除后这些角色将失去此权限。
                    </p>
                )}
                <p>此操作不可恢复。</p>
            </Modal>
        </PageContainer>
    );
};

export default PermissionPage;
