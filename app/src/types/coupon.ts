export interface Coupon {
  id: number
  name: string
  discount: number
  minAmount?: number
  expireAt?: string
  expireDate?: string
  type?: 'fixed' | 'percent'
}
