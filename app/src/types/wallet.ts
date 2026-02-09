export type TransactionType =
  | 'recharge'
  | 'withdraw'
  | 'withdrawal'
  | 'payment'
  | 'refund'
  | 'earning'
  | 'bonus'
  | 'commission'

export type PaymentStatus = 'pending' | 'paid' | 'failed'

export type TransactionStatus = 'pending' | 'success' | 'failed'

export type RechargeMethod = 'wechat' | 'alipay'

export interface WalletInfo {
  balanceCents: number
  frozenCents: number
  totalRechargeCents?: number
  totalSpentCents?: number
  vipLevel?: number
  couponCount?: number
}

export interface TransactionData {
  id: number
  type: TransactionType
  title: string
  description: string
  amount: number // 分
  createdAt: string
}

export interface WalletState {
  balance: number
  frozenBalance: number
  vipLevel: number
  couponCount: number
  totalSpent: number
  totalRecharge: number
}

export interface WalletSummary {
  balanceCents: number
  totalRechargeCents: number
  totalSpentCents: number
}

export interface AmountOption {
  value: number
  bonus?: number
}

export interface PaymentMethod {
  value: string
  name: string
  icon: string
  enabled: boolean
  tip?: string
}
