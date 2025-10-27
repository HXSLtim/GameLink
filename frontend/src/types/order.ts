import { BaseEntity } from './user';

/**
 * 货币类型枚举
 */
export type Currency = 'CNY' | 'USD' | 'EUR';

/**
 * 订单状态枚举
 */
export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'refunded';

/**
 * 订单实体 - 与后端model.Order保持一致
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
  status?: OrderStatus;
  price_cents?: number;
  currency?: Currency;
  scheduled_start?: string;
  scheduled_end?: string;
  cancel_reason?: string;
}

/**
 * 订单状态映射到徽章类型
 */
export const ORDER_STATUS_BADGE: Record<
  OrderStatus,
  'success' | 'warning' | 'error' | 'processing'
> = {
  pending: 'processing',
  confirmed: 'processing',
  in_progress: 'processing',
  completed: 'success',
  cancelled: 'warning',
  refunded: 'warning',
};

/**
 * 订单状态显示文本
 */
export const ORDER_STATUS_TEXT: Record<OrderStatus, string> = {
  pending: '待处理',
  confirmed: '已确认',
  in_progress: '进行中',
  completed: '已完成',
  cancelled: '已取消',
  refunded: '已退款',
};
