import { ref, onMounted, onUnmounted } from 'vue'

/**
 * 网络状态检测 Hook
 * 用于检测当前网络连接状态
 */
export function useNetwork() {
  const isOnline = ref(true)
  const networkType = ref<string>('unknown')

  // 更新网络状态
  const updateNetworkStatus = () => {
    // #ifdef H5
    isOnline.value = navigator.onLine
    // @ts-ignore
    const connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection
    if (connection) {
      networkType.value = connection.effectiveType || 'unknown'
    }
    // #endif

    // #ifndef H5
    uni.getNetworkType({
      success: (res) => {
        networkType.value = res.networkType
        isOnline.value = res.networkType !== 'none'
      },
    })
    // #endif
  }

  // 监听网络变化
  const setupListeners = () => {
    // #ifdef H5
    window.addEventListener('online', updateNetworkStatus)
    window.addEventListener('offline', updateNetworkStatus)
    // #endif

    // #ifndef H5
    uni.onNetworkStatusChange((res) => {
      isOnline.value = res.isConnected
      networkType.value = res.networkType
    })
    // #endif
  }

  const removeListeners = () => {
    // #ifdef H5
    window.removeEventListener('online', updateNetworkStatus)
    window.removeEventListener('offline', updateNetworkStatus)
    // #endif
  }

  onMounted(() => {
    updateNetworkStatus()
    setupListeners()
  })

  onUnmounted(() => {
    removeListeners()
  })

  return {
    isOnline,
    networkType,
    refresh: updateNetworkStatus,
  }
}
