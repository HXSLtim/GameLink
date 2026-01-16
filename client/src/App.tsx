import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "@/components/theme-provider";
import DesktopLayout from "@/layouts/DesktopLayout";
import { AuthProvider } from "@/components/auth/auth-provider";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { OfflineBanner } from "@/components/offline-banner";
import { PageLoading } from "@/components/loading-spinner";
import { InstallPrompt } from "@/components/pwa/install-prompt";
import { Toaster } from "sonner";

// HomePage is loaded synchronously for faster first contentful paint (FCP)
import HomePage from "@/pages/home/HomePage";

// Lazy load other pages to split code
const LoginPage = lazy(() => import("@/pages/auth/LoginPage"));
const RegisterPage = lazy(() => import("@/pages/auth/RegisterPage"));
const ForbiddenPage = lazy(() => import("@/pages/error/ForbiddenPage"));
const PageShowcase = lazy(() => import("@/components/page-showcase"));
const PlayerListPage = lazy(() => import("@/pages/player/PlayerListPage"));
const PlayerDetailPage = lazy(() => import("@/pages/player/PlayerDetailPage"));
const ChatListPage = lazy(() => import("@/pages/chat/ChatListPage"));
const ChatRoomPage = lazy(() => import("@/pages/chat/ChatRoomPage"));
const ProfilePage = lazy(() => import("@/pages/profile/ProfilePage"));
const OrderListPage = lazy(() => import("@/pages/order/OrderListPage"));
const OrderDetailPage = lazy(() => import("@/pages/order/OrderDetailPage"));
const ReviewOrderPage = lazy(() => import("@/pages/order/ReviewOrderPage"));
const CreateOrderPage = lazy(() => import("@/pages/order/CreateOrderPage"));
const PaymentPage = lazy(() => import("@/pages/payment/PaymentPage"));
const PlayerDashboardPage = lazy(() => import("@/pages/player/dashboard/PlayerDashboardPage"));
const PlayerOrderListPage = lazy(() => import("@/pages/player/order/PlayerOrderListPage"));
const EarningsPage = lazy(() => import("@/pages/player/earnings/EarningsPage"));
const BecomePlayerPage = lazy(() => import("@/pages/player/apply/BecomePlayerPage"));
const EditPlayerProfilePage = lazy(() => import("@/pages/player/profile/EditPlayerProfilePage"));
const RechargePage = lazy(() => import("@/pages/wallet/RechargePage"));
const WalletPage = lazy(() => import("@/pages/wallet/WalletPage"));
const VipPage = lazy(() => import("@/pages/vip/VipPage"));
const EditProfilePage = lazy(() => import("@/pages/settings/EditProfilePage"));
const ChangePasswordPage = lazy(() => import("@/pages/settings/ChangePasswordPage"));
const NotificationPage = lazy(() => import("@/pages/notification/NotificationPage"));
const CouponCenterPage = lazy(() => import("@/pages/coupon/CouponCenterPage"));
const FavoritesPage = lazy(() => import("@/pages/profile/FavoritesPage"));
const ReferralPage = lazy(() => import("@/pages/referral/ReferralPage"));
const ActivityListPage = lazy(() => import("@/pages/activity/ActivityListPage"));
const ForgotPasswordPage = lazy(() => import("@/pages/auth/ForgotPasswordPage"));
const NotFoundPage = lazy(() => import("@/pages/error/NotFoundPage"));
const TermsPage = lazy(() => import("@/pages/legal/TermsPage"));
const PrivacyPage = lazy(() => import("@/pages/legal/PrivacyPage"));
const RealNamePage = lazy(() => import("@/pages/player/verification/RealNamePage"));
const SkillAuthPage = lazy(() => import("@/pages/player/verification/SkillAuthPage"));
const TeamPage = lazy(() => import("@/pages/player/team/TeamPage"));

// Room pages
const RoomListPage = lazy(() => import("@/pages/room/RoomListPage"));
const RoomDetailPage = lazy(() => import("@/pages/room/RoomDetailPage"));
const CreateRoomPage = lazy(() => import("@/pages/room/CreateRoomPage"));

// LFG pages
const LFGPage = lazy(() => import("@/pages/lfg/LFGPage"));
const CreateLFGPage = lazy(() => import("@/pages/lfg/CreateLFGPage"));

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <OfflineBanner />
          <InstallPrompt />
          <Suspense fallback={<PageLoading />}>
            <Routes>
              {/* Public Routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} /> {/* Added route */}
              <Route path="/forgot-password" element={<ForgotPasswordPage />} />
              <Route path="/terms" element={<TermsPage />} />
              <Route path="/privacy" element={<PrivacyPage />} />
              <Route path="/403" element={<ForbiddenPage />} />

              {/* Protected Layout Routes */}
              <Route element={<ProtectedRoute />}>
                <Route element={<DesktopLayout />}>
                  <Route path="/" element={<HomePage />} />
                  <Route path="/players" element={<PlayerListPage />} />
                  <Route path="/player/dashboard" element={<PlayerDashboardPage />} />
                  <Route path="/player/orders" element={<PlayerOrderListPage />} />
                  <Route path="/player/earnings" element={<EarningsPage />} />
                  <Route path="/player/apply" element={<BecomePlayerPage />} /> {/* Added route */}
                  <Route path="/player/verification/realname" element={<RealNamePage />} /> {/* Added route */}
                  <Route path="/player/verification/skills" element={<SkillAuthPage />} /> {/* Added route */}
                  <Route path="/player/team" element={<TeamPage />} /> {/* Added route */}
                  <Route path="/player/profile/edit" element={<EditPlayerProfilePage />} />
                  <Route path="/orders" element={<OrderListPage />} />
                  <Route path="/orders/create" element={<CreateOrderPage />} />
                  <Route path="/orders/:id/review" element={<ReviewOrderPage />} /> {/* Added route */}
                  <Route path="/payment/:orderId" element={<PaymentPage />} /> {/* Added route */}
                  <Route path="/orders/:id" element={<OrderDetailPage />} />
                  <Route path="/chat" element={<ChatListPage />} />
                  <Route path="/chat/:id" element={<ChatRoomPage />} />
                  <Route path="/players/:id" element={<PlayerDetailPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="/favorites" element={<FavoritesPage />} /> {/* Added route */}
                  <Route path="/wallet" element={<WalletPage />} />
                  <Route path="/wallet/recharge" element={<RechargePage />} />
                  <Route path="/vip" element={<VipPage />} />
                  <Route path="/settings/profile" element={<EditProfilePage />} />
                  <Route path="/settings/password" element={<ChangePasswordPage />} />
                  <Route path="/notifications" element={<NotificationPage />} />
                  <Route path="/coupons" element={<CouponCenterPage />} /> {/* Added route */}
                  <Route path="/referral" element={<ReferralPage />} /> {/* Added route */}
                  <Route path="/activities" element={<ActivityListPage />} /> {/* Added route */}

                  {/* Room Routes */}
                  <Route path="/rooms" element={<RoomListPage />} />
                  <Route path="/rooms/create" element={<CreateRoomPage />} />
                  <Route path="/rooms/:id" element={<RoomDetailPage />} />

                  {/* LFG Routes */}
                  <Route path="/lfg" element={<LFGPage />} />
                  <Route path="/lfg/create" element={<CreateLFGPage />} />

                  <Route path="/page-structure" element={<PageShowcase />} />
                </Route>
              </Route>

              {/* Fallback */}
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </Suspense>
          <Toaster />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
