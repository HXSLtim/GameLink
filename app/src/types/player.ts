import type { PlayerReviewData } from '@/types/review'
import type { CertStatus, DashboardStatus } from '@/types/status'
import type { Gender } from '@/types/common'

export interface PlayerGame {
  id?: number | string
  name: string
}

export interface PlayerGameTag {
  key: string
  name: string
}

export interface PlayerSummary {
  id: number
  nickname: string
  avatar?: string
  rating?: number
  orderCount?: number
  isOnline?: boolean
}

export interface PlayerOrderInfo extends PlayerSummary {
  avatar: string
  isOnline: boolean
  rating: number
}

export interface PlayerDashboardInfo extends PlayerSummary {
  rating: number
  certificationStatus: CertStatus
}

export interface PlayerTodayStats {
  orders: number
  earnings: number
  duration: number
  rating: string
}

export interface PlayerHeaderData {
  nickname: string
  avatar?: string
  coverImage?: string
  signature?: string
  gender?: Gender
  isOnline: boolean
  isVerified: boolean
  rating: number
  orderCount: number
  favoriteCount: number
  createdAt: string
}

export interface PlayerDetailData extends PlayerHeaderData {
  id: number
  games: PlayerGameData[]
  services: PlayerServiceData[]
  reviews: PlayerReviewData[]
}

export interface PlayerGameData {
  id: number
  name: string
  icon?: string
  rankName: string
  price: number
}

export interface PlayerServiceData {
  id: number
  name: string
  description?: string
  price: number
  unit?: string
}

export interface PlayerServiceCardData {
  id: number
  gameId: number
  gameName: string
  gameIcon?: string
  serviceName: string
  price: number
  unit: string
  rankName: string
  description?: string
  isOnline: boolean
}

export interface PlayerServiceForm {
  gameId?: number
  gameName: string
  serviceType: string
  serviceName: string
  rankId?: number
  rankName: string
  price: number
  unit: string
  description: string
}

export interface PlayerCardData {
  id?: number
  playerId?: number
  nickname: string
  avatar?: string
  bio?: string
  isOnline?: boolean
  isVerified?: boolean
  status?: DashboardStatus
  rating?: number
  orderCount?: number
  minPrice?: number
  hourlyRate?: number
  games?: Array<PlayerGame | string>
  rank?: string
  mainGame?: string
}

export type RecommendPlayerData = PlayerCardData & { id: number }

export interface FavoritePlayerData extends PlayerCardData {
  id: number
  playerId: number
}
