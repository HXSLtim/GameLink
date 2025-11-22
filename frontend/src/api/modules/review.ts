/**
 * 评价API模块
 */

import { apiClient } from '@/api/client';
import type {
  Review,
  CreateReviewRequest,
  UpdateReviewRequest,
  GetReviewsParams,
  GetReviewsResponse,
  ReviewStats,
} from '@/shared/types/review';

/**
 * 创建评价
 */
export const createReview = async (data: CreateReviewRequest): Promise<Review> => {
  const response = await apiClient.post<{ data: Review }>('/reviews', data);
  return response.data;
};

/**
 * 获取评价列表
 */
export const getReviews = async (params?: GetReviewsParams): Promise<GetReviewsResponse> => {
  const response = await apiClient.get<{ data: GetReviewsResponse }>('/reviews', { params });
  return response.data;
};

/**
 * 获取评价详情
 */
export const getReviewById = async (id: number): Promise<Review> => {
  const response = await apiClient.get<{ data: Review }>(`/reviews/${id}`);
  return response.data;
};

/**
 * 更新评价
 */
export const updateReview = async (id: number, data: UpdateReviewRequest): Promise<Review> => {
  const response = await apiClient.put<{ data: Review }>(`/reviews/${id}`, data);
  return response.data;
};

/**
 * 删除评价
 */
export const deleteReview = async (id: number): Promise<void> => {
  await apiClient.delete(`/reviews/${id}`);
};

/**
 * 获取用户收到的评价
 */
export const getUserReviews = async (
  userId: number,
  params?: Omit<GetReviewsParams, 'revieweeId'>
): Promise<GetReviewsResponse> => {
  const response = await apiClient.get<{ data: GetReviewsResponse }>(`/users/${userId}/reviews`, {
    params,
  });
  return response.data;
};

/**
 * 获取陪玩师收到的评价
 */
export const getPlayerReviews = async (
  playerId: number,
  params?: Omit<GetReviewsParams, 'revieweeId'>
): Promise<GetReviewsResponse> => {
  const response = await apiClient.get<{ data: GetReviewsResponse }>(`/players/${playerId}/reviews`, {
    params,
  });
  return response.data;
};

/**
 * 获取订单评价
 */
export const getOrderReviews = async (orderId: number): Promise<Review[]> => {
  const response = await apiClient.get<{ data: Review[] }>(`/orders/${orderId}/reviews`);
  return response.data;
};

/**
 * 获取评价统计
 */
export const getReviewStats = async (params?: {
  userId?: number;
  playerId?: number;
}): Promise<ReviewStats> => {
  const response = await apiClient.get<{ data: ReviewStats }>('/reviews/stats', { params });
  return response.data;
};