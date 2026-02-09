/**
 * 列表页通用 Hook
 * 封装分页加载、刷新、错误处理等逻辑
 */
import { ref, type Ref } from 'vue'
import type { PageStateType, UseListPageOptions, UseListPageReturn } from '@/types/page'
export type { UseListPageOptions, UseListPageReturn } from '@/types/page'

export function useListPage<T, P = Record<string, any>>(
  options: UseListPageOptions<T, P>
): UseListPageReturn<T, P> {
  const {
    fetchFn,
    extractList = (data) => data?.items || data?.list || data || [],
    pageSize = 10,
    immediate = false,
    getCacheFn,
    saveCacheFn,
  } = options

  // 状态
  const list = ref<T[]>([]) as Ref<T[]>
  const pageState = ref<PageStateType>('loading')
  const errorMessage = ref('')
  const loading = ref(false)
  const loadingMore = ref(false)
  const noMore = ref(false)
  const refreshing = ref(false)
  const isOffline = ref(false)
  
  // 分页
  const page = ref(1)
  const total = ref(0)
  
  // 额外参数
  const extraParams = ref<Partial<P>>({}) as Ref<Partial<P>>

  const load = async (params?: Partial<P>, isRefresh = false) => {
    if (isRefresh) {
      page.value = 1
      noMore.value = false
      pageState.value = 'loading'
    }
    
    const currentPage = page.value
    
    if (currentPage > 1) {
      loadingMore.value = true
    } else {
      loading.value = true
    }
    
    try {
      const mergedParams = {
        ...extraParams.value,
        ...params,
        page: currentPage,
        pageSize,
      } as P & { page: number; page_size: number }
      
      const res = await fetchFn(mergedParams)
      const items = extractList(res.data)
      
      if (isRefresh || currentPage === 1) {
        list.value = items
      } else {
        list.value.push(...items)
      }
      
      // 判断是否还有更多
      if (items.length < pageSize) {
        noMore.value = true
      }
      
      // 更新状态
      if (list.value.length === 0) {
        pageState.value = 'empty'
      } else {
        pageState.value = 'content'
        isOffline.value = false
      }
      
      // 保存缓存
      if (saveCacheFn && isRefresh && items.length > 0) {
        saveCacheFn(items)
      }
      
      page.value++
    } catch (error: any) {
      console.error('列表加载失败', error)
      
      if (currentPage === 1 || list.value.length === 0) {
        // 尝试使用缓存
        if (getCacheFn) {
          const cached = getCacheFn()
          if (cached && cached.length > 0) {
            list.value = cached
            isOffline.value = true
            pageState.value = 'content'
            noMore.value = true
            return
          }
        }
        
        pageState.value = 'error'
        errorMessage.value = error?.message || '加载失败，请重试'
      } else {
        uni.showToast({ title: '加载失败', icon: 'none' })
      }
    } finally {
      loading.value = false
      loadingMore.value = false
      refreshing.value = false
    }
  }

  const loadMore = async () => {
    if (loadingMore.value || noMore.value || loading.value) return
    await load()
  }

  const refresh = async (params?: Partial<P>) => {
    if (params) {
      extraParams.value = { ...extraParams.value, ...params }
    }
    refreshing.value = true
    await load(params, true)
  }

  const reset = () => {
    list.value = []
    page.value = 1
    noMore.value = false
    isOffline.value = false
    pageState.value = 'loading'
    extraParams.value = {} as Partial<P>
  }

  // 立即加载
  if (immediate) {
    load(undefined, true)
  }

  return {
    list,
    pageState,
    errorMessage,
    loading,
    loadingMore,
    noMore,
    refreshing,
    isOffline,
    page,
    total,
    load,
    loadMore,
    refresh,
    reset,
    extraParams,
  }
}
