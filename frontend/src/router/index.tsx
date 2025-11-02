import { createBrowserRouter, Navigate } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import { RouteLoading } from '../routes/LazyRoutes';
// 懒加载页面组件（保持命名与原路由一致）
const Login = lazy(() => import('pages/Login').then((m) => ({ default: m.Login })));
const Register = lazy(() => import('pages/Register').then((m) => ({ default: m.Register })));
const Dashboard = lazy(() => import('pages/Dashboard').then((m) => ({ default: m.Dashboard })));
const OrderList = lazy(() => import('pages/Orders').then((m) => ({ default: m.OrderList })));
const OrderDetail = lazy(() => import('pages/Orders').then((m) => ({ default: m.OrderDetail })));
const GameList = lazy(() => import('pages/Games').then((m) => ({ default: m.GameList })));
const GameDetail = lazy(() => import('pages/Games').then((m) => ({ default: m.GameDetail })));
const PlayerList = lazy(() => import('pages/Players').then((m) => ({ default: m.PlayerList })));
const UserList = lazy(() => import('pages/Users').then((m) => ({ default: m.UserList })));
const UserDetail = lazy(() => import('pages/Users').then((m) => ({ default: m.UserDetail })));
const PaymentList = lazy(() => import('pages/Payments').then((m) => ({ default: m.PaymentList })));
const PaymentDetailPage = lazy(() => import('pages/Payments').then((m) => ({ default: m.PaymentDetailPage })));
const ReviewList = lazy(() => import('pages/Reviews').then((m) => ({ default: m.ReviewList })));
const ReportDashboard = lazy(() => import('pages/Reports').then((m) => ({ default: m.ReportDashboard })));
const PermissionList = lazy(() => import('pages/Permissions').then((m) => ({ default: m.PermissionList })));
const SettingsDashboard = lazy(() => import('pages/Settings').then((m) => ({ default: m.SettingsDashboard })));
const ComponentsDemo = lazy(() => import('pages/ComponentsDemo').then((m) => ({ default: m.ComponentsDemo })));
const CacheDemo = lazy(() => import('pages/CacheDemo').then((m) => ({ default: m.CacheDemo })));
const CachePageA = lazy(() => import('pages/CacheDemo').then((m) => ({ default: m.CachePageA })));
const CachePageB = lazy(() => import('pages/CacheDemo').then((m) => ({ default: m.CachePageB })));
import { ProtectedRoute } from './ProtectedRoute';
import { MainLayout } from './layouts/MainLayout';
import { FEATURE_FLAGS } from '../config';

export const router = createBrowserRouter([
  // 公开路由（无需认证）
  { path: '/login', element: <Suspense fallback={<RouteLoading />}> <Login /> </Suspense> },
  { path: '/register', element: <Suspense fallback={<RouteLoading />}> <Register /> </Suspense> },

  // 演示路由（无需认证）
  ...(FEATURE_FLAGS.showcase.enableShowcaseRoute
    ? [{ path: '/showcase', element: <Suspense fallback={<RouteLoading />}> <ComponentsDemo /> </Suspense> }]
    : []),

  ...(FEATURE_FLAGS.cacheDemo.enableCacheDemoRoute
    ? [
        {
          path: '/cache-demo',
          element: (
            <Suspense fallback={<RouteLoading />}>
              <CacheDemo />
            </Suspense>
          ),
          children: [
            { index: true, element: <Navigate to='/cache-demo/a' replace /> },
            { path: 'a', element: <Suspense fallback={<RouteLoading />}> <CachePageA /> </Suspense> },
            { path: 'b', element: <Suspense fallback={<RouteLoading />}> <CachePageB /> </Suspense> },
          ],
        },
      ]
    : []),

  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          { index: true, element: <Navigate to='/dashboard' replace /> },
          { path: 'dashboard', element: <Suspense fallback={<RouteLoading />}> <Dashboard /> </Suspense> },
          { path: 'orders', element: <Suspense fallback={<RouteLoading />}> <OrderList /> </Suspense> },
          { path: 'orders/:id', element: <Suspense fallback={<RouteLoading />}> <OrderDetail /> </Suspense> },
          { path: 'games', element: <Suspense fallback={<RouteLoading />}> <GameList /> </Suspense> },
          { path: 'games/:id', element: <Suspense fallback={<RouteLoading />}> <GameDetail /> </Suspense> },
          { path: 'players', element: <Suspense fallback={<RouteLoading />}> <PlayerList /> </Suspense> },
          { path: 'users', element: <Suspense fallback={<RouteLoading />}> <UserList /> </Suspense> },
          { path: 'users/:id', element: <Suspense fallback={<RouteLoading />}> <UserDetail /> </Suspense> },
          { path: 'payments', element: <Suspense fallback={<RouteLoading />}> <PaymentList /> </Suspense> },
          { path: 'payments/:id', element: <Suspense fallback={<RouteLoading />}> <PaymentDetailPage /> </Suspense> },
          { path: 'reviews', element: <Suspense fallback={<RouteLoading />}> <ReviewList /> </Suspense> },
          { path: 'reports', element: <Suspense fallback={<RouteLoading />}> <ReportDashboard /> </Suspense> },
          { path: 'permissions', element: <Suspense fallback={<RouteLoading />}> <PermissionList /> </Suspense> },
          { path: 'settings', element: <Suspense fallback={<RouteLoading />}> <SettingsDashboard /> </Suspense> },
          ...(FEATURE_FLAGS.showcase.enableComponentsRoute
            ? [{ path: 'components', element: <Suspense fallback={<RouteLoading />}> <ComponentsDemo /> </Suspense> }]
            : []),
        ],
      },
    ],
  },

  { path: '*', element: <Navigate to='/dashboard' replace /> },
]);
