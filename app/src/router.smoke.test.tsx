/* @vitest-environment jsdom */

import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { act } from "react"
import { createRoot, type Root } from "react-dom/client"
import { createMemoryRouter, RouterProvider } from "react-router-dom"
import { appRoutes } from "@/router"
import { listOrders } from "@/services/orders"

vi.mock("@/services/orders", () => ({
  listOrders: vi.fn(),
}))

;(globalThis as typeof globalThis & { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true

interface RenderResult {
  container: HTMLDivElement
  root: Root
}

async function flush(): Promise<void> {
  await act(async () => {
    await Promise.resolve()
    await new Promise((resolve) => setTimeout(resolve, 0))
  })
}

async function renderAt(path: string): Promise<RenderResult> {
  const container = document.createElement("div")
  document.body.appendChild(container)

  const router = createMemoryRouter(appRoutes, {
    initialEntries: [path],
  })

  const root = createRoot(container)
  await act(async () => {
    root.render(<RouterProvider router={router} />)
  })
  await flush()

  return { container, root }
}

describe("router smoke flow", () => {
  beforeEach(() => {
    localStorage.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    document.body.innerHTML = ""
  })

  it("redirects unauthenticated users from /orders to login", async () => {
    vi.mocked(listOrders).mockResolvedValue([])

    const { container, root } = await renderAt("/orders")
    expect(container.textContent).toContain("GameLink 用户端")
    expect(container.textContent).toContain("登录")

    act(() => {
      root.unmount()
    })
  })

  it("renders orders page for authenticated users", async () => {
    localStorage.setItem("token", "test-token")
    vi.mocked(listOrders).mockResolvedValue([])

    const { container, root } = await renderAt("/orders")
    expect(container.textContent).toContain("我的订单")
    expect(container.textContent).toContain("暂无订单")
    expect(listOrders).toHaveBeenCalled()

    act(() => {
      root.unmount()
    })
  })

  it("shows invalid payment id message on malformed payment result route", async () => {
    localStorage.setItem("token", "test-token")

    const { container, root } = await renderAt("/payment/result/not-a-number")
    expect(container.textContent).toContain("支付结果")
    expect(container.textContent).toContain("无效的支付单号")

    act(() => {
      root.unmount()
    })
  })
})
