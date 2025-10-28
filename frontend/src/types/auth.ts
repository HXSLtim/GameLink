import type { UserRole, UserStatus } from './user';

/**
 * 登录请求
 */
export interface LoginRequest {
  username: string;
  password: string;
}

/**
 * 登录结果
 */
export interface LoginResult {
  token: string;
  expires_in?: number;
  user?: CurrentUser;
}

/**
 * 刷新令牌请求
 */
export interface RefreshTokenRequest {
  refresh_token: string;
}

/**
 * 当前用户信息
 */
export interface CurrentUser {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  avatar_url?: string;
  role: UserRole;
  status: UserStatus;
  last_login_at?: string;
  created_at?: string;
  updated_at?: string;
  // 兼容字段
  username?: string;
  avatar?: string;
}
