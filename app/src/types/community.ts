export interface ChannelData {
  id: number
  name: string
  description?: string
  avatar?: string
  memberCount: number
  maxMembers?: number
  isActive: boolean
  isJoined: boolean
  gameId?: number
  gameName?: string
}
