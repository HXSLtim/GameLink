import type { ReactElement } from "react"
import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { getPaymentAdapter } from "@/services/payment"

export function DashboardPage(): ReactElement {
  const navigate = useNavigate()
  const [orderIdInput, setOrderIdInput] = useState("1")
  const [status, setStatus] = useState("未发起支付")

  const startPayment = async (channel: "alipay" | "wechat") => {
    const orderId = Number(orderIdInput)
    if (!Number.isInteger(orderId) || orderId <= 0) {
      setStatus("请输入有效订单号")
      return
    }

    setStatus(`正在发起${channel === "alipay" ? "支付宝" : "微信"}支付...`)
    const adapter = getPaymentAdapter(channel)
    const result = await adapter.createOrder({
      orderId,
      amountFen: 100,
      subject: "测试支付",
    })

    if (!result.success) {
      setStatus(result.message ?? "支付创建失败")
      return
    }

    if (result.redirectUrl) {
      setStatus("支付参数创建成功，正在跳转支付页面")
      window.location.href = result.redirectUrl
      return
    }

    if (result.paymentId) {
      setStatus(`支付参数已创建，paymentId=${result.paymentId}`)
      navigate(`/payment/result/${result.paymentId}`)
      return
    }

    setStatus("支付参数已创建，但未返回 paymentId")
  }

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <h1 className="m-0 text-xl font-semibold text-slate-900">用户端工作台</h1>
      <p className="m-0 text-sm text-slate-600">这里是 React 版本迁移后的临时首页，支付链路已接入真实 `/user/payments` 合约。</p>
      <label className="flex flex-col gap-1 text-sm text-slate-700">
        订单 ID
        <input
          value={orderIdInput}
          onChange={(event) => setOrderIdInput(event.target.value)}
          className="h-10 rounded-md border border-slate-300 px-3 outline-none focus:border-slate-500"
          placeholder="请输入订单ID"
        />
      </label>
      <div className="flex gap-3">
        <Button variant="outline" onClick={() => startPayment("alipay")}>
          发起支付宝支付
        </Button>
        <Button variant="outline" onClick={() => startPayment("wechat")}>
          发起微信支付
        </Button>
      </div>
      <div className="flex flex-wrap gap-3">
        <Button variant="outline" onClick={() => navigate("/players")}>陪玩师列表</Button>
        <Button variant="outline" onClick={() => navigate("/orders")}>我的订单</Button>
        <Button variant="outline" onClick={() => navigate("/chat")}>聊天会话</Button>
        <Button variant="outline" onClick={() => navigate("/wallet")}>我的钱包</Button>
        <Button variant="outline" onClick={() => navigate("/reviews")}>我的评价</Button>
        <Button variant="outline" onClick={() => navigate("/profile")}>个人中心</Button>
      </div>
      <p className="m-0 text-sm text-slate-700">状态：{status}</p>
    </section>
  )
}
