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
 * 认证状态枚举
 */
export enum VerificationStatus {
  PENDING = 'pending', // 待认证
  VERIFIED = 'verified', // 已认证
  REJECTED = 'rejected', // 已拒绝
}

/**
 * 实体 ID 类型
 * 注意：后端使用 uint64，但 JavaScript number 只能安全表示到 2^53-1
 * 对于超大 ID，后端应返回字符串格式
 */
export type EntityId = number;

/**
 * 基础实体接口
 */
export interface BaseEntity {
  /**
   * 实体 ID
   */
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

/**
 * 用户基础信息 - 与后端 model.User 保持一致
 */
export interface User extends BaseEntity {
  phone?: string;
  email?: string;
  name: string;
  avatar_url?: string;
  role: UserRole;
  status: UserStatus;
  last_login_at?: string;
}

/**
 * 陪玩师信息 - 与后端 model.Player 保持一致
 */
export interface Player extends BaseEntity {
  user_id: number;
  nickname?: string;
  bio?: string;
  rating_average: number;
  rating_count: number;
  hourly_rate_cents: number;
  main_game_id?: number;
  verification_status: VerificationStatus;
}

/**
 * 用户详情（包含统计和陪玩师信息）
 */
export interface UserDetail extends User {
  // 统计信息
  order_count?: number; // 订单数量
  total_spent?: number; // 总消费（分）
  review_count?: number; // 评价数量

  // 陪玩师信息（如果角色是 player）
  player?: Player;
}

/**
 * 陪玩师详情（包含用户信息）
 */
export interface PlayerDetail extends Player {
  user?: User;
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
  sort_by?: 'created_at' | 'updated_at' | 'name' | 'last_login_at';
  sort_order?: 'asc' | 'desc';
}

/**
 * 陪玩师列表查询参数
 */
export interface PlayerListQuery {
  page?: number;
  page_size?: number;
  user_id?: number;
  main_game_id?: number;
  verification_status?: VerificationStatus;
  min_rating?: number;
  max_rating?: number;
  min_hourly_rate?: number;
  max_hourly_rate?: number;
  keyword?: string;
  sort_by?: 'created_at' | 'updated_at' | 'rating_average' | 'hourly_rate_cents' | 'rating_count';
  sort_order?: 'asc' | 'desc';
}

/**
 * 创建用户请求
 */
export interface CreateUserRequest {
  phone?: string;
  email?: string;
  name: string;
  password: string;
  role?: UserRole;
  status?: UserStatus;
}

/**
 * 更新用户请求
 */
export interface UpdateUserRequest {
  phone?: string;
  email?: string;
  name?: string;
  avatar_url?: string;
  role?: UserRole;
  status?: UserStatus;
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
 * 创建陪玩师请求
 */
export interface CreatePlayerRequest {
  user_id: number;
  nickname?: string;
  bio?: string;
  hourly_rate_cents: number;
  main_game_id?: number;
}

/**
 * 更新陪玩师请求
 */
export interface UpdatePlayerRequest {
  nickname?: string;
  bio?: string;
  hourly_rate_cents?: number;
  main_game_id?: number;
  verification_status?: VerificationStatus;
}

/**
 * 创建用户及陪玩师请求
 */
export interface CreateUserWithPlayerRequest {
  user: CreateUserRequest;
  player: Omit<CreatePlayerRequest, 'user_id'>;
}
