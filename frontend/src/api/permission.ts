/**
 * RBAC 权限管理 API
 * Requirements: 1.1, 2.1, 2.2, 2.3, 6.3, 6.5, 9.1, 9.2
 */

import apiClient from './client';
import type { ApiResponse } from '../types/api';
import type {
  Permission,
  PermissionTreeNode,
  CreatePermissionDto,
  UpdatePermissionDto,
  PermissionQueryParams,
  Role,
  CreateRoleDto,
  UpdateRoleDto,
  RoleQueryParams,
  BatchAssignPermissionsDto,
  AssignUserRolesDto,
  BatchAssignUserRolesDto,
  PermissionAuditLog,
  AuditLogQueryParams,
  AuditLogExportParams,
  UserEffectivePermissions,
  PaginatedList,
} from '../types/permission';

/**
 * 权限管理 API
 * Requirements: 1.1, 2.1
 */
export const permissionApi = {
  /**
   * 获取权限列表（分页）
   * GET /api/admin/permissions
   */
  list: (params?: PermissionQueryParams) =>
    apiClient.get<ApiResponse<PaginatedList<Permission>>>('/admin/permissions', { params }),

  /**
   * 获取权限详情
   * GET /api/admin/permissions/:id
   */
  get: (id: number) =>
    apiClient.get<ApiResponse<Permission>>(`/admin/permissions/${id}`),

  /**
   * 获取权限树
   * GET /api/admin/permissions/tree
   */
  getTree: () =>
    apiClient.get<ApiResponse<PermissionTreeNode[]>>('/admin/permissions/tree'),

  /**
   * 获取权限分组列表
   * GET /api/admin/permissions/groups
   */
  getGroups: () =>
    apiClient.get<ApiResponse<string[]>>('/admin/permissions/groups'),

  /**
   * 创建权限
   * POST /api/admin/permissions
   */
  create: (data: CreatePermissionDto) =>
    apiClient.post<ApiResponse<Permission>>('/admin/permissions', data),

  /**
   * 全量更新权限
   * PUT /api/admin/permissions/:id
   */
  update: (id: number, data: UpdatePermissionDto) =>
    apiClient.put<ApiResponse<Permission>>(`/admin/permissions/${id}`, data),

  /**
   * 部分更新权限
   * PATCH /api/admin/permissions/:id
   */
  patch: (id: number, data: Partial<UpdatePermissionDto>) =>
    apiClient.patch<ApiResponse<Permission>>(`/admin/permissions/${id}`, data),

  /**
   * 删除权限（软删除）
   * DELETE /api/admin/permissions/:id
   */
  delete: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/permissions/${id}`),

  /**
   * 获取当前用户权限码列表
   * GET /api/admin/me/permissions
   */
  getMyPermissions: () =>
    apiClient.get<ApiResponse<string[]>>('/admin/me/permissions'),

  /**
   * 获取当前用户菜单
   * GET /api/admin/me/menus
   */
  getMyMenus: () =>
    apiClient.get<ApiResponse<unknown[]>>('/admin/me/menus'),
};


/**
 * 角色管理 API
 * Requirements: 2.2, 2.3
 */
export const roleApi = {
  /**
   * 获取角色列表（分页）
   * GET /api/admin/roles
   */
  list: (params?: RoleQueryParams) =>
    apiClient.get<ApiResponse<PaginatedList<Role>>>('/admin/roles', { params }),

  /**
   * 获取角色详情
   * GET /api/admin/roles/:id
   */
  get: (id: number) =>
    apiClient.get<ApiResponse<Role>>(`/admin/roles/${id}`),

  /**
   * 获取角色详情（包含权限）
   * GET /api/admin/roles/:id?with_permissions=true
   */
  getWithPermissions: (id: number) =>
    apiClient.get<ApiResponse<Role>>(`/admin/roles/${id}`, {
      params: { with_permissions: true },
    }),

  /**
   * 创建角色
   * POST /api/admin/roles
   */
  create: (data: CreateRoleDto) =>
    apiClient.post<ApiResponse<Role>>('/admin/roles', data),

  /**
   * 全量更新角色
   * PUT /api/admin/roles/:id
   */
  update: (id: number, data: UpdateRoleDto) =>
    apiClient.put<ApiResponse<Role>>(`/admin/roles/${id}`, data),

  /**
   * 部分更新角色
   * PATCH /api/admin/roles/:id
   */
  patch: (id: number, data: Partial<UpdateRoleDto>) =>
    apiClient.patch<ApiResponse<Role>>(`/admin/roles/${id}`, data),

  /**
   * 删除角色（软删除）
   * DELETE /api/admin/roles/:id
   */
  delete: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/roles/${id}`),

  /**
   * 获取角色权限ID列表
   * GET /api/admin/roles/:id/permissions
   */
  getPermissions: (id: number) =>
    apiClient.get<ApiResponse<number[]>>(`/admin/roles/${id}/permissions`),

  /**
   * 批量分配角色权限（事务保证原子性）
   * PUT /api/admin/roles/:id/permissions/batch
   */
  batchAssignPermissions: (roleId: number, data: BatchAssignPermissionsDto) =>
    apiClient.put<ApiResponse<void>>(
      `/admin/roles/${roleId}/permissions/batch`,
      data
    ),

  /**
   * 单个添加权限
   * POST /api/admin/roles/:id/permissions/:pid
   */
  addPermission: (roleId: number, permissionId: number) =>
    apiClient.post<ApiResponse<void>>(
      `/admin/roles/${roleId}/permissions/${permissionId}`
    ),

  /**
   * 单个移除权限
   * DELETE /api/admin/roles/:id/permissions/:pid
   */
  removePermission: (roleId: number, permissionId: number) =>
    apiClient.delete<ApiResponse<void>>(
      `/admin/roles/${roleId}/permissions/${permissionId}`
    ),

  /**
   * 设置角色继承关系
   * PUT /api/admin/roles/:id/parent
   */
  setParent: (roleId: number, parentId: number | null) =>
    apiClient.put<ApiResponse<Role>>(`/admin/roles/${roleId}/parent`, {
      parentId,
    }),

  /**
   * 获取角色继承链
   * GET /api/admin/roles/:id/inheritance-chain
   */
  getInheritanceChain: (roleId: number) =>
    apiClient.get<ApiResponse<Role[]>>(`/admin/roles/${roleId}/inheritance-chain`),
};


