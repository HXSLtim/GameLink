import { apiClient } from '../../api/client';
import type {
  Order,
  OrderDetail,
  OrderListQuery,
  CreateOrderRequest,
  UpdateOrderRequest,
  AssignOrderRequest,
  ReviewOrderRequest,
  CancelOrderRequest,
  OrderStatsData,
} from '../../types/order';
import type { ListResult } from '../../types/api';

/**
 * 订单列表响应
 */
export type OrderListResponse = ListResult<Order>;

/**
 * 订单API服务
 */
export const orderApi = {
  /**
   * 获取订单列表
   */
  getList: (params: OrderListQuery): Promise<OrderListResponse> => {
    return apiClient.get('/api/v1/admin/orders', { params });
  },

  /**
   * 获取订单详情
   */
  getDetail: (id: number): Promise<OrderDetail> => {
    return apiClient.get(`/api/v1/admin/orders/${id}`);
  },

  /**
   * 创建订单
   */
  create: (data: CreateOrderRequest): Promise<Order> => {
    return apiClient.post('/api/v1/admin/orders', data);
  },

  /**
   * 更新订单信息
   */
  update: (id: number, data: UpdateOrderRequest): Promise<Order> => {
    return apiClient.put(`/api/v1/admin/orders/${id}`, data);
  },

  /**
   * 分配订单给陪玩师
   */
  assign: (id: number, data: AssignOrderRequest): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/assign`, data);
  },

  /**
   * 审核订单
   */
  review: (id: number, data: ReviewOrderRequest): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/review`, data);
  },

  /**
   * 取消订单
   */
  cancel: (id: number, data: CancelOrderRequest): Promise<Order> => {
    return apiClient.post(`/api/v1/admin/orders/${id}/cancel`, data);
  },

  /**
   * 删除订单
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/orders/${id}`);
  },

  /**
   * 获取订单操作日志
   */
  getLogs: (id: number): Promise<OrderDetail['logs']> => {
    return apiClient.get(`/api/v1/admin/orders/${id}/logs`);
  },

  /**
   * 获取订单统计
   */
  getStatistics: (): Promise<OrderStatsData> => {
    return apiClient.get('/api/v1/admin/stats/orders');
  },

  /**
   * 获取用户的订单列表
   */
  getUserOrders: (
    userId: number,
    params?: {
      page?: number;
      page_size?: number;
      status?: string[];
      date_from?: string;
      date_to?: string;
    },
  ): Promise<OrderListResponse> => {
    return apiClient.get(`/api/v1/admin/users/${userId}/orders`, { params });
  },
};

// 导出类型以便组件使用
export type { Order as OrderInfo } from '../../types/order';
export type { OrderDetail as OrderDetailType } from '../../types/order';
