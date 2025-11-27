import { useEffect } from 'react';
import type { ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import type { Role } from './types';

interface RouteGuardProps {
    children: ReactNode;
    roles?: Role[];
    requiresAuth?: boolean;
}

const RouteGuard = ({ children, roles, requiresAuth }: RouteGuardProps) => {
    const navigate = useNavigate();
    const location = useLocation();

    // In a real app, this would come from a context or store
    const rawRole = localStorage.getItem('user_role');
    const userRole = rawRole ? (rawRole.toUpperCase() as Role) : null;
    const isAuthenticated = !!userRole;

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

    if (requiresAuth && !isAuthenticated) {
        return null; // or loading spinner
    }

    if (roles && userRole && !roles.includes(userRole)) {
        return null; // or unauthorized page
    }

    return <>{children}</>;
};

export default RouteGuard;
