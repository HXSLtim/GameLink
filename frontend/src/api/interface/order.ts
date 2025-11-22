/**
 * 订单相关接口定义
 */

import type { BaseEntity } from '@/shared/types/api';

/**
 * 订单状态
 */
export enum OrderStatus {
  PENDING = 'pending', // 待处理
  CONFIRMED = 'confirmed', // 已确认
  IN_PROGRESS = 'in_progress', // 进行中
  COMPLETED = 'completed', // 已完成
  CANCELED = 'canceled', // 已取消
  REFUNDED = 'refunded', // 已退款
}

/**
 * 订单实体
 */
export interface Order extends BaseEntity {
  id: number;
  orderNo: string;
  userId: number;
  playerId: number;
  gameId: number;
  gameName: string;
  serviceType: string;
  duration: number;
  amount: number;
  status: OrderStatus;
  notes?: string;
  requirements?: string;
  startTime?: string;
  endTime?: string;
  paymentStatus: 'pending' | 'paid' | 'refunded';
  paymentMethod?: string;
  user: {
    id: number;
    username: string;
    avatar?: string;
  };
  player: {
    id: number;
    username: string;
    avatar?: string;
  };
}

/**
 * 创建订单请求
 */
export interface CreateOrderRequest {
  playerId: number;
  gameId: number;
  serviceType: string;
  duration: number;
  amount: number;
  notes?: string;
  requirements?: string;
}

/**
 * 更新订单请求
 */
export interface UpdateOrderRequest {
  status?: OrderStatus;
  notes?: string;
  requirements?: string;
  startTime?: string;
  endTime?: string;
}

/**
 * 订单列表查询参数
 */
export interface GetOrdersParams {
  status?: OrderStatus;
  userId?: number;
  playerId?: number;
  page?: number;
  pageSize?: number;
  startDate?: string;
  endDate?: string;
}

/**
 * 订单列表响应
 */
export interface GetOrdersResponse {
  list: Order[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

/**
 * 订单统计
 */
export interface OrderStats {
  total: number;
  pending: number;
  inProgress: number;
  completed: number;
  canceled: number;
  refunded: number;
  totalAmount: number;
}
