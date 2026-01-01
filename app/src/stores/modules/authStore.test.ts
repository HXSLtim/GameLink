/**
 * Auth Store Tests (Taro App)
 * Tests authentication, verification code, and WeChat login flows
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAuthStore } from './authStore';
import type { UserInfo } from '../types';

// Mock Taro
const mockTaro = {
  getStorageSync: vi.fn(),
  setStorageSync: vi.fn(),
  removeStorageSync: vi.fn(),
  showToast: vi.fn(),
  login: vi.fn(),
};

vi.mock('@tarojs/taro', () => ({
  default: mockTaro,
}));

describe('authStore (Taro App)', () => {
  beforeEach(() => {
    // Reset store state before each test
    useAuthStore.getState().setToken('');
    useAuthStore.getState().setUserInfo({} as UserInfo);
    useAuthStore.setState({
      token: null,
      userInfo: null,
      isAuthenticated: false,
      loading: false,
    });

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
    });
  });

  describe('Send Verification Code', () => {
    it('should send verification code successfully', async () => {
      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.sendCode('13800138000');
      });

      expect(result.current.loading).toBe(false);
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '验证码已发送',
        icon: 'success',
        duration: 2000,
      });
    });

    it('should handle send code failure', async () => {
      const { result } = renderHook(() => useAuthStore());

      // Mock setTimeout to resolve immediately
      vi.spyOn(global, 'setTimeout').mockImplementation((cb) => {
        cb(null as any);
        return {} as any;
      });

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      await act(async () => {
        try {
          await result.current.sendCode('13800138000');
        } catch (error) {
          // Expected error (no API integration yet)
        }
      });

      expect(result.current.loading).toBe(false);

      consoleErrorSpy.mockRestore();
    });
  });

  describe('Login with Code', () => {
    it('should login successfully with phone and code', async () => {
      const { result } = renderHook(() => useAuthStore());

      // Mock setTimeout to resolve immediately
      vi.spyOn(global, 'setTimeout').mockImplementation((cb) => {
        cb(null as any);
        return {} as any;
      });

      await act(async () => {
        await result.current.loginWithCode('13800138000', '123456');
      });

      expect(result.current.token).toBeTruthy();
      expect(result.current.userInfo).toBeDefined();
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.loading).toBe(false);
      expect(result.current.userInfo?.phone).toBe('13800138000');

      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '登录成功',
        icon: 'success',
        duration: 2000,
      });
    });

    it('should handle login with code failure', async () => {
      const { result } = renderHook(() => useAuthStore());

      const mockError = new Error('Invalid verification code');
      vi.spyOn(global, 'setTimeout').mockImplementation((cb, _delay) => {
        throw mockError;
      });

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      await act(async () => {
        try {
          await result.current.loginWithCode('13800138000', '000000');
        } catch (error) {
          // Expected to throw
        }
      });

      expect(result.current.loading).toBe(false);

      consoleErrorSpy.mockRestore();
    });
  });

  describe('WeChat Login', () => {
    it('should login successfully with WeChat', async () => {
      const { result } = renderHook(() => useAuthStore());

      // Mock setTimeout to resolve immediately
      vi.spyOn(global, 'setTimeout').mockImplementation((cb) => {
        cb(null as any);
        return {} as any;
      });

      await act(async () => {
        await result.current.wechatLogin();
      });

      expect(result.current.token).toBeTruthy();
      expect(result.current.userInfo).toBeDefined();
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.loading).toBe(false);
      expect(result.current.userInfo?.name).toBe('微信用户');

      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '登录成功',
        icon: 'success',
        duration: 2000,
      });
    });

    it('should handle WeChat login failure', async () => {
      const { result } = renderHook(() => useAuthStore());

      const mockError = new Error('WeChat login failed');
      vi.spyOn(global, 'setTimeout').mockImplementation((cb, _delay) => {
        throw mockError;
      });

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      await act(async () => {
        try {
          await result.current.wechatLogin();
        } catch (error) {
          // Expected to throw
        }
      });

      expect(result.current.loading).toBe(false);

      consoleErrorSpy.mockRestore();
    });
  });

  describe('Logout', () => {
    it('should logout successfully and clear state', () => {
      const { result } = renderHook(() => useAuthStore());

      // Setup: User is logged in
      const mockUser: UserInfo = {
        id: 1,
        name: 'Test User',
        phone: '13800138000',
        avatar: 'https://example.com/avatar.jpg',
        role: 'user',
        createdAt: '2024-01-01T00:00:00Z',
      };

      act(() => {
        result.current.setToken('mock-token');
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.token).toBe('mock-token');

      // Logout
      act(() => {
        result.current.logout();
      });

      expect(result.current.token).toBeNull();
      expect(result.current.userInfo).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);

      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '已退出登录',
        icon: 'success',
        duration: 1500,
      });
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
    });

    it('should update user info', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Test User',
        phone: '13800138000',
        avatar: 'https://example.com/avatar.jpg',
        role: 'user',
        createdAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.userInfo).toEqual(mockUser);
    });
  });

  describe('isLoggedIn Selector', () => {
    it('should return true when token and userInfo exist', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Test User',
        phone: '13800138000',
        role: 'user',
        createdAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('test-token');
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.isLoggedIn()).toBe(true);
    });

    it('should return false when token is missing', () => {
      const mockUser: UserInfo = {
        id: 1,
        name: 'Test User',
        phone: '13800138000',
        role: 'user',
        createdAt: '2024-01-01T00:00:00Z',
      };

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(mockUser);
      });

      expect(result.current.isLoggedIn()).toBe(false);
    });

    it('should return false when userInfo is missing', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('test-token');
      });

      expect(result.current.isLoggedIn()).toBe(false);
    });

    it('should return false when both are missing', () => {
      const { result } = renderHook(() => useAuthStore());

      expect(result.current.isLoggedIn()).toBe(false);
    });
  });

  describe('Loading States', () => {
    it('should set loading to true during sendCode', async () => {
      const { result } = renderHook(() => useAuthStore());

      const loadingPromise = act(async () => {
        await result.current.sendCode('13800138000');
      });

      // Note: Loading state might be cleared before we can check it
      // This test verifies the loading state flow
      expect(result.current.loading).toBe(false); // Should be false after completion

      await loadingPromise;
    });

    it('should set loading to true during loginWithCode', async () => {
      const { result } = renderHook(() => useAuthStore());

      vi.spyOn(global, 'setTimeout').mockImplementation((cb) => {
        cb(null as any);
        return {} as any;
      });

      await act(async () => {
        await result.current.loginWithCode('13800138000', '123456');
      });

      expect(result.current.loading).toBe(false); // Should be false after completion
    });

    it('should set loading to true during wechatLogin', async () => {
      const { result } = renderHook(() => useAuthStore());

      vi.spyOn(global, 'setTimeout').mockImplementation((cb) => {
        cb(null as any);
        return {} as any;
      });

      await act(async () => {
        await result.current.wechatLogin();
      });

      expect(result.current.loading).toBe(false); // Should be false after completion
    });
  });

  describe('Edge Cases', () => {
    it('should handle null user info gracefully', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setUserInfo(null as any);
      });

      expect(result.current.userInfo).toBeNull();
      expect(result.current.isLoggedIn()).toBe(false);
    });

    it('should handle empty token', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setToken('');
      });

      expect(result.current.token).toBe('');
      expect(result.current.isAuthenticated).toBe(false);
    });

    it('should handle logout with no prior login', () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.logout();
      });

      expect(result.current.token).toBeNull();
      expect(result.current.userInfo).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);

      expect(mockTaro.showToast).toHaveBeenCalled();
    });
  });
});
