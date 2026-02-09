export interface GameCertData {
  gameId?: number
  gameName: string
  rankId?: number
  rankName: string
  screenshot?: string
}

import type { Gender } from '@/types/common'

export interface CertificationForm {
  realName: string
  idNumber: string
  gender: Gender | ''
  idCardFront: string
  idCardBack: string
  games: GameCertData[]
  introduction: string
  voiceSample?: string
  voiceDuration?: number
}
