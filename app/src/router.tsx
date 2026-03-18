import { lazy, Suspense } from "react"
import { createBrowserRouter } from "react-router-dom"
import { App } from "@/App"
import { RouteAuthGuard } from "@/components/RouteAuthGuard"

// 懒加载页面组件
const LoginPage = lazy(() => import("@/features/auth/LoginPage").then(m => ({ default: m.LoginPage })))
const ChatPage = lazy(() => import("@/features/chat/ChatPage").then(m => ({ default: m.ChatPage })))
const DashboardPage = lazy(() => import("@/features/dashboard/DashboardPage").then(m => ({ default: m.DashboardPage })))
const OrdersPage = lazy(() => import("@/features/orders/OrdersPage").then(m => ({ default: m.OrdersPage })))
const PaymentCallbackPage = lazy(() => import("@/features/payment/PaymentCallbackPage").then(m => ({ default: m.PaymentCallbackPage })))
const PaymentResultPage = lazy(() => import("@/features/payment/PaymentResultPage").then(m => ({ default: m.PaymentResultPage })))
const PlayersPage = lazy(() => import("@/features/players/PlayersPage").then(m => ({ default: m.PlayersPage })))
const ProfilePage = lazy(() => import("@/features/profile/ProfilePage").then(m => ({ default: m.ProfilePage })))
const ReviewsPage = lazy(() => import("@/features/reviews/ReviewsPage").then(m => ({ default: m.ReviewsPage })))
const WalletPage = lazy(() => import("@/features/wallet/WalletPage").then(m => ({ default: m.WalletPage })))

// Loading 组件
function PageLoader() {
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="text-center">
        <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-slate-300 border-r-transparent"></div>
        <p className="text-sm text-slate-600">加载中...</p>
      </div>
    </div>
  )
}

// Suspense 包装器
function LazyRoute({ children }: { children: React.ReactNode }) {
  return <Suspense fallback={<PageLoader />}>{children}</Suspense>
}

export const appRoutes = [
  {
    path: "/",
    element: <App />,
    children: [
      {
        index: true,
        element: <LazyRoute><LoginPage /></LazyRoute>,
      },
      {
        path: "login",
        element: <LazyRoute><LoginPage /></LazyRoute>,
      },
      {
        path: "payment/callback",
        element: <LazyRoute><PaymentCallbackPage /></LazyRoute>,
      },
      {
        element: <RouteAuthGuard />,
        children: [
          {
            path: "dashboard",
            element: <LazyRoute><DashboardPage /></LazyRoute>,
          },
          {
            path: "players",
            element: <LazyRoute><PlayersPage /></LazyRoute>,
          },
          {
            path: "orders",
            element: <LazyRoute><OrdersPage /></LazyRoute>,
          },
          {
            path: "chat",
            element: <LazyRoute><ChatPage /></LazyRoute>,
          },
          {
            path: "wallet",
            element: <LazyRoute><WalletPage /></LazyRoute>,
          },
          {
            path: "reviews",
            element: <LazyRoute><ReviewsPage /></LazyRoute>,
          },
          {
            path: "profile",
            element: <LazyRoute><ProfilePage /></LazyRoute>,
          },
          {
            path: "payment/result/:paymentId",
            element: <LazyRoute><PaymentResultPage /></LazyRoute>,
          },
          {
            path: "payment/result",
            element: <LazyRoute><PaymentResultPage /></LazyRoute>,
          },
        ],
      },
    ],
  },
]

export const router = createBrowserRouter(appRoutes)
