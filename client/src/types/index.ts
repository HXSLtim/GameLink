// 用户相关
export interface User {
  id: string
  nickname: string
  avatar?: string
  phone?: string
  role: 'user' | 'player'
}

// 陪玩师相关
export interface Player {
  id: string
  nickname: string
  avatar?: string
  games: string[]
  rating: number
  price: number
  status: 'online' | 'busy' | 'offline'
  bio?: string
  orderCount?: number
}

// 订单相关
export type OrderStatus = 'pending' | 'accepted' | 'in_progress' | 'completed' | 'canceled'

export interface Order {
  id: string
  userId: string
  playerId: string
  playerName: string
  game: string
  hours: number
  amount: number
  status: OrderStatus
  createdAt: string
  updatedAt: string
}

// API 响应
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
