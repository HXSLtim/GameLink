import { ref, computed } from 'vue'
import type { PageStateType, PageStateOptions } from '@/types/page'

/**
 * 页面状态管理 Hook
 * 统一处理 loading/error/empty/content 状态
 */
export function usePageState(options: PageStateOptions = {}) {
  const { initialState = 'loading', requireAuth = false } = options
  
  const state = ref<PageStateType>(initialState)
  const errorMessage = ref('')
  const hasData = ref(false)

  // 状态判断
  const isLoading = computed(() => state.value === 'loading')
  const isError = computed(() => state.value === 'error')
  const isEmpty = computed(() => state.value === 'empty')
  const isContent = computed(() => state.value === 'content')

  // 设置加载中
  const setLoading = () => {
    state.value = 'loading'
    errorMessage.value = ''
  }

  // 设置错误
  const setError = (message?: string) => {
    state.value = 'error'
    errorMessage.value = message || '加载失败，请稍后重试'
  }

  // 设置空态
  const setEmpty = () => {
    state.value = 'empty'
    hasData.value = false
  }

  // 设置内容态
  const setContent = () => {
    state.value = 'content'
    hasData.value = true
  }

  // 设置需要登录
  const setNeedLogin = () => {
    state.value = 'login'
  }

  // 设置离线
  const setOffline = () => {
    state.value = 'offline'
  }

  // 根据数据自动判断状态
  const setStateByData = <T>(data: T[] | T | null | undefined) => {
    if (data === null || data === undefined) {
      setEmpty()
      return
    }
    
    if (Array.isArray(data)) {
      if (data.length === 0) {
        setEmpty()
      } else {
        setContent()
      }
    } else {
      setContent()
    }
  }

  // 封装异步请求
  const handleRequest = async <T>(
    request: () => Promise<T>,
    options?: {
      onSuccess?: (data: T) => void
      onError?: (error: any) => void
      checkEmpty?: (data: T) => boolean
    }
  ): Promise<T | null> => {
    try {
      setLoading()
      const data = await request()
      
      // 检查是否为空
      if (options?.checkEmpty) {
        if (options.checkEmpty(data)) {
          setEmpty()
        } else {
          setContent()
        }
      } else {
        // 默认检查逻辑
        if (Array.isArray(data)) {
          data.length === 0 ? setEmpty() : setContent()
        } else if (data === null || data === undefined) {
          setEmpty()
        } else {
          setContent()
        }
      }
      
      options?.onSuccess?.(data)
      return data
    } catch (error: any) {
      console.error('Request failed:', error)
      
      // 判断错误类型
      if (error?.message?.includes('timeout') || error?.message?.includes('超时')) {
        setError('连接超时，请检查网络后重试')
      } else if (error?.message?.includes('network') || error?.message?.includes('网络')) {
        setOffline()
      } else if (error?.statusCode === 401) {
        setNeedLogin()
      } else {
        setError(error?.message || '请求失败')
      }
      
      options?.onError?.(error)
      return null
    }
  }

  return {
    state,
    errorMessage,
    hasData,
    isLoading,
    isError,
    isEmpty,
    isContent,
    setLoading,
    setError,
    setEmpty,
    setContent,
    setNeedLogin,
    setOffline,
    setStateByData,
    handleRequest,
  }
}
