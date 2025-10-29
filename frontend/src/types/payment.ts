import type { BaseEntity } from './user';
import type { Currency } from './order';

/**
 * 支付方式枚举
 */
export type PaymentMethod = 'wechat' | 'alipay' | 'balance';

/**
 * 支付状态枚举
 */
export type PaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded' | 'cancelled';

/**
 * 支付实体 - 与后端 model.Payment 保持一致
 */
export interface Payment extends BaseEntity {
  orderId: number;
  userId: number;
  amountCents: number;
  currency?: Currency;
  method: string;
  status: string;
  transactionId?: string;
  providerTxId?: string;
}

/**
 * 支付详情（包含扩展信息）
 */
export interface PaymentDetail extends Payment {
  order?: {
    id: number;
    orderNo?: string;
    title?: string;
  };
  user?: {
    id: number;
    name: string;
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
 * 退款请求
 */
export interface RefundPaymentRequest {
  reason: string;
  amountCents?: number; // 部分退款金额，不传则全额退款
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
 * 支付方式图标
 */
export const PAYMENT_METHOD_ICON: Record<PaymentMethod, string> = {
  wechat: '💚',
  alipay: '💙',
  balance: '💰',
};

/**
 * 金额格式化（分转元）
 */
export const formatAmount = (cents: number, currency?: Currency): string => {
  const amount = cents / 100;
  const curr = currency || 'CNY';
  return `${amount.toFixed(2)} ${curr}`;
};
