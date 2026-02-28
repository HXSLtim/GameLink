import { useEffect, useState, type ReactElement } from "react"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { listOrders } from "@/services/orders"

function resolveOrderNo(order: UnknownRecord): string {
  return readString(order, "orderNo") ?? readString(order, "order_no") ?? `订单#${readNumber(order, "id") ?? "-"}`
}

export function OrdersPage(): ReactElement {
  const [orders, setOrders] = useState<UnknownRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadOrders = async () => {
      setLoading(true)
      setError(null)
      try {
        setOrders(await listOrders())
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取订单失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadOrders()
  }, [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="flex items-center justify-between">
        <h1 className="m-0 text-xl font-semibold text-slate-900">我的订单</h1>
        <Button variant="outline" onClick={() => void listOrders().then(setOrders).catch(() => setError("刷新失败"))}>
          刷新
        </Button>
      </div>
      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}
      {!loading && orders.length === 0 ? <p className="m-0 text-sm text-slate-600">暂无订单</p> : null}
      <ul className="m-0 list-none space-y-2 p-0">
        {orders.map((order) => {
          const status = readString(order, "status") ?? "unknown"
          const amountCents = readNumber(order, "totalPriceCents") ?? readNumber(order, "amountCents")
          return (
            <li key={resolveOrderNo(order)} className="rounded-lg border border-slate-200 p-3 text-sm">
              <p className="m-0 font-medium text-slate-900">{resolveOrderNo(order)}</p>
              <p className="m-0 text-slate-600">状态: {status} · 金额(分): {amountCents ?? "-"}</p>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
