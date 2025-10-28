/**
 * API 成功响应格式 - 与后端统一
 */
export interface ApiSuccessResponse<T> {
  success: true;
  data: T;
  message?: string;
}

/**
 * API 列表响应格式 - 与后端统一
 */
export interface ApiListResponse<T> {
  success: true;
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

/**
 * API 错误响应格式 - 与后端统一
 */
export interface ApiErrorResponse {
  success: false;
  code: number;
  message: string;
  details?: any;
}

/**
 * API 响应基础类型
 * 使用联合类型严格区分成功和失败响应
 */
export type ApiResponse<T = unknown> = ApiSuccessResponse<T> | ApiErrorResponse;

/**
 * 分页查询参数
 */
export interface PageQuery {
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  keyword?: string;
  [key: string]: unknown; // 索引签名，允许额外的查询参数
}

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
 * 列表响应（不包含 success 字段的数据部分）
 */
export interface ListResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 分页结果（废弃，使用 ListResult 代替）
 * @deprecated 使用 ListResult 代替
 */
export type PageResult<T> = ListResult<T>;

/**
 * API错误类
 */
export class ApiError extends Error {
  code: number;
  details?: unknown;

  constructor(code: number, message: string, details?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }

  /**
   * 判断是否为网络错误
   */
  isNetworkError(): boolean {
    return this.code === 0 || this.code === -1;
  }

  /**
   * 判断是否为认证错误
   */
  isAuthError(): boolean {
    return this.code === 401 || this.code === 403;
  }

  /**
   * 判断是否为客户端错误
   */
  isClientError(): boolean {
    return this.code >= 400 && this.code < 500;
  }

  /**
   * 判断是否为服务器错误
   */
  isServerError(): boolean {
    return this.code >= 500;
  }

  /**
   * 获取用户友好的错误信息
   */
  getFriendlyMessage(): string {
    if (this.isNetworkError()) {
      return '网络连接失败，请检查网络设置';
    }
    if (this.isAuthError()) {
      return '登录已过期，请重新登录';
    }
    if (this.isClientError()) {
      return '请求参数错误，请检查输入信息';
    }
    if (this.isServerError()) {
      return '服务器暂时不可用，请稍后重试';
    }
    return this.message || '操作失败，请重试';
  }
}

/**
 * HTTP状态码常量
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  METHOD_NOT_ALLOWED: 405,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const;

/**
 * 通用错误码
 */
export const ERROR_CODES = {
  // 系统错误
  UNKNOWN_ERROR: -1,
  SUCCESS: 0,

  // 参数错误
  INVALID_PARAMS: 1001,
  MISSING_REQUIRED_FIELD: 1002,
  INVALID_FORMAT: 1003,

  // 认证错误
  UNAUTHORIZED: 2001,
  TOKEN_EXPIRED: 2002,
  INVALID_TOKEN: 2003,

  // 权限错误
  FORBIDDEN: 3001,
  INSUFFICIENT_PERMISSIONS: 3002,

  // 业务错误
  RESOURCE_NOT_FOUND: 4001,
  RESOURCE_ALREADY_EXISTS: 4002,
  OPERATION_NOT_ALLOWED: 4003,

  // 支付错误
  PAYMENT_FAILED: 5001,
  INSUFFICIENT_BALANCE: 5002,
  PAYMENT_TIMEOUT: 5003,
} as const;

/**
 * 请求配置
 */
export interface RequestConfig {
  timeout?: number;
  retries?: number;
  retryDelay?: number;
  signal?: AbortSignal;
}

/**
 * 批量操作结果
 */
export interface BatchResult<T = unknown> {
  success: boolean;
  total: number;
  succeeded: number;
  failed: number;
  results?: T[];
  errors?: Array<{
    index: number;
    error: string;
    data?: unknown;
  }>;
}
