import { http } from '../../../utils/request'

interface Service {
  id: number
  name: string
  price: number
  unit: string
}

Page({
  data: {
    playerId: 0,
    services: [] as Service[],
    selectedService: 0,
    quantity: 1,
    remark: '',
    serviceFee: 0,
    totalPrice: 0,
    submitting: false,
  },
  
  onLoad(options) {
    const playerId = Number(options.playerId) || 0
    this.setData({ playerId })
    this.loadServices()
  },
  
  async loadServices() {
    try {
      const services = await http.get<Service[]>(`/players/${this.data.playerId}/services`)
      if (services.length > 0) {
        this.setData({
          services,
          selectedService: services[0].id,
        })
        this.calculatePrice()
      }
    } catch (err) {
      console.error('加载服务失败:', err)
      // 模拟数据
      const services = [
        { id: 1, name: '陪玩', price: 30, unit: '局' },
        { id: 2, name: '代练', price: 50, unit: '星' },
      ]
      this.setData({
        services,
        selectedService: services[0].id,
      })
      this.calculatePrice()
    }
  },
  
  handleServiceSelect(e: WechatMiniprogram.TouchEvent) {
    const { id } = e.currentTarget.dataset
    this.setData({ selectedService: id })
    this.calculatePrice()
  },
  
  handleQuantityMinus() {
    const { quantity } = this.data
    if (quantity > 1) {
      this.setData({ quantity: quantity - 1 })
      this.calculatePrice()
    }
  },
  
  handleQuantityPlus() {
    const { quantity } = this.data
    if (quantity < 99) {
      this.setData({ quantity: quantity + 1 })
      this.calculatePrice()
    }
  },
  
  handleRemarkInput(e: WechatMiniprogram.Input) {
    this.setData({ remark: e.detail.value })
  },
  
  calculatePrice() {
    const { services, selectedService, quantity } = this.data
    const service = services.find(s => s.id === selectedService)
    if (service) {
      const serviceFee = service.price * quantity
      this.setData({
        serviceFee,
        totalPrice: serviceFee,
      })
    }
  },
  
  async handleSubmit() {
    const { playerId, selectedService, quantity, remark, submitting } = this.data
    
    if (submitting) return
    if (!selectedService) {
      wx.showToast({ title: '请选择服务', icon: 'none' })
      return
    }
    
    this.setData({ submitting: true })
    
    try {
      const order = await http.post<{ id: number }>('/orders', {
        playerId,
        serviceId: selectedService,
        quantity,
        remark,
      }, { showLoading: true })
      
      wx.redirectTo({ url: `/pages/order/detail/index?id=${order.id}` })
    } catch (err) {
      console.error('创建订单失败:', err)
      wx.showToast({ title: '创建订单失败', icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },
})
