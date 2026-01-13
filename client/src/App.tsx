import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ThemeProvider } from "@/components/theme-provider";
import DesktopLayout from "@/layouts/DesktopLayout";
import { AuthProvider } from "@/components/auth/auth-provider";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { Toaster } from "sonner";

// Lazy load pages to split code
const LoginPage = lazy(() => import("@/pages/auth/login-page"));
const ForbiddenPage = lazy(() => import("@/pages/error/403-page"));
const PageShowcase = lazy(() => import("@/components/page-showcase"));
const PlayerListPage = lazy(() => import("@/pages/player/player-list-page"));
const PlayerDetailPage = lazy(() => import("@/pages/player/player-detail-page"));
const ChatListPage = lazy(() => import("@/pages/chat/chat-list-page"));
const ChatRoomPage = lazy(() => import("@/pages/chat/chat-room-page"));
const ProfilePage = lazy(() => import("@/pages/profile/profile-page"));
const OrderListPage = lazy(() => import("@/pages/order/order-list-page"));
const OrderDetailPage = lazy(() => import("@/pages/order/order-detail-page"));
const WalletPage = lazy(() => import("@/pages/wallet/wallet-page"));
const VipPage = lazy(() => import("@/pages/vip/vip-page"));
const EditProfilePage = lazy(() => import("@/pages/settings/edit-profile-page"));
const ChangePasswordPage = lazy(() => import("@/pages/settings/change-password-page"));
const NotificationPage = lazy(() => import("@/pages/notification/notification-page"));

const HomePage = lazy(() => import("@/pages/home/home-page"));

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <Suspense fallback={<div className="flex h-screen items-center justify-center">Loading...</div>}>
            <Routes>
              {/* Public Routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/403" element={<ForbiddenPage />} />

              {/* Protected Layout Routes */}
              <Route element={<ProtectedRoute />}>
                <Route element={<DesktopLayout />}>
                  <Route path="/" element={<HomePage />} />
                  <Route path="/players" element={<PlayerListPage />} />
                  <Route path="/orders" element={<OrderListPage />} />
                  <Route path="/orders/:id" element={<OrderDetailPage />} />
                  <Route path="/chat" element={<ChatListPage />} />
                  <Route path="/chat/:id" element={<ChatRoomPage />} />
                  <Route path="/players/:id" element={<PlayerDetailPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="/wallet" element={<WalletPage />} />
                  <Route path="/vip" element={<VipPage />} />
                  <Route path="/settings/profile" element={<EditProfilePage />} />
                  <Route path="/settings/password" element={<ChangePasswordPage />} />
                  <Route path="/notifications" element={<NotificationPage />} />
                  <Route path="/page-structure" element={<PageShowcase />} />
                </Route>
              </Route>

              {/* Fallback */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </Suspense>
          <Toaster />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
