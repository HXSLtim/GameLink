/**
 * User Store
 * 用户状态管理
 */

import Taro from '@tarojs/taro';
import { create } from 'zustand';
import type {
  User,
  UserRole,
  Wallet,
  VipLevelInfo,
} from '@/types/api';

interface UserState {
  // 状态
  currentUser: User | null;
  currentRole: UserRole;
  token: string | null;
  isLoggedIn: boolean;
  isPlayer: boolean;

  // VIP 信息
  vipLevel?: VipLevelInfo;
  vipExp: number;

  // 钱包信息
  wallet?: Wallet;

  // Actions
  setToken: (token: string | null) => void;
  setCurrentUser: (user: User | null) => void;
  switchRole: (role: UserRole) => void;
  logout: () => void;
  updateWallet: (wallet: Wallet) => void;
  updateVipInfo: (vipLevel: VipLevelInfo, vipExp: number) => void;

  // 辅助方法
  hasRole: (role: UserRole) => boolean;
  canSwitchTo: (role: UserRole) => boolean;
}

export const useUserStore = create<UserState>((set, get) => ({
  // 初始状态
  currentUser: null,
  currentRole: 'user',
  token: null,
  isLoggedIn: false,
  isPlayer: false,
  vipExp: 0,

  // 设置 Token
  setToken: (token) => {
    set({ token, isLoggedIn: !!token });
    // 持久化到本地存储
    if (token) {
      Taro.setStorageSync('token', token);
    } else {
      Taro.removeStorageSync('token');
    }
  },

  // 设置当前用户
  setCurrentUser: (user) => {
    const hasPlayerRole = user?.roles?.some(r => r.slug === 'player') || user?.role === 'player';
    set({
      currentUser: user,
      currentRole: user?.role || 'user',
      isPlayer: hasPlayerRole,
      vipExp: user?.vipExp || 0,
      vipLevel: user?.vipLevel,
      wallet: user?.wallet,
    });

    // 持久化用户信息
    if (user) {
      Taro.setStorageSync('userInfo', user);
    } else {
      Taro.removeStorageSync('userInfo');
    }
  },

  // 切换角色
  switchRole: (role) => {
    const { currentUser } = get();
    if (!currentUser) return;

    // 验证用户是否有该角色
    const hasRole = currentUser.roles?.some(r => r.slug === role) || currentUser.role === role;
    if (!hasRole) {
      console.warn(`User does not have role: ${role}`);
      return;
    }

    set({ currentRole: role, isPlayer: role === 'player' });
  },

  // 登出
  logout: () => {
    set({
      currentUser: null,
      token: null,
      isLoggedIn: false,
      currentRole: 'user',
      isPlayer: false,
      wallet: undefined,
      vipLevel: undefined,
      vipExp: 0,
    });

    // 清除本地存储
    Taro.removeStorageSync('token');
    Taro.removeStorageSync('userInfo');
  },

  // 更新钱包信息
  updateWallet: (wallet) => {
    set({ wallet });
    // 同时更新 currentUser 中的 wallet
    const { currentUser } = get();
    if (currentUser) {
      set({ currentUser: { ...currentUser, wallet } });
    }
  },

  // 更新 VIP 信息
  updateVipInfo: (vipLevel, vipExp) => {
    set({ vipLevel, vipExp });
    // 同时更新 currentUser 中的信息
    const { currentUser } = get();
    if (currentUser) {
      set({
        currentUser: {
          ...currentUser,
          vipLevel,
          vipExp,
        },
      });
    }
  },

  // 检查是否有指定角色
  hasRole: (role) => {
    const { currentUser, currentRole } = get();
    if (!currentUser) return false;
    return currentRole === role ||
           currentUser.role === role ||
           currentUser.roles?.some(r => r.slug === role);
  },

  // 检查是否可以切换到指定角色
  canSwitchTo: (role) => {
    const { currentUser } = get();
    if (!currentUser) return false;
    return currentUser.roles?.some(r => r.slug === role) || currentUser.role === role;
  },
}));

// 初始化：从本地存储恢复状态
export const initializeAuthState = () => {
  const token = Taro.getStorageSync('token');
  const userInfo = Taro.getStorageSync('userInfo');

  if (token && userInfo) {
    useUserStore.getState().setToken(token);
    useUserStore.getState().setCurrentUser(userInfo);
  }
};

// 导出辅助 hooks
export const useAuth = () => useUserStore();
export const useIsLoggedIn = () => useUserStore((state) => state.isLoggedIn);
export const useCurrentUser = () => useUserStore((state) => state.currentUser);
export const useIsPlayer = () => useUserStore((state) => state.isPlayer);
export const useCurrentRole = () => useUserStore((state) => state.currentRole);
export const useWallet = () => useUserStore((state) => state.wallet);
export const useVipInfo = () => useUserStore((state) => ({
  vipLevel: state.vipLevel,
  vipExp: state.vipExp,
}));
