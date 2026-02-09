export interface LoginFormData {
  account: string  // 手机号或邮箱
  password: string
}

export interface RegisterFormData {
  phone: string
  nickname: string
  password: string
  confirmPassword: string
}

import type { AppUserRole } from '@/types/user'

export interface RoleOption {
  value: AppUserRole
  icon: string
  name: string
  desc: string
}
