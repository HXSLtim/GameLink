import type { ReactElement } from "react"
import { useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { usePolling } from "@/hooks/usePolling"
import { getPaymentStatus } from "@/services/payment"
import type { PaymentStatusData } from "@/services/payment/types"

const terminalStatus = new Set(["paid", "failed", "refunded"])
const POLLING_TIMEOUT = 5 * 60 * 1000 // 5 minutes
const INITIAL_DELAY = 3000 // 3 seconds
const MAX_DELAY = 30000 // 30 seconds

function statusText(status: string): string {
  switch (status) {
    case "paid":
      return "支付成功"
    case "failed":
      return "支付失败"
    case "refunded":
      return "已退款"
    default:
      return "支付处理中"
  }
}

export function PaymentResultPage(): ReactElement {
  const navigate = useNavigate()
  const params = useParams<{ paymentId: string }>()
  const [searchParams] = useSearchParams()
  const fallbackPaymentId = Number(searchParams.get("paymentId") ?? searchParams.get("payment_id"))
  const paymentId = Number(params.paymentId ?? fallbackPaymentId)

  const [payment, setPayment] = useState<PaymentStatusData | null>(null)

  const validPaymentId = Number.isInteger(paymentId) && paymentId > 0

  const statusLabel = useMemo(() => {
    if (!payment) {
      return "未查询"
    }
    return statusText(payment.status)
  }, [payment])

  const { state: pollingState, stop: stopPolling } = usePolling<PaymentStatusData>({
    fetcher: async () => {
      const response = await getPaymentStatus(paymentId)
      if (!response.success) {
        throw new Error(response.message || "查询支付状态失败")
      }
      return response.data
    },
    shouldStop: (data) => {
      if (!data) return false
      return terminalStatus.has(data.status)
    },
    initialDelay: INITIAL_DELAY,
    maxDelay: MAX_DELAY,
    backoffFactor: 1.5,
    timeout: POLLING_TIMEOUT,
    immediate: validPaymentId,
    onSuccess: (data) => {
      setPayment(data)
    },
    onError: () => {
      // Error is handled in state
    },
    onStop: (reason) => {
      if (reason === "timeout") {
        // Could show a message to user about timeout
        console.warn("Payment polling timed out after", POLLING_TIMEOUT, "ms")
      }
    },
  })

  const isLoading = pollingState.isLoading && pollingState.attemptCount === 0
  const isPolling = pollingState.isPolling && !terminalStatus.has(payment?.status ?? "")
  const errorMessage = pollingState.error?.message ?? null

  // Invalid payment ID case
  if (!validPaymentId) {
    return (
      <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
        <h1 className="m-0 text-xl font-semibold text-slate-900">支付结果</h1>
        <p className="m-0 text-sm text-red-700">无效的支付单号</p>
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          返回工作台
        </Button>
      </section>
    )
  }

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <h1 className="m-0 text-xl font-semibold text-slate-900">支付结果</h1>
      <p className="m-0 text-sm text-slate-600">支付单号：{paymentId}</p>
      <p className="m-0 text-sm text-slate-600">
        状态：
        {isLoading ? "查询中..." : statusLabel}
        {isPolling && !isLoading ? " (持续查询中...)" : ""}
      </p>
      {payment?.paidAt ? <p className="m-0 text-sm text-slate-600">支付时间：{payment.paidAt}</p> : null}
      {errorMessage ? <p className="m-0 text-sm text-red-700">{errorMessage}</p> : null}
      {pollingState.attemptCount > 1 && !terminalStatus.has(payment?.status ?? "") ? (
        <p className="m-0 text-xs text-slate-400">已查询 {pollingState.attemptCount} 次，下次将在 {Math.round(pollingState.currentDelay / 1000)} 秒后</p>
      ) : null}
      <div className="flex gap-3">
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          返回工作台
        </Button>
        {isPolling ? (
          <Button variant="outline" onClick={stopPolling}>
            停止查询
          </Button>
        ) : (
          <Button variant="outline" onClick={() => window.location.reload()}>
            立即刷新
          </Button>
        )}
      </div>
    </section>
  )
}
