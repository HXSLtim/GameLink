/**
 * RBAC (Role-Based Access Control) 类型定义
 * 对应后端 model.RoleModel 和 model.Permission
 */

import type { BaseEntity } from './user';

/**
 * HTTP方法枚举
 */
export enum HTTPMethod {
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
  PATCH = 'PATCH',
  DELETE = 'DELETE',
}

/**
 * 角色标识枚举
 */
export enum RoleSlug {
  SUPER_ADMIN = 'super_admin',
  ADMIN = 'admin',
  PLAYER = 'player',
  USER = 'user',
}

/**
 * 权限信息
 */
export interface Permission extends BaseEntity {
  method: HTTPMethod;
  path: string;
  code: string;
  group: string;
  description: string;
}

/**
 * 角色信息
 */
export interface Role extends BaseEntity {
  slug: string;
  name: string;
  description: string;
  isSystem: boolean;
  permissions?: Permission[];
}

/**
 * 角色权限关联
 */
export interface RolePermission {
  roleId: number;
  permissionId: number;
  createdAt: string;
  role?: Role;
  permission?: Permission;
}

/**
 * 权限分组信息
 */
export interface PermissionGroup {
  group: string;
  description?: string;
  permissions: Permission[];
}

/**
 * 角色列表查询参数
 */
export interface RoleListQuery {
  page?: number;
  pageSize?: number;
  keyword?: string;
  isSystem?: boolean;
  sortBy?: 'createdAt' | 'updatedAt' | 'name';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 权限列表查询参数
 */
export interface PermissionListQuery {
  page?: number;
  pageSize?: number;
  keyword?: string;
  group?: string;
  method?: HTTPMethod;
  sortBy?: 'createdAt' | 'group' | 'path';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 创建角色请求
 */
export interface CreateRoleRequest {
  slug: string;
  name: string;
  description?: string;
  permissionIds?: number[];
}

/**
 * 更新角色请求
 */
export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

/**
 * 分配权限给角色请求
 */
export interface AssignPermissionsRequest {
  permissionIds: number[];
}

/**
 * 分配角色给用户请求
 */
export interface AssignRolesToUserRequest {
  userId: number;
  roleIds: number[];
}

/**
 * 创建权限请求
 */
export interface CreatePermissionRequest {
  method: HTTPMethod;
  path: string;
  code: string;
  group: string;
  description?: string;
}

/**
 * 更新权限请求
 */
export interface UpdatePermissionRequest {
  code?: string;
  group?: string;
  description?: string;
}


