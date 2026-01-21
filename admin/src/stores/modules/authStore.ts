/**
 * Auth Store - Authentication and Authorization State Management
 *
 * Features:
 * - Token storage in localStorage (persistent across sessions)
 * - Permission checking methods (hasPermission, isAdmin)
 * - Logout clears all stores
 * - Integrated with auth API
 * - Auto-restores state from localStorage on page refresh
 */

import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { UserInfo, LoginRequest } from '../types';
import { authApi } from '@/api/auth';
import { permissionApi } from '@/api/permission';

import { logger } from '@/utils/logger';
interface AuthState {
  // State
  token: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  /** 水合状态标志: Zustand persist 是否已完成从 localStorage 恢复状态 */
  _hydrated: boolean;

  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  setToken: (token: string) => void;
  setUserInfo: (user: UserInfo) => void;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
  refreshPermissions: () => Promise<void>;

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
      _hydrated: false, // 初始为 false，水合完成后设为 true

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

          // Fetch user permissions
          let permissions: string[] = [];
          try {
            const permResponse = await permissionApi.getMyPermissions();
            if (permResponse.data?.success && permResponse.data?.data) {
              permissions = permResponse.data.data;
            }
          } catch (permError) {
            logger.warn('[authStore] Failed to fetch permissions, continuing with empty permissions:', permError);
          }

          // Map API user response to UserInfo interface
          const userInfo: UserInfo = {
            id: user.id,
            name: user.username,
            email: user.email,
            phone: undefined,
            avatar: undefined,
            role: user.role,
            permissions,
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

          logger.info('[authStore] Login successful, permissions:', permissions);

          // 不再手动写入 localStorage，由 persist 中间件自动处理

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

          // 不再手动清理 localStorage，persist 中间件会自动处理
          // 只需要清空状态，中间件会同步到 localStorage

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
        }
      },

      setToken: (token: string) => {
        const isAuthenticated = !!token;
        set({ token, isAuthenticated });
        // 不再手动写入 localStorage，由 persist 中间件自动处理
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

      refreshPermissions: async () => {
        const { userInfo } = get();
        if (!userInfo) return;

        try {
          const permResponse = await permissionApi.getMyPermissions();
          if (permResponse.data?.success && permResponse.data?.data) {
            set({
              userInfo: {
                ...userInfo,
                permissions: permResponse.data.data,
                updatedAt: new Date().toISOString(),
              },
            });
            logger.info('[authStore] Permissions refreshed:', permResponse.data.data);
          }
        } catch (error) {
          logger.error('[authStore] Failed to refresh permissions:', error);
        }
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
      // Use localStorage for persistent auth state
      storage: createJSONStorage(() => localStorage),
      // Persist token and userInfo
      partialize: (state) => ({
        userInfo: state.userInfo,
        isAuthenticated: state.isAuthenticated,
        token: state.token,
        // _hydrated 不持久化，每次刷新都从 false 开始
      }),
      // 水合完成回调
      onRehydrateStorage: () => (state) => {
        logger.info('[authStore] Rehydration complete');
        if (state) {
          state._hydrated = true;
        }
      },
    }
  )
);

// Export type for external use
export type { AuthState };

// Selector hooks for common use cases (Super Dev 最佳实践: 精确订阅)
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useUserInfo = () => useAuthStore((state) => state.userInfo);
export const useIsAdmin = () => useAuthStore((state) => state.isAdmin());
export const useAuthLoading = () => useAuthStore((state) => state.loading);
export const useAuthError = () => useAuthStore((state) => state.error);
export const useAuthToken = () => useAuthStore((state) => state.token);
/** 新增: 水合状态 hook，用于检查 Zustand persist 是否已完成 */
export const useIsHydrated = () => useAuthStore((state) => state._hydrated);
