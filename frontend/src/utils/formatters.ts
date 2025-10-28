import { OrderStatus } from '../types/order';

/**
 * Tag 颜色类型
 */
export type TagColor =
  | 'success'
  | 'warning'
  | 'error'
  | 'info'
  | 'processing'
  | 'pending'
  | 'default';

/**
 * 格式化订单状态文本 - 与后端一致
 */
export const formatOrderStatus = (status: OrderStatus): string => {
  const statusMap: Record<OrderStatus, string> = {
    [OrderStatus.PENDING]: '待处理',
    [OrderStatus.CONFIRMED]: '已确认',
    [OrderStatus.IN_PROGRESS]: '进行中',
    [OrderStatus.COMPLETED]: '已完成',
    [OrderStatus.CANCELED]: '已取消',
    [OrderStatus.REFUNDED]: '已退款',
  };
  return statusMap[status] || '未知';
};

/**
 * 获取订单状态对应的标签颜色
 */
export const getOrderStatusColor = (status: OrderStatus): TagColor => {
  const colorMap: Record<OrderStatus, TagColor> = {
    [OrderStatus.PENDING]: 'warning',
    [OrderStatus.CONFIRMED]: 'info',
    [OrderStatus.IN_PROGRESS]: 'processing',
    [OrderStatus.COMPLETED]: 'success',
    [OrderStatus.CANCELED]: 'default',
    [OrderStatus.REFUNDED]: 'error',
  };
  return colorMap[status] || 'default';
};

/**
 * 格式化金额（分转元）
 */
export const formatCurrency = (cents: number): string => {
  const yuan = cents / 100;
  return `¥${yuan.toFixed(2)}`;
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
  
  const date = new Date(dateStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}`;
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

/**
 * 格式化价格（分转元）
 */
export const formatPrice = (cents: number): string => {
  return formatCurrency(cents);
};
