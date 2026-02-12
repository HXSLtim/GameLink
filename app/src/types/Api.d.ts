/**
 * API 相关类型定义
 */

// 通用分页参数
export interface PaginationParams {
  page?: number
  pageSize?: number
}

// 通用分页响应
export interface Pagination {
  page: number
  pageSize: number
  total: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}

// API 响应格式
export interface ApiResponse<T = any> {
  success: boolean
  code: number
  message: string
  data: T
  pagination?: Pagination
  meta?: any
  traceId?: string
}

// 列表响应
export interface ListResponse<T> {
  items: T[]
  pagination: Pagination
}
