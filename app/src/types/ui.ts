export interface QuickActionItem {
  key: string
  icon: string
  label: string
  badge?: number
}

export interface StatItem {
  value: string | number
  label: string
  unit?: string
  highlight?: boolean
  onClick?: () => void
}

export interface HeaderStatItem {
  key?: string
  label: string
  value: string | number
  clickable?: boolean
}

export interface MenuItem {
  key: string
  label: string
  icon: string
  iconColor?: string
  iconBg?: string
  badge?: number | string
  value?: string
  disabled?: boolean
}

export interface SettingsItem {
  key: string
  label: string
  icon?: string
  iconColor?: string
  value?: string
  type?: 'link' | 'switch'
  checked?: boolean
}

export interface SettingsState {
  pushEnabled: boolean
  messageEnabled: boolean
  orderEnabled: boolean
  promotionEnabled: boolean
  showOnlineStatus: boolean
  allowStrangerMessage: boolean
}

export type ToastType = 'success' | 'error' | 'warning' | 'info' | 'loading'

export interface ToastOptions {
  message: string
  type?: ToastType
  duration?: number
  icon?: string
}

export interface ToastState {
  visible: boolean
  message: string
  type: ToastType
  icon?: string
}

export interface ConfirmOptions {
  title?: string
  content?: string
  confirmText?: string
  cancelText?: string
}

export interface TabItem {
  key: string
  label: string
  badge?: number | string
  disabled?: boolean
}

export interface TabBarItem {
  path: string
  text: string
  icon: string
  iconNormal?: string
  iconActive?: string
  badge?: number
  dot?: boolean
}
