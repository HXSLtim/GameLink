import { useEffect, useRef } from 'react';
import type { ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { message } from 'antd';
import type { Role } from './types';
import { usePermission } from '@/hooks/usePermission';

interface RouteGuardProps {
    children: ReactNode;
    roles?: Role[];
    requiresAuth?: boolean;
    permission?: string;
}

const RouteGuard = ({ children, roles, requiresAuth, permission }: RouteGuardProps) => {
    const navigate = useNavigate();
    const location = useLocation();
    const hasShownError = useRef(false);

    // In a real app, this would come from a context or store
    const rawRole = localStorage.getItem('user_role');
    const userRole = rawRole ? (rawRole.toUpperCase() as Role) : null;
    const isAuthenticated = !!userRole;

    // Check permission only when permission is defined and not empty
    const needsPermissionCheck = !!permission && permission.length > 0;
    const { hasPermission, loading: permissionLoading } = usePermission(needsPermissionCheck ? permission : '');

    useEffect(() => {
        if (requiresAuth && !isAuthenticated) {
            navigate('/login', { state: { from: location } });
            return;
        }

        if (roles && userRole && !roles.includes(userRole)) {
            // Redirect to appropriate home based on role or 403 page
            if (userRole === 'USER') navigate('/');
            else if (userRole === 'COMPANION') navigate('/companion');
            else if (userRole === 'ADMIN') navigate('/admin');
            else navigate('/'); // Fallback
        }
    }, [isAuthenticated, userRole, roles, requiresAuth, navigate, location]);

    // Check permission after loading - only if permission check is needed
    useEffect(() => {
        if (needsPermissionCheck && !permissionLoading && !hasPermission && !hasShownError.current) {
            hasShownError.current = true;
            message.error('您没有访问此页面的权限');
            navigate('/admin', { replace: true });
        }
    }, [needsPermissionCheck, permissionLoading, hasPermission, navigate]);

    if (requiresAuth && !isAuthenticated) {
        return null; // or loading spinner
    }

    if (roles && userRole && !roles.includes(userRole)) {
        return null; // or unauthorized page
    }

    // Wait for permission check if permission is specified
    if (needsPermissionCheck && permissionLoading) {
        return null; // or loading spinner
    }

    if (needsPermissionCheck && !hasPermission) {
        return null; // or unauthorized page
    }

    return <>{children}</>;
};

export default RouteGuard;
