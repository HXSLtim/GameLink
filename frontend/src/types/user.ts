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
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
}

/**
 * 用户基础信息 - 与后端 model.User 保持一致
 */
export interface User extends BaseEntity {
  phone?: string;
  email?: string;
  name: string;
  avatarUrl?: string;
  role: UserRole;
  status: UserStatus;
  lastLoginAt?: string;
}

/**
 * 陪玩师信息 - 与后端 model.Player 保持一致
 */
export interface Player extends BaseEntity {
  userId: number;
  user?: User; // 关联用户信息
  nickname?: string;
  bio?: string;
  rank?: string;
  rating?: number;
  ratingAverage: number;
  ratingCount: number;
  hourlyRateCents: number;
  mainGameId?: number;
  mainGame?: { id: number; name: string }; // 关联游戏信息
  verificationStatus: VerificationStatus;
  isVerified?: boolean;
  isAvailable?: boolean;
}

/**
 * 用户详情（包含统计和陪玩师信息）
 */
export interface UserDetail extends User {
  // 统计信息
  orderCount?: number; // 订单数量
  totalSpent?: number; // 总消费（分）
  reviewCount?: number; // 评价数量

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
  pageSize?: number;
  keyword?: string; // 搜索关键词（姓名/手机/邮箱）
  role?: UserRole;
  status?: UserStatus;
  createdStart?: string;
  createdEnd?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'name' | 'lastLoginAt';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 陪玩师列表查询参数
 */
export interface PlayerListQuery {
  page?: number;
  pageSize?: number;
  userId?: number;
  mainGameId?: number;
  verificationStatus?: VerificationStatus;
  isVerified?: boolean;
  minRating?: number;
  maxRating?: number;
  minHourlyRate?: number;
  maxHourlyRate?: number;
  keyword?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'ratingAverage' | 'hourlyRateCents' | 'ratingCount';
  sortOrder?: 'asc' | 'desc';
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
  avatarUrl?: string;
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
  userId: number;
  nickname?: string;
  bio?: string;
  hourlyRateCents: number;
  mainGameId?: number;
}

/**
 * 更新陪玩师请求
 */
export interface UpdatePlayerRequest {
  nickname?: string;
  bio?: string;
  hourlyRateCents?: number;
  mainGameId?: number;
  verificationStatus?: VerificationStatus;
}

/**
 * 创建用户及陪玩师请求
 */
export interface CreateUserWithPlayerRequest {
  user: CreateUserRequest;
  player: Omit<CreatePlayerRequest, 'user_id'>;
}
