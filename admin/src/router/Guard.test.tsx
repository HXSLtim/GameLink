/**
 * RouteGuard 组件测试
 * 测试路由级别的权限控制
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import RouteGuard from './Guard';

// Mock usePermission hook
vi.mock('@/hooks/usePermission', () => ({
    usePermission: vi.fn(() => ({
        hasPermission: true,
        loading: false,
    })),
}));

// Mock authStore selectors (Guard uses these, not useAuthStore directly)
const mockIsAuthenticated = vi.fn((): boolean => false);
const mockUserInfo = vi.fn((): { id: number; role: string; name: string; email: string; permissions: string[] } | null => null);
const mockIsHydrated = vi.fn((): boolean => true);

vi.mock('@/stores/modules/authStore', () => ({
    useAuthStore: vi.fn(() => ({
        isAuthenticated: false,
        userInfo: null,
    })),
    useIsAuthenticated: () => mockIsAuthenticated(),
    useUserInfo: () => mockUserInfo(),
    useIsHydrated: () => mockIsHydrated(),
}));

// Import mocked modules
import { usePermission } from '@/hooks/usePermission';

const mockedUsePermission = vi.mocked(usePermission);

// Helper to render with router
const renderWithRouter = (
    ui: React.ReactElement,
    { initialEntries = ['/'] } = {}
) => {
    return render(
        <MemoryRouter initialEntries={initialEntries}>
            <Routes>
                <Route path="/admin/login" element={<div>Login Page</div>} />
                <Route path="/403" element={<div>403 Forbidden</div>} />
                <Route path="/" element={<div>Home</div>} />
                <Route path="/admin" element={ui} />
                <Route path="/admin/*" element={ui} />
            </Routes>
        </MemoryRouter>
    );
};

describe('RouteGuard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        localStorage.clear();
        
        // Default mock: not authenticated
        mockIsAuthenticated.mockReturnValue(false);
        mockUserInfo.mockReturnValue(null);
        mockIsHydrated.mockReturnValue(true);
        
        mockedUsePermission.mockReturnValue({
            hasPermission: true,
            loading: false,
        });
    });

    describe('Authentication', () => {
        it('should render children when authenticated', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'admin', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: [] 
            });

            renderWithRouter(
                <RouteGuard requiresAuth>
                    <div>Protected Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Protected Content')).toBeInTheDocument();
        });

        it('should redirect to login when not authenticated', () => {
            mockIsAuthenticated.mockReturnValue(false);
            mockUserInfo.mockReturnValue(null);

            renderWithRouter(
                <RouteGuard requiresAuth>
                    <div>Protected Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Login Page')).toBeInTheDocument();
        });

        it('should render children when requiresAuth is false', () => {
            mockIsAuthenticated.mockReturnValue(false);

            renderWithRouter(
                <RouteGuard requiresAuth={false}>
                    <div>Public Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Public Content')).toBeInTheDocument();
        });
    });

    describe('Role-based Access', () => {
        it('should render children when user has required role', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'ADMIN', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: [] 
            });

            renderWithRouter(
                <RouteGuard roles={['ADMIN']}>
                    <div>Admin Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Admin Content')).toBeInTheDocument();
        });

        it('should redirect when user does not have required role', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'USER', 
                name: 'User', 
                email: 'user@test.com', 
                permissions: [] 
            });

            renderWithRouter(
                <RouteGuard roles={['ADMIN']}>
                    <div>Admin Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Home')).toBeInTheDocument();
        });
    });

    describe('Permission-based Access', () => {
        it('should render children when user has required permission', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'ADMIN', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: ['admin.users.list'] 
            });
            mockedUsePermission.mockReturnValue({
                hasPermission: true,
                loading: false,
            });

            renderWithRouter(
                <RouteGuard requiresAuth permission="admin.users.list">
                    <div>Users List</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Users List')).toBeInTheDocument();
        });

        it('should redirect to 403 when user lacks required permission', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'ADMIN', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: [] 
            });
            mockedUsePermission.mockReturnValue({
                hasPermission: false,
                loading: false,
            });

            renderWithRouter(
                <RouteGuard requiresAuth permission="admin.users.delete">
                    <div>Delete Users</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('403 Forbidden')).toBeInTheDocument();
        });

        it('should show loading spinner while checking permissions', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'ADMIN', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: [] 
            });
            mockedUsePermission.mockReturnValue({
                hasPermission: false,
                loading: true,
            });

            const { container } = renderWithRouter(
                <RouteGuard requiresAuth permission="admin.test">
                    <div>Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            // Check for Spin component
            expect(container.querySelector('.ant-spin')).toBeInTheDocument();
        });
    });

    describe('Combined Auth and Permission', () => {
        it('should check both auth and permission', () => {
            mockIsAuthenticated.mockReturnValue(true);
            mockUserInfo.mockReturnValue({ 
                id: 1, 
                role: 'ADMIN', 
                name: 'Admin', 
                email: 'admin@test.com', 
                permissions: ['*'] 
            });
            mockedUsePermission.mockReturnValue({
                hasPermission: true,
                loading: false,
            });

            renderWithRouter(
                <RouteGuard requiresAuth roles={['ADMIN']} permission="admin.dashboard.view">
                    <div>Dashboard</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            expect(screen.getByText('Dashboard')).toBeInTheDocument();
        });
    });

    describe('Hydration', () => {
        it('should show loading while hydrating', () => {
            mockIsHydrated.mockReturnValue(false);

            const { container } = renderWithRouter(
                <RouteGuard requiresAuth>
                    <div>Content</div>
                </RouteGuard>,
                { initialEntries: ['/admin'] }
            );

            // Check for Spin component during hydration
            expect(container.querySelector('.ant-spin')).toBeInTheDocument();
        });
    });
});
