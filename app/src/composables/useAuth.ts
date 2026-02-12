/**
 * 认证相关组合式函数
 */

import { computed } from 'vue'
import { useUserStore } from '@/store'
import { wechatLogin, refreshToken as refreshTokenApi } from '@/api/auth'

export function useAuth() {
  const userStore = useUserStore()
  
  // 计算属性
  const isLoggedIn = computed(() => userStore.isLoggedIn)
  const isPlayer = computed(() => userStore.isPlayer)
  const isUser = computed(() => userStore.isUser)
  const userInfo = computed(() => userStore.userInfo)
  
  /**
   * 微信登录
   */
  async function loginWithWechat(referralCode?: string) {
    return new Promise((resolve, reject) => {
      // #ifdef MP-WEIXIN
      uni.login({
        success: async (loginRes) => {
          if (!loginRes.code) {
            reject(new Error('微信登录失败：未获取到 code'))
            return
          }
          try {
            const res = await wechatLogin({
              code: loginRes.code,
              referralCode,
            })
            
            if (res.data) {
              const accessToken = res.data.accessToken || res.data.token
              if (!accessToken) {
                reject(new Error('登录失败：缺少 access token'))
                return
              }
              userStore.login({
                accessToken,
                refreshToken: res.data.refreshToken,
                user: res.data.user,
              })
              resolve(res.data)
            }
          } catch (error) {
            reject(error)
          }
        },
        fail: reject,
      })
      // #endif
      
      // #ifndef MP-WEIXIN
      reject(new Error('微信登录仅支持小程序环境'))
      // #endif
    })
  }
  
  /**
   * 刷新 Token
   */
  async function refreshToken() {
    if (!userStore.refreshToken) {
      throw new Error('No refresh token')
    }
    
    try {
      const res = await refreshTokenApi(userStore.refreshToken)
      if (res.data) {
        userStore.setToken(res.data.accessToken, userStore.refreshToken || '')
        return res.data
      }
    } catch (error) {
      userStore.logout()
      throw error
    }
  }
  
  /**
   * 登出
   */
  function logout() {
    userStore.logout()
  }
  
  /**
   * 检查登录状态，未登录则跳转
   */
  function checkLogin(redirectUrl?: string) {
    if (!isLoggedIn.value) {
      const url = redirectUrl 
        ? `/pages/auth/login/index?redirect=${encodeURIComponent(redirectUrl)}`
        : '/pages/auth/login/index'
      uni.navigateTo({ url })
      return false
    }
    return true
  }
  
  /**
   * 检查陪玩师权限
   */
  function checkPlayerRole() {
    if (!isPlayer.value) {
      uni.showToast({ title: '请先成为陪玩师', icon: 'none' })
      return false
    }
    return true
  }
  
  return {
    // 状态
    isLoggedIn,
    isPlayer,
    isUser,
    userInfo,
    
    // 方法
    loginWithWechat,
    refreshToken,
    logout,
    checkLogin,
    checkPlayerRole,
  }
}
