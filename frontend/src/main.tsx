import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { App } from './App';
import { Dashboard } from './pages/Dashboard';
import { Login } from './pages/Login';
import { Users } from './pages/Users';
import { Orders } from './pages/Orders';
import { Permissions } from './pages/Permissions';
import { RequireAuth } from './components/RequireAuth';
import { ErrorBoundary } from './components/ErrorBoundary';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
// 使用 less 入口以支持主题变量覆盖
import '@arco-design/web-react/dist/css/index.less';

const container = document.getElementById('root');

if (!container) {
  throw new Error('Root container not found');
}

const root = createRoot(container);

const router = createBrowserRouter(
  [
    {
      path: '/',
      element: (
        <RequireAuth>
          <App />
        </RequireAuth>
      ),
      handle: { crumb: '控制台' },
      children: [
        { index: true, element: <Dashboard />, handle: { crumb: '总览' } },
        { path: 'settings', element: <div>设置（占位）</div>, handle: { crumb: '设置' } },
        { path: 'users', element: <Users />, handle: { crumb: '用户' } },
        { path: 'orders', element: <Orders />, handle: { crumb: '订单' } },
        { path: 'permissions', element: <Permissions />, handle: { crumb: '权限' } },
      ],
    },
    { path: '/login', element: <Login /> },
  ],
  { future: { v7_relativeSplatPath: true } },
);

root.render(
  <StrictMode>
    <ErrorBoundary>
      <ThemeProvider>
        <AuthProvider>
          <RouterProvider router={router} future={{ v7_startTransition: true }} />
        </AuthProvider>
      </ThemeProvider>
    </ErrorBoundary>
  </StrictMode>,
);
