/**
 * 游戏相关接口定义
 */

import type { BaseEntity } from '@/shared/types/api';

/**
 * 游戏分类
 */
export enum GameCategory {
  MOBA = 'moba',
  FPS = 'fps',
  RPG = 'rpg',
  CARD = 'card',
  SPORTS = 'sports',
  OTHER = 'other',
}

/**
 * 游戏实体
 */
export interface Game extends BaseEntity {
  id: number;
  name: string;
  icon?: string;
  category: GameCategory;
  description?: string;
  status: 'active' | 'inactive';
  popularity: number;
}

/**
 * 游戏服务类型
 */
export interface GameService {
  id: number;
  gameId: number;
  name: string;
  description?: string;
  price: number;
  duration: number;
  status: 'active' | 'inactive';
}

/**
 * 游戏列表查询参数
 */
export interface GetGamesParams {
  category?: GameCategory;
  status?: 'active' | 'inactive';
  keyword?: string;
  page?: number;
  pageSize?: number;
}

/**
 * 游戏列表响应
 */
export interface GetGamesResponse {
  list: Game[];
  total: number;
  page: number;
  pageSize: number;
}