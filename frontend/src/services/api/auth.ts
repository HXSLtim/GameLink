import { apiClient } from '../../api/client';
import type { LoginRequest, LoginResult, CurrentUser } from '../../types/auth';

/**
 * 认证API服务
 */
export const authApi = {
  /**
   * 用户登录
   */
  login: (data: LoginRequest): Promise<LoginResult> => {
    return apiClient.post('/api/v1/auth/login', data);
  },

  /**
   * 刷新令牌
   */
  refresh: (): Promise<LoginResult> => {
    return apiClient.post('/api/v1/auth/refresh', {});
  },

  /**
   * 用户登出
   */
  logout: (): Promise<void> => {
    return apiClient.post('/api/v1/auth/logout', {});
  },

  /**
   * 获取当前用户信息
   */
  getCurrentUser: (): Promise<CurrentUser> => {
    return apiClient.get('/api/v1/auth/me');
  },
};
