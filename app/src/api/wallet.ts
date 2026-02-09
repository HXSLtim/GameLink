/**
 * 钱包相关 API
 */

import { get, post, type RequestConfig } from './request'
import type {
  PaymentStatus,
  RechargeMethod,
  TransactionStatus,
  TransactionType,
  WalletInfo,
} from '@/types/wallet'
import type { OrderPaymentMethod } from '@/types/order'

// 交易记录
export interface Transaction {
  id: number
  type: TransactionType
  title: string
  description?: string
  amountCents: number
  balanceAfterCents: number
  orderId?: number
  status: TransactionStatus
  createdAt: string
}

// 充值参数
export interface RechargeParams {
  amountCents: number
  method: RechargeMethod
}

// 支付参数
export interface PaymentParams {
  orderId: number
  method: OrderPaymentMethod
  requestId: string
  walletAmountCents?: number
  thirdPartyMethod?: RechargeMethod
}

// 支付结果
export interface PaymentResult {
  paymentId: number
  status: PaymentStatus
  payInfo?: {
    prepay_id?: string
    code_url?: string
    trade_no?: string
    qr_code?: string
  }
  newBalanceCents?: number
}

/**
 * 获取钱包余额信息
 */
export function getWalletInfo() {
  return get<WalletInfo>('/users/wallet/balance')
}

/**
 * 获取交易记录
 */
export function getTransactions(
  params?: {
  page?: number
  page_size?: number
    type?: TransactionType
    status?: TransactionStatus
  },
  config?: Partial<RequestConfig>
) {
  return get<Transaction[]>('/users/wallet/transactions', params, config)
}

/**
 * 充值
 */
export function recharge(data: RechargeParams) {
  return post<PaymentResult>('/users/wallet/recharge', data)
}

/**
 * 支付订单
 */
export function payOrder(data: PaymentParams) {
  return post<PaymentResult>('/users/payments', data)
}

/**
 * 查询支付状态
 */
export function getPaymentStatus(paymentId: number) {
  return get<{ status: PaymentStatus; orderId: number }>(`/users/payments/${paymentId}`)
}

/**
 * 取消支付
 */
export function cancelPayment(paymentId: number) {
  return post<void>(`/users/payments/${paymentId}/cancel`)
}

export default {
  getWalletInfo,
  getTransactions,
  recharge,
  payOrder,
  getPaymentStatus,
  cancelPayment,
}
