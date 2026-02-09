export interface GameOption {
  id: number
  name: string
  icon?: string
}

export interface GameTabItem {
  id: number | string
  name: string
  icon?: string
}

export interface GameRankOption {
  id: number | string
  name: string
  icon?: string
  gameId?: number | string
  level?: number
}

export interface GameCardData {
  id: number
  name: string
  coverImage?: string
  categoryId?: string
  isHot?: boolean
  playerCount: number
  minPrice: number
  maxPrice: number
}

export interface HotGameData {
  id: number
  name: string
  icon?: string
  playerCount?: number
}
