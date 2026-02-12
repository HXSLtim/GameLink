/**
 * 收藏相关 API
 */

import { get, post, del, type RequestConfig } from './request'

// 收藏的陪玩师
export interface FavoritePlayer {
  id: number
  playerId: number
  playerName: string
  playerAvatar?: string
  playerRating: number
  nickname?: string
  avatar?: string
  rating?: number
  orderCount?: number
  minPrice?: number
  games?: Array<string | { id?: number; name: string }>
  isOnline: boolean
  gameNames: string[]
  createdAt: string
}

/**
 * 获取收藏列表
 */
export function getFavorites(
  params?: { page?: number; page_size?: number },
  config?: Partial<RequestConfig>
) {
  return get<FavoritePlayer[]>('/user/favorites/players', params, config)
}

/**
 * 添加收藏
 */
export function addFavorite(playerId: number) {
  return post<void>(`/user/favorites/players/${playerId}`)
}

/**
 * 取消收藏
 */
export function removeFavorite(playerId: number) {
  return del<void>(`/user/favorites/players/${playerId}`)
}

/**
 * 检查是否已收藏
 */
export function checkFavorite(playerId: number) {
  return get<{ isFavorite?: boolean; isFavorited?: boolean }>(`/user/favorites/players/${playerId}/check`)
}

/**
 * 批量取消收藏
 */
export function batchRemoveFavorites(playerIds: number[]) {
  return Promise.all(playerIds.map(playerId => del<void>(`/user/favorites/players/${playerId}`)))
}

export default {
  getFavorites,
  addFavorite,
  removeFavorite,
  checkFavorite,
  batchRemoveFavorites,
}
