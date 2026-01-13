import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/stores';

interface ProtectedRouteProps {
    roles?: string[]; // Allowed roles
}

export function ProtectedRoute({ roles }: ProtectedRouteProps) {
    const { isAuthenticated, role } = useAuthStore();
    const location = useLocation();

    if (!isAuthenticated) {
        // Redirect to login page with return url
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    if (roles && !roles.includes(role)) {
        // Role based access control
        return <Navigate to="/403" replace />;
    }

    return <Outlet />;
}
