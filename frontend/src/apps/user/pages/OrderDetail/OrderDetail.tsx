/**
 * 订单详情页面 - Discord风格
 */

import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { DiscordLayout } from '@/shared/components/DiscordLayout';
import { OrderStatus } from '@/shared/types/order';
import styles from './OrderDetail.module.less';

/**
 * 订单详情接口
 */
interface OrderDetailInfo {
  id: number;
  orderNo: string;
  status: OrderStatus;
  gameName: string;
  serviceName: string;
  serviceDescription: string;
  duration: number;
  totalPrice: number;
  note?: string;
  createdAt: string;
  updatedAt: string;
  player?: {
    id: number;
    nickname: string;
    avatar: string;
    contact?: string;
  };
}

const statusConfig: Record<OrderStatus, { label: string; color: string; description: string }> = {
  pending: {
    label: '待接单',
    color: '#faa81a',
    description: '订单已创建，等待陪玩师接单',
  },
  confirmed: {
    label: '已接单',
    color: '#00aff4',
    description: '陪玩师已接单，请联系陪玩师开始服务',
  },
  in_progress: {
    label: '进行中',
    color: '#5865f2',
    description: '订单服务进行中',
  },
  completed: {
    label: '已完成',
    color: '#3ba55d',
    description: '订单已完成，期待您的评价',
  },
  canceled: {
    label: '已取消',
    color: '#72767d',
    description: '订单已取消',
  },
};

/**
 * 订单详情页面组件
 */
export const OrderDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [order, setOrder] = useState<OrderDetailInfo | null>(null);

  useEffect(() => {
    // Mock数据加载（后续替换为真实API调用）
    setLoading(true);
    setTimeout(() => {
      setOrder({
        id: Number(id),
        orderNo: `ORD20251116${String(id).padStart(3, '0')}`,
        status: 'in_progress',
        gameName: '王者荣耀',
        serviceName: '上分陪玩',
        serviceDescription: '带你快速上分，保证胜率',
        duration: 2,
        totalPrice: 100,
        note: '希望能打辅助位，我玩射手',
        createdAt: '2025-11-16 14:30:00',
        updatedAt: '2025-11-16 14:35:00',
        player: {
          id: 1,
          nickname: '王者大神',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player1',
          contact: 'WeChat: player001',
        },
      });
      setLoading(false);
    }, 500);
  }, [id]);

  const handleAction = (action: string) => {
    switch (action) {
      case 'cancel':
        if (confirm('确认要取消订单吗？')) {
          alert('取消订单功能待实现');
        }
        break;
      case 'contact':
        alert(`联系陪玩师：${order?.player?.contact || '暂无联系方式'}`);
        break;
      case 'complete':
        if (confirm('确认服务已完成吗？')) {
          alert('完成订单功能待实现');
        }
        break;
      case 'review':
        navigate(`/user/orders/${id}/review`);
        break;
      default:
        break;
    }
  };

  if (loading) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.loading}>
          <div className={styles.spinner} />
          <p>加载中...</p>
        </div>
      </DiscordLayout>
    );
  }

  if (!order) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.error}>
          <p>❌ 订单不存在</p>
          <button onClick={() => navigate('/user/orders')}>返回订单列表</button>
        </div>
      </DiscordLayout>
    );
  }

  const statusInfo = statusConfig[order.status];

  return (
    <DiscordLayout showMemberPanel={false}>
      <div className={styles.orderDetail}>
        {/* 页面标题 */}
        <div className={styles.header}>
          <button className={styles.backButton} onClick={() => navigate('/user/orders')}>
            ← 返回订单列表
          </button>
          <h1 className={styles.title}>订单详情</h1>
        </div>

        {/* 订单状态 */}
        <div className={styles.statusSection}>
          <div className={styles.statusBadge} style={{ backgroundColor: statusInfo.color }}>
            {statusInfo.label}
          </div>
          <div className={styles.statusDesc}>{statusInfo.description}</div>
        </div>

        {/* 基本信息 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>基本信息</h2>
          <div className={styles.infoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>订单编号</span>
              <span className={styles.infoValue}>{order.orderNo}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>创建时间</span>
              <span className={styles.infoValue}>{order.createdAt}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>更新时间</span>
              <span className={styles.infoValue}>{order.updatedAt}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>游戏</span>
              <span className={styles.infoValue}>{order.gameName}</span>
            </div>
          </div>
        </div>

        {/* 服务信息 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>服务信息</h2>
          <div className={styles.serviceInfo}>
            <div className={styles.serviceName}>{order.serviceName}</div>
            <div className={styles.serviceDesc}>{order.serviceDescription}</div>
            <div className={styles.serviceDetails}>
              <div className={styles.serviceItem}>
                <span className={styles.serviceLabel}>时长</span>
                <span className={styles.serviceValue}>{order.duration}小时</span>
              </div>
              <div className={styles.serviceItem}>
                <span className={styles.serviceLabel}>总价</span>
                <span className={styles.servicePriceValue}>¥{order.totalPrice}</span>
              </div>
            </div>
          </div>
        </div>

        {/* 陪玩师信息 */}
        {order.player && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>陪玩师信息</h2>
            <div className={styles.playerInfo}>
              <img
                src={order.player.avatar}
                alt={order.player.nickname}
                className={styles.playerAvatar}
              />
              <div className={styles.playerDetails}>
                <div className={styles.playerName}>{order.player.nickname}</div>
                {order.player.contact && (
                  <div className={styles.playerContact}>{order.player.contact}</div>
                )}
              </div>
              <button
                className={styles.contactButton}
                onClick={() => handleAction('contact')}
              >
                联系陪玩师
              </button>
            </div>
          </div>
        )}

        {/* 备注 */}
        {order.note && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>备注</h2>
            <div className={styles.noteContent}>{order.note}</div>
          </div>
        )}

        {/* 操作按钮 */}
        <div className={styles.actions}>
          {order.status === 'pending' && (
            <button
              className={`${styles.actionButton} ${styles.dangerButton}`}
              onClick={() => handleAction('cancel')}
            >
              取消订单
            </button>
          )}
          {order.status === 'confirmed' && (
            <>
              <button
                className={`${styles.actionButton} ${styles.primaryButton}`}
                onClick={() => handleAction('contact')}
              >
                联系陪玩师
              </button>
              <button
                className={`${styles.actionButton} ${styles.secondaryButton}`}
                onClick={() => handleAction('cancel')}
              >
                取消订单
              </button>
            </>
          )}
          {order.status === 'in_progress' && (
            <>
              <button
                className={`${styles.actionButton} ${styles.primaryButton}`}
                onClick={() => handleAction('contact')}
              >
                联系陪玩师
              </button>
              <button
                className={`${styles.actionButton} ${styles.successButton}`}
                onClick={() => handleAction('complete')}
              >
                确认完成
              </button>
            </>
          )}
          {order.status === 'completed' && (
            <button
              className={`${styles.actionButton} ${styles.primaryButton}`}
              onClick={() => handleAction('review')}
            >
              评价订单
            </button>
          )}
        </div>
      </div>
    </DiscordLayout>
  );
};

export default OrderDetail;
