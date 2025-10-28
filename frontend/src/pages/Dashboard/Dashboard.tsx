import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Tag, Badge } from '../../components';
import { mockStatistics, mockOrders } from '../../services/mockData';
import { OrderStatus, ReviewStatus } from '../../types/order.types';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatCurrency,
  formatGameType,
  formatServiceType,
  formatRelativeTime,
} from '../../utils/formatters';
import styles from './Dashboard.module.less';

const StatsIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <line x1="12" y1="20" x2="12" y2="10" strokeWidth="2" strokeLinecap="round" />
    <line x1="18" y1="20" x2="18" y2="4" strokeWidth="2" strokeLinecap="round" />
    <line x1="6" y1="20" x2="6" y2="16" strokeWidth="2" strokeLinecap="round" />
  </svg>
);

const OrderIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M9 11L12 14L22 4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M21 12V19C21 19.5304 20.7893 20.0391 20.4142 20.4142C20.0391 20.7893 19.5304 21 19 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H16"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const ReviewIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M9 11L12 14L22 4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M21 12V19C21 19.5304 20.7893 20.0391 20.4142 20.4142C20.0391 20.7893 19.5304 21 19 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H16"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="17" cy="6" r="3" strokeWidth="2" />
  </svg>
);

const MoneyIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <line x1="12" y1="1" x2="12" y2="23" strokeWidth="2" strokeLinecap="round" />
    <path
      d="M17 5H9.5C8.57174 5 7.6815 5.36875 7.02513 6.02513C6.36875 6.6815 6 7.57174 6 8.5C6 9.42826 6.36875 10.3185 7.02513 10.9749C7.6815 11.6313 8.57174 12 9.5 12H14.5C15.4283 12 16.3185 12.3687 16.9749 13.0251C17.6313 13.6815 18 14.5717 18 15.5C18 16.4283 17.6313 17.3185 16.9749 17.9749C16.3185 18.6313 15.4283 19 14.5 19H6"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  // 获取最近订单（前5个）
  const recentOrders = mockOrders.slice(0, 5);

  // 获取待审核订单
  const pendingReviews = mockOrders.filter((order) => order.reviewStatus === ReviewStatus.PENDING);

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>仪表盘</h1>

      {/* 统计卡片 */}
      <div className={styles.statsGrid}>
        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <StatsIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总订单数</div>
            <div className={styles.statValue}>{mockStatistics.total}</div>
            <div className={styles.statTrend}>
              <span className={styles.trendUp}>↑ 12% 较上月</span>
            </div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <OrderIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>今日订单</div>
            <div className={styles.statValue}>{mockStatistics.todayOrders}</div>
            <div className={styles.statTrend}>
              <span className={styles.trendUp}>↑ 8% 较昨日</span>
            </div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <ReviewIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>待审核</div>
            <div className={styles.statValue}>
              <Badge count={mockStatistics.pendingReview} />
              {mockStatistics.pendingReview}
            </div>
            <div className={styles.statTrend}>需要处理</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <MoneyIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>今日收入</div>
            <div className={styles.statValue}>{formatCurrency(mockStatistics.todayRevenue)}</div>
            <div className={styles.statTrend}>
              <span className={styles.trendUp}>↑ 15% 较昨日</span>
            </div>
          </div>
        </Card>
      </div>

      {/* 快捷入口 */}
      <div className={styles.quickActions}>
        <h2 className={styles.sectionTitle}>快捷入口</h2>
        <div className={styles.actionsGrid}>
          <Card className={styles.actionCard} onClick={() => navigate('/orders')}>
            <div className={styles.actionIcon}>
              <OrderIcon />
            </div>
            <div className={styles.actionTitle}>所有订单</div>
            <div className={styles.actionDesc}>查看和管理所有订单</div>
          </Card>

          <Card
            className={styles.actionCard}
            onClick={() => navigate('/orders?status=pending_review')}
          >
            <div className={styles.actionIcon}>
              <ReviewIcon />
            </div>
            <div className={styles.actionTitle}>
              待审核订单
              {mockStatistics.pendingReview > 0 && <Badge count={mockStatistics.pendingReview} />}
            </div>
            <div className={styles.actionDesc}>处理需要审核的订单</div>
          </Card>

          <Card
            className={styles.actionCard}
            onClick={() => navigate('/orders?status=in_progress')}
          >
            <div className={styles.actionIcon}>
              <StatsIcon />
            </div>
            <div className={styles.actionTitle}>进行中订单</div>
            <div className={styles.actionDesc}>监控正在进行的订单</div>
          </Card>

          <Card className={styles.actionCard} onClick={() => navigate('/orders')}>
            <div className={styles.actionIcon}>
              <MoneyIcon />
            </div>
            <div className={styles.actionTitle}>财务报表</div>
            <div className={styles.actionDesc}>查看收入和统计数据</div>
          </Card>
        </div>
      </div>

      {/* 最近订单 */}
      <div className={styles.recentOrders}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>最近订单</h2>
          <button className={styles.viewAllButton} onClick={() => navigate('/orders')}>
            查看全部 →
          </button>
        </div>

        <div className={styles.ordersList}>
          {recentOrders.map((order) => (
            <Card
              key={order.id}
              className={styles.orderCard}
              onClick={() => navigate(`/orders/${order.id}`)}
            >
              <div className={styles.orderHeader}>
                <div className={styles.orderNo}>{order.orderNo}</div>
                <Tag color={getOrderStatusColor(order.status)}>
                  {formatOrderStatus(order.status)}
                </Tag>
              </div>

              <div className={styles.orderInfo}>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>用户:</span>
                  <span className={styles.infoValue}>{order.user.username}</span>
                </div>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>游戏:</span>
                  <span className={styles.infoValue}>{formatGameType(order.gameType)}</span>
                </div>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>服务:</span>
                  <span className={styles.infoValue}>{formatServiceType(order.serviceType)}</span>
                </div>
                <div className={styles.infoRow}>
                  <span className={styles.infoLabel}>金额:</span>
                  <span className={`${styles.infoValue} ${styles.price}`}>
                    {formatCurrency(order.price)}
                  </span>
                </div>
              </div>

              <div className={styles.orderFooter}>
                <span className={styles.orderTime}>{formatRelativeTime(order.createdAt)}</span>
              </div>
            </Card>
          ))}
        </div>
      </div>

      {/* 待审核订单 */}
      {pendingReviews.length > 0 && (
        <div className={styles.pendingReviews}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>
              待审核订单
              <Badge count={pendingReviews.length} />
            </h2>
            <button
              className={styles.viewAllButton}
              onClick={() => navigate('/orders?reviewStatus=pending')}
            >
              查看全部 →
            </button>
          </div>

          <div className={styles.reviewsList}>
            {pendingReviews.slice(0, 3).map((order) => (
              <Card
                key={order.id}
                className={styles.reviewCard}
                onClick={() => navigate(`/orders/${order.id}`)}
              >
                <div className={styles.reviewHeader}>
                  <div className={styles.orderNo}>{order.orderNo}</div>
                  <Tag color="warning">待审核</Tag>
                </div>

                <div className={styles.reviewContent}>
                  <div className={styles.reviewText}>
                    用户 <strong>{order.user.username}</strong> 的订单已完成，等待审核
                  </div>
                  <div className={styles.reviewMeta}>
                    <span>{formatGameType(order.gameType)}</span>
                    <span>•</span>
                    <span>{formatCurrency(order.price)}</span>
                    <span>•</span>
                    <span>{formatRelativeTime(order.updatedAt)}</span>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
