import type { BaseEntity } from './user';
import type { Currency } from './order';
import React from 'react';

/**
 * 图标属性接口 (临时定义，稍后替换为实际图标组件)
 */
export interface IconProps {
  size?: number;
  color?: string;
  className?: string;
}

/**
 * 支付方式枚举
 */
export type PaymentMethod = 'wechat' | 'alipay' | 'balance';

/**
 * 支付状态枚举
 */
export type PaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded' | 'cancelled';

/**
 * 支付实体 - 与后端 model.Payment 保持一致（camelCase）
 */
export interface Payment extends BaseEntity {
  orderId: number;
  userId: number;
  method: PaymentMethod;
  amountCents: number;
  currency?: Currency;
  status: PaymentStatus;
  providerTradeNo?: string;
  providerRaw?: any; // 第三方支付返回的原始数据
  paidAt?: string;
  refundedAt?: string;
}

/**
 * 创建支付响应
 */
export interface CreatePaymentResponse {
  id: number;
  orderId: number;
  method: PaymentMethod;
  amountCents: number;
  currency: Currency;
  status: PaymentStatus;
  providerTradeNo?: string;
  paidAt?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 支付查询参数
 */
export interface GetPaymentsParams {
  page?: number;
  page_size?: number;
  orderId?: number;
  userId?: number;
  method?: PaymentMethod;
  status?: PaymentStatus;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * 支付详情（包含扩展信息）
 */
export interface PaymentDetail extends Payment {
  // 关联订单信息
  order?: {
    id: number;
    title?: string;
    status?: string;
    userId?: number;
    playerId?: number;
  };
  // 关联用户信息
  user?: {
    id: number;
    name: string;
    phone?: string;
    email?: string;
  };
  // 退款信息
  refundInfo?: {
    refundAmount: number;
    refundReason: string;
    refundedAt: string;
  };
}

/**
 * 支付列表查询参数
 */
export interface PaymentListQuery {
  page?: number;
  pageSize?: number;
  keyword?: string; // 搜索关键词（交易号/订单ID）
  orderId?: number;
  userId?: number;
  method?: PaymentMethod;
  status?: PaymentStatus;
  currency?: Currency;
  dateFrom?: string;
  dateTo?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'amountCents';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 创建支付请求
 */
export interface CreatePaymentRequest {
  orderId: number;
  method: PaymentMethod;
  amountCents: number;
  currency?: Currency;
}

/**
 * 更新支付请求
 */
export interface UpdatePaymentRequest {
  status?: PaymentStatus;
  providerTxId?: string;
}

/**
 * 查询支付响应
 */
export interface GetPaymentsResponse {
  list: Payment[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 退款请求
 */
export interface RefundPaymentRequest {
  reason: string;
  amountCents?: number; // 部分退款金额，不传则全额退款
}

/**
 * 退款请求（API使用）
 */
export interface RefundRequest extends RefundPaymentRequest {
  paymentId: number;
}

/**
 * 退款响应
 */
export interface RefundResponse {
  id: number;
  paymentId: number;
  amountCents: number;
  reason: string;
  status: 'pending' | 'completed' | 'failed';
  createdAt: string;
}

/**
 * 支付状态显示文本
 */
export const PAYMENT_STATUS_TEXT: Record<PaymentStatus, string> = {
  pending: '待支付',
  paid: '已支付',
  failed: '支付失败',
  refunded: '已退款',
  cancelled: '已取消',
};

/**
 * 支付状态徽章类型
 */
export const PAYMENT_STATUS_BADGE: Record<
  PaymentStatus,
  'success' | 'warning' | 'error' | 'processing'
> = {
  pending: 'processing',
  paid: 'success',
  failed: 'error',
  refunded: 'warning',
  cancelled: 'error',
};

/**
 * 支付方式显示文本
 */
export const PAYMENT_METHOD_TEXT: Record<PaymentMethod, string> = {
  wechat: '微信支付',
  alipay: '支付宝',
  balance: '余额支付',
};

/**
 * 支付方式图标（SVG 组件）- 占位符实现
 */
export const PAYMENT_METHOD_ICON: Record<PaymentMethod, React.FC<IconProps>> = {
  wechat: ({ size = 24, color = '#07C160' }) =>
    React.createElement('div', {
      style: { width: size, height: size, backgroundColor: color, borderRadius: '50%' },
      title: '微信支付',
    }),
  alipay: ({ size = 24, color = '#1677FF' }) =>
    React.createElement('div', {
      style: { width: size, height: size, backgroundColor: color, borderRadius: '50%' },
      title: '支付宝',
    }),
  balance: ({ size = 24, color = '#FFC107' }) =>
    React.createElement('div', {
      style: { width: size, height: size, backgroundColor: color, borderRadius: '50%' },
      title: '余额支付',
    }),
};

// 为了兼容性，导出空组件对象
export const WechatPayIcon = PAYMENT_METHOD_ICON.wechat;
export const AlipayIcon = PAYMENT_METHOD_ICON.alipay;
export const BalanceIcon = PAYMENT_METHOD_ICON.balance;


/**
 * 金额格式化（分转元）
 */
export const formatAmount = (cents: number, currency?: Currency): string => {
  const amount = cents / 100;
  const curr = currency || 'CNY';
  return `${amount.toFixed(2)} ${curr}`;
};
