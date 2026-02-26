export interface SessionUser {
  id?: number
  role?: string
  username?: string
  avatar?: string
}

export function getToken(): string | null {
  return localStorage.getItem("token")
}

export function isLoggedIn(): boolean {
  return Boolean(getToken())
}

export function getSessionUser(): SessionUser | null {
  const raw = localStorage.getItem("user_info")
  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw) as SessionUser
  } catch {
    return null
  }
}

export function clearSession(): void {
  localStorage.removeItem("token")
  localStorage.removeItem("user_role")
  localStorage.removeItem("user_info")
}
