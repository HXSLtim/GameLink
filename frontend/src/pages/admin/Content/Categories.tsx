/**
 * 内容分类管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Button, Space, Tag, Input, Select,
  Modal, message, Form, InputNumber, Popconfirm,
} from 'antd';
import {
  SearchOutlined, PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { contentCategoryApi } from '@/api/content';
import type { ContentCategory, ContentCategoryStatus } from '@/types/content';
import { CATEGORY_STATUS_TEXT, CATEGORY_STATUS_COLOR } from '@/types/content';

const { TextArea } = Input;

const CategoriesPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<ContentCategory[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingCategory, setEditingCategory] = useState<ContentCategory | null>(null);
  const [form] = Form.useForm();

  // 筛选条件
  const [keyword, setKeyword] = useState('');
  const [status, setStatus] = useState<ContentCategoryStatus | ''>('');

  const fetchCategories = useCallback(async () => {
    setLoading(true);
    try {
      const params: Record<string, unknown> = { page, pageSize };
      if (keyword) params.keyword = keyword;
      if (status) params.status = status;
      const res = await contentCategoryApi.getCategories(params) as unknown as { success: boolean; data: { items: ContentCategory[]; total: number } };
      if (res.success) {
        setCategories(res.data?.items || []);
        setTotal(res.data?.total || 0);
      }
    } catch {
      message.error('获取分类列表失败');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword, status]);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  const handleCreate = () => {
    setEditingCategory(null);
    form.resetFields();
    form.setFieldsValue({ status: 'active', sortOrder: 0 });
    setModalVisible(true);
  };

  const handleEdit = (record: ContentCategory) => {
    setEditingCategory(record);
    form.setFieldsValue({
      name: record.name,
      description: record.description,
      sortOrder: record.sortOrder,
      status: record.status,
    });
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingCategory) {
        await contentCategoryApi.updateCategory(editingCategory.id, values);
        message.success('更新成功');
      } else {
        await contentCategoryApi.createCategory(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      form.resetFields();
      fetchCategories();
    } catch (err: unknown) {
      const error = err as { response?: { status?: number } };
      if (error.response?.status === 409) {
        message.error('分类名称已存在');
      } else {
        message.error('操作失败');
      }
    }
  };

  const handleDelete = async (id: number, feedCount?: number) => {
    if (feedCount && feedCount > 0) {
      // 如果分类下有动态，需要选择迁移目标
      Modal.confirm({
        title: '删除分类',
        content: `该分类下有 ${feedCount} 条动态，请选择迁移目标分类或直接删除（动态将变为无分类）`,
        okText: '直接删除',
        cancelText: '取消',
        onOk: async () => {
          try {
            await contentCategoryApi.deleteCategory(id);
            message.success('删除成功');
            fetchCategories();
          } catch {
            message.error('删除失败');
          }
        },
      });
    } else {
      try {
        await contentCategoryApi.deleteCategory(id);
        message.success('删除成功');
        fetchCategories();
      } catch {
        message.error('删除失败');
      }
    }
  };

  const columns: ColumnsType<ContentCategory> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '分类名称',
      dataIndex: 'name',
      width: 150,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
      render: (desc) => desc || '-',
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      width: 80,
    },
    {
      title: '动态数',
      dataIndex: 'feedCount',
      width: 100,
      render: (count) => count ?? '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: ContentCategoryStatus) => (
        <Tag color={CATEGORY_STATUS_COLOR[status]}>
          {CATEGORY_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      width: 160,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除此分类？"
            onConfirm={() => handleDelete(record.id, record.feedCount)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="内容分类管理"
      extra={
        <Space>
          <Button icon={<ReloadOutlined />} onClick={fetchCategories}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新建分类
          </Button>
        </Space>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Input
          placeholder="搜索分类名称"
          prefix={<SearchOutlined />}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onPressEnter={() => { setPage(1); fetchCategories(); }}
          style={{ width: 200 }}
          allowClear
        />
        <Select
          placeholder="状态"
          value={status}
          onChange={(v) => { setStatus(v); setPage(1); }}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="active">启用</Select.Option>
          <Select.Option value="inactive">禁用</Select.Option>
        </Select>
        <Button type="primary" icon={<SearchOutlined />} onClick={() => { setPage(1); fetchCategories(); }}>
          搜索
        </Button>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={categories}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
        scroll={{ x: 900 }}
      />

      {/* 新建/编辑弹窗 */}
      <Modal
        title={editingCategory ? '编辑分类' : '新建分类'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => { setModalVisible(false); form.resetFields(); }}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="分类名称"
            rules={[{ required: true, message: '请输入分类名称' }]}
          >
            <Input placeholder="请输入分类名称" maxLength={50} />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="请输入分类描述（可选）" maxLength={200} />
          </Form.Item>
          <Form.Item name="sortOrder" label="排序">
            <InputNumber min={0} max={9999} style={{ width: '100%' }} placeholder="数字越小越靠前" />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Select.Option value="active">启用</Select.Option>
              <Select.Option value="inactive">禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default CategoriesPage;
