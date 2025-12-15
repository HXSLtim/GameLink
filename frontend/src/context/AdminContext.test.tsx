/**
 * AdminContext Tests
 * 
 * Tests for admin context provider and permission management
 * Requirements: 3.1, 3.4 - Permission context and checking
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { AdminProvider, useAdmin } from './AdminContext';

// Mock permission store - must be before AdminContext import uses it
const mockSetPermissions = vi.fn();
const mockClearPermissions = vi.fn();

vi.mock('@/utils/permission', () => ({
    permissionStore: {
        setPermissions: mockSetPermissions,
        clearPermissions: mockClearPermissions,
    },
}));

// Mock the admin API
// Note: apiClient interceptor returns response.data, so API returns { success, code, data }
// Then AdminContext extracts .data from that
const mockGetMyPermissions = vi.fn();
const mockGetMyMenus = vi.fn();

vi.mock('@/api/admin', () => ({
    adminApi: {
        getMyPermissions: () => mockGetMyPermissions(),
        getMyMenus: () => mockGetMyMenus(),
    },
}));

// Mock menuPermission filter to return menus as-is for testing
vi.mock('@/utils/menuPermission', () => ({
    filterMenusByPermission: (menus: unknown[]) => menus,
}));

// Test component to access context
const TestConsumer = () => {
    const {
        permissions,
        menus,
        rawMenus,
        loading,
        hasPermission,
        hasAllPermissions,
        hasAnyPermission,
        isSuperAdmin,
        permissionVersion,
    } = useAdmin();

    return (
        <div>
            <div data-testid="loading">{loading.toString()}</div>
            <div data-testid="permissions">{JSON.stringify(permissions)}</div>
            <div data-testid="menus">{JSON.stringify(menus)}</div>
            <div data-testid="raw-menus">{JSON.stringify(rawMenus)}</div>
            <div data-testid="is-super-admin">{isSuperAdmin.toString()}</div>
            <div data-testid="has-users-list">{hasPermission('admin.users.list').toString()}</div>
            <div data-testid="has-all">{hasAllPermissions(['admin.users.list', 'admin.users.create']).toString()}</div>
            <div data-testid="has-any">{hasAnyPermission(['admin.users.list', 'admin.users.delete']).toString()}</div>
            <div data-testid="permission-version">{permissionVersion}</div>
        </div>
    );
};

describe('AdminContext', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        localStorage.clear();
    });

    afterEach(() => {
        localStorage.clear();
    });

    describe('Initial State', () => {
        it('should have empty permissions when no token', async () => {
            // No token set - should clear permissions
            mockGetMyPermissions.mockResolvedValue({ data: [] });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe('[]');
            });
        });

        it('should load permissions when token exists', async () => {
            localStorage.setItem('token', 'test-token');
            const permissions = ['admin.users.list', 'admin.users.create'];
            const menus = [{ id: 1, name: 'Users', path: '/admin/users' }];

            // API returns { success, code, data } format after interceptor
            mockGetMyPermissions.mockResolvedValue({ data: permissions });
            mockGetMyMenus.mockResolvedValue({ data: menus });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(permissions));
            });
        });
    });

    describe('hasPermission', () => {
        it('should return true when user has the permission', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.users.list', 'admin.users.create'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-users-list').textContent).toBe('true');
            });
        });

        it('should return false when user lacks the permission', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.orders.list'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-users-list').textContent).toBe('false');
            });
        });

        it('should return true for super admin regardless of permission', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['*'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('is-super-admin').textContent).toBe('true');
                expect(screen.getByTestId('has-users-list').textContent).toBe('true');
            });
        });
    });

    describe('hasAllPermissions', () => {
        it('should return true when user has all permissions', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.users.list', 'admin.users.create', 'admin.users.delete'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-all').textContent).toBe('true');
            });
        });

        it('should return false when user is missing any permission', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.users.list'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-all').textContent).toBe('false');
            });
        });
    });

    describe('hasAnyPermission', () => {
        it('should return true when user has at least one permission', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.users.list'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-any').textContent).toBe('true');
            });
        });

        it('should return false when user has none of the permissions', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockResolvedValue({
                data: ['admin.orders.list'],
            });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('has-any').textContent).toBe('false');
            });
        });
    });

    describe('Permission Store Sync', () => {
        it('should sync permissions to store on load', async () => {
            localStorage.setItem('token', 'test-token');
            const permissions = ['admin.users.list'];
            mockGetMyPermissions.mockResolvedValue({ data: permissions });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(mockSetPermissions).toHaveBeenCalledWith(permissions);
            });
        });

        it('should clear permissions from store when no token', async () => {
            mockGetMyPermissions.mockResolvedValue({ data: [] });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(mockClearPermissions).toHaveBeenCalled();
            });
        });
    });

    describe('Error Handling', () => {
        it('should handle API errors gracefully', async () => {
            localStorage.setItem('token', 'test-token');
            mockGetMyPermissions.mockRejectedValue(new Error('API Error'));
            mockGetMyMenus.mockRejectedValue(new Error('API Error'));

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe('[]');
                expect(screen.getByTestId('menus').textContent).toBe('[]');
            });
        });
    });

    describe('Permission Change Events', () => {
        it('should refresh permissions when permission change event is triggered', async () => {
            localStorage.setItem('token', 'test-token');
            const initialPermissions = ['admin.users.list'];
            const updatedPermissions = ['admin.users.list', 'admin.users.create'];

            mockGetMyPermissions
                .mockResolvedValueOnce({ data: initialPermissions })
                .mockResolvedValueOnce({ data: updatedPermissions });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            // Wait for initial load
            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(initialPermissions));
            });

            // Trigger permission change event
            await act(async () => {
                window.dispatchEvent(new CustomEvent('gamelink:permission-change'));
            });

            // Wait for refresh
            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(updatedPermissions));
            });
        });

        it('should refresh permissions when storage event is triggered from another tab', async () => {
            localStorage.setItem('token', 'test-token');
            const initialPermissions = ['admin.users.list'];
            const updatedPermissions = ['admin.users.list', 'admin.orders.list'];

            mockGetMyPermissions
                .mockResolvedValueOnce({ data: initialPermissions })
                .mockResolvedValueOnce({ data: updatedPermissions });
            mockGetMyMenus.mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            // Wait for initial load
            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(initialPermissions));
            });

            // Simulate storage event from another tab
            await act(async () => {
                const storageEvent = new StorageEvent('storage', {
                    key: 'permission_change_timestamp',
                    newValue: Date.now().toString(),
                });
                window.dispatchEvent(storageEvent);
            });

            // Wait for refresh
            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(updatedPermissions));
            });
        });
    });
});
