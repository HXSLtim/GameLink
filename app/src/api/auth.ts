/**
 * 认证相关 API
 */

import { post, get, put, type RequestConfig } from './request'
import type { UserInfo } from '@/store/user'
import type { AppUserRole } from '@/types/user'

// ============================================
// 请求类型
// ============================================

export interface LoginRequest {
  phone?: string   // 手机号
  email?: string   // 邮箱
  password: string
}

export interface WeChatLoginRequest {
  code: string
  encryptedData?: string
  iv?: string
  referralCode?: string
}

export interface RegisterRequest {
  phone: string
  password: string
  nickname: string
  role: AppUserRole
  verifyCode?: string
}

// ============================================
// 响应类型
// ============================================

export interface LoginResponse {
  token: string
  expires_at: string
  user: UserInfo
  // 兼容字段（部分接口可能使用）
  accessToken?: string
  refreshToken?: string
}

// ============================================
// API 函数
// ============================================

/**
 * 账号密码登录
 */
export function login(data: LoginRequest, config?: Partial<RequestConfig>) {
  return post<LoginResponse>('/auth/login', data, config)
}

/**
 * 微信小程序登录
 */
export function wechatLogin(data: WeChatLoginRequest) {
  return post<LoginResponse>('/public/auth/wechat/login', data)
}

/**
 * 刷新 Token
 */
export function refreshToken(token: string) {
  return post<{ accessToken: string; expiresIn: number }>('/auth/refresh', { refreshToken: token })
}

/**
 * 获取用户信息
 */
export function getProfile() {
  return get<UserInfo>('/users/me')
}

/**
 * 注册
 */
export function register(data: RegisterRequest) {
  return post<LoginResponse>('/auth/register', data)
}

/**
 * 发送验证码（注册/登录等场景）
 */
export function sendVerifyCode(phone: string) {
  return post('/public/verification/send', { phone })
}

/**
 * 退出登录
 */
export function logout() {
  return post('/auth/logout')
}

/**
 * 修改密码
 */
export function changePassword(data: { oldPassword: string; newPassword: string }) {
  return post('/auth/change-password', data)
}

// ============================================
// 微信登录辅助函数
// ============================================

/**
 * 微信登录流程
 * 1. 调用 wx.login 获取 code
 * 2. 调用后端接口换取 token
 */
export async function doWeChatLogin(): Promise<LoginResponse> {
  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    uni.login({
      provider: 'weixin',
      success: async (loginRes) => {
        try {
          const response = await wechatLogin({ code: loginRes.code })
          resolve(response.data)
        } catch (error) {
          reject(error)
        }
      },
      fail: (err) => {
        console.error('wx.login failed:', err)
        reject(new Error('微信登录失败'))
      },
    })
    // #endif
    
    // #ifndef MP-WEIXIN
    reject(new Error('当前环境不支持微信登录'))
    // #endif
  })
}

export default {
  login,
  wechatLogin,
  refreshToken,
  getProfile,
  register,
  sendVerifyCode,
  logout,
  changePassword,
  doWeChatLogin,
}
