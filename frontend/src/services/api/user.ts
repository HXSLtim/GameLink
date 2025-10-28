import { apiClient } from '../../api/client';
import type {
  User,
  UserDetail,
  UserListQuery,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateUserStatusRequest,
  UpdateUserRoleRequest,
  Player,
  PlayerListQuery,
  CreatePlayerRequest,
  UpdatePlayerRequest,
  CreateUserWithPlayerRequest,
} from '../../types/user';
import type { ListResult } from '../../types/api';

/**
 * 用户列表响应
 */
export type UserListResponse = ListResult<User>;

/**
 * 陪玩师列表响应
 */
export type PlayerListResponse = ListResult<Player>;

/**
 * 用户API服务
 */
export const userApi = {
  /**
   * 获取用户列表
   */
  getList: (params: UserListQuery): Promise<UserListResponse> => {
    return apiClient.get('/api/v1/admin/users', { params });
  },

  /**
   * 获取用户详情
   */
  getDetail: (id: number): Promise<UserDetail> => {
    return apiClient.get(`/api/v1/admin/users/${id}`);
  },

  /**
   * 创建用户
   */
  create: (data: CreateUserRequest): Promise<User> => {
    return apiClient.post('/api/v1/admin/users', data);
  },

  /**
   * 更新用户信息
   */
  update: (id: number, data: UpdateUserRequest): Promise<User> => {
    return apiClient.put(`/api/v1/admin/users/${id}`, data);
  },

  /**
   * 更新用户状态
   */
  updateStatus: (id: number, data: UpdateUserStatusRequest): Promise<User> => {
    return apiClient.put(`/api/v1/admin/users/${id}/status`, data);
  },

  /**
   * 更新用户角色
   */
  updateRole: (id: number, data: UpdateUserRoleRequest): Promise<User> => {
    return apiClient.put(`/api/v1/admin/users/${id}/role`, data);
  },

  /**
   * 删除用户（软删除）
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/users/${id}`);
  },

  /**
   * 获取用户的订单列表
   */
  getUserOrders: (id: number, params: { page?: number; page_size?: number }): Promise<unknown> => {
    return apiClient.get(`/api/v1/admin/users/${id}/orders`, { params });
  },

  /**
   * 创建用户及陪玩师信息
   */
  createWithPlayer: (data: CreateUserWithPlayerRequest): Promise<UserDetail> => {
    return apiClient.post('/api/v1/admin/users/with-player', data);
  },

  /**
   * 获取用户操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/users/${id}/logs`);
  },
};

/**
 * 陪玩师API服务
 */
export const playerApi = {
  /**
   * 获取陪玩师列表
   */
  getList: (params: PlayerListQuery): Promise<PlayerListResponse> => {
    return apiClient.get('/api/v1/admin/players', { params });
  },

  /**
   * 获取陪玩师详情
   */
  getDetail: (id: number): Promise<Player> => {
    return apiClient.get(`/api/v1/admin/players/${id}`);
  },

  /**
   * 创建陪玩师
   */
  create: (data: CreatePlayerRequest): Promise<Player> => {
    return apiClient.post('/api/v1/admin/players', data);
  },

  /**
   * 更新陪玩师信息
   */
  update: (id: number, data: UpdatePlayerRequest): Promise<Player> => {
    return apiClient.put(`/api/v1/admin/players/${id}`, data);
  },

  /**
   * 更新验证状态
   */
  updateVerification: (
    id: number,
    data: { verification_status: 'pending' | 'verified' | 'rejected' },
  ): Promise<Player> => {
    return apiClient.put(`/api/v1/admin/players/${id}/verification`, data);
  },

  /**
   * 删除陪玩师
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/players/${id}`);
  },

  /**
   * 获取陪玩师操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/players/${id}/logs`);
  },
};

// 导出类型以便组件使用
export type { User as UserInfo, UserDetail };
