/**
 * 订单API模块
 */

import { apiClient } from '@/api/client';
import type {
  Order,
  OrderStatus,
  CreateOrderRequest,
  UpdateOrderRequest,
  GetOrdersParams,
  GetOrdersResponse,
  OrderStats,
} from '@/shared/types/order';

/**
 * 创建订单
 */
export const createOrder = async (data: CreateOrderRequest): Promise<Order> => {
  const response = await apiClient.post<{ data: Order }>('/orders', data);
  return response.data;
};

/**
 * 获取订单列表
 */
export const getOrders = async (params?: GetOrdersParams): Promise<GetOrdersResponse> => {
  const response = await apiClient.get<{ data: GetOrdersResponse }>('/orders', { params });
  return response.data;
};

/**
 * 获取订单详情
 */
export const getOrderById = async (id: number): Promise<Order> => {
  const response = await apiClient.get<{ data: Order }>(`/orders/${id}`);
  return response.data;
};

/**
 * 更新订单
 */
export const updateOrder = async (id: number, data: UpdateOrderRequest): Promise<Order> => {
  const response = await apiClient.put<{ data: Order }>(`/orders/${id}`, data);
  return response.data;
};

/**
 * 取消订单
 */
export const cancelOrder = async (id: number, reason?: string): Promise<void> => {
  await apiClient.post(`/orders/${id}/cancel`, { reason });
};

/**
 * 完成订单
 */
export const completeOrder = async (id: number): Promise<void> => {
  await apiClient.post(`/orders/${id}/complete`);
};

/**
 * 获取用户订单列表
 */
export const getUserOrders = async (userId: number, params?: Omit<GetOrdersParams, 'userId'>): Promise<GetOrdersResponse> => {
  const response = await apiClient.get<{ data: GetOrdersResponse }>(`/users/${userId}/orders`, { params });
  return response.data;
};

/**
 * 获取陪玩师订单列表
 */
export const getPlayerOrders = async (playerId: number, params?: Omit<GetOrdersParams, 'playerId'>): Promise<GetOrdersResponse> => {
  const response = await apiClient.get<{ data: GetOrdersResponse }>(`/players/${playerId}/orders`, { params });
  return response.data;
};

/**
 * 获取订单统计
 */
export const getOrderStats = async (params?: {
  startDate?: string;
  endDate?: string;
  userId?: number;
  playerId?: number;
}): Promise<OrderStats> => {
  const response = await apiClient.get<{ data: OrderStats }>('/orders/stats', { params });
  return response.data;
};

/**
 * 开始订单（陪玩师操作）
 */
export const startOrder = async (id: number): Promise<void> => {
  await apiClient.post(`/orders/${id}/start`);
};

/**
 * 确认订单（用户确认陪玩师）
 */
export const confirmOrder = async (id: number): Promise<void> => {
  await apiClient.post(`/orders/${id}/confirm`);
};