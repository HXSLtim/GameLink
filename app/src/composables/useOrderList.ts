/**
 * 订单列表专用 Hook
 * 封装订单列表的数据加载、状态切换等逻辑
 */
import { ref, computed, onMounted } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { useListPage } from './useListPage'
import { getOrders, type Order as ApiOrder, type OrderStatus as ApiOrderStatus } from '@/api/order'
import { useUserStore } from '@/store/user'
import type { Order, OrderStatus } from '@/components/OrderCard/index.vue'
import type { TabItem } from '@/components/TabsBar/index.vue'

// 默认标签配置
export const defaultOrderTabs: TabItem[] = [
  { key: 'all', label: '全部' },
  { key: 'pending', label: '待支付' },
  { key: 'in_progress', label: '进行中' },
  { key: 'completed', label: '已完成' },
  { key: 'canceled', label: '已取消' },
]

// 转换 API 响应为 Order 类型
function transformOrder(o: ApiOrder): Order {
  return {
    id: o.id,
    orderNo: o.orderNo,
    status: o.status.replace('pending_payment', 'pending').replace('pending_accept', 'confirmed') as OrderStatus,
    player: {
      id: o.playerId,
      nickname: o.playerName,
      avatar: o.playerAvatar,
    },
    gameName: o.gameName,
    serviceName: o.serviceType,
    quantity: o.quantity,
    unit: o.unit,
    totalAmount: o.totalPriceCents / 100,
    createdAt: o.createdAt,
    reviewed: false,
  }
}

export function useOrderList() {
  const userStore = useUserStore()
  
  // 当前标签
  const currentTab = ref('all')
  
  // 使用通用列表 Hook
  const listPage = useListPage<Order, { status?: ApiOrderStatus }>({
    fetchFn: async (params) => {
      const res = await getOrders({
        page: params.page,
        pageSize: params.pageSize,
        status: params.status,
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const orderList: ApiOrder[] = data?.orders || data || []
      return orderList.map(transformOrder)
    },
    pageSize: 10,
  })
  
  // 切换标签
  const switchTab = (tab: string) => {
    if (currentTab.value === tab) return
    currentTab.value = tab
    listPage.reset()
    refreshList()
  }
  
  // 刷新列表
  const refreshList = () => {
    // 检查登录状态
    if (!userStore.isLoggedIn) {
      listPage.pageState.value = 'login'
      return
    }
    
    const params: { status?: ApiOrderStatus } = {}
    if (currentTab.value !== 'all') {
      params.status = currentTab.value as ApiOrderStatus
    }
    listPage.refresh(params)
  }
  
  // 订单操作
  const handleOrderAction = (key: string, order: Order) => {
    switch (key) {
      case 'cancel':
        cancelOrder(order)
        break
      case 'pay':
        uni.navigateTo({ url: `/pages/order/detail/index?id=${order.id}&needPay=1` })
        break
      case 'contact':
        uni.navigateTo({ url: `/pages/message/chat/index?playerId=${order.player.id}` })
        break
      case 'complete':
        completeOrder(order)
        break
      case 'review':
        uni.navigateTo({ url: `/pages/order/detail/index?id=${order.id}&action=review` })
        break
      case 'reorder':
        uni.navigateTo({ url: `/pages/order/create/index?playerId=${order.player.id}` })
        break
      case 'viewDispute':
        uni.navigateTo({ url: `/pages/order/detail/index?id=${order.id}` })
        break
    }
  }
  
  const cancelOrder = (order: Order) => {
    uni.showModal({
      title: '确认取消',
      content: '确定要取消这个订单吗？',
      success: async (res) => {
        if (res.confirm) {
          uni.showLoading({ title: '取消中...' })
          await new Promise(resolve => setTimeout(resolve, 500))
          order.status = 'canceled'
          uni.hideLoading()
          uni.showToast({ title: '已取消', icon: 'success' })
        }
      }
    })
  }
  
  const completeOrder = (order: Order) => {
    uni.showModal({
      title: '确认完成',
      content: '确认陪玩师已完成服务吗？',
      success: async (res) => {
        if (res.confirm) {
          uni.showLoading({ title: '确认中...' })
          await new Promise(resolve => setTimeout(resolve, 500))
          order.status = 'completed'
          uni.hideLoading()
          uni.showToast({ title: '已完成', icon: 'success' })
        }
      }
    })
  }
  
  // 跳转详情
  const goToDetail = (orderId: number) => {
    uni.navigateTo({ url: `/pages/order/detail/index?id=${orderId}` })
  }
  
  // 跳转陪玩师列表
  const goToPlayerList = () => {
    uni.switchTab({ url: '/pages/player/list/index' })
  }
  
  // 返回
  const goBack = () => {
    uni.navigateBack()
  }
  
  return {
    // 列表数据
    orders: listPage.list,
    pageState: listPage.pageState,
    errorMessage: listPage.errorMessage,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    
    // 标签
    currentTab,
    tabs: defaultOrderTabs,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: refreshList,
    switchTab,
    handleOrderAction,
    goToDetail,
    goToPlayerList,
    goBack,
  }
}
