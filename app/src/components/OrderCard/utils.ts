import type { OrderStatus, OrderViewMode } from '@/types/order'
import type { OrderStatus as ApiOrderStatus } from '@/api/order'

type RawOrderStatus = ApiOrderStatus | OrderStatus

const mapStatus = (status: RawOrderStatus): OrderStatus => {
  switch (status) {
    case 'pending_payment':
      return 'pending'
    case 'pending_accept':
      return 'pending'
    case 'accepted':
      return 'confirmed'
    case 'in_progress':
    case 'completed':
    case 'disputed':
    case 'refunding':
    case 'refunded':
      return status
    case 'cancelled':
    case 'canceled':
      return 'canceled'
    case 'pending':
    case 'confirmed':
      return status
    default:
      return 'pending'
  }
}

export function normalizeOrderStatus(status: RawOrderStatus, viewMode: OrderViewMode): OrderStatus {
  const normalized = mapStatus(status)
  if (status === 'pending_accept' && viewMode === 'user') {
    return 'confirmed'
  }
  return normalized
}
