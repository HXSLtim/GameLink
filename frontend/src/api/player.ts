/**
 * 陪玩师相关API
 */
import { apiClient } from './client';

/**
 * 陪玩师信息接口
 */
export interface Player {
  id: number;
  userId: number;
  gameId: number;
  game: {
    id: number;
    name: string;
  };
  user: {
    id: number;
    username: string;
    avatar?: string;
  };
  bio?: string;
  pricePerHour: number;
  status: 'available' | 'busy' | 'offline';
  rating: number;
  reviewCount: number;
  totalOrders: number;
  tags?: string[];
  createdAt: string;
  updatedAt: string;
}

/**
 * 陪玩师列表查询参数
 */
export interface GetPlayersParams {
  /**
   * 游戏ID
   */
  gameId?: number;

  /**
   * 状态筛选
   */
  status?: 'available' | 'busy' | 'offline';

  /**
   * 最低价格
   */
  minPrice?: number;

  /**
   * 最高价格
   */
  maxPrice?: number;

  /**
   * 搜索关键词
   */
  keyword?: string;

  /**
   * 排序字段
   */
  sortBy?: 'rating' | 'price' | 'reviews' | 'orders';

  /**
   * 排序方向
   */
  sortOrder?: 'asc' | 'desc';

  /**
   * 页码
   */
  page?: number;

  /**
   * 每页数量
   */
  pageSize?: number;
}

/**
 * 陪玩师列表响应
 */
export interface GetPlayersResponse {
  players: Player[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 获取陪玩师列表
 */
export const getPlayers = async (
  params?: GetPlayersParams
): Promise<GetPlayersResponse> => {
  const response = await apiClient.get<GetPlayersResponse>('/players', {
    params,
  });
  return response.data;
};

/**
 * 获取陪玩师详情
 */
export const getPlayerById = async (id: number): Promise<Player> => {
  const response = await apiClient.get<Player>(`/players/${id}`);
  return response.data;
};

/**
 * 获取推荐陪玩师
 */
export const getRecommendedPlayers = async (
  limit: number = 10
): Promise<Player[]> => {
  const response = await apiClient.get<Player[]>('/players/recommended', {
    params: { limit },
  });
  return response.data;
};
