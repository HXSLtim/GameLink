/**
 * 订单卡片组件
 * 用于展示订单的基本信息、状态、价格等
 */

import { FC } from 'react';
import styles from './OrderCard.module.less';

export type OrderStatus = 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'canceled';

export interface OrderCardProps {
  /** 订单ID */
  id: number;

  /** 订单编号 */
  orderNo: string;

  /** 游戏名称 */
  gameName: string;

  /** 服务项目名称 */
  serviceName: string;

  /** 订单状态 */
  status: OrderStatus;

  /** 陪玩师信息（已接单时） */
  player?: {
    id: number;
    nickname: string;
    avatar: string;
  };

  /** 下单时间 */
  createdAt: string;

  /** 服务时长（小时） */
  duration: number;

  /** 总价格 */
  totalPrice: number;

  /** 点击事件 */
  onClick?: (id: number) => void;

  /** 操作按钮点击事件 */
  onAction?: (action: string, id: number) => void;

  /** 自定义类名 */
  className?: string;
}

const statusConfig: Record<OrderStatus, { label: string; className: string }> = {
  pending: { label: '待接单', className: 'statusPending' },
  confirmed: { label: '已接单', className: 'statusConfirmed' },
  in_progress: { label: '进行中', className: 'statusInProgress' },
  completed: { label: '已完成', className: 'statusCompleted' },
  canceled: { label: '已取消', className: 'statusCanceled' },
};

/**
 * 订单卡片组件
 */
export const OrderCard: FC<OrderCardProps> = ({
  id,
  orderNo,
  gameName,
  serviceName,
  status,
  player,
  createdAt,
  duration,
  totalPrice,
  onClick,
  onAction,
  className,
}) => {
  const handleClick = () => {
    onClick?.(id);
  };

  const handleAction = (action: string) => {
    onAction?.(action, id);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getActionButtons = () => {
    switch (status) {
      case 'pending':
        return (
          <button
            className={`${styles.actionButton} ${styles.buttonDanger}`}
            onClick={(e) => {
              e.stopPropagation();
              handleAction('cancel');
            }}
          >
            取消订单
          </button>
        );
      case 'confirmed':
        return (
          <>
            <button
              className={`${styles.actionButton} ${styles.buttonPrimary}`}
              onClick={(e) => {
                e.stopPropagation();
                handleAction('contact');
              }}
            >
              联系陪玩师
            </button>
            <button
              className={`${styles.actionButton} ${styles.buttonSecondary}`}
              onClick={(e) => {
                e.stopPropagation();
                handleAction('cancel');
              }}
            >
              取消
            </button>
          </>
        );
      case 'in_progress':
        return (
          <button
            className={`${styles.actionButton} ${styles.buttonPrimary}`}
            onClick={(e) => {
              e.stopPropagation();
              handleAction('contact');
            }}
          >
            联系陪玩师
          </button>
        );
      case 'completed':
        return (
          <button
            className={`${styles.actionButton} ${styles.buttonPrimary}`}
            onClick={(e) => {
              e.stopPropagation();
              handleAction('review');
            }}
          >
            评价
          </button>
        );
      default:
        return null;
    }
  };

  const statusInfo = statusConfig[status];

  return (
    <div
      className={`${styles.orderCard} ${className || ''}`}
      onClick={handleClick}
      role="button"
      tabIndex={0}
      onKeyPress={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          handleClick();
        }
      }}
    >
      {/* 头部：订单号和状态 */}
      <div className={styles.header}>
        <div className={styles.orderInfo}>
          <span className={styles.orderNo}>#{orderNo}</span>
          <span className={styles.createdAt}>{formatDate(createdAt)}</span>
        </div>
        <span className={`${styles.status} ${styles[statusInfo.className]}`}>
          {statusInfo.label}
        </span>
      </div>

      {/* 主体内容 */}
      <div className={styles.content}>
        {/* 游戏和服务信息 */}
        <div className={styles.serviceInfo}>
          <h3 className={styles.gameName}>{gameName}</h3>
          <p className={styles.serviceName}>{serviceName}</p>
        </div>

        {/* 陪玩师信息（已接单时显示） */}
        {player && (
          <div className={styles.playerInfo}>
            <img
              src={player.avatar}
              alt={player.nickname}
              className={styles.playerAvatar}
            />
            <span className={styles.playerNickname}>{player.nickname}</span>
          </div>
        )}

        {/* 订单详情 */}
        <div className={styles.details}>
          <div className={styles.detailItem}>
            <span className={styles.detailLabel}>时长</span>
            <span className={styles.detailValue}>{duration}小时</span>
          </div>
          <div className={styles.detailItem}>
            <span className={styles.detailLabel}>总价</span>
            <span className={styles.priceValue}>￥{totalPrice}</span>
          </div>
        </div>
      </div>

      {/* 操作按钮 */}
      {onAction && (
        <div className={styles.actions}>
          {getActionButtons()}
        </div>
      )}
    </div>
  );
};

export default OrderCard;
