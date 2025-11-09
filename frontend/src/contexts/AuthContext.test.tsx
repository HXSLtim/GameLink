import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';
import { authApi } from '../services/api/auth';

// Mock auth API
vi.mock('../services/api/auth', () => ({
  authApi: {
    login: vi.fn(),
    logout: vi.fn(),
    getCurrentUser: vi.fn(),
    refresh: vi.fn(),
  },
}));

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <AuthProvider>{children}</AuthProvider>
);

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('should throw error when useAuth is used outside AuthProvider', () => {
    // Suppress console.error for this test
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      renderHook(() => useAuth());
    }).toThrow('useAuth must be used within AuthProvider');

    consoleSpy.mockRestore();
  });

  it('should initialize with null user and token when no token in localStorage', async () => {
    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toBeNull();
    expect(result.current.token).toBeNull();
  });

  it('should load user from localStorage on mount', async () => {
    const mockUser = {
      id: 1,
      name: 'testuser',
      username: 'testuser',
      role: 'admin' as const,
      status: 'active' as const,
    };

    // 设置localStorage中的token和user
    localStorage.setItem('gamelink_token', 'test-token');
    localStorage.setItem('gamelink_user', JSON.stringify(mockUser));

    const { result } = renderHook(() => useAuth(), { wrapper });

    // Initially loading
    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.token).toBe('test-token');
    expect(result.current.user).toEqual(mockUser);
    // AuthContext不会在初始化时调用API，只从localStorage读取
    expect(authApi.getCurrentUser).not.toHaveBeenCalled();
  });

  it('should load invalid token from localStorage but not clear it automatically', async () => {
    const invalidUser = { id: 0, name: 'invalid', role: 'user' as const, status: 'active' as const };
    localStorage.setItem('gamelink_token', 'invalid-token');
    localStorage.setItem('gamelink_user', JSON.stringify(invalidUser));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // AuthContext会从localStorage读取，即使数据无效也不会自动清理
    // 只有在登录失败时才会清理
    expect(result.current.token).toBe('invalid-token');
    expect(result.current.user).toEqual(invalidUser);
  });

  it('should login and store token', async () => {
    const mockUser = {
      id: 1,
      name: 'testuser',
      username: 'testuser',
      role: 'admin' as const,
      status: 'active' as const,
    };

    const mockLoginResult = {
      token: 'new-token',
      user: mockUser,
    };

    vi.mocked(authApi.login).mockResolvedValue(mockLoginResult);

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.login('testuser', 'password123');
    });

    expect(localStorage.getItem('gamelink_token')).toBe('new-token');
    expect(result.current.token).toBe('new-token');
    expect(result.current.user).toEqual(mockUser);
    expect(authApi.login).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123',
    });
  });

  it('should logout and clear token', async () => {
    const mockUser = {
      id: 1,
      name: 'testuser',
      username: 'testuser',
      role: 'admin' as const,
      status: 'active' as const,
    };
    localStorage.setItem('gamelink_token', 'test-token');
    localStorage.setItem('gamelink_user', JSON.stringify(mockUser));

    vi.mocked(authApi.logout).mockResolvedValue();

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.user).not.toBeNull();
    });

    await act(async () => {
      await result.current.logout();
    });

    expect(result.current.token).toBeNull();
    expect(result.current.user).toBeNull();
    expect(localStorage.getItem('gamelink_token')).toBeNull();
    expect(authApi.logout).toHaveBeenCalled();
  });

  it('should handle login failure gracefully', async () => {
    vi.mocked(authApi.login).mockRejectedValue(new Error('Invalid credentials'));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      try {
        await result.current.login('testuser', 'wrongpassword');
      } catch (error) {
        // 预期会抛出错误
        expect(error).toBeDefined();
      }
    });

    // 登录失败后，token和user应该被清理
    expect(result.current.token).toBeNull();
    expect(result.current.user).toBeNull();
    expect(localStorage.getItem('gamelink_token')).toBeNull();
  });
});
