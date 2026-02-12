/**
 * 陪玩师相关 API
 */

import { get, post, put, del, type RequestConfig } from './request'

// ============ 陪玩师服务 ============

export interface PlayerService {
  id: number
  gameId: number
  gameName: string
  gameIcon?: string
  serviceName: string
  serviceType: string
  rankId?: number
  rankName?: string
  price: number
  unit: string
  description?: string
  isOnline: boolean
  orderCount: number
  createdAt: string
  updatedAt: string
}

export interface CreateServiceData {
  gameId: number
  serviceType?: string
  serviceName?: string
  rankId?: number
  rankName?: string
  price?: number
  priceCents?: number
  unit: string
  description?: string
  isOnline?: boolean
}

/**
 * 获取陪玩师服务列表
 */
export function getPlayerServices() {
  return get<PlayerService[]>('/player/services')
}

/**
 * 创建服务
 */
export function createPlayerService(data: CreateServiceData) {
  return post<PlayerService>('/player/services', data)
}

/**
 * 更新服务
 */
export function updatePlayerService(id: number, data: Partial<CreateServiceData>) {
  return put<PlayerService>(`/player/services/${id}`, data)
}

/**
 * 删除服务
 */
export function deletePlayerService(id: number) {
  return del<void>(`/player/services/${id}`)
}

/**
 * 切换服务上下架状态
 */
export function toggleServiceStatus(id: number, isOnline: boolean) {
  return put<PlayerService>(`/player/services/${id}/status`, { isOnline })
}

// ============ 排班设置 ============

export interface ScheduleSlot {
  dayOfWeek: number // 0-6, 0 = Sunday
  startTime: string // HH:mm
  endTime: string   // HH:mm
  isAvailable: boolean
}

export interface PlayerSchedule {
  slots: ScheduleSlot[]
  timezone: string
  autoOffline: boolean
  updatedAt: string
}

/**
 * 获取排班设置
 */
export function getPlayerSchedule() {
  return get<PlayerSchedule>('/player/schedule')
}

/**
 * 更新排班设置
 */
export function updatePlayerSchedule(data: Partial<PlayerSchedule>) {
  return put<PlayerSchedule>('/player/schedule', data)
}

// ============ 数据统计 ============

export interface TodayStats {
  orderCount: number
  earningsCents: number
  serviceDurationMinutes: number
  serviceDuration?: number
  averageRating: number
  newFollowers: number
}

export interface OverviewStats {
  totalOrders: number
  totalEarningsCents: number
  totalServiceHours: number
  averageRating: number
  followerCount: number
  favoriteCount: number
  completionRate: number
  responseRate: number
  monthlyOrders: number
  monthlyEarningsCents: number
  weeklyTrend: Array<{
    date: string
    orders: number
    earnings: number
  }>
}

/**
 * 获取今日统计
 */
export function getTodayStats() {
  return get<TodayStats>('/player/stats/today')
}

/**
 * 获取总览统计
 */
export function getOverviewStats() {
  return get<OverviewStats>('/player/stats/overview')
}

// ============ 公开评价 ============

export interface PlayerReview {
  id: number
  userId: number
  userName: string
  userAvatar?: string
  rating: number
  content?: string
  tags?: string[]
  images?: string[]
  orderId: number
  gameName: string
  serviceName: string
  reply?: string
  repliedAt?: string
  createdAt: string
}

export interface ReviewListParams {
  page?: number
  page_size?: number
  rating?: number
}

/**
 * 获取陪玩师公开评价列表
 */
export function getPlayerReviews(
  playerId: number,
  params?: ReviewListParams,
  config?: Partial<RequestConfig>
) {
  return get<PlayerReview[]>(`/public/players/${playerId}/reviews`, params, config)
}

/**
 * 获取收益概览
 */
export function getPlayerEarnings(params?: any, config?: Partial<RequestConfig>) {
  return get<any>('/player/earnings/summary', params, config)
}

/**
 * 获取收益趋势
 */
export function getEarningsStats(params?: { period: string }, config?: Partial<RequestConfig>) {
  return get<any>('/player/earnings/trend', params, config)
}

/**
 * 获取实名认证状态
 */
export function getCertificationStatus() {
  return get<any>('/player/certification/identity')
}

/**
 * 提交实名认证
 */
export function submitCertification(data: any) {
  return post<void>('/player/certification/identity', data)
}

// ============ 陪玩师申请与状态 ============

/**
 * 申请成为陪玩师
 */
export function applyPlayer(data: any) {
  return post<void>('/player/apply', data)
}

/**
 * 更新在线状态
 */
export function updateOnlineStatus(status: 'online' | 'offline' | 'busy') {
  return put<void>('/player/online-status', { status })
}

/**
 * 获取在线状态
 */
export function getOnlineStatus() {
  return get<{ status: string }>('/player/online-status')
}

/**
 * 获取可接订单列表
 */
export function getAvailableOrders(params?: { page?: number; page_size?: number }, config?: Partial<RequestConfig>) {
  return get<any[]>('/player/orders/available', params, config)
}

export default {
  // 服务管理
  getPlayerServices,
  createPlayerService,
  updatePlayerService,
  deletePlayerService,
  toggleServiceStatus,
  // 排班
  getPlayerSchedule,
  updatePlayerSchedule,
  // 统计
  getTodayStats,
  getOverviewStats,
  getPlayerEarnings,
  getEarningsStats,
  // 评价
  getPlayerReviews,
  // 认证
  getCertificationStatus,
  submitCertification,
  // 申请与状态
  applyPlayer,
  updateOnlineStatus,
  getOnlineStatus,
  getAvailableOrders,
}
