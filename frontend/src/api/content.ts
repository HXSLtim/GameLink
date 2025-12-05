/**
 * 内容管理 API
 * 包含动态审核、聊天监控、举报管理、内容分类等功能
 */
import apiClient from './client';
import type {
  Feed,
  FeedQueryParams,
  ChatMessage,
  ChatMessageQueryParams,
  MuteUserRequest,
  FeedReport,
  FeedReportQueryParams,
  ProcessReportRequest,
  ContentCategory,
  ContentCategoryQueryParams,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  ContentStatsDTO,
} from '@/types/content';

// 通用响应类型
export interface ApiResponse<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  pagination?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

// ==================== 动态审核 API ====================

export const feedApi = {
  // 获取动态列表
  getFeeds: (params?: FeedQueryParams) =>
    apiClient.get<ApiResponse<{ items: Feed[]; total: number }>>('/admin/content/feeds', { params }),

  // 获取动态详情
  getFeed: (id: number) =>
    apiClient.get<ApiResponse<Feed>>(`/admin/content/feeds/${id}`),

  // 批准动态
  approveFeed: (id: number, note?: string) =>
    apiClient.put<ApiResponse<void>>(`/admin/content/feeds/${id}/approve`, { note }),

  // 拒绝动态
  rejectFeed: (id: number, note: string) =>
    apiClient.put<ApiResponse<void>>(`/admin/content/feeds/${id}/reject`, { note }),

  // 删除动态
  deleteFeed: (id: number, note?: string) =>
    apiClient.delete<ApiResponse<void>>(`/admin/content/feeds/${id}`, { data: { note } }),

  // 批量批准动态
  batchApproveFeed: (feedIds: number[], note?: string) =>
    apiClient.post<ApiResponse<void>>('/admin/content/feeds/batch-approve', { feedIds, note }),

  // 批量拒绝动态
  batchRejectFeed: (feedIds: number[], note: string) =>
    apiClient.post<ApiResponse<void>>('/admin/content/feeds/batch-reject', { feedIds, note }),
};

// ==================== 聊天监控 API ====================

export const chatModerationApi = {
  // 获取聊天消息列表
  getMessages: (params?: ChatMessageQueryParams) =>
    apiClient.get<ApiResponse<{ items: ChatMessage[]; total: number }>>('/admin/content/chat/messages', { params }),

  // 删除消息
  deleteMessage: (id: number, reason?: string) =>
    apiClient.delete<ApiResponse<void>>(`/admin/content/chat/messages/${id}`, { data: { reason } }),

  // 禁言用户
  muteUser: (data: MuteUserRequest) =>
    apiClient.post<ApiResponse<void>>('/admin/content/chat/mute', data),

  // 解除禁言
  unmuteUser: (groupId: number, userId: number) =>
    apiClient.post<ApiResponse<void>>('/admin/content/chat/unmute', null, { params: { groupId, userId } }),
};

// ==================== 动态举报 API ====================

export const feedReportApi = {
  // 获取举报列表
  getReports: (params?: FeedReportQueryParams) =>
    apiClient.get<ApiResponse<{ items: FeedReport[]; total: number }>>('/admin/content/reports', { params }),

  // 获取举报详情
  getReport: (id: number) =>
    apiClient.get<ApiResponse<FeedReport>>(`/admin/content/reports/${id}`),

  // 处理举报
  processReport: (id: number, data: ProcessReportRequest) =>
    apiClient.post<ApiResponse<void>>(`/admin/content/reports/${id}/process`, data),
};

// ==================== 内容分类 API ====================

export const contentCategoryApi = {
  // 获取分类列表
  getCategories: (params?: ContentCategoryQueryParams) =>
    apiClient.get<ApiResponse<{ items: ContentCategory[]; total: number }>>('/admin/content/categories', { params }),

  // 获取分类详情
  getCategory: (id: number) =>
    apiClient.get<ApiResponse<ContentCategory>>(`/admin/content/categories/${id}`),

  // 创建分类
  createCategory: (data: CreateCategoryRequest) =>
    apiClient.post<ApiResponse<ContentCategory>>('/admin/content/categories', data),

  // 更新分类
  updateCategory: (id: number, data: UpdateCategoryRequest) =>
    apiClient.put<ApiResponse<void>>(`/admin/content/categories/${id}`, data),

  // 删除分类
  deleteCategory: (id: number, migrateToCategoryId?: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/content/categories/${id}`, {
      params: migrateToCategoryId ? { migrateToCategoryId } : undefined,
    }),
};

// ==================== 内容统计 API ====================

export const contentStatsApi = {
  // 获取统计数据
  getStats: (days?: number) =>
    apiClient.get<ApiResponse<ContentStatsDTO>>('/admin/content/stats', { params: { days } }),
};
