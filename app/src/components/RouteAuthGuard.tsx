import type { ReactElement } from "react"
import { Navigate, Outlet, useLocation } from "react-router-dom"
import { isLoggedIn } from "@/services/session"

export function RouteAuthGuard(): ReactElement {
  const location = useLocation()

  if (!isLoggedIn()) {
    return <Navigate to="/login" replace state={{ from: `${location.pathname}${location.search}` }} />
  }

  return <Outlet />
}
