/**
 * 评价管理 - 举报管理页面
 * 需求: 3.1, 3.2, 3.3, 3.4, 3.5
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Select,
  Button,
  Space,
  Tag,
  message,
  Modal,
  Input,
  Typography,
  Descriptions,
  Rate,
  Radio,
} from 'antd';
import {
  ReloadOutlined,
  EyeOutlined,
  DeleteOutlined,
  WarningOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { reviewReportApi } from '@/api/review';
import type {
  ReviewReport,
  ReviewReportStatus,
  ReviewReportQueryParams,
  ReportHandleAction,
} from '@/types/review';
import {
  REPORT_STATUS_TEXT,
  REPORT_STATUS_COLOR,
  REPORT_ACTION_TEXT,
} from '@/types/review';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

const ReviewReports: React.FC = () => {
  // 状态
  const [loading, setLoading] = useState(false);
  const [reports, setReports] = useState<ReviewReport[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });

  // 筛选条件
  const [statusFilter, setStatusFilter] = useState<ReviewReportStatus | undefined>();

  // 弹窗状态
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [handleModalVisible, setHandleModalVisible] = useState(false);
  const [currentReport, setCurrentReport] = useState<ReviewReport | null>(null);
  const [handleAction, setHandleAction] = useState<ReportHandleAction>('reject');
  const [handleNote, setHandleNote] = useState('');

  // 加载举报列表
  const fetchReports = useCallback(async () => {
    setLoading(true);
    try {
      const params: ReviewReportQueryParams = {
        page: pagination.current,
        pageSize: pagination.pageSize,
        status: statusFilter,
      };
      const response = await reviewReportApi.getReports(params);
      if (response.data.success) {
        setReports(response.data.data || []);
        if (response.data.pagination) {
          setPagination(prev => ({
            ...prev,
            total: response.data.pagination!.total,
          }));
        }
      }
    } catch {
      message.error('获取举报列表失败');
    } finally {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.current, pagination.pageSize, statusFilter]);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  // 表格变化
  const handleTableChange = (paginationConfig: TablePaginationConfig) => {
    setPagination(prev => ({
      ...prev,
      current: paginationConfig.current || 1,
      pageSize: paginationConfig.pageSize || 20,
    }));
  };

  // 查看详情
  const openDetailModal = (report: ReviewReport) => {
    setCurrentReport(report);
    setDetailModalVisible(true);
  };

  // 打开处理弹窗
  const openHandleModal = (report: ReviewReport) => {
    setCurrentReport(report);
    setHandleAction('reject');
    setHandleNote('');
    setHandleModalVisible(true);
  };

  // 处理举报
  const handleReport = async () => {
    if (!currentReport) return;
    try {
      const response = await reviewReportApi.handleReport(currentReport.id, {
        action: handleAction,
        note: handleNote || undefined,
      });
      if (response.data.success) {
        const actionText = REPORT_ACTION_TEXT[handleAction];
        message.success(`举报已处理: ${actionText}`);
        setHandleModalVisible(false);
        setCurrentReport(null);
        fetchReports();
      }
    } catch {
      message.error('处理失败');
    }
  };

  // 表格列定义
  const columns: ColumnsType<ReviewReport> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '评价ID',
      dataIndex: 'reviewId',
      key: 'reviewId',
      width: 100,
    },
    {
      title: '举报人',
      dataIndex: 'reporterName',
      key: 'reporterName',
      width: 120,
      render: (name: string, record: ReviewReport) => (
        <span>{name || `用户${record.reporterId}`}</span>
      ),
    },
    {
      title: '举报原因',
      dataIndex: 'reason',
      key: 'reason',
      width: 200,
      ellipsis: true,
      render: (reason: string) => (
        <Text ellipsis={{ tooltip: reason }}>{reason}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: ReviewReportStatus) => (
        <Tag color={REPORT_STATUS_COLOR[status]}>
          {REPORT_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '处理人',
      dataIndex: 'handlerName',
      key: 'handlerName',
      width: 120,
      render: (name: string, record: ReviewReport) => (
        record.handledBy ? (name || `用户${record.handledBy}`) : '-'
      ),
    },
    {
      title: '举报时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '处理时间',
      dataIndex: 'handledAt',
      key: 'handledAt',
      width: 180,
      render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      fixed: 'right',
      render: (_: unknown, record: ReviewReport) => (
        <Space size={4}>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => openDetailModal(record)}
          >
            详情
          </Button>
          {record.status === 'pending' && (
            <Button
              type="primary"
              size="small"
              onClick={() => openHandleModal(record)}
            >
              处理
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Card title="举报管理">
      {/* 筛选区域 */}
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="举报状态"
          value={statusFilter}
          onChange={setStatusFilter}
          style={{ width: 120 }}
          allowClear
          options={[
            { value: 'pending', label: '待处理' },
            { value: 'approved', label: '已通过' },
            { value: 'rejected', label: '已驳回' },
          ]}
        />
        <Button icon={<ReloadOutlined />} onClick={fetchReports}>
          刷新
        </Button>
      </Space>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={reports}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1200 }}
      />

      {/* 详情弹窗 */}
      <Modal
        title="举报详情"
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false);
          setCurrentReport(null);
        }}
        footer={
          currentReport?.status === 'pending' ? (
            <Button type="primary" onClick={() => {
              setDetailModalVisible(false);
              openHandleModal(currentReport);
            }}>
              处理举报
            </Button>
          ) : null
        }
        width={700}
      >
        {currentReport && (
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="举报ID">{currentReport.id}</Descriptions.Item>
            <Descriptions.Item label="评价ID">{currentReport.reviewId}</Descriptions.Item>
            <Descriptions.Item label="举报人">
              {currentReport.reporterName || `用户${currentReport.reporterId}`}
            </Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag color={REPORT_STATUS_COLOR[currentReport.status]}>
                {REPORT_STATUS_TEXT[currentReport.status]}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="举报原因" span={2}>
              <Paragraph>{currentReport.reason}</Paragraph>
            </Descriptions.Item>
            {currentReport.evidence && (
              <Descriptions.Item label="举报证据" span={2}>
                <a href={currentReport.evidence} target="_blank" rel="noopener noreferrer">
                  查看证据
                </a>
              </Descriptions.Item>
            )}
            <Descriptions.Item label="举报时间">
              {dayjs(currentReport.createdAt).format('YYYY-MM-DD HH:mm:ss')}
            </Descriptions.Item>
            {currentReport.handledAt && (
              <Descriptions.Item label="处理时间">
                {dayjs(currentReport.handledAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
            )}
            {currentReport.handlingNote && (
              <Descriptions.Item label="处理备注" span={2}>
                <Text>{currentReport.handlingNote}</Text>
              </Descriptions.Item>
            )}
            {currentReport.review && (
              <>
                <Descriptions.Item label="被举报评价" span={2}>
                  <div style={{ padding: '8px 0' }}>
                    <div>
                      <Text strong>评分: </Text>
                      <Rate disabled value={currentReport.review.rating} style={{ fontSize: 14 }} />
                    </div>
                    <div style={{ marginTop: 8 }}>
                      <Text strong>内容: </Text>
                      <Paragraph style={{ marginBottom: 0 }}>
                        {currentReport.review.comment || '-'}
                      </Paragraph>
                    </div>
                  </div>
                </Descriptions.Item>
              </>
            )}
          </Descriptions>
        )}
      </Modal>

      {/* 处理弹窗 */}
      <Modal
        title="处理举报"
        open={handleModalVisible}
        onOk={handleReport}
        onCancel={() => {
          setHandleModalVisible(false);
          setCurrentReport(null);
          setHandleNote('');
        }}
        okText="确定"
        cancelText="取消"
      >
        <div style={{ marginBottom: 16 }}>
          <Text strong>选择处理方式:</Text>
          <Radio.Group
            value={handleAction}
            onChange={e => setHandleAction(e.target.value)}
            style={{ display: 'block', marginTop: 8 }}
          >
            <Space orientation="vertical">
              <Radio value="delete">
                <Space>
                  <DeleteOutlined style={{ color: '#ff4d4f' }} />
                  删除评价 - 删除被举报的评价
                </Space>
              </Radio>
              <Radio value="warn">
                <Space>
                  <WarningOutlined style={{ color: '#faad14' }} />
                  警告评价者 - 保留评价但警告用户
                </Space>
              </Radio>
              <Radio value="reject">
                <Space>
                  <CloseCircleOutlined style={{ color: 'var(--ant-color-primary)' }} />
                  驳回举报 - 举报不成立
                </Space>
              </Radio>
            </Space>
          </Radio.Group>
        </div>
        <div>
          <Text strong>处理备注 (可选):</Text>
          <TextArea
            placeholder="请输入处理备注"
            value={handleNote}
            onChange={e => setHandleNote(e.target.value)}
            rows={3}
            style={{ marginTop: 8 }}
          />
        </div>
      </Modal>
    </Card>
  );
};

export default ReviewReports;
