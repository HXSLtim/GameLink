/**
 * 加载状态管理
 */

import { ref } from 'vue'

export function useLoading(initialValue = false) {
  const loading = ref(initialValue)
  const loadingMore = ref(false)
  const refreshing = ref(false)
  const noMore = ref(false)
  
  /**
   * 开始加载
   */
  function startLoading() {
    loading.value = true
  }
  
  /**
   * 结束加载
   */
  function stopLoading() {
    loading.value = false
  }
  
  /**
   * 开始加载更多
   */
  function startLoadingMore() {
    loadingMore.value = true
  }
  
  /**
   * 结束加载更多
   */
  function stopLoadingMore() {
    loadingMore.value = false
  }
  
  /**
   * 开始刷新
   */
  function startRefreshing() {
    refreshing.value = true
  }
  
  /**
   * 结束刷新
   */
  function stopRefreshing() {
    refreshing.value = false
  }
  
  /**
   * 设置没有更多数据
   */
  function setNoMore(value = true) {
    noMore.value = value
  }
  
  /**
   * 重置状态
   */
  function reset() {
    loading.value = false
    loadingMore.value = false
    refreshing.value = false
    noMore.value = false
  }
  
  /**
   * 包装异步函数，自动管理加载状态
   */
  async function withLoading<T>(fn: () => Promise<T>): Promise<T> {
    startLoading()
    try {
      return await fn()
    } finally {
      stopLoading()
    }
  }
  
  return {
    loading,
    loadingMore,
    refreshing,
    noMore,
    startLoading,
    stopLoading,
    startLoadingMore,
    stopLoadingMore,
    startRefreshing,
    stopRefreshing,
    setNoMore,
    reset,
    withLoading,
  }
}
