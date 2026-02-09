/**
 * 主题切换 Hook
 * 支持日间/夜间主题切换，可根据时间自动切换
 */

import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'gamelink_theme'

// 全局主题状态
const themeMode = ref<ThemeMode>('auto')
const isDark = ref(false)

/**
 * 判断当前是否应该使用夜间模式
 */
function shouldBeDark(): boolean {
  const hour = new Date().getHours()
  // 18:00 - 06:00 使用夜间模式
  return hour >= 18 || hour < 6
}

/**
 * 应用主题到页面
 */
function applyTheme(dark: boolean) {
  isDark.value = dark
  
  // #ifdef H5
  // 设置 HTML 元素的 data-theme 属性
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
    // 同时设置 body 的 class
    document.body.classList.toggle('theme-dark', dark)
  }
  // #endif
  
  // #ifdef MP-WEIXIN
  // 小程序设置状态栏和 TabBar 样式
  try {
    uni.setNavigationBarColor({
      frontColor: dark ? '#ffffff' : '#000000',
      backgroundColor: dark ? '#151618' : '#FFFFFF',
      animation: {
        duration: 300,
        timingFunc: 'easeIn'
      }
    })
  } catch (e) {
    // 静默失败
  }
  
  try {
    uni.setTabBarStyle({
      color: dark ? '#94A3B8' : '#64748B',
      selectedColor: dark ? '#7ACC35' : '#7ACC35',
      backgroundColor: dark ? '#1C1D20' : '#FFFFFF',
      borderStyle: dark ? 'black' : 'black',
    })
  } catch (e) {
    // 静默失败
  }
  // #endif
}

/**
 * 初始化主题
 */
function initTheme() {
  try {
    const stored = uni.getStorageSync(STORAGE_KEY) as ThemeMode
    if (stored) {
      themeMode.value = stored
    }
  } catch (e) {
    console.error('Failed to load theme preference:', e)
  }
  
  updateTheme()
}

/**
 * 根据 themeMode 更新实际主题
 */
function updateTheme() {
  let dark = false
  
  switch (themeMode.value) {
    case 'dark':
      dark = true
      break
    case 'light':
      dark = false
      break
    case 'auto':
    default:
      dark = shouldBeDark()
      break
  }
  
  applyTheme(dark)
}

export function useTheme() {
  // 监听主题模式变化
  watch(themeMode, (mode) => {
    try {
      uni.setStorageSync(STORAGE_KEY, mode)
    } catch (e) {
      console.error('Failed to save theme preference:', e)
    }
    updateTheme()
  })
  
  /**
   * 设置主题模式
   */
  const setThemeMode = (mode: ThemeMode) => {
    themeMode.value = mode
  }
  
  /**
   * 切换日间/夜间
   */
  const toggleTheme = () => {
    if (themeMode.value === 'auto') {
      // 如果是自动模式，切换到当前相反的模式
      themeMode.value = isDark.value ? 'light' : 'dark'
    } else {
      themeMode.value = themeMode.value === 'dark' ? 'light' : 'dark'
    }
  }
  
  /**
   * 循环切换：light -> dark -> auto -> light
   */
  const cycleTheme = () => {
    const modes: ThemeMode[] = ['light', 'dark', 'auto']
    const currentIndex = modes.indexOf(themeMode.value)
    const nextIndex = (currentIndex + 1) % modes.length
    themeMode.value = modes[nextIndex]
  }
  
  return {
    themeMode,
    isDark,
    setThemeMode,
    toggleTheme,
    cycleTheme,
    initTheme,
  }
}

// 导出初始化函数供 App.vue 使用
export { initTheme }
