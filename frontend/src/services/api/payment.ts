import { apiClient } from '../../api/client';
import type {
  Payment,
  PaymentDetail,
  PaymentListQuery,
  CreatePaymentRequest,
  UpdatePaymentRequest,
  RefundPaymentRequest,
} from '../../types/payment';
import type { ListResult } from '../../types/api';

/**
 * 支付列表响应
 */
export type PaymentListResponse = ListResult<Payment>;

/**
 * 支付API服务
 */
export const paymentApi = {
  /**
   * 获取支付列表
   */
  getList: (params: PaymentListQuery): Promise<PaymentListResponse> => {
    return apiClient.get('/api/v1/admin/payments', { params });
  },

  /**
   * 获取支付详情
   */
  getDetail: (id: number): Promise<PaymentDetail> => {
    return apiClient.get(`/api/v1/admin/payments/${id}`);
  },

  /**
   * 创建支付
   */
  create: (data: CreatePaymentRequest): Promise<Payment> => {
    return apiClient.post('/api/v1/admin/payments', data);
  },

  /**
   * 更新支付信息
   */
  update: (id: number, data: UpdatePaymentRequest): Promise<Payment> => {
    return apiClient.put(`/api/v1/admin/payments/${id}`, data);
  },

  /**
   * 申请退款
   */
  refund: (id: number, data: RefundPaymentRequest): Promise<Payment> => {
    return apiClient.post(`/api/v1/admin/payments/${id}/refund`, data);
  },

  /**
   * 确认收款
   */
  capture: (id: number): Promise<Payment> => {
    return apiClient.post(`/api/v1/admin/payments/${id}/capture`, {});
  },

  /**
   * 删除支付记录
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/payments/${id}`);
  },

  /**
   * 获取支付操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/payments/${id}/logs`);
  },
};

