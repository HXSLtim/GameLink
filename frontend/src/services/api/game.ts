import { apiClient } from '../../api/client';
import type {
  Game,
  GameDetail,
  GameListQuery,
  CreateGameRequest,
  UpdateGameRequest,
  BulkGameRequest,
  BulkGameResult,
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
   * 获取游戏详情（包含统计信息）
   */
  getDetail: (id: number): Promise<GameDetail> => {
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
   * 批量操作游戏
   * @param data 批量操作请求
   * @returns 批量操作结果
   */
  bulkOperation: (data: BulkGameRequest): Promise<BulkGameResult> => {
    return apiClient.post('/api/v1/admin/games/bulk', data);
  },

  /**
   * 批量删除游戏
   * @param ids 游戏ID数组
   */
  bulkDelete: (ids: number[]): Promise<BulkGameResult> => {
    return apiClient.post('/api/v1/admin/games/bulk', {
      ids,
      action: 'delete',
    });
  },

  /**
   * 批量激活游戏
   * @param ids 游戏ID数组
   */
  bulkActivate: (ids: number[]): Promise<BulkGameResult> => {
    return apiClient.post('/api/v1/admin/games/bulk', {
      ids,
      action: 'activate',
    });
  },

  /**
   * 批量下架游戏
   * @param ids 游戏ID数组
   */
  bulkDeactivate: (ids: number[]): Promise<BulkGameResult> => {
    return apiClient.post('/api/v1/admin/games/bulk', {
      ids,
      action: 'deactivate',
    });
  },

  /**
   * 获取游戏操作日志
   */
  getLogs: (id: number): Promise<unknown[]> => {
    return apiClient.get(`/api/v1/admin/games/${id}/logs`);
  },

  /**
   * 获取游戏统计信息
   */
  getStats: (id: number): Promise<unknown> => {
    return apiClient.get(`/api/v1/admin/games/${id}/stats`);
  },
};
