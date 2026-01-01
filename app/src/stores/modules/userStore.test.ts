/**
 * User Store Tests (Taro App)
 * Tests user profile, wallet, VIP status, and avatar upload
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useUserStore } from './userStore';
import type { UserInfo } from '../types';

// Mock Taro
const mockTaro = {
  getStorageSync: vi.fn(),
  setStorageSync: vi.fn(),
  removeStorageSync: vi.fn(),
  showToast: vi.fn(),
  showLoading: vi.fn(),
  hideLoading: vi.fn(),
  uploadFile: vi.fn(),
};

vi.mock('@tarojs/taro', () => ({
  default: mockTaro,
}));

// Mock API client
vi.mock('../../api/client', () => ({
  get: vi.fn(),
  put: vi.fn(),
}));

const { get, put } = require('../../api/client');

describe('userStore (Taro App)', () => {
  const mockUserInfo: UserInfo = {
    id: 1,
    name: 'Test User',
    phone: '13800138000',
    email: 'test@gamelink.com',
    avatar: 'https://example.com/avatar.jpg',
    role: 'user',
    status: 'active',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  beforeEach(() => {
    // Reset store state before each test
    useUserStore.setState({
      userInfo: null,
      loading: false,
      wallet: {
        balanceCents: 0,
        frozenCents: 0,
        vipLevel: 0,
        vipExpireAt: null,
      },
    });

    // Clear all mocks
    vi.clearAllMocks();

    // Default Taro.getStorageSync mock
    mockTaro.getStorageSync.mockReturnValue(null);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('State Initialization', () => {
    it('should initialize with default state', () => {
      const { result } = renderHook(() => useUserStore());

      expect(result.current.userInfo).toBeNull();
      expect(result.current.loading).toBe(false);
      expect(result.current.wallet).toEqual({
        balanceCents: 0,
        frozenCents: 0,
        vipLevel: 0,
        vipExpireAt: null,
      });
    });
  });

  describe('Fetch User Info', () => {
    it('should fetch user info successfully', async () => {
      const mockResponse = {
        success: true,
        data: mockUserInfo,
      };

      get.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUserInfo();
      });

      expect(result.current.userInfo).toEqual(mockUserInfo);
      expect(result.current.loading).toBe(false);
      expect(get).toHaveBeenCalledWith('/auth/me');
    });

    it('should handle fetch user info failure', async () => {
      const mockError = new Error('Failed to fetch user info');
      get.mockRejectedValue(mockError);

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUserInfo();
      });

      expect(result.current.loading).toBe(false);
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: 'Failed to fetch user info',
        icon: 'none',
        duration: 2000,
      });

      consoleErrorSpy.mockRestore();
    });

    it('should handle unsuccessful response', async () => {
      const mockResponse = {
        success: false,
        message: 'User not found',
      };

      get.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUserInfo();
      });

      expect(result.current.loading).toBe(false);
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: 'User not found',
        icon: 'none',
        duration: 2000,
      });
    });
  });

  describe('Update Profile', () => {
    it('should update profile successfully', async () => {
      const updatedData = {
        name: 'Updated User',
        email: 'updated@gamelink.com',
      };

      const mockResponse = {
        success: true,
        data: updatedData,
      };

      put.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      // Set initial user info
      act(() => {
        result.current.userInfo = mockUserInfo;
      });

      await act(async () => {
        await result.current.updateProfile(updatedData);
      });

      expect(result.current.userInfo).toEqual({
        ...mockUserInfo,
        ...updatedData,
      });
      expect(result.current.loading).toBe(false);
      expect(put).toHaveBeenCalledWith('/user/profile', updatedData);
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '更新成功',
        icon: 'success',
        duration: 2000,
      });
    });

    it('should handle update profile failure', async () => {
      const mockError = new Error('Update failed');
      put.mockRejectedValue(mockError);

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        try {
          await result.current.updateProfile({ name: 'Updated' });
        } catch (error) {
          // Expected to throw
        }
      });

      expect(result.current.loading).toBe(false);
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: 'Update failed',
        icon: 'none',
        duration: 2000,
      });

      consoleErrorSpy.mockRestore();
    });

    it('should update profile when userInfo is null', async () => {
      const updatedData = {
        name: 'New User',
      };

      const mockResponse = {
        success: true,
        data: {
          ...mockUserInfo,
          ...updatedData,
        },
      };

      put.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.updateProfile(updatedData);
      });

      expect(result.current.userInfo).toEqual(mockResponse.data);
    });
  });

  describe('Upload Avatar', () => {
    it('should upload avatar successfully', async () => {
      const filePath = '/tmp/avatar.jpg';
      const avatarUrl = 'https://example.com/uploaded-avatar.jpg';

      const mockResponse = {
        data: JSON.stringify({
          success: true,
          data: {
            url: avatarUrl,
          },
        }),
      };

      mockTaro.uploadFile.mockResolvedValue(mockResponse);
      mockTaro.getStorageSync.mockReturnValue('mock-token');

      const { result } = renderHook(() => useUserStore());

      // Set initial user info
      act(() => {
        result.current.userInfo = mockUserInfo;
      });

      await act(async () => {
        const uploadedUrl = await result.current.uploadAvatar(filePath);
        expect(uploadedUrl).toBe(avatarUrl);
      });

      expect(result.current.userInfo?.avatar).toBe(avatarUrl);
      expect(mockTaro.uploadFile).toHaveBeenCalledWith({
        url: expect.stringContaining('/upload/avatar'),
        filePath,
        name: 'avatar',
        header: {
          Authorization: 'Bearer mock-token',
        },
      });
      expect(mockTaro.showLoading).toHaveBeenCalled();
      expect(mockTaro.hideLoading).toHaveBeenCalled();
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: '上传成功',
        icon: 'success',
        duration: 2000,
      });
    });

    it('should handle upload avatar failure', async () => {
      const filePath = '/tmp/avatar.jpg';
      const mockError = new Error('Upload failed');

      mockTaro.uploadFile.mockRejectedValue(mockError);

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          try {
            await result.current.uploadAvatar(filePath);
          } catch (error) {
            // Expected to throw
          }
        })
      ).rejects.toThrow();

      expect(mockTaro.hideLoading).toHaveBeenCalled();
      expect(mockTaro.showToast).toHaveBeenCalledWith({
        title: 'Upload failed',
        icon: 'none',
        duration: 2000,
      });

      consoleErrorSpy.mockRestore();
    });

    it('should handle upload response with no data', async () => {
      const filePath = '/tmp/avatar.jpg';

      const mockResponse = {
        data: JSON.stringify({
          success: false,
          message: 'Invalid file',
        }),
      };

      mockTaro.uploadFile.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          try {
            await result.current.uploadAvatar(filePath);
          } catch (error) {
            // Expected to throw
          }
        })
      ).rejects.toThrow();
    });

    it('should update avatar when userInfo is null', async () => {
      const filePath = '/tmp/avatar.jpg';
      const avatarUrl = 'https://example.com/uploaded-avatar.jpg';

      const mockResponse = {
        data: JSON.stringify({
          success: true,
          data: {
            url: avatarUrl,
          },
        }),
      };

      mockTaro.uploadFile.mockResolvedValue(mockResponse);
      mockTaro.getStorageSync.mockReturnValue('mock-token');

      const { result } = renderHook(() => useUserStore());

      // userInfo is null
      expect(result.current.userInfo).toBeNull();

      await act(async () => {
        await result.current.uploadAvatar(filePath);
      });

      // userInfo should still be null (no user to update)
      expect(result.current.userInfo).toBeNull();
    });
  });

  describe('Fetch Wallet', () => {
    it('should fetch wallet successfully', async () => {
      const mockWalletData = {
        balanceCents: 10000, // 100 yuan
        frozenCents: 5000, // 50 yuan
      };

      const mockResponse = {
        success: true,
        data: mockWalletData,
      };

      get.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchWallet();
      });

      expect(result.current.wallet.balanceCents).toBe(10000);
      expect(result.current.wallet.frozenCents).toBe(5000);
      expect(get).toHaveBeenCalledWith('/user/wallet/balance');
    });

    it('should handle fetch wallet failure silently', async () => {
      const mockError = new Error('Failed to fetch wallet');
      get.mockRejectedValue(mockError);

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchWallet();
      });

      // Should not show toast for wallet errors
      expect(mockTaro.showToast).not.toHaveBeenCalled();

      // Wallet should remain unchanged
      expect(result.current.wallet.balanceCents).toBe(0);

      consoleErrorSpy.mockRestore();
    });

    it('should preserve VIP status when fetching wallet', async () => {
      const { result } = renderHook(() => useUserStore());

      // Set initial VIP status
      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 3,
          vipExpireAt: '2024-12-31T23:59:59Z',
        };
      });

      const mockWalletData = {
        balanceCents: 10000,
        frozenCents: 5000,
      };

      const mockResponse = {
        success: true,
        data: mockWalletData,
      };

      get.mockResolvedValue(mockResponse);

      await act(async () => {
        await result.current.fetchWallet();
      });

      // VIP status should be preserved
      expect(result.current.wallet.vipLevel).toBe(3);
      expect(result.current.wallet.vipExpireAt).toBe('2024-12-31T23:59:59Z');
      // Balance should be updated
      expect(result.current.wallet.balanceCents).toBe(10000);
    });
  });

  describe('VIP Status Selectors', () => {
    it('should return false for non-VIP user', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 0,
          vipExpireAt: null,
        };
      });

      expect(result.current.isVip()).toBe(false);
    });

    it('should return true for active VIP user', () => {
      const { result } = renderHook(() => useUserStore());

      // Set VIP to expire in the future
      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 30);

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 3,
          vipExpireAt: futureDate.toISOString(),
        };
      });

      expect(result.current.isVip()).toBe(true);
    });

    it('should return false for expired VIP user', () => {
      const { result } = renderHook(() => useUserStore());

      // Set VIP to expire in the past
      const pastDate = new Date();
      pastDate.setDate(pastDate.getDate() - 10);

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 3,
          vipExpireAt: pastDate.toISOString(),
        };
      });

      expect(result.current.isVip()).toBe(false);
    });

    it('should calculate remaining VIP days correctly', () => {
      const { result } = renderHook(() => useUserStore());

      // Set VIP to expire in 7 days
      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 7);

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 3,
          vipExpireAt: futureDate.toISOString(),
        };
      });

      expect(result.current.vipDaysLeft()).toBe(7);
    });

    it('should return 0 days for expired VIP', () => {
      const { result } = renderHook(() => useUserStore());

      // Set VIP to expire in the past
      const pastDate = new Date();
      pastDate.setDate(pastDate.getDate() - 10);

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 3,
          vipExpireAt: pastDate.toISOString(),
        };
      });

      expect(result.current.vipDaysLeft()).toBe(0);
    });

    it('should return 0 days when vipExpireAt is null', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 0,
          vipExpireAt: null,
        };
      });

      expect(result.current.vipDaysLeft()).toBe(0);
    });
  });

  describe('Balance Selectors', () => {
    it('should convert balance cents to yuan', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.wallet = {
          balanceCents: 10000, // 100 yuan
          frozenCents: 5000, // 50 yuan
          vipLevel: 0,
          vipExpireAt: null,
        };
      });

      expect(result.current.balanceYuan()).toBe(100);
    });

    it('should convert frozen balance cents to yuan', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.wallet = {
          balanceCents: 10000,
          frozenCents: 5000, // 50 yuan
          vipLevel: 0,
          vipExpireAt: null,
        };
      });

      expect(result.current.frozenBalanceYuan()).toBe(50);
    });

    it('should handle zero balance', () => {
      const { result } = renderHook(() => useUserStore());

      expect(result.current.balanceYuan()).toBe(0);
      expect(result.current.frozenBalanceYuan()).toBe(0);
    });
  });

  describe('Edge Cases', () => {
    it('should handle null userInfo gracefully', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.userInfo = null;
      });

      expect(result.current.userInfo).toBeNull();
    });

    it('should handle updateProfile with null userInfo', async () => {
      const mockResponse = {
        success: true,
        data: {
          id: 1,
          name: 'New User',
          createdAt: '2024-01-01T00:00:00Z',
        },
      };

      put.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUserStore());

      // userInfo is null
      expect(result.current.userInfo).toBeNull();

      await act(async () => {
        await result.current.updateProfile({ name: 'New User' });
      });

      // userInfo should be set to the response data
      expect(result.current.userInfo).toEqual(mockResponse.data);
    });

    it('should handle VIP level 0 with future expiration', () => {
      const { result } = renderHook(() => useUserStore());

      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 30);

      act(() => {
        result.current.wallet = {
          balanceCents: 0,
          frozenCents: 0,
          vipLevel: 0, // Not a VIP
          vipExpireAt: futureDate.toISOString(),
        };
      });

      expect(result.current.isVip()).toBe(false);
    });

    it('should handle very large balance values', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.wallet = {
          balanceCents: 99999999, // Almost 1 million yuan
          frozenCents: 50000000,
          vipLevel: 0,
          vipExpireAt: null,
        };
      });

      expect(result.current.balanceYuan()).toBe(999999.99);
      expect(result.current.frozenBalanceYuan()).toBe(500000);
    });

    it('should handle wallet fetch with empty response', async () => {
      const mockResponse = {
        success: false,
      };

      get.mockResolvedValue(mockResponse);

      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchWallet();
      });

      // Wallet should remain unchanged
      expect(result.current.wallet.balanceCents).toBe(0);

      consoleErrorSpy.mockRestore();
    });
  });
});
