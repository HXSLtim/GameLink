import { httpClient } from "@/services/http/client"
import type {
  PaymentChannelAdapter,
  PaymentCreatePayload,
  PaymentCreateResponse,
  PaymentOrderInput,
  PaymentStartResult,
} from "@/services/payment/types"

function createRequestId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID()
  }
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

export const alipayAdapter: PaymentChannelAdapter = {
  async createOrder(input: PaymentOrderInput): Promise<PaymentStartResult> {
    const payload: PaymentCreatePayload = {
      orderId: input.orderId,
      method: "alipay",
      requestId: createRequestId(),
    }

    const response = await httpClient.post<PaymentCreateResponse>("/user/payments", payload)
    const body = response.data

    if (!body.success) {
      return { success: false, message: body.message }
    }

    const payInfo = body.data?.payInfo ?? {}
    const redirectUrl =
      typeof payInfo.qr_code === "string"
        ? payInfo.qr_code
        : typeof payInfo.payUrl === "string"
          ? payInfo.payUrl
          : undefined

    return {
      success: true,
      paymentId: body.data?.paymentId,
      payInfo,
      redirectUrl,
      message: body.message,
    }
  },
}
