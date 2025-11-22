/**
 * 游戏API模块
 */

import { apiClient } from '@/api/client';
import type {
  Game,
  GameCategory,
  GameService,
  GetGamesParams,
  GetGamesResponse,
} from '@/shared/types/game';

/**
 * 获取游戏列表
 */
export const getGames = async (params?: GetGamesParams): Promise<GetGamesResponse> => {
  const response = await apiClient.get<{ data: GetGamesResponse }>('/games', { params });
  return response.data;
};

/**
 * 获取游戏详情
 */
export const getGameById = async (id: number): Promise<Game> => {
  const response = await apiClient.get<{ data: Game }>(`/games/${id}`);
  return response.data;
};

/**
 * 获取热门游戏
 */
export const getPopularGames = async (limit: number = 10): Promise<Game[]> => {
  const response = await apiClient.get<{ data: Game[] }>('/games/popular', {
    params: { limit },
  });
  return response.data;
};

/**
 * 获取游戏分类列表
 */
export const getGameCategories = async (): Promise<{ label: string; value: GameCategory }[]> => {
  const response = await apiClient.get<{ data: { label: string; value: GameCategory }[] }>(
    '/games/categories'
  );
  return response.data;
};

/**
 * 获取游戏服务类型
 */
export const getGameServices = async (gameId: number): Promise<GameService[]> => {
  const response = await apiClient.get<{ data: GameService[] }>(`/games/${gameId}/services`);
  return response.data;
};

/**
 * 搜索游戏
 */
export const searchGames = async (keyword: string): Promise<Game[]> => {
  const response = await apiClient.get<{ data: Game[] }>('/games/search', {
    params: { keyword },
  });
  return response.data;
};