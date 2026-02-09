import type { PlayerInfo } from '@/api/publicPlayer'
import type { PlayerCardData } from '@/types/player'
import type { CachedPlayer } from '@/utils/offlineData'

const DEFAULT_PRICE_CENTS = 2000
const DEFAULT_RATING = 5.0

export function mapPlayerInfoToCard(p: PlayerInfo): PlayerCardData & { id: number } {
  const tags = p.tags || []
  const rank = p.rank || tags[0] || ''
  const mainGame = p.mainGame || ''
  const gameTags = tags.filter(tag => tag && tag !== rank).slice(0, 2)

  const minPriceCents =
    typeof (p as { minPriceCents?: number }).minPriceCents === 'number'
      ? (p as { minPriceCents?: number }).minPriceCents
      : typeof (p as { minPrice?: number }).minPrice === 'number'
        ? (p as { minPrice?: number }).minPrice
        : undefined

  const hourlyRateCents =
    typeof p.hourlyRateCents === 'number'
      ? p.hourlyRateCents
      : typeof (p as { hourlyRate?: number }).hourlyRate === 'number'
        ? (p as { hourlyRate?: number }).hourlyRate
        : typeof minPriceCents === 'number'
          ? minPriceCents
          : DEFAULT_PRICE_CENTS

  const minPrice =
    typeof minPriceCents === 'number'
      ? Math.round(minPriceCents / 100)
      : Math.round(hourlyRateCents / 100)

  // 处理后端返回的字符串状态字段
  const pAny = p as Record<string, unknown>
  const onlineStatus = pAny.onlineStatus ?? pAny.OnlineStatus
  const verificationStatus = pAny.verificationStatus ?? pAny.VerificationStatus
  
  const isOnline = typeof p.isOnline === 'boolean' 
    ? p.isOnline 
    : onlineStatus === 'online'
  
  const isVerified = typeof p.isVerified === 'boolean'
    ? p.isVerified
    : verificationStatus === 'verified'

  // 处理 bio 字段
  const bio = (p.bio || (pAny.bio as string) || '').trim()

  // rating 为 0 时回退为默认值
  const rawRating = p.rating ?? p.ratingAverage ?? 0
  const rating = rawRating > 0 ? rawRating : DEFAULT_RATING

  // orderCount：直接取后端值
  const orderCount = p.orderCount ?? (pAny.ratingCount as number) ?? 0

  return {
    id: p.id,
    nickname: p.nickname,
    avatar: p.avatar || p.avatarUrl || '',
    isOnline,
    isVerified,
    status: isOnline ? 'online' : 'offline',
    rating,
    orderCount,
    hourlyRate: hourlyRateCents,
    minPrice,
    rank,
    mainGame,
    bio,
    games: gameTags.length ? gameTags : (mainGame ? [mainGame] : []),
  }
}

export function mapCachedPlayerToCard(p: CachedPlayer): PlayerCardData & { id: number } {
  const mainGame = p.mainGame || ''
  const hourlyRate = p.hourlyRate ?? DEFAULT_PRICE_CENTS
  return {
    id: p.id,
    nickname: p.nickname,
    avatar: p.avatar || '',
    isOnline: p.isOnline,
    isVerified: false,
    status: p.isOnline ? 'online' : 'offline',
    rating: p.rating ?? DEFAULT_RATING,
    orderCount: p.orderCount ?? 0,
    hourlyRate,
    minPrice: Math.round(hourlyRate / 100),
    rank: p.rank || '',
    mainGame,
    games: mainGame ? [mainGame] : [],
  }
}

export function mapCardToCachedPlayer(card: PlayerCardData & { id: number }): CachedPlayer {
  const firstGame = Array.isArray(card.games) && card.games.length > 0
    ? (typeof card.games[0] === 'string' ? card.games[0] : card.games[0]?.name)
    : ''
  const mainGame = card.mainGame || firstGame || ''
  const hourlyRate = card.hourlyRate ?? ((card.minPrice || DEFAULT_PRICE_CENTS / 100) * 100)

  return {
    id: card.id,
    nickname: card.nickname,
    avatar: card.avatar || '',
    rank: card.rank || '',
    rating: card.rating ?? DEFAULT_RATING,
    hourlyRate,
    isOnline: card.isOnline || false,
    orderCount: card.orderCount || 0,
    mainGame,
  }
}
