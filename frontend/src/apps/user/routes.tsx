/**
 * 用户端路由配置
 * 参考Discord/KOOK风格
 */
import React, { lazy } from 'react';
import type { RouteObject } from 'react-router-dom';

// 懒加载页面组件
const UserLayout = lazy(() => import('./layouts/UserLayout'));
const DiscoverPage = lazy(() => import('./pages/Discover'));
const PlayersPage = lazy(() => import('./pages/Players'));
const PlayerDetailPage = lazy(() => import('./pages/PlayerDetail'));
const OrdersPage = lazy(() => import('./pages/Orders'));
const OrderDetailPage = lazy(() => import('./pages/OrderDetail'));
const CreateOrderPage = lazy(() => import('./pages/CreateOrder'));
const MessagesPage = lazy(() => import('./pages/Messages'));
const ProfilePage = lazy(() => import('./pages/Profile'));

export const userRoutes: RouteObject[] = [
  {
    path: '/user',
    element: <UserLayout />,
    children: [
      {
        path: 'discover',
        element: <DiscoverPage />,
      },
      {
        path: 'players',
        element: <PlayersPage />,
      },
      {
        path: 'players/:id',
        element: <PlayerDetailPage />,
      },
      {
        path: 'orders',
        element: <OrdersPage />,
      },
      {
        path: 'orders/:id',
        element: <OrderDetailPage />,
      },
      {
        path: 'orders/create',
        element: <CreateOrderPage />,
      },
      {
        path: 'messages',
        element: <MessagesPage />,
      },
      {
        path: 'messages/:userId',
        element: <MessagesPage />,
      },
      {
        path: 'profile',
        element: <ProfilePage />,
      },
    ],
  },
];
