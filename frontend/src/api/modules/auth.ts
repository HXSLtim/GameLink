/**
 * 认证API服务
 */
import { request } from '../client';
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from '@/shared/types/auth';

/**
 * 认证服务
 */
export const authService = {
  /**
   * 用户登录
   * @param data 登录信息
   * @returns 登录结果
   */
  login: (data: LoginRequest): Promise<LoginResponse> => {
    return request.post<LoginResponse>('/auth/login', data);
  },

  /**
   * 用户注册
   * @param data 注册信息
   * @returns 注册结果
   */
  register: (data: RegisterRequest): Promise<RegisterResponse> => {
    return request.post<RegisterResponse>('/auth/register', data);
  },

  /**
   * 用户登出
   */
  logout: (): Promise<void> => {
    return request.post<void>('/auth/logout');
  },

  /**
   * 获取当前用户信息
   * @returns 用户信息
   */
  getCurrentUser: (): Promise<User> => {
    return request.get<User>('/auth/me');
  },

  /**
   * 刷新Token
   * @returns 新Token
   */
  refreshToken: (): Promise<{ token: string }> => {
    return request.post<{ token: string }>('/auth/refresh');
  },
};
