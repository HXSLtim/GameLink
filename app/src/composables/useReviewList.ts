/**
 * 评价列表专用 Hook
 */
import { ref, reactive } from 'vue'
import { useListPage } from './useListPage'
import type { ReviewData } from '@/components/ReviewCard/index.vue'

// Mock API - 实际应替换为真实 API
const getReviews = async (params: any) => {
  // TODO: 替换为真实 API
  return { data: { items: [], total: 0 } }
}

export function useReviewList() {
  // 标签
  const tabs = reactive([
    { label: '已评价', value: 'completed', count: 0 },
    { label: '待评价', value: 'pending', count: 0 },
  ])
  const currentTab = ref('completed')
  
  // 使用通用列表 Hook
  const listPage = useListPage<ReviewData>({
    fetchFn: async (params) => {
      const res = await getReviews({
        page: params.page,
        pageSize: params.pageSize,
        type: currentTab.value,
      })
      return res
    },
    extractList: (data: any) => {
      return data?.items || data || []
    },
    pageSize: 20,
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
  const skipReview = (review: ReviewData) => {
    uni.showToast({ title: '已跳过', icon: 'success' })
    listPage.list.value = listPage.list.value.filter(r => r.id !== review.id)
  }
  
  // 去评价
  const writeReview = (review: ReviewData) => {
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
