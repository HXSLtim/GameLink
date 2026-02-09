import { ref } from 'vue'
import type { ToastOptions, ToastState } from '@/types/ui'

const toastState = ref<ToastState>({
  visible: false,
  message: '',
  type: 'info',
  icon: undefined,
})

let hideTimer: number | null = null

/**
 * 游戏风格 Toast 提示
 */
export function useToast() {
  const show = (options: ToastOptions | string) => {
    const opts: ToastOptions = typeof options === 'string' ? { message: options } : options
    
    // 清除之前的定时器
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }
    
    toastState.value = {
      visible: true,
      message: opts.message,
      type: opts.type || 'info',
      icon: opts.icon,
    }
    
    // 自动隐藏（loading 不自动隐藏）
    if (opts.type !== 'loading') {
      const duration = opts.duration ?? 2000
      hideTimer = setTimeout(() => {
        hide()
      }, duration) as unknown as number
    }
  }
  
  const hide = () => {
    toastState.value.visible = false
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }
  }
  
  // 快捷方法
  const success = (message: string, duration?: number) => {
    show({ message, type: 'success', duration })
  }
  
  const error = (message: string, duration?: number) => {
    show({ message, type: 'error', duration })
  }
  
  const warning = (message: string, duration?: number) => {
    show({ message, type: 'warning', duration })
  }
  
  const info = (message: string, duration?: number) => {
    show({ message, type: 'info', duration })
  }
  
  const loading = (message = '加载中...') => {
    show({ message, type: 'loading' })
  }
  
  return {
    state: toastState,
    show,
    hide,
    success,
    error,
    warning,
    info,
    loading,
  }
}

// 全局单例
let globalToast: ReturnType<typeof useToast> | null = null

export function getGlobalToast() {
  if (!globalToast) {
    globalToast = useToast()
  }
  return globalToast
}
