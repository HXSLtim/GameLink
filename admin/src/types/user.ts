/**
 * User Module Types
 * 用户模块类型定义（含批量操作）
 */

export type UserRole = 'user' | 'player' | 'admin';

export type UserStatus = 'active' | 'banned' | 'suspended';

export type BatchTarget = 'users' | 'role' | 'all';

export type BatchNotificationType = 'system' | 'marketing' | 'personal' | 'activity';

export type BatchPointsType = 'admin' | 'activity' | 'compensation';

export interface BatchRoleRequest {
  userIds: number[];
  role: UserRole;
}

export interface BatchStatusRequest {
  userIds: number[];
  status: UserStatus;
}

export interface BatchPointsRequest {
  target: BatchTarget;
  userIds?: number[];
  roles?: UserRole[];
  cents: number;
  reason: string;
  type: BatchPointsType;
}

export interface BatchNotificationRequest {
  target: BatchTarget;
  userIds?: number[];
  roles?: UserRole[];
  title: string;
  content: string;
  type: BatchNotificationType;
}
