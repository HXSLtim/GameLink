// Auth Store - Taro App
// Authentication and authorization state management for Taro mini-program

import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import Taro from '@tarojs/taro';
import type { UserInfo } from '../types';

interface AuthState {
  // State
  token: string | null;
  userInfo: UserInfo | null;
  isAuthenticated: boolean;
  loading: boolean;

  // Actions
  sendCode: (phone: string) => Promise<void>;
  loginWithCode: (phone: string, code: string) => Promise<void>;
  wechatLogin: () => Promise<void>;
  logout: () => void;
  setToken: (token: string) => void;
  setUserInfo: (user: UserInfo) => void;

  // Selectors
  isLoggedIn: () => boolean;
}

// Taro Storage Adapter for zustand persist middleware
const taroStorage = {
  getItem: (name: string): string | null => {
    try {
      const value = Taro.getStorageSync(name);
      return value || null;
    } catch (e) {
      console.error('Taro getStorageSync error:', e);
      return null;
    }
  },
  setItem: (name: string, value: string): void => {
    try {
      Taro.setStorageSync(name, value);
    } catch (e) {
      console.error('Taro setStorageSync error:', e);
    }
  },
  removeItem: (name: string): void => {
    try {
      Taro.removeStorageSync(name);
    } catch (e) {
      console.error('Taro removeStorageSync error:', e);
    }
  },
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial State
      token: null,
      userInfo: null,
      isAuthenticated: false,
      loading: false,

      // Send verification code
      sendCode: async (phone: string) => {
        set({ loading: true });

        try {
          // TODO: Integrate with actual API
          // import { api } from '@/api/client';
          // await api.post('/auth/send-code', { phone });

          console.log('[authStore] Sending verification code to:', phone);

          // Simulate API call
          await new Promise((resolve) => setTimeout(resolve, 500));

          // Show success message
          Taro.showToast({
            title: '验证码已发送',
            icon: 'success',
            duration: 2000,
          });

          set({ loading: false });
        } catch (error: any) {
          console.error('[authStore] Send code error:', error);

          // Show error message
          Taro.showToast({
            title: error.message || '发送失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      // Login with phone and verification code
      loginWithCode: async (phone: string, code: string) => {
        set({ loading: true });

        try {
          // TODO: Integrate with actual API
          // import { api } from '@/api/client';
          // const response = await api.post<LoginResponse>('/auth/login', { phone, code });

          console.log('[authStore] Login with code:', { phone, code });

          // Simulate API call
          await new Promise((resolve) => setTimeout(resolve, 500));

          // Mock user data (replace with actual API response)
          const mockUser: UserInfo = {
            id: Date.now(),
            name: '微信用户',
            phone,
            avatar: '',
            role: 'user',
            createdAt: new Date().toISOString(),
          };

          const mockToken = 'mock-token-' + Date.now();

          set({
            token: mockToken,
            userInfo: mockUser,
            isAuthenticated: true,
            loading: false,
          });

          // Show success message
          Taro.showToast({
            title: '登录成功',
            icon: 'success',
            duration: 2000,
          });
        } catch (error: any) {
          console.error('[authStore] Login error:', error);

          // Show error message
          Taro.showToast({
            title: error.message || '登录失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      // WeChat login flow
      wechatLogin: async () => {
        set({ loading: true });

        try {
          console.log('[authStore] Starting WeChat login flow');

          // TODO: Complete WeChat login flow
          // Step 1: Get wx.login code
          // const loginRes = await Taro.login();
          // console.log('[authStore] wx.login code:', loginRes.code);

          // Step 2: Send code to backend for session
          // import { api } from '@/api/client';
          // const response = await api.post<LoginResponse>('/auth/wechat-login', {
          //   code: loginRes.code,
          // });

          // Simulate WeChat login
          await new Promise((resolve) => setTimeout(resolve, 800));

          const mockUser: UserInfo = {
            id: Date.now(),
            name: '微信用户',
            avatar: '',
            role: 'user',
            createdAt: new Date().toISOString(),
          };

          const mockToken = 'wechat-token-' + Date.now();

          set({
            token: mockToken,
            userInfo: mockUser,
            isAuthenticated: true,
            loading: false,
          });

          // Show success message
          Taro.showToast({
            title: '登录成功',
            icon: 'success',
            duration: 2000,
          });
        } catch (error: any) {
          console.error('[authStore] WeChat login error:', error);

          // Show error message
          Taro.showToast({
            title: error.message || '微信登录失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      // Logout and clear session
      logout: () => {
        set({
          token: null,
          userInfo: null,
          isAuthenticated: false,
        });

        // Clear other stores if needed
        // import { useOrderStore } from './orderStore';
        // useOrderStore.getState().reset();

        // Show logout message
        Taro.showToast({
          title: '已退出登录',
          icon: 'success',
          duration: 1500,
        });

        console.log('[authStore] User logged out');
      },

      // Set token manually (for cases where token is set externally)
      setToken: (token: string) => {
        set({ token, isAuthenticated: !!token });
        console.log('[authStore] Token updated:', !!token);
      },

      // Update user info manually
      setUserInfo: (user: UserInfo) => {
        set({ userInfo: user });
        console.log('[authStore] User info updated:', user.id);
      },

      // Check if user is logged in
      isLoggedIn: () => {
        const { token, userInfo } = get();
        return !!(token && userInfo);
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => taroStorage),
      // Partial persistence: only persist userInfo
      // Token should be obtained from backend on each app start
      partialize: (state) => ({
        userInfo: state.userInfo,
      }),
    }
  )
);

export type { AuthState };
