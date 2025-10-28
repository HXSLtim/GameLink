import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from 'contexts/AuthContext';

/**
 * 路由守卫组件
 * 保护需要登录才能访问的路由
 */
export const ProtectedRoute = () => {
  const { user, loading } = useAuth();

  // 正在加载认证状态
  if (loading) {
    return (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100vh',
          fontSize: '18px',
          fontWeight: 600,
        }}
      >
        Loading...
      </div>
    );
  }

  // 未登录，重定向到登录页
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // 已登录，渲染子路由
  return <Outlet />;
};
