/**
 * 认证工具
 */

import { getStorage, setStorage, removeStorage, StorageKeys } from './storage'
import { http } from './request'

// 用户信息类型
export interface UserInfo {
  id: number
  nickname: string
  avatar: string
  phone?: string
  isPlayer: boolean // 是否为陪玩师
  playerStatus?: 'pending' | 'approved' | 'rejected' // 陪玩师认证状态
}

// 登录响应
interface LoginResponse {
  token: string
  refreshToken: string
  user: UserInfo
}

/**
 * 微信登录
 */
export async function wxLogin(): Promise<LoginResponse> {
  // 获取微信 code
  const { code } = await new Promise<WechatMiniprogram.LoginSuccessCallbackResult>((resolve, reject) => {
    wx.login({ success: resolve, fail: reject })
  })
  
  // 调用后端登录接口
  const data = await http.post<LoginResponse>('/auth/wx-login', { code }, { showLoading: true })
  
  // 保存登录信息
  setStorage(StorageKeys.TOKEN, data.token)
  setStorage(StorageKeys.REFRESH_TOKEN, data.refreshToken)
  setStorage(StorageKeys.USER_INFO, data.user)
  
  return data
}

/**
 * 获取用户信息
 */
export function getUserInfo(): UserInfo | null {
  return getStorage<UserInfo>(StorageKeys.USER_INFO)
}

/**
 * 更新用户信息
 */
export function setUserInfo(info: UserInfo): boolean {
  return setStorage(StorageKeys.USER_INFO, info)
}

/**
 * 是否已登录
 */
export function isLoggedIn(): boolean {
  return !!getStorage<string>(StorageKeys.TOKEN)
}

/**
 * 是否为陪玩师
 */
export function isPlayer(): boolean {
  const user = getUserInfo()
  return user !== null && user.isPlayer === true && user.playerStatus === 'approved'
}

/**
 * 登出
 */
export function logout(): void {
  removeStorage(StorageKeys.TOKEN)
  removeStorage(StorageKeys.REFRESH_TOKEN)
  removeStorage(StorageKeys.USER_INFO)
  removeStorage(StorageKeys.IDENTITY)
  
  wx.reLaunch({ url: '/pages/index/index' })
}

/**
 * 检查登录状态，未登录则跳转
 */
export function checkLogin(redirectUrl?: string): boolean {
  if (isLoggedIn()) return true
  
  const url = redirectUrl ? encodeURIComponent(redirectUrl) : ''
  wx.navigateTo({ url: `/pages/login/index?redirect=${url}` })
  return false
}

/**
 * 刷新用户信息
 */
export async function refreshUserInfo(): Promise<UserInfo> {
  const user = await http.get<UserInfo>('/user/profile')
  setUserInfo(user)
  return user
}
