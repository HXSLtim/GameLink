import type { OrderStatus } from '@/types/order'
import type { CertStatus, DashboardStatus, ResultStatusType } from '@/types/status'

export type { CertStatus, DashboardStatus, ResultStatusType } from '@/types/status'

export interface OrderStatusPreset {
  icon: string
  text: string
  desc: string
  className: string
  borderColor: string
}

export interface CertStatusPreset {
  icon: string
  title: string
  desc: string
  borderColor: string
}

export interface ResultStatusPreset {
  icon: string
  borderColor: string
}

export interface DashboardStatusPreset {
  text: string
  borderColor: string
  borderWidth: string
}

const canceledPreset: OrderStatusPreset = {
  icon: 'close-circle',
  text: '已取消',
  desc: '订单已取消',
  className: 'default',
  borderColor: 'var(--color-text-placeholder)',
}

const orderStatusPresets: Record<OrderStatus | 'cancelled', OrderStatusPreset> = {
  pending: {
    icon: 'clock',
    text: '待支付',
    desc: '请在规定时间内完成支付',
    className: 'pending',
    borderColor: 'var(--color-warning)',
  },
  confirmed: {
    icon: 'checkmark-circle-fill',
    text: '已支付',
    desc: '等待陪玩师接单',
    className: 'confirmed',
    borderColor: 'var(--color-info)',
  },
  in_progress: {
    icon: 'grid-fill',
    text: '服务中',
    desc: '陪玩师正在为您服务',
    className: 'in-progress',
    borderColor: 'var(--color-primary)',
  },
  completed: {
    icon: 'checkmark-circle-fill',
    text: '已完成',
    desc: '服务已完成，欢迎评价',
    className: 'completed',
    borderColor: 'var(--color-success)',
  },
  canceled: canceledPreset,
  cancelled: canceledPreset,
  refunding: {
    icon: 'reload',
    text: '退款中',
    desc: '正在处理退款',
    className: 'refunding',
    borderColor: 'var(--color-info)',
  },
  refunded: {
    icon: 'red-packet',
    text: '已退款',
    desc: '退款已到账',
    className: 'refunded',
    borderColor: 'var(--color-text-placeholder)',
  },
  disputed: {
    icon: 'info-circle',
    text: '争议中',
    desc: '订单存在争议，处理中',
    className: 'error',
    borderColor: 'var(--color-error)',
  },
}

export function getOrderStatusPreset(status: OrderStatus | 'cancelled'): OrderStatusPreset {
  return orderStatusPresets[status] || orderStatusPresets.pending
}

const certStatusPresets: Record<CertStatus, CertStatusPreset> = {
  none: {
    icon: 'edit-pen',
    title: '未认证',
    desc: '完成认证后即可开始接单',
    borderColor: 'var(--color-text-placeholder)',
  },
  pending: {
    icon: 'clock',
    title: '审核中',
    desc: '预计1-3个工作日内完成审核',
    borderColor: 'var(--color-warning)',
  },
  approved: {
    icon: 'checkmark-circle-fill',
    title: '已认证',
    desc: '恭喜您已通过认证',
    borderColor: 'var(--color-success)',
  },
  rejected: {
    icon: 'close-circle',
    title: '认证失败',
    desc: '请根据反馈修改后重新提交',
    borderColor: 'var(--color-error)',
  },
}

export function getCertStatusPreset(status: CertStatus): CertStatusPreset {
  return certStatusPresets[status]
}

const resultStatusPresets: Record<ResultStatusType, ResultStatusPreset> = {
  success: { icon: 'checkmark-circle-fill', borderColor: 'var(--color-success)' },
  pending: { icon: 'clock', borderColor: 'var(--color-warning)' },
  failed: { icon: 'close-circle', borderColor: 'var(--color-error)' },
  warning: { icon: 'info-circle', borderColor: 'var(--color-warning)' },
}

export function getResultStatusPreset(type: ResultStatusType): ResultStatusPreset {
  return resultStatusPresets[type]
}

const dashboardStatusPresets: Record<DashboardStatus, DashboardStatusPreset> = {
  online: {
    text: '在线接单中',
    borderColor: 'var(--color-success)',
    borderWidth: '4rpx',
  },
  busy: {
    text: '忙碌中',
    borderColor: 'var(--color-warning)',
    borderWidth: '4rpx',
  },
  offline: {
    text: '已下线',
    borderColor: 'var(--color-border)',
    borderWidth: '1rpx',
  },
}

export function getDashboardStatusPreset(status: DashboardStatus): DashboardStatusPreset {
  return dashboardStatusPresets[status]
}
