import { apiClient } from '../../api/client';
import type { ListResult } from '../../types/api';
import type {
  Role,
  Permission,
  RoleListQuery,
  PermissionListQuery,
  CreateRoleRequest,
  UpdateRoleRequest,
  AssignPermissionsRequest,
  AssignRolesToUserRequest,
  CreatePermissionRequest,
  UpdatePermissionRequest,
  PermissionGroup,
} from '../../types/rbac';

/**
 * 角色列表响应
 */
export type RoleListResponse = ListResult<Role>;

/**
 * 权限列表响应
 */
export type PermissionListResponse = ListResult<Permission>;

/**
 * 角色API服务
 */
export const roleApi = {
  /**
   * 获取角色列表
   */
  getList: (params: RoleListQuery): Promise<RoleListResponse> => {
    return apiClient.get('/api/v1/admin/roles', { params });
  },

  /**
   * 获取角色详情
   */
  getDetail: (id: number): Promise<Role> => {
    return apiClient.get(`/api/v1/admin/roles/${id}`);
  },

  /**
   * 创建角色
   */
  create: (data: CreateRoleRequest): Promise<Role> => {
    return apiClient.post('/api/v1/admin/roles', data);
  },

  /**
   * 更新角色信息
   */
  update: (id: number, data: UpdateRoleRequest): Promise<Role> => {
    return apiClient.put(`/api/v1/admin/roles/${id}`, data);
  },

  /**
   * 删除角色
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/roles/${id}`);
  },

  /**
   * 分配权限给角色
   */
  assignPermissions: (id: number, data: AssignPermissionsRequest): Promise<Role> => {
    return apiClient.put(`/api/v1/admin/roles/${id}/permissions`, data);
  },

  /**
   * 分配角色给用户
   */
  assignToUser: (data: AssignRolesToUserRequest): Promise<void> => {
    return apiClient.post('/api/v1/admin/roles/assign-user', data);
  },

  /**
   * 获取角色的权限列表
   */
  getPermissions: (id: number): Promise<Permission[]> => {
    return apiClient.get(`/api/v1/admin/roles/${id}/permissions`);
  },

  /**
   * 获取用户的角色列表
   */
  getUserRoles: (userId: number): Promise<Role[]> => {
    return apiClient.get(`/api/v1/admin/users/${userId}/roles`);
  },
};

/**
 * 权限API服务
 */
export const permissionApi = {
  /**
   * 获取权限列表
   */
  getList: (params: PermissionListQuery): Promise<PermissionListResponse> => {
    return apiClient.get('/api/v1/admin/permissions', { params });
  },

  /**
   * 获取权限分组列表
   */
  getGroups: (): Promise<PermissionGroup[]> => {
    return apiClient.get('/api/v1/admin/permissions/groups');
  },

  /**
   * 获取权限详情
   */
  getDetail: (id: number): Promise<Permission> => {
    return apiClient.get(`/api/v1/admin/permissions/${id}`);
  },

  /**
   * 创建权限
   */
  create: (data: CreatePermissionRequest): Promise<Permission> => {
    return apiClient.post('/api/v1/admin/permissions', data);
  },

  /**
   * 更新权限信息
   */
  update: (id: number, data: UpdatePermissionRequest): Promise<Permission> => {
    return apiClient.put(`/api/v1/admin/permissions/${id}`, data);
  },

  /**
   * 删除权限
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/permissions/${id}`);
  },

  /**
   * 获取用户的权限列表
   */
  getUserPermissions: (userId: number): Promise<Permission[]> => {
    return apiClient.get(`/api/v1/admin/users/${userId}/permissions`);
  },
};


