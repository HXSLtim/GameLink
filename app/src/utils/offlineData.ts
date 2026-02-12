/**
 * 离线兜底数据
 * 当接口失败时，提供默认推荐数据作为占位
 */

// 热门游戏默认数据
export const defaultHotGames = [
  { id: 1, name: '王者荣耀', icon: '', category: 'MOBA', playerCount: 5000 },
  { id: 2, name: '英雄联盟', icon: '', category: 'MOBA', playerCount: 3500 },
  { id: 3, name: '和平精英', icon: '', category: 'FPS', playerCount: 2800 },
  { id: 4, name: '原神', icon: '', category: 'RPG', playerCount: 2200 },
  { id: 5, name: '永劫无间', icon: '', category: '动作', playerCount: 1800 },
  { id: 6, name: 'CSGO', icon: '', category: 'FPS', playerCount: 1500 },
]

// 推荐陪玩师默认数据
export const defaultPlayers = [
  {
    id: 1,
    nickname: '甜甜',
    avatar: '',
    rank: '王者',
    rating: 4.9,
    hourlyRate: 3000,
    isOnline: true,
    orderCount: 520,
    mainGame: '王者荣耀',
  },
  {
    id: 2,
    nickname: '小鱼',
    avatar: '',
    rank: '钻石',
    rating: 4.8,
    hourlyRate: 2500,
    isOnline: true,
    orderCount: 380,
    mainGame: '英雄联盟',
  },
  {
    id: 3,
    nickname: '阿星',
    avatar: '',
    rank: '铂金',
    rating: 4.7,
    hourlyRate: 2000,
    isOnline: false,
    orderCount: 210,
    mainGame: '和平精英',
  },
  {
    id: 4,
    nickname: '小月',
    avatar: '',
    rank: '王者',
    rating: 4.9,
    hourlyRate: 3500,
    isOnline: true,
    orderCount: 680,
    mainGame: '王者荣耀',
  },
]

// 公共频道默认数据
export const defaultChannels = [
  {
    id: 1,
    name: '王者荣耀开黑',
    description: '王者玩家聚集地，一起上分！',
    avatarUrl: '',
    currentMembers: 128,
    maxMembers: 500,
    isActive: true,
    gameId: 1,
    isJoined: false,
  },
  {
    id: 2,
    name: 'LOL 峡谷之巅',
    description: '英雄联盟玩家交流频道',
    avatarUrl: '',
    currentMembers: 86,
    maxMembers: 500,
    isActive: true,
    gameId: 2,
    isJoined: false,
  },
  {
    id: 3,
    name: '吃鸡小分队',
    description: '和平精英组队开黑',
    avatarUrl: '',
    currentMembers: 64,
    maxMembers: 300,
    isActive: true,
    gameId: 3,
    isJoined: false,
  },
  {
    id: 4,
    name: '原神交流群',
    description: '原神玩家交流讨论',
    avatarUrl: '',
    currentMembers: 156,
    maxMembers: 500,
    isActive: true,
    gameId: 4,
    isJoined: false,
  },
]

// 缓存键
const CACHE_KEYS = {
  HOT_GAMES: 'cache_hot_games',
  PLAYERS: 'cache_players',
  CHANNELS: 'cache_channels',
}

// 缓存有效期（毫秒）
const CACHE_TTL = 24 * 60 * 60 * 1000 // 24 小时

interface CacheItem<T> {
  data: T
  timestamp: number
}

// 通用缓存读取
function getCache<T>(key: string): T | null {
  try {
    const raw = uni.getStorageSync(key)
    if (!raw) return null
    
    const cache: CacheItem<T> = JSON.parse(raw)
    const now = Date.now()
    
    // 检查是否过期
    if (now - cache.timestamp > CACHE_TTL) {
      uni.removeStorageSync(key)
      return null
    }
    
    return cache.data
  } catch {
    return null
  }
}

// 通用缓存写入
function setCache<T>(key: string, data: T): void {
  try {
    const cache: CacheItem<T> = {
      data,
      timestamp: Date.now(),
    }
    uni.setStorageSync(key, JSON.stringify(cache))
  } catch {
    // 忽略缓存失败
  }
}

// 类型定义
export interface CachedGame {
  id: number
  name: string
  icon?: string
  category?: string
  playerCount?: number
}

export interface CachedPlayer {
  id: number
  nickname: string
  avatar?: string
  rank?: string
  rating: number
  hourlyRate: number
  isOnline: boolean
  orderCount: number
  mainGame?: string
}

export interface CachedChannel {
  id: number
  name: string
  description?: string
  avatarUrl?: string
  avatar?: string
  currentMembers?: number
  memberCount?: number
  maxMembers?: number
  isActive?: boolean
  gameId?: number
  gameName?: string
  isJoined?: boolean
}

/**
 * 获取热门游戏（优先缓存，降级默认）
 */
export function getCachedHotGames(): CachedGame[] {
  return getCache<CachedGame[]>(CACHE_KEYS.HOT_GAMES) || defaultHotGames
}

export function saveCachedHotGames(games: CachedGame[]) {
  setCache(CACHE_KEYS.HOT_GAMES, games)
}

/**
 * 获取推荐陪玩师（优先缓存，降级默认）
 */
export function getCachedPlayers(): CachedPlayer[] {
  return getCache<CachedPlayer[]>(CACHE_KEYS.PLAYERS) || defaultPlayers
}

export function saveCachedPlayers(players: CachedPlayer[]) {
  setCache(CACHE_KEYS.PLAYERS, players)
}

/**
 * 获取公共频道（优先缓存，降级默认）
 */
export function getCachedChannels(): CachedChannel[] {
  return getCache<CachedChannel[]>(CACHE_KEYS.CHANNELS) || defaultChannels
}

export function saveCachedChannels(channels: CachedChannel[]) {
  setCache(CACHE_KEYS.CHANNELS, channels)
}
