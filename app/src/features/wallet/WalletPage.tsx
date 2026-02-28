import { useEffect, useState, type ReactElement } from "react"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { getWalletBalance, listWalletTransactions } from "@/services/wallet"

export function WalletPage(): ReactElement {
  const [balance, setBalance] = useState<UnknownRecord | null>(null)
  const [transactions, setTransactions] = useState<UnknownRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadWallet = async () => {
      setLoading(true)
      setError(null)
      try {
        const [currentBalance, currentTransactions] = await Promise.all([getWalletBalance(), listWalletTransactions()])
        setBalance(currentBalance)
        setTransactions(currentTransactions)
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取钱包信息失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadWallet()
  }, [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="flex items-center justify-between">
        <h1 className="m-0 text-xl font-semibold text-slate-900">我的钱包</h1>
        <Button
          variant="outline"
          onClick={() => {
            void Promise.all([getWalletBalance(), listWalletTransactions()])
              .then(([currentBalance, currentTransactions]) => {
                setBalance(currentBalance)
                setTransactions(currentTransactions)
              })
              .catch(() => setError("刷新失败"))
          }}
        >
          刷新
        </Button>
      </div>

      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}

      <div className="rounded-lg border border-slate-200 p-3 text-sm">
        <p className="m-0 text-slate-700">余额(分): {balance ? readNumber(balance, "balanceCents") ?? "-" : "-"}</p>
        <p className="m-0 text-slate-700">冻结(分): {balance ? readNumber(balance, "frozenCents") ?? 0 : 0}</p>
      </div>

      <h2 className="m-0 text-base font-semibold text-slate-900">交易记录</h2>
      {!loading && transactions.length === 0 ? <p className="m-0 text-sm text-slate-600">暂无交易记录</p> : null}
      <ul className="m-0 list-none space-y-2 p-0">
        {transactions.map((transaction) => {
          const transactionId = readNumber(transaction, "id")
          const transactionType = readString(transaction, "type") ?? "-"
          const amountCents = readNumber(transaction, "amountCents")
          const status = readString(transaction, "status") ?? "-"
          return (
            <li key={transactionId ?? `${transactionType}-${amountCents}`} className="rounded-lg border border-slate-200 p-3 text-sm">
              <p className="m-0 font-medium text-slate-900">#{transactionId ?? "-"} · {transactionType}</p>
              <p className="m-0 text-slate-600">金额(分): {amountCents ?? "-"} · 状态: {status}</p>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
