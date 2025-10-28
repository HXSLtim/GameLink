import {
  OrderStatus,
  ReviewStatus,
  GameType,
  ServiceType,
  type TagColor,
} from '../types/order.types';

/**
 * 格式化订单状态文本
 */
export const formatOrderStatus = (status: OrderStatus): string => {
  const statusMap: Record<OrderStatus, string> = {
    [OrderStatus.PENDING_PAYMENT]: '待支付',
    [OrderStatus.PENDING_ACCEPT]: '待接单',
    [OrderStatus.PENDING_REVIEW]: '待审核',
    [OrderStatus.IN_REVIEW]: '审核中',
    [OrderStatus.REVIEW_APPROVED]: '审核通过',
    [OrderStatus.REVIEW_REJECTED]: '审核拒绝',
    [OrderStatus.IN_PROGRESS]: '进行中',
    [OrderStatus.COMPLETED]: '已完成',
    [OrderStatus.CANCELLED]: '已取消',
    [OrderStatus.REFUNDED]: '已退款',
  };
  return statusMap[status] || '未知';
};

/**
 * 获取订单状态对应的标签颜色
 */
export const getOrderStatusColor = (status: OrderStatus): TagColor => {
  const colorMap: Record<OrderStatus, TagColor> = {
    [OrderStatus.PENDING_PAYMENT]: 'warning',
    [OrderStatus.PENDING_ACCEPT]: 'info',
    [OrderStatus.PENDING_REVIEW]: 'pending',
    [OrderStatus.IN_REVIEW]: 'processing',
    [OrderStatus.REVIEW_APPROVED]: 'success',
    [OrderStatus.REVIEW_REJECTED]: 'error',
    [OrderStatus.IN_PROGRESS]: 'processing',
    [OrderStatus.COMPLETED]: 'success',
    [OrderStatus.CANCELLED]: 'default',
    [OrderStatus.REFUNDED]: 'error',
  };
  return colorMap[status] || 'default';
};

/**
 * 格式化审核状态文本
 */
export const formatReviewStatus = (status: ReviewStatus): string => {
  const statusMap: Record<ReviewStatus, string> = {
    [ReviewStatus.PENDING]: '待审核',
    [ReviewStatus.IN_REVIEW]: '审核中',
    [ReviewStatus.APPROVED]: '已通过',
    [ReviewStatus.REJECTED]: '已拒绝',
  };
  return statusMap[status] || '未知';
};

/**
 * 获取审核状态对应的标签颜色
 */
export const getReviewStatusColor = (status: ReviewStatus): TagColor => {
  const colorMap: Record<ReviewStatus, TagColor> = {
    [ReviewStatus.PENDING]: 'pending',
    [ReviewStatus.IN_REVIEW]: 'processing',
    [ReviewStatus.APPROVED]: 'success',
    [ReviewStatus.REJECTED]: 'error',
  };
  return colorMap[status] || 'default';
};

/**
 * 格式化游戏类型文本
 */
export const formatGameType = (type: GameType): string => {
  const typeMap: Record<GameType, string> = {
    [GameType.HONOR_OF_KINGS]: '王者荣耀',
    [GameType.LEAGUE_OF_LEGENDS]: '英雄联盟',
    [GameType.PEACEKEEPER_ELITE]: '和平精英',
    [GameType.GENSHIN_IMPACT]: '原神',
    [GameType.OTHER]: '其他',
  };
  return typeMap[type] || '未知';
};

/**
 * 格式化服务类型文本
 */
export const formatServiceType = (type: ServiceType): string => {
  const typeMap: Record<ServiceType, string> = {
    [ServiceType.ACCOMPANY]: '陪玩',
    [ServiceType.BOOST]: '代练',
    [ServiceType.RANK_UP]: '上分',
    [ServiceType.ENTERTAINMENT]: '娱乐',
  };
  return typeMap[type] || '未知';
};

/**
 * 格式化金额
 */
export const formatCurrency = (amount: number): string => {
  return `¥${amount.toFixed(2)}`;
};

/**
 * 格式化时长（小时）
 */
export const formatDuration = (hours: number): string => {
  return `${hours}小时`;
};

/**
 * 格式化日期时间
 */
export const formatDateTime = (dateStr?: string): string => {
  if (!dateStr) return '-';
  return dateStr;
};

/**
 * 格式化相对时间
 */
export const formatRelativeTime = (dateStr: string): string => {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (diff < minute) {
    return '刚刚';
  } else if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`;
  } else if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`;
  } else if (diff < 7 * day) {
    return `${Math.floor(diff / day)}天前`;
  } else {
    return formatDateTime(dateStr);
  }
};
