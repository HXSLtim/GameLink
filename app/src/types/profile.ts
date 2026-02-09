import type { ProfileGenderValue } from '@/types/common'

export interface ProfileForm {
  avatar: string
  nickname: string
  gender: ProfileGenderValue
  birthday: string
  region: string
  bio: string
  games: string[]
}

export interface ProfileContactInfo {
  phone: string
  wechatBound: boolean
}

export type ProfileStatKey = 'orders' | 'favorites' | 'wallet'
