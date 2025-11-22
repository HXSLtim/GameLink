/**
 * 认证Hook
 * 提供登录、登出等认证操作
 */
import { useState } from 'react';
import { Message } from '@arco-design/web-react';
import { useNavigate } from 'react-router-dom';
import { authApi } from '@/api';
import { useAuthStore } from '../stores/useAuthStore';
import type { LoginRequest, RegisterRequest } from '@/api';

/**
 * 认证Hook
 */
export const useAuth = () => {
  const navigate = useNavigate();
  const { setAuth, clearAuth, user, isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);

  /**
   * 登录
   */
  const login = async (data: LoginRequest) => {
    try {
      setLoading(true);
      const response = await authApi.login(data);
      setAuth(response.token, response.user);
      Message.success('登录成功');
      navigate('/');
    } catch (error) {
      Message.error('登录失败，请检查用户名和密码');
      throw error;
    } finally {
      setLoading(false);
    }
  };

  /**
   * 注册
   */
  const register = async (data: RegisterRequest) => {
    try {
      setLoading(true);
      await authApi.register(data);
      Message.success('注册成功，请登录');
      navigate('/login');
    } catch (error) {
      Message.error('注册失败');
      throw error;
    } finally {
      setLoading(false);
    }
  };

  /**
   * 登出
   */
  const logout = async () => {
    try {
      await authApi.logout();
      clearAuth();
      Message.success('已退出登录');
      navigate('/login');
    } catch (error) {
      // 即使API失败也清除本地状态
      clearAuth();
      navigate('/login');
    }
  };

  return {
    user,
    isAuthenticated,
    loading,
    login,
    register,
    logout,
  };
};
