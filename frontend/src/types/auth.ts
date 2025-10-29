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
  expiresIn?: number;
  user?: CurrentUser;
}

/**
 * 刷新令牌请求
 */
export interface RefreshTokenRequest {
  refreshToken: string;
}

/**
 * 当前用户信息
 */
export interface CurrentUser {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  avatarUrl?: string;
  role: UserRole;
  status: UserStatus;
  lastLoginAt?: string;
  createdAt?: string;
  updatedAt?: string;
  // 兼容字段
  username?: string;
  avatar?: string;
}
