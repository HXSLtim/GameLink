/**
 * 我的订单页面 - Discord风格
 * 使用新的OrderCard组件
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DiscordLayout } from '@/shared/components/DiscordLayout';
import { OrderCard } from '@/apps/user/components/PlayerCard';
import type { OrderCardProps } from '@/apps/user/components/PlayerCard';
import { OrderStatus } from '@/shared/types/order';
import styles from './MyOrders.module.less';

/**
 * 我的订单页面组件
 */
export const MyOrders: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<OrderStatus | 'all'>('all');

  // Mock数据（后续替换为真实API调用）
  const [orders] = useState<OrderCardProps[]>([
    {
      id: 1,
      orderNo: 'ORD20251116001',
      gameName: '王者荣耀',
      serviceName: '上分陪玩 - 钻石段位',
      status: 'in_progress',
      player: {
        id: 1,
        nickname: '王者大神',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player1',
      },
      createdAt: '2025-11-16 14:30:00',
      duration: 2,
      totalPrice: 100,
    },
    {
      id: 2,
      orderNo: 'ORD20251115023',
      gameName: '英雄联盟',
      serviceName: '技术教学 - 中单位置',
      status: 'confirmed',
      player: {
        id: 3,
        nickname: 'LOL王者',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player3',
      },
      createdAt: '2025-11-15 18:00:00',
      duration: 2,
      totalPrice: 120,
    },
    {
      id: 3,
      orderNo: 'ORD20251115012',
      gameName: '王者荣耀',
      serviceName: '娱乐陪玩',
      status: 'completed',
      player: {
        id: 2,
        nickname: '温柔小姐姐',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player2',
      },
      createdAt: '2025-11-15 10:00:00',
      duration: 1,
      totalPrice: 80,
    },
    {
      id: 4,
      orderNo: 'ORD20251114045',
      gameName: '和平精英',
      serviceName: '吃鸡陪玩',
      status: 'pending',
      createdAt: '2025-11-14 20:30:00',
      duration: 1,
      totalPrice: 45,
    },
    {
      id: 5,
      orderNo: 'ORD20251113028',
      gameName: '原神',
      serviceName: '深渊挑战',
      status: 'canceled',
      createdAt: '2025-11-13 16:00:00',
      duration: 1,
      totalPrice: 40,
    },
    {
      id: 6,
      orderNo: 'ORD20251112067',
      gameName: 'CS:GO',
      serviceName: '竞技陪玩',
      status: 'completed',
      player: {
        id: 6,
        nickname: 'CS高手',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player6',
      },
      createdAt: '2025-11-12 14:00:00',
      duration: 1,
      totalPrice: 70,
    },
  ]);

  const handleOrderClick = (id: number) => {
    navigate(`/user/orders/${id}`);
  };

  const handleOrderAction = (action: string, id: number) => {
    console.log(`Action: ${action}, Order ID: ${id}`);
    // 这里处理各种操作
    switch (action) {
      case 'cancel':
        // 取消订单
        alert('取消订单功能待实现');
        break;
      case 'contact':
        // 联系陪玩师
        navigate(`/user/chat/${id}`);
        break;
      case 'review':
        // 评价订单
        navigate(`/user/orders/${id}/review`);
        break;
      default:
        break;
    }
  };

  // 筛选订单
  const filteredOrders = orders.filter((order) => {
    if (activeTab === 'all') return true;
    return order.status === activeTab;
  });

  const tabs: Array<{ key: OrderStatus | 'all'; label: string }> = [
    { key: 'all', label: '全部' },
    { key: 'pending', label: '待接单' },
    { key: 'confirmed', label: '已接单' },
    { key: 'in_progress', label: '进行中' },
    { key: 'completed', label: '已完成' },
    { key: 'canceled', label: '已取消' },
  ];

  return (
    <DiscordLayout
      serverList={
        <div className={styles.serverList}>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>📋</span>
          </div>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>✅</span>
          </div>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>❌</span>
          </div>
        </div>
      }
      showMemberPanel={false}
    >
      <div className={styles.myOrders}>
        {/* 页面标题 */}
        <div className={styles.header}>
          <h1 className={styles.title}>我的订单</h1>
        </div>

        {/* 订单状态标签 */}
        <div className={styles.tabs}>
          {tabs.map((tab) => (
            <button
              key={tab.key}
              className={`${styles.tab} ${activeTab === tab.key ? styles.active : ''}`}
              onClick={() => setActiveTab(tab.key)}
            >
              {tab.label}
              <span className={styles.tabCount}>
                {tab.key === 'all'
                  ? orders.length
                  : orders.filter((o) => o.status === tab.key).length}
              </span>
            </button>
          ))}
        </div>

        {/* 订单列表 */}
        <div className={styles.content}>
          {loading ? (
            <div className={styles.loading}>
              <div className={styles.spinner} />
              <p>加载中...</p>
            </div>
          ) : filteredOrders.length > 0 ? (
            <div className={styles.orderList}>
              {filteredOrders.map((order) => (
                <OrderCard
                  key={order.id}
                  {...order}
                  onClick={handleOrderClick}
                  onAction={handleOrderAction}
                />
              ))}
            </div>
          ) : (
            <div className={styles.empty}>
              <p>📭 暂无{activeTab !== 'all' ? tabs.find((t) => t.key === activeTab)?.label : ''}订单</p>
              <button className={styles.emptyButton} onClick={() => navigate('/user/players')}>
                去找陪玩师
              </button>
            </div>
          )}
        </div>
      </div>
    </DiscordLayout>
  );
};

export default MyOrders;
