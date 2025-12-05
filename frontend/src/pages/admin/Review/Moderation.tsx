/**
 * 评价管理 - 评价审核页面
 * 需求: 2.1, 2.2, 2.3, 2.4, 2.5
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Button,
  Space,
  Tag,
  Rate,
  message,
  Modal,
  Input,
  Image,
  Typography,
  Alert,
} from 'antd';
import {
  CheckOutlined,
  CloseOutlined,
  ReloadOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { reviewApi, sensitiveWordApi } from '@/api/review';
import type { Review, DetectSensitiveWordsResult } from '@/types/review';
import SensitiveWordHighlight from './components/SensitiveWordHighlight';

const { TextArea } = Input;
const { Text } = Typography;

const ReviewModeration: React.FC = () => {
  // 状态
  const [loading, setLoading] = useState(false);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });

  // 弹窗状态
  const [rejectModalVisible, setRejectModalVisible] = useState(false);
  const [batchRejectModalVisible, setBatchRejectModalVisible] = useState(false);
  const [currentReviewId, setCurrentReviewId] = useState<number | null>(null);
  const [rejectReason, setRejectReason] = useState('');

  // 敏感词检测结果缓存
  const [sensitiveResults, setSensitiveResults] = useState<Record<number, DetectSensitiveWordsResult>>({});

  // 加载待审核评价
  const fetchPendingReviews = useCallback(async () => {
    setLoading(true);
    try {
      const response = await reviewApi.getPendingReviews({
        page: pagination.current,
        pageSize: pagination.pageSize,
      }) as unknown as {
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
        // 检测敏感词
        detectSensitiveWords(response.data || []);
      }
    } catch {
      message.error('获取待审核评价失败');
    } finally {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.current, pagination.pageSize]);

  // 批量检测敏感词
  const detectSensitiveWords = async (reviewList: Review[]) => {
    const results: Record<number, DetectSensitiveWordsResult> = {};
    for (const review of reviewList) {
      if (review.comment) {
        try {
          const response = await sensitiveWordApi.detectWords(review.comment) as unknown as {
            success: boolean;
            data: DetectSensitiveWordsResult;
          };
          if (response.success) {
            results[review.id] = response.data;
          }
        } catch {
          // 忽略检测失败
        }
      }
    }
    setSensitiveResults(results);
  };

  useEffect(() => {
    fetchPendingReviews();
  }, [fetchPendingReviews]);

  // 表格变化
  const handleTableChange = (paginationConfig: TablePaginationConfig) => {
    setPagination(prev => ({
      ...prev,
      current: paginationConfig.current || 1,
      pageSize: paginationConfig.pageSize || 20,
    }));
  };

  // 批准单个评价
  const handleApprove = async (id: number) => {
    try {
      const response = await reviewApi.approveReview(id) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已批准');
        fetchPendingReviews();
      }
    } catch {
      message.error('操作失败');
    }
  };

  // 打开拒绝弹窗
  const openRejectModal = (id: number) => {
    setCurrentReviewId(id);
    setRejectReason('');
    setRejectModalVisible(true);
  };

  // 拒绝单个评价
  const handleReject = async () => {
    if (!currentReviewId || !rejectReason.trim()) {
      message.warning('请输入拒绝原因');
      return;
    }
    try {
      const response = await reviewApi.rejectReview(currentReviewId, rejectReason) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已拒绝');
        setRejectModalVisible(false);
        setRejectReason('');
        setCurrentReviewId(null);
        fetchPendingReviews();
      }
    } catch {
      message.error('操作失败');
    }
  };

  // 批量批准
  const handleBatchApprove = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要批准的评价');
      return;
    }
    Modal.confirm({
      title: '批量批准',
      icon: <ExclamationCircleOutlined />,
      content: `确定要批准选中的 ${selectedRowKeys.length} 条评价吗？`,
      onOk: async () => {
        try {
          const response = await reviewApi.batchApproveReviews(
            selectedRowKeys.map(k => Number(k))
          ) as unknown as { success: boolean; data: { count: number } };
          if (response.success) {
            message.success(`成功批准 ${response.data.count} 条评价`);
            setSelectedRowKeys([]);
            fetchPendingReviews();
          }
        } catch {
          message.error('批量批准失败');
        }
      },
    });
  };

  // 批量拒绝
  const handleBatchReject = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要拒绝的评价');
      return;
    }
    if (!rejectReason.trim()) {
      message.warning('请输入拒绝原因');
      return;
    }
    try {
      const response = await reviewApi.batchRejectReviews(
        selectedRowKeys.map(k => Number(k)),
        rejectReason
      ) as unknown as { success: boolean; data: { count: number } };
      if (response.success) {
        message.success(`成功拒绝 ${response.data.count} 条评价`);
        setSelectedRowKeys([]);
        setBatchRejectModalVisible(false);
        setRejectReason('');
        fetchPendingReviews();
      }
    } catch {
      message.error('批量拒绝失败');
    }
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
      width: 300,
      render: (comment: string, record: Review) => {
        const sensitiveResult = sensitiveResults[record.id];
        if (sensitiveResult?.hasSensitiveWords) {
          return (
            <div>
              <SensitiveWordHighlight
                content={comment}
                detectedWords={sensitiveResult.detectedWords}
              />
              <Tag color="red" style={{ marginTop: 4 }}>
                含敏感词
              </Tag>
            </div>
          );
        }
        return <Text ellipsis={{ tooltip: comment }}>{comment || '-'}</Text>;
      },
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
      fixed: 'right',
      render: (_: unknown, record: Review) => (
        <Space size="small">
          <Button
            type="primary"
            size="small"
            icon={<CheckOutlined />}
            onClick={() => handleApprove(record.id)}
          >
            批准
          </Button>
          <Button
            danger
            size="small"
            icon={<CloseOutlined />}
            onClick={() => openRejectModal(record.id)}
          >
            拒绝
          </Button>
        </Space>
      ),
    },
  ];

  // 行选择配置
  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
  };

  return (
    <Card title="评价审核">
      {/* 提示信息 */}
      <Alert
        title="待审核评价"
        description="以下评价需要审核后才能在前端展示。包含敏感词的评价会被自动标记。"
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
      />

      {/* 批量操作按钮 */}
      <Space style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          icon={<CheckOutlined />}
          onClick={handleBatchApprove}
          disabled={selectedRowKeys.length === 0}
        >
          批量批准 {selectedRowKeys.length > 0 && `(${selectedRowKeys.length})`}
        </Button>
        <Button
          danger
          icon={<CloseOutlined />}
          onClick={() => setBatchRejectModalVisible(true)}
          disabled={selectedRowKeys.length === 0}
        >
          批量拒绝 {selectedRowKeys.length > 0 && `(${selectedRowKeys.length})`}
        </Button>
        <Button icon={<ReloadOutlined />} onClick={fetchPendingReviews}>
          刷新
        </Button>
      </Space>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={reviews}
        rowKey="id"
        loading={loading}
        rowSelection={rowSelection}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条待审核`,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1300 }}
      />

      {/* 单个拒绝弹窗 */}
      <Modal
        title="拒绝评价"
        open={rejectModalVisible}
        onOk={handleReject}
        onCancel={() => {
          setRejectModalVisible(false);
          setRejectReason('');
          setCurrentReviewId(null);
        }}
        okText="确定"
        cancelText="取消"
      >
        <TextArea
          placeholder="请输入拒绝原因"
          value={rejectReason}
          onChange={e => setRejectReason(e.target.value)}
          rows={4}
        />
      </Modal>

      {/* 批量拒绝弹窗 */}
      <Modal
        title={`批量拒绝 (${selectedRowKeys.length} 条)`}
        open={batchRejectModalVisible}
        onOk={handleBatchReject}
        onCancel={() => {
          setBatchRejectModalVisible(false);
          setRejectReason('');
        }}
        okText="确定"
        cancelText="取消"
      >
        <TextArea
          placeholder="请输入拒绝原因（将应用于所有选中的评价）"
          value={rejectReason}
          onChange={e => setRejectReason(e.target.value)}
          rows={4}
        />
      </Modal>
    </Card>
  );
};

export default ReviewModeration;
