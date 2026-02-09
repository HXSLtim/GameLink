/**
 * 用户信息相关 API
 */

import { get, put } from './request'
import { uploadFile } from './request'
import type { ProfileGender } from '@/types/common'
import type { AppUserRole } from '@/types/user'

// 用户资料
export interface UserProfile {
  id: number
  phone: string
  nickname: string
  avatar?: string
  gender: ProfileGender
  birthday?: string
  region?: string
  bio?: string
  role: AppUserRole
  playerId?: number
  vipLevel: number
  vipExpireAt?: string
  createdAt: string
}

// 更新资料参数
export interface UpdateProfileParams {
  nickname?: string
  gender?: ProfileGender
  birthday?: string
  region?: string
  bio?: string
}

/**
 * 获取当前用户资料
 */
export function getUserProfile() {
  return get<UserProfile>('/users/me')
}

/**
 * 更新用户资料
 */
export function updateUserProfile(data: UpdateProfileParams) {
  return put<UserProfile>('/users/me', data)
}

/**
 * 上传头像
 */
export function uploadAvatar(filePath: string) {
  return uploadFile('/users/avatar', filePath, 'avatar')
}

/**
 * 修改密码
 */
export function changePassword(data: {
  oldPassword: string
  newPassword: string
}) {
  return put<void>('/users/password', data)
}

/**
 * 绑定手机号
 */
export function bindPhone(data: {
  phone: string
  code: string
}) {
  return put<void>('/users/phone', data)
}

/**
 * 发送验证码
 */
export function sendVerifyCode(phone: string, type: 'bind' | 'reset' = 'bind') {
  return get<void>('/users/verify-code', { phone, type })
}

export default {
  getUserProfile,
  updateUserProfile,
  uploadAvatar,
  changePassword,
  bindPhone,
  sendVerifyCode,
}
