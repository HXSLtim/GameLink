import React, { useEffect, useState } from 'react';
import { Table, Button, Space, message, Tag, Popconfirm, Card, Modal, Form, Input, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { adminApi } from '@/api/admin';
import type { Permission, ApiResponse } from '@/api/admin';

const PermissionList: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<Permission[]>([]);
    const [total, setTotal] = useState(0);
    const [current, setCurrent] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [form] = Form.useForm();
    const [groups, setGroups] = useState<string[]>([]);

    const fetchData = async (page = current, size = pageSize) => {
        setLoading(true);
        try {
            const res = await adminApi.getPermissions({ page, page_size: size });
            // Handle axios interceptor return type mismatch
            // Runtime: res is ApiResponse
            // TS: res is AxiosResponse<ApiResponse>
            const response = res as unknown as ApiResponse<{ items: Permission[], totalCount: number }>;
            const responseData = response.data;

            if (responseData && responseData.items) {
                setData(responseData.items);
                setTotal(responseData.totalCount || 0);
            } else {
                setData([]);
                setTotal(0);
            }
        } catch (error) {
            message.error('获取权限列表失败');
        } finally {
            setLoading(false);
        }
    };

    const fetchGroups = async () => {
        try {
            const res = await adminApi.getPermissionGroups();
            // @ts-ignore
            const list = res.data?.data || res.data || [];
            setGroups(Array.isArray(list) ? list : []);
        } catch (error) {
            console.error('获取权限分组失败');
        }
    };

    useEffect(() => {
        fetchData();
        fetchGroups();
    }, []);

    const handleDelete = async (id: number) => {
        try {
            await adminApi.deletePermission(id);
            message.success('删除成功');
            fetchData();
        } catch (error) {
            message.error('删除失败');
        }
    };

    const handleEdit = (record: Permission) => {
        setEditingId(record.id);
        form.setFieldsValue(record);
        setIsModalVisible(true);
    };

    const handleCreate = () => {
        setEditingId(null);
        form.resetFields();
        setIsModalVisible(true);
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();
            if (editingId) {
                await adminApi.updatePermission(editingId, values);
                message.success('更新成功');
            } else {
                await adminApi.createPermission(values);
                message.success('创建成功');
            }
            setIsModalVisible(false);
            fetchData();
            fetchGroups();
        } catch (error) {
            message.error('操作失败');
        }
    };

    const handleTableChange = (pagination: any) => {
        setCurrent(pagination.current);
        setPageSize(pagination.pageSize);
        fetchData(pagination.current, pagination.pageSize);
    };

    const columns = [
        {
            title: '权限名称',
            dataIndex: 'name',
            key: 'name',
            width: 200,
            render: (text: string, record: Permission) => text || record.description || record.code,
        },
        {
            title: '请求方法',
            dataIndex: 'method',
            key: 'method',
            width: 100,
            render: (method: string) => {
                let color = 'default';
                if (method === 'GET') color = 'blue';
                if (method === 'POST') color = 'green';
                if (method === 'PUT') color = 'orange';
                if (method === 'DELETE') color = 'red';
                return <Tag color={color}>{method}</Tag>;
            }
        },
        {
            title: '请求路径',
            dataIndex: 'path',
            key: 'path',
            width: 250,
        },
        {
            title: '权限编码',
            dataIndex: 'code',
            key: 'code',
            width: 200,
            render: (code: string) => <Tag color="purple">{code}</Tag>,
        },
        {
            title: '分组',
            dataIndex: 'group',
            key: 'group',
            width: 200,
            filters: groups.map((g: string) => ({ text: g, value: g })),
            onFilter: (value: any, record: Permission) => record.group === value,
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            ellipsis: true,
        },
        {
            title: '操作',
            key: 'action',
            width: 120,
            fixed: 'right' as const,
            render: (_: any, record: Permission) => (
                <Space size="small">
                    <Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} />
                    <Popconfirm title="确定要删除吗？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <Card style={{ border: 'none' }}>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <h2>权限管理</h2>
                <Space>
                    <Button icon={<ReloadOutlined />} onClick={() => fetchData()}>刷新</Button>
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
                        新增权限
                    </Button>
                </Space>
            </div>
            <Table
                columns={columns}
                dataSource={Array.isArray(data) ? data : []}
                rowKey="id"
                loading={loading}
                pagination={{
                    current,
                    pageSize,
                    total,
                    showSizeChanger: true,
                    showTotal: (total) => `共 ${total} 条`,
                }}
                onChange={handleTableChange}
                scroll={{ x: 1300 }}
            />

            <Modal
                title={editingId ? '编辑权限' : '新增权限'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={() => setIsModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="名称"
                        rules={[{ required: true, message: '请输入名称！' }]}
                    >
                        <Input />
                    </Form.Item>
                    <div style={{ display: 'flex', gap: 16 }}>
                        <Form.Item
                            name="method"
                            label="请求方法"
                            style={{ width: 120 }}
                            rules={[{ required: true, message: '请选择请求方法！' }]}
                        >
                            <Select>
                                <Select.Option value="GET">GET</Select.Option>
                                <Select.Option value="POST">POST</Select.Option>
                                <Select.Option value="PUT">PUT</Select.Option>
                                <Select.Option value="DELETE">DELETE</Select.Option>
                                <Select.Option value="PATCH">PATCH</Select.Option>
                            </Select>
                        </Form.Item>
                        <Form.Item
                            name="path"
                            label="请求路径"
                            style={{ flex: 1 }}
                            rules={[{ required: true, message: '请输入请求路径！' }]}
                        >
                            <Input placeholder="/api/v1/..." />
                        </Form.Item>
                    </div>
                    <Form.Item
                        name="code"
                        label="权限编码"
                        rules={[{ required: true, message: '请输入权限编码！' }]}
                    >
                        <Input placeholder="例如 user:create" />
                    </Form.Item>
                    <Form.Item
                        name="group"
                        label="分组"
                        rules={[{ required: true, message: '请输入分组！' }]}
                    >
                        <Select
                            showSearch
                            placeholder="选择一个分组"
                            optionFilterProp="children"
                            allowClear
                        >
                            {groups.map((group: string) => (
                                <Select.Option key={group} value={group}>{group}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item
                        name="description"
                        label="描述"
                    >
                        <Input.TextArea rows={3} />
                    </Form.Item>
                </Form>
            </Modal>
        </Card>
    );
};

export default PermissionList;
