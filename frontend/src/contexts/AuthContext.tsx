import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  useCallback,
} from 'react';
import { STORAGE_KEYS } from '../config';
import type { CurrentUser } from '../types/auth';
import { authService } from '../services/auth';
import { storage } from '../utils/storage';

interface AuthState {
  user: CurrentUser | null;
  token: string | null;
  loading: boolean;
  loginLoading: boolean;
  login: (token: string) => Promise<void>;
  logout: () => void;
}

const Ctx = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [loginLoading, setLoginLoading] = useState(false);

  // 初始化：从 localStorage 恢复 token
  useEffect(() => {
    const t = storage.getItem<string>(STORAGE_KEYS.token);
    setToken(t);

    if (!t) {
      setLoading(false);
      return;
    }

    // 验证 token 并获取用户信息
    authService
      .me()
      .then((u) => setUser(u))
      .catch(() => {
        // Token 无效，清理
        storage.removeItem(STORAGE_KEYS.token);
        setToken(null);
      })
      .finally(() => setLoading(false));
  }, []);

  // 登录方法
  const login = useCallback(async (t: string): Promise<void> => {
    setLoginLoading(true);
    try {
      // 保存 token
      const saved = storage.setItem(STORAGE_KEYS.token, t);
      if (!saved) {
        throw new Error('Failed to save token');
      }
      setToken(t);

      // 获取用户信息
      const u = await authService.me();
      setUser(u);
    } catch (error) {
      // 登录失败，清理状态
      storage.removeItem(STORAGE_KEYS.token);
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
    setToken(null);
    setUser(null);
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
