/**
 * usePermission Hook 测试
 * 测试权限检查逻辑
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook } from '@testing-library/react';
import { usePermission, usePermissions } from './usePermission';

// Mock useAdmin hook
vi.mock('@/context/useAdmin', () => ({
    useAdmin: vi.fn(() => ({
        permissions: [],
        loading: false,
    })),
}));

import { useAdmin } from '@/context/useAdmin';
const mockedUseAdmin = vi.mocked(useAdmin);

describe('usePermission', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Default: no permissions
        mockedUseAdmin.mockReturnValue({
            permissions: [],
            loading: false,
            menus: [],
            refreshMenus: vi.fn(),
            hasPermission: vi.fn(),
        });
    });

    describe('Single Permission Check', () => {
        it('should return false when user has no permissions', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: [],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('admin.users.list'));

            expect(result.current.hasPermission).toBe(false);
            expect(result.current.loading).toBe(false);
        });

        it('should return true when user has the required permission', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.users.list', 'admin.users.create'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('admin.users.list'));

            expect(result.current.hasPermission).toBe(true);
        });

        it('should return false when user lacks the required permission', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.users.list'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('admin.users.delete'));

            expect(result.current.hasPermission).toBe(false);
        });

        it('should return true for empty permission string', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: [],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission(''));

            expect(result.current.hasPermission).toBe(true);
        });
    });

    describe('Super Admin Wildcard', () => {
        it('should return true for any permission when user has wildcard (*)', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['*'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('any.random.permission'));

            expect(result.current.hasPermission).toBe(true);
        });
    });

    describe('Array Permissions - Any Mode', () => {
        it('should return true when user has at least one of the required permissions', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.users.create'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() =>
                usePermission(['admin.users.create', 'admin.users.delete'])
            );

            expect(result.current.hasPermission).toBe(true);
        });

        it('should return false when user has none of the required permissions', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.orders.list'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() =>
                usePermission(['admin.users.create', 'admin.users.delete'])
            );

            expect(result.current.hasPermission).toBe(false);
        });
    });

    describe('Array Permissions - All Mode', () => {
        it('should return true when user has all required permissions', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.users.create', 'admin.users.delete', 'admin.users.list'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() =>
                usePermission(['admin.users.create', 'admin.users.delete'], 'all')
            );

            expect(result.current.hasPermission).toBe(true);
        });

        it('should return false when user lacks any of the required permissions', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.users.create'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() =>
                usePermission(['admin.users.create', 'admin.users.delete'], 'all')
            );

            expect(result.current.hasPermission).toBe(false);
        });
    });

    describe('Loading State', () => {
        it('should return loading true when permissions are being fetched', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: [],
                loading: true,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('admin.test'));

            expect(result.current.loading).toBe(true);
        });

        it('should return loading false when permissions are loaded', () => {
            mockedUseAdmin.mockReturnValue({
                permissions: ['admin.test'],
                loading: false,
                menus: [],
                refreshMenus: vi.fn(),
                hasPermission: vi.fn(),
            });

            const { result } = renderHook(() => usePermission('admin.test'));

            expect(result.current.loading).toBe(false);
        });
    });
});

describe('usePermissions', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should check multiple permissions and return object with results', () => {
        mockedUseAdmin.mockReturnValue({
            permissions: ['admin.users.list', 'admin.orders.list'],
            loading: false,
            menus: [],
            refreshMenus: vi.fn(),
            hasPermission: vi.fn(),
        });

        const { result } = renderHook(() =>
            usePermissions({
                canViewUsers: 'admin.users.list',
                canDeleteUsers: 'admin.users.delete',
                canViewOrders: 'admin.orders.list',
            })
        );

        expect(result.current.canViewUsers).toBe(true);
        expect(result.current.canDeleteUsers).toBe(false);
        expect(result.current.canViewOrders).toBe(true);
    });

    it('should handle wildcard permissions', () => {
        mockedUseAdmin.mockReturnValue({
            permissions: ['*'],
            loading: false,
            menus: [],
            refreshMenus: vi.fn(),
            hasPermission: vi.fn(),
        });

        const { result } = renderHook(() =>
            usePermissions({
                canDoAnything: 'any.permission',
                canDoSomethingElse: 'another.permission',
            })
        );

        expect(result.current.canDoAnything).toBe(true);
        expect(result.current.canDoSomethingElse).toBe(true);
    });

    it('should return all false when user has no permissions', () => {
        mockedUseAdmin.mockReturnValue({
            permissions: [],
            loading: false,
            menus: [],
            refreshMenus: vi.fn(),
            hasPermission: vi.fn(),
        });

        const { result } = renderHook(() =>
            usePermissions({
                canView: 'admin.view',
                canEdit: 'admin.edit',
                canDelete: 'admin.delete',
            })
        );

        expect(result.current.canView).toBe(false);
        expect(result.current.canEdit).toBe(false);
        expect(result.current.canDelete).toBe(false);
    });
});
