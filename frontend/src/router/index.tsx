import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Login } from 'pages/Login';
import { Register } from 'pages/Register';
import { Dashboard } from 'pages/Dashboard';
import { OrderList, OrderDetail } from 'pages/Orders';
import { GameList, GameDetail } from 'pages/Games';
import { PlayerList } from 'pages/Players';
import { UserList, UserDetail } from 'pages/Users';
import { PaymentList, PaymentDetailPage } from 'pages/Payments';
import { ReviewList } from 'pages/Reviews';
import { ReportDashboard } from 'pages/Reports';
import { PermissionList } from 'pages/Permissions';
import { SettingsDashboard } from 'pages/Settings';
import { ComponentsDemo } from 'pages/ComponentsDemo';
import { CacheDemo, CachePageA, CachePageB } from 'pages/CacheDemo';
import { ProtectedRoute } from './ProtectedRoute';
import { MainLayout } from './layouts/MainLayout';
import { FEATURE_FLAGS } from '../config';

export const router = createBrowserRouter([
  // 公开路由（无需认证）
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },

  // 演示路由（无需认证）
  ...(FEATURE_FLAGS.showcase.enableShowcaseRoute
    ? [{ path: '/showcase', element: <ComponentsDemo /> }]
    : []),

  ...(FEATURE_FLAGS.cacheDemo.enableCacheDemoRoute
    ? [
        {
          path: '/cache-demo',
          element: <CacheDemo />,
          children: [
            { index: true, element: <Navigate to='/cache-demo/a' replace /> },
            { path: 'a', element: <CachePageA /> },
            { path: 'b', element: <CachePageB /> },
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
          { path: 'dashboard', element: <Dashboard /> },
          { path: 'orders', element: <OrderList /> },
          { path: 'orders/:id', element: <OrderDetail /> },
          { path: 'games', element: <GameList /> },
          { path: 'games/:id', element: <GameDetail /> },
          { path: 'players', element: <PlayerList /> },
          { path: 'users', element: <UserList /> },
          { path: 'users/:id', element: <UserDetail /> },
          { path: 'payments', element: <PaymentList /> },
          { path: 'payments/:id', element: <PaymentDetailPage /> },
          { path: 'reviews', element: <ReviewList /> },
          { path: 'reports', element: <ReportDashboard /> },
          { path: 'permissions', element: <PermissionList /> },
          { path: 'settings', element: <SettingsDashboard /> },
          ...(FEATURE_FLAGS.showcase.enableComponentsRoute
            ? [{ path: 'components', element: <ComponentsDemo /> }]
            : []),
        ],
      },
    ],
  },

  { path: '*', element: <Navigate to='/dashboard' replace /> },
]);
