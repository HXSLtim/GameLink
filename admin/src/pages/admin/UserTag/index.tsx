/**
 * 用户标签管理页面
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
    ColorPicker,
    Table,
    Card,
    Row,
    Col,
    Statistic,
    Avatar,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { Color } from 'antd/es/color-picker';
import {
    EditOutlined,
    DeleteOutlined,
    TagsOutlined,
    UserOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton } from '@/components';
import type { SearchField } from '@/components';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';
import apiClient from '@/api/client';

import { logger } from '@/utils/logger';
/**
 * 标签数据接口
 */
interface UserTag {
    id: number;
    name: string;
    color: string;
    description: string;
    userCount?: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * 用户数据接口
 */
interface TagUser {
    id: number;
    name: string;
    email: string;
    avatarUrl?: string;
}

/**
 * 导出列配置
 */
const tagExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '标签名称' },
    { key: 'color', title: '颜色' },
    { key: 'description', title: '描述' },
    { key: 'userCount', title: '用户数' },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

const UserTagPage: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [tags, setTags] = useState<UserTag[]>([]);

    // 统计
    const [stats, setStats] = useState({ totalTags: 0, totalTaggedUsers: 0 });

    // 弹窗状态
    const [editVisible, setEditVisible] = useState(false);
    const [usersVisible, setUsersVisible] = useState(false);
    const [currentTag, setCurrentTag] = useState<UserTag | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    // 标签用户列表
    const [tagUsers, setTagUsers] = useState<TagUser[]>([]);
    const [tagUsersLoading, setTagUsersLoading] = useState(false);
    const [tagUsersTotal, setTagUsersTotal] = useState(0);
    const [tagUsersPage, setTagUsersPage] = useState(1);

    // 原始数据用于搜索
    const [allTags, setAllTags] = useState<UserTag[]>([]);

    /**
     * 加载标签数据
     */
    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get('/admin/user-tags');
            if (response.data.success) {
                const data = response.data.data || [];
                setTags(data);
                setAllTags(data);
                setStats({
                    totalTags: data.length,
                    totalTaggedUsers: data.reduce((sum: number, t: UserTag) => sum + (t.userCount || 0), 0),
                });
            }
        } catch (error) {
            logger.error('Load tags error:', error);
            message.error('加载标签列表失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleSearch = (values: Record<string, unknown>) => {
        // 前端过滤
        if (values.keyword) {
            const keyword = String(values.keyword).toLowerCase();
            setTags(allTags.filter(t => t.name.toLowerCase().includes(keyword) || t.description?.toLowerCase().includes(keyword)));
        } else {
            setTags(allTags);
        }
    };

    const handleCreate = () => {
        setCurrentTag(null);
        form.resetFields();
        form.setFieldsValue({ color: '#1890ff' });
        setEditVisible(true);
    };

    const handleEdit = (record: UserTag) => {
        setCurrentTag(record);
        form.setFieldsValue({
            name: record.name,
            color: record.color,
            description: record.description,
        });
        setEditVisible(true);
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);

            const color = typeof values.color === 'string' ? values.color : (values.color as Color).toHexString();
            const data = {
                name: values.name,
                color: color,
                description: values.description || '',
            };

            if (currentTag) {
                await apiClient.put(`/admin/user-tags/${currentTag.id}`, data);
                message.success('更新成功');
            } else {
                await apiClient.post('/admin/user-tags', data);
                message.success('创建成功');
            }
            setEditVisible(false);
            loadData();
        } catch (error) {
            logger.error('Save tag error:', error);
            message.error('保存失败');
        } finally {
            setSubmitting(false);
        }
    };

    const handleDelete = async (record: UserTag) => {
        try {
            await apiClient.delete(`/admin/user-tags/${record.id}`);
            message.success('删除成功');
            loadData();
        } catch (error) {
            logger.error('Delete tag error:', error);
            message.error('删除失败');
        }
    };

    const handleViewUsers = async (record: UserTag) => {
        setCurrentTag(record);
        setTagUsersPage(1);
        setUsersVisible(true);
        loadTagUsers(record.id, 1);
    };

    const loadTagUsers = async (tagId: number, page: number) => {
        setTagUsersLoading(true);
        try {
            const response = await apiClient.get(`/admin/user-tags/${tagId}/users`, {
                params: { page, page_size: 10 },
            });
            if (response.data.success) {
                setTagUsers(response.data.data || []);
                setTagUsersTotal(response.data.pagination?.total || 0);
            }
        } catch {
            message.error('加载用户列表失败');
        } finally {
            setTagUsersLoading(false);
        }
    };

    const handleExport = () => {
        try {
            message.loading({ content: '正在导出...', key: 'export' });
            exportToCSV(tags as unknown as Record<string, unknown>[], tagExportColumns, 'user_tags');
            message.success({ content: '导出成功', key: 'export' });
        } catch {
            message.error({ content: '导出失败', key: 'export' });
        }
    };

    const searchFields: SearchField[] = [
        { name: 'keyword', label: '关键词', type: 'input', placeholder: '标签名称/描述' },
    ];

    const columns: ColumnsType<UserTag> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '标签',
            key: 'tag',
            width: 200,
            render: (_, record) => (
                <Tag color={record.color} style={{ fontSize: 14, padding: '4px 12px' }}>
                    {record.name}
                </Tag>
            ),
        },
        {
            title: '颜色',
            dataIndex: 'color',
            key: 'color',
            width: 100,
            render: (color: string) => (
                <Space>
                    <div style={{ width: 20, height: 20, backgroundColor: color, borderRadius: 4, border: '1px solid #d9d9d9' }} />
                    <span>{color}</span>
                </Space>
            ),
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            width: 200,
            ellipsis: true,
            render: (desc: string) => desc || '-',
        },
        {
            title: '用户数',
            dataIndex: 'userCount',
            key: 'userCount',
            width: 100,
            render: (count: number, record) => (
                <Button type="link" size="small" onClick={() => handleViewUsers(record)}>
                    {count || 0} 人
                </Button>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: '操作',
            key: 'action',
            width: 150,
            fixed: 'right',
            render: (_, record) => (
                <Space size="small">
                    <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                        编辑
                    </Button>
                    <Popconfirm title="确定要删除该标签吗？删除后所有用户的该标签将被移除。" onConfirm={() => handleDelete(record)}>
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    const toolbarButtons: ToolbarButton[] = [
        {
            text: '导出数据',
            icon: <DownloadOutlined />,
            needSelection: false,
            onClick: () => handleExport(),
        },
    ];

    const userColumns: ColumnsType<TagUser> = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
        {
            title: '用户',
            key: 'user',
            render: (_, record) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} src={record.avatarUrl} />
                    <span>{record.name}</span>
                </Space>
            ),
        },
        { title: '邮箱', dataIndex: 'email', key: 'email' },
    ];

    return (
        <PageContainer title="用户标签管理" subTitle="管理用户分群标签">
            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={12}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic title="标签总数" value={stats.totalTags} prefix={<TagsOutlined />} />
                    </Card>
                </Col>
                <Col span={12}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic title="已标记用户" value={stats.totalTaggedUsers} prefix={<UserOutlined />} />
                    </Card>
                </Col>
            </Row>

            <SearchTable
                columns={columns}
                dataSource={tags}
                rowKey="id"
                searchFields={searchFields}
                onSearch={handleSearch}
                onRefresh={() => loadData()}
                loading={loading}
                showCreate={true}
                createText="新增标签"
                onCreate={handleCreate}
                toolbarButtons={toolbarButtons}
                pagination={false}
                scroll={{ x: 1000 }}
            />

            {/* 编辑弹窗 */}
            <Modal
                title={currentTag ? '编辑标签' : '新增标签'}
                open={editVisible}
                onOk={handleSave}
                onCancel={() => setEditVisible(false)}
                confirmLoading={submitting}
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="name" label="标签名称" rules={[{ required: true, message: '请输入标签名称' }]}>
                        <Input placeholder="请输入标签名称" maxLength={20} />
                    </Form.Item>
                    <Form.Item name="color" label="标签颜色" rules={[{ required: true, message: '请选择颜色' }]}>
                        <ColorPicker showText />
                    </Form.Item>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={3} placeholder="请输入标签描述" maxLength={200} />
                    </Form.Item>
                    <Form.Item label="预览">
                        <Form.Item noStyle shouldUpdate>
                            {() => {
                                const name = form.getFieldValue('name') || '标签预览';
                                const colorValue = form.getFieldValue('color');
                                const color = typeof colorValue === 'string' ? colorValue : colorValue?.toHexString?.() || '#1890ff';
                                return <Tag color={color} style={{ fontSize: 14, padding: '4px 12px' }}>{name}</Tag>;
                            }}
                        </Form.Item>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 标签用户列表弹窗 */}
            <Modal
                title={
                    <Space>
                        <span>标签用户列表</span>
                        {currentTag && <Tag color={currentTag.color}>{currentTag.name}</Tag>}
                    </Space>
                }
                open={usersVisible}
                onCancel={() => setUsersVisible(false)}
                footer={null}
                width={600}
            >
                <Table
                    columns={userColumns}
                    dataSource={tagUsers}
                    rowKey="id"
                    loading={tagUsersLoading}
                    pagination={{
                        current: tagUsersPage,
                        pageSize: 10,
                        total: tagUsersTotal,
                        showTotal: t => `共 ${t} 人`,
                        onChange: (page) => {
                            setTagUsersPage(page);
                            if (currentTag) loadTagUsers(currentTag.id, page);
                        },
                    }}
                    size="small"
                />
            </Modal>
        </PageContainer>
    );
};

export default UserTagPage;
