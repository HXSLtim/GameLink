/**
 * RouteGuard Component Tests
 * 
 * Tests for route protection and permission-based access control
 * Requirements: 3.1, 3.4 - Route-level permission checking
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter, Routes, Route } from 'react-router-dom';

// Mock navigate - must be before RouteGuard import
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

// Mock usePermission hook - must be before RouteGuard import
const mockUsePermission = vi.fn();
vi.mock('@/hooks/usePermission', () => ({
    usePermission: (permission: string) => mockUsePermission(permission),
    default: (permission: string) => mockUsePermission(permission),
}));

// Mock antd components
vi.mock('antd', () => ({
    message: {
        error: vi.fn(),
    },
    Spin: ({ children, tip }: { children?: React.ReactNode; tip?: string }) => (
        <div data-testid="loading-spinner">{tip}{children}</div>
    ),
}));

// Import after mocks
import RouteGuard from './Guard';

describe('RouteGuard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        localStorage.clear();
        mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });
    });

    afterEach(() => {
        localStorage.clear();
    });

    describe('Authentication', () => {
        it('should render children when authenticated and requiresAuth is true', async () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard requiresAuth>
                        <div data-testid="protected">Protected Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(screen.getByTestId('protected')).toBeInTheDocument();
            });
        });

        it('should redirect to login when not authenticated and requiresAuth is true', async () => {
            // No user_role means not authenticated
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard requiresAuth>
                        <div data-testid="protected">Protected Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/login', expect.any(Object));
            });
        });

        it('should render children when requiresAuth is false', async () => {
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard requiresAuth={false}>
                        <div data-testid="public">Public Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(screen.getByTestId('public')).toBeInTheDocument();
            });
        });
    });

    describe('Role-based Access', () => {
        it('should render children when user has required role', async () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard roles={['ADMIN']}>
                        <div data-testid="admin-content">Admin Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(screen.getByTestId('admin-content')).toBeInTheDocument();
            });
        });

        it('should redirect USER to home when accessing ADMIN route', async () => {
            localStorage.setItem('user_role', 'USER');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard roles={['ADMIN']}>
                        <div data-testid="admin-content">Admin Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/');
            });
        });

        it('should redirect COMPANION to companion home when accessing ADMIN route', async () => {
            localStorage.setItem('user_role', 'COMPANION');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard roles={['ADMIN']}>
                        <div data-testid="admin-content">Admin Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/companion');
            });
        });

        it('should redirect ADMIN to admin home when accessing USER route', async () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard roles={['USER']}>
                        <div data-testid="user-content">User Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/admin');
            });
        });

        it('should allow access when user has one of multiple allowed roles', async () => {
            localStorage.setItem('user_role', 'COMPANION');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard roles={['USER', 'COMPANION']}>
                        <div data-testid="content">Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(screen.getByTestId('content')).toBeInTheDocument();
            });
        });
    });

    describe('Permission-based Access', () => {
        it('should render children when user has required permission', () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: true, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard permission="admin.users.list">
                        <div data-testid="permission-content">Permission Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            expect(screen.getByTestId('permission-content')).toBeInTheDocument();
        });

        it('should redirect to 403 when user lacks required permission', async () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: false, loading: false });

            render(
                <MemoryRouter initialEntries={['/admin/users']}>
                    <Routes>
                        <Route path="/admin/users" element={
                            <RouteGuard permission="admin.users.delete">
                                <div data-testid="permission-content">Permission Content</div>
                            </RouteGuard>
                        } />
                    </Routes>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/403', expect.objectContaining({
                    replace: true,
                    state: expect.objectContaining({
                        from: '/admin/users',
                        requiredPermission: 'admin.users.delete'
                    })
                }));
            });
        });

        it('should show loading state while checking permissions', () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: false, loading: true });

            render(
                <MemoryRouter>
                    <RouteGuard permission="admin.users.list">
                        <div data-testid="permission-content">Permission Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            // Should not render children while loading
            expect(screen.queryByTestId('permission-content')).not.toBeInTheDocument();
        });

        it('should not check permission when permission prop is empty', () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: false, loading: false });

            render(
                <MemoryRouter>
                    <RouteGuard permission="">
                        <div data-testid="content">Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            // Should render because empty permission means no check needed
            expect(screen.getByTestId('content')).toBeInTheDocument();
        });
    });

    describe('Combined Guards', () => {
        it('should check both auth and role', async () => {
            localStorage.setItem('user_role', 'USER');

            render(
                <MemoryRouter>
                    <RouteGuard requiresAuth roles={['ADMIN']}>
                        <div data-testid="content">Content</div>
                    </RouteGuard>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/');
            });
        });

        it('should check auth, role, and permission', async () => {
            localStorage.setItem('user_role', 'ADMIN');
            mockUsePermission.mockReturnValue({ hasPermission: false, loading: false });

            render(
                <MemoryRouter initialEntries={['/admin/super']}>
                    <Routes>
                        <Route path="/admin/super" element={
                            <RouteGuard requiresAuth roles={['ADMIN']} permission="admin.super.access">
                                <div data-testid="content">Content</div>
                            </RouteGuard>
                        } />
                    </Routes>
                </MemoryRouter>
            );

            await waitFor(() => {
                expect(mockNavigate).toHaveBeenCalledWith('/403', expect.objectContaining({
                    replace: true,
                    state: expect.objectContaining({
                        from: '/admin/super',
                        requiredPermission: 'admin.super.access'
                    })
                }));
            });
        });
    });
});
