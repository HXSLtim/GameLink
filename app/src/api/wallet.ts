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
  WithdrawMethod,
  WithdrawParams,
  WithdrawRecord,
  WithdrawResult,
  VipDetailInfo,
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

// 充值结果（钱包充值）
export interface RechargeResult {
  orderId: number
  paymentId: number
  balanceCents: number
  payInfo?: {
    prepay_id?: string
    code_url?: string
    trade_no?: string
    qr_code?: string
  }
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

// 优惠券统计
export interface CouponStats {
  total: number
  available: number
  used: number
  expired: number
  locked?: number
  deleted?: number
}

// VIP 状态
export interface VipStatus {
  vipUnlocked: boolean
  currentLevel?: {
    level: number
    name?: string
  }
}

/**
 * 获取钱包余额信息
 */
export function getWalletInfo() {
  return get<WalletInfo>('/user/wallet/balance')
}

/**
 * 获取交易记录
 */
export function getTransactions(
  params?: {
  page?: number
  pageSize?: number
    page_size?: number
    type?: TransactionType
    status?: TransactionStatus
  },
  config?: Partial<RequestConfig>
) {
  return get<Transaction[]>('/user/wallet/transactions', params, config)
}

/**
 * 充值
 */
export function recharge(data: RechargeParams) {
  return post<RechargeResult>('/user/wallet/recharge', data)
}

/**
 * 支付订单
 */
export function payOrder(data: PaymentParams) {
  return post<PaymentResult>('/user/payments', data)
}

/**
 * 查询支付状态
 */
export function getPaymentStatus(paymentId: number) {
  return get<{ status: PaymentStatus; orderId: number }>(`/user/payments/${paymentId}`)
}

/**
 * 取消支付
 */
export function cancelPayment(paymentId: number) {
  return post<void>(`/user/payments/${paymentId}/cancel`)
}

/**
 * 获取当前用户优惠券统计
 */
export function getCouponStats() {
  return get<CouponStats>('/user/coupons/stats')
}

/**
 * 获取当前用户 VIP 状态
 */
export function getVipStatus() {
  return get<VipStatus>('/user/vip/status')
}

export default {
  getWalletInfo,
  getTransactions,
  recharge,
  payOrder,
  getPaymentStatus,
  cancelPayment,
  getCouponStats,
  getVipStatus,
}
