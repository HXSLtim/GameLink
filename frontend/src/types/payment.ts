import { BaseEntity } from './user';
import { Currency } from './order';

/**
 * 支付方式枚举
 */
export type PaymentMethod = 'wechat' | 'alipay';

/**
 * 支付状态枚举
 */
export type PaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded';

/**
 * 支付实体 - 与后端model.Payment保持一致
 */
export interface Payment extends BaseEntity {
  order_id: number;
  user_id: number;
  method: PaymentMethod;
  amount_cents: number;
  currency?: Currency;
  status: PaymentStatus;
  provider_trade_no?: string;
  provider_raw?: any; // JSON对象
  paid_at?: string;
  refunded_at?: string;
}

/**
 * 支付列表查询参数
 */
export interface PaymentListQuery {
  page?: number;
  page_size?: number;
  order_id?: number;
  user_id?: number;
  method?: PaymentMethod;
  status?: PaymentStatus;
  currency?: Currency;
  date_from?: string;
  date_to?: string;
  sort_by?: 'created_at' | 'updated_at' | 'amount_cents' | 'paid_at';
  sort_order?: 'asc' | 'desc';
}

/**
 * 创建支付请求
 */
export interface CreatePaymentRequest {
  order_id: number;
  user_id: number;
  method: PaymentMethod;
  amount_cents: number;
  currency?: Currency;
}

/**
 * 支付状态显示文本
 */
export const PAYMENT_STATUS_TEXT: Record<PaymentStatus, string> = {
  pending: '待支付',
  paid: '已支付',
  failed: '支付失败',
  refunded: '已退款',
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
};

/**
 * 支付方式显示文本
 */
export const PAYMENT_METHOD_TEXT: Record<PaymentMethod, string> = {
  wechat: '微信支付',
  alipay: '支付宝',
};

/**
 * 支付方式图标
 */
export const PAYMENT_METHOD_ICON: Record<PaymentMethod, string> = {
  wechat: '💚',
  alipay: '💙',
};

/**
 * 金额格式化（分转元）
 */
export const formatAmount = (cents: number, currency: Currency = 'CNY'): string => {
  const amount = cents / 100;
  return `${amount.toFixed(2)} ${currency}`;
};
