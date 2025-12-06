/**
 * RBAC 权限管理类型定义
 * Requirements: 1.1, 2.1, 6.1
 */

import type { Timestamps } from './api';

/**
 * HTTP 方法类型
 */
export type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | '*';

/**
 * 权限模型
 * 对应后端 Permission 模型
 */
export interface Permission extends Timestamps {
  id: number;
  method: HTTPMethod;
  path: string;
  code: string;
  group: string;
  description?: string;
  parentId?: number | null;
  sortOrder: number;
  isSystem: boolean;
}

/**
 * 权限树节点
 * 用于树形展示权限
 */
export interface PermissionTreeNode extends Permission {
  children?: PermissionTreeNode[];
}

/**
 * 创建权限请求
 */
export interface CreatePermissionDto {
  method: HTTPMethod;
  path: string;
  code: string;
  group: string;
  description?: string;
  parentId?: number | null;
  sortOrder?: number;
}

/**
 * 更新权限请求（权限码不可修改）
 */
export interface UpdatePermissionDto {
  method?: HTTPMethod;
  path?: string;
  group?: string;
  description?: string;
  parentId?: number | null;
  sortOrder?: number;
}

/**
 * 权限查询参数
 */
export interface PermissionQueryParams {
  page?: number;
  page_size?: number;
  group?: string;
  keyword?: string;
  is_system?: boolean;
}

/**
 * 角色模型
 * 对应后端 RoleModel
 */
export interface Role extends Timestamps {
  id: number;
  slug: string;
  name: string;
  description?: string;
  isSystem: boolean;
  parentId?: number | null;
  priority: number;
  level: number;
  permissions?: Permission[];
}

/**
 * 角色树节点（用于继承关系展示）
 */
export interface RoleTreeNode extends Role {
  children?: RoleTreeNode[];
}

/**
 * 创建角色请求
 */
export interface CreateRoleDto {
  slug: string;
  name: string;
  description?: string;
  parentId?: number | null;
  priority?: number;
}

/**
 * 更新角色请求
 */
export interface UpdateRoleDto {
  name?: string;
  description?: string;
  parentId?: number | null;
  priority?: number;
}

/**
 * 角色查询参数
 */
export interface RoleQueryParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  is_system?: boolean;
}

/**
 * 批量权限分配请求
 */
export interface BatchAssignPermissionsDto {
  permissionIds: number[];
}

/**
 * 用户角色分配请求
 */
export interface AssignUserRolesDto {
  roleIds: number[];
}

/**
 * 批量用户角色分配请求
 */
export interface BatchAssignUserRolesDto {
  userIds: number[];
  roleIds: number[];
}

/**
 * 审计日志操作类型
 */
export type AuditAction =
  | 'permission_create'
  | 'permission_update'
  | 'permission_delete'
  | 'role_create'
  | 'role_update'
  | 'role_delete'
  | 'role_permission_assign'
  | 'user_role_assign';

/**
 * 审计日志目标类型
 */
export type AuditTargetType = 'permission' | 'role' | 'user';

/**
 * 权限审计日志模型
 * 对应后端 PermissionAuditLog
 */
export interface PermissionAuditLog {
  id: number;
  operatorId: number;
  operatorName: string;
  targetType: AuditTargetType;
  targetId: number;
  targetName: string;
  action: AuditAction;
  beforeData?: string;
  afterData?: string;
  ipAddress?: string;
  userAgent?: string;
  requestId?: string;
  createdAt: string;
}

/**
 * 审计日志查询参数
 */
export interface AuditLogQueryParams {
  page?: number;
  page_size?: number;
  target_type?: AuditTargetType;
  action?: AuditAction;
  operator_id?: number;
  target_id?: number;
  date_from?: string;
  date_to?: string;
}

/**
 * 审计日志导出参数
 */
export interface AuditLogExportParams {
  target_type?: AuditTargetType;
  action?: AuditAction;
  operator_id?: number;
  date_from?: string;
  date_to?: string;
  format?: 'csv' | 'json';
}

/**
 * 用户有效权限响应
 */
export interface UserEffectivePermissions {
  userId: number;
  permissions: string[];
  roles: Role[];
  isSuperAdmin: boolean;
}

/**
 * 权限检查模式
 */
export type PermissionCheckMode = 'any' | 'all';

/**
 * 分页列表响应
 */
export interface PaginatedList<T> {
  items: T[];
  totalCount: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

/**
 * 权限分组信息
 */
export interface PermissionGroup {
  name: string;
  count: number;
  permissions: Permission[];
}

