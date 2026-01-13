/// <reference path="./types/index.d.ts" />

interface IAppOption {
  globalData: {
    identity: 'user' | 'player'
    themeMode: 'light' | 'dark'
    themeClass: string
    isLoggedIn: boolean
    statusBarHeight: number
    navBarHeight: number
    systemInfo?: WechatMiniprogram.SystemInfo
  }
  switchIdentity: (identity: 'user' | 'player') => void
  toggleTheme: () => void
  checkUpdate: () => void
}
