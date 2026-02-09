/**
 * 陪玩师公开 API（用户端访问）
 */

import { get, post, type RequestConfig } from './request'
import type { Gender, ProfileGender } from '@/types/common'

// 陪玩师基础信息
export interface PlayerInfo {
  id: number
  userId: number
  nickname: string
  avatar?: string
  avatarUrl?: string
  gender: ProfileGender
  age?: number
  bio?: string
  voiceSample?: string
  rating?: number
  ratingAverage?: number
  ratingCount?: number
  orderCount?: number
  followerCount: number
  isOnline: boolean
  isVerified: boolean
  rank?: string
  mainGame?: string
  hourlyRateCents?: number
  tags?: string[]
}

// 陪玩师服务项
export interface PlayerServiceItem {
  id: number
  gameId: number
  gameName: string
  gameIcon?: string
  serviceType: string
  rankId?: number
  rankName?: string
  priceCents: number
  unit: string
  description?: string
}

// 陪玩师详情
export interface PlayerDetail extends PlayerInfo {
  services: PlayerServiceItem[]
  gameRanks: Array<{
    gameId: number
    gameName: string
    rankName: string
    rankIcon?: string
  }>
  photos?: string[]
  createdAt: string
}

// 列表查询参数
export interface PlayerListParams {
  gameId?: number
  gender?: Gender
  minPrice?: number
  maxPrice?: number
  isOnline?: boolean
  keyword?: string
  sortBy?: 'rating' | 'orders' | 'price_asc' | 'price_desc' | 'newest'
  page?: number
  page_size?: number
}

/**
 * 获取陪玩师列表
 */
export function getPlayerList(params?: PlayerListParams, config?: Partial<RequestConfig>) {
  return get<PlayerInfo[]>('/public/players', params, config)
}

/**
 * 获取陪玩师详情
 */
export function getPlayerDetail(playerId: number, config?: Partial<RequestConfig>) {
  return get<PlayerDetail>(`/public/players/${playerId}`, undefined, config)
}

/**
 * 获取陪玩师服务列表
 */
export function getPlayerServiceItems(playerId: number, config?: Partial<RequestConfig>) {
  return get<PlayerServiceItem[]>(`/public/players/${playerId}/services`, undefined, config)
}

/**
 * 获取推荐/精选陪玩师
 */
export function getRecommendedPlayers(limit = 10, config?: Partial<RequestConfig>) {
  return get<PlayerInfo[]>('/public/players/featured', { limit }, config)
}

/**
 * 获取热门陪玩师
 */
export function getHotPlayers(limit = 10, config?: Partial<RequestConfig>) {
  return get<PlayerInfo[]>('/public/players/hot', { limit }, config)
}

/**
 * 搜索陪玩师
 */
export function searchPlayers(
  keyword: string,
  params?: Omit<PlayerListParams, 'keyword'>,
  config?: Partial<RequestConfig>
) {
  return get<PlayerInfo[]>('/public/search/players', { keyword, ...params }, config)
}

/**
 * 获取陪玩师评价列表
 */
export function getPlayerReviews(playerId: number, params?: { page_size?: number }, config?: Partial<RequestConfig>) {
  return get<any[]>(`/public/players/${playerId}/reviews`, params, config)
}

export default {
  getPlayerList,
  getPlayerDetail,
  getPlayerServiceItems,
  getRecommendedPlayers,
  getHotPlayers,
  searchPlayers,
  getPlayerReviews,
}
