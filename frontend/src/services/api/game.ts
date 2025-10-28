import { apiClient } from '../../api/client';
import type {
  Game,
  GameListQuery,
  CreateGameRequest,
  UpdateGameRequest,
} from '../../types/game';
import type { ListResult } from '../../types/api';

/**
 * 游戏列表响应
 */
export type GameListResponse = ListResult<Game>;

/**
 * 游戏API服务
 */
export const gameApi = {
  /**
   * 获取游戏列表
   */
  getList: (params: GameListQuery): Promise<GameListResponse> => {
    return apiClient.get('/api/v1/admin/games', { params });
  },

  /**
   * 获取游戏详情
   */
  getDetail: (id: number): Promise<Game> => {
    return apiClient.get(`/api/v1/admin/games/${id}`);
  },

  /**
   * 创建游戏
   */
  create: (data: CreateGameRequest): Promise<Game> => {
    return apiClient.post('/api/v1/admin/games', data);
  },

  /**
   * 更新游戏信息
   */
  update: (id: number, data: UpdateGameRequest): Promise<Game> => {
    return apiClient.put(`/api/v1/admin/games/${id}`, data);
  },

  /**
   * 删除游戏
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete(`/api/v1/admin/games/${id}`);
  },

  /**
   * 获取游戏操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/games/${id}/logs`);
  },
};

