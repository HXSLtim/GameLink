/**
 * 评价管理 - 评价详情页面
 * 需求: 1.4, 8.1, 8.2, 9.2
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Descriptions,
  Rate,
  Tag,
  Image,
  Timeline,
  Button,
  Space,
  Spin,
  message,
  Modal,
  Input,
  List,
  Avatar,
  Popconfirm,
  Typography,
  Divider,
} from 'antd';
import {
  ArrowLeftOutlined,
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import { reviewApi, reviewReplyApi, type OperationLog } from '@/api/review';
import type { Review, ReviewReply } from '@/types/review';
import {
  REVIEW_STATUS_TEXT,
  REVIEW_STATUS_COLOR,
} from '@/types/review';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

const ReviewDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  // 状态
  const [loading, setLoading] = useState(true);
  const [review, setReview] = useState<Review | null>(null);
  const [logs, setLogs] = useState<OperationLog[]>([]);
  const [replyContent, setReplyContent] = useState('');
  const [replyLoading, setReplyLoading] = useState(false);
  const [rejectModalVisible, setRejectModalVisible] = useState(false);
  const [rejectReason, setRejectReason] = useState('');

  // 加载评价详情
  const fetchReview = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    try {
      const response = await reviewApi.getReview(Number(id)) as unknown as {
        success: boolean;
        data: Review;
      };
      if (response.success) {
        setReview(response.data);
      }
    } catch {
      message.error('获取评价详情失败');
    } finally {
      setLoading(false);
    }
  }, [id]);

  // 加载操作日志
  const fetchLogs = useCallback(async () => {
    if (!id) return;
    try {
      const response = await reviewApi.getReviewLogs(Number(id)) as unknown as {
        success: boolean;
        data: OperationLog[];
      };
      if (response.success) {
        setLogs(response.data || []);
      }
    } catch {
      console.error('Failed to fetch logs');
    }
  }, [id]);

  useEffect(() => {
    fetchReview();
    fetchLogs();
  }, [fetchReview, fetchLogs]);

  // 批准评价
  const handleApprove = async () => {
    if (!id) return;
    try {
      const response = await reviewApi.approveReview(Number(id)) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已批准');
        fetchReview();
        fetchLogs();
      }
    } catch {
      message.error('操作失败');
    }
  };

  // 拒绝评价
  const handleReject = async () => {
    if (!id || !rejectReason.trim()) {
      message.warning('请输入拒绝原因');
      return;
    }
    try {
      const response = await reviewApi.rejectReview(Number(id), rejectReason) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已拒绝');
        setRejectModalVisible(false);
        setRejectReason('');
        fetchReview();
        fetchLogs();
      }
    } catch {
      message.error('操作失败');
    }
  };

  // 删除评价
  const handleDelete = async () => {
    if (!id) return;
    try {
      const response = await reviewApi.deleteReview(Number(id)) as unknown as { success: boolean };
      if (response.success) {
        message.success('评价已删除');
        navigate('/admin/reviews');
      }
    } catch {
      message.error('删除失败');
    }
  };

  // 提交回复
  const handleReply = async () => {
    if (!id || !replyContent.trim()) {
      message.warning('请输入回复内容');
      return;
    }
    setReplyLoading(true);
    try {
      const response = await reviewReplyApi.createReply(Number(id), replyContent) as unknown as { success: boolean };
      if (response.success) {
        message.success('回复成功');
        setReplyContent('');
        fetchReview();
        fetchLogs();
      }
    } catch {
      message.error('回复失败');
    } finally {
      setReplyLoading(false);
    }
  };

  // 删除回复
  const handleDeleteReply = async (replyId: number) => {
    try {
      const response = await reviewReplyApi.deleteReply(replyId) as unknown as { success: boolean };
      if (response.success) {
        message.success('回复已删除');
        fetchReview();
        fetchLogs();
      }
    } catch {
      message.error('删除失败');
    }
  };

  // 获取操作类型文本
  const getActionText = (action: string) => {
    const actionMap: Record<string, string> = {
      create: '创建评价',
      approve: '批准评价',
      reject: '拒绝评价',
      delete: '删除评价',
      reply: '回复评价',
      update_reply: '更新回复',
      delete_reply: '删除回复',
      report: '举报评价',
      handle_report: '处理举报',
    };
    return actionMap[action] || action;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!review) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: 50 }}>
          <Text type="secondary">评价不存在</Text>
          <br />
          <Button type="link" onClick={() => navigate('/admin/reviews')}>
            返回列表
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <div>
      {/* 返回按钮 */}
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/admin/reviews')}
        style={{ marginBottom: 16 }}
      >
        返回列表
      </Button>

      {/* 评价基本信息 */}
      <Card title="评价详情" style={{ marginBottom: 16 }}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="评价ID">{review.id}</Descriptions.Item>
          <Descriptions.Item label="订单ID">{review.orderId}</Descriptions.Item>
          <Descriptions.Item label="评价者">
            {review.reviewerName || `用户${review.reviewerId}`}
          </Descriptions.Item>
          <Descriptions.Item label="被评价者">
            {review.playerName || `陪玩师${review.playerId}`}
          </Descriptions.Item>
          <Descriptions.Item label="评分">
            <Rate disabled value={review.rating} />
            <Text style={{ marginLeft: 8 }}>{review.rating}分</Text>
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={REVIEW_STATUS_COLOR[review.status]}>
              {REVIEW_STATUS_TEXT[review.status]}
            </Tag>
            {review.isReported && (
              <Tag color="red" style={{ marginLeft: 8 }}>被举报</Tag>
            )}
          </Descriptions.Item>
          <Descriptions.Item label="评价内容" span={2}>
            <Paragraph>{review.comment || '-'}</Paragraph>
          </Descriptions.Item>
          {review.rejectionReason && (
            <Descriptions.Item label="拒绝原因" span={2}>
              <Text type="danger">{review.rejectionReason}</Text>
            </Descriptions.Item>
          )}
          <Descriptions.Item label="创建时间">
            {dayjs(review.createdAt).format('YYYY-MM-DD HH:mm:ss')}
          </Descriptions.Item>
          <Descriptions.Item label="更新时间">
            {dayjs(review.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
          </Descriptions.Item>
        </Descriptions>

        {/* 评价图片 */}
        {review.images && review.images.length > 0 && (
          <>
            <Divider>评价图片</Divider>
            <Image.PreviewGroup>
              <Space wrap>
                {review.images.map((url, index) => (
                  <Image
                    key={index}
                    width={120}
                    height={120}
                    src={url}
                    style={{ objectFit: 'cover', borderRadius: 8 }}
                  />
                ))}
              </Space>
            </Image.PreviewGroup>
          </>
        )}

        {/* 操作按钮 */}
        {review.status === 'pending' && (
          <>
            <Divider />
            <Space>
              <Button
                type="primary"
                icon={<CheckOutlined />}
                onClick={handleApprove}
              >
                批准
              </Button>
              <Button
                danger
                icon={<CloseOutlined />}
                onClick={() => setRejectModalVisible(true)}
              >
                拒绝
              </Button>
              <Popconfirm
                title="确定要删除这条评价吗？"
                onConfirm={handleDelete}
                okText="确定"
                cancelText="取消"
              >
                <Button danger icon={<DeleteOutlined />}>
                  删除
                </Button>
              </Popconfirm>
            </Space>
          </>
        )}
      </Card>

      {/* 回复列表 */}
      <Card title="回复记录" style={{ marginBottom: 16 }}>
        <List
          dataSource={review.replies || []}
          locale={{ emptyText: '暂无回复' }}
          renderItem={(reply: ReviewReply) => (
            <List.Item
              actions={[
                <Button
                  key="edit"
                  type="link"
                  size="small"
                  icon={<EditOutlined />}
                  onClick={() => message.info('编辑功能开发中')}
                >
                  编辑
                </Button>,
                <Popconfirm
                  key="delete"
                  title="确定要删除这条回复吗？"
                  onConfirm={() => handleDeleteReply(reply.id)}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                    删除
                  </Button>
                </Popconfirm>,
              ]}
            >
              <List.Item.Meta
                avatar={<Avatar icon={<UserOutlined />} />}
                title={
                  <Space>
                    <Text strong>{reply.userName || `用户${reply.userId}`}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {dayjs(reply.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                    </Text>
                  </Space>
                }
                description={reply.content}
              />
            </List.Item>
          )}
        />

        {/* 回复输入框 */}
        <Divider />
        <Space.Compact style={{ width: '100%' }}>
          <TextArea
            placeholder="输入回复内容..."
            value={replyContent}
            onChange={e => setReplyContent(e.target.value)}
            rows={2}
            style={{ width: 'calc(100% - 80px)' }}
          />
          <Button
            type="primary"
            loading={replyLoading}
            onClick={handleReply}
            style={{ height: 'auto' }}
          >
            回复
          </Button>
        </Space.Compact>
      </Card>

      {/* 操作日志 */}
      <Card title="操作历史">
        <Timeline
          items={logs.map(log => ({
            color: log.action === 'reject' || log.action === 'delete' ? 'red' : 'blue',
            children: (
              <div>
                <Text strong>{getActionText(log.action)}</Text>
                <br />
                <Text type="secondary">
                  操作人: {log.actorName || `用户${log.actorId}`}
                </Text>
                <br />
                <Text type="secondary">
                  时间: {dayjs(log.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                </Text>
                {log.note && (
                  <>
                    <br />
                    <Text type="secondary">备注: {log.note}</Text>
                  </>
                )}
              </div>
            ),
          }))}
        />
        {logs.length === 0 && (
          <Text type="secondary">暂无操作记录</Text>
        )}
      </Card>

      {/* 拒绝弹窗 */}
      <Modal
        title="拒绝评价"
        open={rejectModalVisible}
        onOk={handleReject}
        onCancel={() => {
          setRejectModalVisible(false);
          setRejectReason('');
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
    </div>
  );
};

export default ReviewDetail;
