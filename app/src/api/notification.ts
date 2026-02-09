/**
 * 通知相关 API
 */

import { get, post, type RequestConfig } from './request'

// 通知类型
export type NotificationType = 
  | 'order'           // 订单相关
  | 'chat'            // 聊天消息
  | 'system'          // 系统通知
  | 'promotion'       // 活动推广
  | 'payment'         // 支付相关
  | 'review'          // 评价相关
  | 'certification'   // 认证相关

// 通知项
export interface Notification {
  id: number
  type: NotificationType
  title: string
  content: string
  extra?: {
    orderId?: number
    userId?: number
    chatGroupId?: number
    url?: string
  }
  isRead: boolean
  createdAt: string
}

// 未读数量
export interface UnreadCount {
  total: number
  order: number
  chat: number
  system: number
  promotion: number
}

// 分页参数
export interface NotificationListParams {
  page?: number
  page_size?: number
  type?: NotificationType
  isRead?: boolean
}

/**
 * 获取通知列表
 */
export function getNotifications(params?: NotificationListParams, config?: Partial<RequestConfig>) {
  return get<Notification[]>('/users/notifications', params, config)
}

/**
 * 获取未读通知数量
 */
export function getUnreadCount(config?: Partial<RequestConfig>) {
  return get<UnreadCount>('/users/notifications/unread-count', undefined, config)
}

/**
 * 标记单条通知为已读
 */
export function markNotificationRead(notificationId: number, config?: Partial<RequestConfig>) {
  return post<void>(`/users/notifications/${notificationId}/read`, undefined, config)
}

/**
 * 标记所有通知为已读
 */
export function markAllNotificationsRead(config?: Partial<RequestConfig>) {
  return post<void>('/users/notifications/read-all', undefined, config)
}

export default {
  getNotifications,
  getUnreadCount,
  markNotificationRead,
  markAllNotificationsRead,
}
