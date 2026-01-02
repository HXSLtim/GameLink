/**
 * Auth Store - Authentication and Authorization State Management
 *
 * Features:
 * - Token storage in sessionStorage (cleared on tab close)
 * - Permission checking methods (hasPermission, isAdmin)
 * - Logout clears all stores
 * - Integrated with auth API
 *
 * Security:
 * - Does NOT persist to localStorage (security best practice)
 * - Token stored in sessionStorage only
 * - Auto-clears on logout
 */

import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { UserInfo, LoginRequest } from '../types';
import { authApi } from '@/api/auth';

import { logger } from '@/utils/logger';
interface AuthState {
  // State
  token: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;

  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  setToken: (token: string) => void;
  setUserInfo: (user: UserInfo) => void;
  clearError: () => void;
  setLoading: (loading: boolean) => void;

  // Selectors (computed values)
  isAdmin: () => boolean;
  hasPermission: (permission: string | string[], mode?: 'any' | 'all') => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasRole: (role: string | string[]) => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial State
      token: null,
      userInfo: null,
      isAuthenticated: false,
      loading: false,
      error: null,

      // Actions
      login: async (credentials: LoginRequest) => {
        set({ loading: true, error: null });

        try {
          // Call login API
          const response = await authApi.login(credentials);
          const { token, user } = response.data.data;

          // Validate response
          if (!token || !user) {
            throw new Error('Invalid response from server');
          }

          // Map API user response to UserInfo interface
          const userInfo: UserInfo = {
            id: user.id,
            name: user.username,
            email: user.email,
            phone: undefined,
            avatar: undefined,
            role: user.role,
            permissions: [], // TODO: Fetch permissions from API
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };

          // Update state
          set({
            token,
            userInfo,
            isAuthenticated: true,
            loading: false,
            error: null,
          });

          // Sync to sessionStorage for API interceptor
          sessionStorage.setItem('auth_token', token);

        } catch (error: unknown) {
          const errorMessage = error && typeof error === 'object' && 'response' in error
            ? (error as { response?: { data?: { message?: string } } }).response?.data?.message
            : error instanceof Error
            ? error.message
            : '登录失败，请检查用户名和密码';
          set({
            error: errorMessage,
            loading: false,
            isAuthenticated: false,
            token: null,
            userInfo: null,
          });
          throw error;
        }
      },

      logout: async () => {
        set({ loading: true });

        try {
          // Call logout API (ignore errors for local cleanup)
          await authApi.logout().catch(() => {
            logger.warn('Logout API call failed, proceeding with local cleanup');
          });

          // Clear auth state
          set({
            token: null,
            userInfo: null,
            isAuthenticated: false,
            loading: false,
            error: null,
          });

          // Clear sessionStorage
          sessionStorage.removeItem('auth_token');
          sessionStorage.removeItem('auth-storage');

          // Clear localStorage (legacy compatibility)
          localStorage.removeItem('token');
          localStorage.removeItem('user_role');
          localStorage.removeItem('user_info');

          // Clear other stores
          // Note: Dynamic imports to avoid circular dependencies
          try {
            const { useOrderStore } = await import('./orderStore');
            useOrderStore.getState().reset();
          } catch (e) {
            logger.warn('Failed to clear orderStore:', e);
          }

          try {
            const { clearPlayerStore } = await import('./playerStore');
            clearPlayerStore();
          } catch (e) {
            logger.warn('Failed to clear playerStore:', e);
          }

          try {
            const { useChatStore } = await import('./chatStore');
            useChatStore.getState().reset();
          } catch (e) {
            logger.warn('Failed to clear chatStore:', e);
          }

        } catch (error) {
          logger.error('Logout error:', error);
          // Force logout even on error
          set({
            token: null,
            userInfo: null,
            isAuthenticated: false,
            loading: false,
          });
          sessionStorage.removeItem('auth_token');
          sessionStorage.removeItem('auth-storage');
        }
      },

      setToken: (token: string) => {
        const isAuthenticated = !!token;
        set({ token, isAuthenticated });

        // Sync to sessionStorage
        if (token) {
          sessionStorage.setItem('auth_token', token);
        } else {
          sessionStorage.removeItem('auth_token');
        }
      },

      setUserInfo: (userInfo: UserInfo) => {
        set({ userInfo });
      },

      clearError: () => {
        set({ error: null });
      },

      setLoading: (loading: boolean) => {
        set({ loading });
      },

      // Selectors
      isAdmin: () => {
        const { userInfo } = get();
        if (!userInfo) return false;

        // Check for admin roles
        const adminRoles = ['admin', 'superAdmin', 'superadmin'];
        return adminRoles.includes(userInfo.role.toLowerCase());
      },

      hasPermission: (permission: string | string[], mode: 'any' | 'all' = 'any') => {
        const { userInfo } = get();
        if (!userInfo) return false;

        // Super admin has all permissions
        if (userInfo.permissions.includes('*')) {
          return true;
        }

        // Normalize to array
        const permissions = Array.isArray(permission) ? permission : [permission];

        // Check permissions based on mode
        if (mode === 'all') {
          return permissions.every(p => userInfo.permissions.includes(p));
        } else {
          return permissions.some(p => userInfo.permissions.includes(p));
        }
      },

      hasAllPermissions: (permissions: string[]) => {
        return get().hasPermission(permissions, 'all');
      },

      hasAnyPermission: (permissions: string[]) => {
        return get().hasPermission(permissions, 'any');
      },

      hasRole: (role: string | string[]) => {
        const { userInfo } = get();
        if (!userInfo) return false;

        const roles = Array.isArray(role) ? role : [role];
        return roles.some(r => r.toLowerCase() === userInfo.role.toLowerCase());
      },
    }),
    {
      name: 'auth-storage',
      // Use sessionStorage (cleared on tab close - security best practice)
      storage: createJSONStorage(() => sessionStorage),
      // Only persist userInfo, token is managed by sessionStorage for API interceptor
      partialize: (state) => ({
        userInfo: state.userInfo,
        isAuthenticated: false, // Force re-validation on page load
        token: state.token,
      }),
    }
  )
);

// Export type for external use
export type { AuthState };

// Selector hooks for common use cases
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useUserInfo = () => useAuthStore((state) => state.userInfo);
export const useIsAdmin = () => useAuthStore((state) => state.isAdmin());
export const useAuthLoading = () => useAuthStore((state) => state.loading);
export const useAuthError = () => useAuthStore((state) => state.error);
