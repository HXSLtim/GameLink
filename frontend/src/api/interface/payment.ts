/**
 * 支付相关接口定义
 */

import type { BaseEntity } from '@/shared/types/api';

/**
 * 支付方式
 */
export enum PaymentMethod {
  WECHAT = 'wechat',
  ALIPAY = 'alipay',
  BALANCE = 'balance',
}

/**
 * 支付状态
 */
export enum PaymentStatus {
  PENDING = 'pending',
  PAID = 'paid',
  FAILED = 'failed',
  REFUNDED = 'refunded',
  CANCELLED = 'cancelled',
}

/**
 * 支付实体
 */
export interface Payment extends BaseEntity {
  id: number;
  orderId: number;
  orderNo: string;
  userId: number;
  amount: number;
  method: PaymentMethod;
  status: PaymentStatus;
  transactionId?: string;
  paidAt?: string;
  refundedAt?: string;
  refundReason?: string;
}

/**
 * 创建支付请求
 */
export interface CreatePaymentRequest {
  orderId: number;
  method: PaymentMethod;
}

/**
 * 支付响应
 */
export interface CreatePaymentResponse {
  paymentId: number;
  paymentUrl?: string;
  qrCode?: string;
}

/**
 * 支付回调
 */
export interface PaymentCallback {
  paymentId: number;
  transactionId: string;
  status: PaymentStatus;
  amount: number;
  paidAt: string;
}

/**
 * 支付列表查询参数
 */
export interface GetPaymentsParams {
  orderId?: number;
  userId?: number;
  status?: PaymentStatus;
  method?: PaymentMethod;
  page?: number;
  pageSize?: number;
  startDate?: string;
  endDate?: string;
}

/**
 * 支付列表响应
 */
export interface GetPaymentsResponse {
  list: Payment[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 退款请求
 */
export interface RefundRequest {
  paymentId: number;
  amount?: number;
  reason: string;
}

/**
 * 退款响应
 */
export interface RefundResponse {
  refundId: number;
  status: string;
  message: string;
}