import type { OrderStatus } from '../types/order';
import type { PaymentStatus, PaymentMethod } from '../types/payment';
import type { UserRole, UserStatus } from '../types/user';

/**
 * 通用状态格式化工具函数
 */

// ==================== 订单状态 ====================
export const ORDER_STATUS_MAP: Record<OrderStatus, string> = {
  pending: '待处理',
  confirmed: '已确认',
  in_progress: '进行中',
  completed: '已完成',
  canceled: '已取消',
  refunded: '已退款',
};

export const ORDER_STATUS_COLORS: Record<OrderStatus, string> = {
  pending: 'orange',
  confirmed: 'blue',
  in_progress: 'cyan',
  completed: 'green',
  canceled: 'default',
  refunded: 'purple',
};

export const formatOrderStatus = (status: OrderStatus): string => {
  return ORDER_STATUS_MAP[status] || status;
};

export const getOrderStatusColor = (status: OrderStatus): string => {
  return ORDER_STATUS_COLORS[status] || 'default';
};

// ==================== 支付状态 ====================
export const PAYMENT_STATUS_MAP: Record<PaymentStatus, string> = {
  pending: '待支付',
  paid: '已支付',
  failed: '支付失败',
  refunded: '已退款',
  cancelled: '已取消',
};

export const PAYMENT_STATUS_COLORS: Record<PaymentStatus, string> = {
  pending: 'orange',
  paid: 'green',
  failed: 'red',
  refunded: 'purple',
  cancelled: 'default',
};

export const formatPaymentStatus = (status: PaymentStatus): string => {
  return PAYMENT_STATUS_MAP[status] || status;
};

export const getPaymentStatusColor = (status: PaymentStatus): string => {
  return PAYMENT_STATUS_COLORS[status] || 'default';
};

// ==================== 支付方式 ====================
export const PAYMENT_METHOD_MAP: Record<PaymentMethod, string> = {
  alipay: '支付宝',
  wechat: '微信支付',
  balance: '余额支付',
};

export const formatPaymentMethod = (method: PaymentMethod): string => {
  return PAYMENT_METHOD_MAP[method] || method;
};

// ==================== 用户角色 ====================
export const USER_ROLE_MAP: Record<UserRole, string> = {
  user: '普通用户',
  player: '陪玩师',
  admin: '管理员',
};

export const USER_ROLE_COLORS: Record<UserRole, string> = {
  user: 'blue',
  player: 'green',
  admin: 'red',
};

export const formatUserRole = (role: UserRole): string => {
  return USER_ROLE_MAP[role] || role;
};

export const getUserRoleColor = (role: UserRole): string => {
  return USER_ROLE_COLORS[role] || 'default';
};

// ==================== 用户状态 ====================
export const USER_STATUS_MAP: Record<UserStatus, string> = {
  active: '正常',
  suspended: '暂停',
  banned: '封禁',
};

export const USER_STATUS_COLORS: Record<UserStatus, string> = {
  active: 'green',
  suspended: 'orange',
  banned: 'red',
};

export const formatUserStatus = (status: UserStatus): string => {
  return USER_STATUS_MAP[status] || status;
};

export const getUserStatusColor = (status: UserStatus): string => {
  return USER_STATUS_COLORS[status] || 'default';
};

// ==================== 评分 ====================
export const RATING_MAP: Record<number, string> = {
  1: '非常差',
  2: '较差',
  3: '一般',
  4: '满意',
  5: '非常满意',
};

export const RATING_COLORS: Record<number, string> = {
  1: 'red',
  2: 'orange',
  3: 'blue',
  4: 'green',
  5: 'green',
};

export const formatRating = (rating: number): string => {
  return RATING_MAP[rating] || `${rating} 星`;
};

export const getRatingColor = (rating: number): string => {
  if (rating >= 4) return 'green';
  if (rating >= 3) return 'blue';
  if (rating >= 2) return 'orange';
  return 'red';
};

// ==================== 游戏分类 ====================
export const GAME_CATEGORY_MAP: Record<string, string> = {
  moba: 'MOBA',
  fps: 'FPS',
  rpg: 'RPG',
  strategy: '策略',
  sports: '体育',
  racing: '竞速',
  puzzle: '益智',
  other: '其他',
};

export const GAME_CATEGORY_COLORS: Record<string, string> = {
  moba: 'blue',
  fps: 'red',
  rpg: 'purple',
  strategy: 'orange',
  sports: 'green',
  racing: 'cyan',
  puzzle: 'yellow',
  other: 'default',
};

export const formatGameCategory = (category: string): string => {
  return GAME_CATEGORY_MAP[category] || category;
};

export const getGameCategoryColor = (category: string): string => {
  return GAME_CATEGORY_COLORS[category] || 'default';
};
