/**
 * 分页管理
 */

import { ref, computed } from 'vue'
import type { PaginationOptions } from '@/types/page'

export function usePagination<T>(options: PaginationOptions = {}) {
  const { initialPage = 1, pageSize = 20 } = options
  
  const page = ref(initialPage)
  const size = ref(pageSize)
  const total = ref(0)
  const list = ref<T[]>([]) as { value: T[] }
  const loading = ref(false)
  const loadingMore = ref(false)
  const noMore = ref(false)
  
  // 计算属性
  const totalPages = computed(() => Math.ceil(total.value / size.value))
  const hasMore = computed(() => page.value < totalPages.value)
  const isEmpty = computed(() => !loading.value && list.value.length === 0)
  
  /**
   * 重置分页
   */
  function reset() {
    page.value = initialPage
    total.value = 0
    list.value = []
    noMore.value = false
  }
  
  /**
   * 设置数据（刷新时）
   */
  function setData(data: T[], totalCount?: number) {
    list.value = data
    if (totalCount !== undefined) {
      total.value = totalCount
    }
    noMore.value = data.length < size.value
  }
  
  /**
   * 追加数据（加载更多时）
   */
  function appendData(data: T[], totalCount?: number) {
    list.value = [...list.value, ...data]
    if (totalCount !== undefined) {
      total.value = totalCount
    }
    noMore.value = data.length < size.value
    page.value++
  }
  
  /**
   * 更新列表中的某一项
   */
  function updateItem(predicate: (item: T) => boolean, updater: (item: T) => T) {
    const index = list.value.findIndex(predicate)
    if (index > -1) {
      list.value[index] = updater(list.value[index])
    }
  }
  
  /**
   * 删除列表中的某一项
   */
  function removeItem(predicate: (item: T) => boolean) {
    const index = list.value.findIndex(predicate)
    if (index > -1) {
      list.value.splice(index, 1)
      total.value = Math.max(0, total.value - 1)
    }
  }
  
  /**
   * 在列表开头插入一项
   */
  function prependItem(item: T) {
    list.value.unshift(item)
    total.value++
  }
  
  /**
   * 获取分页参数
   */
  function getParams() {
    return {
      page: page.value,
      page_size: size.value,
    }
  }
  
  return {
    // 状态
    page,
    size,
    total,
    list,
    loading,
    loadingMore,
    noMore,
    
    // 计算属性
    totalPages,
    hasMore,
    isEmpty,
    
    // 方法
    reset,
    setData,
    appendData,
    updateItem,
    removeItem,
    prependItem,
    getParams,
  }
}
