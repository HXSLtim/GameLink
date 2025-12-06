/**
 * AdminContext Tests
 * 
 * Tests for admin context provider and permission management
 * Requirements: 3.1, 3.4 - Permission context and checking
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { AdminProvider, useAdmin } from './AdminContext';

// Mock the admin API
vi.mock('@/api/admin', () => ({
    adminApi: {
        getMyPermissions: vi.fn(),
        getMyMenus: vi.fn(),
    },
}));

// Mock permission store
vi.mock('@/utils/permission', () => ({
    permissionStore: {
        setPermissions: vi.fn(),
        clearPermissions: vi.fn(),
    },
}));

import { adminApi } from '@/api/admin';
import { permissionStore } from '@/utils/permission';

// Test component to access context
const TestConsumer = () => {
    const {
        permissions,
        menus,
        loading,
        hasPermission,
        hasAllPermissions,
        hasAnyPermission,
        isSuperAdmin,
    } = useAdmin();

    return (
        <div>
            <div data-testid="loading">{loading.toString()}</div>
            <div data-testid="permissions">{JSON.stringify(permissions)}</div>
            <div data-testid="menus">{JSON.stringify(menus)}</div>
            <div data-testid="is-super-admin">{isSuperAdmin.toString()}</div>
            <div data-testid="has-users-list">{hasPermission('admin.users.list').toString()}</div>
            <div data-testid="has-all">{hasAllPermissions(['admin.users.list', 'admin.users.create']).toString()}</div>
            <div data-testid="has-any">{hasAnyPermission(['admin.users.list', 'admin.users.delete']).toString()}</div>
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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            const mockPermissions = ['admin.users.list', 'admin.users.create'];
            const mockMenus = [{ id: 1, name: 'Users', path: '/admin/users' }];

            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: mockPermissions,
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: mockMenus,
            });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('permissions').textContent).toBe(JSON.stringify(mockPermissions));
                expect(screen.getByTestId('menus').textContent).toBe(JSON.stringify(mockMenus));
            });
        });
    });

    describe('hasPermission', () => {
        it('should return true when user has the permission', async () => {
            localStorage.setItem('token', 'test-token');
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.users.list', 'admin.users.create'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.orders.list'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['*'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.users.list', 'admin.users.create', 'admin.users.delete'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.users.list'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.users.list'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: ['admin.orders.list'],
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

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
            const mockPermissions = ['admin.users.list'];
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({
                data: mockPermissions,
            });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(permissionStore.setPermissions).toHaveBeenCalledWith(mockPermissions);
            });
        });

        it('should clear permissions from store when no token', async () => {
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] });

            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(permissionStore.clearPermissions).toHaveBeenCalled();
            });
        });
    });

    describe('Error Handling', () => {
        it('should handle API errors gracefully', async () => {
            localStorage.setItem('token', 'test-token');
            (adminApi.getMyPermissions as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));
            (adminApi.getMyMenus as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

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
});
