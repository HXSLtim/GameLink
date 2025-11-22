/**
 * 陪玩师API模块
 */

import { apiClient } from '@/api/client';
import type { Player, GetPlayersParams, GetPlayersResponse } from '@/shared/types/user';

/**
 * 获取陪玩师列表
 */
export const getPlayers = async (params?: GetPlayersParams): Promise<GetPlayersResponse> => {
  const response = await apiClient.get<{ data: GetPlayersResponse }>('/players', { params });
  return response.data;
};

/**
 * 获取陪玩师详情
 */
export const getPlayerById = async (id: number): Promise<Player> => {
  const response = await apiClient.get<{ data: Player }>(`/players/${id}`);
  return response.data;
};

/**
 * 获取推荐陪玩师
 */
export const getRecommendedPlayers = async (limit: number = 10): Promise<Player[]> => {
  const response = await apiClient.get<{ data: Player[] }>('/players/recommended', {
    params: { limit },
  });
  return response.data;
};

/**
 * 获取游戏陪玩师
 */
export const getGamePlayers = async (gameId: number, params?: Omit<GetPlayersParams, 'gameId'>): Promise<GetPlayersResponse> => {
  const response = await apiClient.get<{ data: GetPlayersResponse }>(`/games/${gameId}/players`, {
    params,
  });
  return response.data;
};

/**
 * 搜索陪玩师
 */
export const searchPlayers = async (keyword: string, params?: Omit<GetPlayersParams, 'keyword'>): Promise<GetPlayersResponse> => {
  const response = await apiClient.get<{ data: GetPlayersResponse }>('/players/search', {
    params: { keyword, ...params },
  });
  return response.data;
};

/**
 * 获取陪玩师统计
 */
export const getPlayerStats = async (playerId: number): Promise<{
  totalOrders: number;
  completedOrders: number;
  totalEarnings: number;
  averageRating: number;
  reviewCount: number;
}> => {
  const response = await apiClient.get<{ data: any }>(`/players/${playerId}/stats`);
  return response.data;
};
