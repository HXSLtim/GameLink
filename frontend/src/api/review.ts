/**
 * 评价管理 API
 * 包含评价列表、审核、举报、回复等功能
 */
import apiClient from './client';
import type {
  Review,
  ReviewQueryParams,
  ReviewReport,
  ReviewReportQueryParams,
  SensitiveWord,
  SensitiveWordQueryParams,
  ReviewStats,
  ReviewTrend,
  TopPlayer,
  GameStats,
  ReviewDisplaySettings,
  DetectSensitiveWordsResult,
  ReviewReply,
} from '@/types/review';

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

// ==================== 评价管理 API ====================

export const reviewApi = {
  // 获取评价列表
  getReviews: (params?: ReviewQueryParams) =>
    apiClient.get<ApiResponse<Review[]>>('/admin/reviews', { params }),

  // 获取评价详情
  getReview: (id: number) =>
    apiClient.get<ApiResponse<Review>>(`/admin/reviews/${id}`),

  // 获取待审核评价列表
  getPendingReviews: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<Review[]>>('/admin/reviews/pending', { params }),

  // 批准评价
  approveReview: (id: number) =>
    apiClient.put<ApiResponse<{ message: string }>>(`/admin/reviews/${id}/approve`),

  // 拒绝评价
  rejectReview: (id: number, reason: string) =>
    apiClient.put<ApiResponse<{ message: string }>>(`/admin/reviews/${id}/reject`, { reason }),

  // 批量批准评价
  batchApproveReviews: (ids: number[]) =>
    apiClient.put<ApiResponse<{ message: string; count: number }>>('/admin/reviews/batch-approve', { ids }),

  // 批量拒绝评价
  batchRejectReviews: (ids: number[], reason: string) =>
    apiClient.put<ApiResponse<{ message: string; count: number }>>('/admin/reviews/batch-reject', { ids, reason }),

  // 删除评价
  deleteReview: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/reviews/${id}`),

  // 获取评价操作日志
  getReviewLogs: (id: number) =>
    apiClient.get<ApiResponse<OperationLog[]>>(`/admin/reviews/${id}/logs`),

  // 获取陪玩师的评价列表
  getPlayerReviews: (playerId: number, params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<Review[]>>(`/admin/players/${playerId}/reviews`, { params }),
};

// ==================== 举报管理 API ====================

export const reviewReportApi = {
  // 获取举报列表
  getReports: (params?: ReviewReportQueryParams) =>
    apiClient.get<ApiResponse<ReviewReport[]>>('/admin/review-reports', { params }),

  // 获取举报详情
  getReport: (id: number) =>
    apiClient.get<ApiResponse<ReviewReport>>(`/admin/review-reports/${id}`),

  // 创建举报
  createReport: (reviewId: number, data: { reason: string; evidence?: string }) =>
    apiClient.post<ApiResponse<{ reportId: number }>>(`/admin/reviews/${reviewId}/reports`, data),

  // 处理举报
  handleReport: (id: number, data: { action: 'delete' | 'warn' | 'reject'; note?: string }) =>
    apiClient.put<ApiResponse<{ message: string }>>(`/admin/review-reports/${id}/handle`, data),
};

// ==================== 敏感词管理 API ====================

export const sensitiveWordApi = {
  // 获取敏感词列表
  getWords: (params?: SensitiveWordQueryParams) =>
    apiClient.get<ApiResponse<{ items: SensitiveWord[]; total: number }>>('/admin/sensitive-words', { params }),

  // 添加敏感词
  addWord: (data: { word: string; category: string; severity: string }) =>
    apiClient.post<ApiResponse<SensitiveWord>>('/admin/sensitive-words', data),

  // 更新敏感词
  updateWord: (id: number, data: { word?: string; category?: string; severity?: string }) =>
    apiClient.put<ApiResponse<SensitiveWord>>(`/admin/sensitive-words/${id}`, data),

  // 删除敏感词
  deleteWord: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/sensitive-words/${id}`),

  // 检测敏感词
  detectWords: (content: string) =>
    apiClient.post<ApiResponse<DetectSensitiveWordsResult>>('/admin/reviews/detect-sensitive', { content }),
};

// ==================== 统计分析 API ====================

export const reviewStatsApi = {
  // 获取统计概览
  getStats: () =>
    apiClient.get<ApiResponse<ReviewStats>>('/admin/reviews/stats'),

  // 获取评价趋势
  getTrend: (params?: { days?: number }) =>
    apiClient.get<ApiResponse<ReviewTrend[]>>('/admin/reviews/trend', { params }),

  // 获取陪玩师排行
  getTopPlayers: (params?: { limit?: number; sortBy?: 'count' | 'rating' }) =>
    apiClient.get<ApiResponse<TopPlayer[]>>('/admin/reviews/top-players', { params }),

  // 获取游戏统计
  getGameStats: () =>
    apiClient.get<ApiResponse<GameStats[]>>('/admin/reviews/game-stats'),

  // 导出统计数据
  exportStats: (type: 'overview' | 'trend' | 'players' | 'games') => {
    const url = `/admin/reviews/export?type=${type}`;
    return apiClient.get(url, { responseType: 'blob' });
  },
};

// ==================== 回复管理 API ====================

export const reviewReplyApi = {
  // 创建回复
  createReply: (reviewId: number, content: string) =>
    apiClient.post<ApiResponse<ReviewReply>>(`/admin/reviews/${reviewId}/reply`, { content }),

  // 更新回复
  updateReply: (id: number, content: string) =>
    apiClient.put<ApiResponse<ReviewReply>>(`/admin/review-replies/${id}`, { content }),

  // 删除回复
  deleteReply: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/review-replies/${id}`),
};

// ==================== 展示设置 API ====================

export const reviewSettingsApi = {
  // 获取展示设置
  getSettings: () =>
    apiClient.get<ApiResponse<ReviewDisplaySettings>>('/admin/review-settings'),

  // 更新展示设置
  updateSettings: (data: Partial<ReviewDisplaySettings>) =>
    apiClient.put<ApiResponse<ReviewDisplaySettings>>('/admin/review-settings', data),
};

// ==================== 操作日志类型 ====================

export interface OperationLog {
  id: number;
  entityType: string;
  entityId: number;
  action: string;
  actorId: number;
  actorName?: string;
  beforeState?: string;
  afterState?: string;
  note?: string;
  createdAt: string;
}

// 操作日志 API
export const operationLogApi = {
  // 搜索操作日志
  searchLogs: (params?: {
    page?: number;
    pageSize?: number;
    entityType?: string;
    entityId?: number;
    actorId?: number;
    action?: string;
    startTime?: string;
    endTime?: string;
  }) => apiClient.get<ApiResponse<OperationLog[]>>('/admin/operation-logs', { params }),

  // 导出操作日志
  exportLogs: (params?: {
    entityType?: string;
    entityId?: number;
    actorId?: number;
    action?: string;
    startTime?: string;
    endTime?: string;
  }) => apiClient.get('/admin/operation-logs/export', { params, responseType: 'blob' }),
};
