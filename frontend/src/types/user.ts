/**
 * 用户角色枚举
 */
export type Role = 'user' | 'player' | 'admin';

/**
 * 用户状态枚举
 */
export type UserStatus = 'active' | 'suspended' | 'banned';

/**
 * 实体 ID 类型
 * 注意：后端使用 uint64，但 JavaScript number 只能安全表示到 2^53-1
 * 对于超大 ID，后端应返回字符串格式
 */
export type EntityId = number | string;

/**
 * 基础实体接口
 */
export interface BaseEntity {
  /**
   * 实体 ID
   * 支持 number 和 string 类型以兼容大整数（uint64）
   */
  id: EntityId;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

/**
 * 用户实体 - 与后端model.User保持一致
 */
export interface User extends BaseEntity {
  phone?: string;
  email?: string;
  name: string;
  avatar_url?: string;
  role: Role;
  status: UserStatus;
  last_login_at?: string;
}

/**
 * 用户列表查询参数
 */
export interface UserListQuery {
  page?: number;
  page_size?: number;
  role?: Role;
  status?: UserStatus;
  keyword?: string;
  sort_by?: 'created_at' | 'updated_at' | 'name' | 'last_login_at';
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
  role: Role;
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
  role?: Role;
  status?: UserStatus;
}
