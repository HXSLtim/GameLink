/**
 * 用户角色枚举
 */
export enum UserRole {
  USER = 'user', // 普通用户
  PLAYER = 'player', // 陪玩师
  ADMIN = 'admin', // 管理员
}

/**
 * 用户状态枚举
 */
export enum UserStatus {
  ACTIVE = 'active', // 正常
  SUSPENDED = 'suspended', // 暂停
  BANNED = 'banned', // 封禁
}

/**
 * 用户基础信息
 */
export interface User {
  id: number;
  phone?: string;
  email?: string;
  name: string;
  avatar_url?: string;
  role: UserRole;
  status: UserStatus;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

/**
 * 用户列表查询参数
 */
export interface UserListQuery {
  page?: number;
  page_size?: number;
  keyword?: string; // 搜索关键词（姓名/手机/邮箱）
  role?: UserRole;
  status?: UserStatus;
  created_start?: string;
  created_end?: string;
}

/**
 * 用户详情
 */
export interface UserDetail extends User {
  // 统计信息
  order_count?: number; // 订单数量
  total_spent?: number; // 总消费（分）
  review_count?: number; // 评价数量

  // 陪玩师信息（如果角色是 player）
  player?: PlayerInfo;
}

/**
 * 陪玩师认证状态
 */
export enum VerificationStatus {
  PENDING = 'pending', // 待认证
  VERIFIED = 'verified', // 已认证
  REJECTED = 'rejected', // 已拒绝
}

/**
 * 陪玩师信息
 */
export interface PlayerInfo {
  id: number;
  user_id: number;
  nickname?: string;
  bio?: string;
  rating_average: number;
  rating_count: number;
  hourly_rate_cents: number; // 时薪（分）
  main_game_id?: number;
  verification_status: VerificationStatus;
  created_at: string;
  updated_at: string;
}

/**
 * 用户状态更新请求
 */
export interface UpdateUserStatusRequest {
  status: UserStatus;
  reason?: string; // 封禁/暂停原因
}

/**
 * 用户角色更新请求
 */
export interface UpdateUserRoleRequest {
  role: UserRole;
}

/**
 * 用户列表响应
 */
export interface UserListResponse {
  list: User[];
  total: number;
  page: number;
  page_size: number;
}
