export type { OrderStatus, OrderViewMode, OrderPerson, Order, OrderActionKey } from '@/types/order'

export interface ActionButton {
  key: import('@/types/order').OrderActionKey
  label: string
  type: 'primary' | 'default'
  plain: boolean
}
