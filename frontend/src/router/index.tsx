import { createBrowserRouter, Navigate } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import { RouteLoading } from './LazyRoutes';
import { appRoutes } from './routes';
import { ProtectedRoute } from './ProtectedRoute';
import { FEATURE_FLAGS } from '@/config';

// 懒加载公共页面组件
// TODO: 创建这些页面组件
// const Login = lazy(() => import('@/apps/user/pages/Login').then((m) => ({ default: m.Login })));
// const Register = lazy(() => import('@/apps/user/pages/Register').then((m) => ({ default: m.Register })));
// const ComponentsDemo = lazy(() => import('@/pages/ComponentsDemo').then((m) => ({ default: m.ComponentsDemo })));
// const CacheDemo = lazy(() => import('@/pages/CacheDemo').then((m) => ({ default: m.CacheDemo })));
// const CachePageA = lazy(() => import('@/pages/CacheDemo').then((m) => ({ default: m.CachePageA })));
// const CachePageB = lazy(() => import('@/pages/CacheDemo').then((m) => ({ default: m.CachePageB })));

// 临时：使用空组件
const PlaceholderPage = () => <div>页面开发中...</div>;
const Login = PlaceholderPage;
const Register = PlaceholderPage;
const ComponentsDemo = PlaceholderPage;
const CacheDemo = PlaceholderPage;
const CachePageA = PlaceholderPage;
const CachePageB = PlaceholderPage;

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

  // 受保护的路由（需要认证）- 包含三端应用路由
  {
    path: '/',
    element: <ProtectedRoute />,
    children: appRoutes,
  },

  // 默认重定向
  { path: '*', element: <Navigate to='/user/discover' replace /> },
]);