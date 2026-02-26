export type PaymentChannel = "alipay" | "wechat"

export interface PaymentOrderInput {
  orderId: number
  amountFen: number
  subject: string
}

export interface PaymentCreatePayload {
  orderId: number
  method: PaymentChannel
  requestId: string
}

export interface PaymentCreateData {
  paymentId: number
  payInfo: Record<string, unknown>
  walletDeducted?: number
  thirdPartyAmount?: number
  walletPaidDirect?: boolean
}

export interface PaymentCreateResponse {
  success: boolean
  code: number
  message: string
  data: PaymentCreateData
}

export interface PaymentStartResult {
  success: boolean
  paymentId?: number
  payInfo?: Record<string, unknown>
  redirectUrl?: string
  message?: string
}

export type PaymentStatus = "pending" | "paid" | "failed" | "refunded"

export interface PaymentStatusData {
  paymentId: number
  orderId: number
  status: PaymentStatus
  paidAt?: string
}

export interface PaymentStatusResponse {
  success: boolean
  code: number
  message: string
  data: PaymentStatusData
}

export interface PaymentChannelAdapter {
  createOrder(input: PaymentOrderInput): Promise<PaymentStartResult>
}
