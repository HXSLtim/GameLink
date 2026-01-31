/**
 * 个人中心专用 Hook
 * 封装用户信息加载和菜单配置
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useTheme } from '@/composables/useTheme'
import { getUserProfile } from '@/api/user'
import { getWalletInfo } from '@/api/wallet'
import type { MenuItem } from '@/components/MenuList/index.vue'

export function useProfile() {
  const userStore = useUserStore()
  const { isDark } = useTheme()
  
  // 扩展用户信息
  const userInfoData = ref({
    orderCount: 0,
    favoriteCount: 0,
    balance: 0,
  })
  
  // 订单数量
  const orderCounts = reactive({
    pending: 0,
    inProgress: 0,
    toReview: 0,
  })
  
  // 计算属性
  const isLoggedIn = computed(() => userStore.isLoggedIn)
  
  const userInfo = computed(() => ({
    ...userStore.userInfo,
    ...userInfoData.value,
  }))
  
  // 功能菜单（登录用户）
  const userMenuItems = computed((): MenuItem[] => {
    const items: MenuItem[] = [
      {
        key: 'wallet',
        label: '我的钱包',
        icon: 'red-packet',
        iconColor: 'var(--color-primary)',
        iconBg: 'rgba(0, 210, 106, 0.1)',
      },
      {
        key: 'favorites',
        label: '我的收藏',
        icon: 'heart',
        iconColor: '#EF4444',
        iconBg: 'rgba(239, 68, 68, 0.1)',
      },
    ]
    
    if (userStore.userInfo?.role === 'player') {
      items.push({
        key: 'playerCenter',
        label: '陪玩中心',
        icon: 'grid',
        iconColor: '#8B5CF6',
        iconBg: 'rgba(139, 92, 246, 0.1)',
      })
    } else {
      items.push({
        key: 'becomePlayer',
        label: '成为陪玩师',
        icon: 'grid',
        iconColor: '#8B5CF6',
        iconBg: 'rgba(139, 92, 246, 0.1)',
      })
    }
    
    return items
  })
  
  // 设置菜单
  const settingsMenuItems: MenuItem[] = [
    {
      key: 'settings',
      label: '设置',
      icon: 'setting',
      iconColor: 'var(--color-text-secondary)',
    },
    {
      key: 'help',
      label: '帮助与反馈',
      icon: 'question-circle',
      iconColor: 'var(--color-text-secondary)',
    },
  ]
  
  // 主题菜单
  const themeMenuItem = computed((): MenuItem => ({
    key: 'theme',
    label: '深色模式',
    icon: isDark.value ? 'eye-off' : 'eye',
    iconColor: isDark.value ? '#FFD700' : '#FF8C00',
    iconBg: isDark.value ? 'rgba(255, 215, 0, 0.1)' : 'rgba(255, 140, 0, 0.1)',
  }))
  
  // 加载用户数据
  const loadUserData = async () => {
    if (!isLoggedIn.value) return
    
    try {
      // 加载用户资料
      const profileRes = await getUserProfile()
      if (profileRes.data) {
        userStore.updateUserInfo({
          id: profileRes.data.id,
          nickname: profileRes.data.nickname,
          avatar: profileRes.data.avatar || '',
          role: profileRes.data.role,
          playerId: profileRes.data.playerId,
        })
      }
      
      // 加载钱包余额
      const walletRes = await getWalletInfo()
      if (walletRes.data) {
        userInfoData.value.balance = walletRes.data.balanceCents
      }
    } catch (error) {
      console.error('加载用户数据失败:', error)
    }
  }
  
  // 菜单点击处理
  const handleMenuClick = (item: MenuItem) => {
    const routes: Record<string, string> = {
      wallet: '/pages/wallet/index/index',
      favorites: '/pages/favorite/list/index',
      settings: '/pages/settings/index/index',
    }
    
    const route = routes[item.key]
    if (route) {
      if (!isLoggedIn.value && ['wallet', 'favorites'].includes(item.key)) {
        goToLogin()
        return
      }
      uni.navigateTo({ url: route })
    } else {
      // 功能开发中
      uni.showToast({ title: '功能开发中', icon: 'none' })
    }
  }
  
  // 订单快捷入口点击
  const handleOrderClick = (status: string) => {
    if (!isLoggedIn.value) {
      goToLogin()
      return
    }
    uni.navigateTo({ url: `/pages/order/list/index?status=${status}` })
  }
  
  // 查看全部订单
  const handleViewAllOrders = () => {
    if (!isLoggedIn.value) {
      goToLogin()
      return
    }
    uni.navigateTo({ url: '/pages/order/list/index' })
  }
  
  // 统计点击
  const handleStatClick = (type: 'orders' | 'favorites' | 'wallet') => {
    if (!isLoggedIn.value) {
      goToLogin()
      return
    }
    
    const routes = {
      orders: '/pages/order/list/index',
      favorites: '/pages/favorite/list/index',
      wallet: '/pages/wallet/index/index',
    }
    uni.navigateTo({ url: routes[type] })
  }
  
  // 编辑资料
  const goToEdit = () => {
    if (!isLoggedIn.value) {
      goToLogin()
      return
    }
    uni.navigateTo({ url: '/pages/profile/edit/index' })
  }
  
  // 登录
  const goToLogin = () => {
    uni.navigateTo({ url: '/pages/auth/login/index' })
  }
  
  // 退出登录
  const handleLogout = () => {
    uni.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          userStore.logout()
          uni.showToast({ title: '已退出登录', icon: 'success' })
        }
      }
    })
  }
  
  return {
    // 状态
    isLoggedIn,
    userInfo,
    orderCounts,
    isDark,
    
    // 菜单配置
    userMenuItems,
    settingsMenuItems,
    themeMenuItem,
    
    // 方法
    loadUserData,
    handleMenuClick,
    handleOrderClick,
    handleViewAllOrders,
    handleStatClick,
    goToEdit,
    goToLogin,
    handleLogout,
  }
}
