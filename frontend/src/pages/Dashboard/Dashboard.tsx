import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Tag, Badge } from '../../components';
import { orderApi } from '../../services/api/order';
import type { Order, OrderStatistics } from '../../types/order';
import { OrderStatus } from '../../types/order';
import {
  formatOrderStatus,
  getOrderStatusColor,
  formatCurrency,
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
  const [statistics, setStatistics] = useState<OrderStatistics>({
    total: 0,
    pending: 0,
    confirmed: 0,
    in_progress: 0,
    completed: 0,
    canceled: 0,
    refunded: 0,
    today_orders: 0,
    today_revenue: 0,
  });
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  // 加载数据
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);

        // 加载统计数据
        const stats = await orderApi.getStatistics();
        setStatistics(stats);

        // 加载最近订单
        const ordersResult = await orderApi.getList({
          page: 1,
          page_size: 5,
          sort_by: 'created_at',
          sort_order: 'desc',
        });
        setRecentOrders(ordersResult.list || []);
      } catch (error) {
        console.error('加载仪表盘数据失败:', error);
        // 确保即使出错也设置为空数组
        setRecentOrders([]);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  if (loading) {
    return (
      <div className={styles.container}>
        <h1 className={styles.title}>仪表盘</h1>
        <p>加载中...</p>
      </div>
    );
  }

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
            <div className={styles.statValue}>{statistics.total}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <OrderIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>今日订单</div>
            <div className={styles.statValue}>{statistics.today_orders}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <OrderIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>进行中</div>
            <div className={styles.statValue}>
              {statistics.in_progress > 0 && <Badge count={statistics.in_progress} />}
              {statistics.in_progress}
            </div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <MoneyIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>今日收入</div>
            <div className={styles.statValue}>{formatCurrency(statistics.today_revenue)}</div>
          </div>
        </Card>
      </div>

      {/* 快捷入口 */}
      <div className={styles.quickActions}>
        <h2 className={styles.sectionTitle}>快捷入口</h2>
        <div className={styles.actionsGrid}>
          <Card className={styles.actionCard}>
            <div className={styles.actionContent} onClick={() => navigate('/orders')}>
              <div className={styles.actionIcon}>
                <OrderIcon />
              </div>
              <div className={styles.actionTitle}>所有订单</div>
              <div className={styles.actionDesc}>查看和管理所有订单</div>
            </div>
          </Card>

          <Card className={styles.actionCard}>
            <div
              className={styles.actionContent}
              onClick={() => navigate(`/orders?status=${OrderStatus.PENDING}`)}
            >
              <div className={styles.actionIcon}>
                <OrderIcon />
              </div>
              <div className={styles.actionTitle}>
                待处理订单
                {statistics.pending > 0 && <Badge count={statistics.pending} />}
              </div>
              <div className={styles.actionDesc}>处理需要确认的订单</div>
            </div>
          </Card>

          <Card className={styles.actionCard}>
            <div
              className={styles.actionContent}
              onClick={() => navigate(`/orders?status=${OrderStatus.IN_PROGRESS}`)}
            >
              <div className={styles.actionIcon}>
                <StatsIcon />
              </div>
              <div className={styles.actionTitle}>进行中订单</div>
              <div className={styles.actionDesc}>监控正在进行的订单</div>
            </div>
          </Card>

          <Card className={styles.actionCard}>
            <div className={styles.actionContent} onClick={() => navigate('/users')}>
              <div className={styles.actionIcon}>
                <MoneyIcon />
              </div>
              <div className={styles.actionTitle}>用户管理</div>
              <div className={styles.actionDesc}>管理用户和陪玩师</div>
            </div>
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
          {!recentOrders || recentOrders.length === 0 ? (
            <p>暂无订单</p>
          ) : (
            recentOrders.map((order) => (
              <Card className={styles.orderCard} key={order.id}>
                <div
                  className={styles.orderContent}
                  onClick={() => navigate(`/orders/${order.id}`)}
                >
                  <div className={styles.orderHeader}>
                    <div className={styles.orderTitle}>{order.title}</div>
                    <Tag color={getOrderStatusColor(order.status)}>
                      {formatOrderStatus(order.status)}
                    </Tag>
                  </div>

                  <div className={styles.orderInfo}>
                    {order.user && (
                      <div className={styles.infoRow}>
                        <span className={styles.infoLabel}>用户:</span>
                        <span className={styles.infoValue}>{order.user.name}</span>
                      </div>
                    )}
                    {order.game && (
                      <div className={styles.infoRow}>
                        <span className={styles.infoLabel}>游戏:</span>
                        <span className={styles.infoValue}>{order.game.name}</span>
                      </div>
                    )}
                    <div className={styles.infoRow}>
                      <span className={styles.infoLabel}>金额:</span>
                      <span className={`${styles.infoValue} ${styles.price}`}>
                        {formatCurrency(order.price_cents)}
                      </span>
                    </div>
                  </div>

                  <div className={styles.orderFooter}>
                    <span className={styles.orderTime}>{formatRelativeTime(order.created_at)}</span>
                  </div>
                </div>
              </Card>
            ))
          )}
        </div>
      </div>
    </div>
  );
};
