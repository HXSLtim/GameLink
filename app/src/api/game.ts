/**
 * 游戏相关 API
 */

import { get, type RequestConfig } from './request'

// 游戏分类
export interface GameCategory {
  id: number
  name: string
  icon?: string
  sort: number
}

// 游戏
export interface Game {
  id: number
  name: string
  icon?: string
  coverImage?: string
  categoryId: number
  categoryName: string
  playerCount: number
  isHot: boolean
  sort: number
  minPrice?: number
  maxPrice?: number
}

// 段位
export interface GameRank {
  id: number
  gameId: number
  name: string
  icon?: string
  level: number
  sort: number
}

/**
 * 获取游戏分类列表
 */
export function getGameCategories(config?: Partial<RequestConfig>) {
  return get<GameCategory[]>('/public/games/categories', undefined, config)
}

/**
 * 获取游戏列表
 */
export function getGames(params?: {
  categoryId?: number
  keyword?: string
  page?: number
  page_size?: number
}, config?: Partial<RequestConfig>) {
  return get<Game[]>('/public/games', params, config)
}

/**
 * 获取游戏详情
 */
export function getGameDetail(gameId: number, config?: Partial<RequestConfig>) {
  return get<Game>(`/public/games/${gameId}`, undefined, config)
}

/**
 * 获取游戏段位列表
 */
export function getGameRanks(gameId: number, config?: Partial<RequestConfig>) {
  return get<GameRank[]>(`/public/games/${gameId}/ranks`, undefined, config)
}

/**
 * 获取热门游戏（使用游戏列表接口，按热门排序）
 */
export function getHotGames(limit = 10, config?: Partial<RequestConfig>) {
  return get<Game[]>('/public/games', { page_size: limit, page: 1 }, config)
}

export default {
  getGameCategories,
  getGames,
  getGameDetail,
  getGameRanks,
  getHotGames,
}
