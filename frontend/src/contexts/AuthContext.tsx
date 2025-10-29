import React, { createContext, useContext, useEffect, useMemo, useState, useCallback } from 'react';
import { STORAGE_KEYS } from '../config';
import type { CurrentUser } from '../types/auth';
import { UserRole, UserStatus } from '../types/user';
import { storage } from '../utils/storage';
import { authApi } from '../services/api/auth';

interface AuthState {
  user: CurrentUser | null;
  token: string | null;
  loading: boolean;
  loginLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const Ctx = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [loginLoading, setLoginLoading] = useState(false);

  // 初始化：从 localStorage 恢复 token 和用户信息
  useEffect(() => {
    const t = storage.getItem<string>(STORAGE_KEYS.token);
    const u = storage.getItem<CurrentUser>(STORAGE_KEYS.user);

    if (t && u) {
      setToken(t);
      setUser(u);
    }

    setLoading(false);
  }, []);

  // 登录方法 (真实API实现)
  const login = useCallback(async (username: string, password: string): Promise<void> => {
    setLoginLoading(true);

    try {
      // 调用真实的登录API
      const response = await authApi.login({
        username, // 后端使用 username 字段
        password,
      });

      // 保存 Token
      storage.setItem(STORAGE_KEYS.token, response.token);

      // 转换用户信息格式（兼容前端现有结构）
      const currentUser: CurrentUser = {
        id: response.user?.id || 0,
        name: response.user?.name || username,
        email: response.user?.email,
        phone: response.user?.phone,
        avatarUrl: response.user?.avatarUrl,
        role: response.user?.role || ('user' as UserRole),
        status: response.user?.status || ('active' as UserStatus),
        lastLoginAt: response.user?.lastLoginAt,
        createdAt: response.user?.createdAt,
        updatedAt: response.user?.updatedAt,
        // 兼容字段
        username: response.user?.name || username,
        avatar: response.user?.avatarUrl,
      };

      // 保存用户信息
      storage.setItem(STORAGE_KEYS.user, currentUser);

      // 更新状态
      setToken(response.token);
      setUser(currentUser);

      console.log('✅ 登录成功:', {
        username: currentUser.name,
        role: currentUser.role,
        token: response.token.substring(0, 20) + '...',
      });
    } catch (error) {
      // 登录失败，清理状态
      storage.removeItem(STORAGE_KEYS.token);
      storage.removeItem(STORAGE_KEYS.user);
      setToken(null);
      setUser(null);

      console.error('❌ 登录失败:', error);
      throw error;
    } finally {
      setLoginLoading(false);
    }
  }, []);

  // 登出方法
  const logout = useCallback(async () => {
    try {
      // 调用真实的登出API（即使失败也继续清理本地状态）
      await authApi.logout().catch((err) => {
        console.warn('登出API调用失败（已忽略）:', err);
      });
    } finally {
      // 清理本地状态
      storage.removeItem(STORAGE_KEYS.token);
      storage.removeItem(STORAGE_KEYS.user);
      setToken(null);
      setUser(null);
      console.log('✅ 退出登录成功');
    }
  }, []);

  const value = useMemo<AuthState>(
    () => ({
      user,
      token,
      loading,
      loginLoading,
      login,
      logout,
    }),
    [user, token, loading, loginLoading, login, logout],
  );

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
