import { useEffect, useState, type ReactElement } from "react"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { listChatGroups } from "@/services/chat"

function resolveGroupName(group: UnknownRecord): string {
  return readString(group, "groupName") ?? readString(group, "name") ?? `会话#${readNumber(group, "id") ?? "-"}`
}

export function ChatPage(): ReactElement {
  const [groups, setGroups] = useState<UnknownRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadGroups = async () => {
      setLoading(true)
      setError(null)
      try {
        setGroups(await listChatGroups())
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取聊天会话失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadGroups()
  }, [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="flex items-center justify-between">
        <h1 className="m-0 text-xl font-semibold text-slate-900">聊天会话</h1>
        <Button variant="outline" onClick={() => void listChatGroups().then(setGroups).catch(() => setError("刷新失败"))}>
          刷新
        </Button>
      </div>
      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}
      {!loading && groups.length === 0 ? <p className="m-0 text-sm text-slate-600">暂无会话</p> : null}
      <ul className="m-0 list-none space-y-2 p-0">
        {groups.map((group) => {
          const unreadCount = readNumber(group, "unreadCount")
          const groupType = readString(group, "groupType") ?? readString(group, "type") ?? "-"
          return (
            <li key={resolveGroupName(group)} className="rounded-lg border border-slate-200 p-3 text-sm">
              <p className="m-0 font-medium text-slate-900">{resolveGroupName(group)}</p>
              <p className="m-0 text-slate-600">类型: {groupType} · 未读: {unreadCount ?? 0}</p>
            </li>
          )
        })}
      </ul>
    </section>
  )
}
