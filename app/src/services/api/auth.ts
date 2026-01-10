/**
 * Auth API
 * 认证相关 API
 */
import apiClient from './request';
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  WeChatLoginRequest,
  SwitchRoleRequest,
  SwitchRoleResponse,
  UserRole,
} from '@/types/api';

export const authApi = {
  /**
   * 用户登录
   */
  login: (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    return apiClient.post<LoginResponse>('/auth/login', data);
  },

  /**
   * 用户注册
   */
  register: (data: LoginRequest & { name: string }): Promise<ApiResponse<LoginResponse>> => {
    return apiClient.post<LoginResponse>('/auth/register', data);
  },

  /**
   * 微信小程序登录
   */
  wechatLogin: (data: WeChatLoginRequest): Promise<ApiResponse<LoginResponse>> => {
    return apiClient.post<LoginResponse>('/public/auth/wechat-login', data);
  },

  /**
   * 刷新 Token
   */
  refreshToken: (refreshToken: string): Promise<ApiResponse<{ token: string; refreshToken: string }>> => {
    return apiClient.post('/auth/refresh', { refreshToken }, { skipAuth: true });
  },

  /**
   * 切换角色
   */
  switchRole: (role: UserRole): Promise<ApiResponse<SwitchRoleResponse>> => {
    const data: SwitchRoleRequest = { role };
    return apiClient.post<SwitchRoleResponse>('/user/role/switch', data);
  },

  /**
   * 登出
   */
  logout: (): Promise<ApiResponse<void>> => {
    return apiClient.post('/auth/logout');
  },
};
