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

// ============================================
// 提现相关
// ============================================
export type WithdrawMethod = 'wechat' | 'alipay' | 'bank'

export type WithdrawStatus = 'pending' | 'processing' | 'success' | 'failed'

export interface WithdrawParams {
  amountCents: number
  method: WithdrawMethod
  account?: string
  realName?: string
}

export interface WithdrawResult {
  withdrawId: number
  status: WithdrawStatus
  amountCents: number
  estimatedTime?: string
}

export interface WithdrawRecord {
  id: number
  amountCents: number
  method: WithdrawMethod
  status: WithdrawStatus
  account?: string
  createdAt: string
  completedAt?: string
  remark?: string
}

// ============================================
// VIP 相关
// ============================================
export interface VipPrivilege {
  key: string
  label: string
  icon: string
  description: string
  /** 各等级是否拥有该特权 [普通, VIP1, VIP2, VIP3] */
  levels: boolean[]
}

export interface VipLevelInfo {
  level: number
  name: string
  discount: number
  minSpent: number
  color: string
  gradient: string
}

export interface VipDetailInfo {
  currentLevel: number
  currentSpentCents: number
  nextLevelSpentCents: number
  levels: VipLevelInfo[]
  privileges: VipPrivilege[]
}
