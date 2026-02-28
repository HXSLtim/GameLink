import { useEffect, useState, type ReactElement } from "react"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { listPlayers } from "@/services/players"

function resolvePlayerName(player: UnknownRecord): string {
  return readString(player, "name") ?? readString(player, "nickname") ?? "未命名陪玩师"
}

export function PlayersPage(): ReactElement {
  const [players, setPlayers] = useState<UnknownRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadPlayers = async () => {
      setLoading(true)
      setError(null)
      try {
        setPlayers(await listPlayers())
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取陪玩师列表失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadPlayers()
  }, [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="flex items-center justify-between">
        <h1 className="m-0 text-xl font-semibold text-slate-900">陪玩师列表</h1>
        <Button variant="outline" onClick={() => void listPlayers().then(setPlayers).catch(() => setError("刷新失败"))}>
          刷新
        </Button>
      </div>
      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}
      {!loading && players.length === 0 ? <p className="m-0 text-sm text-slate-600">暂无数据</p> : null}
      <ul className="m-0 list-none space-y-2 p-0">
        {players.map((player) => {
          const playerId = readNumber(player, "id")
          const rating = readNumber(player, "rating")
          const priceCents = readNumber(player, "priceCents") ?? readNumber(player, "price")
          return (
            <li key={playerId ?? resolvePlayerName(player)} className="rounded-lg border border-slate-200 p-3 text-sm">
              <p className="m-0 font-medium text-slate-900">{resolvePlayerName(player)}</p>
              <p className="m-0 text-slate-600">ID: {playerId ?? "-"} · 评分: {rating ?? "-"} · 价格(分): {priceCents ?? "-"}</p>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
