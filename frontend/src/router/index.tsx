import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Login } from 'pages/Login';
import { Register } from 'pages/Register';
import { Dashboard } from 'pages/Dashboard';
import { OrderList, OrderDetail } from 'pages/Orders';
import { GameList, GameDetail } from 'pages/Games';
import { PlayerList } from 'pages/Players';
import { UserList, UserDetail } from 'pages/Users';
import { PaymentList } from 'pages/Payments';
import { ReviewList } from 'pages/Reviews';
import { ReportDashboard } from 'pages/Reports';
import { PermissionList } from 'pages/Permissions';
import { SettingsDashboard } from 'pages/Settings';
import { ProtectedRoute } from './ProtectedRoute';
import { MainLayout } from './layouts/MainLayout';

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          {
            index: true,
            element: <Navigate to="/dashboard" replace />,
          },
          {
            path: 'dashboard',
            element: <Dashboard />,
          },
          {
            path: 'orders',
            element: <OrderList />,
          },
          {
            path: 'orders/:id',
            element: <OrderDetail />,
          },
          {
            path: 'games',
            element: <GameList />,
          },
          {
            path: 'games/:id',
            element: <GameDetail />,
          },
          {
            path: 'players',
            element: <PlayerList />,
          },
          {
            path: 'users',
            element: <UserList />,
          },
          {
            path: 'users/:id',
            element: <UserDetail />,
          },
          {
            path: 'payments',
            element: <PaymentList />,
          },
          {
            path: 'reviews',
            element: <ReviewList />,
          },
          {
            path: 'reports',
            element: <ReportDashboard />,
          },
          {
            path: 'permissions',
            element: <PermissionList />,
          },
          {
            path: 'settings',
            element: <SettingsDashboard />,
          },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/dashboard" replace />,
  },
]);
