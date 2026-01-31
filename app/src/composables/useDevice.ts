import { ref, onMounted, onUnmounted } from 'vue'

/**
 * 设备检测 Hook
 * 用于判断当前设备类型（PC/平板/手机）
 */
export function useDevice() {
  const isPC = ref(false)
  const isTablet = ref(false)
  const isMobile = ref(true)
  const screenWidth = ref(0)
  
  const PC_BREAKPOINT = 1024
  const TABLET_BREAKPOINT = 768
  
  const checkDevice = () => {
    // #ifdef H5
    screenWidth.value = window.innerWidth
    isPC.value = screenWidth.value >= PC_BREAKPOINT
    isTablet.value = screenWidth.value >= TABLET_BREAKPOINT && screenWidth.value < PC_BREAKPOINT
    isMobile.value = screenWidth.value < TABLET_BREAKPOINT
    // #endif
    
    // #ifndef H5
    // 小程序端默认为移动端
    const systemInfo = uni.getSystemInfoSync()
    screenWidth.value = systemInfo.windowWidth
    isPC.value = false
    isTablet.value = false
    isMobile.value = true
    // #endif
  }
  
  const setupListeners = () => {
    // #ifdef H5
    window.addEventListener('resize', checkDevice)
    // #endif
  }
  
  const removeListeners = () => {
    // #ifdef H5
    window.removeEventListener('resize', checkDevice)
    // #endif
  }
  
  onMounted(() => {
    checkDevice()
    setupListeners()
  })
  
  onUnmounted(() => {
    removeListeners()
  })
  
  return {
    isPC,
    isTablet,
    isMobile,
    screenWidth,
    refresh: checkDevice,
  }
}
