import type { ReactElement } from "react"
import { Outlet } from "react-router-dom"

export function App(): ReactElement {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md flex-col px-4 py-8 sm:max-w-lg">
      <Outlet />
    </main>
  )
}
