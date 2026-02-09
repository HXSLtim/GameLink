/**
 * 评价列表专用 Hook
 */
import { ref, reactive } from 'vue'
import { useListPage } from './useListPage'
import { get, type RequestConfig } from '@/api/request'
import type { ReviewCardData } from '@/types/review'
import type { TabItem } from '@/types/ui'

/**
 * 获取我的评价列表
 */
function getMyReviews(params?: Record<string, any>, config?: Partial<RequestConfig>) {
  return get<{ items: ReviewCardData[]; total: number }>('/user/reviews/my', params, config)
}

export function useReviewList() {
  // 标签
  const tabs = reactive<TabItem[]>([
    { key: 'completed', label: '已评价', badge: 0 },
    { key: 'pending', label: '待评价', badge: 0 },
  ])
  const currentTab = ref('completed')
  
  // 使用通用列表 Hook
  const listPage = useListPage<ReviewCardData>({
    fetchFn: async (params) => {
      const res = await getMyReviews({
        page: params.page,
        page_size: params.pageSize,
        type: currentTab.value,
      })
      return res
    },
    extractList: (data: any) => {
      return data?.items || data || []
    },
    page_size: 20,
  })
  
  // 切换标签
  const switchTab = (tab: string) => {
    currentTab.value = tab
    listPage.refresh()
  }
  
  // 跳转到订单
  const goToOrder = (orderId: number) => {
    uni.navigateTo({ url: `/pages/order/detail/index?id=${orderId}` })
  }
  
  // 跳过评价
  const skipReview = (review: ReviewCardData) => {
    uni.showToast({ title: '已跳过', icon: 'success' })
    listPage.list.value = listPage.list.value.filter(r => r.id !== review.id)
  }
  
  // 去评价
  const writeReview = (review: ReviewCardData) => {
    uni.navigateTo({ url: `/pages/order/detail/index?id=${review.orderId}&action=review` })
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  return {
    // 列表数据
    reviews: listPage.list,
    pageState: listPage.pageState,
    loading: listPage.loading,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    
    // 标签
    tabs,
    currentTab,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: listPage.refresh,
    switchTab,
    goToOrder,
    skipReview,
    writeReview,
    goBack,
  }
}
