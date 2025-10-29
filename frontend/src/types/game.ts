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
 * 游戏状态枚举
 */
export type GameStatus = 'active' | 'inactive' | 'maintenance';

/**
 * 游戏实体 - 与后端 model.Game 保持一致
 */
export interface Game extends BaseEntity {
  key: string;
  name: string;
  category: string;
  iconUrl?: string;
  description?: string;
  status?: GameStatus; // 游戏状态
  tags?: string[]; // 游戏标签
  playerCount?: number; // 玩家数量（统计）
  orderCount?: number; // 订单数量（统计）
  popularityScore?: number; // 热度评分
}

/**
 * 游戏详情（包含扩展信息）
 */
export interface GameDetail extends Game {
  // 统计信息
  stats?: {
    totalPlayers: number;
    totalOrders: number;
    totalRevenue: number;
    avgRating: number;
  };
  // 热门陪玩师
  topPlayers?: Array<{
    id: number;
    nickname: string;
    avatarUrl?: string;
    rating: number;
  }>;
}

/**
 * 游戏列表查询参数
 */
export interface GameListQuery {
  page?: number;
  pageSize?: number;
  category?: string;
  status?: GameStatus;
  keyword?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'name' | 'popularityScore';
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
  status?: GameStatus;
  tags?: string[];
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
  status?: GameStatus;
  tags?: string[];
}

/**
 * 批量操作请求
 */
export interface BulkGameRequest {
  ids: number[];
  action: 'delete' | 'activate' | 'deactivate' | 'maintenance';
}

/**
 * 批量操作结果
 */
export interface BulkGameResult {
  success: number;
  failed: number;
  errors?: Array<{ id: number; error: string }>;
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
 * 游戏状态显示文本
 */
export const GAME_STATUS_TEXT: Record<GameStatus, string> = {
  active: '正常',
  inactive: '已下架',
  maintenance: '维护中',
};

/**
 * 游戏状态颜色
 */
export const GAME_STATUS_COLOR: Record<GameStatus, string> = {
  active: 'green',
  inactive: 'default',
  maintenance: 'orange',
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
