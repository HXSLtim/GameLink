import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Tag, Badge } from '../../components';
import { orderApi } from '../../services/api/order';
import { statsApi } from '../../services/api/stats';
import type { Order } from '../../types/order';
import type { DashboardStats } from '../../types/stats';
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

const UserIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="12" cy="7" r="4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

const PlayerIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="9" cy="7" r="4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M23 21V19C23 18.1645 22.7275 17.3645 22.2297 16.7262C21.7318 16.0879 21.0415 15.6474 20.2605 15.4772"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M16 3.13C16.7819 3.29965 17.4733 3.74032 17.9719 4.37917C18.4704 5.01802 18.7429 5.81924 18.7429 6.655C18.7429 7.49076 18.4704 8.29198 17.9719 8.93083C17.4733 9.56968 16.7819 10.0103 16 10.18"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [dashboardStats, setDashboardStats] = useState<DashboardStats | null>(null);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  // 加载数据
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);

        // 加载 Dashboard 完整统计数据
        const stats = await statsApi.getDashboard();
        setDashboardStats(stats);

        // 加载最近订单
        const ordersResult = await orderApi.getList({
          page: 1,
          pageSize: 5,
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

  if (!dashboardStats) {
    return (
      <div className={styles.container}>
        <h1 className={styles.title}>仪表盘</h1>
        <p>加载数据失败</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>仪表盘</h1>

      {/* 基础统计卡片 */}
      <div className={styles.statsGrid}>
        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <UserIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总用户数</div>
            <div className={styles.statValue}>{dashboardStats.totalUsers}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <PlayerIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总陪玩师</div>
            <div className={styles.statValue}>{dashboardStats.totalPlayers}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <StatsIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总游戏数</div>
            <div className={styles.statValue}>{dashboardStats.totalGames}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <OrderIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总订单数</div>
            <div className={styles.statValue}>{dashboardStats.totalOrders}</div>
          </div>
        </Card>

        <Card className={styles.statCard}>
          <div className={styles.statIcon}>
            <MoneyIcon />
          </div>
          <div className={styles.statContent}>
            <div className={styles.statLabel}>总收入</div>
            <div className={styles.statValue}>
              {formatCurrency(dashboardStats.totalPaidAmountCents)}
            </div>
          </div>
        </Card>
      </div>

      {/* 订单状态统计 */}
      <div className={styles.orderStatsSection}>
        <h2 className={styles.sectionTitle}>订单状态</h2>
        <div className={styles.orderStatsGrid}>
          <Card className={`${styles.orderStatCard} ${styles.pending}`}>
            <div
              className={styles.orderStatWrapper}
              onClick={() => navigate(`/orders?status=${OrderStatus.PENDING}`)}
            >
              <div className={styles.orderStatContent}>
                <div className={styles.orderStatLabel}>待处理</div>
                <div className={styles.orderStatValue}>
                  {dashboardStats.ordersByStatus?.pending || 0}
                </div>
              </div>
              <div className={styles.orderStatIcon}>
                <Badge count={dashboardStats.ordersByStatus?.pending || 0} />
              </div>
            </div>
          </Card>

          <Card className={`${styles.orderStatCard} ${styles.inProgress}`}>
            <div
              className={styles.orderStatWrapper}
              onClick={() => navigate(`/orders?status=${OrderStatus.IN_PROGRESS}`)}
            >
              <div className={styles.orderStatContent}>
                <div className={styles.orderStatLabel}>进行中</div>
                <div className={styles.orderStatValue}>
                  {dashboardStats.ordersByStatus?.in_progress || 0}
                </div>
              </div>
              <div className={styles.orderStatIcon}>
                <Badge count={dashboardStats.ordersByStatus?.in_progress || 0} />
              </div>
            </div>
          </Card>

          <Card className={`${styles.orderStatCard} ${styles.completed}`}>
            <div
              className={styles.orderStatWrapper}
              onClick={() => navigate(`/orders?status=${OrderStatus.COMPLETED}`)}
            >
              <div className={styles.orderStatContent}>
                <div className={styles.orderStatLabel}>已完成</div>
                <div className={styles.orderStatValue}>
                  {dashboardStats.ordersByStatus?.completed || 0}
                </div>
              </div>
              <div className={styles.orderStatIcon}>
                <Badge count={dashboardStats.ordersByStatus?.completed || 0} />
              </div>
            </div>
          </Card>

          <Card className={`${styles.orderStatCard} ${styles.canceled}`}>
            <div
              className={styles.orderStatWrapper}
              onClick={() => navigate(`/orders?status=${OrderStatus.CANCELED}`)}
            >
              <div className={styles.orderStatContent}>
                <div className={styles.orderStatLabel}>已取消</div>
                <div className={styles.orderStatValue}>
                  {dashboardStats.ordersByStatus?.canceled || 0}
                </div>
              </div>
              <div className={styles.orderStatIcon}>
                <Badge count={dashboardStats.ordersByStatus?.canceled || 0} />
              </div>
            </div>
          </Card>
        </div>
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
                {dashboardStats.ordersByStatus?.pending > 0 && (
                  <Badge count={dashboardStats.ordersByStatus.pending} />
                )}
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
                        {formatCurrency(order.priceCents)}
                      </span>
                    </div>
                  </div>

                  <div className={styles.orderFooter}>
                    <span className={styles.orderTime}>{formatRelativeTime(order.createdAt)}</span>
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
