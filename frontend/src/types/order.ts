import type { BaseEntity } from './user';

/**
 * 货币类型枚举
 */
export enum Currency {
  CNY = 'CNY',
  USD = 'USD',
}

/**
 * 订单状态枚举 - 与后端 model.Order 保持一致
 */
export enum OrderStatus {
  PENDING = 'pending', // 待处理
  CONFIRMED = 'confirmed', // 已确认
  IN_PROGRESS = 'in_progress', // 进行中
  COMPLETED = 'completed', // 已完成
  CANCELED = 'canceled', // 已取消（拼写与后端一致）
  REFUNDED = 'refunded', // 已退款
}

/**
 * 订单实体 - 与后端 model.Order 保持一致
 */
export interface Order extends BaseEntity {
  userId: number;
  playerId?: number;
  gameId: number;
  title: string;
  description?: string;
  status: OrderStatus;
  priceCents: number;
  currency?: Currency;
  scheduledStart?: string;
  scheduledEnd?: string;
  cancelReason?: string;
  assignmentSource?: string;
  disputeStatus?: string;

  // 关联信息（API 返回时可能包含）
  user?: {
    id: number;
    name: string;
    avatarUrl?: string;
  };
  player?: {
    id: number;
    nickname?: string;
    avatarUrl?: string;
  };
  game?: {
    id: number;
    name: string;
    iconUrl?: string;
  };
}

/**
 * 订单详情（包含扩展信息）
 */
export interface OrderDetail extends Order {
  // 时间节点
  paidAt?: string;
  acceptedAt?: string;
  startedAt?: string;
  completedAt?: string;
  cancelledAt?: string;

  // 操作日志
  logs?: Array<{
    id: number;
    action: string;
    operatorId: number;
    operatorName: string;
    note?: string;
    createdAt: string;
  }>;

  // 审核记录
  reviews?: Array<{
    id: number;
    approved: boolean; // true=通过，false=拒绝
    reason?: string;
    comment?: string;
    reviewerId: number;
    reviewerName: string;
    createdAt: string;
  }>;
}

/**
 * 订单列表查询参数
 */
export interface OrderListQuery {
  page?: number;
  pageSize?: number;
  userId?: number;
  playerId?: number;
  gameId?: number;
  status?: OrderStatus;
  currency?: Currency;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'scheduledStart' | 'priceCents';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 创建订单请求 - 与后端 Swagger 同步
 */
export interface CreateOrderRequest {
  userId: number;
  gameId: number;
  playerId?: number; // 可选：在创建时指定陪玩师
  title?: string;
  description?: string;
  priceCents: number;
  currency: string; // 必填
  scheduledStart?: string;
  scheduledEnd?: string;
}

/**
 * 更新订单请求 - 与后端 Swagger 同步
 */
export interface UpdateOrderRequest {
  currency: string; // 必填
  priceCents: number; // 必填
  status: string; // 必填：订单状态
  scheduledStart?: string;
  scheduledEnd?: string;
  cancelReason?: string; // 取消原因
}

/**
 * 分配订单请求
 */
export interface AssignOrderRequest {
  playerId: number;
  note?: string;
  source?: string;
}

/**
 * 审核订单请求 - 与后端 Swagger 同步
 */
export interface ReviewOrderRequest {
  approved: boolean; // true=通过，false=拒绝
  reason?: string; // 拒绝原因或备注
}

/**
 * 取消订单请求 - 与后端 Swagger 同步
 */
export interface CancelOrderRequest {
  reason?: string; // 取消原因
}

/**
 * 订单统计数据
 */
export interface OrderStatsData {
  total: number;
  pending: number;
  confirmed: number;
  in_progress: number;
  completed: number;
  canceled: number;
  refunded: number;
  today_orders: number;
  today_revenue: number; // 分
}

/**
 * 订单状态显示文本
 */
export const ORDER_STATUS_TEXT: Record<OrderStatus, string> = {
  [OrderStatus.PENDING]: '待处理',
  [OrderStatus.CONFIRMED]: '已确认',
  [OrderStatus.IN_PROGRESS]: '进行中',
  [OrderStatus.COMPLETED]: '已完成',
  [OrderStatus.CANCELED]: '已取消',
  [OrderStatus.REFUNDED]: '已退款',
};

/**
 * 订单状态映射到徽章类型
 */
export const ORDER_STATUS_BADGE: Record<
  OrderStatus,
  'success' | 'warning' | 'error' | 'processing' | 'default'
> = {
  [OrderStatus.PENDING]: 'processing',
  [OrderStatus.CONFIRMED]: 'processing',
  [OrderStatus.IN_PROGRESS]: 'processing',
  [OrderStatus.COMPLETED]: 'success',
  [OrderStatus.CANCELED]: 'default',
  [OrderStatus.REFUNDED]: 'warning',
};

/**
 * 货币符号
 */
export const CURRENCY_SYMBOL: Record<Currency, string> = {
  [Currency.CNY]: '¥',
  [Currency.USD]: '$',
};
