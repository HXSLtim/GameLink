import React, { useState, useEffect } from 'react';
import { Card, Table, Button, Space, Modal, Form, Input, ColorPicker, Popconfirm, Tag, App } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageContainer } from '@/components';
import { adminApi, type UserTag, type CreateTagDto, type ApiResponse } from '@/api/admin';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
const UserTags: React.FC = () => {
    const { message } = App.useApp();
    const [loading, setLoading] = useState(false);
    const [tags, setTags] = useState<UserTag[]>([]);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingTag, setEditingTag] = useState<UserTag | null>(null);
    const [form] = Form.useForm();

    const fetchTags = async () => {
        setLoading(true);
        try {
            const res = await adminApi.getTags() as unknown as ApiResponse<UserTag[]>;
            if (res.success) {
                setTags(res.data || []);
            }
        } catch {
            message.error('获取标签列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTags();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleCreate = () => {
        setEditingTag(null);
        form.resetFields();
        // Set default color
        form.setFieldsValue({ color: '#1677ff' });
        setIsModalVisible(true);
    };

    const handleEdit = (tag: UserTag) => {
        setEditingTag(tag);
        form.setFieldsValue({
            name: tag.name,
            color: tag.color,
            description: tag.description,
        });
        setIsModalVisible(true);
    };

    const handleDelete = async (id: number) => {
        try {
            const res = await adminApi.deleteTag(id) as unknown as ApiResponse<void>;
            if (res.success) {
                message.success('删除成功');
                fetchTags();
            }
        } catch {
            message.error('删除失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();
            // Convert ColorPicker value to hex string if needed
            const color = typeof values.color === 'string' ? values.color : values.color.toHexString();

            const data: CreateTagDto = {
                name: values.name,
                color: color,
                description: values.description,
            };

            if (editingTag) {
                const res = await adminApi.updateTag(editingTag.id, data) as unknown as ApiResponse<UserTag>;
                if (res.success) {
                    message.success('更新成功');
                    setIsModalVisible(false);
                    fetchTags();
                }
            } else {
                const res = await adminApi.createTag(data) as unknown as ApiResponse<UserTag>;
                if (res.success) {
                    message.success('创建成功');
                    setIsModalVisible(false);
                    fetchTags();
                }
            }
        } catch (error) {
            logger.error("Operation failed", error);
            message.error('操作失败');
        }
    };

    const columns = [
        {
            title: '标签名称',
            dataIndex: 'name',
            key: 'name',
            render: (text: string, record: UserTag) => (
                <Tag color={record.color}>{text}</Tag>
            ),
        },
        {
            title: '颜色值',
            dataIndex: 'color',
            key: 'color',
            render: (color: string) => (
                <Space>
                    <div style={{ width: 16, height: 16, backgroundColor: color, borderRadius: 2, border: '1px solid #d9d9d9' }} />
                    {color}
                </Space>
            ),
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: '用户数',
            dataIndex: 'userCount',
            key: 'userCount',
            render: (count: number) => count || 0,
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 160,
            fixed: 'right',
            render: (_: unknown, record: UserTag) => (
                <Space size={4}>
                    <Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} />
                    <Popconfirm title="确定要删除吗？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="text" danger icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <PageContainer title="用户标签管理" subTitle="管理用户标签体系">
            <Card>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                    <Space>
                        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
                            新增标签
                        </Button>
                    </Space>
                    <Button icon={<ReloadOutlined />} onClick={fetchTags}>刷新</Button>
                </div>
                <Table<UserTag>
                    columns={columns as ColumnsType<UserTag>}
                    dataSource={tags}
                    rowKey="id"
                    loading={loading}
                    scroll={{ x: 1000 }}
                />
            </Card>

            <Modal
                title={editingTag ? '编辑标签' : '新增标签'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={() => setIsModalVisible(false)}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="标签名称"
                        rules={[{ required: true, message: '请输入标签名称' }]}
                    >
                        <Input placeholder="请输入标签名称" />
                    </Form.Item>
                    <Form.Item
                        name="color"
                        label="标签颜色"
                        rules={[{ required: true, message: '请选择标签颜色' }]}
                    >
                        <ColorPicker showText />
                    </Form.Item>
                    <Form.Item
                        name="description"
                        label="描述"
                    >
                        <Input.TextArea rows={3} placeholder="请输入描述" />
                    </Form.Item>
                </Form>
            </Modal>
        </PageContainer>
    );
};

export default UserTags;
