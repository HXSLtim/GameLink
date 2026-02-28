import { useEffect, useMemo, useState, type ReactElement } from "react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { readNumber, readString, type UnknownRecord } from "@/services/http/envelope"
import { getMyProfile } from "@/services/profile"
import { clearSession, getSessionUser } from "@/services/session"

export function ProfilePage(): ReactElement {
  const navigate = useNavigate()
  const [profile, setProfile] = useState<UnknownRecord | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadProfile = async () => {
      setLoading(true)
      setError(null)
      try {
        setProfile(await getMyProfile())
      } catch (fetchError) {
        const message = fetchError instanceof Error ? fetchError.message : "获取个人资料失败"
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    void loadProfile()
  }, [])

  const fallbackUser = useMemo(() => getSessionUser(), [])

  return (
    <section className="mt-8 space-y-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <h1 className="m-0 text-xl font-semibold text-slate-900">个人中心</h1>
      {loading ? <p className="m-0 text-sm text-slate-600">加载中...</p> : null}
      {error ? <p className="m-0 text-sm text-red-700">{error}</p> : null}

      <div className="rounded-lg border border-slate-200 p-3 text-sm">
        <p className="m-0 text-slate-700">用户ID: {profile ? readNumber(profile, "id") ?? "-" : fallbackUser?.id ?? "-"}</p>
        <p className="m-0 text-slate-700">昵称: {profile ? readString(profile, "name") ?? "-" : fallbackUser?.username ?? "-"}</p>
        <p className="m-0 text-slate-700">角色: {profile ? readString(profile, "status") ?? "-" : fallbackUser?.role ?? "-"}</p>
      </div>

      <div className="flex gap-3">
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          返回工作台
        </Button>
        <Button
          variant="outline"
          onClick={() => {
            clearSession()
            navigate("/login", { replace: true })
          }}
        >
          退出登录
        </Button>
      </div>
    </section>
  )
}
