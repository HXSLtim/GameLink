import React, { useMemo, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Tag, Button, Badge, ReviewModal } from '../../components';
import type { ReviewFormData } from '../../components/ReviewModal';
import { getOrderDetail } from '../../services/mockData';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatReviewStatus,
  getReviewStatusColor,
  formatGameType,
  formatServiceType,
  formatCurrency,
  formatDuration,
  formatDateTime,
  formatRelativeTime,
} from '../../utils/formatters';
import { OrderStatus, ReviewStatus, OrderActionType } from '../../types/order.types';
import styles from './OrderDetail.module.less';

const UserIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="12" cy="7" r="4" strokeWidth="2" />
  </svg>
);

const ClockIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <circle cx="12" cy="12" r="10" strokeWidth="2" />
    <path d="M12 6V12L16 14" strokeWidth="2" strokeLinecap="round" />
  </svg>
);

export const OrderDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [showReviewModal, setShowReviewModal] = useState(false);

  const orderDetail = useMemo(() => {
    if (!id) return null;
    return getOrderDetail(id);
  }, [id]);

  const handleReviewSubmit = async (data: ReviewFormData) => {
    console.log('审核提交:', data);
    // 这里可以调用API提交审核
    // await orderService.review(id, data);

    // 模拟API调用
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // 提示成功
    alert(data.result === 'approved' ? '审核通过！' : '审核拒绝！');

    // 可以刷新页面数据
    // 这里简单地关闭modal，实际应该重新获取订单数据
  };

  if (!orderDetail) {
    return (
      <div className={styles.container}>
        <Card className={styles.errorCard}>
          <h2>订单未找到</h2>
          <p>订单 ID: {id} 不存在</p>
          <Button onClick={() => navigate('/orders')}>返回订单列表</Button>
        </Card>
      </div>
    );
  }

  const { order, logs, reviews } = orderDetail;

  // 判断是否可以审核
  const canReview =
    order.status === OrderStatus.PENDING_REVIEW && order.reviewStatus === ReviewStatus.PENDING;

  // 获取操作图标
  const getActionIcon = (action: OrderActionType) => {
    const icons: Record<OrderActionType, string> = {
      [OrderActionType.CREATE]: '📝',
      [OrderActionType.PAY]: '💰',
      [OrderActionType.ACCEPT]: '✅',
      [OrderActionType.START]: '🎮',
      [OrderActionType.SUBMIT_REVIEW]: '📋',
      [OrderActionType.APPROVE]: '✔️',
      [OrderActionType.REJECT]: '❌',
      [OrderActionType.COMPLETE]: '🎉',
      [OrderActionType.CANCEL]: '🚫',
      [OrderActionType.REQUEST_REFUND]: '💸',
      [OrderActionType.REFUND]: '💵',
    };
    return icons[action] || '📌';
  };

  return (
    <div className={styles.container}>
      {/* 头部信息 */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <Button variant="outlined" onClick={() => navigate('/orders')}>
            ← 返回列表
          </Button>
          <h1 className={styles.title}>订单详情</h1>
        </div>
        <div className={styles.headerRight}>
          <Tag color={getOrderStatusColor(order.status)}>{formatOrderStatus(order.status)}</Tag>
          {order.reviewStatus && (
            <Tag color={getReviewStatusColor(order.reviewStatus)}>
              {formatReviewStatus(order.reviewStatus)}
            </Tag>
          )}
        </div>
      </div>

      {/* 主要内容区 */}
      <div className={styles.content}>
        {/* 左侧栏 */}
        <div className={styles.leftColumn}>
          {/* 订单基本信息 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>订单信息</h2>
            <div className={styles.infoGrid}>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>订单号</span>
                <span className={styles.infoValue}>{order.orderNo}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>创建时间</span>
                <span className={styles.infoValue}>{formatDateTime(order.createdAt)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>游戏类型</span>
                <span className={styles.infoValue}>{formatGameType(order.gameType)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>服务类型</span>
                <span className={styles.infoValue}>{formatServiceType(order.serviceType)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>服务时长</span>
                <span className={styles.infoValue}>{formatDuration(order.duration)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>订单金额</span>
                <span className={`${styles.infoValue} ${styles.price}`}>
                  {formatCurrency(order.price)}
                </span>
              </div>
            </div>

            {order.description && (
              <div className={styles.descriptionSection}>
                <h3 className={styles.subTitle}>订单描述</h3>
                <p className={styles.description}>{order.description}</p>
              </div>
            )}

            {order.requirements && (
              <div className={styles.descriptionSection}>
                <h3 className={styles.subTitle}>特殊要求</h3>
                <p className={styles.description}>{order.requirements}</p>
              </div>
            )}
          </Card>

          {/* 用户信息 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              <UserIcon />
              用户信息
            </h2>
            <div className={styles.userCard}>
              <div className={styles.userAvatar}>{order.user.username.charAt(0)}</div>
              <div className={styles.userInfo}>
                <div className={styles.userName}>{order.user.username}</div>
                {order.user.phone && <div className={styles.userMeta}>{order.user.phone}</div>}
                <div className={styles.userMeta}>用户 ID: {order.user.id}</div>
              </div>
            </div>
          </Card>

          {/* 陪玩者信息 */}
          {order.player && (
            <Card className={styles.section}>
              <h2 className={styles.sectionTitle}>
                <UserIcon />
                陪玩者信息
              </h2>
              <div className={styles.playerCard}>
                <div className={styles.playerAvatar}>{order.player.username.charAt(0)}</div>
                <div className={styles.playerInfo}>
                  <div className={styles.playerName}>{order.player.username}</div>
                  <div className={styles.playerMeta}>
                    ⭐ {order.player.rating} 分 · {order.player.completedOrders} 单
                  </div>
                  <div className={styles.playerMeta}>等级 {order.player.level}</div>
                  <div className={styles.playerTags}>
                    {order.player.tags.map((tag, index) => (
                      <Tag key={index} color="info">
                        {tag}
                      </Tag>
                    ))}
                  </div>
                </div>
              </div>
            </Card>
          )}

          {/* 时间节点 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              <ClockIcon />
              时间节点
            </h2>
            <div className={styles.timelineGrid}>
              {order.createdAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>创建时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.createdAt)}</span>
                </div>
              )}
              {order.paidAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>支付时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.paidAt)}</span>
                </div>
              )}
              {order.acceptedAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>接单时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.acceptedAt)}</span>
                </div>
              )}
              {order.startedAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>开始时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.startedAt)}</span>
                </div>
              )}
              {order.completedAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>完成时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.completedAt)}</span>
                </div>
              )}
              {order.cancelledAt && (
                <div className={styles.timeItem}>
                  <span className={styles.timeLabel}>取消时间</span>
                  <span className={styles.timeValue}>{formatDateTime(order.cancelledAt)}</span>
                </div>
              )}
            </div>
          </Card>
        </div>

        {/* 右侧栏 */}
        <div className={styles.rightColumn}>
          {/* 操作按钮 */}
          {canReview && (
            <Card className={styles.section}>
              <h2 className={styles.sectionTitle}>订单操作</h2>
              <div className={styles.actions}>
                <Button
                  variant="primary"
                  onClick={() => setShowReviewModal(true)}
                  className={styles.actionButton}
                >
                  📋 开始审核
                </Button>
              </div>
            </Card>
          )}

          {/* 审核记录 */}
          {reviews.length > 0 && (
            <Card className={styles.section}>
              <h2 className={styles.sectionTitle}>
                审核记录
                <Badge count={reviews.length} />
              </h2>
              <div className={styles.reviewList}>
                {reviews.map((review) => (
                  <div key={review.id} className={styles.reviewItem}>
                    <div className={styles.reviewHeader}>
                      <Tag color={getReviewStatusColor(review.status)}>
                        {formatReviewStatus(review.status)}
                      </Tag>
                      <span className={styles.reviewTime}>
                        {formatRelativeTime(review.createdAt)}
                      </span>
                    </div>
                    <div className={styles.reviewBody}>
                      <div className={styles.reviewMeta}>审核人: {review.reviewer}</div>
                      {review.reason && <div className={styles.reviewReason}>{review.reason}</div>}
                    </div>
                  </div>
                ))}
              </div>
            </Card>
          )}

          {/* 操作历史 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              操作历史
              <Badge count={logs.length} />
            </h2>
            <div className={styles.timeline}>
              {logs.map((log, index) => (
                <div key={log.id} className={styles.timelineItem}>
                  <div className={styles.timelineDot}>{getActionIcon(log.action)}</div>
                  <div className={styles.timelineContent}>
                    <div className={styles.timelineHeader}>
                      <span className={styles.timelineAction}>{log.content}</span>
                      <span className={styles.timelineTime}>
                        {formatRelativeTime(log.createdAt)}
                      </span>
                    </div>
                    <div className={styles.timelineMeta}>
                      {log.operator} · {log.operatorRole}
                    </div>
                    {(log.statusBefore || log.statusAfter) && (
                      <div className={styles.timelineStatus}>
                        {log.statusBefore && (
                          <Tag color={getOrderStatusColor(log.statusBefore)}>
                            {formatOrderStatus(log.statusBefore)}
                          </Tag>
                        )}
                        {log.statusBefore && log.statusAfter && <span>→</span>}
                        {log.statusAfter && (
                          <Tag color={getOrderStatusColor(log.statusAfter)}>
                            {formatOrderStatus(log.statusAfter)}
                          </Tag>
                        )}
                      </div>
                    )}
                  </div>
                  {index < logs.length - 1 && <div className={styles.timelineLine} />}
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>

      {/* 审核Modal */}
      <ReviewModal
        visible={showReviewModal}
        orderNo={order.orderNo}
        onClose={() => setShowReviewModal(false)}
        onSubmit={handleReviewSubmit}
      />
    </div>
  );
};
