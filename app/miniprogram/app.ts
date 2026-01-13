// app.ts
import { getIdentity, getThemeClass } from './utils/theme'
import { isLoggedIn } from './utils/auth'

App<IAppOption>({
  globalData: {
    identity: 'user',
    themeMode: 'dark', // 'light' | 'dark'
    themeClass: 'theme-user',
    isLoggedIn: false,
    statusBarHeight: 44,
    navBarHeight: 44,
  },

  onLaunch() {
    // 初始化身份和主题
    const identity = getIdentity()
    const themeClass = getThemeClass()

    this.globalData.identity = identity
    this.globalData.themeClass = themeClass
    this.globalData.isLoggedIn = isLoggedIn()

    // 获取系统信息
    const systemInfo = wx.getSystemInfoSync()
    this.globalData.systemInfo = systemInfo
    this.globalData.statusBarHeight = systemInfo.statusBarHeight || 44

    // 初始化主题
    const themeMode = wx.getStorageSync('theme_mode') || 'dark'
    this.globalData.themeMode = themeMode

    // 检查更新

    this.checkUpdate()
  },

  // 检查小程序更新
  checkUpdate() {
    if (!wx.canIUse('getUpdateManager')) return

    const updateManager = wx.getUpdateManager()

    updateManager.onCheckForUpdate((res) => {
      console.log('检查更新:', res.hasUpdate)
    })

    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已准备好，是否重启应用？',
        success: (res) => {
          if (res.confirm) {
            updateManager.applyUpdate()
          }
        },
      })
    })

    updateManager.onUpdateFailed(() => {
      console.log('更新失败')
    })
  },

  // 切换身份
  switchIdentity(identity: 'user' | 'player') {
    this.globalData.identity = identity
    this.globalData.themeClass = identity === 'player' ? 'theme-player' : 'theme-user'
  },

  // 切换主题
  toggleTheme() {
    const current = this.globalData.themeMode
    const next = current === 'dark' ? 'light' : 'dark'
    this.globalData.themeMode = next
    wx.setStorageSync('theme_mode', next)

    // update status bar style
    if (next === 'light') {
      wx.setNavigationBarColor({
        frontColor: '#000000',
        backgroundColor: '#ffffff',
      })
    } else {
      wx.setNavigationBarColor({
        frontColor: '#ffffff',
        backgroundColor: '#313338',
      })
    }
  },
})
