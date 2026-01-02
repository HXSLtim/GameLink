/**
 * 动态审核页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Button, Space, Tag, Input, Select, DatePicker,
  Modal, Image, Typography, Tooltip, Popconfirm, Form, App,
} from 'antd';
import {
  SearchOutlined, CheckOutlined, CloseOutlined, DeleteOutlined,
  EyeOutlined, ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { feedApi } from '@/api/content';
import type { Feed, FeedModerationStatus } from '@/types/content';
import {
  FEED_MODERATION_STATUS_COLOR,
  FEED_MODERATION_STATUS_TEXT,
} from '@/types/content';

const { RangePicker } = DatePicker;
const { TextArea } = Input;
const { Paragraph } = Typography;

const FeedsPage: React.FC = () => {
  const { message, modal } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [detailVisible, setDetailVisible] = useState(false);
  const [rejectVisible, setRejectVisible] = useState(false);
  const [currentFeed, setCurrentFeed] = useState<Feed | null>(null);
  const [rejectForm] = Form.useForm();

  // 筛选条件
  const [keyword, setKeyword] = useState('');
  const [status, setStatus] = useState<FeedModerationStatus | ''>('');
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);

  const fetchFeeds = useCallback(async () => {
    setLoading(true);
    try {
      const params: Record<string, unknown> = { page, pageSize };
      if (keyword) params.keyword = keyword;
      if (status) params.moderationStatus = status;
      if (dateRange) {
        params.dateFrom = dateRange[0].format('YYYY-MM-DD');
        params.dateTo = dateRange[1].format('YYYY-MM-DD');
      }
      const res = await feedApi.getFeeds(params);
      if (res.data.success) {
        setFeeds(res.data.data?.items || []);
        setTotal(res.data.data?.total || 0);
      }
    } catch {
      message.error('获取动态列表失败');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword, status, dateRange, message]);

  useEffect(() => {
    fetchFeeds();
  }, [fetchFeeds]);

  const handleApprove = async (id: number) => {
    try {
      await feedApi.approveFeed(id);
      message.success('已批准');
      fetchFeeds();
    } catch {
      message.error('操作失败');
    }
  };

  const handleReject = async () => {
    if (!currentFeed) return;
    try {
      const values = await rejectForm.validateFields();
      await feedApi.rejectFeed(currentFeed.id, values.reason);
      message.success('已拒绝');
      setRejectVisible(false);
      rejectForm.resetFields();
      fetchFeeds();
    } catch {
      message.error('操作失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await feedApi.deleteFeed(id);
      message.success('已删除');
      fetchFeeds();
    } catch {
      message.error('操作失败');
    }
  };

  const handleBatchApprove = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要批准的动态');
      return;
    }
    try {
      await feedApi.batchApproveFeed(selectedRowKeys as number[]);
      message.success(`已批准 ${selectedRowKeys.length} 条动态`);
      setSelectedRowKeys([]);
      fetchFeeds();
    } catch {
      message.error('批量操作失败');
    }
  };

  const handleBatchReject = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要拒绝的动态');
      return;
    }
    modal.confirm({
      title: '批量拒绝',
      content: (
        <Form form={rejectForm}>
          <Form.Item name="reason" rules={[{ required: true, message: '请输入拒绝原因' }]}>
            <TextArea rows={3} placeholder="请输入拒绝原因" />
          </Form.Item>
        </Form>
      ),
      onOk: async () => {
        const values = await rejectForm.validateFields();
        await feedApi.batchRejectFeed(selectedRowKeys as number[], values.reason);
        message.success(`已拒绝 ${selectedRowKeys.length} 条动态`);
        setSelectedRowKeys([]);
        rejectForm.resetFields();
        fetchFeeds();
      },
    });
  };

  const columns: ColumnsType<Feed> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '作者',
      dataIndex: 'authorName',
      width: 120,
      render: (name, record) => (
        <Space>
          {record.authorAvatar && (
            <Image src={record.authorAvatar} width={24} height={24} style={{ borderRadius: '50%' }} />
          )}
          <span>{name || `用户${record.authorId}`}</span>
        </Space>
      ),
    },
    {
      title: '内容',
      dataIndex: 'content',
      ellipsis: true,
      render: (content) => (
        <Tooltip title={content}>
          <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 0 }}>
            {content}
          </Paragraph>
        </Tooltip>
      ),
    },
    {
      title: '图片',
      dataIndex: 'images',
      width: 100,
      render: (images: string[]) =>
        images?.length > 0 ? (
          <Image.PreviewGroup>
            <Space>
              {images.slice(0, 3).map((img, idx) => (
                <Image key={idx} src={img} width={40} height={40} style={{ objectFit: 'cover' }} />
              ))}
              {images.length > 3 && <span>+{images.length - 3}</span>}
            </Space>
          </Image.PreviewGroup>
        ) : '-',
    },
    {
      title: '分类',
      dataIndex: 'categoryName',
      width: 100,
      render: (name) => name || '-',
    },
    {
      title: '状态',
      dataIndex: 'moderationStatus',
      width: 100,
      render: (status: FeedModerationStatus) => (
        <Tag color={FEED_MODERATION_STATUS_COLOR[status]}>
          {FEED_MODERATION_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '发布时间',
      dataIndex: 'createdAt',
      width: 160,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => { setCurrentFeed(record); setDetailVisible(true); }}
            />
          </Tooltip>
          {record.moderationStatus === 'pending' && (
            <>
              <Tooltip title="批准">
                <Button
                  type="link"
                  size="small"
                  icon={<CheckOutlined />}
                  onClick={() => handleApprove(record.id)}
                />
              </Tooltip>
              <Tooltip title="拒绝">
                <Button
                  type="link"
                  size="small"
                  danger
                  icon={<CloseOutlined />}
                  onClick={() => { setCurrentFeed(record); setRejectVisible(true); }}
                />
              </Tooltip>
            </>
          )}
          <Popconfirm title="确定删除此动态？" onConfirm={() => handleDelete(record.id)}>
            <Tooltip title="删除">
              <Button type="link" size="small" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="动态审核"
      extra={
        <Space>
          <Button icon={<ReloadOutlined />} onClick={fetchFeeds}>刷新</Button>
          <Button type="primary" onClick={handleBatchApprove} disabled={selectedRowKeys.length === 0}>
            批量批准
          </Button>
          <Button danger onClick={handleBatchReject} disabled={selectedRowKeys.length === 0}>
            批量拒绝
          </Button>
        </Space>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Input
          placeholder="搜索内容"
          prefix={<SearchOutlined />}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onPressEnter={() => { setPage(1); fetchFeeds(); }}
          style={{ width: 200 }}
          allowClear
        />
        <Select
          placeholder="审核状态"
          value={status}
          onChange={(v) => { setStatus(v); setPage(1); }}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="pending">待审核</Select.Option>
          <Select.Option value="approved">已通过</Select.Option>
          <Select.Option value="rejected">已拒绝</Select.Option>
        </Select>
        <RangePicker
          value={dateRange}
          onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)}
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={() => { setPage(1); fetchFeeds(); }}>
          搜索
        </Button>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={feeds}
        loading={loading}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
          getCheckboxProps: (record) => ({
            disabled: record.moderationStatus !== 'pending',
          }),
        }}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
        scroll={{ x: 1200 }}
      />

      {/* 详情弹窗 - 带快捷操作 */}
      <Modal
        title="动态详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={currentFeed?.moderationStatus === 'pending' ? (
          <Space>
            <Button onClick={() => setDetailVisible(false)}>关闭</Button>
            <Button
              type="primary"
              icon={<CheckOutlined />}
              onClick={() => {
                handleApprove(currentFeed.id);
                setDetailVisible(false);
              }}
            >
              批准
            </Button>
            <Button
              danger
              icon={<CloseOutlined />}
              onClick={() => {
                setDetailVisible(false);
                setRejectVisible(true);
              }}
            >
              拒绝
            </Button>
          </Space>
        ) : null}
        width={600}
      >
        {currentFeed && (
          <div>
            <p><strong>作者：</strong>{currentFeed.authorName || `用户${currentFeed.authorId}`}</p>
            <p><strong>内容：</strong></p>
            <Paragraph style={{ whiteSpace: 'pre-wrap' }}>{currentFeed.content}</Paragraph>
            {currentFeed.images && currentFeed.images.length > 0 && (
              <>
                <p><strong>图片：</strong></p>
                <Image.PreviewGroup>
                  <Space wrap>
                    {currentFeed.images.map((img, idx) => (
                      <Image key={idx} src={img} width={100} height={100} style={{ objectFit: 'cover', borderRadius: 4 }} />
                    ))}
                  </Space>
                </Image.PreviewGroup>
              </>
            )}
            <p><strong>分类：</strong>{currentFeed.categoryName || '-'}</p>
            <p><strong>状态：</strong>
              <Tag color={FEED_MODERATION_STATUS_COLOR[currentFeed.moderationStatus]}>
                {FEED_MODERATION_STATUS_TEXT[currentFeed.moderationStatus]}
              </Tag>
            </p>
            {currentFeed.rejectionReason && (
              <p><strong>拒绝原因：</strong>{currentFeed.rejectionReason}</p>
            )}
            <p><strong>发布时间：</strong>{dayjs(currentFeed.createdAt).format('YYYY-MM-DD HH:mm:ss')}</p>
          </div>
        )}
      </Modal>

      {/* 拒绝弹窗 - 显示动态内容预览 */}
      <Modal
        title="拒绝动态"
        open={rejectVisible}
        onOk={handleReject}
        onCancel={() => { setRejectVisible(false); rejectForm.resetFields(); }}
      >
        {currentFeed && (
          <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 8, marginBottom: 16 }}>
            <div style={{ fontSize: 12, color: '#999', marginBottom: 8 }}>待拒绝的动态内容：</div>
            <Paragraph ellipsis={{ rows: 3 }} style={{ marginBottom: 8 }}>
              {currentFeed.content}
            </Paragraph>
            {currentFeed.images && currentFeed.images.length > 0 && (
              <Space>
                {currentFeed.images.slice(0, 3).map((img, idx) => (
                  <Image key={idx} src={img} width={60} height={60} style={{ objectFit: 'cover', borderRadius: 4 }} />
                ))}
                {currentFeed.images.length > 3 && <span style={{ color: '#999' }}>+{currentFeed.images.length - 3}</span>}
              </Space>
            )}
          </div>
        )}
        <Form form={rejectForm}>
          <Form.Item name="reason" label="拒绝原因" rules={[{ required: true, message: '请输入拒绝原因' }]}>
            <TextArea rows={3} placeholder="请输入拒绝原因" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default FeedsPage;
