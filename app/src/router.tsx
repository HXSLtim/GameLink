import { createBrowserRouter } from "react-router-dom"
import { App } from "@/App"
import { RouteAuthGuard } from "@/components/RouteAuthGuard"
import { LoginPage } from "@/features/auth/LoginPage"
import { ChatPage } from "@/features/chat/ChatPage"
import { DashboardPage } from "@/features/dashboard/DashboardPage"
import { OrdersPage } from "@/features/orders/OrdersPage"
import { PaymentCallbackPage } from "@/features/payment/PaymentCallbackPage"
import { PaymentResultPage } from "@/features/payment/PaymentResultPage"
import { PlayersPage } from "@/features/players/PlayersPage"
import { ProfilePage } from "@/features/profile/ProfilePage"
import { ReviewsPage } from "@/features/reviews/ReviewsPage"
import { WalletPage } from "@/features/wallet/WalletPage"

export const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        index: true,
        element: <LoginPage />,
      },
      {
        path: "login",
        element: <LoginPage />,
      },
      {
        path: "payment/callback",
        element: <PaymentCallbackPage />,
      },
      {
        element: <RouteAuthGuard />,
        children: [
          {
            path: "dashboard",
            element: <DashboardPage />,
          },
          {
            path: "players",
            element: <PlayersPage />,
          },
          {
            path: "orders",
            element: <OrdersPage />,
          },
          {
            path: "chat",
            element: <ChatPage />,
          },
          {
            path: "wallet",
            element: <WalletPage />,
          },
          {
            path: "reviews",
            element: <ReviewsPage />,
          },
          {
            path: "profile",
            element: <ProfilePage />,
          },
          {
            path: "payment/result/:paymentId",
            element: <PaymentResultPage />,
          },
          {
            path: "payment/result",
            element: <PaymentResultPage />,
          },
        ],
      },
    ],
  },
])
