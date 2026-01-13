import { http } from '../../../utils/request'

interface OrderDetail {
  id: number
  orderNo: string
  status: string
  statusText: string
  player: {
    id: number
    nickname: string
    avatar: string
  }
  service: {
    name: string
    price: number
    unit: string
  }
  quantity: number
  totalPrice: number
  remark?: string
  createdAt: string
}

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
    orderId: 0,
    order: {} as OrderDetail,
    statusMap,
  },
  
  onLoad(options) {
    const orderId = Number(options.id) || 0
    this.setData({ orderId })
    this.loadOrder()
  },
  
  async loadOrder() {
    this.setData({ loading: true })
    try {
      const order = await http.get<OrderDetail>(`/orders/${this.data.orderId}`)
      this.setData({ order, loading: false })
    } catch (err) {
      console.error('加载订单失败:', err)
      // 模拟数据
      this.setData({
        loading: false,
        order: {
          id: this.data.orderId,
          orderNo: 'GL202501100001',
          status: 'pending',
          statusText: '待支付',
          player: {
            id: 1,
            nickname: '小甜甜',
            avatar: '',
          },
          service: {
            name: '王者荣耀-陪玩',
            price: 30,
            unit: '局',
          },
          quantity: 2,
          totalPrice: 60,
          remark: '希望能带我上分~',
          createdAt: '2025-01-10 12:00:00',
        },
      })
    }
  },
  
  handlePay() {
    wx.showToast({ title: '支付功能开发中', icon: 'none' })
  },
  
  handleCancel() {
    wx.showModal({
      title: '取消订单',
      content: '确定要取消该订单吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await http.post(`/orders/${this.data.orderId}/cancel`)
            wx.showToast({ title: '订单已取消', icon: 'success' })
            this.loadOrder()
          } catch (err) {
            wx.showToast({ title: '取消失败', icon: 'none' })
          }
        }
      },
    })
  },
  
  handleContact() {
    const { order } = this.data
    wx.navigateTo({ url: `/pages/chat/index?orderId=${order.id}` })
  },
})
