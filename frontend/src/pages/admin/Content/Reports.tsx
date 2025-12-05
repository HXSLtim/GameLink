/**
 * 动态举报管理页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Button, Space, Tag, Select, DatePicker,
  Modal, message, Typography, Tooltip, Form, Input, Radio,
} from 'antd';
import {
  SearchOutlined, EyeOutlined, ReloadOutlined,
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
const { Paragraph } = Typography;

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
        setReports(res.data.data.items || []);
        setTotal(res.data.data.total || 0);
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

  const columns: ColumnsType<FeedReport> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '被举报动态',
      dataIndex: 'feedId',
      width: 100,
      render: (feedId) => `动态#${feedId}`,
    },
    {
      title: '举报人',
      dataIndex: 'reporterName',
      width: 120,
      render: (name, record) => name || `用户${record.reporterId}`,
    },
    {
      title: '举报原因',
      dataIndex: 'reason',
      ellipsis: true,
      render: (reason) => (
        <Tooltip title={reason}>
          <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 0 }}>
            {reason}
          </Paragraph>
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: FeedReportStatus) => (
        <Tag color={FEED_REPORT_STATUS_COLOR[status]}>
          {FEED_REPORT_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '处理人',
      dataIndex: 'handlerName',
      width: 120,
      render: (name) => name || '-',
    },
    {
      title: '举报时间',
      dataIndex: 'createdAt',
      width: 160,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => { setCurrentReport(record); setDetailVisible(true); }}
            />
          </Tooltip>
          {record.status === 'pending' && (
            <Button
              type="primary"
              size="small"
              onClick={() => { setCurrentReport(record); setProcessVisible(true); }}
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
      extra={<Button icon={<ReloadOutlined />} onClick={fetchReports}>刷新</Button>}
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          placeholder="状态"
          value={status}
          onChange={(v) => { setStatus(v); setPage(1); }}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="pending">待处理</Select.Option>
          <Select.Option value="approved">已通过</Select.Option>
          <Select.Option value="rejected">已驳回</Select.Option>
        </Select>
        <RangePicker
          value={dateRange}
          onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)}
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={() => { setPage(1); fetchReports(); }}>
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
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
        scroll={{ x: 1000 }}
      />

      {/* 详情弹窗 */}
      <Modal
        title="举报详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={600}
      >
        {currentReport && (
          <div>
            <p><strong>举报ID：</strong>{currentReport.id}</p>
            <p><strong>被举报动态：</strong>动态#{currentReport.feedId}</p>
            <p><strong>举报人：</strong>{currentReport.reporterName || `用户${currentReport.reporterId}`}</p>
            <p><strong>举报原因：</strong></p>
            <Paragraph>{currentReport.reason}</Paragraph>
            <p><strong>状态：</strong>
              <Tag color={FEED_REPORT_STATUS_COLOR[currentReport.status]}>
                {FEED_REPORT_STATUS_TEXT[currentReport.status]}
              </Tag>
            </p>
            {currentReport.handlerName && (
              <p><strong>处理人：</strong>{currentReport.handlerName}</p>
            )}
            {currentReport.handledAt && (
              <p><strong>处理时间：</strong>{dayjs(currentReport.handledAt).format('YYYY-MM-DD HH:mm:ss')}</p>
            )}
            {currentReport.handlingNote && (
              <p><strong>处理备注：</strong>{currentReport.handlingNote}</p>
            )}
            <p><strong>举报时间：</strong>{dayjs(currentReport.createdAt).format('YYYY-MM-DD HH:mm:ss')}</p>
          </div>
        )}
      </Modal>

      {/* 处理弹窗 */}
      <Modal
        title="处理举报"
        open={processVisible}
        onOk={handleProcess}
        onCancel={() => { setProcessVisible(false); processForm.resetFields(); }}
      >
        <Form form={processForm} layout="vertical">
          <Form.Item
            name="action"
            label="处理方式"
            rules={[{ required: true, message: '请选择处理方式' }]}
          >
            <Radio.Group>
              {(Object.keys(FEED_REPORT_ACTION_TEXT) as FeedReportAction[]).map((action) => (
                <Radio key={action} value={action}>
                  {FEED_REPORT_ACTION_TEXT[action]}
                </Radio>
              ))}
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
