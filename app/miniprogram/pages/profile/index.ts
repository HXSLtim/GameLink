import { getStorage, setStorage, StorageKeys } from '../../utils/storage'

Page({
  data: {
    statusBarHeight: 44,
    role: 'user' as 'user' | 'player',
    themeMode: 'dark', // 默认为暗色
  },

  onLoad() {
    const systemInfo = wx.getSystemInfoSync()
    const role = getStorage<'user' | 'player'>(StorageKeys.USER_MODE) || 'user'
    const app = getApp<IAppOption>()

    this.setData({
      statusBarHeight: systemInfo.statusBarHeight || 44,
      role,
      themeMode: app.globalData.themeMode || 'dark',
    })
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 3 })
    }
    // Sync theme in case it changed elsewhere
    const app = getApp<IAppOption>()
    if (this.data.themeMode !== app.globalData.themeMode) {
      this.setData({ themeMode: app.globalData.themeMode })
    }
  },

  onEditProfile() {
    wx.navigateTo({ url: '/pages/profile/edit' })
  },

  onToggleTheme() {
    const app = getApp<IAppOption>()
    app.toggleTheme()
    this.setData({ themeMode: app.globalData.themeMode })
  },

  onSwitchRole() {
    const newRole = this.data.role === 'user' ? 'player' : 'user'
    this.setData({ role: newRole })
    setStorage(StorageKeys.USER_MODE, newRole)

    // Sync to global app state
    const app = getApp<IAppOption>()
    app.switchIdentity(newRole)

    wx.showToast({
      title: `Switched to ${newRole === 'player' ? 'Player' : 'User'} mode`,
      icon: 'none',
    })
  },

  onMenuTap(e: WechatMiniprogram.CustomEvent) {
    const { type } = e.detail
    const routes: Record<string, string> = {
      wallet: '/pages/wallet/index',
      orders: '/pages/order/list',
      favorites: '/pages/favorites/index',
    }
    if (routes[type]) {
      wx.navigateTo({ url: routes[type] })
    }
  },
})
