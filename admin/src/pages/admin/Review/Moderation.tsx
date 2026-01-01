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
  Modal,
  Input,
  Image,
  Typography,
  Alert,
  Switch,
  App,
} from 'antd';
import {
  CheckOutlined,
  CloseOutlined,
  ReloadOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import dayjs from 'dayjs';
import { reviewApi } from '@/api/review';
import type { Review } from '@/types/review';
import SensitiveWordHighlight from './components/SensitiveWordHighlight';

const { TextArea } = Input;
const { Text } = Typography;

const ReviewModeration: React.FC = () => {
  const { message, modal } = App.useApp();
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

  // 只显示含敏感词的评价
  const [onlySensitive, setOnlySensitive] = useState(false);

  // 全部通过加载状态
  const [approveAllLoading, setApproveAllLoading] = useState(false);

  // 加载待审核评价
  const fetchPendingReviews = useCallback(async () => {
    setLoading(true);
    try {
      const response = await reviewApi.getPendingReviews({
        page: pagination.current,
        pageSize: pagination.pageSize,
      });
      if (response.data.success) {
        setReviews(response.data.data || []);
        if (response.data.pagination) {
          setPagination(prev => ({
            ...prev,
            total: response.data.pagination!.total,
          }));
        }
      }
    } catch {
      message.error('获取待审核评价失败');
    } finally {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.current, pagination.pageSize]);

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
      const response = await reviewApi.approveReview(id);
      if (response.data.success) {
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
      const response = await reviewApi.rejectReview(currentReviewId, rejectReason);
      if (response.data.success) {
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
    modal.confirm({
      title: '批量批准',
      icon: <ExclamationCircleOutlined />,
      content: `确定要批准选中的 ${selectedRowKeys.length} 条评价吗？`,
      onOk: async () => {
        try {
          const response = await reviewApi.batchApproveReviews(
            selectedRowKeys.map(k => Number(k))
          );
          if (response.data.success) {
            message.success(`成功批准 ${response.data.data.count} 条评价`);
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
      );
      if (response.data.success) {
        message.success(`成功拒绝 ${response.data.data.count} 条评价`);
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
      render: (comment: string, record: Review & { hasSensitiveWords?: boolean; sensitiveWords?: string[] }) => {
        if (record.hasSensitiveWords && record.sensitiveWords?.length) {
          return (
            <div>
              <SensitiveWordHighlight
                content={comment}
                detectedWords={record.sensitiveWords.map(word => ({
                  word,
                  category: 'other' as const,
                  severity: 'medium' as const,
                  positions: [],
                }))}
              />
              <Tag color="red" style={{ marginTop: 4 }}>
                含敏感词: {record.sensitiveWords.join(', ')}
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

  // 根据筛选条件过滤评价列表
  const filteredReviews = onlySensitive
    ? reviews.filter((review: Review & { hasSensitiveWords?: boolean }) => review.hasSensitiveWords)
    : reviews;

  // 计算不含敏感词的评价数量
  const nonSensitiveCount = reviews.filter((review: Review & { hasSensitiveWords?: boolean }) => !review.hasSensitiveWords).length;

  // 全部通过（不含敏感词的评价）
  const handleApproveAllNonSensitive = async () => {
    if (nonSensitiveCount === 0) {
      message.warning('没有不含敏感词的评价需要批准');
      return;
    }
    modal.confirm({
      title: '全部通过',
      icon: <ExclamationCircleOutlined />,
      content: `确定要批准所有不含敏感词的 ${nonSensitiveCount} 条评价吗？`,
      onOk: async () => {
        setApproveAllLoading(true);
        try {
          const response = await reviewApi.approveAllNonSensitive();
          if (response.data.success) {
            message.success(`成功批准 ${response.data.data?.count || nonSensitiveCount} 条评价`);
            fetchPendingReviews();
          }
        } catch {
          message.error('操作失败');
        } finally {
          setApproveAllLoading(false);
        }
      },
    });
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
      <Space style={{ marginBottom: 16 }} wrap>
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
        <Button
          type="primary"
          ghost
          icon={<CheckOutlined />}
          onClick={handleApproveAllNonSensitive}
          loading={approveAllLoading}
          disabled={nonSensitiveCount === 0}
        >
          全部通过 ({nonSensitiveCount})
        </Button>
        <Button icon={<ReloadOutlined />} onClick={fetchPendingReviews}>
          刷新
        </Button>
        <Space>
          <Switch
            checked={onlySensitive}
            onChange={setOnlySensitive}
          />
          <span>只显示含敏感词的评价</span>
          {onlySensitive && (
            <Tag color="red">
              {filteredReviews.length} 条
            </Tag>
          )}
        </Space>
      </Space>

      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={filteredReviews}
        rowKey="id"
        loading={loading}
        rowSelection={rowSelection}
        pagination={onlySensitive ? {
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条含敏感词`,
        } : {
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
