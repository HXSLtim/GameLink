import React, { createContext, useContext, useEffect, useMemo, useState, useCallback } from 'react';
import { STORAGE_KEYS } from '../config';
import type { CurrentUser } from '../types/auth';
import { storage } from '../utils/storage';

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

  // 登录方法 (Mock 实现)
  const login = useCallback(async (username: string, password: string): Promise<void> => {
    setLoginLoading(true);

    // 模拟网络延迟
    await new Promise((resolve) => setTimeout(resolve, 800));

    try {
      // Mock 登录验证（任何用户名密码都可以登录）
      if (!username || !password) {
        throw new Error('用户名和密码不能为空');
      }

      // 生成 Mock Token
      const mockToken = `mock-token-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

      // 创建 Mock 用户信息
      const mockUser: CurrentUser = {
        id: Math.floor(Math.random() * 1000),
        username: username,
        email: `${username}@gamelink.com`,
        role: username === 'admin' ? 'admin' : 'user',
        avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      // 保存到 localStorage
      storage.setItem(STORAGE_KEYS.token, mockToken);
      storage.setItem(STORAGE_KEYS.user, mockUser);

      // 更新状态
      setToken(mockToken);
      setUser(mockUser);

      console.log('✅ Mock 登录成功:', { username, token: mockToken });
    } catch (error) {
      // 登录失败，清理状态
      storage.removeItem(STORAGE_KEYS.token);
      storage.removeItem(STORAGE_KEYS.user);
      setToken(null);
      setUser(null);
      throw error;
    } finally {
      setLoginLoading(false);
    }
  }, []);

  // 登出方法
  const logout = useCallback(() => {
    storage.removeItem(STORAGE_KEYS.token);
    storage.removeItem(STORAGE_KEYS.user);
    setToken(null);
    setUser(null);
    console.log('✅ 退出登录成功');
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
