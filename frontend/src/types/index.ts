/**
 * 统一类型导出文件
 * 与后端Go模型保持完全一致
 */

// 基础类型
export * from './api';
export * from './auth';

// 实体类型
export * from './user';
export * from './game';
export * from './player';
export * from './order';
export * from './payment';

// 常用类型别名
export type ID = number;
export type Timestamp = string;
export type Optional<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;
export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>;

// 响应类型扩展
export interface SuccessResponse<T> {
  success: true;
  data: T;
  message?: string;
}

export interface ErrorResponse {
  success: false;
  code: number;
  message: string;
  details?: any;
}

export type ApiResponse<T> = SuccessResponse<T> | ErrorResponse;

// 列表响应类型
export interface ListResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  has_next: boolean;
  has_prev: boolean;
}

// 统计数据类型
export interface DashboardStats {
  user_count: number;
  player_count: number;
  order_count: number;
  revenue_cents: number;
  avg_rating: number;
  completion_rate: number;
}

// 搜索和筛选通用类型
export interface DateRange {
  start?: string;
  end?: string;
}

export interface PaginationParams {
  page?: number;
  page_size?: number;
}

export interface SortParams {
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface SearchParams extends PaginationParams, SortParams {
  keyword?: string;
}

// 批量操作类型
export interface BatchOperation<T> {
  action: 'create' | 'update' | 'delete';
  items: T[];
}

export interface BatchResult {
  success_count: number;
  failed_count: number;
  errors?: Array<{
    index: number;
    error: string;
  }>;
}

// 表单通用类型
export interface FormField {
  name: string;
  label: string;
  type: 'text' | 'email' | 'password' | 'number' | 'select' | 'textarea' | 'date' | 'datetime';
  required?: boolean;
  placeholder?: string;
  options?: Array<{ label: string; value: any }>;
  validation?: {
    min?: number;
    max?: number;
    pattern?: string;
    message?: string;
  };
}

export interface FormConfig {
  fields: FormField[];
  submitText?: string;
  resetText?: string;
  validation?: Record<string, any>;
}

// 图表数据类型
export interface ChartDataPoint {
  x: string | number;
  y: number;
  label?: string;
}

export interface ChartSeries {
  name: string;
  data: ChartDataPoint[];
  color?: string;
}

export interface ChartConfig {
  type: 'line' | 'bar' | 'pie' | 'area';
  title?: string;
  xAxis?: string;
  yAxis?: string;
  series: ChartSeries[];
}

// 导出数据类型
export interface ExportParams {
  format: 'csv' | 'excel' | 'json';
  fields?: string[];
  filters?: Record<string, any>;
  date_range?: DateRange;
}

export interface ExportResult {
  filename: string;
  url: string;
  size: number;
  created_at: string;
  expires_at: string;
}

// 通知类型
export type NotificationType = 'info' | 'success' | 'warning' | 'error';

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  duration?: number;
  action?: {
    label: string;
    callback: () => void;
  };
}

// 文件上传类型
export interface FileUpload {
  id: string;
  name: string;
  size: number;
  type: string;
  url?: string;
  status: 'uploading' | 'success' | 'error';
  progress?: number;
  error?: string;
}

export interface UploadConfig {
  accept?: string[];
  maxSize?: number;
  maxCount?: number;
  autoUpload?: boolean;
}

// 错误类型扩展
export interface ValidationError {
  field: string;
  message: string;
  code: string;
}

export interface BusinessError extends Error {
  code: string;
  details?: any;
}

// 工具类型
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

export type DeepRequired<T> = {
  [P in keyof T]-?: T[P] extends object ? DeepRequired<T[P]> : T[P];
};

export type PickByType<T, V> = Pick<T, { [K in keyof T]: T[K] extends V ? K : never }[keyof T]>;

export type OmitByType<T, V> = Omit<T, { [K in keyof T]: T[K] extends V ? K : never }[keyof T]>;

// 环境类型
export type Environment = 'development' | 'staging' | 'production';

export interface AppConfig {
  environment: Environment;
  apiBaseUrl: string;
  version: string;
  features: Record<string, boolean>;
  limits: {
    maxFileSize: number;
    maxUploadCount: number;
    pageSize: number;
  };
}
