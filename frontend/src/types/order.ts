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
  user_id: number;
  player_id?: number;
  game_id: number;
  title: string;
  description?: string;
  status: OrderStatus;
  price_cents: number;
  currency?: Currency;
  scheduled_start?: string;
  scheduled_end?: string;
  cancel_reason?: string;

  // 关联信息（API 返回时可能包含）
  user?: {
    id: number;
    name: string;
    avatar_url?: string;
  };
  player?: {
    id: number;
    nickname?: string;
    avatar_url?: string;
  };
  game?: {
    id: number;
    name: string;
    icon_url?: string;
  };
}

/**
 * 订单详情（包含扩展信息）
 */
export interface OrderDetail extends Order {
  // 时间节点
  paid_at?: string;
  accepted_at?: string;
  started_at?: string;
  completed_at?: string;
  cancelled_at?: string;

  // 操作日志
  logs?: Array<{
    id: number;
    action: string;
    operator_id: number;
    operator_name: string;
    note?: string;
    created_at: string;
  }>;

  // 审核记录
  reviews?: Array<{
    id: number;
    result: 'approved' | 'rejected';
    reason?: string;
    comment?: string;
    reviewer_id: number;
    reviewer_name: string;
    created_at: string;
  }>;
}

/**
 * 订单列表查询参数
 */
export interface OrderListQuery {
  page?: number;
  page_size?: number;
  user_id?: number;
  player_id?: number;
  game_id?: number;
  status?: OrderStatus;
  currency?: Currency;
  keyword?: string;
  date_from?: string;
  date_to?: string;
  sort_by?: 'created_at' | 'updated_at' | 'scheduled_start' | 'price_cents';
  sort_order?: 'asc' | 'desc';
}

/**
 * 创建订单请求
 */
export interface CreateOrderRequest {
  user_id: number;
  game_id: number;
  title: string;
  description?: string;
  price_cents: number;
  currency?: Currency;
  scheduled_start?: string;
  scheduled_end?: string;
}

/**
 * 更新订单请求
 */
export interface UpdateOrderRequest {
  title?: string;
  description?: string;
  price_cents?: number;
  currency?: Currency;
  scheduled_start?: string;
  scheduled_end?: string;
}

/**
 * 分配订单请求
 */
export interface AssignOrderRequest {
  player_id: number;
  note?: string;
}

/**
 * 审核订单请求
 */
export interface ReviewOrderRequest {
  result: 'approved' | 'rejected';
  reason?: string;
  comment?: string;
}

/**
 * 取消订单请求
 */
export interface CancelOrderRequest {
  cancel_reason: string;
}

/**
 * 订单统计数据
 */
export interface OrderStatistics {
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
