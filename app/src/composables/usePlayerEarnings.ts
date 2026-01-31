/**
 * 陪玩师收益专用 Hook
 */
import { ref, computed, reactive } from 'vue'
import { getPlayerEarnings, getEarningsStats } from '@/api/player'
import type { ChartDataItem } from '@/components/EarningsChart/index.vue'
import type { TabItem } from '@/components/TabsBar/index.vue'

interface EarningsItem {
  id: number
  type: string
  title: string
  description: string
  amount: number // 分
  orderId?: number
  createdAt: string
}

export function usePlayerEarnings() {
  // 状态
  const loading = ref(true)
  const loadingMore = ref(false)
  const noMore = ref(false)
  const showFilter = ref(false)
  const page = ref(1)
  
  // 汇总数据
  const summary = reactive({
    withdrawable: 0,
    pending: 0,
    withdrawn: 0,
    total: 0,
  })
  
  // 周期
  const currentPeriod = ref('week')
  const periodTabs: TabItem[] = [
    { key: 'week', label: '本周' },
    { key: 'month', label: '本月' },
    { key: 'year', label: '本年' },
  ]
  
  // 图表数据
  const chartData = ref<ChartDataItem[]>([])
  const periodTotal = ref(0)
  
  // 收益列表
  const earningsList = ref<EarningsItem[]>([])
  
  // 筛选
  const filterType = ref('all')
  const typeOptions = [
    { label: '全部', value: 'all' },
    { label: '订单收益', value: 'order' },
    { label: '提现', value: 'withdraw' },
    { label: '活动奖励', value: 'bonus' },
  ]
  
  // 获取周期标题
  const getPeriodTitle = () => {
    const titles: Record<string, string> = {
      week: '本周收益',
      month: '本月收益',
      year: '本年收益',
    }
    return titles[currentPeriod.value] || '收益统计'
  }
  
  // 加载统计
  const loadStats = async () => {
    try {
      const res = await getEarningsStats({ period: currentPeriod.value })
      const data = res.data as any
      
      summary.withdrawable = data?.withdrawableCents || 0
      summary.pending = data?.pendingCents || 0
      summary.withdrawn = data?.withdrawnCents || 0
      summary.total = data?.totalCents || 0
      
      periodTotal.value = data?.periodTotalCents || 0
      chartData.value = (data?.chartData || []).map((d: any) => ({
        label: d.label,
        value: d.valueCents || 0,
      }))
    } catch (error) {
      console.error('加载统计失败', error)
    }
  }
  
  // 加载列表
  const loadList = async (refresh = true) => {
    if (refresh) {
      page.value = 1
      noMore.value = false
      earningsList.value = []
    }
    
    loading.value = refresh
    loadingMore.value = !refresh
    
    try {
      const res = await getPlayerEarnings({
        page: page.value,
        pageSize: 20,
        type: filterType.value === 'all' ? undefined : filterType.value,
      }, { showError: false })
      
      const items = ((res.data as any)?.items || res.data || []).map((e: any): EarningsItem => ({
        id: e.id,
        type: e.type || 'order',
        title: e.title || getTypeTitle(e.type),
        description: e.description || '',
        amount: e.amountCents || 0,
        orderId: e.orderId,
        createdAt: e.createdAt,
      }))
      
      if (refresh) {
        earningsList.value = items
      } else {
        earningsList.value.push(...items)
      }
      
      if (items.length < 20) noMore.value = true
      page.value++
    } catch (error) {
      console.error('加载列表失败', error)
    } finally {
      loading.value = false
      loadingMore.value = false
    }
  }
  
  const getTypeTitle = (type: string) => {
    const titles: Record<string, string> = {
      order: '订单收益',
      withdraw: '提现',
      bonus: '活动奖励',
      refund: '退款扣除',
    }
    return titles[type] || '收益'
  }
  
  // 切换周期
  const switchPeriod = (period: string) => {
    currentPeriod.value = period
    loadStats()
  }
  
  // 加载更多
  const loadMore = () => {
    if (loadingMore.value || noMore.value) return
    loadList(false)
  }
  
  // 筛选
  const applyFilter = () => {
    showFilter.value = false
    loadList(true)
  }
  
  const resetFilter = () => {
    filterType.value = 'all'
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  const goToWithdraw = () => uni.navigateTo({ url: '/pages/wallet/withdraw/index' })
  const goToOrder = (orderId?: number) => {
    if (orderId) {
      uni.navigateTo({ url: `/pages/order/detail/index?id=${orderId}` })
    }
  }
  
  // 初始化
  const init = () => {
    loadStats()
    loadList()
  }
  
  return {
    // 状态
    loading,
    loadingMore,
    noMore,
    showFilter,
    
    // 数据
    summary,
    currentPeriod,
    periodTabs,
    chartData,
    periodTotal,
    earningsList,
    filterType,
    typeOptions,
    
    // 方法
    getPeriodTitle,
    switchPeriod,
    loadMore,
    applyFilter,
    resetFilter,
    goBack,
    goToWithdraw,
    goToOrder,
    init,
  }
}
