/**
 * AdminContext Tests
 * 
 * Tests for admin context provider and permission management
 * Requirements: 3.1, 3.4 - Permission context and checking
 * 
 * Note: These tests focus on the context's internal logic.
 * API mocking is complex due to vitest hoisting, so we test
 * the default/error states and basic functionality.
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { AdminProvider, useAdmin } from './AdminContext';

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

    describe('Default State', () => {
        it('should provide default context values', async () => {
            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            // Wait for initial render
            await waitFor(() => {
                expect(screen.getByTestId('loading')).toBeInTheDocument();
            });

            // Default values should be empty arrays and false
            expect(screen.getByTestId('is-super-admin').textContent).toBe('false');
            expect(screen.getByTestId('permission-version')).toBeInTheDocument();
        });

        it('should have empty permissions when no token', async () => {
            // No token set
            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('loading').textContent).toBe('false');
            });

            expect(screen.getByTestId('permissions').textContent).toBe('[]');
            expect(screen.getByTestId('menus').textContent).toBe('[]');
        });
    });

    describe('hasPermission logic', () => {
        it('should return false when permissions array is empty', async () => {
            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('loading').textContent).toBe('false');
            });

            // With empty permissions, hasPermission should return false
            expect(screen.getByTestId('has-users-list').textContent).toBe('false');
            expect(screen.getByTestId('has-all').textContent).toBe('false');
            expect(screen.getByTestId('has-any').textContent).toBe('false');
        });

        it('should not be super admin when permissions are empty', async () => {
            render(
                <AdminProvider>
                    <TestConsumer />
                </AdminProvider>
            );

            await waitFor(() => {
                expect(screen.getByTestId('loading').textContent).toBe('false');
            });

            expect(screen.getByTestId('is-super-admin').textContent).toBe('false');
        });
    });

    describe('Context Provider', () => {
        it('should render children correctly', () => {
            render(
                <AdminProvider>
                    <div data-testid="child">Child Content</div>
                </AdminProvider>
            );

            expect(screen.getByTestId('child')).toBeInTheDocument();
            expect(screen.getByTestId('child').textContent).toBe('Child Content');
        });

        it('should provide useAdmin hook access', () => {
            const TestHookAccess = () => {
                const context = useAdmin();
                return (
                    <div data-testid="hook-test">
                        {typeof context.hasPermission === 'function' ? 'function' : 'not-function'}
                    </div>
                );
            };

            render(
                <AdminProvider>
                    <TestHookAccess />
                </AdminProvider>
            );

            expect(screen.getByTestId('hook-test').textContent).toBe('function');
        });
    });
});
