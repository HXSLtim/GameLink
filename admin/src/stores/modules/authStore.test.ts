/**
 * Auth Store Tests
 * Tests authentication, authorization, permission checking, and state management
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useAuthStore } from './authStore';
import type { UserInfo } from '../types';
import * as authApi from '@/api/auth';

// Mock auth API
vi.mock('@/api/auth', () => ({
  authApi: {
    login: vi.fn(),
    logout: vi.fn(),
  },
}));

describe('authStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useAuthStore.getState().setToken(null);
    useAuthStore.getState().setUserInfo(null);
    useAuthStore.getState().clearError();

    // Clear sessionStorage (using the mock from setup.ts)
    sessionStorage.clear();

    // Clear all mocks
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('State Initialization', () => {
    it('should initialize with default state', () => {
      const { result } = renderHook(() => useAuthStore());

      expect(result.current.token).toBeNull();
      expect(result.current.userInfo).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('Login', () => {
    it('should login successfully with valid credentials', async () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin User',
        email: 'admin@gamelink.com',
        phone: '13800138000',
        avatar: 'https://example.com/avatar.jpg',
        role: 'admin',
        permissions: ['admin.dashboard.read', 'admin.users.read'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const mockToken = 'mock-jwt-token';
      const mockResponse = {
        data: {
          data: {
            token: mockToken,
            user: mockUser,
          },
        },
      };

      vi.mocked(authApi.authApi.login).mockResolvedValue(mockResponse as {
        data: { data: { token: string; user: UserInfo } };
      });

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.login({ username: 'admin', password: '123456' });
      });

      expect(result.current.token).toBe(mockToken);
      expect(result.current.userInfo).toEqual(mockUser);
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(sessionStorage.getItem('auth_token')).toBe(mockToken);
    });

    it('should handle login failure with invalid credentials', async () => {
      const mockError = {
        response: {
          data: {
            message: 'Invalid username or password',
          },
        },
      };

      vi.mocked(authApi.authApi.login).mockRejectedValue(mockError as {
        response: { data: { message: string } };
      });

      const { result } = renderHook(() => useAuthStore());

      await expect(
        act(async () => {
          await result.current.login({ username: 'invalid', password: 'wrong' });
        })
      ).rejects.toThrow();

      await waitFor(() => {
        expect(result.current.token).toBeNull();
        expect(result.current.userInfo).toBeNull();
        expect(result.current.isAuthenticated).toBe(false);
        expect(result.current.loading).toBe(false);
      });
    });

    it('should handle login failure with network error', async () => {
      const mockError = new Error('Network Error');

      vi.mocked(authApi.authApi.login).mockRejectedValue(mockError);

      const { result } = renderHook(() => useAuthStore());

      await expect(
        act(async () => {
          await result.current.login({ username: 'admin', password: '123456' });
        })
      ).rejects.toThrow('Network Error');
    });

    it('should handle invalid server response', async () => {
      const mockResponse = {
        data: {
          data: {
            token: '',
            user: null,
          },
        },
      };

      vi.mocked(authApi.authApi.login).mockResolvedValue(mockResponse as {
        data: { data: { token: string; user: UserInfo } };
      });

      const { result } = renderHook(() => useAuthStore());

      await expect(
        act(async () => {
          await result.current.login({ username: 'admin', password: '123456' });
        })
      ).rejects.toThrow('Invalid response from server');
    });
  });

  describe('Logout', () => {
    it('should logout successfully and clear state', async () => {
      // Setup: User is logged in
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin User',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('mock-token');
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.isAuthenticated).toBe(true);

      // Mock logout API
      vi.mocked(authApi.authApi.logout).mockResolvedValue({} as Record<string, never>);

      // Mock other stores' reset methods
      vi.doMock('./orderStore', () => ({
        useOrderStore: {
          getState: () => ({ reset: vi.fn() }),
        },
      }));
      vi.doMock('./playerStore', () => ({
        clearPlayerStore: vi.fn(),
      }));
      vi.doMock('./chatStore', () => ({
        useChatStore: {
          getState: () => ({ reset: vi.fn() }),
        },
      }));

      await act(async () => {
        await result.current.logout();
      });

      expect(result.current.token).toBeNull();
      expect(result.current.userInfo).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.loading).toBe(false);
      expect(sessionStorage.getItem('auth_token')).toBeNull();
    });

    it('should handle logout API failure gracefully', async () => {
      // Setup: User is logged in
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin User',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('mock-token');
        result.current.setUserInfo(mockUser);
      });

      // Mock logout API failure
      vi.mocked(authApi.authApi.logout).mockRejectedValue(new Error('API Error'));

      await act(async () => {
        await result.current.logout();
      });

      // State should still be cleared even if API fails
      expect(result.current.token).toBeNull();
      expect(result.current.userInfo).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe('Permission Checking', () => {
    it('should return true for super admin with * permission', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Super Admin',
        email: 'superadmin@gamelink.com',
        role: 'superAdmin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.hasPermission('admin.users.read')).toBe(true);
      expect(result.current.hasPermission('any.random.permission')).toBe(true);
      expect(result.current.hasPermission(['admin.users.read', 'admin.orders.write'])).toBe(true);
    });

    it('should check single permission correctly', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['admin.dashboard.read', 'admin.users.read'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.hasPermission('admin.dashboard.read')).toBe(true);
      expect(result.current.hasPermission('admin.users.write')).toBe(false);
      expect(result.current.hasPermission('admin.orders.read')).toBe(false);
    });

    it('should check multiple permissions with "any" mode', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['admin.dashboard.read', 'admin.users.read'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.hasPermission(['admin.dashboard.read', 'admin.users.write'], 'any')).toBe(true);
      expect(result.current.hasPermission(['admin.orders.read', 'admin.users.write'], 'any')).toBe(false);
      expect(result.current.hasAnyPermission(['admin.dashboard.read', 'admin.users.write'])).toBe(true);
    });

    it('should check multiple permissions with "all" mode', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['admin.dashboard.read', 'admin.users.read'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.hasPermission(['admin.dashboard.read', 'admin.users.read'], 'all')).toBe(true);
      expect(result.current.hasPermission(['admin.dashboard.read', 'admin.users.write'], 'all')).toBe(false);
      expect(result.current.hasAllPermissions(['admin.dashboard.read', 'admin.users.read'])).toBe(true);
      expect(result.current.hasAllPermissions(['admin.dashboard.read', 'admin.users.write'])).toBe(false);
    });

    it('should return false when user is not authenticated', () => {
      const { result } = renderHook(() => useAuthStore());

      expect(result.current.hasPermission('admin.users.read')).toBe(false);
      expect(result.current.hasPermission(['admin.users.read'])).toBe(false);
      expect(result.current.isAdmin()).toBe(false);
    });
  });

  describe('Role Checking', () => {
    it('should identify admin users correctly', () => {
      const adminUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(adminUser);
      });

      expect(result.current.isAdmin()).toBe(true);
      expect(result.current.hasRole('admin')).toBe(true);
      expect(result.current.hasRole('Admin')).toBe(true); // Case insensitive
    });

    it('should identify super admin users correctly', () => {
      const superAdminUser: UserInfo = {
        id: 1,
        name: 'Super Admin',
        email: 'superadmin@gamelink.com',
        role: 'superAdmin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(superAdminUser);
      });

      expect(result.current.isAdmin()).toBe(true);
      expect(result.current.hasRole('superAdmin')).toBe(true);
      expect(result.current.hasRole('superadmin')).toBe(true); // Case insensitive
    });

    it('should not identify regular users as admin', () => {
      const regularUser: UserInfo = {
        id: 1,
        name: 'User',
        email: 'user@gamelink.com',
        role: 'user',
        permissions: [],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(regularUser);
      });

      expect(result.current.isAdmin()).toBe(false);
      expect(result.current.hasRole('admin')).toBe(false);
      expect(result.current.hasRole('user')).toBe(true);
    });

    it('should check multiple roles correctly', () => {
      const adminUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'admin',
        permissions: ['*'],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(adminUser);
      });

      expect(result.current.hasRole(['admin', 'superAdmin'])).toBe(true);
      expect(result.current.hasRole(['user', 'player'])).toBe(false);
    });
  });

  describe('Token Management', () => {
    it('should set token and update authentication status', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('test-token');
      });

      expect(result.current.token).toBe('test-token');
      expect(result.current.isAuthenticated).toBe(true);
      expect(sessionStorage.getItem('auth_token')).toBe('test-token');
    });

    it('should clear token and update authentication status', () => {
      const { result } = renderHook(() => useAuthStore());

      // Set token first
      act(() => {
        result.current.setToken('test-token');
      });

      expect(result.current.isAuthenticated).toBe(true);

      // Clear token
      act(() => {
        result.current.setToken('');
      });

      expect(result.current.token).toBe('');
      expect(result.current.isAuthenticated).toBe(false);
      expect(sessionStorage.getItem('auth_token')).toBeNull();
    });

    it('should update user info', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Test User',
        email: 'test@gamelink.com',
        role: 'user',
        permissions: [],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.userInfo).toEqual(mockUser);
    });
  });

  describe('Error Handling', () => {
    it('should clear error message', async () => {
      const { result } = renderHook(() => useAuthStore());

      // Simulate error state
      act(() => {
        useAuthStore.setState({ error: 'Test error' });
      });

      expect(useAuthStore.getState().error).toBe('Test error');

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });

    it('should set loading state', () => {
      const { result } = renderHook(() => useAuthStore());

      expect(result.current.loading).toBe(false);

      act(() => {
        result.current.setLoading(true);
      });

      expect(result.current.loading).toBe(true);

      act(() => {
        result.current.setLoading(false);
      });

      expect(result.current.loading).toBe(false);
    });
  });

  describe('Edge Cases', () => {
    it('should handle empty permissions array', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'User',
        email: 'user@gamelink.com',
        role: 'user',
        permissions: [],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.hasPermission('admin.users.read')).toBe(false);
      expect(result.current.hasPermission([])).toBe(false);
    });

    it('should handle null user info for permission checks', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(null);
      });

      expect(result.current.hasPermission('admin.users.read')).toBe(false);
      expect(result.current.isAdmin()).toBe(false);
      expect(result.current.hasRole('admin')).toBe(false);
    });

    it('should handle case-insensitive role matching', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Admin',
        email: 'admin@gamelink.com',
        role: 'Admin', // Capital A
        permissions: [],
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.isAdmin()).toBe(true);
      expect(result.current.hasRole('admin')).toBe(true);
      expect(result.current.hasRole('ADMIN')).toBe(true);
    });
  });
});
