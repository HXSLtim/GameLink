export interface ReviewBase {
  id: number
  rating: number
  content?: string
  tags?: string[]
  images?: string[]
  reply?: string
  createdAt: string
}

export interface ReviewPlayerInfo {
  id: number
  nickname: string
  avatar?: string
}

export interface ReviewCardData extends ReviewBase {
  orderId: number
  player: ReviewPlayerInfo
  gameName: string
  serviceName: string
}

export interface PlayerReviewData extends ReviewBase {
  userId: number
  userName: string
  userAvatar?: string
}

export type OrderReviewData = ReviewBase
