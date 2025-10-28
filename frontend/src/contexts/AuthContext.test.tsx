import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';
import * as authService from '../services/auth';

// Mock auth service
vi.mock('../services/auth', () => ({
  authService: {
    me: vi.fn(),
    login: vi.fn(),
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

  it('should load user from localStorage token on mount', async () => {
    const mockUser = {
      id: 1,
      username: 'testuser',
      role: 'admin',
    };

    localStorage.setItem('gamelink_token', 'test-token');
    vi.mocked(authService.authService.me).mockResolvedValue(mockUser as any);

    const { result } = renderHook(() => useAuth(), { wrapper });

    // Initially loading
    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.token).toBe('test-token');
    expect(result.current.user).toEqual(mockUser);
    expect(authService.authService.me).toHaveBeenCalled();
  });

  it('should clear token when me() fails', async () => {
    localStorage.setItem('gamelink_token', 'invalid-token');
    vi.mocked(authService.authService.me).mockRejectedValue(new Error('Unauthorized'));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.token).toBeNull();
    expect(result.current.user).toBeNull();
    expect(localStorage.getItem('gamelink_token')).toBeNull();
  });

  it('should login and store token', async () => {
    const mockUser = {
      id: 1,
      username: 'testuser',
      role: 'admin',
    };

    vi.mocked(authService.authService.me).mockResolvedValue(mockUser as any);

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.login('new-token');
    });

    expect(localStorage.getItem('gamelink_token')).toBe('new-token');
    expect(result.current.token).toBe('new-token');

    await waitFor(() => {
      expect(result.current.user).toEqual(mockUser);
    });
  });

  it('should logout and clear token', async () => {
    localStorage.setItem('gamelink_token', 'test-token');
    vi.mocked(authService.authService.me).mockResolvedValue({
      id: 1,
      username: 'testuser',
      role: 'admin',
    } as any);

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.user).not.toBeNull();
    });

    act(() => {
      result.current.logout();
    });

    expect(result.current.token).toBeNull();
    expect(result.current.user).toBeNull();
    expect(localStorage.getItem('gamelink_token')).toBeNull();
  });

  it('should handle login failure gracefully', async () => {
    vi.mocked(authService.authService.me).mockRejectedValue(new Error('Failed to fetch user'));

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.login('test-token');
    });

    await waitFor(() => {
      // Token should be stored even if me() fails
      expect(result.current.token).toBe('test-token');
      // But user should remain null
      expect(result.current.user).toBeNull();
    });
  });
});
