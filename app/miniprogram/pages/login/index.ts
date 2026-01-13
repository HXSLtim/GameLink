import { wxLogin } from '../../utils/auth'

Page({
  data: {
    loading: false,
    agreed: false,
    redirect: '',
  },
  
  onLoad(options) {
    const redirect = options.redirect ? decodeURIComponent(options.redirect) : ''
    this.setData({ redirect })
  },
  
  handleAgreeToggle() {
    this.setData({ agreed: !this.data.agreed })
  },
  
  async handleWxLogin() {
    if (!this.data.agreed) {
      wx.showToast({ title: '请先同意用户协议', icon: 'none' })
      return
    }
    
    if (this.data.loading) return
    
    this.setData({ loading: true })
    
    try {
      await wxLogin()
      
      wx.showToast({ title: '登录成功', icon: 'success' })
      
      // 更新全局状态
      const app = getApp<IAppOption>()
      app.globalData.isLoggedIn = true
      
      // 跳转
      setTimeout(() => {
        const { redirect } = this.data
        if (redirect) {
          wx.redirectTo({ url: redirect })
        } else {
          wx.switchTab({ url: '/pages/index/index' })
        }
      }, 1000)
    } catch (err) {
      console.error('登录失败:', err)
      wx.showToast({ title: '登录失败，请重试', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },
  
  handleUserAgreement() {
    wx.navigateTo({ url: '/pages/webview/index?url=https://gamelink.com/agreement' })
  },
  
  handlePrivacyPolicy() {
    wx.navigateTo({ url: '/pages/webview/index?url=https://gamelink.com/privacy' })
  },
})
