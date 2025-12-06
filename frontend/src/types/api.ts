/**
 * 共享的 API 类型定义
 * 所有 API 模块应该从这里导入通用类型，避免重复定义
 */

/**
 * 分页信息
 */
export interface Pagination {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

/**
 * 通用 API 响应结构
 */
export interface ApiResponse<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  pagination?: Pagination;
  traceId?: string;
}

/**
 * 分页列表响应
 */
export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: Pagination;
}

/**
 * 通用查询参数
 */
export interface BaseQueryParams {
  page?: number;
  page_size?: number;
  keyword?: string;
}

/**
 * 日期范围查询参数
 */
export interface DateRangeParams {
  date_from?: string;
  date_to?: string;
}

/**
 * 排序参数
 */
export interface SortParams {
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

/**
 * 通用 ID 类型
 */
export type ID = number | string;

/**
 * 通用状态类型
 */
export type Status = 'active' | 'inactive' | 'pending' | 'deleted';

/**
 * 时间戳字段
 */
export interface Timestamps {
  createdAt: string;
  updatedAt?: string;
  deletedAt?: string | null;
}
