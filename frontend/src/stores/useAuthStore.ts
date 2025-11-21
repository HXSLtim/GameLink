/**
 * 认证状态管理
 * 使用Zustand实现轻量级状态管理
 */
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../api/types/auth';

/**
 * 认证状态
 */
interface AuthState {
  // 状态
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;

  // 操作
  setAuth: (token: string, user: User) => void;
  clearAuth: () => void;
  updateUser: (user: Partial<User>) => void;
}

/**
 * 认证Store
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 初始状态
      token: null,
      user: null,
      isAuthenticated: false,

      // 设置认证信息
      setAuth: (token, user) => {
        localStorage.setItem('auth_token', token);
        set({
          token,
          user,
          isAuthenticated: true,
        });
      },

      // 清除认证信息
      clearAuth: () => {
        localStorage.removeItem('auth_token');
        set({
          token: null,
          user: null,
          isAuthenticated: false,
        });
      },

      // 更新用户信息
      updateUser: (userUpdate) => {
        const currentUser = get().user;
        if (currentUser) {
          set({
            user: { ...currentUser, ...userUpdate },
          });
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
