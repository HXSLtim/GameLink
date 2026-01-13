import { http } from '../../../utils/request'

interface Order {
  id: number
  orderNo: string
  status: string
  player: {
    nickname: string
    avatar: string
  }
  service: {
    name: string
  }
  quantity: number
  totalPrice: number
  createdAt: string
}

const statusTabs = [
  { label: '全部', value: 'all' },
  { label: '待支付', value: 'pending' },
  { label: '进行中', value: 'processing' },
  { label: '已完成', value: 'completed' },
]

const statusMap: Record<string, { text: string; type: string }> = {
  pending: { text: '待支付', type: 'warning' },
  paid: { text: '已支付', type: 'info' },
  accepted: { text: '已接单', type: 'primary' },
  processing: { text: '进行中', type: 'primary' },
  completed: { text: '已完成', type: 'success' },
  cancelled: { text: '已取消', type: 'default' },
  refunded: { text: '已退款', type: 'error' },
}

Page({
  data: {
    loading: true,
    orders: [] as Order[],
    activeStatus: 'all',
    statusTabs,
    statusMap,
    page: 1,
    hasMore: true,
  },
  
  onLoad(options) {
    const status = options.status || 'all'
    this.setData({ activeStatus: status })
    this.loadOrders()
  },
  
  onPullDownRefresh() {
    this.setData({ page: 1, hasMore: true })
    this.loadOrders().finally(() => {
      wx.stopPullDownRefresh()
    })
  },
  
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadMoreOrders()
    }
  },
  
  async loadOrders() {
    this.setData({ loading: true })
    try {
      const { activeStatus, page } = this.data
      const params: Record<string, unknown> = { page, pageSize: 10 }
      if (activeStatus !== 'all') params.status = activeStatus
      
      const orders = await http.get<Order[]>('/orders', { data: params })
      this.setData({ orders, loading: false, hasMore: orders.length >= 10 })
    } catch (err) {
      console.error('加载订单失败:', err)
      this.setData({ loading: false, orders: [] })
    }
  },
  
  async loadMoreOrders() {
    const { page } = this.data
    this.setData({ page: page + 1 })
    await this.loadOrders()
  },
  
  handleStatusTap(e: WechatMiniprogram.TouchEvent) {
    const { status } = e.currentTarget.dataset
    this.setData({ activeStatus: status, page: 1 })
    this.loadOrders()
  },
  
  handleOrderTap(e: WechatMiniprogram.TouchEvent) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/order/detail/index?id=${id}` })
  },
})
