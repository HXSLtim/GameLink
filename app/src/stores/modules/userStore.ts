// User Store - Taro App
// User profile, wallet, and VIP state management

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import Taro from '@tarojs/taro';
import { get, put } from '../../api/client';
import type { UserInfo } from '../types';

/**
 * User wallet information
 */
interface WalletInfo {
  balanceCents: number;       // Available balance (in cents/分)
  frozenCents: number;        // Frozen balance (in cents/分)
  vipLevel: number;           // VIP level (0 = non-VIP)
  vipExpireAt: string | null; // VIP expiration time (ISO 8601)
}

/**
 * User store state and actions
 */
interface UserState {
  // State
  userInfo: UserInfo | null;
  loading: boolean;

  // Wallet
  wallet: WalletInfo;

  // Actions
  fetchUserInfo: () => Promise<void>;
  updateProfile: (data: Partial<UserInfo>) => Promise<void>;
  uploadAvatar: (filePath: string) => Promise<string>;
  fetchWallet: () => Promise<void>;

  // Selectors (computed values)
  isVip: () => boolean;
  vipDaysLeft: () => number;
  balanceYuan: () => number;      // Convert cents to yuan
  frozenBalanceYuan: () => number; // Convert frozen cents to yuan
}

/**
 * Initial wallet state
 */
const initialWallet: WalletInfo = {
  balanceCents: 0,
  frozenCents: 0,
  vipLevel: 0,
  vipExpireAt: null,
};

/**
 * User store with persistence
 */
export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // Initial state
      userInfo: null,
      loading: false,
      wallet: initialWallet,

      /**
       * Fetch current user information from API
       * GET /api/v1/auth/me
       */
      fetchUserInfo: async () => {
        set({ loading: true });

        try {
          const response = await get<any>('/auth/me');

          if (response.success && response.data) {
            set({
              userInfo: response.data,
              loading: false,
            });
          } else {
            throw new Error(response.message || 'Failed to fetch user info');
          }
        } catch (error: any) {
          console.error('Fetch user info error:', error);

          // Show error toast
          Taro.showToast({
            title: error.message || '获取用户信息失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
        }
      },

      /**
       * Update user profile
       * PUT /api/v1/user/profile
       */
      updateProfile: async (data: Partial<UserInfo>) => {
        set({ loading: true });

        try {
          const response = await put<any>('/user/profile', data);

          if (response.success && response.data) {
            set(state => ({
              userInfo: state.userInfo
                ? { ...state.userInfo, ...response.data }
                : response.data,
              loading: false,
            }));

            // Show success toast
            Taro.showToast({
              title: '更新成功',
              icon: 'success',
              duration: 2000,
            });
          } else {
            throw new Error(response.message || 'Failed to update profile');
          }
        } catch (error: any) {
          console.error('Update profile error:', error);

          // Show error toast
          Taro.showToast({
            title: error.message || '更新失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * Upload avatar and update user avatar URL
       * POST /api/v1/upload/avatar
       */
      uploadAvatar: async (filePath: string) => {
        try {
          // Show loading
          Taro.showLoading({
            title: '上传中...',
            mask: true,
          });

          // Get token from storage
          const token = Taro.getStorageSync('token');

          // Upload file using Taro.uploadFile
          const uploadRes = await Taro.uploadFile({
            url: `${process.env.TARO_APP_API_BASE_URL || 'http://localhost:8080/api/v1'}/upload/avatar`,
            filePath,
            name: 'avatar',
            header: {
              Authorization: `Bearer ${token}`,
            },
          });

          // Hide loading
          Taro.hideLoading();

          // Parse response
          const response = JSON.parse(uploadRes.data);

          if (response.success && response.data?.url) {
            const avatarUrl = response.data.url;

            // Update user avatar in store
            set(state => ({
              userInfo: state.userInfo
                ? { ...state.userInfo, avatar: avatarUrl }
                : null,
            }));

            // Show success toast
            Taro.showToast({
              title: '上传成功',
              icon: 'success',
              duration: 2000,
            });

            return avatarUrl;
          } else {
            throw new Error(response.message || 'Upload failed');
          }
        } catch (error: any) {
          console.error('Upload avatar error:', error);

          // Hide loading
          Taro.hideLoading();

          // Show error toast
          Taro.showToast({
            title: error.message || '上传失败',
            icon: 'none',
            duration: 2000,
          });

          throw error;
        }
      },

      /**
       * Fetch wallet balance and VIP status
       * GET /api/v1/user/wallet/balance
       */
      fetchWallet: async () => {
        try {
          const response = await get<any>('/user/wallet/balance');

          if (response.success && response.data) {
            // Update wallet with API response
            set(state => ({
              wallet: {
                ...state.wallet,
                balanceCents: response.data.balanceCents,
                frozenCents: response.data.frozenCents,
              },
            }));
          }
        } catch (error: any) {
          console.error('Fetch wallet error:', error);

          // Silently fail for wallet errors (non-critical)
          // Don't show toast to avoid annoying users
        }
      },

      // Selectors (computed values)

      /**
       * Check if user has active VIP status
       */
      isVip: () => {
        const { wallet } = get();
        if (wallet.vipLevel === 0) return false;
        if (!wallet.vipExpireAt) return false;

        const expireTime = new Date(wallet.vipExpireAt).getTime();
        return expireTime > Date.now();
      },

      /**
       * Get remaining VIP days
       */
      vipDaysLeft: () => {
        const { wallet } = get();
        if (!wallet.vipExpireAt) return 0;

        const expireTime = new Date(wallet.vipExpireAt).getTime();
        const daysLeft = Math.ceil((expireTime - Date.now()) / (24 * 60 * 60 * 1000));
        return Math.max(0, daysLeft);
      },

      /**
       * Get available balance in yuan (元)
       */
      balanceYuan: () => {
        const { wallet } = get();
        return wallet.balanceCents / 100;
      },

      /**
       * Get frozen balance in yuan (元)
       */
      frozenBalanceYuan: () => {
        const { wallet } = get();
        return wallet.frozenCents / 100;
      },
    }),
    {
      name: 'user-storage',
      // Partialize: only persist basic user info and VIP status
      // Don't persist: loading state, balance (fetch from server)
      partialize: (state) => ({
        userInfo: state.userInfo,
        wallet: {
          balanceCents: 0,           // Don't cache balance
          frozenCents: 0,            // Don't cache frozen balance
          vipLevel: state.wallet.vipLevel,
          vipExpireAt: state.wallet.vipExpireAt,
        },
      }),
    }
  )
);

export type { UserState, WalletInfo };
