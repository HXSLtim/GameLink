/**
 * 评价管理 - 敏感词管理页面
 * 需求: 5.1, 5.2, 5.3, 5.4
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Input,
  Select,
  Button,
  Space,
  Tag,
  message,
  Modal,
  Form,
  Popconfirm,
} from 'antd';
import {
  SearchOutlined,
  PlusOutlined,
  ReloadOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { sensitiveWordApi } from '@/api/review';
import type {
  SensitiveWord,
  SensitiveWordCategory,
  SensitiveWordSeverity,
  SensitiveWordQueryParams,
  SensitiveWordFormData,
} from '@/types/review';
import {
  SENSITIVE_WORD_CATEGORY_TEXT,
  SENSITIVE_WORD_CATEGORY_COLOR,
  SENSITIVE_WORD_SEVERITY_TEXT,
  SENSITIVE_WORD_SEVERITY_COLOR,
} from '@/types/review';

const SensitiveWords: React.FC = () => {
  const [form] = Form.useForm<SensitiveWordFormData>();

  // 状态
  const [loading, setLoading] = useState(false);
  const [words, setWords] = useState<SensitiveWord[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });

  // 筛选条件
  const [keyword, setKeyword] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<SensitiveWordCategory | undefined>();

  // 弹窗状态
  const [modalVisible, setModalVisible] = useState(false);
  const [editingWord, setEditingWord] = useState<SensitiveWord | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // 加载敏感词列表
  const fetchWords = useCallback(async () => {
    setLoading(true);
    try {
      const params: SensitiveWordQueryParams = {
        page: pagination.current,
        pageSize: pagination.pageSize,
        keyword: keyword || undefined,
        category: categoryFilter,
      };
      const response = await sensitiveWordApi.getWords(params) as unknown as {
        success: boolean;
        data: { items: SensitiveWord[]; total: number };
      };
      if (response.success) {
        setWords(response.data.items || []);
        setPagination(prev => ({
          ...prev,
          total: response.data.total,
        }));
      }
    } catch {
      message.error('获取敏感词列表失败');
    } finally {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.current, pagination.pageSize, keyword, categoryFilter]);

  useEffect(() => {
    fetchWords();
  }, [fetchWords]);

  // 搜索
  const handleSearch = () => {
    setPagination(prev => ({ ...prev, current: 1 }));
  };

  // 重置
  const handleReset = () => {
    setKeyword('');
    setCategoryFilter(undefined);
    setPagination(prev => ({ ...prev, current: 1 }));
  };

  // 表格变化
  const handleTableChange = (paginationConfig: TablePaginationConfig) => {
    setPagination(prev => ({
      ...prev,
      current: paginationConfig.current || 1,
      pageSize: paginationConfig.pageSize || 20,
    }));
  };

  // 打开新增弹窗
  const openAddModal = () => {
    setEditingWord(null);
    form.resetFields();
    setModalVisible(true);
  };

  // 打开编辑弹窗
  const openEditModal = (word: SensitiveWord) => {
    setEditingWord(word);
    form.setFieldsValue({
      word: word.word,
      category: word.category,
      severity: word.severity,
    });
    setModalVisible(true);
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);

      if (editingWord) {
        // 编辑
        const response = await sensitiveWordApi.updateWord(editingWord.id, values) as unknown as { success: boolean };
        if (response.success) {
          message.success('敏感词已更新');
          setModalVisible(false);
          fetchWords();
        }
      } else {
        // 新增
        const response = await sensitiveWordApi.addWord(values) as unknown as { success: boolean };
        if (response.success) {
          message.success('敏感词已添加');
          setModalVisible(false);
          fetchWords();
        }
      }
    } catch (err) {
      if (err instanceof Error && err.message.includes('已存在')) {
        message.error('该敏感词已存在');
      } else {
        message.error('操作失败');
      }
    } finally {
      setSubmitting(false);
    }
  };

  // 删除敏感词
  const handleDelete = async (id: number) => {
    try {
      const response = await sensitiveWordApi.deleteWord(id) as unknown as { success: boolean };
      if (response.success) {
        message.success('敏感词已删除');
        fetchWords();
      }
    } catch {
      message.error('删除失败');
    }
  };

  // 表格列定义
  const columns: ColumnsType<SensitiveWord> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '敏感词',
      dataIndex: 'word',
      key: 'word',
      width: 200,
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      width: 120,
      render: (category: SensitiveWordCategory) => (
        <Tag color={SENSITIVE_WORD_CATEGORY_COLOR[category]}>
          {SENSITIVE_WORD_CATEGORY_TEXT[category]}
        </Tag>
      ),
    },
    {
      title: '严重程度',
      dataIndex: 'severity',
      key: 'severity',
      width: 120,
      render: (severity: SensitiveWordSeverity) => (
        <Tag color={SENSITIVE_WORD_SEVERITY_COLOR[severity]}>
          {SENSITIVE_WORD_SEVERITY_TEXT[severity]}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: unknown, record: SensitiveWord) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => openEditModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个敏感词吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 分类选项
  const categoryOptions = [
    { value: 'political', label: '政治' },
    { value: 'pornographic', label: '色情' },
    { value: 'violent', label: '暴力' },
    { value: 'advertising', label: '广告' },
    { value: 'other', label: '其他' },
  ];

  // 严重程度选项
  const severityOptions = [
    { value: 'low', label: '低' },
    { value: 'medium', label: '中' },
    { value: 'high', label: '高' },
  ];

  return (
    <Card title="敏感词管理">
      {/* 筛选区域 */}
      <Space wrap style={{ marginBottom: 16 }}>
        <Input
          placeholder="搜索敏感词"
          value={keyword}
          onChange={e => setKeyword(e.target.value)}
          style={{ width: 200 }}
          prefix={<SearchOutlined />}
          allowClear
          onPressEnter={handleSearch}
        />
        <Select
          placeholder="分类"
          value={categoryFilter}
          onChange={setCategoryFilter}
          style={{ width: 120 }}
          allowClear
          options={categoryOptions}
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
          搜索
        </Button>
        <Button icon={<ReloadOutlined />} onClick={handleReset}>
          重置
        </Button>
        <Button type="primary" icon={<PlusOutlined />} onClick={openAddModal}>
          添加敏感词
        </Button>
      </Space>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={words}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onChange={handleTableChange}
      />

      {/* 新增/编辑弹窗 */}
      <Modal
        title={editingWord ? '编辑敏感词' : '添加敏感词'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false);
          setEditingWord(null);
          form.resetFields();
        }}
        confirmLoading={submitting}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            category: 'other',
            severity: 'medium',
          }}
        >
          <Form.Item
            name="word"
            label="敏感词"
            rules={[
              { required: true, message: '请输入敏感词' },
              { max: 50, message: '敏感词不能超过50个字符' },
            ]}
          >
            <Input placeholder="请输入敏感词" />
          </Form.Item>
          <Form.Item
            name="category"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <Select options={categoryOptions} placeholder="请选择分类" />
          </Form.Item>
          <Form.Item
            name="severity"
            label="严重程度"
            rules={[{ required: true, message: '请选择严重程度' }]}
          >
            <Select options={severityOptions} placeholder="请选择严重程度" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default SensitiveWords;
