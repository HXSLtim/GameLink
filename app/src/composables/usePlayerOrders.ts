/**
 * 陪玩师订单列表专用 Hook
 */
import { ref, reactive } from 'vue'
import { useListPage } from './useListPage'
import { getPlayerOrders, acceptOrder as acceptPlayerOrder, finishService as completePlayerOrder } from '@/api/order'
import { confirmDialog } from '@/composables/useConfirmDialog'
import type { Order, OrderActionKey, OrderTabItem, OrderTabKey } from '@/types/order'
import { normalizeOrderStatus } from '@/components/OrderCard/utils'

type PlayerOrderData = Order

export function usePlayerOrders() {
  // 标签
  const tabs = reactive<OrderTabItem[]>([
    { key: 'all', label: '全部', badge: 0 },
    { key: 'pending', label: '待接单', badge: 0 },
    { key: 'in_progress', label: '进行中', badge: 0 },
    { key: 'completed', label: '已完成', badge: 0 },
  ])
  const currentTab = ref<OrderTabKey>('all')

  // 构建参数
  const buildParams = () => {
    const params: Record<string, any> = {}
    if (currentTab.value !== 'all') {
      params.status = currentTab.value
    }
    return params
  }

  // 使用通用列表 Hook
  const listPage = useListPage<PlayerOrderData>({
    fetchFn: async (params) => {
      const res = await getPlayerOrders({
        page: params.page,
        page_size: params.pageSize,
        ...buildParams(),
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const orders = data?.orders || data?.items || data || []
      return orders.map((o: any): PlayerOrderData => ({
        id: o.id,
        orderNo: o.orderNo,
        status: normalizeOrderStatus(o.status, 'player'),
        user: {
          id: o.userId,
          nickname: o.userNickname || '用户',
          avatar: o.userAvatar,
        },
        gameName: o.gameName || '游戏',
        serviceName: o.serviceType || '陪玩',
        quantity: o.quantity || 1,
        unit: o.unit || '局',
        earnings: o.playerEarningsCents || o.totalCents * 0.8 || 0,
        remark: o.remark,
        createdAt: o.createdAt,
      }))
    },
    pageSize: 20,
  })

  // 切换标签
  const switchTab = (tab: OrderTabKey) => {
    currentTab.value = tab
    listPage.refresh(buildParams())
  }

  // 获取空状态文案
  const getEmptyTitle = () => {
    const titles: Record<string, string> = {
      all: '暂无订单',
      pending: '暂无待接单',
      in_progress: '暂无进行中订单',
      completed: '暂无已完成订单',
    }
    return titles[currentTab.value] || '暂无订单'
  }

  const getEmptyDesc = () => {
    const descs: Record<string, string> = {
      all: '努力接单，收益多多',
      pending: '新订单会及时通知您',
      in_progress: '接单后开始服务',
      completed: '完成订单后可查看',
    }
    return descs[currentTab.value] || ''
  }

  // 操作处理
  const handleAction = async (order: PlayerOrderData, action: OrderActionKey) => {
    switch (action) {
      case 'accept':
        await acceptOrder(order)
        break
      case 'reject':
        rejectOrder(order)
        break
      case 'start':
        startService(order)
        break
      case 'complete':
        await completeService(order)
        break
      case 'contact':
        contactUser(order)
        break
      case 'detail':
        viewDetail(order)
        break
    }
  }

  const acceptOrder = async (order: PlayerOrderData) => {
    const confirmed = await confirmDialog({
      title: '确认接单',
      content: '确定接受该订单吗？',
    })
    if (!confirmed) return
    try {
      uni.showLoading({ title: '处理中...' })
      await acceptPlayerOrder(order.id)
      uni.hideLoading()
      uni.showToast({ title: '接单成功', icon: 'success' })
      listPage.refresh(buildParams())
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }

  const rejectOrder = async (order: PlayerOrderData) => {
    const confirmed = await confirmDialog({
      title: '确认拒绝',
      content: '确定拒绝该订单吗？',
    })
    if (!confirmed) return
    uni.showToast({ title: '已拒绝', icon: 'success' })
    listPage.refresh(buildParams())
  }

  const startService = (order: PlayerOrderData) => {
    uni.showToast({ title: '开始服务', icon: 'success' })
    listPage.refresh(buildParams())
  }

  const completeService = async (order: PlayerOrderData) => {
    const confirmed = await confirmDialog({
      title: '确认完成',
      content: '确定服务已完成吗？',
    })
    if (!confirmed) return
    try {
      uni.showLoading({ title: '处理中...' })
      await completePlayerOrder(order.id)
      uni.hideLoading()
      uni.showToast({ title: '服务已完成', icon: 'success' })
      listPage.refresh(buildParams())
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }

  const contactUser = (order: PlayerOrderData) => {
    const userId = order.user?.id
    if (!userId) {
      uni.showToast({ title: '用户信息缺失', icon: 'none' })
      return
    }
    uni.navigateTo({ url: `/pages/message/chat/index?userId=${userId}` })
  }

  const viewDetail = (order: PlayerOrderData) => {
    uni.navigateTo({ url: `/pages/order/detail/index?id=${order.id}` })
  }

  // 导航
  const goBack = () => uni.navigateBack()

  return {
    // 列表数据
    orders: listPage.list,
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
    getEmptyTitle,
    getEmptyDesc,
    handleAction,
    goBack,
  }
}
