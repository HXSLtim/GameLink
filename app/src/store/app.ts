/**
 * 应用全局状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

export const useAppStore = defineStore('app', () => {
  // 主题
  const themeMode = ref<ThemeMode>('system')
  const systemTheme = ref<'light' | 'dark'>('light')
  
  // 网络状态
  const isOnline = ref(true)
  const networkType = ref<string>('unknown')
  
  // 未读消息数
  const unreadCount = ref(0)
  
  // 计算当前实际主题
  const currentTheme = computed(() => {
    if (themeMode.value === 'system') {
      return systemTheme.value
    }
    return themeMode.value
  })
  
  const isDark = computed(() => currentTheme.value === 'dark')
  
  /**
   * 初始化
   */
  function init() {
    // 恢复主题设置
    try {
      const savedTheme = uni.getStorageSync('theme_mode') as ThemeMode
      if (savedTheme) {
        themeMode.value = savedTheme
      }
    } catch (e) {
      console.error('Failed to restore theme:', e)
    }
    
    // 监听系统主题变化
    // #ifdef MP-WEIXIN
    const systemInfo = uni.getSystemInfoSync()
    systemTheme.value = systemInfo.theme === 'dark' ? 'dark' : 'light'
    
    uni.onThemeChange((res) => {
      systemTheme.value = res.theme === 'dark' ? 'dark' : 'light'
    })
    // #endif
    
    // 监听网络状态
    uni.getNetworkType({
      success: (res) => {
        networkType.value = res.networkType
        isOnline.value = res.networkType !== 'none'
      }
    })
    
    uni.onNetworkStatusChange((res) => {
      isOnline.value = res.isConnected
      networkType.value = res.networkType
    })
  }
  
  /**
   * 设置主题模式
   */
  function setThemeMode(mode: ThemeMode) {
    themeMode.value = mode
    uni.setStorageSync('theme_mode', mode)
    
    // 应用主题
    applyTheme()
  }
  
  /**
   * 切换主题
   */
  function toggleTheme() {
    const newMode = isDark.value ? 'light' : 'dark'
    setThemeMode(newMode)
  }
  
  /**
   * 应用主题到页面
   */
  function applyTheme() {
    const theme = currentTheme.value
    
    // #ifdef H5
    document.documentElement.setAttribute('data-theme', theme)
    // #endif
    
    // #ifdef MP-WEIXIN
    uni.setNavigationBarColor({
      frontColor: theme === 'dark' ? '#ffffff' : '#000000',
      backgroundColor: theme === 'dark' ? '#1A1A1A' : '#FFFFFF',
    })
    // #endif
  }
  
  /**
   * 设置未读消息数
   */
  function setUnreadCount(count: number) {
    unreadCount.value = count
    
    // 设置 TabBar 角标
    if (count > 0) {
      uni.setTabBarBadge({
        index: 2, // 消息 Tab 索引
        text: count > 99 ? '99+' : String(count),
      })
    } else {
      uni.removeTabBarBadge({ index: 2 })
    }
  }
  
  return {
    // 状态
    themeMode,
    systemTheme,
    isOnline,
    networkType,
    unreadCount,
    
    // 计算属性
    currentTheme,
    isDark,
    
    // 方法
    init,
    setThemeMode,
    toggleTheme,
    applyTheme,
    setUnreadCount,
  }
})
