/**
 * Public API Service
 * 公开 API（无需登录）
 */

import apiClient from './request';
import type {
  ApiResponse,
  PaginatedResponse,
  Player,
  Game,
  ServiceItem,
  SearchResponse,
} from '@/types/api';

/**
 * 搜索
 */
export const search = (params: {
  q: string;
  type?: 'all' | 'player' | 'game';
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<SearchResponse>> => {
  return apiClient.get('/public/search', params, { skipAuth: true });
};

/**
 * 获取游戏列表
 */
export const getGames = (params?: {
  page?: number;
  pageSize?: number;
  isActive?: boolean;
}): Promise<ApiResponse<PaginatedResponse<Game>>> => {
  return apiClient.get('/public/games', params, { skipAuth: true });
};

/**
 * 获取游戏详情
 */
export const getGameDetail = (gameId: number): Promise<ApiResponse<Game>> => {
  return apiClient.get(`/public/games/${gameId}`, undefined, { skipAuth: true });
};

/**
 * 获取陪玩师列表
 */
export const getPlayers = (params?: {
  page?: number;
  pageSize?: number;
  gameId?: number;
  rank?: string;
  minRating?: number;
  maxPrice?: number;
  keyword?: string;
}): Promise<ApiResponse<PaginatedResponse<Player>>> => {
  return apiClient.get('/public/players', params, { skipAuth: true });
};

/**
 * 获取陪玩师详情
 */
export const getPlayerDetail = (playerId: number): Promise<ApiResponse<Player>> => {
  return apiClient.get(`/public/players/${playerId}`, undefined, { skipAuth: true });
};

/**
 * 获取游戏的服务项目
 */
export const getServiceItems = (gameId: number, params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<ServiceItem>>> => {
  return apiClient.get(`/public/games/${gameId}/services`, params, { skipAuth: true });
};

/**
 * 获取陪玩师的评价
 */
export const getPlayerReviews = (playerId: number, params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<any>>> => {
  return apiClient.get(`/public/players/${playerId}/reviews`, params, { skipAuth: true });
};

/**
 * 获取排行榜
 */
export const getRanking = (params?: {
  type?: 'daily' | 'weekly' | 'monthly';
  gameId?: number;
}): Promise<ApiResponse<Player[]>> => {
  return apiClient.get('/public/ranking', params, { skipAuth: true });
};

/**
 * 获取系统配置
 */
export const getSystemConfig = (): Promise<ApiResponse<{
  minOrderAmount: number;
  maxOrderAmount: number;
  supportUrl: string;
  agreementUrls: {
    privacy: string;
    terms: string;
  };
}>> => {
  return apiClient.get('/public/config', undefined, { skipAuth: true });
};

/**
 * 获取活动列表
 */
export const getActivities = (params?: {
  page?: number;
  pageSize?: number;
  isActive?: boolean;
}): Promise<ApiResponse<PaginatedResponse<any>>> => {
  return apiClient.get('/public/activities', params, { skipAuth: true });
};

/**
 * 获取活动详情
 */
export const getActivityDetail = (activityId: number): Promise<ApiResponse<any>> => {
  return apiClient.get(`/public/activities/${activityId}`, undefined, { skipAuth: true });
};

/**
 * 获取可领取的优惠券
 */
export const getAvailableCoupons = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<any>>> => {
  return apiClient.get('/public/coupons', params, { skipAuth: true });
};
