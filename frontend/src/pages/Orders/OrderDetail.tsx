import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Tag, Button } from '../../components';
import { orderApi, OrderDetailType } from '../../services/api/order';
import { OrderStatus } from '../../types/order';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatCurrency,
  formatDateTime,
  formatRelativeTime,
} from '../../utils/formatters';
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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [orderDetail, setOrderDetail] = useState<OrderDetailType | null>(null);

  // 加载订单详情
  useEffect(() => {
    const loadOrderDetail = async () => {
      if (!id) {
        setError('订单ID无效');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        const data = await orderApi.getDetail(Number(id));
        setOrderDetail(data);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '加载订单详情失败';
        setError(errorMessage);
        console.error('加载订单详情失败:', err);
      } finally {
        setLoading(false);
      }
    };

    loadOrderDetail();
  }, [id]);

  // 加载中状态
  if (loading) {
    return (
      <div className={styles.container}>
        <p>加载中...</p>
      </div>
    );
  }

  // 错误状态
  if (error || !orderDetail) {
    return (
      <div className={styles.container}>
        <Card className={styles.errorCard}>
          <h2>订单未找到</h2>
          <p>{error || `订单 ID: ${id} 不存在`}</p>
          <Button onClick={() => navigate('/orders')}>返回订单列表</Button>
        </Card>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* 头部 */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <Button variant="primary" onClick={() => navigate('/orders')}>
            ← 返回列表
          </Button>
          <h1 className={styles.title}>订单详情</h1>
        </div>
        <div className={styles.headerRight}>
          <Tag color={getOrderStatusColor(orderDetail.status)}>
            {formatOrderStatus(orderDetail.status)}
          </Tag>
        </div>
      </div>

      {/* 主要内容 */}
      <div className={styles.content}>
        {/* 基本信息 */}
        <Card className={styles.section}>
          <h2 className={styles.sectionTitle}>基本信息</h2>
          <div className={styles.infoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>订单标题</span>
              <span className={styles.infoValue}>{orderDetail.title}</span>
            </div>
            {orderDetail.description && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>订单描述</span>
                <span className={styles.infoValue}>{orderDetail.description}</span>
              </div>
            )}
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>订单金额</span>
              <span className={`${styles.infoValue} ${styles.price}`}>
                {formatCurrency(orderDetail.price_cents)}
              </span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>创建时间</span>
              <span className={styles.infoValue}>{formatDateTime(orderDetail.created_at)}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>更新时间</span>
              <span className={styles.infoValue}>{formatDateTime(orderDetail.updated_at)}</span>
            </div>
            {orderDetail.scheduled_start && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>计划开始</span>
                <span className={styles.infoValue}>
                  {formatDateTime(orderDetail.scheduled_start)}
                </span>
              </div>
            )}
            {orderDetail.scheduled_end && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>计划结束</span>
                <span className={styles.infoValue}>
                  {formatDateTime(orderDetail.scheduled_end)}
                </span>
              </div>
            )}
            {orderDetail.cancel_reason && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>取消原因</span>
                <span className={styles.infoValue}>{orderDetail.cancel_reason}</span>
              </div>
            )}
          </div>
        </Card>

        {/* 用户信息 */}
        {orderDetail.user && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              <UserIcon />
              用户信息
            </h2>
            <div className={styles.userInfo}>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>用户名</span>
                <span className={styles.infoValue}>{orderDetail.user.name}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>用户ID</span>
                <span className={styles.infoValue}>{orderDetail.user.id}</span>
              </div>
            </div>
          </Card>
        )}

        {/* 陪玩师信息 */}
        {orderDetail.player && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              <UserIcon />
              陪玩师信息
            </h2>
            <div className={styles.playerInfo}>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>昵称</span>
                <span className={styles.infoValue}>{orderDetail.player.nickname || '-'}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>陪玩师ID</span>
                <span className={styles.infoValue}>{orderDetail.player.id}</span>
              </div>
            </div>
          </Card>
        )}

        {/* 游戏信息 */}
        {orderDetail.game && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>游戏信息</h2>
            <div className={styles.gameInfo}>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>游戏名称</span>
                <span className={styles.infoValue}>{orderDetail.game.name}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>游戏ID</span>
                <span className={styles.infoValue}>{orderDetail.game.id}</span>
              </div>
            </div>
          </Card>
        )}

        {/* 操作日志 */}
        {orderDetail.logs && orderDetail.logs.length > 0 && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>
              <ClockIcon />
              操作日志
            </h2>
            <div className={styles.timeline}>
              {orderDetail.logs.map((log) => (
                <div key={log.id} className={styles.timelineItem}>
                  <div className={styles.timelineDot} />
                  <div className={styles.timelineContent}>
                    <div className={styles.logHeader}>
                      <span className={styles.logAction}>{log.action}</span>
                      <span className={styles.logTime}>
                        {formatRelativeTime(log.created_at)}
                      </span>
                    </div>
                    <div className={styles.logDetails}>
                      <span className={styles.logOperator}>{log.operator_name}</span>
                      {log.note && <span className={styles.logNote}>{log.note}</span>}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        )}

        {/* 审核记录 */}
        {orderDetail.reviews && orderDetail.reviews.length > 0 && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>审核记录</h2>
            <div className={styles.reviewsList}>
              {orderDetail.reviews.map((review) => (
                <div key={review.id} className={styles.reviewItem}>
                  <div className={styles.reviewHeader}>
                    <Tag color={review.result === 'approved' ? 'success' : 'error'}>
                      {review.result === 'approved' ? '审核通过' : '审核拒绝'}
                    </Tag>
                    <span className={styles.reviewTime}>
                      {formatDateTime(review.created_at)}
                    </span>
                  </div>
                  <div className={styles.reviewContent}>
                    <div className={styles.reviewItem}>
                      <span className={styles.infoLabel}>审核人:</span>
                      <span className={styles.infoValue}>{review.reviewer_name}</span>
                    </div>
                    {review.reason && (
                      <div className={styles.reviewItem}>
                        <span className={styles.infoLabel}>原因:</span>
                        <span className={styles.infoValue}>{review.reason}</span>
                      </div>
                    )}
                    {review.comment && (
                      <div className={styles.reviewItem}>
                        <span className={styles.infoLabel}>备注:</span>
                        <span className={styles.infoValue}>{review.comment}</span>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </Card>
        )}

        {/* 操作区域 */}
        <Card className={styles.section}>
          <h2 className={styles.sectionTitle}>订单操作</h2>
          <div className={styles.actions}>
            {orderDetail.status === OrderStatus.PENDING && (
              <Button
                variant="primary"
                onClick={() => console.log('确认订单')}
              >
                确认订单
              </Button>
            )}
            {orderDetail.status === OrderStatus.CONFIRMED && (
              <Button
                variant="primary"
                onClick={() => console.log('开始服务')}
              >
                开始服务
              </Button>
            )}
            {orderDetail.status === OrderStatus.IN_PROGRESS && (
              <Button
                variant="primary"
                onClick={() => console.log('完成订单')}
              >
                完成订单
              </Button>
            )}
            {[OrderStatus.PENDING, OrderStatus.CONFIRMED].includes(orderDetail.status) && (
              <Button
                variant="secondary"
                onClick={() => console.log('取消订单')}
              >
                取消订单
              </Button>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
};
