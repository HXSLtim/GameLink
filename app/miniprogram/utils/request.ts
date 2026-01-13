/**
 * 网络请求封装
 */

import { config } from '../config/index'
import { getStorage, setStorage, removeStorage, StorageKeys } from './storage'

// API 基础地址
const BASE_URL = config.api.baseUrl

// 请求配置
interface RequestConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'OPTIONS' | 'HEAD'
  data?: Record<string, unknown>
  header?: Record<string, string>
  showLoading?: boolean
  showError?: boolean
}

// API 响应结构
interface ApiResponse<T = unknown> {
  success: boolean
  code: number
  message: string
  data: T
}

// 请求错误
class RequestError extends Error {
  code: number
  constructor(message: string, code: number) {
    super(message)
    this.code = code
    this.name = 'RequestError'
  }
}

/**
 * 获取请求头
 */
function getHeaders(): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  
  const token = getStorage<string>(StorageKeys.TOKEN)
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  return headers
}

/**
 * 刷新 Token
 */
async function refreshToken(): Promise<boolean> {
  const refreshToken = getStorage<string>(StorageKeys.REFRESH_TOKEN)
  if (!refreshToken) return false
  
  try {
    const res = await new Promise<WechatMiniprogram.RequestSuccessCallbackResult>((resolve, reject) => {
      wx.request({
        url: `${BASE_URL}/auth/refresh`,
        method: 'POST',
        data: { refreshToken },
        success: resolve,
        fail: reject,
      })
    })
    
    const data = res.data as ApiResponse<{ token: string; refreshToken: string }>
    if (data.success) {
      setStorage(StorageKeys.TOKEN, data.data.token)
      setStorage(StorageKeys.REFRESH_TOKEN, data.data.refreshToken)
      return true
    }
    return false
  } catch {
    return false
  }
}

/**
 * 处理登出
 */
function handleLogout() {
  removeStorage(StorageKeys.TOKEN)
  removeStorage(StorageKeys.REFRESH_TOKEN)
  removeStorage(StorageKeys.USER_INFO)
  
  wx.reLaunch({ url: '/pages/index/index' })
}

/**
 * 发起请求
 */
export async function request<T = unknown>(config: RequestConfig): Promise<T> {
  const {
    url,
    method = 'GET',
    data,
    header = {},
    showLoading = false,
    showError = true,
  } = config
  
  if (showLoading) {
    wx.showLoading({ title: '加载中...', mask: true })
  }
  
  try {
    const res = await new Promise<WechatMiniprogram.RequestSuccessCallbackResult>((resolve, reject) => {
      wx.request({
        url: url.startsWith('http') ? url : `${BASE_URL}${url}`,
        method,
        data,
        header: { ...getHeaders(), ...header },
        success: resolve,
        fail: reject,
      })
    })
    
    const response = res.data as ApiResponse<T>
    
    // 401 尝试刷新 token
    if (res.statusCode === 401) {
      const refreshed = await refreshToken()
      if (refreshed) {
        return request(config) // 重试请求
      }
      handleLogout()
      throw new RequestError('登录已过期', 401)
    }
    
    // 业务错误
    if (!response.success) {
      throw new RequestError(response.message || '请求失败', response.code)
    }
    
    return response.data
  } catch (err) {
    if (showError && err instanceof RequestError) {
      wx.showToast({ title: err.message, icon: 'none' })
    }
    throw err
  } finally {
    if (showLoading) {
      wx.hideLoading()
    }
  }
}

// 便捷方法
export const http = {
  get: <T>(url: string, config?: Omit<RequestConfig, 'url' | 'method'>) =>
    request<T>({ url, method: 'GET', ...config }),
    
  post: <T>(url: string, data?: Record<string, unknown>, config?: Omit<RequestConfig, 'url' | 'method' | 'data'>) =>
    request<T>({ url, method: 'POST', data, ...config }),
    
  put: <T>(url: string, data?: Record<string, unknown>, config?: Omit<RequestConfig, 'url' | 'method' | 'data'>) =>
    request<T>({ url, method: 'PUT', data, ...config }),
    
  delete: <T>(url: string, config?: Omit<RequestConfig, 'url' | 'method'>) =>
    request<T>({ url, method: 'DELETE', ...config }),
}
