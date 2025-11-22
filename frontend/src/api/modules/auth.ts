/**
 * 认证API模块
 */

import { apiClient } from '@/api/client';
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  User,
} from '@/shared/types/auth';

/**
 * 用户登录
 */
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await apiClient.post<{ data: LoginResponse }>('/auth/login', data);
  return response.data;
};

/**
 * 用户注册
 */
export const register = async (data: RegisterRequest): Promise<RegisterResponse> => {
  const response = await apiClient.post<{ data: RegisterResponse }>('/auth/register', data);
  return response.data;
};

/**
 * 用户登出
 */
export const logout = async (): Promise<void> => {
  await apiClient.post('/auth/logout');
};

/**
 * 获取当前用户信息
 */
export const getCurrentUser = async (): Promise<User> => {
  const response = await apiClient.get<{ data: User }>('/auth/me');
  return response.data;
};

/**
 * 刷新Token
 */
export const refreshToken = async (): Promise<{ token: string }> => {
  const response = await apiClient.post<{ data: { token: string } }>('/auth/refresh');
  return response.data;
};

/**
 * 忘记密码
 */
export const forgotPassword = async (email: string): Promise<void> => {
  await apiClient.post('/auth/forgot-password', { email });
};

/**
 * 重置密码
 */
export const resetPassword = async (token: string, newPassword: string): Promise<void> => {
  await apiClient.post('/auth/reset-password', { token, newPassword });
};
