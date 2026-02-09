/**
 * 陪玩师工作台专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { getTodayStats } from '@/api/player'
import type { StatItem, QuickActionItem, MenuItem } from '@/types/ui'
import { formatDateChinese, formatYuan } from '@/utils/format'
import type { PlayerDashboardInfo, PlayerTodayStats } from '@/types/player'
import type { DashboardStatus } from '@/types/status'


export function usePlayerDashboard() {
  const userStore = useUserStore()
  
  // 状态
  const refreshing = ref(false)
  const workStatus = ref<DashboardStatus>('offline')
  const orderBadge = ref(0)
  
  // 数据
  const playerInfo = ref<PlayerDashboardInfo>({
    id: userStore.userInfo?.id || 0,
    nickname: userStore.userInfo?.nickname || '陪玩师',
    avatar: userStore.userInfo?.avatar,
    rating: 5.0,
    certificationStatus: 'none',
  })
  
  const todayStats = reactive<PlayerTodayStats>({
    orders: 0,
    earnings: 0,
    duration: 0,
    rating: '5.0',
  })
  
  // 统计数据项
  const statItems = computed((): StatItem[] => [
    { value: todayStats.orders, label: '接单数', onClick: goToOrders },
    { value: `¥${formatYuan(todayStats.earnings)}`, label: '收益', highlight: true, onClick: goToEarnings },
    { value: `${todayStats.duration}h`, label: '服务时长' },
    { value: todayStats.rating, label: '评分' },
  ])
  
  // 快捷入口（icon 为 uv-icon 名称）
  const quickActions = computed((): QuickActionItem[] => [
    { key: 'orders', icon: 'list', label: '订单管理', badge: orderBadge.value },
    { key: 'earnings', icon: 'red-packet', label: '我的收益' },
    { key: 'services', icon: 'grid-fill', label: '服务管理' },
    { key: 'schedule', icon: 'calendar', label: '排班设置' },
  ])
  
  // 功能菜单
  const menuItems = computed((): MenuItem[] => [
    { key: 'certification', label: '陪玩认证', icon: 'checkbox-mark' },
    { key: 'profile', label: '个人资料', icon: 'account' },
    { key: 'reviews', label: '用户评价', icon: 'star', value: playerInfo.value.rating?.toFixed(1) || '5.0' },
    { key: 'help', label: '帮助中心', icon: 'question-circle' },
  ])
  
  // 加载今日数据
  const loadTodayStats = async () => {
    try {
      const res = await getTodayStats()
      if (res.data) {
        todayStats.orders = res.data.orderCount || 0
        todayStats.earnings = (res.data.earningsCents || 0) / 100
        todayStats.duration = res.data.serviceDuration || 0
        todayStats.rating = res.data.averageRating?.toFixed(1) || '5.0'
      }
    } catch (error) {
      console.error('加载今日数据失败', error)
    }
  }
  
  // 刷新
  const onRefresh = async () => {
    refreshing.value = true
    await loadTodayStats()
    refreshing.value = false
  }
  
  // 切换工作状态
  const toggleWorkStatus = () => {
    if (workStatus.value === 'online') {
      workStatus.value = 'offline'
      uni.showToast({ title: '已下线', icon: 'none' })
    } else {
      workStatus.value = 'online'
      uni.showToast({ title: '已上线', icon: 'success' })
    }
  }
  
  // 快捷入口点击
  const handleQuickAction = (key: string) => {
    const routes: Record<string, string> = {
      orders: '/pages/player/orders/index',
      earnings: '/pages/player/earnings/index',
      services: '/pages/player/services/index',
      schedule: '/pages/player/schedule/index',
    }
    const route = routes[key]
    if (route) {
      uni.navigateTo({ url: route })
    } else {
      uni.showToast({ title: '功能开发中', icon: 'none' })
    }
  }
  
  // 菜单点击
  const handleMenuClick = (item: MenuItem) => {
    const routes: Record<string, string> = {
      certification: '/pages/player/certification/index',
      profile: '/pages/profile/edit/index',
      reviews: '/pages/review/list/index',
      help: '/pages/help/index',
    }
    const route = routes[item.key]
    if (route) {
      uni.navigateTo({ url: route })
    }
  }
  
  // 导航
  const goToSettings = () => uni.navigateTo({ url: '/pages/settings/index/index' })
  const goToOrders = () => uni.navigateTo({ url: '/pages/player/orders/index' })
  const goToEarnings = () => uni.navigateTo({ url: '/pages/player/earnings/index' })
  
  // 格式化日期
  const formatDate = (date: Date) => formatDateChinese(date)
  
  // 初始化
  const init = () => {
    loadTodayStats()
  }
  
  return {
    // 状态
    refreshing,
    workStatus,
    playerInfo,
    
    // 数据
    todayStats,
    statItems,
    quickActions,
    menuItems,
    
    // 方法
    onRefresh,
    toggleWorkStatus,
    handleQuickAction,
    handleMenuClick,
    goToSettings,
    formatDate,
    init,
  }
}
