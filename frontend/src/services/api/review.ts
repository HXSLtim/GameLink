import { apiClient } from '../../api/client';
import type {
  Review,
  ReviewListQuery,
  CreateReviewRequest,
  UpdateReviewRequest,
  ReviewStatistics,
} from '../../types/review';
import type { ListResult } from '../../types/api';

/**
 * 评价列表响应
 */
export type ReviewListResponse = ListResult<Review>;

/**
 * 评价API服务
 */
export const reviewApi = {
  /**
   * 获取评价列表
   */
  getList: (params: ReviewListQuery): Promise<ReviewListResponse> => {
    return apiClient.get('/api/v1/admin/reviews', { params });
  },

  /**
   * 获取评价详情
   */
  getDetail: (id: number): Promise<Review> => {
    return apiClient.get(`/api/v1/admin/reviews/${id}`);
  },

  /**
   * 创建评价
   */
  create: (data: CreateReviewRequest): Promise<Review> => {
    return apiClient.post('/api/v1/admin/reviews', data);
  },

  /**
   * 更新评价
   */
  update: (id: number, data: UpdateReviewRequest): Promise<Review> => {
    return apiClient.put(`/api/v1/admin/reviews/${id}`, data);
  },

  /**
   * 删除评价
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/reviews/${id}`);
  },

  /**
   * 获取陪玩师的评价列表
   */
  getPlayerReviews: (
    playerId: number,
    params?: { page?: number; page_size?: number },
  ): Promise<ReviewListResponse> => {
    return apiClient.get(`/api/v1/admin/players/${playerId}/reviews`, { params });
  },

  /**
   * 获取评价操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/reviews/${id}/logs`);
  },

  /**
   * 获取评价统计
   */
  getStatistics: (playerId?: number): Promise<ReviewStatistics> => {
    const params = playerId ? { player_id: playerId } : undefined;
    return apiClient.get('/api/v1/admin/stats/reviews', { params });
  },
};

