/**
 * 路由权限守卫组件
 * 实现路由级别的权限控制
 * Requirements: 8.3 - 无权限重定向到 403 页面
 */
import { useEffect, useRef, useMemo } from 'react';
import type { ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Spin } from 'antd';
import type { Role } from './types';
import { usePermission } from '@/hooks/usePermission';
import { useUserInfo, useIsAuthenticated, useIsHydrated } from '@/stores/modules/authStore';

interface RouteGuardProps {
    children: ReactNode;
    roles?: Role[];
    requiresAuth?: boolean;
    permission?: string;
}

/**
 * 路由守卫组件
 * - 检查用户是否已登录
 * - 检查用户角色是否有权限访问
 * - 检查用户是否有指定的权限码
 * - 无权限时重定向到 403 页面
 */
const RouteGuard = ({ children, roles, requiresAuth, permission }: RouteGuardProps) => {
    const navigate = useNavigate();
    const location = useLocation();
    const hasRedirected = useRef(false);

    // Super Dev 最佳实践: 使用选择器精确订阅
    const isAuthenticated = useIsAuthenticated();
    const userInfo = useUserInfo();
    const isHydrated = useIsHydrated();

    const userRole = userInfo?.role?.toUpperCase() as Role | null;

    // Always call usePermission with a stable value - never conditionally
    // Pass the permission as-is, usePermission handles empty strings internally
    const { hasPermission, loading: permissionLoading } = usePermission(permission || '');

    // Compute whether we need permission check after hooks
    const needsPermissionCheck = useMemo(
        () => !!permission && permission.length > 0,
        [permission]
    );

    // Handle authentication redirect - must be called unconditionally
    useEffect(() => {
        if (!isHydrated) return; // Skip if not hydrated yet

        if (requiresAuth && !isAuthenticated) {
            // 根据当前路径决定重定向到哪个登录页
            const loginPath = location.pathname.startsWith('/admin') ? '/admin/login' : '/login';
            navigate(loginPath, { state: { from: location } });
            return;
        }

        if (roles && userRole && !roles.includes(userRole)) {
            // Redirect to appropriate home based on role or 403 page
            if (userRole === 'USER') navigate('/');
            else if (userRole === 'PLAYER') navigate('/player');
            else if (userRole === 'ADMIN') navigate('/admin');
            else navigate('/'); // Fallback
        }
    }, [isHydrated, isAuthenticated, userRole, roles, requiresAuth, navigate, location]);

    // Check permission after loading - redirect to 403 if no permission
    useEffect(() => {
        if (!isHydrated) return; // Skip if not hydrated yet

        if (needsPermissionCheck && !permissionLoading && !hasPermission && !hasRedirected.current) {
            hasRedirected.current = true;
            // 重定向到 403 页面，并传递原始路径信息
            navigate('/403', {
                replace: true,
                state: {
                    from: location.pathname,
                    requiredPermission: permission
                }
            });
        }
    }, [isHydrated, needsPermissionCheck, permissionLoading, hasPermission, navigate, location.pathname, permission]);

    // Reset redirect flag when location changes
    useEffect(() => {
        hasRedirected.current = false;
    }, [location.pathname]);

    // 等待 Zustand persist 水合完成
    if (!isHydrated) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                minHeight: '200px'
            }}>
                <Spin size="large" tip="加载中..." />
            </div>
        );
    }

    if (requiresAuth && !isAuthenticated) {
        return null; // Will redirect to login
    }

    if (roles && userRole && !roles.includes(userRole)) {
        return null; // Will redirect based on role
    }

    // Show loading spinner while checking permissions
    if (needsPermissionCheck && permissionLoading) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                minHeight: '200px'
            }}>
                <Spin size="large" tip="验证权限中..." />
            </div>
        );
    }

    if (needsPermissionCheck && !hasPermission) {
        return null; // Will redirect to 403
    }

    return <>{children}</>;
};

export default RouteGuard;
