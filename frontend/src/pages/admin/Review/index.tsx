/**
 * 评价管理 - 评价列表页面
 * 需求: 1.1, 1.2, 1.3, 1.5
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Input,
  Select,
  DatePicker,
  Button,
  Space,
  Tag,
  Rate,
  Badge,
  message,
  Popconfirm,
  Image,
  Typography,
  Drawer,
  Descriptions,
  Divider,
  Spin,
  Timeline,
} from 'antd';
import {
  SearchOutlined,
  ReloadOutlined,
  EyeOutlined,
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  UserOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { reviewApi } from '@/api/review';
import type {
  Review,
  ReviewStatus,
  ReviewQueryParams,
  OperationLog,
} from '@/types/review';
import {
  REVIEW_STATUS_TEXT,
  REVIEW_STATUS_COLOR,
} from '@/types/review';

const { RangePicker } = DatePicker;
const { Text } = Typography;

const ReviewList: React.FC = () => {
  // 状态
  const [loading, setLoading] = useState(false);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });

  // 筛选条件
  const [filters, setFilters] = useState<ReviewQueryParams>({});
  const [keyword, setKeyword] = useState('');
  const [statusFilter, setStatusFilter] = useState<ReviewStatus | undefined>();
  const [ratingFilter, setRatingFilter] = useState<number | undefined>();
  const [reportedFilter, setReportedFilter] = useState<boolean | undefined>();
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);

  // 详情抽屉状态
  const [detailVisible, setDetailVisible] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [currentReview, setCurrentReview] = useState<Review | null>(null);
  const [reviewLogs, setReviewLogs] = useState<OperationLog[]>([]);

  // 加载数据
  const fetchReviews = useCallback(async (params?: ReviewQueryParams) => {
    setLoading(true);
    try {
      const queryParams: ReviewQueryParams = {
        page: pagination.current,
        pageSize: pagination.pageSize,
        ...filters,
        ...params,
      };

      // API 响应已在 client.ts 拦截器中解包，直接返回 data
      const response = await reviewApi.getReviews(queryParams) as unknown as {
        success: boolean;
        data: Review[];
        pagination?: { total: number };
      };
      if (response.success) {
        setReviews(response.data || []);
        if (response.pagination) {
          setPagination(prev => ({
            ...prev,
            total: response.pagination!.total,
          }));
        }
      }
    } catch (err) {
      message.error('获取评价列表失败');
      console.error('Failed to fetch reviews:', err);
    } finally {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.current, pagination.pageSize, filters]);

  useEffect(() => {
    fetchReviews();
  }, [fetchReviews]);

  // 搜索
  const handleSearch = () => {
    const newFilters: ReviewQueryParams = {
      keyword: keyword || undefined,
      status: statusFilter,
      minRating: ratingFilter,
      maxRating: ratingFilter,
      isReported: reportedFilter,
      startTime: dateRange?.[0]?.format('YYYY-MM-DD') || undefined,
      endTime: dateRange?.[1]?.format('YYYY-MM-DD') || undefined,
    };
    setFilters(newFilters);
    setPagination(prev => ({ ...prev, current: 1 }));
  };

  // 重置筛选
  const handleReset = () => {
    setKeyword('');
    setStatusFilter(undefined);
    setRatingFilter(undefined);
    setReportedFilter(undefined);
    setDateRange(null);
    setFilters({});
    setPagination(prev => ({ ...prev, current: 1 }));
  };

  // 表格变化
  const handleTableChange = (
    paginationConfig: TablePaginationConfig,
  ) => {
    setPagination(prev => ({
      ...prev,
      current: paginationConfig.current || 1,
      pageSize: paginationConfig.pageSize || 20,
    }));
  };

  // 批准评价
  const handleApprove = async (id: number) => {
    try {
      const response = await reviewApi.approveReview(id) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已批准');
        fetchReviews();
      }
    } catch {
      message.error('操作失败');
    }
  };

  // 删除评价
  const handleDelete = async (id: number) => {
    try {
      const response = await reviewApi.deleteReview(id) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已删除');
        fetchReviews();
      }
    } catch {
      message.error('删除失败');
    }
  };

  // 查看详情
  const handleViewDetail = async (record: Review) => {
    setCurrentReview(record);
    setDetailVisible(true);
    setDetailLoading(true);
    try {
      // 获取操作日志
      const logsRes = await reviewApi.getReviewLogs(record.id) as unknown as {
        success: boolean;
        data: OperationLog[];
      };
      if (logsRes.success) {
        setReviewLogs(Array.isArray(logsRes.data) ? logsRes.data : []);
      }
    } catch {
      setReviewLogs([]);
    } finally {
      setDetailLoading(false);
    }
  };

  // 关闭详情抽屉
  const handleCloseDetail = () => {
    setDetailVisible(false);
    setCurrentReview(null);
    setReviewLogs([]);
  };

  // 表格列定义
  const columns: ColumnsType<Review> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '订单ID',
      dataIndex: 'orderId',
      key: 'orderId',
      width: 100,
    },
    {
      title: '评价者',
      dataIndex: 'reviewerName',
      key: 'reviewerName',
      width: 120,
      render: (name: string, record: Review) => (
        <span>{name || `用户${record.reviewerId}`}</span>
      ),
    },
    {
      title: '被评价者',
      dataIndex: 'playerName',
      key: 'playerName',
      width: 120,
      render: (name: string, record: Review) => (
        <span>{name || `陪玩师${record.playerId}`}</span>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      width: 150,
      render: (rating: number) => <Rate disabled defaultValue={rating} />,
    },
    {
      title: '评价内容',
      dataIndex: 'comment',
      key: 'comment',
      width: 200,
      ellipsis: true,
      render: (comment: string) => (
        <Text ellipsis={{ tooltip: comment }}>{comment || '-'}</Text>
      ),
    },
    {
      title: '图片',
      dataIndex: 'images',
      key: 'images',
      width: 100,
      render: (images: string[]) => (
        images && images.length > 0 ? (
          <Image.PreviewGroup>
            <Space>
              {images.slice(0, 2).map((url, index) => (
                <Image
                  key={index}
                  width={40}
                  height={40}
                  src={url}
                  style={{ objectFit: 'cover' }}
                />
              ))}
              {images.length > 2 && <Text type="secondary">+{images.length - 2}</Text>}
            </Space>
          </Image.PreviewGroup>
        ) : '-'
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: ReviewStatus, record: Review) => (
        <Space>
          <Tag color={REVIEW_STATUS_COLOR[status]}>
            {REVIEW_STATUS_TEXT[status]}
          </Tag>
          {record.isReported && (
            <Badge status="error" text="被举报" />
          )}
        </Space>
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
      width: 200,
      fixed: 'right',
      render: (_: unknown, record: Review) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          {record.status === 'pending' && (
            <>
              <Button
                type="link"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleApprove(record.id)}
              >
                批准
              </Button>
              <Button
                type="link"
                size="small"
                danger
                icon={<CloseOutlined />}
                onClick={() => {
                  // TODO: 打开拒绝弹窗
                  message.info('请前往审核页面进行拒绝操作');
                }}
              >
                拒绝
              </Button>
            </>
          )}
          <Popconfirm
            title="确定要删除这条评价吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="评价列表">
      {/* 筛选区域 */}
      <Space wrap style={{ marginBottom: 16 }}>
        <Input
          placeholder="搜索订单ID/评价内容"
          value={keyword}
          onChange={e => setKeyword(e.target.value)}
          style={{ width: 200 }}
          prefix={<SearchOutlined />}
          allowClear
        />
        <Select
          placeholder="评价状态"
          value={statusFilter}
          onChange={setStatusFilter}
          style={{ width: 120 }}
          allowClear
          options={[
            { value: 'pending', label: '待审核' },
            { value: 'approved', label: '已通过' },
            { value: 'rejected', label: '已拒绝' },
            { value: 'deleted', label: '已删除' },
          ]}
        />
        <Select
          placeholder="评分"
          value={ratingFilter}
          onChange={setRatingFilter}
          style={{ width: 100 }}
          allowClear
          options={[
            { value: 1, label: '1星' },
            { value: 2, label: '2星' },
            { value: 3, label: '3星' },
            { value: 4, label: '4星' },
            { value: 5, label: '5星' },
          ]}
        />
        <Select
          placeholder="举报状态"
          value={reportedFilter}
          onChange={setReportedFilter}
          style={{ width: 120 }}
          allowClear
          options={[
            { value: true, label: '被举报' },
            { value: false, label: '未举报' },
          ]}
        />
        <RangePicker
          value={dateRange}
          onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)}
          placeholder={['开始日期', '结束日期']}
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
          搜索
        </Button>
        <Button icon={<ReloadOutlined />} onClick={handleReset}>
          重置
        </Button>
      </Space>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={reviews}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1400 }}
      />

      {/* 详情抽屉 */}
      <Drawer
        title="评价详情"
        placement="right"
        size="large"
        onClose={handleCloseDetail}
        open={detailVisible}
      >
        {currentReview && (
          <Spin spinning={detailLoading}>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="评价ID">{currentReview.id}</Descriptions.Item>
              <Descriptions.Item label="订单ID">{currentReview.orderId}</Descriptions.Item>
              <Descriptions.Item label="评价者">
                <Space>
                  <UserOutlined />
                  {currentReview.reviewerName || `用户${currentReview.reviewerId}`}
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="被评价者">
                <Space>
                  <UserOutlined />
                  {currentReview.playerName || `陪玩师${currentReview.playerId}`}
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="评分" span={2}>
                <Rate disabled value={currentReview.rating} />
                <Text style={{ marginLeft: 8 }}>{currentReview.rating} 分</Text>
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={REVIEW_STATUS_COLOR[currentReview.status]}>
                  {REVIEW_STATUS_TEXT[currentReview.status]}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="举报状态">
                {currentReview.isReported ? (
                  <Badge status="error" text="被举报" />
                ) : (
                  <Badge status="success" text="正常" />
                )}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间" span={2}>
                {dayjs(currentReview.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              <Descriptions.Item label="评价内容" span={2}>
                {currentReview.comment || '-'}
              </Descriptions.Item>
            </Descriptions>

            {/* 评价图片 */}
            {currentReview.images && currentReview.images.length > 0 && (
              <>
                <Divider styles={{ content: { margin: 0 } }}>评价图片</Divider>
                <Image.PreviewGroup>
                  <Space wrap>
                    {currentReview.images.map((url, index) => (
                      <Image
                        key={index}
                        width={100}
                        height={100}
                        src={url}
                        style={{ objectFit: 'cover' }}
                      />
                    ))}
                  </Space>
                </Image.PreviewGroup>
              </>
            )}

            {/* 操作日志 */}
            <Divider styles={{ content: { margin: 0 } }}>操作日志</Divider>
            {reviewLogs.length > 0 ? (
              <Timeline
                items={reviewLogs.map(log => {
                  const logData = log as OperationLog & { reason?: string; actorUserId?: number };
                  return {
                    children: (
                      <div>
                        <Text strong>{String(log.action)}</Text>
                        <br />
                        <Text type="secondary">
                          {log.actorName || (logData.actorUserId ? `管理员ID: ${logData.actorUserId}` : '系统')} - {dayjs(log.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                        </Text>
                        {(log.note || logData.reason) && (
                          <div><Text type="secondary">备注: {log.note || logData.reason}</Text></div>
                        )}
                      </div>
                    ),
                  };
                })}
              />
            ) : (
              <Text type="secondary">暂无操作日志</Text>
            )}

            {/* 操作按钮 */}
            <Divider />
            <Space>
              {currentReview.status === 'pending' && (
                <>
                  <Button
                    type="primary"
                    icon={<CheckOutlined />}
                    onClick={() => {
                      handleApprove(currentReview.id);
                      handleCloseDetail();
                    }}
                  >
                    批准
                  </Button>
                  <Button
                    danger
                    icon={<CloseOutlined />}
                    onClick={() => {
                      message.info('请前往审核页面进行拒绝操作');
                    }}
                  >
                    拒绝
                  </Button>
                </>
              )}
              <Popconfirm
                title="确定要删除这条评价吗？"
                onConfirm={() => {
                  handleDelete(currentReview.id);
                  handleCloseDetail();
                }}
                okText="确定"
                cancelText="取消"
              >
                <Button danger icon={<DeleteOutlined />}>
                  删除
                </Button>
              </Popconfirm>
            </Space>
          </Spin>
        )}
      </Drawer>
    </Card>
  );
};

export default ReviewList;
