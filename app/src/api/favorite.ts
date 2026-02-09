/**
 * 收藏相关 API
 */

import { get, post, del } from './request'

// 收藏的陪玩师
export interface FavoritePlayer {
  id: number
  playerId: number
  playerName: string
  playerAvatar?: string
  playerRating: number
  isOnline: boolean
  gameNames: string[]
  createdAt: string
}

/**
 * 获取收藏列表
 */
export function getFavorites(params?: { page?: number; page_size?: number }) {
  return get<FavoritePlayer[]>('/users/favorites/players', params)
}

/**
 * 添加收藏
 */
export function addFavorite(playerId: number) {
  return post<void>(`/users/favorites/players/${playerId}`)
}

/**
 * 取消收藏
 */
export function removeFavorite(playerId: number) {
  return del<void>(`/users/favorites/players/${playerId}`)
}

/**
 * 检查是否已收藏
 */
export function checkFavorite(playerId: number) {
  return get<{ isFavorited: boolean }>(`/users/favorites/players/${playerId}/check`)
}

/**
 * 批量取消收藏
 */
export function batchRemoveFavorites(playerIds: number[]) {
  return post<void>('/users/favorites/players/batch-remove', { playerIds })
}

export default {
  getFavorites,
  addFavorite,
  removeFavorite,
  checkFavorite,
  batchRemoveFavorites,
}
