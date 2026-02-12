import type { PlayerSummary } from '@/types/player'
import type { OrderReviewData } from '@/types/review'
import type { TabItem } from '@/types/ui'

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'canceled'
  | 'refunding'
  | 'refunded'
  | 'disputed'

export type OrderViewMode = 'user' | 'player'

export type OrderTabKey = 'all' | OrderStatus
export type OrderQuickEntryStatus = Extract<OrderStatus, 'pending' | 'in_progress' | 'completed' | 'refunding'>
export type OrderActionKey =
  | 'cancel'
  | 'pay'
  | 'contact'
  | 'complete'
  | 'review'
  | 'reorder'
  | 'refund'
  | 'viewDispute'
  | 'reject'
  | 'accept'
  | 'start'
  | 'detail'

export type OrderTabItem = Omit<TabItem, 'key'> & { key: OrderTabKey }

export interface OrderPerson {
  id: number
  nickname: string
  avatar?: string
}

export interface Order {
  id: number
  orderNo: string
  status: OrderStatus
  player?: OrderPerson
  user?: OrderPerson
  gameName: string
  serviceName: string
  quantity: number
  unit: string
  totalAmount?: number
  earnings?: number
  remark?: string
  createdAt: string
  reviewed?: boolean
}

export interface OrderDetailData {
  id: number
  orderNo: string
  status: OrderStatus
  player: PlayerSummary
  gameName: string
  serviceName: string
  quantity: number
  unit: string
  gameAccount?: string
  remark?: string
  serviceFee: number
  couponDiscount: number
  vipDiscount: number
  totalAmount: number
  paymentMethod?: OrderPaymentMethod
  createdAt: string
  paidAt?: string
  startedAt?: string
  completedAt?: string
  scheduledStart?: string
  review?: OrderReviewData
  refund?: RefundInfo
}

export interface InfoItem {
  label: string
  value: string
  copyable?: boolean
}

export interface FeeItem {
  label: string
  value: number
  isDiscount?: boolean
}

export interface OrderCountSummary {
  pending: number
  inProgress: number
  toReview: number
  refunding: number
}

export interface OrderPaymentInfo {
  orderId?: number
  orderNo?: string
  amount?: number
  method?: OrderPaymentMethod
  paidAt?: string
}

export type OrderPaymentMethod = 'wechat' | 'alipay' | 'wallet' | 'combined'

export type OrderPaymentStatus = 'success' | 'pending' | 'failed' | 'paid'

// ============ 退款相关类型 ============

/** 退款状态 */
export type RefundStatus = 'pending' | 'processing' | 'refunded' | 'rejected'

/** 退款详情 */
export interface RefundInfo {
  id: number
  status: RefundStatus
  reason: string
  amount: number
  createdAt: string
  processedAt?: string
  rejectReason?: string
}

// ============ 订单进度相关类型 ============

/** 订单进度步骤 */
export interface OrderProgressStep {
  key: string
  label: string
  icon: string
  /** 是否已完成 */
  done: boolean
  /** 是否为当前步骤 */
  active: boolean
  /** 完成时间 */
  time?: string
}
