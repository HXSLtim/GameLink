/**
 * 通知相关接口定义
 */

import type { BaseEntity } from '@/shared/types/api';

/**
 * 通知类型
 */
export enum NotificationType {
  ORDER = 'order',
  PAYMENT = 'payment',
  SYSTEM = 'system',
  MESSAGE = 'message',
  REVIEW = 'review',
}

/**
 * 通知状态
 */
export enum NotificationStatus {
  UNREAD = 'unread',
  READ = 'read',
}

/**
 * 通知实体
 */
export interface Notification extends BaseEntity {
  id: number;
  userId: number;
  type: NotificationType;
  title: string;
  content: string;
  data?: Record<string, any>;
  status: NotificationStatus;
  readAt?: string;
}

/**
 * 通知列表查询参数
 */
export interface GetNotificationsParams {
  type?: NotificationType;
  status?: NotificationStatus;
  page?: number;
  pageSize?: number;
}

/**
 * 通知列表响应
 */
export interface GetNotificationsResponse {
  list: Notification[];
  total: number;
  page: number;
  pageSize: number;
  unreadCount: number;
}

/**
 * 未读通知数量
 */
export interface UnreadNotificationCount {
  total: number;
  byType: {
    [key in NotificationType]?: number;
  };
}

/**
 * 创建通知请求（系统通知）
 */
export interface CreateNotificationRequest {
  userId: number;
  type: NotificationType;
  title: string;
  content: string;
  data?: Record<string, any>;
}