/**
 * 用户角色分配 API
 * Requirements: 9.1, 9.2, 10.3
 */
export const userRoleApi = {
  /**
   * 获取用户角色列表
   * GET /api/admin/users/:id/roles
   */
  getUserRoles: (userId: number) =>
    apiClient.get<ApiResponse<Role[]>>(`/admin/users/${userId}/roles`),

  /**
   * 分配用户角色
   * PUT /api/admin/users/:id/roles
   */
  assignRoles: (userId: number, data: AssignUserRolesDto) =>
    apiClient.put<ApiResponse<void>>(`/admin/users/${userId}/roles`, data),

  /**
   * 获取用户有效权限（合并后的完整权限列表）
   * GET /api/admin/users/:id/permissions
   */
  getUserPermissions: (userId: number) =>
    apiClient.get<ApiResponse<UserEffectivePermissions>>(
      `/admin/users/${userId}/permissions`
    ),

  /**
   * 批量分配用户角色
   * POST /api/admin/users/batch/roles
   */
  batchAssignRoles: (data: BatchAssignUserRolesDto) =>
    apiClient.post<ApiResponse<{
      success: number;
      failed: number;
      errors?: Array<{ userId: number; error: string }>;
    }>>('/admin/users/batch/roles', data),

  /**
   * 检查用户是否为超级管理员
   * GET /api/admin/users/:id/is-super-admin
   */
  checkSuperAdmin: (userId: number) =>
    apiClient.get<ApiResponse<boolean>>(`/admin/users/${userId}/is-super-admin`),
};


/**
 * 权限审计日志 API
 * Requirements: 6.3, 6.5
 */
export const auditLogApi = {
  /**
   * 获取权限审计日志列表（分页）
   * GET /api/admin/audit/permissions
   */
  list: (params?: AuditLogQueryParams) =>
    apiClient.get<ApiResponse<PaginatedList<PermissionAuditLog>>>(
      '/admin/audit/permissions',
      { params }
    ),

  /**
   * 获取审计日志详情
   * GET /api/admin/audit/permissions/:id
   */
  get: (id: number) =>
    apiClient.get<ApiResponse<PermissionAuditLog>>(
      `/admin/audit/permissions/${id}`
    ),

  /**
   * 导出审计日志（CSV格式）
   * GET /api/admin/audit/permissions/export
   */
  export: (params?: AuditLogExportParams) =>
    apiClient.get<Blob>('/admin/audit/permissions/export', {
      params,
      responseType: 'blob',
    }),

  /**
   * 获取审计日志统计
   * GET /api/admin/audit/permissions/stats
   */
  getStats: (params?: { date_from?: string; date_to?: string }) =>
    apiClient.get<ApiResponse<{
      total: number;
      byAction: Record<string, number>;
      byTargetType: Record<string, number>;
    }>>('/admin/audit/permissions/stats', { params }),
};

/**
 * 统一导出所有权限相关 API
 */
export const rbacApi = {
  permission: permissionApi,
  role: roleApi,
  userRole: userRoleApi,
  auditLog: auditLogApi,
};

export default rbacApi;
