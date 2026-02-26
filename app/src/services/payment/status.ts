import { httpClient } from "@/services/http/client"
import type { PaymentStatusResponse } from "@/services/payment/types"

export async function getPaymentStatus(paymentId: number): Promise<PaymentStatusResponse> {
  const response = await httpClient.get<PaymentStatusResponse>(`/user/payments/${paymentId}`)
  return response.data
}
