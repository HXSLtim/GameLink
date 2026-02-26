import { alipayAdapter } from "@/services/payment/alipay"
import type { PaymentChannel, PaymentChannelAdapter } from "@/services/payment/types"
import { wechatAdapter } from "@/services/payment/wechat"
export { getPaymentStatus } from "@/services/payment/status"

const adapters: Record<PaymentChannel, PaymentChannelAdapter> = {
  alipay: alipayAdapter,
  wechat: wechatAdapter,
}

export function getPaymentAdapter(channel: PaymentChannel): PaymentChannelAdapter {
  return adapters[channel]
}
