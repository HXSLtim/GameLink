/**
 * Route Guard Component for GameLink Client
 * Handles authentication, authorization, and view mode enforcement
 *
 * Features:
 * - Authentication check (redirect to /login if not authenticated)
 * - Role-based access control (user vs player routes)
 * - View mode enforcement (user/player mode switching)
 * - Zustand hydration handling
 */

import { useEffect, useRef, useMemo, type ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/stores/modules/auth-store';

interface RouteGuardProps {
    children: ReactNode;
    requiresAuth?: boolean;
    roles?: ('user' | 'player')[];
    viewMode?: 'user' | 'player';
}

/**
 * Route Guard Component
 *
 * Usage:
 * ```tsx
 * <RouteGuard requiresAuth roles={['player']} viewMode="player">
 *   <PlayerDashboard />
 * </RouteGuard>
 * ```
 */
export function RouteGuard({
    children,
    requiresAuth = false,
    roles,
    viewMode
}: RouteGuardProps) {
    const navigate = useNavigate();
    const location = useLocation();
    const hasRedirected = useRef(false);

    // Subscribe to auth state
    const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
    const role = useAuthStore((state) => state.role);
    const currentViewMode = useAuthStore((state) => state.viewMode);
    const isPlayer = useAuthStore((state) => state.isPlayer);

    // Check if hydration is complete (Zustand persist)
    const isHydrated = useAuthStore.persist?.hasHydrated() ?? true;

    // Determine if user has required role
    const hasRequiredRole = useMemo(() => {
        if (!roles || roles.length === 0) return true;

        // Map role to check against
        const userRole = role === 'player' ? 'player' : 'user';
        return roles.includes(userRole);
    }, [roles, role]);

    // Determine if view mode matches
    const hasCorrectViewMode = useMemo(() => {
        if (!viewMode) return true;
        return currentViewMode === viewMode;
    }, [viewMode, currentViewMode]);

    // Handle authentication redirect
    useEffect(() => {
        if (!isHydrated) return;

        if (requiresAuth && !isAuthenticated && !hasRedirected.current) {
            hasRedirected.current = true;
            navigate('/login', {
                state: { from: location.pathname },
                replace: true
            });
        }
    }, [isHydrated, requiresAuth, isAuthenticated, navigate, location.pathname]);

    // Handle role-based redirect
    useEffect(() => {
        if (!isHydrated) return;
        if (!isAuthenticated) return;

        if (!hasRequiredRole && !hasRedirected.current) {
            hasRedirected.current = true;

            // Redirect based on user's actual role
            if (role === 'player') {
                navigate('/player/dashboard', { replace: true });
            } else {
                navigate('/', { replace: true });
            }
        }
    }, [isHydrated, isAuthenticated, hasRequiredRole, role, navigate]);

    // Handle view mode enforcement
    useEffect(() => {
        if (!isHydrated) return;
        if (!isAuthenticated) return;
        if (!viewMode) return;

        // If user needs to be in player mode but isn't a player
        if (viewMode === 'player' && !isPlayer && !hasRedirected.current) {
            hasRedirected.current = true;
            navigate('/player/apply', { replace: true });
        }

        // If view mode doesn't match, switch it
        if (!hasCorrectViewMode && isPlayer) {
            if (viewMode === 'player') {
                useAuthStore.getState().switchToPlayerMode();
            } else {
                useAuthStore.getState().switchToUserMode();
            }
        }
    }, [isHydrated, isAuthenticated, viewMode, isPlayer, hasCorrectViewMode, navigate]);

    // Reset redirect flag when location changes
    useEffect(() => {
        hasRedirected.current = false;
    }, [location.pathname]);

    // Wait for hydration
    if (!isHydrated) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                    <p className="text-gray-600">加载中...</p>
                </div>
            </div>
        );
    }

    // Check authentication
    if (requiresAuth && !isAuthenticated) {
        return null; // Will redirect
    }

    // Check role
    if (!hasRequiredRole) {
        return null; // Will redirect
    }

    // Check view mode for player routes
    if (viewMode === 'player' && !isPlayer) {
        return null; // Will redirect to apply page
    }

    return <>{children}</>;
}

/**
 * Public route guard (only accessible when NOT authenticated)
 * Redirects to home if already logged in
 */
export function PublicOnlyGuard({ children }: { children: ReactNode }) {
    const navigate = useNavigate();
    const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
    const role = useAuthStore((state) => state.role);
    const isHydrated = useAuthStore.persist?.hasHydrated() ?? true;

    useEffect(() => {
        if (!isHydrated) return;

        if (isAuthenticated) {
            // Redirect based on role
            if (role === 'player') {
                navigate('/player/dashboard', { replace: true });
            } else {
                navigate('/', { replace: true });
            }
        }
    }, [isHydrated, isAuthenticated, role, navigate]);

    if (!isHydrated) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                    <p className="text-gray-600">加载中...</p>
                </div>
            </div>
        );
    }

    if (isAuthenticated) {
        return null; // Will redirect
    }

    return <>{children}</>;
}
