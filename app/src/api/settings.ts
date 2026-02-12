/**
 * 用户设置相关 API
 */

import { get, put } from './request'

// 用户通用设置
export interface UserSettings {
  language: string
  theme: 'light' | 'dark' | 'system'
  fontSize: 'small' | 'medium' | 'large'
  hideOnlineStatus: boolean
  hideLocation: boolean
}

// 通知设置
export interface NotificationSettings {
  pushEnabled: boolean
  messageEnabled: boolean
  orderEnabled: boolean
  promotionEnabled: boolean
  soundEnabled: boolean
  vibrationEnabled: boolean
}

/**
 * 获取用户通用设置
 */
export function getUserSettings() {
  return get<UserSettings>('/user/settings')
}

/**
 * 更新用户通用设置
 */
export function updateUserSettings(data: Partial<UserSettings>) {
  return put<UserSettings>('/user/settings', data)
}

/**
 * 获取通知设置
 */
export function getNotificationSettings() {
  return get<NotificationSettings>('/user/notification-settings')
}

/**
 * 更新通知设置
 */
export function updateNotificationSettings(data: Partial<NotificationSettings>) {
  return put<NotificationSettings>('/user/notification-settings', data)
}

export default {
  getUserSettings,
  updateUserSettings,
  getNotificationSettings,
  updateNotificationSettings,
}
