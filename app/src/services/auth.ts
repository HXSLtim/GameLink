import { httpClient } from "@/services/http/client"

interface LoginPayload {
  username: string
  password: string
}

interface LoginResponse {
  success?: boolean
  data?: {
    token: string
    user: {
      id: number
      role: string
      name?: string
      username?: string
      avatar?: string
      avatarUrl?: string
    }
  }
  message?: string
}

export async function login(payload: LoginPayload): Promise<void> {
  const response = await httpClient.post<LoginResponse>("/auth/login", payload)
  const body = response.data

  if (!body.success || !body.data) {
    throw new Error(body.message ?? "登录失败")
  }

  localStorage.setItem("token", body.data.token)
  localStorage.setItem("user_role", body.data.user.role)
  localStorage.setItem(
    "user_info",
    JSON.stringify({
      ...body.data.user,
      username: body.data.user.name ?? body.data.user.username,
      avatar: body.data.user.avatarUrl ?? body.data.user.avatar,
    }),
  )
}
