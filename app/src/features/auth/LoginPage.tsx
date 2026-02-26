import type { FormEvent, ReactElement } from "react"
import { useEffect, useState } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import { AlertCircle, ShieldCheck } from "lucide-react"
import { Button } from "@/components/ui/button"
import { login } from "@/services/auth"
import { isLoggedIn } from "@/services/session"

export function LoginPage(): ReactElement {
  const navigate = useNavigate()
  const location = useLocation()
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isLoggedIn()) {
      navigate("/dashboard", { replace: true })
    }
  }, [navigate])

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login({ username: username.trim(), password: password.trim() })
      const from = (location.state as { from?: string } | null)?.from
      navigate(from ?? "/dashboard", { replace: true })
    } catch (err) {
      const message = err instanceof Error ? err.message : "登录失败，请稍后重试"
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="mt-8 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur sm:p-8">
      <div className="mb-6 flex items-center gap-3">
        <div className="rounded-xl bg-slate-900 p-2 text-white">
          <ShieldCheck className="h-5 w-5" />
        </div>
        <div>
          <h1 className="m-0 text-xl font-semibold text-slate-900">GameLink 用户端</h1>
          <p className="m-0 text-sm text-slate-500">React + shadcn 迁移入口</p>
        </div>
      </div>

      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="flex flex-col gap-1 text-sm text-slate-700">
          账号
          <input
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            placeholder="请输入手机号或邮箱"
            className="h-10 rounded-md border border-slate-300 px-3 outline-none focus:border-slate-500"
            required
          />
        </label>

        <label className="flex flex-col gap-1 text-sm text-slate-700">
          密码
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="请输入密码"
            className="h-10 rounded-md border border-slate-300 px-3 outline-none focus:border-slate-500"
            required
          />
        </label>

        {error ? (
          <div className="flex items-center gap-2 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
            <AlertCircle className="h-4 w-4" />
            {error}
          </div>
        ) : null}

        <Button className="w-full" disabled={loading}>
          {loading ? "登录中..." : "登录"}
        </Button>
      </form>
    </section>
  )
}
