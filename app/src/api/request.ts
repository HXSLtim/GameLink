/**
 * API 请求封装
 * 支持 Token 自动携带、错误处理、响应解析
 */

import { useUserStore } from '@/store/user'

// API 基础路径
// 开发环境：自动使用当前页面的 hostname，解决局域网访问问题
// 生产环境：使用环境变量配置
function resolveBaseUrl(): string {
  const envUrl = import.meta.env.VITE_API_BASE_URL
  if (envUrl) return envUrl

  // 未配置时，自动匹配当前访问的主机名（支持 localhost 和 LAN IP）
  // #ifdef H5
  const host = window.location.hostname || 'localhost'
  return `http://${host}:8080/api/v1`
  // #endif

  // 非 H5 平台使用默认值
  return 'http://localhost:8080/api/v1'
}

const BASE_URL = resolveBaseUrl()

// 请求配置
export interface RequestConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: Record<string, any>
  params?: Record<string, any>
  header?: Record<string, string>
  showLoading?: boolean
  showError?: boolean
}

// API 响应格式
interface ApiResponse<T = any> {
  success: boolean
  code: number
  message: string
  data: T
  pagination?: {
    page: number
    page_size: number
    total: number
    totalPages: number
  }
}

// 错误类型
export class ApiError extends Error {
  code: number
  
  constructor(message: string, code: number) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}

/**
 * 构建完整 URL（带查询参数）
 */
function buildUrl(url: string, params?: Record<string, any>): string {
  if (!params) return url
  
  const queryString = Object.entries(params)
    .filter(([, value]) => value !== undefined && value !== null && value !== '')
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    .join('&')
  
  if (!queryString) return url
  return url.includes('?') ? `${url}&${queryString}` : `${url}?${queryString}`
}

/**
 * 通用请求函数
 */
export function request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
  const userStore = useUserStore()
  
  const {
    url,
    method = 'GET',
    data,
    params,
    header = {},
    showLoading = false,
    showError = true,
  } = config
  
  // 显示加载
  if (showLoading) {
    uni.showLoading({ title: '加载中...', mask: true })
  }
  
  // 构建请求头
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...header,
  }
  
  // 添加 Token
  if (userStore.token) {
    headers['Authorization'] = `Bearer ${userStore.token}`
  }
  
  // 构建完整 URL
  const fullUrl = BASE_URL + buildUrl(url, params)
  
  return new Promise((resolve, reject) => {
    uni.request({
      url: fullUrl,
      method,
      data,
      header: headers,
      timeout: 30000,
      success: (res) => {
        if (showLoading) {
          uni.hideLoading()
        }
        
        const statusCode = res.statusCode
        const responseData = res.data as ApiResponse<T>
        
        // HTTP 状态码检查
        if (statusCode === 401) {
          // 未授权，清除登录状态
          userStore.logout()
          reject(new ApiError('登录已过期，请重新登录', 401))
          return
        }
        
        if (statusCode === 403) {
          if (showError) {
            uni.showToast({ title: '无权限访问', icon: 'none' })
          }
          reject(new ApiError('无权限访问', 403))
          return
        }
        
        if (statusCode >= 500) {
          if (showError) {
            uni.showToast({ title: '服务器错误，请稍后重试', icon: 'none' })
          }
          reject(new ApiError('服务器错误', statusCode))
          return
        }
        
        // 业务响应检查
        if (!responseData.success) {
          if (showError) {
            uni.showToast({ title: responseData.message || '请求失败', icon: 'none' })
          }
          reject(new ApiError(responseData.message, responseData.code))
          return
        }
        
        resolve(responseData)
      },
      fail: (err) => {
        if (showLoading) {
          uni.hideLoading()
        }
        
        console.error('Request failed:', err)
        
        if (showError) {
          uni.showToast({ title: '网络请求失败', icon: 'none' })
        }
        
        reject(new ApiError('网络请求失败', -1))
      },
    })
  })
}

/**
 * GET 请求
 */
export function get<T = any>(url: string, params?: Record<string, any>, config?: Partial<RequestConfig>) {
  return request<T>({ url, method: 'GET', params, ...config })
}

/**
 * POST 请求
 */
export function post<T = any>(url: string, data?: Record<string, any>, config?: Partial<RequestConfig>) {
  return request<T>({ url, method: 'POST', data, ...config })
}

/**
 * PUT 请求
 */
export function put<T = any>(url: string, data?: Record<string, any>, config?: Partial<RequestConfig>) {
  return request<T>({ url, method: 'PUT', data, ...config })
}

/**
 * DELETE 请求
 */
export function del<T = any>(url: string, params?: Record<string, any>, config?: Partial<RequestConfig>) {
  return request<T>({ url, method: 'DELETE', params, ...config })
}

/**
 * 上传文件
 */
export function uploadFile(
  url: string, 
  filePath: string, 
  name: string = 'file',
  formData?: Record<string, any>
): Promise<ApiResponse> {
  const userStore = useUserStore()
  
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: BASE_URL + url,
      filePath,
      name,
      formData,
      header: {
        'Authorization': userStore.token ? `Bearer ${userStore.token}` : '',
      },
      success: (res) => {
        try {
          const data = JSON.parse(res.data) as ApiResponse
          if (data.success) {
            resolve(data)
          } else {
            uni.showToast({ title: data.message || '上传失败', icon: 'none' })
            reject(new ApiError(data.message, data.code))
          }
        } catch (e) {
          reject(new ApiError('解析响应失败', -1))
        }
      },
      fail: (err) => {
        console.error('Upload failed:', err)
        uni.showToast({ title: '上传失败', icon: 'none' })
        reject(new ApiError('上传失败', -1))
      },
    })
  })
}

export default {
  request,
  get,
  post,
  put,
  del,
  uploadFile,
}
