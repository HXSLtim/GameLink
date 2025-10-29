import type { BaseEntity } from './user';

/**
 * 游戏分类枚举
 */
export type GameCategory =
  | 'moba'
  | 'fps'
  | 'rpg'
  | 'strategy'
  | 'sports'
  | 'racing'
  | 'puzzle'
  | 'other';

/**
 * 游戏实体 - 与后端 model.Game 保持一致
 */
export interface Game extends BaseEntity {
  key: string;
  name: string;
  category: string;
  iconUrl?: string;
  description?: string;
}

/**
 * 游戏列表查询参数
 */
export interface GameListQuery {
  page?: number;
  pageSize?: number;
  category?: string;
  keyword?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'name';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 创建游戏请求
 */
export interface CreateGameRequest {
  key: string;
  name: string;
  category: string;
  iconUrl?: string;
  description?: string;
}

/**
 * 更新游戏请求
 */
export interface UpdateGameRequest {
  key?: string;
  name?: string;
  category?: string;
  iconUrl?: string;
  description?: string;
}

/**
 * 游戏分类显示文本
 */
export const GAME_CATEGORY_TEXT: Record<GameCategory, string> = {
  moba: 'MOBA',
  fps: '射击',
  rpg: '角色扮演',
  strategy: '策略',
  sports: '体育',
  racing: '竞速',
  puzzle: '益智',
  other: '其他',
};

/**
 * 常见游戏预设
 */
export const POPULAR_GAMES: Array<{
  key: string;
  name: string;
  category: GameCategory;
  iconUrl?: string;
}> = [
  { key: 'lol', name: '英雄联盟', category: 'moba' },
  { key: 'dota2', name: 'Dota 2', category: 'moba' },
  { key: 'csgo', name: 'CS:GO', category: 'fps' },
  { key: 'valorant', name: 'Valorant', category: 'fps' },
  { key: 'pubg', name: '绝地求生', category: 'fps' },
  { key: 'king_glory', name: '王者荣耀', category: 'moba' },
  { key: 'peace', name: '和平精英', category: 'fps' },
  { key: 'genshin', name: '原神', category: 'rpg' },
];
