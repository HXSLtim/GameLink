import type { Ref } from 'vue'

export type PageStateType = 'loading' | 'error' | 'empty' | 'offline' | 'login' | 'content'

export interface PageStateOptions {
  /** 初始状态 */
  initialState?: PageStateType
  /** 是否需要登录 */
  requireAuth?: boolean
}

export interface PaginationOptions {
  initialPage?: number
  pageSize?: number
}

export interface UseListPageOptions<T, P = Record<string, any>> {
  /** 获取数据的 API 函数 */
  fetchFn: (params: P & { page: number; page_size: number; pageSize: number }) => Promise<{ data: any }>
  /** 从响应中提取数据列表 */
  extractList?: (data: any) => T[]
  /** 每页数量 */
  pageSize?: number
  /** 是否立即加载 */
  immediate?: boolean
  /** 缓存 key（用于离线缓存） */
  cacheKey?: string
  /** 获取缓存数据的函数 */
  getCacheFn?: () => T[] | null
  /** 保存缓存数据的函数 */
  saveCacheFn?: (data: T[]) => void
}

export interface UseListPageReturn<T, P = Record<string, any>> {
  // 状态
  list: Ref<T[]>
  pageState: Ref<PageStateType>
  errorMessage: Ref<string>
  loading: Ref<boolean>
  loadingMore: Ref<boolean>
  noMore: Ref<boolean>
  refreshing: Ref<boolean>
  isOffline: Ref<boolean>

  // 分页信息
  page: Ref<number>
  total: Ref<number>

  // 方法
  load: (params?: Partial<P>) => Promise<void>
  loadMore: () => Promise<void>
  refresh: (params?: Partial<P>) => Promise<void>
  reset: () => void

  // 额外参数
  extraParams: Ref<Partial<P>>
}
