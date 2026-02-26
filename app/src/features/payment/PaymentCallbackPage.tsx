import type { ReactElement } from "react"
import { useEffect } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { parsePaymentCallback } from "@/services/payment/callback"
import { isLoggedIn } from "@/services/session"

export function PaymentCallbackPage(): ReactElement {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const callback = parsePaymentCallback(searchParams)
  const paymentId = callback.paymentId

  useEffect(() => {
    if (!paymentId) {
      return
    }

    const target = `/payment/result/${paymentId}`
    if (isLoggedIn()) {
      navigate(target, { replace: true })
      return
    }

    navigate("/login", { replace: true, state: { from: target } })
  }, [navigate, paymentId])

  if (paymentId) {
    return (
      <section className="mt-8 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
        <p className="m-0 text-sm text-slate-600">支付回跳处理中，正在进入结果页...</p>
      </section>
    )
  }

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <h1 className="m-0 text-xl font-semibold text-slate-900">支付回跳参数缺失</h1>
      <p className="m-0 text-sm text-slate-600">未能从回跳地址中识别支付单号，请返回工作台重新发起支付。</p>
      <Button variant="outline" onClick={() => navigate("/dashboard")}>返回工作台</Button>
    </section>
  )
}
