/**
 * 支付API模块
 */

import { apiClient } from '@/api/client';
import type {
  Payment,
  PaymentMethod,
  PaymentStatus,
  CreatePaymentRequest,
  CreatePaymentResponse,
  GetPaymentsParams,
  GetPaymentsResponse,
  RefundRequest,
  RefundResponse,
} from '@/shared/types/payment';

/**
 * 创建支付
 */
export const createPayment = async (data: CreatePaymentRequest): Promise<CreatePaymentResponse> => {
  const response = await apiClient.post<{ data: CreatePaymentResponse }>('/payments', data);
  return response.data;
};

/**
 * 获取支付列表
 */
export const getPayments = async (params?: GetPaymentsParams): Promise<GetPaymentsResponse> => {
  const response = await apiClient.get<{ data: GetPaymentsResponse }>('/payments', { params });
  return response.data;
};

/**
 * 获取支付详情
 */
export const getPaymentById = async (id: number): Promise<Payment> => {
  const response = await apiClient.get<{ data: Payment }>(`/payments/${id}`);
  return response.data;
};

/**
 * 获取订单支付记录
 */
export const getOrderPayments = async (orderId: number): Promise<Payment[]> => {
  const response = await apiClient.get<{ data: Payment[] }>(`/orders/${orderId}/payments`);
  return response.data;
};

/**
 * 处理支付回调
 */
export const handlePaymentCallback = async (callbackData: Record<string, any>): Promise<void> => {
  await apiClient.post('/payments/callback', callbackData);
};

/**
 * 查询支付状态
 */
export const queryPaymentStatus = async (paymentId: number): Promise<PaymentStatus> => {
  const response = await apiClient.get<{ data: { status: PaymentStatus } }>(
    `/payments/${paymentId}/status`
  );
  return response.data.status;
};

/**
 * 申请退款
 */
export const requestRefund = async (data: RefundRequest): Promise<RefundResponse> => {
  const response = await apiClient.post<{ data: RefundResponse }>('/payments/refund', data);
  return response.data;
};

/**
 * 获取退款记录
 */
export const getRefunds = async (params?: {
  paymentId?: number;
  orderId?: number;
  page?: number;
  pageSize?: number;
}): Promise<{
  list: Payment[];
  total: number;
  page: number;
  pageSize: number;
}> => {
  const response = await apiClient.get<{ data: { list: Payment[]; total: number; page: number; pageSize: number } }>(
    '/payments/refunds',
    { params }
  );
  return response.data;
};