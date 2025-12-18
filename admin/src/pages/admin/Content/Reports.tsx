/**
 * 动态举报管理页面
 * 优化：显示被举报动态内容，方便管理员审核
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Button, Space, Tag, Select, DatePicker,
  Modal, message, Typography, Tooltip, Form, Input, Radio,
  Image, Descriptions, Divider, Avatar, Empty,
} from 'antd';
import {
  SearchOutlined, EyeOutlined, ReloadOutlined, UserOutlined,
  FileImageOutlined, ExclamationCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { feedReportApi } from '@/api/content';
import type { FeedReport, FeedReportStatus, FeedReportAction } from '@/types/content';
import {
  FEED_REPORT_STATUS_TEXT,
  FEED_REPORT_STATUS_COLOR,
  FEED_REPORT_ACTION_TEXT,
} from '@/types/content';

const { RangePicker } = DatePicker;
const { TextArea } = Input;
const { Paragraph, Text } = Typography;

const ReportsPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [reports, setReports] = useState<FeedReport[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailVisible, setDetailVisible] = useState(false);
  const [processVisible, setProcessVisible] = useState(false);
  const [currentReport, setCurrentReport] = useState<FeedReport | null>(null);
  const [processForm] = Form.useForm();

  // 筛选条件
  const [status, setStatus] = useState<FeedReportStatus | ''>('');
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);

  const fetchReports = useCallback(async () => {
    setLoading(true);
    try {
      const params: Record<string, unknown> = { page, pageSize };
      if (status) params.status = status;
      if (dateRange) {
        params.dateFrom = dateRange[0].format('YYYY-MM-DD');
        params.dateTo = dateRange[1].format('YYYY-MM-DD');
      }
      const res = await feedReportApi.getReports(params);
      if (res.data.success) {
        setReports(res.data.data?.items || []);
        setTotal(res.data.data?.total || 0);
      }
    } catch {
      message.error('获取举报列表失败');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, status, dateRange]);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  const handleProcess = async () => {
    if (!currentReport) return;
    try {
      const values = await processForm.validateFields();
      await feedReportApi.processReport(currentReport.id, {
        action: values.action,
        note: values.note,
      });
      message.success('处理成功');
      setProcessVisible(false);
      processForm.resetFields();
      fetchReports();
    } catch {
      message.error('操作失败');
    }
  };

  // 渲染动态内容预览
  const renderFeedPreview = (report: FeedReport) => {
    const feed = report.feed;
    if (!feed) {
      return <Text type="secondary">动态已删除或不存在</Text>;
    }
    return (
      <div style={{ maxWidth: 300 }}>
        <Paragraph
          ellipsis={{ rows: 2, expandable: false }}
          style={{ marginBottom: 4 }}
        >
          {feed.content || '无文字内容'}
        </Paragraph>
        {feed.images && feed.images.length > 0 && (
          <Space size={4}>
            <FileImageOutlined />
            <Text type="secondary">{feed.images.length}张图片</Text>
          </Space>
        )}
      </div>
    );
  };

  const columns: ColumnsType<FeedReport> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 70,
    },
    {
      title: '被举报动态',
      key: 'feedContent',
      width: 320,
      render: (_, record) => renderFeedPreview(record),
    },
    {
      title: '举报原因',
      dataIndex: 'reason',
      width: 200,
      ellipsis: true,
      render: (reason) => (
        <Tooltip title={reason}>
          <Space>
            <ExclamationCircleOutlined style={{ color: '#faad14' }} />
            <Text>{reason}</Text>
          </Space>
        </Tooltip>
      ),
    },
    {
      title: '举报人',
      dataIndex: 'reporterName',
      width: 100,
      render: (name, record) => name || `用户${record.reporterId}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (status: FeedReportStatus) => (
        <Tag color={FEED_REPORT_STATUS_COLOR[status]}>
          {FEED_REPORT_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '举报时间',
      dataIndex: 'createdAt',
      width: 150,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => {
              setCurrentReport(record);
              setDetailVisible(true);
            }}
          >
            详情
          </Button>
          {record.status === 'pending' && (
            <Button
              type="primary"
              size="small"
              onClick={() => {
                setCurrentReport(record);
                setProcessVisible(true);
              }}
            >
              处理
            </Button>
          )}
        </Space>
      ),
    },
  ];


  return (
    <Card
      title="动态举报管理"
      extra={
        <Button icon={<ReloadOutlined />} onClick={fetchReports}>
          刷新
        </Button>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          placeholder="状态"
          value={status}
          onChange={(v) => {
            setStatus(v);
            setPage(1);
          }}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="pending">待处理</Select.Option>
          <Select.Option value="approved">已通过</Select.Option>
          <Select.Option value="rejected">已驳回</Select.Option>
        </Select>
        <RangePicker
          value={dateRange}
          onChange={(dates) =>
            setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)
          }
        />
        <Button
          type="primary"
          icon={<SearchOutlined />}
          onClick={() => {
            setPage(1);
            fetchReports();
          }}
        >
          搜索
        </Button>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={reports}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
        }}
        scroll={{ x: 1100 }}
      />

      {/* 详情弹窗 - 包含完整动态内容 */}
      <Modal
        title="举报详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={
          currentReport?.status === 'pending' ? (
            <Space>
              <Button onClick={() => setDetailVisible(false)}>关闭</Button>
              <Button
                type="primary"
                onClick={() => {
                  setDetailVisible(false);
                  setProcessVisible(true);
                }}
              >
                立即处理
              </Button>
            </Space>
          ) : null
        }
        width={700}
      >
        {currentReport && (
          <div>
            {/* 被举报的动态内容 */}
            <Descriptions title="被举报动态" column={1} bordered size="small">
              {currentReport.feed ? (
                <>
                  <Descriptions.Item label="发布者">
                    <Space>
                      <Avatar
                        src={currentReport.feed.authorAvatar}
                        icon={<UserOutlined />}
                        size="small"
                      />
                      <span>
                        {currentReport.feed.authorName ||
                          `用户${currentReport.feed.authorId}`}
                      </span>
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label="动态内容">
                    <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>
                      {currentReport.feed.content || '无文字内容'}
                    </Paragraph>
                  </Descriptions.Item>
                  {currentReport.feed.images &&
                    currentReport.feed.images.length > 0 && (
                      <Descriptions.Item label="图片">
                        <Image.PreviewGroup>
                          <Space wrap>
                            {currentReport.feed.images.map((img, idx) => (
                              <Image
                                key={idx}
                                src={img}
                                width={80}
                                height={80}
                                style={{ objectFit: 'cover', borderRadius: 4 }}
                              />
                            ))}
                          </Space>
                        </Image.PreviewGroup>
                      </Descriptions.Item>
                    )}
                  <Descriptions.Item label="动态状态">
                    <Tag
                      color={
                        currentReport.feed.moderationStatus === 'approved'
                          ? 'green'
                          : currentReport.feed.moderationStatus === 'rejected'
                            ? 'red'
                            : currentReport.feed.moderationStatus === 'deleted'
                              ? 'default'
                              : 'orange'
                      }
                    >
                      {currentReport.feed.moderationStatus}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="发布时间">
                    {currentReport.feed.createdAt}
                  </Descriptions.Item>
                </>
              ) : (
                <Descriptions.Item label="状态">
                  <Empty
                    description="动态已删除或不存在"
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                  />
                </Descriptions.Item>
              )}
            </Descriptions>

            <Divider />

            {/* 举报信息 */}
            <Descriptions title="举报信息" column={2} bordered size="small">
              <Descriptions.Item label="举报ID">{currentReport.id}</Descriptions.Item>
              <Descriptions.Item label="举报人">
                {currentReport.reporterName || `用户${currentReport.reporterId}`}
              </Descriptions.Item>
              <Descriptions.Item label="举报原因" span={2}>
                <Text type="warning">{currentReport.reason}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="举报时间">
                {dayjs(currentReport.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={FEED_REPORT_STATUS_COLOR[currentReport.status]}>
                  {FEED_REPORT_STATUS_TEXT[currentReport.status]}
                </Tag>
              </Descriptions.Item>
              {currentReport.handlerName && (
                <Descriptions.Item label="处理人">
                  {currentReport.handlerName}
                </Descriptions.Item>
              )}
              {currentReport.handledAt && (
                <Descriptions.Item label="处理时间">
                  {dayjs(currentReport.handledAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
              )}
              {currentReport.handlingNote && (
                <Descriptions.Item label="处理备注" span={2}>
                  {currentReport.handlingNote}
                </Descriptions.Item>
              )}
            </Descriptions>
          </div>
        )}
      </Modal>

      {/* 处理弹窗 */}
      <Modal
        title="处理举报"
        open={processVisible}
        onOk={handleProcess}
        onCancel={() => {
          setProcessVisible(false);
          processForm.resetFields();
        }}
        okText="确认处理"
      >
        {currentReport?.feed && (
          <div
            style={{
              background: '#f5f5f5',
              padding: 12,
              borderRadius: 8,
              marginBottom: 16,
            }}
          >
            <Text type="secondary" style={{ fontSize: 12 }}>
              被举报内容预览：
            </Text>
            <Paragraph
              ellipsis={{ rows: 3 }}
              style={{ marginBottom: 0, marginTop: 4 }}
            >
              {currentReport.feed.content}
            </Paragraph>
          </div>
        )}
        <Form form={processForm} layout="vertical">
          <Form.Item
            name="action"
            label="处理方式"
            rules={[{ required: true, message: '请选择处理方式' }]}
          >
            <Radio.Group>
              {(Object.keys(FEED_REPORT_ACTION_TEXT) as FeedReportAction[]).map(
                (action) => (
                  <Radio key={action} value={action} style={{ display: 'block', marginBottom: 8 }}>
                    {FEED_REPORT_ACTION_TEXT[action]}
                  </Radio>
                )
              )}
            </Radio.Group>
          </Form.Item>
          <Form.Item name="note" label="处理备注">
            <TextArea rows={3} placeholder="请输入处理备注（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default ReportsPage;
