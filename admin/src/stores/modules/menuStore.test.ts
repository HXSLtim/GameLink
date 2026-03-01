/**
 * Menu Store Tests
 * Tests menu state management, permission filtering, and caching
 */

import { renderHook, act } from '@testing-library/react';
import { vi, beforeEach, afterEach, describe, it, expect } from 'vitest';
import { useAuthStore } from './authStore';

// Mock the adminApi BEFORE importing menuStore
const mockGetMyMenus = vi.fn();
vi.mock('@/api/admin', () => ({
  adminApi: {
    getMyMenus: mockGetMyMenus,
  },
}));

// Mock logger
vi.mock('@/utils/logger', () => ({
  logger: {
    error: vi.fn(),
  },
}));

// Now import menuStore after the mock is set up
import { useMenuStore } from './menuStore';

describe('menuStore', () => {
  const mockMenus = [
    {
      id: 1,
      name: 'dashboard',
      path: '/admin/dashboard',
      component: 'Dashboard',
      parentId: null,
      order: 1,
      visible: true,
      type: 'menu' as const,
      permission: 'admin.dashboard.read',
      icon: 'DashboardOutlined',
      children: [],
    },
    {
      id: 2,
      name: 'users',
      path: '/admin/users',
      component: 'Users',
      parentId: null,
      order: 2,
      visible: true,
      type: 'menu' as const,
      permission: 'admin.users.read',
      icon: 'UserOutlined',
      children: [],
    },
  ];

  beforeEach(() => {
    // Clear all mocks
    vi.clearAllMocks();

    // Reset store state before each test
    useMenuStore.setState({
      rawMenus: [],
      menus: [],
      loading: false,
      collapsed: false,
      openKeys: [],
      selectedKeys: [],
    });

    // Reset mock to return successful response
    mockGetMyMenus.mockResolvedValue({
      data: {
        success: true,
        code: 0,
        message: 'Success',
        data: mockMenus,
      },
    });

    // Mock authStore.getState to return user permissions
    vi.spyOn(useAuthStore, 'getState').mockReturnValue({
      userInfo: {
        id: 1,
        name: 'Test Admin',
        email: 'test@example.com',
        permissions: [],
      },
    } as { userInfo: { id: number; name: string; email: string; permissions: string[] } });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('fetchMenus', () => {
    it('should fetch and filter menus based on permissions', async () => {
      // Mock authStore with specific permissions
      vi.spyOn(useAuthStore, 'getState').mockReturnValue({
        userInfo: {
          permissions: ['admin.dashboard.read', 'admin.users.read'],
        },
      } as { userInfo: { permissions: string[] } });

      const { result } = renderHook(() => useMenuStore());

      await act(async () => {
        await result.current.fetchMenus();
      });

      // Verify API was called
      expect(mockGetMyMenus).toHaveBeenCalled();
      expect(result.current.loading).toBe(false);
      expect(result.current.rawMenus.length).toBeGreaterThan(0);
      expect(result.current.menus.length).toBeGreaterThan(0);
    });

    it('should show all menus for super admin', async () => {
      vi.spyOn(useAuthStore, 'getState').mockReturnValue({
        userInfo: {
          permissions: ['*'],
        },
      } as { userInfo: { permissions: string[] } });

      const { result } = renderHook(() => useMenuStore());

      await act(async () => {
        await result.current.fetchMenus();
      });

      expect(mockGetMyMenus).toHaveBeenCalled();
      expect(result.current.menus.length).toBe(result.current.rawMenus.length);
    });

    it('should filter menus without permissions', async () => {
      vi.spyOn(useAuthStore, 'getState').mockReturnValue({
        userInfo: {
          permissions: ['admin.dashboard.read'],
        },
      } as { userInfo: { permissions: string[] } });

      const { result } = renderHook(() => useMenuStore());

      await act(async () => {
        await result.current.fetchMenus();
      });

      // Should only show dashboard menu
      expect(mockGetMyMenus).toHaveBeenCalled();
      expect(result.current.menus.length).toBe(1);
      expect(result.current.menus[0].key).toBe('dashboard');
    });

    it('should handle API errors gracefully', async () => {
      // Mock API error
      mockGetMyMenus.mockRejectedValue(new Error('Network error'));

      const { result } = renderHook(() => useMenuStore());

      await act(async () => {
        await result.current.fetchMenus();
      });

      expect(mockGetMyMenus).toHaveBeenCalled();
      expect(result.current.loading).toBe(false);
      expect(result.current.menus.length).toBe(0);
    });
  });

  describe('filterMenusByPermission', () => {
    it('should filter nested menu structure', () => {
      const { result } = renderHook(() => useMenuStore());

      const mockMenus = [
        {
          id: 1,
          key: 'parent',
          label: 'Parent',
          sort: 1,
          permission: 'admin.parent.read',
          children: [
            {
              id: 11,
              key: 'child1',
              label: 'Child 1',
              sort: 1,
              permission: 'admin.child1.read',
            },
            {
              id: 12,
              key: 'child2',
              label: 'Child 2',
              sort: 2,
              permission: 'admin.child2.read',
            },
          ],
        },
      ];

      const filtered = result.current.filterMenusByPermission(mockMenus, ['admin.parent.read', 'admin.child1.read']);

      expect(filtered.length).toBe(1);
      expect(filtered[0].children?.length).toBe(1);
      expect(filtered[0].children?.[0].key).toBe('child1');
    });

    it('should remove parent with no visible children', () => {
      const { result } = renderHook(() => useMenuStore());

      const mockMenus = [
        {
          id: 1,
          key: 'parent',
          label: 'Parent',
          sort: 1,
          permission: 'admin.parent.read',
          children: [
            {
              id: 11,
              key: 'child1',
              label: 'Child 1',
              sort: 1,
              permission: 'admin.child1.read',
            },
          ],
        },
      ];

      const filtered = result.current.filterMenusByPermission(mockMenus, ['admin.parent.read']);

      // Parent should be removed because it has no visible children
      expect(filtered.length).toBe(0);
    });
  });

  describe('getMenuByKey', () => {
    it('should find menu by key', () => {
      const { result } = renderHook(() => useMenuStore());

      act(() => {
        result.current.rawMenus = [
          {
            id: 1,
            key: 'dashboard',
            label: 'Dashboard',
            sort: 1,
          },
        ];
        result.current.menus = result.current.rawMenus;
      });

      const menu = result.current.getMenuByKey('dashboard');
      expect(menu?.key).toBe('dashboard');
    });

    it('should find nested menu by key', () => {
      const { result } = renderHook(() => useMenuStore());

      act(() => {
        result.current.rawMenus = [
          {
            id: 1,
            key: 'parent',
            label: 'Parent',
            sort: 1,
            children: [
              {
                id: 11,
                key: 'child',
                label: 'Child',
                sort: 1,
              },
            ],
          },
        ];
        result.current.menus = result.current.rawMenus;
      });

      const menu = result.current.getMenuByKey('child');
      expect(menu?.key).toBe('child');
    });
  });

  describe('collapsed state', () => {
    it('should toggle collapsed state', () => {
      const { result } = renderHook(() => useMenuStore());

      expect(result.current.collapsed).toBe(false);

      act(() => {
        result.current.toggleCollapsed();
      });

      expect(result.current.collapsed).toBe(true);

      act(() => {
        result.current.toggleCollapsed();
      });

      expect(result.current.collapsed).toBe(false);
    });

    it('should set collapsed state', () => {
      const { result } = renderHook(() => useMenuStore());

      act(() => {
        result.current.setCollapsed(true);
      });

      expect(result.current.collapsed).toBe(true);
    });
  });

  describe('getBreadcrumbs', () => {
    it('should return breadcrumbs for menu key', () => {
      const { result } = renderHook(() => useMenuStore());

      act(() => {
        result.current.rawMenus = [
          {
            id: 1,
            key: 'dashboard',
            label: 'Dashboard',
            path: '/admin/dashboard',
            sort: 1,
          },
        ];
        result.current.menus = result.current.rawMenus;
      });

      const breadcrumbs = result.current.getBreadcrumbs('dashboard');
      expect(breadcrumbs.length).toBe(1);
      expect(breadcrumbs[0].key).toBe('dashboard');
    });
  });

  describe('getMenuPaths', () => {
    it('should return all menu paths', () => {
      const { result } = renderHook(() => useMenuStore());

      act(() => {
        result.current.rawMenus = [
          {
            id: 1,
            key: 'dashboard',
            label: 'Dashboard',
            path: '/admin/dashboard',
            sort: 1,
          },
          {
            id: 2,
            key: 'users',
            label: 'Users',
            path: '/admin/users',
            sort: 2,
            children: [
              {
                id: 21,
                key: 'users-list',
                label: 'User List',
                path: '/admin/users/list',
                sort: 1,
              },
            ],
          },
        ];
        result.current.menus = result.current.rawMenus;
      });

      const paths = result.current.getMenuPaths();
      expect(paths['dashboard']).toBe('/admin/dashboard');
      expect(paths['users']).toBe('/admin/users');
      expect(paths['users-list']).toBe('/admin/users/list');
    });
  });
});
