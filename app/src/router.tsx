import { createBrowserRouter } from "react-router-dom"
import { App } from "@/App"
import { RouteAuthGuard } from "@/components/RouteAuthGuard"
import { LoginPage } from "@/features/auth/LoginPage"
import { DashboardPage } from "@/features/dashboard/DashboardPage"
import { PaymentCallbackPage } from "@/features/payment/PaymentCallbackPage"
import { PaymentResultPage } from "@/features/payment/PaymentResultPage"

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
