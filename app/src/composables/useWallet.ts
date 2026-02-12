/**
 * 钱包专用 Hook
 */
import { ref, computed, reactive } from 'vue'
import { getCouponStats, getTransactions, getVipStatus, getWalletInfo } from '@/api/wallet'
import type { TransactionData, TransactionType, WalletState } from '@/types/wallet'
import type { QuickActionItem, TabItem } from '@/types/ui'

export function useWallet() {
  // 状态
  const loading = ref(true)
  const loadingMore = ref(false)
  const noMore = ref(false)
  const showBalance = ref(true)
  const currentFilter = ref('all')
  const page = ref(1)
  
  // 数据
  const wallet = reactive<WalletState>({
    balance: 0,
    frozenBalance: 0,
    vipLevel: 0,
    couponCount: 0,
    totalSpent: 0,
    totalRecharge: 0,
  })
  
  const records = ref<TransactionData[]>([])
  
  // 筛选标签
  const filterTabs: TabItem[] = [
    { key: 'all', label: '全部' },
    { key: 'income', label: '收入' },
    { key: 'expense', label: '支出' },
  ]
  
  // 快捷入口（icon 为 uv-icon 名称）
  const quickActions = computed((): QuickActionItem[] => [
    { key: 'coupons', icon: 'coupon', label: '优惠券', badge: wallet.couponCount || undefined },
    { key: 'earnings', icon: 'list', label: '收益' },
    { key: 'invite', icon: 'gift', label: '邀请有礼' },
    { key: 'help', icon: 'question-circle', label: '帮助' },
  ])
  
  // 过滤后的记录
  const filteredRecords = computed(() => {
    if (currentFilter.value === 'all') return records.value
    if (currentFilter.value === 'income') {
      return records.value.filter(r => r.amount > 0)
    }
    return records.value.filter(r => r.amount < 0)
  })
  
  // VIP 折扣文案
  const vipDiscountText = computed(() => {
    const discounts = { 1: 9.8, 2: 9.5, 3: 9, 4: 8.5, 5: 8 }
    return (discounts as any)[wallet.vipLevel] || 10
  })
  
  // 加载钱包信息
  const loadWalletInfo = async () => {
    try {
      const [walletRes, couponStatsRes, vipRes] = await Promise.allSettled([
        getWalletInfo(),
        getCouponStats(),
        getVipStatus(),
      ])

      if (walletRes.status === 'fulfilled' && walletRes.value.data) {
        const data = walletRes.value.data
        wallet.balance = data.balanceCents || 0
        wallet.frozenBalance = data.frozenCents || 0
        wallet.totalSpent = data.totalSpentCents || 0
        wallet.totalRecharge = data.totalRechargeCents || 0
      }

      if (couponStatsRes.status === 'fulfilled' && couponStatsRes.value.data) {
        wallet.couponCount = couponStatsRes.value.data.available || 0
      } else {
        wallet.couponCount = 0
      }

      if (vipRes.status === 'fulfilled' && vipRes.value.data) {
        wallet.vipLevel = vipRes.value.data.currentLevel?.level || 0
      } else {
        wallet.vipLevel = 0
      }
    } catch (error) {
      console.error('加载钱包信息失败', error)
    }
  }
  
  // 加载交易记录
  const loadRecords = async (refresh = true) => {
    if (refresh) {
      page.value = 1
      noMore.value = false
      records.value = []
    }
    
    loading.value = refresh
    loadingMore.value = !refresh
    
    try {
      const res = await getTransactions({
        page: page.value,
        pageSize: 20,
      }, { showError: false })
      
      const items = (res.data as any)?.items || res.data || []
      const newRecords = items.map((t: any): TransactionData => ({
        id: t.id,
        type: t.type as TransactionType,
        title: t.title || getTypeTitle(t.type),
        description: t.description || '',
        amount: t.amountCents || 0,
        createdAt: t.createdAt,
      }))
      
      if (refresh) {
        records.value = newRecords
      } else {
        records.value.push(...newRecords)
      }
      
      if (newRecords.length < 20) {
        noMore.value = true
      }
      
      page.value++
    } catch (error) {
      console.error('加载交易记录失败', error)
    } finally {
      loading.value = false
      loadingMore.value = false
    }
  }
  
  const getTypeTitle = (type: TransactionType) => {
    const titles: Record<TransactionType, string> = {
      recharge: '充值',
      withdraw: '提现',
      withdrawal: '提现',
      payment: '订单支付',
      refund: '退款',
      earning: '收益入账',
      bonus: '活动奖励',
      commission: '佣金',
    }
    return titles[type] || '交易'
  }
  
  // 加载更多
  const loadMore = () => {
    if (loadingMore.value || noMore.value) return
    loadRecords(false)
  }
  
  // 快捷入口点击
  const handleQuickAction = (key: string) => {
    const routes: Record<string, string> = {
      coupons: '/pages/coupon/list/index',
      earnings: '/pages/player/earnings/index',
      invite: '/pages/invite/index',
      help: '/pages/help/index',
    }
    const route = routes[key]
    if (route) {
      uni.navigateTo({ url: route })
    } else {
      uni.showToast({ title: '功能开发中', icon: 'none' })
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  const goToRecharge = () => uni.navigateTo({ url: '/pages/wallet/recharge/index' })
  const goToWithdraw = () => uni.navigateTo({ url: '/pages/wallet/withdraw/index' })
  const goToVip = () => uni.showToast({ title: 'VIP 功能开发中', icon: 'none' })
  
  // 初始化
  const init = () => {
    loadWalletInfo()
    loadRecords()
  }
  
  return {
    // 状态
    loading,
    loadingMore,
    noMore,
    showBalance,
    currentFilter,
    
    // 数据
    wallet,
    records,
    filteredRecords,
    filterTabs,
    quickActions,
    vipDiscountText,
    
    // 方法
    loadMore,
    handleQuickAction,
    goBack,
    goToRecharge,
    goToWithdraw,
    goToVip,
    init,
  }
}
