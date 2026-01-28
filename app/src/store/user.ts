/**
 * 用户状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface UserInfo {
  id: number
  phone: string
  nickname: string
  avatar: string
  role: 'user' | 'player' | 'admin'
  status: string
  playerId?: number // 如果是陪玩师，关联的陪玩师ID
}

const TOKEN_KEY = 'gamelink_token'
const REFRESH_TOKEN_KEY = 'gamelink_refresh_token'
const USER_KEY = 'gamelink_user'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)
  const isPlayer = computed(() => userInfo.value?.role === 'player')
  const isUser = computed(() => userInfo.value?.role === 'user')
  const userId = computed(() => userInfo.value?.id)
  const playerId = computed(() => userInfo.value?.playerId)
  
  /**
   * 初始化 - 从本地存储恢复状态
   */
  function init() {
    try {
      const storedToken = uni.getStorageSync(TOKEN_KEY)
      const storedRefreshToken = uni.getStorageSync(REFRESH_TOKEN_KEY)
      const storedUser = uni.getStorageSync(USER_KEY)
      
      if (storedToken) {
        token.value = storedToken
      }
      if (storedRefreshToken) {
        refreshToken.value = storedRefreshToken
      }
      if (storedUser) {
        userInfo.value = JSON.parse(storedUser)
      }
    } catch (e) {
      console.error('Failed to restore user state:', e)
    }
  }
  
  /**
   * 设置 Token
   */
  function setToken(accessToken: string, refresh?: string) {
    token.value = accessToken
    uni.setStorageSync(TOKEN_KEY, accessToken)
    
    if (refresh) {
      refreshToken.value = refresh
      uni.setStorageSync(REFRESH_TOKEN_KEY, refresh)
    }
  }
  
  /**
   * 设置用户信息
   */
  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    uni.setStorageSync(USER_KEY, JSON.stringify(info))
  }
  
  /**
   * 登录成功后调用
   */
  function login(data: { accessToken: string; refreshToken?: string; user: UserInfo }) {
    setToken(data.accessToken, data.refreshToken)
    setUserInfo(data.user)
  }
  
  /**
   * 退出登录
   */
  function logout() {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    
    uni.removeStorageSync(TOKEN_KEY)
    uni.removeStorageSync(REFRESH_TOKEN_KEY)
    uni.removeStorageSync(USER_KEY)
    
    // 跳转到登录页
    uni.reLaunch({ url: '/pages/auth/login/index' })
  }
  
  /**
   * 更新用户信息（部分更新）
   */
  function updateUserInfo(partial: Partial<UserInfo>) {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...partial }
      uni.setStorageSync(USER_KEY, JSON.stringify(userInfo.value))
    }
  }
  
  /**
   * 切换角色（用户/陪玩师）
   */
  function switchRole(role: 'user' | 'player') {
    if (userInfo.value) {
      userInfo.value.role = role
      uni.setStorageSync(USER_KEY, JSON.stringify(userInfo.value))
    }
  }
  
  return {
    // 状态
    token,
    refreshToken,
    userInfo,
    
    // 计算属性
    isLoggedIn,
    isPlayer,
    isUser,
    userId,
    playerId,
    
    // 方法
    init,
    setToken,
    setUserInfo,
    login,
    logout,
    updateUserInfo,
    switchRole,
  }
})
