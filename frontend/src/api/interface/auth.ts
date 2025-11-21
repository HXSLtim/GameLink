/**
 * 认证相关类型定义
 */

/**
 * 用户信息
 */
export interface User {
  id: number;
  username: string;
  email?: string;
  avatar?: string;
  role: 'admin' | 'player' | 'user';
  status: 'active' | 'inactive' | 'banned';
  createdAt: string;
  updatedAt: string;
}

/**
 * 登录请求
 */
export interface LoginRequest {
  username: string;
  password: string;
  remember?: boolean;
}

/**
 * 登录响应
 */
export interface LoginResponse {
  token: string;
  user: User;
  expiresIn: number;
}

/**
 * 注册请求
 */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

/**
 * 注册响应
 */
export interface RegisterResponse {
  userId: number;
  message: string;
}
