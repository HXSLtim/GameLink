/**
 * 订单相关 API
 */

import { get, post, put, type RequestConfig } from './request'
import type { OrderPaymentStatus } from '@/types/order'

// 订单状态
export type OrderStatus =
  | 'pending_payment'   // 待支付
  | 'pending_accept'    // 待接单
  | 'accepted'          // 已接单
  | 'in_progress'       // 进行中
  | 'completed'         // 已完成
  | 'cancelled'         // 已取消
  | 'refunding'         // 退款中
  | 'refunded'          // 已退款

// 订单项
export interface Order {
  id: number
  orderNo: string
  userId: number
  playerId: number
  playerName: string
  playerAvatar?: string
  gameId: number
  gameName: string
  gameIcon?: string
  serviceType: string
  quantity: number
  unit: string
  unitPriceCents: number
  totalPriceCents: number
  platformFeeCents: number
  status: OrderStatus
  statusText: string
  startTime?: string
  endTime?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

// 订单详情
export interface OrderDetail extends Order {
  player: {
    id: number
    nickname: string
    avatar?: string
    rating: number
    gameRank?: string
  }
  user: {
    id: number
    nickname: string
    avatar?: string
  }
  chatGroupId?: number
  review?: {
    id: number
    rating: number
    content?: string
    createdAt: string
  }
}

// 创建订单参数
export interface CreateOrderParams {
  playerId: number
  serviceId: number
  quantity: number
  couponId?: number
  remark?: string
}

// 列表查询参数
export interface OrderListParams {
  page?: number
  page_size?: number
  status?: OrderStatus | 'all'
}

/**
 * 获取订单列表
 */
export function getOrders(params?: OrderListParams, config?: Partial<RequestConfig>) {
  return get<Order[]>('/orders', params, config)
}

/**
 * 获取订单详情
 */
export function getOrderDetail(orderId: number, config?: Partial<RequestConfig>) {
  return get<OrderDetail>(`/orders/${orderId}`, undefined, config)
}

/**
 * 创建订单
 */
export function createOrder(data: CreateOrderParams) {
  return post<Order>('/orders', data)
}

/**
 * 取消订单
 */
export function cancelOrder(orderId: number, reason?: string) {
  return put<void>(`/orders/${orderId}/cancel`, { reason })
}

/**
 * 确认完成订单
 */
export function completeOrder(orderId: number) {
  return put<void>(`/orders/${orderId}/complete`)
}

/**
 * 申请退款
 */
export function refundOrder(orderId: number, reason: string) {
  return post<void>(`/orders/${orderId}/refund`, { reason })
}

/**
 * 提交评价
 */
export function submitReview(orderId: number, data: {
  rating: number
  content?: string
  tags?: string[]
  images?: string[]
}) {
  return post<void>('/users/reviews', { orderId, ...data })
}

/**
 * 获取订单支付状态
 */
export function getOrderPaymentStatus(orderId: number) {
  return get<{ status: OrderPaymentStatus }>(`/orders/${orderId}/payment-status`)
}

// ============ 陪玩师订单 API ============

/**
 * 获取陪玩师订单列表
 */
export function getPlayerOrders(params?: OrderListParams, config?: Partial<RequestConfig>) {
  return get<Order[]>('/player/orders', params, config)
}

/**
 * 接受订单
 */
export function acceptOrder(orderId: number) {
  return post<void>(`/player/orders/${orderId}/accept`)
}

/**
 * 拒绝订单
 */
export function rejectOrder(orderId: number, reason?: string) {
  return post<void>(`/player/orders/${orderId}/reject`, { reason })
}

/**
 * 开始服务
 */
export function startService(orderId: number) {
  return post<void>(`/player/orders/${orderId}/start`)
}

/**
 * 完成服务
 */
export function finishService(orderId: number) {
  return put<void>(`/player/orders/${orderId}/complete`)
}

// ============ 评价 API ============

/**
 * 获取我的评价列表
 */
export function getMyReviews(params?: { page?: number; page_size?: number; type?: string }) {
  return get<any>('/users/reviews/my', params)
}

/**
 * 创建评价（独立接口，不依赖 orderId）
 */
export function createReview(data: {
  orderId: number
  rating: number
  content?: string
  tags?: string[]
  images?: string[]
}) {
  return post<void>('/users/reviews', data)
}

export default {
  getOrders,
  getOrderDetail,
  createOrder,
  cancelOrder,
  completeOrder,
  refundOrder,
  submitReview,
  getOrderPaymentStatus,
  getPlayerOrders,
  acceptOrder,
  rejectOrder,
  startService,
  finishService,
  getMyReviews,
  createReview,
}
